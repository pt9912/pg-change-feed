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
ifdef VERSION
image: ## Baut und pusht das OCI-Image nach GHCR+Docker Hub (VERSION=<semver>, optional LATEST=true — ADR-0051)
	docker buildx build --push \
	  -t ghcr.io/pt9912/pg-change-feed:$(VERSION) \
	  -t docker.io/pt9912/pg-change-feed:$(VERSION) \
	  $(if $(filter true,$(LATEST)),-t ghcr.io/pt9912/pg-change-feed:latest -t docker.io/pt9912/pg-change-feed:latest,) \
	  --metadata-file harness/image-hash.raw . && grep -o '"containerimage.digest":[[:space:]]*"sha256:[0-9a-f]*' harness/image-hash.raw | head -1 | grep -o 'sha256:[0-9a-f]*' > harness/image-hash.txt && rm harness/image-hash.raw
else
image: ## Baut das OCI-Image; Image-Hash nach harness/image-hash.txt (lokal, nicht committet — ADR-0103); VERSION=<semver> fuer den Multi-Registry-Push (ADR-0051)
	docker buildx build --load --metadata-file harness/image-hash.raw -t ghcr.io/pt9912/pg-change-feed:dev . && grep -o '"containerimage.digest":[[:space:]]*"sha256:[0-9a-f]*' harness/image-hash.raw | head -1 | grep -o 'sha256:[0-9a-f]*' > harness/image-hash.txt && rm harness/image-hash.raw
endif

image-stale: ## Advisory: FROM-Digests gegen Registry-Digests (Modul 14, braucht Netz)
	@bash tools/harness/image-stale.sh

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
# Drei eigenständige Skripte (LH-QA-PER-001…003), je ein Beleg, gebündelt
# hinter diesem Ziel — analog `d-check`s `Makefile` Zeile 84
# (`bench: build`); nicht Teil von `gates`/`ci`/`fullbuild`, weil kein
# einzelner Schwellenwert existiert, gegen den Aufwand/Ergebnis
# entscheiden würde (Kontrast zu coverage-gate). Braucht ein zuvor
# geladenes Image (make image) für den Feed-Container der Skripte.
.PHONY: bench
bench: image ## Performance-Benchmarks LH-QA-PER-001…003 (drei Skripte, dokumentiertes Ergebnis, kein Gate; ADR-0054 §(b))
	@bash tools/bench-source-impact.sh
	@bash tools/bench-scaling.sh
	@bash tools/bench-batch-vs-single.sh

# --- Codegenerierung (kein Gate; Docker-only, ADR-0060) ---
# Protobuf-/gRPC-Codegenerierung laeuft ausschliesslich im gepinnten
# Toolchain-Container (`AGENTS.md` §3.1) — kein Host-protoc/-buf. Die
# Dockerfile-Stufe `proto` traegt protoc und die beiden protoc-gen-*-
# Plugins; die Folgestufe `proto-export` (seit slice-104) kopiert die
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

