using System.Net;
using System.Text.Json.Serialization;

namespace CdcExamples.Http;

/// <summary>
/// EnableTableRequest/-Response spiegeln <c>POST /tables/enable</c>
/// (<c>LH-FA-CFG-001</c>): alle sieben Request-Felder Pflicht,
/// <c>AlreadyEnabled</c> trägt die Idempotenz-Antwort ohne eigenen
/// Statuscode. Form-Vorbild: <c>examples/http-client/tables_admin.go</c>.
/// </summary>
public sealed record EnableTableRequest(
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table,
    [property: JsonPropertyName("table_id")] string TableId,
    [property: JsonPropertyName("schema_version_id")] string SchemaVersionId,
    [property: JsonPropertyName("version")] long Version,
    [property: JsonPropertyName("publication")] string Publication);

public sealed record EnableTableResponse(
    [property: JsonPropertyName("table_id")] string TableId,
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table,
    [property: JsonPropertyName("already_enabled")] bool AlreadyEnabled);

/// <summary>
/// DisableTableRequest/-Response spiegeln <c>POST /tables/disable</c>
/// (<c>LH-FA-CFG-002</c>).
/// </summary>
public sealed record DisableTableRequest(
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table,
    [property: JsonPropertyName("publication")] string Publication);

public sealed record DisableTableResponse(
    [property: JsonPropertyName("removed")] bool Removed,
    [property: JsonPropertyName("retained")] bool Retained);

/// <summary>
/// TableStatusResponse spiegelt <c>GET /tables/status</c>
/// (<c>LH-FA-CFG-003</c>): <c>Enabled</c>/<c>Retained</c> trennen
/// Erfassungs-Zustand und Herkunft.
/// </summary>
public sealed record TableStatusResponse(
    [property: JsonPropertyName("enabled")] bool Enabled,
    [property: JsonPropertyName("retained")] bool Retained);

/// <summary>
/// TableStatusUrlBuilder baut die Lese-Adresse von <c>GET /tables/status</c>:
/// alle vier Parameter sind Pflicht. Form-Vorbild: <c>TableStatusURL</c> in
/// <c>examples/http-client/tables_admin.go</c>.
/// </summary>
public static class TableStatusUrlBuilder
{
    public static string Build(string addr, string source, string schema, string table, string publication)
    {
        var query = string.Join(
            '&',
            $"publication={Uri.EscapeDataString(publication)}",
            $"schema={Uri.EscapeDataString(schema)}",
            $"source={Uri.EscapeDataString(source)}",
            $"table={Uri.EscapeDataString(table)}");
        return $"http://{addr}/tables/status?{query}";
    }
}

/// <summary>
/// TablesAdminClient ruft die drei Tabellen-Verwaltungs-Fähigkeiten der
/// HTTP-/JSON-API auf, die <see cref="TablesClient"/> (Auflistung) nicht
/// trägt. Form-Vorbild: <c>examples/http-client/tables_admin.go</c>.
/// </summary>
public static class TablesAdminClient
{
    public static Task<EnableTableResponse> EnableTableAsync(
        System.Net.Http.HttpClient httpClient, Config cfg, CancellationToken cancellationToken = default) =>
        RequestHelper.SendJsonAsync<EnableTableResponse>(
            httpClient, HttpMethod.Post, $"http://{cfg.Addr}/tables/enable", cfg.AdminToken,
            new EnableTableRequest(cfg.Source, cfg.Schema, cfg.Table, cfg.TableId, cfg.SchemaVersionId, cfg.Version, cfg.Publication),
            HttpStatusCode.Created, cancellationToken);

    public static Task<DisableTableResponse> DisableTableAsync(
        System.Net.Http.HttpClient httpClient, Config cfg, CancellationToken cancellationToken = default) =>
        RequestHelper.SendJsonAsync<DisableTableResponse>(
            httpClient, HttpMethod.Post, $"http://{cfg.Addr}/tables/disable", cfg.AdminToken,
            new DisableTableRequest(cfg.Source, cfg.Schema, cfg.Table, cfg.Publication),
            HttpStatusCode.OK, cancellationToken);

    public static Task<TableStatusResponse> TableStatusAsync(
        System.Net.Http.HttpClient httpClient, Config cfg, CancellationToken cancellationToken = default) =>
        RequestHelper.SendJsonAsync<TableStatusResponse>(
            httpClient, HttpMethod.Get, TableStatusUrlBuilder.Build(cfg.Addr, cfg.Source, cfg.Schema, cfg.Table, cfg.Publication),
            cfg.Token, null, HttpStatusCode.OK, cancellationToken);
}
