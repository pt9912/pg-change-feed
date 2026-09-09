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

help: ## Diese Hilfe
	@grep -hE '^[a-z-]+:.*##' $(MAKEFILE_LIST) | sort | awk 'BEGIN{FS=":.*##"}{printf "  %-14s %s\n",$$1,$$2}'

# gates haengt allein an record-gates; record-gates haengt an ALLEN akkumulierten
# Checks — der Nachweis laeuft strikt nach den Checks (Ordnungskante), waehrend make
# -j die Checks parallel faehrt. Das record-gates-Rezept liefert harness/mk/enforce.mk.
gates: record-gates ## Alle Gates (Checks parallel, Nachweis zuletzt)
record-gates: $(GATE_CHECKS)
