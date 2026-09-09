# ADR-0028: Inbound Use Cases

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-CFG-001`](../../../spec/lastenheft.md),
[`LH-FA-REA-002`](../../../spec/lastenheft.md),
[`LH-FA-CON-004`](../../../spec/lastenheft.md),
[`LH-FA-ADM-001`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-003`](../../../spec/architecture.md) (Inbound Ports)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Alle Driving Adapters (Replication Stream, CLI, SQL, später HTTP-/gRPC)
brauchen eine einheitliche Aufnahme ihrer Aufrufe. Das Lastenheft fordert
SQL-basierte Administration, SQL-Lesezugriffe und eine CLI; jede dieser
Schnittstellen soll dieselben fachlichen Operationen aufrufen, ohne dass
Logik in den Adaptern entsteht.

## Entscheidung

Die Use Cases werden als explizite Inbound Ports geführt. Vorgesehen sind:

- CaptureInboundPort
- EnableTableUseCase
- DisableTableUseCase
- ReadChangesUseCase
- AcknowledgeConsumerUseCase
- RegisterConsumerUseCase
- ResetConsumerUseCase
- RunRetentionUseCase
- GetStatusUseCase

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — eine große Engine-Schnittstelle für alle Adapter | eine einzige Schnittstelle zu pflegen | jede Erweiterung ändert alle Adapter; Aufrufer sehen Operationen, die sie nicht nutzen |
| B — Use-Case-Klassen streng je CRUD-Entity | mechanisch ableitbar | fachliche Vorgänge (Aktivierung, ACK, Retention) lassen sich nicht auf CRUD abbilden |
| **C — Use Cases je fachlicher Fähigkeit** | Adapter rufen genau die Fähigkeit auf, die sie brauchen; Fakes je Use Case trivial | mehr Schnittstellen, die benannt und gepflegt sein wollen |

## Konsequenzen

- Positiv: CLI und SQL decken dieselben Use Cases ab; neue Driving Adapters
  (HTTP-/gRPC) kommen ohne Logik-Duplikation aus.
- Negativ: Für jede neue Fähigkeit entsteht ein expliziter Port.
- Folgepflicht: Die Use-Case-Liste wird mit jedem neuen Driving Adapter auf
  Vollständigkeit geprüft.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Fähigkeits-Schnittweise gilt unabhängig von der Zahl der
Adapter.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).