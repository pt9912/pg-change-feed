#!/usr/bin/env bash
# run-store-tests — Adapter-Tests gegen reale PostgreSQL im Testcontainer
# (ADR-0030). Beide Images sind per Digest gepinnt (Modul 14); der Pin der
# Datenbank stammt aus `docker manifest inspect postgres:18-alpine` (amd64).
# Die DB-Daten bleiben im Container (kein Volume in den Arbeitsbaum); der
# Testcontainer und das Docker-Netz werden in jedem Ausgang abgeräumt, das
# Modul-Cache-Volume bleibt als Vorbereitung für netzlose `make test`-Läufe
# bestehen.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

TOOLCHAIN_IMAGE=${TOOLCHAIN_IMAGE:-golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125}
PG_TEST_IMAGE=${PG_TEST_IMAGE:-postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8}
GO_MODCACHE_VOLUME=${GO_MODCACHE_VOLUME:-pg-change-feed-gomodcache}
NETWORK=cdc-store-test
PG_CONTAINER=cdc-store-test-pg
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
  "$PG_TEST_IMAGE" >/dev/null

# Die Bereitschaft verlangt eine echte Abfrage: pg_isready meldet bereit,
# sobald der Server antwortet — auch der temporäre Server der Initdb-Phase
# antwortet (mit „shutting down“). Die Abfrage liest ihre Bereitschaft vom
# Datenbank-Server, nicht vom Initdb-Lauf.
ready=0
for _ in $(seq 1 60); do
  if docker exec "$PG_CONTAINER" pg_isready -U "$PG_USER" -d "$PG_DB" >/dev/null 2>&1 \
    && docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc "SELECT 1" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done
if [ "$ready" -ne 1 ]; then
  echo "run-store-tests: Testcontainer wurde nicht bereit" >&2
  exit 1
fi

DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"

# CDC-Ziel-Schema und search_path der Rollout-Verbindung: dieselbe
# Umgebungs-Vorbedingung, die der Compose-Init der Integration-Umgebung
# trägt (compose-init/01-cdc-schema.sql) — die Tabellen-DDL entsteht über
# den d-migrate-Rollout (ADR-0043), nicht hier. Der Rollout trägt die
# Consumer-State-Tabellen vor dem Testlauf (State-Schema vor der Nutzung):
# die handgeschriebene DDL des Store-Adapters (ApplySchema) rollt sie
# nicht aus.
docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 \
  -c "CREATE SCHEMA IF NOT EXISTS cdc" \
  -c "ALTER ROLE $PG_USER IN DATABASE $PG_DB SET search_path = cdc"

# Schema-Rollout über d-migrate vor dem Testlauf (ADR-0043): der
# Pflicht-Report landet in tools/schema/plan.yaml, das Rollback-Artefakt
# in tools/schema/down.sql.
make schema-rollout SCHEMA_TARGET="db:$DSN" SCHEMA_ROLLOUT_NETWORK="$NETWORK"

# Modul-Cache befüllen (braucht Netz); der Testlauf selbst trägt den
# DSN über das Docker-Netz und braucht sonst kein Netz.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go mod download

# `internal/bootstrap` läuft in einem eigenen, vorgezogenen Aufruf: seine
# Zugriffsweg-Tests (`register-consumer`, `acknowledge-consumer`,
# `LH-FA-CON-001.a`/`LH-FA-CON-004.a`) schreiben gegen dieselben Zeilen
# (`cdc.consumer`/`cdc.consumer_position`), die mehrere Tests im Paket
# `postgresstorage` tabellenweit abräumen bzw. per `DROP SCHEMA cdc
# CASCADE` neu aufsetzen (`store_test.go`, `tableactivation_test.go`) —
# ein paralleler Lauf (Go-Default über Pakete hinweg) kann eine gerade
# registrierte Kennung wegreißen oder das erweiterte, per d-migrate
# ausgerollte Schema unter dem laufenden Test entfernen. Der vorgezogene
# Lauf steht auf dem frisch ausgerollten Schema, bevor ein anderes Paket
# es berührt; der zweite Aufruf deckt alle übrigen Pakete wie zuvor ab.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_STORE_TEST_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test ./internal/bootstrap/...

mapfile -t OTHER_PACKAGES < <(docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go list ./... | grep -v '/internal/bootstrap$')

docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_STORE_TEST_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test "${OTHER_PACKAGES[@]}"