# Der Rollout trägt die CDC-Schema-Form vollständig — der CHECK über der
# Operation ist mit d-migrate 1.3.0 deklarativ konvergent (schema.yaml,
# chk_change_operation) und braucht keine Nacharbeit mehr (slice-015 hat das
# real gegen einen frischen Rollout getestet, Post-Compare grün). Die drei
# SQL-Views (LH-FA-SST-002) sind seit slice-016 ebenfalls deklarativ
# überführt (schema.yaml, `views:`-Knoten mit `source_dialect: postgresql`
# und `columns:`-Signatur) — d-migrate 1.3.1 behebt den Post-Compare-Drift
# auf frisch angelegten Sichten (Changelog: fehlendes
# `ViewDefinition.sourceDialect` im Fingerabdruck-Vergleich), real gegen
# einen frischen Rollout UND einen Folgelauf gegen eine bereits migrierte
# Instanz getestet (slice-016, Exit 0 in beiden Fällen, kein
# `VIEW_SIGNATURE_UNKNOWN`-Blocker dank `columns:`-Signatur). Der
# `nacharbeit-views.sql`-Schritt entfällt damit; die verbleibenden
# psql-Nacharbeit-Schritte tragen andere Objektklassen, die d-migrate nicht
# ausdrückt (Rollen/GRANTs, Observability-/Heartbeat-Views mit GRANT auf
# cdc_reader). Zwei Schritte seit slice-011 (LH-QA-SEC-001…003,
# LH-FA-SST-004): tools/schema/nacharbeit-roles.sql trägt die drei
# Least-Privilege-Rollen, tools/schema/nacharbeit-observability.sql die
# Metriken-Minimum-View cdc.metrics, die auf cdc_reader grantet — deshalb
# läuft sie nach der Rollen-Datei. Ein weiterer Schritt seit slice-012
# (LH-FA-ADM-002, LH-QA-OPS-002): tools/schema/nacharbeit-heartbeat.sql
# trägt die Health-View cdc.heartbeat, die sich aus demselben Grund wie
# cdc.metrics selbst an cdc_reader grantet — sie läuft deshalb ebenfalls
# nach der Rollen-Datei. Ein weiterer Schritt seit slice-036 (ADR-0050,
# LH-FA-ADM-001, LH-FA-CFG-005): tools/schema/nacharbeit-administration.sql
# trägt die vier schreibenden SQL-Funktionen der Antrags-Queue
# (cdc.enable_table/cdc.disable_table sowie seit slice-066
# cdc.exclude_column/cdc.include_column) und die CHECK-Klausel
# chk_administration_request_kind mit den vier Antragsarten — d-migrate 1.3.1
# generiert Funktions-DDL korrekt über den `functions:`-Knoten, aber `schema
# migrate --execute` bricht für jede dort deklarierte Funktion mit
# POST_EXECUTE_DRIFT (Exit 5) ab (real reproduziert, auch mit einer trivialen
# No-Arg-Funktion), und eine neue CHECK-Klausel an einer bestehenden Tabelle
# trägt derselbe Lauf ebenfalls nicht (Exit 5, BEO-PGC/d-migrate-nacharbeit).
# Die Datei grantet EXECUTE an cdc_admin und läuft deshalb ebenfalls nach der
# Rollen-Datei.
#
# Zentrale Idempotenz-Wache (ADR-0043,
# BEO-PGC/schema-rollout-fremdobjekte): Ein zweiter Lauf gegen ein bereits
# migriertes Ziel blockiert sonst mit Exit 8, weil die sechs Fremdobjekte
# aus den vier nacharbeit-*.sql-Dateien außerhalb des neutralen Modells
# liegen und d-migrate ihren Abbau plant — unabhängig davon, ob dieser Lauf
# sonst inhaltlich nichts oder eine echte neue Schema-Änderung trägt. Ein
# vorgelagerter --plan-only-Lauf (kein --execute, liest nur) schreibt
# denselben Report nach tools/schema/rollout-precheck.yaml (eigene Datei,
# damit der committete Pflicht-Report tools/schema/plan.yaml ausschließlich
# echte --execute-Läufe belegt); endet er blockierend (Exit 8), entscheidet
# tools/schema/rolloutguard anhand des strukturierten Reports, ob
# ausschließlich die sechs bekannten Objekte blockieren. Nur dann läuft der
# reguläre --execute-Schritt zusätzlich mit --allow-destructive. --execute
# selbst läuft in jedem Fall, damit jede echte, gleichzeitig anstehende
# Schema-Änderung im selben Lauf wirksam bleibt (Regressionsbeleg:
# tools/harness/run-schema-rollout-guard-test.sh Lauf 3). Die vier
# nacharbeit-*.sql-Schritte laufen danach unverändert und legen die sechs
# bekannten Objekte sofort wieder an (CREATE OR REPLACE, dieselbe
# Idempotenz wie bei jedem anderen Lauf) — ihr kurzes reales Fehlen
# zwischen --execute und dem ersten nacharbeit-Schritt bleibt folgenlos.
# Jeder andere Fall (kein Blocker, ein unbekannter Blocker, eine andere
# Blocker-Klasse) läuft ohne --allow-destructive und bricht bei einer
# echten neuen destruktiven Änderung weiterhin mit Exit 8 ab.
#
# Grenze: Der Precheck- und der --execute-Lauf sind zwei unabhängige,
# sequenzielle docker-run-Aufrufe gegen denselben lebenden Ziel-Zustand —
# kein d-migrate-Flag liest einen zuvor geprüften Plan zur Ausführung
# wieder ein. Ein zwischen beiden Läufen neu entstehender destruktiver
# Blocker würde vom Precheck nicht erfasst, liefe aber unter dem bereits
# gesetzten --allow-destructive durch (enges, aber reales Fenster).
schema-rollout: schema-validate ## d-migrate: Schema-Rollout --execute mit Pflicht-Report und Rollback-Artefakt (braucht DB-Zugang, kein Gate)
	@mkdir -p tools/schema
	@docker run --rm --user "$(D_MIGRATE_RUN_USER)" --network $(SCHEMA_ROLLOUT_NETWORK) -v "$(CURDIR)":/work -w /work $(D_MIGRATE_IMAGE) schema migrate --source $(SCHEMA_SOURCE) --target "$(SCHEMA_TARGET)" --plan-only --report tools/schema/rollout-precheck.yaml; \
	plan_exit=$$?; \
	allow_destructive=""; \
	if [ "$$plan_exit" = "8" ] && docker run --rm --network none -v "$(CURDIR)":/src:ro -v $(GO_MODCACHE_VOLUME):/go/pkg/mod -w /src -e GOCACHE=/tmp/gocache $(TOOLCHAIN_IMAGE) go run ./tools/schema/rolloutguard tools/schema/rollout-precheck.yaml; then \
	  echo "schema-rollout: nur bekannte Fremdobjekt-Blocker (ADR-0043) - --execute laeuft mit --allow-destructive"; \
	  allow_destructive="--allow-destructive"; \
	fi; \
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
