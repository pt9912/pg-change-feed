package cdcexamples.grpc

import cdc.administration.v1.AdministrationOuterClass.AcknowledgeConsumerResponse
import cdc.administration.v1.AdministrationOuterClass.BackfillTableStatus
import cdc.administration.v1.AdministrationOuterClass.ChangeRecord
import cdc.administration.v1.AdministrationOuterClass.ConsumerLag
import cdc.administration.v1.AdministrationOuterClass.DiagnoseResponse
import cdc.administration.v1.AdministrationOuterClass.DisableTableResponse
import cdc.administration.v1.AdministrationOuterClass.EnableTableResponse
import cdc.administration.v1.AdministrationOuterClass.GetConsumerPositionResponse
import cdc.administration.v1.AdministrationOuterClass.GetTableStatusResponse
import cdc.administration.v1.AdministrationOuterClass.HeartbeatStatus
import cdc.administration.v1.AdministrationOuterClass.ListTablesResponse
import cdc.administration.v1.AdministrationOuterClass.ReadChangesResponse
import cdc.administration.v1.AdministrationOuterClass.RegisterConsumerResponse
import cdc.administration.v1.AdministrationOuterClass.RemoveConsumerResponse
import cdc.administration.v1.AdministrationOuterClass.RetentionBlocker
import cdc.administration.v1.AdministrationOuterClass.RunRetentionResponse
import cdc.administration.v1.AdministrationOuterClass.SourceTable
import cdc.stream.v1.Changestream.Change
import com.google.protobuf.ByteString
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue

/**
 * Prüft [Format]: die Stream-Zeile und jede der elf Administration-
 * Antwortzeilen — reine Funktionen, netzlos testbar. Form-Vorbild:
 * `examples/csharp/grpc-client/GrpcClient.Tests/FormatTests.cs`,
 * `examples/grpc-client/format_test.go`.
 */
class FormatTest {
    @Test
    fun formatChangeBuildsExpectedLine() {
        val change = Change.newBuilder()
            .setChangeId("chg-1")
            .setSchema("public")
            .setTable("widgets")
            .setOperation("INSERT")
            .setNewImage(ByteString.copyFromUtf8("{\"id\":1}"))
            .build()

        val got = Format.formatChange(change)

        assertEquals(
            "grpc-client: change_id=chg-1 table=public.widgets operation=INSERT new_image={\"id\":1}",
            got,
        )
    }

    @Test
    fun formatChangeHandlesEmptyNewImage() {
        val change = Change.newBuilder()
            .setChangeId("chg-2")
            .setSchema("public")
            .setTable("widgets")
            .setOperation("DELETE")
            .build()

        val got = Format.formatChange(change)

        assertEquals("grpc-client: change_id=chg-2 table=public.widgets operation=DELETE new_image=", got)
    }

    @Test
    fun formatRegisterConsumerCarriesIdempotencyFlag() {
        val resp = RegisterConsumerResponse.newBuilder().setConsumerId("c-1").setName("Consumer").setAlreadyRegistered(true).build()
        assertEquals("grpc-client: consumer_id=c-1 name=Consumer already_registered=true", Format.formatRegisterConsumer(resp))
    }

    @Test
    fun formatAcknowledgeConsumerCarriesOffset() {
        val resp = AcknowledgeConsumerResponse.newBuilder().setConsumerId("c-1").setSourceId("quelle-1").setOffset(42).build()
        assertEquals("grpc-client: consumer_id=c-1 source_id=quelle-1 offset=42", Format.formatAcknowledgeConsumer(resp))
    }

    @Test
    fun formatConsumerPositionDistinguishesUnacknowledgedFromZeroOffset() {
        val resp = GetConsumerPositionResponse.newBuilder()
            .setConsumerId("c-1").setSourceId("quelle-1").setOffset(0).setAcknowledged(false)
            .build()
        assertEquals("grpc-client: consumer_id=c-1 source_id=quelle-1 offset=0 acknowledged=false", Format.formatConsumerPosition(resp))
    }

    @Test
    fun formatRemoveConsumerCarriesRemovedFlag() {
        val resp = RemoveConsumerResponse.newBuilder().setConsumerId("c-1").setRemoved(false).build()
        assertEquals("grpc-client: consumer_id=c-1 removed=false", Format.formatRemoveConsumer(resp))
    }

    @Test
    fun formatEnableTableCarriesIdempotencyFlag() {
        val resp = EnableTableResponse.newBuilder()
            .setTableId("public.orders").setSource("quelle-1").setSchema("public").setTable("orders").setAlreadyEnabled(false)
            .build()
        assertEquals(
            "grpc-client: table_id=public.orders source=quelle-1 schema=public table=orders already_enabled=false",
            Format.formatEnableTable(resp),
        )
    }

    @Test
    fun formatDisableTableCarriesRemovedAndRetained() {
        val resp = DisableTableResponse.newBuilder().setRemoved(true).setRetained(true).build()
        assertEquals("grpc-client: removed=true retained=true", Format.formatDisableTable(resp))
    }

