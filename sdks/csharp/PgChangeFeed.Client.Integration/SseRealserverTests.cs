using System.Net.Http;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Sse;
using PgChangeFeed.Client.Sse.Models;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server phase for the SSE stream client: opens
/// <c>GET /changes/stream</c> against the running feed container and
/// receives a change committed afterwards; a second open with an unknown
/// token is rejected with HTTP status 401, as a typed
/// <see cref="PgChangeFeedException"/> — the same set the HTTP client
/// throws.
/// </summary>
public sealed class SseRealserverTests
{
    [Fact]
    public async Task ReceivesACommittedChangeOverTheStream()
    {
        using var httpClient = new HttpClient();
        var options = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), PhaseEnvironment.ApiToken);
        var client = new PgChangeFeedSseClient(httpClient, options);
        PhaseEnvironment.Print("READY");

        var received = await ReceiveSentinelAsync(client.StreamChangesAsync(PhaseEnvironment.ReceiveCts.Token));

        Assert.NotEqual(string.Empty, received.ChangeId);
        Assert.NotEqual(string.Empty, received.TransactionId);
        Assert.NotEqual(string.Empty, received.SourceTableId);
        Assert.NotEqual(string.Empty, received.SchemaVersion);
        Assert.Equal("INSERT", received.Operation);
        Assert.Null(received.OldImage);
        Assert.Equal(PhaseEnvironment.Table, received.Table);
        Assert.Equal("public", received.Schema);
        Assert.Contains(PhaseEnvironment.Sentinel, received.NewImage?.ToString() ?? string.Empty);

        PhaseEnvironment.Print(
            $"RECEIVED change_id={received.ChangeId} table={received.Table} " +
            $"operation={received.Operation} new_image={received.NewImage}");
    }

    [Fact]
    public async Task OpenWithUnknownTokenIsRejected401()
    {
        using var httpClient = new HttpClient();
        var options = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), "no-such-token");
        var client = new PgChangeFeedSseClient(httpClient, options);

        var exception = await Assert.ThrowsAsync<PgChangeFeedUnauthorizedException>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync(PhaseEnvironment.RejectCts.Token)) { }
        });

        Assert.Equal(401, exception.StatusCode);
        PhaseEnvironment.Print("REJECTED status=401");
    }

    private static async Task<Change> ReceiveSentinelAsync(IAsyncEnumerable<Change> stream)
    {
        var deadline = DateTime.UtcNow + TimeSpan.FromSeconds(90);
        await foreach (var change in stream)
        {
            if (change.Table == PhaseEnvironment.Table &&
                change.Operation == "INSERT" &&
                (change.NewImage?.ToString() ?? string.Empty).Contains(PhaseEnvironment.Sentinel))
            {
                return change;
            }
            Assert.True(DateTime.UtcNow < deadline,
                $"kein Event mit dem Sentinel {PhaseEnvironment.Sentinel} innerhalb der Frist empfangen");
        }
        throw new InvalidOperationException("Der Stream endete, bevor der Sentinel empfangen wurde.");
    }
}