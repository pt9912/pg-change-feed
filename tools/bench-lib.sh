#!/usr/bin/env bash
# bench-lib.sh — gemeinsame Umgebungs-Bausteine für tools/bench-*.sh
# (LH-QA-PER-001…003, LH-FA-CAP-009, ADR-0054 §(b)): eigene, von
# tools/harness/run-integration-tests.sh unabhängige PostgreSQL-/
# Feed-Umgebung über `docker network`/`docker run` (kein compose.yaml-Bezug,
# damit ein Bench-Lauf nicht mit einem laufenden `make test-integration`
# um dieselben Container-/Netznamen konkurriert), Schema-Rollout über
# d-migrate (ADR-0043) und Cleanup. Geteilt ist die Umgebung, nicht die
# Messung: Jedes der vier Bench-Skripte bleibt ein eigenständiger Beleg
# (ADR-0054 §(b) — "je Beleg ein eigenes Bench-Skript"), quellt diese
# Datei nur für Aufbau/Abbau.
set -euo pipefail

# Dieselben Digest-Pins wie Makefile/compose.yaml (PG_TEST_IMAGE,
# D_MIGRATE_IMAGE) — Pin-Hebung bleibt ein bewusster Commit (Modul 14).
PG_TEST_IMAGE=${PG_TEST_IMAGE:-postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8}
FEED_IMAGE=${FEED_IMAGE:-ghcr.io/pt9912/pg-change-feed:dev}
D_MIGRATE_IMAGE=${D_MIGRATE_IMAGE:-ghcr.io/pt9912/d-migrate@sha256:862dfb04c34dd17278b1bab46961363c12eeb8d464cf1776565d6285603d2c89}

bench::repo_root() { git rev-parse --show-toplevel; }

bench::network_name()  { echo "${BENCH_NETWORK:-pgc-bench-net}"; }
bench::pg_container()  { echo "${BENCH_PG_CONTAINER:-pgc-bench-postgres}"; }
bench::feed_container() { echo "${BENCH_FEED_CONTAINER:-pgc-bench-feed}"; }

# Räumt eine vorherige (z. B. abgebrochene) Umgebung weg und dient als
# EXIT-Trap — idempotent, jeder Schritt toleriert "existiert nicht".
# `-v` entfernt die anonymen Volumes des Containers (das Postgres-Image
# deklariert ein VOLUME).
bench::cleanup() {
  local feed pg net
  feed=$(bench::feed_container); pg=$(bench::pg_container); net=$(bench::network_name)
  docker rm -fv "$feed" >/dev/null 2>&1 || true
  docker rm -fv "$pg" >/dev/null 2>&1 || true
  docker network rm "$net" >/dev/null 2>&1 || true
}

bench::start_postgres() {
  local net pg root
  net=$(bench::network_name); pg=$(bench::pg_container); root=$(bench::repo_root)
  bench::cleanup
  docker network create "$net" >/dev/null
  # wal_level=logical wie compose.yaml — Vorbedingung für den
  # Replication-Slot, den der Feed-Container beim Start anlegt.
  # compose-init mounten wie compose.yaml: legt das leere Schema `cdc`
  # an und setzt den search_path der `postgres`-Rolle in dieser DB — ohne
  # das landet die unqualifizierte Tabellen-DDL des d-migrate-Rollouts in
  # `public`, während die generierten Views explizit `cdc.<table>`
  # referenzieren ("relation cdc.source_table does not exist", real
  # reproduziert).
  docker run -d --name "$pg" --network "$net" \
    -e POSTGRES_DB=cdc -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres \
    -v "$root/tools/schema/compose-init":/docker-entrypoint-initdb.d:ro \
    "$PG_TEST_IMAGE" \
    postgres -c wal_level=logical -c max_wal_senders=10 \
      -c max_replication_slots=10 -c wal_sender_timeout=2000 >/dev/null

  local ready=0
  for _ in $(seq 1 60); do
    if docker exec "$pg" pg_isready -U postgres -d cdc >/dev/null 2>&1; then
      ready=1
      break
    fi
    sleep 1
  done
  if [ "$ready" -ne 1 ]; then
    echo "bench-lib: PostgreSQL wurde nicht bereit" >&2
    exit 1
  fi
}

bench::schema_rollout() {
  local net pg dsn
  net=$(bench::network_name); pg=$(bench::pg_container)
  dsn="postgres://postgres:postgres@$pg:5432/cdc?sslmode=disable"
  make -C "$(bench::repo_root)" schema-rollout \
    SCHEMA_TARGET="db:$dsn" SCHEMA_ROLLOUT_NETWORK="$net" >/dev/null
}

bench::psql() {
  local pg
  pg=$(bench::pg_container)
  docker exec -i "$pg" psql -U postgres -d cdc -v ON_ERROR_STOP=1 "$@"
}

bench::psql_scalar() {
  local pg
  pg=$(bench::pg_container)
  docker exec "$pg" psql -U postgres -d cdc -tAc "$1"
}

