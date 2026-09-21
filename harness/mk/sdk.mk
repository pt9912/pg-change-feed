# harness/mk/sdk.mk — Werkzeug-Fragment für die SDK-Packages
# (`PgChangeFeed.Client`, C#/NuGet, ADR-0106; `pgchangefeed`, Python/PyPI,
# ADR-0107/ADR-0108; `pgchangefeed-kotlin`, Kotlin/GitHub Packages,
# ADR-0109). Kein Gate: `dotnet restore`/PyPI-Paketbezug/GitHub-Packages-
# Paketbezug braucht Netz, `make gates` bleibt netzlos (ADR-0106
# Festlegung 4, ADR-0107 Festlegung 5, ADR-0109 Festlegung 5, dieselbe
# Begründung wie harness/mk/examples.mk) — dieses Fragment hängt deshalb
# NICHT an GATE_CHECKS.
#
# `sdk-pack-csharp` baut/testet/paketiert das C#-SDK Docker-only im
# gepinnten mcr.microsoft.com/dotnet/sdk-Image (sdks/csharp/Dockerfile,
# Stufe `pack-export`): `dotnet test` gegen alle vier Testflächen (HTTP +
# gRPC + SSE + NATS-Vollinhalt, dasselbe Testprojekt
# PgChangeFeed.Client.Tests) läuft VOR `dotnet pack`
# in derselben Docker-Bau-Kette — ein roter Test bricht den `docker build`
# mit Exit != 0 ab, bevor die `pack`-Stufe je erreicht wird (kein stiller
# Fallback, Muster harness/mk/examples.mk). Wie beim gRPC-Bau der Fläche
# trägt der Aufruf zwingend `--build-context proto=proto` (ADR-0090
# Festlegung 2, übernommen aus slice-sdk-csharp-grpc-client-flaeche) — ohne
# ihn bricht der Bau an der `COPY --from=proto`-Zeile in
# sdks/csharp/Dockerfile ab.
#
# Export analog `make proto-generate` (tools/harness/proto-generate.sh):
# tools/harness/sdk-pack-csharp.sh extrahiert das .nupkg host-seitig aus der
# `pack-export`-Stufe (`docker run --rm --network none <image> | tar -x`,
# `set -o pipefail` unter bash, AGENTS.md §3.9) nach sdks/csharp/dist/
# (`.gitignore`t). Erzeugnis: PgChangeFeed.Client.0.2.0.nupkg. Exit-Code des
# Skripts wird wie bei jedem anderen Ziel direkt gelesen.
.PHONY: sdk-pack-csharp
sdk-pack-csharp: ## C#-SDK bauen+testen+paketieren (sdks/csharp, .nupkg nach sdks/csharp/dist/; Werkzeug, kein Gate; ADR-0106)
	@bash tools/harness/sdk-pack-csharp.sh

# `sdk-pack-kotlin` baut/testet/paketiert das Kotlin-SDK Docker-only im
# gepinnten eclipse-temurin:21-jdk-Image (sdks/kotlin/Dockerfile, Stufe
# `pack-export`): `./gradlew test` gegen BEIDE Testflächen (HTTP + gRPC)
# läuft VOR dem Jar-Bau in derselben Docker-Bau-Kette — ein roter Test
# bricht den `docker build` mit Exit != 0 ab, bevor die `pack`-Stufe je
# erreicht wird (kein stiller Fallback, Muster harness/mk/examples.mk). Wie
# beim gRPC-Bau der Fläche trägt der Aufruf zwingend
# `--build-context proto=proto` (ADR-0090 Festlegung 2, übernommen auf den
# Kotlin-SDK-Baum) — ohne ihn bricht der Bau an der `COPY --from=proto`-Zeile
# in sdks/kotlin/Dockerfile ab.
#
# Export analog `make sdk-pack-csharp`/`make sdk-pack-python`
# (tools/harness/sdk-pack-kotlin.sh): das Skript extrahiert das erzeugte
# .jar host-seitig aus der `pack-export`-Stufe (`docker run --rm --network
# none <image> | tar -x`, `set -o pipefail` unter bash, AGENTS.md §3.9) nach
# sdks/kotlin/dist/ (`.gitignore`t). Erzeugnis:
# pgchangefeed-kotlin-0.1.0.jar. Kein Sources-/Javadoc-Jar — GitHub Packages
# verlangt laut offizieller Dokumentation keines (real recherchiert,
# sdks/kotlin/Dockerfile Stufe `pack`). Exit-Code des Skripts wird wie bei
# jedem anderen Ziel direkt gelesen.
.PHONY: sdk-pack-kotlin
sdk-pack-kotlin: ## Kotlin-SDK bauen+testen+paketieren (sdks/kotlin, .jar nach sdks/kotlin/dist/; Werkzeug, kein Gate; ADR-0109)
	@bash tools/harness/sdk-pack-kotlin.sh

# `sdk-pack-python` baut/testet/paketiert das Python-SDK Docker-only im
# gepinnten python:3.14-slim-Image (sdks/python/Dockerfile, Stufe
# `pack-export`): `pytest` läuft VOR `uv build --no-sources` in derselben
# Docker-Bau-Kette — ein roter Test bricht den `docker build` mit
# Exit != 0 ab, bevor die `pack`-Stufe je erreicht wird (kein stiller
# Fallback, Muster harness/mk/examples.mk). `uv` wird per digest-gepinntem
# Multi-Stage-`COPY` aus Astrals eigenem Werkzeug-Image bezogen, nicht per
# `pip install uv` (ADR-0108 §Entscheidung Festlegung 1/3) — das
# `setuptools.build_meta`-Backend in pyproject.toml bleibt unverändert
# (ADR-0108 §Entscheidung Festlegung 2).
#
# Export analog `make sdk-pack-csharp`
# (tools/harness/sdk-pack-python.sh): das Skript extrahiert BEIDE
# Artefakte (.whl UND .tar.gz) host-seitig aus der `pack-export`-Stufe
# (`docker run --rm --network none <image> | tar -x`, `set -o pipefail`
# unter bash, AGENTS.md §3.9) nach sdks/python/dist/ (`.gitignore`t).
# Exit-Code des Skripts wird wie bei jedem anderen Ziel direkt gelesen.
.PHONY: sdk-pack-python
sdk-pack-python: ## Python-SDK bauen+testen+paketieren (sdks/python, .whl+.tar.gz nach sdks/python/dist/; Werkzeug, kein Gate; ADR-0107, ADR-0108)
	@bash tools/harness/sdk-pack-python.sh
