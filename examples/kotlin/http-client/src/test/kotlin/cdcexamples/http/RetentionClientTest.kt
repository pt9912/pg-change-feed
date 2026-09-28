package cdcexamples.http

import java.net.http.HttpClient
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

/**
 * Prüft `POST /retention/run` — netzlos, über [FakeServer]. Form-Vorbild:
 * `examples/http-client/retention_test.go`.
 */
class RetentionClientTest {
    private fun baseConfig(addr: String) = Config(
        addr = addr, token = "", adminToken = "admin-token", verb = "",
        source = "", publication = "", consumerId = "", name = "", offset = 0,
        schema = "", table = "", tableId = "", schemaVersionId = "", version = 1,
        from = "", to = "", limit = "", minAgeNanos = 0,
    )

    @Test
    fun runRetentionParsesZeroDeletedAsValidOutcome() {
        FakeServer.start().use { server ->
            server.respondWith(200, """{"deleted":0}""")
            val httpClient = HttpClient.newHttpClient()
            val cfg = baseConfig(server.addr).copy(source = "quelle-1", minAgeNanos = 3_600_000_000_000)

            val resp = RetentionClient.runRetention(httpClient, cfg)

            assertEquals(0, resp.deleted)
            assertEquals("/retention/run", server.lastPath)
            assertTrue(server.lastBody!!.contains("\"min_age_nanos\":3600000000000"))
        }
    }

    @Test
    fun runRetentionFailsOnNon2xx() {
        FakeServer.start().use { server ->
            server.respondWith(400, """{"error":"min_age_nanos darf nicht negativ sein"}""")
            val httpClient = HttpClient.newHttpClient()
            val cfg = baseConfig(server.addr).copy(source = "quelle-1", minAgeNanos = -1)

            val ex = assertFailsWith<IllegalStateException> {
                RetentionClient.runRetention(httpClient, cfg)
            }
            assertTrue(ex.message!!.contains("400"))
        }
    }
}
