#!/usr/bin/env bash
# run-dockerfile-from-tests.sh — Tabellentest gegen dockerfile-from.sh
# (FROM-Parsing des Base-Image-Drift-Sensors image-stale.sh): eine
# Fixture mit allen FROM-Formen und das reale Dockerfile. Netzlos, kein
# Docker nötig (reine Bash-/awk-Logik).
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

fail=0
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

D1=sha256:1111111111111111111111111111111111111111111111111111111111111111
D2=sha256:2222222222222222222222222222222222222222222222222222222222222222

cat > "$tmp/Dockerfile" <<EOF
# FROM kommentar@sha256:zzz
ARG X=1
FROM golang:1.27-alpine@$D1
FROM --platform=\$BUILDPLATFORM golang:1.27-alpine@$D1 AS deps
FROM --platform=linux/amd64 --other=1 gcr.io/distroless/static-debian12:nonroot@$D2 AS runtime
FROM deps AS proto
FROM scratch
from lowercase:1.0@$D2 as low
RUN echo FROM not-a-from@sha256:zzz
EOF

want=$(printf '%s %s\n' \
  golang:1.27-alpine "$D1" \
  golang:1.27-alpine "$D1" \
  gcr.io/distroless/static-debian12:nonroot "$D2" \
  deps "" \
  scratch "" \
  lowercase:1.0 "$D2")

got=$(bash tools/harness/dockerfile-from.sh "$tmp/Dockerfile")
if [ "$got" != "$want" ]; then
  echo "FEHLER: Fixture-Ausgabe weicht ab." >&2
  echo "--- erwartet" >&2; printf '%s\n' "$want" >&2
  echo "--- erhalten" >&2; printf '%s\n' "$got" >&2
  fail=1
fi

# Reales Dockerfile: je FROM-Zeile genau eine Ausgabezeile, keine Option
# und keine Variable als Referenz, jede Zeile mit `@` trägt einen Digest.
real_from=$(grep -c -E '^FROM ' Dockerfile)
real_out=$(bash tools/harness/dockerfile-from.sh Dockerfile)
real_lines=$(printf '%s\n' "$real_out" | grep -c .)
if [ "$real_from" != "$real_lines" ]; then
  echo "FEHLER: Dockerfile hat $real_from FROM-Zeilen, dockerfile-from.sh gab $real_lines Zeilen." >&2
  fail=1
fi
if printf '%s\n' "$real_out" | grep -q -E '^(--|\$)'; then
  echo "FEHLER: eine Referenz des Dockerfile beginnt mit Option oder Variable." >&2
  fail=1
fi
pinned_from=$(grep -c -E '^FROM .*@sha256:' Dockerfile)
pinned_out=$(printf '%s\n' "$real_out" | grep -c -E ' sha256:[0-9a-f]{64}$')
if [ "$pinned_from" != "$pinned_out" ]; then
  echo "FEHLER: $pinned_from gepinnte FROM-Zeilen im Dockerfile, $pinned_out mit Digest gelesen." >&2
  fail=1
fi

if [ "$fail" -ne 0 ]; then
  exit 1
fi
echo "OK: dockerfile-from.sh (Fixture + reales Dockerfile)"
