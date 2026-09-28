package cdcexamples.grpc

import cdc.administration.v1.AdministrationOuterClass.AcknowledgeConsumerResponse
import cdc.administration.v1.AdministrationOuterClass.DiagnoseResponse
import cdc.administration.v1.AdministrationOuterClass.DisableTableResponse
import cdc.administration.v1.AdministrationOuterClass.EnableTableResponse
import cdc.administration.v1.AdministrationOuterClass.GetConsumerPositionResponse
import cdc.administration.v1.AdministrationOuterClass.GetTableStatusResponse
import cdc.administration.v1.AdministrationOuterClass.ListTablesResponse
import cdc.administration.v1.AdministrationOuterClass.ReadChangesResponse
import cdc.administration.v1.AdministrationOuterClass.RegisterConsumerResponse
import cdc.administration.v1.AdministrationOuterClass.RemoveConsumerResponse
import cdc.administration.v1.AdministrationOuterClass.RunRetentionResponse
import cdc.stream.v1.Changestream.Change

/**
 * Format baut die Ausgabezeile einer empfangenen Stream-Nachricht sowie die
 * Ausgabezeilen jeder der elf Administration-Antworten — reine Funktionen,
 * netzlos testbar. Form-Vorbild: `examples/csharp/grpc-client/Format.cs`,
 * `examples/grpc-client/format.go`, `consumer.go`, `tables_admin.go`,
 * `changes.go`, `retention.go`, `diagnose.go`.
 *
 * `Change` liegt verschachtelt unter dem generierten Java-Deskriptor
 * `cdc.stream.v1.Changestream` — die `.proto` setzt kein
 * `option java_multiple_files`; anders als beim C#-Codegen (flache Klassen
 * im Namensraum `Cdc.Stream.V1`) verschachtelt der Java-/Kotlin-Codegen
 * Nachrichtentypen standardmäßig unter einer Datei-Außenklasse, deren Name
 * sich aus dem `.proto`-Dateinamen ableitet (`changestream.proto` ->
 * `Changestream`). Der Import macht das für den Aufrufer unsichtbar. Die
 * Administration-Nachrichten liegen aus demselben Grund verschachtelt unter
 * `AdministrationOuterClass` statt unter `Administration` — der Dateiname
 * `administration.proto` träfe sonst denselben Namen wie der Dienst
 * `Administration`, protoc weicht dem Namenskonflikt mit dem
 * `OuterClass`-Zusatz aus.
 */
object Format {
    /**
     * formatChange baut Tabelle, Operation und das neue Row Image, dazu die
     * Kennung, an der sich die Nachricht gegen den Lesezugriffsweg
     * `cdc.changes` halten lässt (`change_id`).
     */
    fun formatChange(change: Change): String {
        return "grpc-client: change_id=${change.changeId} table=${change.schema}.${change.table} " +
            "operation=${change.operation} new_image=${change.newImage.toStringUtf8()}"
    }

    fun formatRegisterConsumer(resp: RegisterConsumerResponse): String =
        "grpc-client: consumer_id=${resp.consumerId} name=${resp.name} already_registered=${resp.alreadyRegistered}"

    fun formatAcknowledgeConsumer(resp: AcknowledgeConsumerResponse): String =
        "grpc-client: consumer_id=${resp.consumerId} source_id=${resp.sourceId} offset=${resp.offset}"

    fun formatConsumerPosition(resp: GetConsumerPositionResponse): String =
        "grpc-client: consumer_id=${resp.consumerId} source_id=${resp.sourceId} offset=${resp.offset} acknowledged=${resp.acknowledged}"

    fun formatRemoveConsumer(resp: RemoveConsumerResponse): String =
        "grpc-client: consumer_id=${resp.consumerId} removed=${resp.removed}"

    fun formatEnableTable(resp: EnableTableResponse): String =
        "grpc-client: table_id=${resp.tableId} source=${resp.source} schema=${resp.schema} table=${resp.table} already_enabled=${resp.alreadyEnabled}"

    fun formatDisableTable(resp: DisableTableResponse): String =
        "grpc-client: removed=${resp.removed} retained=${resp.retained}"

    fun formatTableStatus(resp: GetTableStatusResponse): String =
        "grpc-client: enabled=${resp.enabled} retained=${resp.retained}"

    /** formatListTables listet jede Tabelle der Antwort — `tables` und `retained` je eine `SourceTable`-Zeile. */
    fun formatListTables(resp: ListTablesResponse): String {
        val lines = mutableListOf("grpc-client: tables=${resp.tablesList.size} retained=${resp.retainedList.size}")
        for (t in resp.tablesList) {
            lines.add("  table: table_id=${t.tableId} source=${t.source} schema=${t.schema} table=${t.table}")
        }
        for (t in resp.retainedList) {
            lines.add("  retained: table_id=${t.tableId} source=${t.source} schema=${t.schema} table=${t.table}")
        }
        return lines.joinToString("\n")
    }

    fun formatRunRetention(resp: RunRetentionResponse): String =
        "grpc-client: deleted=${resp.deleted}"

    /** formatReadChanges listet jede gelesene Änderung — dieselben Felder wie der Lesezugriffsweg `cdc.changes`. */
    fun formatReadChanges(resp: ReadChangesResponse): String {
        val lines = mutableListOf("grpc-client: changes=${resp.changesList.size}")
        for (c in resp.changesList) {
            lines.add(
                "  change_id=${c.changeId} table=${c.schema}.${c.table} operation=${c.operation} " +
                    "commit_position=${c.commitPosition} origin=${c.origin} new_image=${c.newImage.toStringUtf8()}",
            )
        }
        return lines.joinToString("\n")
    }

    /**
     * formatDiagnose baut denselben Bericht wie der CLI-Sondermodus: ein
     * `known`/`present`-Feld auf `false` trägt den jeweiligen
     * Abwesenheits-Fall (kein Lebenszeichen, kein Blocker, unbekannte
     * Schätzung/unbekannter Rückstand).
     */
    fun formatDiagnose(resp: DiagnoseResponse): String {
        val lines = mutableListOf(
            "grpc-client: heartbeat_known=${resp.heartbeat.known} heartbeat_age_seconds=${resp.heartbeat.ageSeconds} " +
                "heartbeat_error_class=${resp.heartbeat.errorClass} capture_lag=${resp.captureLag} storage_bytes=${resp.storageBytes}",
        )

        for (cl in resp.consumerLagsList) {
            lines.add("  consumer_lag: consumer_id=${cl.consumerId} known=${cl.known} lag=${cl.lag}")
        }

        if (resp.retentionBlocker.present) {
            val rb = resp.retentionBlocker
            lines.add(
                "  retention_blocker: consumer_id=${rb.consumerId} name=${rb.name} " +
                    "acknowledged_position=${rb.acknowledgedPosition} backlog_known=${rb.backlogKnown} backlog=${rb.backlog}",
            )
        } else {
            lines.add("  retention_blocker: kein Blocker")
        }

        for (b in resp.backfillList) {
            lines.add(
                "  backfill: schema=${b.schema} table=${b.table} status=${b.status} rows_copied=${b.rowsCopied} " +
                    "estimated_rows_known=${b.estimatedRowsKnown} estimated_rows=${b.estimatedRows} " +
                    "warn_size=${b.warnEstimatedSize} warn_duration=${b.warnDuration} error_message=${b.errorMessage}",
            )
        }

        return lines.joinToString("\n")
    }
}
