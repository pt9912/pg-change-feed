# PG Change Feed — Python client

Python client library for [PG Change Feed](https://github.com/pt9912/pg-change-feed), a server that records every INSERT, UPDATE and DELETE of selected PostgreSQL tables and makes these changes available over HTTP, gRPC, Server-Sent Events (SSE) and NATS.

With this package (`pgchangefeed`) a Python application can

- read the recorded changes and keep track of how far it has processed them,
- manage which tables the server captures, and
- receive changes live, as they happen,

either over HTTP or over gRPC, without implementing any of the wire protocols itself.

Version 0.x — the API can still change between releases.

## Installation

```
pip install pgchangefeed
```

Requires Python 3.14 or newer.

## Quick start

A client needs the address of the PG Change Feed server and a token. The server knows two token classes: a *reader* token for read-only calls and an *admin* token for calls that change something (registering consumers, acknowledging positions, enabling tables, running the retention). The admin token also covers all reader calls. Address, source id and tokens come from whoever operates the server. The server does not serve TLS itself, so the examples use unencrypted addresses.

### Read changes and remember your position

A *consumer* is a named reader whose progress the server remembers. You register it once, read changes, and acknowledge the position of the last change you have processed. A position is the `commit_position` of a change; `read_changes` reads from `from_` (inclusive) up to `to` (exclusive). A consumer that has never acknowledged reports offset 0.

```python
import httpx

from pgchangefeed import ClientOptions, PgChangeFeedHttpClient
from pgchangefeed.models import AcknowledgeConsumerRequest, RegisterConsumerRequest

http = httpx.Client(timeout=30.0)
client = PgChangeFeedHttpClient(
    http, ClientOptions(address="http://feed.example.com:8090", api_token="<admin token>")
)

client.register_consumer(RegisterConsumerRequest(consumer_id="billing", name="Billing service"))

position = client.get_consumer_position("billing")
result = client.read_changes("my-source", from_=position.offset + 1)
for change in result.changes:
    print(change.commit_position, change.operation, change.schema, change.table, change.new_image)

if result.changes:
    client.acknowledge_consumer(
        AcknowledgeConsumerRequest(
            consumer_id="billing",
            source_id="my-source",
            offset=result.changes[-1].commit_position,
        )
    )
```

`limit` cuts rows, not positions: if one commit position carries more changes than `limit`, continuing from that position plus one skips the rest of it. Leave `limit` out (or set it generously) where a single position can carry many changes.

### Receive changes live

The three live streams deliver every change committed after you connect. Each surface has its own client; all take the same `ClientOptions`.

gRPC (`address` is `host:port`; the messages are the generated `Change` protobuf messages, row images are JSON bytes, empty when there is none):

```python
import grpc

from pgchangefeed import ClientOptions, PgChangeFeedGrpcClient

options = ClientOptions(address="feed.example.com:9090", api_token="<reader token>")
with grpc.insecure_channel(options.address) as channel:
    client = PgChangeFeedGrpcClient(channel, options)
    for change in client.stream_changes(schema="public", table="orders"):
        print(change.operation, change.schema, change.table, change.new_image.decode())
```

`schema`/`table`/`target` are each optional and independent: a set `schema` without `table` delivers every table of that schema, a set `table` without `schema` delivers every table of that name regardless of schema, both set delivers exactly one table, and leaving all out (the original, still valid call) delivers every change of every captured table. `target` selects the delivery target a change is routed to: a set value delivers only the changes routed to that target, combined with `schema`/`table` as a conjunction, and a target no change carries delivers nothing and raises no error: `client.stream_changes(target="eu")`. A server release that predates the parameter is expected to ignore it on the gRPC stream, which then stays unfiltered, and to answer `400` on the HTTP read path and the SSE stream; this follows from how the server reads its parameters and has not been run against such a release.

Server-Sent Events (`address` is the HTTP base URL; the read timeout must be off for a long-lived stream). The SSE stream takes the same three optional, independent parameters, `target`, `schema` and `table`: `stream_changes(target="eu")` delivers only the changes routed to that target, `stream_changes(schema="public", table="orders")` only the changes of that table, and leaving them out delivers every change. Combined, they act as a conjunction:

```python
import httpx

from pgchangefeed import ClientOptions, PgChangeFeedSseClient

http = httpx.Client(timeout=httpx.Timeout(10.0, read=None))
client = PgChangeFeedSseClient(
    http, ClientOptions(address="http://feed.example.com:8090", api_token="<reader token>")
)
for change in client.stream_changes():
    print(change.operation, change.schema, change.table, change.new_image)
```

NATS (`address` is the NATS URL, the token is the NATS stream token checked when the connection is opened; the stream covers all tables of one source). `stream_changes(target="eu")` subscribes to the subject of one delivery target of the source (`cdc.route.<source>.<target>`) instead: a change routed to a target arrives on its table subject and, with the same payload, on the target subject, a change without a target on the table subject only. An empty `target` means no target (all tables of the source); a non-empty target that is whitespace only or contains `.`, `*`, `>` or whitespace raises `ValueError` at the first `next()` on the returned iterator, before a connection is opened:

```python
from pgchangefeed import ClientOptions, PgChangeFeedNatsStreamClient

client = PgChangeFeedNatsStreamClient(
    ClientOptions(address="nats://feed.example.com:4222", api_token="<NATS stream token>"),
    "my-source",
)
for change in client.stream_changes():
    print(change.operation, change.schema, change.table, change.new_image)
```

### Manage tables and consumers over gRPC

`PgChangeFeedAdministrationClient` wraps the same eleven management/read/diagnose capabilities as `PgChangeFeedHttpClient`, over gRPC instead of HTTP -- one method per RPC. Requests and responses are the generated `pgchangefeed.grpc_gen.administration_pb2` protobuf messages, used directly:

```python
import grpc

from pgchangefeed import ClientOptions, PgChangeFeedAdministrationClient
from pgchangefeed.grpc_gen import administration_pb2

options = ClientOptions(address="feed.example.com:9090", api_token="<admin token>")
with grpc.insecure_channel(options.address) as channel:
    admin = PgChangeFeedAdministrationClient(channel, options)

    admin.enable_table(administration_pb2.EnableTableRequest(
        source="my-source", schema="public", table="orders",
        table_id="orders", schema_version_id="v1", version=1, publication="my_publication",
    ))

    tables = admin.list_tables(administration_pb2.ListTablesRequest(source="my-source", publication="my_publication"))
    for table in tables.tables:
        print(table.schema, table.table)
```

## Administration client over gRPC

`PgChangeFeedAdministrationClient(channel, options)` wraps the eleven RPCs of the `Administration` gRPC service -- the same capabilities as the HTTP table below, over gRPC. Requests and responses are the generated `administration_pb2` protobuf messages, used directly (no separate model type):

| Method | What it does | Token |
|---|---|---|
| `register_consumer(request)` | Registers a consumer. Registering an existing consumer changes nothing (`already_registered` is true). | admin |
| `acknowledge_consumer(request)` | Stores the consumer's position. Repeating the same position has no effect; an earlier position is rejected. | admin |
| `get_consumer_position(request)` | Reads the stored position (`offset`, and `acknowledged`, which is false for a consumer that never acknowledged). | reader |
| `remove_consumer(request)` | Removes a consumer (`removed` is false for one that was never registered). | admin |
| `enable_table(request)` | Starts capturing a table. A table that is already captured changes nothing (`already_enabled` is true); a table that does not exist raises `PgChangeFeedGrpcNotFoundError`. | admin |
| `disable_table(request)` | Stops capturing a table. `retained` reports that changes already stored for it remain. | admin |
| `get_table_status(request)` | Tells whether a table is captured (`enabled`) or no longer captured with stored changes remaining (`retained`). | reader |
| `list_tables(request)` | Lists the captured tables (`tables`) and the tables whose stored changes remain (`retained`), each a `SourceTable`. | reader |
| `run_retention(request)` | Deletes stored changes older than `min_age_nanos` that every consumer with a stored position has passed; returns the number deleted (`deleted`). A `min_age_nanos` of `0` means no minimum age and is valid. | admin |
| `read_changes(request)` | Reads stored changes of a source, optionally for one schema and table, for the range `[from, to)` of commit positions and for one delivery target (`target`, empty is no filter); `changes` is a list of `ChangeRecord`. An unset `from`/`to`/`limit` carries `0` (not set); `from` is a Python keyword, reach the field with `getattr(request, "from")`. | reader |
| `diagnose(request)` | Reads the operational diagnose report (heartbeat, capture lag, per-consumer lag, retention blocker, storage, backfill status). A `known`/`present`/`*_known` field of `false` carries the respective absence case. | reader |

A call that the server answers with a non-`OK` gRPC status raises a subclass of `PgChangeFeedGrpcError` instead of a raw `grpc.RpcError`; the original `grpc.RpcError` is always `__cause__` (see [Error handling](#error-handling) for the HTTP-side exceptions):

| Exception | gRPC status |
|---|---|
| `PgChangeFeedGrpcInvalidArgumentError` | `INVALID_ARGUMENT` — invalid request or a violated rule. |
| `PgChangeFeedGrpcUnauthenticatedError` | `UNAUTHENTICATED` — token missing or unknown. |
| `PgChangeFeedGrpcPermissionDeniedError` | `PERMISSION_DENIED` — known token whose class may not call this RPC (for example a reader token on an admin RPC). |
| `PgChangeFeedGrpcNotFoundError` | `NOT_FOUND` — the table does not exist in the source database (`enable_table`, `disable_table`, `get_table_status`). |
| `PgChangeFeedGrpcInternalError` | `INTERNAL` — unexpected error inside the server. |
| `PgChangeFeedGrpcUnexpectedStatusError` | any other non-`OK` status. |

Every one of these errors carries `message_code`, the message code of the server read from the status detail (`None` when the server sent none, as for `UNAUTHENTICATED` and `PERMISSION_DENIED`); for example `PgChangeFeedGrpcNotFoundError.message_code` is `PCF-E8025` for `enable_table` on a table that does not exist.

## API overview

`PgChangeFeedHttpClient(client, options)` wraps the HTTP API. The `httpx.Client` you pass in stays yours; the library never closes it.

| Method | What it does | Token |
|---|---|---|
| `register_consumer(request)` | Registers a consumer. Registering an existing consumer changes nothing (`already_registered` is true). | admin |
| `acknowledge_consumer(request)` | Stores the consumer's position. Repeating the same position has no effect; a position before the stored one is rejected. | admin |
| `get_consumer_position(consumer_id)` | Reads the stored position (`offset`, and `acknowledged`, which is false for a consumer that never acknowledged). | reader |
| `remove_consumer(consumer_id)` | Removes a consumer. | admin |
| `enable_table(request)` | Starts capturing a table (see below). A table that is already captured changes nothing (`already_enabled` is true); a table that does not exist raises `PgChangeFeedNotFoundError`. | admin |
| `disable_table(request)` | Stops capturing a table. `retained` reports that changes already stored for it remain. | admin |
| `get_status(source, schema, table, publication)` | Tells whether a table is captured (`enabled`) or no longer captured with stored changes remaining (`retained`). | reader |
| `list_tables(source, publication)` | Lists the captured tables and the tables whose stored changes remain. | reader |
| `run_retention(request)` | Deletes stored changes older than `min_age_nanos` that every consumer with a stored position has passed; returns the number deleted. | admin |
| `read_changes(source, schema, table, from_, to, limit, target)` | Reads stored changes of a source, optionally for one schema and table, for the range `[from_, to)` of commit positions and for one delivery target (`target`, combined with `schema`/`table` as a conjunction; a target no change carries returns an empty list). Only `source` is required. | reader |

Request and response classes live in `pgchangefeed.models`.

`EnableTableRequest` names the table and the ids it is captured under: `source`, `schema` and `table` identify the table, `publication` is the PostgreSQL publication of the source, `table_id` is the id you give the table (every change of the table carries it as `source_table_id`), and `schema_version_id` and `version` (a number, 1 or higher) identify the table's schema version (changes carry the id as `schema_version`).

The live streams each have one method:

| Class | Method | What it does |
|---|---|---|
| `PgChangeFeedGrpcClient(channel, options)` | `stream_changes(timeout=None, schema=None, table=None, target=None)` | Opens the gRPC stream and yields generated `Change` messages. `timeout` is the deadline of the whole call in seconds; `schema`/`table`/`target` filter the stream, each optional and independent. |
| `PgChangeFeedSseClient(client, options)` | `stream_changes(target=None, schema=None, table=None)` | Opens `GET /changes/stream` and yields `StreamChange` objects. `target`/`schema`/`table` are optional and independent, all left out delivers every change. |
| `PgChangeFeedNatsStreamClient(options, source_id)` | `stream_changes(timeout=None, target=None)` | Subscribes to all tables of one source, or to one delivery target of it when `target` is set, and yields `StreamChange` objects. `timeout` bounds the total consumption in seconds. |

## The change object

`read_changes` returns `Change` objects (`pgchangefeed.models.Change`):

| Field | Meaning |
|---|---|
| `commit_position` | Position of the committed source transaction that carried the change; the value you read ranges by and acknowledge. |
| `change_id` | Unique id of the change. |
| `transaction_id` | Id of the source transaction. |
| `source_table_id` | Id of the captured table. |
| `schema`, `table` | Schema and name of the table. |
| `sequence` | Order of the row change within its transaction. |
| `operation` | `INSERT`, `UPDATE` or `DELETE`. |
| `old_image` | Row values before the change as JSON; `None` for an INSERT. For UPDATE and DELETE it holds what PostgreSQL provides for the table's replica identity. |
| `new_image` | Row values after the change as JSON; `None` for a DELETE. |
| `schema_version` | Version of the table schema the change was captured with. |
| `committed_at` | Commit time of the source transaction (RFC 3339, UTC). |
| `origin` | `wal` for a change captured live from the database, `backfill` for a change that was taken from the existing table contents. A response without the field, or with JSON `null`, reads as `wal`; any other value is passed through unchanged. |

The live streams deliver `StreamChange` objects (gRPC: the generated `Change` message) with ten fields: `change_id`, `transaction_id`, `source_table_id`, `sequence`, `operation`, `old_image`, `new_image`, `schema_version`, `schema`, `table`. They carry no `commit_position`, `committed_at` or `origin`.

## Error handling

An HTTP call that the server answers with an error status raises a subclass of `PgChangeFeedError`, which carries the HTTP `status_code`. The same classes are raised when the SSE stream cannot be opened. Connection failures and timeouts are not converted: they raise the `httpx` exception (for example `httpx.ConnectError` or `httpx.TimeoutException`).

| Exception | When |
|---|---|
| `PgChangeFeedBadRequestError` | 400 — invalid request or a violated rule. |
| `PgChangeFeedUnauthorizedError` | 401 — token missing or unknown. |
| `PgChangeFeedForbiddenError` | 403 — known token whose class may not call this endpoint (for example a reader token on an admin call). |
| `PgChangeFeedNotFoundError` | 404 — the table does not exist in the source database (`enable_table`, `disable_table`, `get_status`). |
| `PgChangeFeedServerError` | 500 — unexpected error inside the server. |
| `PgChangeFeedUnexpectedStatusError` | any other non-success status. |
| `PgChangeFeedMalformedResponseError` | a success response, SSE event or NATS message whose content cannot be read. |

```python
from pgchangefeed import PgChangeFeedError, PgChangeFeedForbiddenError, PgChangeFeedUnauthorizedError

try:
    client.list_tables("my-source", "my_publication")
except PgChangeFeedUnauthorizedError:
    print("token missing or unknown")
except PgChangeFeedForbiddenError:
    print("token is not allowed to call this endpoint")
except PgChangeFeedError as error:
    print(error.status_code, error.message_code, error)
```

`message_code` is the message code the server attaches to an error (`PCF-<E|W|I><4 digits>`, for example `PCF-E8025` for a table that does not exist in the source database), read from the `code` field of the error body and passed through unchanged. Match on it instead of parsing the text. It is `None` when the server sent no code: authentication errors (`401`, `403`) carry none, a server release that predates message codes sends none, and so does a body that is not JSON or whose `code` is not a non-empty string. The error text is the `error` field of the body when that is a JSON string and the raw body otherwise; it does not depend on `code`, and `code` does not depend on the type of `error`.

The gRPC stream reports a missing or unknown token as `grpc.RpcError` with status `UNAUTHENTICATED`. The gRPC stream clients raise the plain `grpc.RpcError` for every failure; a message code sent with such an error stays readable through the status detail API of the gRPC client itself (a `google.rpc.ErrorInfo` entry of domain `pg-change-feed`, its `reason` is the code). The NATS stream raises the connection error of the NATS client library when the server rejects the token.

`PgChangeFeedAdministrationClient` raises a subclass of `PgChangeFeedGrpcError` for a non-`OK` gRPC status instead of a raw `grpc.RpcError` -- see [Administration client over gRPC](#administration-client-over-grpc) for the class per status code; the original `grpc.RpcError` is always `__cause__`. `PgChangeFeedGrpcError` carries the message code the same way as `message_code` (`None` when absent); it is read from the `reason` of the `google.rpc.ErrorInfo` status detail of domain `pg-change-feed`, whatever the status, and `code` stays the `grpc.StatusCode`. The `error_code` field of the `diagnose` response (the message code of the process error state, empty in normal operation) arrives unchanged from the generated message.

## Reading versus streaming

- **Reading over HTTP** (`read_changes`) is asking: you name a range, the server answers from the changes it has stored. You can read the same range again, and with a registered consumer you can carry on after a restart exactly where you stopped. Changes stay readable until the retention removes them.
- **The live streams** (gRPC, SSE, NATS) are pushing: you get every change committed after you connected, in commit order, with the full row content. There is no delivery guarantee and no replay. A change committed while you were disconnected, or while you read too slowly, does not arrive on the stream. Use the stream to react quickly and `read_changes` to catch up on what it missed.
- The gRPC stream and the SSE stream can be filtered by schema, table and delivery target; the NATS stream can be filtered by delivery target only (`stream_changes(target=...)`); without a target the NATS stream covers all tables of one source.

## More

- [Project README](https://github.com/pt9912/pg-change-feed/blob/main/README.md)
- [User manual](https://github.com/pt9912/pg-change-feed/blob/main/docs/user/benutzerhandbuch.md) — setting up and operating the server, tokens, delivery paths (in German)
- [Reference clients](https://github.com/pt9912/pg-change-feed/tree/main/examples) for each delivery path in Go, C# and Kotlin

## License

MIT — see [LICENSE](https://github.com/pt9912/pg-change-feed/blob/main/LICENSE).
