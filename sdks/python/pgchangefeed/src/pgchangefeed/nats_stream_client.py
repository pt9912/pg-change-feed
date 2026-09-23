"""Public entry point for the PG Change Feed live-change stream over NATS
(SPEC-024, LH-FA-SST-008), third delivery path beside gRPC (SPEC-020) and
SSE (SPEC-021).

One method subscribes the full-content namespace
``cdc.stream.<source_id>.>`` (SPEC-024: all tables of one source) and yields
one ``pgchangefeed.models.StreamChange`` per message — the same ten fields as
the SSE event, no third schema. Messages arrive as JSON bytes; embedded row
images stay JSON values, a missing image is ``None``. The typed hand-mapped
layer mirrors the SSE surface: the wire is JSON, and ``models`` carries the
package's JSON-mapping doctrine.

The token is a **connection-level** credential (SPEC-024: the NATS server
rejects a connection without or with a wrong token, as soon as
``CDC_NATS_STREAM_TOKEN`` is configured): the token goes into the connect
call, not per message. The ``ClientOptions.address`` carries the NATS URL
(``nats://host:4222``), ``api_token`` the serverwide stream token.

Delivery guarantee: none (Core NATS, fire-and-forget, no replay) — the
existing read path (``PgChangeFeedHttpClient.read_changes``) stays the
catch-up path; this stream carries no filter.

The NATS Python client is ``nats-py`` (the official asyncio client of the
NATS project); its asyncio model is wrapped behind the package's synchronous
surface vocabulary (``httpx.Client``, synchronous gRPC stub): the connect
and subscription run in a dedicated thread with its own event loop, the
generator consumes a queue — one iterator form across all four surfaces, the
asyncio boundary stays inside this module.

Draht-Kenntnis-Quelle: ``spec/pflichtenheft.md`` SPEC-024 (direkt) und
``internal/adapters/driven/natsstream/publisher.go`` als serverseitige
Gegenprobe (gelesen, nicht importiert — kein Python-Import eines privaten
Baums dieses Repos).
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
    """Client for the NATS full-content stream (SPEC-024)."""

    def __init__(self, options: ClientOptions, source_id: str) -> None:
        if not source_id or not source_id.strip():
            raise ValueError("source_id must not be empty")
        self._options = options
        self._source_id = source_id.strip()

    def stream_changes(self, timeout: float | None = None) -> Iterator[StreamChange]:
        """Subscribes ``cdc.stream.<source_id>.>`` (SPEC-024: all tables of
        one source) and yields one ``StreamChange`` per committed change
        (fire-and-forget, no replay, one message per row change in commit
        order). The bearer token rides the connection (SPEC-024: the server
        rejects a connection without or with a wrong token — that surfaces as
        the connect error of the NATS client library, not a swallowed empty
        stream). ``timeout`` bounds the total consumption in seconds
        (``None`` = unbounded); past it the generator stops."""
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
    """Builds the subscription subject for one source (SPEC-024): the
    wildcard namespace ``cdc.stream.<source_id>.>`` — all tables of one
    source; ``<schema>.<table>`` kommen je Nachricht."""
    return f"{_SUBJECT_PREFIX}{source_id}.>"


def _parse_stream_change(payload: bytes) -> StreamChange:
    try:
        body = json.loads(payload.decode("utf-8"))
    except (ValueError, UnicodeDecodeError) as exc:
        raise PgChangeFeedMalformedResponseError(
            200,
            "PG Change Feed NATS stream carried a message whose payload is "
            "not valid JSON -- a protocol violation outside SPEC-024's "
            "documented shape.",
        ) from exc
    if not isinstance(body, dict):
        raise PgChangeFeedMalformedResponseError(
            200,
            "PG Change Feed NATS stream carried a message whose payload is "
            "not a JSON object -- a protocol violation outside SPEC-024's "
            "documented shape.",
        )
    try:
        return StreamChange.from_json(body)
    except (KeyError, TypeError) as exc:
        raise PgChangeFeedMalformedResponseError(
            200,
            "PG Change Feed NATS stream carried a message missing an "
            "expected SPEC-024 field -- a protocol violation outside the "
            "documented shapes.",
        ) from exc
