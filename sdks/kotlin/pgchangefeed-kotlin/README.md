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
    implementation("io.github.pt9912:pgchangefeed-kotlin:0.6.0")
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

Apart from the Kotlin standard library, the library's own dependencies are not passed on to your compile classpath. Add the ones whose types you use — the coroutines library for the gRPC `Flow`, `grpc-api` for `Status`, `StatusException` and `Channel` (the gRPC error types and both gRPC clients use them in their public signatures), protobuf for the gRPC row images (`ByteString`), Gson for the JSON row images (`JsonElement`) — with the versions the library is built with:

```kotlin
dependencies {
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.11.0")
    implementation("io.grpc:grpc-api:1.84.0")
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

`streamChanges` takes three optional parameters, `schema`, `table` and `target`, each independent: a schema without a table matches every table of that schema, a table without a schema matches every table of that name across schemas, both set matches exactly one table, and all left `null` (the default) delivers every change of every captured table — the same filter form as `readChanges` below. `target` selects the delivery target a change is routed to: a set value delivers only the changes routed to that target, combined with `schema`/`table` as a conjunction, and a target no change carries delivers nothing and raises no error: `streamChanges(target = "eu")`. A server release that predates the parameter is expected to ignore it on the gRPC stream, which then stays unfiltered, and to answer `400` on the HTTP read path and the SSE stream; this follows from how the server reads its parameters and has not been run against such a release.

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

The SSE stream takes three optional, independent parameters, `target`, `schema` and `table`: `client.streamChanges(target = "eu")` delivers only the changes routed to that target, `client.streamChanges(schema = "public", table = "orders")` only the changes of that table, and leaving them out delivers every change. Combined, they act as a conjunction.

NATS (the address is the NATS URL, the token is the NATS stream token checked when the connection is opened, so a rejected token fails in the constructor). The subject selects what you receive: `buildSourceSubject` covers all tables of one source, `buildSubject` one table, `buildTargetSubject` one delivery target of a source (`cdc.route.<source>.<target>`), `buildSourceTargetsSubject` every delivery target of a source, and without a subject you receive every source. A change routed to a target arrives on its table subject and, with the same payload, on the target subject; a change without a target arrives on the table subject only. The tokens are validated like those of `buildSubject`; an empty or blank target is invalid (the builders have no "no filter" form, use `buildSourceSubject` for that):

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

### Manage tables and consumers over gRPC

`PgChangeFeedAdministrationClient` wraps the same eleven management/read/diagnose capabilities as `PgChangeFeedHttpClient`, over gRPC instead of HTTP — one `suspend fun` per RPC, same convenience/advanced constructor pair as `PgChangeFeedGrpcClient` above. Requests and responses are the generated `cdc.administration.v1` protobuf messages directly, not a separate DTO type:

```kotlin
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedAdministrationClient
import cdc.administration.v1.AdministrationOuterClass.EnableTableRequest
import cdc.administration.v1.AdministrationOuterClass.ListTablesRequest
import java.net.URI
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val options = PgChangeFeedClientOptions(URI("http://feed.example.com:9090"), "<admin token>")
    PgChangeFeedAdministrationClient(options).use { admin ->
        admin.enableTable(
            EnableTableRequest.newBuilder()
                .setSource("my-source").setSchema("public").setTable("orders")
                .setTableId("orders").setSchemaVersionId("v1").setVersion(1).setPublication("my_publication")
                .build(),
        )

        val tables = admin.listTables(ListTablesRequest.newBuilder().setSource("my-source").setPublication("my_publication").build())
        for (t in tables.tablesList) {
            println("${t.schema}.${t.table}")
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
| `readChanges(source, schema, table, from, to, limit, target)` | Reads stored changes of a source, optionally for one schema and table, for the range `[from, to)` of commit positions and for one delivery target (`target`, combined with `schema`/`table` as a conjunction; a target no change carries returns an empty list). Only `source` is required. | reader |

Request and response classes live in `io.github.pt9912.pgchangefeed.http.model`.

`EnableTableRequest` names the table and the ids it is captured under: `source`, `schema` and `table` identify the table, `publication` is the PostgreSQL publication of the source, `tableId` is the id you give the table (every change of the table carries it as `sourceTableId`), and `schemaVersionId` and `version` (a number, 1 or higher) identify the table's schema version (changes carry the id as `schemaVersion`).

The live streams each have one method:

| Class | Method | What it does |
|---|---|---|
| `PgChangeFeedGrpcClient(options)` | `streamChanges(schema, table, target)` | Opens the gRPC stream and returns a `Flow` of the generated `Change` messages. `schema`/`table`/`target` are optional and independent, all left `null` delivers every change. `close()` shuts down the channel the client owns. |
| `PgChangeFeedSseClient(httpClient, options)` | `streamChanges(target, schema, table)` | Opens `GET /changes/stream` and returns a `Sequence` of `io.github.pt9912.pgchangefeed.sse.model.Change` objects. `target`/`schema`/`table` are optional and independent, all left `null` delivers every change. |
| `PgChangeFeedNatsStreamClient(options)` | `streamChanges(subject)` | Subscribes to a subject and returns a `Sequence` of `io.github.pt9912.pgchangefeed.nats.model.Change` objects. `close()` closes the connection the client owns. |

`PgChangeFeedAdministrationClient(options)` wraps the eleven RPCs of the `Administration` gRPC service — the same capabilities as the HTTP table above, over gRPC. Requests and responses are the generated `cdc.administration.v1` protobuf messages, used directly (no separate model type):

| Method | What it does | Token |
|---|---|---|
| `registerConsumer(request)` | Registers a consumer. Registering an existing consumer changes nothing (`alreadyRegistered` is true). | admin |
| `acknowledgeConsumer(request)` | Stores the consumer's position. Repeating the same position has no effect; an earlier position is rejected. | admin |
| `getConsumerPosition(request)` | Reads the stored position (`offset`, and `acknowledged`, which is false for a consumer that never acknowledged). | reader |
| `removeConsumer(request)` | Removes a consumer (`removed` is false for one that was never registered). | admin |
| `enableTable(request)` | Starts capturing a table. A table that is already captured changes nothing (`alreadyEnabled` is true); a table that does not exist throws `PgChangeFeedGrpcNotFoundException`. | admin |
| `disableTable(request)` | Stops capturing a table. `retained` reports that changes already stored for it remain. | admin |
| `getTableStatus(request)` | Tells whether a table is captured (`enabled`) or no longer captured with stored changes remaining (`retained`). | reader |
| `listTables(request)` | Lists the captured tables (`tablesList`) and the tables whose stored changes remain (`retainedList`), each a `SourceTable`. | reader |
| `runRetention(request)` | Deletes stored changes older than `minAgeNanos` that every consumer with a stored position has passed; returns the number deleted (`deleted`). | admin |
| `readChanges(request)` | Reads stored changes of a source, optionally for one schema and table, for the range `[from, to)` of commit positions and for one delivery target (`target`, empty is no filter); `changesList` is a list of `ChangeRecord`. | reader |
| `diagnose(request)` | Reads the operational diagnose report (heartbeat, capture lag, per-consumer lag, retention blocker, storage, backfill status). | reader |

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
        println("${e.statusCode} ${e.messageCode} ${e.message}")
    }
}
```

`messageCode` is the message code the server attaches to an error (`PCF-<E|W|I><4 digits>`, for example `PCF-E8025` for a table that does not exist in the source database), read from the `code` field of the error body and passed through unchanged. Match on it instead of parsing the message. It is `null` when the server sent no code: authentication errors (`401`, `403`) carry none, a server release that predates message codes sends none, and so does a body that is not JSON or whose `code` is not a non-empty string. The error text is the `error` field of the body when that is a JSON string and the raw body otherwise; it does not depend on `code`, and `code` does not depend on the type of `error`.

The gRPC stream reports a missing or unknown token as an `io.grpc.StatusException` with status `UNAUTHENTICATED`, thrown while the `Flow` is collected. The gRPC stream client passes on the plain `StatusException` for every failure; a message code sent with such an error stays readable through the status detail of the exception itself (a `google.rpc.ErrorInfo` entry of domain `pg-change-feed`, its `reason` is the code). The NATS client passes on the exception of the NATS client library (`io.nats.client`) when the server rejects the token, and throws `PgChangeFeedNatsMalformedMessageException` when a message cannot be read.

A `PgChangeFeedAdministrationClient` call that the server answers with a non-`OK` gRPC status throws a subclass of the sealed class `PgChangeFeedGrpcException`, which carries the `io.grpc.Status.Code`; the original `io.grpc.StatusException` is always the `cause`.

| Exception | gRPC status |
|---|---|
| `PgChangeFeedGrpcInvalidArgumentException` | `INVALID_ARGUMENT` — invalid request or a violated rule. |
| `PgChangeFeedGrpcUnauthenticatedException` | `UNAUTHENTICATED` — token missing or unknown. |
| `PgChangeFeedGrpcPermissionDeniedException` | `PERMISSION_DENIED` — known token whose class may not call this RPC (for example a reader token on an admin RPC). |
| `PgChangeFeedGrpcNotFoundException` | `NOT_FOUND` — the table does not exist in the source database (`enableTable`, `disableTable`, `getTableStatus`). |
| `PgChangeFeedGrpcInternalException` | `INTERNAL` — unexpected error inside the server. |
| `PgChangeFeedGrpcUnexpectedStatusException` | any other non-`OK` status. |

Every one of these exceptions carries `messageCode`, the message code of the server read from the `reason` of the `google.rpc.ErrorInfo` status detail of domain `pg-change-feed`, whatever the status (`null` when the server sent none, as for `UNAUTHENTICATED` and `PERMISSION_DENIED`); for example `PgChangeFeedGrpcNotFoundException.messageCode` is `PCF-E8025` for `enableTable` on a table that does not exist. The `errorCode` field of the `diagnose` response (the message code of the process error state, empty in normal operation) arrives unchanged from the generated message.

## Reading versus streaming

- **Reading over HTTP** (`readChanges`) is asking: you name a range, the server answers from the changes it has stored. You can read the same range again, and with a registered consumer you can carry on after a restart exactly where you stopped. Changes stay readable until the retention removes them.
- **The live streams** (gRPC, SSE, NATS) are pushing: you get every change committed after you connected, in commit order, with the full row content. There is no delivery guarantee and no replay. A change committed while you were disconnected, or while you read too slowly, does not arrive on the stream. Use the stream to react quickly and `readChanges` to catch up on what it missed.
- The gRPC stream can be filtered by schema, table and delivery target (`streamChanges(schema, table, target)`); the SSE stream can be filtered by delivery target, schema and table (`streamChanges(target, schema, table)`); on the NATS stream the subject chooses the source, the table or the delivery target.

## Upgrading

- **0.5.0** — `streamChanges` of the SSE client gained the optional parameters `schema` and `table` after `target`. Source code that calls it keeps compiling unchanged. A class file that was compiled against an earlier version calls the old method signature, which no longer exists, and has to be recompiled against 0.5.0.

## More

- [Project README](https://github.com/pt9912/pg-change-feed/blob/main/README.md)
- [User manual](https://github.com/pt9912/pg-change-feed/blob/main/docs/user/benutzerhandbuch.md) — setting up and operating the server, tokens, delivery paths (in German)
- [Reference clients](https://github.com/pt9912/pg-change-feed/tree/main/examples/kotlin) for each delivery path

## License

MIT — see [LICENSE](https://github.com/pt9912/pg-change-feed/blob/main/LICENSE).
