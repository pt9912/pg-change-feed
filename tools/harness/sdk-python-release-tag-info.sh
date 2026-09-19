#!/usr/bin/env bash
# sdk-python-release-tag-info.sh — validiert ein Git-Tag der Form
# sdk-python-v<PEP 440> (ADR-0107 Festlegung 4/5, eigener Tag-Namensraum,
# getrennt vom Server-Release-Raum `v*` und vom C#-SDK-Raum
# `sdk-csharp-v*`, siehe ADR-0107 §Entscheidung Festlegung 5 letzter
# Punkt) und gibt auf stdout `version=<Wert>` aus. Kein `latest=`-
# Äquivalent (wie bei sdk-csharp-release-tag-info.sh): PyPI kennt kein
# Analogon zu einem Docker-`:latest`-Tag (slice-sdk-python-publish-workflow
# §1). Exit 1 bei jeder Abweichung, ohne `version=`-Ausgabe.
#
# PEP 440 ist NICHT identisch mit SemVer 2.0 (u. a. andere Pre-/Post-
# Release-Syntax: PEP 440 kennt "1.0.0a1"/"1.0.0.post1"/"1.0.0.dev1" ohne
# Bindestrich, SemVer kennt "1.0.0-alpha.1"). ADR-0107 Festlegung 4 nutzt
# ausdrücklich nur den einfachen Fall MAJOR.MINOR.PATCH (z. B. "0.1.0"),
# und genau dieser einfache Fall ist zu SemVer 2.0 kompatibel (ADR-0107
# §Entscheidung Festlegung 4, erster Punkt). Dieses Skript validiert
# deshalb bewusst NUR diesen einfachen Fall — kein Epoch-, Pre-/Post-/Dev-
# Release-Segment, keine lokale Versionskennung (`+…`). Eine eigenständige,
# in dieser Datei gehaltene Regex statt eines Sourcing aus
# tools/harness/semver-regex.sh (slice-sdk-python-publish-workflow §3
# Ansatz): die dortige Regex akzeptiert SemVer-Pre-Release-/Build-Metadata-
# Syntax, die in dieser Form kein gültiges PEP 440 ist — ein Sourcing würde
# Tags akzeptieren, die als PEP-440-Version ungültig sind (auch wenn der
# nachfolgende Abgleich gegen `pyproject.toml`s `version` einen solchen Tag
# in der Praxis ohnehin verwerfen würde, wäre die Regex selbst dann falsch
# beschriftet: "PEP 440" behauptet, aber tatsächlich SemVer geprüft). Es
# gibt hier — anders als bei SemVer zwischen `release-tag-info.sh` und
# `sdk-csharp-release-tag-info.sh` — nur EINEN Konsumenten dieser Regex,
# ein eigenes `tools/harness/pep440-*-regex.sh` zum Teilen wäre unnötige
# Indirektion für eine einzige Verwendungsstelle.
set -euo pipefail

tag="${1:?Tag als Argument erwartet, z. B. sdk-python-v0.1.0}"

case "$tag" in
  sdk-python-v*) version="${tag#sdk-python-v}" ;;
  *)
    echo "sdk-python-release-tag-info: Tag '$tag' beginnt nicht mit 'sdk-python-v'" >&2
    exit 1
    ;;
esac

# PEP-440-Release-Segment, hier strikt auf genau drei Komponenten
# (MAJOR.MINOR.PATCH) begrenzt, ohne führende Nullen — deckungsgleich mit
# dem in ADR-0107 Festlegung 4 gewählten einfachen Fall, kompatibel zu
# SemVer 2.0s Release-Segment (kein Pre-/Post-/Dev-/Local-Anteil).
PEP440_SIMPLE_RE='^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$'
if ! [[ "$version" =~ $PEP440_SIMPLE_RE ]]; then
  echo "sdk-python-release-tag-info: Tag '$tag' ist kein gültiger sdk-python-v<MAJOR.MINOR.PATCH>-Tag (geprüft ohne Präfix: '$version')" >&2
  exit 1
fi

echo "version=$version"
