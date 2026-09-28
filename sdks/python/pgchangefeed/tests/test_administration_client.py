"""Network-free unit tests for the gRPC Administration client.

The fake channel mirrors the shape of the real grpcio plumbing at exactly the
seam the client touches: the generated stub asks the channel for a
``unary_unary`` multi-callable (one per RPC path) and invokes it with the
request and metadata; a rejected call raises ``grpc.RpcError`` synchronously
from that invocation -- the same observable behavior a real ``grpc`` channel
carries, without a socket.
"""

from __future__ import annotations

from typing import Any, Callable

import grpc
import pytest

from pgchangefeed.administration_client import PgChangeFeedAdministrationClient
from pgchangefeed.exceptions import (
    PgChangeFeedGrpcInternalError,
    PgChangeFeedGrpcInvalidArgumentError,
    PgChangeFeedGrpcNotFoundError,
    PgChangeFeedGrpcPermissionDeniedError,
    PgChangeFeedGrpcUnauthenticatedError,
    PgChangeFeedGrpcUnexpectedStatusError,
)
from pgchangefeed.grpc_gen import administration_pb2
from pgchangefeed.options import ClientOptions

ADDRESS = "pg-change-feed:9090"
TOKEN = "e2e-admin-token"


class _UnaryRpcError(grpc.RpcError):
    """Carries one gRPC status code and detail text, the way a real
    ``grpc._channel._InactiveRpcError`` raises synchronously from a rejected
    unary call."""

    def __init__(self, code: grpc.StatusCode, details: str = "boom") -> None:
        self._code = code
        self._details = details

    def code(self) -> grpc.StatusCode:
        return self._code

    def details(self) -> str:
        return self._details


class _FakeChannel:
    """Fake ``grpc.Channel``: the stub asks it for a ``unary_unary``
    multi-callable once per RPC path at construction time -- this records
    every invocation as ``(method, metadata, request)``, keyed by the exact
    path it was invoked through, and hands back the canned result or raises
    the canned error. Only ``unary_unary`` is implemented -- the one call
    shape every ``Administration`` RPC uses."""

    def __init__(self, result: Any = None, error: grpc.RpcError | None = None) -> None:
        self.invocations: list[tuple[str, tuple[tuple[str, str], ...], Any]] = []
        self._result = result
        self._error = error

    def unary_unary(
        self,
        method: str,
        request_serializer: Any = None,
        response_deserializer: Any = None,
        **kwargs: Any,
    ) -> Callable[..., Any]:
        def call(request: Any, metadata: tuple[tuple[str, str], ...] | None = None, timeout: float | None = None) -> Any:
            self.invocations.append((method, tuple(metadata or ()), request))
            if self._error is not None:
                raise self._error
            return self._result

        return call


def _make_client(result: Any = None, error: grpc.RpcError | None = None) -> tuple[PgChangeFeedAdministrationClient, _FakeChannel]:
    channel = _FakeChannel(result=result, error=error)
    options = ClientOptions(address=ADDRESS, api_token=TOKEN)
    return PgChangeFeedAdministrationClient(channel, options), channel


def _assert_called(channel: _FakeChannel, method: str) -> None:
    assert len(channel.invocations) == 1
    called_method, metadata, _request = channel.invocations[0]
    assert called_method == method
    assert metadata == (("authorization", f"Bearer {TOKEN}"),)


# --- register_consumer ---


def test_register_consumer_happy_path_returns_the_response_unmapped() -> None:
    response = administration_pb2.RegisterConsumerResponse(consumer_id="billing", name="Billing", already_registered=False)
    client, channel = _make_client(result=response)
    result = client.register_consumer(administration_pb2.RegisterConsumerRequest(consumer_id="billing", name="Billing"))
    assert result is response
    _assert_called(channel, "/cdc.administration.v1.Administration/RegisterConsumer")


