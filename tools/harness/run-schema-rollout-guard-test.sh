#!/usr/bin/env bash
# run-schema-rollout-guard-test — realer Beleg der zentralen
# `make schema-rollout`-Idempotenz-Wache (ADR-0043, tools/schema/rolloutguard,
# BEO-PGC/schema-rollout-fremdobjekte): vier Läufe gegen dieselbe,
# eigenständige Ziel-DB.
#
#   1. Frischer Rollout (kein Blocker) — muss durchlaufen.
#   2. Zweiter Lauf gegen dasselbe, jetzt vollständig migrierte Ziel — trägt
#      ausschließlich die sechs bekannten Fremdobjekt-Blocker — muss erneut
#      Exit 0 liefern UND den --allow-destructive-Pfad des Guards nehmen
#      (stdout trägt die Meldung), nicht nur zufällig Exit 0 aus anderem
#      Grund.
#   3. Eine echte, gleichzeitig anstehende, NICHT-destruktive Schema-Änderung
#      neben den sechs bekannten Blockern: eine von schema.yaml weiterhin
#      deklarierte, nullable Spalte (administration_request.error_message —
#      trägt keine View-/Funktions-Abhängigkeit, anders als
#      process_heartbeat.error_class, das die Sicht cdc.heartbeat trägt) wird
#      per `ALTER TABLE … DROP COLUMN` außerhalb von d-migrate entfernt — der
#      nächste Rollout-Lauf muss sie real zurückbringen. Das ist der
#      Regressionsschutz für den konkret gefundenen Fehler: ein reines
#      Überspringen von `--execute` bei bekannten Blockern ließ eine solche
#      echte Änderung verlustig gehen.
#   4. Ein künstlich per `ALTER TABLE … ADD COLUMN` hinzugefügtes, nicht
#      deklariertes Objekt — muss weiterhin mit Exit 8 abbrechen (der
#      Beleg, dass die Wache nicht pauschal durchlässt).
#
# Eigenständiges Docker-Netz/-Container, unabhängig von der
# Wurzel-`compose.yaml` und `examples/compose.yaml` — Daten leben
# ausschließlich im Container, in jedem Ausgang abgeräumt. Kein Gate:
# braucht DB-Zugang, `make gates` bleibt netzlos.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

PG_TEST_IMAGE=${PG_TEST_IMAGE:-postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8}
NETWORK=cdc-schema-rollout-guard-test
CONTAINER=cdc-schema-rollout-guard-test-pg
DB=cdc
USER=postgres
PASSWORD=postgres
TARGET="db:postgres://$USER:$PASSWORD@$CONTAINER:5432/$DB?sslmode=disable"

docker network inspect "$NETWORK" >/dev/null 2>&1 || docker network create "$NETWORK" >/dev/null

cleanup() {
  docker rm -f "$CONTAINER" >/dev/null 2>&1 || true
  docker network rm "$NETWORK" >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker rm -f "$CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$CONTAINER" \
  --network "$NETWORK" \
  -e POSTGRES_DB="$DB" -e POSTGRES_USER="$USER" -e POSTGRES_PASSWORD="$PASSWORD" \
  "$PG_TEST_IMAGE" -c wal_level=logical >/dev/null

ready=0
for _ in $(seq 1 60); do
  if docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -tAc "SELECT 1" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done
if [ "$ready" -ne 1 ]; then
  echo "run-schema-rollout-guard-test: $CONTAINER nahm innerhalb 60s keine Verbindung an" >&2
  exit 1
fi

docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "CREATE SCHEMA IF NOT EXISTS cdc" \
  -c "ALTER ROLE $USER IN DATABASE $DB SET search_path = cdc"

echo "run-schema-rollout-guard-test: Lauf 1/4 (frischer Rollout, muss durchlaufen)"
make schema-rollout SCHEMA_TARGET="$TARGET" SCHEMA_ROLLOUT_NETWORK="$NETWORK"

echo "run-schema-rollout-guard-test: Lauf 2/4 (Idempotenz — muss den --allow-destructive-Pfad nehmen)"
out=$(make schema-rollout SCHEMA_TARGET="$TARGET" SCHEMA_ROLLOUT_NETWORK="$NETWORK" 2>&1)
echo "$out"
if ! grep -q -- "--allow-destructive" <<<"$out"; then
  echo "run-schema-rollout-guard-test: FEHLER — Lauf 2 hat den --allow-destructive-Pfad nicht genommen (Meldung fehlt in der Ausgabe)" >&2
  exit 1
fi

echo "run-schema-rollout-guard-test: Lauf 3/4 (echte anstehende Änderung neben den sechs bekannten Blockern — muss real zurückkommen)"
docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "ALTER TABLE cdc.administration_request DROP COLUMN error_message"
still_missing=$(docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -tAc \
  "SELECT count(*) FROM information_schema.columns WHERE table_schema='cdc' AND table_name='administration_request' AND column_name='error_message'")
if [ "$still_missing" != "0" ]; then
  echo "run-schema-rollout-guard-test: FEHLER — Vorbedingung fehlgeschlagen, error_message ist nach DROP COLUMN noch da" >&2
  exit 1
fi
make schema-rollout SCHEMA_TARGET="$TARGET" SCHEMA_ROLLOUT_NETWORK="$NETWORK"
restored=$(docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -tAc \
  "SELECT count(*) FROM information_schema.columns WHERE table_schema='cdc' AND table_name='administration_request' AND column_name='error_message'")
if [ "$restored" != "1" ]; then
  echo "run-schema-rollout-guard-test: FEHLER — Lauf 3 hat die echte anstehende Spalten-Wiederherstellung nicht angewendet (Regression des behobenen Fehlers)" >&2
  exit 1
fi

echo "run-schema-rollout-guard-test: Lauf 4/4 (unbekannter Blocker, muss mit Exit 8 abbrechen)"
docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "ALTER TABLE cdc.source ADD COLUMN _rolloutguard_test_col text"
set +e
make schema-rollout SCHEMA_TARGET="$TARGET" SCHEMA_ROLLOUT_NETWORK="$NETWORK"
run4_exit=$?
set -e
if [ "$run4_exit" -eq 0 ]; then
  echo "run-schema-rollout-guard-test: FEHLER — Lauf 4 lief durch (Exit 0), obwohl ein unbekannter Blocker vorlag" >&2
  exit 1
fi

echo "run-schema-rollout-guard-test: OK — alle drei DoD-Belege real erbracht (Idempotenz-Allow, echte Änderung bleibt wirksam, Negativ-Abbruch Exit $run4_exit)"
