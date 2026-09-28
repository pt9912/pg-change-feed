using Cdc.Administration.V1;
using Cdc.Stream.V1;

namespace CdcExamples.Grpc;

/// <summary>
/// Format baut die Ausgabezeile einer empfangenen Stream-Nachricht sowie
/// die Ausgabezeilen jeder der elf Administration-Antworten — reine
/// Funktionen, netzlos testbar. Booleans drucken klein geschrieben
/// (<c>true</c>/<c>false</c>), dieselbe Wertform wie Gos <c>%v</c>. Form-Vorbild:
/// <c>examples/grpc-client/format.go</c>, <c>consumer.go</c>,
/// <c>tables_admin.go</c>, <c>changes.go</c>, <c>retention.go</c>,
/// <c>diagnose.go</c>.
/// </summary>
public static class Format
{
    /// <summary>
    /// FormatChange baut Tabelle, Operation und das neue Row Image, dazu die
    /// Kennung, an der sich die Nachricht gegen den Lesezugriffsweg
    /// <c>cdc.changes</c> halten lässt (<c>change_id</c>).
    /// </summary>
    public static string FormatChange(Change change)
    {
        return $"grpc-client: change_id={change.ChangeId} table={change.Schema}.{change.Table} " +
               $"operation={change.Operation} new_image={change.NewImage.ToStringUtf8()}";
    }

    public static string FormatRegisterConsumer(RegisterConsumerResponse resp) =>
        $"grpc-client: consumer_id={resp.ConsumerId} name={resp.Name} already_registered={Bool(resp.AlreadyRegistered)}";

    public static string FormatAcknowledgeConsumer(AcknowledgeConsumerResponse resp) =>
        $"grpc-client: consumer_id={resp.ConsumerId} source_id={resp.SourceId} offset={resp.Offset}";

    public static string FormatConsumerPosition(GetConsumerPositionResponse resp) =>
        $"grpc-client: consumer_id={resp.ConsumerId} source_id={resp.SourceId} offset={resp.Offset} acknowledged={Bool(resp.Acknowledged)}";

    public static string FormatRemoveConsumer(RemoveConsumerResponse resp) =>
        $"grpc-client: consumer_id={resp.ConsumerId} removed={Bool(resp.Removed)}";

    public static string FormatEnableTable(EnableTableResponse resp) =>
        $"grpc-client: table_id={resp.TableId} source={resp.Source} schema={resp.Schema} table={resp.Table} already_enabled={Bool(resp.AlreadyEnabled)}";

    public static string FormatDisableTable(DisableTableResponse resp) =>
        $"grpc-client: removed={Bool(resp.Removed)} retained={Bool(resp.Retained)}";

    public static string FormatTableStatus(GetTableStatusResponse resp) =>
        $"grpc-client: enabled={Bool(resp.Enabled)} retained={Bool(resp.Retained)}";

    /// <summary>
    /// FormatListTables listet jede Tabelle der Antwort — <c>tables</c> und
    /// <c>retained</c> je eine <c>SourceTable</c>-Zeile.
    /// </summary>
    public static string FormatListTables(ListTablesResponse resp)
    {
        var lines = new List<string> { $"grpc-client: tables={resp.Tables.Count} retained={resp.Retained.Count}" };
        foreach (var t in resp.Tables)
        {
            lines.Add($"  table: table_id={t.TableId} source={t.Source} schema={t.Schema} table={t.Table}");
        }
        foreach (var t in resp.Retained)
        {
            lines.Add($"  retained: table_id={t.TableId} source={t.Source} schema={t.Schema} table={t.Table}");
        }
        return string.Join('\n', lines);
    }

    public static string FormatRunRetention(RunRetentionResponse resp) =>
        $"grpc-client: deleted={resp.Deleted}";

    /// <summary>
    /// FormatReadChanges listet jede gelesene Änderung — dieselben Felder wie
    /// der Lesezugriffsweg <c>cdc.changes</c>.
    /// </summary>
    public static string FormatReadChanges(ReadChangesResponse resp)
    {
        var lines = new List<string> { $"grpc-client: changes={resp.Changes.Count}" };
        foreach (var c in resp.Changes)
        {
            lines.Add($"  change_id={c.ChangeId} table={c.Schema}.{c.Table} operation={c.Operation} " +
                      $"commit_position={c.CommitPosition} origin={c.Origin} new_image={c.NewImage.ToStringUtf8()}");
        }
        return string.Join('\n', lines);
    }

    /// <summary>
    /// FormatDiagnose baut denselben Bericht wie der CLI-Sondermodus: ein
    /// <c>known</c>/<c>present</c>-Feld auf <c>false</c> trägt den jeweiligen
    /// Abwesenheits-Fall (kein Lebenszeichen, kein Blocker, unbekannte
    /// Schätzung/unbekannter Rückstand).
    /// </summary>
    public static string FormatDiagnose(DiagnoseResponse resp)
    {
        var lines = new List<string>
        {
            $"grpc-client: heartbeat_known={Bool(resp.Heartbeat.Known)} heartbeat_age_seconds={resp.Heartbeat.AgeSeconds} " +
            $"heartbeat_error_class={resp.Heartbeat.ErrorClass} capture_lag={resp.CaptureLag} storage_bytes={resp.StorageBytes}",
        };

        foreach (var cl in resp.ConsumerLags)
        {
            lines.Add($"  consumer_lag: consumer_id={cl.ConsumerId} known={Bool(cl.Known)} lag={cl.Lag}");
        }

        if (resp.RetentionBlocker.Present)
        {
            var rb = resp.RetentionBlocker;
            lines.Add($"  retention_blocker: consumer_id={rb.ConsumerId} name={rb.Name} " +
                      $"acknowledged_position={rb.AcknowledgedPosition} backlog_known={Bool(rb.BacklogKnown)} backlog={rb.Backlog}");
        }
        else
        {
            lines.Add("  retention_blocker: kein Blocker");
        }

        foreach (var b in resp.Backfill)
        {
            lines.Add($"  backfill: schema={b.Schema} table={b.Table} status={b.Status} rows_copied={b.RowsCopied} " +
                      $"estimated_rows_known={Bool(b.EstimatedRowsKnown)} estimated_rows={b.EstimatedRows} " +
                      $"warn_size={Bool(b.WarnEstimatedSize)} warn_duration={Bool(b.WarnDuration)} error_message={b.ErrorMessage}");
        }

        return string.Join('\n', lines);
    }

    private static string Bool(bool value) => value ? "true" : "false";
}
