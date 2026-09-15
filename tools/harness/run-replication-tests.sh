#!/usr/bin/env bash
# run-replication-tests — Replication-Stream-Tests gegen reale PostgreSQL
# mit Publication und Logical Replication Slot im Testcontainer (ADR-0030).
# wal_sender_timeout=2000 zieht die Keepalive-Antwort-Pflicht auf etwa eine
# Sekunde — der Keepalive-Beleg (LH-QA-REL-001.a) braucht den Antwort-Zug
# in Test-Zeitspanne; der Default (60s) läge jenseits der Testgrenze.
# Beide Images sind per Digest gepinnt (Modul 14); der Pin der Datenbank
# stammt aus `docker manifest inspect postgres:18-alpine` (amd64). Die
# Instanz startet mit wal_level=logical — der Logical-Replication-Slot
# braucht ihn. DB-Daten bleiben im Container (kein Volume in den
# Arbeitsbaum); der Testcontainer und das Docker-Netz werden in jedem
# Ausgang abgeräumt, das Modul-Cache-Volume bleibt als Vorbereitung für
# netzlose `make test`-Läufe bestehen.
#
# Das Skript traegt zwei Phasen, per Argument waehlbar:
#   measure — erzeugt das `-coverprofile` des Replication-Teils des
#             DB-Adapter-Gegenstands (ADR-0071 Punkt 3) als
#             `replication.coverprofile` in DB_COVERAGE_DIR und ruft
#             tools/harness/db-coverage.sh; dessen Exit ist das Verdikt der
#             DB-Adapter-Coverage (kein Gate).
#   tier    — der Tier-weite `go test ./...`.
# Ohne Argument (make test-replication) laufen beide Phasen nacheinander.
# Die Trennung ist noetig, weil der Tier-Lauf einen vorbestehenden roten Beleg
# traegt (internal/bootstrap, TestWALRetentionThresholdEndToEnd) — im selben
# Schritt verschluckte sein Exit das Verdikt der Messung.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

MODE=${1:-both}
case "$MODE" in
  measure|tier|both) ;;
  *) echo "run-replication-tests: unbekannter Modus '$MODE' (erlaubt: measure, tier)" >&2; exit 2 ;;
esac

DB_COVERAGE_DIR=${DB_COVERAGE_DIR:-${TMPDIR:-/tmp}/pg-change-feed-db-coverage}
mkdir -p "$DB_COVERAGE_DIR"
COVER_PKGS="$(bash tools/harness/db-coverage.sh --coverpkg)"

TOOLCHAIN_IMAGE=${TOOLCHAIN_IMAGE:-golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125}
PG_TEST_IMAGE=${PG_TEST_IMAGE:-postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8}
GO_MODCACHE_VOLUME=${GO_MODCACHE_VOLUME:-pg-change-feed-gomodcache}
NETWORK=cdc-repl-test
PG_CONTAINER=cdc-repl-test-pg
PG_DB=cdc_test
PG_USER=cdc
PG_PASSWORD=cdc

docker network inspect "$NETWORK" >/dev/null 2>&1 || docker network create "$NETWORK"

cleanup() {
  docker rm -f "$PG_CONTAINER" >/dev/null 2>&1 || true
  docker network rm "$NETWORK" >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker rm -f "$PG_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$PG_CONTAINER" \
  --network "$NETWORK" \
  -e POSTGRES_DB="$PG_DB" -e POSTGRES_USER="$PG_USER" -e POSTGRES_PASSWORD="$PG_PASSWORD" \
  "$PG_TEST_IMAGE" -c wal_level=logical -c max_wal_senders=10 -c max_replication_slots=10 -c wal_sender_timeout=2000 >/dev/null

ready=0
for _ in $(seq 1 60); do
  if docker exec "$PG_CONTAINER" pg_isready -U "$PG_USER" -d "$PG_DB" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done
if [ "$ready" -ne 1 ]; then
  echo "run-replication-tests: Testcontainer wurde nicht bereit" >&2
  exit 1
fi

DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"

# Modul-Cache befüllen (braucht Netz); der Testlauf selbst trägt den
# DSN über das Docker-Netz und braucht sonst kein Netz.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go mod download

# Phase `measure` — DB-Adapter-Coverage (ADR-0071 Punkt 3): eigener Messlauf
# über den Replication-Teil des Gegenstands (`postgresack`,
# `replication/receive`), mit `-coverpkg` über die ganze Gegenstandsliste. Der
# Exit von db-coverage.sh ist das Verdikt dieser Phase.
if [[ "$MODE" == "measure" || "$MODE" == "both" ]]; then
  docker run --rm --network "$NETWORK" \
    -v "$(pwd)":/src:ro \
    -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
    -v "$DB_COVERAGE_DIR":/cov \
    -w /src \
    -e GOCACHE=/tmp/gocache \
    -e CDC_REPLICATION_TEST_DSN="$DSN" \
    "$TOOLCHAIN_IMAGE" go test \
      -coverpkg="$COVER_PKGS" \
      -coverprofile=/cov/replication.coverprofile \
      -covermode=atomic \
      ./internal/adapters/driven/postgresack \
      ./internal/adapters/driving/replication/receive

  bash tools/harness/db-coverage.sh
fi

# Phase `tier` — der Tier-weite `go test ./...`; der Exit dieses Aufrufs ist
# das Verdikt dieser Phase.
if [[ "$MODE" == "tier" || "$MODE" == "both" ]]; then
  docker run --rm --network "$NETWORK" \
    -v "$(pwd)":/src:ro \
    -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
    -w /src \
    -e GOCACHE=/tmp/gocache \
    -e CDC_REPLICATION_TEST_DSN="$DSN" \
    "$TOOLCHAIN_IMAGE" go test ./...
fi
