#!/usr/bin/env bash
# run-command-guard-tests.sh — Tabellentest gegen den PreToolUse-Guard
# .claude/hooks/pretooluse-command-guard.sh (Vertrag: Kopfkommentar des Guards,
# Grenz-Zeile: harness/conventions/MR-003-guard-inplace-textwerkzeug.md).
# Je Fall ein Kommandostring, der als Hook-JSON auf die stdin des Guards geht;
# erwartet ist ein Block (Ausgabe `"decision": "block"`, Exit 0, gueltiges JSON,
# Begruendungstext der Klasse pkg | inplace | interp) oder der Pass-Fall (keine
# Ausgabe, Exit 0, leeres stderr). Gruppen: Bestandsregeln (Paketmanager,
# Praefixe, Sub-Shell-Rekursion, Tiefe, fail-closed) · in-place Formen von sed,
# perl, awk/gawk je Position (Kopf, nach &&, Praefix, absoluter Pfad, bash -c,
# find -exec) · Nicht-Treffer je Nachbarform · Host-Interpreter auf Repo-Pfaden
# · benannte Falsch-Positiv-Raender (erwarteter Block) · benannte Grenzen
# (erwarteter Pass, was der Guard nicht liest).
# Der Guard laeuft in einem Wegwerf-Repo im Temp-Verzeichnis (oberste Ebene
# docs/, internal/, Makefile, .claude/); der Test schreibt nur dorthin. Der
# Pruefling ist per GUARD uebersteuerbar (Mutationslaeufe gegen eine Kopie).
# Host-Werkzeuge: bash, awk, mktemp, grep, cp, mkdir, ln.
set -uo pipefail
repo=$(git rev-parse --show-toplevel)
guard_src=${GUARD:-$repo/.claude/hooks/pretooluse-command-guard.sh}
[ -f "$guard_src" ] || { echo "run-command-guard-tests: GUARD $guard_src nicht lesbar" >&2; exit 2; }

tmp=$(mktemp -d "${TMPDIR:-/tmp}/command-guard-test.XXXXXX")
trap 'rm -rf "$tmp"' EXIT
root="$tmp/repo"
mkdir -p "$root/.claude/hooks" "$root/tools/harness" "$root/docs" "$root/internal" "$tmp/noawk"
cp "$guard_src" "$root/.claude/hooks/pretooluse-command-guard.sh"
cp "$repo/tools/harness/extract-command.awk" "$root/tools/harness/extract-command.awk"
printf 'all:\n' >"$root/Makefile"
guard="$root/.claude/hooks/pretooluse-command-guard.sh"
# PATH ohne awk: nur die Werkzeuge, die der Guard vor der awk-Pruefung braucht.
ln -s "$(command -v cat)" "$tmp/noawk/cat"
ln -s "$(command -v dirname)" "$tmp/noawk/dirname"
bash_bin=$(command -v bash)

fail=0
n=0
out=""
err=""
rc=0

json_escape() {
  local s=$1
  s=${s//\\/\\\\}
  s=${s//\"/\\\"}
  s=${s//$'\n'/\\n}
  s=${s//$'\t'/\\t}
  printf '%s' "$s"
}
run_raw() { # <Hook-JSON roh>
  out=$(printf '%s' "$1" | "$bash_bin" "$guard" 2>"$tmp/err")
  rc=$?
  err=$(cat "$tmp/err")
}
run_guard() { # <Kommandostring>
  run_raw "{\"tool_input\":{\"command\":\"$(json_escape "$1")\"}}"
}
json_ok() { # die Block-Ausgabe hat die vier Zeilen des Bestands und einen Grund ohne " und \
  printf '%s\n' "$out" | awk '
    NR == 1 && $0 != "{" { bad = 1 }
    NR == 2 && $0 != "  \"decision\": \"block\"," { bad = 1 }
    NR == 3 && $0 !~ /^  "reason": "[^"\\]*"$/ { bad = 1 }
    NR == 4 && $0 != "}" { bad = 1 }
    END { exit (NR != 4 || bad) }'
}
reason_pattern() {
  case $1 in
    inplace) printf '%s' 'In-place text tools' ;;
    interp) printf '%s' 'Host python/perl on repo paths' ;;
    pkg) printf '%s' 'host package managers' ;;
  esac
}
report() { # <Fallname> <Meldung>
  echo "FEHLER: $1 — $2" >&2
  printf 'Ausgabe: %s\nstderr: %s\n' "$out" "$err" >&2
  fail=1
}
judge_block() { # <Klasse> <Fallname>
  if [ "$rc" -ne 0 ]; then
    report "$2" "Exit $rc, erwartet 0"
  elif ! printf '%s' "$out" | grep -qF '"decision": "block"'; then
    report "$2" "kein Block"
  elif ! printf '%s' "$out" | grep -qF -- "$(reason_pattern "$1")"; then
    report "$2" "Begründung der Klasse $1 fehlt"
  elif ! json_ok; then
    report "$2" "Block-Ausgabe ist kein gültiges JSON der Bestandsform"
  fi
}
block() { # <Klasse> <Fallname> <Kommandostring>
  n=$((n + 1))
  run_guard "$3"
  judge_block "$1" "$2"
}
pass() { # <Fallname> <Kommandostring>
  n=$((n + 1))
  run_guard "$2"
  if [ "$rc" -ne 0 ] || [ -n "$out" ] || [ -n "$err" ]; then
    report "$1" "Pass-Fall erwartet (Exit 0, keine Ausgabe, kein stderr)"
  fi
}

