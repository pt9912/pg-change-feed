"""PG Change Feed Python client library.

This release exposes the shared connection configuration
(`ClientOptions`) and the HTTP API client surface (`PgChangeFeedHttpClient`):
the nine SPEC-018 capabilities plus reading persisted changes
(`ReadChanges`, SPEC-022). gRPC, SSE and NATS-Vollinhalt delivery remain
out of scope for this package (ADR-0107 Festlegung 1) -- a follow-up
release would add a separate client surface for them.
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
from pgchangefeed.http_client import PgChangeFeedHttpClient
from pgchangefeed.options import ClientOptions

__all__ = [
    "ClientOptions",
    "PgChangeFeedHttpClient",
    "PgChangeFeedError",
    "PgChangeFeedBadRequestError",
    "PgChangeFeedUnauthorizedError",
    "PgChangeFeedForbiddenError",
    "PgChangeFeedNotFoundError",
    "PgChangeFeedServerError",
    "PgChangeFeedUnexpectedStatusError",
    "PgChangeFeedMalformedResponseError",
]
