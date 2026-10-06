#!/usr/bin/env bash
# zitat-vergleich — Referent-Messung einer Zitat-Korrektur: vergleicht die
# Einheit, die die alte Adresse an ihrem Stand adressiert, mit der Einheit der
# neuen Adresse an deren Stand; roh, oder mit Tag-Paar nach Normalisierung nur
# des bewegten Baseline-Tags. Einheiten, Normalisierung und die Zusagen der
# Messform folgen ADR-0159; Aufruf, Ausgänge, Host-Werkzeuge, Grenzen und die
# Abweichungen von der Befehlsform im Block der ADR stehen im Vertrag
# harness/targets/zitat-vergleich.md.
#
# Aufruf: zitat-vergleich.sh <alt-stand> <alt-pfad> <alt-ref> <neu-stand> <neu-pfad> <neu-ref> [<alt-tag>:<neu-tag>]
#   ref: '' = ganze Datei · '#anker' = Heading-Abschnitt oder HTML-id · 'L<a>-<b>' = Zeilen
# Exit: 0 gleich · 1 verschieden · 2 keine Einheit, Lokator, Tag-Paar,
#   Argumentzahl oder awk ohne Multibyte
# Die gedruckte Zeile auf stdout ist der Beleg; über make kommt jeder Exit
# ungleich 0 als Exit 2 an.
#
# Mit `source` geladen, definiert die Datei nur die Funktionen (Messläufe
# gegen die Befehlsform der ADR); der Aufruf als Skript führt main aus.
set -uo pipefail

# einheit <stand> <pfad> [<ref>] — druckt die Einheit roh auf stdout; leere
# oder nicht lesbare Einheit, ungültiger Lokator und eine `id` in anderer Form
# als `<a id="X">` enden mit 2 und einer Zeile auf stderr.
einheit() {
  local stand="$1" pfad="$2" ref="${3:-}" out a b rc=0
  out=$(git show "$stand:$pfad" 2>/dev/null && printf .) || { echo "einheit: $stand:$pfad nicht lesbar" >&2; return 2; }
  out="${out%.}"
  case "$ref" in
    "") ;;
    L*) a="${ref#L}"; b="${a#*-}"; a="${a%%-*}"
        [[ "$a" =~ ^[0-9]+$ && "$b" =~ ^[0-9]+$ ]] || { echo "einheit: Lokator $ref" >&2; return 2; }
        out=$(printf '%s' "$out" | sed -n "${a},${b}p" 2>/dev/null && printf .) || { echo "einheit: Lokator $ref" >&2; return 2; }
        out="${out%.}" ;;
    \#*) out=$(printf '%s' "$out" | awk -v want="${ref#\#}" '
          function slug(t, s) { s = tolower(t); gsub(/<[^>]*>/, "", s); gsub(/[^[:alnum:] _-]/, "", s); gsub(/ /, "-", s); return s }
          # fencerun(l): Länge des Fence-Laufs am Zeilenanfang (0 bis 3 Leerzeichen
          # Einzug, mindestens drei gleiche Zeichen ` oder ~), sonst 0; setzt fch
          # (Zeichen) und frest (Rest der Zeile).
          function fencerun(l, ind, c, n) {
            ind = 0
            while (ind < 4 && substr(l, ind + 1, 1) == " ") ind++
            if (ind > 3) return 0
            c = substr(l, ind + 1, 1)
            if (c != "`" && c != "~") return 0
            n = 0
            while (substr(l, ind + 1 + n, 1) == c) n++
            if (n < 3) return 0
            fch = c; frest = substr(l, ind + 1 + n)
            return n
          }
          # idform(l): 0 kein <a id="want" außerhalb von Inline-Code, 1 nur die
          # Form <a id="want">, 2 mindestens eine andere Form (Attribut, />).
          function idform(l, off, p, res) {
            res = 0; off = 0
            while ((p = index(substr(l, off + 1), tag)) > 0) {
              p += off
              if (p == 1 || substr(l, p - 1, 1) != "`") {
                if (substr(l, p + length(tag), 1) == ">") { if (res == 0) res = 1 } else res = 2
              }
              off = p
            }
            return res
          }
          # leer(l): die Zeile trägt nach dem Entfernen aller <a id="…"></a> nur Leerraum.
          function leer(l, r) { r = l; gsub(/<a id="[^"]*"><\/a>/, "", r); return r ~ /^[[:space:]]*$/ }
          BEGIN { tag = "<a id=\"" want "\"" }
          {
            inf = fence
            n = fencerun($0)
            if (!fence) {
              if (n && !(fch == "`" && index(frest, "`"))) { fence = 1; fc = fch; fl = n; inf = 1 }
            } else if (n && fch == fc && n >= fl && frest ~ /^[ \t]*$/) { fence = 0 }
            f = inf ? 0 : idform($0)
            if (f == 2) bad = 1
          }
          done { next }
          !inf && match($0, /^#+ /) && RLENGTH <= 7 {
            lvl = RLENGTH - 1
            if ((mode == "abs" && lvl <= ilvl) || mode == "blk") { done = 1; next }
            if (mode == "vor") { mode = "abs"; ilvl = lvl; next }
            if (mode == "") {
              s = slug(substr($0, lvl + 2)); k = seen[s]++; if (k) s = s "-" k
              if (s == want || f == 1) { mode = "abs"; ilvl = lvl; next }
            }
          }
          mode == "vor" && !inf && leer($0) { next }
          mode == "vor" { mode = "blk" }
          mode == "abs" || mode == "blk" { print; next }
          mode == "" && f == 1 {
            if ($0 ~ /^[[:space:]]*[|]/) { print; done = 1; next }
            if (leer($0)) { mode = "vor"; next }
            mode = "blk"; print
          }
          END { if (bad) exit 3 }' && printf .) || rc=$?
        if [ "$rc" -eq 3 ]; then
          echo "einheit: <a id=\"${ref#\#}\" in anderer Form als <a id=\"${ref#\#}\"> (Attribut oder selbstschließend) in $stand:$pfad, nicht gelesen" >&2
          return 2
        fi
        [ "$rc" -eq 0 ] || { echo "einheit: awk endete mit Exit $rc an $stand:$pfad$ref" >&2; return 2; }
        out="${out%.}" ;;
    *) echo "einheit: Lokator $ref" >&2; return 2 ;;
  esac
  [ -n "$out" ] || { echo "einheit: leere Einheit $stand:$pfad$ref" >&2; return 2; }
  printf '%s' "$out"
}

