using Grpc.Core;
using PgChangeFeed.Client.Grpc;
using PgChangeFeed.Client.Http;

namespace Kompat.Gast;

// Fremde Unterklasse der HTTP-Basis: ruft beide geschützten Konstruktoren der
// 0.5.x-Form auf (ohne und mit Ursache).
internal sealed class GastHttpFehler : PgChangeFeedException
{
    public GastHttpFehler(int statusCode, string message) : base(statusCode, message) { }

    public GastHttpFehler(int statusCode, string message, Exception cause) : base(statusCode, message, cause) { }
}

// Fremde Unterklasse der gRPC-Basis: der geschützte Konstruktor der 0.5.x-Form.
internal sealed class GastGrpcFehler : PgChangeFeedGrpcException
{
    public GastGrpcFehler(StatusCode statusCode, string message, RpcException cause) : base(statusCode, message, cause) { }
}

// Benutzt jede öffentliche Konstruktor- und Lesefläche der 0.5.x-Fehlertypen
// beider Hierarchien in der Form, die ein 0.5.0 gebundener Aufrufer trägt.
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

    private static void PruefeHttp(PgChangeFeedException fehler, int status, string text, Exception? ursache, string was)
    {
        Pruefe(fehler.StatusCode == status, was + " StatusCode");
        Pruefe(fehler.Message == text, was + " Message");
        Pruefe(ReferenceEquals(fehler.InnerException, ursache), was + " InnerException");
    }

    private static void PruefeGrpc(PgChangeFeedGrpcException fehler, StatusCode status, string text, RpcException ursache, string was)
    {
        Pruefe(fehler.StatusCode == status, was + " StatusCode");
        Pruefe(fehler.Message == text, was + " Message");
        Pruefe(ReferenceEquals(fehler.InnerException, ursache), was + " InnerException");
    }

    private static void Lauf()
    {
        var ursache = new InvalidOperationException("ursache");
        PruefeHttp(new PgChangeFeedBadRequestException(400, "bad"), 400, "bad", null, "BadRequest(int,string)");
        PruefeHttp(new PgChangeFeedUnauthorizedException(401, "unauth"), 401, "unauth", null, "Unauthorized(int,string)");
        PruefeHttp(new PgChangeFeedForbiddenException(403, "forbidden"), 403, "forbidden", null, "Forbidden(int,string)");
        PruefeHttp(new PgChangeFeedNotFoundException(404, "missing"), 404, "missing", null, "NotFound(int,string)");
        PruefeHttp(new PgChangeFeedServerErrorException(500, "server"), 500, "server", null, "ServerError(int,string)");
        PruefeHttp(new PgChangeFeedUnexpectedStatusException(418, "teapot"), 418, "teapot", null, "UnexpectedStatus(int,string)");
        PruefeHttp(new PgChangeFeedMalformedResponseException(200, "mal"), 200, "mal", null, "Malformed(int,string)");
        PruefeHttp(new PgChangeFeedMalformedResponseException(200, "mal", ursache), 200, "mal", ursache, "Malformed(int,string,Exception)");
        PruefeHttp(new GastHttpFehler(500, "fremd"), 500, "fremd", null, "Unterklasse base(int,string)");
        PruefeHttp(new GastHttpFehler(500, "fremd", ursache), 500, "fremd", ursache, "Unterklasse base(int,string,Exception)");

        var rpc = new RpcException(new Status(StatusCode.NotFound, "rpc"));
        PruefeGrpc(new PgChangeFeedGrpcInvalidArgumentException("m", rpc), StatusCode.InvalidArgument, "m", rpc, "GrpcInvalidArgument(string,RpcException)");
        PruefeGrpc(new PgChangeFeedGrpcUnauthenticatedException("m", rpc), StatusCode.Unauthenticated, "m", rpc, "GrpcUnauthenticated(string,RpcException)");
        PruefeGrpc(new PgChangeFeedGrpcPermissionDeniedException("m", rpc), StatusCode.PermissionDenied, "m", rpc, "GrpcPermissionDenied(string,RpcException)");
        PruefeGrpc(new PgChangeFeedGrpcNotFoundException("m", rpc), StatusCode.NotFound, "m", rpc, "GrpcNotFound(string,RpcException)");
        PruefeGrpc(new PgChangeFeedGrpcInternalException("m", rpc), StatusCode.Internal, "m", rpc, "GrpcInternal(string,RpcException)");
        PruefeGrpc(new PgChangeFeedGrpcUnexpectedStatusException(StatusCode.Unavailable, "m", rpc), StatusCode.Unavailable, "m", rpc, "GrpcUnexpectedStatus(StatusCode,string,RpcException)");
        PruefeGrpc(new GastGrpcFehler(StatusCode.Aborted, "fremd", rpc), StatusCode.Aborted, "fremd", rpc, "Unterklasse base(StatusCode,string,RpcException)");

        // Fangen über die Basistypen: ein 0.5.0 gebundener Aufrufer fängt je Hierarchie über die Basis.
        var gefangenHttp = false;
        try
        {
            throw new PgChangeFeedNotFoundException(404, "wurf");
        }
        catch (PgChangeFeedException fehler)
        {
            gefangenHttp = fehler is PgChangeFeedNotFoundException;
        }

        Pruefe(gefangenHttp, "catch PgChangeFeedException");

        var gefangenGrpc = false;
        try
        {
            throw new PgChangeFeedGrpcNotFoundException("wurf", rpc);
        }
        catch (PgChangeFeedGrpcException fehler)
        {
            gefangenGrpc = fehler is PgChangeFeedGrpcNotFoundException;
        }

        Pruefe(gefangenGrpc, "catch PgChangeFeedGrpcException");
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
