using PgChangeFeed.Client;
using PgChangeFeed.Client.Sse;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server TLS phase for the SSE stream client: with the trust anchor the
/// SDK opens <c>GET /changes/stream</c> over TLS and receives a change
/// committed afterwards (the runner reads its <c>change_id</c> back from
/// <c>cdc.changes</c>); without an anchor, with a foreign anchor and with a
/// server name outside the certificate the stream fails the TLS verification.
/// </summary>
public sealed class SseTlsRealserverTests
{
    [Fact]
    public async Task ReceivesACommittedChangeOverTlsWithTheTrustAnchor()
    {
        using var client = new PgChangeFeedSseClient(
            new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), PhaseEnvironment.ApiToken, PhaseEnvironment.TlsCaFile));
        PhaseEnvironment.Print("READY");

        await foreach (var change in client.StreamChangesAsync(PhaseEnvironment.ReceiveCts.Token))
        {
            if (change.Table == PhaseEnvironment.Table &&
                change.Operation == "INSERT" &&
                (change.NewImage?.ToString() ?? string.Empty).Contains(PhaseEnvironment.Sentinel))
            {
                PhaseEnvironment.Print(
                    $"RECEIVED change_id={change.ChangeId} table={change.Table} " +
                    $"operation={change.Operation} new_image={change.NewImage}");
                return;
            }
        }

        Assert.Fail("The stream ended before the sentinel arrived.");
    }

    [Fact]
    public async Task StreamsWithoutAnchorWithForeignAnchorAndWithAWrongServerNameFailTheTlsCheck()
    {
        await TlsScenarios.ExpectRefusalsAsync(
            new Uri(PhaseEnvironment.HttpAddr),
            PhaseEnvironment.ApiToken,
            async options =>
            {
                using var client = new PgChangeFeedSseClient(options);
                await foreach (var _ in client.StreamChangesAsync(PhaseEnvironment.RejectCts.Token)) { }
            });
    }
}
