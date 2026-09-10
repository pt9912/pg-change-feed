#!/usr/bin/env bash
# run-integration-tests — MVP-Integrationstest gegen die Compose-Umgebung
# (compose.yaml; LH-QA-POR-003). Kette je Lauf: Compose-Umgebung frisch
# hochfahren → Schema-Rollout über d-migrate (ADR-0043) → Toolchain-Container
# gegen das Compose-Netz. Test-Daten bleiben im Container (kein Volume in
# den Arbeitsbaum); Compose-Container und -Netz werden in jedem Ausgang
# abgeräumt, das Modul-Cache-Volume bleibt als Vorbereitung für netzlose
# `make test`-Läufe bestehen.
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

DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"

cleanup() {
  $COMPOSE down -v --remove-orphans >/dev/null 2>&1 || true
}
trap cleanup EXIT

$COMPOSE down -v --remove-orphans >/dev/null 2>&1 || true
$COMPOSE up -d postgres pg-change-feed >/dev/null

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

# Feed-Container-Smoke: der Container trägt das geladene Image mit dem
# --version-Vertrag; sein Ausgang ist der Compose-Anteil des Lauf-Belegs.
feed_exit=$(docker inspect --format '{{.State.ExitCode}}' "$FEED_CONTAINER" 2>/dev/null || echo fehlt)
if [ "$feed_exit" != "0" ]; then
  echo "run-integration-tests: Feed-Container endete mit Exit $feed_exit (Erwartung: --version, Exit 0)" >&2
  exit 1
fi

# Schema-Rollout über d-migrate vor jedem E2E-Lauf (ADR-0043): der
# Pflicht-Report landet in tools/schema/plan.yaml, das Rollback-Artefakt in
# tools/schema/down.sql.
make schema-rollout SCHEMA_TARGET="db:$DSN" SCHEMA_ROLLOUT_NETWORK="$NETWORK"

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
  "$TOOLCHAIN_IMAGE" go test ./test/integration/...
