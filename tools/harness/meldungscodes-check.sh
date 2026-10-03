#!/usr/bin/env bash
# meldungscodes-check.sh — gleicht die Meldungscodes (PCF-<S><NNNN>) zwischen
# Quelltext, Code-Tabelle und Handbuch-Katalog ab. Die Tabelle
# (internal/domain/messagecode/codes.go) ist die Quelle der Wahrheit; der
# Katalog im Benutzerhandbuch fuehrt dieselbe Menge; der Quelltext verwendet nur
# Codes der Tabelle.
#
# Geprueft wird (Menge der Codes, nicht ihr Sinn):
#   1. Jeder Token PCF-... im Quelltext, in der Tabelle und im Handbuch hat die
#      Form PCF-[EWI][0-9]{4}; eine andere Form ist ein Befund.
#   2. Jeder Code im Quelltext (Produktions-Go unter internal/, cmd/ und
#      tools/schema/, Skripte *.sh unter tools/schema/ und examples/) und im
#      Handbuch steht in der Tabelle.
#   3. Die Codes der Tabelle sind gleich den Codes der Katalog-Zeilen des
#      Handbuchs (Tabellenzeilen, deren erste Zelle einen Code nennt).
# Ein zurueckgezogener Code bleibt in Tabelle und Katalog und wird wie jeder
# andere Code verglichen.
#
# Gegenstand ausgenommen: Tests (*_test.go), erzeugter Code (*.pb.go) und das
# Paket internal/domain/messagecode selbst (traegt die Tabelle und die Form).
#
# Aufruf: `make meldungscodes-check` (netzlos, nur bash/git/grep/find/sed/sort/
# comm). Optional eine abweichende Wurzel als erstes Argument (Tabellentest
# run-meldungscodes-check-tests.sh). Exit 0 gleich, Exit 1 mit Befunden (auf
# stderr), Exit 2 bei Lesefehlern (fehlende oder nicht lesbare Datei, fehlende
# Wurzel, leerer Gegenstand, Tabelle ohne Code); nie "gleich" bei einem
# Lesefehler.
set -euo pipefail

if [ "$#" -ge 1 ]; then
  root=$1
else
  root=$(git rev-parse --show-toplevel)
fi

table_dir=internal/domain/messagecode
table_rel=$table_dir/codes.go
catalog_rel=docs/user/benutzerhandbuch.md
go_roots=(internal cmd tools/schema)
sh_roots=(tools/schema examples)

scratch=$(mktemp -d)
trap 'rm -rf "${scratch:?}"' EXIT

die2() {
  echo "meldungscodes-check: Lesefehler: $*" >&2
  exit 2
}

[ -d "$root" ] || die2 "Wurzel fehlt: $root"
[ -f "$root/$table_rel" ] || die2 "Code-Tabelle fehlt: $table_rel"
[ -f "$root/$catalog_rel" ] || die2 "Katalog fehlt: $catalog_rel"

# list_files <Ausgabe> <Muster> <Wurzeln...>: NUL-getrennte, sortierte Liste;
# find laeuft vor der Schleife, damit ein Lesefehler Exit 2 ist.
list_files() {
  local out=$1 name=$2 r frc=0
  shift 2
  for r in "$@"; do
    [ -d "$root/$r" ] || die2 "Verzeichnis fehlt: $r"
  done
  : > "$out.raw"
  for r in "$@"; do
    frc=0
    find "$root/$r" -type f -name "$name" ! -name '*_test.go' ! -name '*.pb.go' \
      ! -path "$root/$table_dir/*" -print0 >> "$out.raw" 2> "$scratch/find.err" || frc=$?
    if [ "$frc" -ne 0 ]; then
      die2 "Durchsuchen von $r ($(head -n 1 "$scratch/find.err"))"
    fi
  done
  sort -zu "$out.raw" > "$out"
}

findings=""
add_finding() { findings+="$1"$'\n'; }

