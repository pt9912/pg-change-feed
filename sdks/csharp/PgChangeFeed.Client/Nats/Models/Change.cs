using System.Text.Json;
using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Nats.Models;

/// <summary>
/// One change delivered on the NATS stream — the ten fields the server writes
/// into the JSON payload of a message on
/// <c>cdc.stream.&lt;source_id&gt;.&lt;schema&gt;.&lt;table&gt;</c>, one JSON
/// object per change. <see cref="OldImage"/>/<see cref="NewImage"/> carry the
/// row image as an embedded JSON value or <c>null</c> when absent; kept as
/// <see cref="JsonElement"/> rather than a fixed shape, because the row
/// image's own shape depends on the captured table (same as
/// <see cref="PgChangeFeed.Client.Sse.Models.Change"/>).
///
/// The ten fields are the same as those of the SSE event, but the type is its
/// own: the HTTP read returns <see cref="PgChangeFeed.Client.Http.Models.Change"/>
/// (thirteen fields including <c>commit_position</c>/<c>committed_at</c>/<c>origin</c>),
/// the gRPC client returns the generated <c>Change</c> message, the SSE client
/// returns <see cref="PgChangeFeed.Client.Sse.Models.Change"/> — four clients,
/// four independent types that share most field names.
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
