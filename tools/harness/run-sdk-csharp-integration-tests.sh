#!/usr/bin/env bash
# run-sdk-csharp-integration-tests.sh — Realserver-Integrationstest der
# C#-SDK-Zustellweg-Flächen (Mechanik-Klasse
# ADR-0110 §Entscheidung Festlegung 2, gespiegelt vom Python-Vorbild
# tools/harness/run-sdk-python-integration-tests.sh): zwölf Phasen — die vier
# Flächen des Packages PgChangeFeed.Client (HTTP, gRPC, SSE
# und NATS-Vollinhalt) prüfen ihre Protokoll-Annahmen je gegen
# eine reale, laufende Server-Instanz ohne Regel, dieselben vier Flächen
# ein zweites Mal gegen eine Tabelle mit aktiver `rename_column`-Regel und
# ein drittes Mal gegen eine Tabelle mit zwei Routing-Regeln (ein Client mit
# `target` empfängt nur die Changes seines Ziels, ein Client ohne `target`
# alle; Vorbereitung und Ablauf über tools/harness/lib-sdk-route-fixture.sh) —
# der Prüfling ist die
# kompilierte Client-Assembly, der Integrationstest importiert
# PgChangeFeed.Client direkt (kein Wegwerf-Duplikat-Client daneben).
#
# Kette je Lauf: compose.yaml-Umgebung hochfahren (PostgreSQL/NATS/
# Feed-Container, Schema-Rollout über d-migrate, Vorbedingungen der
# Aktivierung), Bau der `integration`-Docker-Stufe des C#-SDK
# (sdks/csharp/Dockerfile), Start je Phase im selben Docker-Netz wie der
# Feed-Container; die stdout-Marker (READY/RECEIVED/REJECTED) liest dieser
# Runner über `docker logs`, während der Test noch läuft. Die `change_id`
# der empfangenen Change wird zusätzlich gegen den Lesezugriffsweg
# `cdc.changes` gehalten (gRPC/SSE/NATS), die Consumer-Registrierung gegen
# `cdc.consumer` (HTTP). Die vier Regel-Phasen (eigene Tabelle
# `feed_e2e_sdkrule`, Vorbereitung über tools/harness/lib-sdk-rule-fixture.sh)
# tragen keinen REJECTED-Beleg — das Negativ steht bei den vier
# Server-ohne-Regel-Phasen — und halten die `change_id` gegen `cdc.changes`
# mit dem umbenannten Zielschlüssel.
#
# Der Runner schreibt den Abdeckungs-Träger docs/user/sdk-e2e-abdeckung.md
# (C#-Abschnitt, marker-gegrenzt) aus derselben Messung, die ihn belegt —
# idempotent, nur bei inhaltlicher Abweichung (Muster
# run-integration-tests.shs Abdeckungstabelle). Grenze der Writer-Form:
# dieser Runner erzeugt Kopf und alles bis zu seinem end-Marker neu und
# erhält den Rest; ein Folge-Runner (Kotlin/Python-HTTP) spiegelt dieselbe
# Form NICHT wortgleich — sein Re-Run würde die Abschnitte oberhalb seines
# begin-Markers wegwerfen. Er hält den Kopf stabil und ersetzt nur seinen
# eigenen Abschnitt (Rest-Erhalt hinter dem end-Marker, wie hier).
#
# Voraussetzungen: Docker, ein geladenes :dev-Image (`make image` vorher —
# compose.yaml trägt keinen build:-Block, ADR-0044) und Netz (NuGet-Restore
# im SDK-Bau, d-migrate-Image). Kein Gate (dieselbe Klasse wie
# `make test-integration`). Cleanup je Ausgang: Test-Container wegwerfen,
# compose-Stack inkl. Volume abräumen.
#
# Aufruf: `make test-sdk-csharp-integration`. Override:
# SDK_CSHARP_INTEGRATION_IMAGE.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

# shellcheck source=tools/harness/lib-sdk-rule-fixture.sh
source tools/harness/lib-sdk-rule-fixture.sh
# shellcheck source=tools/harness/lib-sdk-route-fixture.sh
source tools/harness/lib-sdk-route-fixture.sh

