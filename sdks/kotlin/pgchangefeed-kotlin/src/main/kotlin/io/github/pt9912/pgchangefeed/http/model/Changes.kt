package io.github.pt9912.pgchangefeed.http.model

import com.google.gson.JsonElement
import com.google.gson.annotations.SerializedName

/**
 * One persisted change as returned by `GET /changes` (SPEC-022) — thirteen
 * fields: the ten of the Go domain type `model.Change`, plus
 * `commit_position`, `committed_at` and `origin`. `oldImage`/`newImage`
 * carry the row image as an embedded JSON value or a JSON `null` when
 * absent (`INSERT` has no `old_image`, `DELETE` has no `new_image`); kept
 * as [JsonElement] rather than a fixed shape, because the row image's own
 * shape depends on the captured table, not on this wire contract (same
 * choice as the C#/Python siblings, `System.Text.Json.JsonElement`/`Any`).
 *
 * [origin] is `wal` for a change captured from the replication stream and
 * `backfill` for an existing-rows change (LH-FA-CAP-009); it is the
 * server's string, and a response without the field (or with a JSON
 * `null`) reads as `wal`. Gson skips Kotlin default parameters, so the
 * wire value sits in the nullable [wireOrigin] and [origin] applies the
 * default. The live surfaces (gRPC, SSE, NATS) carry no `origin`.
 *
 * Gson's `JsonElement` adapter represents a present-but-`null` JSON value
 * as `JsonElement.JsonNull` (`isJsonNull == true`), not as a Kotlin `null`
 * reference — the field type stays nullable only to cover a value entirely
 * absent from the response body (outside the documented SPEC-022 shape, in
 * which the key is always present). A caller checks "no row image" via
 * `oldImage?.isJsonNull != false`, not via a plain null check.
 */
data class Change(
    @SerializedName("commit_position") val commitPosition: Long,
    @SerializedName("change_id") val changeId: String,
    @SerializedName("transaction_id") val transactionId: String,
    @SerializedName("source_table_id") val sourceTableId: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
    @SerializedName("sequence") val sequence: Long,
    @SerializedName("operation") val operation: String,
    @SerializedName("old_image") val oldImage: JsonElement?,
    @SerializedName("new_image") val newImage: JsonElement?,
    @SerializedName("schema_version") val schemaVersion: String,
    @SerializedName("committed_at") val committedAt: String,
    @SerializedName("origin") val wireOrigin: String? = null,
) {
    val origin: String
        get() = wireOrigin ?: "wal"
}

/** `ReadChanges` response (`200`) — an empty list on no match, never `404` (SPEC-022). */
data class ReadChangesResponse(
    @SerializedName("changes") val changes: List<Change>,
)
