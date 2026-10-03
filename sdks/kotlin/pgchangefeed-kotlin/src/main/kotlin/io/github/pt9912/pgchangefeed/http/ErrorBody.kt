package io.github.pt9912.pgchangefeed.http

import com.google.gson.Gson
import com.google.gson.JsonParseException
import com.google.gson.JsonParser
import com.google.gson.JsonSyntaxException
import io.github.pt9912.pgchangefeed.http.model.ErrorResponse

/** The error text and the message code read out of an error body. */
internal data class ParsedError(val message: String, val messageCode: String?)

/**
 * Reads the error text (the raw body when it carries no `error` field) and the
 * message code of a non-success response body; shared by the HTTP and the SSE
 * client. Only a non-empty JSON string in `code` counts as a code — Gson's
 * lenient coercion of a number to a string is deliberately not used for it.
 * Never throws.
 */
internal fun parseErrorBody(gson: Gson, body: String): ParsedError {
    val message = try {
        gson.fromJson(body, ErrorResponse::class.java)?.error ?: body
    } catch (ex: JsonSyntaxException) {
        body
    }
    return ParsedError(message, messageCodeOf(body))
}

private fun messageCodeOf(body: String): String? =
    try {
        val root = JsonParser.parseString(body)
        val code = if (root.isJsonObject) root.asJsonObject.get("code") else null
        if (code != null && code.isJsonPrimitive && code.asJsonPrimitive.isString) {
            code.asString.ifEmpty { null }
        } else {
            null
        }
    } catch (ex: JsonParseException) {
        null
    }
