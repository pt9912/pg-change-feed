#!/usr/bin/env bash
# run-command-guard-tests.sh — Tabellentest gegen den PreToolUse-Guard
# .claude/hooks/pretooluse-command-guard.sh (Vertrag: Kopfkommentar des Guards,
# Grenz-Zeile: harness/conventions/MR-003-guard-inplace-textwerkzeug.md).
# Je Fall ein Kommandostring, der als Hook-JSON auf die stdin des Guards geht;
# erwartet ist ein Block (Ausgabe `"decision": "block"`, Exit 0, gueltiges JSON,
# Begruendungstext der Klasse pkg | inplace | interp) oder der Pass-Fall (keine
# Ausgabe, Exit 0, leeres stderr). Gruppen: Bestandsregeln (Paketmanager,
# Praefixe, Sub-Shell-Rekursion, Tiefe, fail-closed) · Kopf-Erkennung hinter
# Wrapper-Optionen und Shell-Schluesselwoertern · in-place Formen von sed, perl,
# awk/gawk je Form, je Mitglied der Zeichenklassen und je Position (Kopf, nach
# Trenner, Praefix, absoluter Pfad, bash -c, eval, find -exec) · Nicht-Treffer je
# Nachbarform · Anfuehrungszeichen-Lesung (ein Trenner im Argument) ·
# Host-Interpreter auf Repo-Pfaden samt Pfadzeichen-Klasse · benannte
# Falsch-Positiv-Raender (erwarteter Block) · benannte Grenzen (erwarteter Pass,
# was der Guard nicht liest).
# Der Guard laeuft in einem Wegwerf-Repo im Temp-Verzeichnis (oberste Ebene
# docs/, internal/, Makefile, .claude/); der Test schreibt nur dorthin. Der
# Pruefling ist per GUARD, der Anfuehrungszeichen-Maskierer per MASKER
# uebersteuerbar (Mutationslaeufe gegen Kopien).
# Host-Werkzeuge: bash, awk, mktemp, grep, cp, mkdir, ln.
set -uo pipefail
repo=$(git rev-parse --show-toplevel)
guard_src=${GUARD:-$repo/.claude/hooks/pretooluse-command-guard.sh}
masker_src=${MASKER:-$repo/tools/harness/mask-quotes.awk}
[ -f "$guard_src" ] || { echo "run-command-guard-tests: GUARD $guard_src nicht lesbar" >&2; exit 2; }
[ -f "$masker_src" ] || { echo "run-command-guard-tests: MASKER $masker_src nicht lesbar" >&2; exit 2; }

tmp=$(mktemp -d "${TMPDIR:-/tmp}/command-guard-test.XXXXXX")
trap 'rm -rf "$tmp"' EXIT
root="$tmp/repo"
mkdir -p "$root/.claude/hooks" "$root/tools/harness" "$root/docs" "$root/internal" "$tmp/noawk"
cp "$guard_src" "$root/.claude/hooks/pretooluse-command-guard.sh"
cp "$repo/tools/harness/extract-command.awk" "$root/tools/harness/extract-command.awk"
cp "$masker_src" "$root/tools/harness/mask-quotes.awk"
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
pass "Segmentierung: ls (Gegenstück zum unlesbaren Maskierer)" 'ls'
mv "$root/tools/harness/mask-quotes.awk" "$tmp/mask-quotes.awk.weg"
n=$((n + 1))
run_guard 'ls'
mv "$tmp/mask-quotes.awk.weg" "$root/tools/harness/mask-quotes.awk"
judge_block pkg "Segmentierung: Maskierer nicht lesbar (fail-closed)"

# --- Kopf-Erkennung hinter Wrapper-Optionen und Schluesselwoertern ---------
# (der Guard vor der Erweiterung blockte diese Formen nicht)
block pkg "Wrapper-Option: xargs -n1" 'ls | xargs -n1 pip'
block pkg "Wrapper-Option: env -i" 'env -i pip install x'
block pkg "Wrapper-Option: time -p" 'time -p pip install x'
block pkg "Wrapper-Option: command -p (ausführend)" 'command -p pip install x'
block pkg "Wrapper mit absolutem Pfad: /usr/bin/env" '/usr/bin/env pip install x'
block pkg "Schlüsselwort: for … do" 'for f in a b; do pip install x; done'
block pkg "Schlüsselwort: if … then" 'if true; then pip install x; fi'
block pkg "Schlüsselwort: !" '! pip install x'
pass "command -v zeigt an" 'command -v pip'
pass "command -V zeigt an" 'command -V npm'
pass "command -pv zeigt an (Bündel)" 'command -pv yarn'
pass "type zeigt an" 'type pip'
pass "which zeigt an" 'which pip'

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

