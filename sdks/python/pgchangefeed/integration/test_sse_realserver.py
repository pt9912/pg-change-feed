"""Real-server integration test for the SSE stream client.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. It proves the client's protocol assumptions
against a real, running server instance, not only against a fake. The runner
drives the environment: it
brings the compose stack up, activates the tables and inserts rows while
this test is streaming; the printed markers
(``READY``/``RECEIVED``/``REJECTED``) are what the runner asserts on, so
every ``print`` is flushed and unbuffered.

The stream is fire-and-forget: between the stream opening and the server
registering the subscriber, a committed change can be discarded for this
subscriber. The runner commits
a bounded sequence of unique rows until exactly one arrives; this test keeps
receiving until it sees its sentinel, so a missed row costs nothing.
"""

from __future__ import annotations

import os
import time

import httpx
import pytest

from pgchangefeed.exceptions import PgChangeFeedUnauthorizedError
from pgchangefeed.models import StreamChange
from pgchangefeed.options import ClientOptions
from pgchangefeed.sse_client import PgChangeFeedSseClient

_ADDR = os.environ["PGCHANGEFEED_HTTP_ADDR"]
_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN"]
_TABLE = os.environ["PGCHANGEFEED_E2E_TABLE"]
_SENTINEL = os.environ["PGCHANGEFEED_E2E_SENTINEL"]

_RECEIVE_DEADLINE_SECONDS = 90.0


def test_realserver_receives_a_committed_change_over_the_stream() -> None:
    http_client = httpx.Client(timeout=95.0)
    client = PgChangeFeedSseClient(http_client, ClientOptions(address=_ADDR, api_token=_TOKEN))
    stream = client.stream_changes()
    print("READY", flush=True)

    received: StreamChange | None = None
    deadline = time.monotonic() + _RECEIVE_DEADLINE_SECONDS
    for received in stream:
        if _TABLE == received.table and _SENTINEL in str(received.new_image):
            break
        if time.monotonic() > deadline:
            pytest.fail(
                f"kein Event mit dem Sentinel {_SENTINEL!r} auf dem Stream "
                f"innerhalb der Frist empfangen"
            )

    assert received is not None
    # Feldvollstaendigkeit am realen Wire-Image: alle zehn Felder
    # je in seinem Wire-Typ, Identitaeten non-empty, INSERT traegt kein
    # Alt-Bild.
    for field, kind in (
        ("change_id", str),
        ("transaction_id", str),
        ("source_table_id", str),
        ("sequence", int),
        ("operation", str),
        ("old_image", type(None)),
        ("new_image", dict),
        ("schema_version", str),
        ("schema", str),
        ("table", str),
    ):
        value = getattr(received, field)
        assert isinstance(value, kind) and not isinstance(value, bool), (
            f"Feld {field}: {type(value).__name__}, wollen {kind.__name__}"
        )
    assert received.change_id != ""
    assert received.transaction_id != ""
    assert received.source_table_id != ""
    assert received.schema_version != ""
    assert received.operation == "INSERT"
    assert _SENTINEL in str(received.new_image)
    assert received.table == _TABLE
    assert received.schema == "public"

    print(
        f"RECEIVED change_id={received.change_id} table={received.table} "
        f"operation={received.operation} new_image={received.new_image}",
        flush=True,
    )


def test_realserver_rejects_the_stream_without_token() -> None:
    http_client = httpx.Client(timeout=15.0)
    client = PgChangeFeedSseClient(
        http_client, ClientOptions(address=_ADDR, api_token="no-such-token")
    )
    with pytest.raises(PgChangeFeedUnauthorizedError) as exc_info:
        list(client.stream_changes())
    assert exc_info.value.status_code == 401
    print("REJECTED status=401", flush=True)