"""PG Change Feed Python client library.

This release exposes the shared connection configuration
(`ClientOptions`), the HTTP API client surface (`PgChangeFeedHttpClient`):
the nine SPEC-018 capabilities plus reading persisted changes
(`ReadChanges`, SPEC-022), the gRPC live-change-stream surface
(`PgChangeFeedGrpcClient`): the `StreamChanges` server-streaming RPC
(SPEC-020), the SSE live-change-stream surface (`PgChangeFeedSseClient`):
the same fire-and-forget, no-replay stream over Server-Sent-Events
(SPEC-021), and the NATS full-content stream surface
(`PgChangeFeedNatsStreamClient`): the same message schema and
fire-and-forget, no-replay boundary, delivered over a NATS subject namespace
instead of HTTP, with connection-level (not per-call) token authentication
(SPEC-024, Welle-Plan §6: v2 desselben Packages — ADR-0110 Festlegung 3
delegiert die Struktur-Entscheidung an den umsetzenden Zug).
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
