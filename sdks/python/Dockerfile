# sdks/python/Dockerfile — Docker-only Bau-/Testsetup des Python-SDK-Pakets
# `pgchangefeed` (ADR-0107 Festlegung 1/3/5, slice-sdk-python-projektgeruest).
# Bau-Kontext ist sdks/python/ selbst — NICHT die Repo-Wurzel und NICHT
# examples/ (kein examples/python/ existiert, ADR-0107 §Kontext Ist-Stand).
#
# Keine Runtime-Stufe hier: ein SDK-Package ist keine startbare Anwendung
# (analog sdks/csharp/Dockerfile). Eine `pack`-Stufe (`python -m build`)
# folgt erst mit slice-sdk-python-pack-werkzeug (ADR-0107 §Konsequenzen
# Folgepflicht 1) — dieses Dockerfile trägt nur `pip install -e .[test]`
# und `pytest`.
#
# Basis-Digest real gemessen (`docker buildx imagetools inspect`,
# Manifest-Index-Digest über alle Plattformen) am Zeitpunkt dieses Zuges:
#   python:3.13-slim
#     -> sha256:8d9d0b8bcf6506481eae4907c18f5e3e7902e629f5f6d684f9e7c32e85e3ddf0
# Pin-Hebung ist ein bewusster Commit (AGENTS.md §3.6/Modul 14), keine
# stille Auflösung.
FROM python:3.13-slim@sha256:8d9d0b8bcf6506481eae4907c18f5e3e7902e629f5f6d684f9e7c32e85e3ddf0 AS build
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
