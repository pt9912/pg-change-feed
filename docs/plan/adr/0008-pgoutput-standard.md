# ADR-0008: pgoutput als Standard

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-SST-001`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-008`](../../../spec/architecture.md),
[`SPEC-010`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Der Logical Replication Stream (ADR-0006) braucht ein Output-Plugin, das die
Protokoll-Nachrichten erzeugt. PostgreSQL liefert `pgoutput` mit aus; alle
Alternativen sind externe Extensions mit eigenem Einrichtungs- und
Betriebspfad.

## Entscheidung

Wir wählen **`pgoutput` als Standard-Output-Plugin**. Keine zusätzliche
PostgreSQL-Extension ist für den Capture-Kern erforderlich.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — `wal2json` | Fertiges JSON; einfacher Prototyp | Externe Extension je Instanz nötig (Installations-/Betriebsfriction); Versionierung über Drittes |
| B — `decoderbufs` | Protobuf, maschinenlesbar | Extension-Pflicht; Protobuf-Abhängigkeit zusätzlich |
| C — `test_decoding` | In PostgreSQL enthalten | Nur Debug-Format, keine stabilen Row-Images — für Betrieb ungeeignet |
| **D — `pgoutput`** | In PostgreSQL enthalten; protokollnah; das native Replikationsformat — keine Zusatzinstallation | Binäres Protokoll, das selbst zu dekodieren ist; Änderungen am Protokoll schlagen in den Adapter |

## Konsequenzen

- Positiv: Der Capture-Kern läuft auf jeder unterstützten Instanz ohne
  DDL- oder Extension-Eingriff (passend zu [`LH-FA-CFG-006`](../../../spec/lastenheft.md)).
- Negativ: Der Adapter trägt die `pgoutput`-Dekodierung samt
  Protokollversionen; das ist bewusstes Adapterdetail (ADR-0032).
- Folgepflicht: Protokollversionen je unterstützter Major-Version prüfen
  ([`SPEC-012`](../../../spec/pflichtenheft.md), [`LH-QA-POR-001`](../../../spec/lastenheft.md)).

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Dekodierung bleibt über den Port-Aufruf austauschbar, ohne
dass diese Wahl revidiert werden müsste.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).