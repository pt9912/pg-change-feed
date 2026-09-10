#!/usr/bin/env bash
# commit-traceability.sh — Grenz-Haelfte des Commit-Traceability-Gates
# (ADR-0045): keine Struktur-ID (Nummern-Kennung vom SPEC- oder ARC-Typ)
# im Betreff je Commit der Range. Die positive Haelfte (je Message eine
# Vertrags-Kennung) traegt das d-check-Modul commits im selben Target
# (d-check.mk) — ein zweiter Sensor ueber dieselbe Haelfte waere eine
# zweite Quelle fuer dieselbe Regel. Aufgerufen vom make-Target
# commit-traceability; die Range kommt als $1, Default HEAD~5..HEAD.
set -euo pipefail

RANGE="${1:-HEAD~5..HEAD}"

# Fail-closed: beide Range-Enden muessen aufloesbar sein, sonst liefe der
# Sensor still durch ein leeres git log.
base="${RANGE%%..*}"
head_ref="${RANGE##*..}"
if ! git rev-parse -q --verify "${base}^{commit}" >/dev/null 2>&1; then
  echo "commit-traceability: Range-Basis \"${base}\" nicht aufloesbar (fail-closed)" >&2
  exit 2
fi
if ! git rev-parse -q --verify "${head_ref}^{commit}" >/dev/null 2>&1; then
  echo "commit-traceability: Range-Spitze \"${head_ref}\" nicht aufloesbar (fail-closed)" >&2
  exit 2
fi

if ! commits="$(git log --format='%h %s' "${RANGE}")"; then
  echo "commit-traceability: git log über \"${RANGE}\" fehlgeschlagen (fail-closed)" >&2
  exit 2
fi
count="$(git rev-list --count "${RANGE}")"

status=0
while IFS= read -r entry; do
  [ -z "${entry}" ] && continue
  sha="${entry%% *}"
  subject="${entry#* }"
  if [[ "${subject}" =~ (SPEC|ARC)-[0-9]{3} ]]; then
    echo "commit-traceability: ${sha}: Struktur-ID im Betreff: ${subject}" >&2
    status=1
  fi
done <<< "${commits}"

if [ "${status}" -eq 0 ]; then
  echo "commit-traceability: OK — ${count} Commit(s) in \"${RANGE}\", Betreffs ohne Struktur-ID"
fi
exit "${status}"