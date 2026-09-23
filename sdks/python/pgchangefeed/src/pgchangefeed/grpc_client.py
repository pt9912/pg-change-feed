"""Public entry point for the PG Change Feed live-change stream over gRPC
(SPEC-020, LH-FA-SST-008).

One method opens the ``ChangeStream/StreamChanges`` server-streaming RPC and
yields the generated ``Change`` message — the ten SPEC-020 fields
(``change_id``, ``transaction_id``, ``source_table_id``, ``sequence``,
``operation``, ``old_image``, ``new_image``, ``schema_version``, ``schema``,
``table``) — unmapped, exactly as the wire defines them. There is no separate
DTO layer here (unlike ``pgchangefeed.models`` for the HTTP surface): the
generated stub already is a typed representation of the wire schema; a
second, hand-mapped type would only risk drifting from it.

The bearer token is sent in the ``authorization`` gRPC metadata entry as
``Bearer <token>`` (SPEC-020) on every call. The ``grpc.Channel`` is
injected, not owned — the caller controls its lifetime and sharing; this
type never closes it. No module-level or global state, a process can hold
several independently configured instances at once.

A missing or invalid bearer token ends the call with gRPC status
``Unauthenticated`` (SPEC-020 Negative) — this surfaces as a ``grpc.RpcError``
from the enumeration itself, not a swallowed empty stream.

Boundary (SPEC-020, LH-FA-SST-008): the stream carries no replay and no
table-granular filtering. A consumer that needs either uses the existing read
path (``PgChangeFeedHttpClient.read_changes``), not this stream.

Draht-Kenntnis-Quelle: ``spec/pflichtenheft.md`` SPEC-020 (direkt),
``examples/csharp/grpc-client/Program.cs``/``examples/kotlin/grpc-client``
als funktionierende Fremdsprachen-Referenzen gegen denselben Server (gelesen,
nicht importiert — kein Python-Import eines privaten Baums dieses Repos);
kein ``examples/python/``-Referenz-Client existiert (ADR-0110
§Entscheidung Festlegung 2).
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
    """Client for the ``StreamChanges`` server-streaming RPC (SPEC-020)."""

    def __init__(self, channel: grpc.Channel, options: ClientOptions) -> None:
        self._client = changestream_pb2_grpc.ChangeStreamStub(channel)
        self._options = options

    def stream_changes(self, timeout: float | None = None) -> Iterator[_Change]:
        """Opens the server stream and yields every ``Change`` the server
        sends from connection time onward (SPEC-020: fire-and-forget, no
        replay, one message per row change in commit order). The request
        carries no filter (SPEC-020: table-granular filtering is not part of
        this version). A missing or invalid bearer token ends the call with
        gRPC status ``Unauthenticated`` (SPEC-020 Negative), raised by the
        iterator, not swallowed. ``timeout`` is the overall gRPC call
        deadline in seconds (``None`` = unbounded, the grpcio default); the
        iterator raises ``DEADLINE_EXCEEDED`` once it lapses."""
        metadata = ((_AUTHORIZATION_METADATA_KEY, f"{_BEARER_PREFIX}{self._options.api_token}"),)
        return self._client.StreamChanges(
            _StreamChangesRequest(), metadata=metadata, timeout=timeout
        )