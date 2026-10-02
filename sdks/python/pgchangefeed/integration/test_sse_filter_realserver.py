"""Real-server filter phase for the SSE stream client.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. Three tables are captured: `A` and `B` in the
first schema and a table named like `A` in a second schema. A stream opened
with `schema` and `table` of `A` receives exactly the changes of that table, a
stream opened with the second `schema` alone receives exactly the changes of
that schema, and a stream without a filter receives all three (checked over a
quiet window after the unfiltered stream has them all).
"""

from __future__ import annotations

import os
import threading
import time
from collections.abc import Callable
from dataclasses import dataclass
from typing import Any

import httpx

from pgchangefeed.options import ClientOptions
from pgchangefeed.sse_client import PgChangeFeedSseClient

_ADDR = os.environ["PGCHANGEFEED_HTTP_ADDR"]
_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN"]
_SENTINEL = os.environ["PGCHANGEFEED_E2E_SENTINEL"]
_SCHEMA_A = os.environ["PGCHANGEFEED_FILTER_SCHEMA_A"]
_SCHEMA_OTHER = os.environ["PGCHANGEFEED_FILTER_SCHEMA_OTHER"]
_TABLE_A = os.environ["PGCHANGEFEED_FILTER_TABLE_A"]
_TABLE_B = os.environ["PGCHANGEFEED_FILTER_TABLE_B"]
_QUIET_SECONDS = int(os.environ["PGCHANGEFEED_FILTER_QUIET_SECONDS"])

_POSITIVE_DEADLINE_SECONDS = 90.0


@dataclass(frozen=True)
class FilterRow:
    """One received change reduced to what the filter phase asserts on."""

    change_id: str
    schema: str
    table: str
    name: str | None


def _row(change_id: str, schema: str, table: str, image: Any) -> FilterRow:
    name = image.get("name") if isinstance(image, dict) else None
    return FilterRow(change_id, schema, table, name if isinstance(name, str) else None)


class FilterCollector:
    """Collects the rows of one stream on a daemon thread."""

    def __init__(self) -> None:
        self._lock = threading.Lock()
        self._rows: list[FilterRow] = []
        self._failure: BaseException | None = None
        self._stopped = False

    @property
    def rows(self) -> list[FilterRow]:
        with self._lock:
            return list(self._rows)

    @property
    def failure(self) -> BaseException | None:
        with self._lock:
            return self._failure

    def start(self, consume: Callable[[Callable[[FilterRow], None]], None]) -> None:
        """Runs `consume` on a daemon thread; it hands every received row to the
        callback. A fault before `stop` is kept in `failure`; the read blocks
        while the stream is idle, so the thread ends with the process."""

        def emit(row: FilterRow) -> None:
            with self._lock:
                self._rows.append(row)

        def run() -> None:
            try:
                consume(emit)
            except Exception as exc:
                with self._lock:
                    if not self._stopped:
                        self._failure = exc

        threading.Thread(target=run, daemon=True).start()

    def stop(self) -> None:
        with self._lock:
            self._stopped = True


def _own(rows: list[FilterRow]) -> list[FilterRow]:
    return [r for r in rows if r.name == _SENTINEL]


def _is_table_a(r: FilterRow) -> bool:
    return r.schema == _SCHEMA_A and r.table == _TABLE_A


def _is_table_b(r: FilterRow) -> bool:
    return r.schema == _SCHEMA_A and r.table == _TABLE_B


def _is_other_schema(r: FilterRow) -> bool:
    return r.schema == _SCHEMA_OTHER


def _has_all_three(rows: list[FilterRow]) -> bool:
    own = _own(rows)
    return any(map(_is_table_a, own)) and any(map(_is_table_b, own)) and any(map(_is_other_schema, own))


def _raise_on_failure(*collectors: FilterCollector) -> None:
    for collector in collectors:
        failure = collector.failure
        if failure is not None:
            raise RuntimeError("Ein Stream-Konsument ist ausgefallen.") from failure


