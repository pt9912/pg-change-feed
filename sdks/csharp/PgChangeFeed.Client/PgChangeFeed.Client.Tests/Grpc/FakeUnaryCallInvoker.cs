using Grpc.Core;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>
/// A minimal fake <see cref="CallInvoker"/> for network-free tests of
/// <see cref="PgChangeFeed.Client.Grpc.PgChangeFeedAdministrationClient"/> —
/// it never touches a socket; it hands back a caller-supplied unary response
/// (or a non-OK status) and records the outgoing request and call options
/// (the <c>authorization</c> metadata entry lives in
/// <see cref="CallOptions.Headers"/>) for assertions. Only
/// <see cref="AsyncUnaryCall{TRequest,TResponse}"/> is implemented — the
/// only RPC shape <c>Administration</c> uses; the remaining
/// <see cref="CallInvoker"/> members are unreachable through this client and
/// throw if ever exercised.
/// </summary>
internal sealed class FakeUnaryCallInvoker : CallInvoker
{
    private readonly object? _response;
    private readonly Status? _failureStatus;

    public CallOptions? LastCallOptions { get; private set; }
    public object? LastRequest { get; private set; }

    private FakeUnaryCallInvoker(object? response, Status? failureStatus)
    {
        _response = response;
        _failureStatus = failureStatus;
    }

    /// <summary>Builds a fake invoker whose unary call returns the given response — the happy path.</summary>
    public static FakeUnaryCallInvoker WithResponse<TResponse>(TResponse response) => new(response, failureStatus: null);

    /// <summary>
    /// Builds a fake invoker whose unary call ends with the given non-OK
    /// status — simulating what a real gRPC channel does when a call is
    /// rejected (e.g. the <c>Unauthenticated</c>/<c>PermissionDenied</c>
    /// auth boundary, or a domain error mapped to <c>InvalidArgument</c>/
    /// <c>NotFound</c>/<c>Internal</c>).
    /// </summary>
    public static FakeUnaryCallInvoker WithStatus(Status status) => new(null, status);

    public override AsyncUnaryCall<TResponse> AsyncUnaryCall<TRequest, TResponse>(
        Method<TRequest, TResponse> method, string? host, CallOptions options, TRequest request)
    {
        LastCallOptions = options;
        LastRequest = request;

        if (_failureStatus is { } status && status.StatusCode != StatusCode.OK)
        {
            var failed = Task.FromException<TResponse>(new RpcException(status));
            return new AsyncUnaryCall<TResponse>(
                failed,
                Task.FromResult(new Metadata()),
                () => status,
                () => new Metadata(),
                () => { });
        }

        var response = (TResponse)_response!;
        return new AsyncUnaryCall<TResponse>(
            Task.FromResult(response),
            Task.FromResult(new Metadata()),
            () => new Status(StatusCode.OK, string.Empty),
            () => new Metadata(),
            () => { });
    }

    public override TResponse BlockingUnaryCall<TRequest, TResponse>(
        Method<TRequest, TResponse> method, string? host, CallOptions options, TRequest request)
        => throw new NotSupportedException("Administration uses async unary calls only.");

    public override AsyncServerStreamingCall<TResponse> AsyncServerStreamingCall<TRequest, TResponse>(
        Method<TRequest, TResponse> method, string? host, CallOptions options, TRequest request)
        => throw new NotSupportedException("Administration uses async unary calls only.");

    public override AsyncClientStreamingCall<TRequest, TResponse> AsyncClientStreamingCall<TRequest, TResponse>(
        Method<TRequest, TResponse> method, string? host, CallOptions options)
        => throw new NotSupportedException("Administration uses async unary calls only.");

    public override AsyncDuplexStreamingCall<TRequest, TResponse> AsyncDuplexStreamingCall<TRequest, TResponse>(
        Method<TRequest, TResponse> method, string? host, CallOptions options)
        => throw new NotSupportedException("Administration uses async unary calls only.");
}
