# harness/mk/sdk.mk — Werkzeug-Fragment für das erste SDK-Package
# (`PgChangeFeed.Client`, C#/NuGet, ADR-0106). Kein Gate: `dotnet restore`
# braucht Netz, `make gates` bleibt netzlos (ADR-0106 Festlegung 4, dieselbe
# Begründung wie harness/mk/examples.mk) — dieses Fragment hängt deshalb
# NICHT an GATE_CHECKS.
#
# `sdk-pack-csharp` baut/testet/paketiert das C#-SDK Docker-only im
# gepinnten mcr.microsoft.com/dotnet/sdk-Image (sdks/csharp/Dockerfile,
# Stufe `pack-export`): `dotnet test` gegen BEIDE Testflächen (HTTP + gRPC,
# dasselbe Testprojekt PgChangeFeed.Client.Tests) läuft VOR `dotnet pack`
# in derselben Docker-Bau-Kette — ein roter Test bricht den `docker build`
# mit Exit != 0 ab, bevor die `pack`-Stufe je erreicht wird (kein stiller
# Fallback, Muster harness/mk/examples.mk). Wie beim gRPC-Bau der Fläche
# trägt der Aufruf zwingend `--build-context proto=proto` (ADR-0090
# Festlegung 2, übernommen aus slice-sdk-csharp-grpc-client-flaeche) — ohne
# ihn bricht der Bau an der `COPY --from=proto`-Zeile in
# sdks/csharp/Dockerfile ab.
#
# Export analog `make proto-generate` (tools/harness/proto-generate.sh):
# tools/harness/sdk-pack-csharp.sh extrahiert das .nupkg host-seitig aus der
# `pack-export`-Stufe (`docker run --rm --network none <image> | tar -x`,
# `set -o pipefail` unter bash, AGENTS.md §3.9) nach sdks/csharp/dist/
# (`.gitignore`t). Erzeugnis: PgChangeFeed.Client.0.1.0.nupkg. Exit-Code des
# Skripts wird wie bei jedem anderen Ziel direkt gelesen.
.PHONY: sdk-pack-csharp
sdk-pack-csharp: ## C#-SDK bauen+testen+paketieren (sdks/csharp, .nupkg nach sdks/csharp/dist/; Werkzeug, kein Gate; ADR-0106)
	@bash tools/harness/sdk-pack-csharp.sh
