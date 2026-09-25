#!/usr/bin/env bash
# bench-source-impact.sh — LH-QA-PER-001-Beleg (ADR-0104, Supersedes
# ADR-0054 §(b) teilweise): Schreibdurchsatz/-latenz derselben Insert-Last
# auf dieselbe Quelltabelle, einmal ohne jeden Replication-Slot/Capture-
# Prozess (Phase "ohne CDC") und einmal mit aktivem Slot und laufendem
# Feed-Container (Phase "mit CDC"). Je Phase mehrere Läufe, Median als
# Kennzahl (Vorbild: d-checks bench-fixture.sh, ADR-0054 Kontext).
# N und RUNS sind bewusst groß gewählt: Jede einzelne Insert-Anweisung
# läuft als eigene, separat committete Transaktion (siehe gen_inserts
# unten) — bei zu kurzer Gesamtdauer schlägt gewöhnliches WAL-Fsync-/
# Scheduling-Jitter relativ stark auf die Prozent-Kennzahl durch, real
# unabhängig von der Auslastung anderer Container auf derselben Maschine
# (separat per `docker stats` geprüft: 0–5 % CPU). Die gewählte Größe
# verlängert die absolute Messdauer je Lauf und mittelt über mehrere
# Läufe, damit ein einzelner Jitter-Ausschlag die Kennzahl nicht kippt.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
# shellcheck source=tools/bench-lib.sh
source tools/bench-lib.sh

trap bench::cleanup EXIT

# THRESHOLD_PCT trägt SPEC-025s Quell-Overhead-Schwelle (ADR-0104): der
# Median-Overhead über RUNS Läufe darf sie nicht überschreiten.
THRESHOLD_PCT=35
N=${BENCH_SOURCE_IMPACT_N:-5000}
RUNS=5
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
if [ "$ms_without_cdc" -gt 0 ]; then
  delta_pct=$(LC_ALL=C awk -v a="$ms_without_cdc" -v b="$ms_with_cdc" 'BEGIN { printf "%.1f", (b - a) * 100.0 / a }')
else
  delta_pct="n/a"
fi

echo "bench-source-impact: Ergebnis (LH-QA-PER-001) — ohne CDC ${ms_without_cdc} ms (Median von $RUNS Läufen), mit CDC ${ms_with_cdc} ms (Median von $RUNS Läufen), Differenz ${delta_ms} ms (${delta_pct}%) für je $N Schreibtransaktionen auf public.$TABLE"

bench::record_row "LH-QA-PER-001" \
  "Schreibdurchsatz/-latenz derselben Insert-Last mit/ohne aktivierter CDC im Vergleich" \
  "Overhead (Median von $RUNS Läufen) ≤ ${THRESHOLD_PCT}% (\`SPEC-025\`)" \
  "tools/bench-source-impact.sh"

if [ "$delta_pct" != "n/a" ] && [ "$(LC_ALL=C awk -v v="$delta_pct" -v t="$THRESHOLD_PCT" 'BEGIN { print (v > t) ? 1 : 0 }')" = "1" ]; then
  echo "bench-source-impact: SCHWELLE ÜBERSCHRITTEN (LH-QA-PER-001, SPEC-025) — Overhead ${delta_pct}% liegt über der zulässigen ${THRESHOLD_PCT}%-Schwelle" >&2
  exit 1
fi
