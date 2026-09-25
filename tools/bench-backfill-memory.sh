#!/usr/bin/env bash
# bench-backfill-memory.sh — Speicher des Feed-Containers im Backfill-Run
# (LH-FA-CAP-009, ADR-0111, ADR-0113): Verlauf im Run, Nachlauf, Zeilenzahl,
# Zeilenbreite, Feed-Einstellung. Eine Messung mit gedruckter Zahl, kein
# Pass/Fail. Vertrag, Ablauf und Grenzen: harness/targets/bench-backfill.md.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
# shellcheck source=tools/bench-lib.sh
source tools/bench-lib.sh

STAGES=${BENCH_MEM_STAGES:-"10000 100000 200000"}
RUNS=${BENCH_MEM_RUNS:-3}
WIDTH=${BENCH_MEM_WIDTH:-narrow}
SERIES_DIR=${BENCH_MEM_SERIES_DIR:-}
RESET=${BENCH_MEM_RESET:-1}
RUN_TIMEOUT_S=${BENCH_MEM_RUN_TIMEOUT_S:-3600}
case "$WIDTH" in narrow | wide) ;; *) echo "bench-backfill-memory: BENCH_MEM_WIDTH ist narrow oder wide" >&2; exit 1 ;; esac

SOURCE_ID=src-bench-mem
SLOT=slot_bench_mem
PUBLICATION=pub_bench_mem
RUN_ID=$(date -u +%Y%m%dT%H%M%SZ)
TAG="bench-backfill-memory[$RUN_ID]"
WINDOW_FILE=$(mktemp)
trap 'rm -f "$WINDOW_FILE" "$WINDOW_FILE.cur" "$WINDOW_FILE.anon" "$WINDOW_FILE.series"; bench::cleanup' EXIT
trap 'exit 1' TERM

table_of() { echo "bench_mem_${WIDTH}_$1"; }

fill_table() {
  local table=$1 rows=$2
  if [ "$WIDTH" = narrow ]; then
    bench::psql >/dev/null <<SQL
CREATE TABLE public.$table (id bigint PRIMARY KEY, name text NOT NULL, amount numeric(12,2), created_at timestamptz, note text) WITH (autovacuum_enabled = false);
INSERT INTO public.$table
SELECT g, 'kunde-' || g, (g % 100000) / 100.0, now(), CASE WHEN g % 5 = 0 THEN NULL ELSE 'notiz-' || g END
FROM generate_series(1, $rows) AS g;
SQL
  else
    bench::psql >/dev/null <<SQL
CREATE TABLE public.$table (id bigint PRIMARY KEY, name text NOT NULL, amount numeric(12,2), created_at timestamptz, note text, payload jsonb, body text) WITH (autovacuum_enabled = false);
INSERT INTO public.$table
SELECT g, 'kunde-' || g, (g % 100000) / 100.0, now(), CASE WHEN g % 5 = 0 THEN NULL ELSE 'notiz-' || g END,
       jsonb_build_object('a', repeat(md5(g::text), 8), 'b', repeat(md5((g + 1)::text), 8), 'n', g),
       repeat(md5((g + 2)::text), 20)
FROM generate_series(1, $rows) AS g;
SQL
  fi
}

mib() { LC_ALL=C awk -v b="$1" 'BEGIN { printf "%.1f", b / 1048576 }'; }

# Liest die Speicherzähler des Feed-Containers vom Host: cgroup v2 des
# Containers (memory.current, memory.stat) und /proc/<pid>/status des
# Prozesses; setzt die Variablen S_CUR, S_ANON, S_FILE, S_KERNEL, S_RSSANON.
mem_sample() {
  local k v _
  if ! { S_CUR=$(<"$CGROUP/memory.current"); } 2>/dev/null; then
    echo "$TAG: Feed-Container nicht mehr lesbar: $(docker inspect --format 'Status {{.State.Status}}, Exit {{.State.ExitCode}}, OOMKilled {{.State.OOMKilled}}' "$FEED" 2>&1)" >&2
    kill -TERM $$
    exit 1
  fi
  while read -r k v; do
    case "$k" in anon) S_ANON=$v ;; file) S_FILE=$v ;; kernel) S_KERNEL=$v ;; esac
  done < "$CGROUP/memory.stat"
  while read -r k v _; do
    if [ "$k" = "RssAnon:" ]; then S_RSSANON=$((v * 1024)); fi
  done < "/proc/$FEED_PID/status"
}

hwm_kib() { awk '$1 == "VmHWM:" { print $2 }' "/proc/$FEED_PID/status"; }
peak_bytes() { <"$CGROUP/memory.peak" cat; }

