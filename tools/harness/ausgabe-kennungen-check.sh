#!/usr/bin/env bash
# ausgabe-kennungen-check.sh — prueft, dass kein Ausgabe-Literal des Servers
# und seiner Betriebs-Skripte eine interne Kennung (LH-, ADR-, SPEC-, ARC-)
# traegt. Die Ausgabe des Programms ist die Betreiber-Sicht; Kennungen der
# Spezifikations-, Entscheidungs- und Anforderungsdokumente gehoeren in Plan,
# Spec und Commit.
#
# Gegenstand: Produktions-Go (*.go ohne *_test.go und *.pb.go) unter internal/,
# cmd/ und tools/schema/ sowie die echo-/printf-Zeilen der Skripte (*.sh) unter
# tools/schema/ und examples/. Nicht gelesen: Tests, die Laeufer unter
# tools/harness/ (Entwickler-Ausgabe), Kommentare.
#
# Muster, ERE (kein grep -P), zweistufig:
#   Go:    eine Kennung hinter einem ungeschlossenen Anfuehrungszeichen oder
#          Backtick (Literal), und die Zeile beginnt nicht mit //.
#   Shell: ein echo/printf, hinter dem eine Kennung steht, und die Zeile
#          beginnt nicht mit #.
# Grenzen (Vertrag harness/sensors/ausgabe-kennungen-check.md): ein mehrzeiliges
# Raw-String-Literal und ein Heredoc werden nicht gelesen; ein nachgestellter
# Kommentar, der ein Literal mit Kennung zitiert, trifft.
#
# Aufruf: `make ausgabe-kennungen-check` (netzlos, nur bash/git/grep/find).
# Optional eine abweichende Wurzel als erstes Argument (fuer den Tabellentest
# run-ausgabe-kennungen-check-tests.sh). Exit 0 ohne Treffer, Exit 1 mit
# Treffern (Ausgabe `datei:zeile:text` auf stderr), Exit 2 bei Lesefehlern
# (nicht lesbare Datei oder Verzeichnis, fehlende Wurzel, kein Gegenstand);
# nie "sauber" bei einem Lesefehler.
set -euo pipefail

if [ "$#" -ge 1 ]; then
  root=$1
else
  root=$(git rev-parse --show-toplevel)
fi

go_roots=(internal cmd tools/schema)
sh_roots=(tools/schema examples)

id='(LH|ADR|SPEC|ARC)-[A-Z0-9]'
pat_go='^[^"`]*("[^"]*"[^"`]*)*["`][^"`]*'"$id"
pat_sh='(echo|printf)[[:space:]].*'"$id"
skip_go='^[0-9]+:[[:space:]]*//'
skip_sh='^[0-9]+:[[:space:]]*#'

scratch=$(mktemp -d)
trap 'rm -rf "${scratch:?}"' EXIT

# find laeuft vor der Schleife: ein nicht lesbares oder fehlendes Verzeichnis
# ist Exit 2, nicht eine stillschweigend kuerzere Liste.
list_files() { # <Ausgabe> <Muster> <Wurzeln...>
  local out=$1 name=$2 frc=0
  shift 2
  local r
  for r in "$@"; do
    if [ ! -d "$root/$r" ]; then
      echo "ausgabe-kennungen-check: Lesefehler: Verzeichnis fehlt: $r" >&2
      exit 2
    fi
  done
  : > "$out.raw"
  for r in "$@"; do
    find "$root/$r" -type f -name "$name" ! -name '*_test.go' ! -name '*.pb.go' -print0 >> "$out.raw" 2> "$scratch/find.err" || frc=$?
    if [ "$frc" -ne 0 ]; then
      echo "ausgabe-kennungen-check: Lesefehler beim Durchsuchen von $r ($(head -n 1 "$scratch/find.err"))" >&2
      exit 2
    fi
  done
  sort -zu "$out.raw" > "$out"
}

hits=""
nfiles_go=0
nfiles_sh=0

# scan <Datei> <Muster> <Kommentar-Muster>: haengt Treffer an $hits an
scan() {
  local f=$1 pat=$2 skip=$3 grc=0 h="" frc=0 rel
  rel=${f#"$root"/}
  # -a: NUL-Byte als Text. grep-Exit 1 = kein Treffer, >= 2 = Lesefehler.
  h=$(grep -anE -e "$pat" "$f" 2>"$scratch/grep.err") || grc=$?
  if [ "$grc" -ge 2 ]; then
    echo "ausgabe-kennungen-check: Lesefehler: $rel ($(head -n 1 "$scratch/grep.err"))" >&2
    exit 2
  fi
  if [ -n "$h" ]; then
    h=$(printf '%s\n' "$h" | grep -vE -e "$skip") || frc=$?
    if [ "$frc" -ge 2 ]; then
      echo "ausgabe-kennungen-check: Lesefehler beim Filtern: $rel" >&2
      exit 2
    fi
    if [ -n "$h" ]; then
      hits+=$(printf '%s\n' "$h" | sed "s|^|$rel:|")$'\n'
    fi
  fi
}

list_files "$scratch/go.list" '*.go' "${go_roots[@]}"
while IFS= read -r -d '' f; do
  nfiles_go=$((nfiles_go + 1))
  scan "$f" "$pat_go" "$skip_go"
done < "$scratch/go.list"

list_files "$scratch/sh.list" '*.sh' "${sh_roots[@]}"
while IFS= read -r -d '' f; do
  nfiles_sh=$((nfiles_sh + 1))
  scan "$f" "$pat_sh" "$skip_sh"
done < "$scratch/sh.list"

if [ "$nfiles_go" -eq 0 ] || [ "$nfiles_sh" -eq 0 ]; then
  echo "ausgabe-kennungen-check: Lesefehler: kein Gegenstand ($nfiles_go Go-Dateien, $nfiles_sh Skripte); leer ist nicht bestanden" >&2
  exit 2
fi

if [ -n "$hits" ]; then
  printf '%s' "$hits" >&2
  echo "ausgabe-kennungen-check: interne Kennung in einem Ausgabe-Literal (siehe oben)" >&2
  exit 1
fi
echo "ausgabe-kennungen-check: keine interne Kennung in Ausgabe-Literalen von $nfiles_go Go-Dateien und $nfiles_sh Skripten"
