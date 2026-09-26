#!/usr/bin/env bash
# pretooluse-command-guard — blockt Host-Paketmanager (apt/pip/npm/cargo/...),
# in-place Textwerkzeuge (sed -i, perl -i, awk -i inplace) und Host-python/-perl
# auf Repo-Pfaden; dieses Repo baut make/Docker-only (AGENTS.md Hard Rule 3.1,
# Absatz Durchsetzung). Reines bash + awk, KEIN node/jq/OCI.
#
# Der awk-Extraktor (tools/harness/extract-command.awk) zieht nur das eine Feld
# tool_input.command aus der Hook-stdin-JSON; bei Parse-Zweifel (malformed,
# abgeschnitten, \u-Escape im Befehl) -> fail-closed (block).
#
# Segmentierung: tools/harness/mask-quotes.awk maskiert Trenner, Leerraum und
# Zeilenumbrueche in Anfuehrungszeichen und hinter einem Backslash; getrennt wird
# an ; & && || | $( ` ( und Zeilenenden AUSSERHALB von Anfuehrungszeichen, ein
# Anfuehrungszeichen-Argument ist ein Token (`grep -E 'sed -i|perl -pi'` und
# `git commit -m "a & sed -i b"` blocken nicht). `$(` und Backtick in doppelten
# Anfuehrungszeichen fuehren aus und trennen weiter. Ein unbalanciertes
# Anfuehrungszeichen (auch ein Apostroph im Text eines Heredocs) segmentiert
# ohne Anfuehrungszeichen-Kenntnis: mehr Segmente, nie weniger.
# Geprueft wird die Befehlsposition jedes Segments — `git commit -m "... pip ..."`
# bleibt erlaubt, `/usr/bin/pip`, `\sed` und `sudo pip` werden erkannt.
# Zuweisungs- und Wrapper-Praefixe (VAR=…, sudo/env/command/busybox/…, danach
# deren Optionen wie in `xargs -r sed` und der Wert einer xargs-Option wie in
# `xargs -n 1 sed`; `command -v X` zeigt an und fuehrt nichts aus),
# Shell-Schluesselwoerter (do then else elif if while until !),
# fuehrende Delimiter ({ )) und Case-Labels (`x) cmd` nach `case`) werden
# uebersprungen.
# In-place: ein Segment mit Kopf sed/gsed/perl/awk/gawk blockt, wenn eines seiner
# Flag-Tokens die in-place Form traegt (sed: --in-place und jede eindeutige
# Abkuerzung ab --i, -i, -i.bak, -ni, -bi; perl: -i, -pi, -0777pi bis zum ersten
# Nicht-Options-Token, das -e/-I-Argument uebersprungen; awk/gawk: -i inplace,
# -iinplace, --include=inplace) — ohne Ruecksicht auf das Ziel, auch auf einer
# Scratchpad-Kopie. Anfuehrungszeichen gehoeren nicht zum Flag (`sed -i''`,
# `sed '-i'` blocken; ein Anfuehrungszeichen-Argument mit Leerraum ist ein
# Token und kein Flag). Ein Kopf hinter -exec/-execdir/-ok/-okdir wird als
# eigenes Kommando gelesen (find … -exec sed -i, -exec sh -c '…').
# Host-Interpreter: ein Segment mit Kopf python/python3/python3.N/perl blockt,
# wenn der GANZE Befehlsstring ein Repo-Pfad-Muster traegt (Absolutpfad der
# Repo-Wurzel oder ein Name der obersten Repo-Ebene ohne davorstehendes
# Pfadzeichen, `./` erlaubt).
# Sub-Shell-Strings (`bash -c "…"`, auch in Flag-Buendeln wie -lc/-ec/-cx) und
# der Rest hinter `eval` werden rekursiv geprueft (Tiefe <= 3, darueber
# fail-closed).
# Grenze: der Guard ist ein Stolperdraht, KEINE Sandbox; Vollstaendigkeit ist
# nicht das Ziel. Nicht gelesen werden Umleitungen und flaglose Schreibwege
# (> datei, tee, dd of=, sed … > tmp && mv), ein cd im selben Kommando, Variablen
# (auch `$x -i` und `eval "$cmd"`), Globs, Aliase und ~ als Pfad, ein Skript, das
# ein Interpreter liest (`python3 x.py` mit Text ohne Repo-Namen), `bash skript.sh`,
# jedes andere in-place-faehige Werkzeug (ed, patch, ruby -i, …) und jeder andere
# Interpreter (node, ruby), ein perl-Buendel mit Buchstaben ausserhalb der Klasse
# (`-Wpi`), Optionen mit Wert hinter einem anderen Wrapper als xargs (`sudo -u x`),
# die Rezepte hinter make und die Docker-Baeuten. Die Zeilen eines Heredocs gelten
# als Kommando-Zeilen (`cat <<EOF` mit `sed -i` am Zeilenanfang blockt).
# Grenz-Zeile und Tabellentest:
# harness/conventions/MR-003-guard-inplace-textwerkzeug.md,
# `make test-command-guard`.
#
# Im Pass-Fall: KEINE Ausgabe — "approve" ueberspringt das Permission-System;
# ohne Ausgabe laeuft die normale Permission-Entscheidung.
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
extractor="$here/../../tools/harness/extract-command.awk"
masker="$here/../../tools/harness/mask-quotes.awk"

