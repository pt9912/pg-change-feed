using Cdc.Administration.V1;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>
/// Coverage for <c>Diagnose</c>: the populated report and the absence cases
/// each <c>Known</c>/<c>Present</c>/<c>*Known</c> field carries when the
/// source never produced the respective signal.
/// </summary>
public class PgChangeFeedAdministrationClientDiagnoseTests
{
    [Fact]
    public async Task DiagnoseAsync_HappyPath_ReturnsFullReport()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new DiagnoseResponse
        {
            Heartbeat = new HeartbeatStatus { Known = true, AgeSeconds = 1.5, ErrorClass = "" },
            CaptureLag = 0.2,
            ConsumerLags = { new ConsumerLag { ConsumerId = "c-1", Known = true, Lag = 3 } },
            RetentionBlocker = new RetentionBlocker
            {
                Present = true, ConsumerId = "c-1", Name = "n-1", AcknowledgedPosition = 10,
                BacklogKnown = true, Backlog = 5,
            },
            StorageBytes = 1024,
            Backfill =
            {
                new BackfillTableStatus
                {
                    Schema = "public", Table = "orders", Status = "completed", RowsCopied = 100,
                    EstimatedRowsKnown = true, EstimatedRows = 100, WarnEstimatedSize = false,
                    WarnDuration = false, ErrorMessage = "",
                },
            },
        });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.DiagnoseAsync(new DiagnoseRequest { Source = "src" });

        Assert.True(response.Heartbeat.Known);
        Assert.True(response.RetentionBlocker.Present);
        Assert.Single(response.ConsumerLags);
        Assert.Single(response.Backfill);
    }

    [Fact]
    public async Task DiagnoseAsync_NoHeartbeatEver_ReportsKnownFalse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new DiagnoseResponse
        {
            Heartbeat = new HeartbeatStatus { Known = false },
            RetentionBlocker = new RetentionBlocker { Present = false },
        });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.DiagnoseAsync(new DiagnoseRequest { Source = "src" });

        Assert.False(response.Heartbeat.Known);
        Assert.False(response.RetentionBlocker.Present);
        Assert.Empty(response.ConsumerLags);
        Assert.Empty(response.Backfill);
    }

    [Fact]
    public async Task DiagnoseAsync_UnknownEstimate_ReportsEstimatedRowsKnownFalse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new DiagnoseResponse
        {
            Heartbeat = new HeartbeatStatus { Known = true },
            RetentionBlocker = new RetentionBlocker { Present = false },
            Backfill =
            {
                new BackfillTableStatus
                {
                    Schema = "public", Table = "orders", Status = "running",
                    EstimatedRowsKnown = false, EstimatedRows = 0,
                },
            },
        });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.DiagnoseAsync(new DiagnoseRequest { Source = "src" });

        Assert.False(response.Backfill[0].EstimatedRowsKnown);
    }
}
