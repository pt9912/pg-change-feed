# ADR-0006: Replication Stream als Driving Adapter

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-CFG-001`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-005`](../../../spec/architecture.md),
[`LH-FA-CFG-001.a`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Änderungsquelle muss den Systemkern treiben. Die naheliegende Lesart —
PostgreSQL ist ein „genutztes" System, also Driven — führt in die falsche
Rolle: Der Stream liefert Ereignisse, er wird nicht beauftragt.

## Entscheidung

Wir wählen **den PostgreSQL Logical Replication Stream als Driving
Adapter**. `PostgresReplicationStreamAdapter` ist Driving Adapter und ruft
`CaptureInboundPort` auf.

    PostgreSQL WAL
          |
          v
    Replication Stream
    Driving Adapter
          |
          v
    CaptureInboundPort
          |
          v
    Capture Application Service

Der Adapter empfängt und dekodiert `pgoutput`, entscheidet aber nicht über
Persistenz, Retention oder Source-ACK.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — PostgreSQL als Driven Adapter (Polling) | Passt zur Intuition „Datenbank = infolge" | CDC braucht Push und Flusskontrolle; Polling verliert Transaktionsgrenzen und verzögert den Capture-Pfad |
| B — Stream-Logik im Application-Layer | Weniger Indirektion | Core hängt am Replication-Protokoll; ADR-0002 bricht |
| **C — Replication Stream als Driving Adapter** | Rollen dem Kontrollfluss treu; Adapter bleibt dumm (empfangen, dekodieren, weiterreichen) | Wer mit „Driving = User-Interface" argumentiert, muss die Rollen-Trennung erst verstehen |

## Konsequenzen

- Positiv: Der Adapter kann nicht still Entscheidungen fällen, die dem Core
  zustehen (Persistenz, Retention, ACK) — das ist Invariante.
- Negativ: Die Grenze „dekodieren ja, entscheiden nein" muss im Review
  wachgehalten werden.
- Folgepflicht: Source-ACK als getrennte Outbound-Rolle (ADR-0007);
  Orchestrierung im Application-Layer (ADR-0027).

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).