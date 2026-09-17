using Cdc.Stream.V1;
using Grpc.Core;
using Grpc.Net.Client;

namespace CdcExamples.Grpc;

/// <summary>
/// Command grpc-client ist ein öffentliches Beispiel für den Live-Change-Stream
/// über gRPC (<c>LH-FA-SST-008</c>, <c>ADR-0060</c>, <c>ADR-0090</c>): es
/// öffnet den Server-Streaming-RPC <c>ChangeStream/StreamChanges</c> real
/// gegen den laufenden Feed-Container und gibt jede empfangene Nachricht aus.
/// Startform ist ein Container-Aufruf, kein Host-Aufruf (<c>ADR-0087</c>
/// Festlegung 3); Adresse und Token kommen aus
/// <c>CDC_GRPC_ADDR</c>/<c>CDC_API_TOKEN_READER</c> und lassen sich per Flag
/// übersteuern (<c>--addr</c>, <c>--token</c>).
///
/// Form-Vorbild: <c>examples/grpc-client</c> (Go). Der Stub entsteht **im
/// Bau** aus der über den benannten Zusatzkontext gelesenen <c>.proto</c> —
/// er liegt nicht im committeten Baum (<c>ADR-0090</c> Festlegung 3). Dieses
/// Programm trägt keine Zustandsmaschine: der Stream kennt kein Replay
/// (<c>ADR-0060</c>), verpasste Changes holt der bestehende Lesezugriffsweg
/// nach.
/// </summary>
internal static class Program
{
    private const string AuthorizationMetadataKey = "authorization";
    private const string BearerPrefix = "Bearer ";

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

        if (string.IsNullOrEmpty(cfg.Addr))
        {
            Console.Error.WriteLine("grpc-client: keine gRPC-Adresse gesetzt — CDC_GRPC_ADDR (oder --addr) ist noetig, um den Stream zu oeffnen");
            return 2;
        }
        if (string.IsNullOrEmpty(cfg.Token))
        {
            Console.Error.WriteLine("grpc-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder --token) ist noetig, um ueber die reader-Rechtsklasse zu lesen");
            return 2;
        }

        // Der Feed-Container spricht Klartext-gRPC (kein TLS, dieselbe Wahl
        // wie der Go-Client mit `insecure.NewCredentials()`) — der Switch
        // erlaubt HTTP/2 ohne Transportverschlüsselung (h2c) für den
        // `HttpClient`, den `GrpcChannel.ForAddress` intern benutzt.
        AppContext.SetSwitch("System.Net.Http.SocketsHttpHandler.Http2UnencryptedSupport", true);

        using var channel = GrpcChannel.ForAddress($"http://{cfg.Addr}");
        var client = new ChangeStream.ChangeStreamClient(channel);

        var headers = new Metadata { { AuthorizationMetadataKey, BearerPrefix + cfg.Token } };

        try
        {
            using var call = client.StreamChanges(new StreamChangesRequest(), headers: headers);
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
