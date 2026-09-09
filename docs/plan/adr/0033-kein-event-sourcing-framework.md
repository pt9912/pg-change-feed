# ADR-0033: Kein Event-Sourcing-Framework

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** — (Abgrenzung; deckungsgleich mit dem Out-of-Scope-Abschnitt
des Lastenhefts)

**Schärft:** — (Abgrenzungs-ADR ohne Spec-Stratum)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

CDC-Changes ähneln oberflächlich Events: sie sind chronologisch geordnet,
positioniert und konsumierbar. Diese Ähnlichkeit lädt dazu, ein
Event-Sourcing-Framework (Aggregate, Event Store, Projections, Replays auf
Aggregatzustand) als Fundament zu übernehmen. Das Lastenheft grenzt das
System jedoch ausdrücklich ab: Technische Datenänderungen sind nicht
automatisch semantische Business Events.

## Entscheidung

PG Change Feed ist kein generisches Event-Sourcing-System. CDC-Changes sind
technische Datenänderungen; kein Event-Sourcing-Framework wird als
Fundament übernommen.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Event-Sourcing-Framework als Basis | fertige Event-Store- und Projection-Werkzeuge | Aggregate-/Replay-Semantik passt nicht auf Tabellen-Changes; Abhängigkeit von einem Framework |
| B — generisches Streaming-Backend als Speicher | ausgereifte Verteilung | widerspricht der Forderung, ohne externe Broker betreibbar zu sein |
| **C — CDC-Abstraktion ohne Event-Sourcing-Anspruch** | schlankes Modell entlang des Lastenhefts; Consumer können später semantic events darauf bauen | semantische Events sind Nutzer-Leistung, nicht eingebautes Feature |

## Konsequenzen

- Positiv: Das Datenmodell bleibt an Tabellen-Changes orientiert; keine
  Framework-Abhängigkeit im Grundbetrieb.
- Negativ: Wer Business Events braucht, baut sie selbst über die
  Changes.
- Folgepflicht: Der Out-of-Scope-Abschnitt des Lastenhefts bleibt mit dieser
  Abgrenzung deckungsgleich.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Abgrenzung ist im Lastenheft als Out-of-Scope verankert;
sie kippt nur mit einer Lastenheft-Änderung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).