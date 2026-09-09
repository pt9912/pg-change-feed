# ADR-0010: PostgreSQL als CDC Store

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-ZIE-004`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-009`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Der ChangeStorePort (ADR-0009) braucht eine erste Implementierung. Das
Lastenheft verlangt Grundbetrieb ohne proprietäre Komponenten
(LH-ZIE-004) und die Zielgruppe betreibt bereits PostgreSQL —
Betriebskompetenz ist also vorhanden, eine zweite Technik wäre Zusatzaufwand.

## Entscheidung

Wir wählen **PostgreSQL als CDC-Speicher der Referenzimplementierung**
(`PostgresChangeStoreAdapter`). Ein externer Broker ist für den Grundbetrieb
nicht erforderlich.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Kafka als CDC-Store | Hoher Durchsatz, eingebaute Retention über Consumer-Gruppen | Proprietäre Infrastruktur im Grundbetrieb; widerspricht LH-ZIE-004 und der Systemgrenze (kein Message Broker) |
| B — RabbitMQ als CDC-Store | Weit verbreitet | Queue-Semantik statt positionsbasiertem Lesen; Verbrauch löscht History (LH-FA-REA-005 leidet) |
| C — Dedizierte CDC-Datenbank / dateibasierter Store | Keine Abhängigkeit vom Quellsystem | Neue Betriebswelt für die Zielgruppe; Persistenz-Garantien (LH-QA-REL-001) selbst zu bauen |
| **D — PostgreSQL als CDC-Store** | Transaktional, vorhandene Betriebskompetenz, SQL-Lesen frei (LH-FA-SST-002); dieselbe Technik wie die Quelle | Store und Quelle teilen sich eine Technik — Kapazitätsplanung muss beide Rollen sehen |

## Konsequenzen

- Positiv: Deployment bleibt eine Technik; das CDC-Schema (SPEC-001) nutzt
  PostgreSQL-Mittel (z. B. `jsonb`) direkt.
- Negativ: Quell- und Store-Instanz können sich gegenseitig belasten;
  Trennung der Instanzen ist Betriebsempfehlung.
- Folgepflicht: Store-Wechsel bleibt über den Port möglich (ADR-0009);
  keine Kopplung des Cores an PostgreSQL (ADR-0032).

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Neue Consumer-Klasse mit Speicheranforderungen, die das CDC-Schema
(SPEC-001) nicht abbilden kann — erkennbar am Beobachtungs-Register, wenn
Retentions- oder Mengen-Beobachtungen dreimal die PostgreSQL-Grenze treffen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).