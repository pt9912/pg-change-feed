"""Typed errors of the PG Change Feed clients.

The HTTP error body ``{"error": "<text>"}`` becomes a typed exception instead of
a raw ``httpx`` exception, consistent across every ``PgChangeFeedHttpClient``
method -- one base class a caller can catch regardless of the concrete status
code. ``PgChangeFeedGrpcError`` and its subclasses do the same for
``PgChangeFeedAdministrationClient``'s gRPC status codes -- a separate
hierarchy, because a gRPC status is a ``grpc.StatusCode`` enum member, not an
HTTP status integer.
"""

from __future__ import annotations

import grpc


class PgChangeFeedError(Exception):
    """Base type for every typed error the PG Change Feed clients raise.

    ``status_code`` is the HTTP status of the failed response (``200`` for a
    malformed SSE event or NATS message, which carries no status of its own).
    """

    def __init__(self, status_code: int, message: str) -> None:
        super().__init__(message)
        self.status_code = status_code


class PgChangeFeedBadRequestError(PgChangeFeedError):
    """``400`` -- an invalid request body or a violated rule of the API."""


class PgChangeFeedUnauthorizedError(PgChangeFeedError):
    """``401`` -- a missing bearer token, or one that matches no configured
    token class."""


class PgChangeFeedForbiddenError(PgChangeFeedError):
    """``403`` -- a known token whose class does not reach the called endpoint,
    e.g. a ``reader`` token against an ``admin`` endpoint."""


class PgChangeFeedNotFoundError(PgChangeFeedError):
    """``404`` -- the addressed table does not exist in the source database
    (only ``enable_table``, ``disable_table`` and ``get_status``)."""


class PgChangeFeedServerError(PgChangeFeedError):
    """``500`` -- an unexpected internal error of the server."""


class PgChangeFeedUnexpectedStatusError(PgChangeFeedError):
    """Any other non-success status code than 400, 401, 403, 404 and 500."""


class PgChangeFeedMalformedResponseError(PgChangeFeedError):
    """A response, SSE event or NATS message whose content cannot be read:
    invalid JSON, or valid JSON missing an expected field. It is distinct from
    the status-code errors above so a caller can still catch
    ``PgChangeFeedError`` uniformly across every method."""


class PgChangeFeedGrpcError(Exception):
    """Base type for every typed error ``PgChangeFeedAdministrationClient``
    raises for a non-``OK`` gRPC status.

    ``code`` is the ``grpc.StatusCode`` the server returned; the original
    ``grpc.RpcError`` is always ``__cause__``.
    """

    def __init__(self, code: grpc.StatusCode, message: str) -> None:
        super().__init__(message)
        self.code = code


class PgChangeFeedGrpcInvalidArgumentError(PgChangeFeedGrpcError):
    """``INVALID_ARGUMENT`` -- a violated rule of the request (e.g. an empty
    identifier, an inverted range, a non-positive limit)."""


class PgChangeFeedGrpcUnauthenticatedError(PgChangeFeedGrpcError):
    """``UNAUTHENTICATED`` -- a missing bearer token, or one that matches no
    configured token class."""


class PgChangeFeedGrpcPermissionDeniedError(PgChangeFeedGrpcError):
    """``PERMISSION_DENIED`` -- a known token whose class does not reach the
    called RPC (e.g. a reader token against an admin RPC)."""


class PgChangeFeedGrpcNotFoundError(PgChangeFeedGrpcError):
    """``NOT_FOUND`` -- the addressed table does not exist in the source
    database (``enable_table``, ``disable_table``, ``get_table_status``)."""


class PgChangeFeedGrpcInternalError(PgChangeFeedGrpcError):
    """``INTERNAL`` -- an unexpected internal error of the server."""


class PgChangeFeedGrpcUnexpectedStatusError(PgChangeFeedGrpcError):
    """Any other non-``OK`` gRPC status than ``INVALID_ARGUMENT``,
    ``UNAUTHENTICATED``, ``PERMISSION_DENIED``, ``NOT_FOUND`` and
    ``INTERNAL``. It is distinct from the five status-specific errors above so
    a caller can still catch ``PgChangeFeedGrpcError`` uniformly across every
    method."""
