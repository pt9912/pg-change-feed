# sdks/csharp/Dockerfile — Docker-only Bau-/Testsetup des C#-SDK-Pakets
# `PgChangeFeed.Client` (ADR-0106 Festlegung 2/4, slice-sdk-csharp-projektgeruest).
# Bau-Kontext ist sdks/csharp/ selbst — NICHT die Repo-Wurzel und NICHT
# examples/csharp/; beide bleiben von diesem Bau unberührt (ADR-0106
# Festlegung 2: eigenständiger Baum, kein `ProjectReference` auf einen
# privaten Baum dieses Repos).
#
# Die gRPC-Client-Fläche (slice-sdk-csharp-grpc-client-flaeche, ADR-0106
# Festlegung 1/2) liest die `.proto` NICHT aus diesem Kontext — sie liegt
# außerhalb davon, unter proto/cdc/stream/v1/ an der Repo-Wurzel. Sie kommt
# über einen ZUSÄTZLICHEN, benannten Bau-Kontext namens `proto` herein
# (`COPY --from=proto …` unten); der Aufrufer bringt ihn mit `docker build
# --build-context proto=proto …` mit (Muster examples/csharp/Dockerfile,
# harness/mk/examples.mk). Ohne ihn bricht der Bau an der COPY-Zeile ab —
# kein stiller Fallback.
#
# Keine Runtime-Stufe hier: ein SDK-Package ist keine startbare Anwendung
# (anders als examples/csharp). Die `pack`/`pack-export`-Stufen unten
# (slice-sdk-csharp-pack-werkzeug, ADR-0106 §Konsequenzen Folgepflicht 1)
# bauen auf `build` auf: `dotnet test` (unten) laeuft VOR `dotnet pack`
# (Stufe `pack`) — ein roter Test bricht den `docker build` mit Exit != 0 ab,
# bevor `pack` je erreicht wird (kein stiller Fallback).
#
# Basis-Digest real gemessen (`docker buildx imagetools inspect`, amd64/linux)
# am Zeitpunkt dieses Zuges — identisch mit dem in
# examples/csharp/Dockerfile gepinnten Kandidaten (dieselbe .NET-Generation,
# ADR-0106 §Kontext Ist-Stand):
#   mcr.microsoft.com/dotnet/sdk:10.0
#     -> sha256:60a2b2230a0d052bc54c0d453e97e331219ed503c5411470ee226859420e693c
# Pin-Hebung ist ein bewusster Commit (AGENTS.md §3.6/Modul 14), keine stille
# Auflösung.
FROM mcr.microsoft.com/dotnet/sdk:10.0@sha256:60a2b2230a0d052bc54c0d453e97e331219ed503c5411470ee226859420e693c AS build
WORKDIR /src

# Die zentrale Paket-Pinnung und die Projektdateien zuerst — sie ändern sich
# seltener als der Quellcode und halten den Docker-Layer-Cache über
# Quelltext-Änderungen hinweg stabil (analog examples/csharp/Dockerfile).
COPY Directory.Packages.props ./
COPY PgChangeFeed.Client/PgChangeFeed.Client.csproj PgChangeFeed.Client/PgChangeFeed.Client.csproj
COPY PgChangeFeed.Client/PgChangeFeed.Client.Tests/PgChangeFeed.Client.Tests.csproj PgChangeFeed.Client/PgChangeFeed.Client.Tests/PgChangeFeed.Client.Tests.csproj
RUN dotnet restore PgChangeFeed.Client/PgChangeFeed.Client.csproj && \
    dotnet restore PgChangeFeed.Client/PgChangeFeed.Client.Tests/PgChangeFeed.Client.Tests.csproj

# README.md liegt auf Sprachwurzel-Ebene (sdks/csharp/README.md) — das
# .csproj packt sie über einen relativen `../README.md`-Verweis
# (PgChangeFeed.Client.csproj); sie muss deshalb an derselben relativen
# Stelle im Bau-Kontext liegen.
COPY README.md ./
COPY PgChangeFeed.Client/ PgChangeFeed.Client/

