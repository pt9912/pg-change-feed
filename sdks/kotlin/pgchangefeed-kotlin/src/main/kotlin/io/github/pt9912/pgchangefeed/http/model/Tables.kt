package io.github.pt9912.pgchangefeed.http.model

import com.google.gson.annotations.SerializedName

/**
 * The table to capture, all seven fields mandatory. `source` and `publication`
 * name the source and its PostgreSQL publication. `tableId` is the id the
 * table is captured under (every change of the table carries it as
 * `source_table_id`); `schemaVersionId` and `version` (must be >= 1) identify
 * the table's schema version (changes carry the id as `schema_version`).
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

/**
 * The captured table; `alreadyEnabled` is true when it was captured before. A
 * table that does not exist in the source database ends `404` instead.
 */
data class EnableTableResponse(
    @SerializedName("table_id") val tableId: String,
    @SerializedName("source") val source: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
    @SerializedName("already_enabled") val alreadyEnabled: Boolean,
)

/**
 * The table to stop capturing, all four fields mandatory.
 */
data class DisableTableRequest(
    @SerializedName("source") val source: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
    @SerializedName("publication") val publication: String,
)

/**
 * `retained` is true when changes already stored for the table remain
 * readable. A table that does not exist in the source database ends `404`
 * instead.
 */
data class DisableTableResponse(
    @SerializedName("removed") val removed: Boolean,
    @SerializedName("retained") val retained: Boolean,
)

/**
 * `enabled` is true for a captured table; `retained` is true for a table that
 * is no longer captured but whose stored changes remain. Both `false` reads a
 * table that was never enabled. A table that does not exist in the source
 * database ends `404` instead.
 */
data class TableStatusResponse(
    @SerializedName("enabled") val enabled: Boolean,
    @SerializedName("retained") val retained: Boolean,
)

/** One table of a [ListTablesResponse]. */
data class TableInfo(
    @SerializedName("table_id") val tableId: String,
    @SerializedName("source") val source: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
)

/**
 * `tables` are the captured tables; `retained` are the tables that are no
 * longer captured but whose stored changes remain. Both lists are empty when
 * no table was ever enabled.
 */
data class ListTablesResponse(
    @SerializedName("tables") val tables: List<TableInfo>,
    @SerializedName("retained") val retained: List<TableInfo>,
)
