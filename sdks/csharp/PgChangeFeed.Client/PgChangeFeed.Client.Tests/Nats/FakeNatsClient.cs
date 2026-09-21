using System.Runtime.CompilerServices;
using NATS.Client.Core;

namespace PgChangeFeed.Client.Tests.Nats;

/// <summary>
/// A minimal fake <see cref="INatsClient"/> for network-free tests of
/// <see cref="PgChangeFeed.Client.Nats.PgChangeFeedNatsStreamClient"/> — it
/// never touches a socket; it hands back a caller-supplied sequence of raw
/// message payloads (or a failure) and records the subscribed subject for
/// assertions. Only <see cref="SubscribeAsync{T}"/> is implemented — the
/// only <see cref="INatsClient"/> member <c>PgChangeFeedNatsStreamClient</c>
/// uses; the remaining members are unreachable through it and throw if ever
/// exercised. Muster: <c>PgChangeFeed.Client.Tests.Grpc.FakeCallInvoker</c>.
/// </summary>
internal sealed class FakeNatsClient : INatsClient
{
    private readonly IReadOnlyList<byte[]> _payloads;
    private readonly Exception? _failure;

    /// <summary>The subject the client last subscribed to — asserted against the subject the production client built/passed.</summary>
    public string? LastSubject { get; private set; }

    private FakeNatsClient(IReadOnlyList<byte[]> payloads, Exception? failure)
    {
        _payloads = payloads;
        _failure = failure;
    }

    /// <summary>
    /// Builds a fake client whose subscription yields the given raw payloads
    /// in order, then completes the enumeration cleanly — the happy path
    /// (SPEC-024: one message per row change, in order).
    /// </summary>
    public static FakeNatsClient WithPayloads(params byte[][] payloads) => new(payloads, failure: null);

    /// <summary>
    /// Builds a fake client whose subscription yields no message and then
    /// throws <paramref name="failure"/> from the enumeration itself —
    /// simulating what a real NATS connection does when the server rejects
    /// it (SPEC-024's connection-level auth boundary: a missing or wrong
    /// <c>CDC_NATS_STREAM_TOKEN</c>).
    /// </summary>
    public static FakeNatsClient WithFailure(Exception failure) => new([], failure);

    public IAsyncEnumerable<NatsMsg<T>> SubscribeAsync<T>(
        string subject,
        string? queueGroup = null,
        INatsDeserialize<T>? serializer = null,
        NatsSubOpts? opts = null,
        CancellationToken cancellationToken = default)
    {
        LastSubject = subject;
        return EnumerateAsync<T>(subject, cancellationToken);
    }

    private async IAsyncEnumerable<NatsMsg<T>> EnumerateAsync<T>(
        string subject,
        [EnumeratorCancellation] CancellationToken cancellationToken)
    {
        foreach (var payload in _payloads)
        {
            cancellationToken.ThrowIfCancellationRequested();
            await Task.Yield();
            var data = (T)(object)payload;
            yield return new NatsMsg<T>(subject, null, payload.Length, null, data, null);
        }

        if (_failure is not null)
        {
            throw _failure;
        }
    }

    public INatsConnection Connection
        => throw new NotSupportedException("Not used by PgChangeFeedNatsStreamClient.");

    public ValueTask ConnectAsync()
        => throw new NotSupportedException("Not used by PgChangeFeedNatsStreamClient.");

    public ValueTask<TimeSpan> PingAsync(CancellationToken cancellationToken = default)
        => throw new NotSupportedException("Not used by PgChangeFeedNatsStreamClient.");

    public ValueTask PublishAsync<T>(
        string subject, T data, NatsHeaders? headers = null, string? replyTo = null,
        INatsSerialize<T>? serializer = null, NatsPubOpts? opts = null, CancellationToken cancellationToken = default)
        => throw new NotSupportedException("Not used by PgChangeFeedNatsStreamClient.");

    public ValueTask PublishAsync(
        string subject, NatsHeaders? headers = null, string? replyTo = null,
        NatsPubOpts? opts = null, CancellationToken cancellationToken = default)
        => throw new NotSupportedException("Not used by PgChangeFeedNatsStreamClient.");

    public ValueTask<NatsMsg<TReply>> RequestAsync<TRequest, TReply>(
        string subject, TRequest? data, NatsHeaders? headers = null,
        INatsSerialize<TRequest>? requestSerializer = null, INatsDeserialize<TReply>? replySerializer = null,
        NatsPubOpts? requestOpts = null, NatsSubOpts? replyOpts = null, CancellationToken cancellationToken = default)
        => throw new NotSupportedException("Not used by PgChangeFeedNatsStreamClient.");

    public ValueTask<NatsMsg<TReply>> RequestAsync<TReply>(
        string subject, INatsDeserialize<TReply>? replySerializer = null,
        NatsSubOpts? replyOpts = null, CancellationToken cancellationToken = default)
        => throw new NotSupportedException("Not used by PgChangeFeedNatsStreamClient.");

    public ValueTask ReconnectAsync()
        => throw new NotSupportedException("Not used by PgChangeFeedNatsStreamClient.");

    public ValueTask DisposeAsync() => ValueTask.CompletedTask;
}
