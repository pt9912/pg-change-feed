# a-check.mk — Architektur-Gate via a-check, zum `include` in das Makefile
# dieses Repos. Erzeugt von `a-check --print-mk`; A_CHECK_IMAGE ist manuell
# auf den Release-Digest gepinnt (Pin-Hebung = bewusster Commit, AC-QA-03 /
# Modul 14).
#
# Benutzerhandbuch (aufgabenorientiert, deutsch):
#   https://github.com/pt9912/a-check/blob/main/docs/user/benutzerhandbuch.md
#
A_CHECK_IMAGE ?= ghcr.io/pt9912/a-check@sha256:2368f7b3a84f1dc5d075edccfe2201e19947d12fcbc8eaf4df84ef162d94f422

# Maschinenform der §2-Schichten-Constraints (ADR-0041):
# a-check haengt an GATE_CHECKS und laeuft damit im `make gates`-Buendel mit.

# Container-Runtime ueber eine Indirektion (podman/nerdctl/docker); wer eine
# eigene Runtime nutzt, definiert sie VOR dem `include`.
DOCKER ?= docker

.PHONY: a-check a-check-graph
a-check: ## Architektur: Hexagon-Regeln via a-check (netzlos, read-only).
	$(DOCKER) run --rm --network none -v "$(CURDIR)":/src:ro $(A_CHECK_IMAGE) /src

a-check-graph: ## Architektur-Graph (Mermaid) aus .a-check.yml auf stdout (read-only, kein Scan).
	$(DOCKER) run --rm --network none -v "$(CURDIR)":/src:ro $(A_CHECK_IMAGE) --print-graph /src

GATE_CHECKS += a-check
