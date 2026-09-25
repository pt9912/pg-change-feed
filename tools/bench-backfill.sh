#!/usr/bin/env bash
# bench-backfill.sh — Messung des Backfills (LH-FA-CAP-009, ADR-0111, ADR-0113):
# Kopierdauer je Tabellengröße, Abweichung der geschätzten Zeilenzahl,
# Speicher des Feed-Containers und Wirkung eines Runs auf die Live-Erfassung.
# Eine Messung mit gedruckter Zahl, kein Pass/Fail gegen eine Schwelle.
# Vertrag, Ablauf und Grenzen: harness/targets/bench-backfill.md.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
# shellcheck source=tools/bench-lib.sh
source tools/bench-lib.sh

FULL=0
if [ "${1:-}" = "--full" ]; then
  FULL=1
fi

if [ "$FULL" -eq 1 ]; then
  STAGES=${BENCH_BACKFILL_STAGES:-"10000 100000 1000000"}
else
  STAGES=${BENCH_BACKFILL_STAGES:-"10000 50000 200000"}
fi
RUNS=${BENCH_BACKFILL_RUNS:-3}
EST_ROWS=${BENCH_BACKFILL_EST_ROWS:-100000}
LIVE_RATE=${BENCH_BACKFILL_LIVE_RATE:-100}
RUN_TIMEOUT_S=${BENCH_BACKFILL_RUN_TIMEOUT_S:-1800} # Sekunden je Run (Uhr: bash SECONDS)
# Warnschwelle des WAL-Rückstands des Capture-Slots (SPEC-013: 100 MiB), der
# Vergleichswert der gedruckten Rückstände.
WAL_WARN_BYTES=104857600

SOURCE_ID=src-bench-backfill
SLOT=slot_bench_backfill
PUBLICATION=pub_bench_backfill
LIVE_TABLE=bench_backfill_live
RUN_ID=$(date -u +%Y%m%dT%H%M%SZ)
TAG="bench-backfill[$RUN_ID]"

# Toleranz und Blockgröße stammen aus dem Code, nicht aus einer zweiten Kopie.
TOLERANCE_MIN=$(grep -oE 'copyDurationToleranceMinutes += [0-9]+' internal/application/usecase/backfill/warn.go | grep -oE '[0-9]+$' || true)
BLOCK_SIZE=$(grep -oE 'DefaultBlockSize += [0-9]+' internal/adapters/driven/postgressnapshot/snapshot.go | grep -oE '[0-9]+$' || true)
if [ -z "$TOLERANCE_MIN" ] || [ -z "$BLOCK_SIZE" ]; then
  echo "$TAG: Toleranz (warn.go) oder Blockgröße (snapshot.go) nicht lesbar" >&2
  exit 1
fi

WRITER_PID=""
STOP_FILE=$(mktemp -u)
COUNT_FILE=$(mktemp -u)
cleanup_all() {
  touch "$STOP_FILE"
  if [ -n "$WRITER_PID" ]; then
    wait "$WRITER_PID" 2>/dev/null || true
  fi
  rm -f "$STOP_FILE" "$COUNT_FILE"
  bench::cleanup
}
trap cleanup_all EXIT

table_of() { echo "bench_backfill_s$1"; }

fill_table() {
  local table=$1 rows=$2 options=${3:-}
  bench::psql >/dev/null <<SQL
CREATE TABLE public.$table (id bigint PRIMARY KEY, name text NOT NULL, amount numeric(12,2), created_at timestamptz, note text) $options;
INSERT INTO public.$table
SELECT g, 'kunde-' || g, (g % 100000) / 100.0, now(), CASE WHEN g % 5 = 0 THEN NULL ELSE 'notiz-' || g END
FROM generate_series(1, $rows) AS g;
SQL
}

row_width() {
  bench::psql_scalar "SELECT round(avg(pg_column_size(t.*)))::int FROM public.$1 t"
}

