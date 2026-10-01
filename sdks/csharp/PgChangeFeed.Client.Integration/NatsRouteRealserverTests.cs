using PgChangeFeed.Client;
using PgChangeFeed.Client.Nats;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server routing phase for the NATS stream client: a subscription to
/// the subject of one delivery target receives the changes routed to that
/// target and no change of another target or without a target (checked over a
/// quiet window), while a subscription to the source's namespace receives all
/// of them.
/// </summary>
public sealed class NatsRouteRealserverTests
{
    [Fact]
    public async Task TargetSubjectReceivesOnlyItsTargetAndSourceSubjectReceivesAll()
    {
        var options = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.NatsUrl), PhaseEnvironment.NatsStreamToken);
        await using var targetedClient = new PgChangeFeedNatsStreamClient(options);
        await using var unfilteredClient = new PgChangeFeedNatsStreamClient(options);
        using var stop = new CancellationTokenSource();
        var targeted = new RouteCollector();
        var unfiltered = new RouteCollector();
        var targetSubject = PgChangeFeedNatsStreamClient.BuildTargetSubject(
            PhaseEnvironment.SourceId, PhaseEnvironment.RouteTargetA);
        var sourceSubject = PgChangeFeedNatsStreamClient.BuildSourceSubject(PhaseEnvironment.SourceId);

        var consumers = new[]
        {
            targeted.RunAsync(Rows(targetedClient.StreamChangesAsync(targetSubject, stop.Token)), stop.Token),
            unfiltered.RunAsync(Rows(unfilteredClient.StreamChangesAsync(sourceSubject, stop.Token)), stop.Token),
        };

        await RouteScenario.RunStreamsAsync(targeted, unfiltered, consumers, stop);
    }

    private static async IAsyncEnumerable<RouteRow> Rows(IAsyncEnumerable<Nats.Models.Change> stream)
    {
        await foreach (var change in stream)
        {
            yield return RouteRow.From(change.ChangeId, change.Table, change.NewImage?.ToString());
        }
    }
}
