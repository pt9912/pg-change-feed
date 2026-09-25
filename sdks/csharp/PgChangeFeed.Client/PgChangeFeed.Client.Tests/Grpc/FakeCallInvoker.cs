using Grpc.Core;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>
/// A minimal fake <see cref="CallInvoker"/> for network-free tests of
/// <see cref="PgChangeFeed.Client.Grpc.PgChangeFeedGrpcClient"/> — it never
/// touches a socket; it hands back a caller-supplied server-streaming result
/// and records the outgoing call options (the <c>authorization</c> metadata
/// entry lives in <see cref="CallOptions.Headers"/>) for assertions. Only
/// <see cref="AsyncServerStreamingCall{TRequest,TResponse}"/> is
/// implemented — the only RPC shape <c>ChangeStream</c> uses; the
/// remaining <see cref="CallInvoker"/> members are unreachable through this
/// client and throw if ever exercised.
/// </summary>
internal sealed class FakeCallInvoker : CallInvoker
{
    private readonly object[] _messages;
    private readonly Status? _failureStatus;

    public CallOptions? LastCallOptions { get; private set; }

    private FakeCallInvoker(object[] messages, Status? failureStatus)
    {
        _messages = messages;
        _failureStatus = failureStatus;
    }

    /// <summary>
    /// Builds a fake invoker whose server-streaming call yields the given
    /// messages in order, then completes the stream cleanly (Status.OK) —
    /// the happy path (one message per row change, in order).
    /// </summary>
    public static FakeCallInvoker WithMessages<TResponse>(params TResponse[] messages)
        => new(messages.Cast<object>().ToArray(), failureStatus: null);

    /// <summary>
    /// Builds a fake invoker whose server-streaming call ends immediately
    /// with the given non-OK status — simulating what a real gRPC channel
    /// does when a call is rejected before any message is sent (e.g. the
    /// <c>Unauthenticated</c> auth boundary).
    /// </summary>
    public static FakeCallInvoker WithStatus(Status status)
        => new([], status);

    public override AsyncServerStreamingCall<TResponse> AsyncServerStreamingCall<TRequest, TResponse>(
        Method<TRequest, TResponse> method, string? host, CallOptions options, TRequest request)
    {
        LastCallOptions = options;
        var typedMessages = _messages.Cast<TResponse>().ToArray();
        var reader = new FakeAsyncStreamReader<TResponse>(typedMessages, _failureStatus);
        var finalStatus = _failureStatus ?? new Status(StatusCode.OK, string.Empty);
        return new AsyncServerStreamingCall<TResponse>(
            reader,
            Task.FromResult(new Metadata()),
            () => finalStatus,
            () => new Metadata(),
            () => { });
    }

    public override TResponse BlockingUnaryCall<TRequest, TResponse>(
        Method<TRequest, TResponse> method, string? host, CallOptions options, TRequest request)
        => throw new NotSupportedException("ChangeStream uses server streaming only.");

    public override AsyncUnaryCall<TResponse> AsyncUnaryCall<TRequest, TResponse>(
        Method<TRequest, TResponse> method, string? host, CallOptions options, TRequest request)
        => throw new NotSupportedException("ChangeStream uses server streaming only.");

    public override AsyncClientStreamingCall<TRequest, TResponse> AsyncClientStreamingCall<TRequest, TResponse>(
        Method<TRequest, TResponse> method, string? host, CallOptions options)
        => throw new NotSupportedException("ChangeStream uses server streaming only.");

    public override AsyncDuplexStreamingCall<TRequest, TResponse> AsyncDuplexStreamingCall<TRequest, TResponse>(
        Method<TRequest, TResponse> method, string? host, CallOptions options)
        => throw new NotSupportedException("ChangeStream uses server streaming only.");
}

/// <summary>
/// Yields <paramref name="messages"/> in order, then either throws an
/// <see cref="RpcException"/> carrying <paramref name="failureStatus"/> (if
/// it is a non-OK status — the auth boundary) or ends the stream cleanly
/// (the happy path). This mirrors what a real <c>Grpc.Net.Client</c> stream
/// reader does: a non-OK trailing status surfaces as an exception from
/// <see cref="MoveNext"/>, not as a swallowed empty stream.
/// </summary>
internal sealed class FakeAsyncStreamReader<TResponse> : IAsyncStreamReader<TResponse>
{
    private readonly IReadOnlyList<TResponse> _messages;
    private readonly Status? _failureStatus;
    private int _index = -1;

    public FakeAsyncStreamReader(IReadOnlyList<TResponse> messages, Status? failureStatus)
    {
        _messages = messages;
        _failureStatus = failureStatus;
    }

    public TResponse Current
        => _index >= 0 && _index < _messages.Count
            ? _messages[_index]
            : throw new InvalidOperationException("No current message — call MoveNext first.");

    public Task<bool> MoveNext(CancellationToken cancellationToken)
    {
        _index++;
        if (_index < _messages.Count)
        {
            return Task.FromResult(true);
        }
        if (_failureStatus is { } status && status.StatusCode != StatusCode.OK)
        {
            throw new RpcException(status);
        }
        return Task.FromResult(false);
    }
}
