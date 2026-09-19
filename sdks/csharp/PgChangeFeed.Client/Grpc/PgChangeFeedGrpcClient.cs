using System.Runtime.CompilerServices;
using Cdc.Stream.V1;
using Grpc.Core;
using Grpc.Net.Client;

namespace PgChangeFeed.Client.Grpc;

/// <summary>
/// Public entry point for the PG Change Feed live-change stream
/// (<c>SPEC-020</c>, <c>LH-FA-SST-008</c>): one method opens the
/// <c>ChangeStream/StreamChanges</c> server-streaming RPC and yields the
/// generated <see cref="Change"/> message — the ten SPEC-020 fields
/// (<c>change_id</c>, <c>transaction_id</c>, <c>source_table_id</c>,
/// <c>sequence</c>, <c>operation</c>, <c>old_image</c>, <c>new_image</c>,
/// <c>schema_version</c>, <c>schema</c>, <c>table</c>) — unmapped, exactly as
/// the wire defines them. There is no separate DTO layer here (unlike
/// <c>PgChangeFeed.Client.Http.Models</c>): the generated stub already is a
/// typed, versioned representation of the wire schema; introducing a second,
/// hand-mapped type would only risk drifting from it.
///
/// The bearer token is sent in the <c>authorization</c> gRPC metadata entry
/// as <c>Bearer &lt;token&gt;</c> (SPEC-020) on every call — there is no
/// global or static state; a process can hold several independently
/// configured instances at once.
///
/// <b>Boundary (SPEC-020, LH-FA-SST-008):</b> the stream carries no replay
/// and no table-granular filtering. A consumer that needs either uses the
/// existing read path (<c>PgChangeFeed.Client.Http.PgChangeFeedHttpClient.ReadChangesAsync</c>),
/// not this stream.
///
/// Draht-Kenntnis-Vorbild (gelesen, nicht importiert — ADR-0106 Festlegung 2):
/// <c>examples/csharp/grpc-client/Program.cs</c>.
/// </summary>
public sealed class PgChangeFeedGrpcClient : IDisposable
{
    private const string AuthorizationMetadataKey = "authorization";
    private const string BearerPrefix = "Bearer ";

    private readonly ChangeStream.ChangeStreamClient _client;
    private readonly PgChangeFeedClientOptions _options;
    private readonly GrpcChannel? _ownedChannel;

    /// <summary>
    /// Convenience constructor: opens and owns its own <see cref="GrpcChannel"/>
    /// against <paramref name="options"/>'s address. <see cref="Dispose"/>
    /// disposes that channel.
    ///
    /// The PG Change Feed server speaks plaintext gRPC over HTTP/2 (h2c, no
    /// TLS) by default (same choice as the reference clients). Because
    /// enabling h2c on <see cref="System.Net.Http.SocketsHttpHandler"/> is a
    /// process-wide <see cref="AppContext"/> switch
    /// (<c>System.Net.Http.SocketsHttpHandler.Http2UnencryptedSupport</c>),
    /// this constructor deliberately does not set it as a hidden side
    /// effect of constructing a client instance — the caller enables it once,
    /// the same way <c>examples/csharp/grpc-client/Program.cs</c> does,
    /// before connecting to a plaintext server. A TLS-terminated deployment
    /// needs no such switch.
    /// </summary>
    public PgChangeFeedGrpcClient(PgChangeFeedClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(options);
        _options = options;
        _ownedChannel = GrpcChannel.ForAddress(options.Address);
        _client = new ChangeStream.ChangeStreamClient(_ownedChannel);
    }

    /// <summary>
    /// Advanced constructor: the <see cref="CallInvoker"/> is injected, not
    /// owned — the caller controls channel lifetime and sharing (e.g. one
    /// channel behind several client surfaces), and tests can supply a fake
    /// invoker without a real server (see
    /// <c>PgChangeFeed.Client.Tests.Grpc</c>).
    /// </summary>
    public PgChangeFeedGrpcClient(CallInvoker callInvoker, PgChangeFeedClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(callInvoker);
        ArgumentNullException.ThrowIfNull(options);
        _options = options;
        _client = new ChangeStream.ChangeStreamClient(callInvoker);
    }

    /// <summary>
    /// <c>StreamChanges</c> — opens the server-streaming RPC and yields every
    /// <see cref="Change"/> the server sends from connection time onward
    /// (SPEC-020: fire-and-forget, no replay, one message per row change in
    /// commit order). The request carries no filter (SPEC-020: table-granular
    /// filtering is not part of this version).
    ///
    /// A missing or invalid bearer token ends the call with gRPC status
    /// <see cref="StatusCode.Unauthenticated"/> (SPEC-020 Negative) — this
    /// surfaces as an <see cref="RpcException"/> from the enumeration itself,
    /// not a swallowed empty stream.
    /// </summary>
    public async IAsyncEnumerable<Change> StreamChangesAsync(
        [EnumeratorCancellation] CancellationToken cancellationToken = default)
    {
        var headers = new Metadata { { AuthorizationMetadataKey, BearerPrefix + _options.ApiToken } };
        using var call = _client.StreamChanges(
            new StreamChangesRequest(), headers: headers, cancellationToken: cancellationToken);

        await foreach (var change in call.ResponseStream.ReadAllAsync(cancellationToken).ConfigureAwait(false))
        {
            yield return change;
        }
    }

    /// <summary>Disposes the channel this instance owns, if any (see the convenience constructor).</summary>
    public void Dispose() => _ownedChannel?.Dispose();
}
