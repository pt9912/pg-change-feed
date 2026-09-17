package cdcexamples.http

import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse
import java.time.Duration

/**
 * TablesClient ruft die Tabellen-Auflistung mit dem `reader`-Token ab
 * (`LH-FA-SST-006`, `ADR-0057`) — Form-Vorbild:
 * `examples/http-client/main.go`, `listTables`. Ein Nicht-200-Status ist ein
 * sichtbarer Fehler, kein stiller Leerwert. Der [HttpClient] ist injiziert,
 * damit der Fehlerpfad ohne erreichbaren Host netzlos testbar ist.
 */
object TablesClient {
    private val requestTimeout: Duration = Duration.ofSeconds(10)

    fun listTables(httpClient: HttpClient, cfg: Config): String {
        val url = TablesUrlBuilder.build(cfg.addr, cfg.source, cfg.publication)
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
