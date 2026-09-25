"""Network-free unit tests for the NATS live change stream client.

The fake replaces the ``nats`` module seam: the tests monkeypatch
``nats.connect`` with an async fake that records the connect options, hands
back a fake connection whose ``subscribe`` records the subject and delivers
its pre-arranged payload through the captured callback — the same observable
behavior a real connection carries (subscribe → callback → generator),
without a socket. Message completeness is asserted against the ten stream
fields in JSON form, the same schema the SSE event carries.
"""

from __future__ import annotations

import json
import time
from typing import Any, Callable, Coroutine

import pytest

import pgchangefeed.nats_stream_client as nats_stream_client
from pgchangefeed.exceptions import PgChangeFeedMalformedResponseError
from pgchangefeed.models import StreamChange
from pgchangefeed.options import ClientOptions

NATS_URL = "nats://nats:4222"
TOKEN = "e2e-nats-stream-token"
SOURCE_ID = "src-e2e"

_STREAM_MESSAGE = {
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

_FULL_MESSAGE_BYTES = json.dumps(_STREAM_MESSAGE).encode()


class _Message:
    def __init__(self, payload: bytes) -> None:
        self.data = payload


class _FakeNatsConnection:
    """Records the subscription subject and feeds its pre-arranged payload
    through the captured callback — the moment the subscription is
    registered."""

    def __init__(self, on_subscribe_payload: bytes | None = None) -> None:
        self.subscribe_subjects: list[str] = []
        self.on_subscribe_payload = on_subscribe_payload

    async def subscribe(self, subject: str, cb: Callable[[Any], Any]) -> None:
        self.subscribe_subjects.append(subject)
        if self.on_subscribe_payload is not None:
            await cb(_Message(self.on_subscribe_payload))


def _install_fake_connect(
    monkeypatch,
    connection: _FakeNatsConnection | None,
    error: Exception | None,
) -> list[tuple[list[str], str | None]]:
    """Replaces ``nats.connect`` with an async fake recording the connect
    arguments (servers, token) — the seam the client touches."""
    connect_calls: list[tuple[list[str], str | None]] = []

    async def fake_connect(*, servers: list[str], token: str | None = None, **kwargs: Any):
        connect_calls.append((list(servers), token))
        if error is not None:
            raise error
        return connection

    monkeypatch.setattr(nats_stream_client.nats, "connect", fake_connect)
    return connect_calls


def _make_client() -> nats_stream_client.PgChangeFeedNatsStreamClient:
    options = ClientOptions(address=NATS_URL, api_token=TOKEN)
    return nats_stream_client.PgChangeFeedNatsStreamClient(options, SOURCE_ID)


# --- Subject format: one wildcard namespace per source ---


def test_subject_namespace_is_the_wildcard_of_the_source() -> None:
    assert nats_stream_client._subject_namespace(SOURCE_ID) == f"cdc.stream.{SOURCE_ID}.>"
    assert nats_stream_client._subject_namespace("src-x") == "cdc.stream.src-x.>"


def test_client_rejects_an_empty_source_id() -> None:
    options = ClientOptions(address=NATS_URL, api_token=TOKEN)
    with pytest.raises(ValueError):
        nats_stream_client.PgChangeFeedNatsStreamClient(options, "   ")


# --- Message schema: ten fields, the same as the SSE event ---


def test_parse_stream_change_maps_the_ten_fields() -> None:
    change = nats_stream_client._parse_stream_change(_FULL_MESSAGE_BYTES)
    assert type(change) is StreamChange
    assert change.change_id == "chg-1"
    assert change.transaction_id == "tx-1"
    assert change.source_table_id == "tbl-1"
    assert change.sequence == 1
    assert change.operation == "INSERT"
    assert change.old_image is None
    assert change.new_image == {"id": 1, "name": "first"}
    assert change.schema_version == "sv-1"
    assert change.schema == "public"
    assert change.table == "orders"


def test_parse_stream_change_surfaces_typed_errors_for_schema_violations() -> None:
    for payload in (
        b"not json at all",
        b"[1, 2, 3]",
        b'{"unexpected": "shape"}',
    ):
        with pytest.raises(PgChangeFeedMalformedResponseError):
            nats_stream_client._parse_stream_change(payload)


# --- Wiring: connection-level token, subject subscription ---


def test_stream_changes_subscribes_the_namespace_with_the_connection_token(monkeypatch) -> None:
    connection = _FakeNatsConnection(on_subscribe_payload=_FULL_MESSAGE_BYTES)
    connect_calls = _install_fake_connect(monkeypatch, connection, None)
    client = _make_client()

    # Der Generator faengt erst beim ersten next() an; die vorarrangierte
    # Message landet waehrend der Anmeldung in der Queue.
    iterator = client.stream_changes(timeout=3.0)
    yielded = next(iterator)
    assert connect_calls == [([NATS_URL], TOKEN)]
    assert connection.subscribe_subjects == [f"cdc.stream.{SOURCE_ID}.>"]
    assert type(yielded) is StreamChange
    assert yielded.change_id == "chg-1"


# --- Connection error path: rejection happens when the connection is opened ---


def test_connect_rejection_surfaces_and_is_not_swallowed(monkeypatch) -> None:
    rejection = RuntimeError("authorization violation")
    _install_fake_connect(monkeypatch, None, rejection)
    client = _make_client()
    with pytest.raises(RuntimeError, match="authorization violation"):
        list(client.stream_changes(timeout=2.0))


def test_timeout_bounds_the_consumption(monkeypatch) -> None:
    connection = _FakeNatsConnection()  # keine Message vorarrangiert
    _install_fake_connect(monkeypatch, connection, None)
    client = _make_client()
    start = time.monotonic()
    assert list(client.stream_changes(timeout=0.3)) == []
    assert time.monotonic() - start < 2.0
