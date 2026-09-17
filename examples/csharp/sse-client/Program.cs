using System.Net;
using System.Net.Http.Headers;

namespace CdcExamples.Sse;

/// <summary>
/// Command sse-client ist ein öffentliches Beispiel für den Live-Change-Stream
/// über Server-Sent-Events (<c>LH-FA-SST-008</c>, <c>ADR-0061</c>,
/// <c>ADR-0090</c>): es öffnet den Endpunkt <c>GET /changes/stream</c> real
/// gegen den laufenden Feed-Container und gibt jedes Event aus. Startform ist
/// ein Container-Aufruf, kein Host-Aufruf (<c>ADR-0087</c> Festlegung 3);
/// Adresse und Token kommen aus <c>CDC_HTTP_ADDR</c>/<c>CDC_API_TOKEN_READER</c>
/// und lassen sich per Flag übersteuern (<c>--addr</c>, <c>--token</c>).
///
/// Form-Vorbild: <c>examples/sse-client</c> (Go). Dieses Programm trägt keine
/// Zustandsmaschine: der Stream kennt kein Replay (<c>ADR-0061</c>), verpasste
/// Changes holt der bestehende Lesezugriffsweg nach.
/// </summary>
internal static class Program
{
    private static int Main(string[] args)
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

        if (string.IsNullOrEmpty(cfg.Addr))
        {
            Console.Error.WriteLine("sse-client: keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder --addr) ist noetig, um den Stream zu oeffnen");
            return 2;
        }
        if (string.IsNullOrEmpty(cfg.Token))
        {
            Console.Error.WriteLine("sse-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder --token) ist noetig, um ueber die reader-Rechtsklasse zu lesen");
            return 2;
        }

        var streamUrl = SseStream.StreamUrl(cfg.Addr);

        // Der Stream bleibt offen, bis die Verbindung endet: der Aufruf läuft
        // ohne Antwort-Frist (Muster: examples/sse-client/main.go).
        using var httpClient = new System.Net.Http.HttpClient { Timeout = Timeout.InfiniteTimeSpan };
        using var request = new HttpRequestMessage(HttpMethod.Get, streamUrl);
        request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", cfg.Token);

        HttpResponseMessage response;
        try
        {
            response = httpClient.Send(request, HttpCompletionOption.ResponseHeadersRead);
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine($"sse-client: Stream oeffnen ({streamUrl}): {ex.Message}");
            return 1;
        }

        using (response)
        {
            if (response.StatusCode != HttpStatusCode.OK)
            {
                Console.Error.WriteLine($"sse-client: Stream-Oeffnung endete mit HTTP-Status {(int)response.StatusCode} (Erwartung 200)");
                return 1;
            }

            using var stream = response.Content.ReadAsStream();
            using var reader = new StreamReader(stream);
            string? Next() => reader.ReadLine();

            while (true)
            {
                Event? ev;
                try
                {
                    ev = SseStream.ReadEvent(Next);
                }
                catch (IOException ex)
                {
                    // Das Ende der Quelle trägt zwei Ausgänge — ein Lesefehler
                    // oder eine geschlossene Verbindung. Beide werden
                    // gemeldet; der Prozess endet mit Nicht-Null.
                    Console.Error.WriteLine($"sse-client: Stream endete mit Fehler: {ex.Message}");
                    return 1;
                }
                if (ev is null)
                {
                    Console.Error.WriteLine("sse-client: Stream wurde geschlossen");
                    return 1;
                }
                Console.WriteLine($"sse-client: {ev.Name} {ev.Data}");
            }
        }
    }
}
