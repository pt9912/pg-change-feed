package cdcexamples.http

import java.net.http.HttpClient
import java.time.Duration
import kotlin.test.Test
import kotlin.test.assertFailsWith

/**
 * Prüft den Fehlerpfad ohne erreichbaren Host — netzlos: die Verbindung zu
 * einem Loopback-Port ohne Listener scheitert sofort (Connection refused),
 * ohne echtes Netz zu brauchen. Dieselbe Grenze wie beim `natsnotify`-Adapter
 * und beim C#-Pendant
 * (`examples/csharp/http-client/HttpClient.Tests/TablesClientTests.cs`).
 */
class TablesClientTest {
    @Test
    fun listTablesFailsOnUnreachableHost() {
        val httpClient = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(2))
            .build()
        val cfg = Config("127.0.0.1:1", "token", "quelle", "pub")

        assertFailsWith<Exception> { TablesClient.listTables(httpClient, cfg) }
    }
}
