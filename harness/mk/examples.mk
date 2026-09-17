# harness/mk/examples.mk — Werkzeug-Fragment für die Beispiel-Client-
# Sprachwurzeln (ADR-0087, ADR-0090). Kein Gate: der Paket-Bezug (NuGet/Maven)
# braucht Netz, und `make gates` bleibt vollständig netzlos (ADR-0090
# Festlegung 5) — dieses Fragment hängt deshalb NICHT an GATE_CHECKS.
#
# `examples-csharp` baut das Werkzeugketten-Image der C#-Sprach-Wurzel
# (Bau-Kontext examples/csharp/, ADR-0087 Festlegung 3): `dotnet
# restore`/`build`/`test` laufen in der Docker-Stufe `build`; ein roter Test
# bricht den `docker build` mit Exit != 0 ab — der Exit-Code dieses Ziels
# wird wie jedes andere direkt gelesen, nie durch eine Pipe (AGENTS.md §3.9).
.PHONY: examples-csharp
examples-csharp: ## C#-Sprachwurzel bauen + testen (examples/csharp, Werkzeug, kein Gate; ADR-0087/ADR-0090)
	docker build -t pg-change-feed-examples:csharp examples/csharp
