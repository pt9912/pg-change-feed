# harness/mk/sdk.mk — Werkzeug-Fragment für die SDK-Packages
# (`PgChangeFeed.Client`, C#/NuGet, ADR-0106; `pgchangefeed`, Python/PyPI,
# ADR-0107/ADR-0108; `pgchangefeed-kotlin`, Kotlin/GitHub Packages,
# ADR-0109). Kein Gate: `dotnet restore`/PyPI-Paketbezug/GitHub-Packages-
# Paketbezug braucht Netz, `make gates` bleibt netzlos (ADR-0106
# Festlegung 4, ADR-0107 Festlegung 5, ADR-0109 Festlegung 5, dieselbe
# Begründung wie harness/mk/examples.mk) — dieses Fragment hängt deshalb
# NICHT an GATE_CHECKS.
#
# `sdk-public-doc-check` prueft, dass keine Datei unter sdks/ eine interne
# Kennung (SPEC-/ADR-/ARC-/LH-FA-/LH-QA-, Slice-/Welle-Name) traegt:
# Kommentare, Docstrings, Fehlertexte, README und Build-Dateien der SDKs
# erreichen Anwender ueber die Pakete (Wheel/sdist, nupkg mit XML-Doku,
# Sources-Jar). Reines grep, netzlos und schnell
# (tools/harness/sdk-public-doc-check.sh); die drei `sdk-pack-*`-Ziele haengen
# davon ab, damit ein Rueckfall vor dem Bau auffaellt. Bewusst kein Teil von
# GATE_CHECKS: ein weiteres Gate aendert die Gate-Liste in harness/README.md
# und ihre Sensor-Bindung; das Ziel bleibt Werkzeug mit eigenem Tabellentest
# (`make test-sdk-public-doc-check`).
.PHONY: sdk-public-doc-check
sdk-public-doc-check: ## Keine interne Kennung in den Dateien unter sdks/ (netzlos, grep; Vorstufe der sdk-pack-*-Ziele; Werkzeug, kein Gate)
	@bash tools/harness/sdk-public-doc-check.sh

.PHONY: test-sdk-public-doc-check
test-sdk-public-doc-check: ## Tabellentest gegen tools/harness/sdk-public-doc-check.sh (netzlos)
	@bash tools/harness/run-sdk-public-doc-check-tests.sh

# `test-sdk-dist-clean` prueft die Hilfsfunktionen, mit denen die drei
# `sdk-pack-*`-Skripte sdks/<sprache>/dist/ vor dem Export ersetzen
# (tools/harness/sdk-dist-clean.sh): Altlasten verschwinden, unzulaessige
# Pfade werden abgelehnt. Netzlos, Werkzeug wie die uebrigen `test-sdk-*`.
.PHONY: test-sdk-dist-clean
test-sdk-dist-clean: ## Tabellentest gegen tools/harness/sdk-dist-clean.sh (netzlos)
	@bash tools/harness/run-sdk-dist-clean-tests.sh

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
# (`.gitignore`t). Erzeugnis: PgChangeFeed.Client.0.2.1.nupkg. Exit-Code des
# Skripts wird wie bei jedem anderen Ziel direkt gelesen.
.PHONY: sdk-pack-csharp
sdk-pack-csharp: sdk-public-doc-check ## C#-SDK bauen+testen+paketieren (sdks/csharp, .nupkg nach sdks/csharp/dist/; Werkzeug, kein Gate; ADR-0106)
	@bash tools/harness/sdk-pack-csharp.sh

# `sdk-pack-kotlin` baut/testet/paketiert das Kotlin-SDK Docker-only im
# gepinnten eclipse-temurin:21-jdk-Image (sdks/kotlin/Dockerfile, Stufe
# `pack-export`): `./gradlew test` gegen VIER Testflächen (HTTP + gRPC + SSE
# + NATS-Vollinhalt) läuft VOR dem Jar-Bau in derselben Docker-Bau-Kette —
# ein roter Test bricht den `docker build` mit Exit != 0 ab, bevor die
# `pack`-Stufe je erreicht wird (kein stiller Fallback, Muster
# harness/mk/examples.mk). Wie beim gRPC-Bau der Fläche trägt der Aufruf
# zwingend `--build-context proto=proto` (ADR-0090 Festlegung 2, übernommen
# auf den Kotlin-SDK-Baum) — ohne ihn bricht der Bau an der `COPY
# --from=proto`-Zeile in sdks/kotlin/Dockerfile ab.
#
# Export analog `make sdk-pack-csharp`/`make sdk-pack-python`
# (tools/harness/sdk-pack-kotlin.sh): das Skript extrahiert das erzeugte
# .jar host-seitig aus der `pack-export`-Stufe (`docker run --rm --network
# none <image> | tar -x`, `set -o pipefail` unter bash, AGENTS.md §3.9) nach
# sdks/kotlin/dist/ (`.gitignore`t). Erzeugnisse:
# pgchangefeed-kotlin-0.2.1.jar und pgchangefeed-kotlin-0.2.1-sources.jar
# (`java { withSourcesJar() }` in build.gradle.kts; die Quellen tragen die KDoc).
# Exit-Code des Skripts wird wie bei jedem anderen Ziel direkt gelesen.
.PHONY: sdk-pack-kotlin
sdk-pack-kotlin: sdk-public-doc-check ## Kotlin-SDK bauen+testen+paketieren (sdks/kotlin, .jar nach sdks/kotlin/dist/; Werkzeug, kein Gate; ADR-0109)
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
# Der gRPC-Teil der Flaeche liest die `.proto` ueber den benannten
# Bau-Kontext `proto` (ADR-0090 Festlegung 2, Muster sdks/csharp) —
# tools/harness/sdk-pack-python.sh traegt ihn zwingend.
.PHONY: sdk-pack-python
sdk-pack-python: sdk-public-doc-check ## Python-SDK bauen+testen+paketieren (sdks/python, .whl+.tar.gz nach sdks/python/dist/; Werkzeug, kein Gate; ADR-0107, ADR-0108, ADR-0110)
	@bash tools/harness/sdk-pack-python.sh

