#!/usr/bin/env bash
# record-gates — Nachweis schreiben, dass `make gates` den aktuellen
# Arbeitsbaum-Zustand abgedeckt hat. Laeuft als LETZTER gates-Prerequisite (nur
# bei gruenen Gates). Der Stop-Hook vergleicht denselben Hash — ein Commit oder
# eine Aenderung ohne frischen Gate-Lauf laesst den Stop-Hook rot.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

mkdir -p .harness/state
bash tools/harness/working-tree-hash.sh > .harness/state/gates-passed.diffsha