# --- Liefer-Punkt 1: die Zeichenklassen der Bündel, je Mitglied ------------
# sed: -[nEsrzub]*i
block inplace "sed-Klasse: n" 'sed -ni s/a/b/ f'
block inplace "sed-Klasse: E (Großbuchstabe)" 'sed -Ei s/a/b/ f'
block inplace "sed-Klasse: s" 'sed -si s/a/b/ f'
block inplace "sed-Klasse: r" 'sed -ri s/a/b/ f'
block inplace "sed-Klasse: z" 'sed -zi s/a/b/ f'
block inplace "sed-Klasse: u" 'sed -ui s/a/b/ f'
block inplace "sed-Klasse: b" 'sed -bi s/a/b/ f'
block inplace "sed-Klasse: mehrere Mitglieder" 'sed -nEsrzubi s/a/b/ f'
pass "sed außerhalb der Klasse: -e trägt das Skript i" 'sed -ei f'
pass "sed außerhalb der Klasse: -f trägt die Skriptdatei i" 'sed -fi f'
pass "sed außerhalb der Klasse: -e klein, nicht -E" 'sed -e s/a/b/ f'
pass "sed außerhalb der Klasse: unbekannter Buchstabe" 'sed -xi s/a/b/ f'
# perl: -[0-7lanpsw]*i
block inplace "perl-Klasse: 0" 'perl -0pi -e s/a/b/ f'
block inplace "perl-Klasse: Ziffern 1 bis 7" 'perl -l0127pi -e s/a/b/ f'
block inplace "perl-Klasse: l" 'perl -li -pe s/a/b/ f'
block inplace "perl-Klasse: a" 'perl -ai -pe s/a/b/ f'
block inplace "perl-Klasse: n" 'perl -ni -e print f'
block inplace "perl-Klasse: p" 'perl -pi -e s/a/b/ f'
block inplace "perl-Klasse: s" 'perl -si -pe s/a/b/ f'
block inplace "perl-Klasse: w" 'perl -wpi -e s/a/b/ f'
block inplace "perl-Klasse: mehrere Mitglieder" 'perl -0lanpswi -e s/a/b/ f'
block inplace "perl: -i hinter dem -e-Argument" "perl -pe 's/a/b/' -i f"
block inplace "perl: -I mit Wert vor -pi" 'perl -I lib -pi -e s/a/b/ f'
pass "perl außerhalb der Klasse: -e trägt den Code i" 'perl -ei f'
pass "perl außerhalb der Klasse: -M trägt den Modulnamen" 'perl -Minteger -e 1'
pass "perl außerhalb der Klasse: -I trägt das Verzeichnis" 'perl -Ii -e 1'
pass "perl außerhalb der Klasse: -Mfeature vor -e" 'perl -Mfeature -e 1'
pass "perl: -i hinter dem Skriptnamen ist ein Skript-Argument" 'perl /tmp/x.pl -input a'
pass "perl: -pi hinter dem Skriptnamen ist ein Skript-Argument" 'perl /tmp/x.pl -pi a'

