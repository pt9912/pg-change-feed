namespace PgChangeFeed.Client.Http;

/// <summary>
/// Base type for every typed error <see cref="PgChangeFeedHttpClient"/>
/// throws for a non-success response. The HTTP error body
/// (<c>{"error": "&lt;text&gt;"}</c>) becomes a typed exception instead of a
/// result type mixed with the success path, consistent across every method
/// on the client.
/// </summary>
public abstract class PgChangeFeedException : Exception
{
    /// <summary>The HTTP status code the server returned.</summary>
    public int StatusCode { get; }

    /// <summary>Creates the exception with the HTTP status code and the error text.</summary>
    protected PgChangeFeedException(int statusCode, string message)
        : base(message)
    {
        StatusCode = statusCode;
    }

    /// <summary>Creates the exception with the HTTP status code, the error text and the cause.</summary>
    protected PgChangeFeedException(int statusCode, string message, Exception innerException)
        : base(message, innerException)
    {
        StatusCode = statusCode;
    }
}

/// <summary>
/// <c>400</c> — an invalid request body or a violated rule of the API.
/// </summary>
public sealed class PgChangeFeedBadRequestException : PgChangeFeedException
{
    /// <summary>Creates the exception with the HTTP status code and the error text.</summary>
    public PgChangeFeedBadRequestException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}

/// <summary>
/// <c>401</c> — a missing bearer token, or one that matches no configured
/// token class.
/// </summary>
public sealed class PgChangeFeedUnauthorizedException : PgChangeFeedException
{
    /// <summary>Creates the exception with the HTTP status code and the error text.</summary>
    public PgChangeFeedUnauthorizedException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}

/// <summary>
/// <c>403</c> — a known token whose class does not reach the called endpoint
/// (e.g. a <c>reader</c> token against an <c>admin</c> endpoint).
/// </summary>
public sealed class PgChangeFeedForbiddenException : PgChangeFeedException
{
    /// <summary>Creates the exception with the HTTP status code and the error text.</summary>
    public PgChangeFeedForbiddenException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}

/// <summary>
/// <c>404</c> — the addressed table does not exist in the source database
/// (only <c>EnableTableAsync</c>, <c>DisableTableAsync</c> and
/// <c>GetStatusAsync</c>).
/// </summary>
public sealed class PgChangeFeedNotFoundException : PgChangeFeedException
{
    /// <summary>Creates the exception with the HTTP status code and the error text.</summary>
    public PgChangeFeedNotFoundException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}

/// <summary><c>500</c> — an unexpected internal error of the server.</summary>
public sealed class PgChangeFeedServerErrorException : PgChangeFeedException
{
    /// <summary>Creates the exception with the HTTP status code and the error text.</summary>
    public PgChangeFeedServerErrorException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}

/// <summary>
/// Any other non-success status code than 400, 401, 403, 404 and 500.
/// </summary>
public sealed class PgChangeFeedUnexpectedStatusException : PgChangeFeedException
{
    /// <summary>Creates the exception with the HTTP status code and the error text.</summary>
    public PgChangeFeedUnexpectedStatusException(int statusCode, string message)
        : base(statusCode, message)
    {
    }
}

/// <summary>
/// A success status code (<c>2xx</c>) whose body does not parse as the
/// expected response — either invalid JSON or a valid-but-empty/<c>null</c>
/// body. It is distinct from the status-code exceptions above so a caller can
/// still catch <see cref="PgChangeFeedException"/> uniformly across every
/// method.
/// </summary>
public sealed class PgChangeFeedMalformedResponseException : PgChangeFeedException
{
    /// <summary>Creates the exception with the HTTP status code and a description of the problem.</summary>
    public PgChangeFeedMalformedResponseException(int statusCode, string message)
        : base(statusCode, message)
    {
    }

    /// <summary>Creates the exception with the HTTP status code, a description of the problem and the cause.</summary>
    public PgChangeFeedMalformedResponseException(int statusCode, string message, Exception innerException)
        : base(statusCode, message, innerException)
    {
    }
}
