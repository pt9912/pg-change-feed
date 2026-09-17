# examples/csharp/Dockerfile — Werkzeugketten-Image der C#-Sprach-Wurzel
# (ADR-0087 Festlegung 3, ADR-0090 Festlegung 1/5). Bau-Kontext ist
# examples/csharp/ selbst — NICHT die Repo-Wurzel; die Wurzel-`Dockerfile`
# und `.dockerignore` bleiben davon unberührt (ADR-0087 Festlegung 3,
# Kontext-Bindung 4).
#
# Basis-Digests real gemessen (`docker manifest inspect`, amd64/linux) am
# Zeitpunkt dieses Zuges — deckungsgleich mit dem in ADR-0087 Festlegung 3
# gemessenen Kandidaten:
#   mcr.microsoft.com/dotnet/sdk:10.0
#     -> sha256:60a2b2230a0d052bc54c0d453e97e331219ed503c5411470ee226859420e693c
#   mcr.microsoft.com/dotnet/runtime:10.0
#     -> sha256:e6541e52aeab4e412012178d66a6637b65252cca2ad15ada7cd1849d7d702c3d
# Pin-Hebung ist ein bewusster Commit (AGENTS.md §3.6/Modul 14), keine stille
# Auflösung.

# Stufe `build`: `dotnet restore`/`build`/`test` im gepinnten SDK-Image, für
# **alle drei** Programme dieser Sprach-Wurzel (ADR-0090 Festlegung 1/5 — ein
# `make`-Ziel je Sprache trägt den wachsenden Umfang, hier drei Programme
# statt zweier). Der Paket-Bezug (NuGet) braucht Netz — deshalb bleibt
# `make examples-csharp` ein Werkzeug, kein Gate (ADR-0090 Festlegung 5). Die
# Tests selbst sind reine Funktionen ohne externen Dienst (netzlos prüfbar im
# Sinn von ADR-0087 Festlegung 4); ein `RUN dotnet test`, das fehlschlägt,
# bricht den Bau mit Exit != 0 ab.
FROM mcr.microsoft.com/dotnet/sdk:10.0@sha256:60a2b2230a0d052bc54c0d453e97e331219ed503c5411470ee226859420e693c AS build
WORKDIR /src

# Die zentrale Paket-Pinnung und die Projektdateien zuerst — sie ändern sich
# seltener als der Quellcode und halten den Docker-Layer-Cache über
# Quelltext-Änderungen hinweg stabil.
COPY Directory.Packages.props ./
COPY http-client/http-client.csproj http-client/http-client.csproj
COPY http-client/HttpClient.Tests/HttpClient.Tests.csproj http-client/HttpClient.Tests/HttpClient.Tests.csproj
COPY sse-client/sse-client.csproj sse-client/sse-client.csproj
COPY sse-client/SseClient.Tests/SseClient.Tests.csproj sse-client/SseClient.Tests/SseClient.Tests.csproj
COPY nats-client/nats-client.csproj nats-client/nats-client.csproj
COPY nats-client/NatsClient.Tests/NatsClient.Tests.csproj nats-client/NatsClient.Tests/NatsClient.Tests.csproj
RUN dotnet restore http-client/http-client.csproj && \
    dotnet restore http-client/HttpClient.Tests/HttpClient.Tests.csproj && \
    dotnet restore sse-client/sse-client.csproj && \
    dotnet restore sse-client/SseClient.Tests/SseClient.Tests.csproj && \
    dotnet restore nats-client/nats-client.csproj && \
    dotnet restore nats-client/NatsClient.Tests/NatsClient.Tests.csproj

COPY http-client/ http-client/
COPY sse-client/ sse-client/
COPY nats-client/ nats-client/
RUN dotnet build http-client/http-client.csproj -c Release --no-restore
RUN dotnet test http-client/HttpClient.Tests/HttpClient.Tests.csproj -c Release --no-restore
RUN dotnet publish http-client/http-client.csproj -c Release --no-restore -o /app/publish-http
RUN dotnet build sse-client/sse-client.csproj -c Release --no-restore
RUN dotnet test sse-client/SseClient.Tests/SseClient.Tests.csproj -c Release --no-restore
RUN dotnet publish sse-client/sse-client.csproj -c Release --no-restore -o /app/publish-sse
RUN dotnet build nats-client/nats-client.csproj -c Release --no-restore
RUN dotnet test nats-client/NatsClient.Tests/NatsClient.Tests.csproj -c Release --no-restore
RUN dotnet publish nats-client/nats-client.csproj -c Release --no-restore -o /app/publish-nats

# Stufe `runtime-sse`: startbares Image des SSE-Clients (ADR-0090, LP1). Ein
# eigener, benannter Zielstufen-Tag statt eines dispatchenden ENTRYPOINTs im
# `runtime`-Image — das hält den bereits dokumentierten Container-Aufruf des
# HTTP-Clients (`pg-change-feed-examples:csharp`, unten) byte-gleich stehen.
FROM mcr.microsoft.com/dotnet/runtime:10.0@sha256:e6541e52aeab4e412012178d66a6637b65252cca2ad15ada7cd1849d7d702c3d AS runtime-sse
WORKDIR /app
COPY --from=build /app/publish-sse ./
ENTRYPOINT ["dotnet", "sse-client.dll"]

# Stufe `runtime-nats`: startbares Image des NATS-Clients (ADR-0090, LP1,
# slice-101) — derselbe Grund für den eigenen, benannten Zielstufen-Tag wie
# bei `runtime-sse`.
FROM mcr.microsoft.com/dotnet/runtime:10.0@sha256:e6541e52aeab4e412012178d66a6637b65252cca2ad15ada7cd1849d7d702c3d AS runtime-nats
WORKDIR /app
COPY --from=build /app/publish-nats ./
ENTRYPOINT ["dotnet", "nats-client.dll"]

# Stufe `runtime`: startbares Image des HTTP-Clients (ADR-0087 Festlegung 3 —
# die Startform ist ein Container-Aufruf, kein Host-`dotnet run`). Bleibt die
# **letzte** Stufe der Datei, damit ein `docker build` ohne `--target`
# (unverändert in `harness/mk/examples.mk`) weiterhin dieses Image liefert.
FROM mcr.microsoft.com/dotnet/runtime:10.0@sha256:e6541e52aeab4e412012178d66a6637b65252cca2ad15ada7cd1849d7d702c3d AS runtime
WORKDIR /app
COPY --from=build /app/publish-http ./
ENTRYPOINT ["dotnet", "http-client.dll"]
