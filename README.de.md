# PG Change Feed

*[English](README.md) | Deutsch*

> **Durable change feeds for PostgreSQL, powered by logical replication.**

**Quelltext & Issues:** [github.com/pt9912/pg-change-feed](https://github.com/pt9912/pg-change-feed)

## Was ist PG Change Feed?

PG Change Feed stellt persistente Change Feeds für bestehende PostgreSQL-Tabellen bereit. Es richtet sich an Anwendungen und Integrationen, die Änderungen zuverlässig und unabhängig konsumieren wollen, ohne selbst PostgreSQL WAL oder das Logical-Replication-Protokoll verarbeiten zu müssen.

## Was kann ich heute tun?

Alles unten steht Ende-zu-Ende real getestet zur Verfügung — über
Umgebungsvariablen, `docker compose`/`make` und SQL, ohne grafische
Oberfläche.

| Bereich | Was es tut |
|---|---|
| **Erfassung** | PostgreSQL-Quelle per Logical Replication anbinden; Tabellen per SQL-Administration (`cdc.enable_table`/`disable_table`) oder `CDC_TABLES` live (de)aktivieren; einzelne Spalten vom Erfassen ausschließen; Änderungen (INSERT/UPDATE/DELETE) transaktionsgetreu und dauerhaft speichern; Schemaänderungen an erfassten Tabellen erkennen. |
| **Lesen** | Erfasste Änderungen per SQL (`cdc.changes`) lesen — oder über eine von vier Zustellwegen: HTTP-/JSON-API, gRPC-Stream, Server-Sent-Events oder NATS (Wecksignal oder vollständiger Change-Inhalt tabellen-granular). |
| **Consumer** | Mehrere unabhängige Consumer registrieren, ihre Position bestätigen und fortsetzen — Duplikate werden stillen Lücken vorgezogen. |
| **Aufbewahrung** | Zeit- und consumer-basierte Retention betreiben, blockierende Consumer sichtbar machen, bevor sie die Löschung verhindern. |
| **Betrieb** | Betriebsstatus, CLI-Diagnose, Metriken und WAL-Rückstand abfragen. |
| **Sicherheit** | Rollenspezifische Zugriffsrechte durchsetzen (`cdc_capture`/`cdc_admin`/`cdc_reader`, Least-Privilege — kein eigener Login, Zugriff läuft über die PostgreSQL-Verbindung selbst). |
| **Distribution** | Als OCI-Image für `linux/amd64` **und** `linux/arm64` beziehen (GHCR und Docker Hub, identischer Digest); HTTP-API, gRPC-Stream, SSE-Stream und NATS-Vollinhalts-Stream stehen zusätzlich als offizielle C#-Client-Bibliothek zur Verfügung ([`PgChangeFeed.Client`](https://www.nuget.org/packages/PgChangeFeed.Client) auf NuGet.org), die HTTP-API zusätzlich als offizielle Python-Client-Bibliothek ([`pgchangefeed`](https://pypi.org/project/pgchangefeed/) auf PyPI). |

Details und Beispiele je Zugriffsweg (Go, C#, Kotlin) stehen im
[Benutzerhandbuch](docs/user/benutzerhandbuch.md); der volle Anforderungs-
und Akzeptanzkriterien-Umfang in [`spec/lastenheft.md`](spec/lastenheft.md).

Siehe:

- [`docs/user/benutzerhandbuch.md`](docs/user/benutzerhandbuch.md) für die Bedienung.
- [`sdks/csharp/`](sdks/csharp/) für die offizielle C#-Client-Bibliothek ([`PgChangeFeed.Client`](https://www.nuget.org/packages/PgChangeFeed.Client) auf NuGet.org).
- [`sdks/python/`](sdks/python/) für die offizielle Python-Client-Bibliothek ([`pgchangefeed`](https://pypi.org/project/pgchangefeed/) auf PyPI).
- [`docs/user/releasing.md`](docs/user/releasing.md) für den Release-Prozess.
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
