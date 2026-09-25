package io.github.pt9912.pgchangefeed.nats

/**
 * Thrown when a message on the NATS stream does not decode to the ten-field
 * JSON change — the NATS counterpart of
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedMalformedResponseException].
 *
 * It is its own small exception type, not part of
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedException]'s hierarchy: that
 * hierarchy's `statusCode` is an HTTP status, and a NATS message carries no
 * status code, only a subject and a payload.
 */
class PgChangeFeedNatsMalformedMessageException(
    message: String,
    cause: Throwable? = null,
) : Exception(message, cause)
