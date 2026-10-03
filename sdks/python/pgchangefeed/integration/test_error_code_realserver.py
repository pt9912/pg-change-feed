"""Real-server integration test for the message code on the typed errors.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. Enabling a table that does not exist in the
source database is answered with HTTP status 404 and gRPC status
``NOT_FOUND``; both carry the message code the server assigns to that case,
and the SDK exposes it as ``message_code`` of the typed error: from the
``code`` field of the HTTP error body, and from the ``google.rpc.ErrorInfo``
status detail on the gRPC side. A reader token against the same operation is
rejected with HTTP status 403 and gRPC status ``PERMISSION_DENIED``, which
carry no message code. The unit tests build the wire form themselves; this
test reads it from the real server. The printed markers
(``READY``/``RECEIVED``/``REJECTED``) are what the runner asserts on.
"""

from __future__ import annotations

import os

import grpc
import httpx

from pgchangefeed.administration_client import PgChangeFeedAdministrationClient
from pgchangefeed.exceptions import (
    PgChangeFeedForbiddenError,
    PgChangeFeedGrpcNotFoundError,
    PgChangeFeedGrpcPermissionDeniedError,
    PgChangeFeedNotFoundError,
)
from pgchangefeed.grpc_gen import administration_pb2
from pgchangefeed.http_client import PgChangeFeedHttpClient
from pgchangefeed.models import EnableTableRequest
from pgchangefeed.options import ClientOptions

_HTTP_ADDR = os.environ["PGCHANGEFEED_HTTP_ADDR"]
_GRPC_ADDR = os.environ["PGCHANGEFEED_GRPC_ADDR"]
_ADMIN_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN_ADMIN"]
_READER_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN_READER"]
_SOURCE = os.environ["PGCHANGEFEED_SOURCE_ID"]
_PUBLICATION = os.environ["PGCHANGEFEED_HTTP_PUBLICATION"]

_MISSING_TABLE = "sdk_error_code_missing_table"
# The server's code for enabling a table that is missing at the source
# (request: schema "public", table "sdk_error_code_missing_table").
_MISSING_TABLE_CODE = "PCF-E8025"


def _http_client(token: str) -> PgChangeFeedHttpClient:
    return PgChangeFeedHttpClient(httpx.Client(timeout=30.0), ClientOptions(address=_HTTP_ADDR, api_token=token))


def _administration_client(token: str) -> PgChangeFeedAdministrationClient:
    channel = grpc.insecure_channel(_GRPC_ADDR)
    return PgChangeFeedAdministrationClient(channel, ClientOptions(address=_GRPC_ADDR, api_token=token))


def _http_enable_request() -> EnableTableRequest:
    return EnableTableRequest(
        source=_SOURCE,
        schema="public",
        table=_MISSING_TABLE,
        table_id="sdk-error-code",
        schema_version_id="sdk-error-code-v1",
        version=1,
        publication=_PUBLICATION,
    )


def _grpc_enable_request() -> administration_pb2.EnableTableRequest:
    return administration_pb2.EnableTableRequest(
        source=_SOURCE,
        schema="public",
        table=_MISSING_TABLE,
        table_id="sdk-error-code",
        schema_version_id="sdk-error-code-v1",
        version=1,
        publication=_PUBLICATION,
    )


def test_realserver_error_carries_the_message_code() -> None:
    print("READY", flush=True)

    try:
        _http_client(_ADMIN_TOKEN).enable_table(_http_enable_request())
    except PgChangeFeedNotFoundError as error:
        http_status = error.status_code
        http_code = error.message_code
    else:
        raise AssertionError("enable_table on a missing table did not raise over HTTP")

    try:
        _administration_client(_ADMIN_TOKEN).enable_table(_grpc_enable_request())
    except PgChangeFeedGrpcNotFoundError as error:
        grpc_status = error.code
        grpc_code = error.message_code
    else:
        raise AssertionError("enable_table on a missing table did not raise over gRPC")

    assert http_status == 404
    assert grpc_status == grpc.StatusCode.NOT_FOUND
    assert http_code == _MISSING_TABLE_CODE
    assert grpc_code == _MISSING_TABLE_CODE
    print(f"RECEIVED code={http_code} http={http_status} grpc={grpc_status.name}", flush=True)


def test_realserver_rejected_call_carries_no_message_code() -> None:
    try:
        _http_client(_READER_TOKEN).enable_table(_http_enable_request())
    except PgChangeFeedForbiddenError as error:
        http_status = error.status_code
        http_code = error.message_code
    else:
        raise AssertionError("enable_table with a reader token did not raise over HTTP")

    try:
        _administration_client(_READER_TOKEN).enable_table(_grpc_enable_request())
    except PgChangeFeedGrpcPermissionDeniedError as error:
        grpc_status = error.code
        grpc_code = error.message_code
    else:
        raise AssertionError("enable_table with a reader token did not raise over gRPC")

    assert http_status == 403
    assert grpc_status == grpc.StatusCode.PERMISSION_DENIED
    assert http_code is None
    assert grpc_code is None
    print(f"REJECTED http={http_status} grpc={grpc_status.name} code=none", flush=True)
