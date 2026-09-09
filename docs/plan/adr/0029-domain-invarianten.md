# ADR-0029: Domain-Invarianten

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-REL-001`](../../../spec/lastenheft.md),
[`LH-FA-CON-003`](../../../spec/lastenheft.md),
[`LH-FA-RET-004`](../../../spec/lastenheft.md)

**Schärft:** [§2 Schichten und Constraints](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Das Lastenheft stellt Datenintegrität und Betriebssicherheit über alles
(Zuverlässigkeits-Anforderungen). Die Sicherung dieser Zusagen läuft über
acht Invarianten, die an keiner Adaptergrenze verhandelbar sein dürfen. Die
Frage ist, an welcher Stelle des Schichtmodells sie erzwungen werden.

## Entscheidung

Die Domain-Invarianten werden im Domain Core als Bestandteil des
Domänenmodells erzwungen:

1. Persist-before-ACK: Source-ACK nur nach dauerhafter Persistenz.
2. Consumer-ACK regulär nur vorwärts.
3. Offene Transaktionen sind nicht konsumierbar.
4. Rollbacks erzeugen keine regulären Changes.
5. Retention löscht keine benötigten Changes.
6. Eindeutige Transaktionszuordnung und Reihenfolge innerhalb einer
   Transaktion.
7. Jeder Change referenziert eine Schema-Version.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Invarianten nur in den Adaptern erzwingen | nah an der jeweiligen Technologie | dieselbe Invariante mehrfach je Adapter; ein neuer Adapter kann sie vergessen |
| B — Invarianten nur über Tests sichern | kein Code-Aufwand im Core | Tests beweisen Einzelstände, erzwingen keine Invariante zur Laufzeit |
| **C — Invarianten im Domain Core erzwingen** | einmal definiert, für alle Adapter bindend; unit-testbar ohne Infrastruktur | der Core braucht Konstruktionen, die illegale Zustände gar nicht erst zulassen |

## Konsequenzen

- Positiv: Jeder Adapter erbt die Invarianten; Verstöße sind zur Laufzeit
  unmöglich statt nur getestet.
- Negativ: Domänentypen werden restriktiver (z. B. Positionstyp ohne
  Rückwärts-ACK).
- Folgepflicht: Jede neue Adapter-Implementierung darf die Invarianten nur
  über die Domänentypen berühren.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Invarianten folgen aus den Lastenhefts-Zusagen, nicht aus
einer Technologie.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).