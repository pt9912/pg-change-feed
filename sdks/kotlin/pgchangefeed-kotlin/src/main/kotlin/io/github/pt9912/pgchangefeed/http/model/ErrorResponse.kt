package io.github.pt9912.pgchangefeed.http.model

import com.google.gson.annotations.SerializedName

/**
 * The error body of every non-success response (`{"error": "<text>"}`) —
 * internal, because it is only ever used to build a typed
 * `PgChangeFeedException`; it is never a return value of a public method.
 */
internal data class ErrorResponse(
    @SerializedName("error") val error: String,
)
