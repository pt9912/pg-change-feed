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
# Geprueft wird die Befehlsposition jedes Kommando-Segments (Trennung an
# ; & && || | $( ` ( und Zeilenenden) — `git commit -m "... pip ..."` bleibt
# erlaubt, `/usr/bin/pip` und `sudo pip` werden erkannt. Zuweisungs- und
# Wrapper-Praefixe (VAR=…, sudo/env/command/…, danach deren Optionen wie in
# `xargs -r sed`) sowie fuehrende Brace-Group-Delimiter ({ … }) werden
# uebersprungen.
# In-place: ein Segment mit Kopf sed/perl/awk/gawk blockt, wenn eines seiner
# Flag-Tokens die in-place Form traegt (sed: --in-place, -i, -i.bak, -ni; perl:
# -i, -pi, -0777pi; awk/gawk: -i inplace, -iinplace, --include=inplace) — ohne
# Ruecksicht auf das Ziel, auch auf einer Scratchpad-Kopie und hinter
# -exec/-execdir/-ok/-okdir. Die Flag-Tokens werden ROH gelesen (ohne
# Anfuehrungszeichen-Bereinigung; ein Token mit Anfuehrungszeichen ist kein
# Flag): `grep 'sed -i|perl -pi'` als Text blockt nicht.
# Host-Interpreter: ein Segment mit Kopf python/python3/python3.N/perl blockt,
# wenn der GANZE Befehlsstring ein Repo-Pfad-Muster traegt (Absolutpfad der
# Repo-Wurzel oder ein Name der obersten Repo-Ebene ohne davorstehendes
# Pfadzeichen, `./` erlaubt).
# ACHTUNG quote-BLIND: ein Trenner (; & | …) IN einem Argument (z. B. einer
# Commit-Message `… & gofmt …`) startet ein neues Segment — steht dort ein
# blockiertes Wort am Kopf (oder `sed -i` bzw. ein python-Aufruf mit Repo-Name),
# blockt der Guard (False-Positive). Abhilfe: Commit-Messages via
# `git commit -F <datei>`.
# Sub-Shell-Strings (`bash -c "…"`, auch in Flag-Buendeln wie -lc/-ec/-cx)
# werden rekursiv geprueft (Tiefe <= 3, darueber fail-closed).
# Grenze: der Guard ist ein Stolperdraht, KEINE Sandbox; Vollstaendigkeit ist
# nicht das Ziel. Nicht gelesen werden Umleitungen und flaglose Schreibwege
# (> datei, tee, dd of=, sed … > tmp && mv), ein cd im selben Kommando, Variablen,
# Globs und ~ als Pfad, ein Skript, das ein Interpreter liest (`python3 x.py`
# mit Text ohne Repo-Namen), `bash skript.sh`, jedes andere in-place-faehige
# Werkzeug (ed, patch, ruby -i, …), Optionen mit Wert hinter einem Wrapper
# (`sudo -u x`), die Rezepte hinter make und die Docker-Bauten. Grenz-Zeile und
# Tabellentest: harness/conventions/MR-003-guard-inplace-textwerkzeug.md,
# `make test-command-guard`.
#
# Im Pass-Fall: KEINE Ausgabe — "approve" wuerde das Permission-System
# ueberspringen; ohne Ausgabe laeuft die normale Permission-Entscheidung.
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
extractor="$here/../../tools/harness/extract-command.awk"

# Die Begruendungstexte tragen weder `"` noch `\`, damit die Ausgabe gueltiges
# JSON bleibt.
REASON_PKG="This repository is make/Docker-only (AGENTS.md Hard Rule 3.1). Use make targets; do not install or run host package managers or host toolchains (apt/brew/pip/npm/cargo/go/...). On parse doubt the guard fails closed."
REASON_INPLACE="In-place text tools (sed -i, perl -i, awk -i inplace) rewrite repo files without a trace and are blocked, also on a scratch copy (AGENTS.md Hard Rule 3.1). Change a file with the Edit/Write tools; to try a change on a copy, write to stdout: sed s/a/b/ file > /path/to/scratch-copy."
REASON_INTERP="Host python/perl on repo paths is blocked (AGENTS.md Hard Rule 3.1). Change a repo file with the Edit/Write tools; run a mutation script against an absolute scratch copy: python3 /path/to/scratch/mutate.py /path/to/scratch/copy."
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
PREFIXES="sudo env command exec nice time xargs eval"
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

