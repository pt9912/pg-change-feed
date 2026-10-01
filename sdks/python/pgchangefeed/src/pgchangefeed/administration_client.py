"""Client for the PG Change Feed Administration gRPC service.

One method per RPC -- ``register_consumer``, ``acknowledge_consumer``,
``get_consumer_position``, ``remove_consumer``, ``enable_table``,
``disable_table``, ``get_table_status``, ``list_tables``, ``run_retention``,
``read_changes`` and ``diagnose``. Requests and responses are the generated
protobuf messages of ``pgchangefeed.grpc_gen.administration_pb2`` unchanged --
like ``PgChangeFeedGrpcClient`` for the live stream, there is no separate data
class layer: a generated message is already a typed object, and a hand-written
mirror would only duplicate eleven message shapes without a deserialization
need to justify it (unlike ``pgchangefeed.models``, whose data classes exist
because JSON needs a target type to deserialize into).

Every non-``OK`` gRPC status becomes a typed ``PgChangeFeedGrpcError``
subclass instead of a result mixed with the success path, the same design as
``pgchangefeed.exceptions.PgChangeFeedError`` for the HTTP surface; the
original ``grpc.RpcError`` is always ``__cause__``.

The bearer token is sent in the ``authorization`` call metadata entry as
``Bearer <token>`` on every call. The ``grpc.Channel`` is passed in, not
owned: the caller controls its lifetime and sharing; this class never closes
it. There is no module-level or global state, so a process can hold several
independently configured instances at once.
"""

from __future__ import annotations

from typing import Callable, TypeVar

import grpc

from pgchangefeed.exceptions import (
    PgChangeFeedGrpcError,
    PgChangeFeedGrpcInternalError,
    PgChangeFeedGrpcInvalidArgumentError,
    PgChangeFeedGrpcNotFoundError,
    PgChangeFeedGrpcPermissionDeniedError,
    PgChangeFeedGrpcUnauthenticatedError,
    PgChangeFeedGrpcUnexpectedStatusError,
)
from pgchangefeed.grpc_gen import administration_pb2, administration_pb2_grpc
from pgchangefeed.options import ClientOptions

_AUTHORIZATION_METADATA_KEY = "authorization"
_BEARER_PREFIX = "Bearer "

_CODE_TO_ERROR: dict[grpc.StatusCode, type[PgChangeFeedGrpcError]] = {
    grpc.StatusCode.INVALID_ARGUMENT: PgChangeFeedGrpcInvalidArgumentError,
    grpc.StatusCode.UNAUTHENTICATED: PgChangeFeedGrpcUnauthenticatedError,
    grpc.StatusCode.PERMISSION_DENIED: PgChangeFeedGrpcPermissionDeniedError,
    grpc.StatusCode.NOT_FOUND: PgChangeFeedGrpcNotFoundError,
    grpc.StatusCode.INTERNAL: PgChangeFeedGrpcInternalError,
}

_T = TypeVar("_T")


