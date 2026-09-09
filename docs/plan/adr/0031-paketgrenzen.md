# ADR-0031: Paketgrenzen

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** — (Struktur-Festlegung ohne direkte Lastenheft-Zeile)

**Schärft:** [§2 Schichten und Constraints](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Abhängigkeitsregel (Adapters → Application → Domain) ist auf der
Konzeptebene vereinbart. Ohne physische Trennung bleibt sie eine Lesegewohn-
heit: ein einziger Import macht den Domain-Kern von PostgreSQL-Bibliotheken
abhängig, ohne dass es ein Bauwerkzeug anzeigt. Build- und Paketstruktur
können die Grenze sichtbar machen.

## Entscheidung

Domain, Application und Adapter werden soweit sinnvoll auch auf Build-/
Package-Ebene getrennt — die Paketstruktur spiegelt die Schichten des
Architekturmodells.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — ein flaches Paket für alles | keine Import-Hürden, schneller Start | Grenzen sind reine Konvention; erste Verletzung passiert unbemerkt |
| B — Trennung nur per Review-Prüfpflicht | flexibel, keine Werkzeugpflege | menschliche Prüfung skalierbar mit jedem neuen Adapter schlechter |
| **C — physische Pakettrennung plus Import-Linting** | Verletzungen scheitern am Build-Werkzeug, nicht am Menschen | Paket-Schnittstellen müssen bewusst entworfen werden |

## Konsequenzen

- Positiv: Die Schichtregel ist maschinell sichtbar; zirkuläre Abhängig-
  keiten zwischen den Schichten werden ausgeschlossen.
- Negativ: Querung einer Paketgrenze braucht eine explizite Schnittstelle.
- Folgepflicht: Das Import-Linting-Gate wird mit ADR-0036 eingeführt; bis
  dahin ist die Grenze Review-Prüfpflicht.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| depguard | Domain importiert nicht aus Application/Adapters; Application nicht aus konkreten Adaptern | — (Gate geplant, siehe ADR-0036) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Paket-Spiegelung der Schichten gilt unabhängig von der
gewählten Sprache.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).