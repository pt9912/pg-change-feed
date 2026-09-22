# PG Change Feed — Kotlin SDK

Official Kotlin/JVM client library for [PG Change Feed](https://github.com/pt9912/pg-change-feed), a durable change feed system for PostgreSQL built on logical replication.

This package (`pgchangefeed-kotlin`) lets a Kotlin or Java application consume PG Change Feed's HTTP API, gRPC change stream and SSE change stream without implementing the wire protocol itself — see [`LH-FA-SST-009`](https://github.com/pt9912/pg-change-feed/blob/main/spec/lastenheft.md) for the requirement this SDK fulfills.

## Status

This package is at an early, pre-1.0 stage (`0.x.y`, [ADR-0109](https://github.com/pt9912/pg-change-feed/blob/main/docs/plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md)). The current release provides the shared connection configuration (`PgChangeFeedClientOptions`: server address and bearer token), a full HTTP API client surface (`PgChangeFeedHttpClient`): consumer registration/acknowledgement/position/removal, table enable/disable, status, table listing, retention, and reading changes — the nine `SPEC-018` capabilities plus `GET /changes` (`SPEC-022`) — with typed request/response data classes and a typed exception hierarchy for `400`/`401`/`403`/`404`/`500`, the gRPC live-change-stream surface (`PgChangeFeedGrpcClient`, `SPEC-020`): `streamChanges()` opens the `ChangeStream/StreamChanges` server-streaming RPC and yields a `kotlinx.coroutines.flow.Flow` of the generated `Change` message — fire-and-forget, no replay, no table-granular filtering (`LH-FA-SST-008` Boundary) — and the SSE live-change-stream surface (`PgChangeFeedSseClient`, `SPEC-021`): `streamChanges()` opens `GET /changes/stream` and yields a `Sequence` of the ten `SPEC-021` message fields, reusing the same typed exception hierarchy as `PgChangeFeedHttpClient`, with the same fire-and-forget/no-replay boundary. NATS-Vollinhalt delivery remains out of scope for this package's planned first full release and would be added by a later follow-up (`ADR-0109` Festlegung 1).

If a surface you need isn't covered yet, the direct wire protocol remains fully usable on its own — see the [`examples/kotlin`](https://github.com/pt9912/pg-change-feed/tree/main/examples/kotlin) reference clients in the main repository.

## Installation

This package is distributed via [GitHub Packages](https://maven.pkg.github.com/pt9912/pg-change-feed) (Gradle/Maven registry), not Maven Central (`ADR-0109` Festlegung 2). **GitHub Packages always requires authentication to read a package, even a public one** — unlike Maven Central, NuGet, or PyPI. You need a GitHub account and a classic personal access token (PAT) with the `read:packages` scope.

Add the registry and your credentials in `settings.gradle.kts`:

```kotlin
dependencyResolutionManagement {
    repositories {
        maven {
            name = "GitHubPackages"
            url = uri("https://maven.pkg.github.com/pt9912/pg-change-feed")
            credentials {
                username = System.getenv("GITHUB_ACTOR")
                password = System.getenv("GITHUB_TOKEN") // a classic PAT with read:packages
            }
        }
    }
}
```

Then declare the dependency in `build.gradle.kts`:

```kotlin
dependencies {
    implementation("io.github.pt9912:pgchangefeed-kotlin:0.1.0")
}
```

## Documentation

This README intentionally does not duplicate the wire protocol documentation. For the full picture — server setup, HTTP/gRPC/SSE/NATS delivery paths, and operational guidance — see the [project README](https://github.com/pt9912/pg-change-feed/blob/main/README.md) and the technical specification (`spec/pflichtenheft.md`, `SPEC-018`/`SPEC-020`) in the main repository.

## License

MIT — see [LICENSE](https://github.com/pt9912/pg-change-feed/blob/main/LICENSE).
