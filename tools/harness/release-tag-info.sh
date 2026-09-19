#!/usr/bin/env bash
# release-tag-info.sh — validiert ein Git-Tag der Form v<SemVer-2.0>
# (ADR-0051 Entscheidung 1) und gibt auf stdout `version=<Wert>` sowie
# `latest=<true|false>` aus. `latest=true` genau dann, wenn die Version
# keinen Prerelease-Identifier trägt — eine Build-Metadata nach `+`
# zählt dafür nicht mit (SemVer 2.0 §10: Build-Metadata ändert die
# Präzedenz nicht und ist kein Prerelease). Exit 1 bei jeder Abweichung
# von SemVer 2.0, ohne `version=`/`latest=`-Ausgabe.
set -euo pipefail

tag="${1:?Tag als Argument erwartet, z. B. v1.2.3}"

case "$tag" in
  v*) version="${tag#v}" ;;
  *)
    echo "release-tag-info: Tag '$tag' beginnt nicht mit 'v'" >&2
    exit 1
    ;;
esac

# Offizielle SemVer-2.0-Regex (semver.org), auf POSIX-ERE (Bash `=~`)
# übertragen — keine `(?:...)`-Gruppen, sonst identische Struktur:
# numerische Identifier ohne führende Null (`0|[1-9][0-9]*`), alphanumerische
# Identifier dürfen führende Nullen tragen (`[0-9]*[a-zA-Z-][0-9a-zA-Z-]*`).
semver_re='^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-((0|[1-9][0-9]*|[0-9]*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9][0-9]*|[0-9]*[a-zA-Z-][0-9a-zA-Z-]*))*))?(\+([0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?$'
if ! [[ "$version" =~ $semver_re ]]; then
  echo "release-tag-info: Tag '$tag' ist kein gültiger v<SemVer-2.0>-Tag (geprüft ohne führendes v: '$version')" >&2
  exit 1
fi

# Build-Metadata (alles ab dem ersten `+`) abtrennen, bevor auf einen
# Prerelease-Bindestrich geprüft wird — sonst zählt ein Bindestrich in
# der Build-Metadata (z. B. `1.0.0+exp-sha.5114f85`) fälschlich als
# Prerelease.
core_and_prerelease="${version%%+*}"
case "$core_and_prerelease" in
  *-*) latest=false ;;
  *) latest=true ;;
esac

echo "version=$version"
echo "latest=$latest"
