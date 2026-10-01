package cdcexamples.http

import com.google.gson.JsonElement
import com.google.gson.annotations.SerializedName
import java.net.http.HttpClient

/**
 * ReadChangeResponse spiegelt einen Change in der Antwortform von
 * `GET /changes` (`LH-FA-SST-006`). Die Row Images stehen als eingebettete
 * JSON-Werte; ein fehlendes Bild wird zu `JsonNull` (dieselbe Form wie
 * `StreamMessage` in
 * `examples/kotlin/nats-stream-client/.../Format.kt`). Form-Vorbild:
 * `examples/http-client/changes.go`.
 */
data class ReadChangeResponse(
    @SerializedName("commit_position") val commitPosition: Long,
    @SerializedName("change_id") val changeId: String,
    @SerializedName("transaction_id") val transactionId: String,
    @SerializedName("source_table_id") val sourceTableId: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
    @SerializedName("sequence") val sequence: Long,
    @SerializedName("operation") val operation: String,
    @SerializedName("old_image") val oldImage: JsonElement,
    @SerializedName("new_image") val newImage: JsonElement,
    @SerializedName("schema_version") val schemaVersion: String,
    @SerializedName("committed_at") val committedAt: String,
    @SerializedName("origin") val origin: String,
)

/**
 * ReadChangesResponse trägt den JSON-Response-Body bei Erfolg: ohne
 * Treffer eine leere, gesetzte Liste.
 */
data class ReadChangesResponse(
    @SerializedName("changes") val changes: List<ReadChangeResponse>,
)

/**
 * ChangesUrlBuilder baut die Lese-Adresse von `GET /changes`: `source` ist
 * Pflicht, die übrigen sechs Parameter sind unabhängig optional — ein
 * leerer Wert bleibt weg statt als leerer Query-Parameter zu erscheinen.
 * `target` wählt das Zustellziel einer Change (`ADR-0137`).
 * `from` ist die untere Positions-Grenze **einschließlich**, `to` die
 * obere Grenze **ausschließlich**. `GET /changes` lässt nur einen
 * unbekannten Parameter**namen** mit `400` enden; ein leerer Wert eines
 * bekannten optionalen Parameters (`schema`, `table`, `from`, `to`,
 * `limit`) bleibt für den Server unberücksichtigt und endet mit `200` —
 * dieser Client sendet ihn trotzdem gar nicht erst. Form-Vorbild:
 * `ChangesURL` in `examples/http-client/changes.go`.
 */
object ChangesUrlBuilder {
    fun build(
        addr: String,
        source: String,
        schema: String,
        table: String,
        from: String,
        to: String,
        limit: String,
        target: String = "",
    ): String {
        val params = sortedMapOf("source" to source)
        addIfNotEmpty(params, "schema", schema)
        addIfNotEmpty(params, "table", table)
        addIfNotEmpty(params, "target", target)
        addIfNotEmpty(params, "from", from)
        addIfNotEmpty(params, "to", to)
        addIfNotEmpty(params, "limit", limit)

        val query = params.entries.joinToString("&") { (key, value) -> "$key=${UrlEncoding.encode(value)}" }
        return "http://$addr/changes?$query"
    }

    private fun addIfNotEmpty(params: MutableMap<String, String>, key: String, value: String) {
        if (value.isNotEmpty()) {
            params[key] = value
        }
    }
}

/**
 * ChangesClient ruft den `reader`-Endpunkt `GET /changes` auf. Form-Vorbild:
 * `examples/http-client/changes.go`.
 */
object ChangesClient {
    fun readChanges(httpClient: HttpClient, cfg: Config): ReadChangesResponse =
        RequestHelper.sendJson(
            httpClient, "GET",
            ChangesUrlBuilder.build(cfg.addr, cfg.source, cfg.schema, cfg.table, cfg.from, cfg.to, cfg.limit, cfg.target),
            cfg.token, null, 200,
        )
}
