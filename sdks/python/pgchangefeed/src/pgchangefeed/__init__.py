"""PG Change Feed Python client library.

The package offers four clients and one shared configuration class:

- ``ClientOptions`` -- server address and token, shared by all clients.
- ``PgChangeFeedHttpClient`` -- the HTTP API: manage consumers, enable and
  disable captured tables, run the retention and read stored changes.
- ``PgChangeFeedGrpcClient`` -- a live gRPC stream of changes.
- ``PgChangeFeedSseClient`` -- a live stream of changes over Server-Sent
  Events.
- ``PgChangeFeedNatsStreamClient`` -- a live stream of changes over NATS; the
  token is checked once when the connection is opened.

The three live streams carry the same change messages. They deliver what is
committed while the client is connected; they do not repeat changes that were
missed (use ``PgChangeFeedHttpClient.read_changes`` for catching up).
"""

from pgchangefeed.exceptions import (
    PgChangeFeedBadRequestError,
    PgChangeFeedError,
    PgChangeFeedForbiddenError,
    PgChangeFeedMalformedResponseError,
    PgChangeFeedNotFoundError,
    PgChangeFeedServerError,
    PgChangeFeedUnauthorizedError,
    PgChangeFeedUnexpectedStatusError,
)
from pgchangefeed.grpc_client import PgChangeFeedGrpcClient
from pgchangefeed.http_client import PgChangeFeedHttpClient
from pgchangefeed.models import StreamChange
from pgchangefeed.nats_stream_client import PgChangeFeedNatsStreamClient
from pgchangefeed.options import ClientOptions
from pgchangefeed.sse_client import PgChangeFeedSseClient

__all__ = [
    "ClientOptions",
    "PgChangeFeedHttpClient",
    "PgChangeFeedGrpcClient",
    "PgChangeFeedSseClient",
    "PgChangeFeedNatsStreamClient",
    "StreamChange",
    "PgChangeFeedError",
    "PgChangeFeedBadRequestError",
    "PgChangeFeedUnauthorizedError",
    "PgChangeFeedForbiddenError",
    "PgChangeFeedNotFoundError",
    "PgChangeFeedServerError",
    "PgChangeFeedUnexpectedStatusError",
    "PgChangeFeedMalformedResponseError",
]
