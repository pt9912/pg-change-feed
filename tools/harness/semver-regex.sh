# semver-regex.sh — gemeinsame POSIX-ERE SemVer-2.0-Regex (semver.org),
# geteilt zwischen tools/harness/release-tag-info.sh (`v<SemVer>`, ADR-0051
# Entscheidung 1), tools/harness/sdk-csharp-release-tag-info.sh
# (`sdk-csharp-v<SemVer>`, ADR-0106 Festlegung 4) und
# tools/harness/sdk-kotlin-release-tag-info.sh (`sdk-kotlin-v<SemVer>`,
# ADR-0109 Festlegung 4) — eine Kopie in mehreren Dateien würde diese
# Stellen unabhängig voneinander driften lassen (ursprünglich
# slice-sdk-csharp-publish-workflow §6 Risiko 3, Entscheidung: gemeinsames
# Skript statt getrennter Kopien; Pythons SDK teilt sich diese Regex
# ausdrücklich NICHT, siehe tools/harness/sdk-python-release-tag-info.sh —
# PEP 440 ist nicht identisch mit SemVer 2.0). Kein ausführbares Skript, nur eine
# Variablen-Zuweisung — per `source`/`.` einzubinden, kein eigener
# `set -euo pipefail` hier (der bleibt Sache des sourcenden Skripts).
#
# Numerische Identifier ohne führende Null (`0|[1-9][0-9]*`), alphanumerische
# Identifier dürfen führende Nullen tragen (`[0-9]*[a-zA-Z-][0-9a-zA-Z-]*`).
SEMVER_RE='^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-((0|[1-9][0-9]*|[0-9]*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9][0-9]*|[0-9]*[a-zA-Z-][0-9a-zA-Z-]*))*))?(\+([0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?$'
