#!/usr/bin/env bash
# run-sdk-kotlin-release-tag-info-tests.sh — Tabellentest gegen
# sdk-kotlin-release-tag-info.sh (ADR-0109 Festlegung 4/5). Deckt denselben
# Fällekatalog wie run-sdk-csharp-release-tag-info-tests.sh (geteilte Regex
# über tools/harness/semver-regex.sh), gegen den `sdk-kotlin-v`-Präfix statt
# `sdk-csharp-v` und ohne `latest=`-Rückgabe (kein GitHub-Packages-Äquivalent
# zu einem Docker-`:latest`-Tag). Netzlos, kein Docker nötig (reine
# Bash-Logik).
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

fail=0

assert_ok() {
  local tag="$1" want_version="$2" out got_version
  if ! out=$(bash tools/harness/sdk-kotlin-release-tag-info.sh "$tag" 2>&1); then
    echo "FEHLER: '$tag' sollte gültig sein, scheiterte: $out" >&2
    fail=1
    return
  fi
  got_version=$(printf '%s\n' "$out" | grep '^version=' | cut -d= -f2)
  if [ "$got_version" != "$want_version" ]; then
    echo "FEHLER: '$tag' -> version=$got_version, wollen version=$want_version" >&2
    fail=1
  fi
}

assert_invalid() {
  local tag="$1" out
  if out=$(bash tools/harness/sdk-kotlin-release-tag-info.sh "$tag" 2>&1); then
    echo "FEHLER: '$tag' sollte ungültig sein, wurde aber akzeptiert: $out" >&2
    fail=1
  fi
}

# Gültig.
assert_ok "sdk-kotlin-v0.1.0" "0.1.0"
assert_ok "sdk-kotlin-v1.0.0" "1.0.0"
assert_ok "sdk-kotlin-v1.0.0-alpha.1" "1.0.0-alpha.1"
assert_ok "sdk-kotlin-v1.0.0+build.5" "1.0.0+build.5"
assert_ok "sdk-kotlin-v1.0.0+exp-sha.5114f85" "1.0.0+exp-sha.5114f85"
assert_ok "sdk-kotlin-v2.0.0-rc.1+build.5" "2.0.0-rc.1+build.5"

# Ungültig: falscher/fehlender Präfix.
assert_invalid "0.1.0"
assert_invalid "v0.1.0"
assert_invalid "kotlin-v0.1.0"
assert_invalid "sdk-csharp-v0.1.0"
assert_invalid "sdk-python-v0.1.0"
# Ungültig: kein SemVer 2.0 (dieselben Fälle wie release-tag-info/
# sdk-csharp-release-tag-info, geteilte Regex).
assert_invalid "sdk-kotlin-v1.0"
assert_invalid "sdk-kotlin-v01.0.0"
assert_invalid "sdk-kotlin-v1.0.0-"
assert_invalid "sdk-kotlin-v1.0.0-01"

if [ "$fail" -eq 0 ]; then
  echo "run-sdk-kotlin-release-tag-info-tests: alle Fälle bestanden"
else
  echo "run-sdk-kotlin-release-tag-info-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
