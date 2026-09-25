package io.github.pt9912.pgchangefeed.sse

import io.github.pt9912.pgchangefeed.http.PgChangeFeedBadRequestException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedMalformedResponseException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedServerErrorException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnauthorizedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnexpectedStatusException
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertIs

/**
 * The error-response mapping when opening `GET /changes/stream` fails before
 * any frame is delivered — the same set of typed exceptions as
 * `io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient` (the uniform
 * `{"error": "<text>"}` body), reused rather than a second hierarchy. Same
 * structure as
 * `io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClientAuthBoundaryTest`.
 */
class PgChangeFeedSseClientAuthBoundaryTest {

    // Mutation that turns this test red (checked for real):
    // `PgChangeFeedSseClient.buildException`s
    // `401`-Zweig auf `PgChangeFeedForbiddenException(...)` statt
    // `PgChangeFeedUnauthorizedException(...)` geändert — dieser Test
    // schlägt dann fehl (falscher Exception-Typ), statt still grün zu
    // bleiben.
    @Test
    fun `missing or unknown token throws Unauthorized`() {
        val (client, transport) = TestClientFactory.create(apiToken = "unknown-token") { _ ->
            FakeSseTransport.lineResponse(401, """{"error":"missing or unknown bearer token"}""")
        }

        val ex = assertFailsWith<PgChangeFeedUnauthorizedException> {
            client.streamChanges().toList()
        }

        assertEquals(401, ex.statusCode)
        assertEquals("missing or unknown bearer token", ex.message)
        assertEquals("Bearer unknown-token", transport.lastRequest?.headers?.get("Authorization"))
    }

    @Test
    fun `stream not available throws UnexpectedStatus`() {
        // `503` when the HTTP API is on but the change stream is not
        // available — outside the 400/401/403/404/500 set.
        val (client, _) = TestClientFactory.create { _ ->
            FakeSseTransport.lineResponse(503, """{"error":"change stream not available"}""")
        }

        val ex = assertFailsWith<PgChangeFeedUnexpectedStatusException> {
            client.streamChanges().toList()
        }

        assertEquals(503, ex.statusCode)
    }

    @Test
    fun `an invalid request throws BadRequest`() {
        val (client, _) = TestClientFactory.create { _ ->
            FakeSseTransport.lineResponse(400, """{"error":"unexpected request"}""")
        }

        val ex = assertFailsWith<PgChangeFeedBadRequestException> {
            client.streamChanges().toList()
        }

        assertEquals(400, ex.statusCode)
    }

    @Test
    fun `unexpected internal error throws ServerError`() {
        val (client, _) = TestClientFactory.create { _ ->
            FakeSseTransport.lineResponse(500, """{"error":"unexpected internal error"}""")
        }

        val ex = assertFailsWith<PgChangeFeedServerErrorException> {
            client.streamChanges().toList()
        }

        assertEquals(500, ex.statusCode)
        assertEquals("unexpected internal error", ex.message)
    }

    @Test
    fun `a non-JSON error body falls back to the raw body`() {
        val (client, _) = TestClientFactory.create { _ ->
            FakeSseTransport.lineResponse(500, "plain text failure")
        }

        val ex = assertFailsWith<PgChangeFeedServerErrorException> {
            client.streamChanges().toList()
        }

        assertEquals("plain text failure", ex.message)
    }

    @Test
    fun `a non-JSON frame data payload throws MalformedResponse`() {
        val (client, _) = TestClientFactory.create { _ ->
            FakeSseTransport.lineResponse(200, "event: change\ndata: not json at all\n\n")
        }

        val ex = assertFailsWith<PgChangeFeedMalformedResponseException> {
            client.streamChanges().toList()
        }

        assertEquals(200, ex.statusCode)
        assertIs<PgChangeFeedException>(ex)
    }

    @Test
    fun `an admin token reaches the stream that a bare open call never sends eagerly`() {
        val (client, transport) = TestClientFactory.create(apiToken = "admin-token") { _ ->
            FakeSseTransport.lineResponse(
                200,
                """event: change
data: {"change_id":"c-1","transaction_id":"tx-1","source_table_id":"t-1","sequence":1,"operation":"INSERT","old_image":null,"new_image":{"id":1},"schema_version":"t-1-v1","schema":"public","table":"orders"}

""",
            )
        }

        // Die Sequence ist kalt — bis hier hin wurde noch keine Anfrage
        // gebaut (SseTransportResponse existiert erst nach der ersten
        // Iteration).
        assertEquals(null, transport.lastRequest)

        val changes = client.streamChanges().toList()

        assertEquals(1, changes.size)
        assertEquals("Bearer admin-token", transport.lastRequest?.headers?.get("Authorization"))
    }
}
