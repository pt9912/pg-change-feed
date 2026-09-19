namespace PgChangeFeed.Client.Http;

/// <summary>
/// Base type for every typed error <see cref="PgChangeFeedHttpClient"/>
/// throws for a non-success response — the uniform SPEC-018 error body
/// (<c>{"error": "&lt;text&gt;"}</c>) becomes a typed exception instead of a
/// result type mixed with the success path, consistent across every method
/// on the client.
/// </summary>
public abstract class PgChangeFeedException : Exception
{
    /// <summary>The HTTP status code the server returned.</summary>
    public int StatusCode { get; }

    protected PgChangeFeedException(int statusCode, string message)
        : base(message)
    {
        StatusCode = statusCode;
    }
}

/// <summary>
/// <c>400</c> — an invalid request body or a violated domain invariant
/// (SPEC-018).
/// </summary>
public sealed class PgChangeFeedBadRequestException : PgChangeFeedException
{
    public PgChangeFeedBadRequestException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}

/// <summary>
/// <c>401</c> — a missing bearer token, or one that matches no configured
/// token class (SPEC-018).
/// </summary>
public sealed class PgChangeFeedUnauthorizedException : PgChangeFeedException
{
    public PgChangeFeedUnauthorizedException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}

/// <summary>
/// <c>403</c> — a known token whose rights class does not reach the called
/// endpoint (e.g. a <c>reader</c> token against an <c>admin</c> endpoint,
/// SPEC-018).
/// </summary>
public sealed class PgChangeFeedForbiddenException : PgChangeFeedException
{
    public PgChangeFeedForbiddenException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}

/// <summary>
/// <c>404</c> — the addressed table is physically missing at the source
/// (only <c>EnableTable</c>/<c>DisableTable</c>/<c>GetStatus</c>, SPEC-018).
/// </summary>
public sealed class PgChangeFeedNotFoundException : PgChangeFeedException
{
    public PgChangeFeedNotFoundException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}

/// <summary><c>500</c> — an unexpected internal server error (SPEC-018).</summary>
public sealed class PgChangeFeedServerErrorException : PgChangeFeedException
{
    public PgChangeFeedServerErrorException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}

/// <summary>
/// Any non-success status code outside the five SPEC-018/SPEC-022 document
/// (400/401/403/404/500) — a defensive fallback that is itself not part of
/// the documented wire contract.
/// </summary>
public sealed class PgChangeFeedUnexpectedStatusException : PgChangeFeedException
{
    public PgChangeFeedUnexpectedStatusException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}
