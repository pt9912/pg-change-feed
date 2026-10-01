"""Real-server routing phase for the gRPC stream client.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. A stream opened with `target` receives the
changes routed to that target and no change of another target or without a
target (checked over a quiet window), while a stream opened without `target`
receives all of them.
"""

from __future__ import annotations

import os
from collections.abc import Callable

import grpc
from route_scenario import (
    ROUTE_TARGET_A,
    STREAM_TIMEOUT_SECONDS,
    RouteCollector,
    RouteRow,
    row_from,
    run_streams,
)

from pgchangefeed.grpc_client import PgChangeFeedGrpcClient
from pgchangefeed.options import ClientOptions

_ADDR = os.environ["PGCHANGEFEED_GRPC_ADDR"]
_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN"]


def _consume(client: PgChangeFeedGrpcClient, target: str | None) -> Callable[[Callable[[RouteRow], None]], None]:
    def consume(emit: Callable[[RouteRow], None]) -> None:
        for change in client.stream_changes(timeout=STREAM_TIMEOUT_SECONDS, target=target):
            emit(row_from(change.change_id, change.table, change.new_image))

    return consume


def test_realserver_stream_with_target_receives_only_its_target_and_stream_without_target_receives_all() -> None:
    options = ClientOptions(address=_ADDR, api_token=_TOKEN)
    targeted_channel = grpc.insecure_channel(_ADDR)
    unfiltered_channel = grpc.insecure_channel(_ADDR)
    targeted = RouteCollector()
    unfiltered = RouteCollector()
    targeted.start(_consume(PgChangeFeedGrpcClient(targeted_channel, options), ROUTE_TARGET_A))
    unfiltered.start(_consume(PgChangeFeedGrpcClient(unfiltered_channel, options), None))

    run_streams(targeted, unfiltered)

    targeted_channel.close()
    unfiltered_channel.close()
