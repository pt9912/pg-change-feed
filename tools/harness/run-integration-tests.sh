#!/usr/bin/env bash
# run-integration-tests — Compose-Integrationstest gegen die
# Compose-Umgebung (compose.yaml; LH-QA-POR-003), inhaltlich über den
# ursprünglichen MVP-Zuschnitt hinausgewachsen (Rollen-DSN-Trennung,
# Lasttest-Beleg, Black-Box-CLI-Rundlauf). Kette je Lauf:
#
# Compose-PostgreSQL frisch hochfahren → Schema-Rollout über d-migrate
# (ADR-0043).
#
# Rollen-DSN-Verifikation direkt gegen die Instanz (ADR-0047).
#
# Vorbedingungen der Aktivierung (Quell-Tabellen, Quelle-Zeile).
# Feed-Container als CDC-Runtime starten — seine Verdrahtung aktiviert die
# Tabellen über den EnableTable Use Case (ADR-0028), nicht über Seed-SQL.
# Toolchain-Container gegen das Compose-Netz.
#
# Der Feed-Container streamt dabei selbst — die Verdrahtung liegt am
# Composition Root (ADR-0026); der Test schreibt nur Quelländerungen
# und liest den Store.
#
# Am Ende der Kette steht ein Black-Box-Rundlauf
# (ADR-0030 E2E-Tier): `register-consumer`/`acknowledge-consumer` laufen
# dort ausschließlich als externer `docker exec`-Aufruf gegen den
# laufenden Feed-Container, über einen simulierten Container-Neustart
# hinweg — kein Go-Paket-Import interner Anwendungslogik in diesem
# Abschnitt.
#
# Danach ein CLI-Diagnose-Beleg (LH-FA-SST-003): derselbe externe
# `docker exec`-Zugriffsweg
# gegen den `diagnose`-Sondermodus, einmal im Normalbetrieb und einmal mit
# einem direkt in `cdc.process_heartbeat` geschriebenen Fehlerzustand
# (Boundary).
#
# Ergänzend ein Retention-Sichtbarkeits-Beleg (LH-FA-RET-005) über
# denselben `diagnose`-Aufruf: „kein
# Blocker" vor jeder Consumer-Bestätigung, dann ein real blockierender
# Consumer samt `cdc_storage_bytes`.
#
# Direkt davor ein kombinierter
# Retention-Lebenszyklus-Rundlauf (LH-FA-RET-002):
# eine eigene, isolierte Zeile und ein eigener Consumer durchlaufen real
# Blocker-Sichtbarkeit über `cdc.retention_blockers`, Bestätigung über die
# Position hinweg, die reale Abwesenheit jeder Zeile für `src-e2e` danach,
# die reale Löschung sowie die durchgehende numerische Lesbarkeit von
# `cdc_storage_bytes`.
#
# Ein Upgrade-Sicherheits-Rundlauf (LH-QA-OPS-005)
# folgt danach: ein realer Container-Tausch über
# `$COMPOSE up -d --force-recreate --no-deps pg-change-feed` ersetzt den
# Feed-Container durch eine neue Instanz desselben `:dev`-Images, während
# `postgres`/`nats` unberührt bleiben — Datenstand vor dem Tausch bleibt
# über `cdc.changes` identisch lesbar, eine danach eingefügte Zeile wird
# weiterhin erfasst.
#
# Die Backfill-Rundläufe (LH-FA-CAP-009: Happy Path,
# Schema-Version, Startposition eines frisch registrierten Consumers,
# Boundary, Regelstand mit rename_column, Ausschluss und nicht anwendbarer
# Regel, Replay-Invariante, DDL-Fenster, Negative) laufen vor dem
# Upgrade-Sicherheits-Rundlauf;
# die Haltepunkte der Backfill-Phasen beschreibt der Kopf des Abschnitts
# `Backfill-Rundläufe`.
#
# Die Leerlauf-Bestätigung (ADR-0120: WAL über der Fehlerschwelle, der
# Feed-Container läuft weiter) läuft ebenfalls vor dem
# Upgrade-Sicherheits-Rundlauf.
#
# Ihre Gegenseite: bei gehaltener Persistierung beendet die
# Fehlerschwelle den Feed-Container (ADR-0049).
#
# Die Transformations-Rundläufe (LH-FA-CFG-007: die
# Form der Regeln auf allen fünf Zustellwegen, Neustart, Ausschluss mit
# Regel) laufen nach der Leerlauf-Bestätigung, ebenfalls vor dem
# Upgrade-Sicherheits-Rundlauf.
#
# Nichtanwendbarkeit einer Regel und ihre
# Abhilfe (ADR-0112) laufen als eigener,
# letzter Rundlauf nach TestE2ESchemaChangeIncompatibleTypeChange:
# eine kompatible Spalten-Erweiterung, deren Name den Zielnamen einer aktiven
# rename_column-Regel trifft, beendet den Erfassungspfad sichtbar; nach
# `cdc.remove_transformation` und einem Neustart erscheint die zuvor nicht
# bestätigte Transaktion über `cdc.changes`.
#
# Test-Daten bleiben im Container (kein Volume in den Arbeitsbaum);
# Compose-Container und -Netz werden in jedem Ausgang abgeräumt, das
# Modul-Cache-Volume bleibt als Vorbereitung für netzlose `make test`-Läufe
# bestehen.
#
# Eine Prüfung eines Container-Logs liest `docker logs … | grep -F "…"
# >/dev/null` ohne `-q`: `grep` liest bis zum Dateiende, `docker logs`
# schreibt vollständig, und `pipefail` wertet den Treffer-Status der Pipeline.
set -euo pipefail
# Vor dem cd unten einfangen: BASH_SOURCE[0] ist relativ zum AUFRUF-Verzeichnis
# (z. B. `../tools/harness/run-integration-tests.sh` bei Aufruf aus einem
# Unterverzeichnis) — nach dem cd zur Repo-Wurzel würde derselbe relative
# Ausdruck etwas anderes bezeichnen. `pwd -P` löst Symlinks physisch auf
# (POSIX, GNU wie BSD/macOS identisch).
EIGENER_PFAD_ABS="$(cd "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)/$(basename -- "${BASH_SOURCE[0]}")"
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

# Serverweiter NATS-Verbindungs-Token (ADR-0100 Teilfrage 4/5): derselbe Wert
# wie compose.yaml (nats-Service --auth, pg-change-feed CDC_NATS_STREAM_TOKEN)
# — jede NATS-Verbindung in diesem Lauf braucht ihn jetzt, auch die bislang
# anonyme Wecksignal-Verbindung (natssub).
NATS_STREAM_TOKEN=e2e-nats-stream-token

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
# Portabel statt `realpath --relative-to` (GNU-only, BSD/macOS-realpath kennt
# das Flag nicht): EIGENER_PFAD_ABS wurde oben VOR dem cd zur Repo-Wurzel
# eingefangen (Zeile ~50) — der aktuelle `pwd -P` ist jetzt die Repo-Wurzel,
# der Präfix-Abzug liefert den repo-relativen Pfad dieses Skripts.
ABDECKUNG_QUELLE="${EIGENER_PFAD_ABS#$(pwd -P)/}"
ABDECKUNG_GO_ZEILEN=""
ABDECKUNG_KOPF='# E2E-Abdeckung je Spec-Kennung

Erzeugt von `make test-integration` über `tools/harness/run-integration-tests.sh`:
die Go-Zeilen leitet das Testpaket aus seinem eigenen Quelltext ab
(`TestAbdeckungstabelleZeilen` in `test/integration/integration_test.go`; ob jede Funktion in einem Testmuster dieses Runners steht, prüft `TestRunnerFuehrtJedeE2EFunktionAus`),
die Bash-Zeilen deklariert jede Phase des Runners an Ort und Stelle über
einen Anker. Diese Datei ist eine **stabile Abdeckungs-Deklaration**, kein
Lauf-Beleg: der Runner schreibt sie nur bei inhaltlicher Abweichung. Sie
ändert sich mit den Nachweis-Deklarationen und mit dem Ort ihrer Quellen —
die Spalte `Ort` bindet an `Datei:Zeile`, jede Einfügung oberhalb einer
`func TestE2E*` oder einer Runner-Phase verschiebt die Zeilennummer und
damit diese Tabelle, ohne dass sich an der Abdeckung etwas ändert. Die
Beschreibungsspalte trägt keine Kennungen; die Aussage eines Nachweises steht
über die Spalte `Ort` an ihrer Quelle.
'
abdeckung_bash_zeilen=()

# abdeckung_declare <Nachweis> <Kennungen> <Kurzbeschreibung> <Anker>
# Der Anker ist ein wörtlicher Ausschnitt aus der Phase dahinter — der
# Runner sucht ihn ab der Deklarations-Zeile im eigenen Quelltext und trägt
# die gefundene Zeile als `Ort` ein; ein Anker, der in keiner späteren Zeile
# mehr steht, beendet den Lauf. Die vier Argumente sind Shell-Wörter: ein
# Backtick-Zeichen in Kurzbeschreibung oder Anker gehört escaped (`\``),
# sonst ersetzt die Shell es samt Inhalt durch die Ausgabe eines Kommandos —
# ohne Fehlermeldung, wenn dieses Kommando existiert.
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
    printf '| %s | %s | `%s` | `%s:%s` |\n' "$spalte" "$kurzbeschreibung" "$nachweis" "$quelldatei" "$zeile"
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
    printf '| Spec-Kennung | Kurzbeschreibung | Nachweis | Ort |\n'
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
  rm -rf "${WAL_TMP:-}"
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

# Schema-Rollout über d-migrate vor jedem E2E-Lauf (ADR-0043): Pflicht-Report
# und Rollback-Artefakt liegen in SCHEMA_ARTEFACT_DIR, der Arbeitsbaum bleibt
# unberührt.
make schema-rollout SCHEMA_TARGET="db:$DSN" SCHEMA_ROLLOUT_NETWORK="$NETWORK"

abdeckung_declare "Rollen-DSN-Verifikation" "LH-QA-SEC-001,LH-QA-SEC-002,LH-QA-SEC-003" "die drei Gruppenrollen real gegeneinander geprüft: ein Reader-Login scheitert am schreibenden Aufruf, eine Replication-Verbindung ohne REPLICATION-Attribut scheitert, dieselbe Verbindung mit der Capture-Rolle gelingt" "Rollen-DSN-Verifikation belegt — cdc_reader-Login-Identität schreibender Zugriff abgelehnt"

