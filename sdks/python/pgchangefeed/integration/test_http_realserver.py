"""Real-server integration test for the HTTP API client.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. A roundtrip: the SDK registers a disposable
consumer with the admin token and lists tables with the reader token; the
runner checks the registration against the SQL read path (`cdc.consumer`); a
call with an unknown token is rejected with HTTP status 401. The individual
capabilities stay with the network-free unit tests
(`tests/test_http_client.py` against `httpx.MockTransport`); this roundtrip
proves the wire assumptions (auth header form, JSON mapping, rejection
behavior) at the server.
"""

from __future__ import annotations

import os
import uuid

import httpx
import pytest

from pgchangefeed.exceptions import PgChangeFeedUnauthorizedError
from pgchangefeed.http_client import PgChangeFeedHttpClient
from pgchangefeed.models import RegisterConsumerRequest
from pgchangefeed.options import ClientOptions

_ADDR = os.environ["PGCHANGEFEED_HTTP_ADDR"]
_ADMIN_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN_ADMIN"]
_READER_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN_READER"]
_SOURCE = os.environ["PGCHANGEFEED_SOURCE_ID"]
_PUBLICATION = os.environ["PGCHANGEFEED_HTTP_PUBLICATION"]


def test_realserver_registers_a_consumer_and_lists_tables() -> None:
    http_client = httpx.Client(timeout=30.0)
    admin_client = PgChangeFeedHttpClient(
        http_client, ClientOptions(address=_ADDR, api_token=_ADMIN_TOKEN)
    )
    print("READY", flush=True)

    consumer_id = f"python-sdk-e2e-{uuid.uuid4().hex[:12]}"
    registered = admin_client.register_consumer(
        RegisterConsumerRequest(consumer_id=consumer_id, name=f"Python SDK E2E {consumer_id}")
    )
    assert registered.consumer_id == consumer_id
    assert registered.name == f"Python SDK E2E {consumer_id}"
    assert registered.already_registered is False

    reader_client = PgChangeFeedHttpClient(
        http_client, ClientOptions(address=_ADDR, api_token=_READER_TOKEN)
    )
    tables = reader_client.list_tables(_SOURCE, _PUBLICATION)

    print(
        f"RECEIVED consumer_id={registered.consumer_id} tables={len(tables.tables)}",
        flush=True,
    )


def test_realserver_rejects_the_call_without_token() -> None:
    http_client = httpx.Client(timeout=15.0)
    client = PgChangeFeedHttpClient(
        http_client, ClientOptions(address=_ADDR, api_token="no-such-token")
    )
    with pytest.raises(PgChangeFeedUnauthorizedError) as exc_info:
        client.list_tables(_SOURCE, _PUBLICATION)
    assert exc_info.value.status_code == 401
    print("REJECTED status=401", flush=True)