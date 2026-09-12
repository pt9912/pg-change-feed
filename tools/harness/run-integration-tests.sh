#!/usr/bin/env bash
# run-integration-tests — MVP-Integrationstest gegen die Compose-Umgebung
# (compose.yaml; LH-QA-POR-003). Kette je Lauf: Compose-PostgreSQL frisch
# hochfahren → Schema-Rollout über d-migrate (ADR-0043) → Vorbedingungen
# der Aktivierung (Quell-Tabellen, Quelle-Zeile) → Feed-Container als
# CDC-Runtime starten — seine Verdrahtung aktiviert die Tabellen über den
# EnableTable Use Case (ADR-0028), nicht über Seed-SQL → Toolchain-
# Container gegen das Compose-Netz. Der Feed-Container streamt dabei
# selbst (Verdrahtung je ADR-0026); der Test schreibt nur Quelländerungen
# und liest den Store. Am Ende der Kette steht ein Black-Box-Rundlauf
# (ADR-0030 E2E-Tier): `register-consumer`/`acknowledge-consumer` laufen
# dort ausschließlich als externer `docker exec`-Aufruf gegen den
# laufenden Feed-Container, über einen simulierten Container-Neustart
# hinweg — kein Go-Paket-Import interner Anwendungslogik in diesem
# Abschnitt.
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

# Rollen-DSN-Verifikation gegen den Compose-Stack (LH-QA-SEC-001…003,
# ADR-0047, BEO-PGC/rollen-test-abdeckungsluecken): Die laufende
# CDC_CAPTURE_DSN/CDC_ADMIN_DSN/CDC_READER_DSN-Verdrahtung des
# Feed-Containers verbindet sich testbedingt mit allen drei DSNs über den
# Superuser postgres (siehe compose.yaml-Kommentar) — die drei
# Gruppenrollen selbst liegen seit dem Schema-Rollout oben aber real in
# dieser Instanz vor (nacharbeit-roles.sql). Dieser Abschnitt belegt die
# Rollentrennung direkt gegen sie, mit derselben Fixture-Disziplin wie
# internal/bootstrap/roles_wiring_test.go::newTestLoginRole: anmeldefähige
# Test-Identitäten je Rolle, angelegt und am Ende dieses Abschnitts wieder
# entfernt — kein Bestandteil des Rollen-DDL
# (tools/schema/nacharbeit-roles.sql, unverändert seit slice-011/-023).
#
# Zwei Prüfungen: (1) ein cdc_reader-Login scheitert an einem schreibenden
# Aufruf — dieselbe Fehlerklasse wie
# TestCdcReaderLoginConnectionRejectsWrite. (2) schließt
# BEO-PGC/rollen-test-abdeckungsluecken Punkt (2): Replication-Stream und
# ACK-Adapter (beide cdc_capture-gebunden, ADR-0047 — dieselbe physische
# Verbindung, Option C) waren gegen Rollen-Vertauschung nicht
# testgesichert, weil run-replication-tests.sh keine rollenbeschränkten
# Login-Test-Identitäten bereitstellt. Real geschlossen mit zwei
# Verbindungsversuchen im Replication-Protokoll (`replication=database`):
# eine Login-Identität ohne das REPLICATION-Attribut (hier mit der
# falschen Gruppenrolle cdc_admin) scheitert am Verbindungsaufbau
# ("permission denied to start WAL sender"); eine Login-Identität mit
# cdc_capture-Mitgliedschaft UND direkt gesetztem REPLICATION-Attribut
# (ADR-0047 Kontext-Befund 2: PostgreSQL vererbt Rollen-Attribute nicht
# über Mitgliedschaft) gelingt — der Kontrast belegt, dass die Ablehnung
# oben an der Rolle liegt, nicht an einem allgemeinen Verbindungsfehler.
ROLE_TEST_PASSWORD=pgc-e2e-role-verify
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
DROP ROLE IF EXISTS pgc_e2e_reader_login;
CREATE ROLE pgc_e2e_reader_login LOGIN PASSWORD '$ROLE_TEST_PASSWORD' IN ROLE cdc_reader;
DROP ROLE IF EXISTS pgc_e2e_wrongrole_login;
CREATE ROLE pgc_e2e_wrongrole_login LOGIN PASSWORD '$ROLE_TEST_PASSWORD' IN ROLE cdc_admin;
DROP ROLE IF EXISTS pgc_e2e_capture_login;
CREATE ROLE pgc_e2e_capture_login LOGIN PASSWORD '$ROLE_TEST_PASSWORD' IN ROLE cdc_capture;
ALTER ROLE pgc_e2e_capture_login REPLICATION;
SQL

