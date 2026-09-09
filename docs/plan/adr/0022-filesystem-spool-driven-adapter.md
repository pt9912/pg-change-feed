# ADR-0022: Filesystem Spool als Driven Adapter

**Status:** Proposed

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`ADR-0021`](0021-large-transaction-buffer.md)

**Schärft:** [`ARC-010`](../../../spec/architecture.md),
[`SPEC-005`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der TransactionBufferPort (ADR-0021) braucht eine erste
Implementierung für das Spooling großer offener Transaktionen. Die Wahl
der Spool-Technik soll den CDC-Speicher nicht zum Zirkelbezug führen.

## Entscheidung

Eine erste TransactionBufferPort-Implementierung kann temporäre Dateien
verwenden.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Spool im RAM belassen | keine zweite Technik | verlagert das Kernproblem von ADR-0021 nur, löst es nicht |
| B — Spool in PostgreSQL | nur eine Technologie im Betrieb | belastet gerade den CDC-Speicher, dessen Ausfall die Spool-Situation oft erst erzeugt; Zirkelbezug |
| **C — temporäre Dateien** | entkoppelt vom CDC-Speicher; einfach zu reinigen und zu begrenzen | Dateisystem-Konfiguration, Berechtigungen und Aufräumbetrieb |

## Konsequenzen

- Positiv: Spool und CDC-Speicher stören sich nicht gegenseitig; einfache
  kapazitative Begrenzung über das Dateisystem.
- Negativ: Betrieb braucht eine dokumentierte Spool-Konfiguration.
- Folgepflicht: Crash-Verhalten ist explizit zu testen, bevor dieser ADR
  Accepted werden kann.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Abschließen der Crash-Tests für das Datei-Spooling — beobachtbar am
vorliegenden Testbeleg; bis dahin bleibt dieser ADR Proposed.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Proposed (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0022` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).