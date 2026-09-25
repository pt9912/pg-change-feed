namespace PgChangeFeed.Client.Nats;

/// <summary>
/// Thrown when a message on the NATS stream does not decode to the ten-field
/// JSON change — the NATS counterpart of
/// <c>PgChangeFeed.Client.Http.PgChangeFeedMalformedResponseException</c>.
///
/// It is its own small exception type, not part of
/// <c>PgChangeFeed.Client.Http.PgChangeFeedException</c>'s hierarchy: that
/// hierarchy's <c>StatusCode</c> is an HTTP status, and a NATS message carries
/// no status code, only a subject and a payload.
/// </summary>
public sealed class PgChangeFeedNatsMalformedMessageException : Exception
{
    /// <summary>Creates the exception with a description of the problem.</summary>
    public PgChangeFeedNatsMalformedMessageException(string message)
        : base(message)
    {
    }

    /// <summary>Creates the exception with a description of the problem and the cause.</summary>
    public PgChangeFeedNatsMalformedMessageException(string message, Exception innerException)
        : base(message, innerException)
    {
    }
}
