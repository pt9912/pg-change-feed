#!/usr/bin/env bash
# sdk-csharp-release-tag-info.sh — validiert ein Git-Tag der Form
# sdk-csharp-v<SemVer-2.0> (ADR-0106 Festlegung 4, eigener Tag-Namensraum,
# getrennt vom Server-Release-Raum `v*`, siehe ADR-0106 §Verglichene
# Alternativen E3 verworfen) und gibt auf stdout `version=<Wert>` aus.
# Kein `latest=`-Äquivalent (anders als tools/harness/release-tag-info.sh):
# NuGet kennt kein Analogon zu einem Docker-`:latest`-Tag
# (slice-sdk-csharp-publish-workflow §1). Exit 1 bei jeder Abweichung von
# `sdk-csharp-v<SemVer-2.0>`, ohne `version=`-Ausgabe.
#
# Die SemVer-2.0-Regex selbst teilt sich dieses Skript mit
# tools/harness/release-tag-info.sh über tools/harness/semver-regex.sh
# (slice-sdk-csharp-publish-workflow §6 Risiko 3, Entscheidung: gemeinsames
# Skript statt getrennter, unabhängig driftender Kopien) — Präfix-Handling
# und Rückgabewert bleiben je Skript eigenständig, weil beide Tag-Räume
# unterschiedliche Semantik tragen (Server-Stabilität vs. kein Äquivalent).
set -euo pipefail

tag="${1:?Tag als Argument erwartet, z. B. sdk-csharp-v0.1.0}"

case "$tag" in
  sdk-csharp-v*) version="${tag#sdk-csharp-v}" ;;
  *)
    echo "sdk-csharp-release-tag-info: Tag '$tag' beginnt nicht mit 'sdk-csharp-v'" >&2
    exit 1
    ;;
esac

# shellcheck source=tools/harness/semver-regex.sh
. "$(dirname "${BASH_SOURCE[0]}")/semver-regex.sh"
if ! [[ "$version" =~ $SEMVER_RE ]]; then
  echo "sdk-csharp-release-tag-info: Tag '$tag' ist kein gültiger sdk-csharp-v<SemVer-2.0>-Tag (geprüft ohne Präfix: '$version')" >&2
  exit 1
fi

echo "version=$version"
