#!/usr/bin/env bash
# run-version-injection-test.sh — Regressionstest gegen die
# `VERSION`-Injektion vom `Dockerfile` bis ins Binary (`ADR-0051`,
# Review-Finding F-1 zu `release-binary-version-injektion`): baut real
# über `docker buildx build --build-arg VERSION=...` (derselbe Pfad wie
# `make image VERSION=...`) und prueft `--version` gegen den erwarteten
# Marker. Eine versehentliche Entfernung des `-X`-Flags oder des
# `ARG VERSION` in der `build`-Stufe waere sonst lautlos — der Build
# kompiliert weiterhin fehlerfrei, `main.version` bliebe nur unveraendert
# auf ihrem Default. Netzlos, sofern die Basis-Images bereits lokal
# gecacht sind (z. B. nach einem vorherigen `make image`/`make gates`).
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

MARKER="version-injection-test-marker-42"
TEST_TAG="pg-change-feed-version-injection-test:tmp"

if ! docker buildx build --load --build-arg VERSION="$MARKER" -t "$TEST_TAG" . >/tmp/version-injection-test-build.log 2>&1; then
  echo "run-version-injection-test: Build scheiterte:" >&2
  cat /tmp/version-injection-test-build.log >&2
  exit 1
fi

out=$(docker run --rm "$TEST_TAG" --version 2>&1)
rc=$?
docker rmi -f "$TEST_TAG" >/dev/null 2>&1

if [ "$rc" -ne 0 ]; then
  echo "run-version-injection-test: Container-Lauf scheiterte: $out" >&2
  exit 1
fi

if [ "$out" != "pg-change-feed $MARKER" ]; then
  echo "run-version-injection-test: VERSION-Injektion wirkt nicht — erwartet 'pg-change-feed $MARKER', erhalten '$out'" >&2
  exit 1
fi

echo "run-version-injection-test: VERSION-Injektion wirkt korrekt ($out)"
