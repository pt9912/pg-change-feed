#!/usr/bin/env bash
# run-handbuch-public-doc-check-tests.sh — Tabellentest gegen
# handbuch-public-doc-check.sh: je Kennungsart und je Link-Klasse ein Treffer
# (Exit 1), Wortrand-Faelle ohne Treffer, saubere Datei (Exit 0), die
# Klassifikation (unklassifizierte und fehlende Datei: Exit 2) und die
# ausgenommenen Dateien (Exit 0). Netzlos, kein Docker noetig.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

fail=0
cases=0
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

all_files="benutzerhandbuch.md benutzerhandbuch-standard.md version.md releasing.md bench-abdeckung.md ci-matrix-abdeckung.md e2e-abdeckung.md sdk-e2e-abdeckung.md"

# fresh: Wurzel mit allen acht Dateien, sauberer Inhalt
fresh() {
  rm -rf "$tmp/root"
  mkdir -p "$tmp/root/docs/user"
  local f
  for f in $all_files; do
    printf '%s\n' 'Der Betreiber startet den Feed-Container.' > "$tmp/root/docs/user/$f"
  done
}

run() {
  bash tools/harness/handbuch-public-doc-check.sh "$tmp/root" > /dev/null 2>&1
  echo $?
}

check() {
  local want=$1 what=$2 got=$3
  cases=$((cases + 1))
  if [ "$got" -ne "$want" ]; then
    echo "FEHLER: $what -> Exit $got, erwartet $want" >&2
    fail=1
  fi
}

# expect <Exit> <Beschreibung> <Datei unter docs/user/> <Inhalt>
expect() {
  fresh
  printf '%s\n' "$4" > "$tmp/root/docs/user/$3"
  check "$1" "$2" "$(run)"
}

fresh
check 0 "saubere Wurzel" "$(run)"

expect 1 "LH-Kennung" benutzerhandbuch.md 'Anforderung LH-FA-SST-009 ist erfuellt.'
expect 1 "LH-Kennung (QA)" benutzerhandbuch.md 'Siehe LH-QA-OPS-001.'
expect 1 "ADR-Kennung" benutzerhandbuch.md 'Siehe ADR-0134.'
expect 1 "SPEC-Kennung" benutzerhandbuch-standard.md 'Siehe SPEC-022.'
expect 1 "ARC-Kennung" version.md 'ARC-004'
expect 1 "BEO-Kennung" benutzerhandbuch.md 'Register BEO-PGC/x'
expect 1 "MR-Kennung" benutzerhandbuch.md 'Adaption MR-001'
expect 1 "CO-Kennung" benutzerhandbuch.md 'Carveout CO-123'
expect 1 "Slice-Name" benutzerhandbuch.md 'neu in slice-foo'
expect 1 "Slice-Name in Klammern" benutzerhandbuch.md '(slice-1)'
expect 1 "Welle-Name am Zeilenanfang" benutzerhandbuch.md 'welle-15 brachte das'
expect 1 "Kennung in Fenced-Block" benutzerhandbuch.md $'```\nADR-0134\n```'
expect 1 "Link nach docs/plan/" benutzerhandbuch.md 'Siehe [x](../plan/adr/y.md).'
expect 1 "Link nach docs/reviews/" benutzerhandbuch.md 'Siehe docs/reviews/r.md.'
expect 1 "relativer Link nach ../reviews/" benutzerhandbuch-standard.md '[r](../reviews/r.md)'
expect 1 "Kennung in version.md" version.md 'slice-x'

expect 0 "byte-slice ohne Treffer" benutzerhandbuch.md 'the byte-slice-x is a word'
expect 0 "CO-2 Emissionen ohne Treffer" benutzerhandbuch.md 'CO-2 Emissionen'
expect 0 "Slice-1 (Grossschreibung) ohne Treffer" benutzerhandbuch.md 'Slice-1'
expect 0 "Tag-Muster sdk-csharp-v ohne Treffer" benutzerhandbuch.md 'Tag sdk-csharp-v0.2.0'
expect 0 "Kennung in releasing.md ausgenommen" releasing.md 'ADR-0051 slice-x'
expect 0 "Kennung in e2e-abdeckung.md ausgenommen" e2e-abdeckung.md 'LH-FA-SST-009'
expect 0 "Kennung in bench-abdeckung.md ausgenommen" bench-abdeckung.md 'LH-QA-PER-001 docs/plan/x'
expect 0 "Kennung in ci-matrix-abdeckung.md ausgenommen" ci-matrix-abdeckung.md 'LH-QA-POR-001'
expect 0 "Kennung in sdk-e2e-abdeckung.md ausgenommen" sdk-e2e-abdeckung.md 'LH-FA-SST-009'

fresh
printf '%s\n' 'Neue Datei.' > "$tmp/root/docs/user/neu.md"
check 2 "unklassifizierte .md" "$(run)"

fresh
mkdir -p "$tmp/root/docs/user/unter"
printf '%s\n' 'Neue Datei.' > "$tmp/root/docs/user/unter/neu.md"
check 2 "unklassifizierte .md im Unterverzeichnis" "$(run)"

fresh
printf '%s\n' 'bild' > "$tmp/root/docs/user/bild.png"
check 0 "Nicht-.md-Datei ist nicht Gegenstand" "$(run)"

fresh
rm "$tmp/root/docs/user/version.md"
check 2 "fehlende geprueft-genannte Datei" "$(run)"

fresh
rm "$tmp/root/docs/user/releasing.md"
check 2 "fehlende ausgenommen-genannte Datei" "$(run)"

if [ "$fail" -eq 0 ]; then
  echo "run-handbuch-public-doc-check-tests: alle $cases Fälle bestanden"
else
  echo "run-handbuch-public-doc-check-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
