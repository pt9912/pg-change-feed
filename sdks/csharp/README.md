# PG Change Feed — .NET client

.NET client library for [PG Change Feed](https://github.com/pt9912/pg-change-feed), a server that records every INSERT, UPDATE and DELETE of selected PostgreSQL tables and makes these changes available over HTTP, gRPC, Server-Sent Events (SSE) and NATS.

With this package (`PgChangeFeed.Client`) a .NET application can

- read the recorded changes and keep track of how far it has processed them,
- manage which tables the server captures, and
- receive changes live, as they happen,

without implementing any of the wire protocols itself.

Version 0.x — the API can still change between releases.

## Installation

```
dotnet add package PgChangeFeed.Client
```

Requires .NET 10 or newer. The examples below assume implicit usings, the default of current .NET project templates.

## Quick start

A client needs the address of the PG Change Feed server and a token. The server knows two token classes: a *reader* token for read-only calls and an *admin* token for calls that change something (registering consumers, acknowledging positions, enabling tables, running the retention). The admin token also covers all reader calls. Address, source id and tokens come from whoever operates the server.

### Read changes and remember your position

A *consumer* is a named reader whose progress the server remembers. You register it once, read changes, and acknowledge the position of the last change you have processed. A position is the `CommitPosition` of a change; `ReadChangesAsync` reads from `from` (inclusive) up to `to` (exclusive). A consumer that has never acknowledged reports offset 0.

```csharp
using PgChangeFeed.Client;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Http.Models;

using var http = new HttpClient();
var options = new PgChangeFeedClientOptions(new Uri("http://feed.example.com:8090"), "<admin token>");
var client = new PgChangeFeedHttpClient(http, options);

await client.RegisterConsumerAsync(new RegisterConsumerRequest("billing", "Billing service"));

var position = await client.GetConsumerPositionAsync("billing");
var result = await client.ReadChangesAsync("my-source", from: (long)position.Offset + 1);
foreach (var change in result.Changes)
{
    Console.WriteLine($"{change.CommitPosition} {change.Operation} {change.Schema}.{change.Table} {change.NewImage}");
}

if (result.Changes.Count > 0)
{
    await client.AcknowledgeConsumerAsync(
        new AcknowledgeConsumerRequest("billing", "my-source", (ulong)result.Changes[^1].CommitPosition));
}
```

`limit` cuts rows, not positions: if one commit position carries more changes than `limit`, continuing from that position plus one skips the rest of it. Leave `limit` out (or set it generously) where a single position can carry many changes.

### Receive changes live

The three live streams deliver every change committed after you connect. Each surface has its own client; all take the same `PgChangeFeedClientOptions` and return an `IAsyncEnumerable` you read with `await foreach`. Every `StreamChangesAsync` accepts a `CancellationToken` to end the stream.

gRPC (the address is the `http://host:port` URL of the gRPC endpoint; the messages are the generated `Cdc.Stream.V1.Change` protobuf messages, row images are JSON in a `ByteString`):

```csharp
using PgChangeFeed.Client;
using PgChangeFeed.Client.Grpc;

var options = new PgChangeFeedClientOptions(new Uri("http://feed.example.com:9090"), "<reader token>");
using var client = new PgChangeFeedGrpcClient(options);
await foreach (var change in client.StreamChangesAsync())
{
    Console.WriteLine($"{change.Operation} {change.Schema}.{change.Table} {change.NewImage.ToStringUtf8()}");
}
```

Server-Sent Events (the address is the HTTP base URL; the request timeout must be off for a long-lived stream):

```csharp
using PgChangeFeed.Client;
using PgChangeFeed.Client.Sse;

using var http = new HttpClient { Timeout = Timeout.InfiniteTimeSpan };
var options = new PgChangeFeedClientOptions(new Uri("http://feed.example.com:8090"), "<reader token>");
var client = new PgChangeFeedSseClient(http, options);
await foreach (var change in client.StreamChangesAsync())
{
    Console.WriteLine($"{change.Operation} {change.Schema}.{change.Table} {change.NewImage}");
}
```

NATS (the address is the NATS URL, the token is the NATS stream token checked when the connection is opened). The subject selects what you receive: `BuildSourceSubject` covers all tables of one source, `BuildSubject` one table, and without a subject you receive every source:

```csharp
using PgChangeFeed.Client;
using PgChangeFeed.Client.Nats;

var options = new PgChangeFeedClientOptions(new Uri("nats://feed.example.com:4222"), "<NATS stream token>");
await using var client = new PgChangeFeedNatsStreamClient(options);
var subject = PgChangeFeedNatsStreamClient.BuildSubject("my-source", "public", "orders");
await foreach (var change in client.StreamChangesAsync(subject))
{
    Console.WriteLine($"{change.Operation} {change.Schema}.{change.Table} {change.NewImage}");
}
```

## API overview

`PgChangeFeedHttpClient(httpClient, options)` wraps the HTTP API. The `HttpClient` you pass in stays yours; the library never disposes it. Every method also takes an optional `CancellationToken`.

| Method | What it does | Token |
|---|---|---|
| `RegisterConsumerAsync(request)` | Registers a consumer. Registering an existing consumer changes nothing (`AlreadyRegistered` is true). | admin |
| `AcknowledgeConsumerAsync(request)` | Stores the consumer's position. Repeating the same position has no effect; a position before the stored one is rejected. | admin |
| `GetConsumerPositionAsync(consumerId)` | Reads the stored position (`Offset`, and `Acknowledged`, which is false for a consumer that never acknowledged). | reader |
| `RemoveConsumerAsync(consumerId)` | Removes a consumer. | admin |
| `EnableTableAsync(request)` | Starts capturing a table. | admin |
| `DisableTableAsync(request)` | Stops capturing a table. `Retained` reports that changes already stored for it remain. | admin |
| `GetStatusAsync(source, schema, table, publication)` | Tells whether a table is captured (`Enabled`) or no longer captured with stored changes remaining (`Retained`). | reader |
| `ListTablesAsync(source, publication)` | Lists the captured tables and the tables whose stored changes remain. | reader |
| `RunRetentionAsync(request)` | Deletes stored changes older than `MinAgeNanos` that every consumer with a stored position has passed; returns the number deleted. | admin |
| `ReadChangesAsync(source, schema, table, from, to, limit)` | Reads stored changes of a source, optionally for one schema and table and for the range `[from, to)` of commit positions. Only `source` is required. | reader |

Request and response types live in `PgChangeFeed.Client.Http.Models`.

The live streams each have one method:

| Class | Method | What it does |
|---|---|---|
| `PgChangeFeedGrpcClient(options)` | `StreamChangesAsync(cancellationToken)` | Opens the gRPC stream and yields generated `Cdc.Stream.V1.Change` messages. The client owns its channel; dispose it when done. |
| `PgChangeFeedSseClient(httpClient, options)` | `StreamChangesAsync(cancellationToken)` | Opens `GET /changes/stream` and yields `PgChangeFeed.Client.Sse.Models.Change` objects. |
| `PgChangeFeedNatsStreamClient(options)` | `StreamChangesAsync(subject, cancellationToken)` | Subscribes to a subject and yields `PgChangeFeed.Client.Nats.Models.Change` objects. Dispose the client when done. |

## The change object

`ReadChangesAsync` returns `PgChangeFeed.Client.Http.Models.Change` records:

| Property | Meaning |
|---|---|
| `CommitPosition` | Position of the committed source transaction that carried the change; the value you read ranges by and acknowledge. |
| `ChangeId` | Unique id of the change. |
| `TransactionId` | Id of the source transaction. |
| `SourceTableId` | Id of the captured table. |
| `Schema`, `Table` | Schema and name of the table. |
| `Sequence` | Order of the row change within its transaction. |
| `Operation` | `INSERT`, `UPDATE` or `DELETE`. |
| `OldImage` | Row values before the change as a `JsonElement`; `null` for an INSERT. For UPDATE and DELETE it holds what PostgreSQL provides for the table's replica identity. |
| `NewImage` | Row values after the change as a `JsonElement`; `null` for a DELETE. |
| `SchemaVersion` | Version of the table schema the change was captured with. |
| `CommittedAt` | Commit time of the source transaction (RFC 3339, UTC). |
| `Origin` | `wal` for a change captured live from the database, `backfill` for a change that was taken from the existing table contents. A response without the field, or with JSON `null`, reads as `wal`; any other value is passed through unchanged. |

The live streams deliver change objects with ten properties: `ChangeId`, `TransactionId`, `SourceTableId`, `Sequence`, `Operation`, `OldImage`, `NewImage`, `SchemaVersion`, `Schema`, `Table`. They carry no `CommitPosition`, `CommittedAt` or `Origin`.

## Error handling

Every failing HTTP call throws a subclass of `PgChangeFeedException`, which carries the HTTP `StatusCode`. The same exceptions are thrown when the SSE stream cannot be opened.

| Exception | When |
|---|---|
| `PgChangeFeedBadRequestException` | 400 — invalid request or a violated rule. |
| `PgChangeFeedUnauthorizedException` | 401 — token missing or unknown. |
| `PgChangeFeedForbiddenException` | 403 — known token whose class may not call this endpoint (for example a reader token on an admin call). |
| `PgChangeFeedNotFoundException` | 404 — the table does not exist in the source database (`EnableTableAsync`, `DisableTableAsync`, `GetStatusAsync`). |
| `PgChangeFeedServerErrorException` | 500 — unexpected error inside the server. |
| `PgChangeFeedUnexpectedStatusException` | any other non-success status. |
| `PgChangeFeedMalformedResponseException` | a success response or SSE event whose content cannot be read. |

```csharp
using PgChangeFeed.Client.Http;

try
{
    await client.ListTablesAsync("my-source", "my_publication");
}
catch (PgChangeFeedUnauthorizedException)
{
    Console.WriteLine("token missing or unknown");
}
catch (PgChangeFeedForbiddenException)
{
    Console.WriteLine("token is not allowed to call this endpoint");
}
catch (PgChangeFeedException error)
{
    Console.WriteLine($"{error.StatusCode} {error.Message}");
}
```

The gRPC stream reports a missing or unknown token as a `Grpc.Core.RpcException` with status `Unauthenticated`. The NATS stream throws the connection exception of the NATS client library when the server rejects the token, and `PgChangeFeedNatsMalformedMessageException` when a message cannot be read.

## Reading versus streaming

- **Reading over HTTP** (`ReadChangesAsync`) is asking: you name a range, the server answers from the changes it has stored. You can read the same range again, and with a registered consumer you can carry on after a restart exactly where you stopped. Changes stay readable until the retention removes them.
- **The live streams** (gRPC, SSE, NATS) are pushing: you get every change committed after you connected, in commit order, with the full row content. There is no delivery guarantee and no replay. A change committed while you were disconnected, or while you read too slowly, does not arrive on the stream. Use the stream to react quickly and `ReadChangesAsync` to catch up on what it missed.
- The gRPC and SSE streams cannot be filtered by table; on the NATS stream the subject chooses the source or the table.

## More

- [Project README](https://github.com/pt9912/pg-change-feed/blob/main/README.md)
- [User manual](https://github.com/pt9912/pg-change-feed/blob/main/docs/user/benutzerhandbuch.md) — setting up and operating the server, tokens, delivery paths (in German)
- [Reference clients](https://github.com/pt9912/pg-change-feed/tree/main/examples/csharp) for each delivery path

## License

MIT — see [LICENSE](https://github.com/pt9912/pg-change-feed/blob/main/LICENSE).
