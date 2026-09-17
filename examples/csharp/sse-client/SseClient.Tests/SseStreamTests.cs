using Xunit;

namespace CdcExamples.Sse.Tests;

/// <summary>
/// Prüft den Aufbau der Stream-Adresse und das Zerlegen eines SSE-Frames
/// (<c>LH-FA-SST-008</c>, <c>ADR-0061</c>) — netzlos, reine Funktionen.
/// Form-Vorbild: <c>examples/sse-client/stream_test.go</c>.
/// </summary>
public class SseStreamTests
{
    /// <summary>
    /// Lines liefert eine Zeilenquelle über einen Ausschnitt; jeder Aufruf
    /// gibt die nächste Zeile zurück, nach der letzten <c>null</c> — dieselbe
    /// Form, die <c>Program</c> aus dem <see cref="StreamReader"/> über den
    /// Response-Body bildet.
    /// </summary>
    private static Func<string?> Lines(params string[] all)
    {
        var i = 0;
        return () => i < all.Length ? all[i++] : null;
    }

    [Fact]
    public void StreamUrlBuildsChangesStreamEndpoint()
    {
        var got = SseStream.StreamUrl("feed:8080");
        Assert.Equal("http://feed:8080/changes/stream", got);
    }

    [Fact]
    public void ReadEventReadsNameAndPayload()
    {
        var ev = SseStream.ReadEvent(Lines(
            "event: change",
            """data: {"change_id":"c-1","table":"orders"}""",
            ""));

        Assert.NotNull(ev);
        Assert.Equal("change", ev!.Name);
        Assert.Equal("""{"change_id":"c-1","table":"orders"}""", ev.Data);
    }

    [Fact]
    public void ReadEventStopsAtFrameBoundary()
    {
        var next = Lines(
            "event: change",
            """data: {"change_id":"c-1"}""",
            "",
            "event: change",
            """data: {"change_id":"c-2"}""",
            "");

        var first = SseStream.ReadEvent(next);
        var second = SseStream.ReadEvent(next);

        Assert.Equal("""{"change_id":"c-1"}""", first!.Data);
        Assert.Equal("""{"change_id":"c-2"}""", second!.Data);
    }

    [Fact]
    public void ReadEventReportsExhaustedSource()
    {
        var next = Lines(
            "event: change",
            """data: {"change_id":"c-1"}""");

        var ev = SseStream.ReadEvent(next);

        Assert.Null(ev);
    }
}
