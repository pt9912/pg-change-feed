package io.github.pt9912.pgchangefeed.http.model

import com.google.gson.annotations.SerializedName

/**
 * `RunRetention` request — `POST /retention/run` (SPEC-018); `source`
 * mandatory, `minAgeNanos` must be >= 0.
 */
data class RunRetentionRequest(
    @SerializedName("source") val source: String,
    @SerializedName("min_age_nanos") val minAgeNanos: Long,
)

/** `RunRetention` response (`200`) — the number of changes actually deleted (SPEC-018). */
data class RunRetentionResponse(
    @SerializedName("deleted") val deleted: Int,
)
