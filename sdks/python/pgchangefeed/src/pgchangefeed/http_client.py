"""Client for the PG Change Feed HTTP/JSON API.

One method per capability of the API: ``register_consumer``,
``acknowledge_consumer``, ``get_consumer_position``, ``remove_consumer``,
``enable_table``, ``disable_table``, ``get_status``, ``list_tables``,
``run_retention`` and ``read_changes``. Requests and responses are typed data
classes that mirror the JSON documents of the API (``pgchangefeed.models``).
Every non-success response raises a typed ``PgChangeFeedError`` subclass
instead of a raw ``httpx`` exception, and a success (``2xx``) response whose
body cannot be read -- invalid JSON, or valid JSON missing an expected field --
raises ``PgChangeFeedMalformedResponseError``. Connection errors and timeouts
of the transport are not converted; they reach the caller as the ``httpx``
exception.

The ``httpx.Client`` is passed in, not owned: the caller controls its
lifetime, connection pooling and transport (including ``httpx.MockTransport``
for tests); this class never closes it. The bearer token and server address
come from ``ClientOptions``; there is no module-level or global state, so a
process can hold several independently configured instances at once.
"""

from __future__ import annotations

from typing import Any, Callable, TypeVar

import httpx

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
from pgchangefeed.models import (
    AcknowledgeConsumerRequest,
    AcknowledgeConsumerResponse,
    ConsumerPositionResponse,
    DisableTableRequest,
    DisableTableResponse,
    EnableTableRequest,
    EnableTableResponse,
    ListTablesResponse,
    ReadChangesResponse,
    RegisterConsumerRequest,
    RegisterConsumerResponse,
    RemoveConsumerResponse,
    RunRetentionRequest,
    RunRetentionResponse,
    TableStatusResponse,
)
from pgchangefeed.options import ClientOptions

T = TypeVar("T")

_STATUS_TO_ERROR: dict[int, type[PgChangeFeedError]] = {
    400: PgChangeFeedBadRequestError,
    401: PgChangeFeedUnauthorizedError,
    403: PgChangeFeedForbiddenError,
    404: PgChangeFeedNotFoundError,
    500: PgChangeFeedServerError,
}


