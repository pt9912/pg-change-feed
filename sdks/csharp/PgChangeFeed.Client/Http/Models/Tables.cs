using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Http.Models;

/// <summary>
/// <c>EnableTable</c> request — <c>POST /tables/enable</c> (SPEC-018), all
/// seven fields mandatory, <c>Version</c> must be &gt;= 1.
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
/// <c>EnableTable</c> response (<c>201</c>). A table physically missing at
/// the source ends <c>404</c> instead (SPEC-018).
/// </summary>
public sealed record EnableTableResponse(
    [property: JsonPropertyName("table_id")] string TableId,
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table,
    [property: JsonPropertyName("already_enabled")] bool AlreadyEnabled);

/// <summary>
/// <c>DisableTable</c> request — <c>POST /tables/disable</c> (SPEC-018), all
/// four fields mandatory.
/// </summary>
public sealed record DisableTableRequest(
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table,
    [property: JsonPropertyName("publication")] string Publication);

/// <summary>
/// <c>DisableTable</c> response (<c>200</c>). A table physically missing at
/// the source ends <c>404</c> instead (SPEC-018).
/// </summary>
public sealed record DisableTableResponse(
    [property: JsonPropertyName("removed")] bool Removed,
    [property: JsonPropertyName("retained")] bool Retained);

/// <summary>
/// <c>GetStatus</c> response (<c>200</c>) — both <c>false</c> reads a table
/// that was never enabled. A table physically missing at the source ends
/// <c>404</c> instead (SPEC-018).
/// </summary>
public sealed record TableStatusResponse(
    [property: JsonPropertyName("enabled")] bool Enabled,
    [property: JsonPropertyName("retained")] bool Retained);

/// <summary>One entry of a <see cref="ListTablesResponse"/> list (SPEC-018).</summary>
public sealed record TableInfo(
    [property: JsonPropertyName("table_id")] string TableId,
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table);

/// <summary>
/// <c>ListTables</c> response (<c>200</c>) — both lists are empty without
/// any activation (SPEC-018).
/// </summary>
public sealed record ListTablesResponse(
    [property: JsonPropertyName("tables")] IReadOnlyList<TableInfo> Tables,
    [property: JsonPropertyName("retained")] IReadOnlyList<TableInfo> Retained);
