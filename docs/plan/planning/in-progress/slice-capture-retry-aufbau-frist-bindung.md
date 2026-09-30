# Slice capture-retry-aufbau-frist-bindung: Der Wert der Aufbau-Frist in `Run` hängt an einem Test

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice
verschieden wäre.

**Bezug:** [`ADR-0136`](../../adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md)
Festlegung 2 (Aufbau-Frist 30 s), [`ADR-0135`](../../adr/0135-capture-transient-wiederholung-stream-zyklus.md)
(Gesamtfenster), [`LH-QA-REL-002`](../../../../spec/lastenheft.md).

**Berührte Spec-Stellen:** —

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Closure von
[slice-capture-transient-wiederholung](../done/slice-capture-transient-wiederholung.md).
**Datum:** 2026-09-30.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Aufruf `runStreamCycle(attemptCtx, streamRetryMaxDelay, open, …)`
in `Run` (`internal/bootstrap/wiring.go`) ist an einen Test gebunden: ein anderer
Wert der Aufbau-Frist färbt `make test` rot. Ausgangslage: die Mutation der
Frist auf `10*streamRetryWindow` blieb grün, gefahren vom Reviewer (Fixrunde 2,
N-5) und unabhängig vom Verifier (Verifikations-Report §4).

**Ausdrücklich NICHT in diesem Slice:**

- **Eine Änderung des Werts (30 s).** Der Wert ist Setzung der
  [`ADR-0136`](../../adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md);
  eine Schärfung ist dort an einen Betriebsbefund gebunden (Re-Evaluierungs-Trigger).
- **Die Bindung der Schließ-Pfade und des Streaming-Signals.** Beides ist im
  Schwester-Slice `slice-capture-transient-wiederholung` getragen.

## 2. Definition of Done

- [ ] Der Aufbau-Frist-Wert in `Run` ist als benannte Größe geführt und durch
      einen Quelltext- oder Verdrahtungstest gebunden. *Zu belegen durch:* die
      Mutation des Werts im Aufruf in `Run` färbt `make test` rot (Stelle,
      Instanz und Farbe im Bericht genannt).
- [ ] `make gates` grün — Exit-Code ungefiltert gesichert
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` (Aufruf von `runStreamCycle` in `Run`) | update | Wert als benannte Größe `streamSetupTimeout` (30 s, Wert unverändert); der Aufruf in `Run` übergibt sie statt `streamRetryMaxDelay`. |
| `internal/bootstrap/administration_startorder_internal_test.go` | update | `TestRunSourceTextPassesTheSetupTimeoutToTheStreamCycle` bindet den Wert der Größe (30 s) und ihre Verwendung als zweites Argument von `runStreamCycle` in `Run`; Muster: der Quelltext-Test zur Weitergabe von `OnStreaming`. |

**Suchlauf (`AGENTS.md` §3.13).** Bewegte Eigenschaft: der Wert, den `Run` als
Aufbau-Frist übergibt (vorher `streamRetryMaxDelay`, jetzt `streamSetupTimeout`).
Parent `b56b0e87`.

```suchlauf
b56b0e87 4 -n streamRetryMaxDelay -- internal
diff 3 -n streamRetryMaxDelay -- internal
b56b0e87 0 -n streamSetupTimeout -- internal
diff 8 -n streamSetupTimeout -- internal
b56b0e87 3 -n streamRetryMaxDelay -- docs spec harness :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 3 -n streamRetryMaxDelay -- docs spec harness :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
```

Befund: der Parent führte vier Stellen mit `streamRetryMaxDelay` in `internal`
(Definition, Backoff-Deckel, `streamRetryStableAfter`, der Aufruf in `Run`); der
Diff drei, die Aufbau-Frist trägt nun die eigene Größe. Die drei Treffer in der
Doku sind `ADR-0136` (`Accepted`, unberührbar) und beschreiben den Wert 30 s, der
gleich bleibt. Nichtgefunden: kein weiterer Aufrufer von `runStreamCycle` mit
einem Produktionswert außerhalb von `Run`.

## 4. Trigger

**Start** (`next` → `in-progress`): Priorität durch die Roadmap; WIP-Limit 1.

**Rückführungen:**

- `in-progress` → `next`: die Bindung verlangt einen Umbau von `Run`, der über
  eine benannte Größe hinausgeht.
- `in-progress` → `open`: der Wert ist selbst strittig (dann Architect-Frage).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector) grün +
Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

- **Ein Quelltext-Test bindet die Schreibweise, nicht das Verhalten.**
  *Erwartet, zu belegen durch:* die Mutation des Werts (nicht des Namens) färbt
  den Test rot. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure)*
- **Beobachtungs-Register (`../observations/`):** *(zu tragen bei Closure)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** *(zu tragen bei Closure)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:**
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert) ist
einschlägig: dieser Slice ist die Bindung, die der Review als Lücke fand.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