set +e
reader_write_output=$(docker exec "$PG_CONTAINER" psql \
  "postgres://pgc_e2e_reader_login:$ROLE_TEST_PASSWORD@localhost:5432/$PG_DB?sslmode=disable" \
  -c "INSERT INTO cdc.change (change_id, transaction_id, source_table_id, sequence, operation, old_data, new_data, schema_version) VALUES ('c-role-verify','tx-role-verify','st-role-verify',1,'INSERT',NULL,'{}'::jsonb,'sv-role-verify')" 2>&1)
reader_write_status=$?
set -e
if [ "$reader_write_status" -eq 0 ] || ! printf '%s' "$reader_write_output" | grep -q "permission denied"; then
  echo "run-integration-tests: cdc_reader-Login-Identität INSERT auf cdc.change — erwartet permission denied, erhalten (Status $reader_write_status): $reader_write_output" >&2
  exit 1
fi

set +e
wrongrole_replication_output=$(docker exec "$PG_CONTAINER" psql \
  "postgres://pgc_e2e_wrongrole_login:$ROLE_TEST_PASSWORD@localhost:5432/$PG_DB?sslmode=disable&replication=database" \
  -c "IDENTIFY_SYSTEM;" 2>&1)
wrongrole_replication_status=$?
set -e
if [ "$wrongrole_replication_status" -eq 0 ] || ! printf '%s' "$wrongrole_replication_output" | grep -q "permission denied to start WAL sender"; then
  echo "run-integration-tests: cdc_admin-Login-Identität (falsche Rolle für eine Replication-Verbindung) — erwartet 'permission denied to start WAL sender', erhalten (Status $wrongrole_replication_status): $wrongrole_replication_output" >&2
  exit 1
fi

set +e
capture_replication_output=$(docker exec "$PG_CONTAINER" psql \
  "postgres://pgc_e2e_capture_login:$ROLE_TEST_PASSWORD@localhost:5432/$PG_DB?sslmode=disable&replication=database" \
  -c "IDENTIFY_SYSTEM;" 2>&1)
capture_replication_status=$?
set -e
if [ "$capture_replication_status" -ne 0 ]; then
  echo "run-integration-tests: cdc_capture-Login-Identität (korrekte Rolle + REPLICATION-Attribut) — erwartet Erfolg als Kontrast zur vorigen Ablehnung, erhalten (Status $capture_replication_status): $capture_replication_output" >&2
  exit 1
fi

docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -c \
  "DROP ROLE pgc_e2e_reader_login; DROP ROLE pgc_e2e_wrongrole_login; DROP ROLE pgc_e2e_capture_login;" >/dev/null

echo "run-integration-tests: Rollen-DSN-Verifikation belegt — cdc_reader-Login-Identität schreibender Zugriff abgelehnt, Replication-Verbindung ohne REPLICATION-Attribut abgelehnt, mit cdc_capture-Mitgliedschaft+REPLICATION-Attribut erfolgreich (BEO-PGC/rollen-test-abdeckungsluecken Punkt 2 geschlossen)"

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

