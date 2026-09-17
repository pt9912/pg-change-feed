package cdcexamples.nats

import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse
import java.time.Duration

/**
 * ChangesClient holt die Änderungen eines Quelle/Schema/Tabelle-Filters mit
 * dem `reader`-Token ab (`LH-FA-SST-006`, zweiseitiger Ablauf `ADR-0079`) —
 * Form-Vorbild: `examples/nats-client/main.go` (`fetchChanges`),
 * `examples/kotlin/http-client/TablesClient.kt`. Ein Nicht-200-Status ist
 * ein sichtbarer Fehler, kein stiller Leerwert. Der [HttpClient] ist
 * injiziert, damit der Fehlerpfad ohne erreichbaren Host netzlos testbar ist.
 */
object ChangesClient {
    private val requestTimeout: Duration = Duration.ofSeconds(10)

    fun fetch(httpClient: HttpClient, cfg: Config): String {
        val url = ChangesUrl.build(cfg.addr, cfg.source, cfg.schema, cfg.table)
        val request = HttpRequest.newBuilder(URI.create(url))
            .header("Authorization", "Bearer ${cfg.token}")
            .timeout(requestTimeout)
            .GET()
            .build()

        val response = httpClient.send(request, HttpResponse.BodyHandlers.ofString())
        if (response.statusCode() != 200) {
            throw IllegalStateException("HTTP-Status ${response.statusCode()}: ${response.body()}")
        }
        return response.body()
    }
}
