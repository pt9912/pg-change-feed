#!/usr/bin/env bash
# transformation-demo.sh — zeigt die Wirkung einer Transformationsregel live
# in der Demo-Umgebung (ADR-0112). Gekapselt hinter `make
# example-transformation-demo` (harness/mk/examples.mk).
#
# Voraussetzung: eine laufende Demo-Umgebung (`make example-demo-up`).
#
# Ablauf, host-seitig (dasselbe `docker exec`-Muster wie examples/bootstrap.sh):
#   1. cdc.set_transformation auf public.orders anwenden: die Spalte
#      "customer" wird im Row Image auf "customer_name" umbenannt.
#   2. Auf den Antrags-Status "applied" pollen.
#   3. Eine neue Zeile einfügen — sie erscheint danach mit dem umbenannten
#      Schlüssel; die von examples/bootstrap.sh bereits eingefügte Zeile
#      bleibt in ihrer ursprünglichen Form (die Regel wirkt nicht rückwirkend,
#      siehe docs/user/benutzerhandbuch.md "Transformationsregel konfigurieren").
#   4. Beide Formen über cdc.changes gegenüberstellen.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

CONTAINER=cdc-examples-postgres
DB=cdc
DBUSER=postgres

# examples/.env ist die einzige Quelle der Wahrheit (ADR-0098 Festlegung 3);
# dieses Skript exportiert sie, statt CDC_SOURCE_ID ein zweites Mal zu tragen.
set -a
# shellcheck disable=SC1091
source examples/.env
set +a

if ! docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -tAc "SELECT 1" >/dev/null 2>&1; then
  echo "transformation-demo: $CONTAINER ist nicht erreichbar — zuerst 'make example-demo-up' ausführen" >&2
  exit 1
fi

echo "transformation-demo: Regel anwenden (customer -> customer_name) …"
request_id=$(docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -v ON_ERROR_STOP=1 -tAc \
  "SELECT cdc.set_transformation('${CDC_SOURCE_ID}', 'public', 'orders', 'demo_rename_customer', '{\"kind\": \"rename_column\", \"column\": \"customer\", \"to\": \"customer_name\"}'::json)")

applied=0
for _ in $(seq 1 30); do
  status=$(docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -tAc \
    "SELECT status FROM cdc.administration_request WHERE administration_request_id = '${request_id}'")
  if [ "$status" = "applied" ]; then
    applied=1
    break
  fi
  if [ "$status" = "failed" ]; then
    echo "transformation-demo: Antrag ${request_id} endete failed:" >&2
    docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -tAc \
      "SELECT error_message FROM cdc.administration_request WHERE administration_request_id = '${request_id}'" >&2
    exit 1
  fi
  sleep 1
done
if [ "$applied" -ne 1 ]; then
  echo "transformation-demo: Antrag ${request_id} wurde nicht innerhalb 30s applied" >&2
  exit 1
fi

echo "transformation-demo: Regel aktiv — neue Zeile einfügen …"
docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "INSERT INTO public.orders (customer, amount) VALUES ('Grace Hopper', 17.00)"

echo "transformation-demo: Rohform (vor der Regel) und transformierte Form (danach) über cdc.changes:"
docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -c \
  "SELECT change_id, new_data FROM cdc.changes WHERE table_name = 'orders' ORDER BY change_id"

echo "transformation-demo: fertig — die erste Zeile trägt weiterhin den Schlüssel \"customer\", die soeben eingefügte trägt \"customer_name\" (siehe auch GET /changes?source=${CDC_SOURCE_ID})"
