using PgChangeFeed.Client;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Http.Models;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server TLS phase for the HTTP API client: the feed container serves
/// HTTP over TLS with a certificate the system does not know. With the trust
/// anchor the SDK registers a disposable consumer with the admin token and
/// lists tables with the reader token (the runner reads the registration back
/// from <c>cdc.consumer</c>); without an anchor, with a foreign anchor and
/// with a server name outside the certificate the same call fails the TLS
/// verification.
/// </summary>
public sealed class HttpTlsRealserverTests
{
    [Fact]
    public async Task RegistersAConsumerAndListsTablesOverTlsWithTheTrustAnchor()
    {
        PhaseEnvironment.Print("READY");
        var address = new Uri(PhaseEnvironment.HttpAddr);

        using var admin = new PgChangeFeedHttpClient(
            new PgChangeFeedClientOptions(address, PhaseEnvironment.AdminToken, PhaseEnvironment.TlsCaFile));
        var consumerId = $"csharp-sdk-tls-{DateTime.UtcNow:yyyyMMddHHmmss}";
        var registered = await admin.RegisterConsumerAsync(
            new RegisterConsumerRequest(consumerId, $"C# SDK TLS {consumerId}"), PhaseEnvironment.ReceiveCts.Token);
        Assert.Equal(consumerId, registered.ConsumerId);

        using var reader = new PgChangeFeedHttpClient(
            new PgChangeFeedClientOptions(address, PhaseEnvironment.ReaderToken, PhaseEnvironment.TlsCaFile));
        var tables = await reader.ListTablesAsync(
            PhaseEnvironment.SourceId, PhaseEnvironment.HttpPublication, PhaseEnvironment.ReceiveCts.Token);

        PhaseEnvironment.Print($"RECEIVED consumer_id={registered.ConsumerId} tables={tables.Tables.Count}");
    }

    [Fact]
    public async Task CallsWithoutAnchorWithForeignAnchorAndWithAWrongServerNameFailTheTlsCheck()
    {
        await TlsScenarios.ExpectRefusalsAsync(
            new Uri(PhaseEnvironment.HttpAddr),
            PhaseEnvironment.ReaderToken,
            async options =>
            {
                using var client = new PgChangeFeedHttpClient(options);
                await client.ListTablesAsync(
                    PhaseEnvironment.SourceId, PhaseEnvironment.HttpPublication, PhaseEnvironment.RejectCts.Token);
            });
    }
}
