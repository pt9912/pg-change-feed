# ADR-0039: Paketstruktur-Detaillierung (Go)

**Status:** Superseded by [`ADR-0042`](0042-transport-typen-am-port.md)
**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** —

**Schärft:** [`ARC-001`](../../../spec/architecture.md),
[`ARC-002`](../../../spec/architecture.md),
[`ARC-005`](../../../spec/architecture.md),
[`ARC-006`](../../../spec/architecture.md) (Komponenten der Sicht) und
[architecture.md §2](../../../spec/architecture.md) (Schichten-Constraints)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0038`](0038-implementierungssprache-go.md) legte Go als
Implementierungssprache und eine Grundpaketstruktur fest
(`internal/{domain, application, adapters, bootstrap}`), beließ aber die
innere Gliederung offen: Services flach, keine Command-/Result-Typen je
Use Case, keine adapter-interne Schichtung. Vor dem ersten
Implementierungs-Slice braucht die Struktur Tiefe, damit
Use-Case-Grenzen, Transport-Typen und Mapper adressierbar sind. Anlass
ist die Review einer Use-Case-granularen Hexagon-Vorlage
(Use-Case-Pakete, Adapter-Interne mit request/response bzw.
entity/repository/mapper, Fehler-Pakete je Schicht). Die
Paketstruktur-Änderung ist eine Schärfung — Sprache, Cross-Compilation,
CGO-Ziel und nativer PostgreSQL-Stack aus
[`ADR-0038`](0038-implementierungssprache-go.md) bleiben unverändert
bestehen und werden hiermit nicht berührt.

## Entscheidung

Die Paketstruktur aus [`ADR-0038`](0038-implementierungssprache-go.md)
wird wie folgt detailliert (Go-Adaption, CDC-Domäne):

    cmd/pg-change-feed/main.go

    internal/
    ├── domain/
    │   ├── model/        # Change, ChangeTransaction, SourcePosition,
    │   │                 # Consumer, ConsumerPosition, SchemaVersion,
    │   │                 # RetentionPolicy ([`SPEC-001`](../../../spec/pflichtenheft.md)…004)
    │   ├── event/        # TransactionCaptured, ConsumerAdvanced,
    │   │                 # RetentionBlocked, SchemaChanged,
    │   │                 # CaptureLagExceeded
    │   └── errors/       # Domänen-Fehler (Invarianten-Verletzungen)
    ├── application/
    │   ├── port/
    │   │   ├── inbound/  # CaptureInboundPort, EnableTableUseCase,
    │   │                 # DisableTableUseCase, ReadChangesUseCase,
    │   │                 # AcknowledgeConsumerUseCase,
    │   │                 # RegisterConsumerUseCase,
    │   │                 # ResetConsumerUseCase, RunRetentionUseCase,
    │   │                 # GetStatusUseCase
    │   │   └── outbound/ # ChangeStorePort, ReplicationAckPort,
    │   │                 # ConsumerStatePort, SchemaStorePort,
    │   │                 # TransactionBufferPort, MetricsPort, ClockPort
    │   └── usecase/
    │       ├── capture/      # CaptureService, CaptureCommand, CaptureResult
    │       ├── enable/       # EnableTableService, EnableTableCommand
    │       ├── disable/      # DisableTableService, DisableTableCommand
    │       ├── read/         # ReadChangesService, ReadChangesQuery, ReadChangesResult
    │       ├── acknowledge/  # AcknowledgeConsumerService, AcknowledgeConsumerCommand
    │       ├── register/     # RegisterConsumerService, RegisterConsumerCommand
    │       ├── reset/        # ResetConsumerService, ResetConsumerCommand
    │       ├── retention/    # RunRetentionService, RunRetentionCommand
    │       └── status/       # GetStatusService, GetStatusResult
    ├── adapters/
    │   ├── driving/
    │   │   ├── replication/  # Stream-Adapter: receive/, decode/, mapper/
    │   │   ├── cli/
    │   │   └── sql/
    │   └── driven/
    │       ├── postgresstorage/  # row-Typen, queries/, mapper/
    │       ├── postgresack/
    │       ├── metadata/         # row-Typen, queries/, mapper/
    │       ├── spool/
    │       └── metrics/
    └── bootstrap/

Jeder Use Case trägt sein Service- und Transport-Typ-Paar (Command/Query,
Result) im eigenen Paket; Domain-Events aus [`ADR-0025`](0025-domain-events.md)
liegen in `domain/event/`. Driving-Adapter halten ihre Transport-Typen
(request/response) und die Übersetzung in Domänenmodelle im `mapper/`
des eigenen Adapter-Pakets; Driven-Adapter halten Zeilenabbilder
(`SPEC-002`, Row Images) als row-Typen und übersetzen über den eigenen
`mapper/` — es entstehen keine Java-artigen Klassenhierarchien
([`ADR-0038`](0038-implementierungssprache-go.md) bleibt in dieser
Hinsicht unberührt).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Grundform aus ADR-0038 ohne Detaillierung | weniger Verzeichnisse, weniger Import-Disziplin | Command/Result-Typen geraten in Sammel-Pakete; Use-Case-Grenzen und Mapper-Verantwortung verschwimmen |
| B — Wurzel „hexagon/adapter/main" nach Hexagon-Referenz-Vorlage | konsistente Muster-Sprache mit Lehrbuch-Vorlage | dreistufige Umbenennung unüblich im Go-Ökosystem; bricht mit der verbindlichen Terminologie aus ADR-0001 (Domain/Application/Adapters) |
| C — Fachliche Wurzel (capture/, consumer/, retention/ je Paket vollvertikal) | Use-Case-Perspektive konsequent | bricht den Schichten-Constraint (Adapters → Application → Domain): fachliche Pakete mischen Schichten |
| **D — internal-Baum mit Use-Case-Paketen und adapter-interner Schichtung (gewählt)** | Use-Case-Grenzen sichtbar, Transport-Typen beim Use Case, Row-Bilder isoliert in Adapter-Mappern | mehr Pakete und mehr Import-Disziplin; depguard-Gate ist Folgepflicht |

## Konsequenzen

- Positiv: Command-/Query-/Result-Typen liegen beim Use Case; die
  Row-Image-Übersetzung ([`SPEC-002`](../../../spec/pflichtenheft.md))
  bleibt im jeweiligen Driven-Adapter; `domain/errors/` und
  `domain/event/` geben den zwei Querschnitts-Anliegen je ein Paket.
- Negativ: mehr Pakete, mehr Import-Regeln; ein depguard-Verstoß ist
  leichter möglich und deshalb linting-würdig.
- Folgepflicht: das geplante Import-Linting-Gate
  ([`ADR-0036`](0036-architekturpruefung-ci.md)) deckt die neuen Grenzen
  mit; Slice-Pläne adressieren Pakete über die `ARC-*`-Kennungen der
  Komponenten.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| depguard | `internal/domain` importiert nichts aus `internal/application`, `internal/adapters`, `internal/bootstrap`; `internal/adapters/driving` importiert nur `internal/application/port/inbound` und `internal/domain` | `— (Gate geplant, siehe [`ADR-0036`](0036-architekturpruefung-ci.md))` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

`permanent` — die Struktur wächst mit den Artefakten (Ports, Use Cases,
Adapter) ohne Umkehr-Event; ein künftiger Wandel ist eine neue
Entscheidung, keine erneute Prüfung dieser.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted — Supersedes ADR-0038 (Detaillierung aus Paketstruktur-Ideen, Hexagon-Referenz-Stil) | — |
| 2026-09-09 | Superseded durch ADR-0042 (Transport-Typen am Port: Definitions-Stelle geschärft, Typ-Aliase am Use Case); die übrigen Struktur-Regeln bleiben unverändert verbindlich | [ADR-0042](0042-transport-typen-am-port.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).