package io.github.pt9912.pgchangefeed.sse

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue

/**
 * Feld-für-Feld-Vollständigkeitstest gegen `SPEC-021`s zehn
 * Nachrichtenfelder, Consumer-seitig geprüft: jedes vom Fake gelieferte
 * Feld erreicht den Aufrufer über [PgChangeFeedSseClient.streamChanges]
 * unverändert — analog dem gRPC-Pendant
 * `io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClientMessageSchemaTest`
 * und dem serverseitigen Feldvollständigkeits-Testmuster in
 * `internal/adapters/driving/http/sse_test.go`, das dieser Test nicht
 * importiert (`ADR-0109` §Kontext Bindung „Import-Grenze, hier ohne
 * Ausnahme").
 */
class PgChangeFeedSseClientMessageSchemaTest {

    // Rot färbende Mutation (real geprüft,
    // slice-sdk-kotlin-sse-client-flaeche): `PgChangeFeedSseClient.streamChanges`s
    // `yield(parseChange(...))`-Aufruf um `.copy(schema = "")` ergänzt, das
    // das `schema`-Feld vor der Ausgabe leert — dieser Test schlägt dann
    // bei der `schema`-Assertion fehl, statt still grün zu bleiben.
    @Test
    fun `streamChanges yields all SPEC-021 fields unchanged`() {
        val frameData = """{"change_id":"change-1","transaction_id":"tx-1","source_table_id":"table-1",""" +
            """"sequence":2,"operation":"UPDATE","old_image":{"id":1},""" +
            """"new_image":{"id":1,"bestellstatus":"bezahlt"},"schema_version":"table-1-v1",""" +
            """"schema":"public","table":"orders"}"""
        val (client, _) = TestClientFactory.create { _ ->
            FakeSseTransport.lineResponse(200, "event: change\ndata: $frameData\n\n")
        }

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
        val frameData = """{"change_id":"change-2","transaction_id":"tx-2","source_table_id":"table-1",""" +
            """"sequence":1,"operation":"INSERT","old_image":null,"new_image":{"id":9},""" +
            """"schema_version":"table-1-v1","schema":"public","table":"orders"}"""
        val (client, _) = TestClientFactory.create { _ ->
            FakeSseTransport.lineResponse(200, "event: change\ndata: $frameData\n\n")
        }

        val received = client.streamChanges().single()

        assertTrue(received.oldImage?.isJsonNull != false)
    }
}
