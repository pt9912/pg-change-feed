# sdks/python/Dockerfile — Docker-only Bau-/Testsetup des Python-SDK-Pakets
# `pgchangefeed`. Bau-Kontext ist sdks/python/ selbst — NICHT die Repo-Wurzel
# und NICHT examples/ (es gibt kein examples/python/).
#
# Die gRPC-Client-Fläche liest die `.proto` NICHT aus diesem Kontext — sie
# liegt außerhalb davon, unter proto/cdc/stream/v1/ an der Repo-Wurzel. Sie
# kommt über einen ZUSÄTZLICHEN, benannten Bau-Kontext namens `proto` herein
# (`COPY --from=proto …` unten); der Aufrufer bringt ihn mit `docker build
# --build-context proto=proto …` mit (Muster sdks/csharp/Dockerfile). Ohne ihn
# bricht der Bau an der COPY-Zeile ab — kein stiller Fallback. Die Stub-Module
# entstehen im Bau (`grpcio-tools`/`protoc`, Stufe `build`) und werden nicht
# committet (.gitignore); sie landen als Unterpackage `pgchangefeed.grpc_gen`
# im Wheel — das erzeugte `_pb2_grpc`-Modul trägt den generierten flachen
# Geschwister-Import, den der `sed`-Schritt unten auf den Package-Import
# hebt, bevor pytest ihn importiert.
#
# Keine Runtime-Stufe hier: ein SDK-Package ist keine startbare Anwendung
# (analog sdks/csharp/Dockerfile). Die `pack`/`pack-export`-Stufen unten bauen
# auf `build` auf: `pytest` (unten) läuft VOR `uv build --no-sources` (Stufe
# `pack`) — ein roter Test bricht den `docker build` mit Exit != 0 ab, bevor
# `pack` je erreicht wird (kein stiller Fallback, Muster
# sdks/csharp/Dockerfile). Die `integration`-Stufe baut ebenfalls auf `build`
# auf und startet pytest NICHT selbst — ihr CMD ist der
# Realserver-Integrationstest der SDK-Fläche
# (tools/harness/run-sdk-python-integration-tests.sh); sie braucht Netz und
# DB-Zugang und ist deshalb kein Gate.
#
# Basis-Digest gemessen (`docker buildx imagetools inspect`,
# Manifest-Index-Digest über alle Plattformen):
#   python:3.14-slim (3.14.7-slim-trixie, aktuelle Python-Version)
#     -> sha256:caaf356f40667c496d405780745b9ac25771c189a51dfcc42430d531ea09f8a2
# Pin-Hebung ist ein bewusster Commit, keine stille Auflösung.
FROM python:3.14-slim@sha256:caaf356f40667c496d405780745b9ac25771c189a51dfcc42430d531ea09f8a2 AS build
WORKDIR /src

# pyproject.toml zuerst — es ändert sich seltener als der Quellcode und
# hält den Docker-Layer-Cache über Quelltext-Änderungen hinweg stabil
# (analog sdks/csharp/Dockerfile).
COPY pgchangefeed/pyproject.toml pgchangefeed/pyproject.toml

# README.md liegt im committeten Baum auf Sprachwurzel-Ebene
# (sdks/python/README.md), eine Ebene über pyproject.toml -- setuptools
# erlaubt für `[project] readme` aber keinen Pfad außerhalb des
# Projekt-Wurzelverzeichnisses (`DistutilsOptionError: Cannot access
# '.../README.md' (or anything outside '...')`). Diese COPY-Zeile legt die
# Datei deshalb zusätzlich an die von pyproject.toml erwartete lokale
# Stelle (pgchangefeed/README.md, neben pyproject.toml) -- dasselbe Muster
# wie die `.proto`-Einspeisung in sdks/csharp/Dockerfile (COPY an eine vom
# Projekt erwartete relative Stelle, ohne sie im committeten Baum zu
# duplizieren).
COPY README.md pgchangefeed/README.md

COPY pgchangefeed/src pgchangefeed/src
COPY pgchangefeed/tests pgchangefeed/tests
COPY pgchangefeed/integration pgchangefeed/integration

# Der zusätzliche, benannte Bau-Kontext `proto` trägt NUR proto/ — diese
# COPY-Zeile ist der Punkt, an dem der Bau ohne den Zusatzkontext real
# abbricht (kein Fallback, kein übersprungener Schritt; Muster
# sdks/csharp/Dockerfile). Die flache Ziel-Zeile (`changestream.proto`, ohne
# Paketpfad) hält die erzeugten Stub-Module flach in `grpc_gen/` statt unter
# `cdc/stream/v1/` — genau wie die `.proto`-Einspeisung oben eine vom Aufruf
# erwartete relative Stelle.
COPY --from=proto cdc/stream/v1/changestream.proto changestream.proto

