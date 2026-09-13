# harness/mk/coverage.mk — Coverage-Gate-Fragment (ADR-0054).
# Vierte Docker-Stage `coverage` (Dockerfile, nach `deps`) misst real die
# Go-Test-Coverage ueber internal/...+cmd/... und prueft sie ueber
# tools/coverage-gate.sh gegen THRESHOLD; haengt coverage-gate an
# GATE_CHECKS — der Root-Aggregator faehrt es via make gates.
#
# Kalibrierungs-Bindung (harness/README.md §Sensors, ADR-0054): bootstrap-
# aware Gate — Einstiegsstufe 35 % (realer Ist-Stand beim ersten Lauf:
# 39.6 %, abgerundet auf den naechsten vollen 5-%-Schritt). Endstufe 80 %
# steht fest; Hochschalt-Trigger „naechste Coverage-Verbesserung schliesst
# die Luecke zur naechsten Stufe" bis 80 % erreicht ist. Override: `make
# coverage-gate THRESHOLD=…`; Senkung unter die hier gueltige Stufe nur per
# ADR (AGENTS.md §3.6).
THRESHOLD ?= 35

# `--no-cache-filter coverage`: erzwingt die Neu-Auswertung der
# Coverage-Stage, ohne den deps-Cache zu verlieren — ein stale Layer-Hash
# darf keine rote Stage maskieren (Muster: d-check NO_CACHE_FILTER_COV).
NO_CACHE_FILTER_COV := --no-cache-filter coverage

.PHONY: coverage-gate
coverage-gate: ## Coverage-Schwelle (Kalibrierungs-Bindung: bootstrap-aware 35 % -> 80 %, Historie in harness/README §Sensors)
	docker build $(NO_CACHE_FILTER_COV) \
	    --build-arg COVERAGE_THRESHOLD=$(THRESHOLD) \
	    --target coverage -t pg-change-feed:coverage .

GATE_CHECKS += coverage-gate
