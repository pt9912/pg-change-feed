package cdcexamples.http

import java.net.http.HttpClient
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertTrue

/**
 * Prüft die vier Consumer-Verwaltungs-Fähigkeiten — netzlos, über
 * [FakeServer]. Form-Vorbild: `examples/http-client/consumer_test.go`.
 */
class ConsumerClientTest {
    private fun baseConfig(addr: String) = Config(
        addr = addr, token = "reader-token", adminToken = "admin-token", verb = "",
        source = "", publication = "", consumerId = "", name = "", offset = 0,
        schema = "", table = "", tableId = "", schemaVersionId = "", version = 1,
        from = "", to = "", limit = "", minAgeNanos = 0,
    )

    @Test
    fun consumerPositionUrlBuilderCarriesConsumerId() {
        val got = ConsumerPositionUrlBuilder.build("feed:8080", "consumer-1")
        assertEquals("http://feed:8080/consumers/position?consumer_id=consumer-1", got)
    }

    @Test
    fun consumerPositionUrlBuilderEscapesReservedCharacters() {
        val got = ConsumerPositionUrlBuilder.build("feed:8080", "consumer & 1")
        assertEquals("http://feed:8080/consumers/position?consumer_id=consumer%20%26%201", got)
    }

    @Test
    fun registerConsumerParsesSuccessResponse() {
        FakeServer.start().use { server ->
            server.respondWith(201, """{"consumer_id":"consumer-1","name":"Consumer Eins","already_registered":false}""")
            val httpClient = HttpClient.newHttpClient()
            val cfg = baseConfig(server.addr).copy(consumerId = "consumer-1", name = "Consumer Eins")

            val resp = ConsumerClient.registerConsumer(httpClient, cfg)

            assertEquals("consumer-1", resp.consumerId)
            assertFalse(resp.alreadyRegistered)
            assertEquals("POST", server.lastMethod)
            assertEquals("/consumers", server.lastPath)
            assertTrue(server.lastBody!!.contains("\"consumer_id\":\"consumer-1\""))
        }
    }

    @Test
    fun registerConsumerFailsOnNon2xx() {
        FakeServer.start().use { server ->
            server.respondWith(400, """{"error":"consumer_id und name sind Pflichtfelder"}""")
            val httpClient = HttpClient.newHttpClient()

            val ex = assertFailsWith<IllegalStateException> {
                ConsumerClient.registerConsumer(httpClient, baseConfig(server.addr))
            }
            assertTrue(ex.message!!.contains("400"))
            assertTrue(ex.message!!.contains("Pflichtfelder"))
        }
    }

    @Test
    fun acknowledgeConsumerParsesSuccessResponse() {
        FakeServer.start().use { server ->
            server.respondWith(200, """{"consumer_id":"consumer-1","source_id":"quelle-1","offset":42}""")
            val httpClient = HttpClient.newHttpClient()
            val cfg = baseConfig(server.addr).copy(consumerId = "consumer-1", source = "quelle-1", offset = 42)

            val resp = ConsumerClient.acknowledgeConsumer(httpClient, cfg)

            assertEquals(42L, resp.offset)
            assertEquals("/consumers/acknowledge", server.lastPath)
            assertTrue(server.lastBody!!.contains("\"offset\":42"))
        }
    }

    @Test
    fun acknowledgeConsumerFailsOnNon2xx() {
        FakeServer.start().use { server ->
            server.respondWith(400, """{"error":"consumer_id, source und offset sind Pflichtfelder"}""")
            val httpClient = HttpClient.newHttpClient()

            val ex = assertFailsWith<IllegalStateException> {
                ConsumerClient.acknowledgeConsumer(httpClient, baseConfig(server.addr))
            }
            assertTrue(ex.message!!.contains("400"))
            assertTrue(ex.message!!.contains("Pflichtfelder"))
        }
    }

    @Test
    fun consumerPositionParsesUnacknowledgedBoundary() {
        FakeServer.start().use { server ->
            server.respondWith(200, """{"consumer_id":"consumer-1","source_id":"","offset":0,"acknowledged":false}""")
            val httpClient = HttpClient.newHttpClient()
            val cfg = baseConfig(server.addr).copy(consumerId = "consumer-1")

            val resp = ConsumerClient.consumerPosition(httpClient, cfg)

            assertFalse(resp.acknowledged)
            assertEquals("GET", server.lastMethod)
            assertEquals("consumer_id=consumer-1", server.lastQuery)
        }
    }

    @Test
    fun consumerPositionFailsOnNon2xx() {
        FakeServer.start().use { server ->
            server.respondWith(400, """{"error":"consumer_id ist Pflichtfeld"}""")
            val httpClient = HttpClient.newHttpClient()

            val ex = assertFailsWith<IllegalStateException> {
                ConsumerClient.consumerPosition(httpClient, baseConfig(server.addr))
            }
            assertTrue(ex.message!!.contains("400"))
            assertTrue(ex.message!!.contains("Pflichtfeld"))
        }
    }

    @Test
    fun removeConsumerParsesIdempotentOutcome() {
        FakeServer.start().use { server ->
            server.respondWith(200, """{"consumer_id":"nie-registriert","removed":false}""")
            val httpClient = HttpClient.newHttpClient()
            val cfg = baseConfig(server.addr).copy(consumerId = "nie-registriert")

            val resp = ConsumerClient.removeConsumer(httpClient, cfg)

            assertFalse(resp.removed)
            assertEquals("/consumers/remove", server.lastPath)
        }
    }

    @Test
    fun removeConsumerFailsOnNon2xx() {
        FakeServer.start().use { server ->
            server.respondWith(400, """{"error":"consumer_id ist Pflichtfeld"}""")
            val httpClient = HttpClient.newHttpClient()

            val ex = assertFailsWith<IllegalStateException> {
                ConsumerClient.removeConsumer(httpClient, baseConfig(server.addr))
            }
            assertTrue(ex.message!!.contains("400"))
            assertTrue(ex.message!!.contains("Pflichtfeld"))
        }
    }
}
