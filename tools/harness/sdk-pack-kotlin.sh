#!/usr/bin/env bash
# sdk-pack-kotlin.sh — baut/testet/paketiert das Kotlin-SDK-Package
# `pgchangefeed-kotlin` (ADR-0109 Festlegung 5, slice-sdk-kotlin-pack-werkzeug).
#
# Docker-only: `./gradlew test`/`./gradlew build` laufen ausschliesslich im
# gepinnten `eclipse-temurin:21-jdk`-Image (sdks/kotlin/Dockerfile,
# AGENTS.md §3.1). `./gradlew test` laeuft VOR dem Jar-Bau in derselben
# Docker-Bau-Kette — ein roter Test bricht den `docker build` mit
# Exit != 0 ab, bevor die `pack`-Stufe je erreicht wird (kein stiller
# Fallback).
#
# Die Export-Stufe `pack-export` gibt /out (das .jar) als tar-Stream ueber
# stdout aus; dieses Skript baut die Stufe und extrahiert host-seitig nach
# sdks/kotlin/dist/ — analog tools/harness/sdk-pack-csharp.sh/
# sdk-pack-python.sh: kein Bind-Mount, kein `--user`-Workaround, die
# extrahierte Datei gehoert dadurch automatisch dem aufrufenden Nutzer.
#
# `set -o pipefail` macht die Pipe `docker run | tar -x` sicher: ein
# `docker run`-Fehlschlag wuerde sonst hinter einem erfolgreichen, aber
# leeren `tar -x` verschwinden (AGENTS.md §3.9 — der Exit-Code einer Pipe
# ist sonst der des letzten Glieds). Dieses Skript laeuft, wie die
# uebrigen `tools/harness/*.sh`-Gates/Werkzeuge, explizit unter bash
# (`make` ruft `bash tools/harness/sdk-pack-kotlin.sh`) statt unter dem
# `make`-Default `/bin/sh` (dort `dash`, kein `pipefail`) — deshalb traegt
# die Pipe hier, ohne einen globalen `SHELL`-Override im Makefile zu
# brauchen.
#
# Der gRPC-Teil der Flaeche liest die `.proto` ueber einen zusaetzlichen,
# benannten Bau-Kontext (`--build-context proto=proto`, ADR-0090
# Festlegung 2, uebernommen auf den Kotlin-SDK-Baum) — ohne ihn bricht der
# Bau an der `COPY --from=proto`-Zeile in sdks/kotlin/Dockerfile ab.
#
# Aufruf: `make sdk-pack-kotlin`. Override: SDK_PACK_KOTLIN_IMAGE. Kein
# Gate (ADR-0109 Festlegung 5: GitHub-Packages-Paketbezug braucht Netz,
# `make gates` bleibt netzlos).
set -euo pipefail

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root"

SDK_PACK_KOTLIN_IMAGE=${SDK_PACK_KOTLIN_IMAGE:-pg-change-feed:sdk-kotlin-pack-export}

docker build --build-context proto=proto --target pack-export -t "$SDK_PACK_KOTLIN_IMAGE" sdks/kotlin
mkdir -p sdks/kotlin/dist
docker run --rm --network none "$SDK_PACK_KOTLIN_IMAGE" | tar -x -C sdks/kotlin/dist
