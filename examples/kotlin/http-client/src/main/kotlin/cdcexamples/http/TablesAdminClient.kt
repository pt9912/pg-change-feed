package cdcexamples.http

import com.google.gson.annotations.SerializedName
import java.net.http.HttpClient

/**
 * EnableTableRequest/-Response spiegeln `POST /tables/enable`
 * (`LH-FA-CFG-001`): alle sieben Request-Felder Pflicht, `alreadyEnabled`
 * trägt die Idempotenz-Antwort ohne eigenen Statuscode. Form-Vorbild:
 * `examples/http-client/tables_admin.go`.
 */
data class EnableTableRequest(
    @SerializedName("source") val source: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
    @SerializedName("table_id") val tableId: String,
    @SerializedName("schema_version_id") val schemaVersionId: String,
    @SerializedName("version") val version: Long,
    @SerializedName("publication") val publication: String,
)

data class EnableTableResponse(
    @SerializedName("table_id") val tableId: String,
    @SerializedName("source") val source: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
    @SerializedName("already_enabled") val alreadyEnabled: Boolean,
)

/**
 * DisableTableRequest/-Response spiegeln `POST /tables/disable`
 * (`LH-FA-CFG-002`).
 */
data class DisableTableRequest(
    @SerializedName("source") val source: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
    @SerializedName("publication") val publication: String,
)

data class DisableTableResponse(
    @SerializedName("removed") val removed: Boolean,
    @SerializedName("retained") val retained: Boolean,
)

/**
 * TableStatusResponse spiegelt `GET /tables/status` (`LH-FA-CFG-003`):
 * `enabled`/`retained` trennen Erfassungs-Zustand und Herkunft.
 */
data class TableStatusResponse(
    @SerializedName("enabled") val enabled: Boolean,
    @SerializedName("retained") val retained: Boolean,
)

/**
 * TableStatusUrlBuilder baut die Lese-Adresse von `GET /tables/status`:
 * alle vier Parameter sind Pflicht. Form-Vorbild: `TableStatusURL` in
 * `examples/http-client/tables_admin.go`.
 */
object TableStatusUrlBuilder {
    fun build(addr: String, source: String, schema: String, table: String, publication: String): String {
        val query = listOf(
            "publication" to publication,
            "schema" to schema,
            "source" to source,
            "table" to table,
        ).joinToString("&") { (key, value) -> "$key=${UrlEncoding.encode(value)}" }
        return "http://$addr/tables/status?$query"
    }
}

/**
 * TablesAdminClient ruft die drei Tabellen-Verwaltungs-Fähigkeiten der
 * HTTP-/JSON-API auf, die [TablesClient] (Auflistung) nicht trägt.
 * Form-Vorbild: `examples/http-client/tables_admin.go`.
 */
object TablesAdminClient {
    fun enableTable(httpClient: HttpClient, cfg: Config): EnableTableResponse =
        RequestHelper.sendJson(
            httpClient, "POST", "http://${cfg.addr}/tables/enable", cfg.adminToken,
            EnableTableRequest(cfg.source, cfg.schema, cfg.table, cfg.tableId, cfg.schemaVersionId, cfg.version, cfg.publication),
            201,
        )

    fun disableTable(httpClient: HttpClient, cfg: Config): DisableTableResponse =
        RequestHelper.sendJson(
            httpClient, "POST", "http://${cfg.addr}/tables/disable", cfg.adminToken,
            DisableTableRequest(cfg.source, cfg.schema, cfg.table, cfg.publication), 200,
        )

    fun tableStatus(httpClient: HttpClient, cfg: Config): TableStatusResponse =
        RequestHelper.sendJson(
            httpClient, "GET", TableStatusUrlBuilder.build(cfg.addr, cfg.source, cfg.schema, cfg.table, cfg.publication),
            cfg.token, null, 200,
        )
}
