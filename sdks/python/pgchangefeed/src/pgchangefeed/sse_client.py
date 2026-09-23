"""Public entry point for the PG Change Feed live-change stream over HTTP
Server-Sent-Events (SPEC-021, LH-FA-SST-008).

One method opens ``GET /changes/stream`` and yields one
``pgchangefeed.models.StreamChange`` per committed change — the ten SPEC-021
fields (``change_id``, ``transaction_id``, ``source_table_id``, ``sequence``,
``operation``, ``old_image``, ``new_image``, ``schema_version``, ``schema``,
``table``), embedded row images stay JSON values, a missing image is
``None``. This is a typed, hand-mapped layer (unlike the gRPC surface, whose
generated stub already carries the wire types): the SSE wire is JSON, and
``models`` already carries the package's JSON-mapping doctrine.

The bearer token is sent in the ``Authorization`` header as
``Bearer <token>`` (SPEC-021: ``reader`` or ``admin`` — streaming is
read-only). The ``httpx.Client`` is injected, not owned — the caller
controls its lifetime, connection pooling and transport (including a
``httpx.MockTransport`` for tests); this type never closes it. No
module-level or global state, a process can hold several independently
configured instances at once.

Boundary (SPEC-021, LH-FA-SST-008): the stream carries no replay — the
``Last-Event-ID`` header is neither sent nor evaluated — and no
table-granular filtering. A consumer that needs either uses the existing
read path (``PgChangeFeedHttpClient.read_changes``), not this stream.

Frame form (server side, ``ADR-0061``/``SPEC-021``): ``event: change`` line,
one ``data:`` line carrying the JSON object, blank line ends the frame. The
parser is deliberately line-based over a line iterator (not chunk-boundary
aware) — chunk boundaries need not coincide with line boundaries. A stream
that ends before a frame's blank line discards the begun frame — an
incomplete frame is not an event (same boundary as the C#/Kotlin reference
clients). A frame with a non-``change`` event name is skipped; a ``change``
frame whose JSON body does not match the documented schema raises a typed
``PgChangeFeedMalformedResponseError``. A non-200 response becomes the same
typed status-code error as the HTTP surface.

Draht-Kenntnis-Quelle: ``spec/pflichtenheft.md`` SPEC-021 (direkt) und
``internal/adapters/driving/http/sse.go`` als serverseitige Gegenprobe
(gelesen, nicht importiert — kein Python-Import eines privaten Baums dieses
Repos).
"""

from __future__ import annotations

import json
from typing import Iterator

import httpx

from pgchangefeed.exceptions import PgChangeFeedMalformedResponseError
from pgchangefeed.http_client import _build_error
from pgchangefeed.models import StreamChange
from pgchangefeed.options import ClientOptions

_STREAM_PATH = "/changes/stream"
_EVENT_TYPE_CHANGE = "change"


class PgChangeFeedSseClient:
    """Client for the SSE live-change stream ``GET /changes/stream``
    (SPEC-021)."""

    def __init__(self, client: httpx.Client, options: ClientOptions) -> None:
        self._client = client
        self._options = options

    def stream_changes(self) -> Iterator[StreamChange]:
        """Opens the SSE stream and yields one ``StreamChange`` per committed
        change, from connection time onward (SPEC-021: fire-and-forget, no
        replay, one message per row change in commit order). The request
        carries no filter (SPEC-021 carries no filter). A missing or invalid
        bearer token, or any other non-200 response, raises the same typed
        status-code error as the HTTP surface (SPEC-021 Negative: ``401``,
        not a swallowed empty stream)."""
        url = f"{self._options.address.rstrip('/')}{_STREAM_PATH}"
        headers = {"Authorization": f"Bearer {self._options.api_token}"}
        with self._client.stream("GET", url, headers=headers) as response:
            if not 200 <= response.status_code < 300:
                response.read()
                raise _build_error(response)
            yield from _stream_changes(response.iter_lines())


def _stream_changes(lines: Iterator[str]) -> Iterator[StreamChange]:
    """Yields one ``StreamChange`` per complete ``event: change`` frame.

    Line-based over the caller's line iterator: ``event: ``/``data: ``-Zeilen
    sammeln die Felder, die Leerzeile schließt das Frame ab. Ein Frame, vor
    dessen Leerzeile der Stream endet, wird verworfen — ein unvollständiges
    Frame ist kein Event (Boundary wie bei den Referenz-Clients)."""
    event_name: str | None = None
    event_data: str | None = None
    for line in lines:
        if line.endswith("\r"):
            line = line[:-1]
        if line == "":
            if event_name is None and event_data is None:
                continue
            if event_name == _EVENT_TYPE_CHANGE and event_data is not None:
                yield _parse_stream_change(event_data)
            event_name = None
            event_data = None
        elif line.startswith("event: "):
            event_name = line[len("event: "):]
        elif line.startswith("data: "):
            event_data = line[len("data: "):]


def _parse_stream_change(payload: str) -> StreamChange:
    try:
        body = json.loads(payload)
    except ValueError as exc:
        raise PgChangeFeedMalformedResponseError(
            200,
            "PG Change Feed SSE stream carried an event whose data payload is "
            "not valid JSON -- a protocol violation outside SPEC-021's "
            "documented shapes.",
        ) from exc
    if not isinstance(body, dict):
        raise PgChangeFeedMalformedResponseError(
            200,
            "PG Change Feed SSE stream carried an event whose data payload is "
            "not a JSON object -- a protocol violation outside SPEC-021's "
            "documented shape.",
        )
    try:
        return StreamChange.from_json(body)
    except (KeyError, TypeError) as exc:
        raise PgChangeFeedMalformedResponseError(
            200,
            "PG Change Feed SSE stream carried an event missing an expected "
            "SPEC-021 field -- a protocol violation outside the documented "
            "shapes.",
        ) from exc