# sdks/kotlin/Dockerfile — Docker-only Bau-/Testsetup des Kotlin-SDK-Pakets
# `pgchangefeed-kotlin`. Bau-Kontext ist sdks/kotlin/ selbst — NICHT die
# Repo-Wurzel und NICHT examples/kotlin/; beide bleiben von diesem Bau
# unberührt (eigenständiger Baum, kein `project(":…")`-Abhängigkeitspfad auf
# einen privaten Baum dieses Repos).
#
# Die gRPC-Stream-Client-Fläche liest die `.proto` NICHT aus diesem Kontext —
# sie liegt außerhalb davon, unter proto/cdc/stream/v1/ an der Repo-Wurzel. Sie
# kommt über einen ZUSÄTZLICHEN, benannten Bau-Kontext namens `proto` herein
# (`COPY --from=proto …` unten); der Aufrufer bringt ihn mit `docker build
# --build-context proto=proto …` mit (Muster examples/kotlin/Dockerfile,
# harness/mk/examples.mk). Ohne ihn bricht der Bau an der COPY-Zeile ab —
# kein stiller Fallback.
#
# Keine Runtime-Stufe hier: ein SDK-Package ist keine startbare Anwendung
# (analog sdks/csharp/Dockerfile, sdks/python/Dockerfile). Die
# `pack`/`pack-export`-Stufen unten bauen auf `build` auf: `./gradlew test`
# (unten) läuft VOR dem Jar-Bau — ein roter Test bricht den `docker build`
# mit Exit != 0 ab, bevor `pack` je erreicht wird (kein stiller Fallback).
# `./gradlew publish` läuft ebenfalls in einer eigenen Docker-Stufe
# (`publish`, ganz unten): das Veröffentlichen ist Docker-only wie jeder
# andere Schritt.
#
# JDK-Basis: `eclipse-temurin:21-jdk`, NICHT `25-jdk` (Gradle 8.14 unterstützt
# JDK 25 nur als Toolchain-Ziel ab Gradle 9.1.0; das SDK bleibt bei derselben,
# bereits real gebauten JDK-/Gradle-Kombination wie examples/kotlin/Dockerfile).
# Basis-Digest gemessen (`docker buildx imagetools inspect`, linux/amd64,
# 2026-09-20) — identisch mit dem in examples/kotlin/Dockerfile gepinnten
# Kandidaten (dieselbe Basis):
#   eclipse-temurin:21-jdk
#     -> sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e
# Pin-Hebung ist ein bewusster Commit, keine stille Auflösung.
FROM eclipse-temurin:21-jdk@sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e AS build
WORKDIR /src

# Das Projekt liegt unter pgchangefeed-kotlin/ (Bau-Kontext ist sdks/kotlin/
# selbst, analog sdks/python/Dockerfile: pyproject.toml liegt unter
# pgchangefeed/ und WORKDIR wechselt vor dem Bau dorthin). Der Gradle-Wrapper
# und die Projekt-Manifeste zuerst — sie ändern sich seltener als der
# Quellcode und halten den Docker-Layer-Cache über Quelltext-Änderungen
# hinweg stabil (Muster: examples/kotlin/Dockerfile).
COPY pgchangefeed-kotlin/gradlew pgchangefeed-kotlin/gradlew
COPY pgchangefeed-kotlin/gradle/ pgchangefeed-kotlin/gradle/
COPY pgchangefeed-kotlin/settings.gradle.kts pgchangefeed-kotlin/build.gradle.kts pgchangefeed-kotlin/

WORKDIR /src/pgchangefeed-kotlin
RUN chmod +x ./gradlew && ./gradlew --no-daemon help >/dev/null

COPY pgchangefeed-kotlin/src src

# Der zusätzliche, benannte Bau-Kontext `proto` trägt NUR proto/ — diese
# COPY-Zeile ist der Punkt, an dem der Bau ohne den Zusatzkontext real
# abbricht (kein Fallback, kein übersprungener Schritt). Das Ziel liegt unter
# der von `com.google.protobuf`s Gradle-Plugin standardmäßig gelesenen
# Quellmenge `src/main/proto/` (pgchangefeed-kotlin/build.gradle.kts), dasselbe
# Muster wie examples/kotlin/Dockerfile/grpc-client.
COPY --from=proto cdc/stream/v1/changestream.proto src/main/proto/changestream.proto

RUN ./gradlew --no-daemon test
RUN ./gradlew --no-daemon build

