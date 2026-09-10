#!/usr/bin/env bash
# run-integration-tests — MVP-Integrationstest gegen die Compose-Umgebung
# (compose.yaml; LH-QA-POR-003). Kette je Lauf: Compose-PostgreSQL frisch
# hochfahren → Schema-Rollout über d-migrate (ADR-0043) → Aktivierung der
# Feed-Tabellen (Publication, Bindungs-Zeilen) → Feed-Container als
# CDC-Runtime starten → Toolchain-Container gegen das Compose-Netz. Der
# Feed-Container streamt dabei selbst (Verdrahtung je ADR-0026); der Test
# schreibt nur Quelländerungen und liest den Store.
#
# Test-Daten bleiben im Container (kein Volume in den Arbeitsbaum);
# Compose-Container und -Netz werden in jedem Ausgang abgeräumt, das
# Modul-Cache-Volume bleibt als Vorbereitung für netzlose `make test`-Läufe
# bestehen.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

TOOLCHAIN_IMAGE=${TOOLCHAIN_IMAGE:-golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125}
GO_MODCACHE_VOLUME=${GO_MODCACHE_VOLUME:-pg-change-feed-gomodcache}
COMPOSE=${COMPOSE:-docker compose -f compose.yaml}
NETWORK=${NETWORK:-cdc-feed-test}
PG_CONTAINER=${PG_CONTAINER:-cdc-test-postgres}
FEED_CONTAINER=${FEED_CONTAINER:-cdc-test-feed}
PG_DB=cdc
PG_USER=postgres
PG_PASSWORD=postgres
# Die Aktivierungs-Werte tragen dieselben Kennungen wie der
# Container-Vertrag in compose.yaml (CDC_TABLES): Feed-Tabellen, Bindungs-
# und Versions-Kennungen der MVP-Läufe.
PUBLICATION=pub_pgc_mvp
SLOT=slot_pgc_mvp

DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"

cleanup() {
  $COMPOSE down -v --remove-orphans >/dev/null 2>&1 || true
}
trap cleanup EXIT

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
  echo "run-integration-tests: Compose-PostgreSQL wurde nicht bereit" >&2
  exit 1
fi

# Schema-Rollout über d-migrate vor jedem E2E-Lauf (ADR-0043): der
# Pflicht-Report landet in tools/schema/plan.yaml, das Rollback-Artefakt in
# tools/schema/down.sql.
make schema-rollout SCHEMA_TARGET="db:$DSN" SCHEMA_ROLLOUT_NETWORK="$NETWORK"

# Aktivierung vor dem Feed-Container-Start (LH-FA-CFG-001): Feed-Tabellen,
# Publication über beide Tabellen, Bindungs-Zeilen in den
# CDC-Referenztabellen. Die Publication ist Start-Vorbedingung des
# Stream-Adapters — ohne sie endet der Feed-Container rot (Klasse
# configuration), deshalb liegt dieser Schritt vor `up`.
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE public.feed_mvp_flow (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_mvp_full (id int PRIMARY KEY, name text);
ALTER TABLE public.feed_mvp_full REPLICA IDENTITY FULL;
CREATE PUBLICATION pub_pgc_mvp FOR TABLE public.feed_mvp_flow, public.feed_mvp_full;
INSERT INTO cdc.source (source_id, name) VALUES ('src-mvp', 'MVP-Quelle');
INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES
  ('tbl-mvp-flow', 'src-mvp', 'public', 'feed_mvp_flow'),
  ('tbl-mvp-full', 'src-mvp', 'public', 'feed_mvp_full');
INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES
  ('sv-mvp-flow', 'tbl-mvp-flow', 1),
  ('sv-mvp-full', 'tbl-mvp-full', 1);
SQL

$COMPOSE up -d pg-change-feed >/dev/null

# Stream-Vertrag am laufenden Container: die Verdrahtung steht, wenn der
# Stream seinen Slot angelegt hat (NewStream) und der Container den
# Stream-Lauf trägt. Die beiden Prüfungen hängen zusammen, aber nicht
# einschließend: der Anlage-Schritt des Slots läuft vor der
# Publication-Prüfung — ein Container, dessen Verdrahtung an der
# Publication scheitert (Klasse configuration), hinterlässt den Slot und
# endet trotzdem. Erst Slot-Bestand und Laufender-Status zusammen tragen
# den Vertrag.
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
  echo "run-integration-tests: Feed-Container trägt die Verdrahtung nicht — Slot $SLOT: ${slot_ok:-0}, Feed-Lauf: $feed_running (Feed-Ausgang: $feed_exit)" >&2
  exit 1
fi

# Modul-Cache befüllen (braucht Netz); der Testlauf selbst trägt den DSN
# über das Compose-Netz und braucht sonst kein Netz.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go mod download

docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_INTEGRATION_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test -v ./test/integration/...