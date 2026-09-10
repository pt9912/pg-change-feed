# ADR-0018: SQL als Driving Adapter

**Status:** Superseded by [`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md)

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-SST-002`](../../../spec/lastenheft.md),
[`LH-FA-ADM-001`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-005`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Lastenheft und Pflichtenheft fordern SQL-basierte Administration und
Kernlesezugriffe ([`LH-FA-ADM-001`](../../../spec/lastenheft.md), [`LH-FA-SST-002`](../../../spec/lastenheft.md)). Die Frage ist, ob die
Logik dieser Zugriffe in SQL selbst liegt oder SQL nur als Eintrittspunkt
in dieselben Use Cases dient wie andere Kanäle.

## Entscheidung

SQL-Funktionen/Views sind Driving Adapter und rufen Inbound Ports auf.
Businesslogik wird nicht in SQL dupliziert.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Logik in PL/pgSQL-Funktionen | nah am Datenbestand, keine Indirektion | Duplikat der Domänenlogik in einer zweiten Sprache; nicht domänentestbar; Wartung an zwei Orten |
| B — keinen SQL-Zugriff anbieten | keine Adapter-Arbeit | verletzt [`LH-FA-SST-002`](../../../spec/lastenheft.md) und [`LH-FA-ADM-001`](../../../spec/lastenheft.md) |
| **C — SQL als Driving Adapter** | alle Kanäle gleichberechtigt über dieselben Ports; eine Logikquelle | SQL-Adapter-Schicht muss gepflegt und berechtigbar sein ([`LH-QA-SEC-002`](../../../spec/lastenheft.md)) |

## Konsequenzen

- Positiv: Admin-Funktionen aus SQL und CLI sind gleichwertig; es gibt eine
  einzige Logikquelle.
- Negativ: View-/Funktions-Schicht darf keine eigene Regeln entwickeln.
- Folgepflicht: getrennte Berechtigbarkeit von Administration und Lesezugriff
  ([`LH-QA-SEC-002`](../../../spec/lastenheft.md)) in der Objektanlage berücksichtigen.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Kanal-Gleichbehandlung hängt an keiner beobachtbaren
Bedingung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |
| 2026-09-10 | Superseded durch ADR-0046 (Rollenspaltung nach Datenrichtung: reine Lese-Views dürfen Driven-Adapter-Tabellen direkt lesen; die hier benannte Pflicht — schreibende/aktionsauslösende SQL-Funktionen rufen Inbound Ports auf — gilt fort) | [ADR-0046](0046-sql-driving-adapter-lese-schreib-trennung.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0018` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
