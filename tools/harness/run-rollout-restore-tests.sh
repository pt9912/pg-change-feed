#!/usr/bin/env bash
# run-rollout-restore-tests.sh — Tabellentest gegen tools/schema/rollout-restore.sh:
# die Rücknahme von plan.yaml und down.sql (Kommando schreibt · Kommando
# scheitert · lokale Änderung bleibt · fehlende Datei bleibt fehlend), dazu die
# Prüfung, dass jeder Aufrufer von `make schema-rollout` durch das Skript geht.
# Netzlos, kein Docker (git und bash). Der Prüfling ist per TOOL übersteuerbar
# (Mutationsläufe gegen eine Kopie).
set -uo pipefail
repo=$(git rev-parse --show-toplevel)
tool=${TOOL:-$repo/tools/schema/rollout-restore.sh}

fail=0
err() { echo "FEHLER: $*" >&2; fail=1; }

# Aufrufer-Prüfung im echten Repo: jede Shell-Zeile unter tools/ und examples/,
# die `make … schema-rollout` aufruft (Kommentar- und echo-Zeilen ausgenommen),
# trägt rollout-restore.sh. Ausgenommen sind der Guard-Test (sichert und stellt
# beide Dateien selbst wieder her) und dieses Skript.
callers=$(cd "$repo" && git grep -n -E '\bmake\b.*schema-rollout' -- 'tools/*.sh' 'examples/*.sh' \
  ':!tools/harness/run-schema-rollout-guard-test.sh' ':!tools/harness/run-rollout-restore-tests.sh' \
  | grep -vE '^[^:]+:[0-9]+:[[:space:]]*(#|echo )' || true)
if [ -z "$callers" ]; then
  err "Aufrufer-Prüfung: keine Aufrufer von make schema-rollout gefunden (Suchmuster trägt nicht)"
fi
while IFS= read -r line; do
  [ -n "$line" ] || continue
  case $line in
    *rollout-restore.sh*) ;;
    *) err "Aufrufer ohne rollout-restore.sh: $line" ;;
  esac
done <<<"$callers"

tmp=$(mktemp -d "${TMPDIR:-/tmp}/rollout-restore-test.XXXXXX")
trap 'rm -rf "$tmp"' EXIT
cd "$tmp" || exit 1
git init -q .
git config user.name test
git config user.email test@example.invalid
git config commit.gpgsign false
mkdir -p tools/schema
printf 'plan-alt\n' >tools/schema/plan.yaml
printf 'down-alt\n' >tools/schema/down.sql
git add -A
git commit -q -m init

# Das Kommando ersetzt beide Erzeugnisse, wie ein Rollout es tut.
writer='printf "plan-neu\n" >tools/schema/plan.yaml; printf "down-neu\n" >tools/schema/down.sql'

# Fall 1 — das Kommando schreibt beide Dateien: danach ist der Baum sauber.
bash "$tool" bash -c "$writer"
rc=$?
[ "$rc" -eq 0 ] || err "Fall 1: Exit $rc, erwartet 0"
[ -z "$(git status --short tools/schema)" ] || err "Fall 1: git status nicht sauber: $(git status --short tools/schema)"

# Fall 2 — das Kommando scheitert nach dem Schreiben: der Exit-Code bleibt,
# der Baum ist sauber.
bash "$tool" bash -c "$writer; exit 7"
rc=$?
[ "$rc" -eq 7 ] || err "Fall 2: Exit $rc, erwartet 7"
[ -z "$(git status --short tools/schema)" ] || err "Fall 2: git status nicht sauber: $(git status --short tools/schema)"

# Fall 3 — eine vor dem Lauf lokal geänderte Datei bleibt in diesem Zustand.
printf 'plan-lokal\n' >tools/schema/plan.yaml
bash "$tool" bash -c "$writer"
[ "$(cat tools/schema/plan.yaml)" = "plan-lokal" ] || err "Fall 3: plan.yaml trägt '$(cat tools/schema/plan.yaml)', erwartet 'plan-lokal'"
[ "$(cat tools/schema/down.sql)" = "down-alt" ] || err "Fall 3: down.sql trägt '$(cat tools/schema/down.sql)', erwartet 'down-alt'"
git checkout -q -- tools/schema

# Fall 4 — eine vor dem Lauf fehlende Datei fehlt danach wieder.
rm tools/schema/down.sql
bash "$tool" bash -c "$writer"
[ ! -e tools/schema/down.sql ] || err "Fall 4: down.sql existiert nach dem Lauf"
[ "$(cat tools/schema/plan.yaml)" = "plan-alt" ] || err "Fall 4: plan.yaml trägt '$(cat tools/schema/plan.yaml)', erwartet 'plan-alt'"
git checkout -q -- tools/schema

# Fall 5 — ohne Kommando: Exit 2.
bash "$tool" >/dev/null 2>&1
rc=$?
[ "$rc" -eq 2 ] || err "Fall 5: Exit $rc, erwartet 2"

if [ "$fail" -eq 0 ]; then
  echo "run-rollout-restore-tests: alle Fälle bestanden"
else
  echo "run-rollout-restore-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
