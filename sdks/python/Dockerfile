# sdks/python/Dockerfile — Docker-only Bau-/Testsetup des Python-SDK-Pakets
# `pgchangefeed` (ADR-0107 Festlegung 1/3/5, slice-sdk-python-projektgeruest).
# Bau-Kontext ist sdks/python/ selbst — NICHT die Repo-Wurzel und NICHT
# examples/ (kein examples/python/ existiert, ADR-0107 §Kontext Ist-Stand).
#
# Keine Runtime-Stufe hier: ein SDK-Package ist keine startbare Anwendung
# (analog sdks/csharp/Dockerfile). Die `pack`/`pack-export`-Stufen unten
# (slice-sdk-python-pack-werkzeug, ADR-0107 §Konsequenzen Folgepflicht 1,
# ADR-0108 §Entscheidung Festlegung 1/3) bauen auf `build` auf: `pytest`
# (unten) laeuft VOR `uv build --no-sources` (Stufe `pack`) — ein roter
# Test bricht den `docker build` mit Exit != 0 ab, bevor `pack` je erreicht
# wird (kein stiller Fallback, Muster sdks/csharp/Dockerfile).
#
# Basis-Digest real gemessen (`docker buildx imagetools inspect`,
# Manifest-Index-Digest über alle Plattformen) am Zeitpunkt dieses Zuges:
#   python:3.14-slim (3.14.7-slim-trixie, aktuelle Python-Version)
#     -> sha256:caaf356f40667c496d405780745b9ac25771c189a51dfcc42430d531ea09f8a2
# Pin-Hebung ist ein bewusster Commit (AGENTS.md §3.6/Modul 14), keine
# stille Auflösung.
FROM python:3.14-slim@sha256:caaf356f40667c496d405780745b9ac25771c189a51dfcc42430d531ea09f8a2 AS build
WORKDIR /src

# pyproject.toml zuerst — es ändert sich seltener als der Quellcode und
# hält den Docker-Layer-Cache über Quelltext-Änderungen hinweg stabil
# (analog sdks/csharp/Dockerfile).
COPY pgchangefeed/pyproject.toml pgchangefeed/pyproject.toml

# README.md liegt im committeten Baum auf Sprachwurzel-Ebene
# (sdks/python/README.md), eine Ebene über pyproject.toml -- setuptools
# erlaubt für `[project] readme` aber keinen Pfad außerhalb des
# Projekt-Wurzelverzeichnisses (real geprüft beim Schreiben dieses Zuges,
# Plan-Nachzug §3 der Slice-Datei: `DistutilsOptionError: Cannot access
# '.../README.md' (or anything outside '...')`). Diese COPY-Zeile legt die
# Datei deshalb zusätzlich an die von pyproject.toml erwartete lokale
# Stelle (pgchangefeed/README.md, neben pyproject.toml) -- dasselbe Muster
# wie die `.proto`-Einspeisung in sdks/csharp/Dockerfile (COPY an eine vom
# Projekt erwartete relative Stelle, ohne sie im committeten Baum zu
# duplizieren).
COPY README.md pgchangefeed/README.md

COPY pgchangefeed/src pgchangefeed/src
COPY pgchangefeed/tests pgchangefeed/tests

WORKDIR /src/pgchangefeed
RUN pip install --no-cache-dir -e ".[test]"
RUN pytest

# --- pack: `uv build --no-sources` ueber den bereits gebauten und
# getesteten Stand (slice-sdk-python-pack-werkzeug, ADR-0107 Festlegung 5,
# ADR-0108 §Entscheidung Festlegung 1/2/3). `uv` selbst wird NICHT per
# `pip install` bezogen, sondern per digest-gepinntem Multi-Stage-`COPY`
# aus Astrals eigenem Werkzeug-Image (exakter Wortlaut aus ADR-0108
# §Entscheidung Festlegung 3; Digest real gegen die Registry
# reverifiziert beim Schreiben dieses Zuges — `docker buildx imagetools
# inspect ghcr.io/astral-sh/uv:0.12.17` liefert denselben Index-Digest wie
# in der ADR genannt, keine Abweichung). `uv build` benutzt das
# bestehende `setuptools.build_meta`-Backend ueber die PEP-517-Schnittstelle
# unveraendert (ADR-0108 §Entscheidung Festlegung 2) -- kein Wechsel an
# pyproject.toml. Das Erzeugnis (Source-Distribution `.tar.gz` UND Wheel
# `.whl`) landet unter /out. ---
FROM build AS pack
COPY --from=ghcr.io/astral-sh/uv:0.12.17@sha256:10787c682e4184e4f290de1171fd4703dc63de99221f10fe1c99002ce7fa9acc /uv /uvx /bin/
RUN mkdir -p /out && uv build --no-sources -o /out

# --- pack-export: gibt /out (.whl UND .tar.gz) als tar-Stream ueber stdout
# aus -- dasselbe Extraktionsmuster wie pack-export im
# sdks/csharp/Dockerfile (tools/harness/sdk-pack-csharp.sh, ADR-0060):
# `tools/harness/sdk-pack-python.sh` baut diese Stufe und liest ihr
# ENTRYPOINT per `docker run --rm --network none <image> | tar -x -C
# sdks/python/dist/` aus -- kein Bind-Mount, kein `--user`-Workaround, die
# extrahierten Dateien gehoeren dadurch automatisch dem aufrufenden
# Nutzer. `tar` traegt beide Artefakte in einem Stream. ---
FROM pack AS pack-export
ENTRYPOINT ["tar", "-cf", "-", "-C", "/out", "."]