def test_register_consumer_boundary_already_registered_is_true_for_an_existing_consumer() -> None:
    response = administration_pb2.RegisterConsumerResponse(consumer_id="billing", name="Billing", already_registered=True)
    client, _channel = _make_client(result=response)
    result = client.register_consumer(administration_pb2.RegisterConsumerRequest(consumer_id="billing", name="Billing"))
    assert result.already_registered is True


def test_register_consumer_negative_permission_denied_raises_typed_error() -> None:
    client, _channel = _make_client(error=_UnaryRpcError(grpc.StatusCode.PERMISSION_DENIED))
    with pytest.raises(PgChangeFeedGrpcPermissionDeniedError):
        client.register_consumer(administration_pb2.RegisterConsumerRequest(consumer_id="billing", name="Billing"))


# --- acknowledge_consumer ---


def test_acknowledge_consumer_happy_path_returns_the_response_unmapped() -> None:
    response = administration_pb2.AcknowledgeConsumerResponse(consumer_id="billing", source_id="my-source", offset=42)
    client, channel = _make_client(result=response)
    result = client.acknowledge_consumer(
        administration_pb2.AcknowledgeConsumerRequest(consumer_id="billing", source_id="my-source", offset=42)
    )
    assert result is response
    _assert_called(channel, "/cdc.administration.v1.Administration/AcknowledgeConsumer")


def test_acknowledge_consumer_boundary_offset_carries_a_large_uint64_unchanged() -> None:
    large_offset = 2**63
    response = administration_pb2.AcknowledgeConsumerResponse(consumer_id="billing", source_id="my-source", offset=large_offset)
    client, channel = _make_client(result=response)
    client.acknowledge_consumer(
        administration_pb2.AcknowledgeConsumerRequest(consumer_id="billing", source_id="my-source", offset=large_offset)
    )
    _method, _metadata, request = channel.invocations[0]
    assert request.offset == large_offset


def test_acknowledge_consumer_negative_invalid_argument_raises_typed_error() -> None:
    client, _channel = _make_client(error=_UnaryRpcError(grpc.StatusCode.INVALID_ARGUMENT))
    with pytest.raises(PgChangeFeedGrpcInvalidArgumentError):
        client.acknowledge_consumer(
            administration_pb2.AcknowledgeConsumerRequest(consumer_id="billing", source_id="my-source", offset=1)
        )


# --- get_consumer_position ---


def test_get_consumer_position_happy_path_returns_the_response_unmapped() -> None:
    response = administration_pb2.GetConsumerPositionResponse(
        consumer_id="billing", source_id="my-source", offset=42, acknowledged=True
    )
    client, channel = _make_client(result=response)
    result = client.get_consumer_position(administration_pb2.GetConsumerPositionRequest(consumer_id="billing"))
    assert result is response
    _assert_called(channel, "/cdc.administration.v1.Administration/GetConsumerPosition")


def test_get_consumer_position_boundary_never_acknowledged_reports_the_starting_position() -> None:
    response = administration_pb2.GetConsumerPositionResponse(
        consumer_id="billing", source_id="my-source", offset=0, acknowledged=False
    )
    client, _channel = _make_client(result=response)
    result = client.get_consumer_position(administration_pb2.GetConsumerPositionRequest(consumer_id="billing"))
    assert result.acknowledged is False
    assert result.offset == 0


def test_get_consumer_position_negative_unauthenticated_raises_typed_error() -> None:
    client, _channel = _make_client(error=_UnaryRpcError(grpc.StatusCode.UNAUTHENTICATED))
    with pytest.raises(PgChangeFeedGrpcUnauthenticatedError):
        client.get_consumer_position(administration_pb2.GetConsumerPositionRequest(consumer_id="billing"))


# --- remove_consumer ---


def test_remove_consumer_happy_path_returns_the_response_unmapped() -> None:
    response = administration_pb2.RemoveConsumerResponse(consumer_id="billing", removed=True)
    client, channel = _make_client(result=response)
    result = client.remove_consumer(administration_pb2.RemoveConsumerRequest(consumer_id="billing"))
    assert result is response
    _assert_called(channel, "/cdc.administration.v1.Administration/RemoveConsumer")