# Rollen-DSN-Verifikation gegen den Compose-Stack (ADR-0047,
# BEO-PGC/rollen-test-abdeckungsluecken): Die laufende
# CDC_CAPTURE_DSN/CDC_ADMIN_DSN/CDC_READER_DSN-Verdrahtung des
# Feed-Containers verbindet sich testbedingt mit allen drei DSNs über den
# Superuser postgres (siehe compose.yaml-Kommentar) — die drei
# Gruppenrollen selbst liegen seit dem Schema-Rollout oben aber real in
# dieser Instanz vor (nacharbeit-roles.sql). Dieser Abschnitt belegt die
# Rollentrennung direkt gegen sie, mit derselben Fixture-Disziplin wie
# internal/bootstrap/roles_wiring_test.go::newTestLoginRole: anmeldefähige
# Test-Identitäten je Rolle, angelegt und am Ende dieses Abschnitts wieder
# entfernt — kein Bestandteil des Rollen-DDL
# (tools/schema/nacharbeit-roles.sql).
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
# Metadaten-Registrierung ist Vorbedingung der Aktivierung
# (Fremdschlüssel). Die Aktivierung selbst trägt die Verdrahtung des
# Feed-Containers: die Bindungen aus CDC_TABLES laufen als
# EnableTable-Aufrufe — Publication, Bindungs- und
# Schema-Version-Zeilen entstehen dort, nicht hier. Die REPLICA IDENTITY
# trägt der Runner für die volle Alt-Bild-Prüfung; sie
# bleibt vom Aktivieren unberührt.
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE public.feed_e2e_flow (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_full (id int PRIMARY KEY, name text);
ALTER TABLE public.feed_e2e_full REPLICA IDENTITY FULL;
CREATE TABLE public.feed_e2e_idle (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_schema (id int PRIMARY KEY, name text, amount text);
CREATE TABLE public.feed_e2e_sql_admin (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_walsender_timing (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_startgrenze (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_transform_abhilfe (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_api_enable_grpc (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_api_enable_http (id int PRIMARY KEY, name text);
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

# Compose-Healthcheck-Vertrag (LH-FA-ADM-002):
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
# (`LH-FA-SCH-003`) beendet
# den Erfassungspfad ebenfalls dauerhaft (derselbe
# `relationOther`/`ErrIncompatibleSchemaChange`-Pfad wie eine inkompatible
# Typänderung) und läuft deshalb NICHT hier, sondern als eigener Aufruf nach
# der Container-Ende-Grenze, mit einem expliziten Neustart davor — siehe dort.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_INTEGRATION_DSN="$DSN" \
  -e CDC_INTEGRATION_HTTP_URL="$HTTP_BASE_URL" \
  -e CDC_INTEGRATION_HTTP_READER_TOKEN="$HTTP_TOKEN_READER" \
  "$TOOLCHAIN_IMAGE" go test -v \
  -run '^(TestE2ECaptureFlow|TestE2EUpdateOldImageWithFullReplicaIdentity|TestE2EChangesViewMatchesReadChanges|TestE2ERetentionBlockersViewShowsFurthestBehindConsumer|TestE2EMetricsCarriesStorageBytes|TestE2EActivationState|TestE2EActiveTablesViewMatchesActivationState|TestE2EDisableRetainedState|TestE2ESchemaChangeAddColumn|TestE2EChangeTableMetadataExtensibility|TestE2EHeartbeatHealthy|TestE2ETransformationRulesShapeBothImages|TestE2ETransformationConflictsFailWithSpecText|TestE2ERoutingSelectsTargetsOverChangesView|TestE2ERoutingConflictsFailWithSpecText|TestE2ERoutingReplayKeepsLabelAndBackfillCarriesIt|TestE2ERoutingDeleteWithoutFullReplicaIdentity)$' \
  ./test/integration/...

# Go-Hälfte der E2E-Abdeckungstabelle: derselbe Testlauf führt den Erzeuger
# aus (eigener `-run`-Aufruf, dieselbe Toolchain) und legt seinen Zeilensatz
# in ABDECKUNG_GO_ZEILEN ab; eine `func TestE2E*` ohne Spec-Kennung endet
# hier sichtbar, bevor die Datei geschrieben wird.
abdeckung_go_zeilen_lesen

abdeckung_declare "Spaltenausschluss-Rundlauf" "LH-FA-CFG-005,LH-QA-SEC-004" "eine nicht gelistete Tabelle wird über den SQL-Antragsweg aktiviert, der Spaltenausschluss real verarbeitet — die danach erfasste Change trägt den Schlüssel nicht mehr im Row Image, die davor erfasste bleibt mit ihrem Wert unverändert lesbar" "Spaltenausschluss-Rundlauf (LH-FA-CFG-005 Happy Path, ADR-0059) belegt"

# Spaltenausschluss-Rundlauf (ADR-0059): der reale Pfad
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
# in der Row-Image-Konstruktion (Teilfrage 3; cdc.changes ist eine
# reine Projektion über cdc.change) und schreibt persistierte Changes nicht
# rückwirkend um — sie adressiert „künftige Changes von t".
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

abdeckung_declare "Lasttest-Beleg cdc_capture_lag" "LH-FA-ADM-004,LH-QA-PER-004" "cdc_capture_lag bildet den Abstand zwischen Quelländerung und CDC-Verfügbarkeit ab: eine durch eine pausierte CDC-Runtime künstlich verzögerte Transaktion liegt über der ungehinderten — dieselbe Latenz-Messmethode, die \`LH-QA-PER-004\` über \`LH-FA-ADM-004\` verlangt" "Lasttest-Beleg cdc_capture_lag — Baseline"

# Lasttest-Beleg (LH-FA-ADM-004): cdc_capture_lag
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

abdeckung_declare "Black-Box-CLI-Rundlauf" "LH-QA-POR-003,LH-FA-CON-001,LH-FA-CON-003,LH-FA-CON-004,LH-FA-CON-005" "register-consumer und acknowledge-consumer laufen ausschließlich als externe docker exec-Aufrufe gegen den Produktions-Binary, über einen simulierten Container-Neustart hinweg — Registrierung, Positions-Persistierung (cdc.consumer_position gegen die bestätigte Position gehalten), Bestätigung und Fortsetzen ab der bestätigten Position nach dem Neustart real belegt" "Black-Box-CLI-Rundlauf belegt — register-consumer/acknowledge-consumer extern"

# Black-Box-CLI-Rundlauf (ADR-0030 E2E-Tier): anders als der
# Go-Testlauf oben (`go test ./test/integration/...`, der intern gegen
# `postgresstorage`/`bootstrap` läuft) ruft dieser Abschnitt
# `register-consumer`/`acknowledge-consumer` ausschließlich als externen
# Prozess gegen den laufenden, containerisierten Produktions-Binary auf
# (`docker exec … /pg-change-feed …`, distroless — kein Shell im
# Feed-Container, daher kein `sh -c`-Umweg nötig). Lesen bleibt über den
# bestehenden externen Lesezugriffsweg `cdc.changes` — die
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

# Retention-Lebenszyklus-Rundlauf (kombiniert, LH-FA-RET-002): eine
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
# Entfernung der Position selbst (LH-FA-RET-003).
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
# (LH-FA-SST-003): an dieser Stelle hat noch kein
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
if ! printf '%s' "$diagnose_noblocker_output" | grep -qF "Blockierender Consumer: kein Blocker"; then
  echo "run-integration-tests: diagnose-Ausgabe zeigt vor jeder Consumer-Bestätigung nicht 'kein Blocker' (LH-FA-RET-005): $diagnose_noblocker_output" >&2
  exit 1
fi
if ! printf '%s' "$diagnose_noblocker_output" | grep -qE "Speicherverbrauch cdc_storage_bytes: [0-9]+ Bytes"; then
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

# CLI-Diagnose-Beleg (LH-FA-SST-003):
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
# Blocker: `cdc.retention_blockers` wählt je Quelle den
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
  "Betriebsstatus:" \
  "Fehlerzustand: keiner (Normalbetrieb)" \
  "CDC-Abstand cdc_capture_lag:" \
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
if ! printf '%s' "$diagnose_output" | grep -qE "Blockierender Consumer: .*\($CLI_CONSUMER\), bestätigte Position [0-9]+, Rückstand [0-9]+"; then
  echo "run-integration-tests: diagnose-Ausgabe (Normalbetrieb) zeigt $CLI_CONSUMER nicht als real blockierenden Consumer (LH-FA-RET-005): $diagnose_output" >&2
  exit 1
fi
if ! printf '%s' "$diagnose_output" | grep -qE "Speicherverbrauch cdc_storage_bytes: [0-9]+ Bytes"; then
  echo "run-integration-tests: diagnose-Ausgabe (Normalbetrieb) trägt keine numerische cdc_storage_bytes-Zeile (LH-FA-RET-006): $diagnose_output" >&2
  exit 1
fi

echo "run-integration-tests: CLI-Diagnose-Beleg (Normalbetrieb) — alle vier Signale (LH-FA-ADM-002…005) sowie der reale Blocker $CLI_CONSUMER und cdc_storage_bytes (LH-FA-RET-005/006) in der diagnose-Ausgabe sichtbar"

abdeckung_declare "CLI-Diagnose-Beleg (Fehlerzustand)" "LH-FA-ADM-003,LH-QA-REL-003" "ein direkt in die Heartbeat-Projektion geschriebener Fehlerzustand ist in der diagnose-Ausgabe von Normalbetrieb unterscheidbar, ohne den laufenden Feed-Container zu beenden — derselbe Fehlerinjektionstest, den \`LH-QA-REL-003\`s Messmethode über \`LH-FA-ADM-003\` verlangt" "CLI-Diagnose-Beleg (Fehlerzustand) —"

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
  if [ "$diagnose_error_status" -eq 0 ] && printf '%s' "$diagnose_error_output" | grep -qF "Fehlerzustand: schema"; then
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

abdeckung_declare "SQL-Administration Live-Reload (enable)" "LH-FA-ADM-001,LH-FA-CFG-001,LH-FA-CFG-006" "eine bewusst nicht in der Bindungsliste geführte Tabelle wird über den SQL-Antrag aktiviert und vom bereits laufenden Feed-Container ohne Neustart erfasst — ihre eigene DDL (Spalten, Trigger) bleibt dabei real unverändert, keine Anwendungscode-Anpassung nötig" "SQL-Administration Live-Reload-Beleg (enable)"

abdeckung_declare "SQL-Administration Live-Reload (disable)" "LH-FA-CFG-002" "derselbe Antrags-Weg spiegelbildlich: die Deaktivierung beendet die Erfassung dieser Tabelle, der Prozess läuft weiter" "SQL-Administration Live-Reload-Beleg (disable)"

# SQL-Administration Live-Reload-Beleg (ADR-0050): `feed_e2e_sql_admin` ist bewusst NICHT Teil von
# CDC_TABLES (compose.yaml) — ihre Aktivierung/Deaktivierung läuft
# ausschließlich über die Antrags-Queue (`cdc.enable_table`/
# `cdc.disable_table`), verarbeitet von der Administrations-Goroutine des
# bereits laufenden Feed-Containers (kein Neustart, kein `docker restart`
# wie beim Black-Box-CLI-Rundlauf oben). `SELECT cdc.enable_table(...)`
# bestätigt nur „beantragt" — der Poll unten wartet auf
# `status = 'applied'`, bevor die reale Erfassungswirkung geprüft wird.
ADMIN_TABLE=feed_e2e_sql_admin

# LH-FA-CFG-006-Beleg: `cdc.enable_table` schreibt nur einen Antrag nach
# `cdc.administration_request`; die Administrations-Goroutine wendet ihn
# über `TableActivationAdapter.Publish` an — ausschließlich `ALTER
# PUBLICATION ... ADD TABLE` gegen die Publication, nie eine DDL-Änderung
# an der Quelltabelle selbst. Fingerabdruck aus Spaltenliste und
# Trigger-Anzahl vor dem Antrag, gegen denselben Fingerabdruck nach
# `status = 'applied'` unten geprüft.
admin_table_ddl_before=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT string_agg(column_name || ':' || data_type, ',' ORDER BY ordinal_position) FROM information_schema.columns WHERE table_schema = 'public' AND table_name = '$ADMIN_TABLE'")
admin_table_triggers_before=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM pg_trigger WHERE tgrelid = 'public.$ADMIN_TABLE'::regclass AND NOT tgisinternal")

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

admin_table_ddl_after=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT string_agg(column_name || ':' || data_type, ',' ORDER BY ordinal_position) FROM information_schema.columns WHERE table_schema = 'public' AND table_name = '$ADMIN_TABLE'")
admin_table_triggers_after=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM pg_trigger WHERE tgrelid = 'public.$ADMIN_TABLE'::regclass AND NOT tgisinternal")
if [ "$admin_table_ddl_before" != "$admin_table_ddl_after" ] || [ "$admin_table_triggers_before" != "$admin_table_triggers_after" ]; then
  echo "run-integration-tests: cdc.enable_table($ADMIN_TABLE) veränderte die DDL der Quelltabelle selbst (LH-FA-CFG-006 verletzt)" >&2
  echo "  vorher: columns=$admin_table_ddl_before triggers=$admin_table_triggers_before" >&2
  echo "  nachher: columns=$admin_table_ddl_after triggers=$admin_table_triggers_after" >&2
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

echo "run-integration-tests: SQL-Administration Live-Reload-Beleg (enable) — cdc.enable_table($ADMIN_TABLE) ohne Neustart verarbeitet, Änderung id=1 real erfasst, DDL der Quelltabelle unverändert"

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
# LH-FA-CFG-002, der Architect-Verdikt zur Walsender-Wirksamkeit):
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

# Retention-Beleg (ADR-0014): der Hintergrundzug
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
# erhalten bleibt, solange ein Consumer zurückhängt. Erst
# die zweite Bestätigung beider Consumer über die neue Position hinweg
# gibt beide Zeilen aus Consumer-Sicht frei; der zweite Poll belegt, dass
# danach nur die zurückdatierte Zeile real entfernt wird,
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

# NATS-Happy-Path-Beleg (ADR-0055): compose.yaml
# verdrahtet den Feed-Container mit CDC_NATS_URL=nats://nats:4222 (siehe
# dortigen Kommentar). Ein eigener Wegwerf-Testclient
# (tools/harness/natssub/main.go, per `go run` im Toolchain-Container)
# abonniert das tabellen-granulare Subjekt
# cdc.changes.src-e2e.public.feed_e2e_full real, BEVOR die auslösende
# Change entsteht — Core NATS liefert nichts nach (ADR-0055,
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

docker rm -fv "$NATS_SUBSCRIBER_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$NATS_SUBSCRIBER_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/natssub "nats://nats:4222" "$NATS_SUBJECT" "$NATS_STREAM_TOKEN" >/dev/null

nats_subscriber_ready=0
for _ in $(seq 1 60); do
  if docker logs "$NATS_SUBSCRIBER_CONTAINER" 2>/dev/null | grep -F "READY" >/dev/null; then
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
  docker rm -fv "$NATS_SUBSCRIBER_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.feed_e2e_full (id, name) VALUES (230, 'NatsHappyPath');
SQL

nats_signal_received=0
for _ in $(seq 1 60); do
  if docker logs "$NATS_SUBSCRIBER_CONTAINER" 2>/dev/null | grep -F "RECEIVED" >/dev/null; then
    nats_signal_received=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$NATS_SUBSCRIBER_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 0.5
done
nats_subscriber_output=$(docker logs "$NATS_SUBSCRIBER_CONTAINER" 2>&1 || true)
docker rm -fv "$NATS_SUBSCRIBER_CONTAINER" >/dev/null 2>&1 || true
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

# NATS-Boundary-Beleg (ADR-0055): Gegenstück zum
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

# NATS-Negative-Beleg — Reconnect-Nachholen (ADR-0055): Gegenstück zum Boundary-Beleg oben (dort: nie abonniert
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

docker rm -fv "$NATS_RECONNECT_BEFORE_CONTAINER" "$NATS_RECONNECT_AFTER_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$NATS_RECONNECT_BEFORE_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/natssub "nats://nats:4222" "$NATS_SUBJECT" "$NATS_STREAM_TOKEN" >/dev/null

nats_reconnect_before_ready=0
for _ in $(seq 1 60); do
  if docker logs "$NATS_RECONNECT_BEFORE_CONTAINER" 2>/dev/null | grep -F "READY" >/dev/null; then
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
  docker rm -fv "$NATS_RECONNECT_BEFORE_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

nats_subs_before_disconnect=$(docker exec "$NATS_CONTAINER" wget -q -O - "http://localhost:8222/subsz?subs=1")
if ! echo "$nats_subs_before_disconnect" | grep -qF "\"subject\": \"$NATS_SUBJECT\""; then
  echo "run-integration-tests: NATS-Negative-Beleg — Test-Subscriber ($NATS_SUBJECT) war vor der Trennung entgegen der Erwartung nicht als Abonnent beim Broker gelistet: $nats_subs_before_disconnect" >&2
  docker rm -fv "$NATS_RECONNECT_BEFORE_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

docker network disconnect "$NETWORK" "$NATS_RECONNECT_BEFORE_CONTAINER"

subscriber_networks_after_disconnect=$(docker inspect "$NATS_RECONNECT_BEFORE_CONTAINER" --format '{{json .NetworkSettings.Networks}}')
if echo "$subscriber_networks_after_disconnect" | grep -qF "\"$NETWORK\""; then
  echo "run-integration-tests: NATS-Negative-Beleg — Test-Subscriber-Container trägt laut docker inspect nach dem Trennungsversuch noch die Netzbindung $NETWORK (reale Trennung nicht hergestellt): $subscriber_networks_after_disconnect" >&2
  docker rm -fv "$NATS_RECONNECT_BEFORE_CONTAINER" >/dev/null 2>&1 || true
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
docker rm -fv "$NATS_RECONNECT_BEFORE_CONTAINER" >/dev/null 2>&1 || true
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
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/natssub "nats://nats:4222" "$NATS_SUBJECT" "$NATS_STREAM_TOKEN" >/dev/null

nats_reconnect_after_ready=0
for _ in $(seq 1 60); do
  if docker logs "$NATS_RECONNECT_AFTER_CONTAINER" 2>/dev/null | grep -F "READY" >/dev/null; then
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
  docker rm -fv "$NATS_RECONNECT_AFTER_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.feed_e2e_full (id, name) VALUES (241, 'NatsReconnectResumed');
SQL

nats_reconnect_signal_resumed=0
for _ in $(seq 1 60); do
  if docker logs "$NATS_RECONNECT_AFTER_CONTAINER" 2>/dev/null | grep -F "RECEIVED" >/dev/null; then
    nats_reconnect_signal_resumed=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$NATS_RECONNECT_AFTER_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 0.5
done
nats_reconnect_after_output=$(docker logs "$NATS_RECONNECT_AFTER_CONTAINER" 2>&1 || true)
docker rm -fv "$NATS_RECONNECT_AFTER_CONTAINER" >/dev/null 2>&1 || true
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

abdeckung_declare "HTTP-API-Rundlauf" "LH-FA-SST-005,LH-FA-SST-006,LH-FA-REA-001" "ein Wegwerf-Client ruft RegisterConsumer mit dem admin-Token und ListTables mit dem reader-Token real per HTTP gegen den laufenden Feed-Container auf und liest zusätzlich Changes über \`GET /changes\` mit dem reader-Token; die Registrierung wird gegen cdc.consumer bestätigt, der gelesene Change gegen cdc.changes — die spätere API, deren Ermöglichung \`LH-FA-SST-005\` forderte, ohne das interne CDC-Modell zu verändern; der Lesezugriff trägt einen echten Bereich \`[from, to)\` zweier Positionen (\`LH-FA-REA-001\`)" "GET /changes real per HTTP mit reader-Token (die eigens eingefügte Zeile"

abdeckung_declare "HTTP-Diagnose-Querabgleich" "LH-FA-SST-006,LH-FA-SST-003,LH-FA-ADM-002,LH-FA-ADM-003" "GET /diagnose liefert denselben Betriebsstatus (Lebenszeichen bekannt, Fehlerzustand) wie eine kontemporäre docker exec diagnose-Ausgabe — ein echter Wertabgleich zwischen zwei Zugriffswegen, nicht nur die Feld-Präsenz der Antwort" "HTTP-Diagnose-Querabgleich (ADR-0132) belegt"

# HTTP-API-Rundlauf (LH-FA-SST-006): ein Wegwerf-Client
# (tools/harness/httpclient) ruft RegisterConsumer mit dem admin-Token,
# ListTables mit dem reader-Token und `GET /changes` mit dem reader-Token
# real per HTTP gegen den laufenden Feed-Container auf — ein echter
# Netzwerk-Request, kein `docker exec` und kein Mock. `feed_e2e_full` bleibt
# über den ganzen Lauf aktiviert (siehe Lasttest-Beleg oben), ihr Auftreten
# in der ListTables-Antwort belegt einen realen, nicht leeren Rückgabewert.
#
# Der Changes-Lese-Beleg läuft spät im Ablauf (nach dem NATS-Negative-Beleg);
# `feed_e2e_full` trägt zu diesem Zeitpunkt bereits erfasste Changes, die der
# Client als Bestand lesen kann. Damit der Beleg eine bestimmte, unmittelbar
# zuvor erfasste Änderung prüft, fügt der Runner eine eigene Zeile ein und
# wartet ihre Erfassung über cdc.changes ab, BEVOR der Client sie über den
# Endpunkt liest; ihre commit_position geht als Bereich `[from, to)` in den
# Aufruf ein (ADR-0081 Teilfrage 2), ihre change_id wird unabhängig gegen
# denselben Lesezugriffsweg gehalten.
#
# Muss vor der abschließenden Schema-Negative-Testfunktion unten laufen,
# weil diese den Feed-Container dauerhaft beendet.
HTTP_CONSUMER=http-e2e-consumer
HTTP_READ_SCHEMA=public
HTTP_READ_TABLE=feed_e2e_full
HTTP_READ_SENTINEL=HttpChangesReadE2ESentinel
HTTP_READ_ID=285
HTTP_READ_LIMIT=100

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$HTTP_READ_TABLE (id, name) VALUES ($HTTP_READ_ID, '$HTTP_READ_SENTINEL');
SQL

http_read_position=""
for _ in $(seq 1 60); do
  http_read_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$HTTP_READ_TABLE' AND new_data->>'id' = '$HTTP_READ_ID'" 2>/dev/null || true)
  if [ -n "$http_read_position" ]; then
    break
  fi
  sleep 0.5
done
if [ -z "$http_read_position" ]; then
  echo "run-integration-tests: Changes-Lese-Beleg — die eigens eingefügte Zeile (id=$HTTP_READ_ID, $HTTP_READ_TABLE) wurde nicht innerhalb der Zeitspanne über cdc.changes erfasst" >&2
  exit 1
fi
http_read_to=$((http_read_position + 1))

# Kontemporäre CLI-Diagnose-Baseline (docker exec, unmittelbar vor dem
# HTTP-Aufruf) — der Querabgleich unten hält den über GET /diagnose
# gelesenen Betriebsstatus gegen diese Ausgabe, nicht nur gegen die
# Feld-Präsenz der Antwort (ADR-0132).
set +e
http_diagnose_cli_output=$(exec_feed diagnose)
http_diagnose_cli_status=$?
set -e
if [ "$http_diagnose_cli_status" -ne 0 ]; then
  echo "run-integration-tests: HTTP-Diagnose-Querabgleich — CLI-Diagnose-Baseline (docker exec diagnose) endete mit Ausgang $http_diagnose_cli_status: $http_diagnose_cli_output" >&2
  exit 1
fi

# `set +e` um den Client-Aufruf: der Ausgang des Clients wird von der
# `http_status`-Prüfung ausgewertet und über die folgende Fehlerzeile
# gemeldet; unter `set -e` beendet ein fehlgeschlagenes `docker run` schon
# die Zuweisung selbst und die Phase endet rot ohne Ausgabe (AGENTS.md §3.9:
# Exit-Code direkt ausgewertet, hier mit sichtbarem Grund).
set +e
http_output=$(docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/httpclient \
  "$HTTP_BASE_URL" "$HTTP_TOKEN_ADMIN" "$HTTP_TOKEN_READER" "$HTTP_CONSUMER" "HTTP E2E Consumer" src-e2e pub_pgc_e2e \
  "$HTTP_READ_SCHEMA" "$HTTP_READ_TABLE" "$http_read_position" "$http_read_to" "$HTTP_READ_LIMIT" 2>&1)
http_status=$?
set -e
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
if ! printf '%s' "$http_output" | grep -qE '^READ changes=[0-9]+ table=feed_e2e_full schema=public '; then
  echo "run-integration-tests: HTTP-API-Rundlauf — keine READ-Zeile des Changes-Lese-Aufrufs (reader-Token, GET /changes, Bereich [$http_read_position,$http_read_to)): $http_output" >&2
  exit 1
fi
if ! printf '%s' "$http_output" | grep -qF "$HTTP_READ_SENTINEL"; then
  echo "run-integration-tests: HTTP-API-Rundlauf — die READ-Zeile trägt die eigens eingefügte Zeile (id=$HTTP_READ_ID, $HTTP_READ_SENTINEL) nicht: $http_output" >&2
  exit 1
fi
if ! printf '%s' "$http_output" | grep -qF "DIAGNOSED body="; then
  echo "run-integration-tests: HTTP-API-Rundlauf — keine DIAGNOSED-Zeile (reader-Token, GET /diagnose): $http_output" >&2
  exit 1
fi
if ! printf '%s' "$http_output" | grep -qE 'DIAGNOSED body=.*"heartbeat_age_seconds"'; then
  echo "run-integration-tests: HTTP-API-Rundlauf — die DIAGNOSED-Zeile trägt kein heartbeat_age_seconds-Feld: $http_output" >&2
  exit 1
fi

# Querabgleich (ADR-0132): dieselbe Quelle über zwei unabhängige
# Zugriffswege — der CLI-Sondermodus (Baseline oben) und GET /diagnose —
# tragen an diesem Zeitpunkt denselben Betriebsstatus. Verglichen werden
# zwei robuste, boolesche Ableitungen statt roher Zeitstempel/Sekundenwerte
# (die zwischen zwei sequenziellen Aufrufen nie exakt gleich wären): ob ein
# Lebenszeichen bekannt ist, und ob ein Fehlerzustand vorliegt.
if printf '%s' "$http_diagnose_cli_output" | grep -qF "Betriebsstatus: Lebenszeichen vor"; then
  http_diagnose_cli_known=1
else
  http_diagnose_cli_known=0
fi
if printf '%s' "$http_output" | grep -qE 'DIAGNOSED body=.*"heartbeat_age_seconds":null'; then
  http_diagnose_http_known=0
else
  http_diagnose_http_known=1
fi
if [ "$http_diagnose_cli_known" != "$http_diagnose_http_known" ]; then
  echo "run-integration-tests: HTTP-Diagnose-Querabgleich — Betriebsstatus (Lebenszeichen bekannt) widersprüchlich zwischen CLI (known=$http_diagnose_cli_known) und GET /diagnose (known=$http_diagnose_http_known): CLI=$http_diagnose_cli_output HTTP=$http_output" >&2
  exit 1
fi

if printf '%s' "$http_diagnose_cli_output" | grep -qF "Fehlerzustand: keiner (Normalbetrieb)"; then
  http_diagnose_cli_no_error=1
else
  http_diagnose_cli_no_error=0
fi
if printf '%s' "$http_output" | grep -qE 'DIAGNOSED body=.*"error_class":null'; then
  http_diagnose_http_no_error=1
else
  http_diagnose_http_no_error=0
fi
if [ "$http_diagnose_cli_no_error" != "$http_diagnose_http_no_error" ]; then
  echo "run-integration-tests: HTTP-Diagnose-Querabgleich — Fehlerzustand widersprüchlich zwischen CLI (kein Fehler=$http_diagnose_cli_no_error) und GET /diagnose (kein Fehler=$http_diagnose_http_no_error): CLI=$http_diagnose_cli_output HTTP=$http_output" >&2
  exit 1
fi

echo "run-integration-tests: HTTP-Diagnose-Querabgleich (ADR-0132) belegt — GET /diagnose (reader-Token) und eine kontemporäre docker exec diagnose-Ausgabe tragen denselben Betriebsstatus: Lebenszeichen bekannt=$http_diagnose_http_known, kein Fehlerzustand=$http_diagnose_http_no_error"

registered_via_http=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT consumer_id FROM cdc.consumer WHERE consumer_id = '$HTTP_CONSUMER'")
if [ "$registered_via_http" != "$HTTP_CONSUMER" ]; then
  echo "run-integration-tests: cdc.consumer trägt $HTTP_CONSUMER nicht nach dem realen HTTP-RegisterConsumer-Aufruf" >&2
  exit 1
fi

# Unabhängiger SQL-Beleg, dass genau die über den Endpunkt gelesene Änderung
# real erfasst wurde: die change_id der sentinel-tragenden READ-Zeile steht
# über den bestehenden Lesezugriffsweg in cdc.changes — der Netzwerk-Leseweg
# ist damit keine erfundene Ausgabe des Clients.
http_read_change_id=$(printf '%s' "$http_output" | grep -F "$HTTP_READ_SENTINEL" | grep -oE 'change_id=[^ ]+' | head -n1 | cut -d= -f2 || true)
if [ -z "$http_read_change_id" ]; then
  echo "run-integration-tests: HTTP-API-Rundlauf — die sentinel-tragende READ-Zeile trägt keine change_id: $http_output" >&2
  exit 1
fi
http_read_captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$http_read_change_id' AND table_name = '$HTTP_READ_TABLE' AND new_data->>'name' = '$HTTP_READ_SENTINEL'")
if [ -z "$http_read_captured" ] || [ "$http_read_captured" -lt 1 ]; then
  echo "run-integration-tests: HTTP-API-Rundlauf — die über GET /changes gelesene Änderung (change_id=$http_read_change_id, $HTTP_READ_SENTINEL) ist nicht real über cdc.changes lesbar (count=${http_read_captured:-leer})" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem HTTP-API-Rundlauf nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: HTTP-API-Rundlauf (LH-FA-SST-006, ADR-0057/ADR-0081) belegt — RegisterConsumer real per HTTP mit admin-Token ($HTTP_CONSUMER, cdc.consumer bestätigt), ListTables real per HTTP mit reader-Token (feed_e2e_full in der Antwort), GET /changes real per HTTP mit reader-Token (die eigens eingefügte Zeile id=$HTTP_READ_ID/$HTTP_READ_SENTINEL, Bereich [$http_read_position,$http_read_to), change_id=$http_read_change_id gegen cdc.changes gehalten), GET /diagnose real per HTTP mit reader-Token (HTTP-Diagnose-Querabgleich siehe oben): $http_output"

abdeckung_declare "HTTP-API-Consumer-Entfernung" "LH-FA-CON-006" "derselbe Wegwerf-Client bestätigt real per HTTP eine Position für den zuvor registrierten Consumer (macht ihn zum Retention-Blocker der Quelle, real gegen cdc.retention_blockers geprüft) und entfernt ihn danach über POST /consumers/remove mit dem admin-Token; cdc.consumer und cdc.consumer_position tragen ihn danach beide nicht mehr" "HTTP-API-Consumer-Entfernung (LH-FA-CON-006, ADR-0057) belegt"

# HTTP-API-Consumer-Entfernung (LH-FA-CON-006): zwei getrennte
# Aufrufe desselben Wegwerf-Clients (tools/harness/httpclient), damit der
# DB-Zustand real dazwischen geprüft werden kann — innerhalb eines
# einzigen Prozesslaufs wäre das nicht beobachtbar. Läuft nach dem
# bestehenden HTTP-API-Rundlauf oben, damit dessen eigene
# registered_via_http-Prüfung den noch registrierten Consumer sieht,
# bevor dieser Block ihn entfernt. Offset 1 liegt weit unter jeder realen
# LSN-abgeleiteten Position der übrigen Consumer dieses Laufs
# — real garantiert der kleinste Wert für src-e2e.
HTTP_REMOVE_OFFSET=1

set +e
ack_output=$(docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/httpclient acknowledge \
  "$HTTP_BASE_URL" "$HTTP_TOKEN_ADMIN" "$HTTP_CONSUMER" src-e2e "$HTTP_REMOVE_OFFSET" 2>&1)
ack_status=$?
set -e
if [ "$ack_status" -ne 0 ] || ! printf '%s' "$ack_output" | grep -qF "ACKNOWLEDGED"; then
  echo "run-integration-tests: HTTP-API-Consumer-Entfernung — AcknowledgeConsumer (httpclient acknowledge) endete mit Ausgang $ack_status: $ack_output" >&2
  exit 1
fi

blocker_before=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT consumer_id FROM cdc.retention_blockers WHERE source_id = 'src-e2e'")
if [ "$blocker_before" != "$HTTP_CONSUMER" ]; then
  echo "run-integration-tests: HTTP-API-Consumer-Entfernung — $HTTP_CONSUMER ist nach der realen Bestätigung nicht der Retention-Blocker von src-e2e (gefunden: ${blocker_before:-leer})" >&2
  exit 1
fi

set +e
remove_output=$(docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/httpclient remove \
  "$HTTP_BASE_URL" "$HTTP_TOKEN_ADMIN" "$HTTP_CONSUMER" 2>&1)
remove_status=$?
set -e
if [ "$remove_status" -ne 0 ] || ! printf '%s' "$remove_output" | grep -qE '^REMOVED consumer='"$HTTP_CONSUMER"' body=.*"removed":true'; then
  echo "run-integration-tests: HTTP-API-Consumer-Entfernung — RemoveConsumer (httpclient remove) endete mit Ausgang $remove_status: $remove_output" >&2
  exit 1
fi

consumer_after=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.consumer WHERE consumer_id = '$HTTP_CONSUMER'")
if [ "$consumer_after" != "0" ]; then
  echo "run-integration-tests: HTTP-API-Consumer-Entfernung — cdc.consumer trägt $HTTP_CONSUMER nach der realen Entfernung noch (count=$consumer_after)" >&2
  exit 1
fi

position_after=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.consumer_position WHERE consumer_id = '$HTTP_CONSUMER'")
if [ "$position_after" != "0" ]; then
  echo "run-integration-tests: HTTP-API-Consumer-Entfernung — cdc.consumer_position trägt $HTTP_CONSUMER nach der realen Entfernung noch (count=$position_after)" >&2
  exit 1
fi

blocker_after=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.retention_blockers WHERE source_id = 'src-e2e' AND consumer_id = '$HTTP_CONSUMER'")
if [ "$blocker_after" != "0" ]; then
  echo "run-integration-tests: HTTP-API-Consumer-Entfernung — cdc.retention_blockers nennt $HTTP_CONSUMER nach der realen Entfernung noch als Blocker von src-e2e (count=$blocker_after)" >&2
  exit 1
fi

echo "run-integration-tests: HTTP-API-Consumer-Entfernung (LH-FA-CON-006, ADR-0057) belegt — AcknowledgeConsumer real per HTTP mit admin-Token machte $HTTP_CONSUMER zum Retention-Blocker von src-e2e (cdc.retention_blockers real geprüft), RemoveConsumer real per HTTP mit admin-Token entfernte ihn danach — cdc.consumer, cdc.consumer_position und cdc.retention_blockers tragen ihn danach alle drei nicht mehr"

abdeckung_declare "Strukturiertes-Logging-Beleg" "LH-QA-OPS-004" "ein Wegwerf-Werkzeug liest die stdout-Logzeilen des laufenden Feed-Containers real per docker logs und prüft jede nicht-leere Zeile als eigenständiges JSON-Objekt mit den Feldern time/level/msg" "Strukturiertes-Logging-Beleg (LH-QA-OPS-004, ADR-0024) belegt"

# Strukturiertes-Logging-Beleg (LH-QA-OPS-004): tools/harness/logcheck
# liest die bislang akkumulierten stdout-Logzeilen des laufenden
# Feed-Containers (docker logs, CDC_LOG_LEVEL=debug in compose.yaml —
# reichlich Zeilen bis zu diesem späten Punkt im Lauf) und prüft jede
# nicht-leere Zeile als eigenständiges JSON-Objekt mit den drei vom
# Standard-Handler garantierten Feldern time/level/msg.
log_lines=$(docker logs "$FEED_CONTAINER" 2>&1)
if [ -z "$log_lines" ]; then
  echo "run-integration-tests: Strukturiertes-Logging-Beleg — docker logs $FEED_CONTAINER lieferte keine Zeile" >&2
  exit 1
fi

set +e
logcheck_output=$(printf '%s\n' "$log_lines" | docker run --rm -i --network none \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/logcheck 2>&1)
logcheck_status=$?
set -e
if [ "$logcheck_status" -ne 0 ] || ! printf '%s' "$logcheck_output" | grep -qE '^STRUCTURED_LOG_OK lines=[0-9]+ levels='; then
  echo "run-integration-tests: Strukturiertes-Logging-Beleg (logcheck) endete mit Ausgang $logcheck_status: $logcheck_output" >&2
  exit 1
fi

echo "run-integration-tests: Strukturiertes-Logging-Beleg (LH-QA-OPS-004, ADR-0024) belegt — $logcheck_output"

abdeckung_declare "Metriken-Minimum-Beleg (ausstehende Changes, Fehlerklassen)" "LH-QA-OPS-003" "ein eigener, zurückliegender Consumer und ein direkt in die Heartbeat-Projektion geschriebener Fehlerzustand belegen real über cdc.metrics die beiden bislang fehlenden Dimensionen: ausstehende Changes je Consumer und Fehler je Klasse" "Metriken-Minimum-Beleg (LH-QA-OPS-003, SPEC-009) belegt"

# Metriken-Minimum-Beleg (LH-QA-OPS-003): cdc.metrics trägt
# fünf der sieben im Lastenheft geforderten Dimensionen bereits mit
# eigenen Belegen anderswo (u. a. cdc_capture_lag,
# cdc_storage_bytes). Dieser Beleg schließt die beiden
# zuvor fehlenden real: ein eigener, zurückliegender Consumer belegt
# cdc_changes_pending; ein direkt geschriebener Fehlerzustand belegt
# cdc_errors_total — derselbe Schreibweg und dieselbe Wiederholschleife
# gegen den periodischen Heartbeat-Takt (5s) wie der CLI-Diagnose-
# Fehlerzustand-Beleg oben.
METRICS_CONSUMER=metrics-e2e-consumer
if ! exec_feed register-consumer "$METRICS_CONSUMER"; then
  echo "run-integration-tests: Metriken-Minimum-Beleg — register-consumer ($METRICS_CONSUMER) endete mit einem Fehler" >&2
  exit 1
fi
if ! exec_feed acknowledge-consumer "$METRICS_CONSUMER" 1; then
  echo "run-integration-tests: Metriken-Minimum-Beleg — acknowledge-consumer ($METRICS_CONSUMER, Position 1) endete mit einem Fehler" >&2
  exit 1
fi

pending_value=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_changes_pending' AND label = '$METRICS_CONSUMER'")
if [ -z "$pending_value" ] || [ "$pending_value" = "0" ]; then
  echo "run-integration-tests: Metriken-Minimum-Beleg — cdc.metrics trägt keine positive cdc_changes_pending-Zeile für $METRICS_CONSUMER (Wert: ${pending_value:-leer})" >&2
  exit 1
fi

errors_class_seen=0
errors_value=""
for _ in $(seq 1 20); do
  docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -c \
    "UPDATE cdc.process_heartbeat SET heartbeat_at = current_timestamp, error_class = 'internal' WHERE source_id = 'src-e2e'" >/dev/null
  errors_value=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_errors_total' AND label = 'internal'")
  if [ -n "$errors_value" ] && [ "$errors_value" != "0" ]; then
    errors_class_seen=1
    break
  fi
done
if [ "$errors_class_seen" -ne 1 ]; then
  echo "run-integration-tests: Metriken-Minimum-Beleg — cdc.metrics trägt keine positive cdc_errors_total-Zeile für die Klasse 'internal' innerhalb von 20 Versuchen (letzter Wert: ${errors_value:-leer})" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem Metriken-Minimum-Beleg nicht mehr weiter (der Fehlerzustand wurde direkt in cdc.process_heartbeat geschrieben, nicht vom Erfassungspfad ausgelöst — kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: Metriken-Minimum-Beleg (LH-QA-OPS-003, SPEC-009) belegt — cdc_changes_pending für $METRICS_CONSUMER trägt $pending_value, cdc_errors_total für Klasse 'internal' trägt $errors_value"

abdeckung_declare "gRPC-Stream-Rundlauf" "LH-FA-SST-008" "ein Wegwerf-Client öffnet real über gRPC den Server-Stream gegen den laufenden Feed-Container und empfängt eine danach committete Änderung; ein Öffnungsversuch ohne gültiges Token endet mit gRPC-Status Unauthenticated" "gRPC-Stream-Rundlauf (LH-FA-SST-008, ADR-0060) belegt"

# gRPC-Stream-Rundlauf (ADR-0060): ein Wegwerf-Client
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

docker rm -fv "$GRPC_CLIENT_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$GRPC_CLIENT_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/grpcclient "$GRPC_ADDR" "$HTTP_TOKEN_READER" >/dev/null

grpc_client_ready=0
for _ in $(seq 1 60); do
  if docker logs "$GRPC_CLIENT_CONTAINER" 2>/dev/null | grep -F "READY" >/dev/null; then
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
  docker rm -fv "$GRPC_CLIENT_CONTAINER" >/dev/null 2>&1 || true
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
    if docker logs "$GRPC_CLIENT_CONTAINER" 2>/dev/null | grep -F "RECEIVED" >/dev/null; then
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
  if docker logs "$GRPC_CLIENT_CONTAINER" 2>/dev/null | grep -F "REJECTED code=Unauthenticated" >/dev/null; then
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
docker rm -fv "$GRPC_CLIENT_CONTAINER" >/dev/null 2>&1 || true

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

abdeckung_declare "gRPC-Stream-Filter-Rundlauf" "LH-FA-SST-008" "ein Wegwerf-Client öffnet den Server-Stream mit gesetztem schema/table-Filter; eine Change auf der gefilterten Tabelle wird zugestellt, eine Change auf einer anderen, aktivierten Tabelle bleibt bei aktivem Filter unzugestellt" "gRPC-Stream-Filter-Rundlauf (LH-FA-SST-008, ADR-0133) belegt"

# gRPC-Stream-Filter-Rundlauf (ADR-0133): derselbe Wegwerf-Client
# (tools/harness/grpcclient), jetzt mit gesetztem schema/table-Filter
# (public/feed_e2e_full) geöffnet. Erst eine Change auf einer anderen,
# ebenfalls aktivierten Tabelle (feed_e2e_schema, über CDC_TABLES aktiviert
# und seit TestE2ESchemaChangeAddColumn oben capturebar — ihre reale
# Spaltenentfernung, die die Erfassung dieser Tabelle dauerhaft beendet,
# läuft erst deutlich später als TestE2ESchemaChangeDropColumn) — eine
# begrenzte, reale Wartezeit belegt, dass sie NICHT beim Client ankommt
# (`stream.Recv()` bleibt blockiert, keine RECEIVED-Zeile). Danach eine
# Change auf der gefilterten Tabelle — sie erreicht den Client, geprüft
# über dieselbe RECEIVED-Zeile und ihre `change_id` gegen `cdc.changes`.
# Eigener, von IDs 260ff. und 270ff. getrennter Wertebereich (280ff. auf
# feed_e2e_full, id=9010 auf feed_e2e_schema — vierstellig gewählt, damit
# kein späterer Bestandstest auf derselben Tabelle denselben Schlüssel
# beansprucht (feed_e2e_schema führt sonst nur 1/2/3/4/10/11)).
GRPC_FILTER_CLIENT_CONTAINER=cdc-e2e-grpcclient-filter
GRPC_FILTER_TABLE=feed_e2e_full
GRPC_FILTER_OTHER_TABLE=feed_e2e_schema
GRPC_FILTER_SENTINEL=GrpcFilterMatchE2ESentinel
GRPC_FILTER_OTHER_SENTINEL=GrpcFilterOtherE2ESentinel

docker rm -fv "$GRPC_FILTER_CLIENT_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$GRPC_FILTER_CLIENT_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/grpcclient "$GRPC_ADDR" "$HTTP_TOKEN_READER" public "$GRPC_FILTER_TABLE" >/dev/null

grpc_filter_client_ready=0
for _ in $(seq 1 60); do
  if docker logs "$GRPC_FILTER_CLIENT_CONTAINER" 2>/dev/null | grep -F "READY" >/dev/null; then
    grpc_filter_client_ready=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$GRPC_FILTER_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 1
done
if [ "$grpc_filter_client_ready" -ne 1 ]; then
  echo "run-integration-tests: gRPC-Stream-Filter-Rundlauf — Test-Client wurde nicht innerhalb der Zeitspanne bereit: $(docker logs "$GRPC_FILTER_CLIENT_CONTAINER" 2>&1 || true)" >&2
  docker rm -fv "$GRPC_FILTER_CLIENT_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$GRPC_FILTER_OTHER_TABLE (id, name) VALUES (9010, '$GRPC_FILTER_OTHER_SENTINEL');
SQL

# Keine Zustellung ist kein Ereignis, das ein Poll beobachten kann — eine
# reale, aber begrenzte Wartezeit, bevor geprüft wird, dass die
# nicht passende Change NICHT ankommt.
sleep 3
if docker logs "$GRPC_FILTER_CLIENT_CONTAINER" 2>/dev/null | grep -F "RECEIVED" >/dev/null; then
  grpc_filter_leaked_output=$(docker logs "$GRPC_FILTER_CLIENT_CONTAINER" 2>&1 || true)
  docker rm -fv "$GRPC_FILTER_CLIENT_CONTAINER" >/dev/null 2>&1 || true
  echo "run-integration-tests: gRPC-Stream-Filter-Rundlauf — eine Change auf $GRPC_FILTER_OTHER_TABLE erreichte den auf $GRPC_FILTER_TABLE gefilterten Client (Filter griff nicht): $grpc_filter_leaked_output" >&2
  exit 1
fi

grpc_filter_received=0
for grpc_filter_attempt in $(seq 1 5); do
  grpc_filter_id=$((280 + grpc_filter_attempt))
  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$GRPC_FILTER_TABLE (id, name) VALUES ($grpc_filter_id, '$GRPC_FILTER_SENTINEL');
SQL
  for _ in $(seq 1 20); do
    if docker logs "$GRPC_FILTER_CLIENT_CONTAINER" 2>/dev/null | grep -F "RECEIVED" >/dev/null; then
      grpc_filter_received=1
      break
    fi
    if [ "$(docker inspect --format '{{.State.Running}}' "$GRPC_FILTER_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
      break
    fi
    sleep 0.5
  done
  if [ "$grpc_filter_received" -eq 1 ]; then
    break
  fi
done

grpc_filter_client_output=$(docker logs "$GRPC_FILTER_CLIENT_CONTAINER" 2>&1 || true)
docker rm -fv "$GRPC_FILTER_CLIENT_CONTAINER" >/dev/null 2>&1 || true

if [ "$grpc_filter_received" -ne 1 ]; then
  echo "run-integration-tests: gRPC-Stream-Filter-Rundlauf — Test-Client empfing keine der committeten Änderungen auf der gefilterten Tabelle ($GRPC_FILTER_TABLE, $GRPC_FILTER_SENTINEL): $grpc_filter_client_output" >&2
  exit 1
fi
if ! printf '%s' "$grpc_filter_client_output" | grep -qE "RECEIVED .*table=$GRPC_FILTER_TABLE .*operation=INSERT .*new_image=.*$GRPC_FILTER_SENTINEL"; then
  echo "run-integration-tests: gRPC-Stream-Filter-Rundlauf — die RECEIVED-Zeile trägt nicht die erwartete gefilterte Änderung ($GRPC_FILTER_TABLE, $GRPC_FILTER_SENTINEL): $grpc_filter_client_output" >&2
  exit 1
fi
if printf '%s' "$grpc_filter_client_output" | grep -qF "$GRPC_FILTER_OTHER_SENTINEL"; then
  echo "run-integration-tests: gRPC-Stream-Filter-Rundlauf — die nicht passende Change ($GRPC_FILTER_OTHER_TABLE, $GRPC_FILTER_OTHER_SENTINEL) erscheint dennoch in der Client-Ausgabe: $grpc_filter_client_output" >&2
  exit 1
fi

grpc_filter_change_id=$(printf '%s' "$grpc_filter_client_output" | grep -oE 'RECEIVED change_id=[^ ]+' | head -n1 | cut -d= -f2 || true)
if [ -z "$grpc_filter_change_id" ]; then
  echo "run-integration-tests: gRPC-Stream-Filter-Rundlauf — die RECEIVED-Zeile trägt keine change_id: $grpc_filter_client_output" >&2
  exit 1
fi
grpc_filter_captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$grpc_filter_change_id' AND table_name = '$GRPC_FILTER_TABLE' AND new_data->>'name' = '$GRPC_FILTER_SENTINEL'")
if [ -z "$grpc_filter_captured" ] || [ "$grpc_filter_captured" -lt 1 ]; then
  echo "run-integration-tests: gRPC-Stream-Filter-Rundlauf — die über den gefilterten Stream empfangene Änderung (change_id=$grpc_filter_change_id, $GRPC_FILTER_SENTINEL) ist nicht real über cdc.changes lesbar (count=${grpc_filter_captured:-leer})" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem gRPC-Stream-Filter-Rundlauf nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: gRPC-Stream-Filter-Rundlauf (LH-FA-SST-008, ADR-0133) belegt — ein mit schema=public/table=$GRPC_FILTER_TABLE gefilterter Wegwerf-Client (tools/harness/grpcclient) empfing eine Change auf $GRPC_FILTER_OTHER_TABLE nicht (3s reale Wartezeit ohne RECEIVED-Zeile), aber eine danach committete Änderung auf $GRPC_FILTER_TABLE (change_id=$grpc_filter_change_id, unabhängig über cdc.changes gelesen)"

abdeckung_declare "gRPC-Administration-Rundlauf" "LH-FA-CON-001,LH-FA-CFG-004,LH-FA-REA-001" "ein Wegwerf-Client ruft real über gRPC mindestens eine Administration-RPC je Token-Klasse auf (ListTables und ReadChanges mit dem reader-Token, RegisterConsumer mit dem admin-Token); ein Aufruf ohne Token endet mit gRPC-Status Unauthenticated, ein reader-Token gegen die admin-RPC RegisterConsumer mit PermissionDenied" "gRPC-Administration-Rundlauf (ADR-0131) belegt"

abdeckung_declare "gRPC-Diagnose-Querabgleich" "LH-FA-SST-006,LH-FA-SST-003,LH-FA-ADM-002,LH-FA-ADM-003" "der gRPC-Diagnose-RPC liefert denselben Betriebsstatus (Lebenszeichen bekannt, Fehlerzustand) wie eine kontemporäre docker exec diagnose-Ausgabe — ein echter Wertabgleich zwischen zwei Zugriffswegen, nicht nur die Feld-Präsenz der Antwort" "gRPC-Diagnose-Querabgleich (ADR-0132) belegt"

# gRPC-Administration-Rundlauf (ADR-0131): ein Wegwerf-Client
# (tools/harness/grpcadminclient, per `go run` im Toolchain-Container)
# verbindet sich real über gRPC mit dem laufenden Feed-Container und ruft
# den `Administration`-Service auf — ListTables und ReadChanges mit dem
# reader-Token, RegisterConsumer mit dem admin-Token. Das Ergebnis von
# RegisterConsumer wird unabhängig gegen den Lesezugriffsweg `cdc.consumer`
# gehalten, das Ergebnis von ReadChanges gegen `cdc.changes`. Abschließend
# belegt derselbe Prozess die beiden Negative-Pfade: ein Aufruf ohne Token
# endet mit gRPC-Status `Unauthenticated`, ein reader-Token gegen die
# admin-RPC RegisterConsumer mit `PermissionDenied`. Ein echter
# Netzwerk-Request über denselben Compose-Netz-Alias wie der
# gRPC-Stream-Rundlauf oben, kein `docker exec` und kein Mock.
#
# Der ReadChanges-Beleg braucht eine bestimmte, unmittelbar zuvor erfasste
# Änderung: der Runner fügt eine eigene Zeile in die dauerhaft aktivierte
# Tabelle feed_e2e_full ein und wartet ihre Erfassung über cdc.changes ab,
# BEVOR der Client sie über die RPC liest; ihre commit_position geht als
# Bereich `[from, to)` in den Aufruf ein (ADR-0131 Teilfrage 3), ihre
# change_id wird unabhängig gegen denselben Lesezugriffsweg gehalten.
GRPC_ADMIN_CONSUMER_ID="grpc-admin-e2e-consumer"
GRPC_ADMIN_READ_SCHEMA=public
GRPC_ADMIN_READ_TABLE=feed_e2e_full
GRPC_ADMIN_READ_SENTINEL=GrpcAdminChangesReadE2ESentinel
GRPC_ADMIN_READ_ID=286
GRPC_ADMIN_READ_LIMIT=100

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$GRPC_ADMIN_READ_TABLE (id, name) VALUES ($GRPC_ADMIN_READ_ID, '$GRPC_ADMIN_READ_SENTINEL');
SQL

grpc_admin_read_position=""
for _ in $(seq 1 60); do
  grpc_admin_read_position=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$GRPC_ADMIN_READ_TABLE' AND new_data->>'id' = '$GRPC_ADMIN_READ_ID'" 2>/dev/null || true)
  if [ -n "$grpc_admin_read_position" ]; then
    break
  fi
  sleep 0.5
done
if [ -z "$grpc_admin_read_position" ]; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — die eigens eingefügte Zeile (id=$GRPC_ADMIN_READ_ID, $GRPC_ADMIN_READ_TABLE) wurde nicht innerhalb der Zeitspanne über cdc.changes erfasst" >&2
  exit 1
fi
grpc_admin_read_to=$((grpc_admin_read_position + 1))

# Kontemporäre CLI-Diagnose-Baseline (docker exec, unmittelbar vor dem
# gRPC-Aufruf) — der Querabgleich unten hält den über gRPC gelesenen
# Betriebsstatus gegen diese Ausgabe, nicht nur gegen die Feld-Präsenz der
# Antwort (ADR-0132).
set +e
grpc_admin_diagnose_cli_output=$(exec_feed diagnose)
grpc_admin_diagnose_cli_status=$?
set -e
if [ "$grpc_admin_diagnose_cli_status" -ne 0 ]; then
  echo "run-integration-tests: gRPC-Diagnose-Querabgleich — CLI-Diagnose-Baseline (docker exec diagnose) endete mit Ausgang $grpc_admin_diagnose_cli_status: $grpc_admin_diagnose_cli_output" >&2
  exit 1
fi

# `set +e` um den Client-Aufruf: der Ausgang wird von der folgenden Prüfung
# ausgewertet und über eine sichtbare Fehlerzeile gemeldet; unter `set -e`
# beendet ein fehlgeschlagenes `docker run` schon die Zuweisung selbst und
# die Phase endet rot ohne Ausgabe (AGENTS.md §3.9).
set +e
grpc_admin_output=$(docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/grpcadminclient "$GRPC_ADDR" "$HTTP_TOKEN_READER" "$HTTP_TOKEN_ADMIN" "$GRPC_ADMIN_CONSUMER_ID" src-e2e pub_pgc_e2e \
  "$GRPC_ADMIN_READ_SCHEMA" "$GRPC_ADMIN_READ_TABLE" "$grpc_admin_read_position" "$grpc_admin_read_to" "$GRPC_ADMIN_READ_LIMIT" 2>&1)
grpc_admin_exit=$?
set -e
if [ "$grpc_admin_exit" -ne 0 ]; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — Test-Client endete nicht mit Ausgang 0 (Ausgang: $grpc_admin_exit): $grpc_admin_output" >&2
  exit 1
fi
if ! printf '%s' "$grpc_admin_output" | grep -qE '^LISTED tables=[0-9]+ retained=[0-9]+$'; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — keine LISTED-Zeile (ListTables mit reader-Token): $grpc_admin_output" >&2
  exit 1
fi
if ! printf '%s' "$grpc_admin_output" | grep -qF "REGISTERED consumer_id=$GRPC_ADMIN_CONSUMER_ID"; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — keine REGISTERED-Zeile (RegisterConsumer mit admin-Token): $grpc_admin_output" >&2
  exit 1
fi
if ! printf '%s' "$grpc_admin_output" | grep -qE '^READ changes=[0-9]+ table=feed_e2e_full schema=public '; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — keine READ-Zeile des ReadChanges-Aufrufs (reader-Token, Bereich [$grpc_admin_read_position,$grpc_admin_read_to)): $grpc_admin_output" >&2
  exit 1
fi
if ! printf '%s' "$grpc_admin_output" | grep -qF "$GRPC_ADMIN_READ_SENTINEL"; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — die READ-Zeile trägt die eigens eingefügte Zeile (id=$GRPC_ADMIN_READ_ID, $GRPC_ADMIN_READ_SENTINEL) nicht: $grpc_admin_output" >&2
  exit 1
fi
if ! printf '%s' "$grpc_admin_output" | grep -qE '^DIAGNOSED heartbeat_known='; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — keine DIAGNOSED-Zeile (Diagnose mit reader-Token): $grpc_admin_output" >&2
  exit 1
fi

# Querabgleich (ADR-0132): dieselbe Quelle über zwei unabhängige
# Zugriffswege — der CLI-Sondermodus (Baseline oben) und der gRPC-RPC —
# tragen an diesem Zeitpunkt denselben Betriebsstatus. Verglichen werden
# zwei robuste, boolesche Ableitungen statt roher Zeitstempel/Sekundenwerte
# (die zwischen zwei sequenziellen Aufrufen nie exakt gleich wären): ob ein
# Lebenszeichen bekannt ist, und ob ein Fehlerzustand vorliegt.
if printf '%s' "$grpc_admin_diagnose_cli_output" | grep -qF "Betriebsstatus: Lebenszeichen vor"; then
  grpc_admin_diagnose_cli_known=1
else
  grpc_admin_diagnose_cli_known=0
fi
grpc_admin_diagnose_known_line=$(printf '%s' "$grpc_admin_output" | grep -oE '^DIAGNOSED heartbeat_known=[a-z]+' | head -n1)
case "$grpc_admin_diagnose_known_line" in
  *heartbeat_known=true) grpc_admin_diagnose_grpc_known=1 ;;
  *heartbeat_known=false) grpc_admin_diagnose_grpc_known=0 ;;
  *)
    echo "run-integration-tests: gRPC-Diagnose-Querabgleich — DIAGNOSED-Zeile trägt kein auswertbares heartbeat_known: $grpc_admin_output" >&2
    exit 1
    ;;
esac
if [ "$grpc_admin_diagnose_cli_known" != "$grpc_admin_diagnose_grpc_known" ]; then
  echo "run-integration-tests: gRPC-Diagnose-Querabgleich — Betriebsstatus (Lebenszeichen bekannt) widersprüchlich zwischen CLI (known=$grpc_admin_diagnose_cli_known) und gRPC (known=$grpc_admin_diagnose_grpc_known): CLI=$grpc_admin_diagnose_cli_output GRPC=$grpc_admin_output" >&2
  exit 1
fi

if printf '%s' "$grpc_admin_diagnose_cli_output" | grep -qF "Fehlerzustand: keiner (Normalbetrieb)"; then
  grpc_admin_diagnose_cli_no_error=1
else
  grpc_admin_diagnose_cli_no_error=0
fi
grpc_admin_diagnose_error_class=$(printf '%s' "$grpc_admin_output" | grep -oE 'heartbeat_error_class=[^ ]*' | head -n1 | cut -d= -f2)
if [ -z "$grpc_admin_diagnose_error_class" ]; then
  grpc_admin_diagnose_grpc_no_error=1
else
  grpc_admin_diagnose_grpc_no_error=0
fi
if [ "$grpc_admin_diagnose_cli_no_error" != "$grpc_admin_diagnose_grpc_no_error" ]; then
  echo "run-integration-tests: gRPC-Diagnose-Querabgleich — Fehlerzustand widersprüchlich zwischen CLI (kein Fehler=$grpc_admin_diagnose_cli_no_error) und gRPC (kein Fehler=$grpc_admin_diagnose_grpc_no_error): CLI=$grpc_admin_diagnose_cli_output GRPC=$grpc_admin_output" >&2
  exit 1
fi

echo "run-integration-tests: gRPC-Diagnose-Querabgleich (ADR-0132) belegt — der gRPC-Diagnose-RPC (reader-Token) und eine kontemporäre docker exec diagnose-Ausgabe tragen denselben Betriebsstatus: Lebenszeichen bekannt=$grpc_admin_diagnose_grpc_known, kein Fehlerzustand=$grpc_admin_diagnose_grpc_no_error"

if ! printf '%s' "$grpc_admin_output" | grep -qF "REJECTED code=Unauthenticated"; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — Aufruf ohne Token wurde nicht mit Unauthenticated abgelehnt: $grpc_admin_output" >&2
  exit 1
fi
if ! printf '%s' "$grpc_admin_output" | grep -qF "REJECTED code=PermissionDenied"; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — reader-Token gegen RegisterConsumer wurde nicht mit PermissionDenied abgelehnt: $grpc_admin_output" >&2
  exit 1
fi

# Unabhängiger SQL-Beleg, dass RegisterConsumer real über den bestehenden
# Lesezugriffsweg cdc.consumer sichtbar ist — kein erfundenes REGISTERED
# des Clients.
grpc_admin_registered=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.consumer WHERE consumer_id = '$GRPC_ADMIN_CONSUMER_ID'")
if [ -z "$grpc_admin_registered" ] || [ "$grpc_admin_registered" -lt 1 ]; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — der über RegisterConsumer registrierte Consumer ($GRPC_ADMIN_CONSUMER_ID) ist nicht real über cdc.consumer lesbar (count=${grpc_admin_registered:-leer})" >&2
  exit 1
fi

# Unabhängiger SQL-Beleg, dass genau die über ReadChanges gelesene Änderung
# real erfasst wurde: die change_id der sentinel-tragenden READ-Zeile steht
# über den bestehenden Lesezugriffsweg in cdc.changes — der gRPC-Leseweg ist
# damit keine erfundene Ausgabe des Clients.
grpc_admin_read_change_id=$(printf '%s' "$grpc_admin_output" | grep -F "$GRPC_ADMIN_READ_SENTINEL" | grep -oE 'change_id=[^ ]+' | head -n1 | cut -d= -f2 || true)
if [ -z "$grpc_admin_read_change_id" ]; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — die sentinel-tragende READ-Zeile trägt keine change_id: $grpc_admin_output" >&2
  exit 1
fi
grpc_admin_read_captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$grpc_admin_read_change_id' AND table_name = '$GRPC_ADMIN_READ_TABLE' AND new_data->>'name' = '$GRPC_ADMIN_READ_SENTINEL'")
if [ -z "$grpc_admin_read_captured" ] || [ "$grpc_admin_read_captured" -lt 1 ]; then
  echo "run-integration-tests: gRPC-Administration-Rundlauf — die über ReadChanges gelesene Änderung (change_id=$grpc_admin_read_change_id, $GRPC_ADMIN_READ_SENTINEL) ist nicht real über cdc.changes lesbar (count=${grpc_admin_read_captured:-leer})" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem gRPC-Administration-Rundlauf nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: gRPC-Administration-Rundlauf (ADR-0131) belegt — ein Wegwerf-Client (tools/harness/grpcadminclient) rief real über gRPC ListTables (reader-Token), ReadChanges (reader-Token, Bereich [$grpc_admin_read_position,$grpc_admin_read_to), change_id=$grpc_admin_read_change_id gegen cdc.changes gehalten), Diagnose (reader-Token, gRPC-Diagnose-Querabgleich siehe oben) und RegisterConsumer (admin-Token) gegen den laufenden Feed-Container ($GRPC_ADDR) auf, der registrierte Consumer ($GRPC_ADMIN_CONSUMER_ID) ist unabhängig über cdc.consumer lesbar; ein Aufruf ohne Token wurde mit Unauthenticated, ein reader-Token gegen RegisterConsumer mit PermissionDenied abgelehnt: $grpc_admin_output"

abdeckung_declare "Direkter Zugriffsweg Live-Reload (gRPC enable)" "LH-FA-CFG-001" "EnableTable über den direkten gRPC-Zugriffsweg (ohne die SQL-Antragsqueue) aktualisiert den laufenden Assembler ohne Neustart — derselbe Live-Reload-Vertrag wie der SQL-Antragsqueue-Pfad, wegunabhängig" "Direkter gRPC-Zugriffsweg — EnableTable"
abdeckung_declare "Direkter Zugriffsweg Live-Reload (gRPC disable)" "LH-FA-CFG-002" "derselbe direkte gRPC-Zugriffsweg spiegelbildlich: DisableTable beendet die Erfassung dieser Tabelle ohne Neustart, der Feed-Container läuft weiter" "Direkter gRPC-Zugriffsweg — DisableTable"
abdeckung_declare "Direkter Zugriffsweg Live-Reload (HTTP enable)" "LH-FA-CFG-001" "derselbe Beleg über den direkten HTTP-Zugriffsweg — fachliche Gleichwertigkeit von HTTP und gRPC (\`LH-FA-SST-006\`)" "Direkter HTTP-Zugriffsweg — EnableTable"

# Direkter gRPC-/HTTP-Zugriffsweg — EnableTable/DisableTable ohne die
# SQL-Antragsqueue (LH-FA-CFG-001): bis zu diesem Zug
# aktualisierte nur der SQL-Antragsqueue-Pfad (cdc.enable_table/
# cdc.disable_table, verarbeitet von der Administrations-Goroutine) den
# laufenden Assembler — EnableTable/DisableTable über gRPC oder HTTP riefen
# denselben Inbound Use Case auf, ohne den Nachtrag in den laufenden Prozess
# (internal/bootstrap/assemblersync.go). Eine über diesen Weg aktivierte
# Tabelle verlor jede Änderung zwischen Aktivierung und dem nächsten
# Prozess-Neustart. Dieser Beleg spielt exakt das ursprüngliche
# Reproduktionsszenario nach: enable-table/disable-table real über gRPC bzw.
# HTTP gegen den laufenden Feed-Container, ohne die SQL-Antragsqueue zu
# berühren.
GRPC_API_ENABLE_TABLE=feed_e2e_api_enable_grpc
GRPC_API_ENABLE_TABLE_ID=tbl-e2e-api-enable-grpc
GRPC_API_ENABLE_SCHEMA_VERSION_ID=tbl-e2e-api-enable-grpc-v1

set +e
grpc_api_enable_output=$(docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/grpcenabletableclient "$GRPC_ADDR" enable-table \
  "$HTTP_TOKEN_ADMIN" src-e2e public "$GRPC_API_ENABLE_TABLE" "$GRPC_API_ENABLE_TABLE_ID" "$GRPC_API_ENABLE_SCHEMA_VERSION_ID" pub_pgc_e2e 2>&1)
grpc_api_enable_status=$?
set -e
if [ "$grpc_api_enable_status" -ne 0 ]; then
  echo "run-integration-tests: direkter gRPC-Zugriffsweg — EnableTable($GRPC_API_ENABLE_TABLE) endete mit Ausgang $grpc_api_enable_status: $grpc_api_enable_output" >&2
  exit 1
fi
if ! printf '%s' "$grpc_api_enable_output" | grep -qF "ENABLED"; then
  echo "run-integration-tests: direkter gRPC-Zugriffsweg — keine ENABLED-Zeile: $grpc_api_enable_output" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$GRPC_API_ENABLE_TABLE (id, name) VALUES (1, 'GrpcApiEnabled');
SQL

grpc_api_enable_captured=0
for _ in $(seq 1 120); do
  found=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$GRPC_API_ENABLE_TABLE' AND new_data->>'id' = '1'")
  if [ "$found" = "1" ]; then
    grpc_api_enable_captured=1
    break
  fi
  sleep 0.25
done
if [ "$grpc_api_enable_captured" -ne 1 ]; then
  echo "run-integration-tests: direkter gRPC-Zugriffsweg — über EnableTable aktivierte Tabelle $GRPC_API_ENABLE_TABLE — Änderung (id=1) wurde nicht vom laufenden Feed-Container erfasst (kein Neustart) — genau der Fehler, den dieser Zug behebt" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem direkten gRPC-EnableTable nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: Direkter gRPC-Zugriffsweg — EnableTable($GRPC_API_ENABLE_TABLE) ohne SQL-Antragsqueue und ohne Neustart verarbeitet, Änderung id=1 real erfasst: $grpc_api_enable_output"

# Deaktivierung: derselbe direkte gRPC-Zugriffsweg, spiegelbildlich.
set +e
grpc_api_disable_output=$(docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/grpcenabletableclient "$GRPC_ADDR" disable-table \
  "$HTTP_TOKEN_ADMIN" src-e2e public "$GRPC_API_ENABLE_TABLE" pub_pgc_e2e 2>&1)
grpc_api_disable_status=$?
set -e
if [ "$grpc_api_disable_status" -ne 0 ]; then
  echo "run-integration-tests: direkter gRPC-Zugriffsweg — DisableTable($GRPC_API_ENABLE_TABLE) endete mit Ausgang $grpc_api_disable_status: $grpc_api_disable_output" >&2
  exit 1
fi
if ! printf '%s' "$grpc_api_disable_output" | grep -qF "DISABLED"; then
  echo "run-integration-tests: direkter gRPC-Zugriffsweg — keine DISABLED-Zeile: $grpc_api_disable_output" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$GRPC_API_ENABLE_TABLE (id, name) VALUES (2, 'GrpcApiDisabled');
SQL

# Keine Erfassung ist kein Ereignis, das ein Poll beobachten kann — eine
# reale, aber begrenzte Wartezeit, bevor geprüft wird, dass die Zeile NICHT
# in cdc.changes ankommt (dasselbe Muster wie der SQL-Administration
# Live-Reload-Beleg oben).
sleep 3
grpc_api_disabled_leaked=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$GRPC_API_ENABLE_TABLE' AND new_data->>'id' = '2'")
if [ "$grpc_api_disabled_leaked" != "0" ]; then
  echo "run-integration-tests: direkter gRPC-Zugriffsweg — nach DisableTable($GRPC_API_ENABLE_TABLE) wurde eine Änderung (id=2) dennoch erfasst — Deaktivierung griff nicht" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem direkten gRPC-DisableTable nicht mehr weiter (der Prozess muss weiterlaufen, nur die Tabellen-Erfassung endet)" >&2
  exit 1
fi

echo "run-integration-tests: Direkter gRPC-Zugriffsweg — DisableTable($GRPC_API_ENABLE_TABLE) ohne SQL-Antragsqueue verarbeitet, Änderung id=2 nicht erfasst, Feed-Container läuft unverändert weiter: $grpc_api_disable_output"

# Derselbe Beleg über den direkten HTTP-Zugriffsweg (nur der Enable-Fall,
# LH-FA-SST-006 Boundary „fachliche Gleichwertigkeit über alle
# Zugriffswege"): ein Wegwerf-Client (tools/harness/httpclient enable-table)
# ruft EnableTable real per HTTP mit dem Admin-Token auf, ohne die
# SQL-Antragsqueue zu durchlaufen.
HTTP_API_ENABLE_TABLE=feed_e2e_api_enable_http
HTTP_API_ENABLE_TABLE_ID=tbl-e2e-api-enable-http
HTTP_API_ENABLE_SCHEMA_VERSION_ID=tbl-e2e-api-enable-http-v1

set +e
http_api_enable_output=$(docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/httpclient enable-table \
  "$HTTP_BASE_URL" "$HTTP_TOKEN_ADMIN" src-e2e public "$HTTP_API_ENABLE_TABLE" "$HTTP_API_ENABLE_TABLE_ID" "$HTTP_API_ENABLE_SCHEMA_VERSION_ID" pub_pgc_e2e 2>&1)
http_api_enable_status=$?
set -e
if [ "$http_api_enable_status" -ne 0 ]; then
  echo "run-integration-tests: direkter HTTP-Zugriffsweg — EnableTable($HTTP_API_ENABLE_TABLE) endete mit Ausgang $http_api_enable_status: $http_api_enable_output" >&2
  exit 1
fi
if ! printf '%s' "$http_api_enable_output" | grep -qF "ENABLED"; then
  echo "run-integration-tests: direkter HTTP-Zugriffsweg — keine ENABLED-Zeile: $http_api_enable_output" >&2
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$HTTP_API_ENABLE_TABLE (id, name) VALUES (1, 'HttpApiEnabled');
SQL

http_api_enable_captured=0
for _ in $(seq 1 120); do
  found=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$HTTP_API_ENABLE_TABLE' AND new_data->>'id' = '1'")
  if [ "$found" = "1" ]; then
    http_api_enable_captured=1
    break
  fi
  sleep 0.25
done
if [ "$http_api_enable_captured" -ne 1 ]; then
  echo "run-integration-tests: direkter HTTP-Zugriffsweg — über EnableTable aktivierte Tabelle $HTTP_API_ENABLE_TABLE — Änderung (id=1) wurde nicht vom laufenden Feed-Container erfasst (kein Neustart) — genau der Fehler, den dieser Zug behebt" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem direkten HTTP-EnableTable nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: Direkter HTTP-Zugriffsweg — EnableTable($HTTP_API_ENABLE_TABLE) ohne SQL-Antragsqueue und ohne Neustart verarbeitet, Änderung id=1 real erfasst: $http_api_enable_output"

abdeckung_declare "SSE-Stream-Rundlauf" "LH-FA-SST-008" "ein Wegwerf-Client öffnet real per HTTP den Endpunkt \`GET /changes/stream\` gegen den laufenden Feed-Container und empfängt eine danach committete Änderung; ein Aufruf ohne gültiges Token endet mit HTTP-Status 401" "SSE-Stream-Rundlauf (LH-FA-SST-008, ADR-0061) belegt"

# SSE-Stream-Rundlauf (ADR-0061): ein Wegwerf-Client
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

docker rm -fv "$SSE_CLIENT_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$SSE_CLIENT_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/sseclient "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" >/dev/null

sse_client_ready=0
for _ in $(seq 1 60); do
  if docker logs "$SSE_CLIENT_CONTAINER" 2>/dev/null | grep -F "READY" >/dev/null; then
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
  docker rm -fv "$SSE_CLIENT_CONTAINER" >/dev/null 2>&1 || true
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
    if docker logs "$SSE_CLIENT_CONTAINER" 2>/dev/null | grep -F "RECEIVED" >/dev/null; then
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
  if docker logs "$SSE_CLIENT_CONTAINER" 2>/dev/null | grep -F "REJECTED code=401" >/dev/null; then
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
docker rm -fv "$SSE_CLIENT_CONTAINER" >/dev/null 2>&1 || true

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

abdeckung_declare "SSE-Stream-Filter-Rundlauf" "LH-FA-SST-008" "ein Wegwerf-Client öffnet GET /changes/stream mit gesetztem schema/table-Query-Filter; eine Change auf der gefilterten Tabelle wird zugestellt, eine Change auf einer anderen, aktivierten Tabelle bleibt bei aktivem Filter unzugestellt" "SSE-Stream-Filter-Rundlauf (LH-FA-SST-008, ADR-0133) belegt"

# SSE-Stream-Filter-Rundlauf (ADR-0133): derselbe Wegwerf-Client
# (tools/harness/sseclient), jetzt mit gesetztem schema/table-Query-Filter
# (public/feed_e2e_full) geöffnet — dieselbe Reihenfolge wie beim
# gRPC-Filter-Rundlauf oben: erst eine Change auf feed_e2e_schema (real
# geprüfte Abwesenheit über eine begrenzte Wartezeit; die Tabelle bleibt
# bis TestE2ESchemaChangeDropColumn unten capturebar, siehe gRPC-Block
# oben), danach eine Change auf der gefilterten Tabelle. Eigener
# Wertebereich (290ff. auf feed_e2e_full, id=9011 auf feed_e2e_schema,
# getrennt vom id=9010 des gRPC-Filter-Rundlaufs oben — beide vierstellig,
# damit kein späterer Bestandstest auf derselben Tabelle denselben
# Schlüssel beansprucht).
SSE_FILTER_CLIENT_CONTAINER=cdc-e2e-sseclient-filter
SSE_FILTER_TABLE=feed_e2e_full
SSE_FILTER_OTHER_TABLE=feed_e2e_schema
SSE_FILTER_SENTINEL=SseFilterMatchE2ESentinel
SSE_FILTER_OTHER_SENTINEL=SseFilterOtherE2ESentinel

docker rm -fv "$SSE_FILTER_CLIENT_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$SSE_FILTER_CLIENT_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/sseclient "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" public "$SSE_FILTER_TABLE" >/dev/null

sse_filter_client_ready=0
for _ in $(seq 1 60); do
  if docker logs "$SSE_FILTER_CLIENT_CONTAINER" 2>/dev/null | grep -F "READY" >/dev/null; then
    sse_filter_client_ready=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$SSE_FILTER_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 1
done
if [ "$sse_filter_client_ready" -ne 1 ]; then
  echo "run-integration-tests: SSE-Stream-Filter-Rundlauf — Test-Client wurde nicht innerhalb der Zeitspanne bereit: $(docker logs "$SSE_FILTER_CLIENT_CONTAINER" 2>&1 || true)" >&2
  docker rm -fv "$SSE_FILTER_CLIENT_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$SSE_FILTER_OTHER_TABLE (id, name) VALUES (9011, '$SSE_FILTER_OTHER_SENTINEL');
SQL

# Keine Zustellung ist kein Ereignis, das ein Poll beobachten kann — eine
# reale, aber begrenzte Wartezeit, bevor geprüft wird, dass die
# nicht passende Change NICHT ankommt.
sleep 3
if docker logs "$SSE_FILTER_CLIENT_CONTAINER" 2>/dev/null | grep -F "RECEIVED" >/dev/null; then
  sse_filter_leaked_output=$(docker logs "$SSE_FILTER_CLIENT_CONTAINER" 2>&1 || true)
  docker rm -fv "$SSE_FILTER_CLIENT_CONTAINER" >/dev/null 2>&1 || true
  echo "run-integration-tests: SSE-Stream-Filter-Rundlauf — eine Change auf $SSE_FILTER_OTHER_TABLE erreichte den auf $SSE_FILTER_TABLE gefilterten Client (Filter griff nicht): $sse_filter_leaked_output" >&2
  exit 1
fi

sse_filter_received=0
for sse_filter_attempt in $(seq 1 5); do
  sse_filter_id=$((290 + sse_filter_attempt))
  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$SSE_FILTER_TABLE (id, name) VALUES ($sse_filter_id, '$SSE_FILTER_SENTINEL');
SQL
  for _ in $(seq 1 20); do
    if docker logs "$SSE_FILTER_CLIENT_CONTAINER" 2>/dev/null | grep -F "RECEIVED" >/dev/null; then
      sse_filter_received=1
      break
    fi
    if [ "$(docker inspect --format '{{.State.Running}}' "$SSE_FILTER_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
      break
    fi
    sleep 0.5
  done
  if [ "$sse_filter_received" -eq 1 ]; then
    break
  fi
done

sse_filter_client_output=$(docker logs "$SSE_FILTER_CLIENT_CONTAINER" 2>&1 || true)
docker rm -fv "$SSE_FILTER_CLIENT_CONTAINER" >/dev/null 2>&1 || true

if [ "$sse_filter_received" -ne 1 ]; then
  echo "run-integration-tests: SSE-Stream-Filter-Rundlauf — Test-Client empfing keine der committeten Änderungen auf der gefilterten Tabelle ($SSE_FILTER_TABLE, $SSE_FILTER_SENTINEL): $sse_filter_client_output" >&2
  exit 1
fi
if ! printf '%s' "$sse_filter_client_output" | grep -qE "RECEIVED .*table=$SSE_FILTER_TABLE .*operation=INSERT .*new_image=.*$SSE_FILTER_SENTINEL"; then
  echo "run-integration-tests: SSE-Stream-Filter-Rundlauf — die RECEIVED-Zeile trägt nicht die erwartete gefilterte Änderung ($SSE_FILTER_TABLE, $SSE_FILTER_SENTINEL): $sse_filter_client_output" >&2
  exit 1
fi
if printf '%s' "$sse_filter_client_output" | grep -qF "$SSE_FILTER_OTHER_SENTINEL"; then
  echo "run-integration-tests: SSE-Stream-Filter-Rundlauf — die nicht passende Change ($SSE_FILTER_OTHER_TABLE, $SSE_FILTER_OTHER_SENTINEL) erscheint dennoch in der Client-Ausgabe: $sse_filter_client_output" >&2
  exit 1
fi

sse_filter_change_id=$(printf '%s' "$sse_filter_client_output" | grep -oE 'RECEIVED change_id=[^ ]+' | head -n1 | cut -d= -f2 || true)
if [ -z "$sse_filter_change_id" ]; then
  echo "run-integration-tests: SSE-Stream-Filter-Rundlauf — die RECEIVED-Zeile trägt keine change_id: $sse_filter_client_output" >&2
  exit 1
fi
sse_filter_captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$sse_filter_change_id' AND table_name = '$SSE_FILTER_TABLE' AND new_data->>'name' = '$SSE_FILTER_SENTINEL'")
if [ -z "$sse_filter_captured" ] || [ "$sse_filter_captured" -lt 1 ]; then
  echo "run-integration-tests: SSE-Stream-Filter-Rundlauf — die über den gefilterten Stream empfangene Änderung (change_id=$sse_filter_change_id, $SSE_FILTER_SENTINEL) ist nicht real über cdc.changes lesbar (count=${sse_filter_captured:-leer})" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem SSE-Stream-Filter-Rundlauf nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: SSE-Stream-Filter-Rundlauf (LH-FA-SST-008, ADR-0133) belegt — ein mit schema=public/table=$SSE_FILTER_TABLE gefilterter Wegwerf-Client (tools/harness/sseclient) empfing eine Change auf $SSE_FILTER_OTHER_TABLE nicht (3s reale Wartezeit ohne RECEIVED-Zeile), aber eine danach committete Änderung auf $SSE_FILTER_TABLE (change_id=$sse_filter_change_id, unabhängig über cdc.changes gelesen)"

abdeckung_declare "NATS-Vollinhalts-Stream-Rundlauf" "LH-FA-SST-008" "ein Wegwerf-Client verbindet sich real über NATS mit gültigem Token, abonniert cdc.stream.<...> und empfängt eine danach committete Änderung als vollständiges JSON-Event; ein Verbindungsversuch ohne und einer mit falschem Token werden vom NATS-Server abgelehnt; das bestehende Wecksignal (natssub) funktioniert mit demselben Test-Token unverändert weiter" "NATS-Vollinhalts-Stream-Rundlauf (LH-FA-SST-008, ADR-0100) belegt"

# NATS-Vollinhalts-Stream-Rundlauf (ADR-0100): ein
# Wegwerf-Client (tools/harness/natsstreamsub, per `go run` im
# Toolchain-Container) belegt beide Ablehnungshälften — ein
# Verbindungsversuch ohne Token ("REJECTED-NO-TOKEN") und einer mit einem
# falschen Token ("REJECTED-WRONG-TOKEN") werden vom NATS-Server bereits auf
# Verbindungsebene abgelehnt (ADR-0100 Teilfrage 4) — compose.yaml betreibt
# den nats-Service seit dieser ADR mit `--auth $NATS_STREAM_TOKEN` (serverweit,
# bewusst benannter Nebeneffekt, ADR-0100 §Konsequenzen). Danach verbindet
# er sich real mit dem gültigen Token, abonniert das tabellen-granulare
# Subjekt cdc.stream.src-e2e.public.feed_e2e_full, bevor die auslösende
# Change entsteht, und empfängt sie als vollständiges JSON-Event
# ("RECEIVED") — geprüft über Tabelle, Operation und den Wert der Spalte
# `name` im neuen Row Image sowie über die change_id gegen den bestehenden
# Lesezugriffsweg cdc.changes; die Feldvollständigkeit des
# Nachrichtenschemas trägt internal/adapters/driven/natsstream/publisher_test.go
# auf Unit-Ebene. Derselbe Test-Token trägt zugleich den Beleg, dass das
# bestehende Wecksignal (natssub, oben bereits real gegen denselben Token
# geführt) unverändert funktioniert — dieselbe Verbindung, zwei Fähigkeiten
# (ADR-0100 Teilfrage 5). `feed_e2e_full` bleibt über den ganzen Lauf
# aktiviert (siehe Lasttest-Beleg oben); die Retry-IDs 301–305 belegt keine
# andere Phase derselben Tabelle (261–265 gRPC, 271–275 SSE, 285 der
# HTTP-Lese-Beleg). Läuft vor dem Upgrade-Sicherheits-Container-Tausch
# (der den Feed-Container ersetzt) und vor der Container-Ende-Grenze der
# beiden Schema-Negative-Testfunktionen unten.
NATS_STREAM_CLIENT_CONTAINER=cdc-e2e-natsstreamsub
NATS_STREAM_TABLE=feed_e2e_full
NATS_STREAM_SENTINEL=NatsStreamE2ESentinel
NATS_STREAM_SUBJECT="cdc.stream.src-e2e.public.$NATS_STREAM_TABLE"

docker rm -fv "$NATS_STREAM_CLIENT_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$NATS_STREAM_CLIENT_CONTAINER" --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go run ./tools/harness/natsstreamsub "nats://nats:4222" "$NATS_STREAM_SUBJECT" "$NATS_STREAM_TOKEN" >/dev/null

nats_stream_client_ready=0
for _ in $(seq 1 60); do
  if docker logs "$NATS_STREAM_CLIENT_CONTAINER" 2>/dev/null | grep -F "READY" >/dev/null; then
    nats_stream_client_ready=1
    break
  fi
  if [ "$(docker inspect --format '{{.State.Running}}' "$NATS_STREAM_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
    break
  fi
  sleep 1
done
if [ "$nats_stream_client_ready" -ne 1 ]; then
  echo "run-integration-tests: NATS-Vollinhalts-Stream-Rundlauf — Test-Client wurde nicht innerhalb der Zeitspanne bereit: $(docker logs "$NATS_STREAM_CLIENT_CONTAINER" 2>&1 || true)" >&2
  docker rm -fv "$NATS_STREAM_CLIENT_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

# Die beiden Ablehnungs-Checks pollen statt einmalig zu prüfen: Der
# NATS-Testclient schreibt REJECTED-NO-TOKEN/REJECTED-WRONG-TOKEN/READY
# strikt sequenziell auf stdout, aber `docker logs` liest den Docker-eigenen
# Log-Puffer, dessen Flush minimal hinter dem Container-Stdout zurückbleiben
# kann — real beobachtet im PostgreSQL-17-Leg von e2e.yml-Lauf 35438303264:
# der READY-Wait oben hatte bereits terminiert, der REJECTED-WRONG-TOKEN-Check
# direkt danach sah im selben `docker logs`-Aufruf noch keine der beiden
# Zeilen. Dasselbe Poll-statt-Einzel-Check-Muster wie beim READY-/
# RECEIVED-Warten weiter oben in dieser Datei.
nats_stream_client_rejected_no_token=0
for _ in $(seq 1 20); do
  if docker logs "$NATS_STREAM_CLIENT_CONTAINER" 2>/dev/null | grep -F "REJECTED-NO-TOKEN" >/dev/null; then
    nats_stream_client_rejected_no_token=1
    break
  fi
  sleep 0.2
done
if [ "$nats_stream_client_rejected_no_token" -ne 1 ]; then
  echo "run-integration-tests: NATS-Vollinhalts-Stream-Rundlauf — Verbindungsversuch ohne Token wurde nicht vom NATS-Server abgelehnt: $(docker logs "$NATS_STREAM_CLIENT_CONTAINER" 2>&1 || true)" >&2
  docker rm -fv "$NATS_STREAM_CLIENT_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

nats_stream_client_rejected_wrong_token=0
for _ in $(seq 1 20); do
  if docker logs "$NATS_STREAM_CLIENT_CONTAINER" 2>/dev/null | grep -F "REJECTED-WRONG-TOKEN" >/dev/null; then
    nats_stream_client_rejected_wrong_token=1
    break
  fi
  sleep 0.2
done
if [ "$nats_stream_client_rejected_wrong_token" -ne 1 ]; then
  echo "run-integration-tests: NATS-Vollinhalts-Stream-Rundlauf — Verbindungsversuch mit falschem Token wurde nicht vom NATS-Server abgelehnt: $(docker logs "$NATS_STREAM_CLIENT_CONTAINER" 2>&1 || true)" >&2
  docker rm -fv "$NATS_STREAM_CLIENT_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

# Die Zustellung ist Fire-and-Forget ohne Replay (ADR-0100 Teilfrage 3,
# dieselbe Broadcaster-Grenze wie bei gRPC/SSE oben): zwischen „Client hat
# abonniert" und „Publisher hat den Empfänger am Broadcaster registriert"
# liegt ein kurzes Fenster. Der Rundlauf committet deshalb eine begrenzte
# Folge von Zeilen mit IDs aus dem sonst unbelegten Bereich 301–305 (siehe
# Kommentar oben), bis der Client genau eine davon real empfangen hat.
nats_stream_received=0
for nats_stream_attempt in $(seq 1 5); do
  nats_stream_id=$((300 + nats_stream_attempt))
  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$NATS_STREAM_TABLE (id, name) VALUES ($nats_stream_id, '$NATS_STREAM_SENTINEL');
SQL
  for _ in $(seq 1 20); do
    if docker logs "$NATS_STREAM_CLIENT_CONTAINER" 2>/dev/null | grep -F "RECEIVED" >/dev/null; then
      nats_stream_received=1
      break
    fi
    if [ "$(docker inspect --format '{{.State.Running}}' "$NATS_STREAM_CLIENT_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
      break
    fi
    sleep 0.5
  done
  if [ "$nats_stream_received" -eq 1 ]; then
    break
  fi
done

nats_stream_client_output=$(docker logs "$NATS_STREAM_CLIENT_CONTAINER" 2>&1 || true)
docker rm -fv "$NATS_STREAM_CLIENT_CONTAINER" >/dev/null 2>&1 || true

if [ "$nats_stream_received" -ne 1 ]; then
  echo "run-integration-tests: NATS-Vollinhalts-Stream-Rundlauf — Test-Client empfing keine der committeten Änderungen ($NATS_STREAM_TABLE, $NATS_STREAM_SENTINEL) über das Subjekt $NATS_STREAM_SUBJECT: $nats_stream_client_output" >&2
  exit 1
fi
if ! printf '%s' "$nats_stream_client_output" | grep -qE "RECEIVED .*table=$NATS_STREAM_TABLE .*operation=INSERT .*new_image=.*$NATS_STREAM_SENTINEL"; then
  echo "run-integration-tests: NATS-Vollinhalts-Stream-Rundlauf — die RECEIVED-Zeile trägt nicht die erwartete Änderung ($NATS_STREAM_TABLE, INSERT, vollständiger Inhalt): $nats_stream_client_output" >&2
  exit 1
fi

# Unabhängiger SQL-Beleg, dass genau die empfangene Änderung real erfasst
# wurde: die change_id der RECEIVED-Zeile steht über den bestehenden
# Lesezugriffsweg in cdc.changes — der Stream-Empfang ist damit keine
# erfundene Ausgabe des Clients.
nats_stream_change_id=$(printf '%s' "$nats_stream_client_output" | grep -oE 'RECEIVED change_id=[^ ]+' | head -n1 | cut -d= -f2 || true)
if [ -z "$nats_stream_change_id" ]; then
  echo "run-integration-tests: NATS-Vollinhalts-Stream-Rundlauf — die RECEIVED-Zeile trägt keine change_id: $nats_stream_client_output" >&2
  exit 1
fi
nats_stream_captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
  "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$nats_stream_change_id' AND table_name = '$NATS_STREAM_TABLE' AND new_data->>'name' = '$NATS_STREAM_SENTINEL'")
if [ -z "$nats_stream_captured" ] || [ "$nats_stream_captured" -lt 1 ]; then
  echo "run-integration-tests: NATS-Vollinhalts-Stream-Rundlauf — die über den Stream empfangene Änderung (change_id=$nats_stream_change_id, $NATS_STREAM_SENTINEL) ist nicht real über cdc.changes lesbar (count=${nats_stream_captured:-leer})" >&2
  exit 1
fi

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
if [ "$feed_running" != "true" ]; then
  echo "run-integration-tests: Feed-Container lief nach dem NATS-Vollinhalts-Stream-Rundlauf nicht mehr weiter (kein Neustart erwartet)" >&2
  exit 1
fi

echo "run-integration-tests: NATS-Vollinhalts-Stream-Rundlauf (LH-FA-SST-008, ADR-0100) belegt — ein Wegwerf-Client (tools/harness/natsstreamsub) verband sich real mit gültigem Token über NATS, empfing eine danach committete Änderung als vollständiges JSON-Event über $NATS_STREAM_SUBJECT (Tabelle, Operation und Spaltenwert real am Event; die Feldvollständigkeit trägt publisher_test.go auf Unit-Ebene), deren change_id ($nats_stream_change_id) unabhängig über cdc.changes lesbar ist; Verbindungsversuche ohne und mit falschem Token wurden vom NATS-Server abgelehnt, und das bestehende Wecksignal (natssub) funktionierte mit demselben Test-Token unverändert weiter (siehe NATS-Happy-Path-/Negative-Belege oben): $nats_stream_client_output"

# --- Backfill-Rundläufe (LH-FA-CAP-009) ----------------------------
# Die Phasen laufen am laufenden Feed-Container und lesen und lösen
# ausschließlich über externe Wege: SQL-Funktionen und Views, HTTP,
# `docker exec`, `docker kill`. Drei Haltepunkte machen die Zeitpunkte eines
# Runs deterministisch: eine offene Schreibtransaktion mit
# Transaktionskennung hält die Slot-Anlage an (Run `running`, noch kein
# Snapshot); `docker pause` des Feed-Containers, während sie endet, öffnet
# ein Fenster nach dem Snapshot-Export; ein unbestätigter Schlüssel in
# `cdc.transaction` hält den Schreiber im zweiten Block an. Jede Sitzung
# hält kürzer als das Zeitlimit der Slot-Anlage (30 s) und wird von
# `bf_hold_end` beendet; jede Phase legt ihre eigene Tabelle an.
bf_fail() {
  echo "run-integration-tests: $*" >&2
  exit 1
}

bf_sql() {
  docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -tAc "$1"
}

# bf_expect <gelesen> <erwartet> <Beschreibung>
bf_expect() {
  [ "$1" = "$2" ] || bf_fail "$3 — erwartet '$2', gelesen '$1'"
}

# bf_await_sql <Anfrage> <erwartet> <Sekunden> <Beschreibung> [Takt in Sekunden]
bf_await_sql() {
  local query=$1 want=$2 seconds=$3 what=$4 step=${5:-0.25} got="" deadline
  deadline=$(( $(date +%s) + seconds ))
  while :; do
    got=$(bf_sql "$query")
    if [ "$got" = "$want" ]; then
      return 0
    fi
    if [ "$(date +%s)" -ge "$deadline" ]; then
      break
    fi
    sleep "$step"
  done
  bf_fail "$what — erwartet '$want', gelesen '${got:-leer}' (Anfrage: $query)"
}

# bf_await_applied <Antrags-ID> <Phase>
bf_await_applied() {
  local request_id=$1 phase=$2 status="" message=""
  for _ in $(seq 1 60); do
    status=$(bf_sql "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$request_id'")
    if [ "$status" = "applied" ]; then
      return 0
    fi
    if [ "$status" = "failed" ]; then
      message=$(bf_sql "SELECT error_message FROM cdc.administration_request WHERE administration_request_id = '$request_id'")
      bf_fail "$phase — Antrag $request_id endete failed: $message"
    fi
    sleep 0.5
  done
  bf_fail "$phase — Antrag $request_id wurde nicht innerhalb der Zeitspanne vermerkt (status=${status:-leer})"
}

# bf_enable <Tabelle> <Phase>
bf_enable() {
  local request_id
  request_id=$(bf_sql "SELECT cdc.enable_table('src-e2e', 'public', '$1')")
  [ -n "$request_id" ] || bf_fail "$2 — cdc.enable_table($1) lieferte keine Antrags-ID"
  bf_await_applied "$request_id" "$2"
}

# bf_request <Tabelle> <Phase> setzt bf_run_id auf die Kennung des Runs
# (die Antrags-ID) und wartet auf die Annahme.
bf_request() {
  bf_run_id=$(bf_sql "SELECT cdc.backfill_table('src-e2e', 'public', '$1')")
  [ -n "$bf_run_id" ] || bf_fail "$2 — cdc.backfill_table($1) lieferte keine Antrags-ID"
  bf_await_applied "$bf_run_id" "$2"
}

# bf_await_run <Run-ID> <Status> <Sekunden> <Phase>: ein anderer Endzustand
# als der erwartete endet die Phase.
bf_await_run() {
  local run_id=$1 want=$2 seconds=$3 phase=$4 status="" message="" i
  for ((i = 0; i < seconds * 4; i++)); do
    status=$(bf_sql "SELECT status FROM cdc.backfill_run WHERE run_id = '$run_id'")
    if [ "$status" = "$want" ]; then
      return 0
    fi
    case "$status" in
      completed | failed | interrupted)
        message=$(bf_sql "SELECT coalesce(error_message, '') FROM cdc.backfill_run WHERE run_id = '$run_id'")
        bf_fail "$phase — Run $run_id endete $status statt $want: $message"
        ;;
    esac
    sleep 0.25
  done
  bf_fail "$phase — Run $run_id erreichte $want nicht innerhalb von ${seconds}s (status=${status:-leer})"
}

# bf_await_healthy <Phase>
bf_await_healthy() {
  local health=""
  for _ in $(seq 1 60); do
    health=$(docker inspect --format '{{.State.Health.Status}}' "$FEED_CONTAINER" 2>/dev/null || echo fehlt)
    if [ "$health" = "healthy" ]; then
      return 0
    fi
    sleep 1
  done
  bf_fail "$1 — Feed-Container meldet Health-Status ${health:-fehlt}, wollen healthy"
}

# bf_hold_start <Sitzungsname> <SQL>: eine Hintergrund-Sitzung, die <SQL>
# ausführt (bf_hold_pid trägt den Prozess des Aufrufs); bf_hold_end beendet
# sie über ihren application_name.
bf_hold_start() {
  docker exec -i -e PGAPPNAME="$1" "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 >/dev/null 2>&1 <<<"$2" &
  bf_hold_pid=$!
}

bf_hold_end() {
  bf_sql "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE application_name = '$1'" >/dev/null
}

# bf_http <Modus und Argumente von tools/harness/httpclient>: der Client-Lauf
# im Toolchain-Container; ein Fehlschlag des Clients endet die Phase.
bf_http() {
  local output status
  set +e
  output=$(docker run --rm --network "$NETWORK" \
    -v "$(pwd)":/src:ro \
    -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
    -w /src \
    -e GOCACHE=/tmp/gocache \
    "$TOOLCHAIN_IMAGE" go run ./tools/harness/httpclient "$@" 2>&1)
  status=$?
  set -e
  if [ "$status" -ne 0 ]; then
    bf_fail "httpclient $1 endete mit Ausgang $status: $output"
  fi
  printf '%s\n' "$output"
}

# bf_read_ids <Ausgabe von bf_http changes>: die change_id der READ-Zeilen,
# sortiert und kommagetrennt.
bf_read_ids() {
  printf '%s\n' "$1" | grep '^READ ' | grep -oE 'change_id=[^ ]+' | cut -d= -f2 | LC_ALL=C sort | paste -sd,
}

abdeckung_declare "Backfill-Happy-Path" "LH-FA-CAP-009,LH-FA-REA-001" "eine aktivierte Tabelle mit Bestand wird über cdc.backfill_table als Run übernommen; jede Bestandszeile ist über cdc.changes und GET /changes als INSERT mit origin backfill lesbar, unterscheidbar von einer danach über den WAL-Pfad erfassten Änderung derselben Zeile (origin wal), und die change_id der HTTP-Antwort ist gegen cdc.changes gehalten" "Backfill-Happy-Path (LH-FA-CAP-009) belegt"

BF_PHASE="Backfill-Happy-Path"
BF_TABLE=feed_e2e_backfill

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$BF_TABLE (id int PRIMARY KEY, name text, note text);
INSERT INTO public.$BF_TABLE (id, name, note) VALUES
  (1, 'BackfillAlpha', 'n1'),
  (2, 'BackfillBeta', NULL),
  (3, 'BackfillGamma', 'n3'),
  (4, 'BackfillDelta', 'n4'),
  (5, 'BackfillEpsilon', 'n5');
SQL
bf_enable "$BF_TABLE" "$BF_PHASE"

# Ohne Backfill trägt das Log keine Zeile des Bestands.
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_TABLE'")" 0 "$BF_PHASE — Changes der Tabelle vor dem Backfill"

bf_request "$BF_TABLE" "$BF_PHASE"
bf_happy_run=$bf_run_id
bf_await_run "$bf_happy_run" completed 60 "$BF_PHASE"
bf_expect "$(bf_sql "SELECT status FROM cdc.backfill_status WHERE source_id = 'src-e2e' AND schema_name = 'public' AND table_name = '$BF_TABLE'")" completed "$BF_PHASE — Status über cdc.backfill_status"
bf_expect "$(bf_sql "SELECT rows_copied FROM cdc.backfill_run WHERE run_id = '$bf_happy_run'")" 5 "$BF_PHASE — rows_copied"
bf_position=$(bf_sql "SELECT snapshot_position FROM cdc.backfill_run WHERE run_id = '$bf_happy_run'")
[ -n "$bf_position" ] || bf_fail "$BF_PHASE — der abgeschlossene Run trägt keine Snapshot-Position"

# Lesen über den SQL-Weg: fünf INSERT-Changes der Herkunft backfill ohne
# old_data, ein Block (Transaktion 0bf-<Run>-00000001), alle auf der
# Snapshot-Position, jede Bestandszeile einmal.
bf_where="source_id = 'src-e2e' AND table_name = '$BF_TABLE' AND origin = 'backfill'"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE $bf_where AND operation = 'INSERT' AND old_data IS NULL")" 5 "$BF_PHASE — INSERT-Changes der Herkunft backfill"
bf_expect "$(bf_sql "SELECT string_agg(DISTINCT transaction_id, ',') FROM cdc.changes WHERE $bf_where")" "0bf-$bf_happy_run-00000001" "$BF_PHASE — Transaktions-Kennung des Blocks"
bf_expect "$(bf_sql "SELECT string_agg(DISTINCT commit_position::text, ',') FROM cdc.changes WHERE $bf_where")" "$bf_position" "$BF_PHASE — Commit-Position der Backfill-Changes"
bf_expect "$(bf_sql "SELECT string_agg(new_data->>'id', ',' ORDER BY (new_data->>'id')::int) FROM cdc.changes WHERE $bf_where")" "1,2,3,4,5" "$BF_PHASE — Schlüssel des Bestands"
bf_expected_ids=""
for n in 1 2 3 4 5; do
  bf_expected_ids="${bf_expected_ids:+$bf_expected_ids,}0bf-$bf_happy_run-00000001-$n"
done
bf_expect "$(bf_sql "SELECT string_agg(change_id, ',' ORDER BY sequence) FROM cdc.changes WHERE $bf_where")" "$bf_expected_ids" "$BF_PHASE — change_id nach Sequenz"
bf_expect "$(bf_sql "SELECT jsonb_exists(new_data, 'note') FROM cdc.changes WHERE $bf_where AND new_data->>'id' = '2'")" f "$BF_PHASE — NULL-Spalte im Row Image der Zeile 2"
bf_expect "$(bf_sql "SELECT new_data->>'note' FROM cdc.changes WHERE $bf_where AND new_data->>'id' = '3'")" n3 "$BF_PHASE — Spaltenwert der Zeile 3"

# Lesen über HTTP: der Bereich [X, X+1) trägt genau den Bestand.
bf_http_backfill=$(bf_http changes "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" src-e2e public "$BF_TABLE" "$bf_position" "$((bf_position + 1))")
bf_expect "$(printf '%s\n' "$bf_http_backfill" | grep -c '^READ .* operation=INSERT .* origin=backfill$')" 5 "$BF_PHASE — READ-Zeilen INSERT/backfill über GET /changes"
bf_expect "$(printf '%s\n' "$bf_http_backfill" | grep -c " commit_position=$bf_position ")" 5 "$BF_PHASE — Commit-Position der READ-Zeilen"
bf_expect "$(bf_read_ids "$bf_http_backfill")" "$(bf_sql "SELECT string_agg(change_id, ',' ORDER BY change_id COLLATE \"C\") FROM cdc.changes WHERE $bf_where")" "$BF_PHASE — change_id über GET /changes gegen cdc.changes"

# Eine danach über den WAL-Pfad erfasste Änderung derselben Zeile trägt
# origin wal und liegt hinter der Snapshot-Position.
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
UPDATE public.$BF_TABLE SET name = 'BackfillGammaWal' WHERE id = 3;
SQL
bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_TABLE' AND origin = 'wal'" 1 30 "$BF_PHASE — WAL-Change der Zeile 3"
bf_expect "$(bf_sql "SELECT string_agg(origin || ':' || operation, ',' ORDER BY commit_position, transaction_id, sequence) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_TABLE' AND new_data->>'id' = '3'")" "backfill:INSERT,wal:UPDATE" "$BF_PHASE — Herkunft und Operation der Zeile 3 in Lese-Ordnung"
bf_wal_position=$(bf_sql "SELECT commit_position FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_TABLE' AND origin = 'wal'")
[ "$bf_wal_position" -gt "$bf_position" ] || bf_fail "$BF_PHASE — die WAL-Change liegt auf Position $bf_wal_position, nicht hinter der Snapshot-Position $bf_position"
bf_http_all=$(bf_http changes "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" src-e2e public "$BF_TABLE" "" "")
bf_expect "$(printf '%s\n' "$bf_http_all" | grep -c ' origin=backfill$')" 5 "$BF_PHASE — READ-Zeilen der Herkunft backfill im ganzen Log"
bf_expect "$(printf '%s\n' "$bf_http_all" | grep -c ' origin=wal$')" 1 "$BF_PHASE — READ-Zeilen der Herkunft wal im ganzen Log"
printf '%s\n' "$bf_http_all" | grep '^READ ' | tail -n1 | grep -qF "BackfillGammaWal" || bf_fail "$BF_PHASE — die letzte READ-Zeile ist nicht die WAL-Änderung der Zeile 3: $bf_http_all"

bf_feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
bf_expect "$bf_feed_running" true "$BF_PHASE — Feed-Container läuft weiter"

echo "run-integration-tests: Backfill-Happy-Path (LH-FA-CAP-009) belegt — Run $bf_happy_run übernahm 5 Bestandszeilen von $BF_TABLE auf Position $bf_position (Transaktion 0bf-$bf_happy_run-00000001), lesbar über cdc.changes und GET /changes als INSERT mit origin backfill, change_id gegen cdc.changes gehalten ($bf_expected_ids); die danach erfasste Änderung der Zeile 3 trägt origin wal auf Position $bf_wal_position: $bf_http_all"

abdeckung_declare "Backfill-Schema-Version-Beleg" "LH-FA-CAP-009,LH-FA-SCH-005" "die Schema-Version einer Backfill-Change ist die zum Antrag aktuelle Zeile aus cdc.schema_version; nach einer Spalten-Erweiterung ohne Änderung über den WAL-Pfad bleibt sie unverändert, und das Row Image trägt die neue Spalte" "Backfill-Schema-Version-Beleg (LH-FA-CAP-009) belegt"

BF_PHASE="Backfill-Schema-Version-Beleg"
BF_VERSION_TABLE=feed_e2e_backfill_version

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$BF_VERSION_TABLE (id int PRIMARY KEY, name text);
INSERT INTO public.$BF_VERSION_TABLE (id, name) VALUES (1, 'VersionAlpha'), (2, 'VersionBeta');
SQL
bf_enable "$BF_VERSION_TABLE" "$BF_PHASE"

bf_version_id_sql="SELECT sv.schema_version_id FROM cdc.schema_version sv JOIN cdc.source_table st ON st.source_table_id = sv.source_table_id WHERE st.source_id = 'src-e2e' AND st.schema_name = 'public' AND st.table_name = '$BF_VERSION_TABLE' ORDER BY sv.version DESC LIMIT 1"
bf_version_before=$(bf_sql "$bf_version_id_sql")
[ -n "$bf_version_before" ] || bf_fail "$BF_PHASE — die Aktivierung hat keine Schema-Version für $BF_VERSION_TABLE angelegt"

# Die Erweiterung erzeugt keine Change über den WAL-Pfad: keine Relation-
# Nachricht, keine neue Version.
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
ALTER TABLE public.$BF_VERSION_TABLE ADD COLUMN extra text DEFAULT 'neu';
SQL

bf_request "$BF_VERSION_TABLE" "$BF_PHASE"
bf_version_run=$bf_run_id
bf_await_run "$bf_version_run" completed 60 "$BF_PHASE"
bf_version_where="source_id = 'src-e2e' AND table_name = '$BF_VERSION_TABLE' AND origin = 'backfill'"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE $bf_version_where")" 2 "$BF_PHASE — Backfill-Changes der Tabelle"
bf_expect "$(bf_sql "SELECT string_agg(DISTINCT schema_version, ',') FROM cdc.changes WHERE $bf_version_where")" "$bf_version_before" "$BF_PHASE — Schema-Version der Backfill-Changes"
bf_expect "$(bf_sql "$bf_version_id_sql")" "$bf_version_before" "$BF_PHASE — höchste registrierte Version nach dem Run"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE $bf_version_where AND new_data->>'extra' = 'neu'")" 2 "$BF_PHASE — neue Spalte im Row Image"
bf_version_shape=$(bf_sql "SELECT count(*) FROM cdc.table_schema WHERE schema_version_id = '$bf_version_before'")

echo "run-integration-tests: Backfill-Schema-Version-Beleg (LH-FA-CAP-009) belegt — Run $bf_version_run auf $BF_VERSION_TABLE nach ADD COLUMN ohne WAL-Change: alle Backfill-Changes tragen die zum Antrag aktuelle Version $bf_version_before, das Row Image trägt die Spalte extra; die Version führt $bf_version_shape Zeilen in cdc.table_schema"

abdeckung_declare "Backfill-Startposition" "LH-FA-CON-005,LH-FA-CAP-009" "die Anfangsposition eines frisch registrierten Consumers wird über den CLI-Weg und den HTTP-Weg gemessen und gegen die Position des Bestands gehalten: sie liegt vor der Snapshot-Position, der Bestand ist ab ihr lesbar; ein Consumer mit bestätigter Position hinter dem Bestand liest ihn über seinen Fortschritt nicht" "Backfill-Startposition (LH-FA-CON-005) belegt"

BF_PHASE="Backfill-Startposition"
BF_CONSUMER=bf-e2e-start-consumer

exec_feed register-consumer "$BF_CONSUMER" || bf_fail "$BF_PHASE — register-consumer ($BF_CONSUMER) endete mit einem Fehler"
bf_start_out=$(bf_http position "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" "$BF_CONSUMER")
bf_start_offset=$(printf '%s\n' "$bf_start_out" | grep -oE '"offset":[0-9]+' | cut -d: -f2)
bf_expect "$bf_start_offset" 0 "$BF_PHASE — Anfangsposition (offset) über GET /consumers/position: $bf_start_out"
printf '%s\n' "$bf_start_out" | grep -qF '"acknowledged":false' || bf_fail "$BF_PHASE — der frisch registrierte Consumer trägt eine bestätigte Position: $bf_start_out"
bf_status_row=$(bf_sql "SELECT coalesce(acknowledged_position::text, 'NULL') || '|' || coalesce(latest_commit_position::text, 'NULL') FROM cdc.consumer_status WHERE consumer_id = '$BF_CONSUMER'")
bf_expect "${bf_status_row%%|*}" NULL "$BF_PHASE — bestätigte Position in cdc.consumer_status"
[ "$bf_start_offset" -lt "$bf_position" ] || bf_fail "$BF_PHASE — die Anfangsposition $bf_start_offset liegt nicht vor der Snapshot-Position $bf_position"

# Ab der Anfangsposition ist der Bestand lesbar: jede Position hinter ihr,
# darunter die Snapshot-Position, gehört zum Fortschritt des Consumers.
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_TABLE' AND origin = 'backfill' AND commit_position > $bf_start_offset")" 5 "$BF_PHASE — Bestand hinter der Anfangsposition"

# Bestätigt der Consumer eine Position hinter dem Bestand, liegt der Bestand
# vor seiner Position und erscheint nicht in seinem Fortschritt.
bf_latest_position=$(bf_sql "SELECT max(commit_position) FROM cdc.transaction WHERE source_id = 'src-e2e'")
exec_feed acknowledge-consumer "$BF_CONSUMER" "$bf_latest_position" || bf_fail "$BF_PHASE — acknowledge-consumer ($BF_CONSUMER, $bf_latest_position) endete mit einem Fehler"
bf_ack_out=$(bf_http position "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" "$BF_CONSUMER")
printf '%s\n' "$bf_ack_out" | grep -qF "\"offset\":$bf_latest_position" || bf_fail "$BF_PHASE — die bestätigte Position ist nicht $bf_latest_position: $bf_ack_out"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_TABLE' AND origin = 'backfill' AND commit_position > $bf_latest_position")" 0 "$BF_PHASE — Bestand hinter der bestätigten Position"
bf_http remove "$HTTP_BASE_URL" "$HTTP_TOKEN_ADMIN" "$BF_CONSUMER" >/dev/null
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.consumer WHERE consumer_id = '$BF_CONSUMER'")" 0 "$BF_PHASE — Consumer nach der Entfernung"

echo "run-integration-tests: Backfill-Startposition (LH-FA-CON-005) belegt — der frisch registrierte Consumer $BF_CONSUMER startet an Position $bf_start_offset (acknowledged=false; cdc.consumer_status trägt bestätigte und letzte Commit-Position als $bf_status_row, NULL heißt leer), die Snapshot-Position des Bestands ist $bf_position: ab der Anfangsposition sind 5 Backfill-Changes lesbar; nach der Bestätigung von $bf_latest_position liegen 0 hinter ihr: $bf_start_out"

abdeckung_declare "Backfill-Boundary (leere Tabelle, zweiter Antrag)" "LH-FA-CAP-009" "eine leere Tabelle endet completed mit 0 Zeilen und schreibt keine Transaktion; ein zweiter Antrag für dieselbe Tabelle bei aktivem Run endet failed mit Grund, ohne Run-Zeile, der erste Run läuft davon unberührt zu Ende" "Backfill-Boundary (leere Tabelle, zweiter Antrag) belegt"

BF_PHASE="Backfill-Boundary"
BF_EMPTY_TABLE=feed_e2e_backfill_empty
BF_DUP_TABLE=feed_e2e_backfill_dup

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$BF_EMPTY_TABLE (id int PRIMARY KEY, name text);
CREATE TABLE public.$BF_DUP_TABLE (id int PRIMARY KEY, name text);
INSERT INTO public.$BF_DUP_TABLE (id, name) VALUES (1, 'DupAlpha'), (2, 'DupBeta');
SQL
bf_enable "$BF_EMPTY_TABLE" "$BF_PHASE"
bf_enable "$BF_DUP_TABLE" "$BF_PHASE"

bf_request "$BF_EMPTY_TABLE" "$BF_PHASE"
bf_empty_run=$bf_run_id
bf_await_run "$bf_empty_run" completed 60 "$BF_PHASE"
bf_expect "$(bf_sql "SELECT rows_copied FROM cdc.backfill_run WHERE run_id = '$bf_empty_run'")" 0 "$BF_PHASE — rows_copied der leeren Tabelle"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.transaction WHERE transaction_id LIKE '0bf-$bf_empty_run-%'")" 0 "$BF_PHASE — Transaktionen des Runs der leeren Tabelle"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_EMPTY_TABLE'")" 0 "$BF_PHASE — Changes der leeren Tabelle"

# Haltepunkt 1: die offene Schreibtransaktion hält die Slot-Anlage an. Der
# Run steht dann `running` ohne Snapshot-Position — der Zustandswechsel
# liegt vor dem Öffnen des Snapshots.
bf_hold_start e2e-bf-xid $'BEGIN;\nSELECT pg_current_xact_id();\nSELECT pg_sleep(120);'
bf_await_sql "SELECT count(*) FROM pg_stat_activity WHERE application_name = 'e2e-bf-xid' AND backend_xid IS NOT NULL" 1 30 "$BF_PHASE — Haltetransaktion"
bf_request "$BF_DUP_TABLE" "$BF_PHASE"
bf_dup_first=$bf_run_id
bf_await_run "$bf_dup_first" running 30 "$BF_PHASE"
bf_expect "$(bf_sql "SELECT coalesce(snapshot_position::text, 'NULL') FROM cdc.backfill_run WHERE run_id = '$bf_dup_first'")" NULL "$BF_PHASE — Snapshot-Position des Runs an der Haltetransaktion"

bf_dup_second=$(bf_sql "SELECT cdc.backfill_table('src-e2e', 'public', '$BF_DUP_TABLE')")
[ -n "$bf_dup_second" ] || bf_fail "$BF_PHASE — der zweite Antrag lieferte keine Antrags-ID"
bf_await_sql "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$bf_dup_second'" failed 30 "$BF_PHASE — zweiter Antrag bei aktivem Run"
bf_dup_error=$(bf_sql "SELECT error_message FROM cdc.administration_request WHERE administration_request_id = '$bf_dup_second'")
printf '%s' "$bf_dup_error" | grep -qF "aktiver Backfill-Run" || bf_fail "$BF_PHASE — der Fehlertext des zweiten Antrags nennt den aktiven Run nicht: $bf_dup_error"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.backfill_run WHERE run_id = '$bf_dup_second'")" 0 "$BF_PHASE — Run-Zeile des zweiten Antrags"
bf_expect "$(bf_sql "SELECT status FROM cdc.backfill_run WHERE run_id = '$bf_dup_first'")" running "$BF_PHASE — erster Run nach dem zweiten Antrag"

bf_hold_end e2e-bf-xid
bf_await_run "$bf_dup_first" completed 60 "$BF_PHASE"
bf_expect "$(bf_sql "SELECT rows_copied FROM cdc.backfill_run WHERE run_id = '$bf_dup_first'")" 2 "$BF_PHASE — rows_copied des ersten Runs"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.backfill_run WHERE source_id = 'src-e2e' AND table_name = '$BF_DUP_TABLE'")" 1 "$BF_PHASE — Run-Zeilen der Tabelle"

echo "run-integration-tests: Backfill-Boundary (leere Tabelle, zweiter Antrag) belegt — Run $bf_empty_run auf der leeren Tabelle endete completed mit 0 Zeilen und schrieb keine Transaktion; der zweite Antrag $bf_dup_second bei aktivem Run $bf_dup_first endete failed ($bf_dup_error) ohne Run-Zeile, der erste Run stand an der Haltetransaktion running ohne Snapshot-Position und endete danach completed mit 2 Zeilen"

abdeckung_declare "Backfill-Regelstand (rename_column, Ausschluss, Nichtanwendbarkeit)" "LH-FA-CFG-007,LH-FA-CAP-009,LH-QA-SEC-004" "für eine Tabelle mit rename_column-Regel trägt ein Backfill-Run den Bestand über cdc.changes und GET /changes mit umbenanntem Schlüssel und origin backfill, in derselben Schlüsselmenge wie die danach über den WAL-Pfad erfasste Zeile; mit zusätzlichem exclude_column auf der umbenannten Spalte trägt kein Bild des nächsten Runs Quellnamen, Zielnamen oder Wert; eine im Run nicht anwendbare Regel (Zielname nach der Regel als Spalte angelegt) endet den Run failed mit der Klasse schema ohne Change, der Erfassungspfad läuft weiter, und nach dem Entfernen der Regel übernimmt ein neuer Antrag den Bestand in Rohform" "Backfill-Regelstand (LH-FA-CFG-007) belegt"

BF_PHASE="Backfill-Regelstand"
BF_RULE_TABLE=feed_e2e_backfill_rule
BF_RULE_BAD_TABLE=feed_e2e_backfill_rulebad
BF_RULE_SPEC='{"kind":"rename_column","column":"name","to":"customer_name"}'

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$BF_RULE_TABLE (id int PRIMARY KEY, name text, note text);
INSERT INTO public.$BF_RULE_TABLE (id, name, note) VALUES
  (1, 'RegelAlpha', 'r1'),
  (2, 'RegelBeta', NULL),
  (3, NULL, 'r3');
SQL
bf_enable "$BF_RULE_TABLE" "$BF_PHASE"
bf_rule_request=$(bf_sql "SELECT cdc.set_transformation('src-e2e', 'public', '$BF_RULE_TABLE', 'kundenname', '$BF_RULE_SPEC')")
[ -n "$bf_rule_request" ] || bf_fail "$BF_PHASE — cdc.set_transformation($BF_RULE_TABLE) lieferte keine Antrags-ID"
bf_await_applied "$bf_rule_request" "$BF_PHASE"

# Run mit Regel: der Bestand trägt den umbenannten Schlüssel; die Zeile ohne
# Wert in `name` trägt weder Quell- noch Zielschlüssel.
bf_request "$BF_RULE_TABLE" "$BF_PHASE"
bf_rule_run=$bf_run_id
bf_await_run "$bf_rule_run" completed 60 "$BF_PHASE"
bf_rule_where="source_id = 'src-e2e' AND table_name = '$BF_RULE_TABLE' AND origin = 'backfill' AND transaction_id LIKE '0bf-$bf_rule_run-%'"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE $bf_rule_where AND operation = 'INSERT'")" 3 "$BF_PHASE — INSERT-Changes des Runs mit Regel"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE $bf_rule_where AND jsonb_exists(new_data, 'name')")" 0 "$BF_PHASE — Bilder mit dem Quellschlüssel name"
bf_expect "$(bf_sql "SELECT string_agg(new_data->>'customer_name', ',' ORDER BY (new_data->>'id')::int) FROM cdc.changes WHERE $bf_rule_where AND jsonb_exists(new_data, 'customer_name')")" "RegelAlpha,RegelBeta" "$BF_PHASE — Werte unter dem Zielschlüssel customer_name"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE $bf_rule_where AND jsonb_exists(new_data, 'customer_name')")" 2 "$BF_PHASE — Bilder mit dem Zielschlüssel (die NULL-Zeile trägt keinen)"
bf_expect "$(bf_sql "SELECT string_agg(new_data->>'note', ',' ORDER BY (new_data->>'id')::int) FROM cdc.changes WHERE $bf_rule_where AND jsonb_exists(new_data, 'note')")" "r1,r3" "$BF_PHASE — nicht umbenannte Spalte note"

bf_rule_position=$(bf_sql "SELECT snapshot_position FROM cdc.backfill_run WHERE run_id = '$bf_rule_run'")
bf_rule_http=$(bf_http changes "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" src-e2e public "$BF_RULE_TABLE" "$bf_rule_position" "$((bf_rule_position + 1))")
bf_expect "$(printf '%s\n' "$bf_rule_http" | grep -c '^READ .* operation=INSERT .* origin=backfill$')" 3 "$BF_PHASE — READ-Zeilen INSERT/backfill über GET /changes"
bf_expect "$(printf '%s\n' "$bf_rule_http" | grep -c 'customer_name')" 2 "$BF_PHASE — READ-Zeilen mit dem Zielschlüssel über GET /changes"
bf_expect "$(printf '%s\n' "$bf_rule_http" | grep '^READ ' | grep -c '"name"')" 0 "$BF_PHASE — READ-Zeilen mit dem Quellschlüssel über GET /changes"
bf_expect "$(bf_read_ids "$bf_rule_http")" "$(bf_sql "SELECT string_agg(change_id, ',' ORDER BY change_id COLLATE \"C\") FROM cdc.changes WHERE $bf_rule_where")" "$BF_PHASE — change_id über GET /changes gegen cdc.changes"

# Dieselbe Regel im WAL-Pfad: die danach erfasste Zeile trägt dieselbe
# Schlüsselmenge wie die Backfill-Change der Zeile 1.
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$BF_RULE_TABLE (id, name, note) VALUES (4, 'RegelDelta', 'r4');
SQL
bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_RULE_TABLE' AND origin = 'wal'" 1 30 "$BF_PHASE — WAL-Change der Zeile 4"
bf_rule_keys_wal=$(bf_sql "SELECT string_agg(k, ',' ORDER BY k) FROM cdc.changes c, jsonb_object_keys(c.new_data) k WHERE c.source_id = 'src-e2e' AND c.table_name = '$BF_RULE_TABLE' AND c.origin = 'wal'")
bf_rule_keys_backfill=$(bf_sql "SELECT string_agg(k, ',' ORDER BY k) FROM cdc.changes c, jsonb_object_keys(c.new_data) k WHERE $bf_rule_where AND c.new_data->>'id' = '1'")
bf_expect "$bf_rule_keys_wal" "customer_name,id,note" "$BF_PHASE — Schlüsselmenge der WAL-Change"
bf_expect "$bf_rule_keys_backfill" "$bf_rule_keys_wal" "$BF_PHASE — Schlüsselmenge der Backfill-Change gegen die der WAL-Change"

# Ausschluss der umbenannten Spalte (LH-QA-SEC-004): das Bild des nächsten
# Runs trägt weder den Quellnamen noch den Zielnamen noch einen Wert von
# name; der frühere Run bleibt unverändert lesbar.
bf_exclude_request=$(bf_sql "SELECT cdc.exclude_column('src-e2e', 'public', '$BF_RULE_TABLE', 'name')")
[ -n "$bf_exclude_request" ] || bf_fail "$BF_PHASE — cdc.exclude_column($BF_RULE_TABLE) lieferte keine Antrags-ID"
bf_await_applied "$bf_exclude_request" "$BF_PHASE"
bf_request "$BF_RULE_TABLE" "$BF_PHASE"
bf_rule_run_excluded=$bf_run_id
bf_await_run "$bf_rule_run_excluded" completed 60 "$BF_PHASE"
bf_excluded_where="source_id = 'src-e2e' AND table_name = '$BF_RULE_TABLE' AND origin = 'backfill' AND transaction_id LIKE '0bf-$bf_rule_run_excluded-%'"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE $bf_excluded_where")" 4 "$BF_PHASE — Changes des Runs mit Ausschluss"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE $bf_excluded_where AND (jsonb_exists(new_data, 'name') OR jsonb_exists(new_data, 'customer_name') OR new_data::text LIKE '%Regel%')")" 0 "$BF_PHASE — Bilder mit Quellname, Zielname oder Wert der ausgeschlossenen Spalte"
bf_expect "$(bf_sql "SELECT string_agg(new_data->>'note', ',' ORDER BY (new_data->>'id')::int) FROM cdc.changes WHERE $bf_excluded_where AND jsonb_exists(new_data, 'note')")" "r1,r3,r4" "$BF_PHASE — nicht ausgeschlossene Spalte note im Run mit Ausschluss"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE $bf_rule_where AND jsonb_exists(new_data, 'customer_name')")" 2 "$BF_PHASE — der frühere Run vor dem Ausschluss bleibt unverändert lesbar"

# Nichtanwendbarkeit im Run: die Regel war beim Antrag anwendbar; die danach
# angelegte Spalte customer_name macht den Zielnamen zur Kollision. Der Run
# endet failed mit der Klasse schema, ohne Change; der Erfassungspfad läuft
# weiter (kein Fehlerzustand im Lebenszeichen, eine danach in eine andere
# Tabelle geschriebene Zeile wird erfasst).
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$BF_RULE_BAD_TABLE (id int PRIMARY KEY, name text);
INSERT INTO public.$BF_RULE_BAD_TABLE (id, name) VALUES (1, 'BadAlpha'), (2, 'BadBeta');
SQL
bf_enable "$BF_RULE_BAD_TABLE" "$BF_PHASE"
bf_bad_rule_request=$(bf_sql "SELECT cdc.set_transformation('src-e2e', 'public', '$BF_RULE_BAD_TABLE', 'kundenname', '$BF_RULE_SPEC')")
[ -n "$bf_bad_rule_request" ] || bf_fail "$BF_PHASE — cdc.set_transformation($BF_RULE_BAD_TABLE) lieferte keine Antrags-ID"
bf_await_applied "$bf_bad_rule_request" "$BF_PHASE"
bf_sql "ALTER TABLE public.$BF_RULE_BAD_TABLE ADD COLUMN customer_name text" >/dev/null

bf_request "$BF_RULE_BAD_TABLE" "$BF_PHASE"
bf_bad_run=$bf_run_id
bf_await_sql "SELECT status FROM cdc.backfill_run WHERE run_id = '$bf_bad_run'" failed 60 "$BF_PHASE — Run mit nicht anwendbarer Regel"
bf_bad_error=$(bf_sql "SELECT error_message FROM cdc.backfill_run WHERE run_id = '$bf_bad_run'")
printf '%s' "$bf_bad_error" | grep -q '^schema: ' || bf_fail "$BF_PHASE — der Fehlertext beginnt nicht mit der Klasse schema: $bf_bad_error"
printf '%s' "$bf_bad_error" | grep -qF 'kundenname' || bf_fail "$BF_PHASE — der Fehlertext nennt den Regelnamen nicht: $bf_bad_error"
bf_expect "$(bf_sql "SELECT status FROM cdc.backfill_status WHERE source_id = 'src-e2e' AND schema_name = 'public' AND table_name = '$BF_RULE_BAD_TABLE'")" failed "$BF_PHASE — Status über cdc.backfill_status"
bf_expect "$(bf_sql "SELECT rows_copied FROM cdc.backfill_run WHERE run_id = '$bf_bad_run'")" 0 "$BF_PHASE — rows_copied des fehlgeschlagenen Runs"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_RULE_BAD_TABLE'")" 0 "$BF_PHASE — Changes der Tabelle nach dem fehlgeschlagenen Run"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.transaction WHERE transaction_id LIKE '0bf-$bf_bad_run-%'")" 0 "$BF_PHASE — Transaktionen des fehlgeschlagenen Runs"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.heartbeat WHERE source_id = 'src-e2e' AND error_class IS NULL")" 1 "$BF_PHASE — Lebenszeichen der Quelle ohne Fehlerzustand (der Run-Fehler ist run-lokal)"
bf_expect "$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)" true "$BF_PHASE — Feed-Container läuft weiter"
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$BF_TABLE (id, name, note) VALUES (6, 'RegelSonde', 'n6');
SQL
bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_TABLE' AND origin = 'wal' AND new_data->>'id' = '6'" 1 30 "$BF_PHASE — WAL-Change der Sonden-Zeile nach dem fehlgeschlagenen Run"

# Abhilfe: Regel entfernen, neuer Antrag — ohne Neustart des Feed-Containers.
bf_remove_request=$(bf_sql "SELECT cdc.remove_transformation('src-e2e', 'public', '$BF_RULE_BAD_TABLE', 'kundenname')")
[ -n "$bf_remove_request" ] || bf_fail "$BF_PHASE — cdc.remove_transformation($BF_RULE_BAD_TABLE) lieferte keine Antrags-ID"
bf_await_applied "$bf_remove_request" "$BF_PHASE"
bf_request "$BF_RULE_BAD_TABLE" "$BF_PHASE"
bf_bad_retry=$bf_run_id
bf_await_run "$bf_bad_retry" completed 60 "$BF_PHASE"
bf_expect "$(bf_sql "SELECT rows_copied FROM cdc.backfill_run WHERE run_id = '$bf_bad_retry'")" 2 "$BF_PHASE — rows_copied des neuen Antrags nach der Abhilfe"
bf_expect "$(bf_sql "SELECT string_agg(new_data->>'name', ',' ORDER BY (new_data->>'id')::int) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_RULE_BAD_TABLE' AND origin = 'backfill' AND transaction_id LIKE '0bf-$bf_bad_retry-%'")" "BadAlpha,BadBeta" "$BF_PHASE — Bestand in Rohform nach dem Entfernen der Regel"

echo "run-integration-tests: Backfill-Regelstand (LH-FA-CFG-007) belegt — Run $bf_rule_run übernahm 3 Zeilen von $BF_RULE_TABLE mit umbenanntem Schlüssel customer_name (Werte $(bf_sql "SELECT string_agg(new_data->>'customer_name', ',' ORDER BY (new_data->>'id')::int) FROM cdc.changes WHERE $bf_rule_where AND jsonb_exists(new_data, 'customer_name')"), Schlüsselmenge $bf_rule_keys_backfill gleich der WAL-Change), lesbar über cdc.changes und GET /changes; Run $bf_rule_run_excluded nach exclude_column auf name trägt in 4 Bildern weder name noch customer_name noch einen Wert von name; Run $bf_bad_run auf $BF_RULE_BAD_TABLE endete failed ($bf_bad_error) ohne Change, der Erfassungspfad lief weiter, der neue Antrag $bf_bad_retry nach dem Entfernen der Regel übernahm 2 Zeilen in Rohform"

# Replay-Invariante: die Schreiber, die Haltetransaktion und die Anwendung des
# Logs liegen in TestE2EBackfillReplayInvariant
# (test/integration/backfill_e2e_test.go); die Abdeckungszeile entsteht aus
# dem Doc-Kommentar der Testfunktion.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_INTEGRATION_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test -v -count=1 -run '^TestE2EBackfillReplayInvariant$' ./test/integration/...

echo "run-integration-tests: TestE2EBackfillReplayInvariant (LH-FA-CAP-009) belegt — Zeile 'Replay-Invariante:' der Testausgabe oben"

abdeckung_declare "Backfill-DDL-Fenster" "LH-FA-CAP-009" "ändert eine DDL zwischen Snapshot-Export und Snapshot-Import die Tabelle, endet der Run failed ohne Change und der Erfassungspfad läuft weiter: ein DROP COLUMN mit der Fehlerklasse storage, ein Umschreiben der Tabelle (ALTER COLUMN TYPE) mit der Fehlerklasse transient; ein neuer Antrag übernimmt den Bestand vollständig in der geänderten Form" "Backfill-DDL-Fenster (LH-FA-CAP-009) belegt"

BF_PHASE="Backfill-DDL-Fenster"
BF_DDL_TABLE=feed_e2e_backfill_ddl
BF_REWRITE_TABLE=feed_e2e_backfill_rewrite

# bf_ddl_window <Tabelle> <DDL-Anweisung> <Fehlerklasse>: beantragt einen
# Backfill der Tabelle, führt die DDL im Fenster zwischen Snapshot-Export und
# Snapshot-Import aus und prüft, dass der Run failed mit der Klasse endet,
# ohne Change und ohne Fehlerzustand des Erfassungspfads. Sie setzt
# bf_ddl_run, bf_ddl_error und bf_pause_ms.
#
# Die Haltetransaktion hält die Slot-Anlage an, der Run steht dann hinter der
# Vorbedingungsprüfung. Der Feed-Container ist angehalten (docker pause),
# während die Haltetransaktion endet: der Run liest die Antwort der
# Slot-Anlage erst nach dem Fortsetzen, der Snapshot ist dann exportiert und
# noch nicht importiert. Die DDL-Sitzung wartet, bis die Slot-Anlage fertig
# ist (der Walsender des Slots wartet auf den nächsten Befehl), führt in diesem
# Fenster die DDL aus und committet; danach läuft der Container weiter. Die
# Pause bleibt unter der Hälfte von wal_sender_timeout (2 s, compose.yaml),
# innerhalb derer PostgreSQL die Verbindung des Erfassungspfads beendet.
bf_ddl_window() {
  local table=$1 ddl=$2 class=$3 pause_start
  bf_hold_start e2e-bf-xid $'BEGIN;\nSELECT pg_current_xact_id();\nSELECT pg_sleep(120);'
  bf_await_sql "SELECT count(*) FROM pg_stat_activity WHERE application_name = 'e2e-bf-xid' AND backend_xid IS NOT NULL" 1 30 "$BF_PHASE $table — Haltetransaktion"
  bf_request "$table" "$BF_PHASE"
  bf_ddl_run=$bf_run_id
  bf_await_run "$bf_ddl_run" running 30 "$BF_PHASE $table"
  bf_hold_start e2e-bf-ddl "BEGIN;
DO \$\$
DECLARE deadline timestamptz := clock_timestamp() + interval '30 seconds';
BEGIN
  WHILE NOT EXISTS (SELECT 1 FROM pg_replication_slots s JOIN pg_stat_activity a ON a.pid = s.active_pid WHERE s.slot_name LIKE 'cdc\_bf\_%' AND a.state <> 'active') LOOP
    IF clock_timestamp() > deadline THEN
      RAISE EXCEPTION 'die Slot-Anlage des Runs wurde nicht fertig';
    END IF;
    PERFORM pg_sleep(0.01);
    PERFORM pg_stat_clear_snapshot();
  END LOOP;
END \$\$;
$ddl;
COMMIT;"
  local ddl_pid=$bf_hold_pid
  bf_await_sql "SELECT count(*) FROM pg_stat_activity WHERE application_name = 'e2e-bf-ddl' AND state = 'active'" 1 30 "$BF_PHASE $table — DDL-Sitzung wartet"
  pause_start=$(date +%s%N)
  docker pause "$FEED_CONTAINER" >/dev/null
  bf_hold_end e2e-bf-xid
  wait "$ddl_pid" || bf_fail "$BF_PHASE $table — die DDL-Sitzung endete mit einem Fehler"
  docker unpause "$FEED_CONTAINER" >/dev/null
  bf_pause_ms=$(( ($(date +%s%N) - pause_start) / 1000000 ))
  [ "$bf_pause_ms" -lt 1000 ] || bf_fail "$BF_PHASE $table — die Pause des Feed-Containers dauerte ${bf_pause_ms} ms; erlaubt sind unter 1000 ms, die Hälfte von wal_sender_timeout (2000 ms)"
  bf_await_run "$bf_ddl_run" failed 60 "$BF_PHASE $table"
  bf_ddl_error=$(bf_sql "SELECT error_message FROM cdc.backfill_run WHERE run_id = '$bf_ddl_run'")
  printf '%s' "$bf_ddl_error" | grep -q "^$class: " || bf_fail "$BF_PHASE $table — der Fehlertext beginnt nicht mit der Klasse $class: $bf_ddl_error"
  bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$table'")" 0 "$BF_PHASE $table — Changes der Tabelle nach dem fehlgeschlagenen Run"
  bf_expect "$(bf_sql "SELECT count(*) FROM cdc.transaction WHERE transaction_id LIKE '0bf-$bf_ddl_run-%'")" 0 "$BF_PHASE $table — Transaktionen des fehlgeschlagenen Runs"
  bf_expect "$(bf_sql "SELECT count(*) FROM cdc.heartbeat WHERE source_id = 'src-e2e' AND error_class IS NULL")" 1 "$BF_PHASE $table — Lebenszeichen der Quelle ohne Fehlerzustand (der Run-Fehler ist run-lokal)"
  bf_expect "$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)" true "$BF_PHASE $table — Feed-Container läuft weiter"
}

# Entfernte Spalte: DECLARE scheitert am Katalog (Klasse storage).
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$BF_DDL_TABLE (id int PRIMARY KEY, name text, note text);
INSERT INTO public.$BF_DDL_TABLE (id, name, note) VALUES (1, 'DdlAlpha', 'n1'), (2, 'DdlBeta', 'n2'), (3, 'DdlGamma', 'n3');
CREATE TABLE public.$BF_REWRITE_TABLE (id int PRIMARY KEY, name text);
INSERT INTO public.$BF_REWRITE_TABLE (id, name) VALUES (1, 'RewAlpha'), (2, 'RewBeta'), (3, 'RewGamma');
SQL
bf_enable "$BF_DDL_TABLE" "$BF_PHASE"
bf_enable "$BF_REWRITE_TABLE" "$BF_PHASE"

bf_ddl_window "$BF_DDL_TABLE" "ALTER TABLE public.$BF_DDL_TABLE DROP COLUMN note" storage
bf_ddl_run_drop=$bf_ddl_run
bf_ddl_error_drop=$bf_ddl_error
bf_pause_ms_drop=$bf_pause_ms

# Abhilfe: ein neuer Antrag übernimmt den Bestand in der geänderten Form.
bf_request "$BF_DDL_TABLE" "$BF_PHASE"
bf_ddl_retry=$bf_run_id
bf_await_run "$bf_ddl_retry" completed 60 "$BF_PHASE"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_DDL_TABLE' AND origin = 'backfill' AND NOT jsonb_exists(new_data, 'note')")" 3 "$BF_PHASE — Bestand ohne die entfernte Spalte nach dem neuen Antrag"

# Umgeschriebene Tabelle: der ältere Snapshot sieht die neue Datei leer; die
# Sperre mit dem Filenode-Vergleich beendet den Run (Klasse transient), er
# schließt nicht mit 0 Zeilen ab.
bf_ddl_window "$BF_REWRITE_TABLE" "ALTER TABLE public.$BF_REWRITE_TABLE ALTER COLUMN name TYPE varchar(64)" transient
bf_expect "$(bf_sql "SELECT rows_copied FROM cdc.backfill_run WHERE run_id = '$bf_ddl_run'")" 0 "$BF_PHASE — rows_copied des abgebrochenen Runs nach dem Umschreiben"
bf_ddl_run_rewrite=$bf_ddl_run
bf_ddl_error_rewrite=$bf_ddl_error
bf_pause_ms_rewrite=$bf_pause_ms

bf_request "$BF_REWRITE_TABLE" "$BF_PHASE"
bf_rewrite_retry=$bf_run_id
bf_await_run "$bf_rewrite_retry" completed 60 "$BF_PHASE"
bf_expect "$(bf_sql "SELECT rows_copied FROM cdc.backfill_run WHERE run_id = '$bf_rewrite_retry'")" 3 "$BF_PHASE — rows_copied des neuen Antrags nach dem Umschreiben"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_REWRITE_TABLE' AND origin = 'backfill' AND transaction_id LIKE '0bf-$bf_rewrite_retry-%'")" 3 "$BF_PHASE — Bestand der umgeschriebenen Tabelle nach dem neuen Antrag"
bf_expect "$(bf_sql "SELECT string_agg(new_data->>'name', ',' ORDER BY (new_data->>'id')::int) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_REWRITE_TABLE' AND origin = 'backfill'")" "RewAlpha,RewBeta,RewGamma" "$BF_PHASE — Werte des Bestands der umgeschriebenen Tabelle"

echo "run-integration-tests: Backfill-DDL-Fenster (LH-FA-CAP-009) belegt — DROP COLUMN zwischen Snapshot-Export und Snapshot-Import endete Run $bf_ddl_run_drop failed ($bf_ddl_error_drop) ohne Change (Pause des Feed-Containers ${bf_pause_ms_drop} ms), der neue Antrag $bf_ddl_retry übernahm 3 Zeilen ohne die entfernte Spalte; ALTER COLUMN TYPE endete Run $bf_ddl_run_rewrite failed ($bf_ddl_error_rewrite) ohne Change (Pause ${bf_pause_ms_rewrite} ms), der neue Antrag $bf_rewrite_retry übernahm 3 Zeilen; der Erfassungspfad lief in beiden Fällen weiter"

abdeckung_declare "Backfill-Negative (docker kill, queued-Aufnahme)" "LH-FA-CAP-009" "ein docker kill des Feed-Containers mitten im Run hinterlässt keine sichtbare Change des Runs, kein Slot und keine Sitzung bleiben; nach dem Neustart steht der Run interrupted, ein erneuter Antrag übernimmt den Bestand einmal und vollständig, und ein zum Abbruchzeitpunkt queued wartender Run wird ohne neuen Antrag aufgenommen und ausgeführt" "Backfill-Negative (docker kill, queued-Aufnahme) belegt"

BF_PHASE="Backfill-Negative"
BF_KILL_TABLE=feed_e2e_backfill_kill
BF_QUEUED_TABLE=feed_e2e_backfill_queued
BF_KILL_ROWS=2500

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$BF_KILL_TABLE (id int PRIMARY KEY, name text);
INSERT INTO public.$BF_KILL_TABLE (id, name) SELECT g, 'Kill-' || g FROM generate_series(1, $BF_KILL_ROWS) g;
CREATE TABLE public.$BF_QUEUED_TABLE (id int PRIMARY KEY, name text);
INSERT INTO public.$BF_QUEUED_TABLE (id, name) SELECT g, 'Queued-' || g FROM generate_series(1, 4) g;
SQL
bf_enable "$BF_KILL_TABLE" "$BF_PHASE"
bf_enable "$BF_QUEUED_TABLE" "$BF_PHASE"

# Der Run soll erst im zweiten Block stehen bleiben (Blockgröße 1000). Die
# Haltetransaktion hält die Slot-Anlage an; während sie endet, ist der
# Feed-Container angehalten (docker pause), und die Konflikt-Sitzung legt in
# diesem Fenster die Transaktionskennung des zweiten Blocks unbestätigt an
# (die Slot-Anlage wartet auf jede offene Schreibtransaktion).
# Nach dem Fortsetzen schreibt der Run den ersten Block und wartet im zweiten
# auf die Konflikt-Sitzung. Die Pause bleibt unter der Hälfte von
# wal_sender_timeout (2 s, compose.yaml).
bf_hold_start e2e-bf-xid $'BEGIN;\nSELECT pg_current_xact_id();\nSELECT pg_sleep(120);'
bf_await_sql "SELECT count(*) FROM pg_stat_activity WHERE application_name = 'e2e-bf-xid' AND backend_xid IS NOT NULL" 1 30 "$BF_PHASE — Haltetransaktion"
bf_request "$BF_KILL_TABLE" "$BF_PHASE"
bf_kill_run=$bf_run_id
bf_await_run "$bf_kill_run" running 30 "$BF_PHASE"

# Der zweite Antrag gegen eine andere Tabelle wird angenommen und wartet
# hinter dem ersten Run.
bf_request "$BF_QUEUED_TABLE" "$BF_PHASE"
bf_queued_run=$bf_run_id
bf_expect "$(bf_sql "SELECT status FROM cdc.backfill_run WHERE run_id = '$bf_queued_run'")" queued "$BF_PHASE — Status des zweiten Runs hinter dem ersten"

bf_hold_start e2e-bf-conflict "BEGIN;
DO \$\$
DECLARE deadline timestamptz := clock_timestamp() + interval '30 seconds';
BEGIN
  WHILE NOT EXISTS (SELECT 1 FROM pg_replication_slots s JOIN pg_stat_activity a ON a.pid = s.active_pid WHERE s.slot_name LIKE 'cdc\_bf\_%' AND a.state <> 'active') LOOP
    IF clock_timestamp() > deadline THEN
      RAISE EXCEPTION 'die Slot-Anlage des Runs wurde nicht fertig';
    END IF;
    PERFORM pg_sleep(0.01);
    PERFORM pg_stat_clear_snapshot();
  END LOOP;
END \$\$;
INSERT INTO cdc.transaction (transaction_id, source_id, commit_position, committed_at) VALUES ('0bf-$bf_kill_run-00000002', 'src-e2e', 1, now());
SELECT pg_sleep(120);"
bf_await_sql "SELECT count(*) FROM pg_stat_activity WHERE application_name = 'e2e-bf-conflict' AND state = 'active'" 1 30 "$BF_PHASE — Konflikt-Sitzung wartet"
bf_pause_start=$(date +%s%N)
docker pause "$FEED_CONTAINER" >/dev/null
bf_hold_end e2e-bf-xid
bf_await_sql "SELECT count(*) FROM pg_stat_activity WHERE application_name = 'e2e-bf-conflict' AND backend_xid IS NOT NULL" 1 30 "$BF_PHASE — Konflikt-Zeile des zweiten Blocks" 0.02
docker unpause "$FEED_CONTAINER" >/dev/null
bf_pause_ms=$(( ($(date +%s%N) - bf_pause_start) / 1000000 ))
[ "$bf_pause_ms" -lt 1000 ] || bf_fail "$BF_PHASE — die Pause des Feed-Containers dauerte ${bf_pause_ms} ms; erlaubt sind unter 1000 ms, die Hälfte von wal_sender_timeout (2000 ms)"

bf_await_sql "SELECT rows_copied FROM cdc.backfill_run WHERE run_id = '$bf_kill_run'" 1000 60 "$BF_PHASE — Fortschritt nach dem ersten Block"
bf_await_sql "SELECT count(*) FROM pg_stat_activity WHERE wait_event_type = 'Lock' AND wait_event = 'transactionid'" 1 30 "$BF_PHASE — Schreiber wartet im zweiten Block"
bf_expect "$(bf_sql "SELECT status FROM cdc.backfill_run WHERE run_id = '$bf_kill_run'")" running "$BF_PHASE — Status des Runs im zweiten Block"

docker kill "$FEED_CONTAINER" >/dev/null
bf_expect "$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)" false "$BF_PHASE — Feed-Container nach docker kill"
bf_hold_end e2e-bf-conflict

# Nach dem Abbruch bleibt weder eine Sitzung noch ein Slot des Runs zurück,
# und der Run hat nichts sichtbar gemacht.
bf_await_sql "SELECT count(*) FROM pg_stat_activity WHERE backend_type IN ('client backend', 'walsender') AND pid <> pg_backend_pid()" 0 60 "$BF_PHASE — Sitzungen nach dem Abbruch"
bf_expect "$(bf_sql "SELECT count(*) FROM pg_replication_slots WHERE slot_name LIKE 'cdc\_bf\_%'")" 0 "$BF_PHASE — temporäre Slots des Runs nach dem Abbruch"
bf_await_sql "SELECT active FROM pg_replication_slots WHERE slot_name = '$SLOT'" f 60 "$BF_PHASE — Haupt-Slot nach dem Abbruch"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_KILL_TABLE'")" 0 "$BF_PHASE — sichtbare Changes des abgebrochenen Runs"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.transaction WHERE transaction_id LIKE '0bf-$bf_kill_run-%'")" 0 "$BF_PHASE — Transaktionen des abgebrochenen Runs"
bf_expect "$(bf_sql "SELECT status FROM cdc.backfill_run WHERE run_id = '$bf_kill_run'")" running "$BF_PHASE — Status des abgebrochenen Runs vor dem Neustart"

docker start "$FEED_CONTAINER" >/dev/null
bf_await_healthy "$BF_PHASE"

bf_await_run "$bf_kill_run" interrupted 30 "$BF_PHASE"
bf_expect "$(bf_sql "SELECT finished_at IS NOT NULL FROM cdc.backfill_run WHERE run_id = '$bf_kill_run'")" t "$BF_PHASE — finished_at des unterbrochenen Runs"
bf_expect "$(bf_sql "SELECT rows_copied FROM cdc.backfill_run WHERE run_id = '$bf_kill_run'")" 1000 "$BF_PHASE — rows_copied des unterbrochenen Runs (letzter festgehaltener Fortschritt)"

# Die zum Abbruchzeitpunkt queued wartende Zeile wird ohne neuen Antrag
# aufgenommen und ausgeführt.
bf_await_run "$bf_queued_run" completed 90 "$BF_PHASE"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_QUEUED_TABLE' AND origin = 'backfill' AND transaction_id LIKE '0bf-$bf_queued_run-%'")" 4 "$BF_PHASE — Bestand des aufgenommenen Runs"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.backfill_run WHERE source_id = 'src-e2e' AND table_name = '$BF_KILL_TABLE'")" 1 "$BF_PHASE — kein automatischer Neustart des unterbrochenen Runs"

# Ein erneuter Antrag beginnt neu: der Bestand ist danach einmal und
# vollständig lesbar.
bf_request "$BF_KILL_TABLE" "$BF_PHASE"
bf_retry_run=$bf_run_id
bf_await_run "$bf_retry_run" completed 90 "$BF_PHASE"
bf_kill_where="source_id = 'src-e2e' AND table_name = '$BF_KILL_TABLE' AND origin = 'backfill'"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE $bf_kill_where")" "$BF_KILL_ROWS" "$BF_PHASE — Backfill-Changes nach dem erneuten Antrag"
bf_expect "$(bf_sql "SELECT count(DISTINCT new_data->>'id') FROM cdc.changes WHERE $bf_kill_where")" "$BF_KILL_ROWS" "$BF_PHASE — verschiedene Schlüssel"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE $bf_kill_where AND transaction_id NOT LIKE '0bf-$bf_retry_run-%'")" 0 "$BF_PHASE — Changes fremder Runs"
bf_expect "$(bf_sql "SELECT count(DISTINCT transaction_id) FROM cdc.changes WHERE $bf_kill_where")" 3 "$BF_PHASE — Blöcke des erneuten Runs"

# Der Erfassungspfad setzt nach dem Neustart fort.
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$BF_KILL_TABLE (id, name) VALUES ($((BF_KILL_ROWS + 1)), 'KillAfterRestart');
SQL
bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$BF_KILL_TABLE' AND origin = 'wal'" 1 30 "$BF_PHASE — Erfassung nach dem Neustart"

echo "run-integration-tests: Backfill-Negative (docker kill, queued-Aufnahme) belegt — docker kill im zweiten Block von Run $bf_kill_run (rows_copied 1000, running): keine Sitzung und kein Slot blieben zurück, keine Change des Runs war sichtbar; nach dem Neustart steht der Run interrupted, der queued wartende Run $bf_queued_run wurde ohne neuen Antrag aufgenommen und endete completed (4 Zeilen); der erneute Antrag $bf_retry_run übernahm $BF_KILL_ROWS Zeilen einmal (3 Blöcke), die Erfassung setzte nach dem Neustart fort"

abdeckung_declare "Leerlauf-Bestätigung (WAL ohne Inhalt für die Publication, Fehlerschwelle)" "LH-FA-CAP-009,LH-QA-REL-001" "ein Backfill-Run, dessen WAL die Fehlerschwelle des WAL-Rückstands (Override der Konfigurationsdatei) überschreitet, und ein Schreiber auf eine nicht aktivierte Tabelle, der in Stücken unter der Warnschwelle schreibt (Summe über der Fehlerschwelle, nach jedem Stück bestätigt der Slot), beenden den Feed-Container nicht: der Run endet completed, der Bestand ist über cdc.changes lesbar, der Container läuft weiter" "Leerlauf-Bestätigung (LH-FA-CAP-009, LH-QA-REL-001) belegt"

# Der Feed-Container läuft in dieser Phase mit einer Konfigurationsdatei, die
# die Schwellen des WAL-Rückstands senkt (`wal_retention_*_bytes` haben kein
# Env-Gegenstück, SPEC-013): eine Compose-Override-Datei im Temp-Verzeichnis
# des Runners hängt sie an den Dienst und setzt `CDC_CONFIG_FILE`;
# `compose.yaml` bleibt unverändert, und der Container wird am Phasen-Ende
# ohne Override wiederhergestellt. Die Größe der Last folgt aus der Messung:
# die Phase liest das WAL, das der Run erzeugt hat, und bricht ab, wenn es die
# Fehlerschwelle nicht übersteigt. Der Schreiber auf die nicht aktivierte
# Tabelle schreibt in Stücken: bis die Bestätigung eines Stücks folgt, steht das
# Stück im Rückstand des Slots; das WAL je Stück liegt unter der Warnschwelle,
# die Summe der Stücke über der Fehlerschwelle. Der Runner wartet nach jedem
# Stück, bis `confirmed_flush_lsn` die WAL-Position hinter dem Stück erreicht
# hat.
BF_PHASE="Leerlauf-Bestätigung"
WAL_TABLE=feed_e2e_wal_backfill
WAL_FOREIGN=feed_e2e_wal_foreign
WAL_ROWS=30000
WAL_FOREIGN_CHUNKS=6
WAL_FOREIGN_CHUNK_ROWS=10000
WAL_FOREIGN_CONFIRM_SECONDS=30
WAL_WARN_BYTES=4194304
WAL_ERROR_BYTES=8388608
WAL_WAIT_SECONDS=12

WAL_TMP=$(mktemp -d)
chmod 0755 "$WAL_TMP"
printf 'wal_retention_warn_bytes: %s\nwal_retention_error_bytes: %s\n' "$WAL_WARN_BYTES" "$WAL_ERROR_BYTES" > "$WAL_TMP/config.yaml"
cat > "$WAL_TMP/override.yaml" <<YAML
services:
  pg-change-feed:
    environment:
      CDC_CONFIG_FILE: /etc/cdc/e2e-wal-config.yaml
    volumes:
      - $WAL_TMP/config.yaml:/etc/cdc/e2e-wal-config.yaml:ro
YAML
chmod 0644 "$WAL_TMP/config.yaml" "$WAL_TMP/override.yaml"

$COMPOSE -f "$WAL_TMP/override.yaml" up -d --force-recreate --no-deps pg-change-feed >/dev/null
bf_await_healthy "$BF_PHASE"
bf_expect "$(docker inspect --format '{{range .Config.Env}}{{println .}}{{end}}' "$FEED_CONTAINER" | grep -c '^CDC_CONFIG_FILE=/etc/cdc/e2e-wal-config.yaml$')" 1 "$BF_PHASE — Konfigurationsdatei im Feed-Container"
wal_feed_started=$(docker inspect --format '{{.State.StartedAt}}' "$FEED_CONTAINER")

# bf_wal_hold <Sekunden>: der Feed-Container läuft über die gesamte Zeitspanne
# (mindestens zwei Prüf-Takte des WAL-Rückstands, `heartbeatInterval` 5 s), ohne
# neu gestartet worden zu sein.
bf_wal_hold() {
  local i
  for ((i = 0; i < $1; i++)); do
    [ "$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)" = true ] \
      || bf_fail "$BF_PHASE — $2: der Feed-Container endete (Ausgang $(docker inspect --format '{{.State.ExitCode}}' "$FEED_CONTAINER" 2>/dev/null || echo unbekannt)), Log-Ende: $(docker logs --tail 3 "$FEED_CONTAINER" 2>&1 | tr '\n' ' ')"
    sleep 1
  done
  bf_expect "$(docker inspect --format '{{.State.StartedAt}}' "$FEED_CONTAINER")" "$wal_feed_started" "$BF_PHASE — $2: Startzeitpunkt des Feed-Containers (kein Neustart)"
}

# Backfill-Run über mehr WAL, als die Fehlerschwelle trägt.
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$WAL_TABLE (id int PRIMARY KEY, name text);
INSERT INTO public.$WAL_TABLE (id, name) SELECT g, 'Wal-' || g || repeat('x', 60) FROM generate_series(1, $WAL_ROWS) g;
CREATE TABLE public.$WAL_FOREIGN (id bigint PRIMARY KEY, payload text);
SQL
bf_enable "$WAL_TABLE" "$BF_PHASE"
wal_run_before=$(bf_sql "SELECT pg_current_wal_lsn()")
bf_request "$WAL_TABLE" "$BF_PHASE"
wal_run_id=$bf_run_id
bf_await_run "$wal_run_id" completed 120 "$BF_PHASE"
wal_run_after=$(bf_sql "SELECT pg_current_wal_lsn()")
wal_run_bytes=$(bf_sql "SELECT pg_wal_lsn_diff('$wal_run_after', '$wal_run_before')::bigint")
[ "$wal_run_bytes" -gt "$WAL_ERROR_BYTES" ] || bf_fail "$BF_PHASE — der Run erzeugte $wal_run_bytes B WAL, nicht mehr als die Fehlerschwelle $WAL_ERROR_BYTES B: die Last trägt den Beleg nicht"
bf_wal_hold "$WAL_WAIT_SECONDS" "nach dem Run"
bf_expect "$(bf_sql "SELECT rows_copied FROM cdc.backfill_run WHERE run_id = '$wal_run_id'")" "$WAL_ROWS" "$BF_PHASE — rows_copied"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$WAL_TABLE' AND origin = 'backfill'")" "$WAL_ROWS" "$BF_PHASE — Bestand über cdc.changes"

# Schreiber auf eine nicht aktivierte Tabelle: in Stücken, ohne Inhalt für die
# Publication. Jedes Stück trägt eigene Schlüssel; die Summe des WAL liegt über
# der Fehlerschwelle, jedes Stück unter der Warnschwelle.
wal_foreign_bytes=0
wal_foreign_chunk_max=0
for ((wal_chunk = 1; wal_chunk <= WAL_FOREIGN_CHUNKS; wal_chunk++)); do
  wal_chunk_first=$(((wal_chunk - 1) * WAL_FOREIGN_CHUNK_ROWS + 1))
  wal_chunk_last=$((wal_chunk * WAL_FOREIGN_CHUNK_ROWS))
  wal_chunk_before=$(bf_sql "SELECT pg_current_wal_lsn()")
  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$WAL_FOREIGN SELECT g, repeat('x', 130) FROM generate_series($wal_chunk_first, $wal_chunk_last) g;
SQL
  wal_chunk_after=$(bf_sql "SELECT pg_current_wal_lsn()")
  wal_chunk_bytes=$(bf_sql "SELECT pg_wal_lsn_diff('$wal_chunk_after', '$wal_chunk_before')::bigint")
  [ "$wal_chunk_bytes" -lt "$WAL_WARN_BYTES" ] || bf_fail "$BF_PHASE — Stück $wal_chunk von $WAL_FOREIGN_CHUNKS erzeugte $wal_chunk_bytes B WAL, nicht weniger als die Warnschwelle $WAL_WARN_BYTES B: das Stück ist ein Stoß, keine Bestätigungs-Runde"
  [ "$wal_chunk_bytes" -le "$wal_foreign_chunk_max" ] || wal_foreign_chunk_max=$wal_chunk_bytes
  wal_foreign_bytes=$((wal_foreign_bytes + wal_chunk_bytes))
  bf_await_sql "SELECT pg_wal_lsn_diff(confirmed_flush_lsn, '$wal_chunk_after') >= 0 FROM pg_replication_slots WHERE slot_name = '$SLOT'" t "$WAL_FOREIGN_CONFIRM_SECONDS" "$BF_PHASE — confirmed_flush_lsn des Slots erreichte die Position hinter Stück $wal_chunk von $WAL_FOREIGN_CHUNKS nicht innerhalb von ${WAL_FOREIGN_CONFIRM_SECONDS} s" 0.05
done
[ "$wal_foreign_bytes" -gt "$WAL_ERROR_BYTES" ] || bf_fail "$BF_PHASE — die $WAL_FOREIGN_CHUNKS Stücke erzeugten zusammen $wal_foreign_bytes B WAL, nicht mehr als die Fehlerschwelle $WAL_ERROR_BYTES B: die Last trägt den Beleg nicht"
bf_wal_hold "$WAL_WAIT_SECONDS" "nach dem Schreiber auf die nicht aktivierte Tabelle"

# Der Erfassungspfad nimmt danach weiter Changes an.
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$WAL_TABLE (id, name) VALUES ($((WAL_ROWS + 1)), 'WalAfter');
SQL
bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$WAL_TABLE' AND origin = 'wal'" 1 30 "$BF_PHASE — Erfassung nach beiden Lasten"

# Der Rückstand aus dem Log des Containers: keine Fehlerschwellen-Zeile, die
# gemessene Spitze steht in der Ausgabe.
wal_log=$(docker logs "$FEED_CONTAINER" 2>&1)
bf_expect "$(printf '%s\n' "$wal_log" | grep -c 'WAL-Rückstand über Fehlerschwelle')" 0 "$BF_PHASE — Fehlerschwellen-Zeilen im Log des Feed-Containers"
wal_peak=$(printf '%s\n' "$wal_log" | grep '"metric":"cdc_wal_retention_bytes"' | grep -o '"bytes":[0-9]*' | cut -d: -f2 | sort -n | tail -n1 || true)
[ -n "$wal_peak" ] || bf_fail "$BF_PHASE — das Log des Feed-Containers trägt keine Messung des WAL-Rückstands"

# Der Container läuft danach wieder mit der Konfiguration von compose.yaml.
$COMPOSE up -d --force-recreate --no-deps pg-change-feed >/dev/null
bf_await_healthy "$BF_PHASE"
rm -rf "$WAL_TMP"

echo "run-integration-tests: Leerlauf-Bestätigung (LH-FA-CAP-009, LH-QA-REL-001) belegt — bei den Schwellen Warn $WAL_WARN_BYTES B / Fehler $WAL_ERROR_BYTES B (Konfigurationsdatei) endete der Backfill-Run $wal_run_id über $WAL_ROWS Zeilen completed und erzeugte $wal_run_bytes B WAL, ein Schreiber auf die nicht aktivierte Tabelle $WAL_FOREIGN erzeugte in $WAL_FOREIGN_CHUNKS Stücken zu je $WAL_FOREIGN_CHUNK_ROWS Zeilen zusammen $wal_foreign_bytes B WAL (höchstes Stück $wal_foreign_chunk_max B, unter der Warnschwelle; nach jedem Stück erreichte confirmed_flush_lsn die Position hinter dem Stück, Frist ${WAL_FOREIGN_CONFIRM_SECONDS} s); der Feed-Container lief über beide Lasten weiter (je ${WAL_WAIT_SECONDS} s beobachtet, kein Neustart), der Bestand ist über cdc.changes lesbar, der höchste der im Log des Feed-Containers gemessenen Rückstände (Proben im 5-s-Takt, keine Spitze) liegt bei $wal_peak B"

abdeckung_declare "Fehlerschwelle des WAL-Rückstands beendet den Container" "LH-QA-REL-001,LH-QA-REL-003" "bei gehaltener Persistierung (exklusive Sperre auf cdc.change) und WAL ohne Inhalt für die Publication über der Fehlerschwelle beendet der Feed-Container den Lauf mit Ausgang 1 und der Abbruch-Zeile im Log und meldet einen Fehlerzustand der Klasse replication über cdc.process_heartbeat; die gehaltene Persistierung allein beendet ihn nicht, und die wartende Transaktion ist nach dem Neustart über cdc.changes lesbar" "Fehlerschwelle beendet den Container (LH-QA-REL-001, LH-QA-REL-003) belegt"

# Die Gegenseite der vorigen Phase: dort bestätigt der Stream im Leerlauf und
# entlastet den Rückstand; hier hält eine Sitzung des Runners `cdc.change`
# exklusiv, die Persistierung der nächsten Transaktion einer aktivierten
# Tabelle wartet auf die Sperre (der Runner liest die wartende Anweisung in
# `pg_stat_activity`), und ein Schreiber auf eine nicht aktivierte Tabelle
# bringt den Rückstand über die Fehlerschwelle — der Fehlerzustand trägt
# danach die Klasse `replication` (`ADR-0049`). Der Container läuft wie in
# der vorigen Phase mit der Konfigurationsdatei der gesenkten Schwellen
# (`WAL_WARN_BYTES`/`WAL_ERROR_BYTES`, `bf_wal_hold`, `wal_feed_started`) und
# wird am Phasen-Ende ohne sie wiederhergestellt.
BF_PHASE="Fehlerschwelle beendet den Container"
WAL_STOP_TABLE=feed_e2e_wal_stop
WAL_STOP_FOREIGN=feed_e2e_wal_stop_foreign
WAL_STOP_HOLD=pgc-e2e-persist-hold
WAL_STOP_WAIT_SECONDS=90

WAL_TMP=$(mktemp -d)
chmod 0755 "$WAL_TMP"
printf 'wal_retention_warn_bytes: %s\nwal_retention_error_bytes: %s\n' "$WAL_WARN_BYTES" "$WAL_ERROR_BYTES" > "$WAL_TMP/config.yaml"
cat > "$WAL_TMP/override.yaml" <<YAML
services:
  pg-change-feed:
    environment:
      CDC_CONFIG_FILE: /etc/cdc/e2e-wal-config.yaml
    volumes:
      - $WAL_TMP/config.yaml:/etc/cdc/e2e-wal-config.yaml:ro
YAML
chmod 0644 "$WAL_TMP/config.yaml" "$WAL_TMP/override.yaml"

$COMPOSE -f "$WAL_TMP/override.yaml" up -d --force-recreate --no-deps pg-change-feed >/dev/null
bf_await_healthy "$BF_PHASE"
bf_expect "$(docker inspect --format '{{range .Config.Env}}{{println .}}{{end}}' "$FEED_CONTAINER" | grep -c '^CDC_CONFIG_FILE=/etc/cdc/e2e-wal-config.yaml$')" 1 "$BF_PHASE — Konfigurationsdatei im Feed-Container"
wal_feed_started=$(docker inspect --format '{{.State.StartedAt}}' "$FEED_CONTAINER")

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$WAL_STOP_TABLE (id int PRIMARY KEY, name text);
CREATE TABLE public.$WAL_STOP_FOREIGN (id bigint PRIMARY KEY, payload text);
SQL
bf_enable "$WAL_STOP_TABLE" "$BF_PHASE"

# Die Persistierung hält an: die Sperre steht, bevor die Transaktion eintrifft.
bf_hold_start "$WAL_STOP_HOLD" "BEGIN; LOCK TABLE cdc.change IN ACCESS EXCLUSIVE MODE; SELECT pg_sleep(180);"
bf_await_sql "SELECT count(*) FROM pg_locks l JOIN pg_stat_activity a ON a.pid = l.pid WHERE a.application_name = '$WAL_STOP_HOLD' AND l.mode = 'AccessExclusiveLock' AND l.granted" 1 15 "$BF_PHASE — Sperre auf cdc.change"
bf_sql "INSERT INTO public.$WAL_STOP_TABLE (id, name) VALUES (1, 'WalStopHeld')" >/dev/null
bf_await_sql "SELECT count(*) FROM pg_stat_activity WHERE wait_event_type = 'Lock' AND application_name <> '$WAL_STOP_HOLD' AND query LIKE '%INSERT INTO cdc.change%'" 1 30 "$BF_PHASE — wartende Persistierung des Feed-Containers"

# Nullprobe: die gehaltene Persistierung allein erzeugt keinen Rückstand über
# der Fehlerschwelle, der Container läuft über mindestens zwei Prüf-Takte.
bf_wal_hold "$WAL_WAIT_SECONDS" "gehaltene Persistierung ohne Last"

wal_stop_before=$(bf_sql "SELECT pg_current_wal_lsn()")
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
INSERT INTO public.$WAL_STOP_FOREIGN SELECT g, repeat('x', 130) FROM generate_series(1, 60000) g;
SQL
wal_stop_bytes=$(bf_sql "SELECT pg_wal_lsn_diff(pg_current_wal_lsn(), '$wal_stop_before')::bigint")
[ "$wal_stop_bytes" -gt "$WAL_ERROR_BYTES" ] || bf_fail "$BF_PHASE — der Schreiber erzeugte $wal_stop_bytes B WAL, nicht mehr als die Fehlerschwelle $WAL_ERROR_BYTES B: die Last trägt den Beleg nicht"

wal_stop_loaded=$(date +%s)
wal_stop_ended=0
for ((i = 0; i < WAL_STOP_WAIT_SECONDS; i++)); do
  if [ "$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)" = false ]; then
    wal_stop_ended=1
    break
  fi
  sleep 1
done
[ "$wal_stop_ended" -eq 1 ] || bf_fail "$BF_PHASE — der Feed-Container lief ${WAL_STOP_WAIT_SECONDS} s nach einer Last von $wal_stop_bytes B WAL über der Fehlerschwelle $WAL_ERROR_BYTES B weiter"
wal_stop_seconds=$(( $(date +%s) - wal_stop_loaded ))
bf_expect "$(docker inspect --format '{{.State.ExitCode}}' "$FEED_CONTAINER")" 1 "$BF_PHASE — Ausgang des Feed-Containers"

wal_stop_log=$(docker logs "$FEED_CONTAINER" 2>&1)
bf_expect "$(printf '%s\n' "$wal_stop_log" | grep -c 'WAL-Rückstand über Fehlerschwelle — kontrollierter Abbruch')" 1 "$BF_PHASE — Abbruch-Zeilen im Log des Feed-Containers"
wal_stop_logged=$(printf '%s\n' "$wal_stop_log" | grep 'WAL-Rückstand über Fehlerschwelle — kontrollierter Abbruch' | grep -o '"bytes":[0-9]*' | cut -d: -f2)
[ -n "$wal_stop_logged" ] && [ "$wal_stop_logged" -gt "$WAL_ERROR_BYTES" ] || bf_fail "$BF_PHASE — der gemessene Rückstand der Abbruch-Zeile (${wal_stop_logged:-leer} B) liegt nicht über der Fehlerschwelle $WAL_ERROR_BYTES B"
bf_expect "$(printf '%s\n' "$wal_stop_log" | grep -c 'Lauf beendet mit Fehler')" 1 "$BF_PHASE — Ende-Zeilen im Log des Feed-Containers"

# Der Fehlerzustand ist über die Heartbeat-Projektion sichtbar und trägt die
# Klasse `replication` (`ADR-0049`): die Ausgangs-Klasse-Regel in
# `mergeStreamAndWALFaultOutcome` lässt den WAL-Schwellen-Fehler durch, auch
# wenn der Stream mit einem Fehler endet, dessen Kette denselben Abbruch
# trägt (die gehaltene Persistierung oben).
wal_stop_class=$(bf_sql "SELECT coalesce(error_class, '') FROM cdc.process_heartbeat WHERE source_id = 'src-e2e'")
bf_expect "$wal_stop_class" "replication" "$BF_PHASE — Fehlerklasse in cdc.process_heartbeat"

# Die Persistierung wird frei, der Container läuft wieder mit der
# Konfiguration von compose.yaml, und die wartende Transaktion ist gelesen.
bf_hold_end "$WAL_STOP_HOLD"
$COMPOSE up -d --force-recreate --no-deps pg-change-feed >/dev/null
bf_await_healthy "$BF_PHASE"
rm -rf "$WAL_TMP"
bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$WAL_STOP_TABLE' AND new_data->>'id' = '1'" 1 60 "$BF_PHASE — Change der wartenden Transaktion nach dem Neustart"

echo "run-integration-tests: Fehlerschwelle beendet den Container (LH-QA-REL-001, LH-QA-REL-003) belegt — bei gehaltener Persistierung (cdc.change exklusiv gesperrt, wartende Persistierung in pg_stat_activity) lief der Feed-Container über eine Nullprobe von ${WAL_WAIT_SECONDS} s weiter; ein Schreiber auf die nicht aktivierte Tabelle $WAL_STOP_FOREIGN erzeugte $wal_stop_bytes B WAL, der Container endete $wal_stop_seconds s danach (höchstens ${WAL_STOP_WAIT_SECONDS} s gewartet) mit Ausgang 1, die Abbruch-Zeile nannte einen Rückstand von $wal_stop_logged B (Fehlerschwelle $WAL_ERROR_BYTES B), cdc.process_heartbeat trug die Klasse $wal_stop_class; nach dem Neustart ohne Konfigurationsdatei war die wartende Transaktion über cdc.changes lesbar"

# --- Transformations-Rundläufe (LH-FA-CFG-007) --------------------------------
# Zwei Phasen am laufenden Feed-Container, ausschließlich über externe Wege:
# SQL-Funktionen und Views, HTTP, die drei Stream-Clients und `docker restart`.
# Jede Phase legt ihre eigene Tabelle an und nimmt ihre Regeln am Ende zurück.
# Die Bilder aller Operationen und die Konfliktfreiheit tragen die beiden
# Testfunktionen in test/integration/transformation_e2e_test.go, die der
# erste `go test`-Aufruf dieser Datei fährt.

# tf_set <Tabelle> <Regelname> <Regelform> <Phase>: setzt eine Regel und wartet
# auf den Vermerk.
tf_set() {
  local request_id
  request_id=$(bf_sql "SELECT cdc.set_transformation('src-e2e', 'public', '$1', '$2', '$3')")
  [ -n "$request_id" ] || bf_fail "$4 — cdc.set_transformation($1, $2) lieferte keine Antrags-ID"
  bf_await_applied "$request_id" "$4"
}

# tf_remove <Tabelle> <Regelname> <Phase>
tf_remove() {
  local request_id
  request_id=$(bf_sql "SELECT cdc.remove_transformation('src-e2e', 'public', '$1', '$2')")
  [ -n "$request_id" ] || bf_fail "$3 — cdc.remove_transformation($1, $2) lieferte keine Antrags-ID"
  bf_await_applied "$request_id" "$3"
}

# tf_exclude <Tabelle> <Spalte> <Phase>
tf_exclude() {
  local request_id
  request_id=$(bf_sql "SELECT cdc.exclude_column('src-e2e', 'public', '$1', '$2')")
  [ -n "$request_id" ] || bf_fail "$3 — cdc.exclude_column($1, $2) lieferte keine Antrags-ID"
  bf_await_applied "$request_id" "$3"
}

# tf_restart_feed <Phase>: startet den Feed-Container real neu, wartet auf
# healthy und beendet die Phase, wenn seine Startzeit gleich blieb. Der Beleg
# eines Neustarts ist die Startzeit; tf_started_before und tf_started_after
# tragen sie nach dem Aufruf.
tf_restart_feed() {
  tf_started_before=$(docker inspect --format '{{.State.StartedAt}}' "$FEED_CONTAINER")
  docker restart "$FEED_CONTAINER" >/dev/null
  bf_await_healthy "$1"
  tf_started_after=$(docker inspect --format '{{.State.StartedAt}}' "$FEED_CONTAINER")
  [ "$tf_started_before" != "$tf_started_after" ] || bf_fail "$1 — die Startzeit des Feed-Containers blieb nach docker restart gleich ($tf_started_after)"
}

# tf_row <Tabelle> <id> <name> <Status> <Phase>: fügt eine Zeile ein und wartet,
# bis ihre Change über cdc.changes lesbar ist.
tf_row() {
  bf_sql "INSERT INTO public.$1 (id, name, status) VALUES ($2, '$3', '$4')" >/dev/null
  bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$1' AND new_data->>'id' = '$2'" 1 30 "$5 — Change der Zeile $2"
}

# tf_form <Tabelle> <id> <erwartetes Bild ohne id> <Beschreibung>: das
# persistierte Bild der Zeile trägt ohne den Schlüssel id genau diese Form.
tf_form() {
  bf_expect "$(bf_sql "SELECT (new_data - 'id') = '$3'::jsonb FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$1' AND new_data->>'id' = '$2'")" t "$4"
}

# tf_no_value <Tabelle> <id> <Wert> <Beschreibung>: der Wert steht weder im
# Neu- noch im Alt-Bild der Change.
tf_no_value() {
  bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$1' AND new_data->>'id' = '$2' AND (new_data::text LIKE '%$3%' OR coalesce(old_data::text, '') LIKE '%$3%')")" 0 "$4"
}

# tf_image_of <Zeile>: das Row Image einer RECEIVED- oder READ-Zeile der
# Wegwerf-Clients.
tf_image_of() {
  printf '%s\n' "$1" | sed -E 's/^.* new_image=//; s/ origin=(wal|backfill)$//'
}

# tf_change_id_of <Zeile>: die change_id einer RECEIVED- oder READ-Zeile.
tf_change_id_of() {
  grep -oE 'change_id=[^ ]+' <<<"$1" | cut -d= -f2
}

# tf_expect_form <Bild> <change_id> <erwartete Form ohne id> <Beschreibung>:
# das über einen Zustellweg gelesene Bild trägt die Form und ist gleich dem
# persistierten Bild derselben Change.
tf_expect_form() {
  local image=$1 change_id=$2 want=$3 what=$4
  [ -n "$image" ] || bf_fail "$what — die Zeile trägt kein Row Image"
  bf_expect "$(bf_sql "SELECT ('$image'::jsonb - 'id') = '$want'::jsonb")" t "$what — Form des Bildes ($image)"
  bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$change_id' AND new_data = '$image'::jsonb")" 1 "$what — Bild gleich dem persistierten Bild der Change $change_id"
}

# tf_client_start <Container> <go run-Argumente>: startet einen Wegwerf-Client
# im Toolchain-Container.
tf_client_start() {
  local name=$1
  shift
  docker rm -fv "$name" >/dev/null 2>&1 || true
  docker run -d --name "$name" --network "$NETWORK" \
    -v "$(pwd)":/src:ro \
    -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
    -w /src \
    -e GOCACHE=/tmp/gocache \
    "$TOOLCHAIN_IMAGE" go run "$@" >/dev/null
}

# tf_client_has <Container> <Marker>: die Ausgabe des Clients trägt den Marker.
tf_client_has() {
  local logs
  logs=$(docker logs "$1" 2>&1 || true)
  [[ "$logs" == *"$2"* ]]
}

# tf_client_await <Container> <Marker> <Sekunden> <Phase>: wartet auf eine
# Zeile des Clients; endet der Client vorher ohne sie, endet die Phase.
tf_client_await() {
  local name=$1 marker=$2 seconds=$3 phase=$4 i
  for ((i = 0; i < seconds * 2; i++)); do
    if tf_client_has "$name" "$marker"; then
      return 0
    fi
    if [ "$(docker inspect --format '{{.State.Running}}' "$name" 2>/dev/null || echo false)" != "true" ]; then
      break
    fi
    sleep 0.5
  done
  if tf_client_has "$name" "$marker"; then
    return 0
  fi
  bf_fail "$phase — Client $name meldete '$marker' nicht: $(docker logs "$name" 2>&1 || true)"
}

# tf_client_line <Container>: die erste RECEIVED-Zeile des Clients.
tf_client_line() {
  local logs
  logs=$(docker logs "$1" 2>&1 || true)
  grep -m1 '^RECEIVED ' <<<"$logs" || true
}

abdeckung_declare "Transformationen-Happy-Path (fünf Zustellwege)" "LH-FA-CFG-007,LH-FA-SST-002,LH-FA-SST-006,LH-FA-SST-008" "eine per cdc.set_transformation beantragte rename_column- und map_value-Regel prägt eine danach eingefügte Zeile auf allen fünf Wegen in derselben Form: über cdc.changes, GET /changes, den gRPC-Stream, den SSE-Stream und den NATS-Vollinhalts-Stream trägt das Neu-Bild den Zielnamen und den abgebildeten Wert und nicht den Quellschlüssel; die zuvor erfasste Change bleibt in Rohform lesbar, ein nicht abgebildeter Wert bleibt unverändert und ein fehlender Wert abwesend (beides über cdc.changes und GET /changes), und nach dem Entfernen der Regeln trägt die nächste Change wieder die Rohform; die Wegwerf-Clients der Stream-Wege drucken nur das Neu-Bild, die Bilder von UPDATE und DELETE trägt allein TestE2ETransformationRulesShapeBothImages über cdc.changes" "Transformationen-Happy-Path (LH-FA-CFG-007) belegt"

TF_PHASE="Transformationen-Happy-Path"
TF_TABLE=feed_e2e_transform
# Kopplung: die Phase „Transformationen-Neustart und Ausschluss“ setzt
# dieselben zwei Regelformen (TF_RULE_RENAME, TF_RULE_MAP) auf ihrer eigenen
# Tabelle; sie liest beide Variablen aus diesem Block.
TF_RULE_RENAME='{"kind":"rename_column","column":"name","to":"customer_name"}'
TF_RULE_MAP='{"kind":"map_value","column":"status","values":{"o":"open","c":"closed"}}'
TF_FORM_STREAM='{"customer_name":"TfNeu","status":"open"}'
TF_GRPC_CONTAINER=cdc-e2e-tf-grpc
TF_SSE_CONTAINER=cdc-e2e-tf-sse
TF_NATS_CONTAINER=cdc-e2e-tf-nats

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$TF_TABLE (id int PRIMARY KEY, name text, status text);
SQL
bf_enable "$TF_TABLE" "$TF_PHASE"

# Die Zeile vor den Regeln behält ihre Rohform.
tf_row "$TF_TABLE" 1 TfAlt o "$TF_PHASE"
tf_set "$TF_TABLE" kundenname "$TF_RULE_RENAME" "$TF_PHASE"
tf_set "$TF_TABLE" status_lesbar "$TF_RULE_MAP" "$TF_PHASE"

# Die drei Stream-Clients laufen gleichzeitig und melden ihre Bereitschaft,
# bevor die erste Zeile mit Regeln entsteht.
tf_client_start "$TF_GRPC_CONTAINER" ./tools/harness/grpcclient "pg-change-feed:9090" "$HTTP_TOKEN_READER"
tf_client_start "$TF_SSE_CONTAINER" ./tools/harness/sseclient "$HTTP_BASE_URL" "$HTTP_TOKEN_READER"
tf_client_start "$TF_NATS_CONTAINER" ./tools/harness/natsstreamsub "nats://nats:4222" "cdc.stream.src-e2e.public.$TF_TABLE" "$NATS_STREAM_TOKEN"
tf_client_await "$TF_GRPC_CONTAINER" READY 90 "$TF_PHASE"
tf_client_await "$TF_SSE_CONTAINER" READY 90 "$TF_PHASE"
tf_client_await "$TF_NATS_CONTAINER" READY 90 "$TF_PHASE"

# Die Zustellung über die Streams trägt kein Replay: zwischen „Client ist
# bereit“ und „Empfänger ist am Broadcaster registriert“ liegt ein kurzes
# Fenster. Der Lauf fügt deshalb eine begrenzte Folge gleichförmiger Zeilen
# ein, bis jeder der drei Clients eine davon empfangen hat.
tf_stream_rows=0
tf_all_received=0
for tf_attempt in $(seq 1 6); do
  tf_stream_rows=$tf_attempt
  bf_sql "INSERT INTO public.$TF_TABLE (id, name, status) VALUES ($((10 + tf_attempt)), 'TfNeu', 'o')" >/dev/null
  for _ in $(seq 1 20); do
    if tf_client_has "$TF_GRPC_CONTAINER" "RECEIVED" && tf_client_has "$TF_SSE_CONTAINER" "RECEIVED" && tf_client_has "$TF_NATS_CONTAINER" "RECEIVED"; then
      tf_all_received=1
      break
    fi
    sleep 0.5
  done
  if [ "$tf_all_received" -eq 1 ]; then
    break
  fi
done
[ "$tf_all_received" -eq 1 ] || bf_fail "$TF_PHASE — nicht jeder Stream-Client empfing eine der $tf_stream_rows Zeilen (gRPC: $(tf_client_line "$TF_GRPC_CONTAINER") SSE: $(tf_client_line "$TF_SSE_CONTAINER") NATS: $(tf_client_line "$TF_NATS_CONTAINER"))"

tf_grpc_line=$(tf_client_line "$TF_GRPC_CONTAINER")
tf_sse_line=$(tf_client_line "$TF_SSE_CONTAINER")
tf_nats_line=$(tf_client_line "$TF_NATS_CONTAINER")
docker rm -fv "$TF_GRPC_CONTAINER" "$TF_SSE_CONTAINER" "$TF_NATS_CONTAINER" >/dev/null 2>&1 || true
for tf_way in gRPC-Stream:"$tf_grpc_line" SSE-Stream:"$tf_sse_line" NATS-Vollinhalts-Stream:"$tf_nats_line"; do
  tf_way_name=${tf_way%%:*}
  tf_way_line=${tf_way#*:}
  [[ "$tf_way_line" == *" table=$TF_TABLE "* ]] || bf_fail "$TF_PHASE — $tf_way_name: die RECEIVED-Zeile nennt nicht die Tabelle $TF_TABLE: $tf_way_line"
  [[ "$tf_way_line" == *" operation=INSERT "* ]] || bf_fail "$TF_PHASE — $tf_way_name: die RECEIVED-Zeile nennt nicht die Operation INSERT: $tf_way_line"
  tf_expect_form "$(tf_image_of "$tf_way_line")" "$(tf_change_id_of "$tf_way_line")" "$TF_FORM_STREAM" "$TF_PHASE — $tf_way_name"
done

# Ein nicht abgebildeter Wert und ein fehlender Wert: die Regeln ändern nur,
# was sie nennen.
tf_row "$TF_TABLE" 20 TfX x "$TF_PHASE"
bf_sql "INSERT INTO public.$TF_TABLE (id, name, status) VALUES (21, NULL, 'c')" >/dev/null
tf_total=$((1 + tf_stream_rows + 2))
bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$TF_TABLE'" "$tf_total" 30 "$TF_PHASE — Changes der Tabelle"

# Weg 1, cdc.changes: Rohform vor den Regeln, die Form nach ihnen.
tf_form "$TF_TABLE" 1 '{"name":"TfAlt","status":"o"}' "$TF_PHASE — cdc.changes, Change vor den Regeln"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$TF_TABLE' AND (new_data - 'id') = '$TF_FORM_STREAM'::jsonb")" "$tf_stream_rows" "$TF_PHASE — cdc.changes, Changes mit der Form der Stream-Zeilen"
tf_form "$TF_TABLE" 20 '{"customer_name":"TfX","status":"x"}' "$TF_PHASE — cdc.changes, nicht abgebildeter Wert"
tf_form "$TF_TABLE" 21 '{"status":"closed"}' "$TF_PHASE — cdc.changes, fehlender Wert bleibt abwesend"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$TF_TABLE' AND jsonb_exists(new_data, 'name')")" 1 "$TF_PHASE — cdc.changes, Changes mit dem Quellschlüssel name (nur die vor den Regeln)"

# Weg 2, GET /changes: jede gelesene Change trägt dasselbe Bild wie
# cdc.changes, und die Menge der Kennungen ist dieselbe.
tf_range=$(bf_sql "SELECT min(commit_position) || ' ' || (max(commit_position) + 1) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$TF_TABLE'")
read -r tf_from tf_to <<<"$tf_range"
tf_http=$(bf_http changes "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" src-e2e public "$TF_TABLE" "$tf_from" "$tf_to")
bf_expect "$(printf '%s\n' "$tf_http" | grep -c '^READ ')" "$tf_total" "$TF_PHASE — READ-Zeilen über GET /changes"
bf_expect "$(bf_read_ids "$tf_http")" "$(bf_sql "SELECT string_agg(change_id, ',' ORDER BY change_id COLLATE \"C\") FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$TF_TABLE'")" "$TF_PHASE — change_id über GET /changes gegen cdc.changes"
while IFS= read -r tf_line; do
  tf_http_image=$(tf_image_of "$tf_line")
  tf_http_id=$(tf_change_id_of "$tf_line")
  bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$tf_http_id' AND new_data = '$tf_http_image'::jsonb")" 1 "$TF_PHASE — GET /changes, Bild gleich cdc.changes ($tf_http_id)"
done < <(printf '%s\n' "$tf_http" | grep '^READ ')

# Rücknahme: nach dem Entfernen beider Regeln trägt die nächste Change wieder
# die Rohform, die früheren bleiben unverändert.
tf_remove "$TF_TABLE" kundenname "$TF_PHASE"
tf_remove "$TF_TABLE" status_lesbar "$TF_PHASE"
tf_row "$TF_TABLE" 30 TfRoh o "$TF_PHASE"
tf_form "$TF_TABLE" 30 '{"name":"TfRoh","status":"o"}' "$TF_PHASE — Rohform nach dem Entfernen der Regeln"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$TF_TABLE' AND (new_data - 'id') = '$TF_FORM_STREAM'::jsonb")" "$tf_stream_rows" "$TF_PHASE — frühere Changes nach dem Entfernen der Regeln"
bf_expect "$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)" true "$TF_PHASE — Feed-Container läuft weiter"

echo "run-integration-tests: Transformationen-Happy-Path (LH-FA-CFG-007) belegt — auf $TF_TABLE prägten rename_column (name zu customer_name) und map_value (status) $tf_stream_rows Stream-Zeile(n) auf allen fünf Wegen in derselben Form und zwei weitere Zeilen (nicht abgebildeter und fehlender Wert) über cdc.changes und GET /changes; gRPC: $tf_grpc_line; SSE: $tf_sse_line; NATS: $tf_nats_line; die Zeile vor den Regeln blieb in Rohform, nach dem Entfernen der Regeln trug die nächste Change wieder die Rohform"

abdeckung_declare "Transformationen-Neustart und Ausschluss" "LH-FA-CFG-007,LH-FA-CFG-005,LH-QA-SEC-004" "nach einem realen Container-Neustart leitet der Prozessstart den Regelstand aus den applied-Zeilen ab und die danach erfasste Change trägt die Form; cdc.remove_transformation stellt die Rohform für künftige Changes wieder her; nach cdc.exclude_column auf der Spalte mit Regel trägt die Change weder Quellnamen noch Zielnamen noch Wert im Bild, vor und nach dem Neustart, und die nicht ausgeschlossene Spalte bleibt darin" "Transformationen-Neustart und Ausschluss (LH-FA-CFG-007) belegt"

TF_PHASE="Transformationen-Neustart und Ausschluss"
TF_RESTART_TABLE=feed_e2e_transform_restart

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$TF_RESTART_TABLE (id int PRIMARY KEY, name text, status text);
SQL
bf_enable "$TF_RESTART_TABLE" "$TF_PHASE"
tf_set "$TF_RESTART_TABLE" kundenname "$TF_RULE_RENAME" "$TF_PHASE"
tf_set "$TF_RESTART_TABLE" status_lesbar "$TF_RULE_MAP" "$TF_PHASE"

tf_row "$TF_RESTART_TABLE" 1 TfR1 o "$TF_PHASE"
tf_form "$TF_RESTART_TABLE" 1 '{"customer_name":"TfR1","status":"open"}' "$TF_PHASE — Form vor dem Neustart"

# Der Neustart ist ein echter Prozess-Neustart desselben Containers: der
# Prozessstart leitet den Regelstand aus den applied-Zeilen ab, und die Form
# der nächsten Change ist der Beleg dieser Ableitung. Der Beleg des Neustarts
# selbst ist die Startzeit des Containers (tf_restart_feed).
tf_restart_feed "$TF_PHASE"
tf_first_restart_before=$tf_started_before
tf_first_restart_after=$tf_started_after
tf_row "$TF_RESTART_TABLE" 2 TfR2 o "$TF_PHASE"
tf_form "$TF_RESTART_TABLE" 2 '{"customer_name":"TfR2","status":"open"}' "$TF_PHASE — Form nach dem Neustart"

# Rücknahme: die nächste Change trägt die Rohform, die frühere bleibt.
tf_remove "$TF_RESTART_TABLE" kundenname "$TF_PHASE"
tf_remove "$TF_RESTART_TABLE" status_lesbar "$TF_PHASE"
tf_row "$TF_RESTART_TABLE" 3 TfR3 o "$TF_PHASE"
tf_form "$TF_RESTART_TABLE" 3 '{"name":"TfR3","status":"o"}' "$TF_PHASE — Rohform nach dem Entfernen der Regeln"
tf_form "$TF_RESTART_TABLE" 2 '{"customer_name":"TfR2","status":"open"}' "$TF_PHASE — frühere Change nach dem Entfernen der Regeln"

# Ausschluss der Spalte mit Regel: dasselbe Regelpaar wird neu gesetzt, dann
# fällt die Spalte name aus dem Bild — unter dem Quellnamen wie unter dem
# Zielnamen —, die Spalte status trägt ihre Regel weiter.
tf_set "$TF_RESTART_TABLE" kundenname "$TF_RULE_RENAME" "$TF_PHASE"
tf_set "$TF_RESTART_TABLE" status_lesbar "$TF_RULE_MAP" "$TF_PHASE"
tf_row "$TF_RESTART_TABLE" 4 TfR4 o "$TF_PHASE"
tf_form "$TF_RESTART_TABLE" 4 '{"customer_name":"TfR4","status":"open"}' "$TF_PHASE — Form nach dem erneuten Setzen"
tf_exclude "$TF_RESTART_TABLE" name "$TF_PHASE"
tf_row "$TF_RESTART_TABLE" 5 TfSecretVorNeustart o "$TF_PHASE"
tf_form "$TF_RESTART_TABLE" 5 '{"status":"open"}' "$TF_PHASE — Bild mit Ausschluss und Regeln vor dem Neustart"
tf_no_value "$TF_RESTART_TABLE" 5 TfSecretVorNeustart "$TF_PHASE — Wert der ausgeschlossenen Spalte vor dem Neustart"

tf_restart_feed "$TF_PHASE"
tf_row "$TF_RESTART_TABLE" 6 TfSecretNachNeustart o "$TF_PHASE"
tf_form "$TF_RESTART_TABLE" 6 '{"status":"open"}' "$TF_PHASE — Bild mit Ausschluss und Regeln nach dem Neustart"
tf_no_value "$TF_RESTART_TABLE" 6 TfSecretNachNeustart "$TF_PHASE — Wert der ausgeschlossenen Spalte nach dem Neustart"

# Aufräumen: ohne Regeln bleibt der Ausschluss, die Spalte status trägt die
# Rohform.
tf_remove "$TF_RESTART_TABLE" kundenname "$TF_PHASE"
tf_remove "$TF_RESTART_TABLE" status_lesbar "$TF_PHASE"
tf_row "$TF_RESTART_TABLE" 7 TfSecretOhneRegeln o "$TF_PHASE"
tf_form "$TF_RESTART_TABLE" 7 '{"status":"o"}' "$TF_PHASE — Bild mit Ausschluss ohne Regeln"
tf_no_value "$TF_RESTART_TABLE" 7 TfSecretOhneRegeln "$TF_PHASE — Wert der ausgeschlossenen Spalte ohne Regeln"

echo "run-integration-tests: Transformationen-Neustart und Ausschluss (LH-FA-CFG-007) belegt — auf $TF_RESTART_TABLE trug die Change nach einem realen docker restart (Startzeit $tf_first_restart_before, danach $tf_first_restart_after) die Form beider Regeln; nach dem Entfernen trug die nächste Change die Rohform; nach exclude_column auf name trug das Bild vor und nach dem zweiten docker restart (Startzeit $tf_started_before, danach $tf_started_after) weder name noch customer_name noch den Wert, die Spalte status blieb (mit Regel offen, ohne Regel roh)"

# --- Routing-Rundläufe (LH-FA-CFG-008) ---------------------------------------
# Drei Phasen am laufenden Feed-Container, ausschließlich über externe Wege:
# SQL-Funktionen und Views, HTTP, gRPC, die Wegwerf-Clients und `docker
# restart`; die Nichtanwendbarkeit und ihre Abhilfe laufen als letzter Rundlauf
# dieses Runners. Jede Phase legt eigene Tabellen an. Das Ziel einer Change
# trägt die Spalte route_target der Sicht cdc.changes; die Nachrichten der
# Stream-Wege tragen es nicht, deshalb hält jede Phase das über einen Weg
# gelesene Ziel gegen die persistierte Zeile derselben change_id. Auflösung bei
# zwei treffenden Regeln, R1 bis R6, Replay, Backfill-Label und die
# DELETE-Messung tragen die vier Testfunktionen in
# test/integration/routing_e2e_test.go, die der erste `go test`-Aufruf dieser
# Datei fährt.

# rt_set <Tabelle> <Regelname> <Regelform> <Phase>: setzt eine Routing-Regel
# und wartet auf den Vermerk.
rt_set() {
  local request_id
  request_id=$(bf_sql "SELECT cdc.set_route('src-e2e', 'public', '$1', '$2', '$3')")
  [ -n "$request_id" ] || bf_fail "$4 — cdc.set_route($1, $2) lieferte keine Antrags-ID"
  bf_await_applied "$request_id" "$4"
}

# rt_remove <Tabelle> <Regelname> <Phase>
rt_remove() {
  local request_id
  request_id=$(bf_sql "SELECT cdc.remove_route('src-e2e', 'public', '$1', '$2')")
  [ -n "$request_id" ] || bf_fail "$3 — cdc.remove_route($1, $2) lieferte keine Antrags-ID"
  bf_await_applied "$request_id" "$3"
}

# rt_target_of_change <change_id>: das persistierte Ziel einer Change, NULL als
# Text NULL.
rt_target_of_change() {
  bf_sql "SELECT coalesce(route_target, 'NULL') FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$1'"
}

# rt_target_of_row <Tabelle> <id>: das persistierte Ziel der Change einer Zeile.
rt_target_of_row() {
  bf_sql "SELECT coalesce(route_target, 'NULL') FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$1' AND new_data->>'id' = '$2'"
}

# rt_row <Tabelle> <id> <name> <region als SQL-Literal> <Phase>: fügt eine Zeile
# ein und wartet, bis ihre Change über cdc.changes lesbar ist.
rt_row() {
  bf_sql "INSERT INTO public.$1 (id, name, region) VALUES ($2, '$3', $4)" >/dev/null
  bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$1' AND new_data->>'id' = '$2'" 1 30 "$5 — Change der Zeile $2"
}

# rt_received_count <Container>: die Zahl der RECEIVED-Zeilen eines Clients.
rt_received_count() {
  local logs
  logs=$(docker logs "$1" 2>&1 || true)
  printf '%s\n' "$logs" | grep -c '^RECEIVED ' || true
}

# rt_received_lines <Container>: alle RECEIVED-Zeilen eines Clients.
rt_received_lines() {
  local logs
  logs=$(docker logs "$1" 2>&1 || true)
  printf '%s\n' "$logs" | grep '^RECEIVED ' || true
}

# rt_expect_target_line <Zeile> <Tabelle> <erwartetes Ziel> <Beschreibung>: die
# RECEIVED-Zeile eines auf ein Ziel gewählten Stream-Clients nennt die Tabelle,
# ein INSERT, und die persistierte Change derselben change_id trägt das Ziel.
rt_expect_target_line() {
  local line=$1 table=$2 target=$3 what=$4 change_id
  [[ "$line" == *" table=$table "* ]] || bf_fail "$what — die RECEIVED-Zeile nennt nicht die Tabelle $table: $line"
  [[ "$line" == *" operation=INSERT "* ]] || bf_fail "$what — die RECEIVED-Zeile nennt nicht die Operation INSERT: $line"
  change_id=$(tf_change_id_of "$line")
  [ -n "$change_id" ] || bf_fail "$what — die RECEIVED-Zeile trägt keine change_id: $line"
  bf_expect "$(rt_target_of_change "$change_id")" "$target" "$what — persistiertes Ziel der Change $change_id"
  bf_expect "$(bf_sql "SELECT new_data->>'region' FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$change_id'")" "$target" "$what — Region der Change $change_id"
}

# rt_grpc_read <Ziel oder leer> <Tabelle> <from> <to>: der RPC ReadChanges im
# Toolchain-Container über tools/harness/grpcadminclient, je Aufruf mit einem
# eigenen Consumer; ein Fehlschlag des Clients endet die Phase.
rt_grpc_read() {
  local output status consumer flags=()
  if [ -n "$1" ]; then
    flags=(-target "$1")
  fi
  consumer="rt-e2e-consumer-$(date +%s%N)"
  set +e
  output=$(docker run --rm --network "$NETWORK" \
    -v "$(pwd)":/src:ro \
    -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
    -w /src \
    -e GOCACHE=/tmp/gocache \
    "$TOOLCHAIN_IMAGE" go run ./tools/harness/grpcadminclient ${flags[@]+"${flags[@]}"} \
    "$GRPC_ADDR" "$HTTP_TOKEN_READER" "$HTTP_TOKEN_ADMIN" "$consumer" src-e2e pub_pgc_e2e public "$2" "$3" "$4" 500 2>&1)
  status=$?
  set -e
  if [ "$status" -ne 0 ]; then
    bf_fail "grpcadminclient endete mit Ausgang $status: $output"
  fi
  printf '%s\n' "$output"
}

abdeckung_declare "Routing-Happy-Path (fünf Zustellwege)" "LH-FA-CFG-008,LH-FA-SST-006,LH-FA-SST-008" "zwei per cdc.set_route beantragte Inhaltsregeln (Ziele eu und us) lenken eine danach eingefügte Zeile auf ihr Ziel, und jedes Ziel ist auf allen fünf Wegen auswählbar: über cdc.changes (WHERE route_target), GET /changes?target=, den ReadChanges-RPC mit target, den gRPC-Stream, den SSE-Stream und das NATS-Subjekt cdc.route.<source_id>.<ziel>; jedes über einen Weg gelesene Ziel wird gegen die persistierte Zeile derselben change_id gehalten, ein auf ein Ziel gewählter Stream-Client empfängt von einer festen Menge gemischter Changes (Region NULL, asia, us, eu) genau die Changes seines Ziels und im Ruhefenster keine Change eines anderen Ziels oder ohne Ziel, die Pull-Wege liefern je Ziel genau die Kennungen der SQL-Auswahl, ein ungefilterter Leser sieht eine Change ohne Treffer (route_target NULL) und die gerouteten Changes" "Routing-Happy-Path (LH-FA-CFG-008) belegt"

RT_PHASE="Routing-Happy-Path"
RT_TABLE=feed_e2e_route
RT_RULE_EU='{"target":"eu","order":10,"when":{"column":"region","equals":"eu"}}'
RT_RULE_US='{"target":"us","order":20,"when":{"column":"region","equals":"us"}}'
RT_NATS_URL="nats://nats:4222"

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$RT_TABLE (id int PRIMARY KEY, name text, region text);
SQL
bf_enable "$RT_TABLE" "$RT_PHASE"
rt_set "$RT_TABLE" eu_orders "$RT_RULE_EU" "$RT_PHASE"
rt_set "$RT_TABLE" us_orders "$RT_RULE_US" "$RT_PHASE"

# Neun Stream-Clients laufen gleichzeitig: je Weg einer auf Ziel eu, einer auf
# Ziel us und einer ohne Ziel, der drei Changes der Tabelle empfängt.
RT_CLIENTS="cdc-e2e-rt-grpc-eu cdc-e2e-rt-grpc-us cdc-e2e-rt-grpc-all cdc-e2e-rt-sse-eu cdc-e2e-rt-sse-us cdc-e2e-rt-sse-all cdc-e2e-rt-nats-eu cdc-e2e-rt-nats-us cdc-e2e-rt-nats-all"
RT_WINDOW=15s
tf_client_start cdc-e2e-rt-grpc-eu ./tools/harness/grpcclient -target eu -window "$RT_WINDOW" "$GRPC_ADDR" "$HTTP_TOKEN_READER"
tf_client_start cdc-e2e-rt-grpc-us ./tools/harness/grpcclient -target us -window "$RT_WINDOW" "$GRPC_ADDR" "$HTTP_TOKEN_READER"
tf_client_start cdc-e2e-rt-grpc-all ./tools/harness/grpcclient -count 3 -window "$RT_WINDOW" "$GRPC_ADDR" "$HTTP_TOKEN_READER" public "$RT_TABLE"
tf_client_start cdc-e2e-rt-sse-eu ./tools/harness/sseclient -target eu -window "$RT_WINDOW" "$HTTP_BASE_URL" "$HTTP_TOKEN_READER"
tf_client_start cdc-e2e-rt-sse-us ./tools/harness/sseclient -target us -window "$RT_WINDOW" "$HTTP_BASE_URL" "$HTTP_TOKEN_READER"
tf_client_start cdc-e2e-rt-sse-all ./tools/harness/sseclient -count 3 -window "$RT_WINDOW" "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" public "$RT_TABLE"
tf_client_start cdc-e2e-rt-nats-eu ./tools/harness/natsstreamsub -source src-e2e -target eu -window "$RT_WINDOW" "$RT_NATS_URL" - "$NATS_STREAM_TOKEN"
tf_client_start cdc-e2e-rt-nats-us ./tools/harness/natsstreamsub -source src-e2e -target us -window "$RT_WINDOW" "$RT_NATS_URL" - "$NATS_STREAM_TOKEN"
tf_client_start cdc-e2e-rt-nats-all ./tools/harness/natsstreamsub -count 3 -window "$RT_WINDOW" "$RT_NATS_URL" "cdc.stream.src-e2e.public.$RT_TABLE" "$NATS_STREAM_TOKEN"
for rt_client in $RT_CLIENTS; do
  tf_client_await "$rt_client" READY 120 "$RT_PHASE"
done

# Die Zustellung über die Streams trägt kein Replay: zwischen „Client ist
# bereit“ und „Empfänger ist am Broadcaster registriert“ liegt ein kurzes
# Fenster. Der Lauf fügt deshalb Dreiergruppen ein (Region asia ohne Treffer,
# us, eu), bis die sechs Clients mit Ziel eine Change und die drei Clients ohne
# Ziel drei Changes empfangen haben; drei aufeinanderfolgende Zeilen dieser
# Folge tragen je eine Region.
rt_attempts=0
rt_all_received=0
for rt_attempt in $(seq 1 6); do
  rt_attempts=$rt_attempt
  rt_row "$RT_TABLE" $((rt_attempt * 100 + 1)) RtAsien "'asia'" "$RT_PHASE"
  rt_row "$RT_TABLE" $((rt_attempt * 100 + 2)) RtUsa "'us'" "$RT_PHASE"
  rt_row "$RT_TABLE" $((rt_attempt * 100 + 3)) RtEuropa "'eu'" "$RT_PHASE"
  for _ in $(seq 1 20); do
    rt_ready=1
    for rt_client in $RT_CLIENTS; do
      case "$rt_client" in
        *-all) rt_need=3 ;;
        *) rt_need=1 ;;
      esac
      if [ "$(rt_received_count "$rt_client")" -lt "$rt_need" ]; then
        rt_ready=0
      fi
    done
    if [ "$rt_ready" -eq 1 ]; then
      rt_all_received=1
      break
    fi
    sleep 0.5
  done
  if [ "$rt_all_received" -eq 1 ]; then
    break
  fi
done
if [ "$rt_all_received" -ne 1 ]; then
  rt_report=""
  for rt_client in $RT_CLIENTS; do
    rt_report="$rt_report $rt_client=$(rt_received_count "$rt_client")"
  done
  bf_fail "$RT_PHASE — nicht jeder Stream-Client empfing seine Changes nach $rt_attempts Dreiergruppen (RECEIVED-Zeilen:$rt_report)"
fi

# Jeder Client ist jetzt am Broadcaster registriert. Eine feste Menge
# gemischter Changes (je Gruppe: Region NULL, asia, us, eu; zwei Gruppen) folgt;
# die Clients zählen im Ruhefenster weiter. Ein Client mit Ziel empfängt von
# dieser Menge genau die Changes seines Ziels und keine Change eines anderen
# Ziels oder ohne Ziel, ein Client ohne Ziel alle acht.
rt_final_rows=8
for rt_group in 1 2; do
  rt_row "$RT_TABLE" $((1000 + rt_group * 10 + 1)) RtNull "NULL" "$RT_PHASE"
  rt_row "$RT_TABLE" $((1000 + rt_group * 10 + 2)) RtAsien "'asia'" "$RT_PHASE"
  rt_row "$RT_TABLE" $((1000 + rt_group * 10 + 3)) RtUsa "'us'" "$RT_PHASE"
  rt_row "$RT_TABLE" $((1000 + rt_group * 10 + 4)) RtEuropa "'eu'" "$RT_PHASE"
done
for rt_client in $RT_CLIENTS; do
  tf_client_await "$rt_client" WINDOW-END 90 "$RT_PHASE"
done

# rt_audit_client <Container> <Ziel oder leer> <Beschreibung>: prüft jede
# RECEIVED-Zeile des Clients gegen die persistierte Change derselben change_id.
# Ein Client mit Ziel (<Ziel> gesetzt) empfängt keine Change mit anderem Ziel
# und keine ohne Ziel. Setzt rt_audit_final auf die Zahl der Zeilen aus der
# festen Menge (id >= 1000).
rt_audit_client() {
  local name=$1 want=$2 what=$3 line change_id row
  rt_audit_final=0
  while IFS= read -r line; do
    [ -n "$line" ] || continue
    change_id=$(tf_change_id_of "$line")
    row=$(bf_sql "SELECT coalesce(route_target, 'NULL') || ' ' || (new_data->>'id') FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$change_id'")
    [ -n "$row" ] || bf_fail "$what — keine persistierte Change zur Kennung '$change_id' (Zeile: $line)"
    if [ -n "$want" ] && [ "${row% *}" != "$want" ]; then
      bf_fail "$what — der Client mit Ziel $want empfing die Change $change_id mit Ziel '${row% *}' (Zeile: $line)"
    fi
    if [ "${row#* }" -ge 1000 ]; then
      rt_audit_final=$((rt_audit_final + 1))
    fi
  done < <(rt_received_lines "$name")
}

# Ein auf ein Ziel gewählter Client empfängt als erste Change eine mit diesem
# Ziel, obwohl vor ihr Changes anderer Ziele und ohne Ziel erfasst wurden.
rt_lines_report=""
for rt_way in grpc sse nats; do
  for rt_target in eu us; do
    rt_line=$(rt_received_lines "cdc-e2e-rt-$rt_way-$rt_target" | head -n1)
    rt_expect_target_line "$rt_line" "$RT_TABLE" "$rt_target" "$RT_PHASE — $rt_way, Ziel $rt_target"
    rt_lines_report="$rt_lines_report $rt_way/$rt_target: $rt_line;"
    rt_audit_client "cdc-e2e-rt-$rt_way-$rt_target" "$rt_target" "$RT_PHASE — $rt_way, Ziel $rt_target, alle empfangenen Changes"
    bf_expect "$rt_audit_final" "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RT_TABLE' AND route_target = '$rt_target' AND (new_data->>'id')::int >= 1000")" "$RT_PHASE — $rt_way, Ziel $rt_target, Changes der festen Menge"
  done
  rt_audit_client "cdc-e2e-rt-$rt_way-all" "" "$RT_PHASE — $rt_way, Client ohne Ziel"
  bf_expect "$rt_audit_final" "$rt_final_rows" "$RT_PHASE — $rt_way, Client ohne Ziel, Changes der festen Menge"
  # Ein Client ohne Ziel empfängt drei aufeinanderfolgende Changes: eine ohne
  # Ziel, eine mit us, eine mit eu.
  rt_all_targets=""
  while IFS= read -r rt_line; do
    rt_all_targets="$rt_all_targets$(rt_target_of_change "$(tf_change_id_of "$rt_line")"),"
  done < <(rt_received_lines "cdc-e2e-rt-$rt_way-all" | head -n3)
  rt_all_sorted=$(printf '%s' "$rt_all_targets" | tr ',' '\n' | grep -v '^$' | LC_ALL=C sort | paste -sd,)
  bf_expect "$rt_all_sorted" "NULL,eu,us" "$RT_PHASE — $rt_way, Client ohne Ziel: Ziele seiner drei Changes"
done
docker rm -fv $RT_CLIENTS >/dev/null 2>&1 || true

# cdc.changes, GET /changes und ReadChanges: dieselben Kennungen wie die SQL-
# Auswahl; ein ungefilterter Lesezugriff liefert jede Change.
rt_range=$(bf_sql "SELECT min(commit_position) || ' ' || (max(commit_position) + 1) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RT_TABLE'")
read -r rt_from rt_to <<<"$rt_range"
rt_ids_sql() {
  bf_sql "SELECT string_agg(change_id, ',' ORDER BY change_id COLLATE \"C\") FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RT_TABLE' $1"
}
rt_total=$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RT_TABLE'")
bf_expect "$rt_total" $((rt_attempts * 3 + rt_final_rows)) "$RT_PHASE — Changes der Tabelle"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RT_TABLE' AND route_target IS NULL")" $((rt_attempts + 4)) "$RT_PHASE — Changes ohne Ziel (Region asia oder NULL)"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RT_TABLE' AND route_target = 'eu' AND new_data->>'region' <> 'eu'")" 0 "$RT_PHASE — Ziel eu auf einer Change anderer Region"
for rt_target in eu us; do
  rt_sql_ids=$(rt_ids_sql "AND route_target = '$rt_target'")
  [ -n "$rt_sql_ids" ] || bf_fail "$RT_PHASE — keine Change mit Ziel $rt_target in der Sicht"
  bf_expect "$(bf_read_ids "$(bf_http changes -target "$rt_target" "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" src-e2e public "$RT_TABLE" "$rt_from" "$rt_to")")" "$rt_sql_ids" "$RT_PHASE — GET /changes?target=$rt_target gegen cdc.changes"
  bf_expect "$(bf_read_ids "$(rt_grpc_read "$rt_target" "$RT_TABLE" "$rt_from" "$rt_to")")" "$rt_sql_ids" "$RT_PHASE — ReadChanges mit target $rt_target gegen cdc.changes"
done
rt_all_ids=$(rt_ids_sql "")
bf_expect "$(bf_read_ids "$(bf_http changes "$HTTP_BASE_URL" "$HTTP_TOKEN_READER" src-e2e public "$RT_TABLE" "$rt_from" "$rt_to")")" "$rt_all_ids" "$RT_PHASE — GET /changes ohne target gegen cdc.changes"
bf_expect "$(bf_read_ids "$(rt_grpc_read "" "$RT_TABLE" "$rt_from" "$rt_to")")" "$rt_all_ids" "$RT_PHASE — ReadChanges ohne target gegen cdc.changes"
bf_expect "$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)" true "$RT_PHASE — Feed-Container läuft weiter"

echo "run-integration-tests: Routing-Happy-Path (LH-FA-CFG-008) belegt — auf $RT_TABLE lenkten die Regeln eu_orders und us_orders $rt_attempts Dreiergruppe(n) (asia, us, eu) auf eu, us und kein Ziel; gRPC-Stream, SSE-Stream und NATS-Subjekt cdc.route.src-e2e.<ziel> lieferten je Ziel eine Change des Ziels als erste (persistiertes Ziel je change_id gehalten); danach folgten $rt_final_rows gemischte Changes (Region NULL, asia, us, eu, zweimal), und jeder Client mit Ziel empfing im Ruhefenster von $RT_WINDOW genau die Changes seines Ziels aus dieser Menge und keine mit anderem Ziel oder ohne Ziel (jede empfangene Zeile gegen die persistierte Change geprüft), jeder Client ohne Ziel alle $rt_final_rows und als erste drei Changes mit den Zielen NULL, eu und us; GET /changes?target= und ReadChanges mit target lieferten je Ziel genau die Kennungen der SQL-Auswahl WHERE route_target (Bereich [$rt_from,$rt_to)), ungefiltert alle $rt_total Changes:$rt_lines_report"

abdeckung_declare "Routing-Neustart und Ausschluss-Sperre" "LH-FA-CFG-008,LH-FA-CFG-005,LH-QA-SEC-004" "nach einem realen Container-Neustart leitet der Prozessstart den Routing-Regelstand aus den applied-Zeilen ab und die danach erfasste Change trägt das Ziel; cdc.remove_route stellt für künftige Changes das Verhalten ohne die Regel wieder her; nach cdc.exclude_column auf einer Spalte ohne Bedingung trägt das Bild weder ihren Schlüssel noch ihren Wert, vor und nach einem weiteren Neustart, und der Ausschluss der Bedingungsspalte endet auch nach dem Neustart failed" "Routing-Neustart und Ausschluss-Sperre (LH-FA-CFG-008) belegt"

RS_PHASE="Routing-Neustart und Ausschluss-Sperre"
RS_TABLE=feed_e2e_route_restart

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$RS_TABLE (id int PRIMARY KEY, name text, region text, secret text);
SQL
bf_enable "$RS_TABLE" "$RS_PHASE"
rt_set "$RS_TABLE" eu_orders "$RT_RULE_EU" "$RS_PHASE"
rt_set "$RS_TABLE" rest '{"target":"sonstige","order":100}' "$RS_PHASE"

# rs_row <id> <region als SQL-Literal> <Geheimwert>: eine Zeile samt Spalte
# secret und das Warten auf ihre Change.
rs_row() {
  bf_sql "INSERT INTO public.$RS_TABLE (id, name, region, secret) VALUES ($1, 'Rs$1', $2, '$3')" >/dev/null
  bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RS_TABLE' AND new_data->>'id' = '$1'" 1 30 "$RS_PHASE — Change der Zeile $1"
}

rs_row 1 "'eu'" RsGeheim1
bf_expect "$(rt_target_of_row "$RS_TABLE" 1)" eu "$RS_PHASE — Ziel vor dem Neustart"

# Der Neustart ist ein echter Prozess-Neustart desselben Containers: der
# Prozessstart leitet den Regelstand aus den applied-Zeilen ab, und das Ziel der
# nächsten Change ist der Beleg dieser Ableitung (Beleg des Neustarts selbst ist
# die Startzeit, tf_restart_feed).
tf_restart_feed "$RS_PHASE"
rs_first_before=$tf_started_before
rs_first_after=$tf_started_after
rs_row 2 "'eu'" RsGeheim2
rs_row 3 "'us'" RsGeheim3
bf_expect "$(rt_target_of_row "$RS_TABLE" 2)" eu "$RS_PHASE — Ziel der Inhaltsregel nach dem Neustart"
bf_expect "$(rt_target_of_row "$RS_TABLE" 3)" sonstige "$RS_PHASE — Ziel der Abschlussregel nach dem Neustart"

# Rücknahme: ohne die Inhaltsregel trifft die Abschlussregel, danach wirkt die
# neu gesetzte Regel wieder.
rt_remove "$RS_TABLE" eu_orders "$RS_PHASE"
rs_row 4 "'eu'" RsGeheim4
bf_expect "$(rt_target_of_row "$RS_TABLE" 4)" sonstige "$RS_PHASE — Ziel nach cdc.remove_route"
bf_expect "$(rt_target_of_row "$RS_TABLE" 2)" eu "$RS_PHASE — frühere Change nach cdc.remove_route unverändert"
rt_set "$RS_TABLE" eu_orders "$RT_RULE_EU" "$RS_PHASE"
rs_row 5 "'eu'" RsGeheim5
bf_expect "$(rt_target_of_row "$RS_TABLE" 5)" eu "$RS_PHASE — Ziel nach dem erneuten Setzen"

# Ausschluss der Spalte secret (keine Bedingung): kein Schlüssel und kein Wert im
# Bild, das Ziel bleibt, vor und nach einem weiteren Neustart.
tf_exclude "$RS_TABLE" secret "$RS_PHASE"
rs_row 6 "'eu'" RsGeheimVorNeustart
tf_form "$RS_TABLE" 6 '{"name":"Rs6","region":"eu"}' "$RS_PHASE — Bild mit Ausschluss vor dem Neustart"
tf_no_value "$RS_TABLE" 6 RsGeheimVorNeustart "$RS_PHASE — Wert der ausgeschlossenen Spalte vor dem Neustart"
bf_expect "$(rt_target_of_row "$RS_TABLE" 6)" eu "$RS_PHASE — Ziel mit Ausschluss vor dem Neustart"
tf_restart_feed "$RS_PHASE"
rs_row 7 "'eu'" RsGeheimNachNeustart
tf_form "$RS_TABLE" 7 '{"name":"Rs7","region":"eu"}' "$RS_PHASE — Bild mit Ausschluss nach dem Neustart"
tf_no_value "$RS_TABLE" 7 RsGeheimNachNeustart "$RS_PHASE — Wert der ausgeschlossenen Spalte nach dem Neustart"
bf_expect "$(rt_target_of_row "$RS_TABLE" 7)" eu "$RS_PHASE — Ziel mit Ausschluss nach dem Neustart"

# R3 in die andere Richtung, nach dem Neustart: der Ausschluss der Spalte mit
# Bedingung endet failed, die Spalte region bleibt im Bild.
rs_exclude_request=$(bf_sql "SELECT cdc.exclude_column('src-e2e', 'public', '$RS_TABLE', 'region')")
[ -n "$rs_exclude_request" ] || bf_fail "$RS_PHASE — cdc.exclude_column(region) lieferte keine Antrags-ID"
bf_await_sql "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$rs_exclude_request'" failed 30 "$RS_PHASE — Ausschluss der Bedingungsspalte"
bf_expect "$(bf_sql "SELECT error_message FROM cdc.administration_request WHERE administration_request_id = '$rs_exclude_request'")" "Spalte trägt eine Routing-Bedingung: public.$RS_TABLE.region" "$RS_PHASE — Fehlertext des Ausschlusses der Bedingungsspalte"
rs_row 8 "'eu'" RsGeheimNachAblehnung
tf_form "$RS_TABLE" 8 '{"name":"Rs8","region":"eu"}' "$RS_PHASE — Bild nach dem abgelehnten Ausschluss"
bf_expect "$(rt_target_of_row "$RS_TABLE" 8)" eu "$RS_PHASE — Ziel nach dem abgelehnten Ausschluss"

echo "run-integration-tests: Routing-Neustart und Ausschluss-Sperre (LH-FA-CFG-008) belegt — auf $RS_TABLE trug die Change nach einem realen docker restart (Startzeit $rs_first_before, danach $rs_first_after) das Ziel der abgeleiteten Regeln (eu, Abschlussregel sonstige); nach cdc.remove_route traf die Abschlussregel, nach dem erneuten Setzen wieder die Inhaltsregel; nach cdc.exclude_column auf secret trug das Bild vor und nach dem zweiten docker restart (Startzeit $tf_started_before, danach $tf_started_after) weder secret noch den Wert und weiter das Ziel eu; der Ausschluss der Bedingungsspalte region endete nach dem Neustart failed (Spalte trägt eine Routing-Bedingung)"

abdeckung_declare "Routing-Aktivierung über die API" "LH-FA-CFG-008,LH-FA-CFG-001,LH-FA-SST-006" "eine Tabelle mit applied-Routing-Regel, die weder per cdc.enable_table noch per Prozessstart, sondern über EnableTable des gRPC- bzw. HTTP-Zugriffswegs aktiviert wird, trägt den abgeleiteten Regelstand: die danach erfasste Change trägt das Ziel der Regel; eine auf demselben Weg aktivierte Tabelle ohne Regel trägt kein Ziel" "Routing-Aktivierung über die API (LH-FA-CFG-008) belegt"

RA_PHASE="Routing-Aktivierung über die API"
RA_RULE='{"target":"api","order":10,"when":{"column":"region","equals":"eu"}}'
ra_report=""
for ra_way in grpc http; do
  for ra_kind in rule plain; do
    ra_table="feed_e2e_route_api_${ra_way}_${ra_kind}"
    ra_id="tbl-e2e-route-api-${ra_way}-${ra_kind}"
    docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$ra_table (id int PRIMARY KEY, name text, region text);
SQL
    if [ "$ra_kind" = rule ]; then
      rt_set "$ra_table" api_orders "$RA_RULE" "$RA_PHASE"
    fi
    if [ "$ra_way" = grpc ]; then
      set +e
      ra_output=$(docker run --rm --network "$NETWORK" \
        -v "$(pwd)":/src:ro \
        -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
        -w /src \
        -e GOCACHE=/tmp/gocache \
        "$TOOLCHAIN_IMAGE" go run ./tools/harness/grpcenabletableclient "$GRPC_ADDR" enable-table \
        "$HTTP_TOKEN_ADMIN" src-e2e public "$ra_table" "$ra_id" "$ra_id-v1" pub_pgc_e2e 2>&1)
      ra_status=$?
      set -e
    else
      ra_output=$(bf_http enable-table "$HTTP_BASE_URL" "$HTTP_TOKEN_ADMIN" src-e2e public "$ra_table" "$ra_id" "$ra_id-v1" pub_pgc_e2e) && ra_status=0 || ra_status=$?
    fi
    [ "$ra_status" -eq 0 ] || bf_fail "$RA_PHASE — EnableTable($ra_table) über $ra_way endete mit Ausgang $ra_status: $ra_output"
    [[ "$ra_output" == *ENABLED* ]] || bf_fail "$RA_PHASE — keine ENABLED-Zeile über $ra_way: $ra_output"
    rt_row "$ra_table" 1 RaEins "'eu'" "$RA_PHASE"
    if [ "$ra_kind" = rule ]; then
      bf_expect "$(rt_target_of_row "$ra_table" 1)" api "$RA_PHASE — Ziel auf $ra_table (EnableTable über $ra_way)"
    else
      bf_expect "$(rt_target_of_row "$ra_table" 1)" NULL "$RA_PHASE — Ziel auf $ra_table ohne Regel (EnableTable über $ra_way)"
    fi
    ra_report="$ra_report $ra_table: $(rt_target_of_row "$ra_table" 1);"
  done
done
bf_expect "$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)" true "$RA_PHASE — Feed-Container läuft weiter"

echo "run-integration-tests: Routing-Aktivierung über die API (LH-FA-CFG-008) belegt — Tabellen mit applied-Regel api_orders (Region eu führt auf api), über EnableTable des gRPC- und des HTTP-Wegs ohne cdc.enable_table und ohne Neustart aktiviert, trugen das Ziel der Regel; je eine auf demselben Weg aktivierte Tabelle ohne Regel trug kein Ziel (Ziel je Tabelle:$ra_report)"

abdeckung_declare "Prozessstart-Vorlauf-Frist" "LH-FA-CFG-007,LH-QA-REL-001" "eine Sperre auf einer eigenen, bislang nicht aktivierten Tabelle hält den enable-Antrag über die Frist des Vorlaufs (30s, \`ADR-0128\`) hinaus fest; ein realer Neustart des Feed-Containers startet Stream und Administrations-Goroutine trotzdem — der Healthcheck bleibt während der Wartezeit gesund, eine auf einer bereits aktivierten Tabelle committete Änderung wird binnen Frist plus Toleranz über cdc.changes sichtbar, der Prozess endet nicht mit einer Fehlerklasse, und der Antrag wird erst nach Freigabe der Sperre applied" "Prozessstart-Vorlauf-Frist (ADR-0128) belegt"

# Prozessstart-Vorlauf-Frist (ADR-0128): eine zweite Sitzung sperrt die
# eigens dafür angelegte, bislang nicht aktivierte Tabelle
# $STARTGRENZE_TABLE mit `LOCK TABLE ...
# IN ACCESS EXCLUSIVE MODE` (bf_hold_start/bf_hold_end halten und lösen
# sie gezielt, länger als die Frist von 30s); `cdc.enable_table` legt
# währenddessen den Antrag an — `ALTER PUBLICATION ... ADD TABLE` braucht
# einen mit ACCESS EXCLUSIVE unvereinbaren Lock auf die Zieltabelle und
# hängt deshalb an der Sperre, unabhängig vom folgenden Neustart. Ein
# `docker restart` liest denselben, weiterhin `pending` stehenden Antrag im
# Vorlauf des neuen Prozesses: `Stream.Run` (`START_REPLICATION`) und die
# Administrations-Goroutine starten nach Ablauf der Frist trotzdem
# (`ADR-0128` Festlegung 1 und 3) — belegt über den eigenen
# `--healthcheck`-Ausgang des Containers (unabhängig vom Docker-eigenen
# Health-Status-Takt) und eine auf der dauerhaft aktivierten Tabelle
# feed_e2e_full committete Änderung. Nach Freigabe der Sperre verarbeitet
# die Administrations-Goroutine den Antrag ohne Frist zu Ende (applied).
STARTGRENZE_TABLE=feed_e2e_startgrenze
STARTGRENZE_HOLD=pgc-e2e-startgrenze-hold

bf_hold_start "$STARTGRENZE_HOLD" "BEGIN; LOCK TABLE public.$STARTGRENZE_TABLE IN ACCESS EXCLUSIVE MODE; SELECT pg_sleep(180);"
bf_await_sql "SELECT count(*) FROM pg_locks l JOIN pg_stat_activity a ON a.pid = l.pid WHERE a.application_name = '$STARTGRENZE_HOLD' AND l.mode = 'AccessExclusiveLock' AND l.granted" 1 15 "Prozessstart-Vorlauf-Frist — Sperre auf $STARTGRENZE_TABLE"

startgrenze_request_id=$(bf_sql "SELECT cdc.enable_table('src-e2e', 'public', '$STARTGRENZE_TABLE')")
[ -n "$startgrenze_request_id" ] || bf_fail "Prozessstart-Vorlauf-Frist — cdc.enable_table($STARTGRENZE_TABLE) lieferte keine Antrags-ID"

startgrenze_restart_at=$(date +%s)
docker restart "$FEED_CONTAINER" >/dev/null
sleep 1
bf_sql "INSERT INTO public.feed_e2e_full (id, name) VALUES (9001, 'NachDemVorlauf')" >/dev/null

# Während der Frist (bis zu 30s, ADR-0128 Festlegung 2) bleibt der
# Healthcheck gesund — mehrere Abfragen über das Fenster verteilt, direkt
# über den Container-eigenen `--healthcheck`-Ausgang (compose.yaml), nicht
# über Dockers eigenen, langsameren Health-Status-Takt. Die Änderung oben
# ist während dieser Schleife nicht sichtbar: der Stream liest erst nach
# Ablauf der Frist.
for startgrenze_probe in $(seq 1 9); do
  docker exec "$FEED_CONTAINER" /pg-change-feed --healthcheck \
    || bf_fail "Prozessstart-Vorlauf-Frist — Healthcheck meldet ungesund während der Wartezeit (Versuch $startgrenze_probe)"
  sleep 3
done

bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = 'feed_e2e_full' AND new_data->>'id' = '9001'" 1 20 "Prozessstart-Vorlauf-Frist — Änderung auf feed_e2e_full"
startgrenze_elapsed=$(( $(date +%s) - startgrenze_restart_at ))

feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
[ "$feed_running" = "true" ] || bf_fail "Prozessstart-Vorlauf-Frist — Feed-Container lief nach dem Neustart nicht mehr (der Prozess sollte trotz Frist-Ablauf nicht mit einer Fehlerklasse enden)"

bf_hold_end "$STARTGRENZE_HOLD"
bf_await_applied "$startgrenze_request_id" "Prozessstart-Vorlauf-Frist"

echo "run-integration-tests: Prozessstart-Vorlauf-Frist (ADR-0128) belegt — Neustart bei gesperrtem $STARTGRENZE_TABLE, Healthcheck blieb über $startgrenze_probe Abfragen gesund, Änderung auf feed_e2e_full nach ${startgrenze_elapsed}s seit dem Neustart erfasst, Antrag $startgrenze_request_id nach Freigabe der Sperre applied"

abdeckung_declare "Upgrade-Sicherheits-Rundlauf" "LH-QA-OPS-005,LH-FA-RET-001" "ein realer Container-Tausch ersetzt den Feed-Container durch eine neue Instanz desselben Images, während Datenbank und NATS unberührt bleiben — der Datenstand davor bleibt lesbar (\`count(*)\`-Beleg gegen cdc.changes nach dem Tausch, \`LH-FA-RET-001\`), danach Eingefügtes wird weiter erfasst" "Upgrade-Sicherheits-Rundlauf (LH-QA-OPS-005, ADR-0064) belegt"

# Upgrade-Sicherheits-Rundlauf (LH-QA-OPS-005): bildet den
# Mechanismus eines Anwendungs-Upgrades nach — ein realer Container-Tausch
# über `$COMPOSE up -d --force-recreate --no-deps pg-change-feed` ersetzt den
# Feed-Container durch eine neue Instanz desselben Images. Der Tausch belegt
# das Upgrade des Prozesses, nicht des Schemas: den zweiten
# `make schema-rollout`-Lauf gegen ein migriertes Ziel (idempotent über die
# elf bekannten Fremdobjekte des Idempotenz-Guards in
# `tools/schema/rolloutguard`) belegt
# `tools/harness/run-schema-rollout-guard-test.sh`. Der Rundlauf läuft hier,
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

# TestE2ESchemaChangeDropColumn (LH-FA-SCH-003) läuft als eigener go-test-Aufruf, NACH allem,
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
# Negative-Fall) läuft als eigener go-test-Aufruf: sie meldet eine nicht
# sicher als Obermenge erkennbare Typänderung sichtbar über die
# Fehlerklasse `schema` und beendet damit den Erfassungspfad des
# Feed-Containers dauerhaft (`restart: "no"`, kein Neustart-Vertrag,
# siehe Funktionskommentar). Alles, was den laufenden Container noch
# braucht (Lasttest-Beleg, Black-Box-CLI-Rundlauf oben), lief davor;
# TestE2ESchemaChangeDropColumn direkt davor beendete den Container
# ebenfalls bereits dauerhaft — der Neustart oben stellt einen
# sauberen, healthy Zustand wieder her, bevor diese Funktion ihren eigenen
# Mechanismus real prüft. Die Phase „Transformationen-Nichtanwendbarkeit
# und Abhilfe“ unten stellt den Container nach dieser Funktion noch einmal
# wieder her — sie ist der letzte Rundlauf des Runners.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_INTEGRATION_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test -v -run '^TestE2ESchemaChangeIncompatibleTypeChange$' ./test/integration/...

abdeckung_declare "Transformationen-Nichtanwendbarkeit und Abhilfe" "LH-FA-CFG-007,LH-FA-ADM-003,LH-FA-SCH-004" "eine kompatible Spalten-Erweiterung, deren Name den Zielnamen einer aktiven rename_column-Regel trifft, beendet den Erfassungspfad sichtbar mit der Fehlerklasse schema (Gegenprobe: eine Erweiterung mit anderem Namen tut es nicht); cdc.remove_transformation wird beantragt, während der Prozess steht, nach dem Neustart ist der Antrag applied und die zuvor nicht bestätigte Transaktion erscheint über cdc.changes in Rohform, ohne zweiten schema-Fehler" "Transformationen-Nichtanwendbarkeit und Abhilfe (ADR-0112 Folgepflicht 5) belegt"

# Transformationen-Nichtanwendbarkeit und Abhilfe (ADR-0112 Teilfrage 4 und
# Folgepflicht 5): läuft nach der zweiten Container-Ende-Grenze
# (TestE2ESchemaChangeIncompatibleTypeChange oben), auf einer eigenen
# Tabelle (feed_e2e_transform_abhilfe), unberührt von den beiden
# vorangegangenen Schema-Testfunktionen. Zwei Neustarts, mit
# unterschiedlicher Behandlung des Replication-Slots:
#
# 1. Erster Neustart (unten): stellt einen sauberen, healthy Zustand nach
#    TestE2ESchemaChangeIncompatibleTypeChange wieder her — derselbe
#    Neustart-und-Health-Poll-Rahmen wie zwischen TestE2ESchemaChangeDropColumn
#    und TestE2ESchemaChangeIncompatibleTypeChange, mit Slot-Neuanlage: die
#    von jener Funktion hinterlassene, am Slot unbestätigte Transaktion ist
#    echt inkompatibel — ein Replay ohne Slot-Neuanlage triggert denselben
#    schema-Fehler erneut.
# 2. Zweiter Neustart (nach der Kollision dieser Phase): lässt den Slot
#    unangetastet — die dabei am Slot unbestätigte Transaktion (id=3) ist
#    genau die, die nach der Abhilfe über cdc.changes erscheinen soll
#    (Kriterium (d)); eine Slot-Neuanlage verwirft sie unwiederbringlich.
TA_PHASE="Transformationen-Nichtanwendbarkeit und Abhilfe"
TA_TABLE=feed_e2e_transform_abhilfe
TA_RULE='{"kind":"rename_column","column":"name","to":"label"}'

docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -c \
  "SELECT pg_drop_replication_slot('$SLOT') WHERE EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = '$SLOT');" >/dev/null

docker start "$FEED_CONTAINER" >/dev/null
bf_await_healthy "$TA_PHASE — Neustart nach TestE2ESchemaChangeIncompatibleTypeChange"

bf_enable "$TA_TABLE" "$TA_PHASE"
tf_set "$TA_TABLE" umbenennung "$TA_RULE" "$TA_PHASE"

docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_INTEGRATION_DSN="$DSN" \
  "$TOOLCHAIN_IMAGE" go test -v -run '^TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass$' ./test/integration/...

# Kriterium (a): der Prozess endet sichtbar — Container-Lauf und Log-Sentinel.
ta_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
[ "$ta_running" = "false" ] || bf_fail "$TA_PHASE — Feed-Container lief nach der Kollision entgegen der Erwartung noch (der Prozess soll sichtbar mit Fehlerklasse schema enden)"

ta_log=$(docker logs "$FEED_CONTAINER" 2>&1 || true)
[[ "$ta_log" == *"Transformationsregel auf die Änderung nicht anwendbar"* ]] || bf_fail "$TA_PHASE — Container-Log trägt den Sentinel-Text der Fehlerklasse schema nicht: $ta_log"

# Kriterium (b): der Antrag wird beantragt, während der Prozess steht, und
# bleibt zunächst pending (SPEC-019-Status; „requested“ im ADR-Wortlaut) —
# reine SQL-Aufrufe, kein laufender Feed-Container nötig.
ta_remove_request=$(bf_sql "SELECT cdc.remove_transformation('src-e2e', 'public', '$TA_TABLE', 'umbenennung')")
[ -n "$ta_remove_request" ] || bf_fail "$TA_PHASE — cdc.remove_transformation($TA_TABLE) lieferte keine Antrags-ID"
bf_expect "$(bf_sql "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$ta_remove_request'")" pending "$TA_PHASE — Antrag $ta_remove_request während des Prozess-Stillstands"

# Kriterium (c)/(d): Neustart ohne Slot-Eingriff diesmal — die am Slot
# unbestätigte Transaktion (id=3) wird erneut vorgelegt, diesmal ohne die
# entfernte Regel.
docker start "$FEED_CONTAINER" >/dev/null
bf_await_healthy "$TA_PHASE — Neustart nach der Abhilfe"
bf_await_applied "$ta_remove_request" "$TA_PHASE"

ta_second_fault=$(bf_sql "SELECT coalesce(error_class, '') FROM cdc.heartbeat WHERE source_id = 'src-e2e'")
[ "$ta_second_fault" != "schema" ] || bf_fail "$TA_PHASE — ein zweiter schema-Fehler nach der Abhilfe (error_class=$ta_second_fault)"

bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$TA_TABLE' AND new_data->>'id' = '3'" 1 30 "$TA_PHASE — zuvor nicht bestätigte Transaktion (id=3)"
tf_form "$TA_TABLE" 3 '{"name":"Betroffen"}' "$TA_PHASE — Rohform nach der Abhilfe (Regel entfernt)"

echo "run-integration-tests: Transformationen-Nichtanwendbarkeit und Abhilfe (ADR-0112 Folgepflicht 5) belegt — Kollision auf $TA_TABLE beendete den Erfassungspfad real (error_class=schema, Sentinel im Log), cdc.remove_transformation ($ta_remove_request) während des Stillstands beantragt (pending), nach dem Neustart applied vor der ersten Transaktion der Tabelle, die zuvor nicht bestätigte Zeile (id=3) erscheint über cdc.changes in Rohform ohne zweiten schema-Fehler"

abdeckung_declare "Routing-Nichtanwendbarkeit und Abhilfe" "LH-FA-CFG-008,LH-FA-ADM-003,LH-FA-SCH-004" "eine auf eine Change nicht anwendbare Routing-Regel beendet den Erfassungspfad sichtbar mit der Fehlerklasse schema, ohne dass eine Change der Tabelle persistiert wird; die Erstaktivierung ohne Spaltenform und die Publication mit Spaltenliste erzeugen die Nichtanwendbarkeit am System, cdc.remove_route wird beantragt, während der Prozess steht, nach dem Neustart ist der Antrag applied und die zuvor nicht bestätigte Transaktion erscheint über cdc.changes ohne Ziel und ohne zweiten schema-Fehler; ein Spaltenentfernen mit bekannter Spaltenform endet als inkompatible Schemaänderung ohne Zeile für die Change in cdc.changes" "Routing-Nichtanwendbarkeit und Abhilfe (LH-FA-CFG-008) belegt"

# Routing-Nichtanwendbarkeit und Abhilfe: läuft als letzter Rundlauf dieses
# Runners, nach der Container-Ende-Grenze der Transformations-Abhilfe oben, auf
# eigenen Tabellen. Drei Fälle, jeder beendet den Feed-Container:
#
# (b) Erstaktivierung ohne Spaltenform: die Aktivierung legt die Versionszeile
#     ohne Spaltenform an; die Regel auf region wird beantragt, die Spalte
#     entfernt, bevor eine Change der Tabelle die Spaltenform anlegt. Die
#     erste Relation der Tabelle trägt region nicht mehr, die Change endet als
#     nicht anwendbare Regel.
# (c) Publication mit Spaltenliste, die region nicht trägt: die Spalte besteht
#     an der Quelle, die Relation trägt sie nicht.
#     Beide Fälle tragen die Abhilfe: cdc.remove_route, während der Prozess
#     steht, und ein Neustart ohne Slot-Eingriff legen die zuvor nicht
#     bestätigte Transaktion vor.
# (a) Spaltenform bekannt, Spalte entfernt: der Pfad endet an der inkompatiblen
#     Schemaänderung, vor der Regel; keine Zeile für diese Change. Dieser Fall
#     bleibt am Ende stehen, ohne Wiederanlauf.
RN_PHASE="Routing-Nichtanwendbarkeit und Abhilfe"
RN_RULE='{"target":"eu","order":10,"when":{"column":"region","equals":"eu"}}'
RN_SENTINEL_NOT_APPLICABLE="Routing-Regel auf die Änderung nicht anwendbar"
RN_SENTINEL_INCOMPATIBLE="Relation-Änderung nicht sicher als Obermenge interpretierbar"

# rn_shape_rows <Tabelle>: die Zahl der Spaltenform-Zeilen der Versionen einer
# Tabelle in cdc.table_schema.
rn_shape_rows() {
  bf_sql "SELECT count(*) FROM cdc.table_schema ts JOIN cdc.schema_version sv ON sv.schema_version_id = ts.schema_version_id JOIN cdc.source_table st ON st.source_table_id = sv.source_table_id WHERE st.source_id = 'src-e2e' AND st.schema_name = 'public' AND st.table_name = '$1'"
}

# rn_await_end <Phase>: wartet, bis der Feed-Container beendet ist.
rn_await_end() {
  local i
  for ((i = 0; i < 90; i++)); do
    if [ "$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)" = false ]; then
      return 0
    fi
    sleep 1
  done
  bf_fail "$1 — der Feed-Container lief 90 s nach der auslösenden Änderung weiter"
}

# rn_expect_end <Phase> <Seit-Zeitpunkt> <Log-Sentinel>: der Prozess ist
# beendet, sein Log seit dem Zeitpunkt trägt den Sentinel, und
# cdc.process_heartbeat trägt die Klasse schema.
rn_expect_end() {
  local phase=$1 since=$2 sentinel=$3 log
  rn_await_end "$phase"
  log=$(docker logs --since "$since" "$FEED_CONTAINER" 2>&1 || true)
  [[ "$log" == *"$sentinel"* ]] || bf_fail "$phase — das Log des Feed-Containers seit $since trägt den Sentinel '$sentinel' nicht: $log"
  bf_expect "$(bf_sql "SELECT coalesce(error_class, '') FROM cdc.process_heartbeat WHERE source_id = 'src-e2e'")" schema "$phase — Fehlerklasse in cdc.process_heartbeat"
}

# rn_abhilfe <Tabelle> <Phase>: cdc.remove_route wird beantragt, während der
# Prozess steht (Antrag pending); der Neustart verarbeitet ihn vor der ersten
# Transaktion der Tabelle, und die zuvor nicht bestätigte Zeile (id=1) erscheint
# über cdc.changes ohne Ziel und ohne zweiten schema-Fehler.
rn_abhilfe() {
  local table=$1 phase=$2 request_id fault
  request_id=$(bf_sql "SELECT cdc.remove_route('src-e2e', 'public', '$table', 'eu_orders')")
  [ -n "$request_id" ] || bf_fail "$phase — cdc.remove_route($table) lieferte keine Antrags-ID"
  bf_expect "$(bf_sql "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$request_id'")" pending "$phase — Antrag $request_id während des Prozess-Stillstands"
  docker start "$FEED_CONTAINER" >/dev/null
  bf_await_healthy "$phase — Neustart nach der Abhilfe"
  bf_await_applied "$request_id" "$phase"
  fault=$(bf_sql "SELECT coalesce(error_class, '') FROM cdc.heartbeat WHERE source_id = 'src-e2e'")
  [ "$fault" != "schema" ] || bf_fail "$phase — ein zweiter schema-Fehler nach der Abhilfe (error_class=$fault)"
  bf_await_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$table' AND new_data->>'id' = '1'" 1 30 "$phase — zuvor nicht bestätigte Transaktion (id=1)"
  bf_expect "$(rt_target_of_row "$table" 1)" NULL "$phase — Ziel der Change nach der Abhilfe (Regel entfernt)"
}

RN_B_TABLE=feed_e2e_route_nb
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$RN_B_TABLE (id int PRIMARY KEY, name text, region text);
SQL
bf_enable "$RN_B_TABLE" "$RN_PHASE"
rn_b_shape=$(rn_shape_rows "$RN_B_TABLE")
bf_expect "$rn_b_shape" 0 "$RN_PHASE (b) — Spaltenform-Zeilen der Erstaktivierung"
rt_set "$RN_B_TABLE" eu_orders "$RN_RULE" "$RN_PHASE (b)"
bf_sql "ALTER TABLE public.$RN_B_TABLE DROP COLUMN region" >/dev/null
rn_b_since=$(date -u +%Y-%m-%dT%H:%M:%S.%NZ)
bf_sql "INSERT INTO public.$RN_B_TABLE (id, name) VALUES (1, 'Betroffen')" >/dev/null
rn_expect_end "$RN_PHASE (b)" "$rn_b_since" "$RN_SENTINEL_NOT_APPLICABLE"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RN_B_TABLE'")" 0 "$RN_PHASE (b) — Changes der Tabelle bei beendetem Prozess"
rn_abhilfe "$RN_B_TABLE" "$RN_PHASE (b)"
tf_form "$RN_B_TABLE" 1 '{"name":"Betroffen"}' "$RN_PHASE (b) — Bild nach der Abhilfe"

RN_C_TABLE=feed_e2e_route_pub
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$RN_C_TABLE (id int PRIMARY KEY, name text, region text);
SQL
bf_enable "$RN_C_TABLE" "$RN_PHASE"
rn_c_shape=$(rn_shape_rows "$RN_C_TABLE")
bf_expect "$rn_c_shape" 0 "$RN_PHASE (c) — Spaltenform-Zeilen der Erstaktivierung"
rt_set "$RN_C_TABLE" eu_orders "$RN_RULE" "$RN_PHASE (c)"
bf_sql "ALTER PUBLICATION pub_pgc_e2e DROP TABLE public.$RN_C_TABLE" >/dev/null
bf_sql "ALTER PUBLICATION pub_pgc_e2e ADD TABLE public.$RN_C_TABLE (id, name)" >/dev/null
bf_expect "$(bf_sql "SELECT array_to_string(attnames, ',') FROM pg_publication_tables WHERE pubname = 'pub_pgc_e2e' AND tablename = '$RN_C_TABLE'")" "id,name" "$RN_PHASE (c) — Spaltenliste der Publication"
rn_c_since=$(date -u +%Y-%m-%dT%H:%M:%S.%NZ)
bf_sql "INSERT INTO public.$RN_C_TABLE (id, name, region) VALUES (1, 'Betroffen', 'eu')" >/dev/null
rn_expect_end "$RN_PHASE (c)" "$rn_c_since" "$RN_SENTINEL_NOT_APPLICABLE"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RN_C_TABLE'")" 0 "$RN_PHASE (c) — Changes der Tabelle bei beendetem Prozess"
rn_abhilfe "$RN_C_TABLE" "$RN_PHASE (c)"

RN_A_TABLE=feed_e2e_route_drop
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE public.$RN_A_TABLE (id int PRIMARY KEY, name text, region text);
SQL
bf_enable "$RN_A_TABLE" "$RN_PHASE"
rt_set "$RN_A_TABLE" eu_orders "$RN_RULE" "$RN_PHASE (a)"
rt_row "$RN_A_TABLE" 1 RnEins "'eu'" "$RN_PHASE (a)"
bf_expect "$(rt_target_of_row "$RN_A_TABLE" 1)" eu "$RN_PHASE (a) — Ziel der Change vor der Spaltenentfernung"
rn_a_shape=$(rn_shape_rows "$RN_A_TABLE")
[ "$rn_a_shape" -gt 0 ] || bf_fail "$RN_PHASE (a) — die erste Change legte keine Spaltenform an"
bf_sql "ALTER TABLE public.$RN_A_TABLE DROP COLUMN region" >/dev/null
rn_a_since=$(date -u +%Y-%m-%dT%H:%M:%S.%NZ)
bf_sql "INSERT INTO public.$RN_A_TABLE (id, name) VALUES (2, 'Betroffen')" >/dev/null
rn_expect_end "$RN_PHASE (a)" "$rn_a_since" "$RN_SENTINEL_INCOMPATIBLE"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RN_A_TABLE' AND new_data->>'id' = '2'")" 0 "$RN_PHASE (a) — Zeilen für die Change nach der Spaltenentfernung (weder NULL noch ein Ziel)"
bf_expect "$(bf_sql "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND table_name = '$RN_A_TABLE'")" 1 "$RN_PHASE (a) — Changes der Tabelle insgesamt"
bf_expect "$(rt_target_of_row "$RN_A_TABLE" 1)" eu "$RN_PHASE (a) — Ziel der früheren Change unverändert"

echo "run-integration-tests: Routing-Nichtanwendbarkeit und Abhilfe (LH-FA-CFG-008) belegt — (b) Erstaktivierung ohne Spaltenform ($rn_b_shape Spaltenform-Zeilen nach der Aktivierung) und (c) Publication-Spaltenliste id,name ($rn_c_shape Spaltenform-Zeilen) beendeten den Erfassungspfad real mit Klasse schema und dem Sentinel '$RN_SENTINEL_NOT_APPLICABLE' ohne persistierte Change; cdc.remove_route wurde bei stehendem Prozess beantragt (pending), nach dem Neustart applied, die nicht bestätigte Zeile erschien über cdc.changes ohne Ziel und ohne zweiten schema-Fehler; (a) Spaltenentfernung bei bekannter Spaltenform ($rn_a_shape Spaltenform-Zeilen) beendete den Pfad mit Klasse schema und dem Sentinel '$RN_SENTINEL_INCOMPATIBLE', die Change nach der Entfernung hat keine Zeile in cdc.changes"

# Zusammensetzung der E2E-Abdeckungstabelle: erst Go-Zeilen nach Quelldatei
# und Quellzeile, dann die Bash-Zeilen der deklarierten Phasen nach
# Runner-Zeile. Die Deklarations-Anker sind an dieser Stelle bereits alle
# gelaufen; eine fehlende Datei entsteht hier und bleibt sonst unangetastet.
abdeckung_schreiben "$ABDECKUNG_GO_ZEILEN"
echo "run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus $(( $(printf '%s\n' "$ABDECKUNG_GO_ZEILEN" | grep -c .) )) Go-Zeilen und ${#abdeckung_bash_zeilen[@]} Bash-Zeilen"
