# ADR-0025: Domain Events

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** —

**Schärft:** [`ARC-001`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Application und Betrieb brauchen interne Nachrichten über Zustandswechsel
— eine aufgenommene Transaktion, ein fortgeschrittener Consumer, eine
blockierte Retention, eine erkannte Schemaänderung, überschrittenes
Capture-Lag. Diese Steuernachrichten dürfen nicht mit den
CDC-Datenänderungen selbst vermischt werden.

## Entscheidung

Interne Domain Events wie `TransactionCaptured`, `ConsumerAdvanced`,
`RetentionBlocked`, `SchemaChanged` und `CaptureLagExceeded` sind erlaubt
und nicht mit CDC-Changes gleichzusetzen.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nur Rückgabewerte und Polling | keine Event-Mechanik | Zustandswechsel (Lag, blockierte Retention) sind nur umständlich beobachtbar |
| B — CDC-Changes als Event-Stream missbrauchen | ein einziger Kanal | vermischt Nutzdaten mit Steuerung; Consumer sähen interne Ereignisse als Changes — gefährlich für die Abgrenzung (LH §5) |
| **C — eigene interne Domain Events** | saubere Abgrenzung; Observability ohne Domain-Verletzung | zweite Ereignisart mit eigener Semantik und eigener Weiterleitung |

## Konsequenzen

- Positiv: CDC-Changes und interne Steuernachrichten sind klar getrennt;
  Telemetrie kann über Ports (ADR-0024) subscriben.
- Negativ: Events brauchen eine definierte Weiterleitung, die nichts mit
  dem Change-Pfad vermischt.
- Folgepflicht: —

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Abgrenzung zu CDC-Changes trägt die Systemgrenze des
Lastenhefts (§1, §5).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0025` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).