COMPOSE=${COMPOSE:-docker compose -f compose.yaml}
NETWORK=${NETWORK:-cdc-feed-test}
PG_CONTAINER=${PG_CONTAINER:-cdc-test-postgres}
FEED_CONTAINER=${FEED_CONTAINER:-cdc-test-feed}
PG_DB=cdc
PG_USER=postgres
PG_PASSWORD=postgres
# Derselbe Slot-Name wie der Container-Vertrag in compose.yaml (CDC_SLOT).
SLOT=slot_pgc_e2e
API_TOKEN=e2e-reader-token
API_TOKEN_ADMIN=e2e-admin-token
API_TOKEN_READER=e2e-reader-token
NATS_STREAM_TOKEN=e2e-nats-stream-token
SOURCE_ID=src-e2e
HTTP_PUBLICATION=pub_pgc_e2e
DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"

SDK_INTEGRATION_IMAGE=${SDK_CSHARP_INTEGRATION_IMAGE:-pg-change-feed:sdk-csharp-integration}
SDK_TEST_CONTAINER=cdc-sdk-csharp-client-test

TEST_TABLE=feed_e2e_full
ABDECKUNG_ZIEL=docs/user/sdk-e2e-abdeckung
ABDECKUNG_ZIEL_DATEI=${ABDECKUNG_ZIEL}.md
GRPC_TEST_NAME=GrpcRealserverTests
SSE_TEST_NAME=SseRealserverTests
NATS_TEST_NAME=NatsRealserverTests
HTTP_TEST_NAME=HttpRealserverTests
GRPC_SENTINEL=CsharpGrpcSdkE2ESentinel
SSE_SENTINEL=CsharpSseSdkE2ESentinel
NATS_SENTINEL=CsharpNatsSdkE2ESentinel
HTTP_SENTINEL=CsharpHttpSdkE2ESentinel

