using NATS.Net;

namespace CdcExamples.Nats;

/// <summary>
/// Command nats-client ist ein öffentliches Beispiel für den zweiseitigen
/// Zugriffsweg über das NATS-Wecksignal (<c>LH-FA-SST-007</c>,
/// <c>ADR-0079</c>, <c>ADR-0090</c>): es abonniert das tabellen-granulare
/// Subjekt <c>cdc.changes.&lt;source_id&gt;.&lt;schema&gt;.&lt;table&gt;</c>
/// (<c>SPEC-017</c>, <c>ADR-0056</c>) und holt beim Weckruf die Änderung
/// selbst über die HTTP-/JSON-API (<c>LH-FA-SST-006</c>) — das Signal trägt
/// per Vertrag einen leeren Payload (<c>ADR-0055</c>), es sagt nur „lies
/// erneut über den bestehenden Zugriffsweg". Startform ist ein
/// Container-Aufruf, kein Host-Aufruf (<c>ADR-0087</c> Festlegung 3);
/// NATS-URL, Adresse und Token kommen aus <c>CDC_NATS_URL</c>/
/// <c>CDC_HTTP_ADDR</c>/<c>CDC_API_TOKEN_READER</c> und lassen sich per Flag
/// übersteuern (<c>--nats-url</c>, <c>--addr</c>, <c>--token</c>,
/// <c>--source</c>, <c>--schema</c>, <c>--table</c>).
///
/// Form-Vorbild: <c>examples/nats-client</c> (Go). Dieses Programm trägt
/// keine Zustandsmaschine: eine Reconnect-/Dedup-/Rückstand-Logik ist bewusst
/// nicht seine Aufgabe (<c>SPEC-023</c>).
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

        // `CDC_HTTP_ADDR` ungesetzt heißt: die HTTP-API ist deaktiviert
        // (`SPEC-018`). Ohne sie kann der Weckruf keine Änderung holen — das
        // Beispiel scheitert sichtbar, statt still nichts zu tun (`ADR-0079`
        // Festlegung 2).
        if (string.IsNullOrEmpty(cfg.Addr))
        {
            Console.Error.WriteLine("nats-client: keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder --addr) ist noetig, um beim Weckruf die Aenderung zu holen");
            return 2;
        }
        if (string.IsNullOrEmpty(cfg.NatsUrl))
        {
            Console.Error.WriteLine("nats-client: keine NATS-URL gesetzt — CDC_NATS_URL (oder --nats-url) ist noetig, um das Wecksignal zu abonnieren");
            return 2;
        }
        if (string.IsNullOrEmpty(cfg.Source) || string.IsNullOrEmpty(cfg.Schema) || string.IsNullOrEmpty(cfg.Table))
        {
            Console.Error.WriteLine("nats-client: --source, --schema und --table sind Pflicht");
            return 2;
        }
        if (string.IsNullOrEmpty(cfg.Token))
        {
            Console.Error.WriteLine("nats-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder --token) ist noetig, um ueber die reader-Rechtsklasse zu lesen");
            return 2;
        }

        var subject = Subject.Build(cfg.Source, cfg.Schema, cfg.Table);

        await using var nats = new NatsClient(cfg.NatsUrl);

        Console.WriteLine($"nats-client: lauscht auf {subject} — beim naechsten Weckruf wird die Aenderung ueber {cfg.Addr} geholt");

        await using var subscription = nats.SubscribeAsync<byte[]>(subject).GetAsyncEnumerator();
        bool received;
        try
        {
            received = await subscription.MoveNextAsync();
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine($"nats-client: Empfang des Wecksignals fehlgeschlagen: {ex.Message}");
            return 1;
        }
        if (!received)
        {
            Console.Error.WriteLine("nats-client: Abonnement endete ohne Wecksignal");
            return 1;
        }

        var msg = subscription.Current;

        // Der Payload trägt keine Änderungsdaten (`SPEC-017`); die Änderung
        // kommt ausschließlich über den HTTP-Lesezugriff. Gedruckt wird die
        // Nutzdaten-Länge (`Data`), nicht `NatsMsg.Size` — das ist die
        // Protokoll-Größe des Rahmens (Subjekt + Nutzdaten) und damit für
        // einen leeren Payload nicht 0.
        Console.WriteLine($"nats-client: Weckruf auf {msg.Subject} (Payload {msg.Data?.Length ?? 0} Byte) — hole die Aenderung ueber HTTP");

        using var httpClient = new System.Net.Http.HttpClient { Timeout = RequestTimeout };
        try
        {
            var body = await ChangesClient.FetchAsync(httpClient, cfg);
            Console.WriteLine(body);
            return 0;
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine($"nats-client: HTTP-Abfrage fehlgeschlagen: {ex.Message}");
            return 1;
        }
    }
}
