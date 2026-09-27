using System.Net.Http;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Sse;
using PgChangeFeed.Client.Sse.Models;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server rule phase for the SSE stream client: opens
/// <c>GET /changes/stream</c> against a table carrying an active
/// <c>rename_column</c> rule and receives a change committed afterwards; the
/// row image arrives with the renamed key and without the source key, read
/// through the client's opaque <see cref="System.Text.Json.JsonElement"/>
/// model.
/// </summary>
public sealed class SseRuleRealserverTests
{
    [Fact]
    public async Task ReceivesARenamedRowImageKeyOverTheStream()
    {
        using var httpClient = new HttpClient();
        var options = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), PhaseEnvironment.ApiToken);
        var client = new PgChangeFeedSseClient(httpClient, options);
        PhaseEnvironment.Print("READY");

        var received = await ReceiveSentinelAsync(client.StreamChangesAsync(PhaseEnvironment.ReceiveCts.Token));

        Assert.Equal("INSERT", received.Operation);
        Assert.Equal(PhaseEnvironment.Table, received.Table);
        Assert.NotNull(received.NewImage);
        var newImage = received.NewImage!.Value;
        Assert.True(newImage.TryGetProperty(PhaseEnvironment.RuleTargetKey, out var target));
        Assert.Equal(PhaseEnvironment.Sentinel, target.GetString());
        Assert.False(newImage.TryGetProperty(PhaseEnvironment.RuleSourceKey, out _));

        PhaseEnvironment.Print(
            $"RECEIVED change_id={received.ChangeId} table={received.Table} " +
            $"operation={received.Operation} new_image={received.NewImage}");
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
