# harness/mk/examples.mk — Werkzeug-Fragment für die Beispiel-Client-
# Sprachwurzeln (ADR-0087, ADR-0090). Kein Gate: der Paket-Bezug (NuGet/Maven)
# braucht Netz, und `make gates` bleibt vollständig netzlos (ADR-0090
# Festlegung 5) — dieses Fragment hängt deshalb NICHT an GATE_CHECKS.
#
# `examples-csharp` baut die Werkzeugketten-Images der C#-Sprach-Wurzel
# (Bau-Kontext examples/csharp/, ADR-0087 Festlegung 3): `dotnet
# restore`/`build`/`test` für **beide** Programme (http-client, sse-client)
# laufen in der gemeinsamen Docker-Stufe `build` (ADR-0090 Festlegung 5 — ein
# Ziel je Sprache trägt den wachsenden Umfang); ein roter Test bricht den
# `docker build` mit Exit != 0 ab, bevor der zweite Aufruf überhaupt den
# `build`-Layer-Cache erreicht — der Exit-Code jedes Aufrufs wird wie bei
# jedem anderen Ziel direkt gelesen, nie durch eine Pipe (AGENTS.md §3.9). Der
# zweite Aufruf (`--target runtime-sse`) trifft auf den bereits ausgeführten
# `build`-Layer-Cache und baut nur noch das zweite Runtime-Image.
.PHONY: examples-csharp
examples-csharp: ## C#-Sprachwurzel bauen + testen (examples/csharp, Werkzeug, kein Gate; ADR-0087/ADR-0090)
	docker build -t pg-change-feed-examples:csharp examples/csharp
	docker build --target runtime-sse -t pg-change-feed-examples:csharp-sse examples/csharp

# `examples-kotlin` baut das Werkzeugketten-Image der Kotlin-Sprach-Wurzel
# (Bau-Kontext examples/kotlin/, ADR-0087 Festlegung 3): der Gradle-Wrapper
# fährt `test`/`installDist` in der Docker-Stufe `build`; ein roter Test
# bricht den `docker build` mit Exit != 0 ab — derselbe Exit-Code-Lesepfad
# wie bei jedem anderen Ziel (AGENTS.md §3.9).
.PHONY: examples-kotlin
examples-kotlin: ## Kotlin-Sprachwurzel bauen + testen (examples/kotlin, Werkzeug, kein Gate; ADR-0087/ADR-0090)
	docker build -t pg-change-feed-examples:kotlin examples/kotlin
