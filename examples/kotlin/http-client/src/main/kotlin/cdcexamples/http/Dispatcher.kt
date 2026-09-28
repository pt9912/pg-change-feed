package cdcexamples.http

import com.google.gson.GsonBuilder
import java.net.http.HttpClient

/**
 * Dispatcher ruft die zum Verb gehörende Fähigkeit auf und liefert den
 * Antwort-Body als druckbaren Text. `tables` liest die rohe
 * Server-Antwort ([TablesClient]); die übrigen neun Fähigkeiten drucken
 * ihre typisierte Antwort eingerückt. Form-Vorbild: `dispatch` in
 * `examples/http-client/main.go`. Ein unbekanntes Verb bricht hier ein
 * zweites Mal ab (verteidigend, für den Fall eines Aufrufs ohne vorherige
 * [Validator.validate]).
 */
object Dispatcher {
    private val indentedGson = GsonBuilder().setPrettyPrinting().create()

    fun dispatch(httpClient: HttpClient, cfg: Config): String = when (cfg.verb) {
        "tables" -> TablesClient.listTables(httpClient, cfg)
        "changes" -> serializeIndented(ChangesClient.readChanges(httpClient, cfg))
        "register-consumer" -> serializeIndented(ConsumerClient.registerConsumer(httpClient, cfg))
        "acknowledge" -> serializeIndented(ConsumerClient.acknowledgeConsumer(httpClient, cfg))
        "consumer-position" -> serializeIndented(ConsumerClient.consumerPosition(httpClient, cfg))
        "remove-consumer" -> serializeIndented(ConsumerClient.removeConsumer(httpClient, cfg))
        "enable-table" -> serializeIndented(TablesAdminClient.enableTable(httpClient, applyEnableTableDefaults(cfg)))
        "disable-table" -> serializeIndented(TablesAdminClient.disableTable(httpClient, cfg))
        "table-status" -> serializeIndented(TablesAdminClient.tableStatus(httpClient, cfg))
        "retention-run" -> serializeIndented(RetentionClient.runRetention(httpClient, cfg))
        else -> throw IllegalStateException("unbekanntes --verb \"${cfg.verb}\"")
    }

    /**
     * applyEnableTableDefaults leitet `tableId`/`schemaVersionId` aus
     * Schema und Tabelle ab, wenn sie nicht gesetzt sind — dieselbe
     * Default-Konvention wie `internal/bootstrap/wiring.go`s
     * `administrationTableID`/`administrationSchemaVersionID`.
     */
    private fun applyEnableTableDefaults(cfg: Config): Config {
        val tableId = cfg.tableId.ifEmpty { "${cfg.schema}.${cfg.table}" }
        val schemaVersionId = cfg.schemaVersionId.ifEmpty { "$tableId-v1" }
        return cfg.copy(tableId = tableId, schemaVersionId = schemaVersionId)
    }

    private fun serializeIndented(value: Any): String = indentedGson.toJson(value)
}
