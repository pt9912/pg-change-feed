#!/usr/bin/env bash
# run-sdk-python-integration-tests.sh — Realserver-Integrationstest der
# Python-SDK-Zustellweg-Flaechen (ADR-0110 §Entscheidung Festlegung
# 2/Folgepflicht 1, slice-sdk-python-grpc-client-flaeche): jede neue
# Zustellweg-Flaechen dieses Packages prueft ihre Protokoll-Annahmen
# zusaetzlich zu den netzlosen Unit-Tests gegen eine reale, laufende
# Server-Instanz.
#
# Kette je Lauf — Bring-up-Phase im Muster tools/harness/run-integration-tests.sh,
# aber eigenstaendig und bewusst schmaler: dieses Skript traegt keinen
# Server-E2E-Beleg (LH-QA-POR-003 bleibt Eigentum des bestehenden Runners,
# der unveraendert bleibt); es baut die compose.yaml-Umgebung hoch
# (PostgreSQL/NATS/Feed-Container, Schema-Rollout ueber d-migrate,
# Vorbedingungen der Aktivierung), baut die `integration`-Docker-Stufe des
# Python-SDK (sdks/python/Dockerfile, zwingend mit dem benannten Bau-Kontext
# `proto` — ohne ihn bricht der Bau an der COPY --from=proto-Zeile ab) und
# startet sie im selben Docker-Netz wie den Feed-Container. Der Pruefling
# ist das SDK selbst: der Integrationstest importiert
# pgchangefeed.grpc_client direkt (kein Wegwerf-Duplikat-Client daneben).
# Die stdout-Marker des Testlaufs (READY/RECEIVED/REJECTED) liest dieser
# Runner ueber `docker logs`, waehrend der Test noch laeuft — deshalb
# ungepuffert (`python -u`, Stufe-ENTRYPOINT). Die `change_id` der
# empfangenen Change wird zusaetzlich gegen den Lesezugriffsweg
# `cdc.changes` gehalten (dieselbe Disziplin wie der gRPC-Stream-Rundlauf
# im Server-E2E-Runner).
#
# Voraussetzungen: Docker, ein geladenes :dev-Image (`make image` vorher —
# compose.yaml traegt keinen build:-Block, ADR-0044) und Netz (pip-Paketbezug
# im SDK-Bau, d-migrate-Image). Kein Gate (dieselbe Klasse wie
# `make test-integration`). Cleanup je Ausgang: Test-Container wegwerfen,
# compose-Stack inkl. Volume abrueumen.
#
# Aufruf: `make test-sdk-python-integration`. Override:
# SDK_PYTHON_INTEGRATION_IMAGE.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

COMPOSE=${COMPOSE:-docker compose -f compose.yaml}
NETWORK=${NETWORK:-cdc-feed-test}
PG_CONTAINER=${PG_CONTAINER:-cdc-test-postgres}
FEED_CONTAINER=${FEED_CONTAINER:-cdc-test-feed}
PG_DB=cdc
PG_USER=postgres
PG_PASSWORD=postgres
# Derselbe Slot-Name wie der Container-Vertrag in compose.yaml (CDC_SLOT).
SLOT=slot_pgc_e2e
# Derselbe Reader-Token wie der Container-Vertrag in compose.yaml
# (CDC_API_TOKEN_READER): der gRPC-Stream akzeptiert beide Token-Klassen
# (SPEC-020), der Test traegt die Reader-Klasse.
API_TOKEN=e2e-reader-token
DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"

SDK_INTEGRATION_IMAGE=${SDK_PYTHON_INTEGRATION_IMAGE:-pg-change-feed:sdk-python-integration}
SDK_TEST_CONTAINER=cdc-sdk-python-grpc-client-test

# Dasselbe aktivierte Tisch-Set wie der Container-Vertrag (CDC_TABLES):
# die Aktivierung traegt die Verdrahtung des Feed-Containers (ADR-0028),
# die Tabellen muessen dafuer physisch existieren. Gegriffen wird fuer den
# gRPC-Beleg auf feed_e2e_full (REPLICA IDENTITY FULL).
TEST_TABLE=feed_e2e_full
TEST_SENTINEL=PythonGrpcSdkE2ESentinel

