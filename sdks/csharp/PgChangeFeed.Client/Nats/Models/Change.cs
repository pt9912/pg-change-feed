using System.Text.Json;
using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Nats.Models;

/// <summary>
/// One change delivered on the SPEC-024 NATS full-content stream — the ten
/// fields the publisher writes into the JSON payload of a message on
/// <c>cdc.stream.&lt;source_id&gt;.&lt;schema&gt;.&lt;table&gt;</c>, one JSON
/// object per change. <see cref="OldImage"/>/<see cref="NewImage"/> carry the
/// row image as an embedded JSON value or <c>null</c> when absent; kept as
/// <see cref="JsonElement"/> rather than a fixed shape, because the row
/// image's own shape depends on the captured table, not on this wire contract
/// (same reasoning as <see cref="PgChangeFeed.Client.Sse.Models.Change"/>).
///
/// SPEC-024 documents the same ten fields as SPEC-021's SSE event — "dasselbe
/// Nachrichtenschema, kein drittes" — but this type is still its own,
/// consistent with <see cref="PgChangeFeed.Client.Sse.Models.Change"/>'s
/// established reasoning: it is one of four independent wire contracts that
/// happen to share most field names — <see cref="PgChangeFeed.Client.Http.Models.Change"/>
/// (SPEC-022, eleven fields including <c>commit_position</c>/<c>committed_at</c>),
/// the generated gRPC <c>Change</c> stub (SPEC-020, its own protobuf-generated
/// type), <see cref="PgChangeFeed.Client.Sse.Models.Change"/> (SPEC-021), and
/// this type (SPEC-024) — reusing one across surfaces would only risk
/// drifting one contract's shape into another's.
/// </summary>
public sealed record Change(
    [property: JsonPropertyName("change_id")] string ChangeId,
    [property: JsonPropertyName("transaction_id")] string TransactionId,
    [property: JsonPropertyName("source_table_id")] string SourceTableId,
    [property: JsonPropertyName("sequence")] long Sequence,
    [property: JsonPropertyName("operation")] string Operation,
    [property: JsonPropertyName("old_image")] JsonElement? OldImage,
    [property: JsonPropertyName("new_image")] JsonElement? NewImage,
    [property: JsonPropertyName("schema_version")] string SchemaVersion,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table);
