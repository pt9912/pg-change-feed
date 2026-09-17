namespace CdcExamples.Nats;

/// <summary>
/// ChangesUrlBuilder baut die Lese-Adresse der Änderungen einer
/// Quelle/Schema/Tabelle über die HTTP-/JSON-API (<c>LH-FA-SST-006</c>): ein
/// <c>GET</c> auf den <c>reader</c>-Endpunkt <c>/changes</c> mit den drei
/// Filtern aus dem Wecksignal-Subjekt (<c>source</c>, <c>schema</c>,
/// <c>table</c>) — der zweiseitige Ablauf aus <c>ADR-0079</c>.
/// <paramref name="addr"/> ist die Horch-Adresse des Feed-Containers
/// (<c>CDC_HTTP_ADDR</c>, Form <c>host:port</c>). Reine Funktion,
/// Form-Vorbild: <c>examples/nats-client/subject.go</c> (<c>ChangesURL</c>).
/// </summary>
public static class ChangesUrlBuilder
{
    public static string Build(string addr, string source, string schema, string table)
    {
        var query = string.Join(
            '&',
            $"schema={Uri.EscapeDataString(schema)}",
            $"source={Uri.EscapeDataString(source)}",
            $"table={Uri.EscapeDataString(table)}");
        return $"http://{addr}/changes?{query}";
    }
}