def _consume(
    client: PgChangeFeedSseClient, schema: str | None, table: str | None
) -> Callable[[Callable[[FilterRow], None]], None]:
    def consume(emit: Callable[[FilterRow], None]) -> None:
        for change in client.stream_changes(schema=schema, table=table):
            emit(_row(change.change_id, change.schema, change.table, change.new_image))

    return consume


def test_realserver_streams_with_schema_and_table_filter_receive_only_their_selection_and_stream_without_filter_receives_all() -> None:
    options = ClientOptions(address=_ADDR, api_token=_TOKEN)
    http_clients = [httpx.Client(timeout=None) for _ in range(3)]
    f1 = FilterCollector()
    f2 = FilterCollector()
    unfiltered = FilterCollector()
    f1.start(_consume(PgChangeFeedSseClient(http_clients[0], options), _SCHEMA_A, _TABLE_A))
    f2.start(_consume(PgChangeFeedSseClient(http_clients[1], options), _SCHEMA_OTHER, None))
    unfiltered.start(_consume(PgChangeFeedSseClient(http_clients[2], options), None, None))

    print("READY", flush=True)
    deadline = time.monotonic() + _POSITIVE_DEADLINE_SECONDS
    while not (
        any(map(_is_table_a, _own(f1.rows)))
        and any(map(_is_other_schema, _own(f2.rows)))
        and _has_all_three(unfiltered.rows)
    ):
        _raise_on_failure(f1, f2, unfiltered)
        assert time.monotonic() < deadline, (
            "innerhalb der Frist weder die Change der Tabelle am Client mit Schema und Tabelle, "
            "noch die des zweiten Schemas am Client mit Schema, noch alle drei Gruppen am Client ohne Filter empfangen"
        )
        time.sleep(0.1)
    print("SEEN", flush=True)

    # The window starts only after the client without a filter received all
    # three groups: the same delivery path has then demonstrably dispatched
    # the foreign changes, so their absence at the filtered clients is no
    # early cut-off.
    time.sleep(_QUIET_SECONDS)
    _raise_on_failure(f1, f2, unfiltered)
    f1.stop()
    f2.stop()
    unfiltered.stop()

    f1_rows = f1.rows
    f2_rows = f2.rows
    own_unfiltered = _own(unfiltered.rows)
    f1_foreign = [r for r in f1_rows if not _is_table_a(r)]
    f2_foreign = [r for r in f2_rows if not _is_other_schema(r)]
    for r in f1_rows:
        print(f"RECEIVED_F1 change_id={r.change_id} schema={r.schema} table={r.table}", flush=True)
    for r in f2_rows:
        print(f"RECEIVED_F2 change_id={r.change_id} schema={r.schema} table={r.table}", flush=True)
    for r in own_unfiltered:
        print(f"RECEIVED_U change_id={r.change_id} schema={r.schema} table={r.table}", flush=True)
    # The line carries the measured counts of foreign changes before the checks
    # below, so a foreign receipt is visible in the line itself.
    print(
        f"FILTER_RESULT f1={len(f1_rows)} f1_foreign={len(f1_foreign)} f2={len(f2_rows)} "
        f"f2_foreign={len(f2_foreign)} unfiltered={len(own_unfiltered)} quiet_seconds={_QUIET_SECONDS}",
        flush=True,
    )

    for client in http_clients:
        client.close()
    assert not f1_foreign, "der Client mit Schema und Tabelle empfing fremde Changes: " + ", ".join(
        f"{r.change_id}({r.schema}.{r.table})" for r in f1_foreign
    )
    assert not f2_foreign, "der Client mit Schema allein empfing fremde Changes: " + ", ".join(
        f"{r.change_id}({r.schema}.{r.table})" for r in f2_foreign
    )
    assert _own(f1_rows), "der Client mit Schema und Tabelle empfing keine Change dieser Phase"
    assert _own(f2_rows), "der Client mit Schema allein empfing keine Change dieser Phase"
    assert _has_all_three(unfiltered.rows), "der Client ohne Filter sah nicht alle drei Gruppen"
