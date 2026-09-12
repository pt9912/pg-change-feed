# Welle 6: Consumer-Zugriffsweg

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-6-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-12.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

[`LH-FA-CON-001.a`](../../../spec/pflichtenheft.md) und
[`LH-FA-CON-004.a`](../../../spec/pflichtenheft.md) sind seit `slice-009`
(`welle-3`) als benannte Spec-Lücke geführt: Registrierungs- und
Bestätigungslogik für Consumer sind eigenständig getestete Einheiten
(`internal/application/usecase`), aber ohne einen von außen erreichbaren
Zugriffsweg (kein CLI-Unterbefehl, keine Netzwerkschnittstelle) — ein
externer Consumer kann sich heute nicht selbst registrieren oder eine
Position bestätigen, ohne die CDC-Speichertabellen direkt zu beschreiben
und dabei jede Prüfung (insbesondere die Vorwärts-Invariante aus
`LH-FA-CON-004.a`) zu umgehen. Diese Welle entscheidet den Zugriffsweg
(Architect-ADR) und verdrahtet ihn für beide Fähigkeiten.

**Das *Mehr* gegenüber den einzelnen Slice-DoDs:** Kein Slice-DoD allein
beweist, dass ein externer Consumer den vollen Zyklus — registrieren, lesen,
bestätigen, neu starten und an der bestätigten Position fortsetzen — über
den neuen Zugriffsweg tatsächlich durchläuft, ohne die Speichertabellen
direkt zu berühren. Das zeigt erst ein Ende-zu-Ende-Test, der einen externen
Consumer ausschließlich über den neuen Zugriffsweg fährt.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Welle 5 liegt in `done/` (`welle-5-results.md`) und `make gates` ist grün
  auf `main` — kein Slice belegt aktuell `in-progress/`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Beide Slices in `done/`.
- `make gates` grün.
- Ende-zu-Ende-Test (slice-022 §2) zeigt: Ein externer Consumer
  registriert sich, liest Changes, bestätigt eine Position, startet neu und
  setzt an der bestätigten Position fort — ausschließlich über den neuen
  Zugriffsweg, ohne direktes Schreiben der CDC-Speichertabellen.
- Closure-Notiz in `welle-6-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-021 | Consumer-Registrierung: Zugriffsweg (ADR) + Verdrahtung | [`LH-FA-CON-001.a`](../../../spec/pflichtenheft.md) |
| slice-022 | Positions-Bestätigung über denselben Zugriffsweg | [`LH-FA-CON-004.a`](../../../spec/pflichtenheft.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- slice-022 setzt slice-021 voraus (nutzt denselben, dort entschiedenen und
  verdrahteten Zugriffsweg).
- Keine andere Welle wird blockiert oder blockiert diese.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Administrative Entfernung von Consumern**
  ([`LH-FA-CON-006`](../../../spec/lastenheft.md)) über denselben
  Zugriffsweg — anderer Vorgang, keine benannte Spec-Lücke; folgt bei
  Bedarf mit eigenem Slice.
- **Least-Privilege-Rollen-Adoption des neuen Zugriffswegs**
  (`BEO-PGC/rollen-verdrahtung`) — Bestand bleibt bewusst stehen: Die
  bestehende gemeinsame Instanz-DSN-Verdrahtung in
  `internal/bootstrap/wiring.go` wird von dieser Welle nicht aufgelöst;
  der neue Zugriffsweg reiht sich in die bestehende Verdrahtung ein, statt
  sie vorwegzunehmen.
- **Authentifizierung/Autorisierung des externen Zugriffswegs selbst**
  (wer darf registrieren/bestätigen) — anderer Vorgang: Diese Welle liefert
  den Zugriffsweg für die im Lastenheft geforderten Fähigkeiten, keine neue
  Sicherheitsschicht darüber.

## Vermerk für den Trigger-Audit bei Closure

[`LH-FA-SST-006`](../../../spec/pflichtenheft.md) (konkrete HTTP-/gRPC-API,
Lastenheft-CR außerhalb dieser Welle) erfüllt wörtlich
[`ADR-0020`](../adr/0020-http-grpc-optional.md)s Re-Evaluierungs-Trigger
(„Beobachtbarer Bedarf eines API-Consumers — sichtbar als Anforderung im
Lastenheft-Change"). Der Trigger-Audit dieser Welle (Modul 6,
Closure-Schritt 2) prüft deshalb explizit, ob `ADR-0020` weiterhin gilt oder
ein Folge-ADR fällig wird — unabhängig vom Ausgang ohne Auswirkung auf
`slice-021`/`slice-022`, die bereits auf CLI festgelegt sind
(`docs/plan/adr/architect-review-slice-021.md` §3, kanalgenerisch).

`LH-FA-SST-007` (Benachrichtigung über NATS, ebenfalls ein Lastenheft-CR
außerhalb dieser Welle) landete während `slice-022`s Umsetzung auf `main`
— dasselbe Muster wie oben, ohne dass ein bestehendes ADR einen Trigger
trägt, der wörtlich auf einen Benachrichtigungsweg zielt
(`docs/reviews/verify-slice-022.md` V-2). Beide Funde zusammen (`slice-021`
V-1, `slice-022` V-2) sind das zweite bzw. dritte Auftreten des Musters
„Lastenheft-CR landet im laufenden Slice-Commit-Fenster" — der
Trigger-Audit hält das als Prozessbeobachtung fest, unabhängig vom
jeweiligen ADR-Ausgang.

Der Lese-Schritt dieser Welle (Modul 6, Closure-Schritt 3) liest
außerdem `BEO-PGC/rollen-verdrahtung`: `slice-022` hat den dritten,
unabhängigen Vorgang beigetragen (`evidence/slice-022.md`), die
3×-Schwelle ist damit erreicht — der Eintrag bekommt bei dieser
Welle-Closure seinen Ausgang (verkörpert/geplant/gestrichen), nicht mehr
nur `weiter offen` ohne Zuweisung.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: [`welle-6-results.md`](welle-6-results.md), Geschwister im Ruheort `done/`.
Zähler: [`../observations/`](../observations/)`BEO-PGC/`, eine Ebene über dem Ruheort.
