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

# Das benannte Verzeichnis `neu` trägt die Artefakte von make sdk-pack-kotlin
# (Modus dist) oder ist leer (Modus registry); run.sh liest es beim Lauf.
COPY --from=neu / /neu
RUN cp -r /src/pgchangefeed-kotlin/gradle /src/pgchangefeed-kotlin/gradlew ./ && chmod +x gradlew
COPY settings.gradle.kts build.gradle.kts ./
COPY src src

# Drei Übersetzungen des Gast-Programms mit je der aufgelösten Laufzeit-Menge:
#   alt05  Quelltext alt gegen die veröffentlichte Bibliothek 0.5.0 (A1; Grundlage des Austauschs A2)
#   alt06  Quelltext alt gegen 0.6.0 (A5; liefert die Laufzeit-Menge 0.6.0)
#   neu06  Quelltext neu (Formen mit Meldungscode) gegen 0.6.0 (A3)
RUN set -eu; \
    bauen() { \
        ./gradlew --no-daemon -q -PpgcfVersion="$2" -PgastQuelle="$3" clean compileKotlin kopiereLaufzeit; \
        mkdir -p "bin/$1"; \
        cp -r build/classes/kotlin/main "bin/$1/klassen"; \
        cp -r build/laufzeit "bin/$1/lib"; \
    }; \
    bauen alt05 0.5.0 alt; \
    bauen alt06 0.6.0 alt; \
    bauen neu06 0.6.0 neu

COPY run.sh /kompat/run.sh
ENTRYPOINT ["/bin/bash", "/kompat/run.sh"]
