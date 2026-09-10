# d-check.mk — Doku-Referenz-Gate via d-check. Emittiert von ai-harness-init,
# adaptiert aus `d-check --print-mk`: doc-check -> docs-check (das Befund-Gate,
# einziges als Gate behauptetes Target) und DCHECK_DIGEST auf den erzeugenden
# Image-Digest gepinnt (Reproduzierbarkeit). advisory doc-*-Targets verbatim.
# Einbinden: `include d-check.mk`; eigene .d-check.yml danebenlegen.
DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v0.75.0
DCHECK_DIGEST ?= sha256:18e9cd857f8db3569526d1f9a3cbeba8af51e9f2dd84c17a22444028b977c3da
# TRACE_FLAGS: optionale Flags für die RTM-Targets (z. B. --json).
TRACE_FLAGS ?=

# Ein gesetzter DCHECK_DIGEST sticht den Tag von DCHECK_IMAGE.
ifeq ($(strip $(DCHECK_DIGEST)),)
DCHECK_REF := $(DCHECK_IMAGE)
else
DCHECK_REF := ghcr.io/pt9912/d-check@$(DCHECK_DIGEST)
endif

.PHONY: docs-check
docs-check: ## Doku-Referenzen prüfen (Befund-Gate)
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF)

.PHONY: doc-trace
doc-trace: ## Requirements Traceability Matrix auf stdout (advisory, DC-FA-CLI-009)
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --trace $(TRACE_FLAGS)

.PHONY: doc-complete
doc-complete: ## Vollständigkeits-Gate: Requirements-Waise ⇒ Exit 1 (DC-FA-CLI-011)
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --trace --require-complete $(TRACE_FLAGS)

.PHONY: doc-doctor
doc-doctor: ## erklärende Diagnose mit Fix-Kandidaten (DC-FA-CLI-007)
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --doctor

.PHONY: doc-repair
doc-repair: ## Reparatur-Patch (unified diff) auf stdout, git-apply-rein (DC-FA-CLI-008)
	@docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --repair

.PHONY: doc-immutable
doc-immutable: ## Doc-/ADR-Immutabilität via git-Diff (Modul vcs); RANGE=base..head oder STAGED=1 (DC-FA-VCS-001)
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --enable vcs --disable links --disable anchors --disable ids --disable matrix --disable external --disable codepaths --disable spans --disable hostpaths --disable diagrams --disable versions --disable pins --disable immutable --disable commits --disable planning --disable tracked --disable targets --disable citations --disable sources --disable structure --disable workflows --disable reviews $(if $(STAGED),--staged,--range $(RANGE))

.PHONY: doc-commits
doc-commits: ## Commit-Message-Traceability via Modul commits; RANGE=base..head (DC-FA-COMMITS-001)
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --enable commits --disable links --disable anchors --disable ids --disable matrix --disable external --disable codepaths --disable spans --disable hostpaths --disable diagrams --disable versions --disable pins --disable immutable --disable vcs --disable planning --disable tracked --disable targets --disable citations --disable sources --disable structure --disable workflows --disable reviews --range $(RANGE)

# commit-traceability — Standing-Gate der Commit-Traceability (ADR-0045):
# die Traceability-Regel (harness/README.md §Traceability rules) je Commit-
# Message der letzten fünf Commits, mechanisch getragen. Positive Hälfte
# (je Message eine Vertrags-Kennung): d-check Modul commits, konfiguriert
# im commits-Abschnitt der .d-check.yml (Befund commit-untraceable).
# Grenz-Hälfte (keine Struktur-ID im Betreff): der Shell-Sensor — das
# commits-Modul kennt heute keine verbotene-Muster-Option; ein zweiter
# Sensor über dieselbe Hälfte wäre eine zweite Quelle. Gate statt
# commit-msg-Hook (ADR-0045, Alternativen): das Gate läuft im
# `make gates`-Bündel und deckt Commits jedes Urhebers im Fenster.
# RANGE=base..head überschreibt den Standing-Default (ADR-0045);
# COMMIT_TRACE_RANGE trägt denselben Wert für die Aufrufe, die den
# Standing-Default nicht berühren wollen.
COMMIT_TRACE_RANGE ?= $(if $(RANGE),$(RANGE),HEAD~5..HEAD)

.PHONY: commit-traceability
commit-traceability: ## Commit-Traceability als Standing-Gate (letzte 5 Commits; ADR-0045); RANGE=base..head
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --enable commits --disable links --disable anchors --disable ids --disable matrix --disable external --disable codepaths --disable spans --disable hostpaths --disable diagrams --disable versions --disable pins --disable immutable --disable vcs --disable planning --disable tracked --disable targets --disable citations --disable sources --disable structure --disable workflows --disable reviews --range $(COMMIT_TRACE_RANGE)
	bash tools/harness/commit-traceability.sh "$(COMMIT_TRACE_RANGE)"

.PHONY: doc-planning
doc-planning: ## Planning-Lifecycle-Konsistenz (Roadmap <-> in-progress) via Modul planning; hermetisch, ohne Range (DC-FA-PLAN-001)
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --enable planning --disable links --disable anchors --disable ids --disable matrix --disable external --disable codepaths --disable spans --disable hostpaths --disable diagrams --disable versions --disable pins --disable immutable --disable vcs --disable commits --disable tracked --disable targets --disable citations --disable sources --disable structure --disable workflows --disable reviews

.PHONY: doc-tracked
doc-tracked: ## Getrackt-Status aufloesbarer Referenz-Ziele via Modul tracked; braucht .git im Mount, ohne Range (DC-FA-TRK-001)
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --enable tracked --disable links --disable anchors --disable ids --disable matrix --disable external --disable codepaths --disable spans --disable hostpaths --disable diagrams --disable versions --disable pins --disable immutable --disable vcs --disable commits --disable planning --disable targets --disable citations --disable sources --disable structure --disable workflows --disable reviews

.PHONY: doc-targets
doc-targets: ## Deklarations-Konsistenz Doku<->Build-Targets via Modul targets; hermetisch, ohne Range (DC-FA-TGT-001)
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --enable targets --disable links --disable anchors --disable ids --disable matrix --disable external --disable codepaths --disable spans --disable hostpaths --disable diagrams --disable versions --disable pins --disable immutable --disable vcs --disable commits --disable planning --disable tracked --disable citations --disable sources --disable structure --disable workflows --disable reviews

.PHONY: doc-structure
doc-structure: ## Struktur-Invarianten innerhalb der Dokumente via Modul structure; hermetisch, ohne Range (DC-FA-STRUCT-001)
	docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --enable structure --disable links --disable anchors --disable ids --disable matrix --disable external --disable codepaths --disable spans --disable hostpaths --disable diagrams --disable versions --disable pins --disable immutable --disable vcs --disable commits --disable planning --disable tracked --disable targets --disable citations --disable sources --disable workflows --disable reviews

.PHONY: doc-usage
doc-usage: ## Aufruf und Optionen von d-check selbst (--help)
	@docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --help

.PHONY: doc-help
doc-help: ## diese Liste der doc-*-Targets
	@grep -hE '^docs?-[a-z-]+:.*## ' $(MAKEFILE_LIST) | sort | sed -E 's/:.*## /  /'