"""Real-server routing phase for the NATS live change stream client.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. A subscription opened with `target` receives
the changes routed to that target and no change of another target or without
a target (checked over a quiet window), while a subscription opened without
`target` receives all of them.
"""

from __future__ import annotations

import os
from collections.abc import Callable

import pgchangefeed.nats_stream_client as nats_stream_client
from route_scenario import (
    ROUTE_TARGET_A,
    STREAM_TIMEOUT_SECONDS,
    RouteCollector,
    RouteRow,
    row_from,
    run_streams,
)

from pgchangefeed.options import ClientOptions

_NATS_URL = os.environ["PGCHANGEFEED_NATS_URL"]
_TOKEN = os.environ["PGCHANGEFEED_NATS_STREAM_TOKEN"]
_SOURCE_ID = os.environ["PGCHANGEFEED_SOURCE_ID"]


def _consume(target: str | None) -> Callable[[Callable[[RouteRow], None]], None]:
    client = nats_stream_client.PgChangeFeedNatsStreamClient(
        ClientOptions(address=_NATS_URL, api_token=_TOKEN), _SOURCE_ID
    )

    def consume(emit: Callable[[RouteRow], None]) -> None:
        for change in client.stream_changes(timeout=STREAM_TIMEOUT_SECONDS, target=target):
            emit(row_from(change.change_id, change.table, change.new_image))

    return consume


def test_realserver_target_subject_receives_only_its_target_and_source_subject_receives_all() -> None:
    targeted = RouteCollector()
    unfiltered = RouteCollector()
    targeted.start(_consume(ROUTE_TARGET_A))
    unfiltered.start(_consume(None))

    run_streams(targeted, unfiltered)
