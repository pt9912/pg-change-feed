#!/usr/bin/env bash
# run-release-tag-info-tests.sh — Tabellentest gegen release-tag-info.sh
# (ADR-0051 Entscheidung 1/2). Deckt die beiden real im Review zu
# release-version-und-workflow gefundenen SemVer-2.0-Abweichungen ab:
# eine führende Null in einem rein numerischen Prerelease-Identifier
# (`1.0.0-01`, von SemVer 2.0 §9 verboten) und ein Bindestrich in der
# Build-Metadata (`1.0.0+exp-sha.5114f85`, kein Prerelease trotz
# Bindestrich). Netzlos, kein Docker nötig (reine Bash-Logik, kein
# externes Werkzeug).
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

fail=0

assert_ok() {
  local tag="$1" want_version="$2" want_latest="$3" out got_version got_latest
  if ! out=$(bash tools/harness/release-tag-info.sh "$tag" 2>&1); then
    echo "FEHLER: '$tag' sollte gültig sein, scheiterte: $out" >&2
    fail=1
    return
  fi
  got_version=$(printf '%s\n' "$out" | grep '^version=' | cut -d= -f2)
  got_latest=$(printf '%s\n' "$out" | grep '^latest=' | cut -d= -f2)
  if [ "$got_version" != "$want_version" ] || [ "$got_latest" != "$want_latest" ]; then
    echo "FEHLER: '$tag' -> version=$got_version latest=$got_latest, wollen version=$want_version latest=$want_latest" >&2
    fail=1
  fi
}

assert_invalid() {
  local tag="$1" out
  if out=$(bash tools/harness/release-tag-info.sh "$tag" 2>&1); then
    echo "FEHLER: '$tag' sollte ungültig sein, wurde aber akzeptiert: $out" >&2
    fail=1
  fi
}

# Gültig, stabil.
assert_ok "v1.0.0" "1.0.0" "true"
assert_ok "v0.1.0" "0.1.0" "true"
# Gültig, Prerelease -> latest=false.
assert_ok "v1.0.0-alpha.1" "1.0.0-alpha.1" "false"
# Gültig, Build-Metadata ohne Prerelease -> latest=true trotz Bindestrich
# in der Build-Metadata (Review-Finding F-3).
assert_ok "v1.0.0+build.5" "1.0.0+build.5" "true"
assert_ok "v1.0.0+exp-sha.5114f85" "1.0.0+exp-sha.5114f85" "true"
# Gültig, Prerelease + Build-Metadata.
assert_ok "v2.0.0-rc.1+build.5" "2.0.0-rc.1+build.5" "false"
# Gültig, numerischer Prerelease-Identifier "0" (kein Leading-Zero-Fall).
assert_ok "v1.0.0-0" "1.0.0-0" "false"

# Ungültig.
assert_invalid "1.0.0"
assert_invalid "v1.0"
assert_invalid "v01.0.0"
assert_invalid "v1.0.0-"
# Review-Finding F-2: führende Null in einem rein numerischen
# Prerelease-Identifier ist laut SemVer 2.0 §9 verboten.
assert_invalid "v1.0.0-01"
assert_invalid "v1.0.0-01.2"

if [ "$fail" -eq 0 ]; then
  echo "run-release-tag-info-tests: alle Fälle bestanden"
else
  echo "run-release-tag-info-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
