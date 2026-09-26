#!/usr/bin/env bash
# run-suchlauf-nachmessen-tests.sh — Tabellentest gegen suchlauf-nachmessen.sh
# (harness/sensors/suchlauf-nachmessen.md): elf Fälle gegen ein Wegwerf-Repo —
# stimmt · weicht ab · Selbstverweis ausgeschlossen · HEAD abgelehnt · kein Block ·
# Commit-Stand mit Pathspec · erlaubte Optionen · Optionen und Pathspec-Magic
# außerhalb der Allow-List (die Marker-Datei eines eingeschleusten Kommandos
# entsteht nicht) · git grep-Fehler · nicht geschlossener Block · Zeilenform
# (Soll, Stand, Muster, leeres Argument).
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
mkdir -p docs other plan/next plan/in-progress
printf 'alpha\nalpha beta\ngamma\n' >docs/a.md
printf 'omega\n' >docs/d.md
printf 'omega\n' >other/c.md
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

# Fall 6 — Commit-Stand mit Pathspec: der Suchraum des Parent-Stands folgt dem
# Pathspec (`omega` steht in docs/ und other/, gezählt wird nur docs/), und ein
# Ausschluss-Pathspec zieht seine Datei ab.
write_plan "$parent 1 omega -- docs" \
  "diff 1 omega -- docs" \
  "$parent 2 omega" \
  "$parent 1 omega -- ':!other'" \
  "$parent 1 omega -- ':(exclude,glob)other/**'"
run_tool
expect "Commit-Stand mit Pathspec" 0 'soll=1 ist=1'

# Fall 7 — erlaubte Optionen: Kurz- und Langformen, -e mit Verknüpfung.
write_plan "diff 1 -i -F 'ALPHA GAMMA'" \
  "diff 1 -e alpha --and -e gamma" \
  "diff 1 --ignore-case --word-regexp 'gamma'"
run_tool
expect "erlaubte Optionen" 0 'soll=1 ist=1'

# Fall 8 — Optionen außerhalb der Allow-List werden abgelehnt, bevor ein Befehl
# läuft; -O und --open-files-in-pager führten sonst ein Kommando aus, dessen
# Marker-Datei entstünde.
marker=$tmp/marker-ausgefuehrt
reject() { # <Fallname> <Muster der Meldung> <Plan-Zeile>
  rm -f "$marker"
  write_plan "$3"
  run_tool
  expect "$1" 2 "$2"
  if printf '%s\n' "$out" | grep -q '^OK'; then
    echo "FEHLER: $1 — es lief trotzdem eine Messung:" >&2
    printf '%s\n' "$out" >&2
    fail=1
  fi
  if [ -e "$marker" ]; then
    echo "FEHLER: $1 — die Plan-Zeile hat ein Kommando ausgeführt" >&2
    fail=1
  fi
}
reject "-O abgelehnt" "nicht erlaubt" "diff 0 -O'touch $marker;true' -e alpha"
reject "--open-files-in-pager abgelehnt" "nicht erlaubt" "diff 0 --open-files-in-pager='touch $marker;true' -e alpha"
reject "-O ohne Kommando abgelehnt" "nicht erlaubt" "diff 0 -O -e alpha"
reject "gebündelt mit -O abgelehnt" "nicht erlaubt" "diff 0 -inO -e alpha"
reject "--no-index abgelehnt" "nicht erlaubt" "diff 0 --no-index -e alpha"
reject "-f abgelehnt" "nicht erlaubt" "diff 0 -f docs/a.md"
reject "-e ohne Muster abgelehnt" "ohne Muster" "diff 0 -i -e"
reject "zwei Suchwörter abgelehnt" "mehr als ein Suchwort" "diff 0 alpha gamma"
reject "Suchwort neben -e abgelehnt" "neben -e" "diff 0 -e alpha gamma"
reject "Pathspec-Magic abgelehnt" "Pathspec" "diff 0 -e alpha -- ':(attr:foo)docs'"
reject "Pathspec :/ abgelehnt" "Pathspec" "diff 0 -e alpha -- ':/docs'"

# Fall 9 — ein Fehler von git grep (ungültiges Muster, Exit über 1) endet mit
# Exit 2 und zählt keine Trefferzeile.
write_plan "diff 0 -E 'a('"
run_tool
expect "git grep-Fehler" 2 'git grep endete mit Exit'

# Fall 10 — nicht geschlossener Block: Exit 2, bevor ein Befehl läuft.
printf '# Plan\n\n```suchlauf\ndiff 3 -E alpha\n' >plan/in-progress/slice-x.md
run_tool
expect "Block nicht geschlossen" 2 'nicht geschlossen'
if printf '%s\n' "$out" | grep -q '^OK'; then
  echo "FEHLER: Block nicht geschlossen — es lief trotzdem eine Messung:" >&2
  fail=1
fi

# Fall 11 — Zeilenform: Soll ohne Zahl, unbekannte Stand-Kennung, Zeile ohne
# Suchmuster, leeres Argument; jede endet mit Exit 2 und keiner Messung.
reject "Soll ohne Zahl" "keine Zahl" "diff drei -E 'alpha'"
reject "Soll mit Vorzeichen" "keine Zahl" "diff -1 -E 'alpha'"
reject "unbekannte Stand-Kennung" "kein Commit" "abcdef0 1 -E 'alpha'"
reject "Stand zu kurz" "weder eine Commit-Kennung" "abc 1 -E 'alpha'"
reject "Zeile ohne Suchmuster" "ohne Suchmuster" "diff 2 -- docs"
reject "Zeile ohne Argumente" "ohne <Stand>" "diff 2"
reject "leeres Argument" "leeres Argument" "diff 1 -e '' -- docs"
reject "Anführungszeichen offen" "nicht geschlossen" "diff 1 -e 'alpha"

if [ "$fail" -eq 0 ]; then
  echo "run-suchlauf-nachmessen-tests: alle Fälle bestanden"
else
  echo "run-suchlauf-nachmessen-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
