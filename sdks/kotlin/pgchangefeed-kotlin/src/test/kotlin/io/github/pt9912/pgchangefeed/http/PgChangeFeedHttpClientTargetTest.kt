package io.github.pt9912.pgchangefeed.http

import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * The optional `target` parameter of `readChanges`: it appears in the query
 * only when set, without it the query is identical to the call before the
 * parameter existed.
 */
class PgChangeFeedHttpClientTargetTest {

    private fun queryOf(call: (PgChangeFeedHttpClient) -> Unit): String {
        val (client, transport) = TestClientFactory.create { request ->
            FakeHttpTransport.jsonResponse(200, """{"changes":[]}""")(request)
        }
        call(client)
        return transport.lastRequest!!.url.substringAfter('?')
    }

    @Test
    fun `readChanges without target carries no target`() {
        assertEquals("source=src", queryOf { it.readChanges(source = "src") })
    }

    @Test
    fun `readChanges sends target only when set and escapes it`() {
        assertEquals("source=src", queryOf { it.readChanges(source = "src", target = null) })
        assertEquals("source=src&target=", queryOf { it.readChanges(source = "src", target = "") })
        assertEquals("source=src&target=eu", queryOf { it.readChanges(source = "src", target = "eu") })
        assertEquals("source=src&target=a%26b%3Dc", queryOf { it.readChanges(source = "src", target = "a&b=c") })
    }

    @Test
    fun `readChanges with target schema table and range is a conjunction on the wire`() {
        val query = queryOf {
            it.readChanges(
                source = "src", schema = "public", table = "orders", from = 1L, to = 10L, limit = 5, target = "eu",
            )
        }

        assertEquals("source=src&schema=public&table=orders&target=eu&from=1&to=10&limit=5", query)
    }
}
