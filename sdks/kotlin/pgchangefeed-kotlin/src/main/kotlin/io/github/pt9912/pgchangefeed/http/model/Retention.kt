package io.github.pt9912.pgchangefeed.http.model

import com.google.gson.annotations.SerializedName

/**
 * The retention run to start: deletes the changes of `source` that are older
 * than `minAgeNanos` and that every consumer with a stored position has
 * passed. `source` is mandatory, `minAgeNanos` must be >= 0.
 */
data class RunRetentionRequest(
    @SerializedName("source") val source: String,
    @SerializedName("min_age_nanos") val minAgeNanos: Long,
)

/** The number of changes actually deleted. */
data class RunRetentionResponse(
    @SerializedName("deleted") val deleted: Int,
)
