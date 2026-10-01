"""Shared course of the routing real-server phases.

Runs only in the `integration` Docker stage (sdks/python/Dockerfile), started
by tools/harness/run-sdk-python-integration-tests.sh inside the same Docker
network as the feed container. A client with the target `ROUTE_TARGET_A` and a
client without a target read the same table; the table carries two routing
rules (region A to target A, region B to target B), and the runner commits
groups of three changes (no rule, B, A). The client with the target must
receive the change for A and no change of any other region; the client
without a target must receive all three kinds.
"""

from __future__ import annotations

import json
import os
import threading
import time
from collections.abc import Callable
from dataclasses import dataclass
from typing import Any

TABLE = os.environ["PGCHANGEFEED_E2E_TABLE"]
SENTINEL = os.environ["PGCHANGEFEED_E2E_SENTINEL"]
ROUTE_TARGET_A = os.environ["PGCHANGEFEED_ROUTE_TARGET_A"]
ROUTE_TARGET_B = os.environ["PGCHANGEFEED_ROUTE_TARGET_B"]
ROUTE_REGION_NONE = os.environ["PGCHANGEFEED_ROUTE_REGION_NONE"]
QUIET_SECONDS = int(os.environ["PGCHANGEFEED_ROUTE_QUIET_SECONDS"])

POSITIVE_DEADLINE_SECONDS = 90.0
# Upper bound of one stream consumption: the positive phase, the quiet window
# and a margin; the consumer threads are daemons and end with the process.
STREAM_TIMEOUT_SECONDS = POSITIVE_DEADLINE_SECONDS + QUIET_SECONDS + 30.0


@dataclass(frozen=True)
class RouteRow:
    """One received change reduced to what the routing phases assert on."""

    change_id: str
    table: str
    region: str | None
    name: str | None


def row_from(change_id: str, table: str, image: Any) -> RouteRow:
    """Reads `region` and `name` from a row image (a dict or its JSON text)."""
    if isinstance(image, (bytes, bytearray)):
        image = image.decode()
    if isinstance(image, str):
        image = json.loads(image) if image else {}
    if not isinstance(image, dict):
        image = {}
    region = image.get("region")
    name = image.get("name")
    return RouteRow(
        change_id,
        table,
        region if isinstance(region, str) else None,
        name if isinstance(name, str) else None,
    )


class RouteCollector:
    """Collects the rows of one stream on a daemon thread."""

    def __init__(self) -> None:
        self._lock = threading.Lock()
        self._rows: list[RouteRow] = []
        self._failure: BaseException | None = None
        self._stopped = False

    @property
    def rows(self) -> list[RouteRow]:
        with self._lock:
            return list(self._rows)

    @property
    def failure(self) -> BaseException | None:
        with self._lock:
            return self._failure

    def start(self, consume: Callable[[Callable[[RouteRow], None]], None]) -> None:
        """Runs `consume` on a daemon thread; it hands every received row to the
        callback. A fault before `stop` is kept in `failure`; the read blocks
        while the stream is idle, so the thread ends with the process."""

        def emit(row: RouteRow) -> None:
            with self._lock:
                self._rows.append(row)

        def run() -> None:
            try:
                consume(emit)
            except Exception as exc:
                # Every fault of the consumer is kept in `failure`, so the
                # scenario reports it on the main thread instead of losing it.
                with self._lock:
                    if not self._stopped:
                        self._failure = exc

        threading.Thread(target=run, daemon=True).start()

    def stop(self) -> None:
        with self._lock:
            self._stopped = True


def _own(rows: list[RouteRow]) -> list[RouteRow]:
    return [r for r in rows if r.table == TABLE and r.name == SENTINEL]


def _has_all_three(rows: list[RouteRow]) -> bool:
    regions = {r.region for r in _own(rows)}
    return {ROUTE_TARGET_A, ROUTE_TARGET_B, ROUTE_REGION_NONE} <= regions


