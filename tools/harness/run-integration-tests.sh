#!/usr/bin/env bash
# run-integration-tests — Compose-Integrationstest gegen die
# Compose-Umgebung (compose.yaml; LH-QA-POR-003), inhaltlich über den
# ursprünglichen MVP-Zuschnitt hinausgewachsen (Rollen-DSN-Trennung,
# Lasttest-Beleg, Black-Box-CLI-Rundlauf). Kette je Lauf:
# Compose-PostgreSQL frisch hochfahren → Schema-Rollout über d-migrate
# (ADR-0043) → Rollen-DSN-Verifikation direkt gegen die Instanz
# (LH-QA-SEC-001…003, ADR-0047) → Vorbedingungen der Aktivierung
# (Quell-Tabellen, Quelle-Zeile) → Feed-Container als CDC-Runtime
# starten — seine Verdrahtung aktiviert die Tabellen über den
# EnableTable Use Case (ADR-0028), nicht über Seed-SQL → Toolchain-
# Container gegen das Compose-Netz. Der Feed-Container streamt dabei
# selbst (Verdrahtung je ADR-0026); der Test schreibt nur Quelländerungen
# und liest den Store. Am Ende der Kette steht ein Black-Box-Rundlauf
# (ADR-0030 E2E-Tier): `register-consumer`/`acknowledge-consumer` laufen
# dort ausschließlich als externer `docker exec`-Aufruf gegen den
# laufenden Feed-Container, über einen simulierten Container-Neustart
# hinweg — kein Go-Paket-Import interner Anwendungslogik in diesem
# Abschnitt. Danach ein CLI-Diagnose-Beleg (LH-FA-SST-003, deckt
# LH-FA-ADM-002…005, slice-038): derselbe externe `docker exec`-Zugriffsweg
# gegen den `diagnose`-Sondermodus, einmal im Normalbetrieb und einmal mit
# einem direkt in `cdc.process_heartbeat` geschriebenen Fehlerzustand
# (LH-FA-ADM-003 Boundary). Ergänzend ein Retention-Sichtbarkeits-Beleg
# (deckt LH-FA-RET-005/006) über denselben `diagnose`-Aufruf: „kein
# Blocker" vor jeder Consumer-Bestätigung, dann ein real blockierender
# Consumer samt `cdc_storage_bytes`. Direkt davor ein kombinierter
# Retention-Lebenszyklus-Rundlauf (deckt LH-FA-RET-002…006 in einer Kette):
# eine eigene, isolierte Zeile und ein eigener Consumer durchlaufen real
# Blocker-Sichtbarkeit über `cdc.retention_blockers`, Bestätigung über die
# Position hinweg, die reale Abwesenheit jeder Zeile für `src-e2e` danach,
# die reale Löschung sowie die durchgehende numerische Lesbarkeit von
# `cdc_storage_bytes`. Ein Upgrade-Sicherheits-Rundlauf (LH-QA-OPS-005,
# ADR-0064) folgt danach: ein realer Container-Tausch über
# `$COMPOSE up -d --force-recreate --no-deps pg-change-feed` ersetzt den
# Feed-Container durch eine neue Instanz desselben `:dev`-Images, während
# `postgres`/`nats` unberührt bleiben — Datenstand vor dem Tausch bleibt
# über `cdc.changes` identisch lesbar, eine danach eingefügte Zeile wird
# weiterhin erfasst.
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
NATS_CONTAINER=${NATS_CONTAINER:-cdc-test-nats}
PG_DB=cdc
PG_USER=postgres
PG_PASSWORD=postgres
# Der Slot-Name trägt denselben Wert wie der Container-Vertrag in
# compose.yaml (CDC_SLOT); der Runner liest ihn im Start- und im
# End-Wächter.
SLOT=slot_pgc_e2e

# Die drei HTTP-API-Werte tragen denselben Wert wie der Container-Vertrag
# in compose.yaml (CDC_HTTP_ADDR/CDC_API_TOKEN_READER/CDC_API_TOKEN_ADMIN,
# ADR-0057): der Toolchain-Container erreicht den Feed-Container über den
# Netzwerk-Alias `pg-change-feed`, den Compose aus dem Service-Namen
# vergibt (dasselbe Muster wie `nats://nats:4222` für den NATS-Server).
HTTP_BASE_URL="http://pg-change-feed:8090"
HTTP_TOKEN_READER=e2e-reader-token
HTTP_TOKEN_ADMIN=e2e-admin-token

DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"

# --- E2E-Abdeckungstabelle (docs/user/e2e-abdeckung.md) --------------------
# Das Erzeugnis dieses Laufs entsteht aus zwei Hälften: die Go-Zeilen aus dem
# Quelltext des Testpakets (`TestAbdeckungstabelleZeilen`, go/parser), die
# Bash-Zeilen deklariert jede Phase über ihren Anker (`abdeckung_declare`).
# Findet der Runner einen Deklarations-Anker nicht mehr, bricht er ab, bevor
# er die Datei anfasst; die Datei wird nur bei inhaltlicher Abweichung
# geschrieben (Temp-Datei + cmp). Geschrieben wird auf dem Host — der
# Toolchain-Container läuft gegen ein read-only Bind-Mount.
ABDECKUNG_ZIEL=docs/user/e2e-abdeckung.md
ABDECKUNG_QUELLE=$(realpath --relative-to="$(pwd)" "${BASH_SOURCE[0]}")
ABDECKUNG_GO_ZEILEN=""
ABDECKUNG_KOPF='# E2E-Abdeckung je Spec-Kennung

Erzeugt von `make test-integration` über `tools/harness/run-integration-tests.sh`:
die Go-Zeilen leitet das Testpaket aus seinem eigenen Quelltext ab
(`TestAbdeckungstabelleZeilen` in `test/integration/integration_test.go`),
die Bash-Zeilen deklariert jede Phase des Runners an Ort und Stelle über
einen Anker. Diese Datei ist eine **stabile Abdeckungs-Deklaration**, kein
Lauf-Beleg: sie ändert sich mit den Nachweis-Deklarationen, nicht mit jedem
Lauf — der Runner schreibt sie nur bei inhaltlicher Abweichung. Die
Beschreibungsspalte trägt keine Kennungen; die Aussage eines Nachweises steht
über die Spalte `Ort` an ihrer Quelle.
'
abdeckung_bash_zeilen=()

# abdeckung_declare <Nachweis> <Kennungen> <Kurzbeschreibung> <Anker>
# Der Anker ist ein wörtlicher Ausschnitt aus der Phase dahinter — der
# Runner sucht ihn ab der Deklarations-Zeile im eigenen Quelltext und trägt
# die gefundene Zeile als `Ort` ein; ein Anker, der in keiner späteren Zeile
# mehr steht, beendet den Lauf.
abdeckung_declare() {
  local nachweis=$1 kennungen=$2 kurzbeschreibung=$3 anker=$4
  local treffer
  treffer=$(awk -v ab="${BASH_LINENO[0]}" -v anker="$anker" \
    'NR > ab && index($0, anker) { print NR; exit }' "$ABDECKUNG_QUELLE")
  if [ -z "$treffer" ]; then
    echo "run-integration-tests: Deklarations-Anker der E2E-Abdeckungstabelle nicht gefunden — Phase '$nachweis', Anker '$anker' (der Anker gehört in eine Zeile hinter dem Deklarations-Aufruf)" >&2
    exit 1
  fi
  abdeckung_bash_zeilen+=("$ABDECKUNG_QUELLE|$treffer|$nachweis|$kennungen|$kurzbeschreibung")
}

# abdeckung_render liest `Quelldatei|Quellzeile|Nachweis|Kennungen|
# Kurzbeschreibung` und schreibt daraus die Tabellenzeile; die Kennungsspalte
# trägt je Kennung den Link auf ihr Definitionsdokument (`ids` des
# Doku-Gates).
abdeckung_render() {
  local quelldatei zeile nachweis kennungen kurzbeschreibung kennung ziel spalte trenner
  while IFS='|' read -r quelldatei zeile nachweis kennungen kurzbeschreibung; do
    if [ -z "$nachweis" ]; then
      continue
    fi
    spalte=""
    trenner=""
    for kennung in ${kennungen//,/ }; do
      case "$kennung" in
        LH-*.[a-z]) ziel="../../spec/pflichtenheft.md" ;;
        SPEC-*) ziel="../../spec/pflichtenheft.md" ;;
        *) ziel="../../spec/lastenheft.md" ;;
      esac
      spalte="${spalte}${trenner}[\`$kennung\`]($ziel)"
      trenner=", "
    done
    printf '| %s | `%s` | `%s:%s` | %s |\n' "$spalte" "$nachweis" "$quelldatei" "$zeile" "$kurzbeschreibung"
  done
}

# abdeckung_go_zeilen_lesen setzt ABDECKUNG_GO_ZEILEN; der Erzeuger läuft im
# bestehenden Testpaket und endet sichtbar, wenn eine `func TestE2E*` keine
# Spec-Kennung trägt.
abdeckung_go_zeilen_lesen() {
  local ausgabe
  if ! ausgabe=$(docker run --rm --network "$NETWORK" \
      -v "$(pwd)":/src:ro \
      -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
      -w /src \
      -e GOCACHE=/tmp/gocache \
      "$TOOLCHAIN_IMAGE" go test -v -run '^TestAbdeckungstabelleZeilen$' ./test/integration/... 2>&1); then
    echo "run-integration-tests: Erzeuger der E2E-Abdeckungstabelle (Go-Hälfte) endete mit einem Fehler: $ausgabe" >&2
    exit 1
  fi
  ABDECKUNG_GO_ZEILEN=$(printf '%s\n' "$ausgabe" | sed -n 's/^.*ABDECKUNG|//p' | sort -t'|' -k1,1 -k2,2n)
  if [ -z "$ABDECKUNG_GO_ZEILEN" ]; then
    echo "run-integration-tests: Erzeuger der E2E-Abdeckungstabelle lieferte keine Go-Zeile (Testausgang ohne ABDECKUNG-Zeile)" >&2
    exit 1
  fi
}

# abdeckung_schreiben <Go-Zeilen> setzt beide Hälften zusammen — Go-Zeilen
# nach Quelldatei und Quellzeile, dann Bash-Zeilen nach Runner-Zeile — und
# schreibt nur bei inhaltlicher Abweichung.
abdeckung_schreiben() {
  local bash_zeilen="" temp
  if [ "${#abdeckung_bash_zeilen[@]}" -gt 0 ]; then
    bash_zeilen=$(printf '%s\n' "${abdeckung_bash_zeilen[@]}" | sort -t'|' -k2,2n)
  fi
  temp=$(mktemp)
  {
    printf '%s\n' "$ABDECKUNG_KOPF"
    printf '| Spec-Kennung | Nachweis | Ort | Kurzbeschreibung |\n'
    printf '| --- | --- | --- | --- |\n'
    printf '%s\n' "$1" | abdeckung_render
    printf '%s\n' "$bash_zeilen" | abdeckung_render
  } > "$temp"
  if [ -f "$ABDECKUNG_ZIEL" ] && cmp -s "$temp" "$ABDECKUNG_ZIEL"; then
    rm -f "$temp"
    echo "run-integration-tests: E2E-Abdeckungstabelle unverändert — $ABDECKUNG_ZIEL entspricht dem Quelltext-Stand"
  else
    chmod 0644 "$temp"
    mv "$temp" "$ABDECKUNG_ZIEL"
    echo "run-integration-tests: E2E-Abdeckungstabelle geschrieben — $ABDECKUNG_ZIEL"
  fi
}

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

