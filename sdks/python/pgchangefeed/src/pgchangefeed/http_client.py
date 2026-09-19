"""Public entry point for the PG Change Feed HTTP/JSON API.

One method per wire capability: the nine port-covered capabilities of
SPEC-018 (``RegisterConsumer``, ``AcknowledgeConsumer``,
``GetConsumerPosition``, ``RemoveConsumer``, ``EnableTable``,
``DisableTable``, ``GetStatus``, ``ListTables``, ``RunRetention``) plus
reading persisted changes (``ReadChanges``, SPEC-022). Requests/responses
are typed data classes that mirror the SPEC-018/SPEC-022 JSON schemas
exactly (``pgchangefeed.models``); every non-success response becomes a
typed ``PgChangeFeedError`` subclass instead of a raw ``httpx`` exception,
and a success (``2xx``) response whose body does not match the documented
shape -- invalid JSON, or valid JSON missing an expected field -- becomes
a typed ``PgChangeFeedMalformedResponseError`` instead of letting a raw
parsing exception leak through. Both paths run through the same
``_handle`` helper, so this holds uniformly across every method rather
than being reimplemented ten times.

The ``httpx.Client`` is injected, not owned -- the caller controls its
lifetime, connection pooling and transport (including a
``httpx.MockTransport`` for tests); this type never closes it. The bearer
token and server address come from ``ClientOptions``, supplied at
construction -- no module-level or global state, a process can hold
several independently configured instances at once.

Draht-Kenntnis-Quelle: ``spec/pflichtenheft.md`` SPEC-018/SPEC-022
(direkt), im Zweifel der Go-Server-Adapter selbst (nur gelesen, nicht
importiert -- kein Python-Import eines privaten Baums dieses Repos) --
kein ``examples/python/``-Referenz-Client existiert (ADR-0107 §Kontext).
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
    """Client for the nine SPEC-018 capabilities plus ``ReadChanges`` (SPEC-022)."""

    def __init__(self, client: httpx.Client, options: ClientOptions) -> None:
        self._client = client
        self._options = options

    # --- RegisterConsumer -- POST /consumers (admin, LH-FA-CON-001) ---

    def register_consumer(self, request: RegisterConsumerRequest) -> RegisterConsumerResponse:
        body = {"consumer_id": request.consumer_id, "name": request.name}
        return self._post("/consumers", body, RegisterConsumerResponse.from_json)

    # --- AcknowledgeConsumer -- POST /consumers/acknowledge (admin, LH-FA-CON-004) ---

    def acknowledge_consumer(
        self, request: AcknowledgeConsumerRequest
    ) -> AcknowledgeConsumerResponse:
        body = {
            "consumer_id": request.consumer_id,
            "source_id": request.source_id,
            "offset": request.offset,
        }
        return self._post("/consumers/acknowledge", body, AcknowledgeConsumerResponse.from_json)

    # --- GetConsumerPosition -- GET /consumers/position (reader|admin, LH-FA-CON-003/005) ---

    def get_consumer_position(self, consumer_id: str) -> ConsumerPositionResponse:
        return self._get(
            "/consumers/position",
            {"consumer_id": consumer_id},
            ConsumerPositionResponse.from_json,
        )

    # --- RemoveConsumer -- POST /consumers/remove (admin, LH-FA-CON-006) ---

    def remove_consumer(self, consumer_id: str) -> RemoveConsumerResponse:
        return self._post(
            "/consumers/remove", {"consumer_id": consumer_id}, RemoveConsumerResponse.from_json
        )

    # --- EnableTable -- POST /tables/enable (admin, LH-FA-CFG-001) ---

    def enable_table(self, request: EnableTableRequest) -> EnableTableResponse:
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

    # --- DisableTable -- POST /tables/disable (admin, LH-FA-CFG-002) ---

    def disable_table(self, request: DisableTableRequest) -> DisableTableResponse:
        body = {
            "source": request.source,
            "schema": request.schema,
            "table": request.table,
            "publication": request.publication,
        }
        return self._post("/tables/disable", body, DisableTableResponse.from_json)

    # --- GetStatus -- GET /tables/status (reader|admin, LH-FA-CFG-003) ---

    def get_status(
        self, source: str, schema: str, table: str, publication: str
    ) -> TableStatusResponse:
        return self._get(
            "/tables/status",
            {"source": source, "schema": schema, "table": table, "publication": publication},
            TableStatusResponse.from_json,
        )

    # --- ListTables -- GET /tables (reader|admin, LH-FA-CFG-004) ---

    def list_tables(self, source: str, publication: str) -> ListTablesResponse:
        return self._get(
            "/tables",
            {"source": source, "publication": publication},
            ListTablesResponse.from_json,
        )

    # --- RunRetention -- POST /retention/run (admin, LH-FA-RET-002..004) ---

    def run_retention(self, request: RunRetentionRequest) -> RunRetentionResponse:
        body = {"source": request.source, "min_age_nanos": request.min_age_nanos}
        return self._post("/retention/run", body, RunRetentionResponse.from_json)

    # --- ReadChanges -- GET /changes (reader|admin, SPEC-022) ---

    def read_changes(
        self,
        source: str,
        schema: str | None = None,
        table: str | None = None,
        from_: int | None = None,
        to: int | None = None,
        limit: int | None = None,
    ) -> ReadChangesResponse:
        """``from``/``to`` are ``commit_position`` values, ``from``
        inclusive and ``to`` exclusive (SPEC-022); ``from_`` avoids
        shadowing the Python keyword ``from`` while sending the ``from``
        query parameter the wire contract expects. There is no default
        ``limit``.
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
                "body that is not valid JSON -- a protocol violation outside SPEC-018/SPEC-022's "
                "documented shapes.",
            ) from exc
        try:
            return mapper(payload)
        except (KeyError, TypeError) as exc:
            raise PgChangeFeedMalformedResponseError(
                response.status_code,
                f"PG Change Feed HTTP API returned status {response.status_code} with a response "
                "body missing an expected SPEC-018/SPEC-022 field -- a protocol violation outside "
                "the documented shapes.",
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