cleanup() {
  docker rm -fv "$SDK_TEST_CONTAINER" >/dev/null 2>&1 || true
  $COMPOSE down -v --remove-orphans >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker image inspect ghcr.io/pt9912/pg-change-feed:dev >/dev/null 2>&1 || {
  echo "run-sdk-csharp-integration-tests: kein geladenes :dev-Image — make image vorher (compose.yaml trägt keinen build:-Block, ADR-0044)" >&2
  exit 1
}

$COMPOSE down -v --remove-orphans >/dev/null 2>&1 || true
$COMPOSE up -d postgres >/dev/null

ready=0
for _ in $(seq 1 60); do
  if docker exec "$PG_CONTAINER" pg_isready -U "$PG_USER" -d "$PG_DB" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done
if [ "$ready" -ne 1 ]; then
  echo "run-sdk-csharp-integration-tests: Compose-PostgreSQL wurde nicht bereit" >&2
  exit 1
fi

bash tools/schema/rollout-restore.sh make schema-rollout SCHEMA_TARGET="db:$DSN" SCHEMA_ROLLOUT_NETWORK="$NETWORK"

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 >/dev/null <<'SQL'
CREATE TABLE public.feed_e2e_flow (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_full (id int PRIMARY KEY, name text);
ALTER TABLE public.feed_e2e_full REPLICA IDENTITY FULL;
CREATE TABLE public.feed_e2e_schema (id int PRIMARY KEY, name text, amount text);
INSERT INTO cdc.source (source_id, name) VALUES ('src-e2e', 'E2E-Quelle');
SQL

$COMPOSE up -d pg-change-feed >/dev/null

wired=0
for _ in $(seq 1 60); do
  slot_ok=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT 1 FROM pg_replication_slots WHERE slot_name = '$SLOT'" 2>/dev/null || true)
  feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
  if [ "$slot_ok" = "1" ] && [ "$feed_running" = "true" ]; then
    wired=1
    break
  fi
  sleep 1
done
if [ "$wired" -ne 1 ]; then
  feed_exit=$(docker inspect --format '{{.State.ExitCode}}' "$FEED_CONTAINER" 2>/dev/null || echo läuft)
  echo "run-sdk-csharp-integration-tests: Feed-Container trägt die Verdrahtung nicht — Slot $SLOT: ${slot_ok:-0}, Feed-Lauf: $feed_running (Feed-Ausgang: $feed_exit)" >&2
  exit 1
fi

healthy=0
for _ in $(seq 1 60); do
  health=$(docker inspect --format '{{.State.Health.Status}}' "$FEED_CONTAINER" 2>/dev/null || echo fehlt)
  if [ "$health" = "healthy" ]; then
    healthy=1
    break
  fi
  sleep 1
done
if [ "$healthy" -ne 1 ]; then
  echo "run-sdk-csharp-integration-tests: Feed-Container meldet Compose-Health-Status ${health:-fehlt}, wollen healthy" >&2
  exit 1
fi

docker build --build-context proto=proto -f sdks/csharp/Dockerfile \
  --target integration -t "$SDK_INTEGRATION_IMAGE" sdks/csharp

# run_phase <Name> <Testklasse> <Sentinel> <ID-Basis> <Reject-Marker>
# <Received-Grep> <SQL-Variante: changes|consumer> <Extra-Env> — ein
# run_phase <Name> <Testklasse> <Sentinel> <ID-Basis> <Reject-Marker>
# <Received-Grep> <SQL-Variante: changes|consumer|changes_renamed> <Tabelle>
# <Extra-Env> — ein Realserver-Rundlauf für genau eine SDK-Fläche: Container
# starten, auf READY warten, begrenzte Folge eindeutiger Zeilen committen
# (Fire-and-Forget-Fenster), bei nicht-leerem Reject-Marker auf ihn warten,
# auf das Prozessende warten, die RECEIVED-Zeile prüfen und die Identität
# unabhängig gegen den SQL-Lesezugriffsweg halten. Ein leerer Reject-Marker
# überspringt den Ablehnungs-Beleg (die vier Regel-Phasen tragen ihn nicht —
# das Negativ steht bei den vier Server-ohne-Regel-Phasen). Die SQL-Variante
# "changes" hält die change_id gegen cdc.changes (Streams), "consumer" die
# Consumer-Registrierung gegen cdc.consumer (HTTP), "changes_renamed" wie
# "changes", aber am umbenannten Zielschlüssel ohne den Quellschlüssel
# (tools/harness/lib-sdk-rule-fixture.sh).
run_phase() {
  local phase_name=$1 test_name=$2 sentinel=$3 id_base=$4 reject_marker=$5 \
        received_grep=$6 sql_kind=$7 table=$8 extra_env=$9
  local attempt insert_id captured ident pair

  local env_args=()
  for pair in $extra_env; do
    env_args+=(-e "$pair")
  done

  docker rm -fv "$SDK_TEST_CONTAINER" >/dev/null 2>&1 || true
  docker run -d --name "$SDK_TEST_CONTAINER" --network "$NETWORK" \
    "${env_args[@]}" \
    -e PGCHANGEFEED_TEST_NAME="$test_name" \
    -e PGCHANGEFEED_E2E_TABLE="$table" \
    -e PGCHANGEFEED_E2E_SENTINEL="$sentinel" \
    "$SDK_INTEGRATION_IMAGE" >/dev/null

  test_ready=0
  for _ in $(seq 1 60); do
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
    echo "run-sdk-csharp-integration-tests: $phase_name — Test wurde nicht innerhalb der Zeitspanne bereit (kein READY): $(docker logs "$SDK_TEST_CONTAINER" 2>&1 || true)" >&2
    exit 1
  fi

  local received=0
  for attempt in 1 2 3 4 5; do
    insert_id=$((id_base + attempt))
    # >/dev/null: die psql-INSERT-Echo-Zeile gehört nicht in den
    # Funktions-stdout (der Aufrufer hält hier nur den Rückgabewert).
    docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 >/dev/null <<SQL
INSERT INTO public.$table (id, name) VALUES ($insert_id, '$sentinel');
SQL
    for _ in $(seq 1 20); do
      if docker logs "$SDK_TEST_CONTAINER" 2>/dev/null | grep -qF "RECEIVED"; then
        received=1
        break
      fi
      if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
        break
      fi
      sleep 0.5
    done
    if [ "$received" -eq 1 ]; then
      break
    fi
  done

  local rejected=1
  if [ -n "$reject_marker" ]; then
    rejected=0
    for _ in $(seq 1 40); do
      if docker logs "$SDK_TEST_CONTAINER" 2>/dev/null | grep -qF "$reject_marker"; then
        rejected=1
        break
      fi
      if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
        break
      fi
      sleep 0.5
    done
  fi

  local test_stopped=0
  for _ in $(seq 1 20); do
    if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
      test_stopped=1
      break
    fi
    sleep 0.5
  done
  local test_exit
  test_exit=$(docker inspect --format '{{.State.ExitCode}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo unbekannt)
  local test_output
  test_output=$(docker logs "$SDK_TEST_CONTAINER" 2>&1 || true)

  if [ "$received" -ne 1 ]; then
    echo "run-sdk-csharp-integration-tests: $phase_name — der Test empfing keine der committeten Änderungen ($table, $sentinel): $test_output" >&2
    exit 1
  fi
  if ! printf '%s' "$test_output" | grep -qE "$received_grep"; then
    echo "run-sdk-csharp-integration-tests: $phase_name — die RECEIVED-Zeile trägt nicht die erwartete Form: $test_output" >&2
    exit 1
  fi
  if [ -n "$reject_marker" ] && [ "$rejected" -ne 1 ]; then
    echo "run-sdk-csharp-integration-tests: $phase_name — der Ablehnungs-Beleg blieb aus ($reject_marker fehlt): $test_output" >&2
    exit 1
  fi
  if [ "$test_stopped" -ne 1 ] || [ "$test_exit" != "0" ]; then
    echo "run-sdk-csharp-integration-tests: $phase_name — der Test endete nicht mit Ausgang 0 (gestoppt: $test_stopped, Ausgang: $test_exit): $test_output" >&2
    exit 1
  fi

  ident=$(printf '%s' "$test_output" | grep -oE 'RECEIVED [a-z_]+=[^ ]+' | head -n1 | cut -d= -f2 || true)
  if [ -z "$ident" ]; then
    echo "run-sdk-csharp-integration-tests: $phase_name — die RECEIVED-Zeile trägt keine Identität: $test_output" >&2
    exit 1
  fi
  if [ "$sql_kind" = "changes" ]; then
    captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$ident' AND table_name = '$table' AND new_data->>'name' = '$sentinel'")
  elif [ "$sql_kind" = "changes_renamed" ]; then
    captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$ident' AND table_name = '$table' AND new_data->>'$SDK_RULE_TARGET_KEY' = '$sentinel' AND NOT jsonb_exists(new_data, '$SDK_RULE_SOURCE_KEY')")
  else
    captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT count(*) FROM cdc.consumer WHERE consumer_id = '$ident'")
  fi
  if [ -z "$captured" ] || [ "$captured" -lt 1 ]; then
    echo "run-sdk-csharp-integration-tests: $phase_name — die Identität ($ident) ist nicht real über den SQL-Lesezugriffsweg lesbar (count=${captured:-leer})" >&2
    exit 1
  fi

  feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
  if [ "$feed_running" != "true" ]; then
    echo "run-sdk-csharp-integration-tests: $phase_name — Feed-Container lief nach dem Lauf nicht mehr weiter (kein Neustart erwartet)" >&2
    exit 1
  fi

  printf '%s\n' "$ident"
}

GRPC_IDENT=$(run_phase \
  "gRPC-Fläche (SPEC-020)" "$GRPC_TEST_NAME" \
  "$GRPC_SENTINEL" 400 "REJECTED code=Unauthenticated" \
  "RECEIVED change_id=[^ ]+ table=$TEST_TABLE .*operation=INSERT .*new_image=.*$GRPC_SENTINEL" \
  changes "$TEST_TABLE" "PGCHANGEFEED_GRPC_ADDR=pg-change-feed:9090 PGCHANGEFEED_API_TOKEN=$API_TOKEN")

SSE_IDENT=$(run_phase \
  "SSE-Fläche (SPEC-021)" "$SSE_TEST_NAME" \
  "$SSE_SENTINEL" 410 "REJECTED status=401" \
  "RECEIVED change_id=[^ ]+ table=$TEST_TABLE .*operation=INSERT .*new_image=.*$SSE_SENTINEL" \
  changes "$TEST_TABLE" "PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090 PGCHANGEFEED_API_TOKEN=$API_TOKEN")

NATS_IDENT=$(run_phase \
  "NATS-Vollinhalts-Fläche (SPEC-024)" "$NATS_TEST_NAME" \
  "$NATS_SENTINEL" 420 "REJECTED token-rejected" \
  "RECEIVED change_id=[^ ]+ table=$TEST_TABLE .*operation=INSERT .*new_image=.*$NATS_SENTINEL" \
  changes "$TEST_TABLE" "PGCHANGEFEED_NATS_URL=nats://nats:4222 PGCHANGEFEED_NATS_STREAM_TOKEN=$NATS_STREAM_TOKEN PGCHANGEFEED_SOURCE_ID=$SOURCE_ID")

HTTP_IDENT=$(run_phase \
  "HTTP-Fläche (SPEC-018)" "$HTTP_TEST_NAME" \
  "$HTTP_SENTINEL" 430 "REJECTED status=401" \
  "RECEIVED consumer_id=[^ ]+" \
  consumer "$TEST_TABLE" "PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090 PGCHANGEFEED_API_TOKEN_ADMIN=$API_TOKEN_ADMIN PGCHANGEFEED_API_TOKEN_READER=$API_TOKEN_READER PGCHANGEFEED_SOURCE_ID=$SOURCE_ID PGCHANGEFEED_HTTP_PUBLICATION=$HTTP_PUBLICATION")

# --- Regel-Phasen (slice-sdk-regel-realserver-e2e, ADR-0112) --------------
# Dieselben vier Flächen, ein zweites Mal gegen die eigene Tabelle
# $SDK_RULE_TABLE mit aktiver rename_column-Regel
# (tools/harness/lib-sdk-rule-fixture.sh): kein REJECTED-Beleg (das Negativ
# steht bei den vier Phasen oben), sql_kind changes_renamed hält die
# change_id gegen den umbenannten Zielschlüssel ohne den Quellschlüssel.
sdk_rule_fixture_setup "run-sdk-csharp-integration-tests"

GRPC_RULE_TEST_NAME=GrpcRuleRealserverTests
SSE_RULE_TEST_NAME=SseRuleRealserverTests
NATS_RULE_TEST_NAME=NatsRuleRealserverTests
HTTP_RULE_TEST_NAME=HttpRuleRealserverTests
GRPC_RULE_SENTINEL=CsharpGrpcRuleSdkE2ESentinel
SSE_RULE_SENTINEL=CsharpSseRuleSdkE2ESentinel
NATS_RULE_SENTINEL=CsharpNatsRuleSdkE2ESentinel
HTTP_RULE_SENTINEL=CsharpHttpRuleSdkE2ESentinel

GRPC_RULE_IDENT=$(run_phase \
  "gRPC-Regel-Fläche (SPEC-020, ADR-0112)" "$GRPC_RULE_TEST_NAME" \
  "$GRPC_RULE_SENTINEL" 500 "" \
  "RECEIVED change_id=[^ ]+ table=$SDK_RULE_TABLE .*operation=INSERT .*new_image=.*$GRPC_RULE_SENTINEL" \
  changes_renamed "$SDK_RULE_TABLE" \
  "PGCHANGEFEED_GRPC_ADDR=pg-change-feed:9090 PGCHANGEFEED_API_TOKEN=$API_TOKEN PGCHANGEFEED_RULE_SOURCE_KEY=$SDK_RULE_SOURCE_KEY PGCHANGEFEED_RULE_TARGET_KEY=$SDK_RULE_TARGET_KEY")

SSE_RULE_IDENT=$(run_phase \
  "SSE-Regel-Fläche (SPEC-021, ADR-0112)" "$SSE_RULE_TEST_NAME" \
  "$SSE_RULE_SENTINEL" 510 "" \
  "RECEIVED change_id=[^ ]+ table=$SDK_RULE_TABLE .*operation=INSERT .*new_image=.*$SSE_RULE_SENTINEL" \
  changes_renamed "$SDK_RULE_TABLE" \
  "PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090 PGCHANGEFEED_API_TOKEN=$API_TOKEN PGCHANGEFEED_RULE_SOURCE_KEY=$SDK_RULE_SOURCE_KEY PGCHANGEFEED_RULE_TARGET_KEY=$SDK_RULE_TARGET_KEY")

NATS_RULE_IDENT=$(run_phase \
  "NATS-Vollinhalts-Regel-Fläche (SPEC-024, ADR-0112)" "$NATS_RULE_TEST_NAME" \
  "$NATS_RULE_SENTINEL" 520 "" \
  "RECEIVED change_id=[^ ]+ table=$SDK_RULE_TABLE .*operation=INSERT .*new_image=.*$NATS_RULE_SENTINEL" \
  changes_renamed "$SDK_RULE_TABLE" \
  "PGCHANGEFEED_NATS_URL=nats://nats:4222 PGCHANGEFEED_NATS_STREAM_TOKEN=$NATS_STREAM_TOKEN PGCHANGEFEED_SOURCE_ID=$SOURCE_ID PGCHANGEFEED_RULE_SOURCE_KEY=$SDK_RULE_SOURCE_KEY PGCHANGEFEED_RULE_TARGET_KEY=$SDK_RULE_TARGET_KEY")

HTTP_RULE_IDENT=$(run_phase \
  "HTTP-Regel-Fläche (SPEC-018, ADR-0112)" "$HTTP_RULE_TEST_NAME" \
  "$HTTP_RULE_SENTINEL" 530 "" \
  "RECEIVED change_id=[^ ]+ table=$SDK_RULE_TABLE .*operation=INSERT .*new_image=.*$HTTP_RULE_SENTINEL" \
  changes_renamed "$SDK_RULE_TABLE" \
  "PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090 PGCHANGEFEED_API_TOKEN=$API_TOKEN PGCHANGEFEED_SOURCE_ID=$SOURCE_ID PGCHANGEFEED_RULE_SOURCE_KEY=$SDK_RULE_SOURCE_KEY PGCHANGEFEED_RULE_TARGET_KEY=$SDK_RULE_TARGET_KEY")

# --- Routing-Phasen (ADR-0137 Teilfrage 5) --------------------------------
# Dieselben vier Flächen gegen die eigene Tabelle $SDK_ROUTE_TABLE mit zwei
# Routing-Regeln (tools/harness/lib-sdk-route-fixture.sh): der Client mit
# `target` empfängt nur die Changes seines Ziels, der Client ohne `target`
# alle; die Kennungen werden gegen cdc.changes (route_target) gehalten.
sdk_route_fixture_setup "run-sdk-csharp-integration-tests"

sdk_route_phase "run-sdk-csharp-integration-tests" \
  "gRPC-Routing-Fläche (SPEC-020, ADR-0137)" "CsharpGrpcRouteSdkE2ESentinel" 600 \
  "PGCHANGEFEED_TEST_NAME=GrpcRouteRealserverTests" \
  "PGCHANGEFEED_GRPC_ADDR=pg-change-feed:9090 PGCHANGEFEED_API_TOKEN=$API_TOKEN" \
  "$SDK_ROUTE_QUIET_SECONDS"

sdk_route_phase "run-sdk-csharp-integration-tests" \
  "SSE-Routing-Fläche (SPEC-021, ADR-0137)" "CsharpSseRouteSdkE2ESentinel" 700 \
  "PGCHANGEFEED_TEST_NAME=SseRouteRealserverTests" \
  "PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090 PGCHANGEFEED_API_TOKEN=$API_TOKEN" \
  "$SDK_ROUTE_QUIET_SECONDS"

sdk_route_phase "run-sdk-csharp-integration-tests" \
  "NATS-Routing-Fläche (SPEC-024, ADR-0137)" "CsharpNatsRouteSdkE2ESentinel" 800 \
  "PGCHANGEFEED_TEST_NAME=NatsRouteRealserverTests" \
  "PGCHANGEFEED_NATS_URL=nats://nats:4222 PGCHANGEFEED_NATS_STREAM_TOKEN=$NATS_STREAM_TOKEN PGCHANGEFEED_SOURCE_ID=$SOURCE_ID" \
  "$SDK_ROUTE_QUIET_SECONDS"

sdk_route_phase "run-sdk-csharp-integration-tests" \
  "HTTP-Routing-Fläche (SPEC-018, ADR-0137)" "CsharpHttpRouteSdkE2ESentinel" 900 \
  "PGCHANGEFEED_TEST_NAME=HttpRouteRealserverTests" \
  "PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090 PGCHANGEFEED_API_TOKEN=$API_TOKEN PGCHANGEFEED_SOURCE_ID=$SOURCE_ID" \
  0

# --- Abdeckungs-Träger (docs/user/sdk-e2e-abdeckung.md) -------------------
# Der C#-Abschnitt entsteht aus derselben Messung, die ihn belegt: je Phase
# ihre Kennungen und der ident-Behalt des Laufs; geschrieben wird nur bei
# inhaltlicher Abweichung (Temp-Datei + cmp), marker-gegrenzt, damit die
# Folge-Slices (Kotlin, Python-HTTP) eigene Abschnitte daneben tragen.
# Die Kennungsspalte trägt je Kennung den Link auf ihr Definitionsdokument
# (ids des Doku-Gates).

# trim_blank_edges <mehrzeiliger String> — entfernt führende und
# nachgestellte Leerzeilen, interne Leerzeilen bleiben erhalten; reines awk,
# kein externes Werkzeug.
trim_blank_edges() {
  awk '
    { lines[NR] = $0 }
    END {
      start = 1
      while (start <= NR && lines[start] == "") start++
      e = NR
      while (e >= start && lines[e] == "") e--
      for (i = start; i <= e; i++) print lines[i]
    }
  ' <<<"$1"
}

abdeckung_csharp_abschnitt() {
  printf '%s\n' \
    '<!-- pgchangefeed-sdk-e2e:csharp-begin -->' \
    '| Spec-Kennung | Kurzbeschreibung | Nachweis | Ort |' \
    '| --- | --- | --- | --- |' \
    "| [\`LH-FA-SST-008\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein C#-SDK-Client (\`PgChangeFeedGrpcClient\`) öffnet real den gRPC-Server-Stream gegen den laufenden Feed-Container und empfängt eine danach committete Änderung; ein Öffnungsversuch ohne gültiges Token endet mit gRPC-Status \`Unauthenticated\` | \`GrpcRealserverTests\` | \`tools/harness/run-sdk-csharp-integration-tests.sh\` |" \
    "| [\`LH-FA-SST-008\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein C#-SDK-Client (\`PgChangeFeedSseClient\`) öffnet real \`GET /changes/stream\` und empfängt eine danach committete Änderung; ein Aufruf ohne gültiges Token endet mit HTTP-Status 401 | \`SseRealserverTests\` | \`tools/harness/run-sdk-csharp-integration-tests.sh\` |" \
    "| [\`LH-FA-SST-008\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein C#-SDK-Client (\`PgChangeFeedNatsStreamClient\`) verbindet sich real per NATS und empfängt eine danach committete Änderung als vollständiges JSON-Event; ein Verbindungsversuch mit falschem Token wird vom NATS-Server abgelehnt | \`NatsRealserverTests\` | \`tools/harness/run-sdk-csharp-integration-tests.sh\` |" \
    "| [\`LH-FA-SST-006\`](../../spec/lastenheft.md), [\`LH-FA-CON-001\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein C#-SDK-Client (\`PgChangeFeedHttpClient\`) registriert real einen Consumer (admin-Token) und listet Tabellen (reader-Token); die Registrierung ist unabhängig über \`cdc.consumer\` lesbar; ein Aufruf ohne gültiges Token endet mit HTTP-Status 401 | \`HttpRealserverTests\` | \`tools/harness/run-sdk-csharp-integration-tests.sh\` |" \
    "| [\`LH-FA-CFG-007\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | eine aktive \`rename_column\`-Regel auf einer eigenen Tabelle: alle vier C#-SDK-Clients empfangen die danach erfasste Änderung mit dem umbenannten Schlüssel im opaken Bild-Modell (\`JsonElement\`, gRPC die Bytes als JSON) — Zielschlüssel trägt den Sentinel, Quellschlüssel fehlt; \`change_id\` je unabhängig über \`cdc.changes\` lesbar | \`GrpcRuleRealserverTests\`, \`SseRuleRealserverTests\`, \`NatsRuleRealserverTests\`, \`HttpRuleRealserverTests\` | \`tools/harness/run-sdk-csharp-integration-tests.sh\` |" \
    "| [\`LH-FA-CFG-008\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | zwei aktive Routing-Regeln (Ziele \`eu\` und \`us\`) auf einer eigenen Tabelle: die vier C#-SDK-Flächen (gRPC-Stream, SSE-Stream, NATS-Zusatz-Subjekt \`cdc.route.<source_id>.<ziel>\`, HTTP-Lesezugriff) liefern mit \`target\` genau die Changes des Ziels \`eu\` und im Ruhefenster keine Change eines anderen Ziels oder ohne Ziel (der HTTP-Lesezugriff liefert exakt die Ziel-Teilmenge der ungefilterten Lesung), ohne \`target\` alle; \`change_id\` je unabhängig über \`cdc.changes\` (\`route_target\`) gegengelesen | \`GrpcRouteRealserverTests\`, \`SseRouteRealserverTests\`, \`NatsRouteRealserverTests\`, \`HttpRouteRealserverTests\` | \`tools/harness/run-sdk-csharp-integration-tests.sh\` |" \
    '<!-- pgchangefeed-sdk-e2e:csharp-end -->'
}

abdeckung_schreiben() {
  local temp fremde_abschnitte=""
  local abschnitt
  abschnitt=$(abdeckung_csharp_abschnitt)
  # Fremde, marker-gegrenzte Abschnitte (Kotlin/Python-HTTP der Folge-Slices)
  # hinter dem eigenen end-Marker bleiben erhalten — jeder Runner ersetzt nur
  # seinen eigenen Abschnitt. Jeder Sprachabschnitt trägt seine eigene
  # Kopf-/Trennzeile (eine eigenständige Markdown-Tabelle je Sprache): eine
  # HTML-Kommentarzeile ohne Pipe-Zeichen zwischen einer gemeinsamen
  # Kopfzeile und den Datenzeilen hätte die Tabelle für jeden
  # Markdown-Renderer nach der Kopfzeile beendet.
  local rest=""
  if [ -f "$ABDECKUNG_ZIEL_DATEI" ]; then
    rest=$(awk '/pgchangefeed-sdk-e2e:csharp-end/{gefunden=1; next} gefunden{print}' "$ABDECKUNG_ZIEL_DATEI")
  fi
  rest=$(trim_blank_edges "$rest")
  temp=$(mktemp)
  {
    printf '%s\n' \
      '# SDK-E2E-Abdeckung je Spec-Kennung' \
      '' \
      'Erzeugt von `make test-sdk-csharp-integration` über' \
      '`tools/harness/run-sdk-csharp-integration-tests.sh`; die Sprach-Runner' \
      'der Folge-Slices (Kotlin, Python-HTTP) erweitern dieselbe Datei um ihre' \
      'marker-gegrenzten Abschnitte. Je Sprach-Abschnitt deklariert der' \
      'zuständige Runner seine Realserver-Phasen an Ort und Stelle. Diese' \
      'Datei ist eine **stabile Abdeckungs-Deklaration**, kein Lauf-Beleg: der' \
      'Runner schreibt sie nur bei inhaltlicher Abweichung. Sie trägt nur' \
      'Zeilen real existierender Runner-Phasen — ein Beleg steht hier nie,' \
      'bevor sein Lauf grün lief.'
    printf '\n%s\n' "$abschnitt"
    if [ -n "$rest" ]; then
      printf '\n%s\n' "$rest"
    fi
  } > "$temp"
  if [ -f "$ABDECKUNG_ZIEL_DATEI" ] && cmp -s "$temp" "$ABDECKUNG_ZIEL_DATEI"; then
    rm -f "$temp"
    echo "run-sdk-csharp-integration-tests: Abdeckungs-Träger unverändert — $ABDECKUNG_ZIEL_DATEI entspricht dem Quelltext-Stand"
  else
    chmod 0644 "$temp"
    mv "$temp" "$ABDECKUNG_ZIEL_DATEI"
    echo "run-sdk-csharp-integration-tests: Abdeckungs-Träger geschrieben — $ABDECKUNG_ZIEL_DATEI"
  fi
}

abdeckung_schreiben

echo "run-sdk-csharp-integration-tests: C#-SDK-Realserver-Belege (slice-sdk-csharp-reale2e, Mechanik-Klasse ADR-0110 Festlegung 2) grün — gRPC-Fläche (PgChangeFeedGrpcClient, pg-change-feed:9090, change_id=$GRPC_IDENT), SSE-Fläche (PgChangeFeedSseClient, pg-change-feed:8090, change_id=$SSE_IDENT) und NATS-Vollinhalts-Fläche (PgChangeFeedNatsStreamClient, nats://nats:4222, change_id=$NATS_IDENT) empfingen je eine danach committete Änderung (change_id je unabhängig über cdc.changes lesbar), die HTTP-Fläche (PgChangeFeedHttpClient) registrierte real einen Consumer (consumer_id=$HTTP_IDENT, unabhängig über cdc.consumer lesbar) und listete Tabellen; ein Aufruf ohne gültiges Token endete je mit gRPC-Status Unauthenticated, HTTP-Status 401 bzw. der laut ablehnenden NATS-Verbindungsablehnung"
echo "run-sdk-csharp-integration-tests: Regel-Belege (slice-sdk-regel-realserver-e2e, ADR-0112) grün — eine aktive rename_column-Regel auf $SDK_RULE_TABLE, alle vier Flächen empfingen die danach erfasste Änderung mit dem umbenannten Schlüssel $SDK_RULE_TARGET_KEY (Quellschlüssel $SDK_RULE_SOURCE_KEY fehlt): gRPC change_id=$GRPC_RULE_IDENT, SSE change_id=$SSE_RULE_IDENT, NATS change_id=$NATS_RULE_IDENT, HTTP change_id=$HTTP_RULE_IDENT (je unabhängig über cdc.changes lesbar)"
echo "run-sdk-csharp-integration-tests: Routing-Belege (ADR-0137 Teilfrage 5) grün — zwei Routing-Regeln auf $SDK_ROUTE_TABLE (Ziel $SDK_ROUTE_TARGET_A und $SDK_ROUTE_TARGET_B), je Fläche ein Client mit target und ein Client ohne target${SDK_ROUTE_REPORT}"