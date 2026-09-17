namespace CdcExamples.Http;

/// <summary>
/// TablesUrlBuilder baut die Lese-Adresse der Tabellen-Auflistung über die
/// HTTP-/JSON-API (<c>LH-FA-SST-006</c>): ein <c>GET</c> auf den
/// <c>reader</c>-Endpunkt <c>/tables</c> mit seinen zwei Pflichtfeldern
/// <c>source</c> und <c>publication</c>. <c>addr</c> ist die Horch-Adresse
/// des Feed-Containers (<c>CDC_HTTP_ADDR</c>, Form <c>host:port</c>). Reine
/// Funktion, Form-Vorbild: <c>examples/http-client/tables.go</c>.
/// </summary>
public static class TablesUrlBuilder
{
    public static string Build(string addr, string source, string publication)
    {
        var query = string.Join(
            '&',
            $"publication={Uri.EscapeDataString(publication)}",
            $"source={Uri.EscapeDataString(source)}");
        return $"http://{addr}/tables?{query}";
    }
}
