#!/usr/bin/env bash
# bench-scaling.sh — LH-QA-PER-002-Beleg (ADR-0054 §(b)): durchläuft die
# drei SPEC-014-Lastenstufen (klein ≤10/s, mittel 100/s×30min,
# groß 1000/s×60min) als feste Eingabeparameter gegen einen aktiven
# Feed-Container und liest je Stufe den erreichten Schreibdurchsatz sowie
# cdc_capture_lag. Kein Pass/Fail — dokumentiertes Ergebnis je Stufe.
#
# Slice-Plan §6 Risiko 2: Die "groß"-Stufe (1.000/s × 60 min) macht einen
# schnellen, wiederholten Lauf unpraktikabel lang. Default-Modus fährt
# deshalb stark verkürzte, aber reale Dauern je Stufe (Sekunden statt
# Minuten) bei UNVERÄNDERTER Ziel-Rate; `--full` fährt die tatsächlichen
# SPEC-014-Dauern (mittel/groß). SPEC-014 legt für die "klein"-Stufe keine
# Dauer fest (nur die Rate ≤10/s) — die hier gewählte volle Dauer (60s) ist
# eine dokumentierte Annahme dieses Bench-Skripts, keine Schärfung von
# SPEC-014 selbst.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
# shellcheck source=tools/bench-lib.sh
source tools/bench-lib.sh

trap bench::cleanup EXIT

FULL=0
if [ "${1:-}" = "--full" ]; then
  FULL=1
fi

TABLE=bench_scaling
SOURCE_ID=src-bench-scaling
SLOT=slot_bench_scaling
PUBLICATION=pub_bench_scaling

echo "bench-scaling: Umgebung wird aufgebaut (Modus: $([ "$FULL" -eq 1 ] && echo voll/SPEC-014 || echo verkürzt))…"
bench::start_postgres
bench::schema_rollout

bench::psql <<SQL
CREATE TABLE public.$TABLE (id int PRIMARY KEY, name text);
INSERT INTO cdc.source (source_id, name) VALUES ('$SOURCE_ID', 'Bench-Skalierung');
SQL

bench::start_feed "public.$TABLE=tbl-bench-scaling:sv-bench-scaling" "$SOURCE_ID" "$SLOT" "$PUBLICATION"

# Rate wird server-seitig über eine INSERT…generate_series-Anweisung je
# Sekunde erzeugt (statt N einzelner Schreibtransaktionen je Sekunde) —
# bei 1.000/s wären 1.000 einzelne docker-exec-Aufrufe pro Sekunde selbst
# ohne DB-Last nicht durchhaltbar; die generate_series-Anweisung bleibt ein
# reales, serverseitig verarbeitetes Schreib-Volumen in der Zielrate.
run_tier() {
  local label=$1 rate=$2 duration=$3 id_base=$4
  local inserted=0 second=0
  local start end batch_start batch_end batch_ms sleep_s
  start=$(date +%s)
  while [ "$second" -lt "$duration" ]; do
    batch_start=$(date +%s%N)
    bench::psql -c "INSERT INTO public.$TABLE (id, name) SELECT $id_base + $second * $rate + g, 'bench-$label' FROM generate_series(1, $rate) AS g;" >/dev/null
    inserted=$((inserted + rate))
    batch_end=$(date +%s%N)
    batch_ms=$(( (batch_end - batch_start) / 1000000 ))
    sleep_s=$(awk -v ms="$batch_ms" 'BEGIN { s = (1000 - ms) / 1000.0; if (s < 0) s = 0; printf "%.3f", s }')
    sleep "$sleep_s"
    second=$((second + 1))
  done
  end=$(date +%s)
  local elapsed=$((end - start))
  local lag
  lag=$(bench::psql_scalar "SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_capture_lag'")
  local achieved
  achieved=$(awk -v n="$inserted" -v s="$elapsed" 'BEGIN { if (s <= 0) s = 1; printf "%.1f", n / s }')
  echo "bench-scaling: Stufe $label — Ziel ${rate}/s über ${duration}s, $inserted Zeilen in ${elapsed}s eingefügt (~${achieved}/s), cdc_capture_lag=${lag}s"
}

if [ "$FULL" -eq 1 ]; then
  KLEIN_DURATION=60
  MITTEL_DURATION=1800
  GROSS_DURATION=3600
else
  KLEIN_DURATION=${BENCH_SCALING_KLEIN_SECONDS:-10}
  MITTEL_DURATION=${BENCH_SCALING_MITTEL_SECONDS:-15}
  GROSS_DURATION=${BENCH_SCALING_GROSS_SECONDS:-15}
fi

run_tier klein 10 "$KLEIN_DURATION" 0
run_tier mittel 100 "$MITTEL_DURATION" 1000000
run_tier gross 1000 "$GROSS_DURATION" 2000000

bench::stop_feed
echo "bench-scaling: Ergebnis (LH-QA-PER-002) — alle drei SPEC-014-Lastenstufen real durchlaufen (Modus: $([ "$FULL" -eq 1 ] && echo voll || echo verkürzt))"
