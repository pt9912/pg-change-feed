#!/usr/bin/env bash
# run-sdk-altserver-tests.sh — die Fehlerfälle der 0.6.0-SDKs gegen einen Server
# vor 0.6.0 (make test-sdk-altserver; Vertrag: harness/targets/sdk-altserver.md).
# Dieselbe Compose-Umgebung wie die SDK-Realserver-Runner (PostgreSQL, NATS,
# Schema-Rollout des Arbeitsbaums über d-migrate); nur der Feed-Container ist das
# Image SDK_ALTSERVER_IMAGE (Standard: das Release 0.5.0 mit Tag und Index-Digest).
# Der Tag bleibt in der Referenz, damit `make pin-stale-all` den Digest gegen den
# Tag 0.5.0 vergleicht. Das Image kommt über eine Override-Datei im
# Temp-Verzeichnis; compose.yaml bleibt unverändert.
#
#   B0  der Server läuft gegen das aktuelle Schema: Health healthy, Slot da, eine
#       eingefügte Zeile einer aktivierten Tabelle liegt über cdc.changes
#   B1  Rohdraht-Probe der HTTP-Seite: der Fehlerkörper der Aktivierung einer
#       fehlenden Tabelle trägt kein Feld `code` (ohne SDK, aus einem Container
#       im Compose-Netz)
#   B2  je Sprache (C#, Kotlin, Python) eine Test-Klasse der SDK-Quelle des
#       Arbeitsstands: HTTP und gRPC enden typisiert, der Meldungscode bleibt leer,
#       der Fehlertext ist da (HTTP: gleich dem Text aus B1), das Reader-Token
#       endet mit 403/PermissionDenied, die Diagnose liefert den Bericht ohne
#       error_code — im Normalbetrieb und in einem Fehlerzustand, den dieser Runner
#       in cdc.process_heartbeat ohne error_code schreibt
#   B3  Negativprobe: SDK_ALTSERVER_IMAGE=<geladenes :dev-Image> färbt den Lauf an
#       B1 rot; mit SDK_ALTSERVER_WEITER=1 läuft der Runner nach einem roten B1
#       weiter und zeigt, dass auch die Test-Klassen von B2 am neuen Server rot
#       werden
#
# Der Runner schreibt nichts in den Arbeitsbaum (kein Abdeckungs-Träger). Er
# bricht ab, wenn die Container- oder Netznamen der anderen Runner belegt sind;
# aufgeräumt wird nur, was er selbst angelegt hat. Kein Gate (DB-Zugang, Docker,
# Netz). Aufruf: `make test-sdk-altserver`.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

SDK_ALTSERVER_IMAGE=${SDK_ALTSERVER_IMAGE:-ghcr.io/pt9912/pg-change-feed:0.5.0@sha256:f99a77ff335a5fa2e84937771337d03711b6ad4b4ac9fb038dd6e310778707ce}
SDK_ALTSERVER_WEITER=${SDK_ALTSERVER_WEITER:-0}

NETWORK=cdc-feed-test
PG_CONTAINER=cdc-test-postgres
NATS_CONTAINER=cdc-test-nats
FEED_CONTAINER=cdc-test-feed
PG_DB=cdc
PG_USER=postgres
PG_PASSWORD=postgres
SLOT=slot_pgc_e2e
SOURCE_ID=src-e2e
HTTP_PUBLICATION=pub_pgc_e2e
API_TOKEN_ADMIN=e2e-admin-token
API_TOKEN_READER=e2e-reader-token
DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"
TEST_CONTAINER=cdc-sdk-altserver-test
PREFIX=run-sdk-altserver-tests

CS_IMAGE=pg-change-feed:sdk-altserver-csharp
KT_IMAGE=pg-change-feed:sdk-altserver-kotlin
PY_IMAGE=pg-change-feed:sdk-altserver-python

workdir=$(mktemp -d)
COMPOSE="docker compose -f compose.yaml -f $workdir/override.yaml"
angelegt=0

cleanup() {
  docker rm -fv "$TEST_CONTAINER" >/dev/null 2>&1 || true
  if [ "$angelegt" -eq 1 ]; then
    $COMPOSE down -v --remove-orphans >/dev/null 2>&1 || true
  fi
  rm -rf -- "$workdir"
}
trap cleanup EXIT

psql_exec() {
  docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 "$@"
}

rot=0
befund() {
  echo "$PREFIX: $1" >&2
  rot=1
  if [ "$SDK_ALTSERVER_WEITER" != 1 ]; then
    exit 1
  fi
}

