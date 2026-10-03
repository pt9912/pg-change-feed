"""Network-free tests for the message code on the typed errors.

One input table per wire form, the same rows in every SDK language: the HTTP
error body (HTTP client and SSE client) and the ``grpc-status-details-bin``
trailer of the administration client. The status detail of the gRPC tests is
serialized with the protobuf runtime from independently declared message
descriptors, not with the reader under test, so a wrong field number in the
reader cannot cancel out against the same mistake in the test.
"""

from __future__ import annotations

from typing import Any

import grpc
import httpx
import pytest
from google.protobuf import any_pb2, descriptor_pb2, descriptor_pool, message_factory

from pgchangefeed.administration_client import PgChangeFeedAdministrationClient
from pgchangefeed.exceptions import (
    PgChangeFeedBadRequestError,
    PgChangeFeedError,
    PgChangeFeedForbiddenError,
    PgChangeFeedGrpcError,
    PgChangeFeedGrpcInternalError,
    PgChangeFeedGrpcInvalidArgumentError,
    PgChangeFeedGrpcNotFoundError,
    PgChangeFeedGrpcPermissionDeniedError,
    PgChangeFeedGrpcUnauthenticatedError,
    PgChangeFeedGrpcUnexpectedStatusError,
    PgChangeFeedNotFoundError,
    PgChangeFeedServerError,
    PgChangeFeedUnauthorizedError,
    PgChangeFeedUnexpectedStatusError,
)
from pgchangefeed.grpc_gen import administration_pb2
from pgchangefeed.http_client import PgChangeFeedHttpClient
from pgchangefeed.options import ClientOptions
from pgchangefeed.sse_client import PgChangeFeedSseClient

ADDRESS = "https://feed.example.invalid"
TOKEN = "token"

# --- HTTP error body ---

_HTTP_ROWS = [
    (400, b'{"error":"x","code":"PCF-E8051"}', PgChangeFeedBadRequestError, "PCF-E8051", "x"),
    (400, b'{"error":"x"}', PgChangeFeedBadRequestError, None, "x"),
    (400, b'{"error":"x","code":""}', PgChangeFeedBadRequestError, None, "x"),
    (400, b'{"error":"x","code":null}', PgChangeFeedBadRequestError, None, "x"),
    (400, b'{"error":"x","code":5}', PgChangeFeedBadRequestError, None, "x"),
    (404, b'{"error":"x","code":"PCF-E8025"}', PgChangeFeedNotFoundError, "PCF-E8025", "x"),
    (500, b'{"error":"x","code":"PCF-E7000"}', PgChangeFeedServerError, "PCF-E7000", "x"),
    (503, b'{"error":"x","code":"PCF-E2001"}', PgChangeFeedUnexpectedStatusError, "PCF-E2001", "x"),
    (401, b'{"error":"x"}', PgChangeFeedUnauthorizedError, None, "x"),
    (403, b'{"error":"x"}', PgChangeFeedForbiddenError, None, "x"),
    (400, b'{"error":"x","code":"NO-PCF","extra":1}', PgChangeFeedBadRequestError, "NO-PCF", "x"),
    (502, b"Bad Gateway", PgChangeFeedUnexpectedStatusError, None, "Bad Gateway"),
]


def _http_client(status: int, content: bytes) -> PgChangeFeedHttpClient:
    transport = httpx.MockTransport(lambda request: httpx.Response(status, content=content))
    return PgChangeFeedHttpClient(httpx.Client(transport=transport), ClientOptions(address=ADDRESS, api_token=TOKEN))


def _sse_client(status: int, content: bytes) -> PgChangeFeedSseClient:
    transport = httpx.MockTransport(lambda request: httpx.Response(status, content=content))
    return PgChangeFeedSseClient(httpx.Client(transport=transport), ClientOptions(address=ADDRESS, api_token=TOKEN))


@pytest.mark.parametrize(("status", "content", "expected", "code", "text"), _HTTP_ROWS)
def test_http_client_error_carries_the_message_code(status: int, content: bytes, expected: type, code: str | None, text: str) -> None:
    with pytest.raises(expected) as exc_info:
        _http_client(status, content).list_tables(source="src", publication="pub")
    assert type(exc_info.value) is expected
    assert exc_info.value.message_code == code
    assert str(exc_info.value) == text


