namespace PgChangeFeed.Client.Sse;

/// <summary>
/// One complete SSE frame — the <c>event:</c> name and the <c>data:</c>
/// payload accumulated up to the blank line that terminates it.
/// </summary>
public sealed record SseFrame(string Name, string Data);

/// <summary>
/// Splits Server-Sent-Events frames from a line source into
/// <see cref="SseFrame"/> values — plain line processing over the runtime's
/// standard library, no further dependency.
///
/// A frame ends with the blank line the server writes after the <c>event:</c>
/// and the <c>data:</c> line; it separates two consecutive events. When the
/// source is exhausted before a frame is complete
/// (<see cref="TextReader.ReadLineAsync(System.Threading.CancellationToken)"/>
/// returns <c>null</c>), <see cref="ReadFrameAsync"/> also returns
/// <c>null</c>, and a frame already begun is discarded — an incomplete frame
/// is not an event.
/// </summary>
public static class SseFrameParser
{
    private const string EventPrefix = "event: ";
    private const string DataPrefix = "data: ";

    /// <summary>
    /// Reads the next complete frame, or returns <c>null</c> when the source
    /// ends before a frame is complete.
    /// </summary>
    public static async Task<SseFrame?> ReadFrameAsync(
        TextReader reader, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(reader);

        string? name = null;
        string? data = null;
        while (true)
        {
            var line = await reader.ReadLineAsync(cancellationToken).ConfigureAwait(false);
            if (line is null)
            {
                return null;
            }
            if (line.Length == 0)
            {
                if (name is null && data is null)
                {
                    continue;
                }
                return new SseFrame(name ?? "", data ?? "");
            }
            if (line.StartsWith(EventPrefix, StringComparison.Ordinal))
            {
                name = line[EventPrefix.Length..];
            }
            else if (line.StartsWith(DataPrefix, StringComparison.Ordinal))
            {
                data = line[DataPrefix.Length..];
            }
        }
    }
}
