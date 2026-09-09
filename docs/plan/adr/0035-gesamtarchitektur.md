# ADR-0035: Gesamtarchitektur

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-SST-001`](../../../spec/lastenheft.md)

**Schärft:** [§1 Komponenten-Übersicht](../../../spec/architecture.md)
(ARC-001 bis ARC-011)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Einzelentscheidungen (Schichtung, Ports, Adapter-Rollen, Persist-before-
ACK) brauchen eine Komposition, die zeigt, wie sie zusammen das Gesamtsystem
bilden: Welche Komponente gibt es, wie fließt eine Änderung von der
PostgreSQL-WAL bis zum Consumer?

## Entscheidung

Die Gesamtarchitektur ist die hexagonale Komposition aus Domain Core,
Application, Inbound-/Outbound Ports, Driving- und Driven Adapters und
Bootstrap. Der Capture-Pfad läuft von der WAL über den Replication Stream
(Driving) in die Application; die Bestätigung läuft über den
ReplicationAckPort (Driven) zurück an dieselbe Quelle:

    PostgreSQL WAL
          |
          v
    +-----------------------------+
    | Driving Adapters            |
    | Replication Stream          |
    | CLI / SQL / HTTP            |
    +-------------+---------------+
                  |
                  v
            Inbound Ports
                  |
                  v
    +-----------------------------+
    | Application                 |
    | Capture / Consumer /        |
    | Retention / Configuration   |
    +-------------+---------------+
                  |
                  v
    +-----------------------------+
    | Domain Core                 |
    | Changes / Transactions /    |
    | Positions / Consumers /     |
    | Schemas / Policies          |
    +-------------+---------------+
                  |
             Outbound Ports
          +-------+--------+---------+
          |                |         |
          v                v         v
     ChangeStore     ReplicationAck Metrics
          |                |         |
          v                v         v
    +-------------------------------------+
    | Driven Adapters                     |
    | PostgreSQL Store / PostgreSQL ACK / |
    | File Spool / Telemetry / Metadata   |
    +-------------------------------------+

Die lebhafte Form dieser Komposition — Komponenten mit `ARC-*`-Adressen,
Sequenzen und Schichttabelle — liegt in der Architektur-Sicht (§1, §2, §4).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — klassische Drei-Schichten (UI/Logik/Daten) | vertraute Form | Capture-Pfad und ACK-Pfad wären symmetrisch abgebildet, obwohl ihre Rollen unterschiedlich sind |
| B — verteilter Aufbau mit separaten Diensten für Capture, Store, Retention | Komponenten einzeln skalierbar | Verteilung widerspricht dem Grundbetrieb ohne externe Infrastruktur und macht Persist-before-ACK über Dienstgrenzen hinweg teuer |
| **C — hexagonale Monolith-Komposition** | Rollen (Driving vs. Driven) explizit; ein Prozess trägt die Invariante Persist-before-ACK | Komponenten-Grenzen brauchen Werkzeug (Import-Linting), um physisch zu bleiben |

## Konsequenzen

- Positiv: Eine Komposition als Adresse für alle Slices; dieselbe Form
  trägt MVP und spätere Erweiterungen (HTTP-/gRPC als zusätzlicher
  Driving Adapter).
- Negativ: Prozessgrenzen und Skalierungsgrenzen fallen für den MVP
  zusammen (High Availability bleibt zukünftige Erweiterung).
- Folgepflicht: Neue Komponenten werden gegen das Gesamt-Diagramm und die
  `ARC-*`-Tabelle der Architektur-Sicht eingepflegt.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Neu prüfen, sobald eine Lastenhefts-Anforderung verteilten Betrieb
erzwingt (z. B. High Availability als bindende Anforderung) — beobachtbar
an der Aufnahme in das Lastenheft.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).