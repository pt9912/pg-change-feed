package io.github.pt9912.pgchangefeed.nats

/**
 * Thrown when a message on the `SPEC-024` NATS full-content stream does not
 * decode to the documented ten-field JSON shape — a protocol violation
 * outside `SPEC-024`'s documented shape, the NATS-surface counterpart of
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedMalformedResponseException]/
 * [io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient]'s equivalent
 * frame-parse failure.
 *
 * Deliberately its own, small exception type rather than a reuse of
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedException]'s hierarchy:
 * that hierarchy's `statusCode` property is an HTTP-only concept
 * (`SPEC-018`/`SPEC-021`'s response status code), which has no NATS
 * equivalent — a NATS message carries no status code, only a subject and a
 * payload. Forcing a placeholder status code onto this type to fit that
 * hierarchy would be a worse fit than a small, standalone type (same
 * reasoning as the C# sibling's
 * `PgChangeFeed.Client.Nats.PgChangeFeedNatsMalformedMessageException`).
 */
class PgChangeFeedNatsMalformedMessageException(
    message: String,
    cause: Throwable? = null,
) : Exception(message, cause)
