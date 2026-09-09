# ADR-0017: Generische Change-Tabelle

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-DAT-001`](../../../spec/lastenheft.md)

**Schärft:** [`SPEC-001`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Aktivierbare Tabellen haben unterschiedliche Schemata; der Store muss
Changes einheitlich persistieren, lesen und der Retention unterwerfen —
ohne dass die Aktivierung einer Tabelle (LH-FA-CFG-001) zu einem
Schemaeingriff am CDC-Speicher wird.

## Entscheidung

Der MVP verwendet die generische Tabelle `cdc.change` statt einer
physischen Change-Tabelle je Quelltabelle. Tabellenspezifische Views
können später ergänzt werden.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Change-Tabelle je Quelltabelle | typisierte Spalten je Quelle | DDL je Aktivierung; unübersichtlich bei vielen Tabellen; laufende Migrationen |
| B — Partitionierung je Quelle innerhalb einer Tabelle | physische Trennung bei einem Tabellenvertrag | dynamische Partitionsverwaltung je Aktivierung; MVP-Komplexität |
| **C — generische Tabelle `cdc.change`** | Aktivierung ohne DDL am Store; einheitliche Lese-/Retention-Pfade | Werte nur als Row Images (ADR-0016); kein DB-seitiger Typschutz |

## Konsequenzen

- Positiv: Aktivierung ist eine Metadaten-Operation, kein Schemaeingriff;
  Lesen (LH-FA-REA-006) und Retention laufen über einen Pfad.
- Negativ: physische Datenlokalität je Quelle erst über spätere
  Partitionierung/Views adressierbar.
- Folgepflicht: tabellenspezifische Views als spätere Ergänzung möglich.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — der generische Vertrag trägt die Aktivierbarkeit je Tabelle
(LH-FA-CFG-001) unabhängig von konkreten Volumina.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0017` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).