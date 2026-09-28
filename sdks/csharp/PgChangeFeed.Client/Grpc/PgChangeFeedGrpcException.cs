using Grpc.Core;

namespace PgChangeFeed.Client.Grpc;

/// <summary>
/// Base type for every typed error <see cref="PgChangeFeedAdministrationClient"/>
/// throws for a non-<see cref="StatusCode.OK"/> gRPC status. The status
/// detail text becomes a typed exception instead of a result type mixed
/// with the success path, the same design as
/// <see cref="PgChangeFeed.Client.Http.PgChangeFeedException"/> for the HTTP
/// surface. The original <see cref="RpcException"/> is always the
/// <see cref="Exception.InnerException"/>.
/// </summary>
public abstract class PgChangeFeedGrpcException : Exception
{
    /// <summary>The gRPC status code the server returned.</summary>
    public StatusCode StatusCode { get; }

    /// <summary>Creates the exception with the gRPC status code, the detail text and the causing <see cref="RpcException"/>.</summary>
    protected PgChangeFeedGrpcException(StatusCode statusCode, string message, RpcException innerException)
        : base(message, innerException)
    {
        StatusCode = statusCode;
    }
}

/// <summary>
/// <c>InvalidArgument</c> — a violated rule of the request (e.g. an empty
/// identifier, an inverted range, a non-positive limit).
/// </summary>
public sealed class PgChangeFeedGrpcInvalidArgumentException : PgChangeFeedGrpcException
{
    /// <summary>Creates the exception with the detail text and the causing <see cref="RpcException"/>.</summary>
    public PgChangeFeedGrpcInvalidArgumentException(string message, RpcException innerException)
        : base(StatusCode.InvalidArgument, message, innerException)
    {
    }
}

/// <summary>
/// <c>Unauthenticated</c> — a missing bearer token, or one that matches no
/// configured token class.
/// </summary>
public sealed class PgChangeFeedGrpcUnauthenticatedException : PgChangeFeedGrpcException
{
    /// <summary>Creates the exception with the detail text and the causing <see cref="RpcException"/>.</summary>
    public PgChangeFeedGrpcUnauthenticatedException(string message, RpcException innerException)
        : base(StatusCode.Unauthenticated, message, innerException)
    {
    }
}

/// <summary>
/// <c>PermissionDenied</c> — a known token whose class does not reach the
/// called RPC (e.g. a <c>reader</c> token against an <c>admin</c> RPC).
/// </summary>
public sealed class PgChangeFeedGrpcPermissionDeniedException : PgChangeFeedGrpcException
{
    /// <summary>Creates the exception with the detail text and the causing <see cref="RpcException"/>.</summary>
    public PgChangeFeedGrpcPermissionDeniedException(string message, RpcException innerException)
        : base(StatusCode.PermissionDenied, message, innerException)
    {
    }
}

/// <summary>
/// <c>NotFound</c> — the addressed table does not exist in the source
/// database (<c>EnableTableAsync</c>, <c>DisableTableAsync</c>,
/// <c>GetTableStatusAsync</c>).
/// </summary>
public sealed class PgChangeFeedGrpcNotFoundException : PgChangeFeedGrpcException
{
    /// <summary>Creates the exception with the detail text and the causing <see cref="RpcException"/>.</summary>
    public PgChangeFeedGrpcNotFoundException(string message, RpcException innerException)
        : base(StatusCode.NotFound, message, innerException)
    {
    }
}

/// <summary><c>Internal</c> — an unexpected internal error of the server.</summary>
public sealed class PgChangeFeedGrpcInternalException : PgChangeFeedGrpcException
{
    /// <summary>Creates the exception with the detail text and the causing <see cref="RpcException"/>.</summary>
    public PgChangeFeedGrpcInternalException(string message, RpcException innerException)
        : base(StatusCode.Internal, message, innerException)
    {
    }
}

/// <summary>
/// Any other non-<see cref="StatusCode.OK"/> gRPC status than
/// <see cref="StatusCode.InvalidArgument"/>, <see cref="StatusCode.Unauthenticated"/>,
/// <see cref="StatusCode.PermissionDenied"/>, <see cref="StatusCode.NotFound"/>
/// and <see cref="StatusCode.Internal"/>. It is distinct from the five
/// status-specific exceptions above so a caller can still catch
/// <see cref="PgChangeFeedGrpcException"/> uniformly across every method.
/// </summary>
public sealed class PgChangeFeedGrpcUnexpectedStatusException : PgChangeFeedGrpcException
{
    /// <summary>Creates the exception with the gRPC status code, the detail text and the causing <see cref="RpcException"/>.</summary>
    public PgChangeFeedGrpcUnexpectedStatusException(StatusCode statusCode, string message, RpcException innerException)
        : base(statusCode, message, innerException)
    {
    }
}
