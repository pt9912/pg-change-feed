#!/usr/bin/env bash
# run-sdk-dist-clean-tests.sh — Tabellentest gegen sdk-dist-clean.sh:
# Altlasten in dist/ verschwinden, das neue Erzeugnis liegt danach dort, ein
# abgebrochener Export laesst dist/ unveraendert, unzulaessige Pfade werden
# abgelehnt. Netzlos, kein Docker noetig.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"
# shellcheck source=tools/harness/sdk-dist-clean.sh
. tools/harness/sdk-dist-clean.sh

fail=0
tmp=$(mktemp -d)
trap 'rm -rf -- "$tmp"' EXIT

check() {
  local what=$1
  shift
  if ! "$@"; then
    echo "FEHLER: $what" >&2
    fail=1
  fi
}

# Pfadpruefung: nur <absolut>/sdks/<name>/dist ist zulaessig.
check "gueltiger dist-Pfad" sdk_dist_check "$tmp/sdks/python/dist"
for bad in "" "/" "$tmp" "$tmp/sdks" "$tmp/sdks/python" "$tmp/dist" \
  "$tmp/sdks/dist" "$tmp/sdks/python/x/dist" "$tmp/sdks/../dist" \
  "sdks/python/dist" "sdks/python/dist/x"; do
  check "Pfad '$bad' wird abgelehnt" bash -c '! ( . tools/harness/sdk-dist-clean.sh; sdk_dist_check "$1" )' _ "$bad"
done
check "swap lehnt unzulaessigen Pfad ab" bash -c '! ( . tools/harness/sdk-dist-clean.sh; sdk_dist_swap "$1" "$2" )' _ "$tmp" "$tmp"
check "Repo-Wurzel bleibt nach abgelehntem swap bestehen" test -d "$tmp"

# Austausch: Altlasten weg, neues Erzeugnis da, kein Staging-Rest.
dist="$tmp/sdks/python/dist"
mkdir -p "$dist"
: > "$dist/pkg-0.1.0.whl"
: > "$dist/pkg-0.2.0.whl"
stage=$(sdk_dist_stage "$dist")
: > "$stage/pkg-0.2.1.whl"
sdk_dist_swap "$stage" "$dist"
check "Altlast 0.1.0 entfernt" test ! -e "$dist/pkg-0.1.0.whl"
check "Altlast 0.2.0 entfernt" test ! -e "$dist/pkg-0.2.0.whl"
check "neues Erzeugnis liegt in dist" test -f "$dist/pkg-0.2.1.whl"
check "Staging-Verzeichnis ist verschwunden" test ! -e "$stage"
check "kein Staging-Rest neben dist" test "$(ls "$tmp/sdks/python" | wc -l)" -eq 1

# Erster Lauf ohne vorhandenes dist.
dist2="$tmp/sdks/kotlin/dist"
stage=$(sdk_dist_stage "$dist2")
: > "$stage/a.jar"
sdk_dist_swap "$stage" "$dist2"
check "dist entsteht beim ersten Lauf" test -f "$dist2/a.jar"

# Abgebrochener Export: nur Staging angelegt, dist bleibt unveraendert.
stage=$(sdk_dist_stage "$dist2")
rm -rf -- "$stage"
check "dist bleibt ohne swap unveraendert" test -f "$dist2/a.jar"

if [ "$fail" -eq 0 ]; then
  echo "run-sdk-dist-clean-tests: alle Fälle bestanden"
else
  echo "run-sdk-dist-clean-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
