# ADR-0005: SourcePosition abstrahiert LSN

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-DAT-004`](../../../spec/lastenheft.md)

**Schärft:** [`SPEC-003`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Positionen sind die zentrale Adresse des Systems: Lesen, ACK und Retention
arbeiten alle mit ihnen ([`LH-FA-DAT-004`](../../../spec/lastenheft.md)). PostgreSQL kennt dafür den LSN.
Würde der LSN-Typ direkt im Core leben, wäre die Abhängigkeitsregel
(ADR-0002) an der wichtigsten Stelle verletzt.

## Entscheidung

Wir wählen **ein technologieunabhängiges `SourcePosition`-Modell im Core**.
Der PostgreSQL-Adapter mappt PostgreSQL-LSN darauf; PostgreSQL-spezifische
LSN-Typen bleiben im Adapter.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — LSN-Typ direkt im Core | Kein Mapping, exakte Typen | PostgreSQL-Typen durchziehen Consumer, Store und Retention; ADR-0002 bricht an der heißen Stelle |
| B — Generischer String als Position | Trivial portabel | Keine Ordnung, keine Typsicherheit; Sortierung ([`LH-FA-REA-004`](../../../spec/lastenheft.md)) wäre String-Vergleich |
| **C — Eigenes `SourcePosition`-Modell mit Adapter-Mapping** | Ordnung und Typ im Core, LSN bleibt Adapterdetail | Mapping- und Vergleichslogik doppelt gepflegt (Adapter und Core) |

## Konsequenzen

- Positiv: Sortierbarkeit und Eindeutigkeit der Position ([`LH-FA-DAT-004`](../../../spec/lastenheft.md))
  sind Domänen-Eigenschaft; ein Quellwechsel berührt nur den Adapter.
- Negativ: Der Adapter muss jede LSN-Feinheit sauber abbilden; ein
  Verlust bei der Abbildung gefährdet die Eindeutigkeit.
- Folgepflicht: `SourcePosition` ist als [`SPEC-003`](../../../spec/pflichtenheft.md) verbindlich festgehalten.

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