# --- Liefer-Punkt 1: je Position -------------------------------------------
block inplace "Position: nach &&" 'cd docs && sed -i s/a/b/ f'
block inplace "Position: nach Pipe, hinter xargs" 'ls | xargs sed -i s/a/b/'
block inplace "Position: xargs mit Option" 'ls | xargs -r sed -i s/a/b/'
block inplace "Position: xargs mit -I{}" 'ls | xargs -I{} sed -i s/a/b/ {}'
block inplace "Position: xargs mit Optionswert -n 1" 'ls | xargs -n 1 sed -i s/a/b/'
block inplace "Position: xargs mit Optionswert -I {}" 'ls | xargs -r -I {} sed -i s/a/b/ {}'
block inplace "Position: xargs mit Optionswert -I REPL" 'ls | xargs -I REPL sed -i s/a/b/ REPL'
block inplace "Position: xargs -0 -P4" 'ls | xargs -0 -P4 sed -i s/a/b/'
block inplace "Position: xargs mit Optionswert -P 4" 'ls | xargs -P 4 sed -i s/a/b/'
block inplace "Position: xargs mit Optionswert -L 1" 'ls | xargs -L 1 sed -i s/a/b/'
block inplace "Position: xargs mit Optionswert -s 99" 'ls | xargs -s 99 sed -i s/a/b/'
block inplace "Position: xargs mit Optionswert -d" 'ls | xargs -d x sed -i s/a/b/'
block inplace "Position: xargs mit Optionswert -a" 'xargs -a /tmp/scratch/liste sed -i s/a/b/'
block inplace "Position: xargs mit Optionswert -E" 'ls | xargs -E eof sed -i s/a/b/'
block inplace "Position: sudo -E (Option ohne Wert)" 'sudo -E sed -i s/a/b/ f'
block pkg "Position: xargs mit Optionswert, Paketmanager" 'ls | xargs -n 1 pip install'
block inplace "Position: find -okdir" 'find . -okdir sed -i s/a/b/ {} \;'
block inplace "bash -c: Zeilenumbruch im String" $'bash -c "ls\nsed -i s/a/b/ f"'
block inplace "Zeilenfortsetzung hinter dem Werkzeug" $'sed \\\n-i s/a/b/ f'
block inplace "perl: -M-Argument endet auf e, -pi folgt" 'perl -Mfeature -pi -e s/a/b/ f'
pass "sed: -- beendet die Optionen" 'sed -n -- 1p f'
block inplace "Position: sudo" 'sudo sed -i s/a/b/ f'
block inplace "Position: Zuweisungs-Präfix" 'FOO=1 sed -i s/a/b/ f'
block inplace "Position: env mit Zuweisung" 'env FOO=1 perl -pi -e s/a/b/ f'
block inplace "Position: absoluter Pfad des Werkzeugs" '/usr/bin/sed -i s/a/b/ f'
block inplace "Position: bash -c" 'bash -c "sed -i s/a/b/ f"'
block inplace "Position: find -exec ;" 'find . -name x -exec sed -i s/a/b/ {} \;'
block inplace "Position: find -exec +" 'find . -name x -exec sed -i s/a/b/ {} +'
block inplace "Position: find -execdir perl" 'find . -execdir perl -pi -e s/a/b/ {} +'
block inplace "Position: find -ok gawk" 'find . -ok gawk -i inplace {print} {} \;'
block inplace "Position: find -exec sh -c" "find . -name x -exec sh -c 'sed -i s/a/b/ {}' \\;"
block inplace "Position: find -exec bash -c mit Bündel" "find . -exec bash -lc 'sed -i s/a/b/ \"\$0\"' {} \\;"
block inplace "Position: command" 'command sed -i s/a/b/ f'
block inplace "Position: time mit Option" 'time -p sed -i s/a/b/ f'
block inplace "Position: nice mit Option" 'nice -n10 sed -i s/a/b/ f'
block inplace "Position: eval mit Anführungszeichen" 'eval "sed -i s/a/b/ f"'
block inplace "Position: eval ohne Anführungszeichen" 'eval sed -i s/a/b/ f'
block inplace "Position: Sub-Shell" '(sed -i s/a/b/ f)'
block inplace "Position: Kommando-Ersetzung" 'echo $(sed -i s/a/b/ f)'
block inplace "Position: Backticks" 'echo `sed -i s/a/b/ f`'
block inplace "Position: Kommando-Ersetzung in doppelten Anführungszeichen" 'echo "$(sed -i s/a/b/ f)"'
block inplace "Position: Backticks in doppelten Anführungszeichen" 'echo "`sed -i s/a/b/ f`"'
block inplace "Position: Brace-Group" '{ sed -i s/a/b/ f; }'
block inplace "Position: nach find … |" 'find . -name x | xargs sed -i s/a/b/'
block inplace "Position: nach ;" 'ls; sed -i s/a/b/ f'
block inplace "Position: nach &&, Trenner davor in Anführungszeichen" 'echo "a;b" && sed -i s/a/b/ f'
block inplace "Position: nach ||" 'false || sed -i s/a/b/ f'
block inplace "Position: nach einer einzelnen &" 'sleep 1 & sed -i s/a/b/ f'
block inplace "Position: nach Zeilenumbruch" $'ls\nsed -i s/a/b/ f'
block inplace "Position: nach Pipe" 'cat f | sed -i s/a/b/ f'
block inplace "Kopf in Anführungszeichen" '"sed" -i s/a/b/ f'
block inplace "Kopf mit Backslash" '\sed -i s/a/b/ f'
block inplace "Kopf: busybox" 'busybox sed -i s/a/b/ f'
block inplace "Kopf: gsed" 'gsed -i s/a/b/ f'
block inplace "Wrapper mit absolutem Pfad: /usr/bin/env" '/usr/bin/env sed -i s/a/b/ f'
block inplace "Wrapper mit absolutem Pfad: /usr/bin/sudo" '/usr/bin/sudo sed -i s/a/b/ f'
block inplace "Flag mit leeren Anführungszeichen: -i''" "sed -i'' s/a/b/ f"
block inplace "Flag in Anführungszeichen: '-i'" "sed '-i' s/a/b/ f"
block inplace "perl -i'' -pe" "perl -i'' -pe s/a/b/ f"
block inplace "Abkürzung --in-p" 'sed --in-p s/a/b/ f'
block inplace "Abkürzung --i" 'sed --i s/a/b/ f'
block inplace "Bündel hinter einem Dateinamen: -input.txt liest sed als -i mit Suffix nput.txt" 'sed -n 1p -input.txt'
block inplace "Schlüsselwort: for … do" 'for f in a b; do sed -i s/a/b/ "$f"; done'
block inplace "Schlüsselwort: while … do" 'while read f; do sed -i s/a/b/ "$f"; done < /tmp/scratch/liste'
block inplace "Schlüsselwort: if … then" 'if true; then sed -i s/a/b/ f; fi'
block inplace "Schlüsselwort: else" 'if false; then ls; else sed -i s/a/b/ f; fi'
block inplace "Schlüsselwort: elif" 'if false; then ls; elif sed -i s/a/b/ f; then ls; fi'
block inplace "Schlüsselwort: if, Bedingung" 'if sed -i s/a/b/ f; then ls; fi'
block inplace "Schlüsselwort: while, Bedingung" 'while sed -i s/a/b/ f; do ls; done'
pass "Kommando-Ersetzung, danach Text: kein Case-Label" 'echo $(date) sed -i x'
pass "case ist beendet: Text hinter einer Kommando-Ersetzung" 'case x in x) ls;; esac; echo $(date) sed -i x'
block inplace "Schlüsselwort: until" 'until sed -i s/a/b/ f; do ls; done'
block inplace "Schlüsselwort: !" '! sed -i s/a/b/ f'
block inplace "Schlüsselwort: case, erstes Label" 'case x in x) sed -i s/a/b/ f;; esac'
block inplace "Schlüsselwort: case, zweites Label" 'case x in y) ls;; x) sed -i s/a/b/ f;; esac'
block inplace "Schlüsselwort: case, Label mit Alternative" 'case x in y|x) sed -i s/a/b/ f;; esac'
block inplace "Funktionsdefinition" 'f() { sed -i s/a/b/ f; }'
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
pass "sed -n mit Bündel-Text im Skript" "sed -n '/-i/p' f"
pass "sed: -i nur im Anführungszeichen-Argument mit Leerraum" "sed 'a -i b' f"
pass "find -exec sh -c ohne in-place" "find . -exec sh -c 'sed -n 1p {}' \\;"
pass "Wrapper mit absolutem Pfad, ohne in-place" '/usr/bin/env sed -n 1p f'
pass "find -exec sed -n" 'find . -name x -exec sed -n 1p {} +'
pass "find -exec grep -i" 'find . -exec grep -i x {} +'
pass "find: -iname hinter dem -exec-Ende" 'find . -exec sed -n 1p {} + -iname x'
pass "bash -c mit sed -n" 'bash -c "sed -n 1p f"'
pass "xargs grep mit Option" 'ls | xargs -r grep -nE x'
pass "xargs grep mit Optionswert" 'ls | xargs -n 1 grep -i x'
pass "xargs sed -n mit Optionswert" 'ls | xargs -n 1 sed -n 1p'

