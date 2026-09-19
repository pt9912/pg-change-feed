#!/usr/bin/env bash
# tools/harness/dockerhub-token.sh — extrahiert den access_token aus der
# JSON-Antwort von POST https://hub.docker.com/v2/auth/token (stdin).
# Reiner grep/sed-Parser ohne eigenen Netzzugriff und ohne jq-Abhaengigkeit
# — netzlos und ohne Docker direkt testbar
# (tools/harness/run-dockerhub-token-tests.sh, make test-dockerhub-token),
# analog dem release-tag-info.sh-Muster derselben Welle. Von
# .github/workflows/hub-description.yml verwendet (Review-Finding F-3 zu
# release-hub-description: die einzige nicht-triviale Logik dieses
# Slices bekommt damit dieselbe committete, netzlose Testbarkeit).
#
# Exit: 0 = Token auf stdout · 1 = kein access_token in der Antwort
# (Login-Fehler, unerwartete Antwortform oder leere/fehlerhafte Eingabe).
set -uo pipefail

input="$(cat)"
token=$(printf '%s' "$input" | grep -o '"access_token"[[:space:]]*:[[:space:]]*"[^"]*"' | sed -E 's/.*"access_token"[[:space:]]*:[[:space:]]*"([^"]*)".*/\1/')
if [ -z "$token" ]; then
  echo "dockerhub-token: kein access_token in der Antwort ($input)" >&2
  exit 1
fi
echo "$token"
