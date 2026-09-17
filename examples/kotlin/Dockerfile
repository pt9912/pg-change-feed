# examples/kotlin/Dockerfile — Werkzeugketten-Image der Kotlin-Sprach-Wurzel
# (ADR-0087 Festlegung 3, ADR-0090 Festlegung 1/5). Bau-Kontext ist
# examples/kotlin/ selbst — NICHT die Repo-Wurzel; die Wurzel-`Dockerfile`
# und `.dockerignore` bleiben davon unberührt (ADR-0087 Festlegung 3,
# Kontext-Bindung 4).
#
# Das gRPC-Programm (grpc-client, ADR-0090 Festlegung 2, slice-103 — die in
# slice-102 am C#-Bau gemessene Form übertragen) liest die `.proto` NICHT aus
# diesem Kontext — sie liegt außerhalb davon, unter proto/cdc/stream/v1/ an
# der Repo-Wurzel. Sie kommt über einen ZUSÄTZLICHEN, benannten Bau-Kontext
# namens `proto` herein (`COPY --from=proto …` unten); der Aufrufer bringt
# ihn mit `docker build --build-context proto=proto …` mit
# (harness/mk/examples.mk). Ohne ihn bricht der Bau an der COPY-Zeile ab —
# kein stiller Fallback (ADR-0090 §Fitness Function).
#
# Basis-Digests real gemessen (`docker manifest inspect`, amd64/linux) am
# Zeitpunkt dieses Zuges (2026-09-17):
#   eclipse-temurin:21-jdk
#     -> sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e
#   eclipse-temurin:21-jre
#     -> sha256:ca7551d4f36647207812e3c8c8596c89fd4e5226a061a5b8c8629839fe8969f6
# Pin-Hebung ist ein bewusster Commit (AGENTS.md §3.6/Modul 14), keine stille
# Auflösung.

# Stufe `build`: `./gradlew test`/`installDist` im gepinnten JDK-Image, für
# **alle vier** Module dieser Sprach-Wurzel (ADR-0090 Festlegung 1/5 — ein
# `make`-Ziel je Sprache trägt den wachsenden Umfang, hier vier Programme
# statt dreier seit slice-103). Der Paket-Bezug (Maven Central, Gradle Plugin
# Portal, der Gradle-Distributions-Download selbst) braucht Netz — deshalb
# bleibt `make examples-kotlin` ein Werkzeug, kein Gate (ADR-0090
# Festlegung 5). Die Tests selbst sind reine Funktionen ohne externen Dienst
# (netzlos prüfbar im Sinn von ADR-0087 Festlegung 4); ein `RUN ./gradlew
# test`, das fehlschlägt, bricht den Bau mit Exit != 0 ab. Der
# Gradle-Wrapper verifiziert die Gradle-Distribution selbst gegen
# `distributionSha256Sum` in `gradle-wrapper.properties`. Für grpc-client
# erzeugen `protoc-gen-grpc-java`/`protoc-gen-grpc-kotlin` den Stub aus der
# über den Zusatzkontext gelesenen `.proto` während dieses Builds (ADR-0090
# Festlegung 2/3, slice-103) — der Stub selbst verlässt diese Stufe nie.
FROM eclipse-temurin:21-jdk@sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e AS build
WORKDIR /src

# Der Gradle-Wrapper und die Projekt-/Modul-Manifeste zuerst — sie ändern
# sich seltener als der Quellcode und halten den Docker-Layer-Cache über
# Quelltext-Änderungen hinweg stabil (Muster: examples/csharp/Dockerfile).
COPY gradlew ./
COPY gradle/ gradle/
COPY settings.gradle.kts build.gradle.kts ./
COPY http-client/build.gradle.kts http-client/build.gradle.kts
COPY sse-client/build.gradle.kts sse-client/build.gradle.kts
COPY nats-client/build.gradle.kts nats-client/build.gradle.kts
COPY grpc-client/build.gradle.kts grpc-client/build.gradle.kts
RUN chmod +x ./gradlew && ./gradlew --no-daemon help >/dev/null

