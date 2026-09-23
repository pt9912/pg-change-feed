"""Typed error hierarchy for PG Change Feed HTTP/JSON API responses.

Mirrors the uniform SPEC-018 error body (``{"error": "<text>"}``) as typed
exceptions instead of a raw ``httpx`` exception, consistent across every
``PgChangeFeedHttpClient`` method -- one base class a caller can catch
regardless of the concrete status code.
"""

from __future__ import annotations


class PgChangeFeedError(Exception):
    """Base type for every typed error ``PgChangeFeedHttpClient`` raises."""

    def __init__(self, status_code: int, message: str) -> None:
        super().__init__(message)
        self.status_code = status_code


class PgChangeFeedBadRequestError(PgChangeFeedError):
    """``400`` -- an invalid request body or a violated domain invariant
    (SPEC-018/SPEC-022)."""


class PgChangeFeedUnauthorizedError(PgChangeFeedError):
    """``401`` -- a missing bearer token, or one that matches no configured
    token class (SPEC-018)."""


class PgChangeFeedForbiddenError(PgChangeFeedError):
    """``403`` -- a known token whose rights class does not reach the
    called endpoint, e.g. a ``reader`` token against an ``admin`` endpoint
    (SPEC-018)."""


class PgChangeFeedNotFoundError(PgChangeFeedError):
    """``404`` -- the addressed table is physically missing at the source
    (only ``EnableTable``/``DisableTable``/``GetStatus``, SPEC-018)."""


class PgChangeFeedServerError(PgChangeFeedError):
    """``500`` -- an unexpected internal server error (SPEC-018/SPEC-022)."""


class PgChangeFeedUnexpectedStatusError(PgChangeFeedError):
    """Any non-success status code outside the five SPEC-018/SPEC-022
    document (400/401/403/404/500) -- a defensive fallback that is itself
    not part of the documented wire contract."""


class PgChangeFeedMalformedResponseError(PgChangeFeedError):
    """A success status code (``2xx``) whose body does not match the
    expected SPEC-018/SPEC-022 response shape, or a SPEC-021 SSE event whose
    data payload does not match the documented stream schema -- invalid
    JSON, or valid JSON missing an expected field. Outside every shape the
    HTTP/SSE contracts document; kept typed and distinct from the
    status-code errors above so a caller can still catch
    ``PgChangeFeedError`` uniformly across every method."""
