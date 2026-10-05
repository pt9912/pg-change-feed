# `make pin-stale-*` (P3 bis P9) — Upstream-Pin-Freshness

Ausführliche Fassung der Index-Zeilen aus [`harness/README.md` §Sensors](../README.md#sensors-feedback-gates); die Zeile dort trägt einen Satz und verlinkt hierher.

## `make pin-stale-race`/`make pin-stale-pgtest`/`make pin-stale-dmigrate`/`make pin-stale-acheck`

Upstream-Pin-Freshness P3–P6 ([`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Pin-Inventar): je Achse ein Digest-Vergleich einer Makefile-/`a-check.mk`-Variable (`TOOLCHAIN_RACE_IMAGE`, `PG_TEST_IMAGE`, `D_MIGRATE_IMAGE`, `A_CHECK_IMAGE`) gegen den aktuellen Registry-Digest, über das gemeinsame `tools/harness/pin-stale.sh` (analog `make image-stale`, aber Quelle ist eine Variable statt einer Dockerfile-`FROM`-Zeile); P5/P6 sind reine Digest-Pins ohne eigenen Tag, deshalb Vergleich gegen den jeweiligen `:latest`-Tag der GHCR-Pakete. Real ausgeführt: P3/P4/P5 zeigen echten Drift (Upstream-Tags wurden seit dem Pin neu gebaut), P6 ist aktuell

**Bindung:** kein Gate, [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz

## `make pin-stale-dcheck`

Upstream-Pin-Freshness P7 ([`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Pin-Inventar): zwei Achsen für `d-check.mk`s `DCHECK_IMAGE`/`DCHECK_DIGEST` — Digest-Drift (trägt der gepinnte Tag noch denselben Bau?) und Tag-Frische (existiert ein neuerer `pt9912/d-check`-Release als der gepinnte Tag, über die GitHub-Releases-API?)

**Bindung:** kein Gate, [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz

## `make pin-stale-baseline`

Upstream-Pin-Freshness P8 ([`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Pin-Inventar): die adoptierte Kurs-Baseline-Version (`harness/conventions.md` §Baseline) gegen den neuesten `pt9912/ai-harness-course`-Release (GitHub-Releases-API) — kein Docker-Pin, eine Baseline-Aktualisierung bleibt ein bewusster Bootstrap-Vorgang (Modul 2)

**Bindung:** kein Gate, [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz

## `make pin-stale-actions`

Upstream-Pin-Freshness P9 ([`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Pin-Inventar): jede SHA-gepinnte `uses:`-Zeile über alle `.github/workflows/*.yml` (`AGENTS.md` §3.8) gegen zwei Achsen — Tag-Mutation (`git ls-remote` gegen den im Kommentar genannten Tag; zeigt er noch auf denselben SHA?) und Tag-Frische (neuester Release des Action-Repos über die GitHub-Releases-API). Ein doppelt referenziertes Repo@Tag über mehrere Workflow-Dateien wird nur einmal geprüft

**Bindung:** kein Gate, [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz
