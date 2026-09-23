using System.Net.Http;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Http.Models;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Realserver phase for the HTTP API surface (SPEC-018): a roundtrip in the
/// pattern of the server-E2E HTTP phase — the SDK registers a disposable
/// consumer with the admin token and lists tables with the reader token;
/// the registration is held against the SQL read path (<c>cdc.consumer</c>)
/// by the runner; a call with an unknown token is rejected with HTTP status
/// 401 (SPEC-018 Negative). The nine capabilities individually stay with the
/// network-free unit tests; this roundtrip proves the wire assumptions
/// (auth header form, JSON mapping, rejection behavior) at the server.
/// </summary>
public sealed class HttpRealserverTests
{
    [Fact]
    public async Task RegistersAConsumerAndListsTablesOverTheApi()
    {
        using var httpClient = new HttpClient();
        PhaseEnvironment.Print("READY");

        var adminOptions = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), PhaseEnvironment.AdminToken);
        var adminClient = new PgChangeFeedHttpClient(httpClient, adminOptions);
        var consumerId = $"csharp-sdk-e2e-{DateTime.UtcNow:yyyyMMddHHmmss}";
        var registered = await adminClient.RegisterConsumerAsync(
            new RegisterConsumerRequest(consumerId, $"C# SDK E2E {consumerId}"),
            PhaseEnvironment.ReceiveCts.Token);
        Assert.Equal(consumerId, registered.ConsumerId);

        var readerOptions = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), PhaseEnvironment.ReaderToken);
        var readerClient = new PgChangeFeedHttpClient(httpClient, readerOptions);
        var tables = await readerClient.ListTablesAsync(
            PhaseEnvironment.SourceId, PhaseEnvironment.HttpPublication, PhaseEnvironment.ReceiveCts.Token);

        PhaseEnvironment.Print(
            $"RECEIVED consumer_id={registered.ConsumerId} tables={tables.Tables.Count}");
    }

    [Fact]
    public async Task CallWithUnknownTokenIsRejected401()
    {
        using var httpClient = new HttpClient();
        var options = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), "no-such-token");
        var client = new PgChangeFeedHttpClient(httpClient, options);

        var exception = await Assert.ThrowsAsync<PgChangeFeedUnauthorizedException>(async () =>
            await client.ListTablesAsync(
                PhaseEnvironment.SourceId, PhaseEnvironment.HttpPublication, PhaseEnvironment.RejectCts.Token));

        Assert.Equal(401, exception.StatusCode);
        PhaseEnvironment.Print("REJECTED status=401");
    }
}