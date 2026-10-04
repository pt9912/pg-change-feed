package cdcexamples.grpc

import cdc.administration.v1.AdministrationGrpcKt
import cdc.stream.v1.ChangeStreamGrpcKt
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import io.grpc.StatusException
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.runBlocking
import kotlin.system.exitProcess

/**
 * Command grpc-client ist ein öffentliches Beispiel für die vollständige
 * gRPC-Fläche (`LH-FA-SST-006`): das `--verb`-Flag ruft eine von zwölf
 * dokumentierten Fähigkeiten real gegen den laufenden Feed-Container auf —
 * der Live-Change-Stream (Default-Verb `stream`, unverändert die
 * ursprüngliche Aufrufform) und die elf unären RPCs des
 * `Administration`-Diensts. Startform ist ein Container-Aufruf, kein
 * Host-Aufruf; die Zugriffs-Abschnitte des Benutzerhandbuchs sind „Zugriff
 * über den gRPC-Change-Stream" und „Zugriff über die gRPC-Verwaltungs-API"
 * (`docs/user/benutzerhandbuch.md`).
 *
 * Form-Vorbild: `examples/grpc-client` (Go), `examples/csharp/grpc-client`
 * (C#) — dieselbe Form (benannter Zusatzkontext, Generator-Stufe, zwölf
 * Verben) auf die Kotlin-Werkzeugkette übertragen. Beide Stubs entstehen im
 * Bau aus den über den benannten Zusatzkontext gelesenen `.proto`-Dateien —
 * sie liegen nicht im committeten Baum. Dieses Programm trägt keine
 * Zustandsmaschine: die elf RPCs stellen je eine Anfrage und enden, der
 * Stream bleibt offen und kennt kein Replay — verpasste Changes holt der
 * bestehende Lesezugriffsweg (`ReadChanges`, `GET /changes`) nach.
 *
 * Anders als der C#-Client (`Grpc.Net.Client`, ein `HttpClient`-basierter
 * Kanal) braucht der Kotlin-Client einen `ManagedChannel`
 * (`io.grpc:grpc-netty-shaded`) und die generierten Coroutine-Stubs
 * (`io.grpc:grpc-kotlin-stub`) — beides Laufzeit-Unterschiede der
 * Werkzeugkette, kein Unterschied in der übertragenen Bau-Form: der
 * `main`-Prozess läuft in `runBlocking`, weil sowohl der
 * Server-Streaming-Aufruf (`Flow<Change>`) als auch jeder unäre
 * Administration-Aufruf (`suspend fun`) nur aus einem Coroutine-Scope
 * heraus aufrufbar sind.
 */
fun main(args: Array<String>): Unit = runBlocking {
    val cfg = try {
        Cli.parse(args) { name -> System.getenv(name) }
    } catch (ex: IllegalArgumentException) {
        System.err.println(ex.message)
        exitProcess(2)
    }

    val validationError = Validator.validate(cfg)
    if (validationError != null) {
        System.err.println("grpc-client: $validationError")
        exitProcess(2)
    }

    // Der Feed-Container spricht TLS, wenn ein Zertifikatspaar konfiguriert
    // ist; dieses Beispiel verbindet im Klartext (dieselbe Wahl wie der
    // Go-Client mit `insecure.NewCredentials()` und der C#-Client mit dem
    // h2c-Switch) — `usePlaintext()` erlaubt den Kanal ohne
    // Transportverschlüsselung.
    val channel = ManagedChannelBuilder.forTarget(cfg.addr).usePlaintext().build()
    try {
        if (cfg.verb == "stream") {
            runStream(channel, cfg)
            return@runBlocking
        }

        val client = AdministrationGrpcKt.AdministrationCoroutineStub(channel)
        try {
            println(Dispatcher.dispatchAdmin(client, cfg))
        } catch (ex: StatusException) {
            System.err.println("grpc-client: --verb=${cfg.verb} fehlgeschlagen: ${ex.status}")
            exitProcess(1)
        }
    } finally {
        channel.shutdownNow()
    }
}

/**
 * runStream öffnet den Server-Streaming-RPC `ChangeStream/StreamChanges`
 * real gegen den laufenden Feed-Container und gibt jede empfangene
 * Nachricht aus, bis die Verbindung endet; die Anfrage baut
 * [StreamRequest.build].
 */
private suspend fun runStream(channel: ManagedChannel, cfg: Config) {
    val stub = ChangeStreamGrpcKt.ChangeStreamCoroutineStub(channel)
    val headers = CallMetadata.headers(cfg.token)
    val request = StreamRequest.build(cfg)

    try {
        stub.streamChanges(request, headers).collect { change -> println(Format.formatChange(change)) }
    } catch (ex: StatusException) {
        System.err.println("grpc-client: Stream endete: ${ex.status}")
        exitProcess(1)
    }

    System.err.println("grpc-client: Stream endete")
    exitProcess(1)
}
