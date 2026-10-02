namespace CdcExamples.Sse;

/// <summary>
/// Event trägt ein Frame des Change-Streams (<c>LH-FA-SST-008</c>): den Namen
/// aus der <c>event:</c>-Zeile und die Nutzlast aus der <c>data:</c>-Zeile.
/// Der Server schreibt je Change ein JSON-Objekt mit den zehn
/// Nachrichtenfeldern in die <c>data:</c>-Zeile (<c>ADR-0061</c>).
/// </summary>
public sealed record Event(string Name, string Data);

/// <summary>
/// SseStream baut die Stream-Adresse und zerlegt die SSE-Frames — reine
/// Funktionen, netzlos testbar. Form-Vorbild: <c>examples/sse-client/stream.go</c>.
/// **Keine neue Abhängigkeit** (<c>ADR-0090</c> Festlegung 1): das Zerlegen
/// eines Frames ist Zeilen-Verarbeitung über die Runtime-Standardbibliothek,
/// kein Fremdmodul (etwa <c>System.Net.ServerSentEvents</c>) nötig.
/// </summary>
public static class SseStream
{
    /// <summary>
    /// StreamUrl baut die Adresse des SSE-Endpunkts (<c>LH-FA-SST-008</c>): ein
    /// <c>GET</c> auf <c>/changes/stream</c>. <paramref name="addr"/> ist die
    /// Horch-Adresse des Feed-Containers (<c>CDC_HTTP_ADDR</c>, Form
    /// <c>host:port</c>). Ein gesetztes <paramref name="target"/> erscheint als
    /// maskierter Query-Parameter und wählt das Zustellziel einer Change; ein
    /// leeres lässt die Adresse ohne Query. <paramref name="schema"/> und
    /// <paramref name="table"/> erscheinen ebenso nur gesetzt als
    /// Query-Parameter, je optional und unabhängig, und wirken mit dem Ziel als
    /// Konjunktion.
    /// </summary>
    public static string StreamUrl(string addr, string target = "", string schema = "", string table = "")
    {
        var query = new List<string>();
        if (schema.Length > 0)
        {
            query.Add($"schema={Uri.EscapeDataString(schema)}");
        }
        if (table.Length > 0)
        {
            query.Add($"table={Uri.EscapeDataString(table)}");
        }
        if (target.Length > 0)
        {
            query.Add($"target={Uri.EscapeDataString(target)}");
        }
        return query.Count == 0
            ? $"http://{addr}/changes/stream"
            : $"http://{addr}/changes/stream?{string.Join("&", query)}";
    }

    /// <summary>
    /// ReadEvent liest ein vollständiges Frame über <paramref name="next"/> und
    /// liefert es zurück. Ein Frame endet mit der Leerzeile, die der Server
    /// nach der <c>event:</c>- und der <c>data:</c>-Zeile schreibt
    /// (<c>ADR-0061</c>); sie trennt zwei aufeinanderfolgende Events. Ist die
    /// Quelle vor dem Frame-Abschluss erschöpft (<paramref name="next"/>
    /// liefert <c>null</c>), liefert ReadEvent <c>null</c>, und ein begonnenes
    /// Frame wird verworfen — ein unvollständiges Frame ist kein Event.
    /// </summary>
    public static Event? ReadEvent(Func<string?> next)
    {
        string? name = null;
        string? data = null;
        while (true)
        {
            var line = next();
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
                return new Event(name ?? "", data ?? "");
            }
            if (line.StartsWith("event: ", StringComparison.Ordinal))
            {
                name = line["event: ".Length..];
            }
            else if (line.StartsWith("data: ", StringComparison.Ordinal))
            {
                data = line["data: ".Length..];
            }
        }
    }
}
