"""Real-server routing phase for the SSE stream client.

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

import httpx
from route_scenario import ROUTE_TARGET_A, RouteCollector, RouteRow, row_from, run_streams

from pgchangefeed.options import ClientOptions
from pgchangefeed.sse_client import PgChangeFeedSseClient

_ADDR = os.environ["PGCHANGEFEED_HTTP_ADDR"]
_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN"]


def _consume(client: PgChangeFeedSseClient, target: str | None) -> Callable[[Callable[[RouteRow], None]], None]:
    def consume(emit: Callable[[RouteRow], None]) -> None:
        for change in client.stream_changes(target=target):
            emit(row_from(change.change_id, change.table, change.new_image))

    return consume


def test_realserver_stream_with_target_receives_only_its_target_and_stream_without_target_receives_all() -> None:
    options = ClientOptions(address=_ADDR, api_token=_TOKEN)
    targeted_http = httpx.Client(timeout=None)
    unfiltered_http = httpx.Client(timeout=None)
    targeted = RouteCollector()
    unfiltered = RouteCollector()
    targeted.start(_consume(PgChangeFeedSseClient(targeted_http, options), ROUTE_TARGET_A))
    unfiltered.start(_consume(PgChangeFeedSseClient(unfiltered_http, options), None))

    run_streams(targeted, unfiltered)

    targeted_http.close()
    unfiltered_http.close()
