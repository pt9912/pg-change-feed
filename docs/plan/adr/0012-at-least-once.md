# ADR-0012: At-Least-Once

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-REL-004`](../../../spec/lastenheft.md)

**Schärft:** [`LH-QA-REL-004`](../../../spec/lastenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Persist-before-ACK (ADR-0011) zieht Wiederholung als Normalfall nach sich:
Ein Crash zwischen Persistenz und ACK führt zu erneut verarbeiteten
Changes. Das Lastenheft nimmt genau das in Kauf ([`LH-QA-REL-004`](../../../spec/lastenheft.md)) und
schließt generisches Exactly-Once über externe Systeme hinweg aus.

## Entscheidung

Wir wählen **At-Least-Once-Semantik für Consumer**. Generisches
systemübergreifendes Exactly-Once wird nicht versprochen.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Exactly-Once über alle Consumer-Systeme versprechen | Bequemste Consumer-Verpflichtung | Verteiltes Exactly-Once ist ohne Protokoll am Ziel nicht leistbar; Versprechen wäre unehrlich |
| B — At-Most-Once | Keine Duplikate | Verluste bei Crash — verletzt [`LH-QA-REL-001`](../../../spec/lastenheft.md) |
| **C — At-Least-Once** | Ehrlich leistbar; Verluste sind ausgeschlossen, Duplikate sind sichtbar und deduplizierbar | Consumer müssen idempotent arbeiten — Pflicht, die beim Consumer ankommt |

## Konsequenzen

- Positiv: Das Modell und die Schnittstellen erlauben idempotente
  Verarbeitung ([`LH-QA-REL-004`](../../../spec/lastenheft.md)); Wiederholung ist kein Fehlerfall, sondern
  Betrieb.
- Negativ: Jeder Consumer trägt die Idempotenz-Pflicht; die Doku muss das
  an einer deutlich sichtbaren Stelle sagen.
- Folgepflicht: Erneutes Lesen desselben Bereichs liefert dieselben Changes
  ([`LH-FA-REA-005`](../../../spec/lastenheft.md)) — Grundlage für Consumer-Deduplizierung.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Semantikwahl trägt die Vertragszusage [`LH-QA-REL-004`](../../../spec/lastenheft.md).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).