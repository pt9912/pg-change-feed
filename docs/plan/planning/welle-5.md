# Welle 5: Realer CDC-Capture-Lag

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-5-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-12.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

[`LH-FA-ADM-004`](../../../spec/lastenheft.md) (messbarer CDC-Abstand) ist
bislang nur als Näherung geliefert (`cdc_capture_lag_approx`,
Persistenz-Zeit-Proxy, slice-013) — kein realer Abstand zwischen
Quell-Commit und CDC-Verfügbarkeit. Diese Welle spiegelt den bereits von
`pglogrepl` gelieferten `CommitMessage.CommitTime`-Zeitstempel durch vier
Schichten (Replication-Decoder, Domäne, Anwendungsschicht/Ports,
Store-Adapter) bis in `cdc.transaction.committed_at` und löst die
Beobachtung [`BEO-PGC/cdc-capture-lag-real`](observations/BEO-PGC/cdc-capture-lag-real/observation.md)
auf: `cdc_capture_lag_approx` wird durch den kanonischen
`cdc_capture_lag` ([`SPEC-013`](../../../spec/pflichtenheft.md)) abgelöst.

**Das *Mehr* gegenüber den einzelnen Slice-DoDs:** Keine der drei Slices
allein beweist, dass der Zeitstempel tatsächlich vom Quell-Commit stammt
und nicht bloß durchgereicht wird, ohne anzukommen — das zeigt erst ein
Ende-zu-Ende-Lasttest, der eine künstliche Verzögerung zwischen
Quell-Commit und Verarbeitung einbaut und nachweist, dass `cdc_capture_lag`
diese Verzögerung tatsächlich abbildet (nicht nur die
Persistenz-Zeit-Differenz von vorher).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln.

- Welle 4 liegt in `done/` (`welle-4-results.md`) und `make gates` ist
  grün auf `main` — kein Slice belegt aktuell `in-progress/`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle drei Slices in `done/`.
- `make gates` grün.
- Ende-zu-Ende-Lasttest (slice-019 §2) zeigt: `cdc_capture_lag` bildet
  eine künstlich eingebaute Verzögerung zwischen Quell-Commit und
  Verarbeitung real ab — nicht nur die Persistenz-Zeit-Differenz.
- Closure-Notiz in `welle-5-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-017 | Commit-Zeitstempel: Decoder+Mapper+Domäne | [`LH-FA-ADM-004`](../../../spec/lastenheft.md) |
| slice-018 | Commit-Zeitstempel: Store-Adapter | [`LH-FA-ADM-004`](../../../spec/lastenheft.md) |
| slice-019 | `cdc_capture_lag` ablösen + Lasttest-Beleg | [`LH-FA-ADM-004`](../../../spec/lastenheft.md), [`SPEC-013`](../../../spec/pflichtenheft.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- slice-018 setzt slice-017 voraus (braucht den domänenseitigen
  Zeitstempel-Zugriff).
- slice-019 setzt slice-018 voraus (braucht den real gespeicherten Wert
  in `cdc.transaction.committed_at`).
- Keine andere Welle wird blockiert oder blockiert diese.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- `--provenance-output`/`--migration-overlay`-Einführung für den
  Schema-Rollout — anderer Vorgang (d-migrate-Workflow, unabhängig vom
  Zeitstempel-Wiring dieser Welle).
- Consumer-Registrierungs-/ACK-Zugriffsweg
  ([`LH-FA-CON-001.a`](../../../spec/pflichtenheft.md)/`004.a`) — anderer
  Vorgang, unbetroffen von der Capture-Seite.
- Rollen-Verdrahtung (`BEO-PGC/rollen-verdrahtung`,
  [`LH-QA-SEC-001`](../../../spec/lastenheft.md)/[`002`](../../../spec/lastenheft.md)) —
  anderer Vorgang, eigene Bootstrap-Änderung.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: [`welle-5-results.md`](welle-5-results.md), Geschwister im Ruheort `done/`.
Zähler: [`../observations/`](../observations/)`BEO-PGC/`, eine Ebene über dem Ruheort.