COPY http-client/src http-client/src
COPY sse-client/src sse-client/src
COPY nats-client/src nats-client/src
COPY grpc-client/src grpc-client/src

# Der zusätzliche, benannte Bau-Kontext `proto` trägt NUR proto/ (ADR-0090
# Festlegung 2) — diese COPY-Zeile ist der Punkt, an dem der Bau ohne den
# Zusatzkontext real abbricht (kein Fallback, kein übersprungener Schritt;
# dieselbe Form wie examples/csharp/Dockerfile, slice-102). Das Ziel liegt
# unter der von `com.google.protobuf`s Gradle-Plugin standardmäßig gelesenen
# Quellmenge `src/main/proto/` (grpc-client/build.gradle.kts) — anders als
# beim C#-Ziel `grpc-client/proto/`, das `Grpc.Tools`s `<Protobuf Include=…>`
# frei wählt.
COPY --from=proto cdc/stream/v1/changestream.proto grpc-client/src/main/proto/changestream.proto

RUN ./gradlew --no-daemon :http-client:test :sse-client:test :nats-client:test :grpc-client:test
RUN ./gradlew --no-daemon :http-client:installDist :sse-client:installDist :nats-client:installDist :grpc-client:installDist

# Stufe `runtime-sse`: startbares Image des SSE-Clients (ADR-0090, LP2). Ein
# eigener, benannter Zielstufen-Tag statt eines dispatchenden ENTRYPOINTs im
# `runtime`-Image — das hält den bereits dokumentierten Container-Aufruf des
# HTTP-Clients (`pg-change-feed-examples:kotlin`, unten) byte-gleich stehen.
FROM eclipse-temurin:21-jre@sha256:ca7551d4f36647207812e3c8c8596c89fd4e5226a061a5b8c8629839fe8969f6 AS runtime-sse
WORKDIR /app
COPY --from=build /src/sse-client/build/install/sse-client ./
ENTRYPOINT ["/app/bin/sse-client"]

# Stufe `runtime-nats`: startbares Image des NATS-Clients (ADR-0090, LP2,
# slice-101) — derselbe Grund für den eigenen, benannten Zielstufen-Tag wie
# bei `runtime-sse`.
FROM eclipse-temurin:21-jre@sha256:ca7551d4f36647207812e3c8c8596c89fd4e5226a061a5b8c8629839fe8969f6 AS runtime-nats
WORKDIR /app
COPY --from=build /src/nats-client/build/install/nats-client ./
ENTRYPOINT ["/app/bin/nats-client"]

# Stufe `runtime-grpc`: startbares Image des gRPC-Clients (ADR-0090, LP1,
# slice-103) — derselbe Grund für den eigenen, benannten Zielstufen-Tag wie
# bei `runtime-sse`/`runtime-nats`. Der erzeugte Protobuf-/gRPC-Stub kam nur
# über die `build`-Stufe zustande und liegt bereits kompiliert in
# `/src/grpc-client/build/install/grpc-client` — diese Stufe kopiert kein
# Quell-Erzeugnis, nur das fertig installierte Programm.
FROM eclipse-temurin:21-jre@sha256:ca7551d4f36647207812e3c8c8596c89fd4e5226a061a5b8c8629839fe8969f6 AS runtime-grpc
WORKDIR /app
COPY --from=build /src/grpc-client/build/install/grpc-client ./
ENTRYPOINT ["/app/bin/grpc-client"]

# Stufe `runtime`: startbares Image des HTTP-Clients (ADR-0087 Festlegung 3 —
# die Startform ist ein Container-Aufruf, kein Host-`gradle run`). Bleibt die
# **letzte** Stufe der Datei, damit ein `docker build` ohne `--target`
# (unverändert in `harness/mk/examples.mk`) weiterhin dieses Image liefert.
FROM eclipse-temurin:21-jre@sha256:ca7551d4f36647207812e3c8c8596c89fd4e5226a061a5b8c8629839fe8969f6 AS runtime
WORKDIR /app
COPY --from=build /src/http-client/build/install/http-client ./
ENTRYPOINT ["/app/bin/http-client"]
