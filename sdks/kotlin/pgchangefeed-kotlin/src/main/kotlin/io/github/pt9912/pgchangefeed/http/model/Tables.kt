package io.github.pt9912.pgchangefeed.http.model

import com.google.gson.annotations.SerializedName

/**
 * `EnableTable` request — `POST /tables/enable` (SPEC-018), all seven
 * fields mandatory, `version` must be >= 1.
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
 * `EnableTable` response (`201`). A table physically missing at the source
 * ends `404` instead (SPEC-018).
 */
data class EnableTableResponse(
    @SerializedName("table_id") val tableId: String,
    @SerializedName("source") val source: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
    @SerializedName("already_enabled") val alreadyEnabled: Boolean,
)

/**
 * `DisableTable` request — `POST /tables/disable` (SPEC-018), all four
 * fields mandatory.
 */
data class DisableTableRequest(
    @SerializedName("source") val source: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
    @SerializedName("publication") val publication: String,
)

/**
 * `DisableTable` response (`200`). A table physically missing at the
 * source ends `404` instead (SPEC-018).
 */
data class DisableTableResponse(
    @SerializedName("removed") val removed: Boolean,
    @SerializedName("retained") val retained: Boolean,
)

/**
 * `GetStatus` response (`200`) — both `false` reads a table that was never
 * enabled. A table physically missing at the source ends `404` instead
 * (SPEC-018).
 */
data class TableStatusResponse(
    @SerializedName("enabled") val enabled: Boolean,
    @SerializedName("retained") val retained: Boolean,
)

/** One entry of a [ListTablesResponse] list (SPEC-018). */
data class TableInfo(
    @SerializedName("table_id") val tableId: String,
    @SerializedName("source") val source: String,
    @SerializedName("schema") val schema: String,
    @SerializedName("table") val table: String,
)

/**
 * `ListTables` response (`200`) — both lists are empty without any
 * activation (SPEC-018).
 */
data class ListTablesResponse(
    @SerializedName("tables") val tables: List<TableInfo>,
    @SerializedName("retained") val retained: List<TableInfo>,
)