stats_mib() {
  docker stats --no-stream --format '{{.MemUsage}}' "$(bench::feed_container)" 2>/dev/null | LC_ALL=C awk '{
    v = $1; unit = v; gsub(/[0-9.]/, "", unit); sub(/[A-Za-z]+$/, "", v)
    if (unit == "GiB") printf "%.1f", v * 1024
    else if (unit == "MiB") printf "%.1f", v
    else if (unit == "KiB") printf "%.1f", v / 1024
    else printf "0.0"
  }'
}

probe_line() {
  mem_sample
  echo "cgroup ${1}: memory.current $(mib "$S_CUR") MiB (anon $(mib "$S_ANON"), file $(mib "$S_FILE"), kernel $(mib "$S_KERNEL")), Prozess RssAnon $(mib "$S_RSSANON") MiB, docker stats $(stats_mib) MiB"
}

# Wartet, bis die Sitzung des entfernten Feeds den Slot freigibt; nach 10 s
# beendet sie der Lauf und nennt es.
wait_slot_inactive() {
  local i pid
  for i in $(seq 1 10); do
    if [ "$(bench::psql_scalar "SELECT COALESCE((SELECT active::text FROM pg_replication_slots WHERE slot_name = '$SLOT'), 'f')")" = "f" ]; then
      return 0
    fi
    sleep 1
  done
  pid=$(bench::psql_scalar "SELECT active_pid FROM pg_replication_slots WHERE slot_name = '$SLOT'")
  echo "$TAG: Slot $SLOT nach 10 s noch aktiv (Sitzung $pid: $(bench::psql_scalar "SELECT state || ' ' || COALESCE(wait_event, '-') FROM pg_stat_activity WHERE pid = $pid")) — Sitzung beendet"
  bench::psql_scalar "SELECT pg_terminate_backend($pid)" >/dev/null
  sleep 2
}

store_rows() { bench::psql_scalar "SELECT count(*) FROM cdc.change"; }

# Tastet den Speicher $1 Sekunden lang alle 0,5 s ab und druckt Anfang,
# Maximum, Minimum und Ende von memory.current und anon.
window_line() {
  local seconds=$1 label=$2 i
  local -a curs=() anons=()
  for i in $(seq 1 $((seconds * 2))); do
    mem_sample
    curs+=("$S_CUR")
    anons+=("$S_ANON")
    sleep 0.5
  done
  printf '%s\n' "${curs[@]}" > "$WINDOW_FILE.cur"
  printf '%s\n' "${anons[@]}" > "$WINDOW_FILE.anon"
  echo "$label: memory.current $(mib "${curs[0]}") MiB zu Beginn, Maximum $(mib "$(sort -n "$WINDOW_FILE.cur" | tail -n 1)"), Minimum $(mib "$(sort -n "$WINDOW_FILE.cur" | head -n 1)"), Ende $(mib "${curs[-1]}"); anon Beginn $(mib "${anons[0]}"), Maximum $(mib "$(sort -n "$WINDOW_FILE.anon" | tail -n 1)"), Minimum $(mib "$(sort -n "$WINDOW_FILE.anon" | head -n 1)"), Ende $(mib "${anons[-1]}") MiB"
}

start_feed_stage() {
  local i active
  wait_slot_inactive
  bench::start_feed "$TABLES_ENV" "$SOURCE_ID" "$SLOT" "$PUBLICATION"
  for i in $(seq 1 60); do
    active=$(bench::psql_scalar "SELECT COALESCE((SELECT active::text FROM pg_replication_slots WHERE slot_name = '$SLOT'), 'f')")
    if [ "$active" = "t" ]; then break; fi
    sleep 1
  done
  FEED=$(bench::feed_container)
  FEED_ID=$(docker inspect --format '{{.Id}}' "$FEED")
  FEED_PID=$(docker inspect --format '{{.State.Pid}}' "$FEED")
  CGROUP=/sys/fs/cgroup/system.slice/docker-$FEED_ID.scope
  if [ ! -r "$CGROUP/memory.current" ] || [ ! -r "/proc/$FEED_PID/status" ]; then
    echo "$TAG: cgroup ($CGROUP) oder /proc/$FEED_PID nicht lesbar" >&2
    return 1
  fi
}

