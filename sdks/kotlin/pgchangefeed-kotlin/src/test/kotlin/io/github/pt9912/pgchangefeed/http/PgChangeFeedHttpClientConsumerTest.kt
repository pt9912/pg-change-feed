package io.github.pt9912.pgchangefeed.http

import io.github.pt9912.pgchangefeed.http.model.AcknowledgeConsumerRequest
import io.github.pt9912.pgchangefeed.http.model.RegisterConsumerRequest
import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * Happy-path coverage for the four consumer-management capabilities of
 * SPEC-018 (`RegisterConsumer`, `AcknowledgeConsumer`,
 * `GetConsumerPosition`, `RemoveConsumer`).
 */
class PgChangeFeedHttpClientConsumerTest {

    @Test
    fun `registerConsumer posts to consumers and parses the response`() {
        val (client, transport) = TestClientFactory.create(apiToken = "admin-token") { request ->
            assertEquals("POST", request.method)
            assertEquals("http://example.invalid:8080/consumers", request.url)
            assertEquals("""{"consumer_id":"c-1","name":"n-1"}""", request.body)
            FakeHttpTransport.jsonResponse(
                201,
                """{"consumer_id":"c-1","name":"n-1","already_registered":false}""",
            )(request)
        }

        val response = client.registerConsumer(RegisterConsumerRequest("c-1", "n-1"))

        assertEquals("c-1", response.consumerId)
        assertEquals("n-1", response.name)
        assertEquals(false, response.alreadyRegistered)
        assertEquals("Bearer admin-token", transport.lastRequest?.headers?.get("Authorization"))
        assertEquals("application/json; charset=utf-8", transport.lastRequest?.headers?.get("Content-Type"))
    }

    @Test
    fun `acknowledgeConsumer posts to consumers acknowledge and parses the response`() {
        val (client, transport) = TestClientFactory.create { request ->
            assertEquals("POST", request.method)
            assertEquals("http://example.invalid:8080/consumers/acknowledge", request.url)
            FakeHttpTransport.jsonResponse(
                200,
                """{"consumer_id":"c-1","source_id":"s-1","offset":42}""",
            )(request)
        }

        val response = client.acknowledgeConsumer(AcknowledgeConsumerRequest("c-1", "s-1", 42L))

        assertEquals("c-1", response.consumerId)
        assertEquals("s-1", response.sourceId)
        assertEquals(42L, response.offset)
        assertEquals(true, transport.lastRequest != null)
    }

    @Test
    fun `getConsumerPosition gets consumers position with the consumer_id query parameter`() {
        val (client, transport) = TestClientFactory.create { request ->
            assertEquals("GET", request.method)
            assertEquals("http://example.invalid:8080/consumers/position?consumer_id=c-1", request.url)
            FakeHttpTransport.jsonResponse(
                200,
                """{"consumer_id":"c-1","source_id":"s-1","offset":7,"acknowledged":true}""",
            )(request)
        }

        val response = client.getConsumerPosition("c-1")

        assertEquals("c-1", response.consumerId)
        assertEquals(7L, response.offset)
        assertEquals(true, response.acknowledged)
        assertEquals(null, transport.lastRequest?.headers?.get("Content-Type"))
    }

    @Test
    fun `removeConsumer posts to consumers remove and parses the response`() {
        val (client, _) = TestClientFactory.create { request ->
            assertEquals("POST", request.method)
            assertEquals("http://example.invalid:8080/consumers/remove", request.url)
            assertEquals("""{"consumer_id":"c-1"}""", request.body)
            FakeHttpTransport.jsonResponse(200, """{"consumer_id":"c-1","removed":true}""")(request)
        }

        val response = client.removeConsumer("c-1")

        assertEquals("c-1", response.consumerId)
        assertEquals(true, response.removed)
    }

    @Test
    fun `removeConsumer reports removed false for a never-registered consumer without a 404`() {
        val (client, _) = TestClientFactory.create { request ->
            FakeHttpTransport.jsonResponse(200, """{"consumer_id":"never-seen","removed":false}""")(request)
        }

        val response = client.removeConsumer("never-seen")

        assertEquals(false, response.removed)
    }
}