def test_remove_consumer_boundary_removed_is_false_for_one_never_registered() -> None:
    response = administration_pb2.RemoveConsumerResponse(consumer_id="ghost", removed=False)
    client, _channel = _make_client(result=response)
    result = client.remove_consumer(administration_pb2.RemoveConsumerRequest(consumer_id="ghost"))
    assert result.removed is False


def test_remove_consumer_negative_permission_denied_raises_typed_error() -> None:
    client, _channel = _make_client(error=_UnaryRpcError(grpc.StatusCode.PERMISSION_DENIED))
    with pytest.raises(PgChangeFeedGrpcPermissionDeniedError):
        client.remove_consumer(administration_pb2.RemoveConsumerRequest(consumer_id="billing"))


# --- enable_table ---


def _enable_table_request() -> administration_pb2.EnableTableRequest:
    return administration_pb2.EnableTableRequest(
        source="my-source",
        schema="public",
        table="orders",
        table_id="orders",
        schema_version_id="v1",
        version=1,
        publication="my_publication",
    )


def test_enable_table_happy_path_returns_the_response_unmapped() -> None:
    response = administration_pb2.EnableTableResponse(
        table_id="orders", source="my-source", schema="public", table="orders", already_enabled=False
    )
    client, channel = _make_client(result=response)
    result = client.enable_table(_enable_table_request())
    assert result is response
    _assert_called(channel, "/cdc.administration.v1.Administration/EnableTable")


def test_enable_table_boundary_already_enabled_is_true_for_a_captured_table() -> None:
    response = administration_pb2.EnableTableResponse(
        table_id="orders", source="my-source", schema="public", table="orders", already_enabled=True
    )
    client, _channel = _make_client(result=response)
    result = client.enable_table(_enable_table_request())
    assert result.already_enabled is True


def test_enable_table_negative_not_found_raises_typed_error_for_a_missing_table() -> None:
    client, _channel = _make_client(error=_UnaryRpcError(grpc.StatusCode.NOT_FOUND))
    with pytest.raises(PgChangeFeedGrpcNotFoundError):
        client.enable_table(_enable_table_request())


# --- disable_table ---


def _disable_table_request() -> administration_pb2.DisableTableRequest:
    return administration_pb2.DisableTableRequest(source="my-source", schema="public", table="orders", publication="my_publication")


def test_disable_table_happy_path_returns_the_response_unmapped() -> None:
    response = administration_pb2.DisableTableResponse(removed=True, retained=True)
    client, channel = _make_client(result=response)
    result = client.disable_table(_disable_table_request())
    assert result is response
    _assert_called(channel, "/cdc.administration.v1.Administration/DisableTable")


def test_disable_table_boundary_retained_is_false_when_nothing_remains() -> None:
    response = administration_pb2.DisableTableResponse(removed=True, retained=False)
    client, _channel = _make_client(result=response)
    result = client.disable_table(_disable_table_request())
    assert result.retained is False


def test_disable_table_negative_not_found_raises_typed_error_for_a_missing_table() -> None:
    client, _channel = _make_client(error=_UnaryRpcError(grpc.StatusCode.NOT_FOUND))
    with pytest.raises(PgChangeFeedGrpcNotFoundError):
        client.disable_table(_disable_table_request())


# --- get_table_status ---


def _get_table_status_request() -> administration_pb2.GetTableStatusRequest:
    return administration_pb2.GetTableStatusRequest(source="my-source", schema="public", table="orders", publication="my_publication")


def test_get_table_status_happy_path_returns_the_response_unmapped() -> None:
    response = administration_pb2.GetTableStatusResponse(enabled=True, retained=False)
    client, channel = _make_client(result=response)
    result = client.get_table_status(_get_table_status_request())
    assert result is response
    _assert_called(channel, "/cdc.administration.v1.Administration/GetTableStatus")


