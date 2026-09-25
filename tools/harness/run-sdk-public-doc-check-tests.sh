#!/usr/bin/env bash
# run-sdk-public-doc-check-tests.sh — Tabellentest gegen
# sdk-public-doc-check.sh: je Kennungsart eine Datei mit Treffer (Exit 1) und
# die Ausnahmen (Bau-Ausgaben, erzeugter Code, saubere Datei: Exit 0).
# Netzlos, kein Docker noetig.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

fail=0
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# expect <0|1> <Beschreibung> <relativer Dateipfad> <Inhalt>
expect() {
  local want=$1 what=$2 path=$3 content=$4 got
  rm -rf "$tmp/root"
  mkdir -p "$tmp/root/$(dirname "$path")"
  printf '%s\n' "$content" > "$tmp/root/$path"
  bash tools/harness/sdk-public-doc-check.sh "$tmp/root" > /dev/null 2>&1
  got=$?
  if [ "$got" -ne "$want" ]; then
    echo "FEHLER: $what -> Exit $got, erwartet $want" >&2
    fail=1
  fi
}

expect 0 "saubere Datei" a/client.py '"""Reads changes from the server."""'
expect 1 "SPEC-Kennung im Python-Docstring" a/client.py '"""Reads changes (SPEC-022)."""'
expect 1 "ADR-Kennung im C#-XML-Kommentar" a/Client.cs '/// <summary>See ADR-0106.</summary>'
expect 1 "ARC-Kennung im Kotlin-KDoc" a/Client.kt '/** ARC-004 */'
expect 1 "LH-Kennung in der README" a/README.md 'Requirement LH-FA-SST-009 is met.'
expect 1 "Slice-Name im Build-Kommentar" a/build.gradle.kts '// added by slice-sdk-readme'
expect 1 "Welle-Name im Kommentar" a/Dockerfile '# welle-15 leftover'
expect 0 "Kennung in einer Bau-Ausgabe (obj)" a/obj/x.cs '// SPEC-018'
expect 0 "Kennung in einer Bau-Ausgabe (dist)" a/dist/x.txt 'ADR-0106'
expect 0 "Kennung im erzeugten Python-Code (grpc_gen)" a/grpc_gen/x.py '# SPEC-020'
expect 0 "Wort mit Bindestrich ohne Kennung" a/x.md 'the byte-slice is a spec-like word'

if [ "$fail" -eq 0 ]; then
  echo "run-sdk-public-doc-check-tests: alle Fälle bestanden"
else
  echo "run-sdk-public-doc-check-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
