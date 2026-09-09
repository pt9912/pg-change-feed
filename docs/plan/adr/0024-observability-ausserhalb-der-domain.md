# ADR-0024: Observability außerhalb der Domain

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-OPS-003`](../../../spec/lastenheft.md),
[`LH-QA-OPS-004`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-006`](../../../spec/architecture.md),
[`ARC-011`](../../../spec/architecture.md),
[`SPEC-009`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Lastenheft fordert maschinenlesbare Metriken und strukturiertes
Logging ([`LH-QA-OPS-003`](../../../spec/lastenheft.md), [`LH-QA-OPS-004`](../../../spec/lastenheft.md)). Die Domain soll davon frei
bleiben — sie kennt weder Logging- noch Metrics-Frameworks —, während
Betriebsinformationen trotzdem aus allen Schichten ankommen müssen.

## Entscheidung

Logging-/Metrics-Frameworks bleiben Infrastruktur. MetricsPort und
EventSinkPort werden durch Driven Adapters implementiert.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Logging-/Metrics-Frameworks direkt in der Domain | bequem, keine Weiterleitung | Domain abhängig von Telemetrie-Technik; verletzt §2 der Architektur-Sicht; Testaufwand steigt |
| B — globale Telemetrie-Singletons | kein Port nötig | versteckte Abhängigkeit; Verdrahtung und Testdoubles intransparent |
| **C — Ports + Driven Adapters** | Domain pur testbar; Telemetrie-Backend austauschbar ([`ARC-011`](../../../spec/architecture.md)) | Weiterleitung der Betriebsereignisse muss explizit verdrahtet werden |

## Konsequenzen

- Positiv: Domain-Tests brauchen kein Telemetrie-Framework; Wechsel des
  Backends (Prometheus/OpenTelemetry) berührt nur den Adapter.
- Negativ: jede Betriebsinformation erreicht das Backend über Ports —
  kein bequemer Direktaufruf.
- Folgepflicht: Verdrahtung der Telemetrie-Adapter im Composition Root
  (ADR-0026).

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Domain-Freiheit von Telemetrie trägt die
Abhängigkeitsregel der Architektur-Sicht (§2).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0024` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).