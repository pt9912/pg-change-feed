package io.github.pt9912.pgchangefeed.nats.model

import com.google.gson.JsonElement
import com.google.gson.annotations.SerializedName

/**
 * One change delivered on the `SPEC-024` NATS full-content stream — the ten
 * fields the server writes into a message's JSON payload, identical to the
 * `SPEC-021` SSE message shape ([io.github.pt9912.pgchangefeed.sse.model.Change]),
 * carried over NATS instead of HTTP. [oldImage]/[newImage] carry the row
 * image as an embedded JSON value or a JSON `null` when absent (`INSERT` has
 * no `old_image`, `DELETE` has no `new_image`); kept as [JsonElement] rather
 * than a fixed shape for the same reason as
 * [io.github.pt9912.pgchangefeed.http.model.Change]'s corresponding fields —
 * the row image's own shape depends on the captured table, not on this wire
 * contract.
 *
 * Kept as its own type rather than reused across surfaces: it is a
 * structurally identical, but independently versioned wire contract from
 * [io.github.pt9912.pgchangefeed.sse.model.Change] (`SPEC-024` vs
 * `SPEC-021`) and the generated gRPC
 * `cdc.stream.v1.Changestream.Change` message (`SPEC-020`) — four surfaces,
 * four independent wire contracts that happen to share most field names
 * (same reasoning as the C# sibling's `PgChangeFeed.Client.Nats.Models.Change`).
 */
data class Change(
    @SerializedName("change_id") val changeId: String,
    @SerializedName("transaction_id") val transactionId: String,
    @SerializedName("source_table_id") val sourceTableId: String,
    @SerializedName("sequence") val sequence: Long,
    @SerializedName("operation") val operation: String,
    @SerializedName("old_image") val oldImage: JsonElement?,
    @SerializedName("new_image") val newImage: JsonElement?,
    @SerializedName("schema_version") val schemaVersion: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
)
