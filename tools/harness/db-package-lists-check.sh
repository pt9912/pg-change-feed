#!/usr/bin/env bash
# db-package-lists-check.sh — haelt die namentlichen Paketlisten des
# DB-Adapter-Gegenstands gleich (ADR-0071 Punkt 1). Die tragende Regel ist die
# Eigenschaft "Testlauf setzt einen externen Dienst voraus"; die Listen sind ihre
# namentlichen Traeger:
#   subject    DB_COVERAGE_PKGS in tools/harness/db-coverage.sh (--coverpkg)
#   filter     das Ausschluss-Muster der Dockerfile-Stufe `coverage`
#   store      die Paketliste des Messlaufs in tools/harness/run-store-tests.sh
#   replic.    die Paketliste des Messlaufs in tools/harness/run-replication-tests.sh
# Gleich heisst: der Filter nimmt genau die Pakete des Gegenstands aus, und die
# Messlaeufe testen zusammen genau die Pakete des Gegenstands. Ein Unterschied
# endet mit Exit 1 und nennt die abweichenden Pakete; Vertrag:
# harness/sensors/db-adapter-coverage.md §Gegenstand. Netzlos, read-only.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

subject=$(bash tools/harness/db-coverage.sh --coverpkg | tr ',' '\n' | LC_ALL=C sort -u)

# Dockerfile: `grep -vE '(^|/)(a|b|c)$'` — genau eine solche Zeile.
filter_lines=$(sed -nE "s/.*grep -vE '\(\^\|\/\)\(([^)]*)\)\\$'.*/\1/p" Dockerfile)
if [[ $(printf '%s\n' "$filter_lines" | grep -c .) -ne 1 ]]; then
  echo "db-package-lists-check: Dockerfile traegt nicht genau ein Ausschluss-Muster '(^|/)(...)\$'" >&2
  exit 1
fi
filter=""
while IFS= read -r alt; do
  matches=$(printf '%s\n' "$subject" | grep -E "/${alt}\$" || true)
  if [[ $(printf '%s\n' "$matches" | grep -c .) -ne 1 ]]; then
    echo "db-package-lists-check: Filter-Eintrag '$alt' trifft nicht genau ein Paket von DB_COVERAGE_PKGS" >&2
    exit 1
  fi
  filter+="$matches"$'\n'
done < <(printf '%s\n' "$filter_lines" | tr '|' '\n')
filter=$(printf '%s' "$filter" | LC_ALL=C sort -u)

# `go test`-Paketliste hinter der -coverprofile-Zeile einer Messphase.
test_list() {
  awk -v marker="$2" '
    index($0, marker) { on = 1; next }
    on {
      n = split($0, w, /[ \t\\]+/)
      for (i = 1; i <= n; i++) if (w[i] ~ /^\.\/internal\//) print w[i]
      if ($0 !~ /\\[ \t]*$/) on = 0
    }
  ' "$1"
}
store=$(test_list tools/harness/run-store-tests.sh '-coverprofile=/cov/store.coverprofile')
replication=$(test_list tools/harness/run-replication-tests.sh '-coverprofile=/cov/replication.coverprofile')
tested=$(printf '%s\n%s\n' "$store" "$replication" | grep . | LC_ALL=C sort -u)

status=0
compare() {
  local name=$1 got=$2
  if [[ "$got" != "$subject" ]]; then
    echo "db-package-lists-check: $name weicht von DB_COVERAGE_PKGS ab (< nur im Gegenstand, > nur in $name):" >&2
    diff <(printf '%s\n' "$subject") <(printf '%s\n' "$got") >&2 || true
    status=1
  fi
}
compare "Dockerfile-Filter" "$filter"
compare "Messlaeufe (store + replication)" "$tested"

if [[ $status -eq 0 ]]; then
  echo "db-package-lists-check: OK — Dockerfile-Filter, Messlaeufe und DB_COVERAGE_PKGS nennen dieselben $(printf '%s\n' "$subject" | grep -c .) Pakete"
fi
exit $status
