#!/usr/bin/env bash
# db-coverage.sh — mergt die -coverprofile der beiden DB-gestuetzten
# Traeger-Laeufe (`make test-store`, `make test-replication`) zu EINER Zahl und
# prueft sie gegen DB_COVERAGE_THRESHOLD.
#
# Subject der Zahl sind die drei Pakete, deren Testlauf einen externen
# PostgreSQL voraussetzt (ADR-0071 Punkt 3): postgresstorage OHNE das
# Unterpaket mapper (das bleibt Gegenstand des Unit-Gates), postgresack,
# replication/receive. Die Zahl traegt dieses Subjekt in ihrem Namen
# ("DB-Adapter-Coverage") und heisst nie "die Coverage" des Repos — das ist die
# Gate-getragene Unit-Zahl (ADR-0071 Punkt 4).
#
# Kein Gate: der Traeger ist der nicht-blockierende Workflow
# .github/workflows/e2e.yml, nicht `make gates` — das Gate laeuft bei jedem
# Commit und bekommt keinen PostgreSQL-Container.
#
# Zaehlbasis: -coverpkg instrumentiert nur die in einem Testbinary VERLINKTEN
# Gegenstands-Pakete (der Lauf ueber postgresstorage allein traegt 610
# Statements fuer dieses Paket und keine Zeile fuer postgresack/receive). Im
# Replication-Lauf testet `go test` zwei Pakete (postgresack,
# replication/receive); jedes der beiden Testbinaries instrumentiert beide
# Gegenstands-Pakete, darum erscheint jede Block-Position dort ZWEIMAL —
# einmal mit ihrem count, einmal mit 0. "Gedeckt" heisst: mindestens ein
# Vorkommen traegt count > 0. Der Merge unten dedupliziert ueber die
# Block-Position und traegt je Position 1 (gedeckt) bzw. 0 — dieselbe Basis,
# die harness/sensors/coverage-gate.md §Zaehlbasis fuer die Unit-Zahl
# beschreibt. Ohne diese Regel (nur das erste Vorkommen) faellt
# replication/receive auf 0 von 155.
#
# Die beiden Laeufe messen VERSCHIEDENE Testbestaende und partitionieren das
# Subject: postgresstorage laeuft nur mit CDC_STORE_TEST_DSN (make test-store),
# postgresack/replication/receive nur mit CDC_REPLICATION_TEST_DSN (make
# test-replication). Jeder Lauf instrumentiert dabei seinen Teil; die beiden
# Dateimengen sind DISJUNKT. Das Merge ist deshalb sinnvoll — die Vereinigung
# der gedeckten Positionen beschreibt das ganze Subject; eine der beiden Laeufe
# allein traegt nur seine Haelfte.
#
# Aufruf:
#   db-coverage.sh --coverpkg   druckt die -coverpkg-Liste des Subjects
#   db-coverage.sh              mergt die vorhandenen Profile aus
#                               $DB_COVERAGE_DIR, druckt die Zahl und prueft sie
#                               gegen DB_COVERAGE_THRESHOLD, sobald beide
#                               Profile vorliegen; sonst Teilzahl ohne Exit 1.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

DB_COVERAGE_PKGS="./internal/adapters/driven/postgresstorage,./internal/adapters/driven/postgresack,./internal/adapters/driving/replication/receive"
DB_COVERAGE_DIR=${DB_COVERAGE_DIR:-${TMPDIR:-/tmp}/pg-change-feed-db-coverage}
# Geltende Stufe der Rampe (bootstrap-aware, ADR-0054 §(a), ADR-0071 Punkt 3).
# Der bewegliche Wert steht ausschliesslich hier; README und Sensor nennen die
# Rampe und den Verweis auf diesen Ort. Override: DB_COVERAGE_THRESHOLD=…
# (Hochschalt-Trigger: die naechste Ausbau-Stufe schliesst die Luecke zur
# naechsten vollen 5-%-Stufe).
DB_COVERAGE_THRESHOLD=${DB_COVERAGE_THRESHOLD:-70}

if [[ "${1:-}" == "--coverpkg" ]]; then
  printf '%s\n' "$DB_COVERAGE_PKGS"
  exit 0
fi

mkdir -p "$DB_COVERAGE_DIR"
store="$DB_COVERAGE_DIR/store.coverprofile"
repl="$DB_COVERAGE_DIR/replication.coverprofile"

profiles=()
[[ -f "$store" ]] && profiles+=("$store")
[[ -f "$repl" ]] && profiles+=("$repl")

if [[ ${#profiles[@]} -eq 0 ]]; then
  echo "db-coverage: kein Profil unter $DB_COVERAGE_DIR — zuerst make test-store / make test-replication laufen lassen" >&2
  exit 2
fi

merged="$DB_COVERAGE_DIR/merged.coverprofile"
awk '
  BEGIN { print "mode: atomic" }
  /^mode:/ { next }
  NF < 3 { next }
  {
    pos = $1
    if (!(pos in seen)) { seen[pos] = 1; order[++n] = pos; s[pos] = $2 }
    if ($3 + 0 > 0) cov[pos] = 1
  }
  END { for (i = 1; i <= n; i++) printf "%s %s %d\n", order[i], s[order[i]], (order[i] in cov ? 1 : 0) }
' "${profiles[@]}" > "$merged"

read -r covered total pct < <(awk '
  /^mode:/ { next }
  NF < 3 { next }
  { t += $2; if ($3 + 0 > 0) c += $2 }
  END { printf "%d %d %.2f\n", c, t, (t > 0 ? c * 100 / t : 0) }
' "$merged")

names=""
[[ -f "$store" ]] && names+="store"
[[ -f "$repl" ]] && names+="${names:+,}replication"

if [[ ${#profiles[@]} -lt 2 ]]; then
  printf 'DB-Adapter-Coverage (Teilzahl, nur %s-Bestand): %s%% (gedeckt %s von %s Statements; Profil: %s)\n' \
    "$names" "$pct" "$covered" "$total" "$names"
  echo "db-coverage: beide Profile noetig fuer die Schwelle (store + replication) — dieser Lauf prueft sie nicht" >&2
  exit 0
fi

printf 'DB-Adapter-Coverage: %s%% (gedeckt %s von %s Statements; Profile gemergt: %s)\n' \
  "$pct" "$covered" "$total" "$names"

if awk -v p="$pct" -v t="$DB_COVERAGE_THRESHOLD" 'BEGIN { exit !(p + 0 >= t + 0) }'; then
  printf 'db-coverage: OK — DB-Adapter-Coverage %s%% erfuellt Schwelle %s%%\n' "$pct" "$DB_COVERAGE_THRESHOLD"
else
  printf 'db-coverage: FAIL — DB-Adapter-Coverage %s%% unter Schwelle %s%%\n' "$pct" "$DB_COVERAGE_THRESHOLD" >&2
  exit 1
fi
