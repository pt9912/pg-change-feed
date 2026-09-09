# ADR-0034: Ports nach Fähigkeiten

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-CON-003`](../../../spec/lastenheft.md),
[`LH-FA-RET-002`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-004`](../../../spec/architecture.md) (Outbound Ports)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Outbound Ports bestimmen, an welchen Stellen die Application von
Technologie abstrahiert. Die übliche Gewohnheit, je Entity ein Repository
anzulegen, würde hier Schnittstellen erzeugen, die konsistenzkritische
Vorgänge (Persist-before-ACK, Retention über Consumer-Positionen) über
mehrere Ports verteilen.

## Entscheidung

Outbound Ports werden entlang fachlich benötigter Fähigkeiten und
Konsistenzgrenzen definiert — nicht reflexartig als Repository pro Entity
(vorgesehen: ChangeStorePort, ReplicationAckPort, ConsumerStatePort,
SchemaStorePort, TransactionBufferPort, MetricsPort).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Repository je Entity | bekannte, mechanische Struktur | Invarianten wie Persist-before-ACK verteilen sich über mehrere Ports; Konsistenzgrenzen verschwimmen |
| B — ein einzelner Storage-Mega-Port | alles an einer Stelle | jeder Adapter implementiert Operationen, die er nicht trägt; Aufrufer sehen zu viel |
| **C — Ports je Fähigkeit und Konsistenzgrenze** | Invariante und Port fallen zusammen; Fakes und Adapter bleiben klein | die Port-Liste braucht aktive Pflege mit jedem neuen Fähigkeits-Bedarf |

## Konsequenzen

- Positiv: Konsistenzgrenzen (Persistenz + ACK, Retention) sind je Port
  als Ganzes austauschbar und fakebar.
- Negativ: Entity-zentrierte Zugriffe („gib mir die Tabelle") laufen über
  den Fähigkeits-Port ihrer Operation.
- Folgepflicht: Neue Ports werden gegen die Frage geprüft, ob sie eine
  Konsistenzgrenze komplett abbilden.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Fähigkeits-Orientierung gilt unabhängig von Technologie-
und Speicherwechseln.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).