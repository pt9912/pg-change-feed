# ADR-0013: Consumer als Domänenkonzept

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-CON-001`](../../../spec/lastenheft.md),
[`LH-FA-CON-004`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-001`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Mehrere unabhängige Consumer (LH-FA-CON-001) brauchen persistierte,
voneinander getrennte Positionen. Wie Lesen und Bestätigen zusammenwirken,
entscheidet darüber, ob Consumer-Positionen eine Fachbedeutung bekommen oder
technische Nebensache bleiben.

## Entscheidung

Wir wählen **Consumer als Domänenkonzept**: Benannte Consumer besitzen
unabhängige persistierte Positionen. Reguläre ACKs sind monoton; Reset ist
explizit administrativ.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Positionen nur im Consumer (Client-Seite) | Kein Serverzustand | Positionsverlust bei Client-Verlust; Retention kann nicht auf Positionen bauen (LH-FA-RET-004 bricht) |
| B — Eine globale Leseposition für alle Consumer | Minimaler Zustand | Consumer blockieren sich gegenseitig; LH-FA-CON-002 (Unabhängigkeit) bricht |
| **C — Benannte Consumer mit unabhängigen persistierten Positionen im Core** | Retention, Lag-Metriken und Neustart-Fortsetzung (LH-FA-CON-005) haben einen Fachbegriff; Regeln sind Invarianten | Zustandsverwaltung (ConsumerStatePort) ist eigene Pflicht |

## Konsequenzen

- Positiv: Lesen verändert Positionen nicht; ACK bewegt sie regulär nur
  vorwärts — beide Regeln sind als Domäneninvariante prüfbar (ADR-0029).
- Negativ: Der administrative Reset braucht eine bewusste, protokollierte
  Sonderoperation — er darf nicht als ACK-Tarnhaupt durchkommen.
- Folgepflicht: Consumer-Fortschritt je benanntem Consumer ist
  Persistenzpflicht (LH-FA-CON-003); Entfernen von Consumern ist
  administrativ (LH-FA-CON-006).

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