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

### Bump-Ablauf: Vergleich vor dem Löschen der alten Baseline

Meldet dieses Ziel einen neueren Tag, wird er neben den alten vendored
(`.harness/baseline/<neu>/` neben `.harness/baseline/<alt>/`). Das alte
Verzeichnis fällt erst, wenn der Plan des Bump-Slice die drei Schritte unten mit
Befehl und gedruckter Zahl trägt (Baseline-Regelwerk
`modul-02-harness-bootstrap.md` §Freshness-Audit der vendored Baseline). Die
Schritte sind eine Verfahrensregel, kein Gate: ein ausgefülltes Dokument weicht
im Inhalt beabsichtigt von seiner Vorlage ab, nur Gliederung und feste Klauseln
sind vergleichbar · seit slice-baseline-6-14-0-dokumente-nachziehen.

1. **Delta der Baseline.** `diff -r` alte ↔ neue Baseline, getrennt für
   `templates/` und `regelwerk/`, die Versionsstrings in Kopien im
   Temp-Verzeichnis normalisiert (`sed 's/v<alt>/vX/g' <Datei> > <Kopie>`, je
   Stand). Jede inhaltliche Änderung wird gegen das Repo-Dokument gehalten, das
   aus der geänderten Vorlage oder Regel abgeleitet ist (Skill, Agent-Datei,
   Briefing, README, Konventionen): „übernommen in `<Datei>`“ oder „betrifft das
   Repo nicht“ mit Grund. Dazu der Durchgang durch die Adaptionen in
   `harness/conventions.md` mit ihren fünf Ausgängen.
2. **Stichprobe gegen den Bestand.** Unabhängig vom Delta, weil eine nie
   übernommene, seither unveränderte Vorlagen-Klausel kein Delta erzeugt. Für
   die aus Vorlagen abgeleiteten Repo-Dokumente (`AGENTS.md`, `harness/README.md`,
   `harness/conventions.md`, `docs/plan/planning/README.md`,
   `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/adr/README.md`,
   `docs/plan/carveouts/`) drei Prüfungen:
   - **Gliederung** gegen die Vorlage:
     `diff <(grep -E '^#{1,4} ' <Vorlage>) <(grep -E '^#{1,4} ' <Datei>)`.
   - **Stehengebliebene Platzhalter**, mit zwei Mustern, jede Trefferzeile
     beurteilt: die Platzhalter der Vorlage selbst,
     `grep -n -F -f <(grep -o -E '<[^<>]+>' <Vorlage> | sort -u) <Datei>`, und
     die Formen, die eine ältere Vorlage hinterlassen hat,
     `grep -n -E '<…>|<z\. B\.|<zuerst|<dann|<bei Bedarf|<mover>|<messung>|<vorschau>|<make-target>|<Pfad oder URL>|<Datum>|<repo-spezifischer|<NNN>|<Pfade zu' <Datei>`.
   - **Voller `diff`** gegen die versions-normalisierte Vorlage, jede
     Abweichung beurteilt, für jedes Dokument, dessen Vorlage im Delta aus
     Schritt 1 steht oder von einer Regel des Deltas berührt ist, und immer für
     `docs/plan/planning/README.md`: ihre Zeilen sind feste Klauseln der
     Vorlage, eine fehlende Klausel ändert weder Gliederung noch Platzhalter.

   Grenze: Diese Stichprobe läuft nur beim Bump. Die Stichprobe bei aktuellem
   Pin — je Baseline-Regel, ob sie im ausgefüllten Artefakt steht — trägt
   weiter der Drift-Audit nach `AGENTS.md` §1.
3. **Ergebnis je Dokument.** „entspricht“ oder „Abweichung mit Beleg“; jede
   Abweichung wird im Bump-Slice behoben oder als Folge-Slice benannt. Bestehende
   Instanzen wiederkehrender Vorlagen (ADR, Slice, Welle, Review-Report, `MR`)
   und die Records unter `done/` und `docs/reviews/` werden nicht rückwirkend
   umgeschrieben; neue Instanzen folgen der neuen Form.

## `make pin-stale-actions`

Upstream-Pin-Freshness P9 ([`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Pin-Inventar): jede SHA-gepinnte `uses:`-Zeile über alle `.github/workflows/*.yml` (`AGENTS.md` §3.8) gegen zwei Achsen — Tag-Mutation (`git ls-remote` gegen den im Kommentar genannten Tag; zeigt er noch auf denselben SHA?) und Tag-Frische (neuester Release des Action-Repos über die GitHub-Releases-API). Ein doppelt referenziertes Repo@Tag über mehrere Workflow-Dateien wird nur einmal geprüft

**Bindung:** kein Gate, [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz
