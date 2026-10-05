#!/usr/bin/env bash
# bench-source-impact.sh — Schreiblatenz derselben Insert-Last auf dieselbe
# Quelltabelle, einmal ohne jeden Replication-Slot/Capture-Prozess (Phase
# "ohne CDC") und einmal mit aktivem Slot und laufendem Feed-Container
# (Phase "mit CDC"). Je Phase mehrere Läufe, Median als Kennzahl (Vorbild:
# d-checks bench-fixture.sh).
# Verdikt: die zusätzliche Latenz je Quelltransaktion
# Δ = (mit CDC − ohne CDC) / N darf max(0,10 ms; 1,5 × t_sync) nicht
# überschreiten; t_sync ist der Median von T_SYNC_SAMPLES im selben Lauf mit
# pg_test_fsync gemessenen Festschreib-Latenzen (ADR-0155). Exit 1 bei Δ über
# der Grenze, Exit 2 wenn t_sync nicht messbar ist (kein Rückfall auf eine
# relative Schwelle).
# N und RUNS sind bewusst groß gewählt: Jede einzelne Insert-Anweisung
# läuft als eigene, separat committete Transaktion (siehe gen_inserts
# unten) — bei zu kurzer Gesamtdauer schlägt gewöhnliches WAL-Fsync-/
# Scheduling-Jitter relativ stark auf Δ durch, real
# unabhängig von der Auslastung anderer Container auf derselben Maschine
# (separat per `docker stats` geprüft: 0–5 % CPU). Die gewählte Größe
# verlängert die absolute Messdauer je Lauf und mittelt über mehrere
# Läufe, damit ein einzelner Jitter-Ausschlag die Kennzahl nicht kippt.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
# shellcheck source=tools/bench-lib.sh
source tools/bench-lib.sh

trap bench::cleanup EXIT

# FACTOR und MIN_MS tragen die Grenze von SPEC-025: Faktor auf t_sync und
# absolute Untergrenze der Zusatzlatenz je Commit in ms.
FACTOR=1.5
MIN_MS=0.10
N=${BENCH_SOURCE_IMPACT_N:-5000}
RUNS=5
T_SYNC_SAMPLES=3
TABLE=bench_source_impact
SOURCE_ID=src-bench-impact
SLOT=slot_bench_impact
PUBLICATION=pub_bench_impact

echo "bench-source-impact: Umgebung wird aufgebaut (N=$N Zeilen je Lauf, $RUNS Läufe je Phase) …"
bench::start_postgres
bench::schema_rollout

bench::psql <<SQL
CREATE TABLE public.$TABLE (id int PRIMARY KEY, name text);
INSERT INTO cdc.source (source_id, name) VALUES ('$SOURCE_ID', 'Bench-Quell-Impact');
SQL

# Bewusst N separate INSERT-Anweisungen statt einer Bulk-Anweisung: das
# Ziel ist die Latenz je Schreibtransaktion unter Logical-Decoding-Last,
# nicht der reine Massendurchsatz einer einzelnen Anweisung.
gen_inserts() {
  local offset=$1 count=$2
  local i
  for i in $(seq 1 "$count"); do
    echo "INSERT INTO public.$TABLE (id, name) VALUES ($((offset + i)), 'row-$((offset + i))');"
  done
}

run_once() {
  local offset=$1 count=$2
  local start end
  start=$(date +%s%N)
  gen_inserts "$offset" "$count" | bench::psql >/dev/null
  end=$(date +%s%N)
  echo $(( (end - start) / 1000000 ))
}

# run_phase_median fährt RUNS Läufe derselben Insert-Last gegen disjunkte
# ID-Bereiche (keine PK-Kollision zwischen den Läufen) und liefert
# "Median|Lauf1,Lauf2,…" — der Aufrufer trennt beide Teile selbst, damit
# sowohl die Kennzahl als auch die Einzelwerte gemeldet werden.
run_phase_median() {
  local phase_offset_base=$1 count=$2
  local -a samples=()
  local run joined
  for run in $(seq 1 "$RUNS"); do
    samples+=("$(run_once $(( phase_offset_base + (run - 1) * count )) "$count")")
  done
  joined=$(IFS=,; echo "${samples[*]}")
  echo "$(bench::median_of "${samples[@]}")|${joined}"
}

# measure_t_sync_once liefert ein Sample der Festschreib-Latenz in ms auf
# stdout: Zeile der Methode $2 (aus `SHOW wal_sync_method`) im Abschnitt
# „one 8kB write“ der Ausgabe von `pg_test_fsync -s 2`, ausgeführt im
# PostgreSQL-Container $1 auf einer Datei im Datenverzeichnis. Nicht
# messbar: keine Ausgabe, Exit 1.
measure_t_sync_once() {
  local pg=$1 method=$2 out usecs
  out=$(docker exec "$pg" sh -c 'f="${PGDATA:-/var/lib/postgresql/data}/bench-sync-probe.tmp"; pg_test_fsync -s 2 -f "$f"; rc=$?; rm -f "$f"; exit $rc' 2>/dev/null) || return 1
  usecs=$(printf '%s\n' "$out" | LC_ALL=C awk -v m="$method" '$1 == m && /usecs\/op/ { print $(NF-1); exit }')
  case "$usecs" in
    ''|*[!0-9.]*) return 1 ;;
  esac
  LC_ALL=C awk -v u="$usecs" 'BEGIN { if (u + 0 <= 0) exit 1; printf "%.3f", u / 1000.0 }'
}

