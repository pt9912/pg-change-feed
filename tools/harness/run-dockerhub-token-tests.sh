#!/usr/bin/env bash
# run-dockerhub-token-tests.sh — Tabellentest gegen dockerhub-token.sh
# (ADR-0051 Entscheidung 8, Review-Finding F-3 zu release-hub-description).
# Canned JSON-Antworten, kein Netzzugriff nötig. Deckt die real gegen die
# echte Docker-Hub-API beobachteten Antwortformen ab (Erfolg,
# Nutzungsfehler, Auth-Fehler, leere/fehlerhafte Eingabe).
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

fail=0

assert_ok() {
  local input="$1" want="$2" got
  if ! got=$(printf '%s' "$input" | bash tools/harness/dockerhub-token.sh 2>&1); then
    echo "FEHLER: Eingabe sollte gültig sein, scheiterte: $got" >&2
    fail=1
    return
  fi
  if [ "$got" != "$want" ]; then
    echo "FEHLER: Eingabe -> '$got', wollen '$want'" >&2
    fail=1
  fi
}

assert_invalid() {
  local input="$1" got
  if got=$(printf '%s' "$input" | bash tools/harness/dockerhub-token.sh 2>/dev/null); then
    echo "FEHLER: Eingabe sollte ungültig sein, lieferte aber: '$got'" >&2
    fail=1
  fi
}

# Erfolgsantwort (Form, real gegen die peter-evans/dockerhub-description-
# Quelle verifiziert, nicht selbst gegen die echte API auslösbar ohne
# gültige Zugangsdaten).
assert_ok '{"access_token": "abc123"}' "abc123"
assert_ok '{"access_token":"no-spaces-variant"}' "no-spaces-variant"

# Real gegen die echte API beobachtete Fehlerantworten (Review-Report
# review-slice-release-hub-description.md).
assert_invalid '{"message":"unauthorized","errinfo":{}}'
assert_invalid '{"message":"body is invalid","errinfo":{"secret":"must be at least 9 characters"}}'
assert_invalid '{"message":"body is invalid","errinfo":{"identifier":"value is required","secret":"value is required"}}'

# Degenerierte Eingaben (Netzausfall / unerwartete Antwortform).
assert_invalid ''
assert_invalid 'not json at all'
assert_invalid '{"access_token": null}'
assert_invalid '{"access_token": ""}'

if [ "$fail" -eq 0 ]; then
  echo "run-dockerhub-token-tests: alle Fälle bestanden"
else
  echo "run-dockerhub-token-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
