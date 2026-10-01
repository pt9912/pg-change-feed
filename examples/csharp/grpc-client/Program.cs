using Cdc.Administration.V1;
using Cdc.Stream.V1;
using Grpc.Core;
using Grpc.Net.Client;

namespace CdcExamples.Grpc;

/// <summary>
/// Command grpc-client ist ein öffentliches Beispiel für die vollständige
/// gRPC-Fläche (<c>LH-FA-SST-006</c>): das <c>--verb</c>-Flag ruft eine von
/// zwölf dokumentierten Fähigkeiten real gegen den laufenden Feed-Container
/// auf — der Live-Change-Stream (Default-Verb <c>stream</c>, unverändert die
/// ursprüngliche Aufrufform) und die elf unären RPCs des
/// <c>Administration</c>-Diensts. Startform ist ein Container-Aufruf, kein
/// Host-Aufruf; die Zugriffs-Abschnitte des Benutzerhandbuchs sind „Zugriff
/// über den gRPC-Change-Stream" und „Zugriff über die gRPC-Verwaltungs-API"
/// (<c>docs/user/benutzerhandbuch.md</c>).
///
/// Form-Vorbild: <c>examples/grpc-client</c> (Go). Beide Stubs entstehen im
/// Bau aus den über den benannten Zusatzkontext gelesenen <c>.proto</c>-
/// Dateien — sie liegen nicht im committeten Baum. Dieses Programm trägt
/// keine Zustandsmaschine: die elf RPCs stellen je eine Anfrage und enden,
/// der Stream bleibt offen und kennt kein Replay — verpasste Changes holt
/// der bestehende Lesezugriffsweg (<c>ReadChanges</c>, <c>GET /changes</c>)
/// nach.
/// </summary>
internal static class Program
{
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
            Console.Error.WriteLine($"grpc-client: {validationError}");
            return 2;
        }

        // Der Feed-Container spricht Klartext-gRPC (kein TLS, dieselbe Wahl
        // wie der Go-Client mit `insecure.NewCredentials()`) — der Switch
        // erlaubt HTTP/2 ohne Transportverschlüsselung (h2c) für den
        // `HttpClient`, den `GrpcChannel.ForAddress` intern benutzt.
        AppContext.SetSwitch("System.Net.Http.SocketsHttpHandler.Http2UnencryptedSupport", true);

        using var channel = GrpcChannel.ForAddress($"http://{cfg.Addr}");

        if (cfg.Verb == "stream")
        {
            return await RunStreamAsync(channel, cfg).ConfigureAwait(false);
        }

        var client = new Administration.AdministrationClient(channel);
        try
        {
            var output = await Dispatcher.DispatchAdminAsync(client, cfg).ConfigureAwait(false);
            Console.WriteLine(output);
            return 0;
        }
        catch (RpcException ex)
        {
            Console.Error.WriteLine($"grpc-client: --verb={cfg.Verb} fehlgeschlagen: {ex.Status}");
            return 1;
        }
    }

    /// <summary>
    /// RunStreamAsync öffnet den Server-Streaming-RPC
    /// <c>ChangeStream/StreamChanges</c> real gegen den laufenden
    /// Feed-Container und gibt jede empfangene Nachricht aus, bis die
    /// Verbindung endet; die Anfrage baut <see cref="StreamRequest.Build"/>.
    /// </summary>
    private static async Task<int> RunStreamAsync(GrpcChannel channel, Config cfg)
    {
        var client = new ChangeStream.ChangeStreamClient(channel);
        var headers = CallMetadata.Headers(cfg.Token);

        try
        {
            using var call = client.StreamChanges(StreamRequest.Build(cfg), headers: headers);
            await foreach (var change in call.ResponseStream.ReadAllAsync())
            {
                Console.WriteLine(Format.FormatChange(change));
            }
        }
        catch (RpcException ex)
        {
            Console.Error.WriteLine($"grpc-client: Stream endete: {ex.Status}");
            return 1;
        }

        Console.Error.WriteLine("grpc-client: Stream endete");
        return 1;
    }
}
