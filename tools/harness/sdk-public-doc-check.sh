#!/usr/bin/env bash
# sdk-public-doc-check.sh — prueft, dass die Dateien der SDK-Baeume unter
# sdks/ keine internen Kennungen tragen: keine Kennung der Spezifikations-,
# Entscheidungs- und Anforderungsdokumente des Projekts (SPEC-, ADR-, ARC-,
# LH-FA-/LH-QA-) und keinen Slice-/Welle-Namen. Kommentare, Docstrings, KDoc,
# XML-Doku, Fehlertexte, README und Build-Dateien der SDKs beschreiben, was der
# Code tut; sie erreichen Anwender ueber die Pakete (Docstrings im Wheel, XML-
# Dokumentationsdatei im nupkg, Quellen im Sources-Jar und in der sdist).
#
# Umfang: alle Textdateien unter dem Wurzelverzeichnis (Standard `sdks`), ohne
# Bau-Ausgaben und erzeugten Code (obj, bin, build, dist, .gradle, __pycache__,
# .pytest_cache, *.egg-info, grpc_gen). Die zur Bauzeit erzeugten Python-Stubs
# prueft der Test tests/test_public_text.py des Python-Packages gegen das
# installierte Paket.
#
# Aufruf: `make sdk-public-doc-check` (netzlos, nur grep). Optional ein
# abweichendes Wurzelverzeichnis als erstes Argument (fuer den Tabellentest
# run-sdk-public-doc-check-tests.sh). Exit 0 ohne Treffer, Exit 1 mit Treffern
# (Ausgabe `datei:zeile: text`).
set -euo pipefail

if [ "$#" -ge 1 ]; then
  root=$1
else
  cd "$(git rev-parse --show-toplevel)"
  root=sdks
fi

pattern='\b(SPEC|ADR|ARC)-[0-9]+|\bLH-(FA|QA)-[A-Z]{3}-[0-9]+|\b(slice|welle)-[a-z0-9]'

hits=$(find "$root" \
  \( -name obj -o -name bin -o -name build -o -name dist -o -name .gradle \
     -o -name __pycache__ -o -name .pytest_cache -o -name '*.egg-info' \
     -o -name grpc_gen \) -prune -o -type f -print0 \
  | xargs -0 -r grep -InE "$pattern" || true)

if [ -n "$hits" ]; then
  printf '%s\n' "$hits" >&2
  echo "sdk-public-doc-check: interne Kennung in den SDK-Dateien (siehe oben)" >&2
  exit 1
fi
echo "sdk-public-doc-check: keine interne Kennung unter $root"
