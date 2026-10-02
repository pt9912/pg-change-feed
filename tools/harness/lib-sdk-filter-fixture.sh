# tools/harness/lib-sdk-filter-fixture.sh — gemeinsame Vorbereitung und
# gemeinsamer Phasenablauf der Filter-Phase (schema/table am SSE-Stream) für
# die drei SDK-Realserver-Runner (ADR-0133 Teilfrage 4): eine Aufrufform für
# alle drei Runner (BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall). Wird
# per `source` eingebunden (Schwester von tools/harness/lib-sdk-route-fixture.sh,
# ruft deren Nachbarin sdk_rule_fixture_await_applied aus
# tools/harness/lib-sdk-rule-fixture.sh auf), kein eigenes Executable. Setzt
# PG_CONTAINER, PG_USER, PG_DB, NETWORK, FEED_CONTAINER, SDK_TEST_CONTAINER und
# SDK_INTEGRATION_IMAGE voraus (Konvention der drei Runner).
#
# Drei Tabellen mit (id, name): public.$SDK_FILTER_TABLE_A,
# public.$SDK_FILTER_TABLE_B und $SDK_FILTER_SCHEMA_OTHER.$SDK_FILTER_TABLE_A
# (gleicher Tabellenname in einem zweiten Schema). Je Versuch committet der
# Runner eine Dreiergruppe in der Reihenfolge B, zweites Schema, A: ein Client
# mit schema und table der Tabelle A sieht vor der Change seiner Tabelle zwei
# Changes, die er nicht empfangen darf.

SDK_FILTER_SOURCE=src-e2e
SDK_FILTER_SCHEMA_A=public
SDK_FILTER_SCHEMA_OTHER=feed_e2e_sdkfilt_s
SDK_FILTER_TABLE_A=feed_e2e_sdkfilt_a
SDK_FILTER_TABLE_B=feed_e2e_sdkfilt_b
# Ruhefenster in Sekunden: derselbe Wert wie SDK_ROUTE_QUIET_SECONDS
# (tools/harness/lib-sdk-route-fixture.sh); das Fenster beginnt erst, wenn der
# Client ohne Filter alle drei Gruppen empfangen hat.
SDK_FILTER_QUIET_SECONDS=15
SDK_FILTER_ATTEMPTS=5

# Zusammenfassung der Phasen, vom Runner am Ende gedruckt.
SDK_FILTER_REPORT=""

# sdk_filter_fixture_setup <Aufrufer-Präfix> — legt die drei Tabellen an
# (das zweite Schema inklusive), aktiviert sie je per cdc.enable_table und
# wartet je Antrag mit Frist auf `status = 'applied'`; `failed` beendet den
# Runner mit dem Fehlertext der Zeile.
sdk_filter_fixture_setup() {
  local prefix=$1 request_id schema table

  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 >/dev/null <<SQL
CREATE SCHEMA $SDK_FILTER_SCHEMA_OTHER;
CREATE TABLE $SDK_FILTER_SCHEMA_A.$SDK_FILTER_TABLE_A (id int PRIMARY KEY, name text);
CREATE TABLE $SDK_FILTER_SCHEMA_A.$SDK_FILTER_TABLE_B (id int PRIMARY KEY, name text);
CREATE TABLE $SDK_FILTER_SCHEMA_OTHER.$SDK_FILTER_TABLE_A (id int PRIMARY KEY, name text);
SQL

  for pair in "$SDK_FILTER_SCHEMA_A $SDK_FILTER_TABLE_A" "$SDK_FILTER_SCHEMA_A $SDK_FILTER_TABLE_B" "$SDK_FILTER_SCHEMA_OTHER $SDK_FILTER_TABLE_A"; do
    set -- $pair
    schema=$1
    table=$2
    request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT cdc.enable_table('$SDK_FILTER_SOURCE', '$schema', '$table')")
    if [ -z "$request_id" ]; then
      echo "$prefix: cdc.enable_table($schema.$table) lieferte keine Antrags-ID" >&2
      exit 1
    fi
    sdk_rule_fixture_await_applied "$prefix" "$request_id" "cdc.enable_table($schema.$table)"
  done
}