# --- pack: kopiert den bereits gebauten Stand (`./gradlew build` oben,
# inklusive `./gradlew test` davor) nach /out. Der Standard-`jar`-Task und der
# `sourcesJar`-Task (`java { withSourcesJar() }` in build.gradle.kts) laufen
# bereits als Teil von `build`/`assemble` und erzeugen
# build/libs/pgchangefeed-kotlin-0.2.1.jar und
# build/libs/pgchangefeed-kotlin-0.2.1-sources.jar (Gradle-Namensform aus
# group/artifactId/version in build.gradle.kts) — kein zusätzlicher
# Gradle-Aufruf nötig. `from(components["java"])` (build.gradle.kts) nimmt das
# Sources-Jar in die Veröffentlichung auf. ---
FROM build AS pack
RUN mkdir -p /out && cp build/libs/*.jar /out/

# --- pack-export: gibt /out (die Jars) als tar-Stream über stdout aus —
# dasselbe Extraktionsmuster wie sdk-pack-csharp/sdk-pack-python:
# `tools/harness/sdk-pack-kotlin.sh` baut diese Stufe und liest ihr
# ENTRYPOINT per `docker run --rm --network none <image> | tar -x -C
# sdks/kotlin/dist/` aus — kein Bind-Mount, kein `--user`-Workaround, die
# extrahierten Dateien gehören dadurch automatisch dem aufrufenden Nutzer. ---
FROM pack AS pack-export
ENTRYPOINT ["tar", "-cf", "-", "-C", "/out", "."]

# --- publish: veröffentlicht den bereits gebauten Stand nach GitHub Packages
# (Docker-only bis einschließlich `publish`). Baut auf `build` auf (NICHT auf
# `pack`/`pack-export` — beide sind nur die Export-Kette für das
# lokale/netzlose Werkzeug `make sdk-pack-kotlin`, für `publish` irrelevant):
# `./gradlew --no-daemon build` ist zu diesem Zeitpunkt des Bau-Layers
# bereits gelaufen, `build/libs/pgchangefeed-kotlin-0.2.1.jar` liegt bereits
# vor. `./gradlew publish` läuft NICHT während `docker build` (kein
# RUN-Schritt in dieser Stufe) — der Publish-Schritt braucht
# Registry-Netzwerkzugriff (https://maven.pkg.github.com) erst zur
# LAUFZEIT, über `docker run` ohne `--network none` (anders als jede andere
# Stufe dieses Dockerfiles, die netzlos bau-/testbar bleibt). Zur Laufzeit
# findet Gradles Up-to-date-Prüfung für `generateProto`/`compileKotlin`/
# `jar` den seit dem Bau unveränderten Dateizustand bereits erfüllt vor und
# lädt nur den bereits gebauten Stand hoch — kein zweiter Compile-/
# Codegenerierungs-Lauf außerhalb von Docker.
#
# Zugangsdaten (GITHUB_ACTOR/GITHUB_TOKEN) werden ausschließlich zur
# Laufzeit über `docker run -e GITHUB_ACTOR=... -e GITHUB_TOKEN=...`
# durchgereicht, NICHT als Build-Arg — ein Build-Arg landet im
# Image-Layer-/History-Cache, ein Secret darf dort nicht auftauchen (kein
# `ARG GITHUB_TOKEN`/`ARG GITHUB_ACTOR` in diesem Dockerfile).
#
# CMD statt ENTRYPOINT: Ein ENTRYPOINT im Exec-Format würde ein am
# `docker run`-Aufruf zusätzlich mitgegebenes Kommando an das ENTRYPOINT-Array
# ANHÄNGEN statt es zu ERSETZEN (Docker-Kernverhalten) — der aufrufende
# Workflow gibt das Kommando trotzdem explizit mit
# (`docker run … <image> ./gradlew --no-daemon publish`, Konsistenz mit den
# übrigen Aufrufen dieses Workflows), CMD wird dabei vollständig ersetzt,
# kein Argument-Duplikat.
FROM build AS publish
CMD ["./gradlew", "--no-daemon", "publish"]

# --- integration: der Realserver-Integrationstest der SDK-Flächen. Baut auf
# `build` auf (dieselben Pins, kein neuer Pin), baut die Integrations-Quellmenge
# (build.gradle.kts-SourceSet `integrationTest` — der Unit-Lauf
# `./gradlew test`/`build` bleibt unberührt, der Task hängt bewusst NICHT an
# `check`) und führt die Testläufe je Fläche als
# `./gradlew integrationTest --tests` aus — die Testklasse je Phase kommt als
# Umgebungsvariable, deren Guard (`:?`) blanken Aufruf laut scheitern lässt
# (kein grünes Default ohne jeden Integrationstest, Muster
# sdks/python/Dockerfile). Restore braucht Netz (Gradle-Abhängigkeiten), der
# Testlauf selbst braucht DB-/Netz-Zugang und ist kein Gate. ---
FROM build AS integration
RUN ./gradlew --no-daemon integrationTestClasses
WORKDIR /src/pgchangefeed-kotlin
CMD ["/bin/sh", "-c", "exec ./gradlew --no-daemon integrationTest --tests \"io.github.pt9912.pgchangefeed.integration.${PGCHANGEFEED_TEST_NAME:?Testklasse je Phase erforderlich}\""]
