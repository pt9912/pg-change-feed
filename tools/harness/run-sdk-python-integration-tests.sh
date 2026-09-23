#!/usr/bin/env bash
# run-sdk-python-integration-tests.sh — Realserver-Integrationstest der
# Python-SDK-Zustellweg-Flaechen (ADR-0110 §Entscheidung Festlegung
# 2/Folgepflicht 1, slice-sdk-python-grpc-client-flaeche): jede neue
# Zustellweg-Flaechen dieses Packages prueft ihre Protokoll-Annahmen
# zusaetzlich zu den netzlosen Unit-Tests gegen eine reale, laufende
# Server-Instanz.
#
# Kette je Lauf — Bring-up-Phase im Muster tools/harness/run-integration-tests.sh,
# aber eigenstaendig und bewusst schmaler: dieses Skript traegt keinen
# Server-E2E-Beleg (LH-QA-POR-003 bleibt Eigentum des bestehenden Runners,
# der unveraendert bleibt); es baut die compose.yaml-Umgebung hoch
# (PostgreSQL/NATS/Feed-Container, Schema-Rollout ueber d-migrate,
# Vorbedingungen der Aktivierung), baut die `integration`-Docker-Stufe des
# Python-SDK (sdks/python/Dockerfile, zwingend mit dem benannten Bau-Kontext
# `proto` — ohne ihn bricht der Bau an der COPY --from=proto-Zeile ab) und
# startet sie im selben Docker-Netz wie den Feed-Container. Der Pruefling
# ist das SDK selbst: der Integrationstest importiert die Client-Flaechen
# des Packages direkt (kein Wegwerf-Duplikat-Client daneben). Die
# stdout-Marker der Testlaeufe (READY/RECEIVED/REJECTED) liest dieser
# Runner ueber `docker logs`, waehrend der Test noch laeuft — deshalb
# ungepuffert (`python -u`, Stufe-CMD). Die `change_id` der
# empfangenen Change wird zusaetzlich gegen den Lesezugriffsweg
# `cdc.changes` gehalten (dieselbe Disziplin wie die Rundlaeufe im
# Server-E2E-Runner).
#
# Eine Phase je Flaeche (slice-sdk-python-grpc-client-flaeche: gRPC,
# SPEC-020; slice-sdk-python-sse-client-flaeche: SSE, SPEC-021;
# slice-sdk-python-nats-stream-client-flaeche: NATS-Vollinhalt, SPEC-024;
# slice-sdk-python-http-reale2e: HTTP-API, SPEC-018 — der Rundlauf
# RegisterConsumer/ListTables mit SQL-Gegenpruefung gegen cdc.consumer):
# jede Phase traegt ihre Testdatei explizit als Umgebungsvariable
# `PGCHANGEFEED_TEST_FILE` (kein stiller
# Ausschluss des Rests, Muster der -run-Muster im Server-E2E-Runner) und
# traegt eigenen Sentinel- und ID-Wertebereich, damit sich die Phasen
# nicht in die Quere kommen.
#
# Voraussetzungen: Docker, ein geladenes :dev-Image (`make image` vorher —
# compose.yaml traegt keinen build:-Block, ADR-0044) und Netz (pip-Paketbezug
# im SDK-Bau, d-migrate-Image). Kein Gate (dieselbe Klasse wie
# `make test-integration`). Cleanup je Ausgang: Test-Container wegwerfen,
# compose-Stack inkl. Volume abrueumen.
#
# Aufruf: `make test-sdk-python-integration`. Override:
# SDK_PYTHON_INTEGRATION_IMAGE.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

COMPOSE=${COMPOSE:-docker compose -f compose.yaml}
NETWORK=${NETWORK:-cdc-feed-test}
PG_CONTAINER=${PG_CONTAINER:-cdc-test-postgres}
FEED_CONTAINER=${FEED_CONTAINER:-cdc-test-feed}
PG_DB=cdc
PG_USER=postgres
PG_PASSWORD=postgres
# Derselbe Slot-Name wie der Container-Vertrag in compose.yaml (CDC_SLOT).
SLOT=slot_pgc_e2e
# Derselbe Reader-Token wie der Container-Vertrag in compose.yaml
# (CDC_API_TOKEN_READER): gRPC-Stream und SSE-Endpunkt akzeptieren beide
# Token-Klassen (SPEC-020/SPEC-021), die Tests tragen die Reader-Klasse.
API_TOKEN=e2e-reader-token
API_TOKEN_ADMIN=e2e-admin-token
API_TOKEN_READER=e2e-reader-token
DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"

