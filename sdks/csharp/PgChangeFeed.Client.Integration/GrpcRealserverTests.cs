using Cdc.Stream.V1;
using Grpc.Core;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server phase for the gRPC stream client: opens the
/// server stream against the running feed container and receives a change
/// committed afterwards; a second open with an unknown token is rejected
/// with gRPC status <c>Unauthenticated</c>, not swallowed as an empty
/// stream. The commit-to-delivery window is
/// fire-and-forget — the runner commits a bounded sequence of unique rows
/// until one arrives; this test keeps receiving until it sees its sentinel.
/// </summary>
public sealed class GrpcRealserverTests
{
    [Fact]
    public async Task ReceivesACommittedChangeOverTheStream()
    {
        var options = new PgChangeFeedClientOptions(new Uri($"http://{PhaseEnvironment.GrpcAddr}"), PhaseEnvironment.ApiToken);
        using var client = new PgChangeFeedGrpcClient(options);
        PhaseEnvironment.Print("READY");

        var received = await ReceiveSentinelAsync(client.StreamChangesAsync(PhaseEnvironment.ReceiveCts.Token));

        Assert.NotEqual(string.Empty, received.ChangeId);
        Assert.NotEqual(string.Empty, received.TransactionId);
        Assert.NotEqual(string.Empty, received.SourceTableId);
        Assert.NotEqual(string.Empty, received.SchemaVersion);
        Assert.Equal("INSERT", received.Operation);
        Assert.Equal(PhaseEnvironment.Table, received.Table);
        Assert.Equal("public", received.Schema);
        Assert.Contains(PhaseEnvironment.Sentinel, received.NewImage?.ToStringUtf8() ?? string.Empty);

        PhaseEnvironment.Print(
            $"RECEIVED change_id={received.ChangeId} table={received.Table} " +
            $"operation={received.Operation} new_image={received.NewImage?.ToStringUtf8()}");
    }

    [Fact]
    public async Task OpenWithUnknownTokenIsRejectedUnauthenticated()
    {
        var options = new PgChangeFeedClientOptions(new Uri($"http://{PhaseEnvironment.GrpcAddr}"), "no-such-token");
        using var client = new PgChangeFeedGrpcClient(options);

        var exception = await Assert.ThrowsAnyAsync<RpcException>(async () =>
            await FirstChangeAsync(client.StreamChangesAsync(PhaseEnvironment.RejectCts.Token)));

        Assert.Equal(StatusCode.Unauthenticated, exception.StatusCode);
        PhaseEnvironment.Print("REJECTED code=Unauthenticated");
    }

    internal static async Task<Change> ReceiveSentinelAsync(IAsyncEnumerable<Change> stream)
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

    private static async Task<Change> FirstChangeAsync(IAsyncEnumerable<Change> stream)
    {
        await foreach (var change in stream)
        {
            return change;
        }
        throw new InvalidOperationException("Der Stream endete ohne eine Change.");
    }
}