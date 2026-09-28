using System.Net;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace CdcExamples.Http;

/// <summary>
/// ReadChangeResponse spiegelt einen Change in der Antwortform von
/// <c>GET /changes</c> (<c>LH-FA-SST-006</c>). Die Row Images stehen als
/// eingebettete JSON-Werte; ein fehlendes Bild wird zu <c>null</c>.
/// Form-Vorbild: <c>examples/http-client/changes.go</c>.
/// </summary>
public sealed record ReadChangeResponse(
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
    [property: JsonPropertyName("origin")] string Origin);

/// <summary>
/// ReadChangesResponse trägt den JSON-Response-Body bei Erfolg: ohne Treffer
/// eine leere, gesetzte Liste.
/// </summary>
public sealed record ReadChangesResponse(
    [property: JsonPropertyName("changes")] IReadOnlyList<ReadChangeResponse> Changes);

/// <summary>
/// ChangesUrlBuilder baut die Lese-Adresse von <c>GET /changes</c>:
/// <c>source</c> ist Pflicht, die übrigen fünf Parameter sind unabhängig
/// optional — ein leerer Wert bleibt weg statt als leerer Query-Parameter
/// zu erscheinen. <c>GET /changes</c> lässt nur einen unbekannten
/// Parameter**namen** mit <c>400</c> enden; ein leerer Wert eines bekannten
/// optionalen Parameters (<c>schema</c>, <c>table</c>, <c>from</c>,
/// <c>to</c>, <c>limit</c>) bleibt für den Server unberücksichtigt und
/// endet mit <c>200</c> — dieser Client sendet ihn trotzdem gar nicht
/// erst. Form-Vorbild: <c>ChangesURL</c> in
/// <c>examples/http-client/changes.go</c>.
/// </summary>
public static class ChangesUrlBuilder
{
    public static string Build(string addr, string source, string schema, string table, string from, string to, string limit)
    {
        var parameters = new SortedDictionary<string, string>(StringComparer.Ordinal) { ["source"] = source };
        AddIfNotEmpty(parameters, "schema", schema);
        AddIfNotEmpty(parameters, "table", table);
        AddIfNotEmpty(parameters, "from", from);
        AddIfNotEmpty(parameters, "to", to);
        AddIfNotEmpty(parameters, "limit", limit);

        var query = string.Join('&', parameters.Select(kv => $"{kv.Key}={Uri.EscapeDataString(kv.Value)}"));
        return $"http://{addr}/changes?{query}";
    }

    private static void AddIfNotEmpty(IDictionary<string, string> parameters, string key, string value)
    {
        if (!string.IsNullOrEmpty(value))
        {
            parameters[key] = value;
        }
    }
}

/// <summary>
/// ChangesClient ruft den <c>reader</c>-Endpunkt <c>GET /changes</c> auf.
/// Form-Vorbild: <c>examples/http-client/changes.go</c>.
/// </summary>
public static class ChangesClient
{
    public static Task<ReadChangesResponse> ReadChangesAsync(
        System.Net.Http.HttpClient httpClient, Config cfg, CancellationToken cancellationToken = default) =>
        RequestHelper.SendJsonAsync<ReadChangesResponse>(
            httpClient, HttpMethod.Get,
            ChangesUrlBuilder.Build(cfg.Addr, cfg.Source, cfg.Schema, cfg.Table, cfg.From, cfg.To, cfg.Limit),
            cfg.Token, null, HttpStatusCode.OK, cancellationToken);
}
