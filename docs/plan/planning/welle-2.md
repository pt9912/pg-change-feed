# Welle welle-2: Reale PostgreSQL-Integration — Store, Stream, E2E

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<NN>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** M1 (MVP-Abnahme) — diese Welle liefert den MVP-Integrationstest-Beleg; der Meilenstein schließt mit seinem grünen Lauf.

**Verantwortlich:** pt9912. **Datum:** 2026-09-09.

---

## 1. Welle-Ziel

<!-- BEDIENHINWEIS: Eine Aussage, die sich an einem Lasttest oder
Akzeptanzkriterium spiegelt. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die MVP-Integrationstest-Abnahme ist belegt: PostgreSQL startet in der
Compose-Umgebung, CDC wird je Tabelle aktiviert, INSERT/UPDATE/DELETE
werden Ende-zu-Ende erfasst und deterministisch gelesen
([`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md) am realen Treiber) — der
MVP-Schnitt ([`spec/lastenheft.md` §1](../../../spec/lastenheft.md)) ist belegt;
der Meilenstein M1 kann erreicht werden.

## 2. Trigger (Welle startet)

<!-- BEDIENHINWEIS: Was muss vorher passiert sein? -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- welle-1 done: alle drei Slices in `done/`, `make gates` grün
  (Beleg: `docs/plan/planning/done/welle-1-results.md`, Verifikation) —
  **bereits eingetreten**.>

## 3. Closure-Trigger (Welle schließt)

<!-- BEDIENHINWEIS: Aktion, nicht Termin. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- <z.B. Alle Slices done.>
- <z.B. `make fullbuild` grün.>
- <z.B. Replay-Lauf gegen Golden Set durchläuft.>
- <z.B. Closure-Notiz in `welle-<NN>-results.md`.>

## 4. Slices in dieser Welle

<!-- BEDIENHINWEIS: keine Status-Spalte ergaenzen. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-004 | PostgreSQL-ChangeStore-Adapter | [`LH-QA-REL-001`](../../../spec/lastenheft.md), [`LH-FA-RET-001`](../../../spec/lastenheft.md) |
| slice-005 | Replication-Stream-Adapter (pgoutput, real) | [`LH-FA-CAP-001`](../../../spec/lastenheft.md)…003 |
| slice-006 | Docker-Compose-Umgebung und MVP-Integrationstest | [`LH-QA-POR-003`](../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

<!-- BEDIENHINWEIS: Falls jemand diese Welle aendert — was bricht? -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Meilenstein M1 (MVP-Abnahme) — der Integrationstest-Beleg
  dieser Welle ist sein erster Träger.
- Wird blockiert von: nichts — der Start-Trigger ist mit welle-1
  eingetreten.

## 6. Out-of-Scope für diese Welle

<!-- BEDIENHINWEIS: explizite Nicht-Inhalte, schuetzt vor Scope-Creep. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Performance-Benchmarks mit Laststufen
  ([`LH-QA-PER-002`](../../../spec/lastenheft.md)) — das Benchmark-Design folgt
  nach dem grünen E2E; Messmethode-Delegation bleibt gesetzt
  ([`SPEC-014`](../../../spec/pflichtenheft.md)).
- Consumer-Verwaltung, Retention-Adapter, SQL-/CLI-Adapter — spätere
  Wellen ([`ADR-0028`](../../../docs/plan/adr/README.md)-Rest; nicht MVP).
- `codepaths`-Aktivierung — geprüft und aktiviert, sobald die
  referenzierten Pfade existieren (Bedingung in `.d-check.yml`); hier
  geprüft, nicht zugesagt.
- `image-cve`-Sensor — advisory; Aktivierungsbedingung eingetreten,
  Träger geplant (Steering-Loop-Kandidat).

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

Ergebnis: <Zeiger auf `welle-<NN>-results.md`, Geschwister im Ruheort `done/`>
Zähler: <Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort>
