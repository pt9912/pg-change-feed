#!/usr/bin/env bash
# rollout-restore — führt ein Kommando aus und stellt danach die Erzeugnisse
# von `make schema-rollout` (tools/schema/plan.yaml, tools/schema/down.sql) auf
# den Zustand vor dem Kommando zurück, auch bei einem Fehlschlag. Der Exit-Code
# ist der des Kommandos. Eine vor dem Lauf lokal geänderte Datei bleibt in
# diesem Zustand; eine fehlende bleibt fehlend.
#
# Aufrufer sind die Läufe, die den Rollout nur als Vorbedingung brauchen
# (tools/schema/apply-rollout.sh, die Integrations-, Bench- und
# Beispiel-Skripte); der Betrieb ruft `make schema-rollout` direkt und behält
# Report und Rollback-Artefakt (harness/targets/schema-rollout.md).
#
# Grenzen: zwei gleichzeitig laufende Aufrufe teilen dieselben zwei Dateien —
# der zweite sichert den Zwischenstand des ersten als „Zustand vor dem Lauf“;
# ein SIGKILL lässt die Erzeugnisse verändert (kein trap fängt ihn).
#
# Aufruf: rollout-restore.sh <Kommando> [Argumente …]
set -uo pipefail
if [ $# -eq 0 ]; then
  echo "Aufruf: rollout-restore.sh <Kommando> [Argumente …]" >&2
  exit 2
fi
cd "$(git rev-parse --show-toplevel)" || exit 2

artefacts=(tools/schema/plan.yaml tools/schema/down.sql)
backup=$(mktemp -d "${TMPDIR:-/tmp}/rollout-restore.XXXXXX") || exit 2
for f in "${artefacts[@]}"; do
  if [ -e "$f" ]; then
    cp -p "$f" "$backup/$(basename "$f")" || exit 2
  else
    : >"$backup/$(basename "$f").fehlt"
  fi
done

restore() {
  local f b
  for f in "${artefacts[@]}"; do
    b="$backup/$(basename "$f")"
    if [ -e "$b" ]; then
      cp -p "$b" "$f"
    else
      rm -f "$f"
    fi
  done
  rm -rf "$backup"
}
trap restore EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

"$@"