# --- Bestandsregeln --------------------------------------------------------
block pkg "Bestand: pip am Kopf" 'pip install x'
block pkg "Bestand: apt-get am Kopf" 'apt-get install y'
block pkg "Bestand: npm am Kopf" 'npm i'
block pkg "Bestand: absoluter Pfad des Paketmanagers" '/usr/bin/pip3 install x'
block pkg "Bestand: sudo-Präfix" 'sudo pip install x'
block pkg "Bestand: Zuweisungs-Präfix" 'FOO=1 pip install x'
block pkg "Bestand: nach &&" 'cd docs && pip install x'
block pkg "Bestand: Brace-Group" '{ pip install x; }'
block pkg "Bestand: bash -c" 'bash -c "pip install x"'
block pkg "Bestand: bash -lc (Flag-Bündel)" 'bash -lc "pip install x"'
block pkg "Bestand: Option hinter xargs" 'ls | xargs -n1 pip'
block pkg "Bestand: Option hinter env" 'env -i pip install x'
block pkg "Bestand: Tiefe 4 der sh -c-Rekursion" 'sh -c sh -c sh -c sh -c ls'
pass "Bestand: Tiefe 3 der sh -c-Rekursion" 'sh -c sh -c sh -c ls'
pass "Bestand: Paketmanager-Name in einem Argument" 'git commit -m "fix pip handling"'
pass "Bestand: Paketmanager-Name als echo-Argument" 'echo pip'
pass "Bestand: make gates" 'make gates'
pass "Bestand: git commit -F Datei" 'git commit -F /tmp/scratch/msg.txt'
pass "Bestand: bash -c ohne Treffer" 'bash -c "ls -la"'
n=$((n + 1))
run_raw '{"tool_input":{"command":"pip'
judge_block pkg "Bestand: abgeschnittenes JSON (fail-closed)"
n=$((n + 1))
run_raw ''
judge_block pkg "Bestand: leere Eingabe (fail-closed)"
n=$((n + 1))
run_raw '{"tool_input":{"command":"ls \u0041"}}'
judge_block pkg "Bestand: \\u-Escape im Befehl (fail-closed)"
n=$((n + 1))
out=$(printf '%s' '{"tool_input":{"command":"ls"}}' | PATH="$tmp/noawk" "$bash_bin" "$guard" 2>"$tmp/err")
rc=$?
err=$(cat "$tmp/err")
judge_block pkg "Bestand: fehlendes awk (fail-closed)"

# --- Liefer-Punkt 1: in-place Formen, je Form ------------------------------
block inplace "sed -i" 'sed -i s/a/b/ f'
block inplace "sed -i.bak" 'sed -i.bak s/a/b/ f'
block inplace "sed -ni (Bündel)" 'sed -ni s/a/b/p f'
block inplace "sed -Ei (Bündel)" 'sed -Ei s/a/b/ f'
block inplace "sed -i mit -e" 'sed -i -e s/a/b/ f'
block inplace "sed -n -i (Flag nicht am ersten Token)" 'sed -n -i s/a/b/p f'
block inplace "sed --in-place" 'sed --in-place s/a/b/ f'
block inplace "sed --in-place=.bak" 'sed --in-place=.bak s/a/b/ f'
block inplace "perl -i" 'perl -i -pe s/a/b/ f'
block inplace "perl -pi" 'perl -pi -e s/a/b/ f'
block inplace "perl -i.bak" 'perl -i.bak -pe s/a/b/ f'
block inplace "perl -0777pi" 'perl -0777pi -e s/a/b/ f'
block inplace "gawk -i inplace" "gawk -i inplace '{print}' f"
block inplace "awk -i inplace" "awk -i inplace '{print}' f"
block inplace "gawk -iinplace" "gawk -iinplace '{print}' f"
block inplace "gawk --include=inplace" "gawk --include=inplace '{print}' f"
block inplace "gawk --include inplace" "gawk --include inplace '{print}' f"

