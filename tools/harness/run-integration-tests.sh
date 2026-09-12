#!/usr/bin/env bash
# run-integration-tests — MVP-Integrationstest gegen die Compose-Umgebung
# (compose.yaml; LH-QA-POR-003). Kette je Lauf: Compose-PostgreSQL frisch
# hochfahren → Schema-Rollout über d-migrate (ADR-0043) → Vorbedingungen
# der Aktivierung (Quell-Tabellen, Quelle-Zeile) → Feed-Container als
# CDC-Runtime starten — seine Verdrahtung aktiviert die Tabellen über den
# EnableTable Use Case (ADR-0028), nicht über Seed-SQL → Toolchain-
# Container gegen das Compose-Netz. Der Feed-Container streamt dabei
# selbst (Verdrahtung je ADR-0026); der Test schreibt nur Quelländerungen
# und liest den Store.
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
# Der Slot-Name trägt denselben Wert wie der Container-Vertrag in
# compose.yaml (CDC_SLOT); der Runner liest ihn im Start- und im
# End-Wächter.
SLOT=slot_pgc_mvp

DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"

cleanup() {
  docker unpause "$FEED_CONTAINER" >/dev/null 2>&1 || true
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

# Vorbedingung der Aktivierung (LH-FA-CFG-001.a): die physischen
# Quell-Tabellen und die Zeile der Quelle in cdc.source — die
# Metadaten-Registrierung ist Vorbedingung der Aktivierung (SPEC-001,
# Fremdschlüssel). Die Aktivierung selbst trägt die Verdrahtung des
# Feed-Containers: die Bindungen aus CDC_TABLES laufen als
# EnableTable-Aufrufe (ADR-0028) — Publication, Bindungs- und
# Schema-Version-Zeilen entstehen dort, nicht hier. Die REPLICA IDENTITY
# trägt der Runner für die volle Alt-Bild-Prüfung (LH-FA-CAP-008); sie
# bleibt vom Aktivieren unberührt (LH-FA-CFG-001.a).
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE public.feed_mvp_flow (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_mvp_full (id int PRIMARY KEY, name text);
ALTER TABLE public.feed_mvp_full REPLICA IDENTITY FULL;
CREATE TABLE public.feed_mvp_idle (id int PRIMARY KEY, name text);
INSERT INTO cdc.source (source_id, name) VALUES ('src-mvp', 'MVP-Quelle');
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

# Compose-Healthcheck-Vertrag (LH-FA-ADM-002, LH-QA-OPS-002, slice-012):
# der Feed-Container befragt cdc.heartbeat über seinen eigenen
# `--healthcheck`-Ausgang (compose.yaml); dieser Lauf belegt den
# End-zu-Ende-Vertrag am realen Docker-Health-Status statt nur am
# Binary-Exit-Code (das gepinnte Toolchain-Image trägt Healthcheck-Belege
# nicht direkt).
healthy=0
for _ in $(seq 1 60); do
  health=$(docker inspect --format '{{.State.Health.Status}}' "$FEED_CONTAINER" 2>/dev/null || echo fehlt)
  if [ "$health" = "healthy" ]; then
    healthy=1
    break
  fi
  sleep 1
done
if [ "$healthy" -ne 1 ]; then
  echo "run-integration-tests: Feed-Container meldet Compose-Health-Status ${health:-fehlt}, wollen healthy (cdc.heartbeat-Alter unter der Schwelle)" >&2
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

# Lasttest-Beleg (LH-FA-ADM-004, SPEC-013 CDC_LAG_THRESHOLDS): cdc_capture_lag
# bildet den Abstand zwischen Quelländerung und CDC-Verfügbarkeit ab. Zwei
# Quelltransaktionen auf derselben Tabelle zeigen den Unterschied: eine
# ungehinderte Transaktion liefert einen kleinen Wert; eine Transaktion,
# deren Erfassung durch eine pausierte CDC-Runtime künstlich verzögert wird
# (`docker pause` hält den Feed-Container über die Freezer-Cgroup an, bevor
# die Transaktion verarbeitet ist), liefert einen Wert nahe der
# Pausendauer. `feed_mvp_full` bleibt über den ganzen Lauf aktiviert (anders
# als `feed_mvp_flow`, das der letzte MVP-Testfall deaktiviert). Die
# Pausendauer bleibt unter `wal_sender_timeout=2000` (compose.yaml): eine
# längere Pause ließe die Quelle die Replication-Verbindung selbst beenden,
# bevor der Feed-Container sie fortsetzen kann — der Lasttest misst die
# Erfassungsverzögerung, keine Verbindungsstörung.
LAG_TABLE=feed_mvp_full
LAG_DELAY_SECONDS=${LAG_DELAY_SECONDS:-1}

baseline_count=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.transaction")
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$LAG_TABLE (id, name) VALUES (90, 'LagBaseline');
SQL

baseline_captured=0
for _ in $(seq 1 120); do
  count=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT count(*) FROM cdc.transaction")
  if [ "$count" -gt "$baseline_count" ]; then
    baseline_captured=1
    break
  fi
  sleep 0.25
done
if [ "$baseline_captured" -ne 1 ]; then
  echo "run-integration-tests: Baseline-Transaktion für den Lasttest-Beleg wurde nicht erfasst" >&2
  exit 1
fi
baseline_lag=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_capture_lag'")

delayed_count=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.transaction")
docker pause "$FEED_CONTAINER" >/dev/null
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$LAG_TABLE (id, name) VALUES (91, 'LagDelayed');
SQL
sleep "$LAG_DELAY_SECONDS"
docker unpause "$FEED_CONTAINER" >/dev/null

delayed_captured=0
for _ in $(seq 1 120); do
  count=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT count(*) FROM cdc.transaction")
  if [ "$count" -gt "$delayed_count" ]; then
    delayed_captured=1
    break
  fi
  sleep 0.25
done
if [ "$delayed_captured" -ne 1 ]; then
  echo "run-integration-tests: verzögerte Transaktion für den Lasttest-Beleg wurde nach dem Fortsetzen nicht erfasst" >&2
  exit 1
fi
delayed_lag=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_capture_lag'")

echo "run-integration-tests: Lasttest-Beleg cdc_capture_lag — Baseline ${baseline_lag}s, verzögert ${delayed_lag}s (künstliche Pause ${LAG_DELAY_SECONDS}s)"

# Toleranz proportional zur Pausendauer (nicht als feste Sekundenspanne):
# die Hälfte der künstlichen Verzögerung trennt die beiden Messungen
# unabhängig davon, wie kurz LAG_DELAY_SECONDS gewählt ist.
if ! awk -v v="$delayed_lag" -v d="$LAG_DELAY_SECONDS" 'BEGIN { exit !(v+0 >= d*0.5) }'; then
  echo "run-integration-tests: cdc_capture_lag=$delayed_lag bildet die künstliche Verzögerung von ${LAG_DELAY_SECONDS}s nicht ab" >&2
  exit 1
fi
if ! awk -v v="$baseline_lag" -v d="$LAG_DELAY_SECONDS" 'BEGIN { exit !(v+0 < d*0.5) }'; then
  echo "run-integration-tests: cdc_capture_lag=$baseline_lag der ungehinderten Transaktion liegt nicht unter der künstlichen Verzögerung von ${LAG_DELAY_SECONDS}s — kein Kontrast zur verzögerten Messung" >&2
  exit 1
fi

# End-Beleg des Feed-Containers: der Lauf sieht auch den Ausgang der
# CDC-Runtime — ein Container, der nach der letzten Test-Assertion endet
# (Störung am Stream, Klasse replication), färbt den Lauf rot statt still
# durchzulaufen.
feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  feed_exit=$(docker inspect --format '{{.State.ExitCode}}' "$FEED_CONTAINER" 2>/dev/null || echo fehlt)
  echo "run-integration-tests: Feed-Container endete während des Testlaufs (Lauf: $feed_running, Ausgang: $feed_exit)" >&2
  exit 1
fi