reltuples_of() {
  bench::psql_scalar "SELECT c.reltuples::float8::text FROM pg_class c WHERE c.oid = to_regclass('public.$1')"
}

# Ein Katalogwert unter 0 heißt unbekannt (dieselbe Abbildung wie der Snapshot-Leser).
estimate_line() {
  local state=$1 table=$2 actual=$3 value
  value=$(reltuples_of "$table")
  LC_ALL=C awk -v tag="$TAG" -v state="$state" -v v="$value" -v n="$actual" 'BEGIN {
    if (v < 0) { printf "%s: Schätzung — %s (%d Zeilen tatsächlich): reltuples=%s, geschätzt unbekannt\n", tag, state, n, v }
    else { printf "%s: Schätzung — %s (%d Zeilen tatsächlich): reltuples=%s, geschätzt %d, Abweichung %+.1f%%\n", tag, state, n, v, v + 0.5, (v - n) * 100.0 / n }
  }'
}

# Speicher des Feed-Containers in MiB (docker stats, ein Aufruf).
feed_mem_mib() {
  docker stats --no-stream --format '{{.MemUsage}}' "$(bench::feed_container)" 2>/dev/null | LC_ALL=C awk '{
    v = $1; unit = v; gsub(/[0-9.]/, "", unit); sub(/[A-Za-z]+$/, "", v)
    if (unit == "GiB") printf "%.1f", v * 1024
    else if (unit == "MiB") printf "%.1f", v
    else if (unit == "KiB") printf "%.1f", v / 1024
    else printf "0.0"
  }'
}

store_rows() { bench::psql_scalar "SELECT count(*) FROM cdc.change"; }

max_of() { LC_ALL=C awk -v a="$1" -v b="$2" 'BEGIN { print (b + 0 > a + 0) ? b : a }'; }

# min/Median/max einer Zahlenliste (Argumente); der Median ist bench::median_of.
range_of() {
  local sorted
  sorted=$(printf '%s\n' "$@" | sort -n)
  printf '%s–%s (Median %s, n=%d)' "$(head -n 1 <<< "$sorted")" "$(tail -n 1 <<< "$sorted")" "$(bench::median_of "$@")" "$#"
}

# WAL-Rückstand des Capture-Slots in Bytes: aktuelle WAL-Position minus
# confirmed_flush_lsn, dieselbe Größe wie cdc_wal_retention_bytes (SPEC-009).
wal_backlog_bytes() {
  bench::psql_scalar "SELECT COALESCE((SELECT pg_wal_lsn_diff(pg_current_wal_lsn(), confirmed_flush_lsn)::bigint FROM pg_replication_slots WHERE slot_name = '$SLOT'), 0)"
}

# Vom Slot gehaltenes WAL in Bytes: aktuelle WAL-Position minus restart_lsn
# (das WAL, das die Quelle für den Slot auf der Platte behält).
wal_held_bytes() {
  bench::psql_scalar "SELECT COALESCE((SELECT pg_wal_lsn_diff(pg_current_wal_lsn(), restart_lsn)::bigint FROM pg_replication_slots WHERE slot_name = '$SLOT'), 0)"
}

# Wartet ohne jeden Schreibzugriff höchstens 120 s auf einen Rückstand unter
# der Warnschwelle (die Leerlauf-Bestätigung des Streams senkt ihn, ADR-0120);
# druckt "Rückstand_Bytes|Wartezeit_s".
settle_wal() {
  local waited=0 wal
  while :; do
    wal=$(wal_backlog_bytes)
    if [ "$wal" -lt "$WAL_WARN_BYTES" ] || [ "$waited" -ge 120 ]; then
      break
    fi
    waited=$((waited + 1))
    sleep 1
  done
  echo "$wal|$waited"
}

mib_of() { LC_ALL=C awk -v b="$1" 'BEGIN { printf "%.0f", b / 1048576 }'; }