# Ein Run: beantragt ihn, tastet den Speicher etwa alle 0,2 s ab, bis er endet.
# Die Reihe steht in $series ("ms cur anon file kernel rssanon rows"); druckt
# die Zusammenfassung und setzt RUN_MS.
run_one() {
  local table=$1 n=$2 label=$3 id t0 now status rows sample deadline series
  series=$WINDOW_FILE.series
  : > "$series"
  id=$(bench::psql_scalar "SELECT cdc.backfill_table('$SOURCE_ID', 'public', '$table')")
  t0=$(date +%s%N)
  deadline=$((SECONDS + RUN_TIMEOUT_S))
  while :; do
    mem_sample
    sample=$(bench::psql_scalar "SELECT COALESCE((SELECT status || '|' || rows_copied FROM cdc.backfill_run WHERE run_id = '$id'), '|0')")
    status=${sample%%|*}
    rows=${sample#*|}
    now=$(date +%s%N)
    echo "$(((now - t0) / 1000000)) $S_CUR $S_ANON $S_FILE $S_KERNEL $S_RSSANON $rows" >> "$series"
    case "$status" in
      completed) break ;;
      failed | interrupted)
        echo "$TAG: Run $id von $table endete $status: $(bench::psql_scalar "SELECT error_message FROM cdc.backfill_run WHERE run_id = '$id'")" >&2
        echo "$TAG: Feed-Container: $(docker inspect --format 'Status {{.State.Status}}, Exit {{.State.ExitCode}}, OOMKilled {{.State.OOMKilled}}' "$FEED" 2>&1)" >&2
        return 1
        ;;
    esac
    if [ "$SECONDS" -gt "$deadline" ]; then
      echo "$TAG: Run $id von $table nicht beendet nach ${RUN_TIMEOUT_S} s" >&2
      return 1
    fi
    sleep 0.2
  done
  RUN_MS=$(bench::psql_scalar "SELECT round(EXTRACT(EPOCH FROM (finished_at - started_at)) * 1000)::bigint FROM cdc.backfill_run WHERE run_id = '$id'")
  echo "$TAG: Stufe $n $WIDTH, $label — Kopierdauer ${RUN_MS} ms, $(LC_ALL=C awk -v n="$n" -v ms="$RUN_MS" 'BEGIN { printf "%.0f", n * 1000.0 / ms }') Zeilen/s; $(LC_ALL=C awk -v n="$n" '
    { if ($2 > pc) pc = $2; if ($3 > pa) pa = $3; if (!q25 && $7 >= n * 0.25) { q25 = 1; a25 = $3; c25 = $2 } if (!q50 && $7 >= n * 0.5) { q50 = 1; a50 = $3; c50 = $2 } if (!q75 && $7 >= n * 0.75) { q75 = 1; a75 = $3; c75 = $2 } a0 = (NR == 1) ? $3 : a0; c0 = (NR == 1) ? $2 : c0; al = $3; cl = $2; cnt = NR }
    END { printf "Proben %d; memory.current Beginn %.1f, bei 25%% %.1f, 50%% %.1f, 75%% %.1f, Ende %.1f, Spitze %.1f MiB; anon Beginn %.1f, bei 25%% %.1f, 50%% %.1f, 75%% %.1f, Ende %.1f, Spitze %.1f MiB", cnt, c0/1048576, c25/1048576, c50/1048576, c75/1048576, cl/1048576, pc/1048576, a0/1048576, a25/1048576, a50/1048576, a75/1048576, al/1048576, pa/1048576 }' "$series")"
  if [ -n "$SERIES_DIR" ]; then
    mkdir -p "$SERIES_DIR"
    cp "$series" "$SERIES_DIR/${RUN_ID}_${WIDTH}_${n}_${label//[ \/]/_}.series"
  fi
  rm -f "$series"
}

# GC-Zeilen (GODEBUG=gctrace=1) seit der Zeile $1: Zahl, größter Heap zu
# Beginn, größtes Ziel, letzter Heap nach dem GC (MB). Ohne GC-Zeile im Fenster
# nennt die Ausgabe, ob GODEBUG=gctrace=1 in BENCH_FEED_ENV steht.
gc_summary() {
  local from=$1 traced=0
  case " ${BENCH_FEED_ENV:-} " in *" GODEBUG="*gctrace=1*) traced=1 ;; esac
  docker logs "$FEED" 2>&1 | grep -E '^gc [0-9]+ @' | tail -n +"$((from + 1))" | LC_ALL=C awk -v traced="$traced" '{
    if (match($0, /[0-9]+->[0-9]+->[0-9]+ MB, [0-9]+ MB goal/)) {
      split(substr($0, RSTART, RLENGTH), p, /->| MB, | MB goal/)
      if (p[1] + 0 > s) s = p[1] + 0
      if (p[4] + 0 > g) g = p[4] + 0
      l = p[3] + 0; c++
    }
  } END {
    if (c > 0) printf "GC-Läufe %d, größter Heap zu Beginn %d MB, größtes Ziel %d MB, letzter Heap nach GC %d MB", c, s, g, l
    else if (traced == 1) printf "keine GC-Zeilen im Fenster (GODEBUG=gctrace=1 gesetzt)"
    else printf "keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt)"
  }'
}
gc_count() { docker logs "$FEED" 2>&1 | grep -cE '^gc [0-9]+ @' || true; }

