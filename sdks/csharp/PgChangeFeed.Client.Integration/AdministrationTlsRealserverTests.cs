using Cdc.Administration.V1;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server TLS phase for the gRPC administration client: with the trust
/// anchor the SDK registers a disposable consumer with the admin token and
/// lists tables with the reader token over TLS (the runner reads the
/// registration back from <c>cdc.consumer</c>); without an anchor, with a
/// foreign anchor and with a server name outside the certificate the same call
/// fails the TLS verification.
/// </summary>
public sealed class AdministrationTlsRealserverTests
{
    private static Uri Address => new($"https://{PhaseEnvironment.GrpcAddr}");

    [Fact]
    public async Task RegistersAConsumerAndListsTablesOverTlsWithTheTrustAnchor()
    {
        PhaseEnvironment.Print("READY");

        using var admin = new PgChangeFeedAdministrationClient(
            new PgChangeFeedClientOptions(Address, PhaseEnvironment.AdminToken, PhaseEnvironment.TlsCaFile));
        var consumerId = $"csharp-sdk-tls-admin-{DateTime.UtcNow:yyyyMMddHHmmss}";
        var registered = await admin.RegisterConsumerAsync(
            new RegisterConsumerRequest { ConsumerId = consumerId, Name = $"C# SDK TLS {consumerId}" },
            PhaseEnvironment.ReceiveCts.Token);
        Assert.Equal(consumerId, registered.ConsumerId);

        using var reader = new PgChangeFeedAdministrationClient(
            new PgChangeFeedClientOptions(Address, PhaseEnvironment.ReaderToken, PhaseEnvironment.TlsCaFile));
        var tables = await reader.ListTablesAsync(
            new ListTablesRequest { Source = PhaseEnvironment.SourceId, Publication = PhaseEnvironment.HttpPublication },
            PhaseEnvironment.ReceiveCts.Token);

        PhaseEnvironment.Print($"RECEIVED consumer_id={registered.ConsumerId} tables={tables.Tables.Count}");
    }

    [Fact]
    public async Task CallsWithoutAnchorWithForeignAnchorAndWithAWrongServerNameFailTheTlsCheck()
    {
        await TlsScenarios.ExpectRefusalsAsync(
            Address,
            PhaseEnvironment.ReaderToken,
            async options =>
            {
                using var client = new PgChangeFeedAdministrationClient(options);
                await client.ListTablesAsync(
                    new ListTablesRequest { Source = PhaseEnvironment.SourceId, Publication = PhaseEnvironment.HttpPublication },
                    PhaseEnvironment.RejectCts.Token);
            });
    }
}
