# tools/harness/lib-sdk-route-fixture.sh — gemeinsame Vorbereitung und
# gemeinsamer Phasenablauf der Routing-Phasen für die drei
# SDK-Realserver-Runner (ADR-0137): eine Aufrufform für alle drei Runner
# (BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall). Wird per `source`
# eingebunden (Muster tools/harness/lib-sdk-rule-fixture.sh, deren
# sdk_rule_fixture_await_applied sie aufruft), kein eigenes Executable. Setzt
# PG_CONTAINER, PG_USER, PG_DB, NETWORK, FEED_CONTAINER, SDK_TEST_CONTAINER
# und SDK_INTEGRATION_IMAGE voraus (Konvention der drei Runner).
#
# Die Tabelle trägt (id, region, name). Zwei Inhaltsregeln lenken die Zeilen
# der Region eu auf das Ziel eu (A) und die Zeilen der Region us auf das Ziel
# us (B); eine Zeile der Region asia trifft keine Regel und trägt kein Ziel.
# Je Versuch committet der Runner eine Dreiergruppe in der Reihenfolge asia,
# us, eu: ein Client mit Ziel eu sieht vor der Change seines Ziels zwei
# Changes, die er nicht empfangen darf.

SDK_ROUTE_TABLE=feed_e2e_sdkroute
SDK_ROUTE_SOURCE=src-e2e
SDK_ROUTE_TARGET_A=eu
SDK_ROUTE_TARGET_B=us
SDK_ROUTE_REGION_NONE=asia
# Ruhefenster der Stream-Phasen in Sekunden: derselbe Wert wie das
# Ruhefenster des Server-Rundlaufs (RT_WINDOW in run-integration-tests.sh);
# übernommen, nicht neu gemessen. Das Fenster beginnt erst, wenn der Client
# ohne Ziel die Change des fremden Ziels bereits empfangen hat.
SDK_ROUTE_QUIET_SECONDS=15
SDK_ROUTE_ATTEMPTS=5

# Zusammenfassung der Phasen (eine Zeile je Fläche), vom Runner am Ende
# gedruckt.
SDK_ROUTE_REPORT=""

# sdk_route_fixture_setup <Aufrufer-Präfix> — legt $SDK_ROUTE_TABLE an,
# aktiviert sie und setzt die zwei Routing-Regeln; wartet je Antrag mit Frist
# auf `status = 'applied'`, `failed` beendet den Runner mit dem Fehlertext der
# Zeile.
sdk_route_fixture_setup() {
  local prefix=$1 enable_request_id request_id

  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 >/dev/null <<SQL
CREATE TABLE public.$SDK_ROUTE_TABLE (id int PRIMARY KEY, region text, name text);
SQL

  enable_request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT cdc.enable_table('$SDK_ROUTE_SOURCE', 'public', '$SDK_ROUTE_TABLE')")
  if [ -z "$enable_request_id" ]; then
    echo "$prefix: cdc.enable_table($SDK_ROUTE_TABLE) lieferte keine Antrags-ID" >&2
    exit 1
  fi
  sdk_rule_fixture_await_applied "$prefix" "$enable_request_id" "cdc.enable_table($SDK_ROUTE_TABLE)"

  request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT cdc.set_route('$SDK_ROUTE_SOURCE', 'public', '$SDK_ROUTE_TABLE', 'sdk_route_a', '{\"target\": \"$SDK_ROUTE_TARGET_A\", \"order\": 10, \"when\": {\"column\": \"region\", \"equals\": \"$SDK_ROUTE_TARGET_A\"}}')")
  if [ -z "$request_id" ]; then
    echo "$prefix: cdc.set_route($SDK_ROUTE_TABLE, sdk_route_a) lieferte keine Antrags-ID" >&2
    exit 1
  fi
  sdk_rule_fixture_await_applied "$prefix" "$request_id" "cdc.set_route($SDK_ROUTE_TABLE, Ziel $SDK_ROUTE_TARGET_A)"

  request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT cdc.set_route('$SDK_ROUTE_SOURCE', 'public', '$SDK_ROUTE_TABLE', 'sdk_route_b', '{\"target\": \"$SDK_ROUTE_TARGET_B\", \"order\": 20, \"when\": {\"column\": \"region\", \"equals\": \"$SDK_ROUTE_TARGET_B\"}}')")
  if [ -z "$request_id" ]; then
    echo "$prefix: cdc.set_route($SDK_ROUTE_TABLE, sdk_route_b) lieferte keine Antrags-ID" >&2
    exit 1
  fi
  sdk_rule_fixture_await_applied "$prefix" "$request_id" "cdc.set_route($SDK_ROUTE_TABLE, Ziel $SDK_ROUTE_TARGET_B)"
}

