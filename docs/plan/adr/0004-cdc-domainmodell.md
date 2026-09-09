# ADR-0004: Eigenes CDC-Domänenmodell

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** —

**Schärft:** [`SPEC-002`](../../../spec/pflichtenheft.md),
[`SPEC-003`](../../../spec/pflichtenheft.md),
[`ARC-001`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Der Capture-Pfad erhält PostgreSQL-Nachrichten (`pgoutput`-Relation-Daten,
Transaktions- und Streaming-Grenzen). Würden diese Nachrichtentypen direkt
zum Domänenmodell, würde die PostgreSQL-Kopplung durch den gesamten Core
laufen — gegen die Abhängigkeitsrichtung (ADR-0002).

## Entscheidung

Wir wählen **ein eigenes CDC-Domänenmodell**. PostgreSQL-Nachrichten werden
nicht zum Domänenmodell. Kernobjekte: Source, SourceTable,
ChangeTransaction, Change, SourcePosition, Consumer, ConsumerPosition,
SchemaVersion, RetentionPolicy.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — PostgreSQL-Nachrichtentypen als Domänenmodell | Kein Mapping-Aufwand | PostgreSQL-Kopplung im gesamten Core; andere Quellen scheitern an fremden Typen |
| B — Generisches JSON-Passthrough | Maximale Flexibilität | Keine Invarianten prüfbar (Reihenfolge, Transaktionszugehörigkeit); Logik wandert in Consumer |
| **C — Eigenes CDC-Domänenmodell** | Invarianten sind Domänen-Eigentum (ADR-0029); Quelle und Speicher bleiben tauschbar | Mapping-Code an beiden Port-Grenzen nötig |

## Konsequenzen

- Positiv: Domain-Invarianten (ADR-0029) sind an eigenen Typen prüfbar;
  `cdc.change` (SPEC-002) und `SourcePosition` (SPEC-003) sind die
  verbindlichen Abbilder.
- Negativ: Zwei Modellschichten — Nachricht und Domänenobjekt — müssen
  gepflegt und gegen Drift gehalten werden.
- Folgepflicht: Übersetzung der Relation Metadata in die Schema-Modelle
  (ADR-0015).

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