"""Real-server TLS phase for the gRPC stream client.

With the trust anchor of the options the SDK opens the server stream over TLS
and receives a change committed afterwards (the runner reads its `change_id`
back from `cdc.changes`); without an anchor, with a foreign anchor and with a
server name outside the certificate the stream fails the TLS check. The
address is `host:port`; the anchor in the options makes the channel a TLS
channel.
"""

from __future__ import annotations

import os
import time

from tls_scenarios import CA_FILE, expect_refusals

from pgchangefeed import ClientOptions, PgChangeFeedGrpcClient, create_grpc_channel

_ADDR = os.environ["PGCHANGEFEED_GRPC_ADDR"]
_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN"]
_TABLE = os.environ["PGCHANGEFEED_E2E_TABLE"]
_SENTINEL = os.environ["PGCHANGEFEED_E2E_SENTINEL"]

_RECEIVE_DEADLINE_SECONDS = 90.0


def test_realserver_receives_a_committed_change_over_tls() -> None:
    options = ClientOptions(address=f"https://{_ADDR}", api_token=_TOKEN, trust_anchor_file=CA_FILE)
    with create_grpc_channel(options) as channel:
        stream = PgChangeFeedGrpcClient(channel, options).stream_changes(timeout=_RECEIVE_DEADLINE_SECONDS)
        print("READY", flush=True)
        deadline = time.monotonic() + _RECEIVE_DEADLINE_SECONDS
        for change in stream:
            assert time.monotonic() < deadline, "no change with the sentinel within the deadline"
            if change.table == _TABLE and change.operation == "INSERT" and _SENTINEL in change.new_image.decode():
                print(
                    f"RECEIVED change_id={change.change_id} table={change.table} "
                    f"operation={change.operation} new_image={change.new_image.decode()}",
                    flush=True,
                )
                return
    raise AssertionError("the stream ended before the sentinel arrived")


def test_realserver_without_anchor_foreign_anchor_and_wrong_name_fail_the_tls_check() -> None:
    def call(options: ClientOptions) -> None:
        with create_grpc_channel(options) as channel:
            next(iter(PgChangeFeedGrpcClient(channel, options).stream_changes(timeout=15)), None)

    expect_refusals(f"https://{_ADDR}", _TOKEN, call)
