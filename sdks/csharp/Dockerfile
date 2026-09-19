# sdks/csharp/Dockerfile — Docker-only Bau-/Testsetup des C#-SDK-Pakets
# `PgChangeFeed.Client` (ADR-0106 Festlegung 2/4, slice-sdk-csharp-projektgeruest).
# Bau-Kontext ist sdks/csharp/ selbst — NICHT die Repo-Wurzel und NICHT
# examples/csharp/; beide bleiben von diesem Bau unberührt (ADR-0106
# Festlegung 2: eigenständiger Baum, kein `ProjectReference` auf einen
# privaten Baum dieses Repos).
#
# Kein `pack`-Aufruf und keine Runtime-Stufe hier: ein SDK-Package ist keine
# startbare Anwendung (anders als examples/csharp), und `dotnet pack` bleibt
# Sache von `slice-sdk-csharp-pack-werkzeug` (ADR-0106 §Konsequenzen
# Folgepflicht 1). Dieses Dockerfile trägt ausschließlich `dotnet
# restore`/`build`/`test` — der engste heute verfügbare Sensor für dieses
# Projektgerüst, es gibt noch kein `make`-Ziel dafür.
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

RUN dotnet build PgChangeFeed.Client/PgChangeFeed.Client.csproj -c Release --no-restore
RUN dotnet test PgChangeFeed.Client/PgChangeFeed.Client.Tests/PgChangeFeed.Client.Tests.csproj -c Release --no-restore
