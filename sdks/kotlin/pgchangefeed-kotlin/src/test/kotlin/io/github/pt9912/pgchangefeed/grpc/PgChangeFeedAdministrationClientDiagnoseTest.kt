package io.github.pt9912.pgchangefeed.grpc

import cdc.administration.v1.AdministrationOuterClass.BackfillTableStatus
import cdc.administration.v1.AdministrationOuterClass.ConsumerLag
import cdc.administration.v1.AdministrationOuterClass.DiagnoseRequest
import cdc.administration.v1.AdministrationOuterClass.DiagnoseResponse
import cdc.administration.v1.AdministrationOuterClass.HeartbeatStatus
import cdc.administration.v1.AdministrationOuterClass.RetentionBlocker
import kotlinx.coroutines.runBlocking
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertTrue

/**
 * Coverage for `diagnose`: the populated report and the absence cases each
 * `known`/`present`/`*Known` field carries when the source never produced
 * the respective signal.
 */
class PgChangeFeedAdministrationClientDiagnoseTest {
    @Test
    fun `diagnose happy path returns full report`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            DiagnoseResponse.newBuilder()
                .setHeartbeat(HeartbeatStatus.newBuilder().setKnown(true).setAgeSeconds(1.5).setErrorClass("").build())
                .setCaptureLag(0.2)
                .addConsumerLags(ConsumerLag.newBuilder().setConsumerId("c-1").setKnown(true).setLag(3.0).build())
                .setRetentionBlocker(
                    RetentionBlocker.newBuilder()
                        .setPresent(true).setConsumerId("c-1").setName("n-1").setAcknowledgedPosition(10)
                        .setBacklogKnown(true).setBacklog(5)
                        .build(),
                )
                .setStorageBytes(1024.0)
                .addBackfill(
                    BackfillTableStatus.newBuilder()
                        .setSchema("public").setTable("orders").setStatus("completed").setRowsCopied(100)
                        .setEstimatedRowsKnown(true).setEstimatedRows(100).setWarnEstimatedSize(false)
                        .setWarnDuration(false).setErrorMessage("")
                        .build(),
                )
                .build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.diagnose(DiagnoseRequest.newBuilder().setSource("src").build())

        assertTrue(response.heartbeat.known)
        assertTrue(response.retentionBlocker.present)
        assertEquals(1, response.consumerLagsCount)
        assertEquals(1, response.backfillCount)
    }

    @Test
    fun `diagnose no heartbeat ever reports known false`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            DiagnoseResponse.newBuilder()
                .setHeartbeat(HeartbeatStatus.newBuilder().setKnown(false).build())
                .setRetentionBlocker(RetentionBlocker.newBuilder().setPresent(false).build())
                .build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.diagnose(DiagnoseRequest.newBuilder().setSource("src").build())

        assertFalse(response.heartbeat.known)
        assertFalse(response.retentionBlocker.present)
        assertEquals(0, response.consumerLagsCount)
        assertEquals(0, response.backfillCount)
    }

    @Test
    fun `diagnose unknown estimate reports estimatedRowsKnown false`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            DiagnoseResponse.newBuilder()
                .setHeartbeat(HeartbeatStatus.newBuilder().setKnown(true).build())
                .setRetentionBlocker(RetentionBlocker.newBuilder().setPresent(false).build())
                .addBackfill(
                    BackfillTableStatus.newBuilder()
                        .setSchema("public").setTable("orders").setStatus("running")
                        .setEstimatedRowsKnown(false).setEstimatedRows(0)
                        .build(),
                )
                .build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.diagnose(DiagnoseRequest.newBuilder().setSource("src").build())

        assertFalse(response.backfillList[0].estimatedRowsKnown)
    }
}