# run_backfill <tabelle> — beantragt einen Run, wartet auf sein Ende und
# druckt eine Zeile
# "ms|rows_copied|estimated|warn_size|warn_duration|peak_mib|wal_peak_bytes|held_peak_bytes".
run_backfill() {
  local table=$1 id status peak=0 mem deadline row wal wal_peak=0 held held_peak=0 sample rest
  id=$(bench::psql_scalar "SELECT cdc.backfill_table('$SOURCE_ID', 'public', '$table')")
  if [ -z "$id" ]; then
    echo "$TAG: cdc.backfill_table lieferte keine Antragskennung für $table" >&2
    return 1
  fi
  deadline=$((SECONDS + RUN_TIMEOUT_S))
  while :; do
    sample=$(bench::psql_scalar "SELECT COALESCE((SELECT status FROM cdc.backfill_run WHERE run_id = '$id'), '') || '|' || COALESCE((SELECT pg_wal_lsn_diff(pg_current_wal_lsn(), confirmed_flush_lsn)::bigint || '|' || pg_wal_lsn_diff(pg_current_wal_lsn(), restart_lsn)::bigint FROM pg_replication_slots WHERE slot_name = '$SLOT'), '0|0')")
    status=${sample%%|*}
    rest=${sample#*|}
    wal=${rest%%|*}
    held=${rest#*|}
    wal_peak=$(max_of "$wal_peak" "${wal:-0}")
    held_peak=$(max_of "$held_peak" "${held:-0}")
    mem=$(feed_mem_mib)
    peak=$(max_of "$peak" "${mem:-0}")
    case "$status" in
      completed) break ;;
      failed | interrupted)
        echo "$TAG: Run $id von $table endete $status: $(bench::psql_scalar "SELECT error_message FROM cdc.backfill_run WHERE run_id = '$id'")" >&2
        echo "$TAG: Feed-Container: $(docker inspect --format 'Status {{.State.Status}}, Exit {{.State.ExitCode}}, OOMKilled {{.State.OOMKilled}}' "$(bench::feed_container)" 2>&1)" >&2
        docker logs --tail 30 "$(bench::feed_container)" >&2 2>&1 || true
        return 1
        ;;
    esac
    if [ "$SECONDS" -gt "$deadline" ]; then
      echo "$TAG: Run $id von $table nicht beendet nach ${RUN_TIMEOUT_S} s" >&2
      return 1
    fi
  done
  row=$(bench::psql_scalar "SELECT round(EXTRACT(EPOCH FROM (finished_at - started_at)) * 1000)::bigint, rows_copied, COALESCE(estimated_rows::text, 'unbekannt'), warn_estimated_size, warn_duration FROM cdc.backfill_run WHERE run_id = '$id'")
  echo "$row|$peak|$wal_peak|$held_peak"
}

# Live-Schreiber: LIVE_RATE Zeilen je Sekunde in die Live-Tabelle ab der Kennung
# $1, bis die Stop-Datei existiert.
live_writer() {
  local base=$1 i=0 start end ms sleep_s
  while [ ! -e "$STOP_FILE" ]; do
    start=$(date +%s%N)
    bench::psql -c "INSERT INTO public.$LIVE_TABLE (id, name) SELECT $base + $i * $LIVE_RATE + g, 'live' FROM generate_series(1, $LIVE_RATE) AS g;" >/dev/null
    i=$((i + 1))
    echo $((i * LIVE_RATE)) > "$COUNT_FILE"
    end=$(date +%s%N)
    ms=$(( (end - start) / 1000000 ))
    sleep_s=$(LC_ALL=C awk -v ms="$ms" 'BEGIN { s = (1000 - ms) / 1000.0; if (s < 0) s = 0; printf "%.3f", s }')
    sleep "$sleep_s"
  done
}

lag_sample() {
  bench::psql_scalar "SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_capture_lag'" 2>/dev/null || true
}

start_writer() {
  rm -f "$STOP_FILE" "$COUNT_FILE"
  live_writer "$1" &
  WRITER_PID=$!
}

