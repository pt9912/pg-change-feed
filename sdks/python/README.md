# PG Change Feed — Python SDK

Official Python client library for [PG Change Feed](https://github.com/pt9912/pg-change-feed), a durable change feed system for PostgreSQL built on logical replication.

This package (`pgchangefeed`) lets a Python application consume PG Change Feed's HTTP API without implementing the wire protocol itself — see [`LH-FA-SST-009`](https://github.com/pt9912/pg-change-feed/blob/main/spec/lastenheft.md) for the requirement this SDK fulfills.

## Status

This package is at an early, pre-1.0 stage (`0.x.y`, [ADR-0107](https://github.com/pt9912/pg-change-feed/blob/main/docs/plan/adr/0107-python-pypi-zweites-sdk-package.md)). The current release provides only the shared connection configuration (`ClientOptions`: server address and bearer token) — the HTTP API client surface itself (the nine `SPEC-018` capabilities plus `GET /changes`, `SPEC-022`) is added by a follow-up release. gRPC, SSE and NATS-vollinhalt delivery remain out of scope for this package's first release (`ADR-0107` Festlegung 1).

If a surface you need isn't covered yet, the direct wire protocol remains fully usable on its own — no `examples/python/` reference client exists yet (`ADR-0107` §Kontext); see the [`examples`](https://github.com/pt9912/pg-change-feed/tree/main/examples) reference clients for other languages in the main repository.

## Installation

```
pip install pgchangefeed
```

## Documentation

This README intentionally does not duplicate the wire protocol documentation. For the full picture — server setup, HTTP/gRPC/SSE/NATS delivery paths, and operational guidance — see the [project README](https://github.com/pt9912/pg-change-feed/blob/main/README.md) and the technical specification (`spec/pflichtenheft.md`, `SPEC-018`) in the main repository.

## License

MIT — see [LICENSE](https://github.com/pt9912/pg-change-feed/blob/main/LICENSE).