# tokens <Datei> <Praefix-Ausgabe>: haengt "datei:zeile:code" je wohlgeformtem
# Token an $2.codes und "datei:zeile:token" je fehlgeformtem an $2.bad.
tokens() {
  local f=$1 out=$2 rel hits="" grc=0 line tok lineno
  rel=${f#"$root"/}
  hits=$(grep -aon -e 'PCF-[A-Za-z0-9]*' "$f" 2> "$scratch/grep.err") || grc=$?
  if [ "$grc" -ge 2 ]; then
    die2 "$rel ($(head -n 1 "$scratch/grep.err"))"
  fi
  [ -n "$hits" ] || return 0
  while IFS= read -r line; do
    lineno=${line%%:*}
    tok=${line#*:}
    if [[ $tok =~ ^PCF-[EWI][0-9]{4}$ ]]; then
      printf '%s:%s:%s\n' "$rel" "$lineno" "$tok" >> "$out.codes"
    else
      printf '%s:%s:%s\n' "$rel" "$lineno" "$tok" >> "$out.bad"
    fi
  done <<< "$hits"
}

: > "$scratch/table.codes"; : > "$scratch/table.bad"
: > "$scratch/src.codes"; : > "$scratch/src.bad"
: > "$scratch/doc.codes"; : > "$scratch/doc.bad"

tokens "$root/$table_rel" "$scratch/table"
if [ ! -s "$scratch/table.codes" ]; then
  die2 "Code-Tabelle $table_rel trägt keinen Code; leer ist nicht bestanden"
fi

nfiles_go=0
nfiles_sh=0
list_files "$scratch/go.list" '*.go' "${go_roots[@]}"
while IFS= read -r -d '' f; do
  nfiles_go=$((nfiles_go + 1))
  tokens "$f" "$scratch/src"
done < "$scratch/go.list"
list_files "$scratch/sh.list" '*.sh' "${sh_roots[@]}"
while IFS= read -r -d '' f; do
  nfiles_sh=$((nfiles_sh + 1))
  tokens "$f" "$scratch/src"
done < "$scratch/sh.list"
if [ "$nfiles_go" -eq 0 ] || [ "$nfiles_sh" -eq 0 ]; then
  die2 "kein Gegenstand ($nfiles_go Go-Dateien, $nfiles_sh Skripte); leer ist nicht bestanden"
fi

tokens "$root/$catalog_rel" "$scratch/doc"

# Katalog-Zeilen: Tabellenzeilen, deren erste Zelle (optional in Backticks) mit
# einem Code beginnt; je Zeile gilt der erste Code.
src_rc=0
sed -n 's/^|[[:space:]]*`\{0,1\}\(PCF-[EWI][0-9]\{4\}\).*/\1/p' "$root/$catalog_rel" > "$scratch/catalog.raw" 2> "$scratch/sed.err" || src_rc=$?
[ "$src_rc" -eq 0 ] || die2 "$catalog_rel ($(head -n 1 "$scratch/sed.err"))"

src_rc=0
sed 's/^[^:]*:[0-9]*://' "$scratch/table.codes" > "$scratch/table.raw" 2> "$scratch/sed.err" || src_rc=$?
[ "$src_rc" -eq 0 ] || die2 "Mengenbildung der Tabelle ($(head -n 1 "$scratch/sed.err"))"
sort -u "$scratch/table.raw" > "$scratch/table.set" || die2 "Mengenbildung der Tabelle (sort)"
sort -u "$scratch/catalog.raw" > "$scratch/catalog.set" || die2 "Mengenbildung des Katalogs (sort)"

# 1. Form
for kind in table src doc; do
  while IFS= read -r line; do
    add_finding "Code in falscher Form: $line (erwartet PCF-[EWI][0-9]{4})"
  done < "$scratch/$kind.bad"
done

# 2. Quelltext und Handbuch nur mit Codes der Tabelle
for kind in src doc; do
  while IFS= read -r line; do
    code=${line##*:}
    grc=0
    grep -qxF -e "$code" "$scratch/table.set" || grc=$?
    if [ "$grc" -ge 2 ]; then
      die2 "Mengenvergleich gegen die Tabelle nicht lesbar"
    fi
    if [ "$grc" -eq 1 ]; then
      add_finding "Code ohne Eintrag in der Tabelle $table_rel: $line"
    fi
  done < "$scratch/$kind.codes"
done

# 3. Tabelle gleich Katalog
comm -23 "$scratch/table.set" "$scratch/catalog.set" > "$scratch/only-table" 2> "$scratch/comm.err" ||
  die2 "Mengenvergleich Tabelle gegen Katalog ($(head -n 1 "$scratch/comm.err"))"
comm -13 "$scratch/table.set" "$scratch/catalog.set" > "$scratch/only-catalog" 2> "$scratch/comm.err" ||
  die2 "Mengenvergleich Katalog gegen Tabelle ($(head -n 1 "$scratch/comm.err"))"
while IFS= read -r code; do
  add_finding "Tabellen-Code ohne Katalog-Zeile in $catalog_rel: $code"
done < "$scratch/only-table"
while IFS= read -r code; do
  add_finding "Katalog-Zeile ohne Eintrag in der Tabelle $table_rel: $code"
done < "$scratch/only-catalog"

if [ -n "$findings" ]; then
  printf '%s' "$findings" >&2
  echo "meldungscodes-check: Quelltext, Tabelle und Katalog sind nicht gleich (siehe oben)" >&2
  exit 1
fi
count=$(wc -l < "$scratch/table.set")
echo "meldungscodes-check: $((count)) Codes in Tabelle und Katalog gleich, Quelltext ($nfiles_go Go-Dateien, $nfiles_sh Skripte) nur mit Codes der Tabelle"
