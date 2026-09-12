# Welle 9: E2E-Abdeckung — CDC-Kernpfad

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-9-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-12.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Der Compose-Integrationstest (`test/integration/integration_test.go`,
umbenannt in `slice-028`) prüft Erfassung/Lesen bislang **white-box**: er
schreibt Quelländerungen real extern (SQL-Insert), liest die Changes aber
über den intern instanziierten Go-Adapter
(`postgresstorage.NewChangeStoreAdapter` → `ReadChanges`), nicht über die
tatsächlich ausgelieferte externe SQL-Sicht `cdc.changes`
(`LH-FA-REA-002`). Dieselbe Lücke betrifft Schema-Änderungen
(`LH-FA-SCH-*`, ALTER-TABLE-Szenarien) — bislang nicht im Compose-
Integrationstest abgedeckt.

`BEO-PGC/lese-doppelquelle` (2×) benennt dazu passend: Die SQL-Views im
`cdc`-Schema und die Go-Use-Case-Lesepfade implementieren dieselbe
Lese-Semantik zweimal, ohne dass ein Sensor/Vertrags-Test beide
gegeneinander hält. Ein Black-Box-Test, der über die SQL-Sicht statt über
den Go-Adapter liest, ist genau der fehlende Vertrags-Test — diese Welle
ist die Gelegenheit, ihn zu bauen.

**Das *Mehr* gegenüber den einzelnen Slice-DoDs:** Kein Slice-DoD allein
beweist, dass Erfassung, Lesen über die externe SQL-Sicht und
Schema-Änderungs-Behandlung zusammen — ausschließlich über extern
ausgelieferte Schnittstellen, nicht über interne Go-Pakete — real
funktionieren, und dass SQL-View- und Go-Adapter-Lesepfad dieselbe Change-
Reihenfolge und denselben Inhalt liefern.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `welle-8` liegt in `done/`.
- Kein Slice liegt in `in-progress/`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices (`slice-029`, `slice-030`) liegen in `done/`.
- `make gates` grün.
- Der Compose-Integrationstest liest Changes real über die externe
  SQL-Sicht `cdc.changes` (nicht mehr ausschließlich über den internen
  Go-Adapter) und belegt dabei, dass SQL-View- und Go-Adapter-Lesepfad
  dieselbe Reihenfolge/denselben Inhalt liefern.
- Ein Schema-Änderungs-Szenario (mindestens `ALTER TABLE ADD COLUMN`) läuft
  real gegen den Compose-Stack und wird korrekt erfasst/gelesen.
- Closure-Notiz in `welle-9-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-029 | Black-Box-Lesepfad über `cdc.changes` — Vertragstest gegen `BEO-PGC/lese-doppelquelle` | [`LH-FA-REA-002`](../../../spec/lastenheft.md) |
| slice-030 | Black-Box-E2E für Schema-Änderungen | [`LH-FA-SCH-001`](../../../spec/lastenheft.md)…`005` |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle.
- Wird blockiert von: keine andere Welle.
- Intern: `slice-029` und `slice-030` sind unabhängig voneinander
  (unterschiedliche Testfälle im selben Compose-Integrationstest); keine
  erzwungene Reihenfolge.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **CDC-Verwaltung (`LH-FA-CFG-*`) und Administration/Observability
  (`LH-FA-ADM-*`)** — bereits als eigene, spätere Welle vorgemerkt
  (Roadmap *Nächste Wellen*, „E2E-Abdeckung — Verwaltung & Observability").
- **Consumer-Zugriffsweg und Rollen-DSN-Trennung** — bereits real
  black-box abgedeckt (`welle-8`); diese Welle baut nicht erneut darauf.
- **`BEO-PGC/lese-doppelquelle` vollständig auflösen** (Ausgang
  verkörpert/gestrichen) — anderer Vorgang: Diese Welle liefert den
  fehlenden Vertrags-Test, der Zähler-Stand und der endgültige Ausgang
  werden bei der Closure bewertet, nicht vorab festgelegt.
- **Retention, Performance-Benchmarks** — bereits als eigene, spätere
  Wellen vorgemerkt.

## 7. Closure-Notiz

<!--
BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md §Verwendung,
Schritt 5) und darf deshalb nichts Tragendes halten.

- Erst nach Welle-Abschluss fuellen; nur die Nummer, nicht die volle Welle-ID.
- Ziel-Form der Ergebnis-Notiz: `welle-results.template.md` — Schwester-Vorlage
  im Template-Verzeichnis, kein Artefakt deines Repos. Sie ist von der
  Ruheort-Regel ausgenommen und faellt mit diesem Kommentar ohnehin weg.
-->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: [welle-9-results.md](welle-9-results.md)
Zähler: [../observations/](../observations/)
