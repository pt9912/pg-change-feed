#!/usr/bin/env bash
# run-notify-tests — Adapter-Tests gegen einen echten NATS-Server im
# Testcontainer (ADR-0055, `slice-052`s §6-Risiko: Docker-only-taugliche
# NATS-Testumgebung analog zu `run-store-tests.sh`). Beide Images sind per
# Digest gepinnt (Modul 14); der NATS-Pin steht bereits in `ADR-0055`
# (`docker buildx imagetools inspect nats:2-alpine`, linux/amd64-Manifest).
# Der Testcontainer und das Docker-Netz werden in jedem Ausgang abgeräumt,
# das Modul-Cache-Volume bleibt als Vorbereitung für netzlose `make
# test`-Läufe bestehen.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

TOOLCHAIN_IMAGE=${TOOLCHAIN_IMAGE:-golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125}
NATS_TEST_IMAGE=${NATS_TEST_IMAGE:-nats:2-alpine@sha256:065e8355c20a5575b3c77224be1855e8103fd148b68fba05130b9b8ddfa40ccc}
GO_MODCACHE_VOLUME=${GO_MODCACHE_VOLUME:-pg-change-feed-gomodcache}
NETWORK=cdc-notify-test
NATS_CONTAINER=cdc-notify-test-nats

docker network inspect "$NETWORK" >/dev/null 2>&1 || docker network create "$NETWORK"

cleanup() {
  docker rm -f "$NATS_CONTAINER" >/dev/null 2>&1 || true
  docker network rm "$NETWORK" >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker rm -f "$NATS_CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$NATS_CONTAINER" --network "$NETWORK" "$NATS_TEST_IMAGE" >/dev/null

# Die Bereitschaft verlangt eine echte TCP-Abfrage gegen den Client-Port
# (4222) — der Container-Start allein sagt nichts über die Annahmebereitschaft
# des NATS-Protokolls.
ready=0
for _ in $(seq 1 60); do
  if docker exec "$NATS_CONTAINER" nc -z 127.0.0.1 4222 >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done
if [ "$ready" -ne 1 ]; then
  echo "run-notify-tests: NATS-Testcontainer wurde nicht bereit" >&2
  exit 1
fi

NATS_URL="nats://$NATS_CONTAINER:4222"

# Modul-Cache befüllen (braucht Netz); der Testlauf selbst trägt die URL
# über das Docker-Netz und braucht sonst kein Netz.
docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  "$TOOLCHAIN_IMAGE" go mod download

docker run --rm --network "$NETWORK" \
  -v "$(pwd)":/src:ro \
  -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
  -w /src \
  -e GOCACHE=/tmp/gocache \
  -e CDC_NATS_TEST_URL="$NATS_URL" \
  "$TOOLCHAIN_IMAGE" go test ./internal/adapters/driven/natsnotify/...
