package kompat

import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcNotFoundException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedBadRequestException
import io.grpc.Status
import kotlin.system.exitProcess

// Gegenrichtung: ein gegen die Bibliothek mit Meldungscode übersetzter Aufrufer
// nutzt die Konstruktoren mit Meldungscode und liest die Eigenschaft. Unter der
// Bibliothek 0.5.0 gibt es diese Mitglieder nicht.
private var aufrufe = 0

private fun pruefe(ok: Boolean, was: String) {
    check(ok) { "Abweichung bei $was" }
    aufrufe++
}

private fun lauf() {
    val http = PgChangeFeedBadRequestException(400, "bad", "PCF-E0001")
    pruefe(http.messageCode == "PCF-E0001", "HTTP messageCode")

    val rpc = Status.NOT_FOUND.withDescription("rpc").asException()
    val grpc = PgChangeFeedGrpcNotFoundException("m", rpc, "PCF-E0002")
    pruefe(grpc.messageCode == "PCF-E0002", "gRPC messageCode")
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
