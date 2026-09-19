using System.Text.Json;
using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Http.Models;

/// <summary>
/// One persisted change as returned by <c>GET /changes</c> (SPEC-022) — the
/// same ten fields as the domain type <c>model.Change</c>. <see cref="OldImage"/>/
/// <see cref="NewImage"/> carry the row image as an embedded JSON value or
/// <c>null</c> when absent; kept as <see cref="JsonElement"/> rather than a
/// fixed shape, because the row image's own shape depends on the captured
/// table, not on this wire contract.
/// </summary>
public sealed record Change(
    [property: JsonPropertyName("commit_position")] long CommitPosition,
    [property: JsonPropertyName("change_id")] string ChangeId,
    [property: JsonPropertyName("transaction_id")] string TransactionId,
    [property: JsonPropertyName("source_table_id")] string SourceTableId,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table,
    [property: JsonPropertyName("sequence")] long Sequence,
    [property: JsonPropertyName("operation")] string Operation,
    [property: JsonPropertyName("old_image")] JsonElement? OldImage,
    [property: JsonPropertyName("new_image")] JsonElement? NewImage,
    [property: JsonPropertyName("schema_version")] string SchemaVersion,
    [property: JsonPropertyName("committed_at")] string CommittedAt);

/// <summary>
/// <c>ReadChanges</c> response (<c>200</c>) — an empty list on no match,
/// never <c>404</c> (SPEC-022).
/// </summary>
public sealed record ReadChangesResponse(
    [property: JsonPropertyName("changes")] IReadOnlyList<Change> Changes);