# measure_t_sync liefert "Median|Sample1,Sample2,…" (ms) über T_SYNC_SAMPLES
# Samples; ein nicht messbares Sample macht die ganze Messung nicht messbar
# (Exit 1, keine Ausgabe).
measure_t_sync() {
  local pg method s i
  local -a samples=()
  pg=$(bench::pg_container) || return 1
  method=$(bench::psql_scalar "SHOW wal_sync_method" | tr -d '[:space:]') || return 1
  [ -n "$method" ] || return 1
  for i in $(seq 1 "$T_SYNC_SAMPLES"); do
    s=$(measure_t_sync_once "$pg" "$method") || return 1
    samples+=("$s")
  done
  echo "$(LC_ALL=C bench::median_of "${samples[@]}")|$(IFS=,; echo "${samples[*]}")"
}

if ! t_sync_result=$(measure_t_sync); then
  echo "bench-source-impact: t_sync nicht messbar (pg_test_fsync im PostgreSQL-Container fehlt oder eine Sample-Zeile ist nicht lesbar) — kein Verdikt (LH-QA-PER-001, ADR-0155)" >&2
  exit 2
fi
t_sync_ms=${t_sync_result%%|*}
t_sync_samples=${t_sync_result#*|}
echo "bench-source-impact: t_sync ${t_sync_ms} ms (Median von $T_SYNC_SAMPLES pg_test_fsync-Samples: ${t_sync_samples} ms; Methode aus wal_sync_method)"

echo "bench-source-impact: Phase 1 — ohne CDC (kein Replication-Slot aktiv), $RUNS Läufe …"
without_result=$(run_phase_median 0 "$N")
ms_without_cdc=${without_result%%|*}
without_samples=${without_result#*|}
echo "bench-source-impact: ohne CDC — Median ${ms_without_cdc} ms für $N Zeilen (Einzelläufe: ${without_samples} ms)"

bench::start_feed "public.$TABLE=tbl-bench-impact:sv-bench-impact" "$SOURCE_ID" "$SLOT" "$PUBLICATION"

echo "bench-source-impact: Phase 2 — mit CDC (aktiver Slot + Feed-Container), $RUNS Läufe …"
with_result=$(run_phase_median $(( RUNS * N )) "$N")
ms_with_cdc=${with_result%%|*}
with_samples=${with_result#*|}
echo "bench-source-impact: mit CDC — Median ${ms_with_cdc} ms für $N Zeilen (Einzelläufe: ${with_samples} ms)"

bench::wait_captured "$SOURCE_ID" "$TABLE" $(( RUNS * N )) 300 || true
bench::stop_feed

delta_ms=$(( ms_with_cdc - ms_without_cdc ))
# Δ je Transaktion, Verhältnis zu t_sync und Grenze (Gleitkomma, "delta|ratio|limit|over").
verdict=$(LC_ALL=C awk -v d="$delta_ms" -v n="$N" -v t="$t_sync_ms" -v f="$FACTOR" -v m="$MIN_MS" 'BEGIN {
  delta = d / n
  limit = f * t; if (m + 0 > limit) limit = m + 0
  printf "%.3f|%.3f|%.3f|%d", delta, delta / t, limit, (delta > limit) ? 1 : 0
}')
IFS='|' read -r delta_tx_ms ratio limit_ms over <<<"$verdict"

echo "bench-source-impact: Ergebnis (LH-QA-PER-001) — ohne CDC ${ms_without_cdc} ms (Median von $RUNS Läufen), mit CDC ${ms_with_cdc} ms (Median von $RUNS Läufen), Differenz ${delta_ms} ms für je $N Schreibtransaktionen auf public.$TABLE"
echo "bench-source-impact: Zusatzlatenz Δ ${delta_tx_ms} ms je Transaktion, t_sync ${t_sync_ms} ms, Verhältnis Δ/t_sync ${ratio}, Grenze max(${MIN_MS} ms; ${FACTOR} × t_sync) = ${limit_ms} ms"

bench::record_row "LH-QA-PER-001" \
  "Zusatzlatenz je Schreibtransaktion derselben Insert-Last mit/ohne aktivierter CDC, bezogen auf die gemessene Festschreib-Latenz" \
  "Zusatzlatenz je Commit (Median von $RUNS Läufen) ≤ max(${MIN_MS} ms; ${FACTOR} × Festschreib-Latenz) (\`SPEC-025\`)" \
  "tools/bench-source-impact.sh"

if [ "$over" = "1" ]; then
  echo "bench-source-impact: GRENZE ÜBERSCHRITTEN (LH-QA-PER-001, SPEC-025) — Δ ${delta_tx_ms} ms liegt über der Grenze ${limit_ms} ms (t_sync ${t_sync_ms} ms)" >&2
  exit 1
fi
