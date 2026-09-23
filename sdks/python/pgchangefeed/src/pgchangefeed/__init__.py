"""PG Change Feed Python client library.

This release exposes the shared connection configuration
(`ClientOptions`), the HTTP API client surface (`PgChangeFeedHttpClient`):
the nine SPEC-018 capabilities plus reading persisted changes
(`ReadChanges`, SPEC-022), the gRPC live-change-stream surface
(`PgChangeFeedGrpcClient`): the `StreamChanges` server-streaming RPC
(SPEC-020), and the SSE live-change-stream surface (`PgChangeFeedSseClient`):
the same fire-and-forget, no-replay stream over Server-Sent-Events
(SPEC-021). NATS-Vollinhalt delivery remains out of scope for this release
(ADR-0110 carries the remaining surface of the four-way matrix) -- the
follow-up slice of the same wave adds a separate client surface for it.
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
from pgchangefeed.options import ClientOptions
from pgchangefeed.sse_client import PgChangeFeedSseClient

__all__ = [
    "ClientOptions",
    "PgChangeFeedHttpClient",
    "PgChangeFeedGrpcClient",
    "PgChangeFeedSseClient",
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
