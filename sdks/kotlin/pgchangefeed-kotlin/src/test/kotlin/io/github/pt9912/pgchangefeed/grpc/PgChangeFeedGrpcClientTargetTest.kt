package io.github.pt9912.pgchangefeed.grpc

import cdc.administration.v1.AdministrationOuterClass.ReadChangesRequest
import cdc.administration.v1.AdministrationOuterClass.ReadChangesResponse
import cdc.stream.v1.Changestream.StreamChangesRequest
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.runBlocking
import java.net.URI
import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * The optional `target` of the gRPC stream request and of the `readChanges`
 * request: set it travels in the request message, left out the message equals
 * the one built before the field existed.
 */
class PgChangeFeedGrpcClientTargetTest {
    private fun options() = PgChangeFeedClientOptions(URI("http://localhost:50051"), "reader-token")

    @Test
    fun `streamChanges without target sends the empty request`() = runBlocking {
        val transport = FakeGrpcStreamTransport.withMessages()
        val client = PgChangeFeedGrpcClient(transport, options())

        client.streamChanges().toList()

        assertEquals(StreamChangesRequest.getDefaultInstance(), transport.lastRequest)
    }

    @Test
    fun `streamChanges carries the target in the request`() = runBlocking {
        for ((target, expected) in listOf<Pair<String?, String>>(null to "", "" to "", "eu" to "eu")) {
            val transport = FakeGrpcStreamTransport.withMessages()
            val client = PgChangeFeedGrpcClient(transport, options())

            client.streamChanges(target = target).toList()

            assertEquals(expected, transport.lastRequest?.target)
        }
    }

    @Test
    fun `streamChanges with schema table and target sets all three`() = runBlocking {
        val transport = FakeGrpcStreamTransport.withMessages()
        val client = PgChangeFeedGrpcClient(transport, options())

        client.streamChanges(schema = "public", table = "orders", target = "eu").toList()

        assertEquals("public", transport.lastRequest?.schema)
        assertEquals("orders", transport.lastRequest?.table)
        assertEquals("eu", transport.lastRequest?.target)
    }

    @Test
    fun `administration readChanges reaches the transport with the target unchanged`() = runBlocking {
        for (target in listOf("", "eu")) {
            val transport = FakeAdministrationTransport.withResponse(ReadChangesResponse.getDefaultInstance())
            val client = AdministrationTestClientFactory.create(transport)

            client.readChanges(ReadChangesRequest.newBuilder().setSource("src").setSchema("public").setTarget(target).build())

            val sent = transport.lastRequest as ReadChangesRequest
            assertEquals(target, sent.target)
            assertEquals("public", sent.schema)
        }
    }
}
