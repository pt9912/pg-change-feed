# semver-regex.sh — gemeinsame POSIX-ERE SemVer-2.0-Regex (semver.org),
# geteilt zwischen tools/harness/release-tag-info.sh (`v<SemVer>`, ADR-0051
# Entscheidung 1) und tools/harness/sdk-csharp-release-tag-info.sh
# (`sdk-csharp-v<SemVer>`, ADR-0106 Festlegung 4) — eine Kopie in zwei
# Dateien würde beide Stellen unabhängig voneinander driften lassen
# (slice-sdk-csharp-publish-workflow §6 Risiko 3, Entscheidung: gemeinsames
# Skript statt getrennter Kopien). Kein ausführbares Skript, nur eine
# Variablen-Zuweisung — per `source`/`.` einzubinden, kein eigener
# `set -euo pipefail` hier (der bleibt Sache des sourcenden Skripts).
#
# Numerische Identifier ohne führende Null (`0|[1-9][0-9]*`), alphanumerische
# Identifier dürfen führende Nullen tragen (`[0-9]*[a-zA-Z-][0-9a-zA-Z-]*`).
SEMVER_RE='^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-((0|[1-9][0-9]*|[0-9]*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9][0-9]*|[0-9]*[a-zA-Z-][0-9a-zA-Z-]*))*))?(\+([0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?$'
