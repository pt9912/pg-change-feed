using PgChangeFeed.Client;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server TLS phase for the gRPC stream client: with the trust anchor the
/// SDK opens the server stream over TLS and receives a change committed
/// afterwards (the runner reads its <c>change_id</c> back from
/// <c>cdc.changes</c>); without an anchor, with a foreign anchor and with a
/// server name outside the certificate the stream fails the TLS verification.
/// </summary>
public sealed class GrpcTlsRealserverTests
{
    private static Uri Address => new($"https://{PhaseEnvironment.GrpcAddr}");

    [Fact]
    public async Task ReceivesACommittedChangeOverTlsWithTheTrustAnchor()
    {
        using var client = new PgChangeFeedGrpcClient(
            new PgChangeFeedClientOptions(Address, PhaseEnvironment.ApiToken, PhaseEnvironment.TlsCaFile));
        PhaseEnvironment.Print("READY");

        var received = await GrpcRealserverTests.ReceiveSentinelAsync(
            client.StreamChangesAsync(cancellationToken: PhaseEnvironment.ReceiveCts.Token));

        Assert.Equal("INSERT", received.Operation);
        Assert.Equal(PhaseEnvironment.Table, received.Table);
        PhaseEnvironment.Print(
            $"RECEIVED change_id={received.ChangeId} table={received.Table} " +
            $"operation={received.Operation} new_image={received.NewImage?.ToStringUtf8()}");
    }

    [Fact]
    public async Task StreamsWithoutAnchorWithForeignAnchorAndWithAWrongServerNameFailTheTlsCheck()
    {
        await TlsScenarios.ExpectRefusalsAsync(
            Address,
            PhaseEnvironment.ApiToken,
            async options =>
            {
                using var client = new PgChangeFeedGrpcClient(options);
                await foreach (var _ in client.StreamChangesAsync(cancellationToken: PhaseEnvironment.RejectCts.Token)) { }
            });
    }
}
