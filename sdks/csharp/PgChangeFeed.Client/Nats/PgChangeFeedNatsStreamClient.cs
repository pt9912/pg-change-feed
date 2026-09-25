using System.Runtime.CompilerServices;
using System.Text.Json;
using NATS.Client.Core;
using NATS.Net;
using PgChangeFeed.Client.Nats.Models;

namespace PgChangeFeed.Client.Nats;

/// <summary>
/// Client for the PG Change Feed live change stream over NATS: one method
/// subscribes to the subject namespace
/// <c>cdc.stream.&lt;source_id&gt;.&lt;schema&gt;.&lt;table&gt;</c> (or a
/// wildcard pattern over it) and yields the ten message fields as
/// <see cref="Change"/> — the same form as the gRPC and SSE clients
/// (<see cref="PgChangeFeed.Client.Grpc.PgChangeFeedGrpcClient.StreamChangesAsync"/>,
/// <see cref="PgChangeFeed.Client.Sse.PgChangeFeedSseClient.StreamChangesAsync"/>).
///
/// Authentication is connection-level, not per call: the NATS server rejects
/// the connection itself when it is configured with a token and the client's
/// token is missing or wrong — there is no per-message header to attach.
/// <see cref="PgChangeFeedClientOptions"/> is reused: its
/// <see cref="PgChangeFeedClientOptions.Address"/> is the NATS server URL (e.g.
/// <c>nats://host:4222</c>) and its <see cref="PgChangeFeedClientOptions.ApiToken"/>
/// is the connection token.
///
/// <b>Limits:</b> no delivery guarantee and no in-stream replay
/// (fire-and-forget); a disconnected or slow-reading consumer misses the
/// affected messages permanently. Missed changes remain recoverable through
/// the read path
/// (<see cref="PgChangeFeed.Client.Http.PgChangeFeedHttpClient.ReadChangesAsync"/>).
/// A rejected connection (missing/wrong token) surfaces as a NATS exception
/// (typically <see cref="NatsServerException"/> or
/// <see cref="NatsConnectionFailedException"/>) thrown from the enumeration
/// itself — the stream is never silently empty, the same as the gRPC client's
/// <c>Unauthenticated</c> <see cref="global::Grpc.Core.RpcException"/>. The exception
/// propagates as <c>NATS.Net</c> throws it; it is not mapped into a second
/// exception hierarchy.
/// </summary>
public sealed class PgChangeFeedNatsStreamClient : IAsyncDisposable
{
    /// <summary>
    /// The root wildcard of the stream namespace — every source, every table
    /// (<c>cdc.stream.&gt;</c>). The default <see cref="StreamChangesAsync"/>
    /// subscribes to when no narrower subject is supplied.
    /// </summary>
    public const string AllSourcesSubject = "cdc.stream.>";