    @Test
    fun formatTableStatusCarriesEnabledAndRetained() {
        val resp = GetTableStatusResponse.newBuilder().setEnabled(true).setRetained(false).build()
        assertEquals("grpc-client: enabled=true retained=false", Format.formatTableStatus(resp))
    }

    @Test
    fun formatListTablesListsTablesAndRetainedSeparately() {
        val resp = ListTablesResponse.newBuilder()
            .addTables(SourceTable.newBuilder().setTableId("public.orders").setSource("quelle-1").setSchema("public").setTable("orders"))
            .addRetained(SourceTable.newBuilder().setTableId("public.legacy").setSource("quelle-1").setSchema("public").setTable("legacy"))
            .build()

        val got = Format.formatListTables(resp)

        assertEquals(
            "grpc-client: tables=1 retained=1\n" +
                "  table: table_id=public.orders source=quelle-1 schema=public table=orders\n" +
                "  retained: table_id=public.legacy source=quelle-1 schema=public table=legacy",
            got,
        )
    }

    @Test
    fun formatListTablesEmptyListsCarryZeroCounts() {
        val resp = ListTablesResponse.getDefaultInstance()
        assertEquals("grpc-client: tables=0 retained=0", Format.formatListTables(resp))
    }

    @Test
    fun formatRunRetentionCarriesDeletedCount() {
        val resp = RunRetentionResponse.newBuilder().setDeleted(7).build()
        assertEquals("grpc-client: deleted=7", Format.formatRunRetention(resp))
    }

    @Test
    fun formatReadChangesListsEachChange() {
        val resp = ReadChangesResponse.newBuilder()
            .addChanges(
                ChangeRecord.newBuilder()
                    .setChangeId("c-1").setSchema("public").setTable("orders").setOperation("INSERT")
                    .setCommitPosition(100).setOrigin("wal").setNewImage(ByteString.copyFromUtf8("{\"id\":1}")),
            )
            .build()

        val got = Format.formatReadChanges(resp)

        assertEquals(
            "grpc-client: changes=1\n" +
                "  change_id=c-1 table=public.orders operation=INSERT commit_position=100 origin=wal new_image={\"id\":1}",
            got,
        )
    }

    @Test
    fun formatReadChangesEmptyResultCarriesZeroCount() {
        val resp = ReadChangesResponse.getDefaultInstance()
        assertEquals("grpc-client: changes=0", Format.formatReadChanges(resp))
    }

    @Test
    fun formatDiagnoseNoRetentionBlockerPrintsAbsenceLine() {
        val resp = DiagnoseResponse.newBuilder()
            .setHeartbeat(HeartbeatStatus.newBuilder().setKnown(false))
            .setRetentionBlocker(RetentionBlocker.newBuilder().setPresent(false))
            .build()

        val got = Format.formatDiagnose(resp)

        assertTrue(got.contains("heartbeat_known=false"))
        assertTrue(got.contains("  retention_blocker: kein Blocker"))
    }

    @Test
    fun formatDiagnosePresentRetentionBlockerCarriesConsumerFields() {
        val resp = DiagnoseResponse.newBuilder()
            .setHeartbeat(HeartbeatStatus.newBuilder().setKnown(true).setAgeSeconds(1.5).setErrorClass(""))
            .setRetentionBlocker(
                RetentionBlocker.newBuilder()
                    .setPresent(true).setConsumerId("c-1").setName("Consumer")
                    .setAcknowledgedPosition(10).setBacklogKnown(true).setBacklog(5),
            )
            .build()

        val got = Format.formatDiagnose(resp)

        assertTrue(
            got.contains(
                "retention_blocker: consumer_id=c-1 name=Consumer acknowledged_position=10 backlog_known=true backlog=5",
            ),
        )
    }

    @Test
    fun formatDiagnoseListsConsumerLagsAndBackfillRows() {
        val resp = DiagnoseResponse.newBuilder()
            .setHeartbeat(HeartbeatStatus.newBuilder().setKnown(true))
            .setRetentionBlocker(RetentionBlocker.newBuilder().setPresent(false))
            .addConsumerLags(ConsumerLag.newBuilder().setConsumerId("c-1").setKnown(true).setLag(2.0))
            .addBackfill(
                BackfillTableStatus.newBuilder()
                    .setSchema("public").setTable("orders").setStatus("completed").setRowsCopied(100)
                    .setEstimatedRowsKnown(true).setEstimatedRows(100)
                    .setWarnEstimatedSize(false).setWarnDuration(false).setErrorMessage(""),
            )
            .build()

        val got = Format.formatDiagnose(resp)

        assertTrue(got.contains("consumer_lag: consumer_id=c-1 known=true lag=2.0"))
        assertTrue(got.contains("backfill: schema=public table=orders status=completed rows_copied=100"))
    }
}
