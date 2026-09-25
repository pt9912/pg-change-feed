"""Typed request/response data classes mirroring the SPEC-018/SPEC-022 JSON
schemas exactly.

Field names and types are taken directly from spec/pflichtenheft.md
SPEC-018 (HTTP-API: Endpunkte und Token-Header-Form) and SPEC-022
(HTTP-API: Changes lesen, GET /changes) -- not from the C# sibling
package, which serves only as a structural reference, not a wire
reference (ADR-0107). Every attribute name equals the JSON field name
verbatim, so no case-mapping layer is needed between the two.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any


# --- RegisterConsumer -- POST /consumers (admin, LH-FA-CON-001) ---


@dataclass(frozen=True)
class RegisterConsumerRequest:
    consumer_id: str
    name: str


@dataclass(frozen=True)
class RegisterConsumerResponse:
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


# --- AcknowledgeConsumer -- POST /consumers/acknowledge (admin, LH-FA-CON-004) ---


@dataclass(frozen=True)
class AcknowledgeConsumerRequest:
    consumer_id: str
    source_id: str
    offset: int


@dataclass(frozen=True)
class AcknowledgeConsumerResponse:
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


# --- GetConsumerPosition -- GET /consumers/position (reader|admin, LH-FA-CON-003/005) ---


@dataclass(frozen=True)
class ConsumerPositionResponse:
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


# --- RemoveConsumer -- POST /consumers/remove (admin, LH-FA-CON-006) ---


@dataclass(frozen=True)
class RemoveConsumerResponse:
    consumer_id: str
    removed: bool

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> RemoveConsumerResponse:
        return cls(consumer_id=data["consumer_id"], removed=data["removed"])


# --- EnableTable -- POST /tables/enable (admin, LH-FA-CFG-001) ---


@dataclass(frozen=True)
class EnableTableRequest:
    source: str
    schema: str
    table: str
    table_id: str
    schema_version_id: str
    version: int
    publication: str


@dataclass(frozen=True)
class EnableTableResponse:
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


# --- DisableTable -- POST /tables/disable (admin, LH-FA-CFG-002) ---


@dataclass(frozen=True)
class DisableTableRequest:
    source: str
    schema: str
    table: str
    publication: str


@dataclass(frozen=True)
class DisableTableResponse:
    removed: bool
    retained: bool

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> DisableTableResponse:
        return cls(removed=data["removed"], retained=data["retained"])


# --- GetStatus -- GET /tables/status (reader|admin, LH-FA-CFG-003) ---


@dataclass(frozen=True)
class TableStatusResponse:
    enabled: bool
    retained: bool

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> TableStatusResponse:
        return cls(enabled=data["enabled"], retained=data["retained"])


# --- ListTables -- GET /tables (reader|admin, LH-FA-CFG-004) ---


@dataclass(frozen=True)
class TableInfo:
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
    tables: list[TableInfo] = field(default_factory=list)
    retained: list[TableInfo] = field(default_factory=list)

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> ListTablesResponse:
        return cls(
            tables=[TableInfo.from_json(item) for item in data["tables"]],
            retained=[TableInfo.from_json(item) for item in data["retained"]],
        )


# --- RunRetention -- POST /retention/run (admin, LH-FA-RET-002..004) ---


@dataclass(frozen=True)
class RunRetentionRequest:
    source: str
    min_age_nanos: int


@dataclass(frozen=True)
class RunRetentionResponse:
    deleted: int

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> RunRetentionResponse:
        return cls(deleted=data["deleted"])


# --- ReadChanges -- GET /changes (reader|admin, SPEC-022) ---


@dataclass(frozen=True)
class Change:
    """One persisted change as returned by ``GET /changes`` (SPEC-022).

    ``origin`` is ``wal`` for a change captured from the replication stream
    and ``backfill`` for an existing-rows change (LH-FA-CAP-009). It is the
    server's string (an empty or unknown value included), and a response
    without the field or with a JSON ``null`` reads as ``wal``. The live
    surfaces (gRPC, SSE, NATS) carry no ``origin``.
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
    changes: list[Change] = field(default_factory=list)

    @classmethod
    def from_json(cls, data: dict[str, Any]) -> ReadChangesResponse:
        return cls(changes=[Change.from_json(item) for item in data["changes"]])


# --- Live change stream, SSE (SPEC-021): the ten fields the live
# --- surfaces carry, not the thirteen of the HTTP read (SPEC-022) --


@dataclass(frozen=True)
class StreamChange:
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
