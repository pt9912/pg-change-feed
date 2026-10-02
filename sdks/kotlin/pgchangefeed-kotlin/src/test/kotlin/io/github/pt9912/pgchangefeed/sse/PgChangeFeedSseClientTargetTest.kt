package io.github.pt9912.pgchangefeed.sse

import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * The optional `target` of the SSE stream: it appears as the query parameter
 * `target` only when set; without it the request URL is the bare stream path.
 */
class PgChangeFeedSseClientTargetTest {

    private val base = "http://example.invalid:8080/changes/stream"

    private fun urlOf(target: String?, schema: String? = null, table: String? = null): String {
        val (client, transport) = TestClientFactory.create { _ -> FakeSseTransport.lineResponse(200, "") }
        client.streamChanges(target = target, schema = schema, table = table).toList()
        return transport.lastRequest!!.url
    }

    @Test
    fun `streamChanges without filters requests the bare stream path`() {
        assertEquals(base, urlOf(null, null, null))
    }

    @Test
    fun `streamChanges sends schema as an escaped query parameter`() {
        assertEquals("$base?schema=", urlOf(null, schema = ""))
        assertEquals("$base?schema=eu", urlOf(null, schema = "eu"))
        assertEquals("$base?schema=a%26b%3Dc", urlOf(null, schema = "a&b=c"))
        assertEquals("$base?schema=a%2Bb", urlOf(null, schema = "a+b"))
        assertEquals("$base?schema=100%25", urlOf(null, schema = "100%"))
        assertEquals("$base?schema=%C3%BC", urlOf(null, schema = "ü"))
        // The Python client sends a space as "+"; the server decodes both alike.
        assertEquals("$base?schema=a%20b", urlOf(null, schema = "a b"))
    }

    @Test
    fun `streamChanges sends table as an escaped query parameter`() {
        assertEquals("$base?table=", urlOf(null, table = ""))
        assertEquals("$base?table=eu", urlOf(null, table = "eu"))
        assertEquals("$base?table=a%26b%3Dc", urlOf(null, table = "a&b=c"))
        assertEquals("$base?table=a%2Bb", urlOf(null, table = "a+b"))
        assertEquals("$base?table=100%25", urlOf(null, table = "100%"))
        assertEquals("$base?table=%C3%BC", urlOf(null, table = "ü"))
        // The Python client sends a space as "+"; the server decodes both alike.
        assertEquals("$base?table=a%20b", urlOf(null, table = "a b"))
    }

    @Test
    fun `streamChanges sends schema and table together`() {
        assertEquals("$base?schema=public&table=orders", urlOf(null, schema = "public", table = "orders"))
    }

    @Test
    fun `streamChanges combines table and target as a conjunction`() {
        assertEquals("$base?table=orders&target=eu", urlOf("eu", table = "orders"))
    }

    @Test
    fun `streamChanges combines schema table and target as a conjunction`() {
        assertEquals(
            "$base?schema=public&table=orders&target=eu",
            urlOf("eu", schema = "public", table = "orders"),
        )
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
        assertEquals("$base?target=a%2Bb", urlOf("a+b"))
        assertEquals("$base?target=100%25", urlOf("100%"))
        assertEquals("$base?target=%C3%BC", urlOf("ü"))
        // The Python client sends a space as "+"; the server decodes both alike.
        assertEquals("$base?target=a%20b", urlOf("a b"))
    }
}
