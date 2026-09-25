#!/usr/bin/env bash
# run-replication-tests — Replication-Stream-Tests gegen reale PostgreSQL
# mit Publication und Logical Replication Slot im Testcontainer (ADR-0030).
# wal_sender_timeout=2000 zieht die Keepalive-Antwort-Pflicht auf etwa eine
# Sekunde — der Keepalive-Beleg (LH-QA-REL-001.a) braucht den Antwort-Zug
# in Test-Zeitspanne; der Default (60s) läge jenseits der Testgrenze. Eine
# zweite Instanz mit dem Standardwert und ohne gleichzeitigen fremden Schreiber
# (CDC_REPLICATION_TEST_STANDARD_DSN) trägt die Rückstand-Messungen der
# Leerlauf-Bestätigung (ADR-0120) und, nach dem Tier-Lauf, den Beleg der
# WAL-Rückstand-Schwellen (CDC_WALRETENTION_TEST_DSN).
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
#   tier    — der Tier-weite `go test ./...` und der Slot-Reserve-Lauf des
#             Snapshot-Adapters gegen einen eigenen PostgreSQL.
# Ohne Argument (make test-replication) laufen beide Phasen nacheinander.
# Die Trennung haelt beide Verdikte lesbar: der Tier-Lauf fuehrt `go test ./...`
# ueber den ganzen Baum — im selben Schritt verschluckte sein Exit das Verdikt
# der Messung.
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
PG_EXCLUSIVE_CONTAINER=cdc-repl-test-pg-slot1
PG_STANDARD_CONTAINER=cdc-repl-test-pg-default
PG_DB=cdc_test
PG_USER=cdc
PG_PASSWORD=cdc

docker network inspect "$NETWORK" >/dev/null 2>&1 || docker network create "$NETWORK"

