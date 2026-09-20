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
# `pack`/`pack-export`-Stufen (slice-sdk-kotlin-pack-werkzeug, ADR-0109
# §Konsequenzen Folgepflicht 1) bauen später auf `build` auf — dieser Slice
# liefert nur `./gradlew build`/`test`, kein `publish`-Aufruf.
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