    private static readonly char[] InvalidTokenChars = ['.', '*', '>', ' ', '\t', '\n', '\r'];

    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNameCaseInsensitive = true,
    };

    private readonly INatsClient _client;
    private readonly INatsClient? _ownedClient;

    /// <summary>
    /// Convenience constructor: builds and owns its own <see cref="INatsClient"/>
    /// against <paramref name="options"/>'s address and token.
    /// <see cref="DisposeAsync"/> disposes that client.
    /// </summary>
    public PgChangeFeedNatsStreamClient(PgChangeFeedClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(options);
        var natsOpts = new NatsOpts
        {
            Url = options.Address.ToString(),
            AuthOpts = new NatsAuthOpts { Token = options.ApiToken },
        };
        var owned = new NatsClient(natsOpts);
        _client = owned;
        _ownedClient = owned;
    }

    /// <summary>
    /// Advanced constructor: the <see cref="INatsClient"/> is passed in, not
    /// owned — the caller controls connection lifetime and sharing (e.g. one
    /// connection behind several subjects/subscriptions), and tests can
    /// supply a fake client without a real server (see
    /// <c>PgChangeFeed.Client.Tests.Nats</c>).
    /// </summary>
    public PgChangeFeedNatsStreamClient(INatsClient client)
    {
        ArgumentNullException.ThrowIfNull(client);
        _client = client;
        _ownedClient = null;
    }

    /// <summary>
    /// Builds the four-token subject for one specific table:
    /// <c>cdc.stream.&lt;sourceId&gt;.&lt;schema&gt;.&lt;table&gt;</c>. Each
    /// token is validated to contain none of NATS's own token separator
    /// (<c>.</c>) or wildcard characters (<c>*</c>, <c>&gt;</c>) — a token
    /// carrying one of these would silently change which subjects the
    /// resulting string matches, rather than fail loudly.
    /// </summary>
    public static string BuildSubject(string sourceId, string schema, string table)
    {
        ValidateToken(sourceId, nameof(sourceId));
        ValidateToken(schema, nameof(schema));
        ValidateToken(table, nameof(table));
        return $"cdc.stream.{sourceId}.{schema}.{table}";
    }

    /// <summary>
    /// Builds the three-token wildcard subject for every table of one source:
    /// <c>cdc.stream.&lt;sourceId&gt;.&gt;</c>.
    /// </summary>
    public static string BuildSourceSubject(string sourceId)
    {
        ValidateToken(sourceId, nameof(sourceId));
        return $"cdc.stream.{sourceId}.>";
    }

    private static void ValidateToken(string value, string paramName)
    {
        if (string.IsNullOrWhiteSpace(value))
        {
            throw new ArgumentException("A subject token must not be null, empty, or whitespace.", paramName);
        }
        if (value.IndexOfAny(InvalidTokenChars) >= 0)
        {
            throw new ArgumentException(
                "A subject token must not contain '.', '*', '>', or whitespace — these are NATS " +
                "subject separators/wildcards, not part of a token's own value.", paramName);
        }
    }

    /// <summary>
    /// Subscribes to <paramref name="subject"/> (default: <see cref="AllSourcesSubject"/>,
    /// every source and table) and yields every <see cref="Change"/> the NATS
    /// server delivers from subscription time onward (fire-and-forget, no
    /// replay, one message per row change in commit order). Use
    /// <see cref="BuildSubject"/>/<see cref="BuildSourceSubject"/> to narrow
    /// the subject to one table or one source.
    ///
    /// A connection rejected by the NATS server (missing/wrong token) ends the
    /// enumeration with a NATS exception — see the class-level boundary note.
    /// A message whose payload cannot be read throws
    /// <see cref="PgChangeFeedNatsMalformedMessageException"/>, because that
    /// is not a stream end.
    /// </summary>
    public async IAsyncEnumerable<Change> StreamChangesAsync(
        string subject = AllSourcesSubject,
        [EnumeratorCancellation] CancellationToken cancellationToken = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(subject);

        await foreach (var msg in _client
            .SubscribeAsync<byte[]>(subject, cancellationToken: cancellationToken)
            .ConfigureAwait(false))
        {
            yield return ParseChange(msg.Data);
        }
    }

    private static Change ParseChange(byte[]? data)
    {
        Change? change;
        try
        {
            change = data is null ? null : JsonSerializer.Deserialize<Change>(data, JsonOptions);
        }
        catch (JsonException ex)
        {
            throw new PgChangeFeedNatsMalformedMessageException(
                "PG Change Feed NATS stream delivered a message whose payload is " +
                "not valid JSON.",
                ex);
        }

        return change
            ?? throw new PgChangeFeedNatsMalformedMessageException(
                "PG Change Feed NATS stream delivered a message whose payload is " +
                "empty or null.");
    }

    /// <summary>Disposes the NATS client this instance owns, if any (see the convenience constructor).</summary>
    public async ValueTask DisposeAsync()
    {
        if (_ownedClient is not null)
        {
            await _ownedClient.DisposeAsync().ConfigureAwait(false);
        }
    }
}