# Obergrenze für einen überschriebenen LAG_DELAY_SECONDS-Wert: oberhalb
# von wal_sender_timeout=2000 (compose.yaml) beendet PostgreSQL die
# Replication-Verbindung selbst, bevor der pausierte Feed-Container sie
# fortsetzen kann — der Container trägt keine Restart-Policy und bliebe
# beendet stehen. 1,5s Sicherheitsabstand zum 2s-Timeout hält auch bei
# Docker-Scheduling-Varianz einen Puffer.
LAG_DELAY_SECONDS_MAX=1.5
if ! awk -v d="$LAG_DELAY_SECONDS" -v m="$LAG_DELAY_SECONDS_MAX" 'BEGIN { exit !(d+0 <= m+0) }'; then
  echo "run-integration-tests: LAG_DELAY_SECONDS=$LAG_DELAY_SECONDS überschreitet die Obergrenze ${LAG_DELAY_SECONDS_MAX}s — darüber beendet PostgreSQL die Replication-Verbindung selbst (wal_sender_timeout=2000 aus compose.yaml), und der Feed-Container ohne Restart-Policy bliebe beendet stehen" >&2
  exit 1
fi

# Die IDs 90/91 liegen in einem eigenen Wertebereich, getrennt von den
# MVP-Referenzzeilen in test/integration/mvp_test.go (id=1, id=7 auf
# derselben Tabelle feed_mvp_full) — keine Kollision zwischen den beiden
# Testfall-Gruppen.
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

# Black-Box-CLI-Rundlauf (LH-QA-POR-003, ADR-0030 E2E-Tier): anders als der
# Go-Testlauf oben (`go test ./test/integration/...`, der intern gegen
# `postgresstorage`/`bootstrap` läuft) ruft dieser Abschnitt
# `register-consumer`/`acknowledge-consumer` ausschließlich als externen
# Prozess gegen den laufenden, containerisierten Produktions-Binary auf
# (`docker exec … /pg-change-feed …`, distroless — kein Shell im
# Feed-Container, daher kein `sh -c`-Umweg nötig). Lesen bleibt über den
# bestehenden externen Lesezugriffsweg `cdc.changes` (LH-FA-SST-002) — die
# CLI trägt keinen Lese-Unterbefehl. `feed_mvp_full` trägt denselben Grund
# wie beim Lasttest-Beleg oben: die Tabelle bleibt über den ganzen Lauf
# aktiviert. Die IDs 95/96 liegen in einem eigenen Wertebereich, getrennt
# von den MVP-Referenzzeilen (id=1) und dem Lasttest-Beleg (id=90/91) auf
# derselben Tabelle.
exec_feed() {
  docker exec "$FEED_CONTAINER" /pg-change-feed "$@"
}

CLI_CONSUMER=cli-e2e-consumer
CLI_TABLE=feed_mvp_full

if ! exec_feed register-consumer "$CLI_CONSUMER"; then
  echo "run-integration-tests: register-consumer (extern, docker exec) endete mit einem Fehler" >&2
  exit 1
fi

registered=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.consumer WHERE consumer_id = '$CLI_CONSUMER'")
if [ "$registered" != "1" ]; then
  echo "run-integration-tests: cdc.consumer trägt $CLI_CONSUMER nicht nach dem externen register-consumer-Aufruf" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$CLI_TABLE (id, name) VALUES (95, 'CliE2EFirst');
SQL

first_position=""
for _ in $(seq 1 120); do
  first_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-mvp' AND table_name = '$CLI_TABLE' AND new_data->>'id' = '95'")
  if [ -n "$first_position" ]; then
    break
  fi
  sleep 0.25
done
if [ -z "$first_position" ]; then
  echo "run-integration-tests: erste CLI-Rundlauf-Änderung (id=95) wurde nicht innerhalb der Zeitspanne über cdc.changes gelesen" >&2
  exit 1
fi

if ! exec_feed acknowledge-consumer "$CLI_CONSUMER" "$first_position"; then
  echo "run-integration-tests: acknowledge-consumer (extern, docker exec) endete mit einem Fehler (Position $first_position)" >&2
  exit 1
fi

