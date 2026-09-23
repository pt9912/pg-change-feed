"""Network-free unit tests for the gRPC stream client surface (SPEC-020).

The fake channel mirrors the shape of the real grpcio plumbing at exactly
the seam the client touches: the generated stub asks the channel for a
``unary_stream`` multi-callable (one per RPC path) and invokes it with the
request and call options; the response iterator surfaces errors from
iteration, never at invocation time -- the same observable behavior a real
``grpc`` channel carries, without a socket. Message completeness is asserted
against the generated stub's descriptor, so a drift between this package
and ``proto/cdc/stream/v1/changestream.proto`` (SPEC-020) breaks the test at
build time instead of on a consumer's wire.
"""

from __future__ import annotations

from typing import Any, Callable, Iterator

import grpc
import pytest

from pgchangefeed.grpc_client import PgChangeFeedGrpcClient
from pgchangefeed.grpc_gen import changestream_pb2, changestream_pb2_grpc
from pgchangefeed.options import ClientOptions

ADDRESS = "pg-change-feed:9090"
TOKEN = "e2e-reader-token"

SPEC_020_FIELD_NAMES = (
    "change_id",
    "transaction_id",
    "source_table_id",
    "sequence",
    "operation",
    "old_image",
    "new_image",
    "schema_version",
    "schema",
    "table",
)


class _FakeChannel:
    """Fake ``grpc.Channel``: records the requested RPC method and every
    invocation (request, metadata) and hands back the result the test built.
    Only ``unary_stream`` is implemented -- the one call shape ``ChangeStream``
    (SPEC-020) uses; every other call shape would be a protocol violation the
    fake refuses rather than fakes."""

    def __init__(self, result_factory: Callable[[], Iterator[Any]]) -> None:
        self.method: str | None = None
        # je Aufruf: (metadata, request, timeout) — der Timeout wird mit
        # gezeichnet, damit eine Regression, die ihn still fallen laesst,
        # sichtbar rot macht statt am Fake vorüberzulaufen.
        self.invocations: list[tuple[tuple[tuple[str, str], ...], Any, float | None]] = []
        self._result_factory = result_factory

    def unary_stream(
        self,
        method: str,
        request_serializer: Any = None,
        response_deserializer: Any = None,
        **kwargs: Any,
    ) -> Callable[..., Iterator[Any]]:
        self.method = method

        def call(
            request: Any,
            metadata: tuple[tuple[str, str], ...] | None = None,
            timeout: float | None = None,
        ) -> Iterator[Any]:
            self.invocations.append((tuple(metadata or ()), request, timeout))
            return self._result_factory()

        return call


class _StreamRpcError(grpc.RpcError):
    """Carries one gRPC status code, the way a real channel's response
    iterator raises it from ``__next__`` when the call is rejected or dies
    (the auth boundary among them)."""

    def __init__(self, code: grpc.StatusCode) -> None:
        self._code = code

    def code(self) -> grpc.StatusCode:
        return self._code


class _ErrorRaisingIterator:
    def __init__(self, error: grpc.RpcError) -> None:
        self._error = error

    def __iter__(self) -> "_ErrorRaisingIterator":
        return self

    def __next__(self) -> Any:
        raise self._error


def _make_client(
    messages: list[changestream_pb2.Change] | None = None,
    error: grpc.RpcError | None = None,
) -> tuple[PgChangeFeedGrpcClient, _FakeChannel]:
    def result_factory() -> Iterator[Any]:
        if error is not None:
            return _ErrorRaisingIterator(error)
        return iter(messages or [])

    channel = _FakeChannel(result_factory)
    options = ClientOptions(address=ADDRESS, api_token=TOKEN)
    return PgChangeFeedGrpcClient(channel, options), channel


# --- Message schema completeness (SPEC-020: ten fields, wire order) ---


def test_generated_change_message_carries_the_ten_spec_020_fields_in_order() -> None:
    fields = changestream_pb2.Change.DESCRIPTOR.fields
    assert tuple(field.name for field in fields) == SPEC_020_FIELD_NAMES
    assert tuple(field.number for field in fields) == tuple(range(1, 11))


def test_stub_wires_the_spec_020_rpc_path() -> None:
    client, channel = _make_client()
    list(client.stream_changes())
    assert channel.method == "/cdc.stream.v1.ChangeStream/StreamChanges"


# --- Happy path (SPEC-020: one message per row change, in order, unmapped) ---


def test_stream_changes_yields_the_generated_change_messages_unmapped() -> None:
    first = changestream_pb2.Change(
        change_id="chg-1",
        transaction_id="tx-1",
        source_table_id="tbl-1",
        sequence=1,
        operation="INSERT",
        new_image=b'{"id": 1, "name": "first"}',
        schema_version="sv-1",
        schema="public",
        table="orders",
    )
    second = changestream_pb2.Change(
        change_id="chg-2",
        transaction_id="tx-1",
        source_table_id="tbl-1",
        sequence=2,
        operation="UPDATE",
        old_image=b'{"id": 1, "name": "first"}',
        new_image=b'{"id": 1, "name": "second"}',
        schema_version="sv-1",
        schema="public",
        table="orders",
    )
    client, _ = _make_client(messages=[first, second])
    yielded = list(client.stream_changes())
    assert yielded == [first, second]
    assert type(yielded[0]) is changestream_pb2.Change


# --- Call form (SPEC-020: empty request, Bearer-token metadata) ---


def test_stream_changes_sends_the_bearer_token_in_authorization_metadata() -> None:
    client, channel = _make_client()
    list(client.stream_changes())
    assert len(channel.invocations) == 1
    metadata, _request, _timeout = channel.invocations[0]
    assert metadata == (("authorization", f"Bearer {TOKEN}"),)


def test_stream_changes_forwards_the_call_timeout() -> None:
    client, channel = _make_client()
    list(client.stream_changes(timeout=30.0))
    _metadata, _request, timeout = channel.invocations[0]
    assert timeout == 30.0


def test_stream_changes_leaves_the_call_deadline_unbounded_by_default() -> None:
    client, channel = _make_client()
    list(client.stream_changes())
    _metadata, _request, timeout = channel.invocations[0]
    assert timeout is None


def test_stream_changes_sends_a_filterless_request() -> None:
    client, channel = _make_client()
    list(client.stream_changes())
    _metadata, request, _timeout = channel.invocations[0]
    assert isinstance(request, changestream_pb2.StreamChangesRequest)
    assert request.SerializeToString() == b""


# --- Auth boundary (SPEC-020 Negative: Unauthenticated, not a swallowed
# --- empty stream) ---


def test_unauthenticated_surfaces_from_the_enumeration() -> None:
    client, _channel = _make_client(error=_StreamRpcError(grpc.StatusCode.UNAUTHENTICATED))
    with pytest.raises(grpc.RpcError) as exc_info:
        list(client.stream_changes())
    assert exc_info.value.code() == grpc.StatusCode.UNAUTHENTICATED