# Die Begruendungstexte tragen weder `"` noch `\`, damit die Ausgabe gueltiges
# JSON bleibt.
REASON_PKG="This repository is make/Docker-only (AGENTS.md Hard Rule 3.1). Use make targets; do not install or run host package managers or host toolchains (apt/brew/pip/npm/cargo/go/...). On parse doubt the guard fails closed."
REASON_INPLACE="In-place text tools (sed -i, perl -i, awk -i inplace) rewrite repo files without a trace and are blocked, also on a scratch copy (AGENTS.md Hard Rule 3.1). Change a file with the Edit/Write tools; to try a change on a copy, write to stdout: sed s/a/b/ file > /path/to/scratch-copy."
REASON_INTERP="Host python/perl on repo paths is blocked (AGENTS.md Hard Rule 3.1). Change a repo file with the Edit/Write tools or a repo tool behind make; to try a change on a copy, edit the copy with Edit/Write or write to stdout: sed s/a/b/ file > /path/to/scratch-copy."
REASON=""

emit_block() {
  printf '{\n  "decision": "block",\n  "reason": "%s"\n}\n' "${REASON:-$REASON_PKG}"
}

# BLOCKED = gebackener universeller Boden (fail-safe, NIE fail-open) + optionale
# Sprach-Fragmente aus tools/harness/blocked/* (add-lang droppt blocked/<sprache>). So
# blockt der Guard sprachlos schon apt/pip/npm/cargo; ein geleertes/fehlendes blocked/
# laesst den Boden UNBERUEHRT (der Guard darf nie fail-open sein). Reine bash+cat-Union
# (kein node/jq): die blocked/*-Dateien sind Wortlisten (whitespace-getrennt).
BLOCKED="apt apt-get brew pip pip3 pipx npm pnpm yarn npx corepack cargo rustup gem conda"
blocked_dir="$here/../../tools/harness/blocked"
if [ -d "$blocked_dir" ]; then
  for bf in "$blocked_dir"/*; do
    [ -f "$bf" ] && BLOCKED="$BLOCKED $(cat "$bf")"
  done
fi
PREFIXES="sudo env command exec nice time xargs eval busybox"
SHELLS="bash sh zsh dash ksh"

in_set() {  # in_set <space-getrennte-menge> <wort>
  local w
  for w in $1; do [ "$w" = "$2" ] && return 0; done
  return 1
}

# Ergebnis in der globalen STRIPPED (kein Subshell-Fork je Token; der Guard
# laeuft vor JEDEM Bash-Call, Latenz zaehlt).
strip_quotes() {  # fuehrende/abschliessende " und ' entfernen (wie ^["']+|["']+$)
  local s=$1
  while [ -n "$s" ]; do case $s in \"*|\'*) s=${s#?};; *) break;; esac; done
  while [ -n "$s" ]; do case $s in *\"|*\') s=${s%?};; *) break;; esac; done
  STRIPPED=$s
}

# unmask <text>: die Steuerzeichen des Maskierers (mask-quotes.awk) zurueck in
# Trenner, Leerraum und Zeilenumbruch; Ergebnis in der globalen UNMASKED.
unmask() {
  local s=$1
  s=${s//$'\001'/'|'}; s=${s//$'\002'/'&'}; s=${s//$'\003'/';'}; s=${s//$'\004'/'('}
  s=${s//$'\005'/'`'}; s=${s//$'\006'/' '}; UNMASKED=${s//$'\007'/$'\n'}
}

# inplace_form <werkzeug> <index des Kopf-Tokens>; return 0 = in-place Form.
# Liest die ROHEN Tokens (`toks`) des aufrufenden scan ab dem Index+1 bis zum
# Ende des Segments oder bis zum find-Ende (`+`, `\`). Ein Werkzeug ausserhalb
# von sed/gsed/perl/awk/gawk ist keine in-place Form.
inplace_form() {
  local tool=$1 k t b nxt
  case $tool in sed|gsed|perl|awk|gawk) ;; *) return 1;; esac
  for ((k = $2 + 1; k < ${#toks[@]}; k++)); do
    t=${toks[$k]}
    case $t in '+'|'\'|';') break;; esac
    t=${t//[\"\']/}                       # Anfuehrungszeichen gehoeren nicht zum Flag
    nxt=${toks[$((k + 1))]:-}; nxt=${nxt//[\"\']/}
    case $tool in
      sed|gsed)
        b=${t%%=*}
        [[ ${#b} -ge 3 && --in-place == "$b"* ]] && return 0
        [[ $t =~ ^-[nEsrzub]*i ]] && return 0
        ;;
      perl)
        [[ $t =~ ^-[0-7lanpsw]*i ]] && return 0
        if [[ $t =~ ^-[0-7lanpsw]*[eE]$ || $t == -I ]]; then k=$((k + 1)); continue; fi
        case $t in -*) ;; *) break;; esac  # ab dem Skriptnamen keine perl-Optionen mehr
        ;;
      awk|gawk)
        case $t in
          -iinplace|--include=inplace) return 0;;
          -i|--include) [ "$nxt" = inplace ] && return 0;;
        esac
        ;;
    esac
  done
  return 1
}

# repo_pfad_muster <befehl>; return 0 = der Befehlsstring nennt einen Repo-Pfad:
# den Absolutpfad der Repo-Wurzel (zwei Ebenen ueber .claude/hooks/) oder den
# Namen eines Eintrags der obersten Ebene (Verzeichnis mit folgendem `/`, Datei
# als ganzes Wort), vor dem kein Pfadzeichen [A-Za-z0-9_./~-] steht; ein
# vorangestelltes `./` ist erlaubt.
repo_pfad_muster() {
  local cmd=$1 root e needle rest pre post c
  local pz='[A-Za-z0-9_./~-]'
  root="$(cd "$here/../.." && pwd)"
  case $cmd in *"$root/"*) return 0;; esac
  for e in "$root"/* "$root"/.[!.]*; do
    [ -e "$e" ] || continue
    e=${e##*/}
    [ "$e" = .git ] && continue
    if [ -d "$root/$e" ]; then needle="$e/"; else needle=$e; fi
    rest=" $cmd"
    while [[ $rest == *"$needle"* ]]; do
      pre=${rest%%"$needle"*}
      post=${rest#*"$needle"}
      c=${pre: -1}
      if [ "$c" = "/" ] && [ "${pre: -2:1}" = "." ]; then c=${pre: -3:1}; fi
      case $c in
        $pz) ;;
        *)
          if [ "${needle: -1}" = "/" ]; then return 0; fi
          case ${post:0:1} in
            $pz) ;;
            *) return 0;;
          esac
          ;;
      esac
      rest="/$post"
    done
  done
  return 1
}

