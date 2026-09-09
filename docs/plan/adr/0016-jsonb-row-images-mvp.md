# ADR-0016: JSONB Row Images im MVP

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-CAP-008`](../../../spec/lastenheft.md)

**Schärft:** [`SPEC-002`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Row Images (`old_data`/`new_data`) brauchen ein Speicherformat für den MVP.
Das Format soll den Lese- und Persistenzpfad ohne typbezogene
Mapping-Logik im Core halten und direkt mit SQL inspizierbar sein; über
die MVP-Grenze hinaus wurde es nicht bewertet.

## Entscheidung

`old_data` und `new_data` werden im MVP als `jsonb` gespeichert
([`SPEC-002`](../../../spec/pflichtenheft.md)). Nach Benchmarks wird die Entscheidung überprüft.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — typisierte Spalten je Quelltabelle | DB-seitige Typsicherheit | DDL-Explosion bei der generischen Change-Tabelle; Aktivierung würde zum Schemaeingriff |
| B — kompaktes binäres Eigenformat | geringer Speicherverbrauch | nicht mit SQL inspizierbar, Debugging erschwert, eigener Serializer |
| **C — jsonb Row Images** | flexibel über Schemata, SQL-inspizierbar, PostgreSQL-nativ | kein typisierter Feldschutz; Speicherverbrauch ggf. höher |

## Konsequenzen

- Positiv: schnelle MVP-Realisierung; Changes sind direkt per SQL
  inspizierbar; keine Mapping-Logik im Core.
- Negativ: typbezogene Prüfung verschiebt sich in die Anwendung.
- Folgepflicht: Überprüfung der Entscheidung nach den ersten Benchmarks
  (Re-Evaluierungs-Trigger unten).

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Die ersten abgeschlossenen Benchmark-Läufe auf den Stufen aus [`SPEC-014`](../../../spec/pflichtenheft.md)
(`LOAD_TIERS`) liegen vor — beobachtbar am Vorliegen der Benchmark-Ergebnisse,
nicht an einem Datum.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0016` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).