class PgChangeFeedAdministrationClient:
    """Client for the eleven unary RPCs of the ``Administration`` gRPC service.

    ``channel`` is the ``grpc.Channel`` to the server's gRPC address;
    ``options`` carries the bearer token.
    """

    def __init__(self, channel: grpc.Channel, options: ClientOptions) -> None:
        self._client = administration_pb2_grpc.AdministrationStub(channel)
        self._options = options

    def register_consumer(
        self, request: administration_pb2.RegisterConsumerRequest
    ) -> administration_pb2.RegisterConsumerResponse:
        """Registers a consumer, a named reader whose position the server
        keeps (admin token). Registering an existing consumer changes
        nothing; the response reports it with ``already_registered``."""
        return self._call(self._client.RegisterConsumer, request)

    def acknowledge_consumer(
        self, request: administration_pb2.AcknowledgeConsumerRequest
    ) -> administration_pb2.AcknowledgeConsumerResponse:
        """Stores the position up to which a consumer has processed a source
        (admin token). Repeating the stored position has no effect; an
        earlier position, or a position of another source, raises
        ``PgChangeFeedGrpcInvalidArgumentError``."""
        return self._call(self._client.AcknowledgeConsumer, request)

    def get_consumer_position(
        self, request: administration_pb2.GetConsumerPositionRequest
    ) -> administration_pb2.GetConsumerPositionResponse:
        """Reads the stored position of a consumer (reader or admin token).
        ``acknowledged`` is ``False`` for a consumer that never acknowledged;
        ``offset`` is then the starting position."""
        return self._call(self._client.GetConsumerPosition, request)

    def remove_consumer(
        self, request: administration_pb2.RemoveConsumerRequest
    ) -> administration_pb2.RemoveConsumerResponse:
        """Removes a consumer (admin token); ``removed`` is ``False`` for one
        that was never registered."""
        return self._call(self._client.RemoveConsumer, request)

    def enable_table(
        self, request: administration_pb2.EnableTableRequest
    ) -> administration_pb2.EnableTableResponse:
        """Starts capturing a table of a source (admin token).
        ``already_enabled`` is ``True`` when the table was captured already; a
        table that does not exist in the source database raises
        ``PgChangeFeedGrpcNotFoundError``."""
        return self._call(self._client.EnableTable, request)

    def disable_table(
        self, request: administration_pb2.DisableTableRequest
    ) -> administration_pb2.DisableTableResponse:
        """Stops capturing a table (admin token). ``retained`` is ``True``
        when changes already stored for the table remain readable; a table
        that does not exist in the source database raises
        ``PgChangeFeedGrpcNotFoundError``."""
        return self._call(self._client.DisableTable, request)

    def get_table_status(
        self, request: administration_pb2.GetTableStatusRequest
    ) -> administration_pb2.GetTableStatusResponse:
        """Tells whether a table is captured (``enabled``) or no longer
        captured with stored changes remaining (``retained``) (reader or
        admin token). A table that was never enabled reports both as
        ``False``; a table that does not exist in the source database raises
        ``PgChangeFeedGrpcNotFoundError``."""
        return self._call(self._client.GetTableStatus, request)

    def list_tables(
        self, request: administration_pb2.ListTablesRequest
    ) -> administration_pb2.ListTablesResponse:
        """Lists the captured tables and the tables that are no longer
        captured but whose stored changes remain (reader or admin token)."""
        return self._call(self._client.ListTables, request)

    def run_retention(
        self, request: administration_pb2.RunRetentionRequest
    ) -> administration_pb2.RunRetentionResponse:
        """Deletes the stored changes of a source that are older than
        ``min_age_nanos`` and that every consumer with a stored position has
        passed (admin token); ``deleted`` is the number removed. A
        ``min_age_nanos`` of ``0`` means no minimum age and is valid."""
        return self._call(self._client.RunRetention, request)

    def read_changes(
        self, request: administration_pb2.ReadChangesRequest
    ) -> administration_pb2.ReadChangesResponse:
        """Reads a bounded range of stored changes of a source (reader or
        admin token) -- the same filter and range semantics as
        ``PgChangeFeedHttpClient.read_changes``: an unset ``from``/``to``/
        ``limit`` carries ``0`` (not set). The request's ``target`` selects the
        delivery target: empty (the default) is no filter, a set value returns
        only changes routed to that target, combined with ``schema``/``table``
        as a conjunction; a target no change carries returns an empty list."""
        return self._call(self._client.ReadChanges, request)

    def diagnose(
        self, request: administration_pb2.DiagnoseRequest
    ) -> administration_pb2.DiagnoseResponse:
        """Reads the operational diagnose report of a source (reader or
        admin token) -- the same report as the CLI diagnose mode and
        ``GET /diagnose``. A ``known``/``present``/``*_known`` field of
        ``False`` carries the respective absence case (no heartbeat ever
        written, no blocking consumer, unknown estimate/backlog)."""
        return self._call(self._client.Diagnose, request)

    def _metadata(self) -> tuple[tuple[str, str], ...]:
        return ((_AUTHORIZATION_METADATA_KEY, f"{_BEARER_PREFIX}{self._options.api_token}"),)

    def _call(self, method: Callable[..., _T], request: object) -> _T:
        try:
            return method(request, metadata=self._metadata())
        except grpc.RpcError as error:
            raise _map_error(error) from error


def _map_error(error: grpc.RpcError) -> PgChangeFeedGrpcError:
    code = error.code()
    message = error.details() or ""
    error_cls = _CODE_TO_ERROR.get(code, PgChangeFeedGrpcUnexpectedStatusError)
    return error_cls(code, message)
