package io.github.pt9912.pgchangefeed.http

/**
 * Typed errors for the non-success responses of [PgChangeFeedHttpClient]: the
 * HTTP error body (`{"error": "<text>", "code": "<message code>"}`) becomes a
 * typed exception instead of a raw HTTP status code mixed with the success
 * path, consistent across every method.
 *
 * `PgChangeFeedException` is a Kotlin **sealed class**: a caller can `when`
 * over every subtype exhaustively (compiler-checked, no silent fallthrough),
 * while `catch (e: PgChangeFeedException)` still catches them all. The
 * subtypes group the status codes as follows: 400, 401, 403, 404, 500, any
 * other status, and a success response that cannot be read.
 *
 * `messageCode` is the message code of the server (`PCF-<E|W|I><4 digits>`)
 * from the `code` field of the error body, passed through unchanged; it is
 * `null` when the server sent none (authentication errors, a server without
 * message codes, a body that is not JSON, a `code` that is not a non-empty
 * string).
 */
sealed class PgChangeFeedException @JvmOverloads constructor(
    val statusCode: Int,
    message: String,
    cause: Throwable? = null,
    val messageCode: String? = null,
) : Exception(message, cause)

/**
 * `400` — an invalid request body or a violated rule of the API.
 */
class PgChangeFeedBadRequestException @JvmOverloads constructor(
    statusCode: Int,
    message: String,
    messageCode: String? = null,
) : PgChangeFeedException(statusCode, message, messageCode = messageCode)

/**
 * `401` — a missing bearer token, or one that matches no configured token
 * class.
 */
class PgChangeFeedUnauthorizedException @JvmOverloads constructor(
    statusCode: Int,
    message: String,
    messageCode: String? = null,
) : PgChangeFeedException(statusCode, message, messageCode = messageCode)

/**
 * `403` — a known token whose class does not reach the called endpoint (e.g. a
 * `reader` token against an `admin` endpoint).
 */
class PgChangeFeedForbiddenException @JvmOverloads constructor(
    statusCode: Int,
    message: String,
    messageCode: String? = null,
) : PgChangeFeedException(statusCode, message, messageCode = messageCode)

/**
 * `404` — the addressed table does not exist in the source database (only
 * `enableTable`, `disableTable` and `getStatus`).
 */
class PgChangeFeedNotFoundException @JvmOverloads constructor(
    statusCode: Int,
    message: String,
    messageCode: String? = null,
) : PgChangeFeedException(statusCode, message, messageCode = messageCode)

/** `500` — an unexpected internal error of the server. */
class PgChangeFeedServerErrorException @JvmOverloads constructor(
    statusCode: Int,
    message: String,
    messageCode: String? = null,
) : PgChangeFeedException(statusCode, message, messageCode = messageCode)

/**
 * Any other non-success status code than 400, 401, 403, 404 and 500.
 */
class PgChangeFeedUnexpectedStatusException @JvmOverloads constructor(
    statusCode: Int,
    message: String,
    messageCode: String? = null,
) : PgChangeFeedException(statusCode, message, messageCode = messageCode)

/**
 * A success status code (`2xx`) whose body does not parse as the expected
 * response — either invalid JSON or valid-but-empty/mismatched JSON. It is
 * distinct from the status-code exceptions above so a caller can still catch
 * [PgChangeFeedException] uniformly across every method. It carries no
 * message code.
 */
class PgChangeFeedMalformedResponseException(
    statusCode: Int,
    message: String,
    cause: Throwable? = null,
) : PgChangeFeedException(statusCode, message, cause)
