"""Real-server rule phase for the NATS live change stream client.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. Subscribes the source's namespace against a
table carrying an active rename_column rule and receives a change committed
afterwards as a full JSON event; the row image arrives with the renamed key
and without the source key, read through the client's opaque value.
"""

from __future__ import annotations

import os
import time

import pgchangefeed.nats_stream_client as nats_stream_client
from pgchangefeed.models import StreamChange
from pgchangefeed.options import ClientOptions

_NATS_URL = os.environ["PGCHANGEFEED_NATS_URL"]
_TOKEN = os.environ["PGCHANGEFEED_NATS_STREAM_TOKEN"]
_SOURCE_ID = os.environ["PGCHANGEFEED_SOURCE_ID"]
_TABLE = os.environ["PGCHANGEFEED_E2E_TABLE"]
_SENTINEL = os.environ["PGCHANGEFEED_E2E_SENTINEL"]
_SOURCE_KEY = os.environ["PGCHANGEFEED_RULE_SOURCE_KEY"]
_TARGET_KEY = os.environ["PGCHANGEFEED_RULE_TARGET_KEY"]

_RECEIVE_DEADLINE_SECONDS = 90.0


def test_realserver_receives_a_renamed_row_image_key_over_the_stream() -> None:
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
    new_image = received.new_image
    assert isinstance(new_image, dict)
    assert _TARGET_KEY in new_image, f"Zielschlüssel fehlt: {new_image!r}"
    assert new_image[_TARGET_KEY] == _SENTINEL
    assert _SOURCE_KEY not in new_image, f"Quellschlüssel noch vorhanden: {new_image!r}"

    print(
        f"RECEIVED change_id={received.change_id} table={received.table} "
        f"operation={received.operation} new_image={received.new_image}",
        flush=True,
    )
