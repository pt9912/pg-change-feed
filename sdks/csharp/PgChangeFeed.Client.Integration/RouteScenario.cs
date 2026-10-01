using System.Text.Json;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>One received change reduced to what the routing phases assert on.</summary>
internal sealed record RouteRow(string ChangeId, string Table, string? Region, string? Name)
{
    /// <summary>Reads <c>region</c> and <c>name</c> from the JSON text of a row image.</summary>
    internal static RouteRow From(string changeId, string table, string? imageJson)
    {
        string? region = null;
        string? name = null;
        if (!string.IsNullOrEmpty(imageJson))
        {
            using var doc = JsonDocument.Parse(imageJson);
            if (doc.RootElement.TryGetProperty("region", out var r) && r.ValueKind == JsonValueKind.String)
            {
                region = r.GetString();
            }
            if (doc.RootElement.TryGetProperty("name", out var n) && n.ValueKind == JsonValueKind.String)
            {
                name = n.GetString();
            }
        }
        return new RouteRow(changeId, table, region, name);
    }
}

/// <summary>Collects the rows of one stream on a background task.</summary>
internal sealed class RouteCollector
{
    private readonly object _gate = new();
    private readonly List<RouteRow> _rows = new();
    private Exception? _failure;

    internal IReadOnlyList<RouteRow> Rows
    {
        get { lock (_gate) { return _rows.ToArray(); } }
    }

    internal Exception? Failure
    {
        get { lock (_gate) { return _failure; } }
    }

    /// <summary>Consumes the stream until it ends or <paramref name="stop"/> is cancelled;
    /// a fault before the stop is kept in <see cref="Failure"/>.</summary>
    internal async Task RunAsync(IAsyncEnumerable<RouteRow> stream, CancellationToken stop)
    {
        try
        {
            await foreach (var row in stream)
            {
                lock (_gate) { _rows.Add(row); }
            }
        }
        catch (Exception) when (stop.IsCancellationRequested)
        {
            // Cancelling the consumer ends the enumeration; nothing to keep.
        }
        catch (Exception ex)
        {
            lock (_gate) { _failure = ex; }
        }
    }
}

/// <summary>
/// The shared course of the routing phases. A client with the target
/// <c>RouteTargetA</c> and a client without a target read the same table; the
/// table carries two routing rules (region A to target A, region B to target
/// B), and the runner commits groups of three changes (no rule, B, A). The
/// client with the target must receive the change for A and no change of any
/// other region; the client without a target must receive all three kinds.
/// </summary>
internal static class RouteScenario
{
    private static readonly TimeSpan PositiveDeadline = TimeSpan.FromSeconds(90);

    private static IEnumerable<RouteRow> Own(IEnumerable<RouteRow> rows)
        => rows.Where(r => r.Table == PhaseEnvironment.Table && r.Name == PhaseEnvironment.Sentinel);

    private static bool HasAllThree(IEnumerable<RouteRow> rows)
    {
        var regions = Own(rows).Select(r => r.Region).ToHashSet();
        return regions.Contains(PhaseEnvironment.RouteTargetA)
            && regions.Contains(PhaseEnvironment.RouteTargetB)
            && regions.Contains(PhaseEnvironment.RouteRegionNone);
    }

    private static void ThrowOnFailure(params RouteCollector[] collectors)
    {
        foreach (var collector in collectors)
        {
            if (collector.Failure is { } failure)
            {
                throw new InvalidOperationException("Ein Stream-Konsument ist ausgefallen.", failure);
            }
        }
    }

