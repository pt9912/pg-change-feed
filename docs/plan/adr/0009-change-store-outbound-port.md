# ADR-0009: Change Store als Outbound Port

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-RET-001`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-004`](../../../spec/architecture.md),
[`ARC-009`](../../../spec/architecture.md),
[`SPEC-001`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Persistenz ist die Kernwirkung des Systems (LH-FA-RET-001): Changes müssen
dauerhaft gespeichert und wieder lesbar sein. Die Technik dahinter — welche
Datenbank, welches Schema — soll nicht in die Application sickern.

## Entscheidung

Wir wählen **Persistenz über den `ChangeStorePort` als Outbound Port**. Er
persistiert committed CDC-Transaktionen und stellt persistierte Changes
bereit. Erste Implementierung: `PostgresChangeStoreAdapter`.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Persistenzaufrufe direkt aus der Application | Weniger Indirektion | PostgreSQL-Schemawissen im Core; Store-Wechsel wird Umbau |
| B — Generisches Repository-Muster je Entity | Bekanntes Muster | Ports entstehen reflexartig pro Tabelle statt je Fähigkeit (ADR-0034 widerspricht dem) |
| **C — Ein Fähigkeits-Port `ChangeStorePort`** | Eine Konsistenzgrenze, ein Port; Speichertechnik frei wählbar; Test über Fake Port | Port-Vertrag muss Lesen und Schreiben gleichermaßen gut abdecken und wird damit breiter |

## Konsequenzen

- Positiv: Application-Tests laufen mit Fake Ports; der Persist-before-ACK-
  Kritische Pfad (ADR-0011) hat eine einzige, prüfbare Grenze.
- Negativ: Der Port-Vertrag ist eine öffentliche Stelle — Änderungen
  berühren alle Adapter und Fakes.
- Folgepflicht: Referenzimplementierung in ADR-0010; Portdefinition
  entlang Fähigkeiten in ADR-0034.

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