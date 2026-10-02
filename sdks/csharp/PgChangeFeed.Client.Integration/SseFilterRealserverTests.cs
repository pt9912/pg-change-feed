using System.Net.Http;
using System.Text.Json;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Sse;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>One received change reduced to what the filter phase asserts on.</summary>
internal sealed record FilterRow(string ChangeId, string Schema, string Table, string? Name);

/// <summary>Collects the rows of one stream on a background task.</summary>
internal sealed class FilterCollector
{
    private readonly object _gate = new();
    private readonly List<FilterRow> _rows = new();
    private Exception? _failure;

    internal IReadOnlyList<FilterRow> Rows
    {
        get { lock (_gate) { return _rows.ToArray(); } }
    }

    internal Exception? Failure
    {
        get { lock (_gate) { return _failure; } }
    }

    /// <summary>Consumes the stream until it ends or <paramref name="stop"/> is cancelled;
    /// a fault before the stop is kept in <see cref="Failure"/>.</summary>
    internal async Task RunAsync(IAsyncEnumerable<Sse.Models.Change> stream, CancellationToken stop)
    {
        try
        {
            await foreach (var change in stream)
            {
                string? name = null;
                if (change.NewImage is { } image
                    && image.ValueKind == JsonValueKind.Object
                    && image.TryGetProperty("name", out var n)
                    && n.ValueKind == JsonValueKind.String)
                {
                    name = n.GetString();
                }
                lock (_gate) { _rows.Add(new FilterRow(change.ChangeId, change.Schema, change.Table, name)); }
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
/// Real-server filter phase for the SSE stream client. Three tables are
/// captured: <c>A</c> and <c>B</c> in the first schema and a table named like
/// <c>A</c> in a second schema. A stream opened with <c>schema</c> and
/// <c>table</c> of <c>A</c> receives exactly the changes of that table, a
/// stream opened with the second <c>schema</c> alone receives exactly the
/// changes of that schema, and a stream without a filter receives all three
/// (checked over a quiet window after the unfiltered stream has them all).
/// </summary>
public sealed class SseFilterRealserverTests
{
    private static readonly TimeSpan PositiveDeadline = TimeSpan.FromSeconds(90);

    [Fact]
    public async Task StreamsWithSchemaAndTableFilterReceiveOnlyTheirSelectionAndStreamWithoutFilterReceivesAll()
    {
        using var httpClient = new HttpClient { Timeout = Timeout.InfiniteTimeSpan };
        var options = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), PhaseEnvironment.ApiToken);
        var client = new PgChangeFeedSseClient(httpClient, options);
        using var stop = new CancellationTokenSource();
        var f1 = new FilterCollector();
        var f2 = new FilterCollector();
        var unfiltered = new FilterCollector();

        var consumers = new[]
        {
            f1.RunAsync(
                client.StreamChangesAsync(
                    stop.Token, schema: PhaseEnvironment.FilterSchemaA, table: PhaseEnvironment.FilterTableA),
                stop.Token),
            f2.RunAsync(
                client.StreamChangesAsync(stop.Token, schema: PhaseEnvironment.FilterSchemaOther),
                stop.Token),
            unfiltered.RunAsync(client.StreamChangesAsync(stop.Token), stop.Token),
        };

        PhaseEnvironment.Print("READY");
        var deadline = DateTime.UtcNow + PositiveDeadline;
        while (!(Own(f1.Rows).Any(IsTableA) && Own(f2.Rows).Any(IsOtherSchema) && HasAllThree(unfiltered.Rows)))
        {
            ThrowOnFailure(f1, f2, unfiltered);
            Assert.True(DateTime.UtcNow < deadline,
                "innerhalb der Frist weder die Change der Tabelle am Client mit Schema und Tabelle, " +
                "noch die des zweiten Schemas am Client mit Schema, noch alle drei Gruppen am Client ohne Filter empfangen");
            await Task.Delay(100);
        }
        PhaseEnvironment.Print("SEEN");

        // The window starts only after the client without a filter received
        // all three groups: the same delivery path has then demonstrably
        // dispatched the foreign changes, so their absence at the filtered
        // clients is no early cut-off.
        await Task.Delay(TimeSpan.FromSeconds(PhaseEnvironment.FilterQuietSeconds));
        ThrowOnFailure(f1, f2, unfiltered);
        stop.Cancel();
        await Task.WhenAll(consumers);

        var f1Rows = f1.Rows;
        var f2Rows = f2.Rows;
        var ownUnfiltered = Own(unfiltered.Rows).ToList();
        var f1Foreign = f1Rows.Where(r => !IsTableA(r)).ToList();
        var f2Foreign = f2Rows.Where(r => !IsOtherSchema(r)).ToList();
        foreach (var row in f1Rows)
        {
            PhaseEnvironment.Print($"RECEIVED_F1 change_id={row.ChangeId} schema={row.Schema} table={row.Table}");
        }
        foreach (var row in f2Rows)
        {
            PhaseEnvironment.Print($"RECEIVED_F2 change_id={row.ChangeId} schema={row.Schema} table={row.Table}");
        }
        foreach (var row in ownUnfiltered)
        {
            PhaseEnvironment.Print($"RECEIVED_U change_id={row.ChangeId} schema={row.Schema} table={row.Table}");
        }
        // The line carries the measured counts of foreign changes before the
        // checks below, so a foreign receipt is visible in the line itself.
        PhaseEnvironment.Print(
            $"FILTER_RESULT f1={f1Rows.Count} f1_foreign={f1Foreign.Count} f2={f2Rows.Count} " +
            $"f2_foreign={f2Foreign.Count} unfiltered={ownUnfiltered.Count} " +
            $"quiet_seconds={PhaseEnvironment.FilterQuietSeconds}");

        Assert.True(f1Foreign.Count == 0,
            "der Client mit Schema und Tabelle empfing fremde Changes: " +
            string.Join(", ", f1Foreign.Select(r => $"{r.ChangeId}({r.Schema}.{r.Table})")));
        Assert.True(f2Foreign.Count == 0,
            "der Client mit Schema allein empfing fremde Changes: " +
            string.Join(", ", f2Foreign.Select(r => $"{r.ChangeId}({r.Schema}.{r.Table})")));
        Assert.NotEmpty(Own(f1Rows));
        Assert.NotEmpty(Own(f2Rows));
        Assert.True(HasAllThree(unfiltered.Rows), "der Client ohne Filter sah nicht alle drei Gruppen");
    }

    private static IEnumerable<FilterRow> Own(IEnumerable<FilterRow> rows)
        => rows.Where(r => r.Name == PhaseEnvironment.Sentinel);

    private static bool IsTableA(FilterRow r)
        => r.Schema == PhaseEnvironment.FilterSchemaA && r.Table == PhaseEnvironment.FilterTableA;

    private static bool IsTableB(FilterRow r)
        => r.Schema == PhaseEnvironment.FilterSchemaA && r.Table == PhaseEnvironment.FilterTableB;

    private static bool IsOtherSchema(FilterRow r) => r.Schema == PhaseEnvironment.FilterSchemaOther;

    private static bool HasAllThree(IEnumerable<FilterRow> rows)
    {
        var own = Own(rows).ToList();
        return own.Any(IsTableA) && own.Any(IsTableB) && own.Any(IsOtherSchema);
    }

    private static void ThrowOnFailure(params FilterCollector[] collectors)
    {
        foreach (var collector in collectors)
        {
            if (collector.Failure is { } failure)
            {
                throw new InvalidOperationException("Ein Stream-Konsument ist ausgefallen.", failure);
            }
        }
    }
}
