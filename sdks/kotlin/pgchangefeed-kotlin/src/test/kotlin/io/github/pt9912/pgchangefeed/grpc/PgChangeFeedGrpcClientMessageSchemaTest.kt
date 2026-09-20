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
 * Feld-für-Feld-Vollständigkeitstest gegen `SPEC-020`s zehn
 * Nachrichtenfelder, Consumer-seitig geprüft: jedes vom Fake gelieferte
 * Feld erreicht den Aufrufer über [PgChangeFeedGrpcClient.streamChanges]
 * unverändert — analog dem serverseitigen Feldvollständigkeits-Testmuster
 * in `internal/adapters/driving/grpc/server_test.go`
 * (`traegtTokenOeffnetStreamUndTraegtChange`), das dieser Test nicht
 * importiert (`ADR-0109` §Kontext Bindung „Import-Grenze, hier ohne
 * Ausnahme").
 */
class PgChangeFeedGrpcClientMessageSchemaTest {
    // Rot färbende Mutation (real geprüft, slice-sdk-kotlin-grpc-client-flaeche):
    // `PgChangeFeedGrpcClient.streamChanges()` um
    // `.map { it.toBuilder().clearSchema().build() }` erweitert, das das
    // `schema`-Feld vor der Ausgabe entfernt — dieser Test schlägt dann bei
    // der `schema`-Assertion fehl, statt still grün zu bleiben.
    @Test
    fun `streamChanges yields all SPEC-020 fields unchanged`() = runBlocking {
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
