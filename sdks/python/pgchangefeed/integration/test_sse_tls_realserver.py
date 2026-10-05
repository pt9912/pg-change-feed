"""Real-server TLS phase for the SSE stream client.

With the trust anchor of the options the SDK opens `GET /changes/stream` over
TLS and receives a change committed afterwards (the runner reads its
`change_id` back from `cdc.changes`); without an anchor, with a foreign anchor
and with a server name outside the certificate the stream fails the TLS check.
"""

from __future__ import annotations

import os

import httpx

from tls_scenarios import CA_FILE, expect_refusals

from pgchangefeed import ClientOptions, PgChangeFeedSseClient, create_http_client

_ADDR = os.environ["PGCHANGEFEED_HTTP_ADDR"]
_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN"]
_TABLE = os.environ["PGCHANGEFEED_E2E_TABLE"]
_SENTINEL = os.environ["PGCHANGEFEED_E2E_SENTINEL"]

_NO_READ_TIMEOUT = httpx.Timeout(10.0, read=None)


def test_realserver_receives_a_committed_change_over_tls() -> None:
    options = ClientOptions(address=_ADDR, api_token=_TOKEN, trust_anchor_file=CA_FILE)
    with create_http_client(options, timeout=_NO_READ_TIMEOUT) as http:
        stream = PgChangeFeedSseClient(http, options).stream_changes()
        print("READY", flush=True)
        for change in stream:
            if change.table == _TABLE and change.operation == "INSERT" and _SENTINEL in str(change.new_image):
                print(
                    f"RECEIVED change_id={change.change_id} table={change.table} "
                    f"operation={change.operation} new_image={change.new_image}",
                    flush=True,
                )
                return
    raise AssertionError("the stream ended before the sentinel arrived")


def test_realserver_without_anchor_foreign_anchor_and_wrong_name_fail_the_tls_check() -> None:
    def call(options: ClientOptions) -> None:
        with create_http_client(options, timeout=httpx.Timeout(10.0)) as http:
            next(iter(PgChangeFeedSseClient(http, options).stream_changes()), None)

    expect_refusals(_ADDR, _TOKEN, call)
