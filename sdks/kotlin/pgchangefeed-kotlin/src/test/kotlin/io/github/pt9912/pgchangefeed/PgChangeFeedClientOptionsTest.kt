package io.github.pt9912.pgchangefeed

import java.net.URI
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class PgChangeFeedClientOptionsTest {

    @Test
    fun `constructs with address and token`() {
        val options = PgChangeFeedClientOptions(URI("https://cdc.example.test:8080"), "test-token")

        assertEquals(URI("https://cdc.example.test:8080"), options.address)
        assertEquals("test-token", options.apiToken)
    }

    @Test
    fun `rejects a blank api token`() {
        assertFailsWith<IllegalArgumentException> {
            PgChangeFeedClientOptions(URI("https://cdc.example.test:8080"), "   ")
        }
    }

    @Test
    fun `rejects an empty api token`() {
        assertFailsWith<IllegalArgumentException> {
            PgChangeFeedClientOptions(URI("https://cdc.example.test:8080"), "")
        }
    }
}
