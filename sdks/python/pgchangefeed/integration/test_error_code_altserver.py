"""Real-server integration test against a server older than the message codes.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-altserver-tests.sh inside the same Docker network as
the feed container. The same calls as test_error_code_realserver.py end with
the same typed errors, but the server sends no ``code`` in the HTTP error body,
no ``google.rpc.ErrorInfo`` in the gRPC status detail and no ``error_code`` in
the diagnose heartbeat: every ``message_code`` stays ``None``, the error text
is unchanged and the SDK does not fail. The error state of the diagnose report
is written by the runner without an ``error_code``. The printed markers
(``READY``/``NORMAL_DONE``/``RECEIVED``) are what the runner asserts on.
"""

from __future__ import annotations

import os
import time

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
_RAW_HTTP_TEXT = os.environ["PGCHANGEFEED_ALTSERVER_HTTP_TEXT"]

_MISSING_TABLE = "sdk_error_code_missing_table"


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


def test_an_older_server_carries_no_message_code_and_the_sdk_stays_intact() -> None:
    print("READY", flush=True)

    try:
        _http_client(_ADMIN_TOKEN).enable_table(_http_enable_request())
    except PgChangeFeedNotFoundError as error:
        http_error = error
    else:
        raise AssertionError("enable_table on a missing table did not raise over HTTP")

    admin = _administration_client(_ADMIN_TOKEN)
    try:
        admin.enable_table(_grpc_enable_request())
    except PgChangeFeedGrpcNotFoundError as error:
        grpc_error = error
    else:
        raise AssertionError("enable_table on a missing table did not raise over gRPC")

    assert http_error.status_code == 404
    assert grpc_error.code == grpc.StatusCode.NOT_FOUND
    assert http_error.message_code is None
    assert grpc_error.message_code is None
    assert str(http_error) != ""
    assert str(grpc_error) != ""
    assert str(http_error) == _RAW_HTTP_TEXT

    try:
        _http_client(_READER_TOKEN).enable_table(_http_enable_request())
    except PgChangeFeedForbiddenError as error:
        http_forbidden = error
    else:
        raise AssertionError("enable_table with a reader token did not raise over HTTP")

    try:
        _administration_client(_READER_TOKEN).enable_table(_grpc_enable_request())
    except PgChangeFeedGrpcPermissionDeniedError as error:
        grpc_forbidden = error
    else:
        raise AssertionError("enable_table with a reader token did not raise over gRPC")

    assert http_forbidden.status_code == 403
    assert grpc_forbidden.code == grpc.StatusCode.PERMISSION_DENIED
    assert http_forbidden.message_code is None
    assert grpc_forbidden.message_code is None

    normal = admin.diagnose(administration_pb2.DiagnoseRequest(source=_SOURCE))
    assert normal.heartbeat.known
    assert normal.heartbeat.error_class == ""
    assert normal.heartbeat.error_code == ""
    print("NORMAL_DONE", flush=True)

    # The runner writes the error state (error_class without error_code) after
    # NORMAL_DONE; the periodic heartbeat resets it, so the call is repeated
    # until the report shows it.
    failing = None
    deadline = time.monotonic() + 60
    while time.monotonic() < deadline:
        report = admin.diagnose(administration_pb2.DiagnoseRequest(source=_SOURCE))
        if report.heartbeat.error_class == "schema":
            failing = report.heartbeat
            break
        time.sleep(0.2)

    assert failing is not None, "no error state in the diagnose report within 60 s"
    assert failing.error_code == ""
    print(
        f"RECEIVED code=none http=404 grpc=NOT_FOUND text={len(str(http_error))} diag_error_code=leer",
        flush=True,
    )
