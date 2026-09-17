# harness/mk/examples.mk — Werkzeug-Fragment für die Beispiel-Client-
# Sprachwurzeln (ADR-0087, ADR-0090). Kein Gate: der Paket-Bezug (NuGet/Maven)
# braucht Netz, und `make gates` bleibt vollständig netzlos (ADR-0090
# Festlegung 5) — dieses Fragment hängt deshalb NICHT an GATE_CHECKS.
#
# `examples-csharp` baut die Werkzeugketten-Images der C#-Sprach-Wurzel
# (Bau-Kontext examples/csharp/, ADR-0087 Festlegung 3): `dotnet
# restore`/`build`/`test` für **alle vier** Programme (http-client,
# sse-client, nats-client, grpc-client seit slice-102) laufen in der
# gemeinsamen Docker-Stufe `build` (ADR-0090 Festlegung 5 — ein Ziel je
# Sprache trägt den wachsenden Umfang); ein roter Test bricht den
# `docker build` mit Exit != 0 ab, bevor der nächste Aufruf überhaupt den
# `build`-Layer-Cache erreicht — der Exit-Code jedes Aufrufs wird wie bei
# jedem anderen Ziel direkt gelesen, nie durch eine Pipe (AGENTS.md §3.9).
# Der zweite, dritte und vierte Aufruf
# (`--target runtime-sse`/`runtime-nats`/`runtime-grpc`) treffen auf den
# bereits ausgeführten `build`-Layer-Cache und bauen nur noch das jeweilige
# Runtime-Image.
#
# Jeder Aufruf trägt **zwingend** `--build-context proto=proto`: grpc-client
# liest die `.proto` über diesen zusätzlichen, benannten Bau-Kontext
# (ADR-0090 Festlegung 2, slice-102) — die `build`-Stufe des Dockerfiles
# kopiert daraus mit `COPY --from=proto …`, und zwar unabhängig vom
# angeforderten `--target`, weil jede Ziel-Stufe von `build` abhängt. Ohne
# den Zusatzkontext bricht `docker build` an dieser COPY-Zeile ab — kein
# stiller Fallback (ADR-0090 §Fitness Function).
.PHONY: examples-csharp
examples-csharp: ## C#-Sprachwurzel bauen + testen (examples/csharp, Werkzeug, kein Gate; ADR-0087/ADR-0090)
	docker build --build-context proto=proto -t pg-change-feed-examples:csharp examples/csharp
	docker build --build-context proto=proto --target runtime-sse -t pg-change-feed-examples:csharp-sse examples/csharp
	docker build --build-context proto=proto --target runtime-nats -t pg-change-feed-examples:csharp-nats examples/csharp
	docker build --build-context proto=proto --target runtime-grpc -t pg-change-feed-examples:csharp-grpc examples/csharp

# `examples-kotlin` baut die Werkzeugketten-Images der Kotlin-Sprach-Wurzel
# (Bau-Kontext examples/kotlin/, ADR-0087 Festlegung 3): der Gradle-Wrapper
# fährt `test`/`installDist` für **alle drei** Module (http-client,
# sse-client, nats-client) in der gemeinsamen Docker-Stufe `build`
# (ADR-0090 Festlegung 5); ein roter Test bricht den `docker build` mit
# Exit != 0 ab, bevor der nächste Aufruf überhaupt den `build`-Layer-Cache
# erreicht — derselbe Exit-Code-Lesepfad wie bei jedem anderen Ziel
# (AGENTS.md §3.9). Der zweite und dritte Aufruf
# (`--target runtime-sse`/`runtime-nats`) treffen auf den bereits
# ausgeführten `build`-Layer-Cache und bauen nur noch das jeweilige
# Runtime-Image.
.PHONY: examples-kotlin
examples-kotlin: ## Kotlin-Sprachwurzel bauen + testen (examples/kotlin, Werkzeug, kein Gate; ADR-0087/ADR-0090)
	docker build -t pg-change-feed-examples:kotlin examples/kotlin
	docker build --target runtime-sse -t pg-change-feed-examples:kotlin-sse examples/kotlin
	docker build --target runtime-nats -t pg-change-feed-examples:kotlin-nats examples/kotlin