def test_get_table_status_boundary_never_enabled_reports_both_false() -> None:
    response = administration_pb2.GetTableStatusResponse(enabled=False, retained=False)
    client, _channel = _make_client(result=response)
    result = client.get_table_status(_get_table_status_request())
    assert result.enabled is False
    assert result.retained is False


def test_get_table_status_negative_not_found_raises_typed_error_for_a_missing_table() -> None:
    client, _channel = _make_client(error=_UnaryRpcError(grpc.StatusCode.NOT_FOUND))
    with pytest.raises(PgChangeFeedGrpcNotFoundError):
        client.get_table_status(_get_table_status_request())


# --- list_tables ---


def test_list_tables_happy_path_returns_the_response_unmapped() -> None:
    table = administration_pb2.SourceTable(table_id="orders", source="my-source", schema="public", table="orders")
    response = administration_pb2.ListTablesResponse(tables=[table], retained=[])
    client, channel = _make_client(result=response)
    result = client.list_tables(administration_pb2.ListTablesRequest(source="my-source", publication="my_publication"))
    assert result is response
    _assert_called(channel, "/cdc.administration.v1.Administration/ListTables")


def test_list_tables_boundary_empty_tables_and_retained() -> None:
    response = administration_pb2.ListTablesResponse(tables=[], retained=[])
    client, _channel = _make_client(result=response)
    result = client.list_tables(administration_pb2.ListTablesRequest(source="my-source", publication="my_publication"))
    assert list(result.tables) == []
    assert list(result.retained) == []


def test_list_tables_negative_unauthenticated_raises_typed_error() -> None:
    client, _channel = _make_client(error=_UnaryRpcError(grpc.StatusCode.UNAUTHENTICATED))
    with pytest.raises(PgChangeFeedGrpcUnauthenticatedError):
        client.list_tables(administration_pb2.ListTablesRequest(source="my-source", publication="my_publication"))


# --- run_retention ---


def test_run_retention_happy_path_returns_the_response_unmapped() -> None:
    response = administration_pb2.RunRetentionResponse(deleted=7)
    client, channel = _make_client(result=response)
    result = client.run_retention(administration_pb2.RunRetentionRequest(source="my-source", min_age_nanos=1_000_000_000))
    assert result is response
    _assert_called(channel, "/cdc.administration.v1.Administration/RunRetention")


def test_run_retention_boundary_zero_min_age_nanos_is_valid_and_means_no_minimum() -> None:
    response = administration_pb2.RunRetentionResponse(deleted=0)
    client, channel = _make_client(result=response)
    result = client.run_retention(administration_pb2.RunRetentionRequest(source="my-source", min_age_nanos=0))
    _method, _metadata, request = channel.invocations[0]
    assert request.min_age_nanos == 0
    assert result.deleted == 0


def test_run_retention_negative_permission_denied_raises_typed_error() -> None:
    client, _channel = _make_client(error=_UnaryRpcError(grpc.StatusCode.PERMISSION_DENIED))
    with pytest.raises(PgChangeFeedGrpcPermissionDeniedError):
        client.run_retention(administration_pb2.RunRetentionRequest(source="my-source", min_age_nanos=0))


# --- read_changes ---


def test_read_changes_happy_path_returns_the_response_unmapped() -> None:
    record = administration_pb2.ChangeRecord(
        commit_position=1,
        change_id="chg-1",
        transaction_id="tx-1",
        source_table_id="tbl-1",
        schema="public",
        table="orders",
        sequence=1,
        operation="INSERT",
        new_image=b'{"id": 1}',
        schema_version="sv-1",
        committed_at="2026-09-28T00:00:00Z",
        origin="wal",
    )
    response = administration_pb2.ReadChangesResponse(changes=[record])
    client, channel = _make_client(result=response)
    result = client.read_changes(administration_pb2.ReadChangesRequest(source="my-source"))
    assert result is response
    _assert_called(channel, "/cdc.administration.v1.Administration/ReadChanges")


