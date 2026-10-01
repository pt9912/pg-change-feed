using System.Net.Http;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Http;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server routing phase for the HTTP API client: <c>ReadChangesAsync</c>
/// with <c>target</c> returns exactly the changes routed to that target (all
/// of region A, and every one the unfiltered read assigns to A), while the
/// call without <c>target</c> returns all of them.
/// </summary>
public sealed class HttpRouteRealserverTests
{
    [Fact]
    public async Task ReadWithTargetReturnsOnlyItsTargetAndReadWithoutTargetReturnsAll()
    {
        using var httpClient = new HttpClient();
        var options = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), PhaseEnvironment.ApiToken);
        var client = new PgChangeFeedHttpClient(httpClient, options);

        await RouteScenario.RunPullAsync(async target =>
        {
            var response = await client.ReadChangesAsync(
                PhaseEnvironment.SourceId, "public", PhaseEnvironment.Table,
                cancellationToken: PhaseEnvironment.ReceiveCts.Token, target: target);
            IReadOnlyList<RouteRow> rows = response.Changes
                .Select(change => RouteRow.From(change.ChangeId, change.Table, change.NewImage?.ToString()))
                .ToList();
            return rows;
        });
    }
}
