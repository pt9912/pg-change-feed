# Welle welle-4: Observability-Vervollständigung

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<NN>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug — M1 (MVP-Abnahme) ist bereits
erreicht (welle-2).

**Verantwortlich:** pt9912. **Datum:** 2026-09-10.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die in welle-3 bewusst zurückgestellte Observability-Basis wird
vervollständigt: Health-Endpoint (per Heartbeat-Pattern, Architect-Verdikt
aus slice-011 liegt bereits vor — `docs/plan/adr/architect-review-slice-011.md`),
sichtbare Fehlerzustände, messbarer CDC-Abstand und strukturiertes Logging.
Damit schließen [`LH-FA-ADM-002`](../../../spec/lastenheft.md)…004,
[`LH-QA-OPS-002`](../../../spec/lastenheft.md)…004 sowie der Rest von
[`LH-FA-SST-004`](../../../spec/lastenheft.md).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- welle-3 done: alle vier Slices in `done/`, vier Gates grün, Closure-Notiz
  geschrieben (Beleg: [`docs/plan/planning/done/welle-3-results.md`](done/welle-3-results.md),
  Verifikation) — **bereits eingetreten**.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices dieser Welle liegen in `done/` (slice-012…014).
- `make gates` grün **und** `make test-integration` grün am verdrahteten
  System — der welle-spezifische Beleg: der Compose-Healthcheck liest den
  Heartbeat-Zustand, nicht mehr nur den Prozess-Start.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-012 | Health-Endpoint per Heartbeat | [`LH-FA-ADM-002`](../../../spec/lastenheft.md), [`LH-QA-OPS-002`](../../../spec/lastenheft.md) |
| slice-013 | Fehlerzustände sichtbar, CDC-Abstand messbar | [`LH-FA-ADM-003`](../../../spec/lastenheft.md)/004, [`LH-QA-REL-003`](../../../spec/lastenheft.md) |
| slice-014 | Strukturiertes Logging | [`LH-QA-OPS-004`](../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: nichts — welle-4 liefert reine Observability-Erweiterung ohne
  Rückwirkung auf andere Wellen.
- Wird blockiert von: nichts — der Start-Trigger ist mit welle-3 eingetreten.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **`BEO-PGC/rollen-verdrahtung`** (Rollen-spezifische DSNs in
  `internal/bootstrap/wiring.go`) — anderer Vorgang: Sicherheits-Thema, kein
  Observability-Thema; bleibt Register-Beobachtung bis zum nächsten
  Schneiden.
- **[`LH-QA-OPS-005`](../../../spec/lastenheft.md)** (Upgrade-Sicherheit) —
  anderer Vorgang: ein Upgrade-Testlauf ist eine Test-/Release-Disziplin,
  keine Observability-Lieferung; bleibt offen für eine eigene Welle.
- **Volle [`SPEC-009`](../../../spec/pflichtenheft.md)-Metrik-Abdeckung**
  (`cdc_wal_retention_bytes`, `cdc_storage_bytes`) — Bestand bleibt bewusst
  stehen: diese zwei Kennzahlen hängen an Retention-Arbeit, die in dieser
  Welle nicht ansteht; sie folgen mit einer künftigen Retention-Welle.
- **HTTP-/gRPC-Health-Endpoint** — [`ADR-0020`](../adr/0020-http-grpc-optional.md)
  bleibt unberührt; der Health-Endpoint dieser Welle läuft über den
  bestehenden SQL-Kanal (Heartbeat-Tabelle + View), kein neuer
  Driving-Adapter-Typ.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: <Zeiger auf `welle-4-results.md`, Geschwister im Ruheort `done/`>
Zähler: <Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort>