    /// <summary>Stream surfaces: positive phase, quiet window, then the checks.</summary>
    internal static async Task RunStreamsAsync(
        RouteCollector targeted, RouteCollector unfiltered, Task[] consumers, CancellationTokenSource stop)
    {
        PhaseEnvironment.Print("READY");
        var deadline = DateTime.UtcNow + PositiveDeadline;
        while (!(Own(targeted.Rows).Any(r => r.Region == PhaseEnvironment.RouteTargetA) && HasAllThree(unfiltered.Rows)))
        {
            ThrowOnFailure(targeted, unfiltered);
            Assert.True(DateTime.UtcNow < deadline,
                $"innerhalb der Frist weder die Change des Ziels {PhaseEnvironment.RouteTargetA} am Client mit Ziel " +
                "noch alle drei Gruppen am Client ohne Ziel empfangen");
            await Task.Delay(100);
        }
        PhaseEnvironment.Print("SEEN");

        // The window starts only after the client without a target received the
        // change of the other target: the same delivery path has then
        // demonstrably dispatched it, so its absence at the client with a target
        // is no early cut-off.
        await Task.Delay(TimeSpan.FromSeconds(PhaseEnvironment.RouteQuietSeconds));
        ThrowOnFailure(targeted, unfiltered);
        stop.Cancel();
        await Task.WhenAll(consumers);

        Evaluate(targeted.Rows, unfiltered.Rows, PhaseEnvironment.RouteQuietSeconds);
    }

    /// <summary>Pull surface: the target read equals the region-A part of the unfiltered read.</summary>
    internal static async Task RunPullAsync(Func<string?, Task<IReadOnlyList<RouteRow>>> read)
    {
        PhaseEnvironment.Print("READY");
        var deadline = DateTime.UtcNow + PositiveDeadline;
        while (!HasAllThree(await read(null)))
        {
            Assert.True(DateTime.UtcNow < deadline,
                "innerhalb der Frist nicht alle drei Gruppen über die ungefilterte Lesung gelesen");
            await Task.Delay(300);
        }
        PhaseEnvironment.Print("SEEN");

        var before = await read(null);
        var targeted = await read(PhaseEnvironment.RouteTargetA);
        var after = await read(null);

        var aOf = (IEnumerable<RouteRow> rows) => Own(rows)
            .Where(r => r.Region == PhaseEnvironment.RouteTargetA).Select(r => r.ChangeId).ToHashSet();
        var targetedOwn = Own(targeted).Select(r => r.ChangeId).ToHashSet();
        Assert.True(aOf(before).IsSubsetOf(targetedOwn),
            "die Lesung mit Ziel enthält nicht jede Change des Ziels, die die Lesung davor ohne Ziel lieferte");
        Assert.True(targetedOwn.IsSubsetOf(aOf(after)),
            "die Lesung mit Ziel enthält eine Change, die die Lesung danach ohne Ziel nicht dem Ziel zuordnet");

        Evaluate(targeted, after, 0);
    }

    private static void Evaluate(IReadOnlyList<RouteRow> targeted, IReadOnlyList<RouteRow> unfiltered, int quietSeconds)
    {
        var foreign = targeted
            .Where(r => r.Table != PhaseEnvironment.Table || r.Region != PhaseEnvironment.RouteTargetA)
            .ToList();
        Assert.True(foreign.Count == 0,
            "der Client mit Ziel " + PhaseEnvironment.RouteTargetA + " empfing fremde Changes: " +
            string.Join(", ", foreign.Select(r => $"{r.ChangeId}(region={r.Region})")));
        Assert.NotEmpty(Own(targeted));
        Assert.True(HasAllThree(unfiltered), "der Client ohne Ziel sah nicht alle drei Gruppen");

        var ownUnfiltered = Own(unfiltered).ToList();
        foreach (var row in targeted)
        {
            PhaseEnvironment.Print($"RECEIVED_TARGETED change_id={row.ChangeId} table={row.Table} region={row.Region}");
        }
        foreach (var row in ownUnfiltered)
        {
            PhaseEnvironment.Print($"RECEIVED_UNFILTERED change_id={row.ChangeId} table={row.Table} region={row.Region}");
        }
        PhaseEnvironment.Print(
            $"ROUTE_RESULT target={PhaseEnvironment.RouteTargetA} targeted={targeted.Count} foreign=0 " +
            $"unfiltered={ownUnfiltered.Count} quiet_seconds={quietSeconds}");
    }
}
