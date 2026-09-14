# Review-Report: slice-064 — 2026-09-14

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `slice-064` — Implementierungs-Commit `31e0f60`, Elter
`7a77a54` (reiner `next→in-progress`-Move). `git diff 7a77a54 31e0f60
--stat` bestätigt fünf geänderte Dateien: `.github/workflows/e2e.yml`,
`Makefile`, `compose.yaml`, `docs/plan/planning/in-progress/slice-064-…md`
(DoD-Nachzug), `harness/README.md` — kein Griff in eine Datei außerhalb der
im Slice-Plan §3 genannten.

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, geschärft
2026-09-13: vier repo-spezifische HIGH-Regeln plus Slice-/Wellen-Chronik-
und Handbuch-Versionshistorie-Regel)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-14

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-064-postgresql-versionsmatrix-e2e.md`
  (vollständig gelesen, inkl. §1 Out-of-Scope, §6 Risiken, §8)
- `docs/plan/adr/0058-testansatz-fuenf-luecken.md` (vollständig gelesen,
  bindend für dieses Review: Entscheidung 4 — Matrix-Job, Parametrisierung,
  Betroffene Dateien, Re-Evaluierungs-Trigger 4)
- `docs/plan/adr/0051-cicd-pipeline-github-actions.md` (als Kontext —
  `ci.yml`/`e2e.yml`-Rollenteilung, nicht-blockierend)
- `spec/pflichtenheft.md` §3 (`SPEC-012`, `PG_MAJOR_VERSIONS = 17, 18`)
- `AGENTS.md` §3 Hard Rules (insb. 3.1, 3.3, 3.7, 3.8, 3.9), §5
  Dokumentations-Regeln
- `harness/conventions.md` (MR-000 ID-Schema)
- vorherige Findings am gleichen Modul: `docs/reviews/review-slice-063.md`
  (letzter Diff in `harness/README.md`/Workflow-Umfeld, keine offenen
  HIGH/MEDIUM)

---

## Findings

Keine HIGH- oder MEDIUM-Findings. Ein LOW-Finding zur Commit-Betreff-Form.

### F-1 — Commit-Betreff weicht vom repo-etablierten Conventional-Commit-Muster ab

- `kategorie`: LOW
- `quelle`: Maintainability (kein Hard-Rule-Verstoß — `AGENTS.md` §5 und
  `harness/README.md` §Traceability rules verlangen nur mindestens eine
  `LH-*`- oder `ADR-*`-Kennung und kein `SPEC-*`/`ARC-*` im Betreff; ein
  Präfix-Format wie `feat(...):`/`docs(...):` ist **nirgends** verbindlich
  vorgeschrieben, auch nicht in `tools/harness/commit-traceability.sh` —
  das Skript prüft ausschließlich auf Struktur-ID-Abwesenheit, kein
  Typ-Präfix)
- `pfad`: Commit-Betreff `31e0f60`: `LH-QA-POR-001: PostgreSQL-Versionsmatrix im e2e.yml-Workflow (ADR-0058)`
- `befund`: Die letzten 40 Commits dieses Repos (`git log --oneline -40`)
  folgen durchgängig dem Muster `<typ>(<scope>): <text> (<IDs>)` (`feat`,
  `fix`, `docs`, `test`). Dieser Commit beginnt stattdessen direkt mit der
  `LH-*`-Kennung als Präfix vor dem Doppelpunkt, ohne Conventional-Commit-
  Typ/Scope. Mechanisch unauffällig (`make commit-traceability` lief grün:
  „Betreffs ohne Struktur-ID"), aber eine Abweichung vom durchgängig
  beobachtbaren Repo-Stil.
- `verifizierbar`: ja — `git log --oneline -40` zeigt das etablierte Muster
  in jedem Vorgänger-Commit; `make commit-traceability` bestätigt, dass die
  Traceability-*Hard Rule* selbst nicht verletzt ist (nur der Stil weicht
  ab).
- `klasse`: „Commit-Betreff-Stil weicht von Conventional-Commits ab"

## Negativbefunde

- geprüft, ohne Befund: `compose.yaml` — Interpolation
  `${PG_TEST_IMAGE:-postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8}`
  syntaktisch valide (`docker compose config` selbst ausgeführt, siehe
  Gate-Läufe unten) und der Default ist **exakt** der bisherige
  PostgreSQL-18-Digest — Zeichen-für-Zeichen identisch mit der vorherigen
  literalen `image:`-Zeile (`git diff 7a77a54 31e0f60 -- compose.yaml`
  zeigt keine Digest-Änderung, nur die Interpolationshülle). `docker
  compose config` ohne gesetzte Variable liefert
  `image: postgres:18-alpine@sha256:63bdc97…`, mit
  `PG_TEST_IMAGE=postgres:17-alpine@sha256:7456ef82…` gesetzt liefert es
  exakt diesen 17er-Digest — Regressionsverhalten wie im DoD gefordert
  bestätigt.
- geprüft, ohne Befund: `Makefile` — Diff ist ausschließlich ein
  Kommentar-Zusatz (vier neue Zeilen Prosa), `PG_TEST_IMAGE ?=
  postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8`
  unverändert (Grep bestätigt Zeile 51 identisch zum Vorzustand). Keine
  funktionale Änderung, wie vom Implementer berichtet.
- geprüft, ohne Befund: PostgreSQL-17-Digest —
  `sha256:7456ef82e5f5bc43d997f4781bbd7c0d6389bff397564649a356e206ba473aee`
  selbst über `docker manifest inspect postgres:17-alpine` gegen die
  Registry geprüft (nicht nur die Implementer-Behauptung übernommen): der
  `amd64`/`linux`-Eintrag im Manifest-Index trägt exakt diesen Digest.
  Übereinstimmung bestätigt.
- geprüft, ohne Befund: `.github/workflows/e2e.yml`s `strategy: matrix:` —
  beide `uses:`-fremden Digest-Zeilen sind SHA-gepinnt mit Tag-Kommentar
  (`postgres:17-alpine@sha256:…` und `postgres:18-alpine@sha256:…`, je mit
  erklärendem Kommentar darüber) — `AGENTS.md` §3.8 ist wörtlich für
  `uses:`-Zeilen formuliert (GitHub Actions), gilt hier sinngemäß für
  Image-Pins und ist eingehalten. YAML-Syntax real validiert (`python3
  yaml.safe_load` parst den vollständigen Workflow fehlerfrei;
  `yamllint` meldet nur eine vorbestehende, unveränderte Warnung auf Zeile
  84, die bereits im Elter-Commit `7a77a54` vorlag — kein neuer Befund).
- geprüft, ohne Befund: `fail-fast: false` — im Slice-Plan/`ADR-0058` nicht
  explizit genannt, aber vom Implementer im Workflow-Kommentar begründet
  (verhindert, dass ein rotes Leg das andere abbricht; beide
  Matrix-Ergebnisse bleiben sichtbar). Die Begründung trägt: `e2e.yml` ist
  laut `ADR-0051` ohnehin nicht-blockierend, ein abgebrochenes zweites Leg
  würde bei einer echten Versionsdrift (`slice-064` §6, Risiko 3) genau die
  Diagnose-Information verlieren, die der Matrix-Job liefern soll. Kein
  Widerspruch zu ADR/Hard Rule, keine Scope-Ausweitung über den
  CI-Workflow-Bereich hinaus — bewertet als legitime, begründete
  Implementer-Entscheidung innerhalb des Slice-Umfangs, kein
  Reviewer-Finding.
- geprüft, ohne Befund: Job-Ebene `env: PG_TEST_IMAGE: ${{
  matrix.pg_test_image }}` sitzt auf Job-Ebene (`jobs.e2e.env`, nicht auf
  einem einzelnen Step) und wird von **beiden** bestehenden `run:`-Steps
  (`make image`, `make test-integration`) ohne eigene `env:`-Blöcke geerbt
  — GitHub-Actions-Vererbungsregel: Step ohne eigenes `env:` erbt das
  Job-`env:` vollständig. Bestätigt durch den geparsten YAML-Baum (siehe
  oben): beide Steps tragen nur `name`/`run`, kein überschreibendes `env:`.
- geprüft, ohne Befund: `harness/README.md` — die `e2e.yml`-Zeile trägt
  jetzt `LH-QA-POR-001`, `ADR-0058`-Verweis und den `· seit slice-064`-
  äquivalenten Zusatz „Seit slice-064 trägt der Job eine `strategy:
  matrix:`…" korrekt als Indikativ-Beschreibung des aktuellen Zustands
  (keine Konjunktiv-/Chronik-Sprache, `AGENTS.md` §3.7 eingehalten).
- geprüft, ohne Befund: Kommentar-Disziplin in `e2e.yml`/`compose.yaml`/
  `Makefile` (`AGENTS.md` §3.7) — alle drei neuen Kommentarblöcke
  beschreiben den aktuellen Zustand indikativ, tragen die Klassen Zusage/
  Kopplung (`ADR-0058`/`LH-QA-POR-001`-Bezug), keine Slice-/Wellen-Chronik als
  alleinige Begründung, kein abgebrochener Satz, kein Verweis auf eine
  verworfene Alternative im Konjunktiv. Der Datumszusatz „ermittelt
  2026-09-14" folgt demselben bereits etablierten Muster wie der
  bestehende Makefile-Kommentar „(`go1.27.1`, real geprüft)" — kein neues
  Muster, keine Chronik.
- geprüft, ohne Befund: Scope-Fidelity — `git diff 7a77a54 31e0f60
  --name-only` zeigt ausschließlich die im Slice-Plan §3 genannten
  Dateien plus die Slice-Plan-Datei selbst (DoD-Nachzug). Keine Berührung
  von `tools/harness/run-integration-tests.sh` oder `.github/workflows/ci.yml`
  — beide ausdrücklich in §1 Out-of-Scope als Gegenstand anderer Slices
  genannt.
- geprüft, ohne Befund: Traceability — Commit trägt `LH-QA-POR-001` und
  `ADR-0058`, keine `SPEC-*`/`ARC-*`-Kennung im Betreff; `make
  commit-traceability` (Teil von `make gates`) bestätigt beide Grenzen
  mechanisch. Stilistische Abweichung vom Conventional-Commit-Präfix siehe
  F-1 (LOW).
- geprüft, ohne Befund: `docker manifest inspect`-Ergebnis für
  PostgreSQL 18 — die in `e2e.yml` referenzierte 18er-Zeile trägt denselben
  Digest wie der bestehende Makefile-/Compose-Default
  (`sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8`),
  keine unabsichtliche Divergenz zwischen den drei Stellen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Commit-Betreff-Stil weicht von
Conventional-Commits ab

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 1 LOW-Finding (stilistisch,
keine semantische Auswirkung, keine Hard-Rule-Verletzung). Keine Fixrunde
zwingend nötig; F-1 kann ohne Rückgabe an den Implementer zur Kenntnis
genommen werden.

**Übergabe:** Da dieser Reviewer-Lauf zu 0 HIGH/MEDIUM kommt und keine
Reviewer→Implementer-Rückgabe zwingend stattfindet, wird die DoD-Zeile
„Review durchgeführt, Report unter `docs/reviews/` liegt vor" im Slice-Plan
(`docs/plan/planning/in-progress/slice-064-postgresql-versionsmatrix-e2e.md`)
im selben Commit, der diesen Report anlegt, auf `[x]` nachgezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). Dieser
Report ist ein **Lauf-Beleg** und wird über Läufe hinweg nicht wieder
gelesen. Er ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11).

**Offener Punkt für die Verifikation (kein Finding):** Ein realer
GitHub-Actions-Matrix-Lauf über beide Legs ist von diesem Review-Lauf aus
nicht ausführbar (Commit ist noch nicht Teil eines gepushten, ausgelösten
Laufs, siehe `BEO-PGC/github-actions-unverifizierbar-lokal`, 2× vor diesem
Slice). Das prüft der Planner-Koordinator separat nach dem Push.

## Gate-Läufe dieses Reviews

Beide ungefiltert und mit unmittelbar geprüftem eigenem Exit-Code
(`AGENTS.md` §3.9) — kein Piping/Wrapper zwischen Lauf und
Exit-Code-Prüfung:

- `make gates` — **Exit 0.** Alle inneren Gates grün (`baseline-verify`:
  54 Dateien OK; `docs-check`/d-check: 510 Dateien, 0 Befunde;
  `commit-traceability`: 5 Commits OK, Betreffs ohne Struktur-ID;
  `a-check`: 0 Befunde; `coverage-gate`: 44.40 % über Schwelle 35 %).
- `docker compose config` (zweimal, mit und ohne `PG_TEST_IMAGE`) —
  beide Male Exit 0, Digest-Ausgabe wie in den Negativbefunden dokumentiert.
- `docker manifest inspect postgres:17-alpine` — Exit 0, `amd64`-Digest
  stimmt mit dem in `e2e.yml` gepinnten Wert überein.
- `python3 -c "yaml.safe_load(...)"` gegen `.github/workflows/e2e.yml` —
  Exit 0, vollständig geparst.
