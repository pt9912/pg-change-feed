from __future__ import annotations

from typing import Callable

import httpx
import pytest

from pgchangefeed.exceptions import (
    PgChangeFeedBadRequestError,
    PgChangeFeedForbiddenError,
    PgChangeFeedMalformedResponseError,
    PgChangeFeedNotFoundError,
    PgChangeFeedServerError,
    PgChangeFeedUnauthorizedError,
    PgChangeFeedUnexpectedStatusError,
)
from pgchangefeed.http_client import PgChangeFeedHttpClient
from pgchangefeed.models import (
    AcknowledgeConsumerRequest,
    DisableTableRequest,
    EnableTableRequest,
    RegisterConsumerRequest,
    RunRetentionRequest,
)
from pgchangefeed.options import ClientOptions

ADDRESS = "https://feed.example.invalid"
ADMIN_TOKEN = "admin-token"
READER_TOKEN = "reader-token"


def _make_client(
    handler: Callable[[httpx.Request], httpx.Response], token: str = ADMIN_TOKEN
) -> PgChangeFeedHttpClient:
    transport = httpx.MockTransport(handler)
    http_client = httpx.Client(transport=transport)
    options = ClientOptions(address=ADDRESS, api_token=token)
    return PgChangeFeedHttpClient(http_client, options)


def _handler_returning(status_code: int, payload: object) -> Callable[[httpx.Request], httpx.Response]:
    def handler(request: httpx.Request) -> httpx.Response:
        return httpx.Response(status_code, json=payload)

    return handler


# --- Happy path, one per API capability including read_changes ---


def test_register_consumer_happy_path() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.method == "POST"
        assert request.url.path == "/consumers"
        assert request.headers["authorization"] == f"Bearer {ADMIN_TOKEN}"
        return httpx.Response(
            201,
            json={"consumer_id": "c1", "name": "Consumer 1", "already_registered": False},
        )

    client = _make_client(handler)
    response = client.register_consumer(RegisterConsumerRequest(consumer_id="c1", name="Consumer 1"))
    assert response.consumer_id == "c1"
    assert response.name == "Consumer 1"
    assert response.already_registered is False


def test_acknowledge_consumer_happy_path() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.method == "POST"
        assert request.url.path == "/consumers/acknowledge"
        return httpx.Response(
            200,
            json={"consumer_id": "c1", "source_id": "s1", "offset": 42},
        )

    client = _make_client(handler)
    response = client.acknowledge_consumer(
        AcknowledgeConsumerRequest(consumer_id="c1", source_id="s1", offset=42)
    )
    assert response.consumer_id == "c1"
    assert response.source_id == "s1"
    assert response.offset == 42


def test_get_consumer_position_happy_path() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.method == "GET"
        assert request.url.path == "/consumers/position"
        assert request.url.params["consumer_id"] == "c1"
        return httpx.Response(
            200,
            json={"consumer_id": "c1", "source_id": "s1", "offset": 7, "acknowledged": True},
        )

    client = _make_client(handler)
    response = client.get_consumer_position("c1")
    assert response.offset == 7
    assert response.acknowledged is True


def test_remove_consumer_happy_path() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.method == "POST"
        assert request.url.path == "/consumers/remove"
        return httpx.Response(200, json={"consumer_id": "c1", "removed": True})

    client = _make_client(handler)
    response = client.remove_consumer("c1")
    assert response.removed is True


def test_enable_table_happy_path() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.method == "POST"
        assert request.url.path == "/tables/enable"
        return httpx.Response(
            201,
            json={
                "table_id": "t1",
                "source": "src",
                "schema": "public",
                "table": "orders",
                "already_enabled": False,
            },
        )

    client = _make_client(handler)
    response = client.enable_table(
        EnableTableRequest(
            source="src",
            schema="public",
            table="orders",
            table_id="t1",
            schema_version_id="sv1",
            version=1,
            publication="pub",
        )
    )
    assert response.table_id == "t1"
    assert response.already_enabled is False


