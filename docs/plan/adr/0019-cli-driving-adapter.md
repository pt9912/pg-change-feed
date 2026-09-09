# ADR-0019: CLI als Driving Adapter

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-SST-003`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-005`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Lastenheft fordert eine CLI für Installation, Diagnose und
Administration ([`LH-FA-SST-003`](../../../spec/lastenheft.md)). Offen ist, ob die CLI einen eigenen
Datenzugriffsweg erhält oder wie jeder andere Kanal über die Inbound
Use Cases läuft.

## Entscheidung

Die CLI verwendet dieselben Inbound Use Cases wie andere Driving Adapters.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — CLI mit eigenem Direktdatenzugriff | schnell gebaut, weniger Indirektion | zweite Logikquelle; Drift gegenüber SQL- und Stream-Kanal; Diagnoseaussagen können auseinanderlaufen |
| B — keine CLI | keine Adapter-Arbeit | verletzt [`LH-FA-SST-003`](../../../spec/lastenheft.md) |
| **C — CLI über dieselben Inbound Ports** | alle Kanäle gleichwertig; Diagnose und Administration über eine Logikquelle | CLI darf Ports nicht abkürzen (Disziplin im Adapter) |

## Konsequenzen

- Positiv: einheitliche Semantik über alle Kanäle; CLI-Ausgaben und
  Statusabfragen ([`LH-FA-ADM-002`](../../../spec/lastenheft.md) ff.) sind vergleichbar.
- Negativ: CLI-Features sind auf die Use-Case-Formen beschränkt.
- Folgepflicht: —

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Kanal-Gleichbehandlung hängt an keiner beobachtbaren
Bedingung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0019` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).