cleanup() {
  docker rm -f "$SDK_TEST_CONTAINER" >/dev/null 2>&1 || true
  $COMPOSE down -v --remove-orphans >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker image inspect ghcr.io/pt9912/pg-change-feed:dev >/dev/null 2>&1 || {
  echo "run-sdk-python-integration-tests: kein geladenes :dev-Image — make image vorher (compose.yaml traegt keinen build:-Block, ADR-0044)" >&2
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
  echo "run-sdk-python-integration-tests: Compose-PostgreSQL wurde nicht bereit" >&2
  exit 1
fi

# Schema-Rollout über d-migrate vor dem Lauf (ADR-0043) — dieselbe Kette
# wie der Server-E2E-Runner.
make schema-rollout SCHEMA_TARGET="db:$DSN" SCHEMA_ROLLOUT_NETWORK="$NETWORK"

# Vorbedingungen der Aktivierung: die drei CDC_TABLES-Tabellen existieren
# physisch, die Quelle-Zeile traegt die Fremdschluessel-Registrierung
# (SPEC-001); die Aktivierung selbst (Publication, Bindungs- und
# Schema-Version-Zeilen) traegt die Verdrahtung des Feed-Containers
# (ADR-0028, LH-FA-CFG-001.a). Die REPLICA IDENTITY FULL bleibt vom
# Aktivieren unberuehrt (LH-FA-CFG-001.a).
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE public.feed_e2e_flow (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_full (id int PRIMARY KEY, name text);
ALTER TABLE public.feed_e2e_full REPLICA IDENTITY FULL;
CREATE TABLE public.feed_e2e_schema (id int PRIMARY KEY, name text, amount text);
INSERT INTO cdc.source (source_id, name) VALUES ('src-e2e', 'E2E-Quelle');
SQL

$COMPOSE up -d pg-change-feed >/dev/null

# Stream-Vertrag am laufenden Container (Muster run-integration-tests.sh):
# die Verdrahtung steht, wenn der Stream seinen Slot angelegt hat und der
# Container laeuft; der Compose-Healthcheck belegt den Prozesszustand
# darueber hinaus.
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
  feed_exit=$(docker inspect --format '{{.State.ExitCode}}' "$FEED_CONTAINER" 2>/dev/null || echo laeuft)
  echo "run-sdk-python-integration-tests: Feed-Container trägt die Verdrahtung nicht — Slot $SLOT: ${slot_ok:-0}, Feed-Lauf: $feed_running (Feed-Ausgang: $feed_exit)" >&2
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
  echo "run-sdk-python-integration-tests: Feed-Container meldet Compose-Health-Status ${health:-fehlt}, wollen healthy" >&2
  exit 1
fi

# Die `integration`-Stufe des Python-SDK bauen — zwingend mit dem benannten
# Bau-Kontext `proto` (COPY --from=proto im Dockerfile, ADR-0090
# Festlegung 2).
docker build --build-context proto=proto -f sdks/python/Dockerfile \
  --target integration -t "$SDK_INTEGRATION_IMAGE" sdks/python

docker rm -f "$SDK_TEST_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$SDK_TEST_CONTAINER" --network "$NETWORK" \
  -e PGCHANGEFEED_GRPC_ADDR="pg-change-feed:9090" \
  -e PGCHANGEFEED_API_TOKEN="$API_TOKEN" \
  -e PGCHANGEFEED_E2E_TABLE="$TEST_TABLE" \
  -e PGCHANGEFEED_E2E_SENTINEL="$TEST_SENTINEL" \
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
  echo "run-sdk-python-integration-tests: Integrationstest wurde nicht innerhalb der Zeitspanne bereit (kein READY): $(docker logs "$SDK_TEST_CONTAINER" 2>&1 || true)" >&2
  exit 1
fi

# Die Zustellung ist Fire-and-Forget ohne Replay (SPEC-020): zwischen
# "Test hat den Stream geoeffnet" und "Server hat den Empfaenger am
# Broadcaster registriert" liegt ein kurzes Fenster, in dem eine committete
# Aenderung fuer diesen Empfaenger verworfen wird. Der Lauf committet
# deshalb eine begrenzte Folge eindeutiger Zeilen, bis der Test genau eine
# davon real empfangen hat (Muster run-integration-tests.shs
# gRPC-Stream-Rundlauf). Die IDs 300ff. liegen in einem eigenen Wertebereich.
received=0
for attempt in 1 2 3 4 5; do
  insert_id=$((300 + attempt))
  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$TEST_TABLE (id, name) VALUES ($insert_id, '$TEST_SENTINEL');
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

# Der Test faehrt nach dem Empfang im selben Lauf noch die Negative-Pruefung
# aus (Stream-Oeffnungsversuch ohne Token); auf deren Ergebnis und auf das
# Prozessende wird separat gewartet, damit die ausgewerteten Zeilen unten
# aus einem abgeschlossenen Lauf stammen.
rejected=0
for _ in $(seq 1 40); do
  if docker logs "$SDK_TEST_CONTAINER" 2>/dev/null | grep -qF "REJECTED code=Unauthenticated"; then
    rejected=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 0.5
done

test_stopped=0
for _ in $(seq 1 20); do
  if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    test_stopped=1
    break
  fi
  sleep 0.5
done
test_exit=$(docker inspect --format '{{.State.ExitCode}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo unbekannt)
test_output=$(docker logs "$SDK_TEST_CONTAINER" 2>&1 || true)

if [ "$received" -ne 1 ]; then
  echo "run-sdk-python-integration-tests: der Test empfing keine der committeten Aenderungen ($TEST_TABLE, $TEST_SENTINEL) ueber den Stream: $test_output" >&2
  exit 1
fi
if ! printf '%s' "$test_output" | grep -qE "RECEIVED .*table=$TEST_TABLE .*operation=INSERT .*new_image=.*$TEST_SENTINEL"; then
  echo "run-sdk-python-integration-tests: die RECEIVED-Zeile traegt nicht die erwartete Aenderung ($TEST_TABLE, INSERT, Sentinel im Row Image): $test_output" >&2
  exit 1
fi
if [ "$rejected" -ne 1 ]; then
  echo "run-sdk-python-integration-tests: Stream-Oeffnungsversuch ohne gueltiges Token wurde nicht mit Unauthenticated abgelehnt: $test_output" >&2
  exit 1
fi
if [ "$test_stopped" -ne 1 ] || [ "$test_exit" != "0" ]; then
  echo "run-sdk-python-integration-tests: der Test endete nicht mit Ausgang 0 (gestoppt: $test_stopped, Ausgang: $test_exit): $test_output" >&2
  exit 1
fi

# Unabhaengiger SQL-Beleg, dass genau die empfangene Aenderung real erfasst
# wurde: die change_id der RECEIVED-Zeile steht ueber den bestehenden
# Lesezugriffsweg in cdc.changes — der Stream-Empfang ist damit keine
# erfundene Ausgabe des Testlaufs.
grpc_change_id=$(printf '%s' "$test_output" | grep -oE 'RECEIVED change_id=[^ ]+' | head -n1 | cut -d= -f2 || true)
if [ -z "$grpc_change_id" ]; then
  echo "run-sdk-python-integration-tests: die RECEIVED-Zeile traegt keine change_id: $test_output" >&2
  exit 1
fi
captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$grpc_change_id' AND table_name = '$TEST_TABLE' AND new_data->>'name' = '$TEST_SENTINEL'")
if [ -z "$captured" ] || [ "$captured" -lt 1 ]; then
  echo "run-sdk-python-integration-tests: die ueber den Stream empfangene Aenderung (change_id=$grpc_change_id, $TEST_SENTINEL) ist nicht real ueber cdc.changes lesbar (count=${captured:-leer})" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-sdk-python-integration-tests: Feed-Container lief nach dem Lauf nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-sdk-python-integration-tests: gRPC-SDK-Realserver-Beleg (ADR-0110 Festlegung 2/Folgepflicht 1) gruen — der Integrationstest (pgchangefeed.grpc_client im $SDK_INTEGRATION_IMAGE-Container) oeffnete real den Server-Stream gegen den laufenden Feed-Container (pg-change-feed:9090) und empfing eine danach committete Aenderung (Tabelle, Operation und Sentinel real am Wire; change_id=$grpc_change_id unabhaengig ueber cdc.changes lesbar); ein Stream-Oeffnungsversuch ohne Token endete mit gRPC-Status Unauthenticated"