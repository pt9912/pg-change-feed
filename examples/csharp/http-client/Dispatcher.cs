using System.Text.Json;

namespace CdcExamples.Http;

/// <summary>
/// Dispatcher ruft die zum Verb gehörende Fähigkeit auf und liefert den
/// Antwort-Body als druckbaren Text. <c>tables</c> liest die rohe
/// Server-Antwort (<see cref="TablesClient"/>); die übrigen neun Fähigkeiten
/// drucken ihre typisierte Antwort eingerückt. Form-Vorbild: <c>dispatch</c>
/// in <c>examples/http-client/main.go</c>. Ein unbekanntes Verb bricht hier
/// ein zweites Mal ab (verteidigend, für den Fall eines Aufrufs ohne
/// vorherige <see cref="Validator.Validate"/>).
/// </summary>
public static class Dispatcher
{
    private static readonly JsonSerializerOptions IndentedOptions = new() { WriteIndented = true };

    public static async Task<string> DispatchAsync(
        System.Net.Http.HttpClient httpClient, Config cfg, CancellationToken cancellationToken = default)
    {
        switch (cfg.Verb)
        {
            case "tables":
                return await TablesClient.ListTablesAsync(httpClient, cfg, cancellationToken).ConfigureAwait(false);
            case "changes":
                return SerializeIndented(await ChangesClient.ReadChangesAsync(httpClient, cfg, cancellationToken).ConfigureAwait(false));
            case "register-consumer":
                return SerializeIndented(await ConsumerClient.RegisterConsumerAsync(httpClient, cfg, cancellationToken).ConfigureAwait(false));
            case "acknowledge":
                return SerializeIndented(await ConsumerClient.AcknowledgeConsumerAsync(httpClient, cfg, cancellationToken).ConfigureAwait(false));
            case "consumer-position":
                return SerializeIndented(await ConsumerClient.ConsumerPositionAsync(httpClient, cfg, cancellationToken).ConfigureAwait(false));
            case "remove-consumer":
                return SerializeIndented(await ConsumerClient.RemoveConsumerAsync(httpClient, cfg, cancellationToken).ConfigureAwait(false));
            case "enable-table":
                var enableCfg = ApplyEnableTableDefaults(cfg);
                return SerializeIndented(await TablesAdminClient.EnableTableAsync(httpClient, enableCfg, cancellationToken).ConfigureAwait(false));
            case "disable-table":
                return SerializeIndented(await TablesAdminClient.DisableTableAsync(httpClient, cfg, cancellationToken).ConfigureAwait(false));
            case "table-status":
                return SerializeIndented(await TablesAdminClient.TableStatusAsync(httpClient, cfg, cancellationToken).ConfigureAwait(false));
            case "retention-run":
                return SerializeIndented(await RetentionClient.RunRetentionAsync(httpClient, cfg, cancellationToken).ConfigureAwait(false));
            default:
                throw new InvalidOperationException($"unbekanntes --verb \"{cfg.Verb}\"");
        }
    }

    /// <summary>
    /// ApplyEnableTableDefaults leitet <c>TableId</c>/<c>SchemaVersionId</c>
    /// aus Schema und Tabelle ab, wenn sie nicht gesetzt sind — dieselbe
    /// Default-Konvention wie <c>internal/bootstrap/wiring.go</c>s
    /// <c>administrationTableID</c>/<c>administrationSchemaVersionID</c>.
    /// </summary>
    private static Config ApplyEnableTableDefaults(Config cfg)
    {
        var tableId = string.IsNullOrEmpty(cfg.TableId) ? $"{cfg.Schema}.{cfg.Table}" : cfg.TableId;
        var schemaVersionId = string.IsNullOrEmpty(cfg.SchemaVersionId) ? $"{tableId}-v1" : cfg.SchemaVersionId;
        return cfg with { TableId = tableId, SchemaVersionId = schemaVersionId };
    }

    private static string SerializeIndented<T>(T value) => JsonSerializer.Serialize(value, IndentedOptions);
}
