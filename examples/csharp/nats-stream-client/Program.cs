using System.Text.Json;

using NATS.Client.Core;
using NATS.Net;

namespace CdcExamples.NatsStream;

/// <summary>
/// Command nats-stream-client ist ein öffentliches Beispiel für den
/// dritten, vollinhaltstragenden NATS-Zustellweg (<c>LH-FA-SST-008</c>,
/// <c>ADR-0100</c>, <c>SPEC-024</c>): es verbindet mit einem gültigen
/// <c>CDC_NATS_STREAM_TOKEN</c>, abonniert den Vollinhalts-Namensraum
/// <c>cdc.stream.&gt;</c> real gegen den laufenden Feed-Container und gibt
/// jede empfangene Change aus. Startform ist ein Container-Aufruf, kein
/// Host-Aufruf (<c>ADR-0087</c> Festlegung 3); NATS-URL und Token kommen aus
/// <c>CDC_NATS_URL</c>/<c>CDC_NATS_STREAM_TOKEN</c> und lassen sich per Flag
/// übersteuern (<c>--nats-url</c>, <c>--token</c>).
///
/// Form-Vorbild: <c>examples/nats-stream-client</c> (Go). Dieses Programm
/// trägt keine Zustandsmaschine: der Stream kennt kein Replay
/// (<c>ADR-0100</c>), verpasste Changes holt der bestehende Lesezugriffsweg
/// nach.
/// </summary>
internal static class Program
{
    private const string Subject = "cdc.stream.>";

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

        if (string.IsNullOrEmpty(cfg.NatsUrl))
        {
            Console.Error.WriteLine("nats-stream-client: keine NATS-URL gesetzt — CDC_NATS_URL (oder --nats-url) ist noetig, um den Vollinhalts-Stream zu oeffnen");
            return 2;
        }
        if (string.IsNullOrEmpty(cfg.Token))
        {
            Console.Error.WriteLine("nats-stream-client: kein Token gesetzt — CDC_NATS_STREAM_TOKEN (oder --token) ist noetig, um den dritten Zustellweg zu abonnieren");
            return 2;
        }

        var opts = new NatsOpts
        {
            Url = cfg.NatsUrl,
            AuthOpts = new NatsAuthOpts { Token = cfg.Token },
        };
        await using var nats = new NatsClient(opts);

        Console.WriteLine($"nats-stream-client: lauscht auf {Subject}");

        try
        {
            await foreach (var msg in nats.SubscribeAsync<byte[]>(Subject))
            {
                var change = JsonSerializer.Deserialize<StreamMessage>(msg.Data ?? []);
                if (change is null)
                {
                    Console.Error.WriteLine("nats-stream-client: Event nicht dekodierbar (leerer Payload)");
                    return 1;
                }
                Console.WriteLine(Format.FormatChange(change));
            }
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine($"nats-stream-client: Empfang beendet: {ex.Message}");
            return 1;
        }

        Console.Error.WriteLine("nats-stream-client: Empfang beendet");
        return 1;
    }
}