acked_first=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT acknowledged_position FROM cdc.consumer_position WHERE consumer_id = '$CLI_CONSUMER' AND source_id = 'src-mvp'")
if [ "$acked_first" != "$first_position" ]; then
  echo "run-integration-tests: cdc.consumer_position trägt $acked_first, wollen $first_position (erste externe Bestätigung)" >&2
  exit 1
fi

# Simulierter Container-Neustart: `docker restart` sendet SIGTERM und
# startet denselben Container neu (kein `compose down`/`up`, Slot und
# Publication bleiben auf der PostgreSQL-Seite unberührt) — die CDC-Runtime
# durchläuft dabei einen echten Prozess-Neustart, nicht nur eine
# `docker pause`/`unpause`-Unterbrechung wie beim Lasttest-Beleg oben.
docker restart "$FEED_CONTAINER" >/dev/null

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
  echo "run-integration-tests: Feed-Container meldet nach dem simulierten Neustart Health-Status ${health:-fehlt}, wollen healthy" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$CLI_TABLE (id, name) VALUES (96, 'CliE2ESecond');
SQL

second_position=""
for _ in $(seq 1 120); do
  second_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-mvp' AND table_name = '$CLI_TABLE' AND new_data->>'id' = '96'")
  if [ -n "$second_position" ]; then
    break
  fi
  sleep 0.25
done
if [ -z "$second_position" ]; then
  echo "run-integration-tests: zweite CLI-Rundlauf-Änderung (id=96, nach dem simulierten Neustart) wurde nicht innerhalb der Zeitspanne über cdc.changes gelesen — die Erfassung hat den Neustart nicht fortgesetzt" >&2
  exit 1
fi

# Neustart-Beleg: ein unabhängiger Lesezugriff fragt die bestätigte
# Position erneut aus cdc.consumer_position ab, statt sie im Skript
# weiterzutragen — genau das, was ein neu gestarteter Consumer-Prozess
# täte. Der anschließende Lesezugriff über cdc.changes zeigt das
# Fortsetzen ab dieser Position: die bereits bestätigte Änderung (id=95)
# wiederholt sich nicht, nur die neue (id=96) erscheint.
restored_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT acknowledged_position FROM cdc.consumer_position WHERE consumer_id = '$CLI_CONSUMER' AND source_id = 'src-mvp'")
if [ "$restored_position" != "$first_position" ]; then
  echo "run-integration-tests: aus cdc.consumer_position zurückgelesene Position nach dem Neustart = $restored_position, wollen $first_position (die zuvor bestätigte)" >&2
  exit 1
fi

resumed_ids=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT string_agg(new_data->>'id', ',' ORDER BY commit_position) FROM cdc.changes WHERE source_id = 'src-mvp' AND table_name = '$CLI_TABLE' AND commit_position > $restored_position")
if [ "$resumed_ids" != "96" ]; then
  echo "run-integration-tests: Fortsetzen ab der bestätigten Position (Neustart-Beleg) — ab commit_position=$restored_position gelesene id-Folge '$resumed_ids', wollen '96' (id=95 bleibt hinter der bestätigten Position, keine Wiederholung)" >&2
  exit 1
fi

if ! exec_feed acknowledge-consumer "$CLI_CONSUMER" "$second_position"; then
  echo "run-integration-tests: zweite externe Bestätigung (acknowledge-consumer, docker exec) endete mit einem Fehler (Position $second_position)" >&2
  exit 1
fi

acked_second=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT acknowledged_position FROM cdc.consumer_position WHERE consumer_id = '$CLI_CONSUMER' AND source_id = 'src-mvp'")
if [ "$acked_second" != "$second_position" ]; then
  echo "run-integration-tests: cdc.consumer_position trägt $acked_second, wollen $second_position (zweite externe Bestätigung, nach dem simulierten Neustart)" >&2
  exit 1
fi

echo "run-integration-tests: Black-Box-CLI-Rundlauf belegt — register-consumer/acknowledge-consumer extern (docker exec), Fortsetzen nach simuliertem Neustart ab Position $first_position, Endposition $second_position"
