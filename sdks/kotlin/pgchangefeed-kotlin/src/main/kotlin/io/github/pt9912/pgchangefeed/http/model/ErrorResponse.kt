package io.github.pt9912.pgchangefeed.http.model

import com.google.gson.annotations.SerializedName

/**
 * The uniform error body of every non-success response
 * (`{"error": "<text>"}`, SPEC-018) — internal, because it is only ever
 * used to build a typed `PgChangeFeedException`; it is never a return
 * value of a public method.
 */
internal data class ErrorResponse(
    @SerializedName("error") val error: String,
)
