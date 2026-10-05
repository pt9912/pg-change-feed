#!/usr/bin/env bash
# run-sdk-kompat-tests.sh — Kompatibilitätsmessung der SDK-Packages mit
# Meldungscode gegenüber 0.5.0 (make test-sdk-kompat; Vertrag:
# harness/targets/sdk-kompat.md). Je Sprache (C#, Kotlin, Python) fährt ein
# Gast-Programm unter tools/harness/sdk-kompat/<sprache>/ jede öffentliche
# Konstruktor- und Aufrufform der 0.5.x-Fehlertypen beider Hierarchien (HTTP, gRPC):
#   A1  Gast gegen das veröffentlichte 0.5.0 gebaut und gelaufen (Grundlinie)
#   A2  dieselben Gast-Binärdateien gegen die neue Bibliothek (Austausch der
#       Bibliotheksdatei, kein Neubau)
#   A3  Gegenrichtung: ein gegen die Bibliothek mit Meldungscode gebauter Gast
#       gegen 0.5.0 muss scheitern
#   A5  Quellseite: derselbe Quelltext gegen beide Versionen übersetzt bzw. die
#       Signaturen gelesen; in C# die null-Matrix
# A4 (Mutationsprobe) ist dieselbe Messung mit einer mutierten neuen Bibliothek
# über SDK_KOMPAT_DIST_<SPRACHE>.
#
# Versionen: Im Modus dist ist die neue Bibliothek die Version aus der
# Version-Datei des SDK (version_datei unten); liegt im Artefakt-Verzeichnis
# kein Artefakt genau dieser Version, endet der Lauf mit Exit 2 vor jedem Bau.
# REGISTRY_VERSION ist die veröffentlichte Version, gegen die A3 und A5 (C#,
# Kotlin) übersetzen und die der Modus registry in A2 misst.
#
# Die Basis jeder Sprache ist das lokal getaggte Image der Stufe `build` des
# SDK-Dockerfiles (kein eigenes Digest-Literal). Der Gast-Bau übersetzt nur; die
# Läufe geschehen per `docker run`: der Docker-Cache überspringt einen Lauf in
# einer RUN-Schicht still. Ein Lauf ohne gedruckte Zeile eines Schritts ist rot.
#
# Aufruf: `make test-sdk-kompat`. Overrides: SDK_KOMPAT_NEU (dist = Artefakte
# von make sdk-pack-*, Standard; registry = veröffentlichte Pakete REGISTRY_VERSION),
# SDK_KOMPAT_DIST_CSHARP / _KOTLIN / _PYTHON (Verzeichnis statt sdks/<sprache>/dist),
# SDK_KOMPAT_SPRACHEN (Teilmenge, Standard "csharp kotlin python"). Exit 0 nur,
# wenn jede Sprache jeden Schritt wie erwartet bestanden hat. Kein Gate:
# Paketbezug braucht Netz.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

REGISTRY_VERSION=0.6.0

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

# version_datei <sprache> — die Datei, deren Versionszeile die Version des Packages führt.
version_datei() {
  case "$1" in
    csharp) printf 'sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj' ;;
    kotlin) printf 'sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts' ;;
    python) printf 'sdks/python/pgchangefeed/pyproject.toml' ;;
  esac
}

# version_lesen <sprache> — die Version aus der Version-Datei. Exit 2, wenn die
# Datei fehlt oder nicht genau eine Versionszeile in der Form X.Y.Z (drei
# Ziffernfolgen, kein Suffix wie -rc.1) trägt; es
# gibt keinen Ersatzwert.
version_lesen() {
  local sprache=$1 datei muster werte
  datei=$(version_datei "$sprache")
  case "$sprache" in
    csharp) muster='s:^ *<Version>\(.*\)</Version> *$:\1:p' ;;
    *) muster='s/^version = "\(.*\)"$/\1/p' ;;
  esac
  if [ ! -f "$datei" ]; then
    echo "run-sdk-kompat-tests: Version-Datei $datei der Sprache $sprache fehlt" >&2
    exit 2
  fi
  werte=$(sed -n "$muster" "$datei")
  if [ "$(printf '%s\n' "$werte" | grep -c .)" -ne 1 ] ||
    ! printf '%s' "$werte" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$'; then
    echo "run-sdk-kompat-tests: $datei trägt nicht genau eine Version der Form X.Y.Z, nur Ziffern ohne Suffix (gelesen: $(printf '%s' "${werte:-keine}" | tr '\n' ' '))" >&2
    exit 2
  fi
  printf '%s' "$werte"
}

# artefakt_name <sprache> <version> — der Dateiname, den make sdk-pack-<sprache> für die Version schreibt.
artefakt_name() {
  case "$1" in
    csharp) printf 'PgChangeFeed.Client.%s.nupkg' "$2" ;;
    kotlin) printf 'pgchangefeed-kotlin-%s.jar' "$2" ;;
    python) printf 'pgchangefeed-%s-py3-none-any.whl' "$2" ;;
  esac
}

# neu_verzeichnis <sprache> <version> — das Verzeichnis, das als Bau-Kontext
# `neu` in das Gast-Image kommt; im Modus dist muss es das Artefakt genau dieser
# Version tragen.
neu_verzeichnis() {
  local sprache=$1 version=$2 var dir name
  if [ "$SDK_KOMPAT_NEU" = registry ]; then
    dir="$scratch/leer-$sprache"
    mkdir -p "$dir"
    : >"$dir/leer"
    printf '%s' "$dir"
    return 0
  fi
  var="SDK_KOMPAT_DIST_$(printf '%s' "$sprache" | tr '[:lower:]' '[:upper:]')"
  dir=${!var:-sdks/$sprache/dist}
  name=$(artefakt_name "$sprache" "$version")
  if [ ! -f "$dir/$name" ]; then
    echo "run-sdk-kompat-tests: $dir trägt kein Artefakt $name (Version $version aus $(version_datei "$sprache")); vorhanden: $(ls -A "$dir" 2>/dev/null | tr '\n' ' ') — make sdk-pack-$sprache vorher (oder SDK_KOMPAT_NEU=registry)" >&2
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
  if [ "$SDK_KOMPAT_NEU" = registry ]; then
    neu_version=$REGISTRY_VERSION
  else
    neu_version=$(version_lesen "$sprache")
  fi
  neu=$(neu_verzeichnis "$sprache" "$neu_version")
  basis="pg-change-feed:sdk-kompat-base-$sprache"
  gast="pg-change-feed:sdk-kompat-$sprache"

  echo "run-sdk-kompat-tests: $sprache — Basis (Stufe build von sdks/$sprache), Gast-Bau, Modus $SDK_KOMPAT_NEU, neue Bibliothek $neu_version, veröffentlicht $REGISTRY_VERSION"
  docker build --build-context proto=proto --target build -t "$basis" "sdks/$sprache" >/dev/null
  docker build --build-arg BASE="$basis" --build-arg NEU_VERSION="$neu_version" \
    --build-arg REGISTRY_VERSION="$REGISTRY_VERSION" \
    --build-context "neu=$neu" -t "$gast" "tools/harness/sdk-kompat/$sprache" >/dev/null

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
