#!/usr/bin/env bash
# proto-generate.sh — Protobuf-/gRPC-Go-Code aus proto/cdc/stream/v1/changestream.proto
# erzeugen (Docker-only, Build-Zeit-Erzeugung, Host-Extraktion; ADR-0060).
#
# Der Generator ist die Dockerfile-Stufe `proto-export` (seit slice-104): sie
# kopiert die `.proto`-Quelle per COPY hinein (kein Bind-Mount) und erzeugt
# den Code zur Build-Zeit; ihr ENTRYPOINT gibt das Erzeugnis als `tar`-Stream
# ueber stdout aus. Dieses Skript baut die Stufe und extrahiert host-seitig,
# damit die erzeugten Dateien dem aufrufenden Nutzer gehoeren (kein
# `--user`-Workaround noetig). Kein Gate: der Generator laeuft nur, wenn sich
# die `.proto`-Quelle aendert; der erzeugte Go-Code liegt committet im Baum
# und wird von `make test`/`make image` mitkompiliert.
#
# `set -o pipefail` macht die Pipe `docker run | tar -x` sicher: ein
# `docker run`-Fehlschlag wuerde sonst hinter einem erfolgreichen, aber
# leeren `tar -x` verschwinden (`AGENTS.md` §3.9 — der Exit-Code einer Pipe
# ist sonst der des letzten Glieds). Dieses Skript laeuft, wie die uebrigen
# `tools/harness/*.sh`-Gates, explizit unter bash (`make` ruft `bash
# tools/harness/proto-generate.sh`) statt unter dem `make`-Default `/bin/sh`
# (dort `dash`, kein `pipefail`) — deshalb traegt die Pipe hier, ohne einen
# globalen `SHELL`-Override im Makefile zu brauchen.
#
# Aufruf: `make proto-generate`. Override: PROTO_IMAGE.
set -euo pipefail

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root"

PROTO_IMAGE=${PROTO_IMAGE:-pg-change-feed:proto-export}

docker build --target proto-export -t "$PROTO_IMAGE" .
docker run --rm --network none "$PROTO_IMAGE" | tar -x -C .
