using Cdc.Administration.V1;
using Cdc.Stream.V1;
using Google.Protobuf;
using Xunit;

namespace CdcExamples.Grpc.Tests;

/// <summary>
/// Prüft <see cref="Format"/>: die Stream-Zeile und jede der elf
/// Administration-Antwortzeilen — reine Funktionen, netzlos testbar. Form-
/// Vorbild: <c>examples/grpc-client/format_test.go</c>, die Sprintf-Aufrufe
/// in <c>consumer.go</c>/<c>tables_admin.go</c>/<c>changes.go</c>/
/// <c>retention.go</c>/<c>diagnose.go</c>.
/// </summary>
public class FormatTests
{
    [Fact]
    public void FormatChange_CarriesIdentityAndPayload()
    {
        var change = new Change
        {
            ChangeId = "c-1",
            Schema = "public",
            Table = "orders",
            Operation = "INSERT",
            NewImage = ByteString.CopyFromUtf8("{\"id\":1}"),
        };

        var got = Format.FormatChange(change);

        Assert.Equal("grpc-client: change_id=c-1 table=public.orders operation=INSERT new_image={\"id\":1}", got);
    }

    [Fact]
    public void FormatChange_HandlesEmptyNewImage()
    {
        var change = new Change
        {
            ChangeId = "c-2",
            Schema = "public",
            Table = "orders",
            Operation = "DELETE",
        };

        var got = Format.FormatChange(change);

        Assert.Equal("grpc-client: change_id=c-2 table=public.orders operation=DELETE new_image=", got);
    }

    [Fact]
    public void FormatRegisterConsumer_CarriesIdempotencyFlag()
    {
        var resp = new RegisterConsumerResponse { ConsumerId = "c-1", Name = "Consumer", AlreadyRegistered = true };

        Assert.Equal("grpc-client: consumer_id=c-1 name=Consumer already_registered=true", Format.FormatRegisterConsumer(resp));
    }

    [Fact]
    public void FormatAcknowledgeConsumer_CarriesOffset()
    {
        var resp = new AcknowledgeConsumerResponse { ConsumerId = "c-1", SourceId = "quelle-1", Offset = 42 };

        Assert.Equal("grpc-client: consumer_id=c-1 source_id=quelle-1 offset=42", Format.FormatAcknowledgeConsumer(resp));
    }

    [Fact]
    public void FormatConsumerPosition_DistinguishesUnacknowledgedFromZeroOffset()
    {
        var resp = new GetConsumerPositionResponse { ConsumerId = "c-1", SourceId = "quelle-1", Offset = 0, Acknowledged = false };

        Assert.Equal("grpc-client: consumer_id=c-1 source_id=quelle-1 offset=0 acknowledged=false", Format.FormatConsumerPosition(resp));
    }

    [Fact]
    public void FormatRemoveConsumer_CarriesRemovedFlag()
    {
        var resp = new RemoveConsumerResponse { ConsumerId = "c-1", Removed = false };

        Assert.Equal("grpc-client: consumer_id=c-1 removed=false", Format.FormatRemoveConsumer(resp));
    }

    [Fact]
    public void FormatEnableTable_CarriesIdempotencyFlag()
    {
        var resp = new EnableTableResponse
        {
            TableId = "public.orders", Source = "quelle-1", Schema = "public", Table = "orders", AlreadyEnabled = false,
        };

        Assert.Equal(
            "grpc-client: table_id=public.orders source=quelle-1 schema=public table=orders already_enabled=false",
            Format.FormatEnableTable(resp));
    }

    [Fact]
    public void FormatDisableTable_CarriesRemovedAndRetained()
    {
        var resp = new DisableTableResponse { Removed = true, Retained = true };

        Assert.Equal("grpc-client: removed=true retained=true", Format.FormatDisableTable(resp));
    }

    [Fact]
    public void FormatTableStatus_CarriesEnabledAndRetained()
    {
        var resp = new GetTableStatusResponse { Enabled = true, Retained = false };

        Assert.Equal("grpc-client: enabled=true retained=false", Format.FormatTableStatus(resp));
    }

