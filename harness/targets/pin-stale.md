# `make pin-stale-*` (P3 bis P9) — Upstream-Pin-Freshness

Ausführliche Fassung der Index-Zeilen aus [`harness/README.md` §Sensors](../README.md#sensors-feedback-gates); die Zeile dort trägt einen Satz und verlinkt hierher.

## `make pin-stale-race`/`make pin-stale-pgtest`/`make pin-stale-dmigrate`/`make pin-stale-acheck`

Upstream-Pin-Freshness P3–P6 ([`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Pin-Inventar): je Achse ein Digest-Vergleich einer Makefile-/`a-check.mk`-Variable (`TOOLCHAIN_RACE_IMAGE`, `PG_TEST_IMAGE`, `D_MIGRATE_IMAGE`, `A_CHECK_IMAGE`) gegen den aktuellen Registry-Digest, über das gemeinsame `tools/harness/pin-stale.sh` (analog `make image-stale`, aber Quelle ist eine Variable statt einer Dockerfile-`FROM`-Zeile); P5/P6 sind reine Digest-Pins ohne eigenen Tag, deshalb Vergleich gegen den jeweiligen `:latest`-Tag der GHCR-Pakete. Je Achse eine Zeile `OK`, `DRIFT` oder `UNBESTIMMT` mit Variable, Tag und Digest. Gemessen 2026-10-05 mit den vier Zielen einzeln, je Exit 0 und gedruckte Zeile `OK` für `TOOLCHAIN_RACE_IMAGE` (`golang:1.27`), `PG_TEST_IMAGE` (`postgres:18-alpine`), `D_MIGRATE_IMAGE` und `A_CHECK_IMAGE` (je `:latest`); das Ergebnis ist der Stand dieses Laufs und bewegt sich mit jedem Neubau eines Upstream-Tags

Ein Bump von `A_CHECK_IMAGE` folgt [§Bump eines Gate-Werkzeugs](#bump-eines-gate-werkzeugs-reichweite-vor-dem-pin-commit).

**Bindung:** kein Gate, [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz

## `make pin-stale-dcheck`

Upstream-Pin-Freshness P7 ([`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Pin-Inventar): zwei Achsen für `d-check.mk`s `DCHECK_IMAGE`/`DCHECK_DIGEST` — Digest-Drift (trägt der gepinnte Tag noch denselben Bau?) und Tag-Frische (existiert ein neuerer `pt9912/d-check`-Release als der gepinnte Tag, über die GitHub-Releases-API?)

**Bindung:** kein Gate, [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz

### Bump eines Gate-Werkzeugs: Reichweite vor dem Pin-Commit

Hebt ein Slice den Pin eines Werkzeugs, dessen Lauf in `make gates` blockiert —
d-check (`DCHECK_IMAGE`/`DCHECK_DIGEST` in `d-check.mk`), a-check
(`A_CHECK_IMAGE`) —, liest er vor dem Pin-Commit das Changelog jeder
übersprungenen Version auf Änderungen dessen, was ein aktives Modul meldet oder
durchlässt. Jede Änderung wird gegen die ADR gehalten, die die Reichweite der
betroffenen Regel trägt. Eine **Erweiterung** — das Gate meldet eine Klasse, die
keine Regel verbietet — bekommt vor dem Pin-Commit eine ADR oder ein Verdikt des
Architect; ohne dieses Artefakt landet der Pin nicht. Der Plan des Bumps trägt je
Version „keine Reichweiten-Änderung“ oder die Änderung mit ihrem Artefakt. Null
Befunde am Arbeitsbaum belegen keine gleiche Reichweite: sie sagen nur, dass der
Bestand die neue Klasse heute nicht trifft. Herkunft:
`BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` · seit slice-dcheck-v0-82-0.

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
   `harness/conventions.md` mit ihren fünf Ausgängen; je MR-Eintrag misst er
   die Einheit seines Verweises in `Ersetzt-Baseline-Regel` **zwischen den
   Tags des Bumps** mit `make zitat-vergleich`; Stände, Normalisierung und die
   Lesung der Ausgänge stehen in
   [`SPEC-038`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge),
   ein Prüfauftrag bekommt seinen Ausgang mit Grund im Plan des Bumps
   ([`ADR-0161`](../../docs/plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
   Entscheidung 5). Dazu die getrackten
   Symlinks auf den alten Tag — `git grep` liest kein Symlink-Ziel (Modus
   120000), `make docs-check` ebenfalls nicht:
   `git ls-files -s | awk '$1==120000{print $4}' | while read -r p; do printf '%s -> %s\n' "$p" "$(readlink "$p")"; done | grep -F '/<alt>/'`
   — jeder Treffer wird auf den neuen Tag umgestellt; nach dem Löschen des alten
   Verzeichnisses meldet `find . -path ./.git -prune -o -xtype l -print` keinen
   hängenden Symlink · seit slice-harness-baseline-v6-14-1
   (`BEO-PGC/bump-ablauf-ohne-symlink-ziele`).
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

**Pins in ADRs und MR-Einträgen bleiben stehen.** Baseline-Pins in
`docs/plan/adr/[0-9]*.md` und `harness/conventions/MR-[0-9]*.md` nennen den
Stand ihrer Abfassung; ein Bump zieht sie nicht nach, `versions` nimmt beide
Pfade aus, `links` und `anchors` nicht
([`ADR-0161`](../../docs/plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
Entscheidung 1).

**Form-Commit vor dem Löschen.** Zeigt ein Link in einem MR-Eintrag in das
Tag, das der Bump löscht, wird er einmal zu Inline-Code (Pfad samt alter
Version und Anker, der Linktext bleibt Prosa). Das steht in einem eigenen
Commit `F`, der **nur** MR-Dateien ändert, dessen Message `ADR-0073` und
`ADR-0161` nennt und der **vor** dem Lösch-Commit liegt (Entscheidung 4). Der
Verifier prüft `make doc-immutable` in den Teil-Ranges `B..F~1` und `F..H`,
je mit dem Leer-Test der Teil-Range
([`SPEC-039`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge);
Befehlsform: Funktion `teilrange` in
[`ADR-0160`](../../docs/plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
Entscheidung 1), und am Form-Commit je MR-Datei den `formnorm`-`cmp`
aus `ADR-0161` Entscheidung 4: eine Zeile je Datei, jede mit `cmp 0`. Ein
Link aus einer ADR in ein gelöschtes Tag bekommt dieselbe Form-Korrektur in
einem eigenen Commit mit `ADR-0073` und einer §Geschichte-Zeile.

## `make pin-stale-actions`

Upstream-Pin-Freshness P9 ([`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Pin-Inventar): jede SHA-gepinnte `uses:`-Zeile über alle `.github/workflows/*.yml` (`AGENTS.md` §3.8) gegen zwei Achsen — Tag-Mutation (`git ls-remote` gegen den im Kommentar genannten Tag; zeigt er noch auf denselben SHA?) und Tag-Frische (neuester Release des Action-Repos über die GitHub-Releases-API). Ein doppelt referenziertes Repo@Tag über mehrere Workflow-Dateien wird nur einmal geprüft

**Bindung:** kein Gate, [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz
