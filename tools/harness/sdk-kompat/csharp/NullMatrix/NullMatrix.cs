using Grpc.Core;
using PgChangeFeed.Client.Grpc;
using PgChangeFeed.Client.Http;

namespace Kompat.NullMatrix;

// Jede Zeile mit der Marke `CASE <Name>` ist ein Übersetzungsfall; das
// Treiber-Skript run.sh ordnet die Fehlerzeilen des Compilers nach Zeilennummer
// dem Fall zu. Alle Fälle nutzen nur Formen, die unter 0.5.0 existieren.

// Fremde Unterklassen der abstrakten HTTP-Basis.
internal sealed class H1 : PgChangeFeedException { public H1() : base(500, "x") { } } // CASE http-basis-2-argumente
internal sealed class H2 : PgChangeFeedException { public H2() : base(500, "x", new Exception()) { } } // CASE http-basis-3-argumente-exception
internal sealed class H3 : PgChangeFeedException { public H3() : base(500, "x", null) { } } // CASE http-basis-3-argumente-null-literal
internal sealed class H4 : PgChangeFeedException { public H4() : base(500, "x", (Exception)null) { } } // CASE http-basis-3-argumente-null-cast-exception
internal sealed class H5 : PgChangeFeedException { public H5() : base(500, "x", innerException: null) { } } // CASE http-basis-3-argumente-benannt-innerException-null
internal sealed class H6 : PgChangeFeedException { public H6() : base(statusCode: 500, message: "x") { } } // CASE http-basis-benannte-argumente

// Fremde Unterklassen der abstrakten gRPC-Basis.
internal sealed class G1 : PgChangeFeedGrpcException { public G1() : base(StatusCode.Unknown, "x", null) { } } // CASE grpc-basis-null-literal
internal sealed class G2 : PgChangeFeedGrpcException { public G2() : base(StatusCode.Unknown, "x", (RpcException)null) { } } // CASE grpc-basis-null-cast-rpcexception
internal sealed class G3 : PgChangeFeedGrpcException { public G3() : base(StatusCode.Unknown, "x", innerException: null) { } } // CASE grpc-basis-benannt-innerException-null

internal static class Aufrufe
{
    public static void Alle()
    {
        _ = new PgChangeFeedBadRequestException(400, "x"); // CASE http-blatt-2-argumente
        _ = new PgChangeFeedMalformedResponseException(200, "x", null); // CASE http-malformed-null-literal
        _ = new PgChangeFeedMalformedResponseException(200, "x", (Exception)null); // CASE http-malformed-null-cast-exception
        _ = new PgChangeFeedGrpcNotFoundException("x", null); // CASE grpc-blatt-null-literal
        _ = new PgChangeFeedGrpcUnexpectedStatusException(StatusCode.Unavailable, "x", null); // CASE grpc-blatt-unexpected-null-literal
    }
}
