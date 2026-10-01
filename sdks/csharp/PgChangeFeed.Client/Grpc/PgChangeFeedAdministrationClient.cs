using Cdc.Administration.V1;
using Grpc.Core;
using Grpc.Net.Client;

namespace PgChangeFeed.Client.Grpc;

/// <summary>
/// Client for the PG Change Feed <c>Administration</c> gRPC service: one
/// method per RPC — <see cref="RegisterConsumerAsync"/>,
/// <see cref="AcknowledgeConsumerAsync"/>, <see cref="GetConsumerPositionAsync"/>,
/// <see cref="RemoveConsumerAsync"/>, <see cref="EnableTableAsync"/>,
/// <see cref="DisableTableAsync"/>, <see cref="GetTableStatusAsync"/>,
/// <see cref="ListTablesAsync"/>, <see cref="RunRetentionAsync"/>,
/// <see cref="ReadChangesAsync"/> and <see cref="DiagnoseAsync"/>. Requests
/// and responses are the generated protobuf messages of
/// <c>Cdc.Administration.V1</c> unchanged — like <see cref="PgChangeFeedGrpcClient"/>,
/// there is no separate DTO layer: a protobuf message is already a typed C#
/// class, and a hand-written mirror would only duplicate eleven message
/// shapes without a deserialization need to justify it (unlike
/// <see cref="PgChangeFeed.Client.Http.PgChangeFeedHttpClient"/>, whose DTOs
/// exist because JSON needs a target type to deserialize into). Every
/// non-<see cref="StatusCode.OK"/> gRPC status becomes a typed
/// <see cref="PgChangeFeedGrpcException"/> subtype instead of a result type
/// mixed with the success path, consistent with
/// <see cref="PgChangeFeed.Client.Http.PgChangeFeedHttpClient"/>'s exception
/// design for the HTTP surface.
///
/// The bearer token is sent in the <c>authorization</c> call metadata entry
/// as <c>Bearer &lt;token&gt;</c> on every call, from <see cref="PgChangeFeedClientOptions"/>
/// supplied at construction — there is no global or static state; a process
/// can hold several independently configured instances at once.
/// </summary>
public sealed class PgChangeFeedAdministrationClient : IDisposable
{
    private const string AuthorizationMetadataKey = "authorization";
    private const string BearerPrefix = "Bearer ";

    private readonly Administration.AdministrationClient _client;
    private readonly PgChangeFeedClientOptions _options;
    private readonly GrpcChannel? _ownedChannel;

    /// <summary>
    /// Convenience constructor: opens and owns its own <see cref="GrpcChannel"/>
    /// against <paramref name="options"/>'s address. <see cref="Dispose"/>
    /// disposes that channel.
    /// </summary>
    public PgChangeFeedAdministrationClient(PgChangeFeedClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(options);
        _options = options;
        _ownedChannel = GrpcChannel.ForAddress(options.Address);
        _client = new Administration.AdministrationClient(_ownedChannel);
    }

    /// <summary>
    /// Advanced constructor: the <see cref="CallInvoker"/> is injected, not
    /// owned — the caller controls channel lifetime and sharing (e.g. one
    /// channel behind several client surfaces), and tests can supply a fake
    /// invoker without a real server.
    /// </summary>
    public PgChangeFeedAdministrationClient(CallInvoker callInvoker, PgChangeFeedClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(callInvoker);
        ArgumentNullException.ThrowIfNull(options);
        _options = options;
        _client = new Administration.AdministrationClient(callInvoker);
    }

    /// <summary>
    /// Registers a consumer, a named reader whose position the server keeps
    /// (admin token). Registering an existing consumer changes nothing; the
    /// response reports it with <c>AlreadyRegistered</c>.
    /// </summary>
    public Task<RegisterConsumerResponse> RegisterConsumerAsync(
        RegisterConsumerRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return CallAsync(() => _client.RegisterConsumerAsync(
            request, headers: Headers(), cancellationToken: cancellationToken).ResponseAsync);
    }

    /// <summary>
    /// Stores the position up to which a consumer has processed a source
    /// (admin token). Repeating the stored position has no effect; an
    /// earlier position, or a position of another source, raises
    /// <see cref="PgChangeFeedGrpcInvalidArgumentException"/>.
    /// </summary>
    public Task<AcknowledgeConsumerResponse> AcknowledgeConsumerAsync(
        AcknowledgeConsumerRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return CallAsync(() => _client.AcknowledgeConsumerAsync(
            request, headers: Headers(), cancellationToken: cancellationToken).ResponseAsync);
    }

    /// <summary>
    /// Reads the stored position of a consumer (reader or admin token).
    /// <c>Acknowledged</c> is <c>false</c> for a consumer that never
    /// acknowledged; <c>Offset</c> is then the starting position.
    /// </summary>
    public Task<GetConsumerPositionResponse> GetConsumerPositionAsync(
        GetConsumerPositionRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return CallAsync(() => _client.GetConsumerPositionAsync(
            request, headers: Headers(), cancellationToken: cancellationToken).ResponseAsync);
    }

    /// <summary>
    /// Removes a consumer (admin token); <c>Removed</c> is <c>false</c> for
    /// one that was never registered.
    /// </summary>
    public Task<RemoveConsumerResponse> RemoveConsumerAsync(
        RemoveConsumerRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return CallAsync(() => _client.RemoveConsumerAsync(
            request, headers: Headers(), cancellationToken: cancellationToken).ResponseAsync);
    }

