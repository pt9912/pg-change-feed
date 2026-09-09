# ADR-0021: Large Transaction Buffer

**Status:** Proposed

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** —

**Schärft:** [`SPEC-005`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Große Quelltransaktionen dürfen nicht unbegrenzt im RAM gehalten werden;
der Replication Stream streamt große Transaktionen in Teilen. Der Core
braucht eine Abstraktion für offene Transaktionen, die beide Fälle
abbildet, ohne den Persistenzpfad zu verändern.

## Entscheidung

`TransactionBufferPort` abstrahiert offene Transaktionen. Kleine
Transaktionen können im RAM liegen, große werden gespult.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — alle offenen Transaktionen im RAM | einfach, schnellster Pfad | unkontrolliertes OOM-Risiko bei großen Transaktionen |
| B — alle offenen Transaktionen sofort auf Platte | gleichmäßiges Verhalten | unnötige I/O im häufigen Fall kleiner Transaktionen |
| **C — Port mit Schwellwert-Schaltung** | kontrollierter Speicherverbrauch; Implementierung austauschbar | Schwellwert ist ein Tuning-Parameter |

## Konsequenzen

- Positiv: Speicherverbrauch ist kontrollierbar; die Spool-Technik ist
  kein Core-Thema ([`SPEC-005`](../../../spec/pflichtenheft.md)).
- Negativ: die Schwellwerte zwischen RAM-Haltung und Spool brauchen
  Konfiguration und Erklärung.
- Folgepflicht: Implementierungs-ADR (ADR-0022); Test mit großer
  Transaktion ohne unbegrenzten RAM-Verbrauch.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — solange der Replication Stream große Transaktionen in Teilen
liefert, braucht der Core diese Abstraktion.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Proposed (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0021` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).