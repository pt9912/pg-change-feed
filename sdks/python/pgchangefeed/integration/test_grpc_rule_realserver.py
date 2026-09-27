"""Real-server rule phase for the gRPC stream client.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. Opens the server stream against a table
carrying an active rename_column rule and receives a change committed
afterwards; the row image arrives with the renamed key and without the
source key, read through the client's opaque wire form (the raw bytes,
parsed as JSON here only to assert on the property under test — the client
itself never interprets the row image shape).
"""

from __future__ import annotations

import json
import os
import time

import grpc

from pgchangefeed.grpc_client import PgChangeFeedGrpcClient
from pgchangefeed.grpc_gen import changestream_pb2
from pgchangefeed.options import ClientOptions

_ADDR = os.environ["PGCHANGEFEED_GRPC_ADDR"]
_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN"]
_TABLE = os.environ["PGCHANGEFEED_E2E_TABLE"]
_SENTINEL = os.environ["PGCHANGEFEED_E2E_SENTINEL"]
_SOURCE_KEY = os.environ["PGCHANGEFEED_RULE_SOURCE_KEY"]
_TARGET_KEY = os.environ["PGCHANGEFEED_RULE_TARGET_KEY"]

_RECEIVE_DEADLINE_SECONDS = 90.0


def test_realserver_receives_a_renamed_row_image_key_over_the_stream() -> None:
    channel = grpc.insecure_channel(_ADDR)
    client = PgChangeFeedGrpcClient(channel, ClientOptions(address=_ADDR, api_token=_TOKEN))
    stream = client.stream_changes(timeout=_RECEIVE_DEADLINE_SECONDS)
    print("READY", flush=True)

    received: changestream_pb2.Change | None = None
    deadline = time.monotonic() + _RECEIVE_DEADLINE_SECONDS
    while received is None:
        assert time.monotonic() < deadline, (
            f"kein Change mit dem Sentinel {_SENTINEL!r} auf dem Stream "
            f"innerhalb der Frist empfangen"
        )
        received = next(stream)
        if not (_TABLE == received.table and _SENTINEL in received.new_image.decode()):
            received = None

    assert received is not None
    new_image = json.loads(received.new_image.decode())
    assert _TARGET_KEY in new_image, f"Zielschlüssel fehlt: {new_image!r}"
    assert new_image[_TARGET_KEY] == _SENTINEL
    assert _SOURCE_KEY not in new_image, f"Quellschlüssel noch vorhanden: {new_image!r}"

    print(
        f"RECEIVED change_id={received.change_id} table={received.table} "
        f"operation={received.operation} new_image={received.new_image.decode()}",
        flush=True,
    )