SDK_INTEGRATION_IMAGE=${SDK_PYTHON_INTEGRATION_IMAGE:-pg-change-feed:sdk-python-integration}
SDK_TEST_CONTAINER=cdc-sdk-python-client-test

# Dasselbe aktivierte Tisch-Set wie der Container-Vertrag (CDC_TABLES):
# die Aktivierung traegt die Verdrahtung des Feed-Containers (ADR-0028),
# die Tabellen muessen dafuer physisch existieren. Gegriffen wird fuer die
# Belege auf feed_e2e_full (REPLICA IDENTITY FULL).
TEST_TABLE=feed_e2e_full
GRPC_TEST_FILE=integration/test_grpc_realserver.py
SSE_TEST_FILE=integration/test_sse_realserver.py
HTTP_TEST_FILE=integration/test_http_realserver.py
NATS_TEST_FILE=integration/test_nats_realserver.py
GRPC_SENTINEL=PythonGrpcSdkE2ESentinel
SSE_SENTINEL=PythonSseSdkE2ESentinel
HTTP_SENTINEL=PythonHttpSdkE2ESentinel
HTTP_PUBLICATION=pub_pgc_e2e
NATS_SENTINEL=PythonNatsSdkE2ESentinel
NATS_STREAM_TOKEN=e2e-nats-stream-token