abdeckung_declare "Rollen-DSN-Verifikation" "LH-QA-SEC-001,LH-QA-SEC-002,LH-QA-SEC-003" "die drei Gruppenrollen real gegeneinander geprüft: ein Reader-Login scheitert am schreibenden Aufruf, eine Replication-Verbindung ohne REPLICATION-Attribut scheitert, dieselbe Verbindung mit der Capture-Rolle gelingt" "Rollen-DSN-Verifikation belegt — cdc_reader-Login-Identität schreibender Zugriff abgelehnt"

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
CREATE TABLE public.feed_e2e_flow (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_full (id int PRIMARY KEY, name text);
ALTER TABLE public.feed_e2e_full REPLICA IDENTITY FULL;
CREATE TABLE public.feed_e2e_idle (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_schema (id int PRIMARY KEY, name text, amount text);
CREATE TABLE public.feed_e2e_sql_admin (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_walsender_timing (id int PRIMARY KEY, name text);
INSERT INTO cdc.source (source_id, name) VALUES ('src-e2e', 'E2E-Quelle');
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

# TestE2ESchemaChangeIncompatibleTypeChange meldet ihren
# Negative-Fall sichtbar über die Fehlerklasse `schema`
# (`mapper.ErrIncompatibleSchemaChange`) — der Erfassungspfad des
# Feed-Containers endet darüber (`bootstrap.Run` -> `os.Exit(1)`), und
# `restart: "no"` in compose.yaml trägt danach keinen Neustart-Vertrag:
# der Container bleibt beendet stehen. Diese Testfunktion läuft deshalb
# separat und zuletzt (siehe unten, nach Lasttest-Beleg und
# Black-Box-CLI-Rundlauf) — jede andere Testfunktion dieses Pakets läuft
# hier, solange der Feed-Container noch gebraucht wird. Eine künftig
# ergänzte Testfunktion gehört in dieses `-run`-Muster, sofern sie den
# laufenden Container nicht ebenfalls beendet. `TestE2ESchemaChangeDropColumn`
# (`LH-FA-SCH-003`, `ADR-0063` Supersedes `ADR-0058` Entscheidung 1) beendet
# den Erfassungspfad seit der Testform-Korrektur ebenfalls dauerhaft (derselbe
# `relationOther`/`ErrIncompatibleSchemaChange`-Pfad wie eine inkompatible
# Typänderung) und läuft deshalb NICHT hier, sondern als eigener Aufruf nach
# der Container-Ende-Grenze, mit einem expliziten Neustart davor — siehe dort.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_INTEGRATION_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test -v \
  -run '^(TestE2ECaptureFlow|TestE2EUpdateOldImageWithFullReplicaIdentity|TestE2EChangesViewMatchesReadChanges|TestE2ERetentionBlockersViewShowsFurthestBehindConsumer|TestE2EMetricsCarriesStorageBytes|TestE2EActivationState|TestE2EActiveTablesViewMatchesActivationState|TestE2EDisableRetainedState|TestE2ESchemaChangeAddColumn|TestE2EChangeTableMetadataExtensibility|TestE2EHeartbeatHealthy)$' \
  ./test/integration/...

# Go-Hälfte der E2E-Abdeckungstabelle: derselbe Testlauf führt den Erzeuger
# aus (eigener `-run`-Aufruf, dieselbe Toolchain) und legt seinen Zeilensatz
# in ABDECKUNG_GO_ZEILEN ab; eine `func TestE2E*` ohne Spec-Kennung endet
# hier sichtbar, bevor die Datei geschrieben wird.
abdeckung_go_zeilen_lesen

abdeckung_declare "Spaltenausschluss-Rundlauf" "LH-FA-CFG-005,LH-QA-SEC-004" "eine nicht gelistete Tabelle wird über den SQL-Antragsweg aktiviert, der Spaltenausschluss real verarbeitet — die danach erfasste Change trägt den Schlüssel nicht mehr im Row Image, die davor erfasste bleibt mit ihrem Wert unverändert lesbar" "Spaltenausschluss-Rundlauf (LH-FA-CFG-005 Happy Path, ADR-0059) belegt"

# Spaltenausschluss-Rundlauf (LH-FA-CFG-005, ADR-0059): der reale Pfad
# SQL-Antrag → Live-Reload → gefiltertes Row Image am laufenden
# Feed-Container. Eine eigene, dedizierte Tabelle wird über dieselbe
# Antrags-Queue aktiviert wie im SQL-Administration-Live-Reload-Beleg
# (kein CDC_TABLES-Eintrag, kein Neustart); `SELECT cdc.exclude_column(...)`
# bestätigt nur „beantragt" — der Poll wartet auf `status = 'applied'`,
# bevor die Filterwirkung geprüft wird. Läuft vor der Container-Ende-Grenze
# unten (TestE2ESchemaChangeDropColumn/-IncompatibleTypeChange), die den
# Erfassungspfad dauerhaft beendet.
COLUMN_TABLE=feed_e2e_column_exclusion
COLUMN_NAME=secret
COLUMN_VALUE_BEFORE=ColumnBeforeExclusionSentinel
COLUMN_VALUE_AFTER=ColumnAfterExclusionSentinel

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$COLUMN_TABLE (id int PRIMARY KEY, name text, $COLUMN_NAME text);
SQL

column_enable_request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT cdc.enable_table('src-e2e', 'public', '$COLUMN_TABLE')")
if [ -z "$column_enable_request_id" ]; then
  echo "run-integration-tests: cdc.enable_table($COLUMN_TABLE) lieferte keine Antrags-ID" >&2
  exit 1
fi

column_enable_status=""
column_enable_applied=0
for _ in $(seq 1 60); do
  column_enable_status=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$column_enable_request_id'")
  if [ "$column_enable_status" = "applied" ]; then
    column_enable_applied=1
    break
  fi
  if [ "$column_enable_status" = "failed" ]; then
    error_message=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT error_message FROM cdc.administration_request WHERE administration_request_id = '$column_enable_request_id'")
    echo "run-integration-tests: cdc.enable_table($COLUMN_TABLE)-Antrag $column_enable_request_id scheiterte: $error_message" >&2
    exit 1
  fi
  sleep 0.5
done
if [ "$column_enable_applied" -ne 1 ]; then
  echo "run-integration-tests: cdc.enable_table($COLUMN_TABLE)-Antrag $column_enable_request_id wurde nicht innerhalb der Zeitspanne von der Administrations-Goroutine verarbeitet (status=${column_enable_status:-leer})" >&2
  exit 1
fi

# Baseline vor dem Ausschluss: die Spalte wird ohne Ausschluss real
# erfasst — dieser Beleg trennt „Wert abwesend" nach dem Ausschluss von
# „die Spalte war nie Teil des Row Image".
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$COLUMN_TABLE (id, name, $COLUMN_NAME) VALUES (1, 'ColumnBefore', '$COLUMN_VALUE_BEFORE');
SQL

column_before_value=""
for _ in $(seq 1 120); do
  column_before_value=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT new_data->>'$COLUMN_NAME' FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$COLUMN_TABLE' AND new_data->>'id' = '1'")
  if [ -n "$column_before_value" ]; then
    break
  fi
  sleep 0.25
done
if [ "$column_before_value" != "$COLUMN_VALUE_BEFORE" ]; then
  echo "run-integration-tests: Spaltenausschluss-Rundlauf — die vor dem Ausschluss erfasste Change (id=1, $COLUMN_TABLE) trägt $COLUMN_NAME='${column_before_value:-leer}', wollen '$COLUMN_VALUE_BEFORE' (Baseline: ohne Ausschluss wird die Spalte real erfasst)" >&2
  exit 1
fi

column_exclude_request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT cdc.exclude_column('src-e2e', 'public', '$COLUMN_TABLE', '$COLUMN_NAME')")
if [ -z "$column_exclude_request_id" ]; then
  echo "run-integration-tests: cdc.exclude_column($COLUMN_TABLE.$COLUMN_NAME) lieferte keine Antrags-ID" >&2
  exit 1
fi

column_exclude_status=""
column_exclude_applied=0
for _ in $(seq 1 60); do
  column_exclude_status=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$column_exclude_request_id'")
  if [ "$column_exclude_status" = "applied" ]; then
    column_exclude_applied=1
    break
  fi
  if [ "$column_exclude_status" = "failed" ]; then
    error_message=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT error_message FROM cdc.administration_request WHERE administration_request_id = '$column_exclude_request_id'")
    echo "run-integration-tests: cdc.exclude_column($COLUMN_TABLE.$COLUMN_NAME)-Antrag $column_exclude_request_id scheiterte: $error_message" >&2
    exit 1
  fi
  sleep 0.5
done
if [ "$column_exclude_applied" -ne 1 ]; then
  echo "run-integration-tests: cdc.exclude_column($COLUMN_TABLE.$COLUMN_NAME)-Antrag $column_exclude_request_id wurde nicht innerhalb der Zeitspanne von der Administrations-Goroutine verarbeitet (status=${column_exclude_status:-leer})" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$COLUMN_TABLE (id, name, $COLUMN_NAME) VALUES (2, 'ColumnAfter', '$COLUMN_VALUE_AFTER');
SQL

column_after_present=""
for _ in $(seq 1 120); do
  column_after_present=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$COLUMN_TABLE' AND new_data->>'id' = '2'")
  if [ "$column_after_present" = "1" ]; then
    break
  fi
  sleep 0.25
done
if [ "$column_after_present" != "1" ]; then
  echo "run-integration-tests: Spaltenausschluss-Rundlauf — die nach dem Ausschluss erfasste Change (id=2, $COLUMN_TABLE) wurde nicht innerhalb der Zeitspanne über cdc.changes gelesen" >&2
  exit 1
fi

# Drei Aussagen zur danach erfassten Change: der ausgeschlossene
# Spaltenschlüssel fehlt im Row Image; die nicht ausgeschlossene Spalte
# bleibt darin (Kontrolle gegen ein leeres Row Image als Alternativerklärung);
# der ausgeschlossene Wert kommt im ganzen persistierten Change nicht vor
# (LH-QA-SEC-004: der Wert erreicht die Persistenzschicht nie).
column_after_key=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT jsonb_exists(new_data, '$COLUMN_NAME') FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$COLUMN_TABLE' AND new_data->>'id' = '2'")
column_after_name=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT new_data->>'name' FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$COLUMN_TABLE' AND new_data->>'id' = '2'")
column_after_leak=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT (coalesce(new_data::text, '') LIKE '%$COLUMN_VALUE_AFTER%') OR (coalesce(old_data::text, '') LIKE '%$COLUMN_VALUE_AFTER%') FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$COLUMN_TABLE' AND new_data->>'id' = '2'")
if [ "$column_after_key" != "f" ]; then
  echo "run-integration-tests: Spaltenausschluss-Rundlauf — die nach dem Ausschluss erfasste Change (id=2, $COLUMN_TABLE) trägt den ausgeschlossenen Spaltenschlüssel $COLUMN_NAME weiterhin im Row Image" >&2
  exit 1
fi
if [ "$column_after_name" != "ColumnAfter" ]; then
  echo "run-integration-tests: Spaltenausschluss-Rundlauf — die nach dem Ausschluss erfasste Change (id=2, $COLUMN_TABLE) trägt die nicht ausgeschlossene Spalte name='${column_after_name:-leer}', wollen 'ColumnAfter'" >&2
  exit 1
fi
if [ "$column_after_leak" != "f" ]; then
  echo "run-integration-tests: Spaltenausschluss-Rundlauf — der ausgeschlossene Wert '$COLUMN_VALUE_AFTER' steht im persistierten Change (id=2, $COLUMN_TABLE) (LH-QA-SEC-004)" >&2
  exit 1
fi

# Der vor dem Ausschluss erfasste Change bleibt über cdc.changes unverändert
# lesbar und trägt seinen damals erfassten Wert weiter: die Filterung liegt
# in der Row-Image-Konstruktion (ADR-0059 Teilfrage 3; cdc.changes ist eine
# reine Projektion über cdc.change) und schreibt persistierte Changes nicht
# rückwirkend um — LH-FA-CFG-005 adressiert „künftige Changes von t".
column_history_value=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT new_data->>'$COLUMN_NAME' FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$COLUMN_TABLE' AND new_data->>'id' = '1'")
if [ "$column_history_value" != "$COLUMN_VALUE_BEFORE" ]; then
  echo "run-integration-tests: Spaltenausschluss-Rundlauf — die vor dem Ausschluss erfasste Change (id=1, $COLUMN_TABLE) liest sich nach dem Ausschluss nicht mehr unverändert (${COLUMN_NAME}='${column_history_value:-leer}', wollen '$COLUMN_VALUE_BEFORE')" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem Spaltenausschluss-Rundlauf nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: Spaltenausschluss-Rundlauf (LH-FA-CFG-005 Happy Path, ADR-0059) belegt — cdc.exclude_column($COLUMN_TABLE.$COLUMN_NAME) wurde ohne Neustart verarbeitet (status=applied), die danach erfasste Change (id=2) trägt $COLUMN_NAME nicht im Row Image und den Wert '$COLUMN_VALUE_AFTER' nirgends, die davor erfasste Change (id=1) bleibt mit '$COLUMN_VALUE_BEFORE' unverändert lesbar"

abdeckung_declare "Spaltenausschluss-Neustart-Beleg" "LH-FA-CFG-005,LH-QA-SEC-004" "nach einem realen Container-Neustart leitet der Prozessstart den dauerhaften Ausschlussstand aus den applied-Zeilen der Spalten-Anträge ab — die danach erfasste Change trägt den Schlüssel nicht mehr, die nicht ausgeschlossene Spalte bleibt" "Spaltenausschluss-Neustart-Beleg (ADR-0065, LH-QA-SEC-004)"

# Neustart-Beleg des dauerhaften Ausschlussstandes (ADR-0065): der
# Prozessstart leitet den Stand aus den `applied`-Zeilen der beiden
# Spalten-Antragsarten ab (`activatedTableBindings` → `ColumnExclusionPort.
# ExcludedColumns`) und trägt ihn in die Bindungen des Assemblers; ohne
# diesen Schritt erfasst der neu gestartete Prozess die Spalte wieder. Der
# `docker restart` durchläuft einen echten Prozess-Neustart (SIGTERM,
# derselbe Container, Slot und Publication auf der PostgreSQL-Seite
# unberührt). Die Tabelle ist in keinem `CDC_TABLES`-Eintrag — ihr
# Bindungs-Stand kommt aus `cdc.source_table`, ihr Ausschlussstand allein
# über den hier geprüften Weg. Läuft vor der Container-Ende-Grenze unten
# (TestE2ESchemaChangeDropColumn/-IncompatibleTypeChange).
COLUMN_VALUE_RESTART=ColumnAfterRestartSentinel

docker restart "$FEED_CONTAINER" >/dev/null

column_restart_healthy=0
for _ in $(seq 1 60); do
  health=$(docker inspect --format '{{.State.Health.Status}}' "$FEED_CONTAINER" 2>/dev/null || echo fehlt)
  if [ "$health" = "healthy" ]; then
    column_restart_healthy=1
    break
  fi
  sleep 1
done
if [ "$column_restart_healthy" -ne 1 ]; then
  echo "run-integration-tests: Spaltenausschluss-Neustart-Beleg — Feed-Container meldet nach dem Neustart Health-Status ${health:-fehlt}, wollen healthy" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$COLUMN_TABLE (id, name, $COLUMN_NAME) VALUES (3, 'ColumnAfterRestart', '$COLUMN_VALUE_RESTART');
SQL

column_restart_present=""
for _ in $(seq 1 120); do
  column_restart_present=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$COLUMN_TABLE' AND new_data->>'id' = '3'")
  if [ "$column_restart_present" = "1" ]; then
    break
  fi
  sleep 0.25
done
if [ "$column_restart_present" != "1" ]; then
  echo "run-integration-tests: Spaltenausschluss-Neustart-Beleg — die nach dem Neustart erfasste Change (id=3, $COLUMN_TABLE) wurde nicht innerhalb der Zeitspanne über cdc.changes gelesen — die Erfassung hat den Neustart nicht fortgesetzt" >&2
  exit 1
fi

# Dieselben drei Aussagen wie nach dem Live-Reload, jetzt über den
# Neustart hinweg: der ausgeschlossene Schlüssel fehlt, die nicht
# ausgeschlossene Spalte bleibt (Kontrolle gegen ein leeres Row Image), und
# der ausgeschlossene Wert steht im ganzen persistierten Change nicht
# (LH-QA-SEC-004).
column_restart_key=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT jsonb_exists(new_data, '$COLUMN_NAME') FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$COLUMN_TABLE' AND new_data->>'id' = '3'")
column_restart_name=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT new_data->>'name' FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$COLUMN_TABLE' AND new_data->>'id' = '3'")
column_restart_leak=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT (coalesce(new_data::text, '') LIKE '%$COLUMN_VALUE_RESTART%') OR (coalesce(old_data::text, '') LIKE '%$COLUMN_VALUE_RESTART%') FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$COLUMN_TABLE' AND new_data->>'id' = '3'")
if [ "$column_restart_key" != "f" ]; then
  echo "run-integration-tests: Spaltenausschluss-Neustart-Beleg — die nach dem Neustart erfasste Change (id=3, $COLUMN_TABLE) trägt den ausgeschlossenen Spaltenschlüssel $COLUMN_NAME weiterhin im Row Image (der Ausschlussstand wurde beim Prozessstart nicht wiederhergestellt)" >&2
  exit 1
fi
if [ "$column_restart_name" != "ColumnAfterRestart" ]; then
  echo "run-integration-tests: Spaltenausschluss-Neustart-Beleg — die nach dem Neustart erfasste Change (id=3, $COLUMN_TABLE) trägt die nicht ausgeschlossene Spalte name='${column_restart_name:-leer}', wollen 'ColumnAfterRestart'" >&2
  exit 1
fi
if [ "$column_restart_leak" != "f" ]; then
  echo "run-integration-tests: Spaltenausschluss-Neustart-Beleg — der ausgeschlossene Wert '$COLUMN_VALUE_RESTART' steht im nach dem Neustart persistierten Change (id=3, $COLUMN_TABLE) (LH-QA-SEC-004)" >&2
  exit 1
fi

echo "run-integration-tests: Spaltenausschluss-Neustart-Beleg (ADR-0065, LH-QA-SEC-004) — der reale Container-Neustart ließ den dauerhaften Ausschlussstand wirksam werden: die danach erfasste Change (id=3) trägt $COLUMN_NAME nicht im Row Image und den Wert '$COLUMN_VALUE_RESTART' nirgends"

abdeckung_declare "Spaltenausschluss-Negative-Beleg" "LH-FA-CFG-005" "der Ausschluss einer an der Quelle fehlenden Spalte endet real failed samt Fehlertext, nicht still applied" "Spaltenausschluss-Negative-Beleg (LH-FA-CFG-005) — cdc.exclude_column"

# Negative-Beleg (LH-FA-CFG-005 Negative): der Ausschluss einer an der Quelle
# nicht existierenden Spalte endet real `failed` samt Fehlertext — nicht
# still `applied` (ErrSourceColumnMissing-Pfad).
column_missing_request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT cdc.exclude_column('src-e2e', 'public', '$COLUMN_TABLE', 'nicht_vorhandene_spalte')")
if [ -z "$column_missing_request_id" ]; then
  echo "run-integration-tests: cdc.exclude_column gegen die nicht existierende Spalte lieferte keine Antrags-ID" >&2
  exit 1
fi

column_missing_status=""
column_missing_failed=0
column_missing_error=""
for _ in $(seq 1 60); do
  column_missing_status=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$column_missing_request_id'")
  if [ "$column_missing_status" = "failed" ]; then
    column_missing_error=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT error_message FROM cdc.administration_request WHERE administration_request_id = '$column_missing_request_id'")
    column_missing_failed=1
    break
  fi
  if [ "$column_missing_status" = "applied" ]; then
    break
  fi
  sleep 0.5
done
if [ "$column_missing_failed" -ne 1 ]; then
  echo "run-integration-tests: cdc.exclude_column gegen die nicht existierende Spalte endete mit status=${column_missing_status:-leer}, wollen failed (LH-FA-CFG-005 Negative: expliziter Fehlerpfad)" >&2
  exit 1
fi
if ! printf '%s' "$column_missing_error" | grep -qF "Spalte existiert nicht an der Quelle"; then
  echo "run-integration-tests: cdc.exclude_column-Antrag $column_missing_request_id endete failed, trägt aber nicht den Fehlertext 'Spalte existiert nicht an der Quelle': $column_missing_error" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem Spaltenausschluss-Negative-Beleg nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: Spaltenausschluss-Negative-Beleg (LH-FA-CFG-005) — cdc.exclude_column($COLUMN_TABLE.nicht_vorhandene_spalte) endete real failed mit Fehlertext, der Feed-Container lief unverändert weiter"

abdeckung_declare "Lasttest-Beleg cdc_capture_lag" "LH-FA-ADM-004" "cdc_capture_lag bildet den Abstand zwischen Quelländerung und CDC-Verfügbarkeit ab: eine durch eine pausierte CDC-Runtime künstlich verzögerte Transaktion liegt über der ungehinderten" "Lasttest-Beleg cdc_capture_lag — Baseline"

# Lasttest-Beleg (LH-FA-ADM-004, SPEC-013 CDC_LAG_THRESHOLDS): cdc_capture_lag
# bildet den Abstand zwischen Quelländerung und CDC-Verfügbarkeit ab. Zwei
# Quelltransaktionen auf derselben Tabelle zeigen den Unterschied: eine
# ungehinderte Transaktion liefert einen kleinen Wert; eine Transaktion,
# deren Erfassung durch eine pausierte CDC-Runtime künstlich verzögert wird
# (`docker pause` hält den Feed-Container über die Freezer-Cgroup an, bevor
# die Transaktion verarbeitet ist), liefert einen Wert nahe der
# Pausendauer. `feed_e2e_full` bleibt über den ganzen Lauf aktiviert (anders
# als `feed_e2e_flow`, das der letzte E2E-Testfall deaktiviert). Die
# Pausendauer bleibt unter `wal_sender_timeout=2000` (compose.yaml): eine
# längere Pause ließe die Quelle die Replication-Verbindung selbst beenden,
# bevor der Feed-Container sie fortsetzen kann — der Lasttest misst die
# Erfassungsverzögerung, keine Verbindungsstörung.
LAG_TABLE=feed_e2e_full
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
# Referenzzeilen in test/integration/integration_test.go (id=1, id=7 auf
# derselben Tabelle feed_e2e_full) — keine Kollision zwischen den beiden
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

abdeckung_declare "Black-Box-CLI-Rundlauf" "LH-QA-POR-003" "register-consumer und acknowledge-consumer laufen ausschließlich als externe docker exec-Aufrufe gegen den Produktions-Binary, über einen simulierten Container-Neustart hinweg" "Black-Box-CLI-Rundlauf belegt — register-consumer/acknowledge-consumer extern"

# Black-Box-CLI-Rundlauf (LH-QA-POR-003, ADR-0030 E2E-Tier): anders als der
# Go-Testlauf oben (`go test ./test/integration/...`, der intern gegen
# `postgresstorage`/`bootstrap` läuft) ruft dieser Abschnitt
# `register-consumer`/`acknowledge-consumer` ausschließlich als externen
# Prozess gegen den laufenden, containerisierten Produktions-Binary auf
# (`docker exec … /pg-change-feed …`, distroless — kein Shell im
# Feed-Container, daher kein `sh -c`-Umweg nötig). Lesen bleibt über den
# bestehenden externen Lesezugriffsweg `cdc.changes` (LH-FA-SST-002) — die
# CLI trägt keinen Lese-Unterbefehl. `feed_e2e_full` trägt denselben Grund
# wie beim Lasttest-Beleg oben: die Tabelle bleibt über den ganzen Lauf
# aktiviert. Die IDs 95/96 liegen in einem eigenen Wertebereich, getrennt
# von den E2E-Referenzzeilen (id=1) und dem Lasttest-Beleg (id=90/91) auf
# derselben Tabelle.
exec_feed() {
  docker exec "$FEED_CONTAINER" /pg-change-feed "$@"
}

abdeckung_declare "Retention-Lebenszyklus-Rundlauf (Blocker-Sichtbarkeit)" "LH-FA-RET-005" "cdc.retention_blockers zeigt real den Consumer mit der am weitesten zurückliegenden bestätigten Position als aktuellen Blocker der Quelle — und nach der Bestätigung über die Position hinweg keinen mehr" "Retention-Lebenszyklus-Rundlauf — cdc.retention_blockers zeigt real"

abdeckung_declare "Retention-Lebenszyklus-Rundlauf (kombiniert)" "LH-FA-RET-002,LH-FA-RET-003,LH-FA-RET-004,LH-FA-RET-005,LH-FA-RET-006" "eine isolierte Zeile und ein eigener Consumer durchlaufen real Blocker-Sichtbarkeit, Bestätigung über die Position hinweg, Löschung sowie die durchgehende numerische Lesbarkeit der Speichergröße" "Retention-Lebenszyklus-Rundlauf (kombiniert) — 'RetentionLifecycle'"

# Retention-Lebenszyklus-Rundlauf (kombiniert, LH-FA-RET-002…006): eine
# eigene, isolierte Zeile (id=210 auf feed_e2e_full) und ein eigener, neu
# registrierter Consumer durchlaufen real die vollständige Kette in einer
# Kette, statt sie wie in den folgenden Abschnitten über mehrere getrennte
# Belege zu prüfen — Blocker-Sichtbarkeit, Bestätigung über die Position
# hinweg, die reale Löschung (während die bestätigte Position noch real
# vorhanden ist), danach die reale Abwesenheit jeder Zeile für `src-e2e`,
# `cdc_storage_bytes` durchgehend numerisch. Läuft an dieser
# Stelle, weil hier noch kein über register-consumer/acknowledge-consumer
# geführter Consumer gegen `src-e2e` bestätigt hat (dieselbe Ausgangslage,
# die der folgende Abschnitt „Zustand 1" voraussetzt) — der hier
# registrierte Consumer wird am Ende real wieder entfernt, damit diese
# Ausgangslage für den folgenden Abschnitt unverändert gilt.
#
# `cdc.retention_blockers` trägt je Quelle höchstens eine Zeile
# (`DISTINCT ON`) und tut das unabhängig vom Rückstandswert, solange
# irgendein Consumer eine bestätigte Position gegen die Quelle trägt
# (`INNER JOIN cdc.consumer_position`) — „kein Blocker" ist deshalb
# ausschließlich die Abwesenheit jeder bestätigten Position, nicht ein
# Rückstand von 0. Die reale Abwesenheit nach der Bestätigung über die
# Position hinweg verlangt deshalb die reale Entfernung dieser Position,
# nicht nur ihr Fortschreiten — dieselbe Abwesenheits-Lesart, die die View
# für einen nie bestätigenden Consumer bereits trägt.
LIFECYCLE_CONSUMER=cli-e2e-lifecycle-consumer
LIFECYCLE_TABLE=feed_e2e_full

storage_bytes_before=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_storage_bytes'")
if ! printf '%s' "$storage_bytes_before" | grep -qE '^[0-9]+$'; then
  echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — cdc_storage_bytes vor der Kette nicht real numerisch abfragbar: '$storage_bytes_before'" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$LIFECYCLE_TABLE (id, name) VALUES (210, 'RetentionLifecycle');
SQL

lifecycle_position=""
for _ in $(seq 1 120); do
  lifecycle_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$LIFECYCLE_TABLE' AND new_data->>'id' = '210'")
  if [ -n "$lifecycle_position" ]; then
    break
  fi
  sleep 0.25
done
if [ -z "$lifecycle_position" ]; then
  echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — 'RetentionLifecycle' (id=210) wurde nicht innerhalb der Zeitspanne über cdc.changes gelesen" >&2
  exit 1
fi

docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -c \
  "UPDATE cdc.transaction SET committed_at = current_timestamp - interval '25 hours' WHERE transaction_id = (SELECT transaction_id FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$LIFECYCLE_TABLE' AND new_data->>'id' = '210')" >/dev/null

earliest_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT min(commit_position) FROM cdc.transaction WHERE source_id = 'src-e2e'")
if [ -z "$earliest_position" ] || [ "$earliest_position" -ge "$lifecycle_position" ]; then
  echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — keine reale, frühere Position als $lifecycle_position gefunden (min=${earliest_position:-leer})" >&2
  exit 1
fi

if ! exec_feed register-consumer "$LIFECYCLE_CONSUMER"; then
  echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — register-consumer ($LIFECYCLE_CONSUMER) endete mit einem Fehler" >&2
  exit 1
fi

if ! exec_feed acknowledge-consumer "$LIFECYCLE_CONSUMER" "$earliest_position"; then
  echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — erste Bestätigung ($LIFECYCLE_CONSUMER, Position $earliest_position) endete mit einem Fehler" >&2
  exit 1
fi

blocker_before=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT consumer_id FROM cdc.retention_blockers WHERE source_id = 'src-e2e'")
if [ "$blocker_before" != "$LIFECYCLE_CONSUMER" ]; then
  echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — cdc.retention_blockers zeigt $LIFECYCLE_CONSUMER nicht als aktuellen Blocker für src-e2e, sondern '${blocker_before:-leer}'" >&2
  exit 1
fi

echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — cdc.retention_blockers zeigt real $LIFECYCLE_CONSUMER als aktuellen Blocker für src-e2e (LH-FA-RET-005)"

# Erste Phase (Consumer-Block, LH-FA-RET-004): eine großzügige, feste
# Wartezeit über mehr als zwei Lösch-Takte (retentionInterval=10s) hinweg
# belegt die Abwesenheit einer Löschung über mehrere reale Takte, nicht nur
# einen einzelnen zu frühen Blick.
sleep 25
lifecycle_present_blocked=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$LIFECYCLE_TABLE' AND new_data->>'id' = '210'")
if [ "$lifecycle_present_blocked" != "1" ]; then
  echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — 'RetentionLifecycle' (id=210, bereits alt genug) wurde entfernt, obwohl $LIFECYCLE_CONSUMER seine Position noch nicht bestätigt hatte (LH-FA-RET-004 Consumer-Block)" >&2
  exit 1
fi

if ! exec_feed acknowledge-consumer "$LIFECYCLE_CONSUMER" "$lifecycle_position"; then
  echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — zweite Bestätigung ($LIFECYCLE_CONSUMER, über die Position von id=210 hinweg, $lifecycle_position) endete mit einem Fehler" >&2
  exit 1
fi

# Die bestätigte Position von $LIFECYCLE_CONSUMER bleibt an dieser Stelle
# real in cdc.consumer_position bestehen (kein DELETE vorher) — der
# folgende Poll belegt dadurch real, dass die Bestätigung über die Position
# hinweg allein die Löschung freigibt, unabhängig von einer späteren
# Entfernung der Position selbst (LH-FA-RET-003/004).
lifecycle_deleted=0
for _ in $(seq 1 60); do
  remaining=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$LIFECYCLE_TABLE' AND new_data->>'id' = '210'")
  if [ "$remaining" = "0" ]; then
    lifecycle_deleted=1
    break
  fi
  sleep 1
done
if [ "$lifecycle_deleted" -ne 1 ]; then
  echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — 'RetentionLifecycle' (id=210) wurde nach Freigabe nicht innerhalb der Zeitspanne real entfernt" >&2
  exit 1
fi

echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — 'RetentionLifecycle' (id=210) real entfernt, während die bestätigte Position von $LIFECYCLE_CONSUMER noch real in cdc.consumer_position vorhanden ist (LH-FA-RET-003/004)"

# Die Zeile ist an dieser Stelle bereits real gelöscht; die Entfernung der
# Position dient ausschließlich der Sichtbarkeits-Prüfung von
# cdc.retention_blockers (DoD-Punkt 1) und der Wiederherstellung der
# Ausgangslage für den nachfolgenden „Zustand 1"-Abschnitt — sie löst keine
# Löschung mehr aus, die nicht bereits erfolgt wäre.
docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -c \
  "DELETE FROM cdc.consumer_position WHERE consumer_id = '$LIFECYCLE_CONSUMER'" >/dev/null

blocker_after=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.retention_blockers WHERE source_id = 'src-e2e'")
if [ "$blocker_after" != "0" ]; then
  echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — cdc.retention_blockers trägt nach der Bestätigung über die Position hinweg noch eine Zeile für src-e2e ($blocker_after), erwartet leer" >&2
  exit 1
fi

echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — cdc.retention_blockers zeigt real keinen Blocker mehr für src-e2e, auch nach der bereits erfolgten realen Löschung (LH-FA-RET-005)"

storage_bytes_after=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_storage_bytes'")
if ! printf '%s' "$storage_bytes_after" | grep -qE '^[0-9]+$'; then
  echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf — cdc_storage_bytes nach der Löschung nicht real numerisch abfragbar: '$storage_bytes_after'" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem Retention-Lebenszyklus-Rundlauf nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: Retention-Lebenszyklus-Rundlauf (kombiniert) — 'RetentionLifecycle' (id=210) real entfernt, cdc_storage_bytes durchgehend numerisch abfragbar (vorher $storage_bytes_before Bytes, nachher $storage_bytes_after Bytes; LH-FA-RET-002…006)"

abdeckung_declare "Retention-Sichtbarkeits-Beleg (kein Blocker)" "LH-FA-RET-005,LH-FA-RET-006" "die diagnose-Ausgabe zeigt vor jeder Consumer-Bestätigung 'kein Blocker' und eine numerische Speichergröße" "Retention-Sichtbarkeits-Beleg (CLI, Zustand 1)"

# Retention-Sichtbarkeits-Beleg (CLI), Zustand 1 — kein Blocker
# (LH-FA-SST-003, deckt LH-FA-RET-005/006): an dieser Stelle hat noch kein
# über register-consumer/acknowledge-consumer geführter Consumer gegen
# `src-e2e` bestätigt — die beiden Consumer, die
# TestE2ERetentionBlockersViewShowsFurthestBehindConsumer weiter oben direkt
# über den ConsumerStatePort-Adapter registrierte, sind bereits per
# `t.Cleanup` entfernt (siehe deren Funktionskommentar in
# test/integration/integration_test.go), und der oben real durchlaufene
# Retention-Lebenszyklus-Consumer ist am Ende dieses Abschnitts ebenfalls
# real entfernt. `cdc.retention_blockers` trägt deshalb real keine Zeile für
# `src-e2e`, und die `diagnose`-Ausgabe muss das als „kein Blocker" zeigen,
# nicht als Fehlerzustand.
set +e
diagnose_noblocker_output=$(exec_feed diagnose)
diagnose_noblocker_status=$?
set -e
if [ "$diagnose_noblocker_status" -ne 0 ]; then
  echo "run-integration-tests: diagnose (Retention-Beleg, kein Blocker, docker exec) endete mit Ausgang $diagnose_noblocker_status: $diagnose_noblocker_output" >&2
  exit 1
fi
if ! printf '%s' "$diagnose_noblocker_output" | grep -qF "Blockierender Consumer (LH-FA-RET-005): kein Blocker"; then
  echo "run-integration-tests: diagnose-Ausgabe zeigt vor jeder Consumer-Bestätigung nicht 'kein Blocker' (LH-FA-RET-005): $diagnose_noblocker_output" >&2
  exit 1
fi
if ! printf '%s' "$diagnose_noblocker_output" | grep -qE "Speicherverbrauch cdc_storage_bytes \(LH-FA-RET-006\): [0-9]+ Bytes"; then
  echo "run-integration-tests: diagnose-Ausgabe trägt keine numerische cdc_storage_bytes-Zeile (LH-FA-RET-006): $diagnose_noblocker_output" >&2
  exit 1
fi

echo "run-integration-tests: Retention-Sichtbarkeits-Beleg (CLI, Zustand 1) — 'kein Blocker' vor jeder Consumer-Bestätigung, cdc_storage_bytes numerisch sichtbar (LH-FA-RET-005/006)"

CLI_CONSUMER=cli-e2e-consumer
CLI_TABLE=feed_e2e_full

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
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$CLI_TABLE' AND new_data->>'id' = '95'")
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
  "SELECT acknowledged_position FROM cdc.consumer_position WHERE consumer_id = '$CLI_CONSUMER' AND source_id = 'src-e2e'")
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
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$CLI_TABLE' AND new_data->>'id' = '96'")
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
  "SELECT acknowledged_position FROM cdc.consumer_position WHERE consumer_id = '$CLI_CONSUMER' AND source_id = 'src-e2e'")
if [ "$restored_position" != "$first_position" ]; then
  echo "run-integration-tests: aus cdc.consumer_position zurückgelesene Position nach dem Neustart = $restored_position, wollen $first_position (die zuvor bestätigte)" >&2
  exit 1
fi

resumed_ids=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT string_agg(new_data->>'id', ',' ORDER BY commit_position) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$CLI_TABLE' AND commit_position > $restored_position")
if [ "$resumed_ids" != "96" ]; then
  echo "run-integration-tests: Fortsetzen ab der bestätigten Position (Neustart-Beleg) — ab commit_position=$restored_position gelesene id-Folge '$resumed_ids', wollen '96' (id=95 bleibt hinter der bestätigten Position, keine Wiederholung)" >&2
  exit 1
fi

if ! exec_feed acknowledge-consumer "$CLI_CONSUMER" "$second_position"; then
  echo "run-integration-tests: zweite externe Bestätigung (acknowledge-consumer, docker exec) endete mit einem Fehler (Position $second_position)" >&2
  exit 1
fi

acked_second=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT acknowledged_position FROM cdc.consumer_position WHERE consumer_id = '$CLI_CONSUMER' AND source_id = 'src-e2e'")
if [ "$acked_second" != "$second_position" ]; then
  echo "run-integration-tests: cdc.consumer_position trägt $acked_second, wollen $second_position (zweite externe Bestätigung, nach dem simulierten Neustart)" >&2
  exit 1
fi

echo "run-integration-tests: Black-Box-CLI-Rundlauf belegt — register-consumer/acknowledge-consumer extern (docker exec), Fortsetzen nach simuliertem Neustart ab Position $first_position, Endposition $second_position"

abdeckung_declare "Verarbeitungsrückstand-Beleg" "LH-FA-ADM-005" "cdc.consumer_status zeigt den realen Rückstand über eine reine SQL-Lesung und nach der zweiten Bestätigung null" "Verarbeitungsrückstand-Beleg cdc.consumer_status"

# Verarbeitungsrückstand-Beleg (LH-FA-ADM-005, cdc.consumer_status): ein
# über register-consumer/acknowledge-consumer geführter Consumer (extern,
# docker exec — dasselbe Muster wie der Black-Box-CLI-Rundlauf oben)
# bestätigt zunächst eine frühe, real erfasste Position; das bindet
# consumer_status.source_id an diese Quelle (die Sicht liest
# latest_commit_position über eine auf cp.source_id korrelierte
# Unterabfrage — ohne jede vorherige Bestätigung liefert sie dafür NULL,
# keinen Rückstand). Eine danach real erfasste weitere Änderung hebt
# latest_commit_position über die bestätigte Position an: der Rückstand
# wird über eine reine SQL-Lesung von cdc.consumer_status sichtbar. Die
# zweite externe Bestätigung auf dieselbe (jetzt aktuelle) Position senkt
# ihn auf 0. Die IDs 120/121 liegen in einem eigenen Wertebereich,
# getrennt von den übrigen Testfall-Gruppen auf derselben Tabelle (id=1,
# id=90/91, id=95/96 oben).
BACKLOG_CONSUMER=cli-e2e-backlog-consumer
BACKLOG_TABLE=feed_e2e_full

if ! exec_feed register-consumer "$BACKLOG_CONSUMER"; then
  echo "run-integration-tests: register-consumer (Rückstands-Beleg, docker exec) endete mit einem Fehler" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$BACKLOG_TABLE (id, name) VALUES (120, 'BacklogBaseline');
SQL

baseline_backlog_position=""
for _ in $(seq 1 120); do
  baseline_backlog_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BACKLOG_TABLE' AND new_data->>'id' = '120'")
  if [ -n "$baseline_backlog_position" ]; then
    break
  fi
  sleep 0.25
done
if [ -z "$baseline_backlog_position" ]; then
  echo "run-integration-tests: Baseline-Änderung des Rückstands-Belegs (id=120) wurde nicht innerhalb der Zeitspanne über cdc.changes gelesen" >&2
  exit 1
fi

if ! exec_feed acknowledge-consumer "$BACKLOG_CONSUMER" "$baseline_backlog_position"; then
  echo "run-integration-tests: erste externe Bestätigung des Rückstands-Belegs (acknowledge-consumer, docker exec) endete mit einem Fehler (Position $baseline_backlog_position)" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$BACKLOG_TABLE (id, name) VALUES (121, 'BacklogAhead');
SQL

latest_backlog_position=""
for _ in $(seq 1 120); do
  latest_backlog_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BACKLOG_TABLE' AND new_data->>'id' = '121'")
  if [ -n "$latest_backlog_position" ]; then
    break
  fi
  sleep 0.25
done
if [ -z "$latest_backlog_position" ]; then
  echo "run-integration-tests: nachfolgende Änderung des Rückstands-Belegs (id=121) wurde nicht innerhalb der Zeitspanne über cdc.changes gelesen" >&2
  exit 1
fi

backlog_before=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT latest_commit_position - acknowledged_position FROM cdc.consumer_status WHERE consumer_id = '$BACKLOG_CONSUMER' AND source_id = 'src-e2e'")
if [ -z "$backlog_before" ] || [ "$backlog_before" -le 0 ]; then
  echo "run-integration-tests: cdc.consumer_status-Rückstand vor der zweiten Bestätigung: '${backlog_before:-leer}' (Erwartung: > 0, acknowledged_position=$baseline_backlog_position, danach erfasste Position=$latest_backlog_position)" >&2
  exit 1
fi

if ! exec_feed acknowledge-consumer "$BACKLOG_CONSUMER" "$latest_backlog_position"; then
  echo "run-integration-tests: zweite externe Bestätigung des Rückstands-Belegs (acknowledge-consumer, docker exec) endete mit einem Fehler (Position $latest_backlog_position)" >&2
  exit 1
fi

backlog_after=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT latest_commit_position - acknowledged_position FROM cdc.consumer_status WHERE consumer_id = '$BACKLOG_CONSUMER' AND source_id = 'src-e2e'")
if [ "$backlog_after" != "0" ]; then
  echo "run-integration-tests: cdc.consumer_status-Rückstand nach der zweiten Bestätigung: $backlog_after, wollen 0 (bestätigte Position=$latest_backlog_position)" >&2
  exit 1
fi

echo "run-integration-tests: Verarbeitungsrückstand-Beleg cdc.consumer_status — Rückstand vor der zweiten Bestätigung $backlog_before, danach $backlog_after (LH-FA-ADM-005)"

abdeckung_declare "CLI-Diagnose-Beleg (Normalbetrieb)" "LH-FA-SST-003,LH-FA-ADM-002,LH-FA-ADM-003,LH-FA-ADM-004,LH-FA-ADM-005" "der diagnose-Sondermodus trägt im Normalbetrieb alle vier Signale: Betriebsstatus, fehlenden Fehlerzustand, numerischen CDC-Abstand und die Rückstände der geführten Consumer" "CLI-Diagnose-Beleg (Normalbetrieb) —"

# CLI-Diagnose-Beleg (LH-FA-SST-003, deckt LH-FA-ADM-002…005, slice-038):
# der neue `diagnose`-Sondermodus liest denselben SQL-Zugriffsweg wie
# `--healthcheck` oben (`cfg.ReaderDSN`) plus `cdc.metrics` und gibt alle
# vier Signale menschenlesbar aus — extern per `docker exec` gegen den
# laufenden Feed-Container (`exec_feed`, dasselbe Muster wie der
# Black-Box-CLI-Rundlauf oben). Normalbetrieb zuerst: `cdc.heartbeat` trägt
# an dieser Stelle bereits ein aktuelles Lebenszeichen ohne Fehlerzustand
# (Compose-Healthcheck-Vertrag oben) und `cdc.metrics` einen gemessenen
# `cdc_capture_lag` sowie die beiden oben geführten Consumer. `latest_commit_position`
# (`cdc.consumer_status`) ist quellenweit, nicht consumer- oder
# tabellenspezifisch — CLI_CONSUMER trägt hier real einen Rückstand
# ungleich 0 (die BACKLOG_CONSUMER-Belegschreibungen liefen auf derselben
# Tabelle nach seiner letzten Bestätigung); nur BACKLOG_CONSUMER, dessen
# zweite Bestätigung unmittelbar davor lief, ist an dieser Stelle
# verlässlich 0. Retention-Sichtbarkeits-Beleg (CLI), Zustand 2 — realer
# Blocker (LH-FA-RET-005/006): `cdc.retention_blockers` wählt je Quelle den
# Consumer mit der kleinsten bestätigten Position — das ist an dieser
# Stelle real CLI_CONSUMER (Rückstand ungleich 0, siehe oben), nicht
# BACKLOG_CONSUMER (Rückstand 0). `cdc_storage_bytes` trägt zu diesem
# Zeitpunkt bereits einen positiven Wert aus der bisherigen CDC-Erfassung.
set +e
diagnose_output=$(exec_feed diagnose)
diagnose_status=$?
set -e
if [ "$diagnose_status" -ne 0 ]; then
  echo "run-integration-tests: diagnose (Normalbetrieb, docker exec) endete mit Ausgang $diagnose_status: $diagnose_output" >&2
  exit 1
fi
for expected in \
  "Betriebsstatus (LH-FA-ADM-002):" \
  "Fehlerzustand (LH-FA-ADM-003): keiner (Normalbetrieb)" \
  "CDC-Abstand cdc_capture_lag (LH-FA-ADM-004):" \
  "$BACKLOG_CONSUMER: 0"
do
  if ! printf '%s' "$diagnose_output" | grep -qF "$expected"; then
    echo "run-integration-tests: diagnose-Ausgabe (Normalbetrieb) trägt nicht den erwarteten Text '$expected': $diagnose_output" >&2
    exit 1
  fi
done
if ! printf '%s' "$diagnose_output" | grep -qE "$CLI_CONSUMER: [0-9]+"; then
  echo "run-integration-tests: diagnose-Ausgabe (Normalbetrieb) trägt keine numerische Rückstands-Zeile für $CLI_CONSUMER: $diagnose_output" >&2
  exit 1
fi
if ! printf '%s' "$diagnose_output" | grep -qE "Blockierender Consumer \(LH-FA-RET-005\): .*\($CLI_CONSUMER\), bestätigte Position [0-9]+, Rückstand [0-9]+"; then
  echo "run-integration-tests: diagnose-Ausgabe (Normalbetrieb) zeigt $CLI_CONSUMER nicht als real blockierenden Consumer (LH-FA-RET-005): $diagnose_output" >&2
  exit 1
fi
if ! printf '%s' "$diagnose_output" | grep -qE "Speicherverbrauch cdc_storage_bytes \(LH-FA-RET-006\): [0-9]+ Bytes"; then
  echo "run-integration-tests: diagnose-Ausgabe (Normalbetrieb) trägt keine numerische cdc_storage_bytes-Zeile (LH-FA-RET-006): $diagnose_output" >&2
  exit 1
fi

echo "run-integration-tests: CLI-Diagnose-Beleg (Normalbetrieb) — alle vier Signale (LH-FA-ADM-002…005) sowie der reale Blocker $CLI_CONSUMER und cdc_storage_bytes (LH-FA-RET-005/006) in der diagnose-Ausgabe sichtbar"

abdeckung_declare "CLI-Diagnose-Beleg (Fehlerzustand)" "LH-FA-ADM-003" "ein direkt in die Heartbeat-Projektion geschriebener Fehlerzustand ist in der diagnose-Ausgabe von Normalbetrieb unterscheidbar, ohne den laufenden Feed-Container zu beenden" "CLI-Diagnose-Beleg (Fehlerzustand) —"

# Fehlerzustand-Beleg (LH-FA-ADM-003 Boundary: „erkennbar von Normalbetrieb
# unterscheidbar"). Ein real vom Erfassungspfad ausgelöster Fehlerzustand
# beendet den Feed-Container-Prozess dauerhaft (`reportFault` läuft nur
# unmittelbar vor `os.Exit`, siehe `internal/bootstrap/wiring.go`
# Funktionskommentar) — ein `docker exec` gegen einen bereits beendeten,
# `restart: "no"`-Container ist danach nicht mehr möglich (slice-038 §6
# Risiko 1). Dieser Beleg schreibt denselben Spaltenwert
# (`cdc.process_heartbeat.error_class`), den der reale Fehlerpfad schriebe,
# direkt über SQL — derselbe Lesepfad (View -> diagnose-Ausgabe), ohne den
# laufenden Container zu beenden. Der periodische Heartbeat-Takt (5s,
# `heartbeatInterval`) überschreibt `error_class` beim nächsten
# erfolgreichen Beat wieder auf NULL, unabhängig von unserem
# SQL-Schreibzug — die Schleife schreibt den Fehlerzustand deshalb
# wiederholt, bis ein `diagnose`-Aufruf ihn innerhalb desselben kurzen
# Fensters liest (bei 20 Versuchen weit innerhalb eines einzigen 5s-Takts).
error_state_seen=0
diagnose_error_output=""
diagnose_error_status=1
for _ in $(seq 1 20); do
  docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -c \
    "UPDATE cdc.process_heartbeat SET heartbeat_at = current_timestamp, error_class = 'schema' WHERE source_id = 'src-e2e'" >/dev/null
  set +e
  diagnose_error_output=$(exec_feed diagnose)
  diagnose_error_status=$?
  set -e
  if [ "$diagnose_error_status" -eq 0 ] && printf '%s' "$diagnose_error_output" | grep -qF "Fehlerzustand (LH-FA-ADM-003): schema"; then
    error_state_seen=1
    break
  fi
done
if [ "$error_state_seen" -ne 1 ]; then
  echo "run-integration-tests: diagnose-Ausgabe zeigt den gesetzten Fehlerzustand 'schema' nicht innerhalb von 20 Versuchen (Boundary LH-FA-ADM-003, letzter Ausgang $diagnose_error_status): $diagnose_error_output" >&2
  exit 1
fi
if printf '%s' "$diagnose_error_output" | grep -qF "keiner (Normalbetrieb)"; then
  echo "run-integration-tests: diagnose-Ausgabe trägt fälschlich die Normalbetrieb-Zeile trotz gesetztem Fehlerzustand: $diagnose_error_output" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem diagnose-Fehlerzustand-Beleg nicht mehr weiter (der Fehlerzustand wurde direkt in cdc.process_heartbeat geschrieben, nicht vom Erfassungspfad ausgelöst — kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: CLI-Diagnose-Beleg (Fehlerzustand) — 'schema' sichtbar und von Normalbetrieb unterscheidbar (LH-FA-ADM-003 Boundary), Feed-Container läuft unverändert weiter"

abdeckung_declare "SQL-Administration Live-Reload (enable)" "LH-FA-ADM-001,LH-FA-CFG-001" "eine bewusst nicht in der Bindungsliste geführte Tabelle wird über den SQL-Antrag aktiviert und vom bereits laufenden Feed-Container ohne Neustart erfasst" "SQL-Administration Live-Reload-Beleg (enable)"

abdeckung_declare "SQL-Administration Live-Reload (disable)" "LH-FA-CFG-002" "derselbe Antrags-Weg spiegelbildlich: die Deaktivierung beendet die Erfassung dieser Tabelle, der Prozess läuft weiter" "SQL-Administration Live-Reload-Beleg (disable)"

# SQL-Administration Live-Reload-Beleg (ADR-0050, LH-FA-ADM-001,
# LH-FA-CFG-001/002): `feed_e2e_sql_admin` ist bewusst NICHT Teil von
# CDC_TABLES (compose.yaml) — ihre Aktivierung/Deaktivierung läuft
# ausschließlich über die Antrags-Queue (`cdc.enable_table`/
# `cdc.disable_table`), verarbeitet von der Administrations-Goroutine des
# bereits laufenden Feed-Containers (kein Neustart, kein `docker restart`
# wie beim Black-Box-CLI-Rundlauf oben). `SELECT cdc.enable_table(...)`
# bestätigt nur „beantragt" — der Poll unten wartet auf
# `status = 'applied'`, bevor die reale Erfassungswirkung geprüft wird.
ADMIN_TABLE=feed_e2e_sql_admin

enable_request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT cdc.enable_table('src-e2e', 'public', '$ADMIN_TABLE')")
if [ -z "$enable_request_id" ]; then
  echo "run-integration-tests: cdc.enable_table($ADMIN_TABLE) lieferte keine Antrags-ID" >&2
  exit 1
fi

enable_status=""
enable_applied=0
for _ in $(seq 1 60); do
  enable_status=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$enable_request_id'")
  if [ "$enable_status" = "applied" ]; then
    enable_applied=1
    break
  fi
  if [ "$enable_status" = "failed" ]; then
    error_message=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT error_message FROM cdc.administration_request WHERE administration_request_id = '$enable_request_id'")
    echo "run-integration-tests: cdc.enable_table($ADMIN_TABLE)-Antrag $enable_request_id scheiterte: $error_message" >&2
    exit 1
  fi
  sleep 0.5
done
if [ "$enable_applied" -ne 1 ]; then
  echo "run-integration-tests: cdc.enable_table($ADMIN_TABLE)-Antrag $enable_request_id wurde nicht innerhalb der Zeitspanne von der Administrations-Goroutine verarbeitet (status=${enable_status:-leer})" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$ADMIN_TABLE (id, name) VALUES (1, 'SqlAdminEnabled');
SQL

admin_enabled_captured=0
for _ in $(seq 1 120); do
  found=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$ADMIN_TABLE' AND new_data->>'id' = '1'")
  if [ "$found" = "1" ]; then
    admin_enabled_captured=1
    break
  fi
  sleep 0.25
done
if [ "$admin_enabled_captured" -ne 1 ]; then
  echo "run-integration-tests: über SQL aktivierte Tabelle $ADMIN_TABLE — Änderung (id=1) wurde nicht vom laufenden Feed-Container erfasst (kein Neustart)" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach der SQL-Aktivierung nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: SQL-Administration Live-Reload-Beleg (enable) — cdc.enable_table($ADMIN_TABLE) ohne Neustart verarbeitet, Änderung id=1 real erfasst"

# Deaktivierung: derselbe Antrags-Weg, spiegelbildlich. Der Feed-Container
# bleibt danach am Leben — nur die Erfassung dieser einen Tabelle endet
# (`LH-FA-CFG-002`), der Prozess selbst nicht.
disable_request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT cdc.disable_table('src-e2e', 'public', '$ADMIN_TABLE')")
if [ -z "$disable_request_id" ]; then
  echo "run-integration-tests: cdc.disable_table($ADMIN_TABLE) lieferte keine Antrags-ID" >&2
  exit 1
fi

disable_status=""
disable_applied=0
for _ in $(seq 1 60); do
  disable_status=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$disable_request_id'")
  if [ "$disable_status" = "applied" ]; then
    disable_applied=1
    break
  fi
  if [ "$disable_status" = "failed" ]; then
    error_message=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT error_message FROM cdc.administration_request WHERE administration_request_id = '$disable_request_id'")
    echo "run-integration-tests: cdc.disable_table($ADMIN_TABLE)-Antrag $disable_request_id scheiterte: $error_message" >&2
    exit 1
  fi
  sleep 0.5
done
if [ "$disable_applied" -ne 1 ]; then
  echo "run-integration-tests: cdc.disable_table($ADMIN_TABLE)-Antrag $disable_request_id wurde nicht innerhalb der Zeitspanne von der Administrations-Goroutine verarbeitet (status=${disable_status:-leer})" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$ADMIN_TABLE (id, name) VALUES (2, 'SqlAdminDisabled');
SQL

# Keine Erfassung ist kein Ereignis, das ein Poll beobachten kann — eine
# reale, aber begrenzte Wartezeit, bevor geprüft wird, dass die Zeile
# NICHT in cdc.changes ankommt.
sleep 3
disabled_leaked=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$ADMIN_TABLE' AND new_data->>'id' = '2'")
if [ "$disabled_leaked" != "0" ]; then
  echo "run-integration-tests: nach cdc.disable_table($ADMIN_TABLE) wurde eine Änderung (id=2) dennoch erfasst — Deaktivierung griff nicht" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach der SQL-Deaktivierung nicht mehr weiter (der Prozess muss weiterlaufen, nur die Tabellen-Erfassung endet)" >&2
  exit 1
fi

echo "run-integration-tests: SQL-Administration Live-Reload-Beleg (disable) — cdc.disable_table($ADMIN_TABLE) verarbeitet, Änderung id=2 nicht erfasst, Feed-Container läuft unverändert weiter"

abdeckung_declare "Publication-Entzug-Wirksamkeit" "LH-FA-CFG-002" "ein direkter Publication-Entzug per DDL trennt die PostgreSQL-seitige Filterung von der App-seitigen Assembler-Filterung: die Assembler-Bindung bleibt über den ganzen Beleg aktiv" "Publication-Entzug-Wirksamkeit — nach ALTER PUBLICATION"

# Publication-Entzug-Wirksamkeit — isolierter Beleg (BEO-PGC/walsender-wirksamkeit,
# LH-FA-CFG-002, ADR-0050, docs/reviews/architect-verdict-walsender-wirksamkeit.md):
# Der SQL-Administration Live-Reload-Beleg (disable) oben prüft nur die
# App-seitige Assembler-Filterung — RemoveBinding läuft synchron mit
# cdc.disable_table und verwirft jede Änderung der Tabelle, unabhängig
# davon, ob PostgreSQLs bereits laufende Decoding-Session die
# WAL-Änderung selbst noch gesendet hätte. Dieser Abschnitt trennt beide
# Ebenen: eine eigene, dedizierte Tabelle (nicht $ADMIN_TABLE) wird über
# die reguläre cdc.enable_table-Antragskette aktiviert und gebunden;
# danach wird die Publication-Mitgliedschaft NICHT über
# cdc.disable_table entzogen, sondern direkt per
# `ALTER PUBLICATION ... DROP TABLE` als $PG_USER gesetzt — es entsteht
# keine Zeile in cdc.administration_request, die Administrations-
# Goroutine wird nicht tätig, Assembler.RemoveBinding läuft für diese
# Tabelle nicht. Die Assembler-Bindung bleibt damit über den gesamten
# Testverlauf aktiv; jede Abwesenheit einer danach eingefügten Zeile in
# cdc.changes ist deshalb ausschließlich durch PostgreSQLs eigene
# Dekodierung erklärbar. Die Tabelle ist eine Wegwerf-Tabelle — kein
# späterer Abschnitt dieses Skripts liest oder schreibt sie —, eine
# Wiederherstellung der Publication-Mitgliedschaft entfällt deshalb.
WALSENDER_TABLE=feed_e2e_walsender_timing

walsender_enable_request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT cdc.enable_table('src-e2e', 'public', '$WALSENDER_TABLE')")
if [ -z "$walsender_enable_request_id" ]; then
  echo "run-integration-tests: cdc.enable_table($WALSENDER_TABLE) lieferte keine Antrags-ID" >&2
  exit 1
fi

walsender_enable_status=""
walsender_enable_applied=0
for _ in $(seq 1 60); do
  walsender_enable_status=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$walsender_enable_request_id'")
  if [ "$walsender_enable_status" = "applied" ]; then
    walsender_enable_applied=1
    break
  fi
  if [ "$walsender_enable_status" = "failed" ]; then
    error_message=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT error_message FROM cdc.administration_request WHERE administration_request_id = '$walsender_enable_request_id'")
    echo "run-integration-tests: cdc.enable_table($WALSENDER_TABLE)-Antrag $walsender_enable_request_id scheiterte: $error_message" >&2
    exit 1
  fi
  sleep 0.5
done
if [ "$walsender_enable_applied" -ne 1 ]; then
  echo "run-integration-tests: cdc.enable_table($WALSENDER_TABLE)-Antrag $walsender_enable_request_id wurde nicht innerhalb der Zeitspanne von der Administrations-Goroutine verarbeitet (status=${walsender_enable_status:-leer})" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$WALSENDER_TABLE (id, name) VALUES (1, 'WalsenderTimingEnabled');
SQL

walsender_enabled_captured=0
for _ in $(seq 1 120); do
  found=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$WALSENDER_TABLE' AND new_data->>'id' = '1'")
  if [ "$found" = "1" ]; then
    walsender_enabled_captured=1
    break
  fi
  sleep 0.25
done
if [ "$walsender_enabled_captured" -ne 1 ]; then
  echo "run-integration-tests: über SQL aktivierte Tabelle $WALSENDER_TABLE — Änderung (id=1) wurde nicht vom laufenden Feed-Container erfasst (kein Neustart)" >&2
  exit 1
fi

# Publication-Entzug DIREKT per DDL, ohne cdc.disable_table — siehe
# Begründung oben. Die Assembler-Bindung dieser Tabelle bleibt dadurch
# unverändert aktiv (activated == true).
docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -c \
  "ALTER PUBLICATION pub_pgc_e2e DROP TABLE public.$WALSENDER_TABLE;"

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$WALSENDER_TABLE (id, name) VALUES (2, 'WalsenderTimingAfterDrop');
SQL

# Keine Erfassung ist kein Ereignis, das ein Poll beobachten kann — eine
# reale, aber begrenzte Wartezeit (Muster wie beim SQL-Administration
# Live-Reload-Beleg oben), bevor das reale Ergebnis gegen cdc.changes
# ausgewertet wird. Beide Ausgänge sind ein gültiges, dokumentiertes
# Testergebnis dieses isolierten Belegs (kein Abbruch in beiden Fällen):
# Abwesenheit widerlegt die Verzögerungs-Annahme empirisch für die
# geprüfte PostgreSQL-Version; Anwesenheit bestätigt sie real und macht
# die App-seitige Assembler-Filterung zur tragenden Ebene.
sleep 3
walsender_after_drop_captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$WALSENDER_TABLE' AND new_data->>'id' = '2'")
if [ "$walsender_after_drop_captured" = "0" ]; then
  echo "run-integration-tests: Publication-Entzug-Wirksamkeit — nach ALTER PUBLICATION ... DROP TABLE $WALSENDER_TABLE (Assembler-Bindung blieb aktiv) wurde Änderung id=2 NICHT erfasst: PostgreSQLs bereits laufende Decoding-Session filtert eine entzogene Tabelle real sofort aus (BEO-PGC/walsender-wirksamkeit widerlegt)"
else
  echo "run-integration-tests: Publication-Entzug-Wirksamkeit — nach ALTER PUBLICATION ... DROP TABLE $WALSENDER_TABLE (Assembler-Bindung blieb aktiv) wurde Änderung id=2 DENNOCH erfasst (Anzahl: $walsender_after_drop_captured): PostgreSQLs Walsender liefert für eine bereits laufende Session real verzögert weiter (BEO-PGC/walsender-wirksamkeit bestätigt) — die App-seitige Assembler-Filterung ist die tragende Ebene"
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem Publication-Entzug-Wirksamkeit-Beleg nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

abdeckung_declare "Retention-Beleg (Alters- und Consumer-Freigabe)" "LH-FA-RET-002,LH-FA-RET-003,LH-FA-RET-004" "eine zurückdatierte Zeile bleibt erhalten, solange ein Consumer zurückhängt, und wird nach der Freigabe durch beide Consumer real entfernt; die zu junge Zeile bleibt durchgehend erhalten" "Retention-Beleg — 'RetentionOld' (id=200) blieb erhalten, solange ein Consumer zurückhing"

# Retention-Beleg (LH-FA-RET-002…004, ADR-0014): der Hintergrundzug
# runRetentionCleanup ruft RunRetentionUseCase periodisch real auf
# (retentionInterval, wiring.go). Zwei Zeilen auf feed_e2e_full (bereits
# über den ganzen Lauf aktiviert): 'RetentionOld' (id=200), deren
# Quelltransaktion direkt über SQL auf ein Alter über der konfigurierten
# RetentionPolicy.MinAge (retentionMinAge = 24h, wiring.go) zurückdatiert
# wird — derselbe Ansatz wie beim Fehlerzustand-Beleg oben, der
# cdc.process_heartbeat direkt schreibt, statt auf reale Zeit zu warten —,
# und 'RetentionYoung' (id=201), die real jung bleibt. Beide Zeilen liegen
# zunächst hinter der zuletzt bestätigten Position beider bereits
# geführten Consumer (CLI_CONSUMER, BACKLOG_CONSUMER) und sind damit
# zusätzlich zur Alters-Prüfung durch die Consumer-Positionen blockiert:
# der erste Poll unten belegt, dass die bereits alte Zeile trotzdem
# erhalten bleibt, solange ein Consumer zurückhängt (LH-FA-RET-004). Erst
# die zweite Bestätigung beider Consumer über die neue Position hinweg
# gibt beide Zeilen aus Consumer-Sicht frei; der zweite Poll belegt, dass
# danach nur die zurückdatierte Zeile real entfernt wird (LH-FA-RET-003),
# die junge nicht — beides ohne Neustart des Feed-Containers.
RETENTION_TABLE=feed_e2e_full

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$RETENTION_TABLE (id, name) VALUES (200, 'RetentionOld');
SQL

retention_old_position=""
for _ in $(seq 1 120); do
  retention_old_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RETENTION_TABLE' AND new_data->>'id' = '200'")
  if [ -n "$retention_old_position" ]; then
    break
  fi
  sleep 0.25
done
if [ -z "$retention_old_position" ]; then
  echo "run-integration-tests: Retention-Beleg — 'RetentionOld' (id=200) wurde nicht innerhalb der Zeitspanne über cdc.changes gelesen" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$RETENTION_TABLE (id, name) VALUES (201, 'RetentionYoung');
SQL

retention_young_position=""
for _ in $(seq 1 120); do
  retention_young_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RETENTION_TABLE' AND new_data->>'id' = '201'")
  if [ -n "$retention_young_position" ]; then
    break
  fi
  sleep 0.25
done
if [ -z "$retention_young_position" ]; then
  echo "run-integration-tests: Retention-Beleg — 'RetentionYoung' (id=201) wurde nicht innerhalb der Zeitspanne über cdc.changes gelesen" >&2
  exit 1
fi

docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -c \
  "UPDATE cdc.transaction SET committed_at = current_timestamp - interval '25 hours' WHERE transaction_id = (SELECT transaction_id FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RETENTION_TABLE' AND new_data->>'id' = '200')" >/dev/null

# Erste Phase (Consumer-Block, LH-FA-RET-004): eine großzügige, feste
# Wartezeit über mehr als zwei Lösch-Takte (retentionInterval=10s) hinweg
# — anders als ein Poll-auf-Verschwinden, das es hier nicht geben soll,
# belegt eine feste Wartezeit die Abwesenheit einer Löschung über
# mehrere reale Takte, nicht nur einen einzelnen zu frühen Blick.
sleep 25
old_present_blocked=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RETENTION_TABLE' AND new_data->>'id' = '200'")
if [ "$old_present_blocked" != "1" ]; then
  echo "run-integration-tests: Retention-Beleg — 'RetentionOld' (id=200, bereits alt genug) wurde entfernt, obwohl CLI_CONSUMER/BACKLOG_CONSUMER seine Position noch nicht bestätigt hatten (LH-FA-RET-004 Consumer-Block)" >&2
  exit 1
fi

if ! exec_feed acknowledge-consumer "$CLI_CONSUMER" "$retention_young_position"; then
  echo "run-integration-tests: Retention-Beleg — CLI_CONSUMER-Bestätigung über die Retention-Positionen hinweg endete mit einem Fehler" >&2
  exit 1
fi
if ! exec_feed acknowledge-consumer "$BACKLOG_CONSUMER" "$retention_young_position"; then
  echo "run-integration-tests: Retention-Beleg — BACKLOG_CONSUMER-Bestätigung über die Retention-Positionen hinweg endete mit einem Fehler" >&2
  exit 1
fi

# Zweite Phase (Alters-Freigabe, LH-FA-RET-003): Beide Zeilen sind ab
# dieser Bestätigung aus Consumer-Sicht frei — ein Poll auf das
# Verschwinden von id=200.
old_deleted=0
for _ in $(seq 1 60); do
  remaining=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RETENTION_TABLE' AND new_data->>'id' = '200'")
  if [ "$remaining" = "0" ]; then
    old_deleted=1
    break
  fi
  sleep 1
done
if [ "$old_deleted" -ne 1 ]; then
  echo "run-integration-tests: Retention-Beleg — 'RetentionOld' (id=200) wurde nach Freigabe durch beide Consumer nicht innerhalb der Zeitspanne real entfernt" >&2
  exit 1
fi

young_present=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RETENTION_TABLE' AND new_data->>'id' = '201'")
if [ "$young_present" != "1" ]; then
  echo "run-integration-tests: Retention-Beleg — 'RetentionYoung' (id=201, zu jung) wurde fälschlich real entfernt (LH-FA-RET-003 Mindestalter)" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem Retention-Beleg nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: Retention-Beleg — 'RetentionOld' (id=200) blieb erhalten, solange ein Consumer zurückhing (LH-FA-RET-004), und wurde nach Freigabe durch beide Consumer real entfernt (LH-FA-RET-003); 'RetentionYoung' (id=201, zu jung) blieb durchgehend erhalten"

abdeckung_declare "NATS-Happy-Path-Beleg" "LH-FA-SST-007" "ein Test-Subscriber abonniert das tabellen-granulare Subjekt real vor der Change und empfängt danach real das leere Wecksignal" "NATS-Happy-Path-Beleg (LH-FA-SST-007)"

abdeckung_declare "NATS-Boundary-Beleg" "LH-FA-SST-007" "eine Change ohne jeden auf dem Subjekt abonnierten Client bleibt vollständig über cdc.changes lesbar, und der Feed-Container läuft unverändert weiter" "NATS-Boundary-Beleg (LH-FA-SST-007)"

abdeckung_declare "NATS-Negative-Beleg (Reconnect-Nachholen)" "LH-FA-SST-007" "ein real vom Compose-Netz getrennter Subscriber verpasst die Change ohne jedes Wecksignal; ein frischer Wiederverbindungs-Subscriber empfängt die nächste" "NATS-Negative-Beleg (LH-FA-SST-007, Reconnect-Nachholen)"

# NATS-Happy-Path-Beleg (LH-FA-SST-007, ADR-0055, ADR-0056): compose.yaml
# verdrahtet den Feed-Container mit CDC_NATS_URL=nats://nats:4222 (siehe
# dortigen Kommentar). Ein eigener Wegwerf-Testclient
# (tools/harness/natssub/main.go, per `go run` im Toolchain-Container)
# abonniert das tabellen-granulare Subjekt
# cdc.changes.src-e2e.public.feed_e2e_full real, BEVOR die auslösende
# Change entsteht — Core NATS liefert nichts nach (ADR-0055 Punkt 1,
# Fire-and-Forget); ein Subscriber, der erst danach abonniert, verpasst das
# Signal strukturell. Der Subscriber läuft dazu als eigener, per Name
# adressierter Container (nicht `--rm` vor dem Poll): `docker logs` trägt
# die Zeile "READY", sobald die Subscription server-seitig bestätigt ist
# (Flush im Tool), erst danach folgt die Change. `feed_e2e_full` bleibt
# über den ganzen Lauf aktiviert (siehe Lasttest-/Retention-Belege oben);
# die ID 230 liegt in einem eigenen, bisher unbenutzten Wertebereich auf
# derselben Tabelle.
NATS_SUBSCRIBER_CONTAINER=cdc-e2e-natssub
NATS_SUBJECT="cdc.changes.src-e2e.public.feed_e2e_full"

docker rm -f "$NATS_SUBSCRIBER_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$NATS_SUBSCRIBER_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/natssub "nats://nats:4222" "$NATS_SUBJECT" >/dev/null

nats_subscriber_ready=0
for _ in $(seq 1 60); do
  if docker logs "$NATS_SUBSCRIBER_CONTAINER" 2>/dev/null | grep -qF "READY"; then
    nats_subscriber_ready=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$NATS_SUBSCRIBER_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 1
done
if [ "$nats_subscriber_ready" -ne 1 ]; then
  echo "run-integration-tests: NATS-Test-Subscriber ($NATS_SUBJECT) wurde nicht innerhalb der Zeitspanne bereit: $(docker logs "$NATS_SUBSCRIBER_CONTAINER" 2>&1 || true)" >&2
  docker rm -f "$NATS_SUBSCRIBER_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.feed_e2e_full (id, name) VALUES (230, 'NatsHappyPath');
SQL

nats_signal_received=0
for _ in $(seq 1 60); do
  if docker logs "$NATS_SUBSCRIBER_CONTAINER" 2>/dev/null | grep -qF "RECEIVED"; then
    nats_signal_received=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$NATS_SUBSCRIBER_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 0.5
done
nats_subscriber_output=$(docker logs "$NATS_SUBSCRIBER_CONTAINER" 2>&1 || true)
docker rm -f "$NATS_SUBSCRIBER_CONTAINER" >/dev/null 2>&1 || true
if [ "$nats_signal_received" -ne 1 ]; then
  echo "run-integration-tests: NATS-Test-Subscriber ($NATS_SUBJECT) empfing nach der Change (id=230) kein Wecksignal innerhalb der Zeitspanne: $nats_subscriber_output" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem NATS-Happy-Path-Beleg nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: NATS-Happy-Path-Beleg (LH-FA-SST-007) — Test-Subscriber ($NATS_SUBJECT) abonnierte real vor der Change (id=230, feed_e2e_full) und empfing danach real das leere Wecksignal: $nats_subscriber_output"

# NATS-Boundary-Beleg (LH-FA-SST-007, ADR-0055): Gegenstück zum
# Happy-Path-Beleg oben — hier abonniert **kein** Client das Subjekt
# cdc.changes.src-e2e.public.feed_e2e_full. Der reale Beleg, dass zum
# Zeitpunkt der Change tatsächlich niemand verbunden ist, kommt aus dem
# NATS-Server selbst: dessen HTTP-Monitor-Port 8222 (bereits für den
# Healthcheck aktiviert, siehe compose.yaml) trägt den Endpoint
# /subsz?subs=1, der jede aktive Subscription mit ihrem Subjekt auflistet
# (auch die internen $SYS-Subscriptions des Servers selbst — deren Subjekte
# beginnen mit "$SYS." und überschneiden sich nie mit "cdc.changes."). Ein
# `docker exec` gegen den NATS-Container liest diese Liste unmittelbar vor
# der Change und stellt sicher, dass unser Subjekt darin nicht vorkommt —
# kein Rückschluss aus "wir haben keinen Subscriber gestartet", sondern eine
# reale Server-seitige Momentaufnahme. Nach der Change bleibt der
# Feed-Container weiter am Leben (CaptureService.Capture()s Notify-Aufruf
# darf einen fehlenden Empfänger nicht als Fehler werten, ADR-0055 Punkt 4)
# und die Change ist über den bestehenden Lesezugriffsweg cdc.changes
# vollständig auffindbar.
nats_subs_before_boundary=$(docker exec "$NATS_CONTAINER" wget -q -O - "http://localhost:8222/subsz?subs=1")
if echo "$nats_subs_before_boundary" | grep -qF "\"subject\": \"$NATS_SUBJECT\""; then
  echo "run-integration-tests: NATS-Boundary-Beleg — vor der auslösenden Change war entgegen der Erwartung bereits ein Client auf $NATS_SUBJECT abonniert (Testablauf hat den Boundary-Fall nicht real hergestellt): $nats_subs_before_boundary" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.feed_e2e_full (id, name) VALUES (231, 'NatsBoundary');
SQL

# Kleine reale Wartezeit, damit ein (fälschlich doch aktiver) Notify-Versuch
# und ein etwaiger Absturz des Feed-Containers Zeit hätten, sich zu zeigen,
# bevor die Belege unten gezogen werden — kein Poll auf ein Ereignis, das
# hier per Definition nicht eintreten soll.
sleep 2

boundary_change_present=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = 'feed_e2e_full' AND new_data->>'id' = '231'")
if [ "$boundary_change_present" != "1" ]; then
  echo "run-integration-tests: NATS-Boundary-Beleg — Change (id=231, feed_e2e_full) war trotz fehlendem NATS-Subscriber nicht vollständig über cdc.changes lesbar (LH-FA-SST-007 Boundary verletzt)" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem NATS-Boundary-Beleg nicht mehr weiter — ein Notify ohne Empfänger hätte CaptureService.Capture() nicht blockieren oder beenden dürfen (ADR-0055 Punkt 4)" >&2
  exit 1
fi

nats_subs_after_boundary=$(docker exec "$NATS_CONTAINER" wget -q -O - "http://localhost:8222/subsz?subs=1")
if echo "$nats_subs_after_boundary" | grep -qF "\"subject\": \"$NATS_SUBJECT\""; then
  echo "run-integration-tests: NATS-Boundary-Beleg — nach der Change war entgegen der Erwartung ein Client auf $NATS_SUBJECT abonniert (Testablauf hat den Boundary-Fall nicht real hergestellt): $nats_subs_after_boundary" >&2
  exit 1
fi

echo "run-integration-tests: NATS-Boundary-Beleg (LH-FA-SST-007) — Change (id=231, feed_e2e_full) entstand real ohne einen auf $NATS_SUBJECT abonnierten Client (belegt über NATS-Server-Monitor /subsz vor und nach der Change), blieb vollständig über cdc.changes lesbar, und der Feed-Container lief unverändert weiter"

# NATS-Negative-Beleg — Reconnect-Nachholen (LH-FA-SST-007, ADR-0055,
# ADR-0056): Gegenstück zum Boundary-Beleg oben (dort: nie abonniert
# gewesen) — hier ist der Test-Subscriber bereits real verbunden und
# abonniert, bevor er real vom Compose-Netz getrennt wird
# (`docker network disconnect "$NETWORK" "$NATS_RECONNECT_BEFORE_CONTAINER"`),
# NICHT der NATS-Server-Container selbst: ein gestoppter Server träfe auch
# den Feed-Container (Publisher) und würde den Testfall verfälschen. Die
# reale Trennung wird über `docker inspect` (Feld `NetworkSettings.Networks`)
# belegt — eine sofortige, von Docker selbst geführte Zustandsauskunft,
# unabhängig davon, wie schnell der NATS-Server eine unterbrochene
# TCP-Verbindung als tot erkennt (server-seitige Ping-Erkennung liegt in
# einer Größenordnung von Minuten, nicht Sekunden, und ist für die reale
# Trennung selbst kein Beleg). `docker logs` des Subscribers belegt
# zusätzlich explizit das Fehlen jedes "RECEIVED"-Frames während der
# Trennung, nicht nur, dass die währenddessen entstandene Change über
# `cdc.changes` sichtbar ist. Die Wiederverbindung wird über einen
# frischen, zweiten Subscriber-Prozess hergestellt statt über
# `docker network connect` auf demselben Container: Eine während der
# Trennung vom Broker noch nicht als tot erkannte TCP-Verbindung könnte
# nach dem Heilen des Netzpfads eine zwischenzeitlich im Sendepuffer
# hängengebliebene Zustellung nachholen — genau die Zweideutigkeit, die
# dieser Beleg ausschließen soll. Ein neuer Subscriber-Prozess hat
# dagegen strukturell keine Möglichkeit, ein vor seiner eigenen
# Subscription publiziertes Signal zu empfangen (dieselbe Begründung wie
# beim Boundary-Beleg oben), was den Negativ-Beleg für die verpasste
# Change eindeutig macht. Die verpasste Change wird ausschließlich über
# den bestehenden SQL-Lesezugriffsweg `cdc.changes` nachgeholt; der neue
# Subscriber-Prozess belegt zusätzlich, dass das Wecksignal für eine
# danach entstehende Change unverändert funktioniert.
NATS_RECONNECT_BEFORE_CONTAINER=cdc-e2e-natssub-reconnect-before
NATS_RECONNECT_AFTER_CONTAINER=cdc-e2e-natssub-reconnect-after

docker rm -f "$NATS_RECONNECT_BEFORE_CONTAINER" "$NATS_RECONNECT_AFTER_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$NATS_RECONNECT_BEFORE_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/natssub "nats://nats:4222" "$NATS_SUBJECT" >/dev/null

nats_reconnect_before_ready=0
for _ in $(seq 1 60); do
  if docker logs "$NATS_RECONNECT_BEFORE_CONTAINER" 2>/dev/null | grep -qF "READY"; then
    nats_reconnect_before_ready=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$NATS_RECONNECT_BEFORE_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 1
done
if [ "$nats_reconnect_before_ready" -ne 1 ]; then
  echo "run-integration-tests: NATS-Negative-Beleg — Test-Subscriber ($NATS_SUBJECT) wurde nicht innerhalb der Zeitspanne bereit: $(docker logs "$NATS_RECONNECT_BEFORE_CONTAINER" 2>&1 || true)" >&2
  docker rm -f "$NATS_RECONNECT_BEFORE_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

nats_subs_before_disconnect=$(docker exec "$NATS_CONTAINER" wget -q -O - "http://localhost:8222/subsz?subs=1")
if ! echo "$nats_subs_before_disconnect" | grep -qF "\"subject\": \"$NATS_SUBJECT\""; then
  echo "run-integration-tests: NATS-Negative-Beleg — Test-Subscriber ($NATS_SUBJECT) war vor der Trennung entgegen der Erwartung nicht als Abonnent beim Broker gelistet: $nats_subs_before_disconnect" >&2
  docker rm -f "$NATS_RECONNECT_BEFORE_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

docker network disconnect "$NETWORK" "$NATS_RECONNECT_BEFORE_CONTAINER"

subscriber_networks_after_disconnect=$(docker inspect "$NATS_RECONNECT_BEFORE_CONTAINER" --format '{{json .NetworkSettings.Networks}}')
if echo "$subscriber_networks_after_disconnect" | grep -qF "\"$NETWORK\""; then
  echo "run-integration-tests: NATS-Negative-Beleg — Test-Subscriber-Container trägt laut docker inspect nach dem Trennungsversuch noch die Netzbindung $NETWORK (reale Trennung nicht hergestellt): $subscriber_networks_after_disconnect" >&2
  docker rm -f "$NATS_RECONNECT_BEFORE_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.feed_e2e_full (id, name) VALUES (240, 'NatsReconnectMissed');
SQL

# Kleine reale Wartezeit: ein (fälschlich doch zugestelltes) Signal hätte
# hier Zeit, im Log aufzutauchen, bevor der Negativ-Beleg unten gezogen
# wird — kein Poll auf ein Ereignis, das hier per Definition nicht
# eintreten soll.
sleep 5

nats_reconnect_before_output=$(docker logs "$NATS_RECONNECT_BEFORE_CONTAINER" 2>&1 || true)
docker rm -f "$NATS_RECONNECT_BEFORE_CONTAINER" >/dev/null 2>&1 || true
if printf '%s' "$nats_reconnect_before_output" | grep -qF "RECEIVED"; then
  echo "run-integration-tests: NATS-Negative-Beleg — Test-Subscriber empfing trotz realer Trennung ein Wecksignal für die verpasste Change (id=240): $nats_reconnect_before_output" >&2
  exit 1
fi

reconnect_missed_present=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = 'feed_e2e_full' AND new_data->>'id' = '240'")
if [ "$reconnect_missed_present" != "1" ]; then
  echo "run-integration-tests: NATS-Negative-Beleg — verpasste Change (id=240) war nach der Trennung nicht über den bestehenden SQL-Lesezugriffsweg cdc.changes vollständig sichtbar" >&2
  exit 1
fi

# Wiederverbindung als frischer Subscriber-Prozess (siehe Begründung oben):
# er abonniert dasselbe Subjekt normal verbunden, bevor die neue Change
# entsteht — dieselbe Reihenfolge-Disziplin wie beim Happy-Path-Beleg.
docker run -d --name "$NATS_RECONNECT_AFTER_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/natssub "nats://nats:4222" "$NATS_SUBJECT" >/dev/null

nats_reconnect_after_ready=0
for _ in $(seq 1 60); do
  if docker logs "$NATS_RECONNECT_AFTER_CONTAINER" 2>/dev/null | grep -qF "READY"; then
    nats_reconnect_after_ready=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$NATS_RECONNECT_AFTER_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 1
done
if [ "$nats_reconnect_after_ready" -ne 1 ]; then
  echo "run-integration-tests: NATS-Negative-Beleg — Wiederverbindungs-Subscriber ($NATS_SUBJECT) wurde nicht innerhalb der Zeitspanne bereit: $(docker logs "$NATS_RECONNECT_AFTER_CONTAINER" 2>&1 || true)" >&2
  docker rm -f "$NATS_RECONNECT_AFTER_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.feed_e2e_full (id, name) VALUES (241, 'NatsReconnectResumed');
SQL

nats_reconnect_signal_resumed=0
for _ in $(seq 1 60); do
  if docker logs "$NATS_RECONNECT_AFTER_CONTAINER" 2>/dev/null | grep -qF "RECEIVED"; then
    nats_reconnect_signal_resumed=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$NATS_RECONNECT_AFTER_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 0.5
done
nats_reconnect_after_output=$(docker logs "$NATS_RECONNECT_AFTER_CONTAINER" 2>&1 || true)
docker rm -f "$NATS_RECONNECT_AFTER_CONTAINER" >/dev/null 2>&1 || true
if [ "$nats_reconnect_signal_resumed" -ne 1 ]; then
  echo "run-integration-tests: NATS-Negative-Beleg — Wiederverbindungs-Subscriber empfing für eine neue Change (id=241) kein Wecksignal innerhalb der Zeitspanne — die Verbindung wäre damit nicht real wiederhergestellt gewesen: $nats_reconnect_after_output" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem NATS-Negative-Beleg (Reconnect-Nachholen) nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: NATS-Negative-Beleg (LH-FA-SST-007, Reconnect-Nachholen) — Test-Subscriber real vom Compose-Netz getrennt (belegt über docker inspect), verpasste Change (id=240) blieb ohne jedes Wecksignal (Log-Beleg) und wurde ausschließlich über cdc.changes nachgeholt; ein frischer Wiederverbindungs-Subscriber empfing für eine neue Change (id=241) real ein Signal, ohne dass die verpasste Change nachträglich zugestellt wurde: $nats_reconnect_after_output"

abdeckung_declare "HTTP-API-Rundlauf" "LH-FA-SST-006" "ein Wegwerf-Client ruft RegisterConsumer mit dem admin-Token und ListTables mit dem reader-Token real per HTTP gegen den laufenden Feed-Container auf, die Registrierung wird gegen cdc.consumer bestätigt" "HTTP-API-Rundlauf (LH-FA-SST-006, ADR-0057) belegt"

# HTTP-API-Rundlauf (LH-FA-SST-006, ADR-0057): ein Wegwerf-Client
# (tools/harness/httpclient) ruft RegisterConsumer mit dem admin-Token und
# ListTables mit dem reader-Token real per HTTP gegen den laufenden
# Feed-Container auf — ein echter Netzwerk-Request, kein `docker exec` und
# kein Mock. `feed_e2e_full` bleibt über den ganzen Lauf aktiviert (siehe
# Lasttest-Beleg oben), ihr Auftreten in der ListTables-Antwort belegt
# einen realen, nicht leeren Rückgabewert. Muss vor der abschließenden
# Schema-Negative-Testfunktion unten laufen, weil diese den Feed-Container
# dauerhaft beendet.
HTTP_CONSUMER=http-e2e-consumer
http_output=$(docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/httpclient \
  "$HTTP_BASE_URL" "$HTTP_TOKEN_ADMIN" "$HTTP_TOKEN_READER" "$HTTP_CONSUMER" "HTTP E2E Consumer" src-e2e pub_pgc_e2e 2>&1)
http_status=$?
if [ "$http_status" -ne 0 ]; then
  echo "run-integration-tests: HTTP-API-Rundlauf (httpclient) endete mit Ausgang $http_status: $http_output" >&2
  exit 1
fi
if ! printf '%s' "$http_output" | grep -qF "REGISTERED"; then
  echo "run-integration-tests: HTTP-API-Rundlauf — keine REGISTERED-Zeile (admin-Token, RegisterConsumer): $http_output" >&2
  exit 1
fi
if ! printf '%s' "$http_output" | grep -qE '"table":"feed_e2e_full"'; then
  echo "run-integration-tests: HTTP-API-Rundlauf — ListTables (reader-Token) trägt die dauerhaft aktivierte Tabelle feed_e2e_full nicht: $http_output" >&2
  exit 1
fi

registered_via_http=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT consumer_id FROM cdc.consumer WHERE consumer_id = '$HTTP_CONSUMER'")
if [ "$registered_via_http" != "$HTTP_CONSUMER" ]; then
  echo "run-integration-tests: cdc.consumer trägt $HTTP_CONSUMER nicht nach dem realen HTTP-RegisterConsumer-Aufruf" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem HTTP-API-Rundlauf nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: HTTP-API-Rundlauf (LH-FA-SST-006, ADR-0057) belegt — RegisterConsumer real per HTTP mit admin-Token ($HTTP_CONSUMER, cdc.consumer bestätigt), ListTables real per HTTP mit reader-Token (feed_e2e_full in der Antwort): $http_output"

abdeckung_declare "gRPC-Stream-Rundlauf" "LH-FA-SST-008" "ein Wegwerf-Client öffnet real über gRPC den Server-Stream gegen den laufenden Feed-Container und empfängt eine danach committete Änderung; ein Öffnungsversuch ohne gültiges Token endet mit gRPC-Status Unauthenticated" "gRPC-Stream-Rundlauf (LH-FA-SST-008, ADR-0060) belegt"

# gRPC-Stream-Rundlauf (LH-FA-SST-008, ADR-0060): ein Wegwerf-Client
# (tools/harness/grpcclient, per `go run` im Toolchain-Container) verbindet
# sich real über gRPC mit dem laufenden Feed-Container, öffnet den
# Server-Stream und empfängt eine danach committete Änderung. Geprüft wird
# die RECEIVED-Zeile über Tabelle, Operation und den Wert der Spalte `name`
# im Stream-Image und ihre `change_id` gegen den Lesezugriffsweg
# `cdc.changes`; die Feldvollständigkeit des Nachrichtenschemas trägt
# `internal/adapters/driving/grpc/server_test.go` auf Unit-Ebene.
# Abschließend belegt derselbe Prozess, dass ein
# Stream-Öffnungsversuch ohne gültiges Token über gRPC-Status
# `Unauthenticated` abgelehnt wird. Ein echter Netzwerk-Request über den
# Compose-Netz-Alias `pg-change-feed:9090` (CDC_GRPC_ADDR im
# Container-Vertrag), kein `docker exec` und kein Mock; der Client trägt
# dieselbe `reader`-Token-Klasse wie der HTTP-Adapter (CDC_API_TOKEN_READER).
# `feed_e2e_full` bleibt über den ganzen Lauf aktiviert (siehe
# Lasttest-Beleg oben); die IDs 260ff. liegen in einem eigenen Wertebereich.
# Läuft vor dem Upgrade-Sicherheits-Container-Tausch (der den Feed-Container
# ersetzt) und vor der Container-Ende-Grenze der beiden
# Schema-Negative-Testfunktionen unten.
GRPC_CLIENT_CONTAINER=cdc-e2e-grpcclient
GRPC_STREAM_TABLE=feed_e2e_full
GRPC_STREAM_SENTINEL=GrpcStreamE2ESentinel
GRPC_ADDR="pg-change-feed:9090"

docker rm -f "$GRPC_CLIENT_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$GRPC_CLIENT_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/grpcclient "$GRPC_ADDR" "$HTTP_TOKEN_READER" >/dev/null

grpc_client_ready=0
for _ in $(seq 1 60); do
  if docker logs "$GRPC_CLIENT_CONTAINER" 2>/dev/null | grep -qF "READY"; then
    grpc_client_ready=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$GRPC_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 1
done
if [ "$grpc_client_ready" -ne 1 ]; then
  echo "run-integration-tests: gRPC-Stream-Rundlauf — Test-Client wurde nicht innerhalb der Zeitspanne bereit: $(docker logs "$GRPC_CLIENT_CONTAINER" 2>&1 || true)" >&2
  docker rm -f "$GRPC_CLIENT_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

# Die Zustellung ist Fire-and-Forget ohne Replay (ADR-0066): zwischen
# „Client hat den Stream geöffnet" und „Server hat den Empfänger am
# Broadcaster registriert" liegt ein kurzes Fenster, in dem eine committete
# Änderung für diesen Empfänger verworfen wird. Der Rundlauf committet
# deshalb eine begrenzte Folge eindeutiger Zeilen, bis der Client genau eine
# davon real empfangen hat.
grpc_received=0
for grpc_attempt in $(seq 1 5); do
  grpc_id=$((260 + grpc_attempt))
  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$GRPC_STREAM_TABLE (id, name) VALUES ($grpc_id, '$GRPC_STREAM_SENTINEL');
SQL
  for _ in $(seq 1 20); do
    if docker logs "$GRPC_CLIENT_CONTAINER" 2>/dev/null | grep -qF "RECEIVED"; then
      grpc_received=1
      break
    fi
    if [ "$(docker inspect --format '{{.State.Running}}' "$GRPC_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
      break
    fi
    sleep 0.5
  done
  if [ "$grpc_received" -eq 1 ]; then
    break
  fi
done

# Der Client führt nach dem Empfang im selben Prozess noch die
# Negative-Prüfung aus (Stream-Öffnungsversuch ohne Token); auf deren
# Ergebnis und auf das Prozessende wird separat gewartet, damit die
# ausgewerteten Zeilen unten aus einem abgeschlossenen Lauf stammen.
grpc_rejected=0
for _ in $(seq 1 40); do
  if docker logs "$GRPC_CLIENT_CONTAINER" 2>/dev/null | grep -qF "REJECTED code=Unauthenticated"; then
    grpc_rejected=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$GRPC_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 0.5
done

grpc_client_stopped=0
for _ in $(seq 1 20); do
  if [ "$(docker inspect --format '{{.State.Running}}' "$GRPC_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    grpc_client_stopped=1
    break
  fi
  sleep 0.5
done
grpc_client_exit=$(docker inspect --format '{{.State.ExitCode}}' "$GRPC_CLIENT_CONTAINER" 2>/dev/null || echo unbekannt)
grpc_client_output=$(docker logs "$GRPC_CLIENT_CONTAINER" 2>&1 || true)
docker rm -f "$GRPC_CLIENT_CONTAINER" >/dev/null 2>&1 || true

if [ "$grpc_received" -ne 1 ]; then
  echo "run-integration-tests: gRPC-Stream-Rundlauf — Test-Client empfing keine der committeten Änderungen ($GRPC_STREAM_TABLE, $GRPC_STREAM_SENTINEL) über den Stream: $grpc_client_output" >&2
  exit 1
fi
if ! printf '%s' "$grpc_client_output" | grep -qE "RECEIVED .*table=$GRPC_STREAM_TABLE .*operation=INSERT .*new_image=.*$GRPC_STREAM_SENTINEL"; then
  echo "run-integration-tests: gRPC-Stream-Rundlauf — die RECEIVED-Zeile trägt nicht die erwartete Änderung ($GRPC_STREAM_TABLE, INSERT, vollständiger Inhalt): $grpc_client_output" >&2
  exit 1
fi
if [ "$grpc_rejected" -ne 1 ]; then
  echo "run-integration-tests: gRPC-Stream-Rundlauf — Stream-Öffnungsversuch ohne gültiges Token wurde nicht mit Unauthenticated abgelehnt: $grpc_client_output" >&2
  exit 1
fi
if [ "$grpc_client_stopped" -ne 1 ] || [ "$grpc_client_exit" != "0" ]; then
  echo "run-integration-tests: gRPC-Stream-Rundlauf — Test-Client endete nicht mit Ausgang 0 (gestoppt: $grpc_client_stopped, Ausgang: $grpc_client_exit): $grpc_client_output" >&2
  exit 1
fi

# Unabhängiger SQL-Beleg, dass genau die empfangene Änderung real erfasst
# wurde: die change_id der RECEIVED-Zeile steht über den bestehenden
# Lesezugriffsweg in cdc.changes — der Stream-Empfang ist damit keine
# erfundene Ausgabe des Clients.
grpc_change_id=$(printf '%s' "$grpc_client_output" | grep -oE 'RECEIVED change_id=[^ ]+' | head -n1 | cut -d= -f2 || true)
if [ -z "$grpc_change_id" ]; then
  echo "run-integration-tests: gRPC-Stream-Rundlauf — die RECEIVED-Zeile trägt keine change_id: $grpc_client_output" >&2
  exit 1
fi
grpc_captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$grpc_change_id' AND table_name = '$GRPC_STREAM_TABLE' AND new_data->>'name' = '$GRPC_STREAM_SENTINEL'")
if [ -z "$grpc_captured" ] || [ "$grpc_captured" -lt 1 ]; then
  echo "run-integration-tests: gRPC-Stream-Rundlauf — die über den Stream empfangene Änderung (change_id=$grpc_change_id, $GRPC_STREAM_SENTINEL) ist nicht real über cdc.changes lesbar (count=${grpc_captured:-leer})" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem gRPC-Stream-Rundlauf nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: gRPC-Stream-Rundlauf (LH-FA-SST-008, ADR-0060) belegt — ein Wegwerf-Client (tools/harness/grpcclient) öffnete real über gRPC den Server-Stream gegen den laufenden Feed-Container ($GRPC_ADDR) und empfing eine danach committete Änderung (Tabelle, Operation und Spaltenwert real am Stream; die Feldvollständigkeit trägt server_test.go auf Unit-Ebene), deren change_id ($grpc_change_id) unabhängig über cdc.changes lesbar ist; ein Stream-Öffnungsversuch ohne gültiges Token wurde mit gRPC-Status Unauthenticated abgelehnt: $grpc_client_output"

abdeckung_declare "SSE-Stream-Rundlauf" "LH-FA-SST-008" "ein Wegwerf-Client öffnet real per HTTP den Endpunkt `GET /changes/stream` gegen den laufenden Feed-Container und empfängt eine danach committete Änderung; ein Aufruf ohne gültiges Token endet mit HTTP-Status 401" "SSE-Stream-Rundlauf (LH-FA-SST-008, ADR-0061) belegt"

# SSE-Stream-Rundlauf (LH-FA-SST-008, ADR-0061): ein Wegwerf-Client
# (tools/harness/sseclient, per `go run` im Toolchain-Container) verbindet
# sich real per HTTP mit dem laufenden Feed-Container, öffnet den
# SSE-Stream `GET /changes/stream` und empfängt eine danach committete
# Änderung. Geprüft wird die RECEIVED-Zeile über Tabelle, Operation und den
# Wert der Spalte `name` im neuen Row Image und ihre `change_id` gegen den
# Lesezugriffsweg `cdc.changes`; die Feldvollständigkeit des
# Nachrichtenschemas trägt `internal/adapters/driving/http/sse_test.go` auf
# Unit-Ebene. Abschließend belegt derselbe Prozess, dass ein
# Öffnungsversuch ohne gültiges Token mit HTTP-Status 401 abgelehnt wird.
# Ein echter Netzwerk-Request über den Compose-Netz-Alias
# `pg-change-feed:8090` (CDC_HTTP_ADDR im Container-Vertrag), kein
# `docker exec` und kein Mock; der Client trägt dieselbe `reader`-Klasse
# wie der übrige HTTP-Zugriff (CDC_API_TOKEN_READER). `feed_e2e_full` bleibt
# über den ganzen Lauf aktiviert (siehe Lasttest-Beleg oben); die IDs 270ff.
# liegen in einem eigenen Wertebereich. Läuft vor dem
# Upgrade-Sicherheits-Container-Tausch (der den Feed-Container ersetzt) und
# vor der Container-Ende-Grenze der beiden Schema-Negative-Testfunktionen
# unten.
SSE_CLIENT_CONTAINER=cdc-e2e-sseclient
SSE_STREAM_TABLE=feed_e2e_full
SSE_STREAM_SENTINEL=SseStreamE2ESentinel

docker rm -f "$SSE_CLIENT_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$SSE_CLIENT_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/sseclient "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" >/dev/null

sse_client_ready=0
for _ in $(seq 1 60); do
  if docker logs "$SSE_CLIENT_CONTAINER" 2>/dev/null | grep -qF "READY"; then
    sse_client_ready=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$SSE_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 1
done
if [ "$sse_client_ready" -ne 1 ]; then
  echo "run-integration-tests: SSE-Stream-Rundlauf — Test-Client wurde nicht innerhalb der Zeitspanne bereit: $(docker logs "$SSE_CLIENT_CONTAINER" 2>&1 || true)" >&2
  docker rm -f "$SSE_CLIENT_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

# Die Zustellung ist Fire-and-Forget ohne Replay (ADR-0066): der Rundlauf
# committet eine begrenzte Folge eindeutiger Zeilen, bis der Client genau
# eine davon real empfangen hat.
sse_received=0
for sse_attempt in $(seq 1 5); do
  sse_id=$((270 + sse_attempt))
  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$SSE_STREAM_TABLE (id, name) VALUES ($sse_id, '$SSE_STREAM_SENTINEL');
SQL
  for _ in $(seq 1 20); do
    if docker logs "$SSE_CLIENT_CONTAINER" 2>/dev/null | grep -qF "RECEIVED"; then
      sse_received=1
      break
    fi
    if [ "$(docker inspect --format '{{.State.Running}}' "$SSE_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
      break
    fi
    sleep 0.5
  done
  if [ "$sse_received" -eq 1 ]; then
    break
  fi
done

# Der Client führt nach dem Empfang im selben Prozess noch die
# Negative-Prüfung aus (Stream-Öffnungsversuch ohne Token); auf deren
# Ergebnis und auf das Prozessende wird separat gewartet, damit die
# ausgewerteten Zeilen unten aus einem abgeschlossenen Lauf stammen.
sse_rejected=0
for _ in $(seq 1 40); do
  if docker logs "$SSE_CLIENT_CONTAINER" 2>/dev/null | grep -qF "REJECTED code=401"; then
    sse_rejected=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$SSE_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 0.5
done

sse_client_stopped=0
for _ in $(seq 1 20); do
  if [ "$(docker inspect --format '{{.State.Running}}' "$SSE_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    sse_client_stopped=1
    break
  fi
  sleep 0.5
done
sse_client_exit=$(docker inspect --format '{{.State.ExitCode}}' "$SSE_CLIENT_CONTAINER" 2>/dev/null || echo unbekannt)
sse_client_output=$(docker logs "$SSE_CLIENT_CONTAINER" 2>&1 || true)
docker rm -f "$SSE_CLIENT_CONTAINER" >/dev/null 2>&1 || true

if [ "$sse_received" -ne 1 ]; then
  echo "run-integration-tests: SSE-Stream-Rundlauf — Test-Client empfing keine der committeten Änderungen ($SSE_STREAM_TABLE, $SSE_STREAM_SENTINEL) über den Stream: $sse_client_output" >&2
  exit 1
fi
if ! printf '%s' "$sse_client_output" | grep -qE "RECEIVED .*table=$SSE_STREAM_TABLE .*operation=INSERT .*new_image=.*$SSE_STREAM_SENTINEL"; then
  echo "run-integration-tests: SSE-Stream-Rundlauf — die RECEIVED-Zeile trägt nicht die erwartete Änderung ($SSE_STREAM_TABLE, INSERT, vollständiger Inhalt): $sse_client_output" >&2
  exit 1
fi
if [ "$sse_rejected" -ne 1 ]; then
  echo "run-integration-tests: SSE-Stream-Rundlauf — Stream-Öffnungsversuch ohne gültiges Token wurde nicht mit 401 abgelehnt: $sse_client_output" >&2
  exit 1
fi
if [ "$sse_client_stopped" -ne 1 ] || [ "$sse_client_exit" != "0" ]; then
  echo "run-integration-tests: SSE-Stream-Rundlauf — Test-Client endete nicht mit Ausgang 0 (gestoppt: $sse_client_stopped, Ausgang: $sse_client_exit): $sse_client_output" >&2
  exit 1
fi

# Unabhängiger SQL-Beleg, dass genau die empfangene Änderung real erfasst
# wurde: die change_id der RECEIVED-Zeile steht über den bestehenden
# Lesezugriffsweg in cdc.changes — der Stream-Empfang ist damit keine
# erfundene Ausgabe des Clients.
sse_change_id=$(printf '%s' "$sse_client_output" | grep -oE 'RECEIVED change_id=[^ ]+' | head -n1 | cut -d= -f2 || true)
if [ -z "$sse_change_id" ]; then
  echo "run-integration-tests: SSE-Stream-Rundlauf — die RECEIVED-Zeile trägt keine change_id: $sse_client_output" >&2
  exit 1
fi
sse_captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$sse_change_id' AND table_name = '$SSE_STREAM_TABLE' AND new_data->>'name' = '$SSE_STREAM_SENTINEL'")
if [ -z "$sse_captured" ] || [ "$sse_captured" -lt 1 ]; then
  echo "run-integration-tests: SSE-Stream-Rundlauf — die über den Stream empfangene Änderung (change_id=$sse_change_id, $SSE_STREAM_SENTINEL) ist nicht real über cdc.changes lesbar (count=${sse_captured:-leer})" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem SSE-Stream-Rundlauf nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: SSE-Stream-Rundlauf (LH-FA-SST-008, ADR-0061) belegt — ein Wegwerf-Client (tools/harness/sseclient) öffnete real per HTTP den SSE-Stream GET /changes/stream gegen den laufenden Feed-Container ($HTTP_BASE_URL) und empfing eine danach committete Änderung (Tabelle, Operation und Spaltenwert real am Stream; die Feldvollständigkeit trägt sse_test.go auf Unit-Ebene), deren change_id ($sse_change_id) unabhängig über cdc.changes lesbar ist; ein Stream-Öffnungsversuch ohne gültiges Token wurde mit HTTP-Status 401 abgelehnt: $sse_client_output"

abdeckung_declare "Upgrade-Sicherheits-Rundlauf" "LH-QA-OPS-005" "ein realer Container-Tausch ersetzt den Feed-Container durch eine neue Instanz desselben Images, während Datenbank und NATS unberührt bleiben — der Datenstand davor bleibt lesbar, danach Eingefügtes wird weiter erfasst" "Upgrade-Sicherheits-Rundlauf (LH-QA-OPS-005, ADR-0064) belegt"

# Upgrade-Sicherheits-Rundlauf (LH-QA-OPS-005, ADR-0064 Supersedes
# ADR-0058 Entscheidung 3): bildet den Mechanismus eines
# Anwendungs-Upgrades nach — ein realer Container-Tausch über
# `$COMPOSE up -d --force-recreate --no-deps pg-change-feed` ersetzt den
# in ADR-0058 vorgesehenen, real blockierten zweiten
# `make schema-rollout`-Lauf (BEO-PGC/schema-rollout-fremdobjekte, Exit 8
# auf vier Fremdobjekten, docs/reviews/blocker-slice-063.md). Läuft hier,
# solange der Feed-Container noch unversehrt und gesund ist — vor
# TestE2ESchemaChangeDropColumn/TestE2ESchemaChangeIncompatibleTypeChange
# unten, die ihn beide dauerhaft beenden. `feed_e2e_full` bleibt über den
# ganzen Lauf aktiviert (siehe Lasttest-Beleg oben); die IDs 250/251 liegen
# in einem eigenen Wertebereich, getrennt von den übrigen Testfall-Gruppen
# auf derselben Tabelle.
UPGRADE_TABLE=feed_e2e_full

upgrade_feed_id_before=$(docker inspect --format '{{.Id}}' "$FEED_CONTAINER")
upgrade_pg_id_before=$(docker inspect --format '{{.Id}}' "$PG_CONTAINER")
upgrade_nats_id_before=$(docker inspect --format '{{.Id}}' "$NATS_CONTAINER")

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$UPGRADE_TABLE (id, name) VALUES (250, 'UpgradeBeforeSwap');
SQL

upgrade_before_position=""
for _ in $(seq 1 120); do
  upgrade_before_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$UPGRADE_TABLE' AND new_data->>'id' = '250'")
  if [ -n "$upgrade_before_position" ]; then
    break
  fi
  sleep 0.25
done
if [ -z "$upgrade_before_position" ]; then
  echo "run-integration-tests: Upgrade-Sicherheits-Rundlauf — Vorher-Datenstand (id=250) wurde nicht innerhalb der Zeitspanne über cdc.changes gelesen" >&2
  exit 1
fi

# Realer Container-Tausch (ADR-0064): --force-recreate stoppt und entfernt
# die bestehende Feed-Container-Instanz und legt eine neue aus demselben
# :dev-Image an (container_name bleibt cdc-test-feed); --no-deps lässt
# postgres/nats unberührt, die bereits laufen und healthy sind.
$COMPOSE up -d --force-recreate --no-deps pg-change-feed >/dev/null

upgrade_feed_id_after=$(docker inspect --format '{{.Id}}' "$FEED_CONTAINER")
if [ "$upgrade_feed_id_after" = "$upgrade_feed_id_before" ]; then
  echo "run-integration-tests: Upgrade-Sicherheits-Rundlauf — Feed-Container-ID nach --force-recreate unverändert ($upgrade_feed_id_after), kein realer Container-Tausch" >&2
  exit 1
fi

upgrade_pg_id_after=$(docker inspect --format '{{.Id}}' "$PG_CONTAINER")
upgrade_nats_id_after=$(docker inspect --format '{{.Id}}' "$NATS_CONTAINER")
if [ "$upgrade_pg_id_after" != "$upgrade_pg_id_before" ] || [ "$upgrade_nats_id_after" != "$upgrade_nats_id_before" ]; then
  echo "run-integration-tests: Upgrade-Sicherheits-Rundlauf — postgres/nats-Container-ID änderte sich trotz --no-deps (postgres: $upgrade_pg_id_before -> $upgrade_pg_id_after, nats: $upgrade_nats_id_before -> $upgrade_nats_id_after)" >&2
  exit 1
fi

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
  echo "run-integration-tests: Feed-Container meldet nach dem Upgrade-Sicherheits-Container-Tausch Health-Status ${health:-fehlt}, wollen healthy" >&2
  exit 1
fi

upgrade_before_still=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$UPGRADE_TABLE' AND new_data->>'id' = '250'")
if [ "$upgrade_before_still" != "1" ]; then
  echo "run-integration-tests: Upgrade-Sicherheits-Rundlauf — Datenstand vor dem Tausch (id=250) nach dem Container-Tausch nicht mehr identisch über cdc.changes lesbar (count=$upgrade_before_still)" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$UPGRADE_TABLE (id, name) VALUES (251, 'UpgradeAfterSwap');
SQL

upgrade_after_position=""
for _ in $(seq 1 120); do
  upgrade_after_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$UPGRADE_TABLE' AND new_data->>'id' = '251'")
  if [ -n "$upgrade_after_position" ]; then
    break
  fi
  sleep 0.25
done
if [ -z "$upgrade_after_position" ]; then
  echo "run-integration-tests: Upgrade-Sicherheits-Rundlauf — nach dem Container-Tausch eingefügte Zeile (id=251) wurde nicht innerhalb der Zeitspanne über cdc.changes gelesen — die Erfassung wurde nach dem Tausch nicht fortgesetzt" >&2
  exit 1
fi

echo "run-integration-tests: Upgrade-Sicherheits-Rundlauf (LH-QA-OPS-005, ADR-0064) belegt — realer Container-Tausch über \$COMPOSE up -d --force-recreate --no-deps ($upgrade_feed_id_before -> $upgrade_feed_id_after), postgres/nats unberührt (--no-deps), Feed-Container danach healthy, Datenstand vor dem Tausch (id=250, Position $upgrade_before_position) identisch lesbar, danach eingefügte Zeile (id=251) weiterhin erfasst (Position $upgrade_after_position)"

abdeckung_declare "Schema-Wiederanlauf nach der Spaltenentfernung" "LH-FA-SCH-003" "nach dem dauerhaften Ende des Erfassungspfads stellt der Runner einen sauberen Zustand wieder her: Replication-Slot neu angelegt, aktuelle Schema-Version nachgetragen, Feed-Container wieder healthy" "TestE2ESchemaChangeDropColumn (LH-FA-SCH-003, ADR-0063) belegt"

# TestE2ESchemaChangeDropColumn (LH-FA-SCH-003, ADR-0063 Supersedes
# ADR-0058 Entscheidung 1) läuft als eigener go-test-Aufruf, NACH allem,
# was den bislang unversehrten Feed-Container noch braucht (Lasttest-Beleg,
# Black-Box-CLI-Rundlauf, Retention-/NATS-/HTTP-Belege oben) und VOR
# TestE2ESchemaChangeIncompatibleTypeChange: eine real entfernte Spalte löst
# denselben relationOther/ErrIncompatibleSchemaChange-Pfad aus wie eine
# inkompatible Typänderung und beendet den Erfassungspfad des
# Feed-Containers ebenso dauerhaft (restart: "no", kein Neustart-Vertrag).
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_INTEGRATION_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test -v -run '^TestE2ESchemaChangeDropColumn$' ./test/integration/...

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "false" ]; then
  echo "run-integration-tests: Feed-Container lief nach TestE2ESchemaChangeDropColumn entgegen der Erwartung noch (der Test belegt gerade, dass eine reale Spaltenentfernung den Erfassungspfad dauerhaft beendet — ADR-0063)" >&2
  exit 1
fi

# Zwei Poison-Zustände trennen TestE2ESchemaChangeDropColumn von einem
# sauberen Neustart, beide real durch Testen dieses Mechanismus gefunden
# (ADR-0063 delegiert den konkreten Mechanismus an diesen Slice):
#
# 1. Der Replication-Slot: Der Prozess stirbt, bevor sein letzter
#    Standby-Status-Update-Zyklus die bereits verarbeitete Position
#    durabel im Slot bestätigt — ein bloßer `docker start` liest daher
#    nicht nur die eine noch offene Transaktion erneut, sondern auch die
#    bereits erfolgreich verarbeitete `ADD COLUMN removable`-Transaktion
#    ein zweites Mal, die dann (real getestet) erneut als kompatible
#    Erweiterung über der zwischenzeitlich korrigierten Spaltenform
#    hinweg registriert wird — und die anschließend wiederholte
#    `DROP COLUMN`-Transaktion triggert denselben Fehler ein zweites Mal.
#    Der Slot wird deshalb real verworfen: ein danach neu angelegter Slot
#    beginnt beim aktuellen WAL-Stand, ohne jede der beiden bereits
#    verarbeiteten Transaktionen erneut vorzulegen.
# 2. Die zuletzt registrierte Schema-Version (`cdc.schema_version`,
#    `cdc.table_schema`): Sie trägt nach `TestE2ESchemaChangeDropColumn`
#    weiterhin `removable`, weil genau diese Spalte real gelöscht wurde,
#    nachdem die kompatible Erweiterung sie bereits registriert hatte —
#    ohne Korrektur vergliche `TestE2ESchemaChangeIncompatibleTypeChange`s
#    eigene erste Änderung (id=10, ohne `removable`) gegen diese veraltete
#    Spaltenform und triggerte denselben `relationOther`-Fehler sofort
#    erneut, bevor die eigentliche Typänderungs-Prüfung dieser Funktion
#    überhaupt beginnt. Eine neue, nachgetragene Version mit der real
#    aktuellen Spaltenform (ohne `removable`) schließt diese Lücke, ohne
#    die bereits registrierte, von `TestE2ESchemaChangeDropColumn`s
#    Boundary-Assertion referenzierte Version `tbl-e2e-schema-v3`
#    anzutasten (Fremdschlüssel aus `cdc.change.schema_version`, die
#    Version bleibt für die historische Change von id=3 unverändert
#    gültig).
SCHEMA_TABLE_ID=tbl-e2e-schema

docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -c \
  "SELECT pg_drop_replication_slot('$SLOT') WHERE EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = '$SLOT');" >/dev/null

stale_schema_version_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT schema_version_id FROM cdc.schema_version WHERE source_table_id = '$SCHEMA_TABLE_ID' ORDER BY version DESC LIMIT 1")
next_schema_version=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT max(version) + 1 FROM cdc.schema_version WHERE source_table_id = '$SCHEMA_TABLE_ID'")
if [ -z "$stale_schema_version_id" ] || [ -z "$next_schema_version" ]; then
  echo "run-integration-tests: keine registrierte Schema-Version für $SCHEMA_TABLE_ID nach TestE2ESchemaChangeDropColumn gefunden" >&2
  exit 1
fi
next_schema_version_id="${SCHEMA_TABLE_ID}-v${next_schema_version}"
docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -c "
INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ('$next_schema_version_id', '$SCHEMA_TABLE_ID', $next_schema_version);
INSERT INTO cdc.table_schema (schema_version_id, ordinal_position, column_name, column_oid)
  SELECT '$next_schema_version_id', ordinal_position, column_name, column_oid
  FROM cdc.table_schema WHERE schema_version_id = '$stale_schema_version_id' AND column_name <> 'removable';
" >/dev/null

docker start "$FEED_CONTAINER" >/dev/null

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
  echo "run-integration-tests: Feed-Container meldet nach dem Neustart zwischen TestE2ESchemaChangeDropColumn und TestE2ESchemaChangeIncompatibleTypeChange Health-Status ${health:-fehlt}, wollen healthy" >&2
  exit 1
fi

echo "run-integration-tests: TestE2ESchemaChangeDropColumn (LH-FA-SCH-003, ADR-0063) belegt — reale Spaltenentfernung beendete den Erfassungspfad real (error_class=schema); Replication-Slot neu angelegt (kein Replay der beiden bereits verarbeiteten Transaktionen) und Schema-Version $next_schema_version_id nachgetragen (real aktuelle Spaltenform ohne removable), Feed-Container real neu gestartet und wieder healthy vor TestE2ESchemaChangeIncompatibleTypeChange"

# TestE2ESchemaChangeIncompatibleTypeChange (LH-FA-SCH-004
# Negative-Fall) läuft als eigener, letzter go-test-Aufruf: sie meldet
# eine nicht sicher als Obermenge erkennbare Typänderung sichtbar über die
# Fehlerklasse `schema` und beendet damit den Erfassungspfad des
# Feed-Containers dauerhaft (`restart: "no"`, kein Neustart-Vertrag,
# siehe Funktionskommentar). Alles, was den laufenden Container noch
# braucht (Lasttest-Beleg, Black-Box-CLI-Rundlauf oben), lief davor;
# TestE2ESchemaChangeDropColumn direkt davor beendete den Container
# ebenfalls bereits dauerhaft (ADR-0063) — der Neustart oben stellt einen
# sauberen, healthy Zustand wieder her, bevor diese Funktion ihren eigenen
# Mechanismus real prüft.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_INTEGRATION_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test -v -run '^TestE2ESchemaChangeIncompatibleTypeChange$' ./test/integration/...

# Zusammensetzung der E2E-Abdeckungstabelle: erst Go-Zeilen nach Quelldatei
# und Quellzeile, dann die Bash-Zeilen der deklarierten Phasen nach
# Runner-Zeile. Die Deklarations-Anker sind an dieser Stelle bereits alle
# gelaufen; eine fehlende Datei entsteht hier und bleibt sonst unangetastet.
abdeckung_schreiben "$ABDECKUNG_GO_ZEILEN"
echo "run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus $(( $(printf '%s\n' "$ABDECKUNG_GO_ZEILEN" | grep -c .) )) Go-Zeilen und ${#abdeckung_bash_zeilen[@]} Bash-Zeilen"
