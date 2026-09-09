# PG Change Feed

> **Durable change feeds for PostgreSQL, powered by logical replication.**

## Was ist PG Change Feed?

PG Change Feed stellt persistente Change Feeds für bestehende PostgreSQL-Tabellen bereit. Es richtet sich an Anwendungen und Integrationen, die Änderungen zuverlässig und unabhängig konsumieren wollen, ohne selbst PostgreSQL WAL oder das Logical-Replication-Protokoll verarbeiten zu müssen.

## Was kann ich heute tun?

PG Change Feed befindet sich derzeit in der Architektur- und frühen Implementierungsphase.

Es gibt noch keinen produktionsreifen Daemon und keine stabile öffentliche API. Der aktuelle Stand besteht aus spezifizierten Anforderungen, Architekturentscheidungen und dem geplanten MVP.

Siehe:

- [`spec/lastenheft.md`](spec/lastenheft.md) für Anforderungen und Akzeptanzkriterien.
- [`spec/pflichtenheft.md`](spec/pflichtenheft.md) für die technische Spezifikation.
- [`docs/plan/adr/`](docs/plan/adr/) für Architekturentscheidungen.

## Warum PG Change Feed?

PostgreSQL stellt mit WAL, Logical Decoding und Logical Replication leistungsfähige CDC-Grundlagen bereit, aber keine allgemeine persistente CDC-Abstraktion mit Change-Historie, stabilen Positionen, unabhängigen Consumern und Retention. PG Change Feed schließt diese Lücke, ohne eine neue Datenbank, einen Message Broker oder Änderungen an der Quellanwendung vorauszusetzen.

## Kerngedanke

**PostgreSQL WAL → durable change feed → independent consumers.**

Eine Quellposition wird erst bestätigt, nachdem die zugehörigen Änderungen dauerhaft gespeichert wurden. Im Fehlerfall werden mögliche Duplikate gegenüber stillen Lücken bevorzugt.

## Was macht es vertrauenswürdig?

- **Prozess:** [`AGENTS.md`](AGENTS.md) für verbindliche Entwicklungsregeln und [`harness/README.md`](harness/README.md) für Source Precedence und vorhandene Gates.
- **Verträge:** [`spec/lastenheft.md`](spec/lastenheft.md) mit nachvollziehbaren `LH-*`-Anforderungen und Akzeptanzkriterien.
- **Technische Spezifikation:** [`spec/pflichtenheft.md`](spec/pflichtenheft.md).
- **Gates:** Nur tatsächlich implementierte und ausführbare Quality Gates werden hier als erfolgreich aufgeführt.
- **Auditierbarkeit:** Architekturentscheidungen liegen in [`docs/plan/adr/`](docs/plan/adr/), Planung in [`docs/plan/planning/`](docs/plan/planning/).

## Lizenz

MIT — siehe [`LICENSE`](LICENSE).
