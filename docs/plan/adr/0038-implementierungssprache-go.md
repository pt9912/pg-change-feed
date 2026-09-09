# ADR-0038: Implementierungssprache Go

**Status:** Superseded by [`ADR-0039`](0039-paketstruktur-detaillierung-go.md)

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-POR-002`](../../../spec/lastenheft.md),
[`LH-QA-POR-003`](../../../spec/lastenheft.md),
[`LH-QA-OPS-001`](../../../spec/lastenheft.md)

**Schärft:** — (Sprachwahl ohne Spec-Stratum; wirkt auf Betrieb und
Deployment des Sicht-Stratums)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

PG Change Feed ist ein langlebiger Infrastruktur- und Streaming-Daemon mit
den Lastenhefts-Zusagen zu Portabilität (Linux primär, reproduzierbares
OCI-Deployment, mehrere PostgreSQL-Major-Versionen) und zu containerisiertem
Betrieb. Die Implementierungssprache entscheidet über Deployment-Aufwand,
Binärgröße, Nebenläufigkeit und die Eignung für das
Replication-Wire-Protokoll. Kotlin wurde als vorheriger Kandidat erwogen
und verworfen.

## Entscheidung

PG Change Feed wird in **Go** implementiert — gegenüber Kotlin als
primäre Implementierungssprache gewählt.

### Begründung

Für diesen Einsatzzweck werden insbesondere folgende Eigenschaften
priorisiert:

- einfache Cross-Compilation
- eigenständige Binaries
- geringe Runtime-Komplexität
- gute Eignung für langlebige Netzwerk- und Streaming-Prozesse
- geringer Deployment-Aufwand
- kleine und reproduzierbare Container-Images
- gute Unterstützung für PostgreSQL und das Wire-/Replication-Protokoll
- einfache Nebenläufigkeit
- schneller Prozessstart
- gute Eignung für Linux-basierte Infrastruktur

### Cross-Compilation

Mindestens folgende Targets sollen unterstützt werden:

    linux/amd64
    linux/arm64

Perspektivisch:

    darwin/arm64
    darwin/amd64
    windows/amd64

### CGO

Der Core und die Referenzimplementierung sollen nach Möglichkeit ohne CGO
auskommen. Ziel:

    CGO_ENABLED=0

Dadurch werden statisch bzw. weitgehend eigenständig auslieferbare Binaries
und einfache Cross-Builds ermöglicht.

### PostgreSQL

Für PostgreSQL wird ein nativer Go-Stack bevorzugt. Eine Abhängigkeit von
`libpq` soll vermieden werden. Die konkrete PostgreSQL-/Replication-
Bibliothek wird als Infrastrukturdetail behandelt und in einem separaten
ADR festgelegt.

### Hexagonale Architektur in Go

Go-Interfaces werden klein und capability-orientiert definiert:

    type ChangeStore interface {
        PersistTransaction(
            ctx context.Context,
            tx domain.ChangeTransaction,
        ) error
    }

    type ReplicationAcknowledger interface {
        Acknowledge(
            ctx context.Context,
            position domain.SourcePosition,
        ) error
    }

Es werden keine künstlichen Java-/Kotlin-artigen Klassenhierarchien
nachgebildet.

### Paketstruktur

Die Zielstruktur orientiert sich an:

    cmd/
    └── pg-change-feed/
        └── main.go

    internal/
    ├── domain/
    ├── application/
    │   ├── inbound/
    │   ├── outbound/
    │   └── service/
    ├── adapters/
    │   ├── driving/
    │   │   ├── replication/
    │   │   └── cli/
    │   └── driven/
    │       ├── postgres/
    │       ├── spool/
    │       └── metrics/
    └── bootstrap/

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Kotlin | starke Typen, JVM-Ökosystem, koroutinenbasierte Streams | JVM-Runtime im Container, größere Images, schwerere Cross-Compilation zu eigenständigen Binaries |
| B — Rust | maximale Kontrolle, keine Runtime, statische Binaries | höhere Einstiegs- und Wartungskosten für ein OSS-Projekt mit Contributor-Ziel |
| **C — Go (gewählt)** | einfache Cross-Compilation, eigenständige Binaries, native Nebenläufigkeit, kleine Images, gutes PostgreSQL-Replication-Ökosystem | weniger Ausdruckskraft im Typsystem; Invarianten liegen in Konstruktoren statt im Typsystem |

## Konsequenzen

- Go wird für Daemon, CLI und Kernimplementierung verwendet.
- Domain und Application bleiben von PostgreSQL-Bibliotheken unabhängig.
- PostgreSQL-spezifischer Code bleibt in Adaptern.
- Cross-Compilation wird Bestandteil der CI/CD-Pipeline.
- Kotlin wird für die Referenzimplementierung nicht weiter verfolgt.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Neu prüfen, wenn eine zwingende Plattform- oder Laufzeitanforderung den
Go-Weg ausschließt — beobachtbar an einer mit Go unrealisierbaren
Lastenhefts-Anforderung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |
| 2026-09-09 | Superseded durch ADR-0039 (Paketstruktur-Detaillierung); Sprache, Cross-Compilation, CGO-Ziel und nativer PostgreSQL-Stack bleiben verbindlich | [ADR-0039](0039-paketstruktur-detaillierung-go.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).