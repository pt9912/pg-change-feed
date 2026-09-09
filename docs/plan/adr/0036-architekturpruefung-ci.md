# ADR-0036: Architekturprüfung in CI

**Status:** Superseded by [`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md)

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-POR-002`](../../../spec/lastenheft.md)

**Schärft:** [§2 Schichten und Constraints](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Abhängigkeitsregel (Adapters → Application → Domain, Domain kennt keine
Treiber) ist vereinbart und wird durch physische Paketgrenzen sichtbar.
Solange kein Werkzeug sie prüft, bleibt sie Review-Prüfpflicht — und genau
bei wachsender Adapterzahl ist menschliche Prüfung die unzuverlässigste
Verteidigung.

## Entscheidung

CI soll verhindern, dass Domain von Application/Adapters und Application
von konkreten Adaptern abhängt — als statisches Import-Linting im
Gate-Lauf. (Status: Proposed — das Gate existiert noch nicht; die Regel
bleibt bis dahin Review-Prüfpflicht.)

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nur Review-Prüfpflicht | kein Werkzeugaufwand | skaliert schlecht; Verstöße passieren unbemerkt im Alltag |
| B — Laufzeit-Reflexion/Tests auf Abhängigkeiten | prüft das laufende Verhalten | teuer, lückenhaft; Abhängigkeiten sind eine statische Eigenschaft |
| **C — statisches Import-Linting im CI-Gate** | Verletzung scheitert am Gate, deterministisch und billig | Werkzeug-Regeln wollen mit der Paketstruktur gepflegt werden |

## Konsequenzen

- Positiv: Die Schichtregel aus der Architektur-Sicht wird zur maschinell
  geprüften Größe; Reviewer prüfen Semantik, das Gate prüft Deckung.
- Negativ: Ein zusätzliches Gate im CI-Lauf; Ausnahmen brauchen einen
  dokumentierten Ort (Suppression-Verbot, AGENTS.md §3.2).
- Folgepflicht: Gate einführen, sobald die Paketstruktur steht; vorher ist
  die Regel Review-Prüfpflicht.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| depguard (statisches Import-Linting) | Domain ⊬ Application/Adapters; Application ⊬ konkrete Adapter | — (Gate geplant) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Hochschalten von Proposed zu Accepted, sobald das Import-Linting-Gate im
Gate-Lauf existiert und grün läuft; sonst Erinnerung bei jeder
Welle-Closure (Trigger-Audit).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Proposed (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |
| 2026-09-09 | Superseded durch ADR-0041 (a-check als Maschinenform der Architektur-Prüfung; die depguard-Benennung ist ersetzt, die benannten Regeln gelten fort) | [ADR-0041](0041-a-check-maschinenform-architekturpruefung.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).