WORKDIR /src/pgchangefeed
# grpcio-tools vor der Stub-Erzeugung: ihr `protoc`-Modul (pip-Extra
# `test`) traegt den Plugin-Aufruf. Die Stub-Erzeugung schreibt in
# src/pgchangefeed/grpc_gen/ (Package-Marker committet, Stub-Module
# gitignore't); die Modul-Namen folgen dem `.proto`-Dateinamen
# (`changestream_pb2.py`/`changestream_pb2_grpc.py`, protoc-Konvention).
# Der `sed`-Schritt zieht den generierten flachen Geschwister-Import in
# `_pb2_grpc.py` auf den Package-Import — die dokumentierte Anpassung
# generierter Python-Stubs (grpc-Handbuch, "Generated Code"-Leitfaden,
# import statement), damit `pgchangefeed.grpc_gen.changestream_pb2_grpc`
# als Package-Modul laedt.
RUN pip install --no-cache-dir -e ".[test]"
RUN python -m grpc_tools.protoc -I.. --python_out=src/pgchangefeed/grpc_gen --grpc_python_out=src/pgchangefeed/grpc_gen changestream.proto \
    && sed -i 's/^import changestream_pb2 as changestream__pb2$/from pgchangefeed.grpc_gen import changestream_pb2 as changestream__pb2/' src/pgchangefeed/grpc_gen/changestream_pb2_grpc.py
# Der blanke `pytest`-Lauf sammelt nur die Unit-Tests (pyproject.toml
# testpaths = tests) — die Realserver-Tests unter pgchangefeed/integration/
# laufen nur in der `integration`-Stufe (expliziter Pfad-Aufruf dort).
RUN pytest

# --- pack: `uv build --no-sources` ueber den bereits gebauten und
# getesteten Stand. `uv` selbst wird NICHT per `pip install` bezogen, sondern
# per digest-gepinntem Multi-Stage-`COPY` aus Astrals eigenem Werkzeug-Image
# (Digest gegen die Registry geprueft: `docker buildx imagetools inspect
# ghcr.io/astral-sh/uv:0.12.17` liefert denselben Index-Digest). `uv build`
# benutzt das bestehende `setuptools.build_meta`-Backend ueber die
# PEP-517-Schnittstelle unveraendert -- kein Wechsel an pyproject.toml. Das
# Erzeugnis (Source-Distribution `.tar.gz` UND Wheel `.whl`) landet unter
# /out; das Wheel traegt die zur Bauzeit erzeugten Stub-Module mit
# (Regular-Unterpackage, setuptools packages.find). ---
FROM build AS pack
COPY --from=ghcr.io/astral-sh/uv:0.12.17@sha256:10787c682e4184e4f290de1171fd4703dc63de99221f10fe1c99002ce7fa9acc /uv /uvx /bin/
RUN mkdir -p /out && uv build --no-sources -o /out

# --- pack-export: gibt /out (.whl UND .tar.gz) als tar-Stream ueber stdout
# aus -- dasselbe Extraktionsmuster wie pack-export im
# sdks/csharp/Dockerfile (tools/harness/sdk-pack-csharp.sh):
# `tools/harness/sdk-pack-python.sh` baut diese Stufe und liest ihr
# ENTRYPOINT per `docker run --rm --network none <image> | tar -x -C
# sdks/python/dist/` aus -- kein Bind-Mount, kein `--user`-Workaround, die
# extrahierten Dateien gehoeren dadurch automatisch dem aufrufenden
# Nutzer. `tar` traegt beide Artefakte in einem Stream. ---
FROM pack AS pack-export
ENTRYPOINT ["tar", "-cf", "-", "-C", "/out", "."]

# --- integration: der Realserver-Integrationstest der SDK-Flaechen. Baut auf
# `build` auf (dasselbe installierte SDK samt `grpcio` und dieselben
# Stub-Module, kein zweiter Install-Pfad) und fuert pytest gegen eine reale,
# laufende Server-Instanz: `tools/harness/run-sdk-python-integration-tests.sh`
# baut diese Stufe (`--target integration`, zwingend mit dem benannten
# Bau-Kontext `proto` — derselbe Abbruchpunkt wie oben), startet sie im
# selben Docker-Netz wie den Feed-Container und liest die stdout-Marker
# (READY/RECEIVED/REJECTED) ueber `docker logs`. Der Aufruf traegt die
# Testdatei je Flaeche als Umgebungsvariable `PGCHANGEFEED_TEST_FILE` —
# der shell-Form-Guard unten (`:?`) scheitert laut bei blankem Aufruf,
# statt still die Unit-Tests ueber testpaths zu fahren (kein gruenes
# Default ohne jeden Integrationstest). `-u` haelt die Marker ungepuffert
# sichtbar; `--capture=no` reicht pytest-Prints ungefiltert nach stdout
# (der Runner pollt die Marker, waehrend der Test laeuft). ---
FROM build AS integration
WORKDIR /src/pgchangefeed
CMD ["/bin/sh", "-c", "python -u -m pytest \"${PGCHANGEFEED_TEST_FILE:?Testdatei je Flaeche erforderlich}\" -v --capture=no"]