# Bereinigungs-Takte des Feeds seit seinem Start (Zeilen des Protokolls):
# Zahl der gelaufenen, Zahl der fehlgeschlagenen und die letzte Fehlermeldung.
retention_summary() {
  docker logs "$FEED" 2>&1 | LC_ALL=C awk '
    /retention: Bereinigung gelaufen/ { ok++ }
    /retention: Bereinigung fehlgeschlagen/ { bad++; last = $0 }
    END { printf "Bereinigung gelaufen %d, fehlgeschlagen %d", ok, bad; if (bad > 0) printf ", letzte Meldung: %s", substr(last, 1, 300) }'
}

echo "$TAG: Host — $(uname -sr), Docker $(docker info --format '{{.ServerVersion}}, {{.NCPU}} CPU, {{.MemTotal}} Byte RAM'), PostgreSQL-Image ${PG_TEST_IMAGE%%@*}, Feed-Image-ID $(docker image inspect --format '{{.Id}}' "$FEED_IMAGE")"
echo "$TAG: Stufen $STAGES, $RUNS Runs je Stufe, Zeilenbreite $WIDTH, Feed-Umgebung '${BENCH_FEED_ENV:-}', Feed-Docker-Argumente '${BENCH_FEED_DOCKER_ARGS:-}', Feed-Image $FEED_IMAGE"
bench::start_postgres
bench::schema_rollout
bench::psql -c "INSERT INTO cdc.source (source_id, name) VALUES ('$SOURCE_ID', 'Bench-Speicher');" >/dev/null

TABLES_ENV=""
for n in $STAGES; do
  table=$(table_of "$n")
  fill_table "$table" "$n"
  if [ -n "$TABLES_ENV" ]; then TABLES_ENV="$TABLES_ENV,"; fi
  TABLES_ENV="${TABLES_ENV}public.$table=tbl-mem-$n:sv-mem-$n"
  echo "$TAG: Tabelle $table — $n Zeilen, mittlere Zeilenbreite $(bench::psql_scalar "SELECT round(avg(pg_column_size(t.*)))::int FROM public.$table t") B (pg_column_size), $(bench::psql_scalar "SELECT round(avg(octet_length(t::text)))::int FROM public.$table t") B (Textform)"
done

for n in $STAGES; do
  table=$(table_of "$n")
  if [ "$RESET" = 1 ]; then
    bench::psql -c "TRUNCATE cdc.change, cdc.transaction" >/dev/null
  fi
  for run in $(seq 1 "$RUNS"); do
    unit="Stufe $n $WIDTH, Run $run/$RUNS"
    start_feed_stage
    echo "$TAG: $unit — Zeilen in cdc.change vor dem Run: $(store_rows); $(window_line 10 'Grundlinie, Ruhe 10 s nach dem Start')"
    gc0=$(gc_count)
    run_one "$table" "$n" "Run $run/$RUNS"
    echo "$TAG: $unit — unmittelbar nach dem Run: $(probe_line 'Run-Ende'), VmHWM $(hwm_kib) KiB, memory.peak seit Start des Feeds $(mib "$(peak_bytes)") MiB; im Run: $(gc_summary "$gc0")"
    gc1=$(gc_count)
    echo "$TAG: $unit — Nachlauf 60 s ohne Schreibzugriff: $(window_line 60 'Nachlauf')"
    echo "$TAG: $unit — Zeilen in cdc.change nach dem Nachlauf: $(store_rows), VmHWM $(hwm_kib) KiB, memory.peak seit Start des Feeds $(mib "$(peak_bytes)") MiB, $(docker inspect --format 'OOMKilled {{.State.OOMKilled}}, Neustarts {{.RestartCount}}' "$FEED"); im Nachlauf: $(gc_summary "$gc1"); seit dem Start des Feeds: $(retention_summary)"
    if [ "${BENCH_MEM_EMPTY_AFTER:-0}" = 1 ]; then
      gc2=$(gc_count)
      bench::psql -c "TRUNCATE cdc.change, cdc.transaction" >/dev/null
      echo "$TAG: $unit — cdc.change geleert, Nachlauf 120 s: $(window_line 120 'Nachlauf nach der Leerung'); im Fenster: $(gc_summary "$gc2")"
    fi
  done
done
echo "$TAG: Ende — Lauf $RUN_ID"
