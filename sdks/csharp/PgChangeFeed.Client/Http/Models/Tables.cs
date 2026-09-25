using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Http.Models;

/// <summary>
/// The table to capture, all seven fields mandatory. <c>Source</c> and
/// <c>Publication</c> name the source and its PostgreSQL publication.
/// <c>TableId</c> is the id the table is captured under (every change of the
/// table carries it as <c>source_table_id</c>); <c>SchemaVersionId</c> and
/// <c>Version</c> (must be &gt;= 1) identify the table's schema version (changes
/// carry the id as <c>schema_version</c>).
/// </summary>
public sealed record EnableTableRequest(
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table,
    [property: JsonPropertyName("table_id")] string TableId,
    [property: JsonPropertyName("schema_version_id")] string SchemaVersionId,
    [property: JsonPropertyName("version")] long Version,
    [property: JsonPropertyName("publication")] string Publication);

/// <summary>
/// The captured table; <c>already_enabled</c> is true when it was captured
/// before. A table that does not exist in the source database ends <c>404</c>
/// instead.
/// </summary>
public sealed record EnableTableResponse(
    [property: JsonPropertyName("table_id")] string TableId,
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table,
    [property: JsonPropertyName("already_enabled")] bool AlreadyEnabled);

/// <summary>
/// The table to stop capturing, all four fields mandatory.
/// </summary>
public sealed record DisableTableRequest(
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table,
    [property: JsonPropertyName("publication")] string Publication);

/// <summary>
/// <c>Retained</c> is true when changes already stored for the table remain
/// readable. A table that does not exist in the source database ends
/// <c>404</c> instead.
/// </summary>
public sealed record DisableTableResponse(
    [property: JsonPropertyName("removed")] bool Removed,
    [property: JsonPropertyName("retained")] bool Retained);

/// <summary>
/// <c>Enabled</c> is true for a captured table; <c>Retained</c> is true for a
/// table that is no longer captured but whose stored changes remain. Both
/// <c>false</c> reads a table that was never enabled. A table that does not
/// exist in the source database ends <c>404</c> instead.
/// </summary>
public sealed record TableStatusResponse(
    [property: JsonPropertyName("enabled")] bool Enabled,
    [property: JsonPropertyName("retained")] bool Retained);

/// <summary>One table of a <see cref="ListTablesResponse"/>.</summary>
public sealed record TableInfo(
    [property: JsonPropertyName("table_id")] string TableId,
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table);

/// <summary>
/// <c>Tables</c> are the captured tables; <c>Retained</c> are the tables that
/// are no longer captured but whose stored changes remain. Both lists are empty
/// when no table was ever enabled.
/// </summary>
public sealed record ListTablesResponse(
    [property: JsonPropertyName("tables")] IReadOnlyList<TableInfo> Tables,
    [property: JsonPropertyName("retained")] IReadOnlyList<TableInfo> Retained);
