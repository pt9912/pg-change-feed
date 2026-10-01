using PgChangeFeed.Client;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server routing phase for the gRPC stream client: a stream opened with
/// <c>target</c> receives the changes routed to that target and no change of
/// another target or without a target (checked over a quiet window), while a
/// stream opened without <c>target</c> receives all of them.
/// </summary>
public sealed class GrpcRouteRealserverTests
{
    [Fact]
    public async Task StreamWithTargetReceivesOnlyItsTargetAndStreamWithoutTargetReceivesAll()
    {
        var options = new PgChangeFeedClientOptions(new Uri($"http://{PhaseEnvironment.GrpcAddr}"), PhaseEnvironment.ApiToken);
        using var targetedClient = new PgChangeFeedGrpcClient(options);
        using var unfilteredClient = new PgChangeFeedGrpcClient(options);
        using var stop = new CancellationTokenSource();
        var targeted = new RouteCollector();
        var unfiltered = new RouteCollector();

        var consumers = new[]
        {
            targeted.RunAsync(
                Rows(targetedClient.StreamChangesAsync(cancellationToken: stop.Token, target: PhaseEnvironment.RouteTargetA)),
                stop.Token),
            unfiltered.RunAsync(
                Rows(unfilteredClient.StreamChangesAsync(cancellationToken: stop.Token)),
                stop.Token),
        };

        await RouteScenario.RunStreamsAsync(targeted, unfiltered, consumers, stop);
    }

    private static async IAsyncEnumerable<RouteRow> Rows(IAsyncEnumerable<Cdc.Stream.V1.Change> stream)
    {
        await foreach (var change in stream)
        {
            yield return RouteRow.From(change.ChangeId, change.Table, change.NewImage?.ToStringUtf8());
        }
    }
}
