# ADR-0030: Testpyramide

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-POR-003`](../../../spec/lastenheft.md),
[`LH-QA-REL-001`](../../../spec/lastenheft.md)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die kritischen Zusagen des Lastenhefts (keine stillen Verluste, korrekte
Reihenfolge, reproduzierbare Testumgebung) sind nur prüfbar, wenn Tests auf
mehreren Ebenen laufen: reine Domänenlogik braucht keine reale PostgreSQL-
Instanz, die Persist-before-ACK-Invariante braucht mindestens eine. Die
Teststrategie entscheidet, wie viel Infrastruktur welcher Prüfungsstufe
zugestanden wird.

## Entscheidung

Die Teststrategie folgt einer Pyramide: viele Domain-Tests, Application-
Tests mit Fake Ports, gezielte Adaptertests, reale
PostgreSQL-Integrationstests und wenige kritische E2E-Tests.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nur E2E über reale PostgreSQL | prüft genau das ausgelieferte Verhalten | langsam, flaky-prone; Domänenlogik wird über Infrastrukturfehler unlesbar |
| B — nur Unit-Tests mit Mocks | schnell, vollständig deterministisch | beweist nichts über das Zusammenspiel mit realer Replikation |
| **C — Pyramide über alle Stufen** | schnelle Rückmeldung unten, Realitätsnähe oben in kleiner Zahl | Fake Ports und Adaptertests wollen gepflegt sein |

## Konsequenzen

- Positiv: Die MVP-Integrationstests (Lastenheft §1) laufen auf der Spitze
  der Pyramide, während Invarianten unten in Masse geprüft werden.
- Negativ: Fake Ports müssen parallel zu den echten Ports weiterentwickelt
  werden.
- Folgepflicht: Die reproduzierbare Testumgebung (Docker-Compose) trägt die
  Integrations- und E2E-Stufe.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go test-Layout | Stufen als getrennte Pakete/Suiten erkennbar; Integrationstests benötigen reale Instanz | — (Gate geplant) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Stufenfolge folgt aus dem Schichtmodell, nicht aus einer
Technologie.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).