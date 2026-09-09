# ADR-0032: PostgreSQL bleibt Adapterdetail

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-POR-001`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-001`](../../../spec/architecture.md) (Domain Core),
[`ARC-008`](../../../spec/architecture.md) (PostgreSQL Logical Replication),
[`ARC-009`](../../../spec/architecture.md) (PostgreSQL Store)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

PostgreSQL tritt in diesem System doppelt auf: als Änderungsquelle
(Replication-Protokoll, LSN, `pgoutput`-Nachrichten) und als persistenter
Speicher. Beide Berührungspunkte könnten verlockend sein, ihre Typen und
Protokolldetails in die Domänenobjekte zu ziehen — das würde die Quelle zur
architektonischen Abhängigkeit des Kerns machen.

## Entscheidung

PostgreSQL-spezifische Typen und Protokollnachrichten bleiben in Adaptern.
Der Core verwendet eigene Modelle (SourcePosition, ChangeTransaction,
TableSchema).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — PostgreSQL-Typen direkt im Core | keine Mapping-Arbeit, direkter Zugriff auf LSN | der Kern hängt am Wire-Protokoll; Domänenlogik ist nicht mehr ohne reale Instanz testbar |
| B — generischer ORM-/Treiber-Layer im Core | einheitlicher Zugriff auf beliebige Stores | Abstraktion ohne konkreten zweiten Nutzen; der Core bekommt Fremdtypen über den Umweg |
| **C — eigene Domänenmodelle, Mapping im Adapter** | Kern bleibt technologie- und testbar; PostgreSQL-Versionen sind Adapterdetail | Mapping-Code je Berührungspunkt (Replication, Store, Metadata) |

## Konsequenzen

- Positiv: Mehrere PostgreSQL-Major-Versionen sind Adapterfrage, keine
  Kernfrage; der Core testet ohne Docker.
- Negativ: Jede PostgreSQL-Eigenheit braucht eine explizite Übersetzung.
- Folgepflicht: Neue PostgreSQL-Berührungspunkte werden als Adapter, nicht
  als Core-Erweiterung angelegt.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Technologie-Freiheit des Kerns ist Zweck der Schichtung
selbst.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).