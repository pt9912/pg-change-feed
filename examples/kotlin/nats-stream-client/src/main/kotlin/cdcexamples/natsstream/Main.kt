package cdcexamples.natsstream

import com.google.gson.Gson
import io.nats.client.Message
import io.nats.client.Nats
import io.nats.client.Options
import java.time.Duration
import kotlin.system.exitProcess

/**
 * Command nats-stream-client ist ein öffentliches Beispiel für den
 * dritten, vollinhaltstragenden NATS-Zustellweg (`LH-FA-SST-008`,
 * `ADR-0100`, `SPEC-024`): es verbindet mit einem gültigen
 * `CDC_NATS_STREAM_TOKEN`, abonniert den Vollinhalts-Namensraum
 * `cdc.stream.>` real gegen den laufenden Feed-Container und gibt jede
 * empfangene Change aus. Startform ist ein Container-Aufruf, kein
 * Host-Aufruf (`ADR-0087` Festlegung 3); NATS-URL und Token kommen aus
 * `CDC_NATS_URL`/`CDC_NATS_STREAM_TOKEN` und lassen sich per Flag
 * übersteuern (`--nats-url`, `--token`).
 *
 * Form-Vorbild: `examples/nats-stream-client` (Go),
 * `examples/csharp/nats-stream-client` (C#). Dieses Programm trägt keine
 * Zustandsmaschine: der Stream kennt kein Replay (`ADR-0100`), verpasste
 * Changes holt der bestehende Lesezugriffsweg nach. Das Warten auf das
 * nächste Ereignis benutzt [io.nats.client.Subscription.nextMessage] mit
 * [Duration.ZERO] — jnats wartet damit real unbegrenzt (dieselbe Semantik
 * wie `examples/kotlin/nats-client`, dort ausführlich begründet).
 */
private const val subject = "cdc.stream.>"

fun main(args: Array<String>) {
    val cfg = try {
        Cli.parse(args) { name -> System.getenv(name) }
    } catch (ex: IllegalArgumentException) {
        System.err.println(ex.message)
        exitProcess(2)
    }

    if (cfg.natsUrl.isEmpty()) {
        System.err.println(
            "nats-stream-client: keine NATS-URL gesetzt — CDC_NATS_URL (oder --nats-url) ist noetig, um den Vollinhalts-Stream zu oeffnen",
        )
        exitProcess(2)
    }
    if (cfg.token.isEmpty()) {
        System.err.println(
            "nats-stream-client: kein Token gesetzt — CDC_NATS_STREAM_TOKEN (oder --token) ist noetig, um den dritten Zustellweg zu abonnieren",
        )
        exitProcess(2)
    }

    val options = Options.Builder().server(cfg.natsUrl).token(cfg.token).build()
    val nc = try {
        Nats.connect(options)
    } catch (ex: Exception) {
        System.err.println("nats-stream-client: Verbindung zu NATS (${cfg.natsUrl}) fehlgeschlagen: $ex")
        exitProcess(1)
    }

    val sub = try {
        nc.subscribe(subject)
    } catch (ex: Exception) {
        System.err.println("nats-stream-client: Abonnement auf $subject fehlgeschlagen: $ex")
        exitProcess(1)
    }

    println("nats-stream-client: lauscht auf $subject")

    val gson = Gson()
    try {
        while (true) {
            val msg: Message = sub.nextMessage(Duration.ZERO) ?: break
            val change = gson.fromJson(String(msg.data, Charsets.UTF_8), StreamMessage::class.java)
            println(Format.formatChange(change))
        }
    } catch (ex: Exception) {
        System.err.println("nats-stream-client: Empfang beendet: $ex")
        exitProcess(1)
    } finally {
        nc.close()
    }

    System.err.println("nats-stream-client: Empfang beendet")
    exitProcess(1)
}