# Die Namen der anderen Runner dürfen nicht belegt sein.
for name in "$PG_CONTAINER" "$NATS_CONTAINER" "$FEED_CONTAINER" "$TEST_CONTAINER"; do
  if docker container inspect "$name" >/dev/null 2>&1; then
    echo "$PREFIX: der Container $name existiert bereits — ein anderer Runner läuft oder ein früherer Lauf ist nicht abgeräumt" >&2
    exit 1
  fi
done
if docker network inspect "$NETWORK" >/dev/null 2>&1; then
  echo "$PREFIX: das Netz $NETWORK existiert bereits — ein anderer Runner läuft oder ein früherer Lauf ist nicht abgeräumt" >&2
  exit 1
fi

printf 'services:\n  pg-change-feed:\n    image: %s\n' "$SDK_ALTSERVER_IMAGE" >"$workdir/override.yaml"

angelegt=1
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
  echo "$PREFIX: Compose-PostgreSQL wurde nicht bereit" >&2
  exit 1
fi

make schema-rollout SCHEMA_TARGET="db:$DSN" SCHEMA_ROLLOUT_NETWORK="$NETWORK"

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
  echo "$PREFIX: B0 ROT — der Feed-Container trägt die Verdrahtung nicht (Slot $SLOT: ${slot_ok:-0}, Lauf: $feed_running, Ausgang: $feed_exit); Log: $(docker logs "$FEED_CONTAINER" 2>&1 | tail -n 5)" >&2
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
  echo "$PREFIX: B0 ROT — der Feed-Container meldet Health ${health:-fehlt}, wollen healthy; Log: $(docker logs "$FEED_CONTAINER" 2>&1 | tail -n 5)" >&2
  exit 1
fi

psql_exec -c "INSERT INTO public.feed_e2e_full (id, name) VALUES (9001, 'altserver-b0')" >/dev/null
change_id=""
for _ in $(seq 1 40); do
  change_id=$(psql_exec -tAc \
    "SELECT change_id FROM cdc.changes WHERE source_id = '$SOURCE_ID' AND table_name = 'feed_e2e_full' AND new_data->>'name' = 'altserver-b0'")
  if [ -n "$change_id" ]; then
    break
  fi
  sleep 0.5
done
if [ -z "$change_id" ]; then
  echo "$PREFIX: B0 ROT — die eingefügte Zeile von feed_e2e_full erscheint nicht über cdc.changes" >&2
  exit 1
fi
laufende_referenz=$(docker inspect --format '{{.Config.Image}}' "$FEED_CONTAINER")
echo "ALTSERVER B0: healthy, Slot $SLOT, change_id=$change_id, Feed-Image $laufende_referenz"

docker build --build-context proto=proto -f sdks/csharp/Dockerfile --target integration -t "$CS_IMAGE" sdks/csharp >/dev/null
docker build --build-context proto=proto -f sdks/kotlin/Dockerfile --target integration -t "$KT_IMAGE" sdks/kotlin >/dev/null
docker build --build-context proto=proto -f sdks/python/Dockerfile --target integration -t "$PY_IMAGE" sdks/python >/dev/null

# B1: Rohdraht-Probe, kein SDK.
b1_rc=0
b1_out=$(docker run --rm -i --network "$NETWORK" \
  -e ALTSERVER_HTTP_ADDR=http://pg-change-feed:8090 \
  -e ALTSERVER_ADMIN_TOKEN="$API_TOKEN_ADMIN" \
  -e ALTSERVER_SOURCE_ID="$SOURCE_ID" \
  -e ALTSERVER_PUBLICATION="$HTTP_PUBLICATION" \
  "$PY_IMAGE" python - <tools/harness/sdk-kompat/altserver/rohdraht.py 2>&1) || b1_rc=$?
printf '%s\n' "$b1_out" | sed 's/^/ALTSERVER B1: /'
raw_text=$(printf '%s\n' "$b1_out" | sed -n 's/^TEXT //p')
if [ "$b1_rc" -ne 0 ]; then
  befund "B1 ROT — der Fehlerkörper der Aktivierung einer fehlenden Tabelle ist nicht der eines Servers ohne Meldungscodes (Ausgang $b1_rc)"
fi

