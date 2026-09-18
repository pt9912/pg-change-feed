package cdcexamples.natsstream

import com.google.gson.JsonElement
import com.google.gson.annotations.SerializedName

/**
 * StreamMessage trägt dasselbe Nachrichtenschema wie der Publisher
 * (`internal/adapters/driven/natsstream`, `SPEC-024`) — hier eigenständig
 * geführt: dieses Beispiel importiert keinen privaten Paketbaum dieses
 * Repositories (`SPEC-023`).
 */
data class StreamMessage(
    @SerializedName("change_id") val changeId: String,
    @SerializedName("transaction_id") val transactionId: String,
    @SerializedName("source_table_id") val sourceTableId: String,
    @SerializedName("sequence") val sequence: Long,
    @SerializedName("operation") val operation: String,
    @SerializedName("old_image") val oldImage: JsonElement,
    @SerializedName("new_image") val newImage: JsonElement,
    @SerializedName("schema_version") val schemaVersion: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
)

/**
 * Format baut die Ausgabezeile einer empfangenen Stream-Nachricht
 * (`LH-FA-SST-008`): Tabelle, Operation und das neue Row Image, dazu die
 * Kennung, an der sich die Nachricht gegen den Lesezugriffsweg
 * `cdc.changes` halten lässt (`change_id`) — dieselbe Form wie beim
 * gRPC-Beispiel
 * (`examples/kotlin/grpc-client/src/main/kotlin/cdcexamples/grpc/Format.kt`).
 */
object Format {
    fun formatChange(change: StreamMessage): String =
        "nats-stream-client: change_id=${change.changeId} table=${change.schema}.${change.table} " +
            "operation=${change.operation} new_image=${change.newImage}"
}