# --- Liefer-Punkt 1: je Position -------------------------------------------
block inplace "Position: nach &&" 'cd docs && sed -i s/a/b/ f'
block inplace "Position: nach Pipe, hinter xargs" 'ls | xargs sed -i s/a/b/'
block inplace "Position: xargs mit Option" 'ls | xargs -r sed -i s/a/b/'
block inplace "Position: xargs mit -I{}" 'ls | xargs -I{} sed -i s/a/b/ {}'
block inplace "Position: sudo" 'sudo sed -i s/a/b/ f'
block inplace "Position: Zuweisungs-Präfix" 'FOO=1 sed -i s/a/b/ f'
block inplace "Position: env mit Zuweisung" 'env FOO=1 perl -pi -e s/a/b/ f'
block inplace "Position: absoluter Pfad des Werkzeugs" '/usr/bin/sed -i s/a/b/ f'
block inplace "Position: bash -c" 'bash -c "sed -i s/a/b/ f"'
block inplace "Position: find -exec ;" 'find . -name x -exec sed -i s/a/b/ {} \;'
block inplace "Position: find -exec +" 'find . -name x -exec sed -i s/a/b/ {} +'
block inplace "Position: find -execdir perl" 'find . -execdir perl -pi -e s/a/b/ {} +'
block inplace "Position: find -ok gawk" 'find . -ok gawk -i inplace {print} {} \;'
block inplace "Ziel unabhängig: Scratchpad-Kopie" 'sed -i s/a/b/ /tmp/scratch/copy.go'
block inplace "Ziel unabhängig: nicht vorhandene Datei" 'sed -i s/a/b/ /nicht-vorhanden-scratch'
block inplace "Ziel unabhängig: Eingabe-Umleitung" 'xargs sed -i s/a/b/ < /tmp/scratch/liste'

# --- Liefer-Punkt 1: Nicht-Treffer je Nachbarform --------------------------
pass "sed -n" 'sed -n 1p Makefile'
pass "sed -E" 'sed -E s/a/b/ f'
pass "sed ohne Flag" 'sed s/a/b/ f'
pass "sed --version" 'sed --version'
pass "sed -n -e" 'sed -n -e s/a/b/p f'
pass "sed Skript mit -i im Text" 'sed -n /-i/p f'
pass "sed auf Scratchpad-Kopie nach stdout" 'sed s/a/b/ f > /tmp/scratch/copy.go'
pass "awk lesend" "awk '{print \$1}' f"
pass "awk -F" "awk -F, '{print \$1}' f"
pass "awk -v i=1" "awk -v i=1 '{print i}' f"
pass "gawk -i mit anderer Bibliothek" "gawk -i /tmp/lib.awk '{print}' f"
pass "gawk --include mit anderer Bibliothek" "gawk --include /tmp/lib.awk '{print}' f"
pass "perl -MList::Util" "perl -MList::Util -e 'print 1'"
pass "perl -e" "perl -e 'print 1'"
pass "perl -ne auf Scratchpad-Datei" "perl -ne print /tmp/x"
pass "perl -lne" "perl -lne print /tmp/x"
pass "perl -I" "perl -I lib -e 1"
pass "grep -i" 'grep -i x f'
pass "cp -i" 'cp -i a b'
pass "git grep -E mit Alternation (Text, roh gelesen)" "git grep -E 'sed -i|perl -pi'"
pass "git grep -E mit Alternation, perl vor sed" "git grep -E 'perl -pi|sed -i'"
pass "gawk -i inplace im Muster (Text, roh gelesen)" "grep -E 'a|gawk -i inplace' f"
pass "find -exec sed -n" 'find . -name x -exec sed -n 1p {} +'
pass "find -exec grep -i" 'find . -exec grep -i x {} +'
pass "find: -iname hinter dem -exec-Ende" 'find . -exec sed -n 1p {} + -iname x'
pass "git log --grep" "git log --grep='sed -i'"
pass "bash -c mit sed -n" 'bash -c "sed -n 1p f"'
pass "xargs grep mit Option" 'ls | xargs -r grep -nE x'