# --- Anführungszeichen: ein Trenner in einem Argument startet kein Segment --
pass "Muster mit | (Suchlauf-Form, drei Glieder)" "git grep -E 'sed -i|perl -pi|awk -i' -- x"
pass "Muster mit | ohne Pfadtrenner, mittleres Glied perl -pi" "grep -E 'a|perl -pi|b' f"
pass "Muster mit | in doppelten Anführungszeichen" 'grep -E "sed -i|perl -pi" f'
pass "Muster mit | und awk -i inplace als Glied" "git grep -E 'sed -i|awk -i inplace|x'"
pass "Muster mit | und Leerzeichen vor dem Anführungszeichen" "grep -E 'a|sed -i ' f"
pass "Muster mit | und python-Glied auf einem Repo-Pfad" "grep -E 'sed -i|python' docs/x"
pass "Muster mit | am Ende (perl -pi als letztes Glied)" "grep -E 'a|perl -pi' f"
pass "echo mit | in einfachen Anführungszeichen" "echo 'a|sed -i x'"
pass "echo mit | und Apostroph in doppelten Anführungszeichen" 'echo "it'"'"'s|sed -i x"'
pass "echo mit doppelten Anführungszeichen in einfachen" "echo 'say \"a|sed -i x\"'"
pass "escaptes Anführungszeichen in doppelten Anführungszeichen" 'echo "a\"|sed -i x"'
pass "Trenner ; in Anführungszeichen" "echo 'a; sed -i x'"
pass "Trenner & in Anführungszeichen (Commit-Message)" 'git commit -m "a & sed -i b"'
pass "Trenner ( in Anführungszeichen" "echo 'a (sed -i x'"
pass "Zeilenumbruch in Anführungszeichen" $'echo "a\nsed -i x"'
pass "Backslash vor |" 'grep a\|sed\ -i\|b f'
pass "Backslash vor Leerzeichen" 'cp a\ b\ sed\ -i c'
pass "git log --grep mit Alternation" "git log --grep='sed -i|perl -pi'"
pass "Tab in Anführungszeichen" $'sed \'a\t-i\tb\' f'
pass "Trenner | ohne Leerraum in doppelten Anführungszeichen" 'echo "x|sed" -i f'
pass "Trenner & ohne Leerraum in doppelten Anführungszeichen" 'echo "x&sed" -i f'
pass "Trenner ; ohne Leerraum in doppelten Anführungszeichen" 'echo "x;sed" -i f'
pass "Trenner ( ohne Leerraum in doppelten Anführungszeichen" 'echo "x(sed" -i f'
pass "Trenner | ohne Leerraum in einfachen Anführungszeichen" "echo 'x|sed' -i f"
pass "Backtick ohne Leerraum in einfachen Anführungszeichen" "echo 'x\`sed' -i f"
pass "Kommando-Ersetzung endet in doppelten Anführungszeichen" 'echo "$(date)|sed" -i f'
pass "Backtick-Ersetzung endet in doppelten Anführungszeichen" 'echo "`date`|sed" -i f'
pass "Heredoc-Text mit Apostroph (unbalanciert, ohne Kommando)" $'cat <<EOF\nit\'s fine\nEOF'
pass "make mit Variable in Anführungszeichen" "make suchlauf-nachmessen PLAN='a|sed -i b'"
pass "Muster mit | zweigliedrig" "git grep -E 'sed -i|perl -pi'"
pass "git log --grep" "git log --grep='sed -i'"

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
block interp "Wrapper mit absolutem Pfad: /usr/bin/env python3" '/usr/bin/env python3 tools/x.py'
block interp "python3 hinter env mit Option" 'env -i python3 tools/x.py'
block interp "Kopf mit Backslash" '\python3 tools/x.py'
block interp "python3 in find -exec" 'find . -exec python3 tools/x.py {} \;'
block interp "python3 in eval" 'eval "python3 tools/x.py"'
pass "python3 in einem Muster mit | auf einem Scratchpad-Pfad" "grep -E 'a|python3 tools/x' /tmp/scratch/f"

