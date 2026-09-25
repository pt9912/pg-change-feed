package io.github.pt9912.pgchangefeed.sse.model

import com.google.gson.JsonElement
import com.google.gson.annotations.SerializedName

/**
 * One change delivered on the `SPEC-021` SSE stream — the ten fields the
 * server writes into a frame's `data:` line, one JSON object per change
 * (`event: change`). [oldImage]/[newImage] carry the row image as an
 * embedded JSON value or a JSON `null` when absent (`INSERT` has no
 * `old_image`, `DELETE` has no `new_image`); kept as [JsonElement] rather
 * than a fixed shape for the same reason as
 * [io.github.pt9912.pgchangefeed.http.model.Change]'s corresponding
 * fields — the row image's own shape depends on the captured table, not on
 * this wire contract.
 *
 * Kept as its own type rather than reused across surfaces: it is a
 * different wire contract than
 * [io.github.pt9912.pgchangefeed.http.model.Change] (`SPEC-022`, thirteen
 * fields including `commit_position`/`committed_at`/`origin`) and the generated
 * gRPC `cdc.stream.v1.Changestream.Change` message (`SPEC-020`, its own
 * protobuf-generated type,
 * [io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient]) — three
 * surfaces, three independent wire contracts that happen to share most
 * field names (same reasoning as the C# sibling's
 * `PgChangeFeed.Client.Sse.Models.Change`).
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
