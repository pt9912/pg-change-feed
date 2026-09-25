"""Client for the PG Change Feed live change stream over gRPC.

``stream_changes`` opens the ``ChangeStream/StreamChanges`` server-streaming
call and yields the generated ``Change`` message with the ten fields
``change_id``, ``transaction_id``, ``source_table_id``, ``sequence``,
``operation``, ``old_image``, ``new_image``, ``schema_version``, ``schema`` and
``table``, exactly as the server sends them. Row images are JSON bytes. There
is no separate data class layer (unlike ``pgchangefeed.models`` for the HTTP
client): the generated message already is the typed form of the stream schema.

The bearer token is sent in the ``authorization`` call metadata entry as
``Bearer <token>``. The ``grpc.Channel`` is passed in, not owned: the caller
controls its lifetime and sharing; this class never closes it. There is no
module-level or global state, so a process can hold several independently
configured instances at once.

A missing or invalid bearer token ends the call with gRPC status
``UNAUTHENTICATED``; this is raised as ``grpc.RpcError`` while iterating, the
stream is never silently empty.

The stream is fire-and-forget: it has no replay and cannot be filtered by
table. ``PgChangeFeedHttpClient.read_changes`` is the way to catch up.
"""

from __future__ import annotations

from typing import Iterator

import grpc

from pgchangefeed.grpc_gen import changestream_pb2, changestream_pb2_grpc
from pgchangefeed.options import ClientOptions

_AUTHORIZATION_METADATA_KEY = "authorization"
_BEARER_PREFIX = "Bearer "

_Change = changestream_pb2.Change
_StreamChangesRequest = changestream_pb2.StreamChangesRequest


class PgChangeFeedGrpcClient:
    """Client for the gRPC ``StreamChanges`` server-streaming call.

    ``channel`` is the ``grpc.Channel`` to the server's gRPC address;
    ``options`` carries the bearer token.
    """

    def __init__(self, channel: grpc.Channel, options: ClientOptions) -> None:
        self._client = changestream_pb2_grpc.ChangeStreamStub(channel)
        self._options = options

    def stream_changes(self, timeout: float | None = None) -> Iterator[_Change]:
        """Opens the server stream and yields every ``Change`` the server sends
        from connection time onward: fire-and-forget, no replay, one message
        per row change in commit order. The request carries no filter.

        A missing or invalid bearer token ends the call with gRPC status
        ``UNAUTHENTICATED``, raised by the iterator. ``timeout`` is the
        deadline of the whole call in seconds (``None`` = unbounded, the grpcio
        default); the iterator raises ``DEADLINE_EXCEEDED`` once it lapses.
        """
        metadata = ((_AUTHORIZATION_METADATA_KEY, f"{_BEARER_PREFIX}{self._options.api_token}"),)
        return self._client.StreamChanges(
            _StreamChangesRequest(), metadata=metadata, timeout=timeout
        )