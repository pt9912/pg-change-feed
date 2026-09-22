# sdks/kotlin/Dockerfile — Docker-only Bau-/Testsetup des Kotlin-SDK-Pakets
# `pgchangefeed-kotlin` (ADR-0109 Festlegung 3/5, slice-sdk-kotlin-projektgeruest).
# Bau-Kontext ist sdks/kotlin/ selbst — NICHT die Repo-Wurzel und NICHT
# examples/kotlin/; beide bleiben von diesem Bau unberührt (ADR-0109
# Festlegung 3: eigenständiger Baum, kein `project(":…")`-Abhängigkeitspfad
# auf einen privaten Baum dieses Repos).
#
# Die gRPC-Stream-Client-Fläche (slice-sdk-kotlin-grpc-client-flaeche,
# ADR-0109 Festlegung 1/3, SPEC-020) liest die `.proto` NICHT aus diesem
# Kontext — sie liegt außerhalb davon, unter proto/cdc/stream/v1/ an der
# Repo-Wurzel. Sie kommt über einen ZUSÄTZLICHEN, benannten Bau-Kontext
# namens `proto` herein (`COPY --from=proto …` unten); der Aufrufer bringt
# ihn mit `docker build --build-context proto=proto …` mit (Muster
# examples/kotlin/Dockerfile, harness/mk/examples.mk). Ohne ihn bricht der
# Bau an der COPY-Zeile ab — kein stiller Fallback.
#
# Keine Runtime-Stufe hier: ein SDK-Package ist keine startbare Anwendung
# (analog sdks/csharp/Dockerfile, sdks/python/Dockerfile). Die
# `pack`/`pack-export`-Stufen unten (slice-sdk-kotlin-pack-werkzeug, ADR-0109
# §Konsequenzen Folgepflicht 1) bauen auf `build` auf: `./gradlew test`
# (unten) laeuft VOR dem Jar-Bau — ein roter Test bricht den `docker build`
# mit Exit != 0 ab, bevor `pack` je erreicht wird (kein stiller Fallback).
# `./gradlew publish` laeuft ebenfalls in einer eigenen Docker-Stufe
# (`publish`, ganz unten) — ADR-0109 Festlegung 5 legt woertlich Docker-only
# bis EINSCHLIESSLICH `publish` fest; eine fruehere Fassung dieses Workflows
# rief `./gradlew publish` direkt auf dem Runner auf, ausserhalb jedes
# Containers (Review-Finding F-1,
# docs/reviews/review-slice-sdk-kotlin-publish-workflow.md) — korrigiert in
# der Fixrunde zu slice-sdk-kotlin-publish-workflow.
#
# JDK-Basis: `eclipse-temurin:21-jdk`, NICHT `25-jdk` (ADR-0109
# Festlegung 5 — Gradle 8.14 unterstützt JDK 25 nur als Toolchain-Ziel ab
# Gradle 9.1.0; das SDK bleibt bei derselben, bereits real gebauten
# JDK-/Gradle-Kombination wie examples/kotlin/Dockerfile). Basis-Digest real
# gemessen (`docker buildx imagetools inspect`, linux/amd64) am Zeitpunkt
# dieses Zuges (2026-09-20) — identisch mit dem in
# examples/kotlin/Dockerfile gepinnten Kandidaten (dieselbe Basis, real
# re-verifiziert statt blind übernommen, AGENTS.md §3.12):
#   eclipse-temurin:21-jdk
#     -> sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e
# Pin-Hebung ist ein bewusster Commit (AGENTS.md §3.6/Modul 14), keine stille
# Auflösung.
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

# Der zusätzliche, benannte Bau-Kontext `proto` trägt NUR proto/ (Muster
# ADR-0090 Festlegung 2, hier auf den SDK-Baum übertragen) — diese
# COPY-Zeile ist der Punkt, an dem der Bau ohne den Zusatzkontext real
# abbricht (kein Fallback, kein übersprungener Schritt). Das Ziel liegt
# unter der von `com.google.protobuf`s Gradle-Plugin standardmäßig
# gelesenen Quellmenge `src/main/proto/` (pgchangefeed-kotlin/build.gradle.kts),
# dasselbe Muster wie examples/kotlin/Dockerfile/grpc-client.
COPY --from=proto cdc/stream/v1/changestream.proto src/main/proto/changestream.proto

