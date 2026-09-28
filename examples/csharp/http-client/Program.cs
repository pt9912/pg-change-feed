namespace CdcExamples.Http;

/// <summary>
/// Command http-client ist ein öffentliches Beispiel für den
/// Anfrage/Antwort-Zugriff über die HTTP-/JSON-API (<c>LH-FA-SST-006</c>):
/// das <c>--verb</c>-Flag ruft eine von zehn dokumentierten Fähigkeiten real
/// gegen den laufenden Feed-Container auf und gibt die Antwort aus. Startform
/// ist ein Container-Aufruf, kein Host-Aufruf (<c>ADR-0087</c> Festlegung 3);
/// Default-Verb ist <c>tables</c> — die bestehende Startform
/// <c>ARGS="--source &lt;quelle&gt; --publication &lt;publication&gt;"</c>
/// bleibt unverändert funktionsfähig.
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

        var validationError = Validator.Validate(cfg);
        if (validationError is not null)
        {
            Console.Error.WriteLine($"http-client: {validationError}");
            return 2;
        }

        using var httpClient = new System.Net.Http.HttpClient { Timeout = RequestTimeout };
        try
        {
            var body = await Dispatcher.DispatchAsync(httpClient, cfg).ConfigureAwait(false);
            Console.WriteLine(body);
            return 0;
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine($"http-client: --verb={cfg.Verb} fehlgeschlagen: {ex.Message}");
            return 1;
        }
    }
}
