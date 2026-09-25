#!/usr/bin/env bash
# sdk-pack-csharp.sh — baut/testet/paketiert das C#-SDK-Package
# `PgChangeFeed.Client` (ADR-0106 Festlegung 4, slice-sdk-csharp-pack-werkzeug).
#
# Docker-only: `dotnet build`/`dotnet test`/`dotnet pack` laufen
# ausschliesslich im gepinnten mcr.microsoft.com/dotnet/sdk-Image
# (sdks/csharp/Dockerfile, AGENTS.md §3.1). `dotnet test` deckt alle vier
# Testflaechen (HTTP + gRPC + SSE + NATS-Vollinhalt, dasselbe Testprojekt
# PgChangeFeed.Client.Tests) und laeuft VOR `dotnet pack` in derselben
# Docker-Bau-Kette — ein roter Test bricht den `docker build` mit Exit != 0
# ab, bevor die `pack`-Stufe je erreicht wird (kein stiller Fallback).
#
# Die Export-Stufe `pack-export` gibt /out (das .nupkg) als tar-Stream ueber
# stdout aus; dieses Skript baut die Stufe und extrahiert host-seitig nach
# sdks/csharp/dist/ — analog tools/harness/proto-generate.sh (ADR-0060):
# kein Bind-Mount, kein `--user`-Workaround, die extrahierte Datei gehoert
# dadurch automatisch dem aufrufenden Nutzer.
#
# `set -o pipefail` macht die Pipe `docker run | tar -x` sicher: ein
# `docker run`-Fehlschlag wuerde sonst hinter einem erfolgreichen, aber
# leeren `tar -x` verschwinden (AGENTS.md §3.9 — der Exit-Code einer Pipe
# ist sonst der des letzten Glieds). Dieses Skript laeuft, wie die uebrigen
# `tools/harness/*.sh`-Gates/Werkzeuge, explizit unter bash (`make` ruft
# `bash tools/harness/sdk-pack-csharp.sh`) statt unter dem `make`-Default
# `/bin/sh` (dort `dash`, kein `pipefail`) — deshalb traegt die Pipe hier,
# ohne einen globalen `SHELL`-Override im Makefile zu brauchen.
#
# Der gRPC-Teil der Flaeche liest die `.proto` ueber einen zusaetzlichen,
# benannten Bau-Kontext (`--build-context proto=proto`, ADR-0090
# Festlegung 2, uebernommen aus slice-sdk-csharp-grpc-client-flaeche) — ohne
# ihn bricht der Bau an der `COPY --from=proto`-Zeile in
# sdks/csharp/Dockerfile ab.
#
# Aufruf: `make sdk-pack-csharp`. Override: SDK_PACK_IMAGE. Kein Gate
# (ADR-0106 Festlegung 4: NuGet-Restore braucht Netz, `make gates` bleibt
# netzlos).
set -euo pipefail

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root"

SDK_PACK_IMAGE=${SDK_PACK_IMAGE:-pg-change-feed:sdk-csharp-pack-export}

docker build --build-context proto=proto --target pack-export -t "$SDK_PACK_IMAGE" sdks/csharp

# Nach erfolgreichem Export ersetzt das Ergebnis sdks/csharp/dist/; dort
# liegen danach nur die Artefakte dieses Baus.
# shellcheck source=tools/harness/sdk-dist-clean.sh
. tools/harness/sdk-dist-clean.sh
dist="$repo_root/sdks/csharp/dist"
stage=$(sdk_dist_stage "$dist")
trap 'rm -rf -- "$stage"' EXIT
docker run --rm --network none "$SDK_PACK_IMAGE" | tar -x -C "$stage"
sdk_dist_swap "$stage" "$dist"
