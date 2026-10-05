# tools/harness/sdk-kompat/kotlin/Dockerfile — Gast-Image der Kotlin-Kompatibilitäts-
# messung (make test-sdk-kompat). Die Basis ist das lokal getaggte Image der
# Stufe `build` von sdks/kotlin/Dockerfile (BASE, vom Runner übergeben): dieses
# Dockerfile trägt kein eigenes Digest-Literal; der Gradle-Wrapper stammt aus
# dieser Basis. Der Bau übersetzt das Gast-Programm und löst die Laufzeit-Jars
# auf; die Läufe der Schritte geschehen beim `docker run` (run.sh), nie in
# einer RUN-Schicht: der Schicht-Cache überspringt einen Lauf im Bau still.
ARG BASE=pg-change-feed:sdk-kompat-base-kotlin
FROM ${BASE}
WORKDIR /kompat

# Versionen, vom Runner übergeben (tools/harness/run-sdk-kompat-tests.sh): die
# neue Bibliothek (Modus dist: aus der Version-Datei des SDK) und die
# veröffentlichte, gegen die A3 und A5 übersetzen. Ohne Wert bricht der Bau ab.
ARG NEU_VERSION
ARG REGISTRY_VERSION
ENV NEU_VERSION=${NEU_VERSION} REGISTRY_VERSION=${REGISTRY_VERSION}
RUN test -n "$NEU_VERSION" && test -n "$REGISTRY_VERSION"

# Das benannte Verzeichnis `neu` trägt die Artefakte von make sdk-pack-kotlin
# (Modus dist) oder ist leer (Modus registry); run.sh liest es beim Lauf.
COPY --from=neu / /neu
RUN cp -r /src/pgchangefeed-kotlin/gradle /src/pgchangefeed-kotlin/gradlew ./ && chmod +x gradlew
COPY settings.gradle.kts build.gradle.kts ./
COPY src src

# Drei Übersetzungen des Gast-Programms mit je der aufgelösten Laufzeit-Menge:
#   alt05  Quelltext alt gegen die veröffentlichte Bibliothek 0.5.0 (A1; Grundlage des Austauschs A2)
#   altreg  Quelltext alt gegen REGISTRY_VERSION (A5; liefert die veröffentlichte Laufzeit-Menge)
#   neureg  Quelltext neu (Formen mit Meldungscode) gegen REGISTRY_VERSION (A3)
RUN set -eu; \
    bauen() { \
        ./gradlew --no-daemon -q -PpgcfVersion="$2" -PgastQuelle="$3" clean compileKotlin kopiereLaufzeit; \
        mkdir -p "bin/$1"; \
        cp -r build/classes/kotlin/main "bin/$1/klassen"; \
        cp -r build/laufzeit "bin/$1/lib"; \
    }; \
    bauen alt05 0.5.0 alt; \
    bauen altreg "$REGISTRY_VERSION" alt; \
    bauen neureg "$REGISTRY_VERSION" neu

COPY run.sh /kompat/run.sh
ENTRYPOINT ["/bin/bash", "/kompat/run.sh"]
