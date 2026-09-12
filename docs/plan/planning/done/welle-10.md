# Welle 10: Schema-Evolution-Nachlieferung (ADR-0015)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-10-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-12.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

`welle-9`s `slice-030` fand real, dass eine Folgepflicht aus
[`ADR-0015`](../../adr/0015-schema-evolution.md) (Accepted, `permanent`,
Option C: `TableSchema`-/`SchemaVersion`-Modelle je Change,
`SchemaStorePort` als Outbound Port, Fehlerklasse `schema` für nicht
sicher interpretierbare Änderungen) nie umgesetzt wurde: der ausgelieferte
Code verhält sich wie das in `ADR-0015` explizit verworfene Option A
(Relation Metadata 1:1 durchreichen, keine stabile historische
Interpretation). Zwei Lastenheft-Akzeptanzkriterien sind dadurch
strukturell unerfüllbar — [`LH-FA-SCH-005`](../../../../spec/lastenheft.md)s
Boundary (zwei Changes vor/nach einer Schemaänderung müssen sich anhand
ihrer Schema-Version unterscheiden lassen) und
[`LH-FA-SCH-004`](../../../../spec/lastenheft.md)s Negative-Fall (eine
inkompatible Typänderung muss erkennbar gemeldet werden, keine stille
Fehlinterpretation). Ein Architect-Verdikt
(`docs/reviews/architect-verdict-slice-030-adr-0015.md`) bestätigte:
`ADR-0015` gilt unverändert fort, die Auflösung ist eine fehlende
Umsetzung nachzuliefern, nicht die ADR zu korrigieren.

**Das *Mehr* gegenüber den einzelnen Slice-DoDs:** Kein Slice-DoD allein
beweist, dass die vier berührten Schichten (Domänen-Modell, neuer
Outbound-Port, neuer Adapter samt DB-Schema-Migration, Decoder/Mapper der
Adapter-Schicht) zusammen real gegen den Compose-Stack funktionieren und
`ADR-0015`s Folgepflicht tatsächlich einlösen — insbesondere, dass die
bereits in `slice-030` geschriebenen Black-Box-Tests
(`TestMVPSchemaChangeAddColumn`, `TestMVPSchemaChangeIncompatibleTypeChange`)
danach ohne Konzeptänderung grün laufen.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `welle-9` liegt in `done/`.
- Kein Slice liegt in `in-progress/`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices dieser Welle liegen in `done/`.
- `make gates` grün.
- `TestMVPSchemaChangeAddColumn` beweist real, dass die `schema_version`
  der nach einem `ALTER TABLE ADD COLUMN` erfassten Changes sich von der
  davor unterscheidet (schließt `LH-FA-SCH-005`s Boundary real).
- `TestMVPSchemaChangeIncompatibleTypeChange` beweist real, dass eine von
  PostgreSQL zugelassene, aber nicht verlustfrei decodierbare
  Typänderung als Fehlerklasse `schema` sichtbar gemeldet wird (schließt
  `LH-FA-SCH-004`s Negative-Fall real).
- `BEO-PGC/schema-evolution-nicht-dynamisch` erreicht Ausgang
  *verkörpert*.
- Closure-Notiz in `welle-10-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-031 | Schema-Persistenz-Fähigkeit — `TableSchema`-Modell, `SchemaStorePort`, Adapter (ohne Live-Verdrahtung) | [`SPEC-004`](../../../../spec/pflichtenheft.md) |
| slice-032 | Dynamische Re-Versionierung im Consume-Pfad | [`LH-FA-SCH-005`](../../../../spec/lastenheft.md) |
| slice-033 | Typ-Auswertung und Fehlerklasse `schema` für inkompatible Typänderungen | [`LH-FA-SCH-004`](../../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: „E2E-Abdeckung — Verwaltung & Observability" (in der Roadmap
  direkt nachfolgend eingereiht — kein inhaltlicher Zusammenhang, reine
  Sequenz-Priorisierung, weil diese Welle eine reale Spec-Nichterfüllung
  schließt).
- Wird blockiert von: keine andere Welle.
- Intern: `slice-031` → `slice-032`/`slice-033` (beide brauchen den
  `SchemaStorePort` aus `slice-031`); `slice-032` und `slice-033` sind
  untereinander unabhängig (verschiedene Zweige in `Assembler.Consume`).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Entfernte Spalten (`LH-FA-SCH-003`)** — bereits als eigener,
  ausgeschlossener Punkt in `slice-030`s §1 benannt; kein Bestandteil
  dieser Welle, eigener Vorgang bei Bedarf.
- **`BEO-PGC/rollen-test-abdeckungsluecken`, Retention, Performance** —
  bereits als eigene, spätere Wellen vorgemerkt.
- **Black-Box-E2E-Erweiterung über die bereits in `slice-030`
  geschriebenen Tests hinaus** — diese Welle macht die bestehenden Tests
  grün, erfindet keine neuen Testszenarien; eine eigene
  „E2E-Abdeckung — Schema-Evolution"-Welle bleibt bei Bedarf ein späterer,
  eigener Vorgang (optional, siehe Architect-Verdikt).
- **`ALTER TABLE … DROP COLUMN` / weitere DDL-Formen** — Bestand bleibt
  bewusst außen vor: Die Architect-Skizze deckt die in `slice-030` bereits
  gebauten Testfälle ab (`ADD COLUMN`, Typänderung); zusätzliche DDL-Formen
  wären ein eigener Fund, kein Bestandteil dieser Nachlieferung.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: [welle-10-results.md](welle-10-results.md)
Zähler: [../observations/](../observations/)
