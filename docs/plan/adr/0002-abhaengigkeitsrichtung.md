# ADR-0002: Abhängigkeitsrichtung

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** —

**Schärft:** [`architecture.md §2 (Schichten und Constraints)`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die hexagonale Architektur (ADR-0001) trägt nur, wenn Abhängigkeiten in eine
Richtung zeigen. Andernfalls entsteht die Kopplung über Umwege wieder —
etwa über geteilte Infrastruktur-Typen zwischen Adaptern.

## Entscheidung

Wir wählen **Abhängigkeiten ausschließlich nach innen**:

    Adapters -> Application -> Domain

Die Domain kennt keine konkreten Treiber, Frameworks, Datenbanken,
Dateisysteme oder Telemetriesysteme.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Freie Importe zwischen allen Paketen | Geringste Reibung im Alltag | Hexagonale Architektur (ADR-0001) wird zur Namenskonvention ohne Wirkung |
| B — Nur „keine Adapter-zu-Adapter-Importe" | Verhindert die schlimmsten Verfilzungen | Application darf weiterhin von PostgreSQL-Typen abhängen; Quelle bleibt un austauschbar |
| **C — Abhängigkeiten ausschließlich nach innen** | Domain und Application bleiben treiberfrei; Tests ohne echte Infrastruktur; Regel ist mechanisch prüfbar | Disziplin bei der Modell-Trennung nötig; gelegentlich Mapper-Code an der Grenze |

## Konsequenzen

- Positiv: Die Domänenlogik wird in reinen Unit-Tests ohne Datenbank und
  ohne Framework geprüft; Adapter sind einzeln ersetzbar.
- Negativ: Technische Typen (z. B. LSN) dürfen nicht durchsickern —
  Mapping-Code an den Port-Grenzen ist die Folge (ADR-0005).
- Folgepflicht: Maschinelle Prüfung der Importregeln (siehe Fitness
  Function); physische Grenzen in ADR-0003.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| depguard | Domain importiert nichts aus Application/Adapters; Application importiert nichts aus Adaptern | `— (Gate geplant, siehe ADR-0036)` |

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