"""Typed errors of the PG Change Feed clients.

The HTTP error body ``{"error": "<text>"}`` becomes a typed exception instead of
a raw ``httpx`` exception, consistent across every ``PgChangeFeedHttpClient``
method -- one base class a caller can catch regardless of the concrete status
code.
"""

from __future__ import annotations


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
