using System.Runtime.CompilerServices;
using Cdc.Stream.V1;
using Grpc.Core;
using Grpc.Net.Client;

namespace PgChangeFeed.Client.Grpc;

/// <summary>
/// Client for the PG Change Feed live change stream over gRPC: one method
/// opens the <c>ChangeStream/StreamChanges</c> server-streaming call and
/// yields the generated <see cref="Change"/> message with the ten fields
/// <c>change_id</c>, <c>transaction_id</c>, <c>source_table_id</c>,
/// <c>sequence</c>, <c>operation</c>, <c>old_image</c>, <c>new_image</c>,
/// <c>schema_version</c>, <c>schema</c> and <c>table</c>, exactly as the
/// server sends them. There is no separate DTO layer (unlike
/// <c>PgChangeFeed.Client.Http.Models</c>): the generated message already is
/// the typed form of the stream schema.
///
/// The bearer token is sent in the <c>authorization</c> call metadata entry
/// as <c>Bearer &lt;token&gt;</c> on every call — there is no global or static
/// state; a process can hold several independently configured instances at
/// once.
///
/// <b>Limits:</b> the stream carries no replay and cannot be filtered by
/// table. A consumer that needs either uses the read path
/// (<c>PgChangeFeed.Client.Http.PgChangeFeedHttpClient.ReadChangesAsync</c>),
/// not this stream.
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
    /// The PG Change Feed server speaks plaintext gRPC over HTTP/2 (no TLS).
    /// An <c>http://</c> address connects without further setup on .NET 10;
    /// this constructor sets no process-wide <see cref="AppContext"/> switch.
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
    /// Opens the server-streaming call and yields every <see cref="Change"/>
    /// the server sends from connection time onward: fire-and-forget, no
    /// replay, one message per row change in commit order. The request carries
    /// no filter.
    ///
    /// A missing or invalid bearer token ends the call with gRPC status
    /// <see cref="StatusCode.Unauthenticated"/> — this surfaces as an
    /// <see cref="RpcException"/> from the enumeration itself, not as a
    /// silently empty stream.
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
