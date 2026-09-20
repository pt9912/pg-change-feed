package io.github.pt9912.pgchangefeed.http

import io.github.pt9912.pgchangefeed.http.model.RegisterConsumerRequest
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertIs

/**
 * The SPEC-018 error-response mapping — the auth boundary (`401`
 * missing/unknown token, `403` a `reader` token against an `admin`
 * endpoint) plus the remaining documented statuses (`400`/`500`) and the
 * defensive fallback for anything outside that closed set. Same structure
 * as the C# sibling's `PgChangeFeedHttpClientAuthBoundaryTests`.
 */
class PgChangeFeedHttpClientAuthBoundaryTest {

    @Test
    fun `missing or unknown token throws Unauthorized`() {
        val (client, transport) = TestClientFactory.create(apiToken = "unknown-token") { request ->
            FakeHttpTransport.jsonResponse(401, """{"error":"missing or unknown bearer token"}""")(request)
        }

        val ex = assertFailsWith<PgChangeFeedUnauthorizedException> {
            client.listTables("src", "pub")
        }

        assertEquals(401, ex.statusCode)
        assertEquals("missing or unknown bearer token", ex.message)
        assertEquals("Bearer unknown-token", transport.lastRequest?.headers?.get("Authorization"))
    }

    @Test
    fun `reader token against an admin endpoint throws Forbidden`() {
        val (client, transport) = TestClientFactory.create(apiToken = "reader-token") { request ->
            FakeHttpTransport.jsonResponse(403, """{"error":"reader token cannot reach an admin endpoint"}""")(request)
        }

        val ex = assertFailsWith<PgChangeFeedForbiddenException> {
            client.registerConsumer(RegisterConsumerRequest("c-1", "n-1"))
        }

        assertEquals(403, ex.statusCode)
        assertEquals("reader token cannot reach an admin endpoint", ex.message)
        assertEquals("Bearer reader-token", transport.lastRequest?.headers?.get("Authorization"))
    }

    @Test
    fun `invalid request body throws BadRequest`() {
        val (client, _) = TestClientFactory.create { request ->
            FakeHttpTransport.jsonResponse(400, """{"error":"consumer_id must not be empty"}""")(request)
        }

        val ex = assertFailsWith<PgChangeFeedBadRequestException> {
            client.registerConsumer(RegisterConsumerRequest("", "n-1"))
        }

        assertEquals(400, ex.statusCode)
        assertEquals("consumer_id must not be empty", ex.message)
    }

    @Test
    fun `a physically missing table throws NotFound`() {
        val (client, _) = TestClientFactory.create { request ->
            FakeHttpTransport.jsonResponse(404, """{"error":"table does not exist at the source"}""")(request)
        }

        val ex = assertFailsWith<PgChangeFeedNotFoundException> {
            client.getStatus("src", "public", "missing", "pub")
        }

        assertEquals(404, ex.statusCode)
    }

    @Test
    fun `unexpected internal error throws ServerError`() {
        val (client, _) = TestClientFactory.create { request ->
            FakeHttpTransport.jsonResponse(500, """{"error":"unexpected internal error"}""")(request)
        }

        val ex = assertFailsWith<PgChangeFeedServerErrorException> {
            client.listTables("src", "pub")
        }

        assertEquals(500, ex.statusCode)
        assertEquals("unexpected internal error", ex.message)
    }

    @Test
    fun `a status code outside the documented set throws UnexpectedStatus`() {
        val (client, _) = TestClientFactory.create { request ->
            FakeHttpTransport.jsonResponse(409, """{"error":"unexpected"}""")(request)
        }

        val ex = assertFailsWith<PgChangeFeedUnexpectedStatusException> {
            client.listTables("src", "pub")
        }

        assertEquals(409, ex.statusCode)
    }

    @Test
    fun `a non-JSON error body falls back to the raw body`() {
        val (client, _) = TestClientFactory.create { request ->
            TransportResponse(500, "plain text failure")
        }

        val ex = assertFailsWith<PgChangeFeedServerErrorException> {
            client.listTables("src", "pub")
        }

        assertEquals("plain text failure", ex.message)
    }

    @Test
    fun `a non-JSON success body throws MalformedResponse`() {
        val (client, _) = TestClientFactory.create { request ->
            TransportResponse(200, "not json at all")
        }

        val ex = assertFailsWith<PgChangeFeedMalformedResponseException> {
            client.listTables("src", "pub")
        }

        assertEquals(200, ex.statusCode)
        assertIs<PgChangeFeedException>(ex)
    }

    @Test
    fun `an admin token reaches an admin endpoint that a reader token cannot`() {
        val (client, transport) = TestClientFactory.create(apiToken = "admin-token") { request ->
            FakeHttpTransport.jsonResponse(200, """{"consumer_id":"c-1","removed":true}""")(request)
        }

        val response = client.removeConsumer("c-1")

        assertEquals(true, response.removed)
        assertEquals("Bearer admin-token", transport.lastRequest?.headers?.get("Authorization"))
    }
}
