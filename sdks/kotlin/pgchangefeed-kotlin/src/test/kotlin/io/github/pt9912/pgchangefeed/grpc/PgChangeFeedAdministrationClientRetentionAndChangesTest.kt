package io.github.pt9912.pgchangefeed.grpc

import cdc.administration.v1.AdministrationOuterClass.ChangeRecord
import cdc.administration.v1.AdministrationOuterClass.ReadChangesRequest
import cdc.administration.v1.AdministrationOuterClass.ReadChangesResponse
import cdc.administration.v1.AdministrationOuterClass.RunRetentionRequest
import cdc.administration.v1.AdministrationOuterClass.RunRetentionResponse
import com.google.protobuf.ByteString
import kotlinx.coroutines.runBlocking
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue

/**
 * Happy-path coverage for `runRetention` and `readChanges`, including the
 * zero-value boundary of each request's optional fields.
 */
class PgChangeFeedAdministrationClientRetentionAndChangesTest {
    @Test
    fun `runRetention happy path returns typed response`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(RunRetentionResponse.newBuilder().setDeleted(3).build())
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.runRetention(
            RunRetentionRequest.newBuilder().setSource("src").setMinAgeNanos(1_000_000_000).build(),
        )

        assertEquals(3L, response.deleted)
    }

    @Test
    fun `runRetention zero minAgeNanos is valid`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(RunRetentionResponse.newBuilder().setDeleted(0).build())
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.runRetention(RunRetentionRequest.newBuilder().setSource("src").setMinAgeNanos(0).build())

        assertEquals(0L, response.deleted)
        assertEquals(0L, (transport.lastRequest as RunRetentionRequest).minAgeNanos)
    }

    @Test
    fun `readChanges happy path returns typed response with images`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            ReadChangesResponse.newBuilder()
                .addChanges(
                    ChangeRecord.newBuilder()
                        .setCommitPosition(1).setChangeId("ch-1").setTransactionId("tx-1").setSourceTableId("t-1")
                        .setSchema("public").setTable("orders").setSequence(0).setOperation("INSERT")
                        .setNewImage(ByteString.copyFromUtf8("""{"id":1}"""))
                        .setSchemaVersion("sv-1").setCommittedAt("2026-09-19T00:00:00Z").setOrigin("wal")
                        .build(),
                )
                .build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.readChanges(
            ReadChangesRequest.newBuilder()
                .setSource("src").setSchema("public").setTable("orders").setFrom(1).setTo(10).setLimit(5)
                .build(),
        )

        assertEquals(1, response.changesCount)
        val change = response.changesList[0]
        assertEquals("ch-1", change.changeId)
        assertEquals("INSERT", change.operation)
        assertEquals("wal", change.origin)
        assertTrue(change.oldImage.isEmpty())
    }

    @Test
    fun `readChanges without optional fields leaves them at zero`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(ReadChangesResponse.getDefaultInstance())
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.readChanges(ReadChangesRequest.newBuilder().setSource("src").build())

        assertEquals(0, response.changesCount)
        val sent = transport.lastRequest as ReadChangesRequest
        assertEquals("", sent.schema)
        assertEquals("", sent.table)
        assertEquals(0L, sent.from)
        assertEquals(0L, sent.to)
        assertEquals(0L, sent.limit)
    }
}