# Der zusätzliche, benannte Bau-Kontext `proto` trägt NUR proto/ (Muster
# ADR-0090 Festlegung 2) — diese COPY-Zeile ist der Punkt, an dem der Bau
# ohne den Zusatzkontext real abbricht (kein Fallback, kein übersprungener
# Schritt; Beleg im Bericht von slice-sdk-csharp-grpc-client-flaeche). Das
# Ziel spiegelt den relativen `Protobuf Include`-Pfad in
# PgChangeFeed.Client.csproj (`Grpc/proto/changestream.proto`).
COPY --from=proto cdc/stream/v1/changestream.proto PgChangeFeed.Client/Grpc/proto/changestream.proto

RUN dotnet build PgChangeFeed.Client/PgChangeFeed.Client.csproj -c Release --no-restore
RUN dotnet test PgChangeFeed.Client/PgChangeFeed.Client.Tests/PgChangeFeed.Client.Tests.csproj -c Release --no-restore

# --- pack: `dotnet pack` ueber den bereits gebauten und getesteten Stand
# (slice-sdk-csharp-pack-werkzeug, ADR-0106 Festlegung 4). `--no-build`
# verlangt, dass die Build-Stufe oben bereits mit derselben Konfiguration
# (Release) gelaufen ist — sie ist es, als vorhergehender RUN-Schritt
# dieses selben Dockerfiles. Das Erzeugnis (PgChangeFeed.Client.0.2.1.nupkg,
# Version aus PgChangeFeed.Client.csproj) landet unter /out. ---
FROM build AS pack
RUN mkdir -p /out && \
    dotnet pack PgChangeFeed.Client/PgChangeFeed.Client.csproj -c Release --no-restore --no-build -o /out

# --- pack-export: gibt /out (das .nupkg) als tar-Stream ueber stdout aus —
# dasselbe Extraktionsmuster wie proto-export im Wurzel-Dockerfile
# (tools/harness/proto-generate.sh, ADR-0060): `tools/harness/sdk-pack-csharp.sh`
# baut diese Stufe und liest ihr ENTRYPOINT per
# `docker run --rm --network none <image> | tar -x -C sdks/csharp/dist/`
# aus — kein Bind-Mount, kein `--user`-Workaround, die extrahierte Datei
# gehoert dadurch automatisch dem aufrufenden Nutzer. ---
FROM pack AS pack-export
ENTRYPOINT ["tar", "-cf", "-", "-C", "/out", "."]

# --- integration: der Realserver-Integrationstest der SDK-Flaechen
# (slice-sdk-csharp-reale2e, Mechanik-Klasse ADR-0110 §Entscheidung
# Festlegung 2, gespiegelt vom Python-Vorbild). Baut auf `build` auf
# (dieselben Pins, kein neuer Pin), kopiert das Integrationsprojekt
# (ProjectReference auf die Client-Assembly — der Pruefling ist das SDK
# selbst, kein Wegwerf-Duplikat-Client) und baut es Release; die
# Testlaeufe fuehrt der Runner je Flaeche als `dotnet test --filter` aus —
# die Testklasse je Phase kommt als Umgebungsvariable, deren Guard
# (`:?`) blanken Aufruf laut scheitern laesst (kein gruenes Default ohne
# jeden Integrationstest, Muster sdks/python/Dockerfile). Restore braucht
# Netz (NuGet), der Testlauf selbst braucht DB-/Netz-Zugang und ist kein
# Gate. ---
FROM build AS integration
COPY PgChangeFeed.Client.Integration/ PgChangeFeed.Client.Integration/
RUN dotnet restore PgChangeFeed.Client.Integration/PgChangeFeed.Client.Integration.csproj && \
    dotnet build PgChangeFeed.Client.Integration/PgChangeFeed.Client.Integration.csproj -c Release --no-restore
WORKDIR /src
CMD ["/bin/sh", "-c", "exec dotnet test PgChangeFeed.Client.Integration/PgChangeFeed.Client.Integration.csproj -c Release --no-build --filter \"FullyQualifiedName~${PGCHANGEFEED_TEST_NAME:?Testklasse je Phase erforderlich}\" --logger \"console;verbosity=detailed\""]