# Pfadzeichen-Klasse vor dem Namen: [A-Za-z0-9_./~-] ist keine Repo-Ebene
pass "Pfadzeichen davor: Kleinbuchstabe" 'python3 -c x/adocs/y'
pass "Pfadzeichen davor: Großbuchstabe" 'python3 -c xAdocs/y'
pass "Pfadzeichen davor: Ziffer" 'python3 -c x9docs/y'
pass "Pfadzeichen davor: Unterstrich" 'python3 -c x_docs/y'
pass "Pfadzeichen davor: Bindestrich" 'python3 -c x-docs/y'
pass "Pfadzeichen davor: Tilde" 'python3 -c x~docs/y'
pass "Pfadzeichen davor: Punkt" 'python3 -c x.docs/y'
pass "Pfadzeichen davor: Schrägstrich" 'python3 -c x/docs/y'
block interp "kein Pfadzeichen davor: Gleichheitszeichen" 'python3 -c x=docs/y'
block interp "kein Pfadzeichen davor: Doppelpunkt" 'python3 -c x:docs/y'
block interp "kein Pfadzeichen davor: Komma" 'python3 -c x,docs/y'
block interp "kein Pfadzeichen davor: Anführungszeichen" "python3 -c \"open(\\\"docs/y\\\")\""
pass "Pfadzeichen dahinter: Kleinbuchstabe" 'python3 -c Makefilex'
pass "Pfadzeichen dahinter: Großbuchstabe" 'python3 -c MakefileX'
pass "Pfadzeichen dahinter: Ziffer" 'python3 -c Makefile9'
pass "Pfadzeichen dahinter: Unterstrich" 'python3 -c Makefile_x'
pass "Pfadzeichen dahinter: Bindestrich" 'python3 -c Makefile-x'
pass "Pfadzeichen dahinter: Tilde" 'python3 -c Makefile~'
pass "Pfadzeichen dahinter: Punkt" 'python3 -c Makefile.x'
pass "Pfadzeichen dahinter: Schrägstrich" 'python3 -c Makefile/x'
block interp "kein Pfadzeichen dahinter: Komma" 'python3 -c Makefile,x'
block interp "kein Pfadzeichen dahinter: Klammer" 'python3 -c Makefile)'
block interp "kein Pfadzeichen dahinter: Gleichheitszeichen" 'python3 -c Makefile=x'
block interp "Datei-Name am Ende des Befehls" 'python3 -c Makefile'

