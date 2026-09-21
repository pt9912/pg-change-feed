using PgChangeFeed.Client.Sse;
using Xunit;

namespace PgChangeFeed.Client.Tests.Sse;

/// <summary>
/// <see cref="SseFrameParser"/> Grenzfälle (SPEC-021) — netzlos, reine
/// Zeilen-Verarbeitung über einen <see cref="StringReader"/>. Muster:
/// <c>examples/csharp/sse-client/SseClient.Tests/SseStreamTests.cs</c>.
/// </summary>
public class SseFrameParserTests
{
    private static StringReader Lines(params string[] lines) => new(string.Join('\n', lines) + "\n");

    [Fact]
    public async Task ReadFrameAsync_ReadsNameAndData()
    {
        using var reader = Lines(
            "event: change",
            """data: {"change_id":"c-1","table":"orders"}""",
            "");

        var frame = await SseFrameParser.ReadFrameAsync(reader);

        Assert.NotNull(frame);
        Assert.Equal("change", frame!.Name);
        Assert.Equal("""{"change_id":"c-1","table":"orders"}""", frame.Data);
    }

    [Fact]
    public async Task ReadFrameAsync_StopsAtBlankLineBoundary()
    {
        using var reader = Lines(
            "event: change",
            """data: {"change_id":"c-1"}""",
            "",
            "event: change",
            """data: {"change_id":"c-2"}""",
            "");

        var first = await SseFrameParser.ReadFrameAsync(reader);
        var second = await SseFrameParser.ReadFrameAsync(reader);

        Assert.Equal("""{"change_id":"c-1"}""", first!.Data);
        Assert.Equal("""{"change_id":"c-2"}""", second!.Data);
    }

    [Fact]
    public async Task ReadFrameAsync_IncompleteFrameAtSourceExhaustion_ReturnsNull()
    {
        // Keine abschließende Leerzeile — die Quelle endet mitten im Frame.
        using var reader = new StringReader("event: change\ndata: {\"change_id\":\"c-1\"}");

        var frame = await SseFrameParser.ReadFrameAsync(reader);

        Assert.Null(frame);
    }

    [Fact]
    public async Task ReadFrameAsync_ExhaustedSourceWithNoBufferedContent_ReturnsNull()
    {
        using var reader = new StringReader("");

        var frame = await SseFrameParser.ReadFrameAsync(reader);

        Assert.Null(frame);
    }

    [Fact]
    public async Task ReadFrameAsync_LeadingBlankLinesAreSkipped()
    {
        using var reader = Lines("", "", "event: change", "data: {}", "");

        var frame = await SseFrameParser.ReadFrameAsync(reader);

        Assert.NotNull(frame);
        Assert.Equal("change", frame!.Name);
        Assert.Equal("{}", frame.Data);
    }

    [Fact]
    public async Task ReadFrameAsync_MissingEventNameDefaultsToEmptyString()
    {
        using var reader = Lines("""data: {"change_id":"c-1"}""", "");

        var frame = await SseFrameParser.ReadFrameAsync(reader);

        Assert.NotNull(frame);
        Assert.Equal("", frame!.Name);
        Assert.Equal("""{"change_id":"c-1"}""", frame.Data);
    }
}
