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
# fährt `test`/`installDist` für **alle vier** Module (http-client,
# sse-client, nats-client, grpc-client seit slice-103) in der gemeinsamen
# Docker-Stufe `build` (ADR-0090 Festlegung 5); ein roter Test bricht den
# `docker build` mit Exit != 0 ab, bevor der nächste Aufruf überhaupt den
# `build`-Layer-Cache erreicht — derselbe Exit-Code-Lesepfad wie bei jedem
# anderen Ziel (AGENTS.md §3.9). Der zweite, dritte und vierte Aufruf
# (`--target runtime-sse`/`runtime-nats`/`runtime-grpc`) treffen auf den
# bereits ausgeführten `build`-Layer-Cache und bauen nur noch das jeweilige
# Runtime-Image.
#
# Der vierte Aufruf trägt **zwingend** `--build-context proto=proto`:
# grpc-client liest die `.proto` über diesen zusätzlichen, benannten
# Bau-Kontext (ADR-0090 Festlegung 2, übertragen aus slice-102, hier
# slice-103) — die `build`-Stufe des Dockerfiles kopiert daraus mit
# `COPY --from=proto …`, und zwar unabhängig vom angeforderten `--target`,
# weil jede Ziel-Stufe von `build` abhängt (dieselbe Struktur wie bei
# `examples-csharp` oben). Ohne den Zusatzkontext bricht `docker build` an
# dieser COPY-Zeile ab — kein stiller Fallback (ADR-0090 §Fitness Function).
.PHONY: examples-kotlin
examples-kotlin: ## Kotlin-Sprachwurzel bauen + testen (examples/kotlin, Werkzeug, kein Gate; ADR-0087/ADR-0090)
	docker build --build-context proto=proto -t pg-change-feed-examples:kotlin examples/kotlin
	docker build --build-context proto=proto --target runtime-sse -t pg-change-feed-examples:kotlin-sse examples/kotlin
	docker build --build-context proto=proto --target runtime-nats -t pg-change-feed-examples:kotlin-nats examples/kotlin
	docker build --build-context proto=proto --target runtime-grpc -t pg-change-feed-examples:kotlin-grpc examples/kotlin

# `example-run-go` startet real ein Go-Beispiel gegen die Demo-Umgebung
# (ADR-0098 Festlegung 1/2, Supersedes ADR-0076 Festlegung 1/Startform-Bullet
# in genau dieser Klausel): baut (falls nötig — Docker-Layer-Cache greift bei
# unverändertem `examples/Dockerfile`-Kontext) das je Oberfläche passende
# Image aus `examples/Dockerfile` (Wurzel-Bau-Kontext, isoliert über
# `examples/Dockerfile.dockerignore`) und startet es real per
# `docker run --rm --network cdc-examples --env-file examples/.env`. `SURFACE=`
# ist ein Pflicht-Argument (`http`, `sse`, `grpc`, `nats` oder `nats-stream`,
# der vierte Zugriffsweg seit ADR-0100); ein fehlendes oder unbekanntes
# `SURFACE` bricht mit `$(error …)` ab, BEVOR ein `docker
# build`/`docker run` versucht wird (kein halb gestarteter Zustand). `ARGS=`
# trägt die Flag-Übersteuerung, die ADR-0076 Festlegung 1 bereits für alle
# Beispiele vorsieht (z. B. `ARGS="-source demo -publication demo_pub"`).
# `SURFACE=http` liefert das `runtime`-Image ohne `--target` (mirror der
# ADR-0087-Konvention); die anderen vier Oberflächen adressieren ihre
# `runtime-<surface>`-Stufe explizit. Das Docker-Netzwerk `cdc-examples` und
# `examples/.env` legt `slice-beispiele-compose-bootstrap` an — dieses Ziel
# referenziert beide nur, ohne sie zu erzeugen; ein Aufruf ohne sie schlägt
# real und sichtbar am `docker run` fehl (kein stiller Fallback). Exit-Code
# jedes Aufrufs wird direkt gelesen, wie bei jedem anderen Ziel (AGENTS.md
# §3.9); kein Host-`go build`/`go run` (`AGENTS.md` §3.1) — auch der
# Go-Startweg läuft über ein Image. Kein Gate (Werkzeug, wie
# `examples-csharp`/`examples-kotlin`): der Start braucht das benannte
# Docker-Netzwerk, `make gates` bleibt netzlos.
.PHONY: example-run-go
example-run-go: ## Go-Beispiel bauen+starten (Pflicht: SURFACE=http|sse|grpc|nats|nats-stream, optional ARGS=…; Werkzeug, kein Gate; ADR-0098)
ifeq ($(filter $(SURFACE),http sse grpc nats nats-stream),)
	$(error SURFACE muss http, sse, grpc, nats oder nats-stream sein, z.B. make example-run-go SURFACE=http)
endif
	docker build -f examples/Dockerfile $(if $(filter $(SURFACE),http),,--target runtime-$(SURFACE)) -t pg-change-feed-examples:go$(if $(filter $(SURFACE),http),,-$(SURFACE)) .
	docker run --rm --network cdc-examples --env-file examples/.env pg-change-feed-examples:go$(if $(filter $(SURFACE),http),,-$(SURFACE)) $(ARGS)