    [Fact]
    public void FormatListTables_ListsTablesAndRetainedSeparately()
    {
        var resp = new ListTablesResponse();
        resp.Tables.Add(new SourceTable { TableId = "public.orders", Source = "quelle-1", Schema = "public", Table = "orders" });
        resp.Retained.Add(new SourceTable { TableId = "public.legacy", Source = "quelle-1", Schema = "public", Table = "legacy" });

        var got = Format.FormatListTables(resp);

        Assert.Equal(
            "grpc-client: tables=1 retained=1\n" +
            "  table: table_id=public.orders source=quelle-1 schema=public table=orders\n" +
            "  retained: table_id=public.legacy source=quelle-1 schema=public table=legacy",
            got);
    }

    [Fact]
    public void FormatListTables_EmptyListsCarryZeroCounts()
    {
        var resp = new ListTablesResponse();

        Assert.Equal("grpc-client: tables=0 retained=0", Format.FormatListTables(resp));
    }

    [Fact]
    public void FormatRunRetention_CarriesDeletedCount()
    {
        var resp = new RunRetentionResponse { Deleted = 7 };

        Assert.Equal("grpc-client: deleted=7", Format.FormatRunRetention(resp));
    }

    [Fact]
    public void FormatReadChanges_ListsEachChange()
    {
        var resp = new ReadChangesResponse();
        resp.Changes.Add(new ChangeRecord
        {
            ChangeId = "c-1", Schema = "public", Table = "orders", Operation = "INSERT",
            CommitPosition = 100, Origin = "wal", NewImage = ByteString.CopyFromUtf8("{\"id\":1}"),
        });

        var got = Format.FormatReadChanges(resp);

        Assert.Equal(
            "grpc-client: changes=1\n" +
            "  change_id=c-1 table=public.orders operation=INSERT commit_position=100 origin=wal new_image={\"id\":1}",
            got);
    }

    [Fact]
    public void FormatReadChanges_EmptyResultCarriesZeroCount()
    {
        var resp = new ReadChangesResponse();

        Assert.Equal("grpc-client: changes=0", Format.FormatReadChanges(resp));
    }

    [Fact]
    public void FormatDiagnose_NoRetentionBlockerPrintsAbsenceLine()
    {
        var resp = new DiagnoseResponse
        {
            Heartbeat = new HeartbeatStatus { Known = false },
            RetentionBlocker = new RetentionBlocker { Present = false },
        };

        var got = Format.FormatDiagnose(resp);

        Assert.Contains("heartbeat_known=false", got);
        Assert.Contains("  retention_blocker: kein Blocker", got);
    }

    [Fact]
    public void FormatDiagnose_PresentRetentionBlockerCarriesConsumerFields()
    {
        var resp = new DiagnoseResponse
        {
            Heartbeat = new HeartbeatStatus { Known = true, AgeSeconds = 1.5, ErrorClass = "" },
            RetentionBlocker = new RetentionBlocker
            {
                Present = true, ConsumerId = "c-1", Name = "Consumer", AcknowledgedPosition = 10, BacklogKnown = true, Backlog = 5,
            },
        };

        var got = Format.FormatDiagnose(resp);

        Assert.Contains(
            "retention_blocker: consumer_id=c-1 name=Consumer acknowledged_position=10 backlog_known=true backlog=5",
            got);
    }

    [Fact]
    public void FormatDiagnose_ListsConsumerLagsAndBackfillRows()
    {
        var resp = new DiagnoseResponse
        {
            Heartbeat = new HeartbeatStatus { Known = true },
            RetentionBlocker = new RetentionBlocker { Present = false },
        };
        resp.ConsumerLags.Add(new ConsumerLag { ConsumerId = "c-1", Known = true, Lag = 2.0 });
        resp.Backfill.Add(new BackfillTableStatus
        {
            Schema = "public", Table = "orders", Status = "completed", RowsCopied = 100,
            EstimatedRowsKnown = true, EstimatedRows = 100, WarnEstimatedSize = false, WarnDuration = false, ErrorMessage = "",
        });

        var got = Format.FormatDiagnose(resp);

        Assert.Contains("consumer_lag: consumer_id=c-1 known=true lag=2", got);
        Assert.Contains("backfill: schema=public table=orders status=completed rows_copied=100", got);
    }
}