class PgChangeFeedHttpClient:
    """Client for the HTTP API: consumers, captured tables, retention and reading changes.

    ``client`` is the ``httpx.Client`` used for every call; ``options`` carries
    the server address and the bearer token.
    """

    def __init__(self, client: httpx.Client, options: ClientOptions) -> None:
        self._client = client
        self._options = options

    # --- POST /consumers (admin token) ---

    def register_consumer(self, request: RegisterConsumerRequest) -> RegisterConsumerResponse:
        """Registers a consumer, a named reader whose position the server keeps.

        Registering a consumer that already exists changes nothing; the
        response reports it with ``already_registered=True``.
        """
        body = {"consumer_id": request.consumer_id, "name": request.name}
        return self._post("/consumers", body, RegisterConsumerResponse.from_json)

    # --- POST /consumers/acknowledge (admin token) ---

    def acknowledge_consumer(
        self, request: AcknowledgeConsumerRequest
    ) -> AcknowledgeConsumerResponse:
        """Stores the position (``offset``) up to which a consumer has processed a source.

        Repeating the stored position has no effect; a position before the
        stored one, or a position of another source, is rejected with
        ``PgChangeFeedBadRequestError``.
        """
        body = {
            "consumer_id": request.consumer_id,
            "source_id": request.source_id,
            "offset": request.offset,
        }
        return self._post("/consumers/acknowledge", body, AcknowledgeConsumerResponse.from_json)

    # --- GET /consumers/position (reader or admin token) ---

    def get_consumer_position(self, consumer_id: str) -> ConsumerPositionResponse:
        """Reads the stored position of a consumer.

        ``acknowledged`` is ``False`` for a consumer that has never
        acknowledged; ``offset`` is then the defined starting position.
        """
        return self._get(
            "/consumers/position",
            {"consumer_id": consumer_id},
            ConsumerPositionResponse.from_json,
        )

    # --- POST /consumers/remove (admin token) ---

    def remove_consumer(self, consumer_id: str) -> RemoveConsumerResponse:
        """Removes a consumer; ``removed`` is ``False`` for one that was never registered."""
        return self._post(
            "/consumers/remove", {"consumer_id": consumer_id}, RemoveConsumerResponse.from_json
        )

    # --- POST /tables/enable (admin token) ---

    def enable_table(self, request: EnableTableRequest) -> EnableTableResponse:
        """Starts capturing a table of a source.

        ``already_enabled`` in the response is ``True`` when the table was
        captured already. A table that does not exist in the source database
        raises ``PgChangeFeedNotFoundError``.
        """
        body = {
            "source": request.source,
            "schema": request.schema,
            "table": request.table,
            "table_id": request.table_id,
            "schema_version_id": request.schema_version_id,
            "version": request.version,
            "publication": request.publication,
        }
        return self._post("/tables/enable", body, EnableTableResponse.from_json)

    # --- POST /tables/disable (admin token) ---

    def disable_table(self, request: DisableTableRequest) -> DisableTableResponse:
        """Stops capturing a table.

        ``retained`` in the response is ``True`` when changes already stored
        for the table remain readable. A table that does not exist in the
        source database raises ``PgChangeFeedNotFoundError``.
        """
        body = {
            "source": request.source,
            "schema": request.schema,
            "table": request.table,
            "publication": request.publication,
        }
        return self._post("/tables/disable", body, DisableTableResponse.from_json)

    # --- GET /tables/status (reader or admin token) ---

    def get_status(
        self, source: str, schema: str, table: str, publication: str
    ) -> TableStatusResponse:
        """Tells whether a table is captured (``enabled``) or no longer captured
        with stored changes remaining (``retained``).

        A table that was never enabled reports both as ``False``; a table that
        does not exist in the source database raises ``PgChangeFeedNotFoundError``.
        """
        return self._get(
            "/tables/status",
            {"source": source, "schema": schema, "table": table, "publication": publication},
            TableStatusResponse.from_json,
        )

    # --- GET /tables (reader or admin token) ---

    def list_tables(self, source: str, publication: str) -> ListTablesResponse:
        """Lists the captured tables (``tables``) and the tables that are no
        longer captured but whose stored changes remain (``retained``)."""
        return self._get(
            "/tables",
            {"source": source, "publication": publication},
            ListTablesResponse.from_json,
        )

    # --- POST /retention/run (admin token) ---

    def run_retention(self, request: RunRetentionRequest) -> RunRetentionResponse:
        """Deletes the stored changes of a source that are older than
        ``min_age_nanos`` and that every consumer with a stored position has
        already passed; ``deleted`` in the response is the number removed."""
        body = {"source": request.source, "min_age_nanos": request.min_age_nanos}
        return self._post("/retention/run", body, RunRetentionResponse.from_json)

    # --- GET /changes (reader or admin token) ---

    def read_changes(
        self,
        source: str,
        schema: str | None = None,
        table: str | None = None,
        from_: int | None = None,
        to: int | None = None,
        limit: int | None = None,
    ) -> ReadChangesResponse:
        """Reads stored changes of a source.

        ``schema`` and ``table`` each narrow the result independently.
        ``from_`` and ``to`` are ``commit_position`` values: ``from_`` is
        inclusive, ``to`` is exclusive. ``from_`` is spelled with a trailing
        underscore because ``from`` is a Python keyword; it is sent as the
        ``from`` query parameter. ``limit`` cuts rows, not positions, and there
        is no default limit. A range without changes returns an empty list.
        """
        return self._get(
            "/changes",
            {
                "source": source,
                "schema": schema,
                "table": table,
                "from": from_,
                "to": to,
                "limit": limit,
            },
            ReadChangesResponse.from_json,
        )

    # --- transport plumbing ---

    def _url(self, path: str) -> str:
        return f"{self._options.address.rstrip('/')}{path}"

    def _headers(self) -> dict[str, str]:
        return {"Authorization": f"Bearer {self._options.api_token}"}

    def _get(self, path: str, params: dict[str, Any], mapper: Callable[[Any], T]) -> T:
        response = self._client.get(
            self._url(path), params=_drop_none(params), headers=self._headers()
        )
        return self._handle(response, mapper)

    def _post(self, path: str, body: dict[str, Any], mapper: Callable[[Any], T]) -> T:
        response = self._client.post(self._url(path), json=body, headers=self._headers())
        return self._handle(response, mapper)

    def _handle(self, response: httpx.Response, mapper: Callable[[Any], T]) -> T:
        if not 200 <= response.status_code < 300:
            raise _build_error(response)
        try:
            payload = response.json()
        except ValueError as exc:
            raise PgChangeFeedMalformedResponseError(
                response.status_code,
                f"PG Change Feed HTTP API returned status {response.status_code} with a response "
                "body that is not valid JSON.",
            ) from exc
        try:
            return mapper(payload)
        except (KeyError, TypeError) as exc:
            raise PgChangeFeedMalformedResponseError(
                response.status_code,
                f"PG Change Feed HTTP API returned status {response.status_code} with a response "
                "body missing an expected field.",
            ) from exc


def _drop_none(params: dict[str, Any]) -> dict[str, Any]:
    return {key: value for key, value in params.items() if value is not None}


def _build_error(response: httpx.Response) -> PgChangeFeedError:
    message = _extract_error_message(response)
    error_cls = _STATUS_TO_ERROR.get(response.status_code, PgChangeFeedUnexpectedStatusError)
    return error_cls(response.status_code, message)


def _extract_error_message(response: httpx.Response) -> str:
    try:
        body = response.json()
    except ValueError:
        return response.text
    if isinstance(body, dict) and "error" in body:
        return str(body["error"])
    return response.text
