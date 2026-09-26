#!/usr/bin/env bash
# suchlauf-nachmessen — misst die `suchlauf`-Blöcke eines Slice-Plans nach:
# je Zeile `<Stand> <Soll> <Argumente von git grep>` führt es
# `git grep -n <Argumente>` am Stand aus (Commit-Kennung, oder `diff` für den
# Arbeitsbaum), zählt die Trefferzeilen und vergleicht sie mit dem Soll.
# Die Plan-Datei ist immer aus dem Suchraum ausgeschlossen, unter ihrem
# Dateinamen in jedem Verzeichnis (der Plan wandert durch die
# Lifecycle-Verzeichnisse). Vertrag, Grenze und Exit-Codes:
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

# Erster Durchgang: jede Zeile prüfen, bevor ein Befehl läuft.
stands=(); solls=(); optss=(); pathss=()
bad=0
for line in "${lines[@]}"; do
  if ! split_words "$line"; then
    echo "suchlauf-nachmessen: Anführungszeichen nicht geschlossen: $line" >&2
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
for i in "${!lines[@]}"; do
  IFS=$'\x1f' read -r -a opts <<<"${optss[$i]}"
  paths=()
  if [ -n "${pathss[$i]}" ]; then IFS=$'\x1f' read -r -a paths <<<"${pathss[$i]}"; fi
  stand=${stands[$i]}
  if [ "$stand" = "diff" ]; then
    out=$(git grep -n "${opts[@]}" -- ${paths[@]+"${paths[@]}"} "${excludes[@]}" 2>&1)
  else
    out=$(git grep -n "${opts[@]}" "$stand" -- ${paths[@]+"${paths[@]}"} "${excludes[@]}" 2>&1)
  fi
  rc=$?
  if [ "$rc" -gt 1 ]; then
    echo "suchlauf-nachmessen: git grep endete mit Exit $rc: ${lines[$i]}" >&2
    printf '%s\n' "$out" >&2
    exit 2
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
