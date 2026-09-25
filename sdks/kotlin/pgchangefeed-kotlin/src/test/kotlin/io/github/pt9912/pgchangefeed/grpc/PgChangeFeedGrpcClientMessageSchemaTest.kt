package io.github.pt9912.pgchangefeed.grpc

import cdc.stream.v1.Changestream.Change
import com.google.protobuf.ByteString
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import kotlinx.coroutines.flow.single
import kotlinx.coroutines.runBlocking
import java.net.URI
import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * Field-by-field completeness test of the ten message fields, checked from
 * the consumer side: every field the fake delivers reaches the caller through
 * [PgChangeFeedGrpcClient.streamChanges] unchanged.
 */
class PgChangeFeedGrpcClientMessageSchemaTest {
    // Mutation that turns this test red (checked for real):
    // `PgChangeFeedGrpcClient.streamChanges()` um
    // `.map { it.toBuilder().clearSchema().build() }` erweitert, das das
    // `schema`-Feld vor der Ausgabe entfernt — dieser Test schlägt dann bei
    // der `schema`-Assertion fehl, statt still grün zu bleiben.
    @Test
    fun `streamChanges yields all ten fields unchanged`() = runBlocking {
        val change = Change.newBuilder()
            .setChangeId("change-1")
            .setTransactionId("tx-1")
            .setSourceTableId("table-1")
            .setSequence(2)
            .setOperation("UPDATE")
            .setOldImage(ByteString.copyFromUtf8("""{"id":1}"""))
            .setNewImage(ByteString.copyFromUtf8("""{"id":1,"bestellstatus":"bezahlt"}"""))
            .setSchemaVersion("table-1-v1")
            .setSchema("public")
            .setTable("orders")
            .build()
        val transport = FakeGrpcStreamTransport.withMessages(change)
        val options = PgChangeFeedClientOptions(URI("http://localhost:50051"), "reader-token")
        val client = PgChangeFeedGrpcClient(transport, options)

        val received = client.streamChanges().single()

        assertEquals("change-1", received.changeId)
        assertEquals("tx-1", received.transactionId)
        assertEquals("table-1", received.sourceTableId)
        assertEquals(2L, received.sequence)
        assertEquals("UPDATE", received.operation)
        assertEquals("""{"id":1}""", received.oldImage.toStringUtf8())
        assertEquals("""{"id":1,"bestellstatus":"bezahlt"}""", received.newImage.toStringUtf8())
        assertEquals("table-1-v1", received.schemaVersion)
        assertEquals("public", received.schema)
        assertEquals("orders", received.table)
    }
}
