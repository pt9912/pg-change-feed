#!/usr/bin/env bash
# fmt-check — Aufrufer von `make fmt-check`: meldet jede Go-Datei unter dem
# Verzeichnis, die `gofmt -l` im gepinnten Toolchain-Image nicht als formatiert
# führt (der Container läuft ohne Netz, das Verzeichnis lesend gemountet, es wird
# nichts geschrieben).
# `gofmt -l` endet auch bei Abweichung mit Exit 0; das Werkzeug wertet deshalb
# die Ausgabe aus, nicht den Exit-Code des Formatierers. Vertrag und Exit-Codes:
# harness/sensors/fmt-check.md.
#
# Aufruf: fmt-check.sh [Verzeichnis]   (Standard: die Repo-Wurzel)
# Umgebung: TOOLCHAIN_IMAGE (Pflicht, setzt das Makefile)
# Exit: 0 jede Go-Datei formatiert · 1 mindestens eine Datei weicht ab (Ausgabe:
#       ein Pfad je Zeile) · 2 Eingabe- oder Formatierer-Fehler (Syntaxfehler,
#       Docker-Fehler, Verzeichnis ohne Go-Datei)
set -uo pipefail

if [ $# -gt 1 ]; then
  echo "Aufruf: fmt-check.sh [Verzeichnis]" >&2
  exit 2
fi
if [ -z "${TOOLCHAIN_IMAGE:-}" ]; then
  echo "fmt-check: TOOLCHAIN_IMAGE fehlt (das Makefile setzt es)" >&2
  exit 2
fi

if [ $# -eq 1 ]; then
  dir=$1
else
  dir=$(git rev-parse --show-toplevel) || exit 2
fi
if [ ! -d "$dir" ]; then
  echo "fmt-check: '$dir' ist kein Verzeichnis" >&2
  exit 2
fi
dir=$(realpath -- "$dir") || exit 2

# Das Skript im Container zählt zuerst die Go-Dateien: ein Verzeichnis ohne
# Go-Datei endet mit Exit 2, nicht mit Exit 0, und ein Exit 0 sagt zu, dass
# mindestens eine Go-Datei geprüft wurde. Danach formatiert `gofmt -l` den
# ganzen Baum. Die Zählung nimmt dieselben Dateien wie `gofmt`: jeden Eintrag,
# der kein Verzeichnis ist (auch einen Symlink), mit Endung `.go` und einem
# Namen ohne Punkt am Anfang.
inner='
n=$(find . ! -type d -name "*.go" ! -name ".*" | wc -l)
if [ "$n" -eq 0 ]; then
  echo "fmt-check: keine Go-Datei im Verzeichnis (leer ist nicht bestanden)" >&2
  exit 2
fi
out=$(gofmt -l .)
rc=$?
if [ "$rc" -ne 0 ]; then
  exit 2
fi
if [ -n "$out" ]; then
  printf "%s\n" "$out"
  echo "fmt-check: $n Go-Dateien geprüft, mindestens eine weicht ab" >&2
  exit 1
fi
echo "fmt-check: $n Go-Dateien geprüft, alle formatiert"
'

docker run --rm --network none -v "$dir":/src:ro -w /src "$TOOLCHAIN_IMAGE" sh -c "$inner"
rc=$?
case "$rc" in
  0 | 1) exit "$rc" ;;
  *)
    echo "fmt-check: Lauf endet mit Exit $rc (Syntaxfehler, Verzeichnis ohne Go-Datei oder Docker-Fehler)" >&2
    exit 2
    ;;
esac
