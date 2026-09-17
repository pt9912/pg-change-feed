namespace CdcExamples.Http;

/// <summary>
/// Command http-client ist ein öffentliches Beispiel für den
/// Anfrage/Antwort-Zugriff über die HTTP-/JSON-API (<c>LH-FA-SST-006</c>,
/// <c>ADR-0057</c>, <c>ADR-0090</c>): es ruft den <c>reader</c>-Endpunkt
/// <c>GET /tables</c> real gegen den laufenden Feed-Container auf und gibt
/// die Antwort aus. Startform ist ein Container-Aufruf, kein Host-Aufruf
/// (<c>ADR-0087</c> Festlegung 3); Adresse und Token kommen aus
/// <c>CDC_HTTP_ADDR</c>/<c>CDC_API_TOKEN_READER</c> und lassen sich per Flag
/// übersteuern (<c>--addr</c>, <c>--token</c>, <c>--source</c>,
/// <c>--publication</c>).
///
/// Form-Vorbild: <c>examples/http-client</c> (Go). Dieses Programm trägt
/// keine Zustandsmaschine: es stellt eine Anfrage und endet.
/// </summary>
internal static class Program
{
    private static readonly TimeSpan RequestTimeout = TimeSpan.FromSeconds(10);

    private static async Task<int> Main(string[] args)
    {
        Config cfg;
        try
        {
            cfg = Cli.Parse(args, Environment.GetEnvironmentVariable);
        }
        catch (ArgumentException ex)
        {
            Console.Error.WriteLine(ex.Message);
            return 2;
        }

        // `CDC_HTTP_ADDR` ungesetzt heißt: die HTTP-API ist deaktiviert — die
        // Zeile dazu steht in §5 *Konfiguration* des Handbuchs. Es gibt dann
        // keinen Endpunkt, und das Beispiel endet mit dieser Meldung.
        if (string.IsNullOrEmpty(cfg.Addr))
        {
            Console.Error.WriteLine("http-client: keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder --addr) ist noetig, um die Verwaltungs-API zu erreichen");
            return 2;
        }
        if (string.IsNullOrEmpty(cfg.Token))
        {
            Console.Error.WriteLine("http-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder --token) ist noetig, um ueber die reader-Rechtsklasse zu lesen");
            return 2;
        }
        if (string.IsNullOrEmpty(cfg.Source) || string.IsNullOrEmpty(cfg.Publication))
        {
            Console.Error.WriteLine("http-client: --source und --publication sind Pflicht — sie sind die zwei Pflichtfelder von GET /tables");
            return 2;
        }

        using var httpClient = new System.Net.Http.HttpClient { Timeout = RequestTimeout };
        try
        {
            var body = await TablesClient.ListTablesAsync(httpClient, cfg);
            Console.WriteLine(body);
            return 0;
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine($"http-client: GET /tables fehlgeschlagen: {ex.Message}");
            return 1;
        }
    }
}