stop_writer() {
  touch "$STOP_FILE"
  wait "$WRITER_PID" 2>/dev/null || true
  WRITER_PID=""
  rm -f "$STOP_FILE"
}

echo "$TAG: Host — $(uname -sr), Docker $(docker info --format '{{.ServerVersion}}, {{.NCPU}} CPU, {{.MemTotal}} Byte RAM'), PostgreSQL-Image ${PG_TEST_IMAGE%%@*}"
echo "$TAG: Umgebung wird aufgebaut (Stufen: $STAGES Zeilen, $RUNS Läufe je Stufe, Blockgröße B=$BLOCK_SIZE Zeilen, Toleranz ${TOLERANCE_MIN} min) …"
bench::start_postgres
bench::schema_rollout
bench::psql -c "INSERT INTO cdc.source (source_id, name) VALUES ('$SOURCE_ID', 'Bench-Backfill');" >/dev/null

# --- Schätzung der Zeilenzahl (pg_class.reltuples) ---------------------------
echo "$TAG: Schätz-Messung — Tabellen mit $EST_ROWS Zeilen …"
fill_table bench_backfill_est_auto "$EST_ROWS"
fill_table bench_backfill_est_never "$EST_ROWS" "WITH (autovacuum_enabled = false)"
estimate_line "frisch befüllt, Autovacuum an" bench_backfill_est_auto "$EST_ROWS"
estimate_line "frisch befüllt, Autovacuum aus (nie analysiert)" bench_backfill_est_never "$EST_ROWS"
waited=0
until [ "$(LC_ALL=C awk -v v="$(reltuples_of bench_backfill_est_auto)" 'BEGIN { print (v >= 0) ? 1 : 0 }')" = "1" ] || [ "$waited" -ge 150 ]; do
  sleep 5
  waited=$((waited + 5))
done
estimate_line "Autovacuum an, ${waited} s nach dem Befüllen" bench_backfill_est_auto "$EST_ROWS"
estimate_line "Autovacuum aus, ${waited} s nach dem Befüllen (nie analysiert)" bench_backfill_est_never "$EST_ROWS"
bench::psql -c "ANALYZE public.bench_backfill_est_auto;" >/dev/null
estimate_line "nach ANALYZE" bench_backfill_est_auto "$EST_ROWS"
extra=$((EST_ROWS / 5))
bench::psql >/dev/null <<SQL
ALTER TABLE public.bench_backfill_est_auto SET (autovacuum_enabled = false);
INSERT INTO public.bench_backfill_est_auto SELECT g, 'kunde-' || g, 0, now(), NULL FROM generate_series($EST_ROWS + 1, $EST_ROWS + $extra) AS g;
SQL
estimate_line "nach ANALYZE und $extra weiteren Zeilen (+20 %) ohne erneutes ANALYZE" bench_backfill_est_auto $((EST_ROWS + extra))

# --- Stufen anlegen, Feed starten -----------------------------------------
tables_env="public.$LIVE_TABLE=tbl-bf-live:sv-bf-live"
bench::psql -c "CREATE TABLE public.$LIVE_TABLE (id bigint PRIMARY KEY, name text);" >/dev/null
for n in $STAGES; do
  fill_table "$(table_of "$n")" "$n" "WITH (autovacuum_enabled = false)"
  tables_env="$tables_env,public.$(table_of "$n")=tbl-bf-$n:sv-bf-$n"
done
bench::start_feed "$tables_env" "$SOURCE_ID" "$SLOT" "$PUBLICATION"
sleep 30
echo "$TAG: Grundlinie ohne Run — Feed-Container in Ruhe 30 s nach dem Start $(feed_mem_mib) MiB, Zeilen in cdc.change: $(store_rows)"

