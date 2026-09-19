#!/usr/bin/env bash
# tools/harness/pin-stale-baseline.sh — advisory Sensor für P8 (ADR-0051
# Pin-Inventar): die adoptierte Kurs-Baseline-Version
# (`harness/conventions.md` §Baseline, kein Docker-Pin) gegen den
# neuesten GitHub-Release von pt9912/ai-harness-course. Kein Gate; braucht
# Netz (GitHub-Releases-API); meldet, hebt nicht an — eine Baseline-
# Aktualisierung ist ein bewusster Bootstrap-Vorgang (Modul 2), kein
# automatischer Nachlauf.
#
# Exit: 0 = keine Drift · 1 = neuerer Release existiert · 2 = API nicht
# erreichbar oder Version nicht auffindbar (Ergebnis unbestimmbar).
set -uo pipefail

file="harness/conventions.md"
pinned=$(grep -E '^\- \*\*Stand:\*\*' "$file" | head -1 | sed -E 's/^- \*\*Stand:\*\* v?//')
if [ -z "$pinned" ]; then
  echo "UNBESTIMMT  Kurs-Baseline — Stand-Zeile in $file nicht gefunden"
  exit 2
fi

latest=$(curl -fsS -m 15 https://api.github.com/repos/pt9912/ai-harness-course/releases/latest 2>/dev/null | grep -m1 '"tag_name"' | sed -E 's/.*"tag_name":[[:space:]]*"v?([^"]+)".*/\1/')
if [ -z "$latest" ]; then
  echo "UNBESTIMMT  Kurs-Baseline — GitHub-Releases-API nicht erreichbar"
  exit 2
fi

if [ "$latest" = "$pinned" ]; then
  echo "OK          Kurs-Baseline v$pinned == neuester Release"
  exit 0
fi

echo "DRIFT       Kurs-Baseline: adoptiert v$pinned, neuester Release v$latest"
exit 1
