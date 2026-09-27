"""Real-server rule phase for the HTTP API client.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. Reads `GET /changes` (source and table as
filter) against a table carrying an active rename_column rule; the row
image of the change committed while the test polls arrives with the renamed
key and without the source key, read through the client's opaque value.
"""

from __future__ import annotations

import os
import time

import httpx

from pgchangefeed.http_client import PgChangeFeedHttpClient
from pgchangefeed.models import Change
from pgchangefeed.options import ClientOptions

_ADDR = os.environ["PGCHANGEFEED_HTTP_ADDR"]
_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN"]
_SOURCE = os.environ["PGCHANGEFEED_SOURCE_ID"]
_TABLE = os.environ["PGCHANGEFEED_E2E_TABLE"]
_SENTINEL = os.environ["PGCHANGEFEED_E2E_SENTINEL"]
_SOURCE_KEY = os.environ["PGCHANGEFEED_RULE_SOURCE_KEY"]
_TARGET_KEY = os.environ["PGCHANGEFEED_RULE_TARGET_KEY"]

_RECEIVE_DEADLINE_SECONDS = 90.0


def test_realserver_reads_a_renamed_row_image_key_over_the_api() -> None:
    http_client = httpx.Client(timeout=30.0)
    client = PgChangeFeedHttpClient(http_client, ClientOptions(address=_ADDR, api_token=_TOKEN))
    print("READY", flush=True)

    received = _receive_sentinel(client)

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


def _receive_sentinel(client: PgChangeFeedHttpClient) -> Change:
    deadline = time.monotonic() + _RECEIVE_DEADLINE_SECONDS
    while time.monotonic() < deadline:
        response = client.read_changes(_SOURCE, "public", _TABLE)
        for change in response.changes:
            if (
                change.table == _TABLE
                and change.operation == "INSERT"
                and _SENTINEL in str(change.new_image)
            ):
                return change
        time.sleep(0.5)
    raise AssertionError(
        f"kein Change mit dem Sentinel {_SENTINEL!r} innerhalb der Frist über "
        f"GET /changes gelesen"
    )