RUN ./gradlew --no-daemon test
RUN ./gradlew --no-daemon build

# --- pack: kopiert den bereits gebauten Jar-Stand (`./gradlew build` oben,
# inklusive `./gradlew test` davor) nach /out (slice-sdk-kotlin-pack-werkzeug,
# ADR-0109 §Konsequenzen Folgepflicht 1). Der Standard-`jar`-Task läuft
# bereits als Teil von `build`/`assemble` und erzeugt
# build/libs/pgchangefeed-kotlin-0.2.0.jar (Gradle-Standard-Namensform aus
# group/artifactId/version in build.gradle.kts) — kein zusätzlicher
# Gradle-Aufruf nötig. Kein Sources-/Javadoc-Jar: GitHub Packages verlangt
# laut offizieller Dokumentation
# (docs.github.com/.../publishing-java-packages-with-gradle, real abgerufen)
# keines, und `from(components["java"])` (build.gradle.kts) packt
# ausschließlich den Haupt-Jar, solange kein `withSourcesJar()`/
# `withJavadocJar()` im `java{}`-Block ergänzt ist (Gradle-Kern-Verhalten
# des `maven-publish`-Plugins). ---
FROM build AS pack
RUN mkdir -p /out && cp build/libs/*.jar /out/

# --- pack-export: gibt /out (den Jar) als tar-Stream über stdout aus —
# dasselbe Extraktionsmuster wie sdk-pack-csharp/sdk-pack-python:
# `tools/harness/sdk-pack-kotlin.sh` baut diese Stufe und liest ihr
# ENTRYPOINT per `docker run --rm --network none <image> | tar -x -C
# sdks/kotlin/dist/` aus — kein Bind-Mount, kein `--user`-Workaround, die
# extrahierte Datei gehört dadurch automatisch dem aufrufenden Nutzer. ---
FROM pack AS pack-export
ENTRYPOINT ["tar", "-cf", "-", "-C", "/out", "."]

# --- publish: veröffentlicht den bereits gebauten Jar nach GitHub Packages
# (slice-sdk-kotlin-publish-workflow, ADR-0109 Festlegung 5 wörtlich:
# Docker-only bis EINSCHLIESSLICH `publish`). Baut auf `build` auf (NICHT
# auf `pack`/`pack-export` — beide sind nur die Export-Kette für das
# lokale/netzlose Werkzeug `make sdk-pack-kotlin`, für `publish` irrelevant):
# `./gradlew --no-daemon build` ist zu diesem Zeitpunkt des Bau-Layers
# bereits gelaufen, `build/libs/pgchangefeed-kotlin-0.2.0.jar` liegt bereits
# vor. `./gradlew publish` läuft NICHT während `docker build` (kein
# RUN-Schritt in dieser Stufe) — der Publish-Schritt braucht
# Registry-Netzwerkzugriff (https://maven.pkg.github.com) erst zur
# LAUFZEIT, über `docker run` ohne `--network none` (anders als jede andere
# Stufe dieses Dockerfiles, die netzlos bau-/testbar bleibt). Zur Laufzeit
# findet Gradles Up-to-date-Prüfung für `generateProto`/`compileKotlin`/
# `jar` den seit dem Bau unveränderten Dateizustand bereits erfüllt vor und
# lädt nur den bereits gebauten Jar hoch — kein zweiter Compile-/
# Codegenerierungs-Lauf außerhalb von Docker (Korrektur von
# docs/reviews/review-slice-sdk-kotlin-publish-workflow.md Finding F-1: eine
# frühere Fassung dieses Workflows rief `./gradlew publish` direkt auf dem
# GitHub-hosted Runner auf, außerhalb jedes Containers, und wiederholte
# dabei die gesamte Bau-Kette ein zweites Mal, ohne den bereits erzeugten
# Jar zu verwenden).
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