@pytest.mark.parametrize(("status", "content", "expected", "code", "text"), _HTTP_ROWS)
def test_sse_client_error_carries_the_message_code(status: int, content: bytes, expected: type, code: str | None, text: str) -> None:
    with pytest.raises(expected) as exc_info:
        list(_sse_client(status, content).stream_changes())
    assert type(exc_info.value) is expected
    assert exc_info.value.message_code == code
    assert str(exc_info.value) == text


def test_http_error_constructed_without_message_code_keeps_working_and_has_none() -> None:
    error = PgChangeFeedError(400, "x")
    assert error.status_code == 400
    assert str(error) == "x"
    assert error.message_code is None


# --- gRPC status detail ---


def _declare_classes() -> tuple[Any, Any]:
    pool = descriptor_pool.DescriptorPool()
    any_file = descriptor_pb2.FileDescriptorProto()
    any_pb2.DESCRIPTOR.CopyToProto(any_file)
    pool.AddSerializedFile(any_file.SerializeToString())

    file = descriptor_pb2.FileDescriptorProto(
        name="test_rpc.proto",
        package="google.rpc",
        syntax="proto3",
        dependency=["google/protobuf/any.proto"],
    )
    field = descriptor_pb2.FieldDescriptorProto
    error_info = file.message_type.add(name="ErrorInfo")
    error_info.field.add(name="reason", number=1, type=field.TYPE_STRING, label=field.LABEL_OPTIONAL)
    error_info.field.add(name="domain", number=2, type=field.TYPE_STRING, label=field.LABEL_OPTIONAL)
    status = file.message_type.add(name="Status")
    status.field.add(name="code", number=1, type=field.TYPE_INT32, label=field.LABEL_OPTIONAL)
    status.field.add(name="message", number=2, type=field.TYPE_STRING, label=field.LABEL_OPTIONAL)
    status.field.add(
        name="details", number=3, type=field.TYPE_MESSAGE, type_name=".google.protobuf.Any", label=field.LABEL_REPEATED
    )
    pool.AddSerializedFile(file.SerializeToString())
    return (
        message_factory.GetMessageClass(pool.FindMessageTypeByName("google.rpc.Status")),
        message_factory.GetMessageClass(pool.FindMessageTypeByName("google.rpc.ErrorInfo")),
    )


_STATUS, _ERROR_INFO = _declare_classes()
_ERROR_INFO_URL = "type.googleapis.com/google.rpc.ErrorInfo"
_RETRY_INFO_URL = "type.googleapis.com/google.rpc.RetryInfo"
_SERVER_DOMAIN = "pg-change-feed"


def _detail(url: str, reason: str, domain: str) -> tuple[str, bytes]:
    return url, _ERROR_INFO(reason=reason, domain=domain).SerializeToString()


def _status(*details: tuple[str, bytes]) -> bytes:
    status = _STATUS()
    for url, value in details:
        entry = status.details.add()
        entry.type_url = url
        entry.value = value
    return status.SerializeToString()


def _trailers(status: bytes) -> tuple[tuple[str, bytes], ...]:
    return (("grpc-status-details-bin", status),)


class _RpcError(grpc.RpcError):
    def __init__(
        self,
        code: grpc.StatusCode,
        trailers: tuple[tuple[str, str | bytes], ...] | None,
        details: str = "boom",
    ) -> None:
        self._code = code
        self._trailers = trailers
        self._details = details

    def code(self) -> grpc.StatusCode:
        return self._code

    def details(self) -> str:
        return self._details

    def trailing_metadata(self) -> tuple[tuple[str, str | bytes], ...] | None:
        return self._trailers


class _Channel:
    def __init__(self, error: grpc.RpcError | None = None, result: Any = None) -> None:
        self._error = error
        self._result = result

    def unary_unary(self, method: str, request_serializer: Any = None, response_deserializer: Any = None, **kwargs: Any) -> Any:
        def call(request: Any, metadata: Any = None, timeout: float | None = None) -> Any:
            if self._error is not None:
                raise self._error
            return self._result

        return call


def _administration_client(error: grpc.RpcError | None = None, result: Any = None) -> PgChangeFeedAdministrationClient:
    return PgChangeFeedAdministrationClient(_Channel(error, result), ClientOptions(address="feed:9090", api_token=TOKEN))


