#!/usr/bin/env bash
# apply-rollout — die eine Schema-Anwendung der Test-Läufe: sie rollt das
# CDC-Schema mit d-migrate aus derselben Quelle aus, die der Betrieb nutzt
# (`tools/schema/schema.yaml`, `make schema-rollout`, ADR-0043), und trägt
# die Umgebungs-Vorbedingung des Rollouts (leeres Ziel-Schema `cdc`, search_path
# der Rollout-Verbindung — derselbe Schritt wie compose-init/01-cdc-schema.sql
# für die Integration-Umgebung).
#
# Aufrufer sind tools/harness/run-store-tests.sh und
# tools/harness/run-replication-tests.sh: beide Test-Läufe stehen damit auf
# demselben Schema-Stand wie der Betrieb, statt ihn je Lauf selbst
# zusammenzustellen. Pflicht-Report (tools/schema/plan.yaml) und
# Rollback-Artefakt (tools/schema/down.sql) des make-Targets stellt
# tools/schema/rollout-restore.sh nach dem Rollout wieder her.
#
# Aufruf: apply-rollout.sh <container> <db> <user> <network> <dsn>
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

CONTAINER=$1
DB=$2
USER=$3
NETWORK=$4
DSN=$5

# Der Server des Ziel-Containers trägt seine Verbindung erst nach dem Ende der
# Initdb-Phase: `pg_isready` antwortet schon am temporären Server dieses
# Abschnitts. Die Anwendung wartet deshalb auf eine echte Abfrage, statt sich
# auf den Bereitschafts-Melder des Aufrufers zu verlassen.
ready=0
for _ in $(seq 1 60); do
  if docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -tAc "SELECT 1" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done
if [ "$ready" -ne 1 ]; then
  echo "apply-rollout: $CONTAINER nahm innerhalb 60s keine Verbindung an" >&2
  exit 1
fi

docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "CREATE SCHEMA IF NOT EXISTS cdc" \
  -c "ALTER ROLE $USER IN DATABASE $DB SET search_path = cdc"

bash tools/schema/rollout-restore.sh make schema-rollout SCHEMA_TARGET="db:$DSN" SCHEMA_ROLLOUT_NETWORK="$NETWORK"