# b2_phase <Sprache> <Image> <Umgebung Testauswahl> — eine Test-Klasse gegen den Server: auf READY und
# NORMAL_DONE warten, danach den Fehlerzustand ohne error_code in
# cdc.process_heartbeat schreiben, bis der Test endet (der periodische
# Heartbeat setzt ihn zurück), dann Ausgang und RECEIVED-Zeile prüfen.
b2_phase() {
  local sprache=$1 image=$2 auswahl=$3 gesehen state rc out zeile
  docker rm -fv "$TEST_CONTAINER" >/dev/null 2>&1 || true
  # Ein Fehlerzustand der vorigen Phase wird zurückgenommen, damit der Test
  # den Normalbetrieb liest.
  psql_exec -c "UPDATE cdc.process_heartbeat SET error_class = NULL WHERE source_id = '$SOURCE_ID'" >/dev/null
  docker run -d --name "$TEST_CONTAINER" --network "$NETWORK" \
    -e PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090 \
    -e PGCHANGEFEED_GRPC_ADDR=pg-change-feed:9090 \
    -e PGCHANGEFEED_API_TOKEN_ADMIN="$API_TOKEN_ADMIN" \
    -e PGCHANGEFEED_API_TOKEN_READER="$API_TOKEN_READER" \
    -e PGCHANGEFEED_SOURCE_ID="$SOURCE_ID" \
    -e PGCHANGEFEED_HTTP_PUBLICATION="$HTTP_PUBLICATION" \
    -e PGCHANGEFEED_ALTSERVER_HTTP_TEXT="$raw_text" \
    -e "$auswahl" \
    "$image" >/dev/null

  for marker in READY NORMAL_DONE; do
    gesehen=0
    for _ in $(seq 1 180); do
      if docker logs "$TEST_CONTAINER" 2>&1 | grep -F "$marker" >/dev/null; then
        gesehen=1
        break
      fi
      state=$(docker inspect --format '{{.State.Running}}' "$TEST_CONTAINER" 2>/dev/null || echo false)
      if [ "$state" != "true" ]; then
        break
      fi
      sleep 1
    done
    if [ "$gesehen" -ne 1 ]; then
      out=$(docker logs "$TEST_CONTAINER" 2>&1 | tail -n 40)
      befund "B2 $sprache ROT — Marker $marker blieb aus: $out"
      return 0
    fi
  done

  for _ in $(seq 1 200); do
    state=$(docker inspect --format '{{.State.Running}}' "$TEST_CONTAINER" 2>/dev/null || echo false)
    if [ "$state" != "true" ]; then
      break
    fi
    psql_exec -c "UPDATE cdc.process_heartbeat SET heartbeat_at = current_timestamp, error_class = 'schema' WHERE source_id = '$SOURCE_ID'" >/dev/null
    sleep 0.3
  done

  psql_exec -c "UPDATE cdc.process_heartbeat SET error_class = NULL WHERE source_id = '$SOURCE_ID'" >/dev/null
  rc=$(docker inspect --format '{{.State.ExitCode}}' "$TEST_CONTAINER" 2>/dev/null || echo unbekannt)
  out=$(docker logs "$TEST_CONTAINER" 2>&1 || true)
  zeile=$(printf '%s\n' "$out" | grep -oE 'RECEIVED code=none http=404 grpc=NOT_FOUND text=[0-9]+ diag_error_code=leer' | head -n 1 || true)
  if [ "$rc" != "0" ] || [ -z "$zeile" ]; then
    befund "B2 $sprache ROT — Ausgang $rc, RECEIVED-Zeile: ${zeile:-fehlt}; Log: $(printf '%s\n' "$out" | tail -n 30)"
    return 0
  fi
  echo "ALTSERVER B2 $sprache: $zeile (Exit $rc)"
}

b2_phase csharp "$CS_IMAGE" PGCHANGEFEED_TEST_NAME=ErrorCodeAltServerTests
b2_phase kotlin "$KT_IMAGE" PGCHANGEFEED_TEST_NAME=ErrorCodeAltServerTest
b2_phase python "$PY_IMAGE" PGCHANGEFEED_TEST_FILE=integration/test_error_code_altserver.py

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  befund "der Feed-Container lief nach den Phasen nicht mehr weiter (kein Neustart erwartet)"
fi

if [ "$rot" -ne 0 ]; then
  echo "$PREFIX: Altserver-Messung ROT (Feed-Image $laufende_referenz)" >&2
  exit 1
fi
echo "$PREFIX: Altserver-Messung grün — Feed-Image $laufende_referenz: kein code im Fehlerkörper (B1), die Eigenschaft bleibt in allen drei Sprachen leer (B2)"
