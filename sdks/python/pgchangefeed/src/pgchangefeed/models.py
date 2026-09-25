"""Typed request and response data classes of the HTTP API and the live streams.

Every attribute name equals the JSON field name verbatim, so no case-mapping
layer sits between the wire and these classes. Requests are plain frozen data
classes; responses add a ``from_json`` constructor that raises ``KeyError`` or
``TypeError`` for a document that lacks an expected field (the clients turn
that into ``PgChangeFeedMalformedResponseError``).
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any


# --- register_consumer -- POST /consumers ---


@dataclass(frozen=True)
class RegisterConsumerRequest:
    """The consumer to register: a unique ``consumer_id`` and a display ``name``."""

    consumer_id: str
    name: str


@dataclass(frozen=True)
class RegisterConsumerResponse:
    """The registered consumer; ``already_registered`` is true when it existed before."""

    consumer_id: str
    name: str
    already_registered: bool

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> RegisterConsumerResponse:
        return cls(
            consumer_id=data["consumer_id"],
            name=data["name"],
            already_registered=data["already_registered"],
        )


# --- acknowledge_consumer -- POST /consumers/acknowledge ---


@dataclass(frozen=True)
class AcknowledgeConsumerRequest:
    """The position to store: ``offset`` is the ``commit_position`` of the last
    change the consumer has processed for ``source_id``."""

    consumer_id: str
    source_id: str
    offset: int


@dataclass(frozen=True)
class AcknowledgeConsumerResponse:
    """The position now stored for the consumer."""

    consumer_id: str
    source_id: str
    offset: int

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> AcknowledgeConsumerResponse:
        return cls(
            consumer_id=data["consumer_id"],
            source_id=data["source_id"],
            offset=data["offset"],
        )


# --- get_consumer_position -- GET /consumers/position ---


@dataclass(frozen=True)
class ConsumerPositionResponse:
    """The stored position of a consumer; ``acknowledged`` is false for a
    consumer that has never acknowledged (``offset`` is then the starting
    position)."""

    consumer_id: str
    source_id: str
    offset: int
    acknowledged: bool

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> ConsumerPositionResponse:
        return cls(
            consumer_id=data["consumer_id"],
            source_id=data["source_id"],
            offset=data["offset"],
            acknowledged=data["acknowledged"],
        )


# --- remove_consumer -- POST /consumers/remove ---


@dataclass(frozen=True)
class RemoveConsumerResponse:
    """``removed`` is false for a consumer that was never registered."""

    consumer_id: str
    removed: bool

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> RemoveConsumerResponse:
        return cls(consumer_id=data["consumer_id"], removed=data["removed"])


# --- enable_table -- POST /tables/enable ---


@dataclass(frozen=True)
class EnableTableRequest:
    """The table to capture.

    ``source`` and ``publication`` name the source and its PostgreSQL
    publication. ``table_id`` is the id the table is captured under (it appears
    as ``source_table_id`` on every change of the table); ``schema_version_id``
    and ``version`` (1 or higher) identify the table schema version the changes
    are recorded with (it appears as ``schema_version``).
    """

    source: str
    schema: str
    table: str
    table_id: str
    schema_version_id: str
    version: int
    publication: str


@dataclass(frozen=True)
class EnableTableResponse:
    """The captured table; ``already_enabled`` is true when it was captured before."""

    table_id: str
    source: str
    schema: str
    table: str
    already_enabled: bool

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> EnableTableResponse:
        return cls(
            table_id=data["table_id"],
            source=data["source"],
            schema=data["schema"],
            table=data["table"],
            already_enabled=data["already_enabled"],
        )


# --- disable_table -- POST /tables/disable ---


@dataclass(frozen=True)
class DisableTableRequest:
    """The table to stop capturing, in the given source and publication."""

    source: str
    schema: str
    table: str
    publication: str


@dataclass(frozen=True)
class DisableTableResponse:
    """``removed`` reports that the table is no longer captured; ``retained`` is
    true when changes already stored for it remain readable."""

    removed: bool
    retained: bool

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> DisableTableResponse:
        return cls(removed=data["removed"], retained=data["retained"])


# --- get_status -- GET /tables/status ---


@dataclass(frozen=True)
class TableStatusResponse:
    """``enabled`` is true for a captured table; ``retained`` is true for a
    table that is no longer captured but whose stored changes remain. Both are
    false for a table that was never enabled."""

    enabled: bool
    retained: bool

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> TableStatusResponse:
        return cls(enabled=data["enabled"], retained=data["retained"])


# --- list_tables -- GET /tables ---


@dataclass(frozen=True)
class TableInfo:
    """One table of a ``ListTablesResponse``."""

    table_id: str
    source: str
    schema: str
    table: str

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> TableInfo:
        return cls(
            table_id=data["table_id"],
            source=data["source"],
            schema=data["schema"],
            table=data["table"],
        )


@dataclass(frozen=True)
class ListTablesResponse:
    """``tables`` are the captured tables; ``retained`` are the tables that are
    no longer captured but whose stored changes remain."""

    tables: list[TableInfo] = field(default_factory=list)
    retained: list[TableInfo] = field(default_factory=list)

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> ListTablesResponse:
        return cls(
            tables=[TableInfo.from_json(item) for item in data["tables"]],
            retained=[TableInfo.from_json(item) for item in data["retained"]],
        )


# --- run_retention -- POST /retention/run ---


@dataclass(frozen=True)
class RunRetentionRequest:
    """Deletes the changes of ``source`` that are older than ``min_age_nanos``
    (0 or higher) and that every consumer with a stored position has passed."""

    source: str
    min_age_nanos: int


@dataclass(frozen=True)
class RunRetentionResponse:
    """``deleted`` is the number of changes removed."""

    deleted: int

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> RunRetentionResponse:
        return cls(deleted=data["deleted"])


# --- read_changes -- GET /changes ---


@dataclass(frozen=True)
class Change:
    """One stored change as returned by ``read_changes``.

    ``origin`` is ``wal`` for a change captured live from the database and
    ``backfill`` for a change taken from the existing table contents. It is
    the server's string (an empty or unknown value included), and a response
    without the field or with a JSON ``null`` reads as ``wal``. The live
    streams (gRPC, SSE, NATS) carry no ``origin``.
    """

    commit_position: int
    change_id: str
    transaction_id: str
    source_table_id: str
    schema: str
    table: str
    sequence: int
    operation: str
    old_image: Any | None
    new_image: Any | None
    schema_version: str
    committed_at: str
    origin: str = "wal"

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> Change:
        return cls(
            commit_position=data["commit_position"],
            change_id=data["change_id"],
            transaction_id=data["transaction_id"],
            source_table_id=data["source_table_id"],
            schema=data["schema"],
            table=data["table"],
            sequence=data["sequence"],
            operation=data["operation"],
            old_image=data["old_image"],
            new_image=data["new_image"],
            schema_version=data["schema_version"],
            committed_at=data["committed_at"],
            origin="wal" if data.get("origin") is None else data["origin"],
        )


@dataclass(frozen=True)
class ReadChangesResponse:
    """The changes of the requested range in a fixed order; empty when none match."""

    changes: list[Change] = field(default_factory=list)

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> ReadChangesResponse:
        return cls(changes=[Change.from_json(item) for item in data["changes"]])


# --- live change streams (SSE and NATS) --


@dataclass(frozen=True)
class StreamChange:
    """One change delivered by a live stream (SSE or NATS).

    A stream change has ten fields; unlike ``Change`` it carries no
    ``commit_position``, ``committed_at`` or ``origin``. The row images are
    JSON values; a missing image is ``None``.
    """

    change_id: str
    transaction_id: str
    source_table_id: str
    sequence: int
    operation: str
    old_image: Any | None
    new_image: Any | None
    schema_version: str
    schema: str
    table: str

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> StreamChange:
        return cls(
            change_id=data["change_id"],
            transaction_id=data["transaction_id"],
            source_table_id=data["source_table_id"],
            sequence=data["sequence"],
            operation=data["operation"],
            old_image=data["old_image"],
            new_image=data["new_image"],
            schema_version=data["schema_version"],
            schema=data["schema"],
            table=data["table"],
        )
