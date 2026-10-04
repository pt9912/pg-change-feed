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

# Quelle der Bibliothek 0.6.0: das benannte Verzeichnis `neu` (Artefakte von
# make sdk-pack-csharp, Modus dist; im Modus registry ein leeres Verzeichnis).
COPY --from=neu / /neu
COPY nuget-registry.config nuget-dist.config ./
COPY Gast src/Gast
COPY NullMatrix src/NullMatrix

# Drei Übersetzungen des Gast-Programms:
#   alt05  Quelltext alt gegen Package 0.5.0 (A1; Grundlage des Austauschs A2)
#   alt06  Quelltext alt gegen Package 0.6.0 (A5 Quelle; liefert die Bibliothek 0.6.0)
#   neu06  Quelltext neu (Formen mit Meldungscode) gegen Package 0.6.0 (A3)
RUN cp -r src/Gast src/Gast-alt05 && cp -r src/Gast src/Gast-alt06 && cp -r src/Gast src/Gast-neu06
RUN dotnet publish src/Gast-alt05/Gast.csproj -c Release -p:PgcfVersion=0.5.0 -p:GastQuelle=alt \
      -p:RestoreConfigFile=/kompat/nuget-registry.config -o /kompat/bin/alt05
RUN dotnet publish src/Gast-alt06/Gast.csproj -c Release -p:PgcfVersion=0.6.0 -p:GastQuelle=alt \
      -p:RestoreConfigFile=/kompat/nuget-registry.config -o /kompat/bin/alt06
RUN dotnet publish src/Gast-neu06/Gast.csproj -c Release -p:PgcfVersion=0.6.0 -p:GastQuelle=neu \
      -p:RestoreConfigFile=/kompat/nuget-registry.config -o /kompat/bin/neu06

# Modus dist: die Bibliotheksdatei für den Austausch A2 stammt aus dem Artefakt
# in /neu, über einen eigenen Paketordner (NUGET_PACKAGES), damit sie das
# veröffentlichte Package 0.6.0 der Übersetzungen oben nicht überschreibt. Die
# Übersetzungen (A3, A5) lesen immer das veröffentlichte Package.
RUN set -eu; \
    mkdir -p /kompat/lib-dist; \
    if ls /neu/*.nupkg >/dev/null 2>&1; then \
        cp -r src/Gast src/Gast-dist; \
        NUGET_PACKAGES=/kompat/pk-dist dotnet restore src/Gast-dist/Gast.csproj -p:PgcfVersion=0.6.0 \
          -p:RestoreConfigFile=/kompat/nuget-dist.config; \
        cp /kompat/pk-dist/pgchangefeed.client/0.6.0/lib/net10.0/PgChangeFeed.Client.dll /kompat/lib-dist/; \
    fi

COPY run.sh /kompat/run.sh
ENTRYPOINT ["/bin/bash", "/kompat/run.sh"]