# --- Liefer-Punkt 2: Host-Interpreter auf Repo-Pfaden ----------------------
block interp "python3 tools/x.py" 'python3 tools/x.py'
block interp "python tools/x.py" 'python tools/x.py'
block interp "python3.12 tools/x.py" 'python3.12 tools/x.py'
block interp "python3 ./tools/x.py" 'python3 ./tools/x.py'
block interp "python3 mit Punkt-Verzeichnis" 'python3 .claude/hooks/x.py'
block interp "python3 Heredoc mit docs/ im Text" $'python3 - <<\'EOF\'\nopen(\'docs/a.md\', \'w\')\nEOF'
block interp "python3 -c mit Makefile" "python3 -c \"open('Makefile', 'w')\""
block interp "perl -e mit Makefile" "perl -e 'unlink q(Makefile)'"
block interp "python3 mit Absolutpfad der Repo-Wurzel" "python3 $root/docs/a.md"
block interp "python3 in bash -c" 'bash -c "python3 tools/x.py"'
block interp "perl -ne auf Datei der obersten Ebene" 'perl -ne print Makefile'
pass "python3 --version" 'python3 --version'
pass "python3 -c ohne Repo-Pfad" "python3 -c 'print(1)'"
pass "python3 auf Scratchpad-Kopie" 'python3 /tmp/scratch/mutate.py /tmp/scratch/copy.go'
pass "python3: docs/ hinter einem Pfadzeichen" 'python3 /tmp/x/docs/a.py'
pass "python3: internal/ in Scratchpad-Kopie" 'python3 /tmp/scratch/mutate.py /tmp/scratch/internal/copy.go'
pass "python3: Datei-Name hinter Pfadzeichen" 'python3 /tmp/x/Makefile'
pass "python3: Verzeichnis-Name ohne Schrägstrich" "python3 -c 'print(\"docs\")'"
pass "python3: Datei-Name als Wortanfang" "python3 -c 'x = Makefile2'"
pass "python3: Datei-Name mit Endung" "python3 -c 'x = Makefile.bak'"
pass "perl -ne auf Scratchpad-Datei ohne Repo-Namen" "perl -ne 'print' /tmp/x"
pass "Kopf ist nicht python: echo" 'echo python3 tools/x.py'
pass "Kopf ist nicht python: grep" 'grep python3 tools/x.sh'
pass "Repo-Pfad ohne Interpreter" 'git grep -n docs/ Makefile'

# --- benannte Falsch-Positiv-Raender (der Guard blockt, weil quote-blind) --
block inplace "Rand: Trenner im Muster, Leerzeichen vor dem Anführungszeichen" "grep -E 'a|sed -i ' f"
block inplace "Rand: Heredoc-Zeile mit sed -i am Kopf" $'cat <<EOF\nsed -i s/a/b/ f\nEOF'
block inplace "Rand: sed -i nach & in einer Commit-Message" 'git commit -m "a & sed -i b"'
block interp "Rand: python3 nennt einen Repo-Namen, schreibt nichts" "python3 -c \"print('tools/x')\""
block interp "Rand: cd im selben Kommando" 'cd /tmp/x && python3 tools/x.py'

# --- benannte Grenzen (der Guard liest es nicht: erwarteter Pass) ----------
pass "Grenze: Umleitung und mv" 'sed s/a/b/ f > /tmp/scratch/tmp && mv /tmp/scratch/tmp f'
pass "Grenze: tee" 'echo x | tee f'
pass "Grenze: dd of=" 'dd if=/dev/null of=f'
pass "Grenze: Skript, das ein Interpreter liest" 'python3 /tmp/scratch/x.py'
pass "Grenze: bash skript.sh" 'bash /tmp/scratch/skript.sh'
pass "Grenze: cd, dann relativer Name ohne Repo-Namen" 'cd /nirgends && python3 x.py'
pass "Grenze: anderes in-place-fähiges Werkzeug (patch)" 'patch f < /tmp/scratch/x.diff'
pass "Grenze: anderes in-place-fähiges Werkzeug (ruby -i)" 'ruby -i -pe 1 f'
pass "Grenze: Flag in Anführungszeichen" "sed '-i' s/a/b/ f"
pass "Grenze: Optionswert hinter einem Wrapper" 'sudo -u x sed -i s/a/b/ f'

if [ "$fail" -eq 0 ]; then
  echo "run-command-guard-tests: alle $n Fälle bestanden"
else
  echo "run-command-guard-tests: mindestens ein Fall von $n gescheitert" >&2
  exit 1
fi
