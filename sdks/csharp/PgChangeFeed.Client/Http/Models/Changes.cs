using System.Text.Json;
using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Http.Models;

/// <summary>
/// One persisted change as returned by <c>GET /changes</c> (SPEC-022) — thirteen
/// fields: the ten the live surfaces (gRPC, SSE, NATS) carry, plus
/// <c>commit_position</c>, <c>committed_at</c> and <c>origin</c>.
/// <see cref="OldImage"/>/<see cref="NewImage"/> carry the row image as an
/// embedded JSON value or <c>null</c> when absent; kept as
/// <see cref="JsonElement"/> rather than a fixed shape, because the row
/// image's own shape depends on the captured table, not on this wire contract.
/// <see cref="Origin"/> is <c>wal</c> for a change captured from the
/// replication stream and <c>backfill</c> for an existing-rows change
/// (LH-FA-CAP-009); it is carried as the server's string (an empty or
/// unknown value included), and a response without the field or with a JSON
/// <c>null</c> reads as <c>wal</c>. The live surfaces carry no <c>origin</c>.
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
    [property: JsonPropertyName("committed_at")] string CommittedAt,
    [property: JsonPropertyName("origin"), JsonConverter(typeof(OriginConverter))] string Origin = "wal");

/// <summary>
/// Reads the <c>origin</c> string as the server sent it and a JSON
/// <c>null</c> as <c>wal</c>; <c>System.Text.Json</c> skips a converter for
/// <c>null</c> unless <see cref="HandleNull"/> is set.
/// </summary>
internal sealed class OriginConverter : JsonConverter<string>
{
    public override bool HandleNull => true;

    public override string Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options) =>
        reader.TokenType == JsonTokenType.Null ? "wal" : reader.GetString()!;

    public override void Write(Utf8JsonWriter writer, string value, JsonSerializerOptions options) =>
        writer.WriteStringValue(value);
}

/// <summary>
/// <c>ReadChanges</c> response (<c>200</c>) — an empty list on no match,
/// never <c>404</c> (SPEC-022).
/// </summary>
public sealed record ReadChangesResponse(
    [property: JsonPropertyName("changes")] IReadOnlyList<Change> Changes);
