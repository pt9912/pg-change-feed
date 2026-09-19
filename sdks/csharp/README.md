# PG Change Feed — C# SDK

Official .NET client library for [PG Change Feed](https://github.com/pt9912/pg-change-feed), a durable change feed system for PostgreSQL built on logical replication.

This package (`PgChangeFeed.Client`) lets a .NET application consume PG Change Feed's HTTP API and gRPC change stream without implementing the wire protocol itself — see [`LH-FA-SST-009`](https://github.com/pt9912/pg-change-feed/blob/main/spec/lastenheft.md) for the requirement this SDK fulfills.

## Status

This package is at an early, pre-1.0 stage (`0.x.y`, [ADR-0106](https://github.com/pt9912/pg-change-feed/blob/main/docs/plan/adr/0106-csharp-nuget-erstes-sdk-package.md)). The current release provides the shared connection configuration (`PgChangeFeedClientOptions`: server address and bearer token). Client surfaces for the HTTP API and the gRPC change stream are added by follow-up releases.

If a surface you need isn't covered yet, the direct wire protocol remains fully usable on its own — see the [`examples/csharp`](https://github.com/pt9912/pg-change-feed/tree/main/examples/csharp) reference clients in the main repository.

## Installation

```
dotnet add package PgChangeFeed.Client
```

## Documentation

This README intentionally does not duplicate the wire protocol documentation. For the full picture — server setup, HTTP/gRPC/SSE/NATS delivery paths, and operational guidance — see the [project README](https://github.com/pt9912/pg-change-feed/blob/main/README.md) and the technical specification (`spec/pflichtenheft.md`, `SPEC-018`/`SPEC-020`) in the main repository.

## License

MIT — see [LICENSE](https://github.com/pt9912/pg-change-feed/blob/main/LICENSE).
