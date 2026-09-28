using Cdc.Administration.V1;

namespace CdcExamples.Grpc;

/// <summary>
/// Dispatcher ruft die zum Verb gehörende Administration-RPC auf und
/// liefert die formatierte Ausgabe. <c>stream</c> läuft nicht hier durch —
/// er hat einen eigenen, dauerhaften Ablauf (<c>Program.RunStreamAsync</c>).
/// Ein unbekanntes Verb bricht hier ein zweites Mal ab (verteidigend, für
/// den Fall eines Aufrufs ohne vorherige <see cref="Validator.Validate"/>).
/// Form-Vorbild: <c>dispatchAdmin</c> in <c>examples/grpc-client/main.go</c>.
/// </summary>
public static class Dispatcher
{
    public static Task<string> DispatchAdminAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        return cfg.Verb switch
        {
            "register-consumer" => ConsumerClient.RegisterConsumerAsync(client, cfg, cancellationToken),
            "acknowledge-consumer" => ConsumerClient.AcknowledgeConsumerAsync(client, cfg, cancellationToken),
            "get-consumer-position" => ConsumerClient.GetConsumerPositionAsync(client, cfg, cancellationToken),
            "remove-consumer" => ConsumerClient.RemoveConsumerAsync(client, cfg, cancellationToken),
            "enable-table" => TablesAdminClient.EnableTableAsync(client, ApplyEnableTableDefaults(cfg), cancellationToken),
            "disable-table" => TablesAdminClient.DisableTableAsync(client, cfg, cancellationToken),
            "get-table-status" => TablesAdminClient.GetTableStatusAsync(client, cfg, cancellationToken),
            "list-tables" => TablesAdminClient.ListTablesAsync(client, cfg, cancellationToken),
            "run-retention" => RetentionClient.RunRetentionAsync(client, cfg, cancellationToken),
            "read-changes" => ChangesClient.ReadChangesAsync(client, cfg, cancellationToken),
            "diagnose" => DiagnoseClient.DiagnoseAsync(client, cfg, cancellationToken),
            _ => throw new InvalidOperationException($"unbekanntes --verb \"{cfg.Verb}\""),
        };
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
}
