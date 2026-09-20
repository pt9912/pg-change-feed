#!/usr/bin/env bash
# sdk-kotlin-release-tag-info.sh — validiert ein Git-Tag der Form
# sdk-kotlin-v<SemVer-2.0> (ADR-0109 Festlegung 4/5, eigener Tag-Namensraum,
# getrennt vom Server-Release-Raum `v*` und den beiden anderen SDK-Räumen
# `sdk-csharp-v*`/`sdk-python-v*`) und gibt auf stdout `version=<Wert>` aus.
# Kein `latest=`-Äquivalent (wie bei sdk-csharp-release-tag-info.sh): GitHub
# Packages kennt kein Analogon zu einem Docker-`:latest`-Tag
# (slice-sdk-kotlin-publish-workflow §1). Exit 1 bei jeder Abweichung von
# `sdk-kotlin-v<SemVer-2.0>`, ohne `version=`-Ausgabe.
#
# Die SemVer-2.0-Regex selbst teilt sich dieses Skript mit
# tools/harness/release-tag-info.sh und
# tools/harness/sdk-csharp-release-tag-info.sh über
# tools/harness/semver-regex.sh (ADR-0109 Festlegung 4: das Kotlin-SDK führt
# reines SemVer 2.0, nicht Pythons PEP-440-Sonderfall — dieselbe Grammatik
# wie beim C#-SDK) — Präfix-Handling und Rückgabewert bleiben je Skript
# eigenständig.
set -euo pipefail

tag="${1:?Tag als Argument erwartet, z. B. sdk-kotlin-v0.1.0}"

case "$tag" in
  sdk-kotlin-v*) version="${tag#sdk-kotlin-v}" ;;
  *)
    echo "sdk-kotlin-release-tag-info: Tag '$tag' beginnt nicht mit 'sdk-kotlin-v'" >&2
    exit 1
    ;;
esac

# shellcheck source=tools/harness/semver-regex.sh
. "$(dirname "${BASH_SOURCE[0]}")/semver-regex.sh"
if ! [[ "$version" =~ $SEMVER_RE ]]; then
  echo "sdk-kotlin-release-tag-info: Tag '$tag' ist kein gültiger sdk-kotlin-v<SemVer-2.0>-Tag (geprüft ohne Präfix: '$version')" >&2
  exit 1
fi

echo "version=$version"