# `test-sdk-kotlin-integration` ist der Realserver-Integrationstest der
# Kotlin-SDK-Zustellweg-Flaechen (slice-sdk-kotlin-reale2e, Mechanik-Klasse
# ADR-0110 Festlegung 2, gespiegelt vom C#-Werkzeug): das Skript
# tools/harness/run-sdk-kotlin-integration-tests.sh faehrt die
# compose.yaml-Umgebung hoch, baut die `integration`-Docker-Stufe des
# Kotlin-SDK (sdks/kotlin/Dockerfile, zwingend mit dem benannten
# Bau-Kontext `proto`) und startet je Flaeche eine Phase
# (`./gradlew integrationTest --tests` je Testklasse, Testklasse als
# Umgebungsvariable mit `:?`-Guard gegen blanken Aufruf) im selben
# Docker-Netz wie den Feed-Container — der Pruefling ist die kompilierte
# Client-Assembly. Kein Gate (braucht DB-Zugang/Docker/Netz, dieselbe
# Klasse wie `make test-integration`); setzt ein geladenes :dev-Image
# voraus (`make image` vorher, compose.yaml traegt keinen build:-Block,
# ADR-0044). Der Runner schreibt den Kotlin-Abschnitt des Abdeckungs-
# Traegers docs/user/sdk-e2e-abdeckung.md aus derselben Messung.
.PHONY: test-sdk-kotlin-integration
test-sdk-kotlin-integration: ## Kotlin-SDK-Realserver-Integrationstest (compose + integration-Stufe, vier Phasen; Werkzeug, kein Gate; slice-sdk-kotlin-reale2e)
	@bash tools/harness/run-sdk-kotlin-integration-tests.sh

# `test-sdk-csharp-integration` ist der Realserver-Integrationstest der
# C#-SDK-Zustellweg-Flaechen (slice-sdk-csharp-reale2e, Mechanik-Klasse
# ADR-0110 Festlegung 2, gespiegelt vom Python-Werkzeug): das Skript
# tools/harness/run-sdk-csharp-integration-tests.sh faehrt die
# compose.yaml-Umgebung hoch (PostgreSQL/NATS/Feed-Container, Schema-Rollout
# ueber d-migrate, Vorbedingungen der Aktivierung), baut die
# `integration`-Docker-Stufe des C#-SDK (sdks/csharp/Dockerfile, zwingend
# mit dem benannten Bau-Kontext `proto`) und startet je Flaeche eine Phase
# im selben Docker-Netz wie den Feed-Container — der Pruefling ist die
# kompilierte Client-Assembly. Kein Gate (braucht DB-Zugang/Docker/Netz,
# dieselbe Klasse wie `make test-integration`); setzt ein geladenes
# :dev-Image voraus (`make image` vorher, compose.yaml traegt keinen
# build:-Block, ADR-0044). Der Runner schreibt den C#-Abschnitt des
# Abdeckungs-Traegers docs/user/sdk-e2e-abdeckung.md aus derselben Messung.
.PHONY: test-sdk-csharp-integration
test-sdk-csharp-integration: ## C#-SDK-Realserver-Integrationstest (compose + integration-Stufe, vier Phasen; Werkzeug, kein Gate; slice-sdk-csharp-reale2e)
	@bash tools/harness/run-sdk-csharp-integration-tests.sh

# `test-sdk-python-integration` ist der Realserver-Integrationstest der
# Python-SDK-Zustellweg-Flaechen (ADR-0110 §Entscheidung Festlegung
# 2/Folgepflicht 1, slice-sdk-python-grpc-client-flaeche): das Skript
# tools/harness/run-sdk-python-integration-tests.sh faehrt die
# compose.yaml-Umgebung hoch (PostgreSQL/NATS/Feed-Container, Schema-Rollout
# ueber d-migrate, Vorbedingungen der Aktivierung), baut die
# `integration`-Docker-Stufe des Python-SDK (sdks/python/Dockerfile,
# zwingend mit dem benannten Bau-Kontext `proto`) und startet sie im selben
# Docker-Netz wie den Feed-Container — der Pruefling ist das SDK selbst.
# Kein Gate (braucht DB-Zugang/Docker/Netz, dieselbe Klasse wie
# `make test-integration`); setzt ein geladenes :dev-Image voraus
# (`make image` vorher, compose.yaml traegt keinen build:-Block, ADR-0044).
.PHONY: test-sdk-python-integration
test-sdk-python-integration: ## Python-SDK-Realserver-Integrationstest (compose + integration-Stufe; Werkzeug, kein Gate; ADR-0110)
	@bash tools/harness/run-sdk-python-integration-tests.sh
