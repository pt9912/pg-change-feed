"""Client for the PG Change Feed live change stream over NATS.

``stream_changes`` subscribes to the subjects ``cdc.stream.<source_id>.>`` --
all tables of one source -- and yields one ``pgchangefeed.models.StreamChange``
per message: the same ten fields as the SSE event. Messages arrive as JSON
bytes; embedded row images stay JSON values, a missing image is ``None``.

The token is a connection-level credential: the NATS server rejects a
connection with a missing or wrong token when it is configured with a stream
token. The token is passed once when the connection is opened, not per
message. ``ClientOptions.address`` is the NATS URL (``nats://host:4222``),
``api_token`` the stream token.

The stream is fire-and-forget: a change committed while the client is not
subscribed is not delivered later. ``PgChangeFeedHttpClient.read_changes`` is
the way to catch up; the stream itself carries no filter.

The client uses ``nats-py``, the asyncio client of the NATS project, behind a
synchronous iterator like the other clients of this package: connecting and
subscribing run in a dedicated thread with its own event loop, and the
generator consumes a queue. The asyncio part stays inside this module.
"""

from __future__ import annotations

import asyncio
import json
import queue
import threading
import time
from typing import Iterator

import nats

from pgchangefeed.exceptions import PgChangeFeedMalformedResponseError
from pgchangefeed.models import StreamChange
from pgchangefeed.options import ClientOptions

_SUBJECT_PREFIX = "cdc.stream."

_SUBSCRIPTION_READY_TIMEOUT_SECONDS = 15.0
_STREAM_POLL_INTERVAL_SECONDS = 0.1


class PgChangeFeedNatsStreamClient:
    """Client for the NATS live change stream of one source.

    ``options`` carries the NATS URL and the stream token; ``source_id`` is the
    source whose tables the stream covers (it must not be empty).
    """

    def __init__(self, options: ClientOptions, source_id: str) -> None:
        if not source_id or not source_id.strip():
            raise ValueError("source_id must not be empty")
        self._options = options
        self._source_id = source_id.strip()

    def stream_changes(self, timeout: float | None = None) -> Iterator[StreamChange]:
        """Subscribes to all tables of the source and yields one ``StreamChange``
        per committed change, from the moment the subscription is ready
        onward: fire-and-forget, no replay, one message per row change.

        The token is sent when the connection is opened. If the server rejects
        it, the connect error of the NATS client library is raised (the stream
        is never silently empty). ``timeout`` bounds the total consumption in
        seconds (``None`` = unbounded); past it the generator stops.
        """
        events: queue.Queue[bytes] = queue.Queue()
        ready = threading.Event()
        connection_error: list[BaseException] = []
        stop = threading.Event()

        def consume() -> None:
            async def run() -> None:
                try:
                    nc = await nats.connect(
                        servers=[self._options.address],
                        token=self._options.api_token,
                        connect_timeout=_SUBSCRIPTION_READY_TIMEOUT_SECONDS,
                        allow_reconnect=False,
                    )
                except BaseException as exc:
                    connection_error.append(exc)
                    ready.set()
                    return

                async def callback(message) -> None:
                    events.put(message.data)

                await nc.subscribe(_subject_namespace(self._source_id), cb=callback)
                ready.set()
                while not stop.is_set():
                    await asyncio.sleep(_STREAM_POLL_INTERVAL_SECONDS)
                await nc.drain()

            asyncio.run(run())

        thread = threading.Thread(target=consume, daemon=True)
        thread.start()

        deadline = None if timeout is None else time.monotonic() + timeout
        if not ready.wait(_SUBSCRIPTION_READY_TIMEOUT_SECONDS):
            raise TimeoutError("NATS subscription did not become ready in time")
        if connection_error:
            raise connection_error[0]

        try:
            while True:
                if deadline is not None and time.monotonic() > deadline:
                    return
                try:
                    payload = events.get(timeout=_STREAM_POLL_INTERVAL_SECONDS)
                except queue.Empty:
                    continue
                yield _parse_stream_change(payload)
        finally:
            stop.set()


def _subject_namespace(source_id: str) -> str:
    """Builds the subscription subject of one source: the wildcard
    ``cdc.stream.<source_id>.>``, which matches every ``<schema>.<table>`` of
    the source."""
    return f"{_SUBJECT_PREFIX}{source_id}.>"


def _parse_stream_change(payload: bytes) -> StreamChange:
    try:
        body = json.loads(payload.decode("utf-8"))
    except (ValueError, UnicodeDecodeError) as exc:
        raise PgChangeFeedMalformedResponseError(
            200,
            "PG Change Feed NATS stream carried a message whose payload is "
            "not valid JSON.",
        ) from exc
    if not isinstance(body, dict):
        raise PgChangeFeedMalformedResponseError(
            200,
            "PG Change Feed NATS stream carried a message whose payload is "
            "not a JSON object.",
        )
    try:
        return StreamChange.from_json(body)
    except (KeyError, TypeError) as exc:
        raise PgChangeFeedMalformedResponseError(
            200,
            "PG Change Feed NATS stream carried a message missing an "
            "expected field.",
        ) from exc
