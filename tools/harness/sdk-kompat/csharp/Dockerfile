# tools/harness/sdk-kompat/csharp/Dockerfile — Gast-Image der C#-Kompatibilitäts-
# messung (make test-sdk-kompat). Die Basis ist das lokal getaggte Image der
# Stufe `build` von sdks/csharp/Dockerfile (BASE, vom Runner übergeben): dieses
# Dockerfile trägt kein eigenes Digest-Literal. Der Bau übersetzt und
# veröffentlicht nur die Gast-Binärdateien; die Läufe der Schritte geschehen
# beim `docker run` (run.sh), nie in einer RUN-Schicht: der Schicht-Cache
# überspringt einen Lauf im Bau still.
ARG BASE=pg-change-feed:sdk-kompat-base-csharp
FROM ${BASE}
WORKDIR /kompat

# Versionen, vom Runner übergeben (tools/harness/run-sdk-kompat-tests.sh): die
# neue Bibliothek (Modus dist: aus der Version-Datei des SDK) und die
# veröffentlichte, gegen die A3 und A5 übersetzen. Ohne Wert bricht der Bau ab.
ARG NEU_VERSION
ARG REGISTRY_VERSION
ENV NEU_VERSION=${NEU_VERSION} REGISTRY_VERSION=${REGISTRY_VERSION}
RUN test -n "$NEU_VERSION" && test -n "$REGISTRY_VERSION"

# Quelle der neuen Bibliothek: das benannte Verzeichnis `neu` (Artefakte von
# make sdk-pack-csharp, Modus dist; im Modus registry ein leeres Verzeichnis).
COPY --from=neu / /neu
COPY nuget-registry.config nuget-dist.config ./
COPY Gast src/Gast
COPY NullMatrix src/NullMatrix

# Drei Übersetzungen des Gast-Programms:
#   alt05  Quelltext alt gegen Package 0.5.0 (A1; Grundlage des Austauschs A2)
#   alt06  Quelltext alt gegen Package REGISTRY_VERSION (A5 Quelle; liefert die veröffentlichte Bibliothek)
#   neu06  Quelltext neu (Formen mit Meldungscode) gegen Package REGISTRY_VERSION (A3)
RUN cp -r src/Gast src/Gast-alt05 && cp -r src/Gast src/Gast-alt06 && cp -r src/Gast src/Gast-neu06
RUN dotnet publish src/Gast-alt05/Gast.csproj -c Release -p:PgcfVersion=0.5.0 -p:GastQuelle=alt \
      -p:RestoreConfigFile=/kompat/nuget-registry.config -o /kompat/bin/alt05
RUN dotnet publish src/Gast-alt06/Gast.csproj -c Release -p:PgcfVersion="$REGISTRY_VERSION" -p:GastQuelle=alt \
      -p:RestoreConfigFile=/kompat/nuget-registry.config -o /kompat/bin/alt06
RUN dotnet publish src/Gast-neu06/Gast.csproj -c Release -p:PgcfVersion="$REGISTRY_VERSION" -p:GastQuelle=neu \
      -p:RestoreConfigFile=/kompat/nuget-registry.config -o /kompat/bin/neu06

# Modus dist: die Bibliotheksdatei für den Austausch A2 stammt aus dem Artefakt
# PgChangeFeed.Client.<NEU_VERSION>.nupkg in /neu, über einen eigenen
# Paketordner (NUGET_PACKAGES), damit sie das veröffentlichte Package der
# Übersetzungen oben nicht überschreibt. Die Übersetzungen (A3, A5) lesen immer
# das veröffentlichte Package.
RUN set -eu; \
    mkdir -p /kompat/lib-dist; \
    if [ -f "/neu/PgChangeFeed.Client.$NEU_VERSION.nupkg" ]; then \
        cp -r src/Gast src/Gast-dist; \
        NUGET_PACKAGES=/kompat/pk-dist dotnet restore src/Gast-dist/Gast.csproj -p:PgcfVersion="$NEU_VERSION" \
          -p:RestoreConfigFile=/kompat/nuget-dist.config; \
        cp "/kompat/pk-dist/pgchangefeed.client/$NEU_VERSION/lib/net10.0/PgChangeFeed.Client.dll" /kompat/lib-dist/; \
    fi

COPY run.sh /kompat/run.sh
ENTRYPOINT ["/bin/bash", "/kompat/run.sh"]