def test_read_changes_boundary_unset_range_fields_carry_zero() -> None:
    # `from` is a Python keyword and unusable in dotted attribute syntax
    # (`request.from`); the generated message still carries the field under
    # that exact proto name, reachable through `getattr`.
    response = administration_pb2.ReadChangesResponse(changes=[])
    client, channel = _make_client(result=response)
    client.read_changes(administration_pb2.ReadChangesRequest(source="my-source"))
    _method, _metadata, request = channel.invocations[0]
    assert getattr(request, "from") == 0
    assert request.to == 0
    assert request.limit == 0


def test_read_changes_negative_invalid_argument_raises_typed_error() -> None:
    client, _channel = _make_client(error=_UnaryRpcError(grpc.StatusCode.INVALID_ARGUMENT))
    with pytest.raises(PgChangeFeedGrpcInvalidArgumentError):
        client.read_changes(administration_pb2.ReadChangesRequest(source="my-source"))


# --- diagnose ---


def test_diagnose_happy_path_returns_the_response_unmapped() -> None:
    response = administration_pb2.DiagnoseResponse(
        heartbeat=administration_pb2.HeartbeatStatus(known=True, age_seconds=1.5, error_class=""),
        capture_lag=0.5,
        consumer_lags=[administration_pb2.ConsumerLag(consumer_id="billing", known=True, lag=0.1)],
        retention_blocker=administration_pb2.RetentionBlocker(present=False),
        storage_bytes=1024.0,
        backfill=[],
    )
    client, channel = _make_client(result=response)
    result = client.diagnose(administration_pb2.DiagnoseRequest(source="my-source"))
    assert result is response
    _assert_called(channel, "/cdc.administration.v1.Administration/Diagnose")


def test_diagnose_boundary_absence_cases_are_false_and_unknown_carries_no_meaning() -> None:
    response = administration_pb2.DiagnoseResponse(
        heartbeat=administration_pb2.HeartbeatStatus(known=False),
        capture_lag=0.0,
        consumer_lags=[],
        retention_blocker=administration_pb2.RetentionBlocker(present=False),
        storage_bytes=0.0,
        backfill=[],
    )
    client, _channel = _make_client(result=response)
    result = client.diagnose(administration_pb2.DiagnoseRequest(source="my-source"))
    assert result.heartbeat.known is False
    assert result.retention_blocker.present is False
    assert list(result.consumer_lags) == []
    assert list(result.backfill) == []


def test_diagnose_negative_internal_raises_typed_error() -> None:
    client, _channel = _make_client(error=_UnaryRpcError(grpc.StatusCode.INTERNAL))
    with pytest.raises(PgChangeFeedGrpcInternalError):
        client.diagnose(administration_pb2.DiagnoseRequest(source="my-source"))


# --- Cross-cutting: every gRPC status maps to its typed error, the original
# --- error is always the cause ---


@pytest.mark.parametrize(
    ("code", "expected"),
    [
        (grpc.StatusCode.INVALID_ARGUMENT, PgChangeFeedGrpcInvalidArgumentError),
        (grpc.StatusCode.UNAUTHENTICATED, PgChangeFeedGrpcUnauthenticatedError),
        (grpc.StatusCode.PERMISSION_DENIED, PgChangeFeedGrpcPermissionDeniedError),
        (grpc.StatusCode.NOT_FOUND, PgChangeFeedGrpcNotFoundError),
        (grpc.StatusCode.INTERNAL, PgChangeFeedGrpcInternalError),
        (grpc.StatusCode.UNAVAILABLE, PgChangeFeedGrpcUnexpectedStatusError),
    ],
)
def test_every_grpc_status_maps_to_its_typed_error(code: grpc.StatusCode, expected: type) -> None:
    error = _UnaryRpcError(code, details="boom")
    client, _channel = _make_client(error=error)
    with pytest.raises(expected) as exc_info:
        client.register_consumer(administration_pb2.RegisterConsumerRequest(consumer_id="billing", name="Billing"))
    assert exc_info.value.code == code
    assert str(exc_info.value) == "boom"
    assert exc_info.value.__cause__ is error
