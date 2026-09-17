# examples/kotlin/Dockerfile — Werkzeugketten-Image der Kotlin-Sprach-Wurzel
# (ADR-0087 Festlegung 3, ADR-0090 Festlegung 1/5). Bau-Kontext ist
# examples/kotlin/ selbst — NICHT die Repo-Wurzel; die Wurzel-`Dockerfile`
# und `.dockerignore` bleiben davon unberührt (ADR-0087 Festlegung 3,
# Kontext-Bindung 4).
#
# Basis-Digests real gemessen (`docker manifest inspect`, amd64/linux) am
# Zeitpunkt dieses Zuges (2026-09-17):
#   eclipse-temurin:21-jdk
#     -> sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e
#   eclipse-temurin:21-jre
#     -> sha256:ca7551d4f36647207812e3c8c8596c89fd4e5226a061a5b8c8629839fe8969f6
# Pin-Hebung ist ein bewusster Commit (AGENTS.md §3.6/Modul 14), keine stille
# Auflösung.

# Stufe `build`: `./gradlew test`/`installDist` im gepinnten JDK-Image. Der
# Paket-Bezug (Maven Central, Gradle Plugin Portal, der Gradle-Distributions-
# Download selbst) braucht Netz — deshalb bleibt `make examples-kotlin` ein
# Werkzeug, kein Gate (ADR-0090 Festlegung 5). Die Tests selbst sind reine
# Funktionen ohne externen Dienst (netzlos prüfbar im Sinn von ADR-0087
# Festlegung 4); ein `RUN ./gradlew test`, das fehlschlägt, bricht den Bau
# mit Exit != 0 ab. Der Gradle-Wrapper verifiziert die Gradle-Distribution
# selbst gegen `distributionSha256Sum` in `gradle-wrapper.properties`.
FROM eclipse-temurin:21-jdk@sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e AS build
WORKDIR /src

# Der Gradle-Wrapper und die Projekt-/Modul-Manifeste zuerst — sie ändern
# sich seltener als der Quellcode und halten den Docker-Layer-Cache über
# Quelltext-Änderungen hinweg stabil (Muster: examples/csharp/Dockerfile).
COPY gradlew ./
COPY gradle/ gradle/
COPY settings.gradle.kts build.gradle.kts ./
COPY http-client/build.gradle.kts http-client/build.gradle.kts
RUN chmod +x ./gradlew && ./gradlew --no-daemon help >/dev/null

COPY http-client/src http-client/src
RUN ./gradlew --no-daemon :http-client:test
RUN ./gradlew --no-daemon :http-client:installDist

# Stufe `runtime`: startbares Image (ADR-0087 Festlegung 3 — die Startform
# ist ein Container-Aufruf, kein Host-`gradle run`).
FROM eclipse-temurin:21-jre@sha256:ca7551d4f36647207812e3c8c8596c89fd4e5226a061a5b8c8629839fe8969f6 AS runtime
WORKDIR /app
COPY --from=build /src/http-client/build/install/http-client ./
ENTRYPOINT ["/app/bin/http-client"]