# tagnorm <tag> — ersetzt nur das Segment /<tag>/ durch /<tag>/-Platzhalter.
tagnorm() { sed -E "s#/${1//./\\.}/#/<tag>/#g"; }

# baumda <stand> <tag> — der Stand trägt das Verzeichnis .harness/baseline/<tag>.
baumda() {
  local treffer
  treffer=$(git ls-tree -d --name-only "$1" -- ".harness/baseline/$2" 2>/dev/null) || return 1
  [ "$treffer" = ".harness/baseline/$2" ]
}

# vergleich <alt-stand> <alt-pfad> <alt-ref> <neu-stand> <neu-pfad> <neu-ref> [<alt-tag>:<neu-tag>]
# — druckt die Zeile vor dem return; Exit 0 gleich, 1 verschieden, 2 keine Einheit
# oder ungültiges Tag-Paar.
vergleich() {
  local x y ec art=roh alt neu
  x=$(einheit "$1" "$2" "$3" && printf .) || { echo "vergleich: $1:$2$3 keine Einheit, Exit 2"; return 2; }
  y=$(einheit "$4" "$5" "$6" && printf .) || { echo "vergleich: $4:$5$6 keine Einheit, Exit 2"; return 2; }
  x="${x%.}"; y="${y%.}"
  if [ $# -ge 7 ]; then
    [[ "$7" =~ ^v[0-9]+\.[0-9]+\.[0-9]+:v[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "vergleich: Tag-Paar $7, Exit 2"; return 2; }
    alt="${7%%:*}"; neu="${7#*:}"
    [ "$alt" != "$neu" ] || { echo "vergleich: Tag-Paar $7 mit gleichen Tags, Exit 2"; return 2; }
    baumda "$1" "$alt" || { echo "vergleich: Tag-Paar $7, $1 trägt .harness/baseline/$alt nicht, Exit 2"; return 2; }
    baumda "$4" "$neu" || { echo "vergleich: Tag-Paar $7, $4 trägt .harness/baseline/$neu nicht, Exit 2"; return 2; }
    art="norm $7"
    x=$(printf '%s' "$x" | tagnorm "$alt" && printf .); x="${x%.}"
    y=$(printf '%s' "$y" | tagnorm "$neu" && printf .); y="${y%.}"
  fi
  if cmp <(printf '%s' "$x") <(printf '%s' "$y"); then ec=0; else ec=$?; fi
  echo "vergleich $art: $1:$2$3 <-> $4:$5$6 cmp $ec"
  return "$ec"
}

main() {
  if [ $# -ne 6 ] && [ $# -ne 7 ]; then
    echo "Aufruf: zitat-vergleich.sh <alt-stand> <alt-pfad> <alt-ref> <neu-stand> <neu-pfad> <neu-ref> [<alt-tag>:<neu-tag>] — $# Argumente, erwartet 6 oder 7; eine Referenz '#anker' steht in Anführungszeichen, Exit 2"
    return 2
  fi
  export LC_ALL=C.UTF-8
  if [ "$(printf 'Ä' | awk '{print tolower($0), length($0)}' 2>/dev/null)" != "ä 1" ]; then
    echo "zitat-vergleich: awk ohne Multibyte/UTF-8 unter LC_ALL=C.UTF-8, Exit 2"
    return 2
  fi
  local root
  root=$(git rev-parse --show-toplevel) || { echo "zitat-vergleich: kein git-Repo, Exit 2"; return 2; }
  cd "$root" || return 2
  vergleich "$@"
}

if [ "${BASH_SOURCE[0]}" = "$0" ]; then
  main "$@"
  exit $?
fi
