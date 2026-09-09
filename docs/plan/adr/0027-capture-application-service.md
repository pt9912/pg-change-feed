# ADR-0027: Capture Application Service

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-REL-001`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-002`](../../../spec/architecture.md) (Application),
[`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md) (Persist-before-ACK)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Der Replication-Stream-Adapter empfängt und dekodiert `pgoutput` und
übersetzt technische Nachrichten in Aufrufe des CaptureInboundPort. Damit
entsteht die Frage, wer über den Zeitpunkt entscheidet, ab dem eine
SourcePosition als dauerhaft verarbeitet gilt — und damit, wann die Quelle
bestätigt werden darf. Diese Entscheidung liegt auf dem kritischen Pfad der
Persist-before-ACK-Invariante.

## Entscheidung

Die Orchestrierung von Persistenz und Source-ACK liegt im Application
Layer (Capture Application Service). Der Stream-Adapter entscheidet nicht
selbst, wann eine SourcePosition dauerhaft verarbeitet ist.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Stream-Adapter orchestriert Persistenz und ACK selbst | kürzester Weg, keine zusätzliche Schicht | Adapter trifft Persistenz-Entscheidungen; schwer testbar ohne reale Quelle |
| B — Persistenz im Adapter, ACK im Application Layer | weniger Aufrufe über Ports | gespaltene Verantwortung; die Invariante liegt an keiner Stelle vollständig |
| **C — Application Layer orchestriert beides** | Invariante an einer Stelle prüfbar; Adapter austauschbar und fakebar | ein zusätzlicher Indirection-Schritt auf dem heißen Pfad |

## Konsequenzen

- Positiv: Persist-before-ACK ist als Use-Case-Logik unit-testbar; der
  Stream-Adapter bleibt austauschbares Infrastrukturdetail.
- Negativ: Der Application Service ist für jeden Aufruf über zwei Ports
  (ChangeStorePort, ReplicationAckPort) hinweg entworfen.
- Folgepflicht: Domain-/Application-Tests decken die Invariante ab, bevor
  der reale Stream-Adapter angeschlossen wird.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Schichtung gilt unabhängig von Technologie- oder
Volumen-Änderungen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).