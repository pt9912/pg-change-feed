# Makefile — generiert von ai-harness-init (Aggregator). Die Gate-Belange leben als
# Fragmente unter harness/mk/*.mk; jedes haengt seine Checks an GATE_CHECKS. Der
# Gate-Nachweis (record-gates) laeuft strikt ZULETZT via Ordnungskante auf GATE_CHECKS
# — waehrend make -j die Checks parallelisiert; .NOTPARALLEL ist bewusst NICHT gewaehlt
# (das serialisierte das ganze Makefile). Sprach-agnostisch: ohne --lang matchen nur
# baseline/doc-gate/enforce, mit --lang zusaetzlich das Code-Gate-Fragment.
GATE_CHECKS :=

.PHONY: gates help

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
image: ## Baut das OCI-Image; Image-Hash nach harness/image-hash.txt (Modul 14)
	docker buildx build --load --metadata-file harness/image-hash.raw -t ghcr.io/pt9912/pg-change-feed:dev . && grep -o '"containerimage.digest":[[:space:]]*"sha256:[0-9a-f]*' harness/image-hash.raw | head -1 | grep -o 'sha256:[0-9a-f]*' > harness/image-hash.txt && rm harness/image-hash.raw

image-stale: ## Advisory: FROM-Digests gegen Registry-Digests (Modul 14, braucht Netz)
	@bash tools/harness/image-stale.sh

# --- Tests (kein Gate; Docker-only, gepinnte Images) ---
# Toolchain-Container = derselbe gepinnte Digest wie im Dockerfile; der
# PostgreSQL-Testcontainer trägt seinen Digest aus `docker manifest inspect
# postgres:18-alpine` (amd64). Caches leben in Docker-Volumes, Daten im
# Container — nichts davon im Arbeitsbaum.
TOOLCHAIN_IMAGE ?= golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125
PG_TEST_IMAGE ?= postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8
GO_MODCACHE_VOLUME ?= pg-change-feed-gomodcache

mod-download: ## Go-Module in den Volume-Cache laden (braucht Netz, Vorbereitung für netzlose Test-Läufe)
	docker run --rm -v "$(CURDIR)":/src:ro -v $(GO_MODCACHE_VOLUME):/go/pkg/mod \
	  -w /src -e GOCACHE=/tmp/gocache $(TOOLCHAIN_IMAGE) go mod download

test: ## Unit-Tests im gepinnten Toolchain-Container (netzlos)
	docker run --rm --network none -v "$(CURDIR)":/src:ro \
	  -v $(GO_MODCACHE_VOLUME):/go/pkg/mod \
	  -w /src -e GOCACHE=/tmp/gocache $(TOOLCHAIN_IMAGE) go test ./...

test-store: ## Adapter-Tests gegen reale PostgreSQL (Testcontainer, gepinnt)
	@bash tools/harness/run-store-tests.sh

# --- Schemamigrationen (kein Gate; d-migrate, ADR-0043) ---
# Das neutrale Schema-YAML (tools/schema/schema.yaml) ist die Quelle der
# CDC-Schema-Form; SQL wird erzeugt, Rollouts laufen mit Pflicht-Report und
# Rollback-Artefakt. Die Quelle der Überführung ist die handgeschriebene DDL
# (internal/adapters/driven/postgresstorage/schema.sql); bis zur Überführung
# bleibt die DDL die Quelle und der Testloader die benannte Grenze (ADR-0043).
# Image per Digest gepinnt — Pin-Hebung = bewusster Commit (Modul 14).
# `schema migrate --execute` braucht DB-Zugang: kein Target davon hängt an
# GATE_CHECKS.
D_MIGRATE_IMAGE ?= ghcr.io/pt9912/d-migrate@sha256:8d1433990ee4dd6a975b29d8db18356d1204ae6e45f96f312872ba5eca1230ea
SCHEMA_SOURCE ?= tools/schema/schema.yaml
SCHEMA_TARGET ?= db:postgres://postgres:postgres@localhost:5432/cdc?sslmode=disable

.PHONY: schema-validate schema-rollout
schema-validate: ## d-migrate: neutrales Schema prüfen (netzlos; Vorlauf vor generate/migrate, kein Gate)
	@if [ ! -f "$(SCHEMA_SOURCE)" ]; then \
	  echo "FEHLER: $(SCHEMA_SOURCE) fehlt — das neutrale Schema-YAML ist die Erstlieferung des d-migrate-Einbaus (ADR-0043); Überführungsquelle ist internal/adapters/driven/postgresstorage/schema.sql" >&2; \
	  exit 2; \
	fi
	docker run --rm --network none -v "$(CURDIR)":/work -w /work $(D_MIGRATE_IMAGE) schema validate --source $(SCHEMA_SOURCE)

schema-rollout: schema-validate ## d-migrate: Schema-Rollout --execute mit Pflicht-Report und Rollback-Artefakt (braucht DB-Zugang, kein Gate)
	@mkdir -p tools/schema
	docker run --rm -v "$(CURDIR)":/work -w /work $(D_MIGRATE_IMAGE) schema migrate --source $(SCHEMA_SOURCE) --target $(SCHEMA_TARGET) --execute --report tools/schema/plan.yaml --generate-rollback --rollback-output tools/schema/down.sql

help: ## Diese Hilfe
	@grep -hE '^[a-z-]+:.*##' $(MAKEFILE_LIST) | sort | awk 'BEGIN{FS=":.*##"}{printf "  %-14s %s\n",$$1,$$2}'

# gates haengt allein an record-gates; record-gates haengt an ALLEN akkumulierten
# Checks — der Nachweis laeuft strikt nach den Checks (Ordnungskante), waehrend make
# -j die Checks parallel faehrt. Das record-gates-Rezept liefert harness/mk/enforce.mk.
gates: record-gates ## Alle Gates (Checks parallel, Nachweis zuletzt)
record-gates: $(GATE_CHECKS)