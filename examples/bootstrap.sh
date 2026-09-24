#!/usr/bin/env bash
# bootstrap.sh — bringt die Demo-Umgebung unter examples/ real und ohne
# manuellen Zwischenschritt hoch (LH-QA-OPS-001, ADR-0098). Gekapselt hinter
# `make example-demo-up` (harness/mk/examples.mk).
#
# Ablauf, host-seitig (kein Compose-Hook — Compose kennt kein "nach healthy
# einmalig SQL ausführen"):
#   1. postgres+nats hochfahren, auf reale Verbindungsbereitschaft warten
#      (dasselbe Muster wie tools/schema/apply-rollout.sh).
#   2. Schema-Rollout über d-migrate treiben (tools/schema/schema.yaml,
#      ADR-0043) — derselbe Weg wie Testläufe und Betrieb.
#   3. Beispiel-Quelle (`cdc.source`) registrieren, Beispiel-Tabelle
#      (`public.orders`) leer anlegen — die physische Tabelle muss vor dem
#      Feed-Start existieren, damit `ALTER PUBLICATION ... ADD TABLE` sie
#      erfasst.
#   4. Feed-Container starten — er aktiviert `public.orders` beim eigenen
#      Start automatisch über CDC_TABLES (ADR-0028, EnableTable Use Case,
#      legt Publication + Replication-Slot an); kein `cdc.enable_table`-
#      Aufruf durch dieses Skript nötig.
#   5. Erst NACH aktivem Slot eine Demo-Zeile in `public.orders` einfügen:
#      Logische Replikation erfasst nur Änderungen, die nach der
#      Slot-Erzeugung committen (kein initialer Snapshot-Export dieses
#      CDC-Wegs) — eine vor dem Feed-Start eingefügte Zeile bliebe über
#      `GET /changes` unsichtbar (real geprüft, siehe Bericht).
#
# Idempotent: CREATE TABLE IF NOT EXISTS / ON CONFLICT DO NOTHING tragen
# einen zweiten Aufruf ohne Fehler (ein bereits laufendes `postgres` bzw.
# `pg-change-feed` bleibt über `docker compose up -d` unverändert stehen;
# die Demo-Zeile wird nur eingefügt, wenn `public.orders` noch leer ist).
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

COMPOSE="docker compose -f examples/compose.yaml"
NETWORK=cdc-examples
CONTAINER=cdc-examples-postgres
FEED_CONTAINER=cdc-examples-feed
DB=cdc
DBUSER=postgres

# examples/.env ist die einzige Quelle der Wahrheit (ADR-0098 Festlegung 3);
# dieses Skript exportiert sie, statt die Werte ein zweites Mal zu tragen.
set -a
# shellcheck disable=SC1091
source examples/.env
set +a

echo "bootstrap: postgres+nats hochfahren …"
$COMPOSE up -d postgres nats

ready=0
for _ in $(seq 1 60); do
  if docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -tAc "SELECT 1" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done
if [ "$ready" -ne 1 ]; then
  echo "bootstrap: $CONTAINER nahm innerhalb 60s keine Verbindung an" >&2
  exit 1
fi

# Vorbedingung des Rollouts: leeres Ziel-Schema `cdc`, search_path der
# Rollout-Verbindung (bereits durch tools/schema/compose-init gesetzt;
# idempotent wiederholt, wie tools/schema/apply-rollout.sh es tut).
docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "CREATE SCHEMA IF NOT EXISTS cdc" \
  -c "ALTER ROLE $DBUSER IN DATABASE $DB SET search_path = cdc"

# `make schema-rollout` ist gegen ein bereits migriertes Ziel idempotent
# (zentrale Wache `tools/schema/rolloutguard`, ADR-0043). Die Wache hier ist
# trotzdem ein einfacher Existenz-Check: eine bereits vom ersten Rollout
# angelegte Tabelle bedeutet "schon migriert", der zweite `example-demo-up`-Lauf
# überspringt den Rollout-Schritt. Ein weiterlebendes Demo-Volume mit einem
# älteren Schema-Stand bekommt das neue Schema dadurch nicht; `make
# example-demo-down` (mit `down -v`) baut es neu auf, die Demo trägt keine
# schützenswerten Daten (ADR-0114).
already_migrated=$(docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -tAc "SELECT to_regclass('cdc.source_table') IS NOT NULL")
if [ "$already_migrated" = "t" ]; then
  echo "bootstrap: Schema bereits ausgerollt — überspringe make schema-rollout (Existenz-Check, siehe Kommentar)"
else
  echo "bootstrap: Schema-Rollout über d-migrate …"
  make schema-rollout SCHEMA_TARGET="db:${CDC_ADMIN_DSN}" SCHEMA_ROLLOUT_NETWORK="$NETWORK"
fi

echo "bootstrap: Beispiel-Quelle registrieren, Beispiel-Tabelle leer anlegen …"
docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "CREATE TABLE IF NOT EXISTS public.orders (id serial PRIMARY KEY, customer text NOT NULL, amount numeric(10,2) NOT NULL)" \
  -c "INSERT INTO cdc.source (source_id, name) VALUES ('${CDC_SOURCE_ID}', 'Demo-Quelle') ON CONFLICT (source_id) DO NOTHING"

echo "bootstrap: Feed-Container starten …"
$COMPOSE up -d pg-change-feed

ready=0
for _ in $(seq 1 60); do
  status=$(docker inspect --format='{{.State.Health.Status}}' "$FEED_CONTAINER" 2>/dev/null || echo "")
  if [ "$status" = "healthy" ]; then
    ready=1
    break
  fi
  sleep 1
done
if [ "$ready" -ne 1 ]; then
  echo "bootstrap: $FEED_CONTAINER wurde nicht healthy" >&2
  docker logs "$FEED_CONTAINER" 2>&1 | tail -40 >&2 || true
  exit 1
fi

# Zusätzliche, über die reine Healthcheck-Antwort hinausgehende
# Vorbedingung: der Replication-Slot muss real existieren, bevor eine
# danach eingefügte Zeile sicher erfasst wird (derselbe Slot-Poll wie
# tools/bench-lib.sh bench::start_feed).
slot_ok=0
for _ in $(seq 1 60); do
  if [ "$(docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -tAc "SELECT 1 FROM pg_replication_slots WHERE slot_name = '${CDC_SLOT}'" 2>/dev/null)" = "1" ]; then
    slot_ok=1
    break
  fi
  sleep 1
done
if [ "$slot_ok" -ne 1 ]; then
  echo "bootstrap: Replication-Slot ${CDC_SLOT} wurde nicht angelegt" >&2
  exit 1
fi

echo "bootstrap: Demo-Zeile einfügen (nach Slot-Erzeugung, damit sie real erfasst wird) …"
docker exec "$CONTAINER" psql -U "$DBUSER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "INSERT INTO public.orders (customer, amount) SELECT 'Ada Lovelace', 42.50 WHERE NOT EXISTS (SELECT 1 FROM public.orders)"

echo "bootstrap: Demo-Umgebung bereit — GET /changes?source=${CDC_SOURCE_ID} liefert die Demo-Zeile aus public.orders"
