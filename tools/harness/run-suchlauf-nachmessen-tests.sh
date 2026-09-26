#!/usr/bin/env bash
# run-suchlauf-nachmessen-tests.sh — Tabellentest gegen suchlauf-nachmessen.sh
# (harness/sensors/suchlauf-nachmessen.md): fünf Fälle gegen ein Wegwerf-Repo —
# stimmt · weicht ab · Selbstverweis ausgeschlossen · HEAD abgelehnt · kein Block.
# Netzlos, kein Docker (git und bash). Der Prüfling ist per TOOL übersteuerbar
# (Mutationsläufe gegen eine Kopie).
set -uo pipefail
repo=$(git rev-parse --show-toplevel)
tool=${TOOL:-$repo/tools/harness/suchlauf-nachmessen.sh}

tmp=$(mktemp -d "${TMPDIR:-/tmp}/suchlauf-nachmessen-test.XXXXXX")
trap 'rm -rf "$tmp"' EXIT
cd "$tmp" || exit 1

git init -q .
git config user.name test
git config user.email test@example.invalid
git config commit.gpgsign false

# Parent-Stand: zwei Treffer auf `alpha` in docs/a.md; der Plan liegt unter
# next/ und trägt das Suchwort selbst.
mkdir -p docs plan/next plan/in-progress
printf 'alpha\nalpha beta\ngamma\n' >docs/a.md
printf '# Plan\nBeschreibung alpha\n' >plan/next/slice-x.md
git add -A
git commit -q -m parent
parent=$(git rev-parse --short HEAD)

# Arbeitsbaum-Stand: eine dritte `alpha`-Zeile (`alpha gamma`), der Plan ist nach
# in-progress/ gewandert.
printf 'alpha\nalpha beta\nalpha gamma\n' >docs/a.md
git mv plan/next/slice-x.md plan/in-progress/slice-x.md

fail=0
out=""
rc=0

# write_plan <block-zeilen>: schreibt den Plan mit einem suchlauf-Block.
write_plan() {
  {
    printf '# Plan\nBeschreibung alpha\n\n```suchlauf\n'
    printf '%s\n' "$@"
    printf '```\n'
  } >plan/in-progress/slice-x.md
}

# run_tool: Ausgabe und Exit des Prüflings.
run_tool() {
  out=$(bash "$tool" plan/in-progress/slice-x.md 2>&1)
  rc=$?
}

expect() { # <Fallname> <erwarteter Exit> <Muster in der Ausgabe oder ->
  local name=$1 want_rc=$2 pattern=$3
  if [ "$rc" -ne "$want_rc" ]; then
    echo "FEHLER: $name — Exit $rc, erwartet $want_rc; Ausgabe:" >&2
    printf '%s\n' "$out" >&2
    fail=1
  elif [ "$pattern" != "-" ] && ! printf '%s\n' "$out" | grep -qE -- "$pattern"; then
    echo "FEHLER: $name — Ausgabe trägt '$pattern' nicht:" >&2
    printf '%s\n' "$out" >&2
    fail=1
  fi
}

# Fall 1 — stimmt: Parent 2 Zeilen, Arbeitsbaum 3; ein Muster mit Leerzeichen in
# Anführungszeichen; ein Pathspec, der den Suchraum einschränkt.
write_plan "$parent 2 -E 'alpha'" \
  "diff 3 -E 'alpha'" \
  "diff 1 'alpha gamma' -- docs" \
  "diff 0 'alpha' -- docs/none"
run_tool
expect "stimmt" 0 'soll=3 ist=3'

# Fall 2 — weicht ab: das Soll des Parent-Stands und das des Arbeitsbaums je
# um eins daneben.
write_plan "$parent 3 -E 'alpha'"
run_tool
expect "weicht ab (Parent)" 1 'ABWEICHUNG  soll=3 ist=2'
write_plan "diff 2 -E 'alpha'"
run_tool
expect "weicht ab (diff)" 1 'ABWEICHUNG  soll=2 ist=3'

# Fall 3 — Selbstverweis ausgeschlossen: der Plan trägt `alpha` in Beschreibung
# und Block-Zeile, am Parent-Stand unter einem anderen Verzeichnis (next/); die
# Zahlen zählen ihn nicht.
write_plan "$parent 2 -E 'alpha'" "diff 3 -E 'alpha'"
run_tool
expect "Selbstverweis" 0 'soll=2 ist=2'

# Fall 4 — HEAD (und jeder Name statt einer Commit-Kennung) wird abgelehnt,
# bevor ein Befehl läuft.
write_plan "$parent 2 -E 'alpha'" "HEAD 2 -E 'alpha'"
run_tool
expect "HEAD abgelehnt" 2 'HEAD'
if printf '%s\n' "$out" | grep -q '^OK'; then
  echo "FEHLER: HEAD abgelehnt — es lief trotzdem eine Messung:" >&2
  printf '%s\n' "$out" >&2
  fail=1
fi
write_plan "HEAD~1 2 -E 'alpha'"
run_tool
expect "HEAD~1 abgelehnt" 2 'HEAD~1'

# Fall 5 — kein Block: ein Plan ohne suchlauf-Block, mit nur einem text-Block und
# mit leerem suchlauf-Block endet mit Exit 2.
printf '# Plan\n| Befehl |\n|---|\n| `git grep alpha` |\n' >plan/in-progress/slice-x.md
run_tool
expect "kein Block (Tabelle)" 2 'keine Zeile'
printf '# Plan\n\n```text\n%s 2 -E alpha\n```\n' "$parent" >plan/in-progress/slice-x.md
run_tool
expect "kein Block (text)" 2 'keine Zeile'
printf '# Plan\n\n```suchlauf\n```\n' >plan/in-progress/slice-x.md
run_tool
expect "kein Block (leer)" 2 'keine Zeile'

if [ "$fail" -eq 0 ]; then
  echo "run-suchlauf-nachmessen-tests: alle Fälle bestanden"
else
  echo "run-suchlauf-nachmessen-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
