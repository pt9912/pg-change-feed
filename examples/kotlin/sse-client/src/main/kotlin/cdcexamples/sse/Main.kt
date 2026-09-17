package cdcexamples.sse

import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse
import kotlin.system.exitProcess

/**
 * Command sse-client ist ein öffentliches Beispiel für den Live-Change-Stream
 * über Server-Sent-Events (`LH-FA-SST-008`, `ADR-0061`, `ADR-0090`): es
 * öffnet den Endpunkt `GET /changes/stream` real gegen den laufenden
 * Feed-Container und gibt jedes Event aus. Startform ist ein
 * Container-Aufruf, kein Host-Aufruf (`ADR-0087` Festlegung 3); Adresse und
 * Token kommen aus `CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER` und lassen sich per
 * Flag übersteuern (`--addr`, `--token`).
 *
 * Form-Vorbild: `examples/sse-client` (Go), `examples/csharp/sse-client`
 * (C#). Dieses Programm trägt keine Zustandsmaschine: der Stream kennt kein
 * Replay (`ADR-0061`), verpasste Changes holt der bestehende
 * Lesezugriffsweg nach. Der Stream wird über
 * [HttpResponse.BodyHandlers.ofLines] gelesen — derselbe Streaming-Body-
 * Handler, der seit JDK 11 zur Standardbibliothek gehört (**keine neue
 * Abhängigkeit**, `ADR-0090` Festlegung 1).
 */
fun main(args: Array<String>) {
    val cfg = try {
        Cli.parse(args) { name -> System.getenv(name) }
    } catch (ex: IllegalArgumentException) {
        System.err.println(ex.message)
        exitProcess(2)
    }

    if (cfg.addr.isEmpty()) {
        System.err.println(
            "sse-client: keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder --addr) ist noetig, um den Stream zu oeffnen",
        )
        exitProcess(2)
    }
    if (cfg.token.isEmpty()) {
        System.err.println(
            "sse-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder --token) ist noetig, um ueber die reader-Rechtsklasse zu lesen",
        )
        exitProcess(2)
    }

    val streamUrl = SseStream.streamUrl(cfg.addr)
    val httpClient = HttpClient.newBuilder().build()
    val request = HttpRequest.newBuilder(URI.create(streamUrl))
        .header("Authorization", "Bearer ${cfg.token}")
        .GET()
        .build()

    // Der Stream bleibt offen, bis die Verbindung endet: der Aufruf läuft
    // ohne `HttpRequest.timeout()` (Muster: examples/sse-client/main.go).
    val response = try {
        httpClient.send(request, HttpResponse.BodyHandlers.ofLines())
    } catch (ex: Exception) {
        System.err.println("sse-client: Stream oeffnen ($streamUrl): $ex")
        exitProcess(1)
    }

    if (response.statusCode() != 200) {
        System.err.println(
            "sse-client: Stream-Oeffnung endete mit HTTP-Status ${response.statusCode()} (Erwartung 200)",
        )
        exitProcess(1)
    }

    val lines = response.body().iterator()
    val next: () -> String? = { if (lines.hasNext()) lines.next() else null }

    while (true) {
        val ev = try {
            SseStream.readEvent(next)
        } catch (ex: Exception) {
            // Das Ende der Quelle trägt zwei Ausgänge — ein Lesefehler oder
            // eine geschlossene Verbindung. Beide werden gemeldet; der
            // Prozess endet mit Nicht-Null.
            System.err.println("sse-client: Stream endete mit Fehler: $ex")
            exitProcess(1)
        }
        if (ev == null) {
            System.err.println("sse-client: Stream wurde geschlossen")
            exitProcess(1)
        }
        println("sse-client: ${ev.name} ${ev.data}")
    }
}
