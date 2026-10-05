"""PG Change Feed Python client library.

The package offers five clients and one shared configuration class:

- ``ClientOptions`` -- server address, token and optional TLS trust anchor,
  shared by all clients.
- ``create_http_client`` and ``create_grpc_channel`` -- the HTTP client and
  the gRPC channel that connect over TLS with the trust anchor of the options.
- ``PgChangeFeedHttpClient`` -- the HTTP API: manage consumers, enable and
  disable captured tables, run the retention and read stored changes.
- ``PgChangeFeedGrpcClient`` -- a live gRPC stream of changes.
- ``PgChangeFeedSseClient`` -- a live stream of changes over Server-Sent
  Events.
- ``PgChangeFeedNatsStreamClient`` -- a live stream of changes over NATS; the
  token is checked once when the connection is opened.
- ``PgChangeFeedAdministrationClient`` -- the same management/read/diagnose
  capabilities as ``PgChangeFeedHttpClient``, over gRPC instead of HTTP.

The three live streams carry the same change messages. They deliver what is
committed while the client is connected; they do not repeat changes that were
missed (use ``PgChangeFeedHttpClient.read_changes`` for catching up).
"""

from pgchangefeed.administration_client import PgChangeFeedAdministrationClient
from pgchangefeed.exceptions import (
    PgChangeFeedBadRequestError,
    PgChangeFeedError,
    PgChangeFeedForbiddenError,
    PgChangeFeedGrpcError,
    PgChangeFeedGrpcInternalError,
    PgChangeFeedGrpcInvalidArgumentError,
    PgChangeFeedGrpcNotFoundError,
    PgChangeFeedGrpcPermissionDeniedError,
    PgChangeFeedGrpcUnauthenticatedError,
    PgChangeFeedGrpcUnexpectedStatusError,
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
from pgchangefeed.tls import create_grpc_channel, create_http_client

__all__ = [
    "ClientOptions",
    "create_http_client",
    "create_grpc_channel",
    "PgChangeFeedHttpClient",
    "PgChangeFeedGrpcClient",
    "PgChangeFeedSseClient",
    "PgChangeFeedNatsStreamClient",
    "PgChangeFeedAdministrationClient",
    "StreamChange",
    "PgChangeFeedError",
    "PgChangeFeedBadRequestError",
    "PgChangeFeedUnauthorizedError",
    "PgChangeFeedForbiddenError",
    "PgChangeFeedNotFoundError",
    "PgChangeFeedServerError",
    "PgChangeFeedUnexpectedStatusError",
    "PgChangeFeedMalformedResponseError",
    "PgChangeFeedGrpcError",
    "PgChangeFeedGrpcInvalidArgumentError",
    "PgChangeFeedGrpcUnauthenticatedError",
    "PgChangeFeedGrpcPermissionDeniedError",
    "PgChangeFeedGrpcNotFoundError",
    "PgChangeFeedGrpcInternalError",
    "PgChangeFeedGrpcUnexpectedStatusError",
]
