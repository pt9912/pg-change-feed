"""Realserver integration test for the gRPC stream client surface (SPEC-020).

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container (ADR-0110 §Entscheidung Festlegung
2/Folgepflicht 1: every new delivery-path surface of this package proves its
protocol assumptions against a real, running server instance, not only
against the specification read). The runner drives the environment: it brings
the compose stack up, activates the tables and inserts rows while this test
is streaming; the printed markers (``READY``/``RECEIVED``/``REJECTED``) are
what the runner asserts on, so every ``print`` is flushed and unbuffered.

The commit-to-delivery window is fire-and-forget (SPEC-020): between the
stream opening and the server registering the subscriber at the Broadcaster,
a committed change can be discarded for this subscriber. The runner commits a
bounded sequence of unique rows until exactly one arrives; this test keeps
receiving until it sees its sentinel, so a missed row costs nothing.
"""

from __future__ import annotations

import os
import time
from typing import Iterator

import grpc
import pytest

from pgchangefeed.grpc_client import PgChangeFeedGrpcClient
from pgchangefeed.grpc_gen import changestream_pb2, changestream_pb2_grpc
from pgchangefeed.options import ClientOptions

_ADDR = os.environ["PGCHANGEFEED_GRPC_ADDR"]
_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN"]
_TABLE = os.environ["PGCHANGEFEED_E2E_TABLE"]
_SENTINEL = os.environ["PGCHANGEFEED_E2E_SENTINEL"]

_RECEIVE_DEADLINE_SECONDS = 90.0


def test_realserver_receives_a_committed_change_over_the_stream() -> None:
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
    # Feldvollstaendigkeit am realen Wire-Image (SPEC-020): alle zehn Felder
    # getragen, je in seinem Wire-Typ; getypte Pruefung, weil ein `!= ""`
    # gegen bytes/int-Felder vakuum waere (nie ungleich-rot).
    for field, kind in (
        ("change_id", str),
        ("transaction_id", str),
        ("source_table_id", str),
        ("sequence", int),
        ("operation", str),
        ("old_image", bytes),
        ("new_image", bytes),
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
    assert received.old_image == b"", "ein INSERT traegt kein Alt-Bild am Wire"
    assert _SENTINEL in received.new_image.decode()
    assert received.table == _TABLE
    assert received.schema == "public"

    print(
        f"RECEIVED change_id={received.change_id} table={received.table} "
        f"operation={received.operation} new_image={received.new_image.decode()}",
        flush=True,
    )


def test_realserver_rejects_the_stream_without_token() -> None:
    channel = grpc.insecure_channel(_ADDR)
    stub = changestream_pb2_grpc.ChangeStreamStub(channel)
    opened: Iterator[changestream_pb2.Change] = stub.StreamChanges(
        changestream_pb2.StreamChangesRequest()
    )
    with pytest.raises(grpc.RpcError) as exc_info:
        next(opened)
    assert exc_info.value.code() == grpc.StatusCode.UNAUTHENTICATED
    print("REJECTED code=Unauthenticated", flush=True)