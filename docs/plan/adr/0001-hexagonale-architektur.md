# ADR-0001: Hexagonale Architektur

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** —

**Schärft:** [`architecture.md §1 (ARC-001…007)`](../../../spec/architecture.md),
[`architecture.md §2 (Schichten)`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

PG Change Feed koppelt eine technische Änderungsquelle (PostgreSQL Logical
Replication) an Consumer über eine persistente CDC-Abstraktion. Die Quelle
soll austauschbar bleiben, die Domänenlogik soll ohne Treiber testbar sein.
Die Terminologie muss vorab verbindlich sein, damit Ports und Adapter nicht
je Implementierung neu ausgehandelt werden.

Verbindliche Terminologie:

    Ports
    ├── Inbound Ports
    └── Outbound Ports

    Adapters
    ├── Driving Adapters
    └── Driven Adapters

Kontrollfluss:

    Driving Adapter
          |
          v
      Inbound Port
          |
          v
    Application / Domain
          |
          v
      Outbound Port
          |
          v
      Driven Adapter

## Entscheidung

Wir wählen **die hexagonale Architektur (Ports & Adapters)**. Driving
Adapters rufen Inbound Ports auf; Application Services verwenden Outbound
Ports; Driven Adapters implementieren Outbound Ports.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Geschichtete Monolith-Architektur (UI/Service/Repo) | Weit verbreitet, niedrige Einstiegshürde | Domäne verfilmt mit Technik; Austausch der Änderungsquelle und Speicher-Implementierung wird zuRefactoring über Schichtgrenzen hinweg |
| B — Mikroservices | Harte Prozessgrenzen, unabhängige Deployments | Betriebskomplexität für einen einzigen Daemon; Transaktionsgrenzen (Persist-before-ACK) über Prozesse hinweg teuer |
| **C — Hexagonale Architektur (Ports & Adapters)** | Domäne pur testbar; Quelle, Speicher und Telemetrie über Ports austauschbar; Kontrolfluss explizit | Ports/Adapter-Disziplin erfordert Eigenleistung; etwas mehr Struktur als ein Schichten-Monolith |

## Konsequenzen

- Positiv: PostgreSQL bleibt Adapterdetail; Domain- und Application-Tests
  laufen ohne echte Datenbank (Fake Ports).
- Negativ: Ein Adapter-Boilerplate bleibt dauerhafter Aufwand; jede neue
  technische Fähigkeit braucht eine Port-Erweiterung.
- Folgepflicht: Abhängigkeitsrichtung in ADR-0002, physische Modulgrenzen
  in ADR-0003, Gesamtzusammensetzung in ADR-0035.

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