package io.github.pt9912.pgchangefeed.grpc

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.runBlocking
import java.net.URI
import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * The optional `schema`/`table` filter of [PgChangeFeedGrpcClient.streamChanges]:
 * the parameterless call still sends the original, filter-less request
 * (regression), and each filter combination reaches the transport unchanged.
 */
class PgChangeFeedGrpcClientFilterTest {
    private fun options() = PgChangeFeedClientOptions(URI("http://localhost:50051"), "reader-token")

    // Mutation that turns this test red (checked for real): `streamChanges()`
    // changed to always set `schema`/`table` to a non-empty placeholder —
    // this test's assertion on both fields being empty then fails, instead of
    // staying silently green.
    @Test
    fun `streamChanges without arguments sends an empty request`() = runBlocking {
        val transport = FakeGrpcStreamTransport.withMessages()
        val client = PgChangeFeedGrpcClient(transport, options())

        client.streamChanges().toList()

        assertEquals("", transport.lastRequest?.schema)
        assertEquals("", transport.lastRequest?.table)
    }

    @Test
    fun `streamChanges with schema and table sends both in the request`() = runBlocking {
        val transport = FakeGrpcStreamTransport.withMessages()
        val client = PgChangeFeedGrpcClient(transport, options())

        client.streamChanges(schema = "public", table = "orders").toList()

        assertEquals("public", transport.lastRequest?.schema)
        assertEquals("orders", transport.lastRequest?.table)
    }

    @Test
    fun `streamChanges with schema only leaves table empty`() = runBlocking {
        val transport = FakeGrpcStreamTransport.withMessages()
        val client = PgChangeFeedGrpcClient(transport, options())

        client.streamChanges(schema = "public").toList()

        assertEquals("public", transport.lastRequest?.schema)
        assertEquals("", transport.lastRequest?.table)
    }
}