    /// <summary>
    /// Starts capturing a table of a source (admin token). <c>AlreadyEnabled</c>
    /// is <c>true</c> when the table was captured already; a table that does
    /// not exist in the source database raises
    /// <see cref="PgChangeFeedGrpcNotFoundException"/>.
    /// </summary>
    public Task<EnableTableResponse> EnableTableAsync(
        EnableTableRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return CallAsync(() => _client.EnableTableAsync(
            request, headers: Headers(), cancellationToken: cancellationToken).ResponseAsync);
    }

    /// <summary>
    /// Stops capturing a table (admin token). <c>Retained</c> is <c>true</c>
    /// when changes already stored for the table remain readable; a table
    /// that does not exist in the source database raises
    /// <see cref="PgChangeFeedGrpcNotFoundException"/>.
    /// </summary>
    public Task<DisableTableResponse> DisableTableAsync(
        DisableTableRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return CallAsync(() => _client.DisableTableAsync(
            request, headers: Headers(), cancellationToken: cancellationToken).ResponseAsync);
    }

    /// <summary>
    /// Tells whether a table is captured (<c>Enabled</c>) or no longer
    /// captured with stored changes remaining (<c>Retained</c>) (reader or
    /// admin token). A table that was never enabled reports both as
    /// <c>false</c>; a table that does not exist in the source database
    /// raises <see cref="PgChangeFeedGrpcNotFoundException"/>.
    /// </summary>
    public Task<GetTableStatusResponse> GetTableStatusAsync(
        GetTableStatusRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return CallAsync(() => _client.GetTableStatusAsync(
            request, headers: Headers(), cancellationToken: cancellationToken).ResponseAsync);
    }

    /// <summary>
    /// Lists the captured tables and the tables that are no longer captured
    /// but whose stored changes remain (reader or admin token).
    /// </summary>
    public Task<ListTablesResponse> ListTablesAsync(
        ListTablesRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return CallAsync(() => _client.ListTablesAsync(
            request, headers: Headers(), cancellationToken: cancellationToken).ResponseAsync);
    }

    /// <summary>
    /// Deletes the stored changes of a source that are older than
    /// <c>MinAgeNanos</c> and that every consumer with a stored position has
    /// passed (admin token); <c>Deleted</c> is the number removed. A
    /// <c>MinAgeNanos</c> of <c>0</c> means no minimum age and is valid.
    /// </summary>
    public Task<RunRetentionResponse> RunRetentionAsync(
        RunRetentionRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return CallAsync(() => _client.RunRetentionAsync(
            request, headers: Headers(), cancellationToken: cancellationToken).ResponseAsync);
    }

    /// <summary>
    /// Reads a bounded range of stored changes of a source (reader or admin
    /// token) — the same filter and range semantics as
    /// <see cref="PgChangeFeed.Client.Http.PgChangeFeedHttpClient.ReadChangesAsync"/>:
    /// an unset <c>From</c>/<c>To</c>/<c>Limit</c> carries <c>0</c> (not set).
    /// The request's <c>Target</c> selects the delivery target: empty (the
    /// default) is no filter, a set value returns only changes routed to that
    /// target, combined with <c>Schema</c>/<c>Table</c> as a conjunction; a
    /// target no change carries returns an empty list.
    /// </summary>
    public Task<ReadChangesResponse> ReadChangesAsync(
        ReadChangesRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return CallAsync(() => _client.ReadChangesAsync(
            request, headers: Headers(), cancellationToken: cancellationToken).ResponseAsync);
    }

    /// <summary>
    /// Reads the operational diagnose report of a source (reader or admin
    /// token) — the same report as the CLI diagnose mode and
    /// <c>GET /diagnose</c>. A <c>Known</c>/<c>Present</c>/<c>*Known</c>
    /// field of <c>false</c> carries the respective absence case (no
    /// heartbeat ever written, no blocking consumer, unknown estimate/
    /// backlog).
    /// </summary>
    public Task<DiagnoseResponse> DiagnoseAsync(
        DiagnoseRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return CallAsync(() => _client.DiagnoseAsync(
            request, headers: Headers(), cancellationToken: cancellationToken).ResponseAsync);
    }

    /// <summary>Disposes the channel this instance owns, if any (see the convenience constructor).</summary>
    public void Dispose() => _ownedChannel?.Dispose();

    private Metadata Headers() => new() { { AuthorizationMetadataKey, BearerPrefix + _options.ApiToken } };

    private static async Task<TResponse> CallAsync<TResponse>(Func<Task<TResponse>> call)
    {
        try
        {
            return await call().ConfigureAwait(false);
        }
        catch (RpcException ex)
        {
            throw MapException(ex);
        }
    }

    private static PgChangeFeedGrpcException MapException(RpcException ex) => ex.StatusCode switch
    {
        StatusCode.InvalidArgument => new PgChangeFeedGrpcInvalidArgumentException(ex.Status.Detail, ex),
        StatusCode.Unauthenticated => new PgChangeFeedGrpcUnauthenticatedException(ex.Status.Detail, ex),
        StatusCode.PermissionDenied => new PgChangeFeedGrpcPermissionDeniedException(ex.Status.Detail, ex),
        StatusCode.NotFound => new PgChangeFeedGrpcNotFoundException(ex.Status.Detail, ex),
        StatusCode.Internal => new PgChangeFeedGrpcInternalException(ex.Status.Detail, ex),
        _ => new PgChangeFeedGrpcUnexpectedStatusException(ex.StatusCode, ex.Status.Detail, ex),
    };
}