# sdk_filter_psql <SQL> — eine Abfrage gegen die Datenbank, Ergebnis auf stdout.
sdk_filter_psql() {
  docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc "$1"
}

# sdk_filter_phase <Aufrufer-Präfix> <Phasenname> <Sentinel> <ID-Basis>
# <Test-Auswahl KEY=WERT> <Extra-Env> — ein Realserver-Rundlauf für den
# SSE-Client mit Filter: Container starten, auf READY warten, je Versuch eine
# Dreiergruppe (B, zweites Schema, A) committen, bis der Test SEEN druckt
# (F1 hat die Change der Tabelle A, F2 die des zweiten Schemas, U alle drei),
# dann auf das Prozessende warten. Die Testausgabe trägt je empfangener Change
# eine Zeile `RECEIVED_<F1|F2|U> change_id=<id> schema=<s> table=<t>` und eine
# Abschlusszeile `FILTER_RESULT f1=<n> f1_foreign=<k> f2=<n> f2_foreign=<k>
# unfiltered=<m> quiet_seconds=<s>`; die Fremdwerte sind am Empfang gezählt, der
# Runner verlangt an beiden 0, hält jede Kennung gegen cdc.changes
# (schema_name, table_name, Sentinel) und setzt SDK_FILTER_REPORT fort.
sdk_filter_phase() {
  local prefix=$1 phase_name=$2 sentinel=$3 id_base=$4 select_env=$5 extra_env=$6
  local quiet=$SDK_FILTER_QUIET_SECONDS
  local attempt group_id seen=0 seen_after_ms=0 pair started finished

  local env_args=()
  for pair in $extra_env; do
    env_args+=(-e "$pair")
  done

  docker rm -fv "$SDK_TEST_CONTAINER" >/dev/null 2>&1 || true
  docker run -d --name "$SDK_TEST_CONTAINER" --network "$NETWORK" \
    "${env_args[@]}" \
    -e "$select_env" \
    -e PGCHANGEFEED_E2E_SENTINEL="$sentinel" \
    -e PGCHANGEFEED_FILTER_SCHEMA_A="$SDK_FILTER_SCHEMA_A" \
    -e PGCHANGEFEED_FILTER_SCHEMA_OTHER="$SDK_FILTER_SCHEMA_OTHER" \
    -e PGCHANGEFEED_FILTER_TABLE_A="$SDK_FILTER_TABLE_A" \
    -e PGCHANGEFEED_FILTER_TABLE_B="$SDK_FILTER_TABLE_B" \
    -e PGCHANGEFEED_FILTER_QUIET_SECONDS="$quiet" \
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

  for attempt in $(seq 1 "$SDK_FILTER_ATTEMPTS"); do
    group_id=$((id_base + attempt * 3))
    started=$(date +%s%N)
    docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 >/dev/null <<SQL
INSERT INTO $SDK_FILTER_SCHEMA_A.$SDK_FILTER_TABLE_B (id, name) VALUES ($((group_id + 1)), '$sentinel');
INSERT INTO $SDK_FILTER_SCHEMA_OTHER.$SDK_FILTER_TABLE_A (id, name) VALUES ($((group_id + 2)), '$sentinel');
INSERT INTO $SDK_FILTER_SCHEMA_A.$SDK_FILTER_TABLE_A (id, name) VALUES ($((group_id + 3)), '$sentinel');
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

  # Die Abschlusszeile steht vor jeder Auswertung im Fehlertext: ein Fremdwert
  # ist in der Zeile selbst sichtbar.
  local result_line
  result_line=$(printf '%s\n' "$test_output" | grep -E '^[[:space:]]*FILTER_RESULT ' | sed -E 's/^[[:space:]]+//' | head -n1 || true)

  if [ "$seen" -ne 1 ]; then
    echo "$prefix: $phase_name — der Test meldete nach $SDK_FILTER_ATTEMPTS Dreiergruppen kein SEEN: $test_output" >&2
    exit 1
  fi
  if ! printf '%s' "$result_line" | grep -qE "^FILTER_RESULT f1=[0-9]+ f1_foreign=[0-9]+ f2=[0-9]+ f2_foreign=[0-9]+ unfiltered=[0-9]+ quiet_seconds=$quiet([^0-9]|\$)"; then
    echo "$prefix: $phase_name — die Abschlusszeile FILTER_RESULT fehlt oder trägt nicht die erwartete Form (quiet_seconds=$quiet): ${result_line:-leer}; Ausgang des Tests: $test_exit; Ausgabe: $test_output" >&2
    exit 1
  fi
  local f1_count f1_foreign f2_count f2_foreign unfiltered_count
  f1_count=$(printf '%s' "$result_line" | sed -E 's/.* f1=([0-9]+) .*/\1/')
  f1_foreign=$(printf '%s' "$result_line" | sed -E 's/.* f1_foreign=([0-9]+) .*/\1/')
  f2_count=$(printf '%s' "$result_line" | sed -E 's/.* f2=([0-9]+) .*/\1/')
  f2_foreign=$(printf '%s' "$result_line" | sed -E 's/.* f2_foreign=([0-9]+) .*/\1/')
  unfiltered_count=$(printf '%s' "$result_line" | sed -E 's/.* unfiltered=([0-9]+) .*/\1/')
  if [ "$f1_foreign" != "0" ] || [ "$f2_foreign" != "0" ]; then
    echo "$prefix: $phase_name — fremde Changes am Filter-Client (f1_foreign=$f1_foreign, f2_foreign=$f2_foreign): $result_line" >&2
    exit 1
  fi
  if [ "$test_stopped" -ne 1 ] || [ "$test_exit" != "0" ]; then
    echo "$prefix: $phase_name — der Test endete nicht mit Ausgang 0 (gestoppt: $test_stopped, Ausgang: $test_exit): $test_output" >&2
    exit 1
  fi

  # sdk_filter_check_lines <Marker> <erwartetes Schema> <erwartete Tabelle oder
  # leer = beliebig> <erwartete Zahl> — jede Kennung der Marker-Zeilen trägt in
  # cdc.changes Schema, Tabelle und Sentinel dieser Phase.
  local ident row
  sdk_filter_check_lines() {
    local marker=$1 want_schema=$2 want_table=$3 want_count=$4 lines=0 ident row
    while IFS= read -r ident; do
      [ -n "$ident" ] || continue
      lines=$((lines + 1))
      row=$(sdk_filter_psql "SELECT schema_name || '.' || table_name FROM cdc.changes WHERE source_id = '$SDK_FILTER_SOURCE' AND change_id = '$ident' AND new_data->>'name' = '$sentinel'")
      if [ -z "$row" ]; then
        echo "$prefix: $phase_name — die Kennung $ident ($marker) ist in cdc.changes mit dem Sentinel dieser Phase nicht lesbar" >&2
        exit 1
      fi
      if [ -n "$want_table" ] && [ "$row" != "$want_schema.$want_table" ]; then
        echo "$prefix: $phase_name — $marker empfing $ident, das in cdc.changes (schema_name.table_name) '$row' trägt, erwartet $want_schema.$want_table" >&2
        exit 1
      fi
      if [ -z "$want_table" ] && [ "${row%%.*}" != "$want_schema" ]; then
        echo "$prefix: $phase_name — $marker empfing $ident, das in cdc.changes (schema_name.table_name) '$row' trägt, erwartet Schema $want_schema" >&2
        exit 1
      fi
    done < <(printf '%s\n' "$test_output" | sed -nE "s/^[[:space:]]*$marker change_id=([^ ]+).*/\1/p")
    if [ "$lines" -lt 1 ] || [ "$lines" != "$want_count" ]; then
      echo "$prefix: $phase_name — $marker-Zeilen ($lines) und FILTER_RESULT-Zahl ($want_count) passen nicht zusammen oder sind leer" >&2
      exit 1
    fi
  }
  sdk_filter_check_lines RECEIVED_F1 "$SDK_FILTER_SCHEMA_A" "$SDK_FILTER_TABLE_A" "$f1_count"
  sdk_filter_check_lines RECEIVED_F2 "$SDK_FILTER_SCHEMA_OTHER" "" "$f2_count"

  # Der Client ohne Filter hat alle drei Gruppen gesehen, jede Kennung in der
  # Sicht mit dem Sentinel dieser Phase.
  local unfiltered_lines=0 seen_tables=""
  while IFS= read -r ident; do
    [ -n "$ident" ] || continue
    unfiltered_lines=$((unfiltered_lines + 1))
    row=$(sdk_filter_psql "SELECT schema_name || '.' || table_name FROM cdc.changes WHERE source_id = '$SDK_FILTER_SOURCE' AND change_id = '$ident' AND new_data->>'name' = '$sentinel'")
    case "$row" in
      "$SDK_FILTER_SCHEMA_A.$SDK_FILTER_TABLE_A" | "$SDK_FILTER_SCHEMA_A.$SDK_FILTER_TABLE_B" | "$SDK_FILTER_SCHEMA_OTHER.$SDK_FILTER_TABLE_A") ;;
      *)
        echo "$prefix: $phase_name — die Kennung $ident des Clients ohne Filter trägt in cdc.changes (schema_name.table_name) '${row:-nicht gefunden}'" >&2
        exit 1
        ;;
    esac
    seen_tables="$seen_tables|$row"
  done < <(printf '%s\n' "$test_output" | sed -nE 's/^[[:space:]]*RECEIVED_U change_id=([^ ]+).*/\1/p')
  for want in "|$SDK_FILTER_SCHEMA_A.$SDK_FILTER_TABLE_A" "|$SDK_FILTER_SCHEMA_A.$SDK_FILTER_TABLE_B" "|$SDK_FILTER_SCHEMA_OTHER.$SDK_FILTER_TABLE_A"; do
    case "$seen_tables|" in
      *"$want|"*) ;;
      *)
        echo "$prefix: $phase_name — der Client ohne Filter empfing keine Change von '${want#|}' (Gegenlesung über cdc.changes)" >&2
        exit 1
        ;;
    esac
  done
  if [ "$unfiltered_lines" != "$unfiltered_count" ]; then
    echo "$prefix: $phase_name — RECEIVED_U-Zeilen ($unfiltered_lines) und FILTER_RESULT unfiltered=$unfiltered_count passen nicht zusammen" >&2
    exit 1
  fi

  local feed_running
  feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
  if [ "$feed_running" != "true" ]; then
    echo "$prefix: $phase_name — Feed-Container lief nach dem Lauf nicht mehr weiter (kein Neustart erwartet)" >&2
    exit 1
  fi

  SDK_FILTER_REPORT="$SDK_FILTER_REPORT; $phase_name: F1 ($SDK_FILTER_SCHEMA_A.$SDK_FILTER_TABLE_A) $f1_count Change(s), gemessen f1_foreign=$f1_foreign, F2 (Schema $SDK_FILTER_SCHEMA_OTHER) $f2_count Change(s), gemessen f2_foreign=$f2_foreign, keine fremde im Ruhefenster von ${quiet}s, $unfiltered_count Change(s) dieser Phase am Client ohne Filter (alle drei Gruppen), SEEN nach ${seen_after_ms}ms ab dem letzten Commit (Versuch $attempt); gedruckte Zeile: $result_line"
}
