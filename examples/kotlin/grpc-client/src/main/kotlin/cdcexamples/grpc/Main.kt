package cdcexamples.grpc

import cdc.stream.v1.ChangeStreamGrpcKt
import cdc.stream.v1.Changestream.StreamChangesRequest
import io.grpc.ManagedChannelBuilder
import io.grpc.Metadata
import io.grpc.StatusException
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.runBlocking
import kotlin.system.exitProcess

/**
 * Command grpc-client ist ein öffentliches Beispiel für den Live-Change-Stream
 * über gRPC (`LH-FA-SST-008`, `ADR-0060`, `ADR-0090`): es öffnet den
 * Server-Streaming-RPC `ChangeStream/StreamChanges` real gegen den laufenden
 * Feed-Container und gibt jede empfangene Nachricht aus. Startform ist ein
 * Container-Aufruf, kein Host-Aufruf (`ADR-0087` Festlegung 3); Adresse und
 * Token kommen aus `CDC_GRPC_ADDR`/`CDC_API_TOKEN_READER` und lassen sich per
 * Flag übersteuern (`--addr`, `--token`).
 *
 * Form-Vorbild: `examples/grpc-client` (Go), `examples/csharp/grpc-client`
 * (C#, `slice-102`) — dieselbe Form (benannter Zusatzkontext, Generator-Stufe)
 * auf die Kotlin-Werkzeugkette übertragen (`slice-103`). Der Stub entsteht
 * **im Bau** aus der über den benannten Zusatzkontext gelesenen `.proto` — er
 * liegt nicht im committeten Baum (`ADR-0090` Festlegung 3). Dieses Programm
 * trägt keine Zustandsmaschine: der Stream kennt kein Replay (`ADR-0060`),
 * verpasste Changes holt der bestehende Lesezugriffsweg nach.
 *
 * Anders als der C#-Client (`Grpc.Net.Client`, ein `HttpClient`-basierter
 * Kanal) braucht der Kotlin-Client einen `ManagedChannel`
 * (`io.grpc:grpc-netty-shaded`) und den generierten Coroutine-Stub
 * (`io.grpc:grpc-kotlin-stub`) — beides Laufzeit-Unterschiede der
 * Werkzeugkette, kein Unterschied in der übertragenen Bau-Form (benannter
 * Zusatzkontext, Generator-Stufe): der `main`-Prozess läuft in
 * `runBlocking`, weil der generierte Server-Streaming-Aufruf einen kalten
 * `Flow<Change>` liefert, den nur ein Coroutine-Scope einsammeln kann.
 */
private const val AUTHORIZATION_METADATA_KEY = "authorization"
private const val BEARER_PREFIX = "Bearer "
private val AUTHORIZATION_METADATA_ENTRY: Metadata.Key<String> =
    Metadata.Key.of(AUTHORIZATION_METADATA_KEY, Metadata.ASCII_STRING_MARSHALLER)

fun main(args: Array<String>): Unit = runBlocking {
    val cfg = try {
        Cli.parse(args) { name -> System.getenv(name) }
    } catch (ex: IllegalArgumentException) {
        System.err.println(ex.message)
        exitProcess(2)
    }

    if (cfg.addr.isEmpty()) {
        System.err.println(
            "grpc-client: keine gRPC-Adresse gesetzt — CDC_GRPC_ADDR (oder --addr) ist noetig, um den Stream zu oeffnen",
        )
        exitProcess(2)
    }
    if (cfg.token.isEmpty()) {
        System.err.println(
            "grpc-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder --token) ist noetig, um ueber die reader-Rechtsklasse zu lesen",
        )
        exitProcess(2)
    }

    // Der Feed-Container spricht Klartext-gRPC (kein TLS, dieselbe Wahl wie
    // der Go-Client mit `insecure.NewCredentials()` und der C#-Client mit dem
    // h2c-Switch) — `usePlaintext()` erlaubt den Kanal ohne
    // Transportverschlüsselung.
    val channel = ManagedChannelBuilder.forTarget(cfg.addr).usePlaintext().build()
    try {
        val stub = ChangeStreamGrpcKt.ChangeStreamCoroutineStub(channel)
        val headers = Metadata().apply { put(AUTHORIZATION_METADATA_ENTRY, BEARER_PREFIX + cfg.token) }

        try {
            stub.streamChanges(StreamChangesRequest.getDefaultInstance(), headers)
                .collect { change -> println(Format.formatChange(change)) }
        } catch (ex: StatusException) {
            System.err.println("grpc-client: Stream endete: ${ex.status}")
            exitProcess(1)
        }

        System.err.println("grpc-client: Stream endete")
        exitProcess(1)
    } finally {
        channel.shutdownNow()
    }
}
