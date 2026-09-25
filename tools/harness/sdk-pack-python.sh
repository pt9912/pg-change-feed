#!/usr/bin/env bash
# sdk-pack-python.sh — baut/testet/paketiert das Python-SDK-Package
# `pgchangefeed` (ADR-0107 Festlegung 5, ADR-0108 §Entscheidung
# Festlegung 1/3, slice-sdk-python-pack-werkzeug).
#
# Docker-only: `pip install`/`pytest`/`uv build --no-sources` laufen
# ausschliesslich im gepinnten `python:3.14-slim`-Image
# (sdks/python/Dockerfile, AGENTS.md §3.1). `pytest` laeuft VOR
# `uv build --no-sources` in derselben Docker-Bau-Kette — ein roter Test
# bricht den `docker build` mit Exit != 0 ab, bevor die `pack`-Stufe je
# erreicht wird (kein stiller Fallback).
#
# Die Export-Stufe `pack-export` gibt /out (.whl UND .tar.gz) als
# tar-Stream ueber stdout aus; dieses Skript baut die Stufe und
# extrahiert host-seitig nach sdks/python/dist/ — analog
# tools/harness/sdk-pack-csharp.sh (ADR-0106): kein Bind-Mount, kein
# `--user`-Workaround, die extrahierten Dateien gehoeren dadurch
# automatisch dem aufrufenden Nutzer.
#
# `set -o pipefail` macht die Pipe `docker run | tar -x` sicher: ein
# `docker run`-Fehlschlag wuerde sonst hinter einem erfolgreichen, aber
# leeren `tar -x` verschwinden (AGENTS.md §3.9 — der Exit-Code einer Pipe
# ist sonst der des letzten Glieds). Dieses Skript laeuft, wie die
# uebrigen `tools/harness/*.sh`-Gates/Werkzeuge, explizit unter bash
# (`make` ruft `bash tools/harness/sdk-pack-python.sh`) statt unter dem
# `make`-Default `/bin/sh` (dort `dash`, kein `pipefail`) — deshalb traegt
# die Pipe hier, ohne einen globalen `SHELL`-Override im Makefile zu
# brauchen.
#
# Der gRPC-Teil der Flaeche liest die `.proto` ueber einen zusaetzlichen,
# benannten Bau-Kontext (`--build-context proto=proto`, ADR-0090
# Festlegung 2, Muster sdks/csharp/Dockerfile) — ohne ihn bricht der Bau
# an der `COPY --from=proto`-Zeile in sdks/python/Dockerfile ab.
#
# Aufruf: `make sdk-pack-python`. Override: SDK_PACK_PYTHON_IMAGE. Kein
# Gate (ADR-0107 Festlegung 5: PyPI-Paketbezug fuer Test-Abhaengigkeiten
# braucht Netz, `make gates` bleibt netzlos).
set -euo pipefail

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root"

SDK_PACK_PYTHON_IMAGE=${SDK_PACK_PYTHON_IMAGE:-pg-change-feed:sdk-python-pack-export}

docker build --build-context proto=proto --target pack-export -t "$SDK_PACK_PYTHON_IMAGE" sdks/python

# Nach erfolgreichem Export ersetzt das Ergebnis sdks/python/dist/; dort
# liegen danach nur die Artefakte dieses Baus.
# shellcheck source=tools/harness/sdk-dist-clean.sh
. tools/harness/sdk-dist-clean.sh
dist="$repo_root/sdks/python/dist"
stage=$(sdk_dist_stage "$dist")
trap 'rm -rf -- "$stage"' EXIT
docker run --rm --network none "$SDK_PACK_PYTHON_IMAGE" | tar -x -C "$stage"
sdk_dist_swap "$stage" "$dist"
