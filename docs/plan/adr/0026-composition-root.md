# ADR-0026: Composition Root

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** —

**Schärft:** [`ARC-007`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die Verdrahtung konkreter Adapter mit Ports und Services soll nicht über
die Codebasis verteilt werden — sonst entsteht jede Abhängigkeit an
beliebiger Stelle und die Abhängigkeitsregel (§2 der Architektur-Sicht)
lässt sich nicht mehr lokal einhalten.

## Entscheidung

Nur `bootstrap` kennt konkrete Adapter und verdrahtet Driving Adapters,
Inbound Ports, Application Services, Outbound Ports und Driven Adapters.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Dependency-Injection-Framework | weniger Handarbeit, konventionelle Verdrahtung | Framework-Abhängigkeit auch im Bootstrap; Verdrahtung versteckt sich hinter Reflection |
| B — globale Singletons | am wenigsten Aufwand | versteckte Abhängigkeiten, unklare Lebenszyklen, Tests brauchen globalen Zustand |
| **C — expliziter Composition Root** | Verdrahtung an einer Stelle lesbar und testbar; Bootstrap darf alles (§2) | die Bootstrap-Datei wächst mit jeder Komponente |

## Konsequenzen

- Positiv: Abhängigkeiten sind an genau einer Stelle sichtbar; der
  Aufbau ist als Test reproduzierbar.
- Negativ: kein Framework übernimmt die Verdrahtungsarbeit.
- Folgepflicht: —

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — der Composition Root trägt die lokale Einhaltbarkeit der
Abhängigkeitsregel (§2).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0026` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).