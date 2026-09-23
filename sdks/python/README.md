# PG Change Feed — Python SDK

Official Python client library for [PG Change Feed](https://github.com/pt9912/pg-change-feed), a durable change feed system for PostgreSQL built on logical replication.

This package (`pgchangefeed`) lets a Python application consume PG Change Feed's delivery paths (HTTP API, gRPC stream, SSE stream, NATS full-content stream) without implementing the wire protocol itself — see [`LH-FA-SST-009`](https://github.com/pt9912/pg-change-feed/blob/main/spec/lastenheft.md) for the requirement this SDK fulfills.

## Status

This package is at an early, pre-1.0 stage (`0.x.y`, [ADR-0107](https://github.com/pt9912/pg-change-feed/blob/main/docs/plan/adr/0107-python-pypi-zweites-sdk-package.md), extended by [ADR-0110](https://github.com/pt9912/pg-change-feed/blob/main/docs/plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)). The current release provides the shared connection configuration (`ClientOptions`: server address and bearer token), a full HTTP API client surface (`PgChangeFeedHttpClient`): consumer registration/acknowledgement/position/removal, table enable/disable, status, table listing, retention, and reading changes — the nine `SPEC-018` capabilities plus `GET /changes` (`SPEC-022`) —, the gRPC live-change-stream surface (`PgChangeFeedGrpcClient`): a server-streaming read of every committed change from connection time onward, fire-and-forget with no in-stream replay (`SPEC-020`, `LH-FA-SST-008` boundary), the SSE live-change-stream surface (`PgChangeFeedSseClient`): the same fire-and-forget, no-replay stream over Server-Sent-Events instead of gRPC (`SPEC-021`, `LH-FA-SST-008` boundary), and the NATS full-content stream surface (`PgChangeFeedNatsStreamClient`): the same message schema and fire-and-forget, no-replay boundary, delivered over a NATS subject namespace instead of HTTP, with connection-level (not per-call) token authentication (`SPEC-024`, `LH-FA-SST-008` boundary) — the full four-way matrix beside C# (`PgChangeFeed.Client`) and Kotlin (`pgchangefeed-kotlin`).

If a surface you need isn't covered yet, the direct wire protocol remains fully usable on its own — no `examples/python/` reference client exists yet (`ADR-0107` §Kontext); see the [`examples`](https://github.com/pt9912/pg-change-feed/tree/main/examples) reference clients for other languages in the main repository.

## Installation

```
pip install pgchangefeed
```

## Documentation

This README intentionally does not duplicate the wire protocol documentation. For the full picture — server setup, HTTP/gRPC/SSE/NATS delivery paths, and operational guidance — see the [project README](https://github.com/pt9912/pg-change-feed/blob/main/README.md) and the technical specification (`spec/pflichtenheft.md`, `SPEC-018`) in the main repository.

## License

MIT — see [LICENSE](https://github.com/pt9912/pg-change-feed/blob/main/LICENSE).
