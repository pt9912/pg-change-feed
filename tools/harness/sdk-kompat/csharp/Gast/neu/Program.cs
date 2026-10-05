using Grpc.Core;
using PgChangeFeed.Client.Grpc;
using PgChangeFeed.Client.Http;

namespace Kompat.Gast;

// Gegenrichtung: ein gegen die Bibliothek mit Meldungscode gebauter Aufrufer
// nutzt die Konstruktoren mit Meldungscode und liest die Eigenschaft. Unter der
// Bibliothek 0.5.0 gibt es diese Mitglieder nicht.
internal static class Program
{
    private static int aufrufe;

    private static void Pruefe(bool ok, string was)
    {
        if (!ok)
        {
            throw new InvalidOperationException("Abweichung bei " + was);
        }

        aufrufe++;
    }

    private static void Lauf()
    {
        var http = new PgChangeFeedBadRequestException(400, "bad", "PCF-E0001");
        Pruefe(http.MessageCode == "PCF-E0001", "HTTP MessageCode");

        var rpc = new RpcException(new Status(StatusCode.NotFound, "rpc"));
        var grpc = new PgChangeFeedGrpcNotFoundException("m", rpc, "PCF-E0002");
        Pruefe(grpc.MessageCode == "PCF-E0002", "gRPC MessageCode");
    }

    private static int Main(string[] args)
    {
        var schritt = args.Length > 0 ? args[0] : "?";
        try
        {
            Lauf();
            Console.WriteLine($"KOMPAT csharp {schritt}: {aufrufe} Aufrufe ok");
            return 0;
        }
        catch (Exception ausnahme)
        {
            Console.WriteLine($"KOMPAT csharp {schritt}: Ausnahme {ausnahme.GetType().Name}: {ausnahme.Message}");
            return 1;
        }
    }
}
