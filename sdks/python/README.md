# PG Change Feed — Python client

Python client library for [PG Change Feed](https://github.com/pt9912/pg-change-feed), a server that records every INSERT, UPDATE and DELETE of selected PostgreSQL tables and makes these changes available over HTTP, gRPC, Server-Sent Events (SSE) and NATS.

With this package (`pgchangefeed`) a Python application can

- read the recorded changes and keep track of how far it has processed them,
- manage which tables the server captures, and
- receive changes live, as they happen,

without implementing any of the wire protocols itself.

Version 0.x — the API can still change between releases.

## Installation

```
pip install pgchangefeed
```

Requires Python 3.14 or newer.

## Quick start

A client needs the address of the PG Change Feed server and a token. The server knows two token classes: a *reader* token for read-only calls and an *admin* token for calls that change something (registering consumers, acknowledging positions, enabling tables, running the retention). The admin token also covers all reader calls. Address, source id and tokens come from whoever operates the server.

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

gRPC (`address` is `host:port`; the messages are the generated `Change` protobuf messages, row images are JSON bytes):

```python
import grpc

from pgchangefeed import ClientOptions, PgChangeFeedGrpcClient

options = ClientOptions(address="feed.example.com:9090", api_token="<reader token>")
with grpc.insecure_channel(options.address) as channel:
    client = PgChangeFeedGrpcClient(channel, options)
    for change in client.stream_changes():
        print(change.operation, change.schema, change.table, change.new_image.decode())
```

Server-Sent Events (`address` is the HTTP base URL; the read timeout must be off for a long-lived stream):

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

NATS (`address` is the NATS URL, the token is the NATS stream token checked when the connection is opened; the stream covers all tables of one source):

```python
from pgchangefeed import ClientOptions, PgChangeFeedNatsStreamClient

client = PgChangeFeedNatsStreamClient(
    ClientOptions(address="nats://feed.example.com:4222", api_token="<NATS stream token>"),
    "my-source",
)
for change in client.stream_changes():
    print(change.operation, change.schema, change.table, change.new_image)
```

## API overview

`PgChangeFeedHttpClient(client, options)` wraps the HTTP API. The `httpx.Client` you pass in stays yours; the library never closes it.

| Method | What it does | Token |
|---|---|---|
| `register_consumer(request)` | Registers a consumer. Registering an existing consumer changes nothing (`already_registered` is true). | admin |
| `acknowledge_consumer(request)` | Stores the consumer's position. Repeating the same position has no effect; a position before the stored one is rejected. | admin |
| `get_consumer_position(consumer_id)` | Reads the stored position (`offset`, and `acknowledged`, which is false for a consumer that never acknowledged). | reader |
| `remove_consumer(consumer_id)` | Removes a consumer. | admin |
| `enable_table(request)` | Starts capturing a table. | admin |
| `disable_table(request)` | Stops capturing a table. `retained` reports that changes already stored for it remain. | admin |
| `get_status(source, schema, table, publication)` | Tells whether a table is captured (`enabled`) or no longer captured with stored changes remaining (`retained`). | reader |
| `list_tables(source, publication)` | Lists the captured tables and the tables whose stored changes remain. | reader |
| `run_retention(request)` | Deletes stored changes older than `min_age_nanos` that every consumer with a stored position has passed; returns the number deleted. | admin |
| `read_changes(source, schema, table, from_, to, limit)` | Reads stored changes of a source, optionally for one schema and table and for the range `[from_, to)` of commit positions. Only `source` is required. | reader |

Request and response classes live in `pgchangefeed.models`.

The live streams each have one method:

| Class | Method | What it does |
|---|---|---|
| `PgChangeFeedGrpcClient(channel, options)` | `stream_changes(timeout=None)` | Opens the gRPC stream and yields generated `Change` messages. `timeout` is the deadline of the whole call in seconds. |
| `PgChangeFeedSseClient(client, options)` | `stream_changes()` | Opens `GET /changes/stream` and yields `StreamChange` objects. |
| `PgChangeFeedNatsStreamClient(options, source_id)` | `stream_changes(timeout=None)` | Subscribes to all tables of one source and yields `StreamChange` objects. `timeout` bounds the total consumption in seconds. |

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

Every failing HTTP call raises a subclass of `PgChangeFeedError`, which carries the HTTP `status_code`. The same classes are raised when the SSE stream cannot be opened.

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
    print(error.status_code, error)
```

The gRPC stream reports a missing or unknown token as `grpc.RpcError` with status `UNAUTHENTICATED`. The NATS stream raises the connection error of the NATS client library when the server rejects the token.

## Reading versus streaming

- **Reading over HTTP** (`read_changes`) is asking: you name a range, the server answers from the changes it has stored. You can read the same range again, and with a registered consumer you can carry on after a restart exactly where you stopped. Changes stay readable until the retention removes them.
- **The live streams** (gRPC, SSE, NATS) are pushing: you get every change committed after you connected, in commit order, with the full row content. There is no delivery guarantee and no replay. A change committed while you were disconnected, or while you read too slowly, does not arrive on the stream. Use the stream to react quickly and `read_changes` to catch up on what it missed.
- The gRPC and SSE streams cannot be filtered by table; the NATS stream of this package covers all tables of one source.

## More

- [Project README](https://github.com/pt9912/pg-change-feed/blob/main/README.md)
- [User manual](https://github.com/pt9912/pg-change-feed/blob/main/docs/user/benutzerhandbuch.md) — setting up and operating the server, tokens, delivery paths (in German)
- [Reference clients](https://github.com/pt9912/pg-change-feed/tree/main/examples) for each delivery path in Go, C# and Kotlin

## License

MIT — see [LICENSE](https://github.com/pt9912/pg-change-feed/blob/main/LICENSE).