# `example-run-csharp`/`example-run-kotlin` starten real eines der vier
# bereits gebauten C#-/Kotlin-Beispiele gegen die Demo-Umgebung (ADR-0098
# Festlegung 2). Anders als `example-run-go` bauen sie **nichts**: der Bau
# liegt bei `make examples-csharp`/`make examples-kotlin` (ADR-0087
# Festlegung 3, ADR-0090) — dieses Ziel startet nur den bereits vorhandenen
# Image-Tag `pg-change-feed-examples:csharp[-<surface>]`/`:kotlin[-<surface>]`.
# `SURFACE=` ist wie bei `example-run-go` ein Pflicht-Argument (`http`,
# `sse`, `grpc` oder `nats`); ein fehlendes oder unbekanntes `SURFACE` bricht
# mit `$(error …)` ab, BEVOR ein `docker run` versucht wird. `ARGS=` trägt
# dieselbe Flag-Übersteuerung wie bei den anderen beiden Sprachen (ADR-0076
# Festlegung 1). Läuft der Bau-Schritt noch nicht (kein `make
# examples-csharp`/`make examples-kotlin` zuvor), scheitert `docker run` real
# und sichtbar mit einem Docker-eigenen "image not found" — kein stiller
# Vorab-Bau und keine stille Fallback-Meldung dieses Ziels selbst; das
# Docker-Netzwerk `cdc-examples` und `examples/.env` legt
# `slice-beispiele-compose-bootstrap` an, dieses Ziel referenziert beide nur,
# ohne sie zu erzeugen. Exit-Code jedes Aufrufs wird direkt gelesen, wie bei
# jedem anderen Ziel (AGENTS.md §3.9). Kein Gate (Werkzeug, wie
# `examples-csharp`/`examples-kotlin`/`example-run-go`): der Start braucht
# das benannte Docker-Netzwerk, `make gates` bleibt netzlos.
.PHONY: example-run-csharp
example-run-csharp: ## C#-Beispiel starten (Pflicht: SURFACE=http|sse|grpc|nats, optional ARGS=…; braucht vorherigen make examples-csharp; Werkzeug, kein Gate; ADR-0098)
ifeq ($(filter $(SURFACE),http sse grpc nats),)
	$(error SURFACE muss http, sse, grpc oder nats sein, z.B. make example-run-csharp SURFACE=http)
endif
	docker run --rm --network cdc-examples --env-file examples/.env pg-change-feed-examples:csharp$(if $(filter $(SURFACE),http),,-$(SURFACE)) $(ARGS)

.PHONY: example-run-kotlin
example-run-kotlin: ## Kotlin-Beispiel starten (Pflicht: SURFACE=http|sse|grpc|nats, optional ARGS=…; braucht vorherigen make examples-kotlin; Werkzeug, kein Gate; ADR-0098)
ifeq ($(filter $(SURFACE),http sse grpc nats),)
	$(error SURFACE muss http, sse, grpc oder nats sein, z.B. make example-run-kotlin SURFACE=http)
endif
	docker run --rm --network cdc-examples --env-file examples/.env pg-change-feed-examples:kotlin$(if $(filter $(SURFACE),http),,-$(SURFACE)) $(ARGS)

# `example-demo-up`/`example-demo-down` kapseln die Demo-/Quickstart-
# Umgebung unter examples/ (LH-QA-OPS-001, ADR-0098 Festlegung 3/4,
# slice-beispiele-compose-bootstrap): ein Nutzer fährt `postgres`+`nats`+
# den Feed-Container hoch und bekommt ohne manuellen SQL-/`make
# schema-rollout`-Zwischenschritt reale, sofort lesbare Demo-Daten
# (`public.orders`, über `CDC_TABLES` beim Feed-Start automatisch aktiviert,
# ADR-0028). Der eigentliche Ablauf steht in examples/bootstrap.sh
# (Schema-Rollout über d-migrate, Beispiel-Quelle/-Tabelle-Registrierung,
# Health-Poll) — dieses Ziel ruft es nur auf; Exit-Code direkt gelesen wie
# bei jedem anderen Ziel (AGENTS.md §3.9). Kein Gate: die Demo-Umgebung
# braucht DB-Zugang und das benannte Docker-Netzwerk `cdc-examples`,
# `make gates` bleibt netzlos — dasselbe Werkzeug-statt-Gate-Argument wie
# `examples-csharp`/`examples-kotlin`/`example-run-go`.
.PHONY: example-demo-up
example-demo-up: ## Demo-Umgebung hochfahren + bootstrappen (postgres+nats+feed, reale Demo-Daten sofort lesbar; Werkzeug, kein Gate; ADR-0098)
	@bash examples/bootstrap.sh

# Symmetrischer Abbau: `-v --remove-orphans` wie
# tools/harness/run-integration-tests.sh — Demo-Daten leben ausschließlich
# im Container, ein zweiter `example-demo-up`-Lauf beginnt wieder frisch.
.PHONY: example-demo-down
example-demo-down: ## Demo-Umgebung abräumen (Container+Netzwerk; Werkzeug, kein Gate; ADR-0098)
	docker compose -f examples/compose.yaml down -v --remove-orphans