def test_disable_table_happy_path() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.method == "POST"
        assert request.url.path == "/tables/disable"
        return httpx.Response(200, json={"removed": True, "retained": False})

    client = _make_client(handler)
    response = client.disable_table(
        DisableTableRequest(source="src", schema="public", table="orders", publication="pub")
    )
    assert response.removed is True
    assert response.retained is False


def test_get_status_happy_path() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.method == "GET"
        assert request.url.path == "/tables/status"
        assert request.url.params["source"] == "src"
        assert request.url.params["publication"] == "pub"
        return httpx.Response(200, json={"enabled": True, "retained": False})

    client = _make_client(handler)
    response = client.get_status(source="src", schema="public", table="orders", publication="pub")
    assert response.enabled is True
    assert response.retained is False


def test_list_tables_happy_path() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.method == "GET"
        assert request.url.path == "/tables"
        return httpx.Response(
            200,
            json={
                "tables": [{"table_id": "t1", "source": "src", "schema": "public", "table": "orders"}],
                "retained": [],
            },
        )

    client = _make_client(handler)
    response = client.list_tables(source="src", publication="pub")
    assert len(response.tables) == 1
    assert response.tables[0].table_id == "t1"
    assert response.retained == []


def test_run_retention_happy_path() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.method == "POST"
        assert request.url.path == "/retention/run"
        return httpx.Response(200, json={"deleted": 3})

    client = _make_client(handler)
    response = client.run_retention(RunRetentionRequest(source="src", min_age_nanos=0))
    assert response.deleted == 3


def test_read_changes_happy_path() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.method == "GET"
        assert request.url.path == "/changes"
        assert request.url.params["source"] == "src"
        assert request.url.params["from"] == "1"
        assert "limit" not in request.url.params
        return httpx.Response(
            200,
            json={
                "changes": [
                    {
                        "commit_position": 1,
                        "change_id": "chg1",
                        "transaction_id": "tx1",
                        "source_table_id": "t1",
                        "schema": "public",
                        "table": "orders",
                        "sequence": 1,
                        "operation": "INSERT",
                        "old_image": None,
                        "new_image": {"id": 1},
                        "schema_version": "sv1",
                        "committed_at": "2026-09-19T00:00:00Z",
                    }
                ]
            },
        )

    client = _make_client(handler)
    response = client.read_changes(source="src", from_=1)
    assert len(response.changes) == 1
    change = response.changes[0]
    assert change.operation == "INSERT"
    assert change.new_image == {"id": 1}
    assert change.old_image is None


# Body shape: the output of the server's GET /changes handler, where `origin`
# is the last field of each change (Go test TestReadChangesTraegtOriginAlsLetztesFeld,
# internal/adapters/driving/http/readchanges_test.go). Not captured from a
# running server.
_CHANGE_WITHOUT_ORIGIN = {
    "commit_position": 1,
    "change_id": "chg1",
    "transaction_id": "tx1",
    "source_table_id": "t1",
    "schema": "public",
    "table": "orders",
    "sequence": 0,
    "operation": "INSERT",
    "old_image": None,
    "new_image": {"id": 1},
    "schema_version": "sv1",
    "committed_at": "2026-09-19T00:00:00Z",
}


@pytest.mark.parametrize(
    ("origin_fields", "expected"),
    [
        ({"origin": "wal"}, "wal"),
        ({"origin": "backfill"}, "backfill"),
        ({"origin": "future-kind"}, "future-kind"),
        ({}, "wal"),
        ({"origin": None}, "wal"),
        ({"origin": ""}, ""),
    ],
    ids=["wal", "backfill", "unknown-value", "field-absent", "json-null", "empty-value"],
)
def test_read_changes_origin_reads_server_value_and_defaults_to_wal(
    origin_fields: dict[str, object], expected: str
) -> None:
    payload = {"changes": [{**_CHANGE_WITHOUT_ORIGIN, **origin_fields}]}
    client = _make_client(_handler_returning(200, payload))
    response = client.read_changes(source="src")
    assert response.changes[0].origin == expected


def test_read_changes_empty_result_is_not_an_error() -> None:
    client = _make_client(_handler_returning(200, {"changes": []}))
    response = client.read_changes(source="src")
    assert response.changes == []