# $1=CDC_TABLES  $2=CDC_SOURCE_ID  $3=CDC_SLOT  $4=CDC_PUBLICATION
bench::start_feed() {
  local net feed pg dsn extra=() entry
  net=$(bench::network_name); feed=$(bench::feed_container); pg=$(bench::pg_container)
  dsn="postgres://postgres:postgres@$pg:5432/cdc?sslmode=disable"
  # BENCH_FEED_ENV: zusätzliche Umgebungsvariablen des Feed-Containers
  # (leerzeichengetrennt, KEY=VALUE); BENCH_FEED_DOCKER_ARGS: zusätzliche
  # Argumente von `docker run` (etwa `--memory 512m`).
  for entry in ${BENCH_FEED_ENV:-}; do
    extra+=(-e "$entry")
  done
  docker rm -fv "$feed" >/dev/null 2>&1 || true
  docker run -d --name "$feed" --network "$net" \
    -e CDC_CAPTURE_DSN="$dsn" -e CDC_ADMIN_DSN="$dsn" -e CDC_READER_DSN="$dsn" \
    -e CDC_SOURCE_ID="$2" -e CDC_PUBLICATION="$4" -e CDC_SLOT="$3" \
    -e CDC_TABLES="$1" -e CDC_LOG_LEVEL=info \
    "${extra[@]}" ${BENCH_FEED_DOCKER_ARGS:-} \
    "$FEED_IMAGE" >/dev/null

  local wired=0 slot_ok feed_running
  for _ in $(seq 1 60); do
    slot_ok=$(bench::psql_scalar "SELECT 1 FROM pg_replication_slots WHERE slot_name = '$3'" 2>/dev/null || true)
    feed_running=$(docker inspect --format '{{.State.Running}}' "$feed" 2>/dev/null || echo false)
    if [ "$slot_ok" = "1" ] && [ "$feed_running" = "true" ]; then
      wired=1
      break
    fi
    sleep 1
  done
  if [ "$wired" -ne 1 ]; then
    echo "bench-lib: Feed-Container trägt die Verdrahtung nicht (Slot $3)" >&2
    docker logs "$feed" 2>&1 | tail -40 >&2 || true
    exit 1
  fi
}

bench::stop_feed() {
  local feed
  feed=$(bench::feed_container)
  docker rm -fv "$feed" >/dev/null 2>&1 || true
}

# Wartet, bis mindestens $2 Zeilen für Quelle/Tabelle in cdc.changes stehen
# (Replikations-Verzögerung zwischen Schreiben und Sichtbarkeit).
# $1=source_id $2=table_name $3=erwartete Mindestzahl $4=Timeout in s (Default 120)
bench::wait_captured() {
  local source=$1 table=$2 want=$3 timeout=${4:-120} count
  for _ in $(seq 1 "$timeout"); do
    count=$(bench::psql_scalar "SELECT count(*) FROM cdc.changes WHERE source_id = '$source' AND table_name = '$table'")
    if [ "${count:-0}" -ge "$want" ]; then
      return 0
    fi
    sleep 1
  done
  echo "bench-lib: nur ${count:-0}/$want Zeilen für $source/$table nach ${timeout}s erfasst" >&2
  return 1
}

# bench::median_of liefert den mittleren Wert einer ungeraden Anzahl
# übergebener Zahlen — robuster gegen einen einzelnen Ausreißer als eine
# Einzelmessung.
bench::median_of() {
  local count=$# mid
  mid=$(( (count + 1) / 2 ))
  printf '%s\n' "$@" | sort -n | sed -n "${mid}p"
}

# bench::record_row hält eine Kennungs-Zeile für docs/user/bench-abdeckung.md
# als eigene Datei unter .tmp/bench-abdeckung-rows/ fest (ADR-0104) — ein
# Schwellen-Skript schreibt nur seine eigene Zeile, ohne die der anderen
# Skripte zu lesen oder zu sperren. bench::render_abdeckung liest am Ende
# des letzten Schwellen-Skripts alle vorhandenen Zeilen und schreibt die
# Tabelle gesammelt.
# $1=Lastenheft-Kennung $2=Kurzbeschreibung $3=Schwelle-Text $4=Ort (Datei:Zeile)
bench::record_row() {
  local id=$1 kurzbeschreibung=$2 schwelle=$3 ort=$4 dir
  dir="$(bench::repo_root)/.tmp/bench-abdeckung-rows"
  mkdir -p "$dir"
  printf '| [`%s`](../../spec/lastenheft.md) | %s | %s | `%s` |\n' \
    "$id" "$kurzbeschreibung" "$schwelle" "$ort" > "$dir/$id.row"
}

# bench::render_abdeckung schreibt docs/user/bench-abdeckung.md aus allen
# bislang unter .tmp/bench-abdeckung-rows/ abgelegten Zeilen (sortiert nach
# Kennung, damit der Diff stabil bleibt) — aufgerufen vom letzten der drei
# Schwellen-Skripte in Ausführungsreihenfolge (tools/bench-batch-vs-single.sh);
# tools/bench-backfill.sh misst ohne Schwelle und schreibt keine Zeile.
bench::render_abdeckung() {
  local root file dir
  root=$(bench::repo_root)
  dir="$root/.tmp/bench-abdeckung-rows"
  file="$root/docs/user/bench-abdeckung.md"
  {
    cat <<'HEADER'
# Bench-Abdeckung je Lastenheft-Kennung

Erzeugt von den drei `tools/bench-*.sh`-Skripten mit Pass/Fail-Schwelle
(`make bench`,
[`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md));
`tools/bench-backfill.sh` misst ohne Schwelle und trägt keine Zeile.
Jede Zeile bindet eine Kennung an ihre real durchgesetzte Pass/Fail-
Schwelle. Diese Datei ist eine **stabile Abdeckungs-Deklaration**, kein
Lauf-Beleg — der zuletzt gemessene Wert steht in der stdout-Ausgabe des
jeweiligen Laufs, nicht hier.

| Lastenheft-Kennung | Kurzbeschreibung | Schwelle (SPEC) | Ort |
| --- | --- | --- | --- |
HEADER
    if [ -d "$dir" ]; then
      find "$dir" -name '*.row' -print0 | sort -z | xargs -0 cat 2>/dev/null || true
    fi
  } > "$file"
}
