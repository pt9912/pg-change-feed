# Makefile — generiert von ai-harness-init (Aggregator). Die Gate-Belange leben als
# Fragmente unter harness/mk/*.mk; jedes haengt seine Checks an GATE_CHECKS. Der
# Gate-Nachweis (record-gates) laeuft strikt ZULETZT via Ordnungskante auf GATE_CHECKS
# — waehrend make -j die Checks parallelisiert; .NOTPARALLEL ist bewusst NICHT gewaehlt
# (das serialisierte das ganze Makefile). Sprach-agnostisch: ohne --lang matchen nur
# baseline/doc-gate/enforce, mit --lang zusaetzlich das Code-Gate-Fragment.
GATE_CHECKS :=

# .PHONY trägt die Targets, deren Name ein Baum-Verzeichnis schattieren kann
# (`test:` gegen das Verzeichnis `test/` — make meldet "bereits aktuell",
# ohne das Rezept zu fahren).
.PHONY: gates help mod-download test test-store test-replication test-integration test-notify image image-stale

# Gate-Fragmente je Belang (baseline/doc-gate/enforce + Sprach-Code-Gates) einbinden.
# Alphabetisch (baseline < doc-gate < enforce < <lang>); die Ordnungskante unten steht
# NACH dem Include und sieht GATE_CHECKS damit vollstaendig.
include harness/mk/*.mk

# Architektur-Gate (a-check) — eingebunden seit dem ersten Binary im Baum
# (cmd/pg-change-feed); das Fragment trägt den gepinnten Release-Digest,
# Pin-Hebung = bewusster Commit (a-check.mk). a-check hängt an GATE_CHECKS
# (ADR-0041) und läuft damit im `make gates`-Bündel mit.
include a-check.mk

.PHONY: image
# Mit VERSION=<semver> (ADR-0051): ein Build, Push nach GHCR UND Docker Hub
# (Content-Mirror statt zweitem Build); LATEST=true setzt zusaetzlich
# :latest auf beiden Registries (nur fuer stabile, nicht-Prerelease Tags —
# das entscheidet der Aufrufer, release.yml). Ohne VERSION unveraendertes
# Verhalten: lokal geladen, nur :dev (ADR-0044/ADR-0103).
#
# --platform linux/amd64,linux/arm64 nur im VERSION=-Zweig (Multi-Arch,
# ADR-0051 additiv): der VERSION-lose :dev-Zweig (--load) bleibt bewusst
# einplattformig fuer eine konsistente, schnelle lokale Dev-Iteration ueber
# unterschiedliche Docker-Setups hinweg — ob --load ueberhaupt eine
# Multi-Platform-Manifestliste laden kann, haengt vom Storage-Treiber ab
# (real bestaetigt: mit aktiviertem containerd-Image-Store gelingt es,
# es ist keine grundsaetzliche buildx-Grenze), und dieses Verhalten bei
# jedem Entwicklerrechner vorauszusetzen waere keine verlaessliche Basis.
# Die bestehende Digest-Extraktion (grep gegen containerimage.digest)
# braucht fuer den VERSION=-Zweig keine Aenderung: bei Multi-Platform-Builds
# liefert --metadata-file exakt einen containerimage.digest-Eintrag, den
# Index-/Manifestlisten-Digest (real gegen eine lokale Test-Registry mit
# --platform linux/amd64,linux/arm64 geprueft: grep-Ergebnis == sha256 des
# von der Registry abgerufenen rohen Index-Manifests).
ifdef VERSION
image: ## Baut und pusht das Multi-Arch-OCI-Image (linux/amd64+linux/arm64) nach GHCR+Docker Hub (VERSION=<semver>, optional LATEST=true — ADR-0051)
	docker buildx build --push --platform linux/amd64,linux/arm64 \
	  --build-arg VERSION=$(VERSION) \
	  -t ghcr.io/pt9912/pg-change-feed:$(VERSION) \
	  -t docker.io/pt9912/pg-change-feed:$(VERSION) \
	  $(if $(filter true,$(LATEST)),-t ghcr.io/pt9912/pg-change-feed:latest -t docker.io/pt9912/pg-change-feed:latest,) \
	  --metadata-file harness/image-hash.raw . && grep -o '"containerimage.digest":[[:space:]]*"sha256:[0-9a-f]*' harness/image-hash.raw | head -1 | grep -o 'sha256:[0-9a-f]*' > harness/image-hash.txt && rm harness/image-hash.raw
else
image: ## Baut das OCI-Image (lokale Host-Plattform); Image-Hash nach harness/image-hash.txt (lokal, nicht committet — ADR-0103); VERSION=<semver> fuer den Multi-Arch-Multi-Registry-Push (ADR-0051)
	docker buildx build --load --metadata-file harness/image-hash.raw -t ghcr.io/pt9912/pg-change-feed:dev . && grep -o '"containerimage.digest":[[:space:]]*"sha256:[0-9a-f]*' harness/image-hash.raw | head -1 | grep -o 'sha256:[0-9a-f]*' > harness/image-hash.txt && rm harness/image-hash.raw
endif

image-stale: ## Advisory: FROM-Digests gegen Registry-Digests (Modul 14, braucht Netz)
	@bash tools/harness/image-stale.sh

# --- Upstream-Pin-Freshness P3-P9 (kein Gate; ADR-0051 Entscheidung 7, braucht Netz) ---
# P1/P2 deckt das bestehende image-stale (Dockerfile-FROM-Zeilen); P3-P6
# sind Makefile-/mk-Variablen mit demselben Digest-Pin-Muster, ueber das
# gemeinsame tools/harness/pin-stale.sh; P7 (zwei Achsen), P8 und P9
# tragen wegen abweichender Pin-Form (Kurs-Baseline-Version, GitHub-
# Action-SHAs) je ein eigenes Skript.
.PHONY: pin-stale-race pin-stale-pgtest pin-stale-dmigrate pin-stale-acheck pin-stale-dcheck pin-stale-baseline pin-stale-actions
pin-stale-race: ## Advisory P3: TOOLCHAIN_RACE_IMAGE gegen Registry-Digest (ADR-0051, braucht Netz)
	@bash tools/harness/pin-stale.sh Makefile TOOLCHAIN_RACE_IMAGE

pin-stale-pgtest: ## Advisory P4: PG_TEST_IMAGE gegen Registry-Digest (ADR-0051, braucht Netz)
	@bash tools/harness/pin-stale.sh Makefile PG_TEST_IMAGE

pin-stale-dmigrate: ## Advisory P5: D_MIGRATE_IMAGE gegen den aktuellen :latest-Digest (ADR-0051, braucht Netz)
	@bash tools/harness/pin-stale.sh Makefile D_MIGRATE_IMAGE ghcr.io/pt9912/d-migrate:latest

pin-stale-acheck: ## Advisory P6: A_CHECK_IMAGE gegen den aktuellen :latest-Digest (ADR-0051, braucht Netz)
	@bash tools/harness/pin-stale.sh a-check.mk A_CHECK_IMAGE ghcr.io/pt9912/a-check:latest

pin-stale-dcheck: ## Advisory P7: DCHECK_IMAGE/DCHECK_DIGEST — Tag-Frische UND Digest-Drift (ADR-0051, braucht Netz)
	@bash tools/harness/pin-stale-dcheck.sh

pin-stale-baseline: ## Advisory P8: adoptierte Kurs-Baseline-Version gegen den neuesten Kurs-Release (ADR-0051, braucht Netz)
	@bash tools/harness/pin-stale-baseline.sh

pin-stale-actions: ## Advisory P9: alle uses:-SHA-Pins ueber .github/workflows/*.yml — Tag-Mutation UND Tag-Frische (ADR-0051, braucht Netz)
	@bash tools/harness/pin-stale-actions.sh

.PHONY: image-cve
# TRIVY_IMAGE traegt aquasec/trivy v0.74.0 (Digest-Pin, Modul 14).
TRIVY_IMAGE ?= aquasec/trivy@sha256:62b1e65e8869bc4b4c6aa4fa2b21595256c7c2f6018a9d9ad61caf87187c1969
# GHCR-Pakete, die per GITHUB_TOKEN gepusht werden, entstehen unabhaengig
# von der Sichtbarkeit des Repos privat (GitHub-Verhalten) — GHCR_USERNAME/
# GHCR_PASSWORD (optional) authentifizieren ausschliesslich den
# gescannten ghcr.io-Pull ueber Trivys eigene --username/TRIVY_PASSWORD-
# Mechanik, getrennt von Trivys eigenem, unauthentifiziertem Bezug seiner
# Vulnerability-DB: ein allgemeines Docker-Credential-Mount (`~/.docker/
# config.json`) faerbte stattdessen jeden Registry-Zugriff des Containers
# ein, einschliesslich des DB-Bezugs, sobald die Config einen
# Credential-Helper referenziert, den der Container nicht ausfuehren kann.
image-cve: ## Advisory: Trivy CRITICAL/HIGH gegen das publizierte GHCR-:latest-Image (ADR-0051, kein Gate, braucht Netz; optional GHCR_USERNAME/GHCR_PASSWORD für ein privates Paket)
	docker run --rm $(if $(GHCR_PASSWORD),-e TRIVY_PASSWORD="$(GHCR_PASSWORD)",) $(TRIVY_IMAGE) image --image-src remote $(if $(GHCR_USERNAME),--username "$(GHCR_USERNAME)",) --severity CRITICAL,HIGH --exit-code 1 ghcr.io/pt9912/pg-change-feed:latest

.PHONY: test-release-tag-info
test-release-tag-info: ## Tabellentest gegen tools/harness/release-tag-info.sh (SemVer-2.0-Validierung, ADR-0051, netzlos)
	@bash tools/harness/run-release-tag-info-tests.sh

.PHONY: suchlauf-nachmessen
suchlauf-nachmessen: ## Misst die suchlauf-Blöcke eines Slice-Plans nach: make suchlauf-nachmessen PLAN=<Datei> (netzlos, kein Gate; harness/sensors/suchlauf-nachmessen.md)
	$(if $(PLAN),,$(error PLAN fehlt, z.B. make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-x.md))
	@bash tools/harness/suchlauf-nachmessen.sh "$(PLAN)"

.PHONY: test-suchlauf-nachmessen
test-suchlauf-nachmessen: ## Tabellentest gegen tools/harness/suchlauf-nachmessen.sh (elf Fälle gegen ein Wegwerf-Repo, netzlos)
	@bash tools/harness/run-suchlauf-nachmessen-tests.sh

.PHONY: kommentar-kennungen
kommentar-kennungen: ## Listet Go-Kommentarblöcke mit mehr als einer Kennung oder "ff." (Kandidaten): make kommentar-kennungen [PATHS=<Pfade>] [COUNT=1] [TESTS=exclude|only] [DIFF=<Basis>] (Docker-only, netzlos, kein Gate; harness/sensors/kommentar-kennungen.md)
	@TOOLCHAIN_IMAGE="$(TOOLCHAIN_IMAGE)" bash tools/harness/kommentar-kennungen.sh "$(PATHS)" "$(COUNT)" "$(TESTS)" "$(DIFF)"

.PHONY: test-kommentar-kennungen
test-kommentar-kennungen: ## Tabellentests des Programms tools/harness/kommentar-kennungen (Go-Test, läuft auch unter make test) und seines Aufrufers (bash, Stub-docker plus vier Docker-Läufe; Docker-only, netzlos)
	docker run --rm --network none -v "$(CURDIR)":/src:ro -w /src -e GOCACHE=/tmp/gocache $(TOOLCHAIN_IMAGE) go test ./tools/harness/kommentar-kennungen/
	@TOOLCHAIN_IMAGE="$(TOOLCHAIN_IMAGE)" bash tools/harness/run-kommentar-kennungen-tests.sh

.PHONY: test-rollout-restore
test-rollout-restore: ## Tabellentest gegen tools/schema/rollout-restore.sh (Rücknahme von plan.yaml/down.sql, Aufrufer-Prüfung, netzlos)
	@bash tools/harness/run-rollout-restore-tests.sh

.PHONY: test-sdk-csharp-release-tag-info
test-sdk-csharp-release-tag-info: ## Tabellentest gegen tools/harness/sdk-csharp-release-tag-info.sh (SemVer-2.0-Validierung, ADR-0106, netzlos)
	@bash tools/harness/run-sdk-csharp-release-tag-info-tests.sh

.PHONY: test-sdk-python-release-tag-info
test-sdk-python-release-tag-info: ## Tabellentest gegen tools/harness/sdk-python-release-tag-info.sh (PEP-440-Einfachfall-Validierung, ADR-0107, netzlos)
	@bash tools/harness/run-sdk-python-release-tag-info-tests.sh

.PHONY: test-sdk-kotlin-release-tag-info
test-sdk-kotlin-release-tag-info: ## Tabellentest gegen tools/harness/sdk-kotlin-release-tag-info.sh (SemVer-2.0-Validierung, ADR-0109, netzlos)
	@bash tools/harness/run-sdk-kotlin-release-tag-info-tests.sh

.PHONY: test-image-stale-parse
test-image-stale-parse: ## Tabellentest gegen tools/harness/dockerfile-from.sh (FROM-Parsing von image-stale, ADR-0039, netzlos)
	@bash tools/harness/run-dockerfile-from-tests.sh

.PHONY: test-dockerhub-token
test-dockerhub-token: ## Tabellentest gegen tools/harness/dockerhub-token.sh (Docker-Hub-Login-Antwort-Parsing, ADR-0051, netzlos)
	@bash tools/harness/run-dockerhub-token-tests.sh

.PHONY: test-version-injection
test-version-injection: ## Regressionstest: VERSION-Injektion vom Dockerfile bis --version (ADR-0051, netzlos sofern Basis-Images bereits gecacht)
	@bash tools/harness/run-version-injection-test.sh

.PHONY: doc-ci-matrix
doc-ci-matrix: ## LH-QA-POR-001/002-Beleg: reale GitHub-Actions-Läufe abfragen, docs/user/ci-matrix-abdeckung.md schreiben (ADR-0105, kein Gate, braucht Netz)
	@bash tools/harness/ci-matrix-abdeckung.sh

# --- Tests (kein Gate; Docker-only, gepinnte Images) ---
# Toolchain-Container = derselbe gepinnte Digest wie im Dockerfile; der
# PostgreSQL-Testcontainer trägt seinen Digest aus `docker manifest inspect
# postgres:18-alpine` (amd64). Caches leben in Docker-Volumes, Daten im
# Container — nichts davon im Arbeitsbaum. PG_TEST_IMAGE bleibt hier auf
# PostgreSQL 18 für lokale/manuelle Läufe (ADR-0058 Entscheidung 4); die
# CI-Versionsmatrix (LH-QA-POR-001, .github/workflows/e2e.yml) überschreibt
# dieselbe Variable je Leg, die compose.yaml per `${PG_TEST_IMAGE}`-
# Interpolation liest — kein zweiter Mechanismus.
TOOLCHAIN_IMAGE ?= golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125
# TOOLCHAIN_RACE_IMAGE trägt denselben Go-Toolchain-Stand wie TOOLCHAIN_IMAGE
# (`go1.27.1`, real geprüft), aber Debian statt Alpine: der Race-Detector
# braucht einen C-Compiler zum Linken (`gcc`), den das Alpine-Image nicht
# trägt (`CGO_ENABLED=0` dort ohne ihn) — `go test -race` bricht sonst vor
# dem ersten Testlauf ab (`ADR-0050` Fitness Function: „Assembler.tables …
# ohne Mutex/Kommando-Kanal ist das ein Data Race, go test -race"). Die
# Produktions-Kompilierung (Dockerfile) bleibt CGO-frei — dieses Image trägt
# ausschließlich den Testlauf.
TOOLCHAIN_RACE_IMAGE ?= golang:1.27@sha256:b475798fb16158e6c38e8b5ca2d870fbeaa8b7fec0fc8ec64b3dc20966040635
PG_TEST_IMAGE ?= postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8
GO_MODCACHE_VOLUME ?= pg-change-feed-gomodcache

mod-download: ## Go-Module in den Volume-Cache laden (braucht Netz, Vorbereitung für netzlose Test-Läufe)
	docker run --rm -v "$(CURDIR)":/src:ro -v $(GO_MODCACHE_VOLUME):/go/pkg/mod \
	  -w /src -e GOCACHE=/tmp/gocache $(TOOLCHAIN_IMAGE) go mod download

test: ## Unit-Tests im gepinnten Toolchain-Container (netzlos, mit Race-Detector)
	docker run --rm --network none -v "$(CURDIR)":/src:ro \
	  -v $(GO_MODCACHE_VOLUME):/go/pkg/mod \
	  -w /src -e GOCACHE=/tmp/gocache -e CGO_ENABLED=1 $(TOOLCHAIN_RACE_IMAGE) go test -race ./...

test-store: ## Adapter-Tests gegen reale PostgreSQL (Testcontainer, gepinnt)
	@bash tools/harness/run-store-tests.sh

test-replication: ## Replication-Stream-Tests gegen reale PostgreSQL mit Publication/Slot (wal_level=logical, gepinnt)
	@bash tools/harness/run-replication-tests.sh

test-notify: ## natsnotify-Adapter-Tests gegen einen echten NATS-Server (Testcontainer, gepinnt, ADR-0055)
	@bash tools/harness/run-notify-tests.sh

test-integration: ## Compose-Integrationstest — Kern-CDC-Pfad, Rollen-DSN-Verifikation, Black-Box-CLI-Rundlauf (Compose + schema-rollout + Toolchain-Container, kein Gate)
	@bash tools/harness/run-integration-tests.sh

# --- Performance-Benchmarks (kein Gate; ADR-0054 §(b)) ---
# Vier eigenständige Skripte, je ein Beleg, gebündelt hinter diesem Ziel —
# analog `d-check`s `Makefile` Zeile 84 (`bench: build`): drei mit Schwelle
# (LH-QA-PER-001…003) und `tools/bench-backfill.sh`, eine Messung ohne
# Schwelle (LH-FA-CAP-009, Vertrag: harness/targets/bench-backfill.md); sie
# läuft zuletzt, damit ihr Abbruch keine Schwellen-Prüfung und nicht die
# Erzeugung von docs/user/bench-abdeckung.md verdeckt.
# Nicht Teil von `gates`/`ci`/`fullbuild`, weil kein einzelner
# Schwellenwert existiert, gegen den Aufwand/Ergebnis entscheiden würde
# (Kontrast zu coverage-gate). Braucht ein geladenes Image (make image)
# für den Feed-Container der Skripte.
.PHONY: bench
bench: image ## Performance-Benchmarks (vier Skripte: LH-QA-PER-001…003 mit Schwelle, Backfill-Messung LH-FA-CAP-009 ohne; dokumentiertes Ergebnis, kein Gate; ADR-0054 §(b))
	@bash tools/bench-source-impact.sh
	@bash tools/bench-scaling.sh
	@bash tools/bench-batch-vs-single.sh
	@bash tools/bench-backfill.sh

# --- Codegenerierung (kein Gate; Docker-only, ADR-0060) ---
# Protobuf-/gRPC-Codegenerierung laeuft ausschliesslich im gepinnten
# Toolchain-Container (`AGENTS.md` §3.1) — kein Host-protoc/-buf. Die
# Dockerfile-Stufe `proto` traegt protoc und die beiden protoc-gen-*-
# Plugins; die Folgestufe `proto-export` kopiert die
# `.proto`-Quelle per COPY hinein und erzeugt den Code **zur Build-Zeit**
# (kein Bind-Mount, kein `--user`-Workaround) — ihr ENTRYPOINT gibt das
# Erzeugnis als `tar`-Stream ueber stdout aus. Der erzeugte Go-Code liegt
# committet im Baum und wird von `make test`/`make image` mitkompiliert.
# Kein Gate: der Generator laeuft nur, wenn sich die `.proto`-Quelle aendert.
#
# Die Logik steht in `tools/harness/proto-generate.sh` (wie die uebrigen
# `tools/harness/*.sh`-Gates) statt inline im Rezept: das Skript laeuft
# explizit unter bash und setzt `pipefail`, wodurch die Pipe
# `docker run | tar -x` sicher wird (`AGENTS.md` §3.9 — ohne `pipefail`
# verschwindet ein `docker run`-Fehlschlag hinter einem erfolgreichen, aber
# leeren `tar -x`). Das Makefile-Rezept selbst liefe unter `/bin/sh`
# (`make`-Default, hier `dash`, kein `pipefail`) — der explizite
# `bash`-Aufruf loest das, ohne einen globalen `SHELL`-Override im Makefile
# zu brauchen.
PROTO_IMAGE ?= pg-change-feed:proto-export

.PHONY: proto-generate
proto-generate: ## Protobuf-/gRPC-Go-Code aus proto/cdc/stream/v1/changestream.proto erzeugen (Docker-only, Build-Zeit-Erzeugung, Host-Extraktion)
	PROTO_IMAGE=$(PROTO_IMAGE) bash tools/harness/proto-generate.sh

# --- Schemamigrationen (kein Gate; d-migrate, ADR-0043) ---
# Das neutrale Schema-YAML (tools/schema/schema.yaml) ist die Quelle der
# CDC-Schema-Form; SQL wird erzeugt, Rollouts laufen mit Pflicht-Report und
# Rollback-Artefakt. Die Quelle der Überführung ist die handgeschriebene DDL
# (internal/adapters/driven/postgresstorage/schema.sql); bis zur Überführung
# bleibt die DDL die Quelle und der Testloader die benannte Grenze (ADR-0043).
# Image per Digest gepinnt — Pin-Hebung = bewusster Commit (Modul 14).
# `schema migrate --execute` braucht DB-Zugang: kein Target davon hängt an
# GATE_CHECKS.
# SCHEMA_ROLLOUT_NETWORK trägt das Docker-Netz des Rollout-Ziels: der
# d-migrate-Container erreicht die Compose-Test-DB über den Dienstnamen
# (Netz cdc-feed-test aus compose.yaml); für ein Host-seitiges localhost-Ziel
# trägt der Aufrufer `host`. Der Default `bridge` passt für ein DB-Ziel, das
# selbst im Default-Brückennetz liegt.
D_MIGRATE_IMAGE ?= ghcr.io/pt9912/d-migrate@sha256:862dfb04c34dd17278b1bab46961363c12eeb8d464cf1776565d6285603d2c89
SCHEMA_SOURCE ?= tools/schema/schema.yaml
SCHEMA_TARGET ?= db:postgres://postgres:postgres@localhost:5432/cdc?sslmode=disable
SCHEMA_ROLLOUT_NETWORK ?= bridge
# Das d-migrate-Image läuft als uid 10001 (dmigrate) und kann in den
# Bind-Mount des Arbeitsbaums nicht schreiben — der Report- und
# Rollback-Schreibpfad trägt den Host-Nutzer (uid/gid des make-Laufs).
D_MIGRATE_RUN_USER ?= $(shell id -u):$(shell id -g)

.PHONY: schema-validate schema-rollout
schema-validate: ## d-migrate: neutrales Schema prüfen (netzlos; Vorlauf vor generate/migrate, kein Gate)
	@if [ ! -f "$(SCHEMA_SOURCE)" ]; then \
	  echo "FEHLER: $(SCHEMA_SOURCE) fehlt — das neutrale Schema-YAML ist die Erstlieferung des d-migrate-Einbaus (ADR-0043); Überführungsquelle ist internal/adapters/driven/postgresstorage/schema.sql" >&2; \
	  exit 2; \
	fi
	docker run --rm --user "$(D_MIGRATE_RUN_USER)" --network none -v "$(CURDIR)":/work -w /work $(D_MIGRATE_IMAGE) schema validate --source $(SCHEMA_SOURCE)

# Rollout des neutralen Schemas: Precheck (--plan-only), Wache
# tools/schema/rolloutguard, Vorlauf `DROP VIEW cdc.<name>` bei einer
# View-Signatur-Änderung (ADR-0114), `schema migrate --execute` mit
# Pflicht-Report und Rollback-Artefakt, danach die vier psql-Nacharbeit-
# Dateien tools/schema/nacharbeit-*.sql (Rollen, Views cdc.metrics und
# cdc.heartbeat, Administrations-Funktionen), die d-migrate nicht ausdrückt
# (ADR-0043). Vertrag, Reihenfolge, Exit-Codes und Grenzen:
# harness/targets/schema-rollout.md.
schema-rollout: schema-validate ## d-migrate: Schema-Rollout --execute mit Pflicht-Report und Rollback-Artefakt (braucht DB-Zugang, kein Gate)
	@mkdir -p tools/schema
	@docker run --rm --user "$(D_MIGRATE_RUN_USER)" --network $(SCHEMA_ROLLOUT_NETWORK) -v "$(CURDIR)":/work -w /work $(D_MIGRATE_IMAGE) schema migrate --source $(SCHEMA_SOURCE) --target "$(SCHEMA_TARGET)" --plan-only --report tools/schema/rollout-precheck.yaml; \
	plan_exit=$$?; \
	allow_destructive=""; \
	drop_views=""; \
	if [ "$$plan_exit" = "8" ]; then \
	  guard_out=$$(docker run --rm --network none -v "$(CURDIR)":/src:ro -v $(GO_MODCACHE_VOLUME):/go/pkg/mod -w /src -e GOCACHE=/tmp/gocache $(TOOLCHAIN_IMAGE) go run ./tools/schema/rolloutguard tools/schema/rollout-precheck.yaml) && guard_ok=1 || guard_ok=0; \
	  if [ "$$guard_ok" = "1" ]; then \
	    if printf '%s\n' "$$guard_out" | grep -qx 'allow-destructive'; then \
	      echo "schema-rollout: bekannte Fremdobjekt-Blocker (ADR-0043) - --execute laeuft mit --allow-destructive"; \
	      allow_destructive="--allow-destructive"; \
	    fi; \
	    drop_views=$$(printf '%s\n' "$$guard_out" | sed -n 's/^drop-view //p'); \
	  fi; \
	fi; \
	for v in $$drop_views; do \
	  echo "schema-rollout: Vorlauf (ADR-0114) - View-Signatur-Aenderung, DROP VIEW cdc.$$v"; \
	  docker run --rm --network $(SCHEMA_ROLLOUT_NETWORK) $(PG_TEST_IMAGE) psql "$(SCHEMA_TARGET:db:%=%)" -v ON_ERROR_STOP=1 -c "DROP VIEW cdc.$$v" || exit 1; \
	done; \
	docker run --rm --user "$(D_MIGRATE_RUN_USER)" --network $(SCHEMA_ROLLOUT_NETWORK) -v "$(CURDIR)":/work -w /work $(D_MIGRATE_IMAGE) schema migrate --source $(SCHEMA_SOURCE) --target "$(SCHEMA_TARGET)" --execute $$allow_destructive --report tools/schema/plan.yaml --generate-rollback --rollback-output tools/schema/down.sql
	docker run --rm --network $(SCHEMA_ROLLOUT_NETWORK) -v "$(CURDIR)":/work:ro $(PG_TEST_IMAGE) psql "$(SCHEMA_TARGET:db:%=%)" -v ON_ERROR_STOP=1 -f /work/tools/schema/nacharbeit-roles.sql
	docker run --rm --network $(SCHEMA_ROLLOUT_NETWORK) -v "$(CURDIR)":/work:ro $(PG_TEST_IMAGE) psql "$(SCHEMA_TARGET:db:%=%)" -v ON_ERROR_STOP=1 -f /work/tools/schema/nacharbeit-observability.sql
	docker run --rm --network $(SCHEMA_ROLLOUT_NETWORK) -v "$(CURDIR)":/work:ro $(PG_TEST_IMAGE) psql "$(SCHEMA_TARGET:db:%=%)" -v ON_ERROR_STOP=1 -f /work/tools/schema/nacharbeit-heartbeat.sql
	docker run --rm --network $(SCHEMA_ROLLOUT_NETWORK) -v "$(CURDIR)":/work:ro $(PG_TEST_IMAGE) psql "$(SCHEMA_TARGET:db:%=%)" -v ON_ERROR_STOP=1 -f /work/tools/schema/nacharbeit-administration.sql

help: ## Diese Hilfe
	@grep -hE '^[a-z-]+:.*##' $(MAKEFILE_LIST) | sort | awk 'BEGIN{FS=":.*##"}{printf "  %-14s %s\n",$$1,$$2}'

# gates haengt allein an record-gates; record-gates haengt an ALLEN akkumulierten
# Checks — der Nachweis laeuft strikt nach den Checks (Ordnungskante), waehrend make
# -j die Checks parallel faehrt. Das record-gates-Rezept liefert harness/mk/enforce.mk.
gates: record-gates ## Alle Gates (Checks parallel, Nachweis zuletzt)
record-gates: $(GATE_CHECKS)
