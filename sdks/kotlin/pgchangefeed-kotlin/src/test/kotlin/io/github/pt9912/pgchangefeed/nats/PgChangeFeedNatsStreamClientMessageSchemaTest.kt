package io.github.pt9912.pgchangefeed.nats

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue

/**
 * Field-by-field completeness test of the ten message fields, checked from
 * the consumer side: every field the fake delivers reaches the caller through
 * [PgChangeFeedNatsStreamClient.streamChanges] unchanged — analogous to the
 * SSE counterpart
 * `io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClientMessageSchemaTest`
 * and the gRPC counterpart
 * `io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClientMessageSchemaTest`.
 */
class PgChangeFeedNatsStreamClientMessageSchemaTest {

    // Mutation that turns this test red (checked for real):
    // `PgChangeFeedNatsStreamClient.streamChanges`s `yield(parseChange(payload))`-Aufruf
    // um `.copy(schema = "")` ergänzt, das das `schema`-Feld vor der Ausgabe
    // leert — dieser Test schlägt dann bei der `schema`-Assertion fehl,
    // statt still grün zu bleiben.
    @Test
    fun `streamChanges yields all ten fields unchanged`() {
        val payload = """{"change_id":"change-1","transaction_id":"tx-1","source_table_id":"table-1",""" +
            """"sequence":2,"operation":"UPDATE","old_image":{"id":1},""" +
            """"new_image":{"id":1,"bestellstatus":"bezahlt"},"schema_version":"table-1-v1",""" +
            """"schema":"public","table":"orders"}"""
        val transport = FakeNatsStreamTransport.withPayloads(payload.toByteArray(Charsets.UTF_8))
        val client = PgChangeFeedNatsStreamClient(transport)

        val received = client.streamChanges().single()

        assertEquals("change-1", received.changeId)
        assertEquals("tx-1", received.transactionId)
        assertEquals("table-1", received.sourceTableId)
        assertEquals(2L, received.sequence)
        assertEquals("UPDATE", received.operation)
        assertEquals("""{"id":1}""", received.oldImage.toString())
        assertEquals("""{"id":1,"bestellstatus":"bezahlt"}""", received.newImage.toString())
        assertEquals("table-1-v1", received.schemaVersion)
        assertEquals("public", received.schema)
        assertEquals("orders", received.table)
    }

    @Test
    fun `an absent row image is delivered as a JSON null, not a missing field`() {
        val payload = """{"change_id":"change-2","transaction_id":"tx-2","source_table_id":"table-1",""" +
            """"sequence":1,"operation":"INSERT","old_image":null,"new_image":{"id":9},""" +
            """"schema_version":"table-1-v1","schema":"public","table":"orders"}"""
        val transport = FakeNatsStreamTransport.withPayloads(payload.toByteArray(Charsets.UTF_8))
        val client = PgChangeFeedNatsStreamClient(transport)

        val received = client.streamChanges().single()

        assertTrue(received.oldImage?.isJsonNull != false)
    }

    @Test
    fun `streamChanges yields every message in order`() {
        fun payload(changeId: String, operation: String): ByteArray {
            val json = """{"change_id":"$changeId",""" +
                """"transaction_id":"tx-1","source_table_id":"table-1","sequence":1,""" +
                """"operation":"$operation","old_image":null,"new_image":{"id":1},""" +
                """"schema_version":"table-1-v1","schema":"public","table":"orders"}"""
            return json.toByteArray(Charsets.UTF_8)
        }

        val transport = FakeNatsStreamTransport.withPayloads(
            payload("change-1", "INSERT"),
            payload("change-2", "UPDATE"),
        )
        val client = PgChangeFeedNatsStreamClient(transport)

        val received = client.streamChanges().map { "${it.changeId}:${it.operation}" }.toList()

        assertEquals(listOf("change-1:INSERT", "change-2:UPDATE"), received)
    }
}