# inplace_form <werkzeug> <index des Kopf-Tokens>; return 0 = in-place Form.
# Liest die ROHEN Tokens (`toks`) des aufrufenden scan ab dem Index+1 bis zum
# Ende des Segments oder bis zum find-Ende (`+`, `\`). Ein Werkzeug ausserhalb
# von sed/perl/awk/gawk ist keine in-place Form.
inplace_form() {
  local tool=$1 k t nxt
  case $tool in sed|perl|awk|gawk) ;; *) return 1;; esac
  for ((k = $2 + 1; k < ${#toks[@]}; k++)); do
    t=${toks[$k]}
    case $t in '+'|'\'|';') break;; esac
    case $t in *\"*|*\'*) continue;; esac   # ein Token mit Anfuehrungszeichen ist kein Flag
    nxt=${toks[$((k + 1))]:-}
    case $tool in
      sed)
        case $t in --in-place|--in-place=*) return 0;; esac
        [[ $t =~ ^-[nEsrzu]*i ]] && return 0
        ;;
      perl)
        [[ $t =~ ^-[0-9lanpsw]*i ]] && return 0
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
        [A-Za-z0-9_./~-]) ;;
        *)
          if [ "${needle: -1}" = "/" ]; then return 0; fi
          case ${post:0:1} in
            [A-Za-z0-9_./~-]) ;;
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
  local s=$cmd
  s=${s//'&&'/$'\n'}; s=${s//'&'/$'\n'}; s=${s//'||'/$'\n'}; s=${s//'|'/$'\n'}
  s=${s//';'/$'\n'};  s=${s//\$\(/$'\n'};  s=${s//'`'/$'\n'}
  s=${s//'('/$'\n'};  s=${s//$'\r'/$'\n'}
  local seg head i j k rest x prefixed tool
  local -a toks stoks
  while IFS= read -r seg; do
    read -ra toks <<< "$seg"
    [ "${#toks[@]}" -eq 0 ] && continue
    stoks=()
    for x in "${toks[@]}"; do strip_quotes "$x"; stoks+=("$STRIPPED"); done
    i=0
    prefixed=0
    while [ "$i" -lt "${#stoks[@]}" ]; do
      if [[ "${stoks[$i]}" =~ ^[A-Za-z_][A-Za-z0-9_]*= ]]; then i=$((i+1)); continue; fi
      if in_set "$PREFIXES" "${stoks[$i]}"; then prefixed=1; i=$((i+1)); continue; fi
      # Optionen hinter einem Wrapper-Praefix (`xargs -r sed`, `xargs -I{} sed`)
      # gehoeren nicht zum Kopf; ein Optionswert (`sudo -u x`) wird nicht erkannt.
      if [ "$prefixed" -eq 1 ]; then
        case "${stoks[$i]}" in -*|"{}") i=$((i+1)); continue;; esac
      fi
      # fuehrende Brace-Group-Delimiter ueberspringen: `{ go build; }` ->
      # Kopf waere sonst `{` und das Tool an Position 2 entkaeme der Pruefung.
      case "${stoks[$i]}" in "{"|"}") i=$((i+1)); continue;; esac
      break
    done
    [ "$i" -ge "${#stoks[@]}" ] && continue
    head=${stoks[$i]}; head=${head##*/}    # /usr/bin/pip -> pip
    in_set "$BLOCKED" "$head" && return 0
    if inplace_form "$head" "$i"; then REASON=$REASON_INPLACE; return 0; fi
    for ((k = i + 1; k < ${#stoks[@]} - 1; k++)); do
      case "${stoks[$k]}" in
        -exec|-execdir|-ok|-okdir)
          tool=${stoks[$((k+1))]}; tool=${tool##*/}
          if inplace_form "$tool" "$((k+1))"; then REASON=$REASON_INPLACE; return 0; fi
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
          rest="${stoks[*]:$((j+1))}"
          scan "$rest" "$((depth+1))" && return 0
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