# `-v` entfernt die anonymen Volumes des Containers (das Postgres-Image
# deklariert ein VOLUME).
cleanup() {
  docker rm -fv "$PG_CONTAINER" "$PG_EXCLUSIVE_CONTAINER" "$PG_STANDARD_CONTAINER" >/dev/null 2>&1 || true
  docker network rm "$NETWORK" >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker rm -fv "$PG_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$PG_CONTAINER" \
  --network "$NETWORK" \
  -e POSTGRES_DB="$PG_DB" -e POSTGRES_USER="$PG_USER" -e POSTGRES_PASSWORD="$PG_PASSWORD" \
  "$PG_TEST_IMAGE" -c wal_level=logical -c max_wal_senders=10 -c max_replication_slots=10 -c wal_sender_timeout=2000 >/dev/null

wait_ready() {
  local container=$1
  for _ in $(seq 1 60); do
    if docker exec "$container" pg_isready -U "$PG_USER" -d "$PG_DB" >/dev/null 2>&1; then
      return 0
    fi
    sleep 1
  done
  echo "run-replication-tests: Testcontainer $container wurde nicht bereit" >&2
  return 1
}
wait_ready "$PG_CONTAINER"

DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"

# Instanz mit dem Standardwert von wal_sender_timeout (60s), ohne Schema und
# ohne gleichzeitigen fremden Schreiber: die Tests der Leerlauf-Bestätigung
# (Rückstand unter einem Zehntel der Last) und der WAL-Rückstand-Schwellen
# messen das WAL der ganzen Instanz und bringen ihre Tabellen selbst mit.
docker rm -fv "$PG_STANDARD_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$PG_STANDARD_CONTAINER" \
  --network "$NETWORK" \
  -e POSTGRES_DB="$PG_DB" -e POSTGRES_USER="$PG_USER" -e POSTGRES_PASSWORD="$PG_PASSWORD" \
  "$PG_TEST_IMAGE" -c wal_level=logical -c max_wal_senders=10 -c max_replication_slots=10 >/dev/null
wait_ready "$PG_STANDARD_CONTAINER"
STANDARD_DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_STANDARD_CONTAINER:5432/$PG_DB?sslmode=disable"

# Modul-Cache befüllen (braucht Netz); der Testlauf selbst trägt den
# DSN über das Docker-Netz und braucht sonst kein Netz.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go mod download

# Phase `measure` — DB-Adapter-Coverage (ADR-0071 Punkt 3): eigener Messlauf
# über den Replication-Teil des Gegenstands (`postgresack`, `postgressnapshot`,
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
    -e CDC_REPLICATION_TEST_STANDARD_DSN="$STANDARD_DSN" \
    "$TOOLCHAIN_IMAGE" go test \
      -coverpkg="$COVER_PKGS" \
      -coverprofile=/cov/replication.coverprofile \
      -covermode=atomic \
      ./internal/adapters/driven/postgresack \
      ./internal/adapters/driven/postgressnapshot \
      ./internal/adapters/driving/replication/receive

  bash tools/harness/db-coverage.sh
fi

# Phase `tier` — drei Läufe nacheinander: der Tier-weite `go test ./...` gegen
# den gemeinsamen Container, der Slot-Reserve-Lauf des Snapshot-Adapters
# gegen einen eigenen Container und der Schwellen-Beleg auf der
# Standard-Instanz (unten); jeder Exit ist ein Verdikt dieser Phase, und die
# beiden letzten Läufe brauchen zusätzlich das `--- PASS` ihres Tests.
# Der Schema-Stand dieses Laufs kommt aus derselben
# Schema-Anwendung wie der Betrieb (tools/schema/apply-rollout.sh,
# `make test-store`): die `internal/bootstrap`-Fixtures starten
# `bootstrap.Run` real, und der liest `cdc.table_schema`
# (`SchemaStorePort.CurrentVersion`), `cdc.administration_request`
# (`ColumnExclusionPort.ExcludedColumns`, ADR-0050) und
# `cdc.process_heartbeat` (HeartbeatPort). Die `measure`-Phase oben läuft
# ohne diesen Schritt: ihr Gegenstand (`postgresack`, `postgressnapshot`,
# `replication/receive`) bringt seine Tabellen selbst mit.
if [[ "$MODE" == "tier" || "$MODE" == "both" ]]; then
  bash tools/schema/apply-rollout.sh "$PG_CONTAINER" "$PG_DB" "$PG_USER" "$NETWORK" "$DSN"

  docker run --rm --network "$NETWORK" \
    -v "$(pwd)":/src:ro \
    -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
    -w /src \
    -e GOCACHE=/tmp/gocache \
    -e CDC_REPLICATION_TEST_DSN="$DSN" \
    -e CDC_REPLICATION_TEST_STANDARD_DSN="$STANDARD_DSN" \
    "$TOOLCHAIN_IMAGE" go test ./...

  # Slot-Reserve des Snapshot-Adapters: `TestSlotReserveExhaustedIsConfiguration`
  # braucht einen PostgreSQL mit `max_replication_slots=1` und läuft deshalb
  # gegen einen eigenen Container, nicht gegen den gemeinsamen der Pakete
  # oben. Ein Lauf, in dem der Test nicht als PASS erscheint (Name
  # verschoben, Variable ungelesen), ist rot.
  docker rm -fv "$PG_EXCLUSIVE_CONTAINER" >/dev/null 2>&1 || true
  docker run -d --name "$PG_EXCLUSIVE_CONTAINER" \
    --network "$NETWORK" \
    -e POSTGRES_DB="$PG_DB" -e POSTGRES_USER="$PG_USER" -e POSTGRES_PASSWORD="$PG_PASSWORD" \
    "$PG_TEST_IMAGE" -c wal_level=logical -c max_wal_senders=10 -c max_replication_slots=1 >/dev/null
  wait_ready "$PG_EXCLUSIVE_CONTAINER"
  reserve_out=$(docker run --rm --network "$NETWORK" \
    -v "$(pwd)":/src:ro \
    -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
    -w /src \
    -e GOCACHE=/tmp/gocache \
    -e CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_EXCLUSIVE_CONTAINER:5432/$PG_DB?sslmode=disable" \
    "$TOOLCHAIN_IMAGE" go test -count=1 -v -run '^TestSlotReserveExhaustedIsConfiguration$' \
      ./internal/adapters/driven/postgressnapshot 2>&1) || { printf '%s\n' "$reserve_out"; exit 1; }
  printf '%s\n' "$reserve_out"
  if ! grep -q -- '--- PASS: TestSlotReserveExhaustedIsConfiguration' <<<"$reserve_out"; then
    echo "run-replication-tests: TestSlotReserveExhaustedIsConfiguration ist nicht als PASS gelaufen" >&2
    exit 1
  fi

  # WAL-Rückstand-Schwellen (ADR-0049): `TestWALRetentionThresholdEndToEnd`
  # misst das WAL der ganzen Instanz gegen Schwellen im KiB-Bereich und läuft
  # deshalb allein auf der Standard-Instanz, nachdem der Tier-Lauf ihre
  # übrigen Tests beendet hat. Ein Lauf ohne PASS des Tests ist rot.
  threshold_out=$(docker run --rm --network "$NETWORK" \
    -v "$(pwd)":/src:ro \
    -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
    -w /src \
    -e GOCACHE=/tmp/gocache \
    -e CDC_WALRETENTION_TEST_DSN="$STANDARD_DSN" \
    "$TOOLCHAIN_IMAGE" go test -count=1 -v -run '^TestWALRetentionThresholdEndToEnd$' \
      ./internal/bootstrap 2>&1) || { printf '%s\n' "$threshold_out"; exit 1; }
  printf '%s\n' "$threshold_out"
  if ! grep -q -- '--- PASS: TestWALRetentionThresholdEndToEnd' <<<"$threshold_out"; then
    echo "run-replication-tests: TestWALRetentionThresholdEndToEnd ist nicht als PASS gelaufen" >&2
    exit 1
  fi
fi