def _raise_on_failure(*collectors: RouteCollector) -> None:
    for collector in collectors:
        failure = collector.failure
        if failure is not None:
            raise RuntimeError("Ein Stream-Konsument ist ausgefallen.") from failure


def run_streams(targeted: RouteCollector, unfiltered: RouteCollector) -> None:
    """Stream surfaces: positive phase, quiet window, then the checks."""
    print("READY", flush=True)
    deadline = time.monotonic() + POSITIVE_DEADLINE_SECONDS
    while not (
        any(r.region == ROUTE_TARGET_A for r in _own(targeted.rows))
        and _has_all_three(unfiltered.rows)
    ):
        _raise_on_failure(targeted, unfiltered)
        assert time.monotonic() < deadline, (
            f"innerhalb der Frist weder die Change des Ziels {ROUTE_TARGET_A} am Client mit Ziel "
            "noch alle drei Gruppen am Client ohne Ziel empfangen"
        )
        time.sleep(0.1)
    print("SEEN", flush=True)

    # The window starts only after the client without a target received the
    # change of the other target: the same delivery path has then demonstrably
    # dispatched it, so its absence at the client with a target is no early
    # cut-off.
    time.sleep(QUIET_SECONDS)
    _raise_on_failure(targeted, unfiltered)
    targeted.stop()
    unfiltered.stop()

    _evaluate(targeted.rows, unfiltered.rows, QUIET_SECONDS)


def run_pull(read: Callable[[str | None], list[RouteRow]]) -> None:
    """Pull surface: the target read equals the region-A part of the unfiltered read."""
    print("READY", flush=True)
    deadline = time.monotonic() + POSITIVE_DEADLINE_SECONDS
    while not _has_all_three(read(None)):
        assert time.monotonic() < deadline, (
            "innerhalb der Frist nicht alle drei Gruppen über die ungefilterte Lesung gelesen"
        )
        time.sleep(0.3)
    print("SEEN", flush=True)

    before = read(None)
    targeted = read(ROUTE_TARGET_A)
    after = read(None)

    def region_a(rows: list[RouteRow]) -> set[str]:
        return {r.change_id for r in _own(rows) if r.region == ROUTE_TARGET_A}

    targeted_own = {r.change_id for r in _own(targeted)}
    assert region_a(before) <= targeted_own, (
        "die Lesung mit Ziel enthält nicht jede Change des Ziels, die die Lesung davor ohne Ziel lieferte"
    )
    assert targeted_own <= region_a(after), (
        "die Lesung mit Ziel enthält eine Change, die die Lesung danach ohne Ziel nicht dem Ziel zuordnet"
    )

    _evaluate(targeted, after, 0)


def _evaluate(targeted: list[RouteRow], unfiltered: list[RouteRow], quiet_seconds: int) -> None:
    foreign = [r for r in targeted if r.table != TABLE or r.region != ROUTE_TARGET_A]
    own_unfiltered = _own(unfiltered)
    for r in targeted:
        print(f"RECEIVED_TARGETED change_id={r.change_id} table={r.table} region={r.region}", flush=True)
    for r in own_unfiltered:
        print(f"RECEIVED_UNFILTERED change_id={r.change_id} table={r.table} region={r.region}", flush=True)
    # The line carries the measured count of foreign changes before the checks
    # below, so a foreign receipt is visible in the line itself.
    print(
        f"ROUTE_RESULT target={ROUTE_TARGET_A} targeted={len(targeted)} foreign={len(foreign)} "
        f"unfiltered={len(own_unfiltered)} quiet_seconds={quiet_seconds}",
        flush=True,
    )

    assert not foreign, (
        f"der Client mit Ziel {ROUTE_TARGET_A} empfing fremde Changes: "
        + ", ".join(f"{r.change_id}(region={r.region})" for r in foreign)
    )
    assert _own(targeted), "der Client mit Ziel empfing keine Change dieser Phase"
    assert _has_all_three(unfiltered), "der Client ohne Ziel sah nicht alle drei Gruppen"