_GRPC_ROWS = [
    (grpc.StatusCode.INVALID_ARGUMENT, _trailers(_status(_detail(_ERROR_INFO_URL, "PCF-E8051", _SERVER_DOMAIN))),
     PgChangeFeedGrpcInvalidArgumentError, "PCF-E8051"),
    (grpc.StatusCode.NOT_FOUND, _trailers(_status(_detail(_ERROR_INFO_URL, "PCF-E8025", _SERVER_DOMAIN))),
     PgChangeFeedGrpcNotFoundError, "PCF-E8025"),
    (grpc.StatusCode.INTERNAL, _trailers(_status(_detail(_ERROR_INFO_URL, "PCF-E7000", _SERVER_DOMAIN))),
     PgChangeFeedGrpcInternalError, "PCF-E7000"),
    (grpc.StatusCode.INVALID_ARGUMENT, (), PgChangeFeedGrpcInvalidArgumentError, None),
    (grpc.StatusCode.INVALID_ARGUMENT, None, PgChangeFeedGrpcInvalidArgumentError, None),
    (grpc.StatusCode.UNAUTHENTICATED, (), PgChangeFeedGrpcUnauthenticatedError, None),
    (grpc.StatusCode.PERMISSION_DENIED, (), PgChangeFeedGrpcPermissionDeniedError, None),
    (grpc.StatusCode.INVALID_ARGUMENT, _trailers(_status(_detail(_ERROR_INFO_URL, "PCF-E8051", "example.com"))),
     PgChangeFeedGrpcInvalidArgumentError, None),
    (grpc.StatusCode.INVALID_ARGUMENT, _trailers(_status(_detail(_RETRY_INFO_URL, "PCF-E8051", _SERVER_DOMAIN))),
     PgChangeFeedGrpcInvalidArgumentError, None),
    (grpc.StatusCode.INVALID_ARGUMENT, _trailers(_status(_detail(_ERROR_INFO_URL, "", _SERVER_DOMAIN))),
     PgChangeFeedGrpcInvalidArgumentError, None),
    (grpc.StatusCode.INVALID_ARGUMENT, _trailers(b"\xff\xff"), PgChangeFeedGrpcInvalidArgumentError, None),
    (grpc.StatusCode.INVALID_ARGUMENT,
     _trailers(_status(_detail(_ERROR_INFO_URL, "PCF-E0001", "example.com"), _detail(_ERROR_INFO_URL, "PCF-E8051", _SERVER_DOMAIN))),
     PgChangeFeedGrpcInvalidArgumentError, "PCF-E8051"),
    (grpc.StatusCode.UNAVAILABLE, _trailers(_status(_detail(_ERROR_INFO_URL, "PCF-E8051", _SERVER_DOMAIN))),
     PgChangeFeedGrpcUnexpectedStatusError, "PCF-E8051"),
]


@pytest.mark.parametrize(("status", "trailers", "expected", "code"), _GRPC_ROWS)
def test_administration_client_error_carries_the_message_code(
    status: grpc.StatusCode, trailers: Any, expected: type, code: str | None
) -> None:
    error = _RpcError(status, trailers)
    with pytest.raises(expected) as exc_info:
        _administration_client(error=error).list_tables(administration_pb2.ListTablesRequest(source="s"))
    assert type(exc_info.value) is expected
    assert exc_info.value.message_code == code
    assert exc_info.value.code == status
    assert str(exc_info.value) == "boom"
    assert exc_info.value.__cause__ is error


def test_administration_client_error_of_a_channel_error_without_trailers_has_no_message_code() -> None:
    class _Bare(grpc.RpcError):
        def code(self) -> grpc.StatusCode:
            return grpc.StatusCode.NOT_FOUND

        def details(self) -> str:
            return "gone"

    with pytest.raises(PgChangeFeedGrpcNotFoundError) as exc_info:
        _administration_client(error=_Bare()).list_tables(administration_pb2.ListTablesRequest(source="s"))
    assert exc_info.value.message_code is None


def test_grpc_error_constructed_without_message_code_keeps_working_and_has_none() -> None:
    error = PgChangeFeedGrpcError(grpc.StatusCode.INTERNAL, "x")
    assert error.code == grpc.StatusCode.INTERNAL
    assert str(error) == "x"
    assert error.message_code is None


# --- diagnose error_code arrives unchanged ---


@pytest.mark.parametrize("error_code", ["PCF-E4003", ""])
def test_diagnose_error_code_arrives_unchanged(error_code: str) -> None:
    response = administration_pb2.DiagnoseResponse(
        heartbeat=administration_pb2.HeartbeatStatus(
            known=True, age_seconds=1.0, error_class="schema" if error_code else "", error_code=error_code
        ),
    )
    result = _administration_client(result=response).diagnose(administration_pb2.DiagnoseRequest(source="s"))
    assert result.heartbeat.error_code == error_code
