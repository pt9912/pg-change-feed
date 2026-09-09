# ADR-0015: Schema Evolution

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-SCH-004`](../../../spec/lastenheft.md),
[`LH-FA-SCH-005`](../../../spec/lastenheft.md)

**Schärft:** [`LH-FA-SCH-004.a`](../../../spec/lastenheft.md),
[`SPEC-004`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Replication Stream liefert Relation Metadata, die sich über die Zeit
ändern. Changes, die unter älteren Schemata erfasst wurden, müssen
weiterhin lesbar bleiben; eine stille Fehlinterpretation ist nach
LH-FA-SCH-004 unzulässig. Der Core braucht daher eine
technologieunabhängige Repräsentation der Schemata.

## Entscheidung

Relation Metadata wird in TableSchema-/SchemaVersion-Modelle übersetzt;
jeder Change referenziert seine Schema-Version. Nicht sicher
interpretierbare Schemaänderungen führen zu einem sichtbaren Fehler.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Relation Metadata 1:1 als Rohdaten durchreichen | kein Übersetzungsaufwand | PostgreSQL-Details sickern in den Core; keine stabile historische Interpretation |
| B — Schemata beim Lesen on-the-fly am aktuellen Stand auflösen | keine Versionshaltung | ältere Changes werden am aktuellen Schema interpretiert — genau die stille Fehlinterpretation, die LH-FA-SCH-004 ausschließt |
| **C — TableSchema-/SchemaVersion-Modelle je Change** | historisch stabile Interpretation; inkompatible Änderungen erkennbar meldbar | Versionshaltung und Übersetzungsaufwand |

## Konsequenzen

- Positiv: Changes sind über Schema-Versionen historisch stabil
  interpretierbar; inkompatible Typänderungen werden sichtbar gemeldet.
- Negativ: Mehraufwand für Übersetzung und Versionsverwaltung.
- Folgepflicht: SchemaStorePort als Outbound Port; Fehlerklasse `schema`
  (SPEC-008) für nicht sicher interpretierbare Änderungen.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Übersetzung in technologieunabhängige Modelle trägt die
Core-Unabhängigkeit (Abhängigkeitsregel, §2 der Architektur-Sicht).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0015` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).