# Die Meldung der Klasse interp nennt die Wege aus AGENTS.md 3.1
n=$((n + 1))
run_guard 'python3 tools/x.py'
judge_block interp "Meldung interp: Edit/Write"
if ! printf '%s' "$out" | grep -qF 'Edit/Write' || ! printf '%s' "$out" | grep -qF 'sed s/a/b/ file > ' \
  || printf '%s' "$out" | grep -qF 'mutate'; then
  report "Meldung interp: Weg nach AGENTS.md 3.1" "nennt nicht Edit/Write und sed nach stdout, oder empfiehlt ein Host-Skript"
fi

# --- benannte Falsch-Positiv-Raender (der Guard blockt trotz Anführungszeichen-Lesung) --
block inplace "Rand: Heredoc-Zeile mit sed -i am Kopf" $'cat <<EOF\nsed -i s/a/b/ f\nEOF'
block inplace "Rand: sed -i am Zeilenanfang hinter einem Heredoc-Apostroph (unbalanciert)" $'cat <<EOF\nit\'s\nsed -i s/a/b/ f\nEOF'
block inplace "Rand: unbalanciertes Anführungszeichen, Trenner im Muster" "grep -E 'a|sed -i f"
block inplace "Rand: sed -i nach & in Text mit unbalanciertem Anführungszeichen" 'echo "a & sed -i b'
block inplace "Rand: sed -i hinter \\; (find-Ende, kein Trenner der Shell)" 'echo a\;sed -i s/a/b/ f'
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
pass "Grenze: Optionswert hinter einem Wrapper (sudo)" 'sudo -u x sed -i s/a/b/ f'
pass "Grenze: Optionswert hinter einem Wrapper (env)" 'env -u FOO sed -i s/a/b/ f'
pass "Grenze: Optionswert hinter einem Wrapper (nice)" 'nice -n 10 sed -i s/a/b/ f'
pass "Grenze: Werkzeug aus einer Variablen" 'x=sed; $x -i s/a/b/ f'
pass "Grenze: eval mit Variable" 'eval "$CMD"'
pass "Grenze: Abkürzung der awk-Langoption" "gawk --inc=inplace '{print}' f"
pass "Grenze: perl-Bündel mit Buchstaben außerhalb der Klasse" 'perl -Wpi -e s/a/b/ f'

if [ "$fail" -eq 0 ]; then
  echo "run-command-guard-tests: alle $n Fälle bestanden"
else
  echo "run-command-guard-tests: mindestens ein Fall von $n gescheitert" >&2
  exit 1
fi
