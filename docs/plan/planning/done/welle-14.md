# Welle 14: Performance-Benchmarks & Test-Coverage-Gate

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-14-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-13.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

`pg-change-feed` hat aktuell **keinen** Test-Coverage-Beleg und **keine**
systematische Performance-Mess-Infrastruktur: [LH-QA-PER-001](../../../../spec/lastenheft.md)…`003`
(Quell-Impact mit/ohne CDC, Skalierbarkeit über die in
[`SPEC-014`](../../../../spec/pflichtenheft.md) fixierten Lastenstufen,
Batch- vs. Einzelabruf-Effizienz) haben keinen einzigen Beleg außer dem
bereits bestehenden, punktuellen `cdc_capture_lag`-Lasttest (der nur
`LH-QA-PER-004` teilweise deckt). `AGENTS.md` §3.2 (Suppression-Verbot)
war bis zu dieser Welle der unausgefüllte Template-Platzhalter.
[`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
hat vorab entschieden: ein scope-eingeschränktes Coverage-Gate
(`internal/...`+`cmd/...`, Endstufe 80 %, bedingte Eskalationsklausel je
nach realem Ist-Stand) und eine Bench-Skript-Familie (drei eigenständige
Belege, kein Gate) — beide nach dem real geprüften Vorbild aus
`/Development/d-check`.

**Das *Mehr* gegenüber den einzelnen Slice-DoDs:** Kein Slice-DoD allein
beweist, dass `make gates` jetzt einen echten Coverage-Beleg trägt **und**
`make bench` alle drei PER-Anforderungen gleichzeitig dokumentiert
belegt — das ist erst die Summe beider Slices.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `slice-047` liegt in `done/`.
- [`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
  liegt vor (Accepted).
- Kein Slice liegt in `in-progress/`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices dieser Welle liegen in `done/`.
- `make gates` grün, inklusive eines real ausgeführten
  `make coverage-gate`-Laufs (Endstufe 80 % oder eine dokumentierte,
  real gemessene Eskalationsstufe gemäß `ADR-0054`).
- `make bench` liefert real alle drei Belege
  (`LH-QA-PER-001`…`003`) mit dokumentiertem Ergebnis.
- Closure-Notiz in `welle-14-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-049 | Test-Coverage-Gate | [ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md) |
| slice-050 | Performance-Benchmark-Infrastruktur | [LH-QA-PER-001](../../../../spec/lastenheft.md), [LH-QA-PER-002](../../../../spec/lastenheft.md), [LH-QA-PER-003](../../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: „Publication-Entzug-Wirksamkeit am laufenden Stream" (in der
  Roadmap direkt nachfolgend eingereiht).
- Wird blockiert von: keine andere Welle.
- Intern: `slice-049` (Coverage-Gate) und `slice-050` (Benchmark-
  Infrastruktur) sind voneinander unabhängig (unterschiedliche Schichten
  — Build-/Gate-Infrastruktur vs. eigenständige Mess-Skripte, `ADR-0054`
  §Empfehlung zum Schnitt) und können parallel oder in beliebiger
  Reihenfolge laufen.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Einführung eines Linters** (`golangci-lint` o. ä.) — `ADR-0054` §(a)
  schließt das ausdrücklich aus: kein Linter existiert, `AGENTS.md` §3.2
  ist bewusst schmal (Suppression-Vollverbot) ausgefüllt; eine künftige
  Linter-Einführung wäre ein eigener, deutlich größerer Vorgang mit
  eigener ADR.
- **`LH-QA-PER-004`-Ausbau** (Commit→CDC-Latenz über eine systematische
  Infrastruktur hinaus) — bereits teilweise über `cdc_capture_lag` und
  den bestehenden Lasttest-Beleg gedeckt; `ADR-0054`s
  Re-Evaluierungs-Trigger (c) benennt einen eigenen Folge-Vorgang, falls
  das je nötig wird.
- **Eine vorab fixierte Ramp-Stufenfolge für das Coverage-Gate** —
  `ADR-0054` entscheidet bewusst gegen eine erfundene Zwischenstufe;
  die Eskalationsstufe (falls nötig) misst der Implementer real beim
  ersten Lauf, nicht diese Welle vorab.

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

Ergebnis: [welle-14-results.md](welle-14-results.md), Geschwister im Ruheort `done/`
Zähler: [../observations/](../observations/)`BEO-PGC/`, eine Ebene über dem Ruheort
