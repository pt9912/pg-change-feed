"""PG Change Feed Python client library.

This release exposes the shared connection configuration
(`ClientOptions`), the HTTP API client surface (`PgChangeFeedHttpClient`):
the nine SPEC-018 capabilities plus reading persisted changes
(`ReadChanges`, SPEC-022), and the gRPC live-change-stream surface
(`PgChangeFeedGrpcClient`): the `StreamChanges` server-streaming RPC
(SPEC-020). SSE and NATS-Vollinhalt delivery remain out of scope for this
package (ADR-0110 carries the remaining surfaces of the four-way matrix) --
the follow-up slices of the same wave add a separate client surface for
them.
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
from pgchangefeed.options import ClientOptions

__all__ = [
    "ClientOptions",
    "PgChangeFeedHttpClient",
    "PgChangeFeedGrpcClient",
    "PgChangeFeedError",
    "PgChangeFeedBadRequestError",
    "PgChangeFeedUnauthorizedError",
    "PgChangeFeedForbiddenError",
    "PgChangeFeedNotFoundError",
    "PgChangeFeedServerError",
    "PgChangeFeedUnexpectedStatusError",
    "PgChangeFeedMalformedResponseError",
]
