"""Real-server integration test for the NATS live change stream client.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. It proves the client's protocol assumptions
against a real, running server instance, not only against a fake. The runner
drives the environment: it
brings the compose stack up, activates the tables and inserts rows while
this test is subscribed; the printed markers
(``READY``/``RECEIVED``/``REJECTED``) are what the runner asserts on, so
every ``print`` is flushed and unbuffered.

The stream is fire-and-forget: between the subscription and the server's
publish of a change, a message can be missed —
the runner commits a bounded sequence of unique rows until exactly one
arrives; this test keeps consuming until it sees its sentinel, so a missed
row costs nothing.
"""

from __future__ import annotations

import asyncio
import os
import time

import pytest

import nats
from nats import errors as nats_errors

import pgchangefeed.nats_stream_client as nats_stream_client
from pgchangefeed.models import StreamChange
from pgchangefeed.options import ClientOptions

_NATS_URL = os.environ["PGCHANGEFEED_NATS_URL"]
_TOKEN = os.environ["PGCHANGEFEED_NATS_STREAM_TOKEN"]
_SOURCE_ID = os.environ["PGCHANGEFEED_SOURCE_ID"]
_TABLE = os.environ["PGCHANGEFEED_E2E_TABLE"]
_SENTINEL = os.environ["PGCHANGEFEED_E2E_SENTINEL"]

_RECEIVE_DEADLINE_SECONDS = 90.0


def test_realserver_receives_a_committed_change_over_the_stream() -> None:
    client = nats_stream_client.PgChangeFeedNatsStreamClient(
        ClientOptions(address=_NATS_URL, api_token=_TOKEN), _SOURCE_ID
    )
    stream = client.stream_changes(timeout=_RECEIVE_DEADLINE_SECONDS)
    print("READY", flush=True)

    received: StreamChange | None = None
    deadline = time.monotonic() + _RECEIVE_DEADLINE_SECONDS
    for candidate in stream:
        if _TABLE == candidate.table and _SENTINEL in str(candidate.new_image):
            received = candidate
            break
        assert time.monotonic() < deadline, (
            f"kein Event mit dem Sentinel {_SENTINEL!r} auf dem Subjekt "
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


def test_realserver_rejects_the_connection_with_wrong_token() -> None:
    async def attempt() -> None:
        await nats.connect(
            servers=[_NATS_URL],
            token="wrong-token",
            connect_timeout=5.0,
            allow_reconnect=False,
        )

    with pytest.raises(nats_errors.Error, match="Authorization Violation") as exc_info:
        asyncio.run(attempt())
    print(f"REJECTED token-rejected: {exc_info.value}", flush=True)
