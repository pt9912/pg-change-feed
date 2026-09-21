using System.Runtime.CompilerServices;
using System.Text.Json;
using NATS.Client.Core;
using NATS.Net;
using PgChangeFeed.Client.Nats.Models;

namespace PgChangeFeed.Client.Nats;

/// <summary>
/// Public entry point for the PG Change Feed NATS full-content stream
/// (<c>SPEC-024</c>, <c>LH-FA-SST-008</c>): one method subscribes to the
/// four-token subject namespace <c>cdc.stream.&lt;source_id&gt;.&lt;schema&gt;.&lt;table&gt;</c>
/// (or a wildcard pattern over it) and yields the ten SPEC-024 message
/// fields as <see cref="Change"/> — the same idiomatic form as the existing
/// gRPC/SSE surfaces (<see cref="PgChangeFeed.Client.Grpc.PgChangeFeedGrpcClient.StreamChangesAsync"/>,
/// <see cref="PgChangeFeed.Client.Sse.PgChangeFeedSseClient.StreamChangesAsync"/>).
///
/// Authentication is connection-level, not per-call (SPEC-024): the NATS
/// server rejects the connection itself when a server-wide token is
/// configured and the client's token is missing or wrong — there is no
/// per-message header to attach, unlike the HTTP/gRPC/SSE surfaces' bearer
/// token. <see cref="PgChangeFeedClientOptions"/> is still the shared
/// denominator this class reuses: <see cref="PgChangeFeedClientOptions.Address"/>
/// becomes the NATS server URL (e.g. <c>nats://host:4222</c>) and
/// <see cref="PgChangeFeedClientOptions.ApiToken"/> becomes the connection
/// token — the same two-value shape as every other surface, carried over a
/// different transport.
///
/// <b>Boundary (SPEC-024, LH-FA-SST-008):</b> no delivery guarantee and no
/// in-stream replay (Core NATS, fire-and-forget); a disconnected or slow-reading
/// consumer misses the affected messages permanently. Missed changes remain
/// recoverable through the existing read path
/// (<see cref="PgChangeFeed.Client.Http.PgChangeFeedHttpClient.ReadChangesAsync"/>).
/// A rejected connection (missing/wrong token) surfaces as a NATS-native
/// exception (typically <see cref="NatsServerException"/> or
/// <see cref="NatsConnectionFailedException"/>) thrown from the enumeration itself
/// — the same "not a swallowed empty stream" boundary as the gRPC surface's
/// <c>Unauthenticated</c> <see cref="Grpc.Core.RpcException"/>, propagated
/// unwrapped rather than remapped into a second exception hierarchy: NATS
/// connection errors are not HTTP status codes, and inventing a parallel
/// mapping here would only risk drifting from what <c>NATS.Net</c> itself
/// already throws.
///
/// Draht-Kenntnis-Vorbild (gelesen, nicht importiert — ADR-0106 Festlegung 2):
/// <c>examples/csharp/nats-stream-client/Format.cs</c>, <c>Program.cs</c>.
/// </summary>
public sealed class PgChangeFeedNatsStreamClient : IAsyncDisposable
{
    /// <summary>
    /// The full-content namespace's root wildcard — every source, every
    /// table (SPEC-024: <c>cdc.stream.&gt;</c>). The default
    /// <see cref="StreamChangesAsync"/> subscribes to when no narrower
    /// subject is supplied.
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
    /// Advanced constructor: the <see cref="INatsClient"/> is injected, not
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
    /// Builds the four-token SPEC-024 subject for one specific table:
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
    /// <c>cdc.stream.&lt;sourceId&gt;.&gt;</c> (SPEC-024).
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
    /// server delivers from subscription time onward (SPEC-024: fire-and-forget,
    /// no replay, one message per row change in commit order). Use
    /// <see cref="BuildSubject"/>/<see cref="BuildSourceSubject"/> to narrow
    /// the subject to one table or one source.
    ///
    /// A connection rejected by the NATS server (missing/wrong token, SPEC-024
    /// Negative) ends the enumeration with a NATS-native exception — see the
    /// class-level boundary note. A message whose payload does not parse as
    /// the documented SPEC-024 shape throws
    /// <see cref="PgChangeFeedNatsMalformedMessageException"/>, because that
    /// is a protocol violation, not a stream end.
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
                "PG Change Feed NATS full-content stream delivered a message whose payload is " +
                "not valid JSON — a protocol violation outside SPEC-024's documented shape.",
                ex);
        }

        return change
            ?? throw new PgChangeFeedNatsMalformedMessageException(
                "PG Change Feed NATS full-content stream delivered a message whose payload is " +
                "empty or null — a protocol violation outside SPEC-024's documented shape.");
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
