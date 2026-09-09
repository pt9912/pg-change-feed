# ADR-0007: Source ACK als Outbound Port

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-REL-001`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-004`](../../../spec/architecture.md),
[`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die WAL-Bestätigung ist die gefährlichste Nebennhandlung des Capture-Pfads:
Wer zu früh bestätigt, verliert still Changes ([`LH-QA-REL-001`](../../../spec/lastenheft.md)). Die
Entscheidung, wann bestätigt wird, darf deshalb nicht beim Adapter liegen,
der denselben Stream liest.

## Entscheidung

Wir wählen **die Source-Bestätigung als Core-gesteuerte Wirkung über einen
Outbound Port**. `ReplicationAckPort` ist Outbound Port;
`PostgresReplicationAckAdapter` ist Driven Adapter. Stream und ACK dürfen
dieselbe technische Verbindung oder Bibliothek verwenden, sind
architektonisch jedoch getrennte Rollen.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — ACK im Stream-Adapter selbst | Weniger Struktur; schnellste Umsetzung | Der Empfänger bestätigt sich selbst; Persist-before-ACK (ADR-0011) ist nicht erzwingbar |
| B — ACK durch eine zweite, unabhängige Verbindung | Härtet die Rollen physisch | Zwei Replication-Verbindungen kosten Slot- und Betriebsoverhead ohne fachlichen Gewinn |
| **C — ACK als Outbound Port, Core-gesteuert** | Die Entscheidung „wann bestätigt" liegt in der Application; technisch darf dieselbe Verbindung dienen | Rollen-Trennung auf einer Verbindung muss im Review wachgehalten werden |

## Konsequenzen

- Positiv: Persist-before-ACK (ADR-0011) ist an der Port-Grenze erzwingbar;
  der Stream-Adapter kann frühestens wünschen, aber nicht bestätigen.
- Negativ: Dieselbe Bibliothek bedient zwei Rollen — Fehlbelegungen sind
  nur durch Benennung (Adapter-Namen) sichtbar.
- Folgepflicht: Invariante in ADR-0011 und ADR-0029; Abnahme im
  Persist-before-ACK-Kritischen Test.

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