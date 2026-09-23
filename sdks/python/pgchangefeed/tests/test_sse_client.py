"""Network-free unit tests for the SSE stream client surface (SPEC-021).

The fake is the same seam the HTTP surface uses: an ``httpx.MockTransport``
carries the raw response bytes, so the tests exercise the real streaming
pipeline (response stream, line iterator, frame parser, JSON mapping)
without a socket. Chunk-boundary independence is asserted explicitly: the
transport hands out the frame bytes split at arbitrary boundaries, because
chunk boundaries need not coincide with line boundaries (SPEC-021 risk).
"""

from __future__ import annotations

import json
from typing import Callable

import httpx
import pytest

from pgchangefeed.exceptions import (
    PgChangeFeedMalformedResponseError,
    PgChangeFeedUnexpectedStatusError,
    PgChangeFeedUnauthorizedError,
)
from pgchangefeed.models import StreamChange
from pgchangefeed.options import ClientOptions
from pgchangefeed.sse_client import PgChangeFeedSseClient

ADDRESS = "https://feed.example.invalid"
TOKEN = "e2e-reader-token"

_EVENT_BODY = {
    "change_id": "chg-1",
    "transaction_id": "tx-1",
    "source_table_id": "tbl-1",
    "sequence": 1,
    "operation": "INSERT",
    "old_image": None,
    "new_image": {"id": 1, "name": "first"},
    "schema_version": "sv-1",
    "schema": "public",
    "table": "orders",
}


def _change_event(body: object, event_type: str = "change") -> str:
    return f"event: {event_type}\ndata: {json.dumps(body)}\n\n"


def _make_client(handler: Callable[[httpx.Request], httpx.Response]) -> PgChangeFeedSseClient:
    transport = httpx.MockTransport(handler)
    http_client = httpx.Client(transport=transport)
    options = ClientOptions(address=ADDRESS, api_token=TOKEN)
    return PgChangeFeedSseClient(http_client, options)


# --- Happy path (SPEC-021: one event per row change, ten fields) ---


def test_stream_changes_yields_typed_events_with_the_ten_fields() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.method == "GET"
        assert request.url.path == "/changes/stream"
        assert request.headers["authorization"] == f"Bearer {TOKEN}"
        return httpx.Response(
            200,
            headers={"Content-Type": "text/event-stream"},
            content=(
                _change_event(_EVENT_BODY).encode() + _change_event(_EVENT_BODY).encode()
            ),
        )

    client = _make_client(handler)
    events = list(client.stream_changes())
    assert len(events) == 2
    for event in events:
        assert type(event) is StreamChange
        assert event.change_id == "chg-1"
        assert event.transaction_id == "tx-1"
        assert event.source_table_id == "tbl-1"
        assert event.sequence == 1
        assert event.operation == "INSERT"
        assert event.old_image is None
        assert event.new_image == {"id": 1, "name": "first"}
        assert event.schema_version == "sv-1"
        assert event.schema == "public"
        assert event.table == "orders"


def test_missing_image_is_none_and_update_carries_both_images() -> None:
    update_body = dict(_EVENT_BODY)
    update_body["operation"] = "UPDATE"
    update_body["old_image"] = {"id": 1, "name": "first"}
    update_body["new_image"] = {"id": 1, "name": "second"}

    def handler(request: httpx.Request) -> httpx.Response:
        return httpx.Response(
            200,
            headers={"Content-Type": "text/event-stream"},
            content=_change_event(update_body).encode(),
        )

    client = _make_client(handler)
    events = list(client.stream_changes())
    assert len(events) == 1
    assert events[0].old_image == {"id": 1, "name": "first"}
    assert events[0].new_image == {"id": 1, "name": "second"}


# --- Frame parser (chunk boundaries, incomplete frames, non-change names) ---


class _ChunkedStream(httpx.SyncByteStream):
    """Hands out the frame bytes split at arbitrary boundaries — chunk
    Grenzen müssen nicht mit Zeilen-Grenzen zusammenfallen (SPEC-021
    Risiko)."""

    def __init__(self, chunks: list[bytes]) -> None:
        self._chunks = chunks

    def __iter__(self):
        yield from self._chunks


def test_chunk_boundaries_need_not_coincide_with_lines() -> None:
    encoded = _change_event(_EVENT_BODY).encode() + _change_event(_EVENT_BODY).encode()
    split = len(encoded) // 2

    def handler(request: httpx.Request) -> httpx.Response:
        return httpx.Response(
            200,
            headers={"Content-Type": "text/event-stream"},
            stream=_ChunkedStream([encoded[:split], encoded[split:]]),
        )

    client = _make_client(handler)
    events = list(client.stream_changes())
    assert len(events) == 2


def test_stream_ending_before_the_blank_line_discards_the_begun_frame() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        return httpx.Response(
            200,
            headers={"Content-Type": "text/event-stream"},
            content=f"event: change\ndata: {json.dumps(_EVENT_BODY)}\n".encode(),
        )

    client = _make_client(handler)
    assert list(client.stream_changes()) == []


def test_non_change_event_names_are_skipped() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        return httpx.Response(
            200,
            headers={"Content-Type": "text/event-stream"},
            content=(
                b"event: keepalive\ndata: {}\n\n" + _change_event(_EVENT_BODY).encode()
            ),
        )

    client = _make_client(handler)
    events = list(client.stream_changes())
    assert len(events) == 1
    assert events[0].change_id == "chg-1"


# --- Malformed payloads: typed error, never a raw parsing exception ---


def test_non_json_data_payload_raises_typed_malformed_error() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        return httpx.Response(
            200,
            headers={"Content-Type": "text/event-stream"},
            content=b"event: change\ndata: not json at all\n\n",
        )

    client = _make_client(handler)
    with pytest.raises(PgChangeFeedMalformedResponseError):
        list(client.stream_changes())


def test_data_payload_missing_expected_field_raises_typed_malformed_error() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        return httpx.Response(
            200,
            headers={"Content-Type": "text/event-stream"},
            content=_change_event({"unexpected": "shape"}).encode(),
        )

    client = _make_client(handler)
    with pytest.raises(PgChangeFeedMalformedResponseError):
        list(client.stream_changes())


def test_data_payload_not_a_json_object_raises_typed_malformed_error() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        return httpx.Response(
            200,
            headers={"Content-Type": "text/event-stream"},
            content=_change_event([1, 2, 3]).encode(),
        )

    client = _make_client(handler)
    with pytest.raises(PgChangeFeedMalformedResponseError):
        list(client.stream_changes())


# --- Status-code mapping (SPEC-021: 401 ohne Token, 503 ohne Broadcaster) ---


def test_unknown_token_surfaces_as_typed_unauthorized_error() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        auth = request.headers.get("authorization", "")
        if auth != f"Bearer {TOKEN}":
            return httpx.Response(401, json={"error": "unknown bearer token"})
        return httpx.Response(200, headers={"Content-Type": "text/event-stream"})

    transport = httpx.MockTransport(handler)
    http_client = httpx.Client(transport=transport)
    options = ClientOptions(address=ADDRESS, api_token="unknown-token")
    client = PgChangeFeedSseClient(http_client, options)
    with pytest.raises(PgChangeFeedUnauthorizedError) as exc_info:
        list(client.stream_changes())
    assert exc_info.value.status_code == 401


def test_503_without_broadcaster_surfaces_as_typed_error() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        return httpx.Response(503, json={"error": "kein Broadcaster verdrahtet"})

    client = _make_client(handler)
    with pytest.raises(PgChangeFeedUnexpectedStatusError) as exc_info:
        list(client.stream_changes())
    assert exc_info.value.status_code == 503