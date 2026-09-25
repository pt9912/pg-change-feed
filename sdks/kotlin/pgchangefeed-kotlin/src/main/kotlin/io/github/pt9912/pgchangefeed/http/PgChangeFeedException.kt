package io.github.pt9912.pgchangefeed.http

/**
 * Typed errors for the non-success responses of [PgChangeFeedHttpClient]: the
 * HTTP error body (`{"error": "<text>"}`) becomes a typed exception instead of
 * a raw HTTP status code mixed with the success path, consistent across every
 * method.
 *
 * `PgChangeFeedException` is a Kotlin **sealed class**: a caller can `when`
 * over every subtype exhaustively (compiler-checked, no silent fallthrough),
 * while `catch (e: PgChangeFeedException)` still catches them all. The
 * subtypes group the status codes as follows: 400, 401, 403, 404, 500, any
 * other status, and a success response that cannot be read.
 */
sealed class PgChangeFeedException(
    val statusCode: Int,
    message: String,
    cause: Throwable? = null,
) : Exception(message, cause)

/**
 * `400` — an invalid request body or a violated rule of the API.
 */
class PgChangeFeedBadRequestException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/**
 * `401` — a missing bearer token, or one that matches no configured token
 * class.
 */
class PgChangeFeedUnauthorizedException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/**
 * `403` — a known token whose class does not reach the called endpoint (e.g. a
 * `reader` token against an `admin` endpoint).
 */
class PgChangeFeedForbiddenException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/**
 * `404` — the addressed table does not exist in the source database (only
 * `enableTable`, `disableTable` and `getStatus`).
 */
class PgChangeFeedNotFoundException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/** `500` — an unexpected internal error of the server. */
class PgChangeFeedServerErrorException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/**
 * Any other non-success status code than 400, 401, 403, 404 and 500.
 */
class PgChangeFeedUnexpectedStatusException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/**
 * A success status code (`2xx`) whose body does not parse as the expected
 * response — either invalid JSON or valid-but-empty/mismatched JSON. It is
 * distinct from the status-code exceptions above so a caller can still catch
 * [PgChangeFeedException] uniformly across every method.
 */
class PgChangeFeedMalformedResponseException(
    statusCode: Int,
    message: String,
    cause: Throwable? = null,
) : PgChangeFeedException(statusCode, message, cause)
