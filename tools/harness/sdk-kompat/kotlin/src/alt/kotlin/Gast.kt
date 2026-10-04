package kompat

import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcInternalException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcInvalidArgumentException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcNotFoundException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcPermissionDeniedException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcUnauthenticatedException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcUnexpectedStatusException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedBadRequestException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedForbiddenException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedMalformedResponseException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedNotFoundException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedServerErrorException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnauthorizedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnexpectedStatusException
import io.grpc.Status
import io.grpc.StatusException
import kotlin.system.exitProcess

// Benutzt jede öffentliche Konstruktor- und Lesefläche der 0.5.x-Fehlertypen
// beider Hierarchien in der Form, die ein gegen 0.5.0 übersetzter Aufrufer trägt.
private var aufrufe = 0

private fun pruefe(ok: Boolean, was: String) {
    check(ok) { "Abweichung bei $was" }
    aufrufe++
}

private fun pruefeHttp(fehler: PgChangeFeedException, status: Int, text: String, ursache: Throwable?, was: String) {
    pruefe(fehler.statusCode == status, "$was statusCode")
    pruefe(fehler.message == text, "$was message")
    pruefe(fehler.cause === ursache, "$was cause")
}

private fun pruefeGrpc(fehler: PgChangeFeedGrpcException, status: Status.Code, text: String, ursache: StatusException, was: String) {
    pruefe(fehler.statusCode == status, "$was statusCode")
    pruefe(fehler.message == text, "$was message")
    pruefe(fehler.cause === ursache, "$was cause")
}

// Erschöpfendes when über die versiegelte Basis: ein neuer Untertyp bräche die Übersetzung.
private fun gruppe(fehler: PgChangeFeedException): String = when (fehler) {
    is PgChangeFeedBadRequestException -> "400"
    is PgChangeFeedUnauthorizedException -> "401"
    is PgChangeFeedForbiddenException -> "403"
    is PgChangeFeedNotFoundException -> "404"
    is PgChangeFeedServerErrorException -> "500"
    is PgChangeFeedUnexpectedStatusException -> "sonst"
    is PgChangeFeedMalformedResponseException -> "malformed"
}

private fun grpcGruppe(fehler: PgChangeFeedGrpcException): String = when (fehler) {
    is PgChangeFeedGrpcInvalidArgumentException -> "invalid"
    is PgChangeFeedGrpcUnauthenticatedException -> "unauth"
    is PgChangeFeedGrpcPermissionDeniedException -> "denied"
    is PgChangeFeedGrpcNotFoundException -> "notfound"
    is PgChangeFeedGrpcInternalException -> "internal"
    is PgChangeFeedGrpcUnexpectedStatusException -> "sonst"
}

private fun lauf() {
    val ursache = IllegalStateException("ursache")
    pruefeHttp(PgChangeFeedBadRequestException(400, "bad"), 400, "bad", null, "BadRequest(Int,String)")
    pruefeHttp(PgChangeFeedUnauthorizedException(401, "unauth"), 401, "unauth", null, "Unauthorized(Int,String)")
    pruefeHttp(PgChangeFeedForbiddenException(403, "forbidden"), 403, "forbidden", null, "Forbidden(Int,String)")
    pruefeHttp(PgChangeFeedNotFoundException(404, "missing"), 404, "missing", null, "NotFound(Int,String)")
    pruefeHttp(PgChangeFeedServerErrorException(500, "server"), 500, "server", null, "ServerError(Int,String)")
    pruefeHttp(PgChangeFeedUnexpectedStatusException(418, "teapot"), 418, "teapot", null, "UnexpectedStatus(Int,String)")
    pruefeHttp(PgChangeFeedMalformedResponseException(200, "mal"), 200, "mal", null, "Malformed(Int,String)")
    pruefeHttp(PgChangeFeedMalformedResponseException(200, "mal", ursache), 200, "mal", ursache, "Malformed(Int,String,Throwable)")
    pruefe(gruppe(PgChangeFeedNotFoundException(404, "x")) == "404", "when über die Basis (HTTP)")

    val rpc = Status.NOT_FOUND.withDescription("rpc").asException()
    pruefeGrpc(PgChangeFeedGrpcInvalidArgumentException("m", rpc), Status.Code.INVALID_ARGUMENT, "m", rpc, "GrpcInvalidArgument(String,StatusException)")
    pruefeGrpc(PgChangeFeedGrpcUnauthenticatedException("m", rpc), Status.Code.UNAUTHENTICATED, "m", rpc, "GrpcUnauthenticated(String,StatusException)")
    pruefeGrpc(PgChangeFeedGrpcPermissionDeniedException("m", rpc), Status.Code.PERMISSION_DENIED, "m", rpc, "GrpcPermissionDenied(String,StatusException)")
    pruefeGrpc(PgChangeFeedGrpcNotFoundException("m", rpc), Status.Code.NOT_FOUND, "m", rpc, "GrpcNotFound(String,StatusException)")
    pruefeGrpc(PgChangeFeedGrpcInternalException("m", rpc), Status.Code.INTERNAL, "m", rpc, "GrpcInternal(String,StatusException)")
    pruefeGrpc(PgChangeFeedGrpcUnexpectedStatusException(Status.Code.UNAVAILABLE, "m", rpc), Status.Code.UNAVAILABLE, "m", rpc, "GrpcUnexpectedStatus(Code,String,StatusException)")
    pruefe(grpcGruppe(PgChangeFeedGrpcNotFoundException("x", rpc)) == "notfound", "when über die Basis (gRPC)")

    val gefangen = try {
        throw PgChangeFeedNotFoundException(404, "wurf")
    } catch (fehler: PgChangeFeedException) {
        fehler is PgChangeFeedNotFoundException
    }
    pruefe(gefangen, "catch PgChangeFeedException")
}

fun main(args: Array<String>) {
    val schritt = args.firstOrNull() ?: "?"
    try {
        lauf()
        println("KOMPAT kotlin $schritt: $aufrufe Aufrufe ok")
    } catch (ausnahme: Throwable) {
        println("KOMPAT kotlin $schritt: Ausnahme ${ausnahme.javaClass.simpleName}: ${ausnahme.message}")
        exitProcess(1)
    }
}
