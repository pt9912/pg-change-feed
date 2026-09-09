# ADR-0003: Physische Modulgrenzen

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** —

**Schärft:** [`architecture.md §2 (Schichten und Constraints)`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Abhängigkeitsregel (ADR-0002) wirkt nur, wenn sie an physischen
Grenzen hängt — sonst wird sie von der Build-Werkzeugkasse nicht gesehen
und driftet zur Absichtserklärung. Der Paketbaum muss die Schichten
widerspiegeln.

## Entscheidung

Wir wählen **physische Trennung in `domain`, `application`,
`adapters/driving`, `adapters/driven` und `bootstrap`**. Architekturgrenzen
sollen in CI geprüft werden.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Flache Paketstruktur nach technischen Themen (models, services, utils) | Schneller Einstieg | Schichten existieren nur im Kopf; Importregeln nicht prüfbar |
| B — Grenzen nur über Code-Review sichern | Kein Build-Aufwand | Einzige Wächter sind Menschen; Drift ist die Regel, der Ausnahmefall die Schärfe |
| **C — Physische Trennung domain/application/adapters/bootstrap mit CI-Prüfung** | Grenzen sind sichtbar und mechanisch prüfbar; ADR-0002 wird einklagbar | Paket-Umzug beim ersten Wachstum nötig, wenn die Anfangs-Zuschnitte falsch gewählt sind |

## Konsequenzen

- Positiv: Die Import-Regeln aus ADR-0002 sind an Paketen prüfbar; neue
  Adapter haben einen festen Ablageort.
- Negativ: Kleine Querschnitts-Typen brauchen eine bewusste Heimatentscheidung.
- Folgepflicht: CI-Prüfung der Grenzen (ADR-0036); Paket-Feinschnitt in
  ADR-0031.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| depguard / go-arch-lint | Pakete unter `adapters/` importieren nichts aus `bootstrap`; `domain` importiert nichts aus `application` oder `adapters` | `— (Gate geplant, siehe ADR-0036)` |

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