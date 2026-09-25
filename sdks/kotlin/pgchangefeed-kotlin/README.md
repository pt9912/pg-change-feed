# PG Change Feed — Kotlin client

Kotlin/JVM client library for [PG Change Feed](https://github.com/pt9912/pg-change-feed), a server that records every INSERT, UPDATE and DELETE of selected PostgreSQL tables and makes these changes available over HTTP, gRPC, Server-Sent Events (SSE) and NATS.

With this package (`pgchangefeed-kotlin`) a Kotlin application can

- read the recorded changes and keep track of how far it has processed them,
- manage which tables the server captures, and
- receive changes live, as they happen,

without implementing any of the wire protocols itself.

Version 0.x — the API can still change between releases.

## Installation

Requires Java 21 or newer. The package is not on Maven Central; it is published on Cloudsmith and on GitHub Packages.

### From Cloudsmith (no account, no token)

The Cloudsmith repository is public: no account and no token are needed to read it. Add it in `settings.gradle.kts`, together with Maven Central, where the library's own dependencies come from:

```kotlin
dependencyResolutionManagement {
    repositories {
        mavenCentral()
        maven { url = uri("https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/") }
    }
}
```

Then declare the dependency in `build.gradle.kts`:

```kotlin
dependencies {
    implementation("io.github.pt9912:pgchangefeed-kotlin:0.2.2")
}
```

Package repository hosting is graciously provided by [Cloudsmith](https://cloudsmith.com), free of charge for open-source projects.

### From GitHub Packages (needs a token)

The package is also published on [GitHub Packages](https://maven.pkg.github.com/pt9912/pg-change-feed). **GitHub Packages always requires authentication to read a package, even a public one.** You need a GitHub account and a classic personal access token (PAT) with the `read:packages` scope. Use the same dependency declaration as above and add the registry with your credentials in `settings.gradle.kts` instead of the Cloudsmith repository (keep `mavenCentral()`):

```kotlin
dependencyResolutionManagement {
    repositories {
        mavenCentral()
        maven {
            name = "GitHubPackages"
            url = uri("https://maven.pkg.github.com/pt9912/pg-change-feed")
            credentials {
                username = System.getenv("GITHUB_ACTOR") // your GitHub user name
                password = System.getenv("GITHUB_TOKEN") // a classic PAT with read:packages
            }
        }
    }
}
```

### Dependencies of the library

Apart from the Kotlin standard library, the library's own dependencies are not passed on to your compile classpath. Add the ones whose types you use — the coroutines library for the gRPC `Flow`, protobuf for the gRPC row images (`ByteString`), Gson for the JSON row images (`JsonElement`) — with the versions the library is built with:

```kotlin
dependencies {
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.11.0")
    implementation("com.google.protobuf:protobuf-java:4.36.2")
    implementation("com.google.code.gson:gson:2.14.0")
}
```

## Quick start

A client needs the address of the PG Change Feed server and a token. The server knows two token classes: a *reader* token for read-only calls and an *admin* token for calls that change something (registering consumers, acknowledging positions, enabling tables, running the retention). The admin token also covers all reader calls. Address, source id and tokens come from whoever operates the server. The server does not serve TLS itself, so the examples use unencrypted addresses.

### Read changes and remember your position

A *consumer* is a named reader whose progress the server remembers. You register it once, read changes, and acknowledge the position of the last change you have processed. A position is the `commitPosition` of a change; `readChanges` reads from `from` (inclusive) up to `to` (exclusive). A consumer that has never acknowledged reports offset 0.

```kotlin
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient
import io.github.pt9912.pgchangefeed.http.model.AcknowledgeConsumerRequest
import io.github.pt9912.pgchangefeed.http.model.RegisterConsumerRequest
import java.net.URI
import java.net.http.HttpClient

fun main() {
    val options = PgChangeFeedClientOptions(URI("http://feed.example.com:8090"), "<admin token>")
    val client = PgChangeFeedHttpClient(HttpClient.newHttpClient(), options)

    client.registerConsumer(RegisterConsumerRequest("billing", "Billing service"))

    val position = client.getConsumerPosition("billing")
    val result = client.readChanges("my-source", from = position.offset + 1)
    for (change in result.changes) {
        println("${change.commitPosition} ${change.operation} ${change.schema}.${change.table} ${change.newImage}")
    }

    if (result.changes.isNotEmpty()) {
        client.acknowledgeConsumer(
            AcknowledgeConsumerRequest("billing", "my-source", result.changes.last().commitPosition),
        )
    }
}
```

`limit` cuts rows, not positions: if one commit position carries more changes than `limit`, continuing from that position plus one skips the rest of it. Leave `limit` out (or set it generously) where a single position can carry many changes.

### Receive changes live

The three live streams deliver every change committed after you connect. Each surface has its own client; all take the same `PgChangeFeedClientOptions`.

gRPC (the address is the `http://host:port` URL of the gRPC endpoint; `streamChanges()` returns a `kotlinx.coroutines.flow.Flow` of the generated `Change` protobuf messages, row images are JSON in a `ByteString`, empty when there is none; the convenience constructor opens its own plaintext channel, `PgChangeFeedGrpcClient(channel, options)` takes an `io.grpc.Channel` you configure yourself):

```kotlin
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient
import java.net.URI
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val options = PgChangeFeedClientOptions(URI("http://feed.example.com:9090"), "<reader token>")
    PgChangeFeedGrpcClient(options).use { client ->
        client.streamChanges().collect { change ->
            println("${change.operation} ${change.schema}.${change.table} ${change.newImage.toStringUtf8()}")
        }
    }
}
```

Server-Sent Events (the address is the HTTP base URL; `streamChanges()` returns a `Sequence` that blocks while it waits for the next change):

```kotlin
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient
import java.net.URI
import java.net.http.HttpClient

fun main() {
    val options = PgChangeFeedClientOptions(URI("http://feed.example.com:8090"), "<reader token>")
    val client = PgChangeFeedSseClient(HttpClient.newHttpClient(), options)
    for (change in client.streamChanges()) {
        println("${change.operation} ${change.schema}.${change.table} ${change.newImage}")
    }
}
```

NATS (the address is the NATS URL, the token is the NATS stream token checked when the connection is opened, so a rejected token fails in the constructor). The subject selects what you receive: `buildSourceSubject` covers all tables of one source, `buildSubject` one table, and without a subject you receive every source:

```kotlin
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.nats.PgChangeFeedNatsStreamClient
import java.net.URI

fun main() {
    val options = PgChangeFeedClientOptions(URI("nats://feed.example.com:4222"), "<NATS stream token>")
    PgChangeFeedNatsStreamClient(options).use { client ->
        val subject = PgChangeFeedNatsStreamClient.buildSubject("my-source", "public", "orders")
        for (change in client.streamChanges(subject)) {
            println("${change.operation} ${change.schema}.${change.table} ${change.newImage}")
        }
    }
}
```

## API overview

`PgChangeFeedHttpClient(httpClient, options)` wraps the HTTP API. The `java.net.http.HttpClient` you pass in stays yours; the library never closes it.

| Method | What it does | Token |
|---|---|---|
| `registerConsumer(request)` | Registers a consumer. Registering an existing consumer changes nothing (`alreadyRegistered` is true). | admin |
| `acknowledgeConsumer(request)` | Stores the consumer's position. Repeating the same position has no effect; a position before the stored one is rejected. | admin |
| `getConsumerPosition(consumerId)` | Reads the stored position (`offset`, and `acknowledged`, which is false for a consumer that never acknowledged). | reader |
| `removeConsumer(consumerId)` | Removes a consumer. | admin |
| `enableTable(request)` | Starts capturing a table (see below). A table that is already captured changes nothing (`alreadyEnabled` is true); a table that does not exist throws `PgChangeFeedNotFoundException`. | admin |
| `disableTable(request)` | Stops capturing a table. `retained` reports that changes already stored for it remain. | admin |
| `getStatus(source, schema, table, publication)` | Tells whether a table is captured (`enabled`) or no longer captured with stored changes remaining (`retained`). | reader |
| `listTables(source, publication)` | Lists the captured tables and the tables whose stored changes remain. | reader |
| `runRetention(request)` | Deletes stored changes older than `minAgeNanos` that every consumer with a stored position has passed; returns the number deleted. | admin |
| `readChanges(source, schema, table, from, to, limit)` | Reads stored changes of a source, optionally for one schema and table and for the range `[from, to)` of commit positions. Only `source` is required. | reader |

Request and response classes live in `io.github.pt9912.pgchangefeed.http.model`.

`EnableTableRequest` names the table and the ids it is captured under: `source`, `schema` and `table` identify the table, `publication` is the PostgreSQL publication of the source, `tableId` is the id you give the table (every change of the table carries it as `sourceTableId`), and `schemaVersionId` and `version` (a number, 1 or higher) identify the table's schema version (changes carry the id as `schemaVersion`).

The live streams each have one method:

| Class | Method | What it does |
|---|---|---|
| `PgChangeFeedGrpcClient(options)` | `streamChanges()` | Opens the gRPC stream and returns a `Flow` of the generated `Change` messages. `close()` shuts down the channel the client owns. |
| `PgChangeFeedSseClient(httpClient, options)` | `streamChanges()` | Opens `GET /changes/stream` and returns a `Sequence` of `io.github.pt9912.pgchangefeed.sse.model.Change` objects. |
| `PgChangeFeedNatsStreamClient(options)` | `streamChanges(subject)` | Subscribes to a subject and returns a `Sequence` of `io.github.pt9912.pgchangefeed.nats.model.Change` objects. `close()` closes the connection the client owns. |

## The change object

`readChanges` returns `io.github.pt9912.pgchangefeed.http.model.Change` data classes:

| Property | Meaning |
|---|---|
| `commitPosition` | Position of the committed source transaction that carried the change; the value you read ranges by and acknowledge. |
| `changeId` | Unique id of the change. |
| `transactionId` | Id of the source transaction. |
| `sourceTableId` | Id of the captured table. |
| `schema`, `table` | Schema and name of the table. |
| `sequence` | Order of the row change within its transaction. |
| `operation` | `INSERT`, `UPDATE` or `DELETE`. |
| `oldImage` | Row values before the change as a Gson `JsonElement`; `JsonNull` for an INSERT (`oldImage?.isJsonNull == true`, not a Kotlin `null`). For UPDATE and DELETE it holds what PostgreSQL provides for the table's replica identity. |
| `newImage` | Row values after the change as a `JsonElement`; `JsonNull` for a DELETE. |
| `schemaVersion` | Version of the table schema the change was captured with. |
| `committedAt` | Commit time of the source transaction (RFC 3339, UTC). |
| `origin` | `wal` for a change captured live from the database, `backfill` for a change that was taken from the existing table contents. A response without the field, or with JSON `null`, reads as `wal`; any other value is passed through unchanged. |

The live streams deliver change objects with ten properties: `changeId`, `transactionId`, `sourceTableId`, `sequence`, `operation`, `oldImage`, `newImage`, `schemaVersion`, `schema`, `table`. They carry no `commitPosition`, `committedAt` or `origin`.

## Error handling

An HTTP call that the server answers with an error status throws a subclass of the sealed class `PgChangeFeedException`, which carries the HTTP `statusCode`. The same exceptions are thrown when the SSE stream cannot be opened. All live in `io.github.pt9912.pgchangefeed.http`. Connection failures and timeouts are not converted: they surface as the `java.io.IOException` of `java.net.http.HttpClient` (for example `java.net.http.HttpTimeoutException`). Every HTTP call has a timeout of ten seconds; the SSE stream has none.

| Exception | When |
|---|---|
| `PgChangeFeedBadRequestException` | 400 — invalid request or a violated rule. |
| `PgChangeFeedUnauthorizedException` | 401 — token missing or unknown. |
| `PgChangeFeedForbiddenException` | 403 — known token whose class may not call this endpoint (for example a reader token on an admin call). |
| `PgChangeFeedNotFoundException` | 404 — the table does not exist in the source database (`enableTable`, `disableTable`, `getStatus`). |
| `PgChangeFeedServerErrorException` | 500 — unexpected error inside the server. |
| `PgChangeFeedUnexpectedStatusException` | any other non-success status. |
| `PgChangeFeedMalformedResponseException` | a success response or SSE event whose content cannot be read. |

```kotlin
import io.github.pt9912.pgchangefeed.http.PgChangeFeedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedForbiddenException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnauthorizedException

fun listTablesReportingErrors(client: PgChangeFeedHttpClient) {
    try {
        client.listTables("my-source", "my_publication")
    } catch (e: PgChangeFeedUnauthorizedException) {
        println("token missing or unknown")
    } catch (e: PgChangeFeedForbiddenException) {
        println("token is not allowed to call this endpoint")
    } catch (e: PgChangeFeedException) {
        println("${e.statusCode} ${e.message}")
    }
}
```

The gRPC stream reports a missing or unknown token as an `io.grpc.StatusException` with status `UNAUTHENTICATED`, thrown while the `Flow` is collected. The NATS client passes on the exception of the NATS client library (`io.nats.client`) when the server rejects the token, and throws `PgChangeFeedNatsMalformedMessageException` when a message cannot be read.

## Reading versus streaming

- **Reading over HTTP** (`readChanges`) is asking: you name a range, the server answers from the changes it has stored. You can read the same range again, and with a registered consumer you can carry on after a restart exactly where you stopped. Changes stay readable until the retention removes them.
- **The live streams** (gRPC, SSE, NATS) are pushing: you get every change committed after you connected, in commit order, with the full row content. There is no delivery guarantee and no replay. A change committed while you were disconnected, or while you read too slowly, does not arrive on the stream. Use the stream to react quickly and `readChanges` to catch up on what it missed.
- The gRPC and SSE streams cannot be filtered by table; on the NATS stream the subject chooses the source or the table.

## More

- [Project README](https://github.com/pt9912/pg-change-feed/blob/main/README.md)
- [User manual](https://github.com/pt9912/pg-change-feed/blob/main/docs/user/benutzerhandbuch.md) — setting up and operating the server, tokens, delivery paths (in German)
- [Reference clients](https://github.com/pt9912/pg-change-feed/tree/main/examples/kotlin) for each delivery path

## License

MIT — see [LICENSE](https://github.com/pt9912/pg-change-feed/blob/main/LICENSE).
