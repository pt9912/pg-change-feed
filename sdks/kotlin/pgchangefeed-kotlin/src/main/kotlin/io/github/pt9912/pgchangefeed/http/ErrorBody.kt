package io.github.pt9912.pgchangefeed.http

import com.google.gson.JsonElement
import com.google.gson.JsonParseException
import com.google.gson.JsonParser

/** The error text and the message code read out of an error body. */
internal data class ParsedError(val message: String, val messageCode: String?)

/**
 * Reads the error text (the raw body when `error` is not a JSON string) and the
 * message code of a non-success response body; shared by the HTTP and the SSE
 * client. Only a JSON string counts for either field — Gson's lenient coercion
 * of a number to a string is deliberately not used — and a non-empty string in
 * `code` is a code whatever the type of `error`. Never throws.
 */
internal fun parseErrorBody(body: String): ParsedError =
    try {
        val root = JsonParser.parseString(body)
        if (root.isJsonObject) {
            val fields = root.asJsonObject
            ParsedError(stringOf(fields.get("error")) ?: body, stringOf(fields.get("code"))?.ifEmpty { null })
        } else {
            ParsedError(body, null)
        }
    } catch (ex: JsonParseException) {
        ParsedError(body, null)
    }

private fun stringOf(element: JsonElement?): String? =
    if (element != null && element.isJsonPrimitive && element.asJsonPrimitive.isString) element.asString else null
