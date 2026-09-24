# harness/mk/coverage.mk — Coverage-Gate-Fragment (ADR-0054, ADR-0071).
# Vierte Docker-Stage `coverage` (Dockerfile, nach `deps`) misst real die
# Go-Test-Coverage ueber die netzlos pruefbare Flaeche (internal/...+cmd/...
# +gen/... ohne die Pakete, deren Testlauf einen externen Dienst voraussetzt)
# und prueft sie ueber tools/coverage-gate.sh gegen THRESHOLD; haengt
# coverage-gate an GATE_CHECKS — der Root-Aggregator faehrt es via make gates.
# Vor dem Bau haelt tools/harness/db-package-lists-check.sh die namentlichen
# Paketlisten des ausgenommenen Gegenstands gleich (harness/sensors/coverage-gate.md
# §Grenze Nr. 4).
#
# Kalibrierungs-Bindung (harness/README.md §Sensors, ADR-0054 §(a)): bootstrap-
# aware Gate. Die geltende Stufe ist THRESHOLD unten und steht ausschliesslich
# an diesem Ort — die uebrigen Traeger nennen nur die Rampe (Einstieg 70 %,
# Endstufe 80 % fest, Nutzer-Entscheidung). Hochschalt-Trigger „naechste
# Coverage-Verbesserung schliesst die Luecke zur naechsten Stufe" — mit der
# Endstufe ist er ausgeschoepft. Override: `make coverage-gate THRESHOLD=…`;
# Senkung unter die hier geltende Stufe nur per ADR (AGENTS.md §3.6).
THRESHOLD ?= 80

# `--no-cache-filter coverage`: erzwingt die Neu-Auswertung der
# Coverage-Stage, ohne den deps-Cache zu verlieren — ein stale Layer-Hash
# darf keine rote Stage maskieren (Muster: d-check NO_CACHE_FILTER_COV).
NO_CACHE_FILTER_COV := --no-cache-filter coverage

.PHONY: coverage-gate
coverage-gate: ## Coverage-Schwelle (bootstrap-aware Rampe Einstieg 70 % -> Endstufe 80 %; Bindung in harness/README §Sensors)
	bash tools/harness/db-package-lists-check.sh
	docker build $(NO_CACHE_FILTER_COV) \
	    --build-arg COVERAGE_THRESHOLD=$(THRESHOLD) \
	    --target coverage -t pg-change-feed:coverage .

GATE_CHECKS += coverage-gate
