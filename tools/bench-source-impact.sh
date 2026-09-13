#!/usr/bin/env bash
# bench-source-impact.sh — LH-QA-PER-001-Beleg (ADR-0054 §(b)):
# Schreibdurchsatz/-latenz derselben Insert-Last auf dieselbe Quelltabelle,
# einmal ohne jeden Replication-Slot/Capture-Prozess (Phase "ohne CDC") und
# einmal mit aktivem Slot und laufendem Feed-Container (Phase "mit CDC").
# Dokumentiertes Ergebnis auf stdout, kein Pass/Fail-Schwellenwert — anders
# als d-checks bench-fixture.sh (eine Kennzahl gegen eine feste Schwelle,
# ADR-0054 Kontext) trägt dieser Beleg keine Schwelle.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
# shellcheck source=tools/bench-lib.sh
source tools/bench-lib.sh

trap bench::cleanup EXIT

N=${BENCH_SOURCE_IMPACT_N:-1000}
TABLE=bench_source_impact
SOURCE_ID=src-bench-impact
SLOT=slot_bench_impact
PUBLICATION=pub_bench_impact

echo "bench-source-impact: Umgebung wird aufgebaut (N=$N Zeilen je Phase) …"
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

run_phase() {
  local offset=$1 count=$2
  local start end
  start=$(date +%s%N)
  gen_inserts "$offset" "$count" | bench::psql >/dev/null
  end=$(date +%s%N)
  echo $(( (end - start) / 1000000 ))
}

echo "bench-source-impact: Phase 1 — ohne CDC (kein Replication-Slot aktiv) …"
ms_without_cdc=$(run_phase 0 "$N")
echo "bench-source-impact: ohne CDC — ${ms_without_cdc} ms für $N Zeilen"

bench::start_feed "public.$TABLE=tbl-bench-impact:sv-bench-impact" "$SOURCE_ID" "$SLOT" "$PUBLICATION"

echo "bench-source-impact: Phase 2 — mit CDC (aktiver Slot + Feed-Container) …"
ms_with_cdc=$(run_phase "$N" "$N")
echo "bench-source-impact: mit CDC — ${ms_with_cdc} ms für $N Zeilen"

bench::wait_captured "$SOURCE_ID" "$TABLE" "$N" 120 || true
bench::stop_feed

delta_ms=$(( ms_with_cdc - ms_without_cdc ))
if [ "$ms_without_cdc" -gt 0 ]; then
  delta_pct=$(awk -v a="$ms_without_cdc" -v b="$ms_with_cdc" 'BEGIN { printf "%.1f", (b - a) * 100.0 / a }')
else
  delta_pct="n/a"
fi

echo "bench-source-impact: Ergebnis (LH-QA-PER-001) — ohne CDC ${ms_without_cdc} ms, mit CDC ${ms_with_cdc} ms, Differenz ${delta_ms} ms (${delta_pct}%) für je $N Schreibtransaktionen auf public.$TABLE"
