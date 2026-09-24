#!/usr/bin/env bash
# run-sdk-csharp-integration-tests.sh — Realserver-Integrationstest der
# C#-SDK-Zustellweg-Flaechen (slice-sdk-csharp-reale2e; Mechanik-Klasse
# ADR-0110 §Entscheidung Festlegung 2, gespiegelt vom Python-Vorbild
# tools/harness/run-sdk-python-integration-tests.sh): die vier FlaecheN des
# Packages PgChangeFeed.Client (HTTP SPEC-018, gRPC SPEC-020, SSE SPEC-021,
# NATS-Vollinhalt SPEC-024) pruefen ihre Protokoll-Annahmen je gegen eine
# reale, laufende Server-Instanz — der Pruefling ist die kompilierte
# Client-Assembly, der Integrationstest importiert PgChangeFeed.Client
# direkt (kein Wegwerf-Duplikat-Client daneben).
#
# Kette je Lauf: compose.yaml-Umgebung hochfahren (PostgreSQL/NATS/
# Feed-Container, Schema-Rollout ueber d-migrate, Vorbedingungen der
# Aktivierung), Bau der `integration`-Docker-Stufe des C#-SDK
# (sdks/csharp/Dockerfile), Start je Phase im selben Docker-Netz wie der
# Feed-Container; die stdout-Marker (READY/RECEIVED/REJECTED) liest dieser
# Runner ueber `docker logs`, waehrend der Test noch laeuft. Die `change_id`
# der empfangenen Change wird zusaetzlich gegen den Lesezugriffsweg
# `cdc.changes` gehalten (gRPC/SSE/NATS), die Consumer-Registrierung gegen
# `cdc.consumer` (HTTP).
#
# Der Runner schreibt den Abdeckungs-Traeger docs/user/sdk-e2e-abdeckung.md
# (C#-Abschnitt, marker-gegrenzt) aus derselben Messung, die ihn belegt —
# idempotent, nur bei inhaltlicher Abweichung (Muster
# run-integration-tests.shs Abdeckungstabelle). Grenze der Writer-Form:
# dieser Runner erzeugt Kopf und alles bis zu seinem end-Marker neu und
# erhaelt den Rest; ein Folge-Runner (Kotlin/Python-HTTP) spiegelt dieselbe
# Form NICHT wortgleich — sein Re-Run wuerde die Abschnitte oberhalb seines
# begin-Markers wegwerfen. Er haelt den Kopf stabil und ersetzt nur seinen
# eigenen Abschnitt (Rest-Erhalt hinter dem end-Marker, wie hier).
#
# Voraussetzungen: Docker, ein geladenes :dev-Image (`make image` vorher —
# compose.yaml traegt keinen build:-Block, ADR-0044) und Netz (NuGet-Restore
# im SDK-Bau, d-migrate-Image). Kein Gate (dieselbe Klasse wie
# `make test-integration`). Cleanup je Ausgang: Test-Container wegwerfen,
# compose-Stack inkl. Volume abrueumen.
#
# Aufruf: `make test-sdk-csharp-integration`. Override:
# SDK_CSHARP_INTEGRATION_IMAGE.
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
API_TOKEN=e2e-reader-token
API_TOKEN_ADMIN=e2e-admin-token
API_TOKEN_READER=e2e-reader-token
NATS_STREAM_TOKEN=e2e-nats-stream-token
SOURCE_ID=src-e2e
HTTP_PUBLICATION=pub_pgc_e2e
DSN="postgres://$PG_USER:$PG_PASSWORD@$PG_CONTAINER:5432/$PG_DB?sslmode=disable"

SDK_INTEGRATION_IMAGE=${SDK_CSHARP_INTEGRATION_IMAGE:-pg-change-feed:sdk-csharp-integration}
SDK_TEST_CONTAINER=cdc-sdk-csharp-client-test

