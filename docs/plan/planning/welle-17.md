# Welle 17: E2E-Testbelege für fünf testfreie Lastenheft-Kennungen (`ADR-0058`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-17-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** —. **Datum:** 2026-09-14.

---

## 1. Welle-Ziel

Fünf Lastenheft-Kennungen ohne jeden Testbeleg (`LH-FA-SCH-003`,
`LH-FA-DAT-006`, `LH-QA-OPS-005`, `LH-QA-POR-001`, `LH-QA-POR-002`) bekommen
je den in [ADR-0058](../adr/0058-testansatz-fuenf-luecken.md) entschiedenen,
konkreten Testansatz umgesetzt — zwei neue E2E-Testfunktionen am laufenden
Feed-Container, eine neue Orchestrierungs-Phase für einen simulierten
Upgrade-Zyklus, eine PostgreSQL-Versionsmatrix im nicht-blockierenden
E2E-Workflow und eine Plattform-Assertion im blockierenden Gate-Workflow.
`ADR-0058` (Accepted) trifft bereits alle fünf Teilentscheidungen inklusive
Alternativen-Vergleich; diese Welle setzt sie in vier Slices um.

## 2. Trigger (Welle startet)

- `ADR-0058` liegt Accepted vor (bereits erfüllt — Architect-Entscheidung
  vom 2026-09-14).

## 3. Closure-Trigger (Welle schließt)

- `slice-062`, `slice-063`, `slice-064`, `slice-065` liegen in `done/`.
- `make gates` grün.
- Ein real belegter grüner `e2e.yml`-Matrix-Lauf über beide PostgreSQL-Legs
  (17, 18) nach einem echten Push auf den Hauptzweig — das *Mehr* gegenüber
  den einzelnen Slice-DoDs: `slice-064` allein parametrisiert nur
  `compose.yaml`/`Makefile`/`e2e.yml`, kann den vollständigen E2E-Lauf aber
  nicht grün belegen, solange `slice-062`/`slice-063` (dieselben
  Testfunktionen/Phasen laufen in jedem Matrix-Leg mit) noch nicht in
  `done/` liegen — alle E2E-Testfälle laufen im selben Compose-Stack-Lauf.
- Closure-Notiz in `welle-17-results.md`.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| slice-062 | E2E-Testfälle Schema-Verhalten & Struktur — entfernte Spalten, Metadaten-Erweiterbarkeit | [LH-FA-SCH-003](../../../spec/lastenheft.md), [LH-FA-DAT-006](../../../spec/lastenheft.md) |
| slice-063 | Upgrade-Sicherheit — simulierter Container-Stopp/Rollout/Start-Zyklus | [LH-QA-OPS-005](../../../spec/lastenheft.md) |
| slice-064 | PostgreSQL-Versionsmatrix im E2E-Workflow | [LH-QA-POR-001](../../../spec/lastenheft.md) |
| slice-065 | Linux-Plattform-Assertion im Gate-Workflow | [LH-QA-POR-002](../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

- `slice-064` (Versionsmatrix) prüft real erst dann beide PostgreSQL-Legs
  vollständig grün, wenn `slice-062`/`slice-063` bereits in `done/` liegen —
  ohne sie fehlen zwei der fünf Testbelege in jedem Matrix-Leg. Kein
  Implementierungs-Blocker (die drei Slices ändern getrennte Dateien), aber
  ein Beleg-Abhängigkeit für den Welle-Closure-Trigger.
- `slice-065` ist unabhängig von den anderen drei (eigener Workflow,
  `ci.yml` statt `e2e.yml`).
- Blockiert: keine andere offene Welle (`Nächste Wellen` ist derzeit leer).
- Wird blockiert von: keiner (`ADR-0058` liegt bereits Accepted vor).

## 6. Out-of-Scope für diese Welle

- **Echter Zwei-Image-Versionsvergleich für `LH-QA-OPS-005`** —
  `ADR-0058` Entscheidung 3/Re-Evaluierungs-Trigger 3 verschiebt das
  ausdrücklich auf den Zeitpunkt, sobald eine echte Release-Historie
  existiert (`slice-040`, Release-Pipeline/Tags — noch nicht angelegt, siehe
  [`ADR-0051`](../adr/0051-cicd-pipeline-github-actions.md)); `slice-039`
  (CI-Workflow/Dependabot) liegt bereits in `done/`. Diese Welle bildet nur
  den Mechanismus (Stopp/Rollout/Start) nach, keinen echten
  Versionswechsel.
- **Lokales Makefile-Default für `PG_TEST_IMAGE`** — `slice-064` ändert die
  Compose-Interpolation und den CI-Matrix-Wert, nicht den
  `Makefile`-Default für lokale/manuelle Läufe (`ADR-0058` Entscheidung 4);
  ein Wechsel des lokalen Default-Digests bleibt ein eigener, bewusster
  Pin-Hebungs-Commit (`AGENTS.md` §3.8).
- **Plattform-Assertion in `e2e.yml`** — `slice-065` fasst ausschließlich
  `ci.yml` an; `ADR-0058` Entscheidung 5 bindet die Messmethode bewusst an
  den blockierenden Gate-Workflow, nicht an den langsameren,
  nicht-blockierenden Compose-Stack.
- **Änderung der `SPEC-012`-Versionsliste selbst** (`PG_MAJOR_VERSIONS`) —
  diese Welle setzt die bereits in `SPEC-012` festgelegten zwei Versionen
  (17, 18) technisch um, ändert die Liste nicht. Eine künftige Ausweitung
  (z. B. PostgreSQL 19) bleibt eine eigene Pflichtenheft-Änderung.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: <Zeiger auf `welle-17-results.md`, Geschwister im Ruheort `done/`>
Zähler: <Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort>