# sdk_route_psql <SQL> — eine Abfrage gegen die Datenbank, Ergebnis auf stdout.
sdk_route_psql() {
  docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc "$1"
}

# sdk_route_phase <Aufrufer-Präfix> <Phasenname> <Sentinel> <ID-Basis>
# <Test-Auswahl KEY=WERT> <Extra-Env> <Ruhefenster in Sekunden> — ein
# Realserver-Rundlauf für genau eine SDK-Fläche mit Ziel: Container starten,
# auf READY warten, je Versuch eine Dreiergruppe (asia, us, eu) committen, bis
# der Test SEEN druckt (Client mit Ziel hat eine Change seines Ziels, Client
# ohne Ziel alle drei Regionen), dann auf das Prozessende warten. Die
# Testausgabe trägt je empfangener Change eine Zeile
# `RECEIVED_TARGETED change_id=<id>` (Client mit Ziel) bzw.
# `RECEIVED_UNFILTERED change_id=<id>` (Client ohne Ziel, nur Zeilen dieser
# Phase) und eine Abschlusszeile `ROUTE_RESULT target=<Ziel> targeted=<n>
# foreign=0 unfiltered=<m> quiet_seconds=<s>`. Der Runner hält jede Kennung
# gegen cdc.changes (route_target) und setzt SDK_ROUTE_REPORT fort. Das
# Ruhefenster 0 kennzeichnet eine Pull-Fläche (der Test vergleicht dort die
# Lesung mit Ziel mit zwei ungefilterten Lesungen und wartet nicht).
sdk_route_phase() {
  local prefix=$1 phase_name=$2 sentinel=$3 id_base=$4 select_env=$5 extra_env=$6 quiet=$7
  local attempt group_id seen=0 seen_after_ms=0 pair started finished

  local env_args=()
  for pair in $extra_env; do
    env_args+=(-e "$pair")
  done

  docker rm -fv "$SDK_TEST_CONTAINER" >/dev/null 2>&1 || true
  docker run -d --name "$SDK_TEST_CONTAINER" --network "$NETWORK" \
    "${env_args[@]}" \
    -e "$select_env" \
    -e PGCHANGEFEED_E2E_TABLE="$SDK_ROUTE_TABLE" \
    -e PGCHANGEFEED_E2E_SENTINEL="$sentinel" \
    -e PGCHANGEFEED_ROUTE_TARGET_A="$SDK_ROUTE_TARGET_A" \
    -e PGCHANGEFEED_ROUTE_TARGET_B="$SDK_ROUTE_TARGET_B" \
    -e PGCHANGEFEED_ROUTE_REGION_NONE="$SDK_ROUTE_REGION_NONE" \
    -e PGCHANGEFEED_ROUTE_QUIET_SECONDS="$quiet" \
    "$SDK_INTEGRATION_IMAGE" >/dev/null

  local test_ready=0
  for _ in $(seq 1 120); do
    if docker logs "$SDK_TEST_CONTAINER" 2>/dev/null | grep -qF "READY"; then
      test_ready=1
      break
    fi
    if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
      break
    fi
    sleep 1
  done
  if [ "$test_ready" -ne 1 ]; then
    echo "$prefix: $phase_name — Test wurde nicht innerhalb der Zeitspanne bereit (kein READY): $(docker logs "$SDK_TEST_CONTAINER" 2>&1 || true)" >&2
    exit 1
  fi

  for attempt in $(seq 1 "$SDK_ROUTE_ATTEMPTS"); do
    group_id=$((id_base + attempt * 3))
    started=$(date +%s%N)
    docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 >/dev/null <<SQL
INSERT INTO public.$SDK_ROUTE_TABLE (id, region, name) VALUES
  ($((group_id + 1)), '$SDK_ROUTE_REGION_NONE', '$sentinel');
INSERT INTO public.$SDK_ROUTE_TABLE (id, region, name) VALUES
  ($((group_id + 2)), '$SDK_ROUTE_TARGET_B', '$sentinel');
INSERT INTO public.$SDK_ROUTE_TABLE (id, region, name) VALUES
  ($((group_id + 3)), '$SDK_ROUTE_TARGET_A', '$sentinel');
SQL
    for _ in $(seq 1 60); do
      if docker logs "$SDK_TEST_CONTAINER" 2>/dev/null | grep -qE '^[[:space:]]*SEEN[[:space:]]*$'; then
        seen=1
        finished=$(date +%s%N)
        seen_after_ms=$(((finished - started) / 1000000))
        break
      fi
      if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
        break
      fi
      sleep 0.2
    done
    if [ "$seen" -eq 1 ]; then
      break
    fi
  done

  local test_stopped=0
  for _ in $(seq 1 $((quiet + 120))); do
    if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
      test_stopped=1
      break
    fi
    sleep 1
  done
  local test_exit test_output
  test_exit=$(docker inspect --format '{{.State.ExitCode}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo unbekannt)
  test_output=$(docker logs "$SDK_TEST_CONTAINER" 2>&1 || true)

  if [ "$seen" -ne 1 ]; then
    echo "$prefix: $phase_name — der Test meldete nach $SDK_ROUTE_ATTEMPTS Dreiergruppen kein SEEN (Ziel $SDK_ROUTE_TARGET_A und Gegenseite nicht beide empfangen): $test_output" >&2
    exit 1
  fi
  if [ "$test_stopped" -ne 1 ] || [ "$test_exit" != "0" ]; then
    echo "$prefix: $phase_name — der Test endete nicht mit Ausgang 0 (gestoppt: $test_stopped, Ausgang: $test_exit): $test_output" >&2
    exit 1
  fi

  local result_line
  result_line=$(printf '%s\n' "$test_output" | grep -E '^[[:space:]]*ROUTE_RESULT ' | sed -E 's/^[[:space:]]+//' | head -n1 || true)
  if ! printf '%s' "$result_line" | grep -qE "^ROUTE_RESULT target=$SDK_ROUTE_TARGET_A targeted=[0-9]+ foreign=0 unfiltered=[0-9]+ quiet_seconds=$quiet([^0-9]|\$)"; then
    echo "$prefix: $phase_name — die Abschlusszeile ROUTE_RESULT fehlt oder trägt nicht die erwartete Form (foreign=0, quiet_seconds=$quiet): ${result_line:-leer}" >&2
    exit 1
  fi
  local targeted_count unfiltered_count
  targeted_count=$(printf '%s' "$result_line" | sed -E 's/.* targeted=([0-9]+) .*/\1/')
  unfiltered_count=$(printf '%s' "$result_line" | sed -E 's/.* unfiltered=([0-9]+) .*/\1/')

  # Jede vom Client mit Ziel empfangene Kennung trägt in der Sicht das Ziel A
  # (und die Region A); eine Kennung mit anderem Ziel oder ohne Ziel ist ein
  # Befund, den der Test (foreign) und diese Gegenlesung unabhängig voneinander
  # tragen.
  local ident route targeted_lines=0
  while IFS= read -r ident; do
    [ -n "$ident" ] || continue
    targeted_lines=$((targeted_lines + 1))
    route=$(sdk_route_psql "SELECT coalesce(route_target, 'NULL') || ' ' || coalesce(new_data->>'region', 'NULL') FROM cdc.changes WHERE source_id = '$SDK_ROUTE_SOURCE' AND change_id = '$ident' AND table_name = '$SDK_ROUTE_TABLE'")
    if [ "$route" != "$SDK_ROUTE_TARGET_A $SDK_ROUTE_TARGET_A" ]; then
      echo "$prefix: $phase_name — der Client mit Ziel $SDK_ROUTE_TARGET_A empfing die Change $ident, die in cdc.changes (route_target region) '${route:-nicht gefunden}' trägt" >&2
      exit 1
    fi
  done < <(printf '%s\n' "$test_output" | sed -nE 's/^[[:space:]]*RECEIVED_TARGETED change_id=([^ ]+).*/\1/p')
  if [ "$targeted_lines" -lt 1 ] || [ "$targeted_lines" != "$targeted_count" ]; then
    echo "$prefix: $phase_name — RECEIVED_TARGETED-Zeilen ($targeted_lines) und ROUTE_RESULT targeted=$targeted_count passen nicht zusammen oder sind leer" >&2
    exit 1
  fi

  # Der Client ohne Ziel hat für diese Phase alle drei Ziele gesehen: Ziel A,
  # Ziel B und eine Change ohne Ziel (route_target NULL), jede Kennung in der
  # Sicht mit dem Sentinel dieser Phase.
  local unfiltered_lines=0 seen_targets=""
  while IFS= read -r ident; do
    [ -n "$ident" ] || continue
    unfiltered_lines=$((unfiltered_lines + 1))
    route=$(sdk_route_psql "SELECT coalesce(route_target, 'NULL') || ' ' || coalesce(new_data->>'region', 'NULL') FROM cdc.changes WHERE source_id = '$SDK_ROUTE_SOURCE' AND change_id = '$ident' AND table_name = '$SDK_ROUTE_TABLE' AND new_data->>'name' = '$sentinel'")
    case "$route" in
      "$SDK_ROUTE_TARGET_A $SDK_ROUTE_TARGET_A" | "$SDK_ROUTE_TARGET_B $SDK_ROUTE_TARGET_B" | "NULL $SDK_ROUTE_REGION_NONE") ;;
      *)
        echo "$prefix: $phase_name — die Kennung $ident des Clients ohne Ziel trägt in cdc.changes (route_target region) '${route:-nicht gefunden}'" >&2
        exit 1
        ;;
    esac
    seen_targets="$seen_targets|$route"
  done < <(printf '%s\n' "$test_output" | sed -nE 's/^[[:space:]]*RECEIVED_UNFILTERED change_id=([^ ]+).*/\1/p')
  for want in "|$SDK_ROUTE_TARGET_A $SDK_ROUTE_TARGET_A" "|$SDK_ROUTE_TARGET_B $SDK_ROUTE_TARGET_B" "|NULL $SDK_ROUTE_REGION_NONE"; do
    case "$seen_targets|" in
      *"$want|"*) ;;
      *)
        echo "$prefix: $phase_name — der Client ohne Ziel empfing keine Change der Gruppe '${want#|}' (Gegenlesung über cdc.changes)" >&2
        exit 1
        ;;
    esac
  done
  if [ "$unfiltered_lines" != "$unfiltered_count" ]; then
    echo "$prefix: $phase_name — RECEIVED_UNFILTERED-Zeilen ($unfiltered_lines) und ROUTE_RESULT unfiltered=$unfiltered_count passen nicht zusammen" >&2
    exit 1
  fi

  local feed_running
  feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
  if [ "$feed_running" != "true" ]; then
    echo "$prefix: $phase_name — Feed-Container lief nach dem Lauf nicht mehr weiter (kein Neustart erwartet)" >&2
    exit 1
  fi

  local window_text="keine fremde im Ruhefenster von ${quiet}s"
  if [ "$quiet" -eq 0 ]; then
    window_text="keine fremde (Lesung mit Ziel = Ziel-Teilmenge der ungefilterten Lesung davor und danach)"
  fi
  SDK_ROUTE_REPORT="$SDK_ROUTE_REPORT; $phase_name: Ziel $SDK_ROUTE_TARGET_A, $targeted_count Change(s) mit Ziel $SDK_ROUTE_TARGET_A und $window_text, $unfiltered_count Change(s) dieser Phase am Client ohne Ziel (alle Gruppen $SDK_ROUTE_TARGET_A/$SDK_ROUTE_TARGET_B/ohne Ziel), SEEN nach ${seen_after_ms}ms ab dem letzten Commit (Versuch $attempt)"
}
