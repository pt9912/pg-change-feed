package io.github.pt9912.pgchangefeed.http

/**
 * Typed error hierarchy for [PgChangeFeedHttpClient]'s non-success
 * responses — mirrors the uniform SPEC-018 error body
 * (`{"error": "<text>"}`) as a typed exception instead of a raw HTTP
 * status code mixed with the success path, consistent across every
 * method.
 *
 * `PgChangeFeedException` is a Kotlin **sealed class**, not a plain open
 * exception hierarchy like the C#/Python siblings
 * (`PgChangeFeed.Client.Http.PgChangeFeedException`,
 * `pgchangefeed.exceptions.PgChangeFeedError`). Field-for-field the seven
 * subtypes below mirror those two exactly (same status-code grouping, same
 * `statusCode` property, same base-type catch-all); the sealed modifier is
 * additive Kotlin idiom, not a behavior change — it lets a caller `when`
 * over every subtype exhaustively (compiler-checked, no silent fallthrough)
 * while `catch (e: PgChangeFeedException)` still works uniformly, exactly
 * like the C#/Python base class. `ADR-0109` Festlegung 4 binds the SemVer
 * major boundary to wire changes, not to this internal design choice — a
 * redesign of this hierarchy stays possible before `1.0.0` without a major
 * bump.
 */
sealed class PgChangeFeedException(
    val statusCode: Int,
    message: String,
    cause: Throwable? = null,
) : Exception(message, cause)

/**
 * `400` — an invalid request body or a violated domain invariant
 * (SPEC-018/SPEC-022).
 */
class PgChangeFeedBadRequestException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/**
 * `401` — a missing bearer token, or one that matches no configured token
 * class (SPEC-018).
 */
class PgChangeFeedUnauthorizedException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/**
 * `403` — a known token whose rights class does not reach the called
 * endpoint (e.g. a `reader` token against an `admin` endpoint, SPEC-018).
 */
class PgChangeFeedForbiddenException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/**
 * `404` — the addressed table is physically missing at the source (only
 * `EnableTable`/`DisableTable`/`GetStatus`, SPEC-018).
 */
class PgChangeFeedNotFoundException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/** `500` — an unexpected internal server error (SPEC-018/SPEC-022). */
class PgChangeFeedServerErrorException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/**
 * Any non-success status code outside the five SPEC-018/SPEC-022 document
 * (400/401/403/404/500) — a defensive fallback that is itself not part of
 * the documented wire contract.
 */
class PgChangeFeedUnexpectedStatusException(statusCode: Int, message: String) :
    PgChangeFeedException(statusCode, message)

/**
 * A success status code (`2xx`) whose body does not parse as the expected
 * response shape — either invalid JSON or valid-but-empty/mismatched JSON.
 * Outside every shape SPEC-018/SPEC-022 document; kept typed and distinct
 * from the status-code exceptions above so a caller can still catch
 * [PgChangeFeedException] uniformly across every method.
 */
class PgChangeFeedMalformedResponseException(
    statusCode: Int,
    message: String,
    cause: Throwable? = null,
) : PgChangeFeedException(statusCode, message, cause)
