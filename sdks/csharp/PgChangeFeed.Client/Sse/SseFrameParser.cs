namespace PgChangeFeed.Client.Sse;

/// <summary>
/// One complete SSE frame — the <c>event:</c> name and the <c>data:</c>
/// payload accumulated up to the blank line that terminates it (SPEC-021).
/// </summary>
public sealed record SseFrame(string Name, string Data);

/// <summary>
/// Zerlegt Server-Sent-Events-Frames aus einer Zeilenquelle in
/// <see cref="SseFrame"/>-Werte — reine Zeilen-Verarbeitung über die
/// Runtime-Standardbibliothek, kein Fremdmodul nötig (dasselbe
/// „keine neue Abhängigkeit"-Muster wie das gelesene, nicht importierte
/// Draht-Kenntnis-Vorbild <c>examples/csharp/sse-client/SseStream.cs</c>,
/// ADR-0106 Festlegung 2).
///
/// Ein Frame endet mit der Leerzeile, die der Server nach der <c>event:</c>-
/// und der <c>data:</c>-Zeile schreibt (SPEC-021); sie trennt zwei
/// aufeinanderfolgende Events. Ist die Quelle vor dem Frame-Abschluss
/// erschöpft (<see cref="TextReader.ReadLineAsync(System.Threading.CancellationToken)"/>
/// liefert <c>null</c>), liefert <see cref="ReadFrameAsync"/> ebenfalls
/// <c>null</c>, und ein begonnenes Frame wird verworfen — ein
/// unvollständiges Frame ist kein Event.
/// </summary>
public static class SseFrameParser
{
    private const string EventPrefix = "event: ";
    private const string DataPrefix = "data: ";

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
