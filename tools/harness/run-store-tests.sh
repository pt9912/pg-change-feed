#!/usr/bin/env bash
# run-store-tests — Adapter-Tests gegen reale PostgreSQL im Testcontainer
# (ADR-0030). Beide Images sind per Digest gepinnt (Modul 14); der Pin der
# Datenbank stammt aus `docker manifest inspect postgres:18-alpine` (amd64).
# Die DB-Daten bleiben im Container (kein Volume in den Arbeitsbaum); der
# Testcontainer und das Docker-Netz werden in jedem Ausgang abgeräumt, das
# Modul-Cache-Volume bleibt als Vorbereitung für netzlose `make test`-Läufe
# bestehen.
#
# Der Lauf erzeugt ein `-coverprofile` über den DB-Adapter-Gegenstand
# (ADR-0071 Punkt 3) und legt es als `store.coverprofile` in DB_COVERAGE_DIR
# ab; die gemergte, subjekt-qualifizierte Zahl entsteht mit dem Replication-Lauf
# (tools/harness/db-coverage.sh). Kein Gate.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

DB_COVERAGE_DIR=${DB_COVERAGE_DIR:-${TMPDIR:-/tmp}/pg-change-feed-db-coverage}
mkdir -p "$DB_COVERAGE_DIR"
COVER_PKGS="$(bash tools/harness/db-coverage.sh --coverpkg)"

TOOLCHAIN_IMAGE=${TOOLCHAIN_IMAGE:-golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125}
PG_TEST_IMAGE=${PG_TEST_IMAGE:-postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8}
GO_MODCACHE_VOLUME=${GO_MODCACHE_VOLUME:-pg-change-feed-gomodcache}
NETWORK=cdc-store-test
PG_CONTAINER=cdc-store-test-pg
PG_DB=cdc_test
PG_USER=cdc
PG_PASSWORD=cdc

docker network inspect "$NETWORK" >/dev/null 2>&1 || docker network create "$NETWORK"

# `-v` entfernt die anonymen Volumes des Containers (das Postgres-Image
# deklariert ein VOLUME).
cleanup() {
  docker rm -fv "$PG_CONTAINER" >/dev/null 2>&1 || true
  docker network rm "$NETWORK" >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker rm -fv "$PG_CONTAINER" >/dev/null 2>&1 || true
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

# Der Schema-Stand vor dem Testlauf: die eine Schema-Anwendung der
# Test-Läufe (tools/schema/apply-rollout.sh), die derselbe d-migrate-Rollout
# trägt wie der Betrieb (ADR-0043) und die auch der Replication-Tier-Lauf
# aufruft.
bash tools/schema/apply-rollout.sh "$PG_CONTAINER" "$PG_DB" "$PG_USER" "$NETWORK" "$DSN"

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
# es berührt; der zweite Aufruf deckt die übrigen Pakete ab, `postgresstorage`
# ist aus diesem Sammelaufruf ausgenommen und läuft weiter unten als eigener
# Messaufruf.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_STORE_TEST_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test ./internal/bootstrap/...

# `postgresstorage` läuft in einem eigenen, letzten Aufruf (die Messung
# unten): seine Tests räumen das `cdc`-Schema über `DROP SCHEMA cdc CASCADE`
# ab und setzen nur die `ApplySchema`-Tabellen neu auf — liefen sie mit den
# übrigen Paketen (oder vor dem vorgezogenen bootstrap-Aufruf), fehlte diesen
# das per d-migrate ausgerollte Schema. Der Aufruf trägt deshalb zugleich das
# `-coverprofile` des Store-Teils der DB-Adapter-Coverage.
mapfile -t OTHER_PACKAGES < <(docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go list ./... \
    | grep -v '/internal/bootstrap$' \
    | grep -v '/internal/adapters/driven/postgresstorage$')

docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_STORE_TEST_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test "${OTHER_PACKAGES[@]}"

# DB-Adapter-Coverage (ADR-0071 Punkt 3): Messlauf über den Store-Teil des
# Gegenstands (`postgresstorage` ohne das Unterpaket `mapper`), mit `-coverpkg`
# über die ganze Gegenstandsliste. Die gemergte, subjekt-qualifizierte Zahl
# entsteht mit dem Replication-Lauf (tools/harness/db-coverage.sh).
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -v "$DB_COVERAGE_DIR":/cov \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_STORE_TEST_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test \
    -coverpkg="$COVER_PKGS" \
    -coverprofile=/cov/store.coverprofile \
    -covermode=atomic \
    ./internal/adapters/driven/postgresstorage

bash tools/harness/db-coverage.sh