scan() {  # scan <cmd> <tiefe>; return 0 = BLOCK, 1 = ok
  local cmd=$1 depth=$2
  [ "$depth" -gt 3 ] && return 0          # zu tief verschachtelt -> fail-closed
  local s
  s=$(CMD=$cmd awk -f "$masker") || return 0   # ohne Maskierer keine Pruefung -> fail-closed
  s=${s//'&&'/$'\n'}; s=${s//'&'/$'\n'}; s=${s//'||'/$'\n'}; s=${s//'|'/$'\n'}
  s=${s//';'/$'\n'};  s=${s//\$\(/$'\n'};  s=${s//'`'/$'\n'}
  s=${s//'('/$'\n'};  s=${s//$'\r'/$'\n'}
  local seg head w i j k x prefixed wrapper= incase=0
  local -a toks stoks
  while IFS= read -r seg; do
    read -ra toks <<< "$seg"
    [ "${#toks[@]}" -eq 0 ] && continue
    stoks=()
    for x in "${toks[@]}"; do strip_quotes "$x"; stoks+=("$STRIPPED"); done
    i=0
    prefixed=0
    while [ "$i" -lt "${#stoks[@]}" ]; do
      w=${stoks[$i]}; w=${w#\\}            # \sed -> sed
      if [[ $w =~ ^[A-Za-z_][A-Za-z0-9_]*= ]]; then i=$((i+1)); continue; fi
      if in_set "$PREFIXES" "${w##*/}"; then
        # `command -v X` / `command -V X` zeigen an und fuehren nichts aus
        if [ "$w" = command ] && [[ ${stoks[$((i+1))]:-} =~ ^-[a-zA-Z]*[vV] ]]; then continue 2; fi
        if [ "$w" = eval ]; then   # der Rest ist ein Kommando-String
          unmask "${stoks[*]:$((i+1))}"
          scan "$UNMASKED" "$((depth+1))" && return 0
          continue 2
        fi
        prefixed=1; wrapper=${w##*/}; i=$((i+1)); continue
      fi
      # Optionen hinter einem Wrapper-Praefix (`xargs -r sed`, `xargs -I{} sed`)
      # gehoeren nicht zum Kopf; der Wert einer xargs-Option (`xargs -n 1 sed`)
      # ebenso, der Wert einer Option anderer Wrapper (`sudo -u x`) nicht.
      if [ "$prefixed" -eq 1 ]; then
        case $w in -n|-P|-L|-I|-s|-d|-a|-E) [ "$wrapper" = xargs ] && i=$((i+1));; esac
        case $w in -*|"{}") i=$((i+1)); continue;; esac
      fi
      # Shell-Schluesselwoerter und fuehrende Delimiter ueberspringen: `{ go build; }`,
      # `do sed -i …` und `) { sed -i …` tragen das Kommando an Position 2.
      [ "$w" = esac ] && incase=0
      case $w in
        "{"|")"|do|then|else|elif|if|while|until|"!") i=$((i+1)); continue;;
        case)
          incase=1; j=$i
          while [ "$j" -lt "${#stoks[@]}" ] && [[ ${stoks[$j]} != *')' ]]; do j=$((j+1)); done
          i=$j; continue;;
        *')') if [ "$incase" -eq 1 ]; then i=$((i+1)); continue; fi;;   # Case-Label
      esac
      break
    done
    [ "$i" -ge "${#stoks[@]}" ] && continue
    head=${w##*/}                           # /usr/bin/pip -> pip
    in_set "$BLOCKED" "$head" && return 0
    if inplace_form "$head" "$i"; then REASON=$REASON_INPLACE; return 0; fi
    for ((k = i + 1; k < ${#stoks[@]} - 1; k++)); do
      case "${stoks[$k]}" in
        -exec|-execdir|-ok|-okdir)
          unmask "${stoks[*]:$((k+1))}"
          scan "$UNMASKED" "$((depth+1))" && return 0
          ;;
      esac
    done
    case $head in
      python|python[0-9]*|perl)
        if repo_pfad_muster "$cmd"; then REASON=$REASON_INTERP; return 0; fi
        ;;
    esac
    if in_set "$SHELLS" "$head"; then
      # -c auch in Flag-Buendeln (-lc, -ec, -cx, …): bei sh/bash ist c das
      # einzige Single-Letter-Flag mit Kommando-String-Semantik.
      j=$((i+1))
      while [ "$j" -lt "${#stoks[@]}" ]; do
        if [[ "${stoks[$j]}" =~ ^-[a-z]*c[a-z]*$ ]]; then
          unmask "${stoks[*]:$((j+1))}"
          scan "$UNMASKED" "$((depth+1))" && return 0
          break
        fi
        j=$((j+1))
      done
    fi
  done <<< "$s"
  return 1
}

input="$(cat)"

# Ohne awk keine Pruefung -> fail-closed. (awk ist POSIX-Basis.)
command -v awk >/dev/null 2>&1 || { emit_block; exit 0; }

set +e
cmd="$(printf '%s' "$input" | awk -f "$extractor")"
rc=$?
set -e
[ "$rc" -ne 0 ] && { emit_block; exit 0; }   # Parse-Zweifel -> fail-closed

scan "$cmd" 0 && emit_block
# Pass-Fall: keine Ausgabe — normale Permission-Pruefung uebernimmt.
exit 0
