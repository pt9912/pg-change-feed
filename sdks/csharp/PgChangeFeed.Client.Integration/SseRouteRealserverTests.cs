using System.Net.Http;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Sse;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server routing phase for the SSE stream client: a stream opened with
/// <c>target</c> receives the changes routed to that target and no change of
/// another target or without a target (checked over a quiet window), while a
/// stream opened without <c>target</c> receives all of them.
/// </summary>
public sealed class SseRouteRealserverTests
{
    [Fact]
    public async Task StreamWithTargetReceivesOnlyItsTargetAndStreamWithoutTargetReceivesAll()
    {
        using var httpClient = new HttpClient { Timeout = Timeout.InfiniteTimeSpan };
        var options = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), PhaseEnvironment.ApiToken);
        var client = new PgChangeFeedSseClient(httpClient, options);
        using var stop = new CancellationTokenSource();
        var targeted = new RouteCollector();
        var unfiltered = new RouteCollector();

        var consumers = new[]
        {
            targeted.RunAsync(
                Rows(client.StreamChangesAsync(stop.Token, target: PhaseEnvironment.RouteTargetA)),
                stop.Token),
            unfiltered.RunAsync(
                Rows(client.StreamChangesAsync(stop.Token)),
                stop.Token),
        };

        await RouteScenario.RunStreamsAsync(targeted, unfiltered, consumers, stop);
    }

    private static async IAsyncEnumerable<RouteRow> Rows(IAsyncEnumerable<Sse.Models.Change> stream)
    {
        await foreach (var change in stream)
        {
            yield return RouteRow.From(change.ChangeId, change.Table, change.NewImage?.ToString());
        }
    }
}