# --- Kopierdauer je Tabellengröße --------------------------------------------
last_rate=""
last_stage=""
for n in $STAGES; do
  table=$(table_of "$n")
  width=$(row_width "$table")
  idle=$(feed_mem_mib)
  echo "$TAG: Stufe $n Zeilen (Zeilenbreite ~${width} B gemittelt über pg_column_size, B=$BLOCK_SIZE, Einfügeform zeilenweise in einer Transaktion), Feed-Container in Ruhe ${idle} MiB, Zeilen in cdc.change vor der Stufe: $(store_rows)"
  ms_list=(); rate_list=(); wal_list=(); held_list=(); peak_stage=0
  for run in $(seq 1 "$RUNS"); do
    result=$(run_backfill "$table")
    IFS='|' read -r ms copied estimated warn_size warn_duration peak wal_peak held_peak <<< "$result"
    if [ "$copied" -ne "$n" ]; then
      echo "$TAG: Run $run von Stufe $n kopierte $copied statt $n Zeilen — Messung ungültig" >&2
      exit 1
    fi
    wal_end=$(wal_backlog_bytes)
    IFS='|' read -r wal_settled wal_wait <<< "$(settle_wal)"
    rate=$(LC_ALL=C awk -v n="$n" -v ms="$ms" 'BEGIN { if (ms <= 0) ms = 1; printf "%.0f", n * 1000.0 / ms }')
    ms_list+=("$ms"); rate_list+=("$rate"); wal_list+=("$wal_peak"); held_list+=("$held_peak"); peak_stage=$(max_of "$peak_stage" "$peak")
    echo "$TAG: Stufe $n, Lauf $run/$RUNS — Kopierdauer ${ms} ms (finished_at − started_at), ${rate} Zeilen/s, geschätzt $estimated, Warnung Größe $warn_size, Warnung Dauer $warn_duration, Feed-Speicher-Spitze ${peak} MiB, WAL-Rückstand des Slots (confirmed_flush_lsn): Spitze $(mib_of "$wal_peak") MiB im Run, $(mib_of "$wal_end") MiB unmittelbar danach, $(mib_of "$wal_settled") MiB ${wal_wait} s danach ohne Schreibzugriff; vom Slot gehaltenes WAL (restart_lsn): Spitze $(mib_of "$held_peak") MiB im Run"
  done
  sleep 20
  after=$(feed_mem_mib)
  sleep 40
  after60=$(feed_mem_mib)
  blocks=$(( (n + BLOCK_SIZE - 1) / BLOCK_SIZE ))
  median_ms=$(bench::median_of "${ms_list[@]}")
  median_rate=$(bench::median_of "${rate_list[@]}")
  median_wal=$(bench::median_of "${wal_list[@]}")
  median_held=$(bench::median_of "${held_list[@]}")
  block_ms=$(LC_ALL=C awk -v ms="$median_ms" -v b="$blocks" 'BEGIN { printf "%.0f", ms / b }')
  wal_per_row=$(LC_ALL=C awk -v w="$median_wal" -v n="$n" 'BEGIN { printf "%.0f", w / n }')
  echo "$TAG: Stufe $n Ergebnis — Kopierdauer $(range_of "${ms_list[@]}") ms, Durchsatz $(range_of "${rate_list[@]}") Zeilen/s, Feed-Speicher-Spitze ${peak_stage} MiB (Ruhe vor der Stufe ${idle} MiB, 20 s nach dem letzten Lauf ${after} MiB, 60 s nach dem letzten Lauf ${after60} MiB, Zeilen in cdc.change danach: $(store_rows)), mittlere Blockdauer ${block_ms} ms (abgeleitet: Median-Dauer / $blocks Blöcke; die längste Einzeldauer ist nicht gemessen), WAL-Rückstand-Spitze im Run Median $(mib_of "$median_wal") MiB (~${wal_per_row} B je Zeile, abgeleitet), vom Slot gehaltenes WAL Median $(mib_of "$median_held") MiB"
  last_rate=$median_rate
  last_stage=$n
  last_wal_peak=$(printf '%s\n' "${wal_list[@]}" | sort -n | tail -n 1)
