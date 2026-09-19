#!/usr/bin/env bash
# run-sdk-python-release-tag-info-tests.sh — Tabellentest gegen
# sdk-python-release-tag-info.sh (ADR-0107 Festlegung 4/5). Anders als
# run-sdk-csharp-release-tag-info-tests.sh deckt dieser Test NUR den
# einfachen PEP-440-Fall MAJOR.MINOR.PATCH (keine geteilte Regex mit
# tools/harness/semver-regex.sh, siehe Kopf-Kommentar von
# sdk-python-release-tag-info.sh). Netzlos, kein Docker nötig (reine
# Bash-Logik).
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

fail=0

assert_ok() {
  local tag="$1" want_version="$2" out got_version
  if ! out=$(bash tools/harness/sdk-python-release-tag-info.sh "$tag" 2>&1); then
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
  if out=$(bash tools/harness/sdk-python-release-tag-info.sh "$tag" 2>&1); then
    echo "FEHLER: '$tag' sollte ungültig sein, wurde aber akzeptiert: $out" >&2
    fail=1
  fi
}

# Gültig: einfacher PEP-440-Fall MAJOR.MINOR.PATCH.
assert_ok "sdk-python-v0.1.0" "0.1.0"
assert_ok "sdk-python-v1.0.0" "1.0.0"
assert_ok "sdk-python-v10.20.30" "10.20.30"

# Ungültig: falscher/fehlender Präfix.
assert_invalid "0.1.0"
assert_invalid "v0.1.0"
assert_invalid "python-v0.1.0"
assert_invalid "sdk-csharp-v0.1.0"

# Ungültig: kein einfacher PEP-440-Fall MAJOR.MINOR.PATCH.
assert_invalid "sdk-python-v1.0"
assert_invalid "sdk-python-v1.0.0.0"
assert_invalid "sdk-python-v01.0.0"
assert_invalid "sdk-python-v1.0.0a1"
assert_invalid "sdk-python-v1.0.0-alpha.1"
assert_invalid "sdk-python-v1.0.0.post1"
assert_invalid "sdk-python-v1.0.0.dev1"
assert_invalid "sdk-python-v1.0.0+build.5"
assert_invalid "sdk-python-v1!1.0.0"

if [ "$fail" -eq 0 ]; then
  echo "run-sdk-python-release-tag-info-tests: alle Fälle bestanden"
else
  echo "run-sdk-python-release-tag-info-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
