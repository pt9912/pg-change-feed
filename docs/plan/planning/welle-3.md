# Welle welle-3: CDC-Verwaltung, Lesen-Vollabdeckung, Sicherheit und Observability-Basis

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<NN>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug — M1 (MVP-Abnahme) ist erreicht (welle-2).

**Verantwortlich:** pt9912. **Datum:** 2026-09-09.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die MVP-lückigen Bereiche sind geschlossen: CDC-Verwaltung als Use
Cases (CFG-001…004 am Inbound-Port), Lesen-Vollabdeckung am
Store-Adapter (REA-001…006), Consumer-Verwaltung (CON-001…006),
Sicherheit (Least-Privilege je SEC-001…003) und
Observability-Basis (Health/Metriken-Minimum).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- welle-2 done: alle drei Slices in `done/`, vier Gates inkl.
  commit-traceability grün (Beleg:
  `docs/plan/planning/done/welle-2-results.md`, Verifikation) —
  **bereits eingetreten**.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices dieser Welle liegen in `done/` (slice-008…011).
- `make gates` grün (vier Gates inkl. commit-traceability und a-check)
  **und** `make test-integration` grün am verdrahteten System — der
  welle-spezifische Beleg: die Use-Case-Fähigkeit ersetzt die Runner-
  Seed-SQL-Aktivierung (der MVP-Test fährt die Use-Case-Abfolge).

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-008 | CDC-Verwaltung als Use Cases | [`LH-FA-CFG-001`](../../../spec/lastenheft.md)…004 |
| slice-009 | Consumer-Verwaltung | [`LH-FA-CON-001`](../../../spec/lastenheft.md)…006 |
| slice-010 | Lesen-Vollabdeckung und SQL-Schnittstelle | [`LH-FA-REA-002`](../../../spec/lastenheft.md)…006 |
| slice-011 | Sicherheit und Observability-Basis | [`LH-QA-SEC-001`](../../../spec/lastenheft.md)…003 |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: nichts — die Welle-3-Slices sind die MVP-Rest-Lücken
  selbst; die Abhängigkeitskette läuft über die Slice-Trigger
  (slice-008 → 009 → 010 → 011, je WIP-Limit 1).
- Wird blockiert von: nichts — der Start-Trigger ist mit welle-2
  eingetreten.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Vollständige Konfigurationsschicht (TOML/YAML) — ENV-Minimal-
  verdrahtung bleibt (slice-007-Muster); die vollständige Schicht folgt
  in späteren Wellen.
- HTTP-/gRPC-API, CLI-Adapter, Exportadapter — bleiben optional ohne
  beobachtbaren Bedarf.
- Vollständige Observability-Abdeckung — Basis hier, volle
  LH-QA-OPS-Metriken folgen.
- Performance-Benchmarks — folgen nach dem grünen E2E
  (welle-2-§6-Rest).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: <Zeiger auf `welle-<NN>-results.md`, Geschwister im Ruheort `done/`>
Zähler: <Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort>
