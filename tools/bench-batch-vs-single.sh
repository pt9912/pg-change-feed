#!/usr/bin/env bash
# bench-batch-vs-single.sh — LH-QA-PER-003-Beleg (ADR-0054 §(b)): vergleicht
# das Lesen von M bereits erfassten Change-Zeilen über cdc.changes
# (LH-FA-SST-002) einmal als ein einziger Batch-Abruf (eine Anweisung, ein
# Roundtrip) und einmal als M Einzelabrufe (eine Anweisung je Zeile, M
# Roundtrips) — dokumentiertes Ergebnis, kein Pass/Fail-Schwellenwert.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
# shellcheck source=tools/bench-lib.sh
source tools/bench-lib.sh

trap bench::cleanup EXIT

# THRESHOLD_FACTOR trägt SPEC-025s Batch-Vorteil-Mindestschwelle
# (ADR-0104): der Einzelabruf muss mindestens so viel langsamer sein.
THRESHOLD_FACTOR=10
M=${BENCH_BATCH_VS_SINGLE_M:-200}
TABLE=bench_batch_vs_single
SOURCE_ID=src-bench-batch
SLOT=slot_bench_batch
PUBLICATION=pub_bench_batch

echo "bench-batch-vs-single: Umgebung wird aufgebaut (M=$M Zeilen) …"
bench::start_postgres
bench::schema_rollout

bench::psql <<SQL
CREATE TABLE public.$TABLE (id int PRIMARY KEY, name text);
INSERT INTO cdc.source (source_id, name) VALUES ('$SOURCE_ID', 'Bench-Batch-vs-Einzelabruf');
SQL

bench::start_feed "public.$TABLE=tbl-bench-batch:sv-bench-batch" "$SOURCE_ID" "$SLOT" "$PUBLICATION"

bench::psql -c "INSERT INTO public.$TABLE (id, name) SELECT g, 'row-' || g FROM generate_series(1, $M) AS g;" >/dev/null

echo "bench-batch-vs-single: warte auf Erfassung von $M Zeilen über cdc.changes …"
bench::wait_captured "$SOURCE_ID" "$TABLE" "$M" 120

echo "bench-batch-vs-single: Batch-Abruf (eine Anweisung, LIMIT $M) …"
batch_start=$(date +%s%N)
bench::psql_scalar "SELECT count(*) FROM (SELECT * FROM cdc.changes WHERE source_id = '$SOURCE_ID' AND table_name = '$TABLE' ORDER BY commit_position LIMIT $M) t" >/dev/null
batch_end=$(date +%s%N)
batch_ms=$(( (batch_end - batch_start) / 1000000 ))
echo "bench-batch-vs-single: Batch-Abruf — ${batch_ms} ms für $M Zeilen (1 Roundtrip)"

echo "bench-batch-vs-single: Einzelabruf ($M Anweisungen, je 1 Zeile) …"
single_start=$(date +%s%N)
i=0
while [ "$i" -lt "$M" ]; do
  bench::psql_scalar "SELECT commit_position FROM cdc.changes WHERE source_id = '$SOURCE_ID' AND table_name = '$TABLE' ORDER BY commit_position LIMIT 1 OFFSET $i" >/dev/null
  i=$((i + 1))
done
single_end=$(date +%s%N)
single_ms=$(( (single_end - single_start) / 1000000 ))
echo "bench-batch-vs-single: Einzelabruf — ${single_ms} ms für $M Zeilen ($M Roundtrips)"

bench::stop_feed

factor=$(LC_ALL=C awk -v b="$batch_ms" -v s="$single_ms" 'BEGIN { if (b <= 0) b = 1; printf "%.1f", s / b }')
echo "bench-batch-vs-single: Ergebnis (LH-QA-PER-003) — Batch ${batch_ms} ms vs. Einzelabruf ${single_ms} ms für $M Zeilen (Einzelabruf ${factor}x langsamer)"

bench::record_row "LH-QA-PER-003" \
  "Lesevorgang über größere Change-Mengen im Batch gegen Einzelabruf" \
  "Batch-Vorteil ≥ ${THRESHOLD_FACTOR}× (\`SPEC-025\`)" \
  "tools/bench-batch-vs-single.sh"
bench::render_abdeckung

if [ "$(LC_ALL=C awk -v v="$factor" -v t="$THRESHOLD_FACTOR" 'BEGIN { print (v < t) ? 1 : 0 }')" = "1" ]; then
  echo "bench-batch-vs-single: SCHWELLE UNTERSCHRITTEN (LH-QA-PER-003, SPEC-025) — Batch-Vorteil ${factor}x liegt unter der geforderten ${THRESHOLD_FACTOR}x-Schwelle" >&2
  exit 1
fi
