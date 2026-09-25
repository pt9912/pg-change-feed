"""Client for the PG Change Feed live change stream over Server-Sent Events.

``stream_changes`` opens ``GET /changes/stream`` and yields one
``pgchangefeed.models.StreamChange`` per committed change, with the ten fields
``change_id``, ``transaction_id``, ``source_table_id``, ``sequence``,
``operation``, ``old_image``, ``new_image``, ``schema_version``, ``schema`` and
``table``. Row images are JSON values; a missing image is ``None``.

The bearer token is sent in the ``Authorization`` header as ``Bearer <token>``;
a ``reader`` or an ``admin`` token works. The ``httpx.Client`` is passed in,
not owned: the caller controls its lifetime, connection pooling and transport
(including ``httpx.MockTransport`` for tests); this class never closes it.
There is no module-level or global state, so a process can hold several
independently configured instances at once.

The stream is fire-and-forget and has no replay: the ``Last-Event-ID`` header
is neither sent nor evaluated, and the stream cannot be filtered by table.
``PgChangeFeedHttpClient.read_changes`` is the way to catch up.

Each event is a frame: an ``event: change`` line, one ``data:`` line with the
JSON object, and a blank line that ends the frame. A stream that ends before a
frame's blank line discards that incomplete frame. A frame with another event
name is skipped; a ``change`` frame whose JSON body cannot be read raises
``PgChangeFeedMalformedResponseError``. A non-200 response raises the same
typed status-code error as the HTTP client.
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
    """Client for the SSE live change stream ``GET /changes/stream``.

    ``client`` is the ``httpx.Client`` used for the request; ``options``
    carries the HTTP base URL and the bearer token. For a long-lived stream the
    client's read timeout must be off (``httpx.Timeout(10.0, read=None)``).
    """

    def __init__(self, client: httpx.Client, options: ClientOptions) -> None:
        self._client = client
        self._options = options

    def stream_changes(self) -> Iterator[StreamChange]:
        """Opens the SSE stream and yields one ``StreamChange`` per committed
        change, from connection time onward: fire-and-forget, no replay, one
        message per row change in commit order. The request carries no filter.

        A missing or invalid bearer token, or any other non-success response,
        raises the same typed status-code error as the HTTP client (``401`` for
        a bad token); the stream is never silently empty.
        """
        url = f"{self._options.address.rstrip('/')}{_STREAM_PATH}"
        headers = {"Authorization": f"Bearer {self._options.api_token}"}
        with self._client.stream("GET", url, headers=headers) as response:
            if not 200 <= response.status_code < 300:
                response.read()
                raise _build_error(response)
            yield from _stream_changes(response.iter_lines())


def _stream_changes(lines: Iterator[str]) -> Iterator[StreamChange]:
    """Yields one ``StreamChange`` per complete ``event: change`` frame.

    Reads line by line from the iterator: ``event: `` and ``data: `` lines
    collect the fields and the blank line closes the frame. A frame whose blank
    line never arrives is discarded; an incomplete frame is not an event."""
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
            "not valid JSON.",
        ) from exc
    if not isinstance(body, dict):
        raise PgChangeFeedMalformedResponseError(
            200,
            "PG Change Feed SSE stream carried an event whose data payload is "
            "not a JSON object.",
        )
    try:
        return StreamChange.from_json(body)
    except (KeyError, TypeError) as exc:
        raise PgChangeFeedMalformedResponseError(
            200,
            "PG Change Feed SSE stream carried an event missing an expected "
            "field.",
        ) from exc
