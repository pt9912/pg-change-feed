"""Real-server routing phase for the HTTP API client.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. `read_changes` with `target` returns exactly
the changes routed to that target (all of region A, and every one the
unfiltered read assigns to A), while the call without `target` returns all of
them.
"""

from __future__ import annotations

import os

import httpx
from route_scenario import TABLE, RouteRow, row_from, run_pull

from pgchangefeed.http_client import PgChangeFeedHttpClient
from pgchangefeed.options import ClientOptions

_ADDR = os.environ["PGCHANGEFEED_HTTP_ADDR"]
_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN"]
_SOURCE = os.environ["PGCHANGEFEED_SOURCE_ID"]


def test_realserver_read_with_target_returns_only_its_target_and_read_without_target_returns_all() -> None:
    http_client = httpx.Client(timeout=30.0)
    client = PgChangeFeedHttpClient(http_client, ClientOptions(address=_ADDR, api_token=_TOKEN))

    def read(target: str | None) -> list[RouteRow]:
        response = client.read_changes(_SOURCE, "public", TABLE, target=target)
        return [row_from(c.change_id, c.table, c.new_image) for c in response.changes]

    run_pull(read)
