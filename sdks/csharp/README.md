# PG Change Feed — C# SDK

Official .NET client library for [PG Change Feed](https://github.com/pt9912/pg-change-feed), a durable change feed system for PostgreSQL built on logical replication.

This package (`PgChangeFeed.Client`) lets a .NET application consume PG Change Feed's HTTP API, gRPC change stream, SSE stream, and NATS full-content stream without implementing the wire protocol itself — see [`LH-FA-SST-009`](https://github.com/pt9912/pg-change-feed/blob/main/spec/lastenheft.md) for the requirement this SDK fulfills.

## Status

This package is at an early, pre-1.0 stage (`0.x.y`, [ADR-0106](https://github.com/pt9912/pg-change-feed/blob/main/docs/plan/adr/0106-csharp-nuget-erstes-sdk-package.md)). The current release provides the shared connection configuration (`PgChangeFeedClientOptions`: server address and bearer token), a full HTTP API client surface (`PgChangeFeedHttpClient`): consumer registration/acknowledgement/position/removal, table enable/disable, status, table listing, retention, and reading changes — the nine `SPEC-018` capabilities plus `GET /changes` (`SPEC-022`; each change carries `origin`, `wal` or `backfill`, read as `wal` when the response has none or carries JSON `null`, any other value passed through unchanged) —, the gRPC live-change-stream surface (`PgChangeFeedGrpcClient`, `StreamChangesAsync`): a server-streaming read of every committed change from connection time onward, fire-and-forget with no in-stream replay (`SPEC-020`, `LH-FA-SST-008` boundary), the SSE live-change-stream surface (`PgChangeFeedSseClient`, `StreamChangesAsync`): the same fire-and-forget, no-replay stream over Server-Sent-Events instead of gRPC (`SPEC-021`, `LH-FA-SST-008` boundary), and the NATS full-content stream surface (`PgChangeFeedNatsStreamClient`, `StreamChangesAsync`): the same message schema and fire-and-forget, no-replay boundary, delivered over a NATS subject namespace instead of HTTP, with connection-level (not per-call) token authentication (`SPEC-024`, `LH-FA-SST-008` boundary).

If a surface you need isn't covered yet, the direct wire protocol remains fully usable on its own — see the [`examples/csharp`](https://github.com/pt9912/pg-change-feed/tree/main/examples/csharp) reference clients in the main repository.

## Installation

```
dotnet add package PgChangeFeed.Client
```

## Documentation

This README intentionally does not duplicate the wire protocol documentation. For the full picture — server setup, HTTP/gRPC/SSE/NATS delivery paths, and operational guidance — see the [project README](https://github.com/pt9912/pg-change-feed/blob/main/README.md) and the technical specification (`spec/pflichtenheft.md`, `SPEC-018`/`SPEC-020`) in the main repository.

## License

MIT — see [LICENSE](https://github.com/pt9912/pg-change-feed/blob/main/LICENSE).
