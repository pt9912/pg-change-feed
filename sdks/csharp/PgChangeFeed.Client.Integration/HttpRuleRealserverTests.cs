using System.Net.Http;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Http.Models;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server rule phase for the HTTP API client: reads
/// <c>GET /changes</c> (source and table as filter) against a table
/// carrying an active <c>rename_column</c> rule; the row image of the
/// change committed while the test polls arrives with the renamed key and
/// without the source key, read through the client's opaque
/// <see cref="System.Text.Json.JsonElement"/> model.
/// </summary>
public sealed class HttpRuleRealserverTests
{
    [Fact]
    public async Task ReadsARenamedRowImageKeyOverTheApi()
    {
        using var httpClient = new HttpClient();
        var options = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), PhaseEnvironment.ApiToken);
        var client = new PgChangeFeedHttpClient(httpClient, options);
        PhaseEnvironment.Print("READY");

        var received = await ReceiveSentinelAsync(client);

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

    private static async Task<Change> ReceiveSentinelAsync(PgChangeFeedHttpClient client)
    {
        var deadline = DateTime.UtcNow + TimeSpan.FromSeconds(90);
        while (DateTime.UtcNow < deadline)
        {
            var response = await client.ReadChangesAsync(
                PhaseEnvironment.SourceId, "public", PhaseEnvironment.Table,
                cancellationToken: PhaseEnvironment.ReceiveCts.Token);
            foreach (var change in response.Changes)
            {
                if (change.Table == PhaseEnvironment.Table &&
                    change.Operation == "INSERT" &&
                    (change.NewImage?.ToString() ?? string.Empty).Contains(PhaseEnvironment.Sentinel))
                {
                    return change;
                }
            }
            await Task.Delay(TimeSpan.FromMilliseconds(500));
        }
        throw new InvalidOperationException(
            $"kein Change mit dem Sentinel {PhaseEnvironment.Sentinel} innerhalb der Frist über GET /changes gelesen");
    }
}
