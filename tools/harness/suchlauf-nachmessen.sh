#!/usr/bin/env bash
# suchlauf-nachmessen — misst die `suchlauf`-Blöcke eines Slice-Plans nach:
# je Zeile `<Stand> <Soll> <Argumente von git grep>` führt es
# `git grep -n <Argumente>` am Stand aus (Commit-Kennung, oder `diff` für den
# Arbeitsbaum), zählt die Trefferzeilen und vergleicht sie mit dem Soll.
# Die Plan-Datei ist immer aus dem Suchraum ausgeschlossen, unter ihrem
# Dateinamen in jedem Verzeichnis (der Plan wandert durch die
# Lifecycle-Verzeichnisse). Eine Plan-Zeile führt keinen Shell-Code aus: der
# Zerleger expandiert nichts, Optionen und Pathspec-Magic stehen auf einer
# Allow-List (check_opts, check_paths). Vertrag, Grenze und Exit-Codes:
# harness/sensors/suchlauf-nachmessen.md (AGENTS.md §3.13).
#
# Aufruf: suchlauf-nachmessen.sh <Plan-Datei>
# Exit: 0 alle Zeilen stimmen · 1 mindestens eine Abweichung · 2 Eingabefehler
set -uo pipefail

if [ $# -ne 1 ]; then
  echo "Aufruf: suchlauf-nachmessen.sh <Plan-Datei>" >&2
  exit 2
fi
plan=$1
if [ ! -f "$plan" ]; then
  echo "suchlauf-nachmessen: Plan-Datei '$plan' nicht gefunden" >&2
  exit 2
fi
plan_abs=$(realpath "$plan")
root=$(git rev-parse --show-toplevel) || exit 2
cd "$root" || exit 2

# Ausschluss der Plan-Datei über ihren Dateinamen in jedem Verzeichnis.
excludes=(":(exclude,glob)**/$(basename "$plan_abs")")

# Die Zeilen der `suchlauf`-Blöcke (eingerückte Zäune erlaubt).
lines=()
in_block=0
while IFS= read -r line || [ -n "$line" ]; do
  if [ "$in_block" -eq 0 ]; then
    if [[ $line =~ ^[[:space:]]*\`\`\`suchlauf[[:space:]]*$ ]]; then
      in_block=1
    fi
  elif [[ $line =~ ^[[:space:]]*\`\`\`[[:space:]]*$ ]]; then
    in_block=0
  elif [ -n "${line//[[:space:]]/}" ]; then
    lines+=("$line")
  fi
done <"$plan_abs"
if [ "$in_block" -ne 0 ]; then
  echo "suchlauf-nachmessen: ein suchlauf-Block in '$plan' ist nicht geschlossen" >&2
  exit 2
fi
if [ ${#lines[@]} -eq 0 ]; then
  echo "suchlauf-nachmessen: '$plan' trägt keine Zeile in einem suchlauf-Block (leer ist nicht bestanden)" >&2
  exit 2
fi

# split_words <text>: zerlegt an Leerraum; ' und " gruppieren, ohne
# Expansion und ohne Backslash-Escape. Ergebnis im Array WORDS.
split_words() {
  WORDS=()
  local s=$1 i c cur='' quote='' have=0
  for ((i = 0; i < ${#s}; i++)); do
    c=${s:i:1}
    if [ -n "$quote" ]; then
      if [ "$c" = "$quote" ]; then quote=''; else cur+=$c; fi
      continue
    fi
    case $c in
      "'" | '"') quote=$c; have=1 ;;
      ' ' | $'\t')
        if [ "$have" -eq 1 ]; then WORDS+=("$cur"); cur=''; have=0; fi
        ;;
      *) cur+=$c; have=1 ;;
    esac
  done
  [ -z "$quote" ] || return 1
  if [ "$have" -eq 1 ]; then WORDS+=("$cur"); fi
  return 0
}

# check_opts: prüft die Optionswörter einer Zeile (Array opts) gegen die
# Allow-List; jede andere Option kann Befehle ausführen (-O, --open-files-in-pager),
# Dateien lesen (-f, --no-index) oder ändert den Suchraum. Erlaubt: -e <Muster>,
# Klammern und die Kurzoptionen aus [iwEFGPvIahHocLln] (auch gebündelt), die
# Langformen der Verknüpfung und Muster-Art; höchstens ein Suchwort, und keines
# neben -e. Bei einem Verstoß steht die Meldung in ERR.
check_opts() {
  local i=0 n=${#opts[@]} w positional=0 have_e=0
  ERR=''
  while [ "$i" -lt "$n" ]; do
    w=${opts[$i]}
    case $w in
      -e)
        i=$((i + 1))
        if [ "$i" -ge "$n" ]; then ERR="Option '-e' ohne Muster"; return 1; fi
        have_e=1
        ;;
      --and | --or | --not | --all-match | --ignore-case | --word-regexp | \
        --extended-regexp | --fixed-strings | --basic-regexp | --perl-regexp | \
        --invert-match | --count | --files-with-matches | --files-without-match | \
        --text | '(' | ')') ;;
      --*)
        ERR="Option '$w' ist nicht erlaubt"; return 1
        ;;
      -*)
        if [[ ! $w =~ ^-[iwEFGPvIahHocLln]+$ ]]; then
          ERR="Option '$w' ist nicht erlaubt (erlaubt: -e <Muster>, -i -w -E -F -G -P -v -I -a -h -H -o -c -l -L -n, --and --or --not --all-match)"
          return 1
        fi
        ;;
      *) positional=$((positional + 1)) ;;
    esac
    i=$((i + 1))
  done
  if [ "$have_e" -eq 1 ] && [ "$positional" -gt 0 ]; then
    ERR="Suchwort neben -e (jedes Muster steht hinter einem -e)"; return 1
  fi
  if [ "$positional" -gt 1 ]; then
    ERR="mehr als ein Suchwort vor '--' (weitere Muster stehen hinter -e)"; return 1
  fi
  return 0
}

# check_paths: prüft die Pathspec-Wörter (Array paths); ein Wort mit `:` am
# Anfang ist nur als `:!`, `:^` oder `:(exclude|glob|literal|icase|top,…)` erlaubt.
magic_re='^:\((exclude|glob|literal|icase|top)(,(exclude|glob|literal|icase|top))*\)'
short_re='^:[!^]'
check_paths() {
  local p
  ERR=''
  for p in ${paths[@]+"${paths[@]}"}; do
    if [[ $p == :* ]] && [[ ! $p =~ $short_re ]] && [[ ! $p =~ $magic_re ]]; then
      ERR="Pathspec '$p' ist nicht erlaubt (erlaubt: :! :^ :(exclude|glob|literal|icase|top))"
      return 1
    fi
  done
  return 0
}

# Erster Durchgang: jede Zeile prüfen, bevor ein Befehl läuft.
stands=(); solls=(); optss=(); pathss=()
bad=0
for line in "${lines[@]}"; do
  if ! split_words "$line"; then
    echo "suchlauf-nachmessen: Anführungszeichen nicht geschlossen: $line" >&2
    bad=1; continue
  fi
  empty=0
  for w in "${WORDS[@]}"; do
    if [ -z "$w" ]; then empty=1; fi
  done
  if [ "$empty" -eq 1 ]; then
    echo "suchlauf-nachmessen: leeres Argument in der Zeile: $line" >&2
    bad=1; continue
  fi
  if [ ${#WORDS[@]} -lt 3 ]; then
    echo "suchlauf-nachmessen: Zeile ohne <Stand> <Soll> <Argumente>: $line" >&2
    bad=1; continue
  fi
  stand=${WORDS[0]}
  soll=${WORDS[1]}
  if [[ ! $soll =~ ^[0-9]+$ ]]; then
    echo "suchlauf-nachmessen: Soll '$soll' ist keine Zahl: $line" >&2
    bad=1; continue
  fi
  if [ "$stand" != "diff" ]; then
    if [[ ! $stand =~ ^[0-9a-f]{7,40}$ ]]; then
      echo "suchlauf-nachmessen: Stand '$stand' ist weder eine Commit-Kennung noch 'diff' (HEAD und Namen bewegen sich): $line" >&2
      bad=1; continue
    fi
    if ! git rev-parse --verify --quiet "$stand^{commit}" >/dev/null; then
      echo "suchlauf-nachmessen: Stand '$stand' ist kein Commit dieses Repos: $line" >&2
      bad=1; continue
    fi
  fi
  opts=(); paths=(); seen_sep=0
  for w in "${WORDS[@]:2}"; do
    if [ "$seen_sep" -eq 0 ] && [ "$w" = "--" ]; then
      seen_sep=1
    elif [ "$seen_sep" -eq 0 ]; then
      opts+=("$w")
    else
      paths+=("$w")
    fi
  done
  if [ ${#opts[@]} -eq 0 ]; then
    echo "suchlauf-nachmessen: Zeile ohne Suchmuster: $line" >&2
    bad=1; continue
  fi
  if ! check_opts; then
    echo "suchlauf-nachmessen: $ERR: $line" >&2
    bad=1; continue
  fi
  if ! check_paths; then
    echo "suchlauf-nachmessen: $ERR: $line" >&2
    bad=1; continue
  fi
  stands+=("$stand"); solls+=("$soll")
  # Argumente als eine mit US (0x1f) getrennte Zeichenkette merken.
  optss+=("$(IFS=$'\x1f'; printf '%s' "${opts[*]}")")
  pathss+=("$(IFS=$'\x1f'; printf '%s' "${paths[*]-}")")
done
if [ "$bad" -ne 0 ]; then
  exit 2
fi

# Zweiter Durchgang: messen.
failures=0
errf=$(mktemp "${TMPDIR:-/tmp}/suchlauf-nachmessen.XXXXXX") || exit 2
trap 'rm -f "$errf"' EXIT
for i in "${!lines[@]}"; do
  IFS=$'\x1f' read -r -a opts <<<"${optss[$i]}"
  paths=()
  if [ -n "${pathss[$i]}" ]; then IFS=$'\x1f' read -r -a paths <<<"${pathss[$i]}"; fi
  stand=${stands[$i]}
  # stdout trägt die Trefferzeilen, stderr wird getrennt gehalten und nicht gezählt.
  if [ "$stand" = "diff" ]; then
    out=$(git grep -n "${opts[@]}" -- ${paths[@]+"${paths[@]}"} "${excludes[@]}" 2>"$errf")
  else
    out=$(git grep -n "${opts[@]}" "$stand" -- ${paths[@]+"${paths[@]}"} "${excludes[@]}" 2>"$errf")
  fi
  rc=$?
  if [ "$rc" -gt 1 ]; then
    echo "suchlauf-nachmessen: git grep endete mit Exit $rc: ${lines[$i]}" >&2
    cat "$errf" >&2
    exit 2
  fi
  if [ -s "$errf" ]; then
    echo "suchlauf-nachmessen: Hinweis von git grep (nicht gezählt): ${lines[$i]}" >&2
    cat "$errf" >&2
  fi
  if [ -z "$out" ]; then ist=0; else ist=$(printf '%s\n' "$out" | wc -l); fi
  if [ "$ist" -eq "${solls[$i]}" ]; then
    verdict=OK
  else
    verdict=ABWEICHUNG
    failures=$((failures + 1))
  fi
  printf '%s  soll=%s ist=%s  %s\n' "$verdict" "${solls[$i]}" "$ist" "${lines[$i]}"
done

if [ "$failures" -ne 0 ]; then
  echo "suchlauf-nachmessen: $failures von ${#lines[@]} Zeilen weichen ab" >&2
  exit 1
fi
echo "suchlauf-nachmessen: ${#lines[@]} Zeilen stimmen"
