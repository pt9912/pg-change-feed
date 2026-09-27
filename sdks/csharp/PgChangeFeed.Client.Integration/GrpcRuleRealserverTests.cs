using System.Text.Json;
using Cdc.Stream.V1;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server rule phase for the gRPC stream client: opens the server
/// stream against a table carrying an active <c>rename_column</c> rule and
/// receives a change committed afterwards; the row image arrives with the
/// renamed key and without the source key, read through the client's opaque
/// wire form (the raw bytes, parsed as JSON here only to assert on the
/// property under test — the client itself never interprets the row image
/// shape).
/// </summary>
public sealed class GrpcRuleRealserverTests
{
    [Fact]
    public async Task ReceivesARenamedRowImageKeyOverTheStream()
    {
        var options = new PgChangeFeedClientOptions(new Uri($"http://{PhaseEnvironment.GrpcAddr}"), PhaseEnvironment.ApiToken);
        using var client = new PgChangeFeedGrpcClient(options);
        PhaseEnvironment.Print("READY");

        var received = await ReceiveSentinelAsync(client.StreamChangesAsync(PhaseEnvironment.ReceiveCts.Token));

        Assert.Equal("INSERT", received.Operation);
        Assert.Equal(PhaseEnvironment.Table, received.Table);
        using var doc = JsonDocument.Parse(received.NewImage?.ToStringUtf8() ?? "{}");
        Assert.True(doc.RootElement.TryGetProperty(PhaseEnvironment.RuleTargetKey, out var target));
        Assert.Equal(PhaseEnvironment.Sentinel, target.GetString());
        Assert.False(doc.RootElement.TryGetProperty(PhaseEnvironment.RuleSourceKey, out _));

        PhaseEnvironment.Print(
            $"RECEIVED change_id={received.ChangeId} table={received.Table} " +
            $"operation={received.Operation} new_image={received.NewImage?.ToStringUtf8()}");
    }

    private static async Task<Change> ReceiveSentinelAsync(IAsyncEnumerable<Change> stream)
    {
        var deadline = DateTime.UtcNow + TimeSpan.FromSeconds(90);
        await foreach (var change in stream)
        {
            if (change.Table == PhaseEnvironment.Table &&
                change.Operation == "INSERT" &&
                (change.NewImage?.ToStringUtf8() ?? string.Empty).Contains(PhaseEnvironment.Sentinel))
            {
                return change;
            }
            Assert.True(DateTime.UtcNow < deadline,
                $"kein Change mit dem Sentinel {PhaseEnvironment.Sentinel} innerhalb der Frist empfangen");
        }
        throw new InvalidOperationException("Der Stream endete, bevor der Sentinel empfangen wurde.");
    }
}