# --- Auth boundary: 401 missing/unknown token, 403 reader token on an admin call ---


def _auth_aware_handler(admin_only_paths: set[str]) -> Callable[[httpx.Request], httpx.Response]:
    def handler(request: httpx.Request) -> httpx.Response:
        auth = request.headers.get("authorization", "")
        if not auth.startswith("Bearer "):
            return httpx.Response(401, json={"error": "missing bearer token"})
        token = auth.removeprefix("Bearer ")
        if token not in (ADMIN_TOKEN, READER_TOKEN):
            return httpx.Response(401, json={"error": "unknown bearer token"})
        if request.url.path in admin_only_paths and token != ADMIN_TOKEN:
            return httpx.Response(403, json={"error": "insufficient rights"})
        return httpx.Response(200, json={"consumer_id": "c1", "removed": False})

    return handler


def test_unknown_token_returns_typed_unauthorized_error() -> None:
    client = _make_client(_auth_aware_handler({"/consumers/remove"}), token="unknown-token")
    with pytest.raises(PgChangeFeedUnauthorizedError) as exc_info:
        client.remove_consumer("c1")
    assert exc_info.value.status_code == 401


def test_reader_token_against_admin_endpoint_returns_typed_forbidden_error() -> None:
    client = _make_client(_auth_aware_handler({"/consumers/remove"}), token=READER_TOKEN)
    with pytest.raises(PgChangeFeedForbiddenError) as exc_info:
        client.remove_consumer("c1")
    assert exc_info.value.status_code == 403


# --- Remaining status-code mapping: 400/404/500 and the fallback for other codes ---


def test_400_returns_typed_bad_request_error_with_error_message() -> None:
    client = _make_client(_handler_returning(400, {"error": "source ist Pflichtfeld"}))
    with pytest.raises(PgChangeFeedBadRequestError) as exc_info:
        client.read_changes(source="src")
    assert exc_info.value.status_code == 400
    assert "Pflichtfeld" in str(exc_info.value)


def test_404_returns_typed_not_found_error() -> None:
    client = _make_client(_handler_returning(404, {"error": "Tabelle existiert nicht an der Quelle"}))
    with pytest.raises(PgChangeFeedNotFoundError) as exc_info:
        client.get_status(source="src", schema="public", table="missing", publication="pub")
    assert exc_info.value.status_code == 404


def test_500_returns_typed_server_error() -> None:
    client = _make_client(_handler_returning(500, {"error": "interner Fehler"}))
    with pytest.raises(PgChangeFeedServerError) as exc_info:
        client.list_tables(source="src", publication="pub")
    assert exc_info.value.status_code == 500


def test_undocumented_status_returns_typed_unexpected_status_error() -> None:
    client = _make_client(_handler_returning(409, {"error": "conflict"}))
    with pytest.raises(PgChangeFeedUnexpectedStatusError) as exc_info:
        client.list_tables(source="src", publication="pub")
    assert exc_info.value.status_code == 409


def test_error_body_without_error_field_falls_back_to_raw_text() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        return httpx.Response(400, content=b"not json at all")

    client = _make_client(handler)
    with pytest.raises(PgChangeFeedBadRequestError) as exc_info:
        client.list_tables(source="src", publication="pub")
    assert "not json at all" in str(exc_info.value)


# --- Malformed 2xx body: must raise a typed error, never a raw parsing exception ---


def test_non_json_success_body_raises_typed_malformed_response_error() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        return httpx.Response(200, content=b"not json at all")

    client = _make_client(handler)
    with pytest.raises(PgChangeFeedMalformedResponseError) as exc_info:
        client.list_tables(source="src", publication="pub")
    assert exc_info.value.status_code == 200


def test_success_body_missing_expected_field_raises_typed_malformed_response_error() -> None:
    client = _make_client(_handler_returning(200, {"unexpected": "shape"}))
    with pytest.raises(PgChangeFeedMalformedResponseError) as exc_info:
        client.list_tables(source="src", publication="pub")
    assert exc_info.value.status_code == 200
