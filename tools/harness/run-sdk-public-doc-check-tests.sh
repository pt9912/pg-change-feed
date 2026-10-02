#!/usr/bin/env bash
# run-sdk-public-doc-check-tests.sh — Tabellentest gegen
# sdk-public-doc-check.sh: je Kennungsart eine Datei mit Treffer (Exit 1), die
# Ausnahmen (Bau-Ausgaben, erzeugter Code, saubere Datei: Exit 0), eine Datei
# mit NUL-Byte und Lesefehler (Exit 2). Je Exit-Klasse prueft mindestens ein
# Fall den Meldungstext. Die Lesefehler-Faelle laufen nur bei uid != 0, weil
# chmod 000 root nicht hindert. Netzlos, kein Docker noetig.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

fail=0
cases=0
tmp=$(mktemp -d)
: "${tmp:?}"
trap 'chmod -R u+rwx "$tmp" 2>/dev/null; rm -rf "${tmp:?}"' EXIT

# fresh: leere Wurzel
fresh() {
  chmod -R u+rwx "${tmp:?}/root" 2>/dev/null
  rm -rf "${tmp:?}/root"
  mkdir -p "$tmp/root"
}

# run: gibt den Exit aus, legt stdout+stderr in $tmp/out ab
run() {
  bash tools/harness/sdk-public-doc-check.sh "$tmp/root" > "$tmp/out" 2>&1
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

# msg <Beschreibung> <ERE>: die letzte Ausgabe traegt den Meldungstext
msg() {
  cases=$((cases + 1))
  if ! grep -qE -e "$2" "$tmp/out"; then
    echo "FEHLER: $1 -> Meldung ohne /$2/: $(head -c 300 "$tmp/out")" >&2
    fail=1
  fi
}

# expect <0|1> <Beschreibung> <relativer Dateipfad> <Inhalt> [<Meldungs-ERE>]
expect() {
  local want=$1 what=$2 path=$3 content=$4
  fresh
  mkdir -p "$tmp/root/$(dirname "$path")"
  printf '%s\n' "$content" > "$tmp/root/$path"
  check "$want" "$what" "$(run)"
  if [ "$#" -ge 5 ]; then msg "$what" "$5"; fi
}

expect 0 "saubere Datei" a/client.py '"""Reads changes from the server."""' \
  '^sdk-public-doc-check: keine interne Kennung unter '
expect 1 "SPEC-Kennung im Python-Docstring" a/client.py '"""Reads changes (SPEC-022)."""' \
  '/a/client\.py:1:.*SPEC-022'
expect 1 "ADR-Kennung im C#-XML-Kommentar" a/Client.cs '/// <summary>See ADR-0106.</summary>'
expect 1 "ARC-Kennung im Kotlin-KDoc" a/Client.kt '/** ARC-004 */'
expect 1 "LH-Kennung in der README" a/README.md 'Requirement LH-FA-SST-009 is met.'
expect 1 "Slice-Name im Build-Kommentar" a/build.gradle.kts '// added by slice-sdk-readme' \
  '^sdk-public-doc-check: interne Kennung in den SDK-Dateien \(siehe oben\)$'
expect 1 "Welle-Name im Kommentar" a/Dockerfile '# welle-15 leftover'
expect 0 "Kennung in einer Bau-Ausgabe (obj)" a/obj/x.cs '// SPEC-018'
expect 0 "Kennung in einer Bau-Ausgabe (dist)" a/dist/x.txt 'ADR-0106'
expect 0 "Kennung im erzeugten Python-Code (grpc_gen)" a/grpc_gen/x.py '# SPEC-020'
expect 0 "Wort mit Bindestrich ohne Kennung" a/x.md 'the byte-slice is a spec-like word'

# Datei mit NUL-Byte: wird als Text gelesen, der Treffer bleibt sichtbar
fresh
mkdir -p "$tmp/root/a"
printf 'a\0b\nSiehe ADR-0134.\n' > "$tmp/root/a/data.txt"
check 1 "Kennung in Datei mit NUL-Byte" "$(run)"

# Lesefehler: als root sind Rechte wirkungslos, der Fall entfaellt dort
if [ "$(id -u)" -ne 0 ]; then
  fresh
  mkdir -p "$tmp/root/a"
  printf '%s\n' 'Siehe ADR-0134.' > "$tmp/root/a/README.md"
  chmod 000 "$tmp/root/a/README.md"
  check 2 "nicht lesbare Datei" "$(run)"
  msg "nicht lesbare Datei" '^sdk-public-doc-check: Lesefehler: .*/a/README\.md '

  fresh
  mkdir -p "$tmp/root/a/unter"
  chmod 000 "$tmp/root/a/unter"
  check 2 "nicht lesbares Unterverzeichnis" "$(run)"
  msg "nicht lesbares Unterverzeichnis" '^sdk-public-doc-check: Lesefehler beim Durchsuchen von .* \('
else
  echo "run-sdk-public-doc-check-tests: Lesefehler-Faelle uebersprungen (root liest trotz chmod 000)"
fi

if [ "$fail" -eq 0 ]; then
  echo "run-sdk-public-doc-check-tests: alle $cases Fälle bestanden"
else
  echo "run-sdk-public-doc-check-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
