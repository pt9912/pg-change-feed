namespace PgChangeFeed.Client.Nats;

/// <summary>
/// Thrown when a message on the SPEC-024 NATS full-content stream does not
/// decode to the documented ten-field JSON shape — a protocol violation
/// outside SPEC-024's documented shape, the NATS-surface counterpart of
/// <c>PgChangeFeed.Client.Http.PgChangeFeedMalformedResponseException</c>.
///
/// Deliberately its own, small exception type rather than a reuse of
/// <c>PgChangeFeed.Client.Http.PgChangeFeedException</c>'s hierarchy: that
/// hierarchy's <c>StatusCode</c> property is an HTTP-only concept (SPEC-018/
/// SPEC-021's response status code), which has no NATS equivalent — a NATS
/// message carries no status code, only a subject and a payload. Forcing a
/// placeholder status code onto this type to fit that hierarchy would be a
/// worse fit than a small, standalone type.
/// </summary>
public sealed class PgChangeFeedNatsMalformedMessageException : Exception
{
    public PgChangeFeedNatsMalformedMessageException(string message)
        : base(message)
    {
    }

    public PgChangeFeedNatsMalformedMessageException(string message, Exception innerException)
        : base(message, innerException)
    {
    }
}
