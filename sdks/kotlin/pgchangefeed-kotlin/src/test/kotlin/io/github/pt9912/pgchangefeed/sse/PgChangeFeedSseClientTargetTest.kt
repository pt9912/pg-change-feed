package io.github.pt9912.pgchangefeed.sse

import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * The optional `target` of the SSE stream: it appears as the query parameter
 * `target` only when set; without it the request URL is the bare stream path.
 */
class PgChangeFeedSseClientTargetTest {

    private fun urlOf(target: String?): String {
        val (client, transport) = TestClientFactory.create { _ -> FakeSseTransport.lineResponse(200, "") }
        if (target == null) client.streamChanges().toList() else client.streamChanges(target = target).toList()
        return transport.lastRequest!!.url
    }

    @Test
    fun `streamChanges without target requests the bare stream path`() {
        assertEquals("http://example.invalid:8080/changes/stream", urlOf(null))
    }

    @Test
    fun `streamChanges sends target as an escaped query parameter`() {
        assertEquals("http://example.invalid:8080/changes/stream?target=", urlOf(""))
        assertEquals("http://example.invalid:8080/changes/stream?target=eu", urlOf("eu"))
        assertEquals("http://example.invalid:8080/changes/stream?target=a%26b%3Dc", urlOf("a&b=c"))
    }
}