cleanup() {
  docker rm -f "$SDK_TEST_CONTAINER" >/dev/null 2>&1 || true
  $COMPOSE down -v --remove-orphans >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker image inspect ghcr.io/pt9912/pg-change-feed:dev >/dev/null 2>&1 || {
  echo "run-sdk-python-integration-tests: kein geladenes :dev-Image — make image vorher (compose.yaml traegt keinen build:-Block, ADR-0044)" >&2
  exit 1
}

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
  echo "run-sdk-python-integration-tests: Compose-PostgreSQL wurde nicht bereit" >&2
  exit 1
fi

# Schema-Rollout über d-migrate vor dem Lauf (ADR-0043) — dieselbe Kette
# wie der Server-E2E-Runner.
make schema-rollout SCHEMA_TARGET="db:$DSN" SCHEMA_ROLLOUT_NETWORK="$NETWORK"

# Vorbedingungen der Aktivierung: die drei CDC_TABLES-Tabellen existieren
# physisch, die Quelle-Zeile traegt die Fremdschluessel-Registrierung
# (SPEC-001); die Aktivierung selbst (Publication, Bindungs- und
# Schema-Version-Zeilen) traegt die Verdrahtung des Feed-Containers
# (ADR-0028, LH-FA-CFG-001.a). Die REPLICA IDENTITY FULL bleibt vom
# Aktivieren unberuehrt (LH-FA-CFG-001.a).
docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE public.feed_e2e_flow (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_full (id int PRIMARY KEY, name text);
ALTER TABLE public.feed_e2e_full REPLICA IDENTITY FULL;
CREATE TABLE public.feed_e2e_schema (id int PRIMARY KEY, name text, amount text);
INSERT INTO cdc.source (source_id, name) VALUES ('src-e2e', 'E2E-Quelle');
SQL

$COMPOSE up -d pg-change-feed >/dev/null

# Stream-Vertrag am laufenden Container (Muster run-integration-tests.sh):
# die Verdrahtung steht, wenn der Stream seinen Slot angelegt hat und der
# Container laeuft; der Compose-Healthcheck belegt den Prozesszustand
# darueber hinaus.
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
  echo "run-sdk-python-integration-tests: Feed-Container trägt die Verdrahtung nicht — Slot $SLOT: ${slot_ok:-0}, Feed-Lauf: $feed_running (Feed-Ausgang: $feed_exit)" >&2
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
  echo "run-sdk-python-integration-tests: Feed-Container meldet Compose-Health-Status ${health:-fehlt}, wollen healthy" >&2
  exit 1
fi

# Die `integration`-Stufe des Python-SDK bauen — zwingend mit dem benannten
# Bau-Kontext `proto` (COPY --from=proto im Dockerfile, ADR-0090
# Festlegung 2).
docker build --build-context proto=proto -f sdks/python/Dockerfile \
  --target integration -t "$SDK_INTEGRATION_IMAGE" sdks/python

# run_surface_phase <Name> <Testdatei> <Sentinel> <ID-Basis> <Reject-Marker>
# <Extra-Env> <Received-Grep> <SQL-Variante: changes|consumer> — ein
# Realserver-Rundlauf fuer genau eine SDK-Flaeche: Container starten (die
# Adress-/Auth-Variablen je Flaeche kommen als Leerzeichen-getrennte
# Extra-Env-Liste herein), auf READY warten, begrenzte Folge eindeutiger
# Zeilen committen (Fire-and-Forget-Fenster, SPEC-020/SPEC-021/SPEC-024),
# auf den Reject-Marker und das Prozessende warten, die RECEIVED-Zeile
# pruefen und die Identitaet unabhaengig gegen den SQL-Lesezugriffsweg
# halten. Die SQL-Variante "changes" haelt die change_id gegen cdc.changes
# (Streams), "consumer" die Consumer-Registrierung gegen cdc.consumer
# (HTTP — Muster run-sdk-csharp-integration-tests.sh).
run_surface_phase() {
  local phase_name=$1 test_file=$2 sentinel=$3 id_base=$4 reject_marker=$5 extra_env=$6 \
        received_grep=$7 sql_kind=$8
  local attempt insert_id captured ident pair

  local env_args=()
  for pair in $extra_env; do
    env_args+=(-e "$pair")
  done

  docker rm -f "$SDK_TEST_CONTAINER" >/dev/null 2>&1 || true
  docker run -d --name "$SDK_TEST_CONTAINER" --network "$NETWORK" \
    "${env_args[@]}" \
    -e PGCHANGEFEED_API_TOKEN="$API_TOKEN" \
    -e PGCHANGEFEED_E2E_TABLE="$TEST_TABLE" \
    -e PGCHANGEFEED_E2E_SENTINEL="$sentinel" \
    -e PGCHANGEFEED_TEST_FILE="$test_file" \
    "$SDK_INTEGRATION_IMAGE" >/dev/null

  test_ready=0
  for _ in $(seq 1 60); do
    if docker logs "$SDK_TEST_CONTAINER" 2>/dev/null | grep -qF "READY"; then
      test_ready=1
      break
    fi
    if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
      break
    fi
    sleep 1
  done
  if [ "$test_ready" -ne 1 ]; then
    echo "run-sdk-python-integration-tests: $phase_name — Test wurde nicht innerhalb der Zeitspanne bereit (kein READY): $(docker logs "$SDK_TEST_CONTAINER" 2>&1 || true)" >&2
    exit 1
  fi

  local received=0
  for attempt in 1 2 3 4 5; do
    insert_id=$((id_base + attempt))
    # >/dev/null: die psql-INSERT-Echo-Zeile („INSERT 0 1") gehoert nicht in
    # den Funktions-stdout — der Aufrufer haelt hier nur den Rueckgabewert
    # (die change_id).
    docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 >/dev/null <<SQL
INSERT INTO public.$TEST_TABLE (id, name) VALUES ($insert_id, '$sentinel');
SQL
    for _ in $(seq 1 20); do
      if docker logs "$SDK_TEST_CONTAINER" 2>/dev/null | grep -qF "RECEIVED"; then
        received=1
        break
      fi
      if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
        break
      fi
      sleep 0.5
    done
    if [ "$received" -eq 1 ]; then
      break
    fi
  done

  local rejected=0
  for _ in $(seq 1 40); do
    if docker logs "$SDK_TEST_CONTAINER" 2>/dev/null | grep -qF "$reject_marker"; then
      rejected=1
      break
    fi
    if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
      break
    fi
    sleep 0.5
  done

  local test_stopped=0
  for _ in $(seq 1 20); do
    if [ "$(docker inspect --format '{{.State.Running}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo false)" != "true" ]; then
      test_stopped=1
      break
    fi
    sleep 0.5
  done
  local test_exit
  test_exit=$(docker inspect --format '{{.State.ExitCode}}' "$SDK_TEST_CONTAINER" 2>/dev/null || echo unbekannt)
  local test_output
  test_output=$(docker logs "$SDK_TEST_CONTAINER" 2>&1 || true)

  if [ "$received" -ne 1 ]; then
    echo "run-sdk-python-integration-tests: $phase_name — der Test lieferte den Happy-Path-Marker nicht ($TEST_TABLE, $sentinel): $test_output" >&2
    exit 1
  fi
  if ! printf '%s' "$test_output" | grep -qE "$received_grep"; then
    echo "run-sdk-python-integration-tests: $phase_name — die RECEIVED-Zeile traegt nicht die erwartete Form: $test_output" >&2
    exit 1
  fi
  if [ "$rejected" -ne 1 ]; then
    echo "run-sdk-python-integration-tests: $phase_name — der Ablehnungs-Beleg blieb aus ($reject_marker fehlt): $test_output" >&2
    exit 1
  fi
  if [ "$test_stopped" -ne 1 ] || [ "$test_exit" != "0" ]; then
    echo "run-sdk-python-integration-tests: $phase_name — der Test endete nicht mit Ausgang 0 (gestoppt: $test_stopped, Ausgang: $test_exit): $test_output" >&2
    exit 1
  fi

  ident=$(printf '%s' "$test_output" | grep -oE 'RECEIVED [a-z_]+=[^ ]+' | head -n1 | cut -d= -f2 || true)
  if [ -z "$ident" ]; then
    echo "run-sdk-python-integration-tests: $phase_name — die RECEIVED-Zeile traegt keine Identitaet: $test_output" >&2
    exit 1
  fi
  if [ "$sql_kind" = "changes" ]; then
    captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT count(*) FROM cdc.changes WHERE source_id = 'src-e2e' AND change_id = '$ident' AND table_name = '$TEST_TABLE' AND new_data->>'name' = '$sentinel'")
  else
    captured=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT count(*) FROM cdc.consumer WHERE consumer_id = '$ident'")
  fi
  if [ -z "$captured" ] || [ "$captured" -lt 1 ]; then
    echo "run-sdk-python-integration-tests: $phase_name — die Identitaet ($ident) ist nicht real ueber den SQL-Lesezugriffsweg lesbar (count=${captured:-leer})" >&2
    exit 1
  fi

  feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
  if [ "$feed_running" != "true" ]; then
    echo "run-sdk-python-integration-tests: $phase_name — Feed-Container lief nach dem Lauf nicht mehr weiter (kein Neustart erwartet)" >&2
    exit 1
  fi

  printf '%s\n' "$ident"
}

GRPC_CHANGE_ID=$(run_surface_phase \
  "gRPC-Flaeche (SPEC-020)" \
  "$GRPC_TEST_FILE" \
  "$GRPC_SENTINEL" 300 \
  "REJECTED code=Unauthenticated" \
  "PGCHANGEFEED_GRPC_ADDR=pg-change-feed:9090" \
  "RECEIVED change_id=[^ ]+ table=$TEST_TABLE .*operation=INSERT .*new_image=.*$GRPC_SENTINEL" \
  changes)

SSE_CHANGE_ID=$(run_surface_phase \
  "SSE-Flaeche (SPEC-021)" \
  "$SSE_TEST_FILE" \
  "$SSE_SENTINEL" 310 \
  "REJECTED status=401" \
  "PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090" \
  "RECEIVED change_id=[^ ]+ table=$TEST_TABLE .*operation=INSERT .*new_image=.*$SSE_SENTINEL" \
  changes)

NATS_CHANGE_ID=$(run_surface_phase \
  "NATS-Vollinhalts-Flaeche (SPEC-024)" \
  "$NATS_TEST_FILE" \
  "$NATS_SENTINEL" 320 \
  "REJECTED token-rejected" \
  "PGCHANGEFEED_NATS_URL=nats://nats:4222 PGCHANGEFEED_NATS_STREAM_TOKEN=$NATS_STREAM_TOKEN PGCHANGEFEED_SOURCE_ID=src-e2e" \
  "RECEIVED change_id=[^ ]+ table=$TEST_TABLE .*operation=INSERT .*new_image=.*$NATS_SENTINEL" \
  changes)

HTTP_IDENT=$(run_surface_phase \
  "HTTP-Flaeche (SPEC-018)" \
  "$HTTP_TEST_FILE" \
  "$HTTP_SENTINEL" 330 \
  "REJECTED status=401" \
  "PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090 PGCHANGEFEED_API_TOKEN_ADMIN=$API_TOKEN_ADMIN PGCHANGEFEED_API_TOKEN_READER=$API_TOKEN_READER PGCHANGEFEED_SOURCE_ID=src-e2e PGCHANGEFEED_HTTP_PUBLICATION=$HTTP_PUBLICATION" \
  "RECEIVED consumer_id=[^ ]+" \
  consumer)

# --- Abdeckungs-Traeger (docs/user/sdk-e2e-abdeckung.md) -------------------
# Der Python-Abschnitt entsteht aus derselben Messung, die ihn belegt;
# geschrieben wird nur bei inhaltlicher Abweichung (Temp-Datei + cmp),
# marker-gegrenzt. Writer-Form-Grenze (Muster csharp-/kotlin-Runner): dieser
# Runner haelt Kopf und fremde Abschnitte (C# und Kotlin) auf BEIDEN Seiten
# stabil und ersetzt nur seinen eigenen marker-gegrenzten Abschnitt; fehlt
# die Traeger-Datei ganz, regeneriert er Kopf und Tabellenkopf (kein
# nackter Abschnitt).
ABDECKUNG_ZIEL_DATEI=docs/user/sdk-e2e-abdeckung.md

abdeckung_python_abschnitt() {
  printf '%s\n' \
    '<!-- pgchangefeed-sdk-e2e:python-begin -->' \
    "| [\`LH-FA-SST-006\`](../../spec/lastenheft.md), [\`LH-FA-CON-001\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein Python-SDK-Client (\`PgChangeFeedHttpClient\`) registriert real einen Consumer (admin-Token) und listet Tabellen (reader-Token); die Registrierung ist unabhaengig ueber \`cdc.consumer\` lesbar; ein Aufruf ohne gueltiges Token endet mit HTTP-Status 401 | \`test_http_realserver.py\` | \`tools/harness/run-sdk-python-integration-tests.sh\` |" \
    "| [\`LH-FA-SST-008\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein Python-SDK-Client (\`PgChangeFeedGrpcClient\`) oeffnet real den gRPC-Server-Stream gegen den laufenden Feed-Container und empfaengt eine danach committete Aenderung; ein Oeffnungsversuch ohne gueltiges Token endet mit gRPC-Status \`Unauthenticated\` | \`test_grpc_realserver.py\` | \`tools/harness/run-sdk-python-integration-tests.sh\` |" \
    "| [\`LH-FA-SST-008\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein Python-SDK-Client (\`PgChangeFeedSseClient\`) oeffnet real \`GET /changes/stream\` und empfaengt eine danach committete Aenderung; ein Aufruf ohne gueltiges Token endet mit HTTP-Status 401 | \`test_sse_realserver.py\` | \`tools/harness/run-sdk-python-integration-tests.sh\` |" \
    "| [\`LH-FA-SST-008\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein Python-SDK-Client (\`PgChangeFeedNatsStreamClient\`) verbindet sich real per NATS und empfaengt eine danach committete Aenderung als vollstaendiges JSON-Event; ein Verbindungsversuch mit falschem Token wird vom NATS-Server abgelehnt | \`test_nats_realserver.py\` | \`tools/harness/run-sdk-python-integration-tests.sh\` |" \
    '<!-- pgchangefeed-sdk-e2e:python-end -->'
}

abdeckung_schreiben() {
  local temp vor="" nach=""
  local abschnitt
  abschnitt=$(abdeckung_python_abschnitt)
  if [ -f "$ABDECKUNG_ZIEL_DATEI" ]; then
    if grep -q "pgchangefeed-sdk-e2e:python-begin" "$ABDECKUNG_ZIEL_DATEI"; then
      vor=$(awk '/pgchangefeed-sdk-e2e:python-begin/{ausgang=1} !ausgang{print}' "$ABDECKUNG_ZIEL_DATEI")
      nach=$(awk 'f{print} /pgchangefeed-sdk-e2e:python-end/{f=1}' "$ABDECKUNG_ZIEL_DATEI")
    else
      vor=$(cat "$ABDECKUNG_ZIEL_DATEI")
    fi
  fi
  temp=$(mktemp)
  {
    if [ -z "$vor" ]; then
      printf '%s\n' \
        '# SDK-E2E-Abdeckung je Spec-Kennung' \
        '' \
        'Erzeugt von `make test-sdk-csharp-integration` über' \
        '`tools/harness/run-sdk-csharp-integration-tests.sh`; die Sprach-Runner' \
        'der Folge-Slices (Kotlin, Python-HTTP) erweitern dieselbe Datei um ihre' \
        'marker-gegrenzten Abschnitte. Je Sprach-Abschnitt deklariert der' \
        'zustaendige Runner seine Realserver-Phasen an Ort und Stelle. Diese' \
        'Datei ist eine **stabile Abdeckungs-Deklaration**, kein Lauf-Beleg: der' \
        'Runner schreibt sie nur bei inhaltlicher Abweichung. Sie trägt nur' \
        'Zeilen real existierender Runner-Phasen — ein Beleg steht hier nie,' \
        'bevor sein Lauf grün lief.' \
        '' \
        '| Spec-Kennung | Kurzbeschreibung | Nachweis | Ort |' \
        '| --- | --- | --- | --- |'
    fi
    if [ -n "$vor" ]; then
      printf '%s\n' "$vor"
    fi
    printf '%s\n' "$abschnitt"
    if [ -n "$nach" ]; then
      printf '%s\n' "$nach"
    fi
  } > "$temp"
  if [ -f "$ABDECKUNG_ZIEL_DATEI" ] && cmp -s "$temp" "$ABDECKUNG_ZIEL_DATEI"; then
    rm -f "$temp"
    echo "run-sdk-python-integration-tests: Abdeckungs-Traeger unveraendert — $ABDECKUNG_ZIEL_DATEI entspricht dem Quelltext-Stand"
  else
    chmod 0644 "$temp"
    mv "$temp" "$ABDECKUNG_ZIEL_DATEI"
    echo "run-sdk-python-integration-tests: Abdeckungs-Traeger geschrieben — $ABDECKUNG_ZIEL_DATEI"
  fi
}

abdeckung_schreiben

echo "run-sdk-python-integration-tests: SDK-Realserver-Belege (ADR-0110 Festlegung 2/Folgepflicht 1) gruen — gRPC-Flaeche (pgchangefeed.grpc_client, pg-change-feed:9090, change_id=$GRPC_CHANGE_ID), SSE-Flaeche (pgchangefeed.sse_client, pg-change-feed:8090, change_id=$SSE_CHANGE_ID) und NATS-Vollinhalts-Flaeche (pgchangefeed.nats_stream_client, nats://nats:4222, change_id=$NATS_CHANGE_ID) oeffneten real ihre Server-Streams gegen den laufenden Feed-Container und empfingen je eine danach committete Aenderung (Tabelle, Operation und Sentinel real am Wire; change_id je unabhaengig ueber cdc.changes lesbar), die HTTP-Flaeche (pgchangefeed.http_client) registrierte real einen Consumer (consumer_id=$HTTP_IDENT, unabhaengig ueber cdc.consumer lesbar) und listete Tabellen; ein Aufruf ohne gueltiges Token endete je mit gRPC-Status Unauthenticated, HTTP-Status 401 bzw. der laut ablehnenden NATS-Verbindungsablehnung"