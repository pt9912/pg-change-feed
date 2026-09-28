package cdcexamples.http

import java.net.http.HttpClient
import kotlin.test.Test
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

/**
 * Prüft den defensiven Schutz von [Dispatcher.dispatch] gegen ein
 * unbekanntes Verb — für den Fall eines Aufrufs ohne vorherige
 * [Validator.validate]. Form-Vorbild: `examples/http-client/main_test.go`,
 * `TestDispatchRejectsUnknownVerb`.
 */
class DispatcherTest {
    @Test
    fun dispatchRejectsUnknownVerbBeforeAnyNetworkCall() {
        val httpClient = HttpClient.newHttpClient()
        val cfg = Config(
            addr = "feed:8080", token = "", adminToken = "", verb = "unbekannt",
            source = "", publication = "", consumerId = "", name = "", offset = 0,
            schema = "", table = "", tableId = "", schemaVersionId = "", version = 1,
            from = "", to = "", limit = "", minAgeNanos = 0,
        )

        val ex = assertFailsWith<IllegalStateException> { Dispatcher.dispatch(httpClient, cfg) }
        assertTrue(ex.message!!.contains("unbekanntes --verb"))
    }
}
