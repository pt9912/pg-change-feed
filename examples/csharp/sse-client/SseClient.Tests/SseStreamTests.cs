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

    [Theory]
    [InlineData("eu", "http://feed:8080/changes/stream?target=eu")]
    [InlineData("a&b=c", "http://feed:8080/changes/stream?target=a%26b%3Dc")]
    [InlineData("", "http://feed:8080/changes/stream")]
    public void StreamUrlCarriesTheTargetAsEscapedQueryParameter(string target, string want)
    {
        Assert.Equal(want, SseStream.StreamUrl("feed:8080", target));
    }

    [Theory]
    [InlineData("", "", "", "http://feed:8080/changes/stream")]
    [InlineData("", "eu", "", "http://feed:8080/changes/stream?schema=eu")]
    [InlineData("", "a&b=c", "", "http://feed:8080/changes/stream?schema=a%26b%3Dc")]
    [InlineData("", "", "eu", "http://feed:8080/changes/stream?table=eu")]
    [InlineData("", "", "a&b=c", "http://feed:8080/changes/stream?table=a%26b%3Dc")]
    [InlineData("", "public", "orders", "http://feed:8080/changes/stream?schema=public&table=orders")]
    [InlineData("eu", "", "orders", "http://feed:8080/changes/stream?table=orders&target=eu")]
    [InlineData("eu", "public", "orders", "http://feed:8080/changes/stream?schema=public&table=orders&target=eu")]
    public void StreamUrlCarriesSchemaAndTableAsEscapedQueryParameters(
        string target, string schema, string table, string want)
    {
        Assert.Equal(want, SseStream.StreamUrl("feed:8080", target, schema, table));
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
