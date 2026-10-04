#!/usr/bin/env bash
# run-sdk-kompat-tests.sh — Kompatibilitätsmessung der SDK-Packages 0.6.0
# gegenüber 0.5.0 (make test-sdk-kompat; Vertrag: harness/targets/sdk-kompat.md).
# Je Sprache (C#, Kotlin, Python) fährt ein Gast-Programm unter
# tools/harness/sdk-kompat/<sprache>/ jede öffentliche Konstruktor- und
# Aufrufform der 0.5.x-Fehlertypen beider Hierarchien (HTTP, gRPC):
#   A1  Gast gegen das veröffentlichte 0.5.0 gebaut und gelaufen (Grundlinie)
#   A2  dieselben Gast-Binärdateien gegen die Bibliothek 0.6.0 (Austausch der
#       Bibliotheksdatei, kein Neubau)
#   A3  Gegenrichtung: ein gegen 0.6.0 gebauter Gast gegen 0.5.0 muss scheitern
#   A5  Quellseite: derselbe Quelltext gegen beide Versionen übersetzt bzw. die
#       Signaturen gelesen; in C# die null-Matrix
# A4 (Mutationsprobe) ist dieselbe Messung mit einer mutierten Bibliothek 0.6.0
# über SDK_KOMPAT_DIST_<SPRACHE>.
#
# Die Basis jeder Sprache ist das lokal getaggte Image der Stufe `build` des
# SDK-Dockerfiles (kein eigenes Digest-Literal). Der Gast-Bau übersetzt nur; die
# Läufe geschehen per `docker run`: der Docker-Cache überspringt einen Lauf in
# einer RUN-Schicht still. Ein Lauf ohne gedruckte Zeile eines Schritts ist rot.
#
# Aufruf: `make test-sdk-kompat`. Overrides: SDK_KOMPAT_NEU (dist = Artefakte
# von make sdk-pack-*, Standard; registry = veröffentlichte 0.6.0-Pakete),
# SDK_KOMPAT_DIST_CSHARP / _KOTLIN / _PYTHON (Verzeichnis statt sdks/<sprache>/dist),
# SDK_KOMPAT_SPRACHEN (Teilmenge, Standard "csharp kotlin python"). Exit 0 nur,
# wenn jede Sprache jeden Schritt wie erwartet bestanden hat. Kein Gate:
# Paketbezug braucht Netz.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

SDK_KOMPAT_NEU=${SDK_KOMPAT_NEU:-dist}
SDK_KOMPAT_SPRACHEN=${SDK_KOMPAT_SPRACHEN:-csharp kotlin python}

case "$SDK_KOMPAT_NEU" in
  dist | registry) ;;
  *)
    echo "run-sdk-kompat-tests: SDK_KOMPAT_NEU=$SDK_KOMPAT_NEU unbekannt (dist oder registry)" >&2
    exit 2
    ;;
esac

scratch=$(mktemp -d)
cleanup() {
  rm -rf -- "$scratch"
}
trap cleanup EXIT

artefakt_muster() {
  case "$1" in
    csharp) printf '*.nupkg' ;;
    kotlin) printf 'pgchangefeed-kotlin-*.jar' ;;
    python) printf 'pgchangefeed-*.whl' ;;
  esac
}

# neu_verzeichnis <sprache> — das Verzeichnis, das als Bau-Kontext `neu` in das Gast-Image kommt.
neu_verzeichnis() {
  local sprache=$1 var dir
  if [ "$SDK_KOMPAT_NEU" = registry ]; then
    dir="$scratch/leer-$sprache"
    mkdir -p "$dir"
    : >"$dir/leer"
    printf '%s' "$dir"
    return 0
  fi
  var="SDK_KOMPAT_DIST_$(printf '%s' "$sprache" | tr '[:lower:]' '[:upper:]')"
  dir=${!var:-sdks/$sprache/dist}
  if ! compgen -G "$dir/$(artefakt_muster "$sprache")" >/dev/null; then
    echo "run-sdk-kompat-tests: $dir trägt kein Artefakt der Sprache $sprache — make sdk-pack-$sprache vorher (oder SDK_KOMPAT_NEU=registry)" >&2
    exit 2
  fi
  printf '%s' "$dir"
}

alle_ok=1
for sprache in $SDK_KOMPAT_SPRACHEN; do
  case "$sprache" in
    csharp | kotlin | python) ;;
    *)
      echo "run-sdk-kompat-tests: Sprache $sprache unbekannt (csharp, kotlin, python)" >&2
      exit 2
      ;;
  esac
  neu=$(neu_verzeichnis "$sprache")
  basis="pg-change-feed:sdk-kompat-base-$sprache"
  gast="pg-change-feed:sdk-kompat-$sprache"

  echo "run-sdk-kompat-tests: $sprache — Basis (Stufe build von sdks/$sprache), Gast-Bau, Modus $SDK_KOMPAT_NEU"
  docker build --build-context proto=proto --target build -t "$basis" "sdks/$sprache" >/dev/null
  docker build --build-arg BASE="$basis" --build-context "neu=$neu" -t "$gast" "tools/harness/sdk-kompat/$sprache" >/dev/null

  log="$scratch/lauf-$sprache.log"
  rc=0
  docker run --rm "$gast" >"$log" 2>&1 || rc=$?
  cat "$log"

  # Jeder Schritt muss mindestens eine gedruckte Zeile tragen.
  for schritt in A1 A2 A3 A5; do
    if ! grep -q "^KOMPAT $sprache $schritt[ :-]" "$log"; then
      echo "run-sdk-kompat-tests: $sprache $schritt — keine gedruckte Zeile" >&2
      rc=1
    fi
  done
  if [ "$rc" -ne 0 ]; then
    echo "run-sdk-kompat-tests: $sprache ROT (Ausgang $rc)" >&2
    alle_ok=0
  else
    echo "run-sdk-kompat-tests: $sprache grün — A1, A2, A3, A5 wie erwartet"
  fi
done

if [ "$alle_ok" -ne 1 ]; then
  echo "run-sdk-kompat-tests: mindestens eine Sprache weicht von der Erwartung ab" >&2
  exit 1
fi
echo "run-sdk-kompat-tests: Kompatibilitätsmessung ($SDK_KOMPAT_NEU) grün für: $SDK_KOMPAT_SPRACHEN"
