# PG Change Feed

*English | [Deutsch](README.de.md)*

> **Durable change feeds for PostgreSQL, powered by logical replication.**

## What is PG Change Feed?

PG Change Feed provides persistent change feeds for existing PostgreSQL tables. It targets applications and integrations that need to consume changes reliably and independently, without having to process PostgreSQL WAL or the logical replication protocol themselves.

## What can I do today?

Everything below is available end-to-end, real-world tested — via
environment variables, `docker compose`/`make` and SQL, with no
graphical interface.

| Area | What it does |
|---|---|
| **Capture** | Attach a PostgreSQL source via logical replication; (de)activate tables via SQL administration (`cdc.enable_table`/`disable_table`) or `CDC_TABLES` live; exclude individual columns from capture; store changes (INSERT/UPDATE/DELETE) transactionally and durably; detect schema changes on captured tables. |
| **Reading** | Read captured changes via SQL (`cdc.changes`) — or through one of four delivery paths: HTTP/JSON API, gRPC stream, Server-Sent Events or NATS (wake-up signal or full change content, table-granular). |
| **Consumers** | Register multiple independent consumers, acknowledge and resume their position — duplicates are preferred over silent gaps. |
| **Retention** | Run time- and consumer-based retention, surface blocking consumers before they prevent deletion. |
| **Operations** | Query operational status, CLI diagnostics, metrics and WAL lag. |
| **Security** | Enforce role-specific access rights (`cdc_capture`/`cdc_admin`/`cdc_reader`, least privilege — no separate login, access runs through the PostgreSQL connection itself). |
| **Distribution** | Available as an OCI image for `linux/amd64` **and** `linux/arm64` (GHCR and Docker Hub, identical digest). |

Details and examples for each access path (Go, C#, Kotlin) are in the
[user manual](docs/user/benutzerhandbuch.md); the full scope of
requirements and acceptance criteria is in [`spec/lastenheft.md`](spec/lastenheft.md).

See:

- [`docs/user/benutzerhandbuch.md`](docs/user/benutzerhandbuch.md) for usage (German).
- [`docs/user/releasing.md`](docs/user/releasing.md) for the release process.
- [`spec/lastenheft.md`](spec/lastenheft.md) for requirements and acceptance criteria.
- [`spec/pflichtenheft.md`](spec/pflichtenheft.md) for the technical specification.
- [`docs/plan/adr/`](docs/plan/adr/) for architecture decisions.

## Why PG Change Feed?

PostgreSQL provides powerful CDC foundations with WAL, logical decoding and logical replication, but no general-purpose, persistent CDC abstraction with change history, stable positions, independent consumers and retention. PG Change Feed closes that gap without requiring a new database, a message broker, or changes to the source application.

## Core idea

**PostgreSQL WAL → durable change feed → independent consumers.**

A source position is only acknowledged after its associated changes have been stored durably. In failure scenarios, possible duplicates are preferred over silent gaps.

## What makes it trustworthy?

- **Process:** [`AGENTS.md`](AGENTS.md) for binding development rules and [`harness/README.md`](harness/README.md) for source precedence and existing gates.
- **Contracts:** [`spec/lastenheft.md`](spec/lastenheft.md) with traceable `LH-*` requirements and acceptance criteria.
- **Technical specification:** [`spec/pflichtenheft.md`](spec/pflichtenheft.md).
- **Gates:** only actually implemented and executable quality gates are listed here as passing.
- **Auditability:** architecture decisions live in [`docs/plan/adr/`](docs/plan/adr/), planning in [`docs/plan/planning/`](docs/plan/planning/).

## License

MIT — see [`LICENSE`](LICENSE).