done
LC_ALL=C awk -v tag="$TAG" -v stage="$last_stage" -v peak="$last_wal_peak" -v warn="$WAL_WARN_BYTES" 'BEGIN {
  printf "%s: WAL-Rückstand der größten Stufe (%s Zeilen) — höchste Spitze im Run %.0f MiB, %s der Warnschwelle von %.0f MiB (SPEC-013)\n", tag, stage, peak / 1048576, (peak < warn) ? "unter" : "nicht unter", warn / 1048576
}'

# --- Wirkung auf die Live-Erfassung -------------------------------------------
echo "$TAG: Live-Erfassung — ${LIVE_RATE} Zeilen/s in public.$LIVE_TABLE, Run über Stufe $last_stage …"
start_writer 0
sleep 10
lag_run=(); lag_after=(); live_wal_peak=0; live_held_peak=0
run_start=$(date +%s)
id=$(bench::psql_scalar "SELECT cdc.backfill_table('$SOURCE_ID', 'public', '$(table_of "$last_stage")')")
while :; do
  status=$(bench::psql_scalar "SELECT status FROM cdc.backfill_run WHERE run_id = '$id'")
  lag=$(lag_sample)
  live_wal_peak=$(max_of "$live_wal_peak" "$(wal_backlog_bytes)")
  live_held_peak=$(max_of "$live_held_peak" "$(wal_held_bytes)")
  [ -n "$lag" ] && lag_run+=("$lag")
  case "$status" in
    completed) break ;;
    failed | interrupted) echo "$TAG: Live-Phase: Run $id endete $status" >&2; exit 1 ;;
  esac
  sleep 1
done
run_seconds=$(( $(date +%s) - run_start ))
for _ in $(seq 1 10); do
  sleep 1
  lag=$(lag_sample)
  [ -n "$lag" ] && lag_after+=("$lag")
done
stop_writer
live_run_rows=$(cat "$COUNT_FILE" 2>/dev/null || echo 0)

echo "$TAG: Live-Erfassung, Referenz ohne Run (${run_seconds} s) …"
start_writer 100000000
lag_ref=()
for _ in $(seq 1 "$run_seconds"); do
  sleep 1
  lag=$(lag_sample)
  [ -n "$lag" ] && lag_ref+=("$lag")
done
stop_writer
live_ref_rows=$(cat "$COUNT_FILE" 2>/dev/null || echo 0)

echo "$TAG: cdc_capture_lag während des Runs (${run_seconds} s, Stufe $last_stage, Live-Last ${LIVE_RATE}/s): ${lag_run[*]:+$(range_of "${lag_run[@]}")} s; in den 10 s danach: ${lag_after[*]:+$(range_of "${lag_after[@]}")} s; Referenz ohne Run: ${lag_ref[*]:+$(range_of "${lag_ref[@]}")} s (Live-Zeilen eingefügt: ${live_run_rows} im Lauf mit Run, ${live_ref_rows} in der Referenz); im Lauf mit Run: WAL-Rückstand des Slots (confirmed_flush_lsn) Spitze $(mib_of "$live_wal_peak") MiB, vom Slot gehaltenes WAL (restart_lsn) Spitze $(mib_of "$live_held_peak") MiB"

# --- Richtgröße aus der Messung ------------------------------------------------
tolerance_s=$((TOLERANCE_MIN * 60))
LC_ALL=C awk -v tag="$TAG" -v rate="$last_rate" -v t="$tolerance_s" -v stage="$last_stage" 'BEGIN {
  x = int(rate * t); d = length(x); p = 10 ^ (d - 1); r = int(x / p) * p
  printf "%s: Richtgröße (abgeleitet) — %s Zeilen/s (Median, Stufe %s) × Toleranz %d s = %d Zeilen, abgerundet auf eine Stelle: %d Zeilen (Startwert der Toleranz: Setzung ohne Messung; die Rate ist über die Stufe hinaus hochgerechnet)\n", tag, rate, stage, t, x, r
}'
echo "$TAG: Ende — Lauf $RUN_ID"