TEST_TABLE=feed_e2e_full
ABDECKUNG_ZIEL=docs/user/sdk-e2e-abdeckung
ABDECKUNG_ZIEL_DATEI=${ABDECKUNG_ZIEL}.md
GRPC_TEST_NAME=GrpcRealserverTests
SSE_TEST_NAME=SseRealserverTests
NATS_TEST_NAME=NatsRealserverTests
HTTP_TEST_NAME=HttpRealserverTests
GRPC_SENTINEL=CsharpGrpcSdkE2ESentinel
SSE_SENTINEL=CsharpSseSdkE2ESentinel
NATS_SENTINEL=CsharpNatsSdkE2ESentinel
HTTP_SENTINEL=CsharpHttpSdkE2ESentinel

cleanup() {
  docker rm -fv "$SDK_TEST_CONTAINER" >/dev/null 2>&1 || true
  $COMPOSE down -v --remove-orphans >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker image inspect ghcr.io/pt9912/pg-change-feed:dev >/dev/null 2>&1 || {
  echo "run-sdk-csharp-integration-tests: kein geladenes :dev-Image — make image vorher (compose.yaml traegt keinen build:-Block, ADR-0044)" >&2
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
  echo "run-sdk-csharp-integration-tests: Compose-PostgreSQL wurde nicht bereit" >&2
  exit 1
fi

make schema-rollout SCHEMA_TARGET="db:$DSN" SCHEMA_ROLLOUT_NETWORK="$NETWORK"

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 >/dev/null <<'SQL'
CREATE TABLE public.feed_e2e_flow (id int PRIMARY KEY, name text);
CREATE TABLE public.feed_e2e_full (id int PRIMARY KEY, name text);
ALTER TABLE public.feed_e2e_full REPLICA IDENTITY FULL;
CREATE TABLE public.feed_e2e_schema (id int PRIMARY KEY, name text, amount text);
INSERT INTO cdc.source (source_id, name) VALUES ('src-e2e', 'E2E-Quelle');
SQL

$COMPOSE up -d pg-change-feed >/dev/null

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
  echo "run-sdk-csharp-integration-tests: Feed-Container trägt die Verdrahtung nicht — Slot $SLOT: ${slot_ok:-0}, Feed-Lauf: $feed_running (Feed-Ausgang: $feed_exit)" >&2
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
  echo "run-sdk-csharp-integration-tests: Feed-Container meldet Compose-Health-Status ${health:-fehlt}, wollen healthy" >&2
  exit 1
fi

docker build --build-context proto=proto -f sdks/csharp/Dockerfile \
  --target integration -t "$SDK_INTEGRATION_IMAGE" sdks/csharp

# run_phase <Name> <Testklasse> <Sentinel> <ID-Basis> <Reject-Marker>
# <Received-Grep> <SQL-Variante: changes|consumer> <Extra-Env> — ein
# Realserver-Rundlauf fuer genau eine SDK-Flaeche: Container starten, auf
# READY warten, begrenzte Folge eindeutiger Zeilen committen
# (Fire-and-Forget-Fenster), auf den Reject-Marker und das Prozessende
# warten, die RECEIVED-Zeile pruefen und die Identitaet unabhaengig gegen
# den SQL-Lesezugriffsweg halten. Die SQL-Variante "changes" haelt die
# change_id gegen cdc.changes (Streams), "consumer" die Consumer-Registrierung
# gegen cdc.consumer (HTTP).
run_phase() {
  local phase_name=$1 test_name=$2 sentinel=$3 id_base=$4 reject_marker=$5 \
        received_grep=$6 sql_kind=$7 extra_env=$8
  local attempt insert_id captured ident pair

  local env_args=()
  for pair in $extra_env; do
    env_args+=(-e "$pair")
  done

  docker rm -fv "$SDK_TEST_CONTAINER" >/dev/null 2>&1 || true
  docker run -d --name "$SDK_TEST_CONTAINER" --network "$NETWORK" \
    "${env_args[@]}" \
    -e PGCHANGEFEED_TEST_NAME="$test_name" \
    -e PGCHANGEFEED_E2E_TABLE="$TEST_TABLE" \
    -e PGCHANGEFEED_E2E_SENTINEL="$sentinel" \
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
    echo "run-sdk-csharp-integration-tests: $phase_name — Test wurde nicht innerhalb der Zeitspanne bereit (kein READY): $(docker logs "$SDK_TEST_CONTAINER" 2>&1 || true)" >&2
    exit 1
  fi

  local received=0
  for attempt in 1 2 3 4 5; do
    insert_id=$((id_base + attempt))
    # >/dev/null: die psql-INSERT-Echo-Zeile gehoert nicht in den
    # Funktions-stdout (der Aufrufer haelt hier nur den Rueckgabewert).
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
    echo "run-sdk-csharp-integration-tests: $phase_name — der Test empfing keine der committeten Aenderungen ($TEST_TABLE, $sentinel): $test_output" >&2
    exit 1
  fi
  if ! printf '%s' "$test_output" | grep -qE "$received_grep"; then
    echo "run-sdk-csharp-integration-tests: $phase_name — die RECEIVED-Zeile traegt nicht die erwartete Form: $test_output" >&2
    exit 1
  fi
  if [ "$rejected" -ne 1 ]; then
    echo "run-sdk-csharp-integration-tests: $phase_name — der Ablehnungs-Beleg blieb aus ($reject_marker fehlt): $test_output" >&2
    exit 1
  fi
  if [ "$test_stopped" -ne 1 ] || [ "$test_exit" != "0" ]; then
    echo "run-sdk-csharp-integration-tests: $phase_name — der Test endete nicht mit Ausgang 0 (gestoppt: $test_stopped, Ausgang: $test_exit): $test_output" >&2
    exit 1
  fi

  ident=$(printf '%s' "$test_output" | grep -oE 'RECEIVED [a-z_]+=[^ ]+' | head -n1 | cut -d= -f2 || true)
  if [ -z "$ident" ]; then
    echo "run-sdk-csharp-integration-tests: $phase_name — die RECEIVED-Zeile traegt keine Identitaet: $test_output" >&2
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
    echo "run-sdk-csharp-integration-tests: $phase_name — die Identitaet ($ident) ist nicht real ueber den SQL-Lesezugriffsweg lesbar (count=${captured:-leer})" >&2
    exit 1
  fi

  feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" 2>/dev/null || echo false)
  if [ "$feed_running" != "true" ]; then
    echo "run-sdk-csharp-integration-tests: $phase_name — Feed-Container lief nach dem Lauf nicht mehr weiter (kein Neustart erwartet)" >&2
    exit 1
  fi

  printf '%s\n' "$ident"
}

GRPC_IDENT=$(run_phase \
  "gRPC-Flaeche (SPEC-020)" "$GRPC_TEST_NAME" \
  "$GRPC_SENTINEL" 400 "REJECTED code=Unauthenticated" \
  "RECEIVED change_id=[^ ]+ table=$TEST_TABLE .*operation=INSERT .*new_image=.*$GRPC_SENTINEL" \
  changes "PGCHANGEFEED_GRPC_ADDR=pg-change-feed:9090 PGCHANGEFEED_API_TOKEN=$API_TOKEN")

SSE_IDENT=$(run_phase \
  "SSE-Flaeche (SPEC-021)" "$SSE_TEST_NAME" \
  "$SSE_SENTINEL" 410 "REJECTED status=401" \
  "RECEIVED change_id=[^ ]+ table=$TEST_TABLE .*operation=INSERT .*new_image=.*$SSE_SENTINEL" \
  changes "PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090 PGCHANGEFEED_API_TOKEN=$API_TOKEN")

NATS_IDENT=$(run_phase \
  "NATS-Vollinhalts-Flaeche (SPEC-024)" "$NATS_TEST_NAME" \
  "$NATS_SENTINEL" 420 "REJECTED token-rejected" \
  "RECEIVED change_id=[^ ]+ table=$TEST_TABLE .*operation=INSERT .*new_image=.*$NATS_SENTINEL" \
  changes "PGCHANGEFEED_NATS_URL=nats://nats:4222 PGCHANGEFEED_NATS_STREAM_TOKEN=$NATS_STREAM_TOKEN PGCHANGEFEED_SOURCE_ID=$SOURCE_ID")

HTTP_IDENT=$(run_phase \
  "HTTP-Flaeche (SPEC-018)" "$HTTP_TEST_NAME" \
  "$HTTP_SENTINEL" 430 "REJECTED status=401" \
  "RECEIVED consumer_id=[^ ]+" \
  consumer "PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090 PGCHANGEFEED_API_TOKEN_ADMIN=$API_TOKEN_ADMIN PGCHANGEFEED_API_TOKEN_READER=$API_TOKEN_READER PGCHANGEFEED_SOURCE_ID=$SOURCE_ID PGCHANGEFEED_HTTP_PUBLICATION=$HTTP_PUBLICATION")

# --- Abdeckungs-Traeger (docs/user/sdk-e2e-abdeckung.md) -------------------
# Der C#-Abschnitt entsteht aus derselben Messung, die ihn belegt: je Phase
# ihre Kennungen und der ident-Behalt des Laufs; geschrieben wird nur bei
# inhaltlicher Abweichung (Temp-Datei + cmp), marker-gegrenzt, damit die
# Folge-Slices (Kotlin, Python-HTTP) eigene Abschnitte daneben tragen.
# Die Kennungsspalte traegt je Kennung den Link auf ihr Definitionsdokument
# (ids des Doku-Gates).
abdeckung_csharp_abschnitt() {
  printf '%s\n' \
    '<!-- pgchangefeed-sdk-e2e:csharp-begin -->' \
    "| [\`LH-FA-SST-008\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein C#-SDK-Client (\`PgChangeFeedGrpcClient\`) oeffnet real den gRPC-Server-Stream gegen den laufenden Feed-Container und empfaengt eine danach committete Aenderung; ein Oeffnungsversuch ohne gueltiges Token endet mit gRPC-Status \`Unauthenticated\` | \`GrpcRealserverTests\` | \`tools/harness/run-sdk-csharp-integration-tests.sh\` |" \
    "| [\`LH-FA-SST-008\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein C#-SDK-Client (\`PgChangeFeedSseClient\`) oeffnet real \`GET /changes/stream\` und empfaengt eine danach committete Aenderung; ein Aufruf ohne gueltiges Token endet mit HTTP-Status 401 | \`SseRealserverTests\` | \`tools/harness/run-sdk-csharp-integration-tests.sh\` |" \
    "| [\`LH-FA-SST-008\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein C#-SDK-Client (\`PgChangeFeedNatsStreamClient\`) verbindet sich real per NATS und empfaengt eine danach committete Aenderung als vollstaendiges JSON-Event; ein Verbindungsversuch mit falschem Token wird vom NATS-Server abgelehnt | \`NatsRealserverTests\` | \`tools/harness/run-sdk-csharp-integration-tests.sh\` |" \
    "| [\`LH-FA-SST-006\`](../../spec/lastenheft.md), [\`LH-FA-CON-001\`](../../spec/lastenheft.md), [\`LH-FA-SST-009\`](../../spec/lastenheft.md) | ein C#-SDK-Client (\`PgChangeFeedHttpClient\`) registriert real einen Consumer (admin-Token) und listet Tabellen (reader-Token); die Registrierung ist unabhaengig ueber \`cdc.consumer\` lesbar; ein Aufruf ohne gueltiges Token endet mit HTTP-Status 401 | \`HttpRealserverTests\` | \`tools/harness/run-sdk-csharp-integration-tests.sh\` |" \
    '<!-- pgchangefeed-sdk-e2e:csharp-end -->'
}

abdeckung_schreiben() {
  local temp fremde_abschnitte=""
  local abschnitt
  abschnitt=$(abdeckung_csharp_abschnitt)
  # Fremde, marker-gegrenzte Abschnitte (Kotlin/Python-HTTP der Folge-Slices)
  # hinter dem eigenen end-Marker bleiben erhalten — jeder Runner ersetzt nur
  # seinen eigenen Abschnitt.
  local rest=""
  if [ -f "$ABDECKUNG_ZIEL_DATEI" ]; then
    rest=$(awk '/pgchangefeed-sdk-e2e:csharp-end/{gefunden=1; next} gefunden{print}' "$ABDECKUNG_ZIEL_DATEI")
  fi
  temp=$(mktemp)
  {
    printf '%s\n' \
      '# SDK-E2E-Abdeckung je Spec-Kennung' \
      '' \
      'Erzeugt von `make test-sdk-csharp-integration` über' \
      '`tools/harness/run-sdk-csharp-integration-tests.sh`; die Sprach-Runner' \
      'der Folge-Slices (Kotlin, Python-HTTP) erweitern dieselbe Datei um ihre' \
      'marker-gegrenzten Abschnitte. Je Sprach-Abschnitt deklariert der' \
      'zuständige Runner seine Realserver-Phasen an Ort und Stelle. Diese' \
      'Datei ist eine **stabile Abdeckungs-Deklaration**, kein Lauf-Beleg: der' \
      'Runner schreibt sie nur bei inhaltlicher Abweichung. Sie trägt nur' \
      'Zeilen real existierender Runner-Phasen — ein Beleg steht hier nie,' \
      'bevor sein Lauf grün lief.' \
      '' \
      '| Spec-Kennung | Kurzbeschreibung | Nachweis | Ort |' \
      '| --- | --- | --- | --- |'
    printf '%s\n' "$abschnitt"
    if [ -n "$rest" ]; then
      printf '%s\n' "$rest"
    fi
  } > "$temp"
  if [ -f "$ABDECKUNG_ZIEL_DATEI" ] && cmp -s "$temp" "$ABDECKUNG_ZIEL_DATEI"; then
    rm -f "$temp"
    echo "run-sdk-csharp-integration-tests: Abdeckungs-Traeger unveraendert — $ABDECKUNG_ZIEL_DATEI entspricht dem Quelltext-Stand"
  else
    chmod 0644 "$temp"
    mv "$temp" "$ABDECKUNG_ZIEL_DATEI"
    echo "run-sdk-csharp-integration-tests: Abdeckungs-Traeger geschrieben — $ABDECKUNG_ZIEL_DATEI"
  fi
}

abdeckung_schreiben

echo "run-sdk-csharp-integration-tests: C#-SDK-Realserver-Belege (slice-sdk-csharp-reale2e, Mechanik-Klasse ADR-0110 Festlegung 2) gruen — gRPC-Flaeche (PgChangeFeedGrpcClient, pg-change-feed:9090, change_id=$GRPC_IDENT), SSE-Flaeche (PgChangeFeedSseClient, pg-change-feed:8090, change_id=$SSE_IDENT) und NATS-Vollinhalts-Flaeche (PgChangeFeedNatsStreamClient, nats://nats:4222, change_id=$NATS_IDENT) empfingen je eine danach committete Aenderung (change_id je unabhaengig ueber cdc.changes lesbar), die HTTP-Flaeche (PgChangeFeedHttpClient) registrierte real einen Consumer (consumer_id=$HTTP_IDENT, unabhaengig ueber cdc.consumer lesbar) und listete Tabellen; ein Aufruf ohne gueltiges Token endete je mit gRPC-Status Unauthenticated, HTTP-Status 401 bzw. der laut ablehnenden NATS-Verbindungsablehnung"