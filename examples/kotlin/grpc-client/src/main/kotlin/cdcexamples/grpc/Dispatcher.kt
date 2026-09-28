package cdcexamples.grpc

import cdc.administration.v1.AdministrationGrpcKt

/**
 * Dispatcher ruft die zum Verb gehörende Administration-RPC auf und liefert
 * die formatierte Ausgabe. `stream` läuft nicht hier durch — er hat eine
 * eigene, dauerhafte Schleife (`runStream` in `Main.kt`). Ein unbekanntes
 * Verb bricht hier ein zweites Mal ab (verteidigend, für den Fall eines
 * Aufrufs ohne vorherige [Validator.validate]). Form-Vorbild:
 * `examples/csharp/grpc-client/Dispatcher.cs`, `dispatchAdmin` in
 * `examples/grpc-client/main.go`.
 */
object Dispatcher {
    suspend fun dispatchAdmin(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String = when (cfg.verb) {
        "register-consumer" -> ConsumerClient.registerConsumer(client, cfg)
        "acknowledge-consumer" -> ConsumerClient.acknowledgeConsumer(client, cfg)
        "get-consumer-position" -> ConsumerClient.getConsumerPosition(client, cfg)
        "remove-consumer" -> ConsumerClient.removeConsumer(client, cfg)
        "enable-table" -> TablesAdminClient.enableTable(client, applyEnableTableDefaults(cfg))
        "disable-table" -> TablesAdminClient.disableTable(client, cfg)
        "get-table-status" -> TablesAdminClient.getTableStatus(client, cfg)
        "list-tables" -> TablesAdminClient.listTables(client, cfg)
        "run-retention" -> RetentionClient.runRetention(client, cfg)
        "read-changes" -> ChangesClient.readChanges(client, cfg)
        "diagnose" -> DiagnoseClient.diagnose(client, cfg)
        else -> throw IllegalStateException("unbekanntes --verb \"${cfg.verb}\"")
    }

    /**
     * applyEnableTableDefaults leitet `tableId`/`schemaVersionId` aus Schema
     * und Tabelle ab, wenn sie nicht gesetzt sind — dieselbe
     * Default-Konvention wie `internal/bootstrap/wiring.go`s
     * `administrationTableID`/`administrationSchemaVersionID`.
     */
    private fun applyEnableTableDefaults(cfg: Config): Config {
        val tableId = cfg.tableId.ifEmpty { "${cfg.schema}.${cfg.table}" }
        val schemaVersionId = cfg.schemaVersionId.ifEmpty { "$tableId-v1" }
        return cfg.copy(tableId = tableId, schemaVersionId = schemaVersionId)
    }
}
