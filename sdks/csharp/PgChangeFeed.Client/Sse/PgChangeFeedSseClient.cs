using System.Net.Http.Headers;
using System.Runtime.CompilerServices;
using System.Text.Json;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Sse.Models;

namespace PgChangeFeed.Client.Sse;

/// <summary>
/// Client for the PG Change Feed live change stream over Server-Sent Events:
/// one method opens <c>GET /changes/stream</c>, assembles each SSE frame via
/// <see cref="SseFrameParser"/> and yields the ten message fields as
/// <see cref="Change"/> — the same form as the gRPC client
/// (<see cref="PgChangeFeed.Client.Grpc.PgChangeFeedGrpcClient.StreamChangesAsync"/>).
///
/// A <see cref="System.Net.Http.HttpClient"/> passed to the
/// <c>(httpClient, options)</c> constructor is passed in, not owned —
/// same as <see cref="PgChangeFeedHttpClient"/>, with one extra obligation
/// for this client: <see cref="System.Net.Http.HttpClient.Timeout"/> covers
/// the *entire* request/response lifetime for a streamed response, not just
/// the time to the response headers — a long-lived stream therefore needs
/// <c>Timeout = <see cref="Timeout.InfiniteTimeSpan"/></c> on the passed-in
/// client, otherwise the connection ends after the default timeout even
/// without a connection problem. The bearer token is sent in the
/// <c>Authorization</c> header as <c>Bearer &lt;token&gt;</c>, same as
/// <see cref="PgChangeFeedHttpClient"/>. There is no global or static state —
/// a process can hold several independently configured instances at once.
///
/// <b>Limits:</b> no delivery guarantee and no in-stream replay; a
/// disconnected or slow-reading consumer misses the affected messages
/// permanently. Missed changes remain recoverable through the read path
/// (<see cref="PgChangeFeedHttpClient.ReadChangesAsync"/>).
/// </summary>
public sealed class PgChangeFeedSseClient : IDisposable
{
    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNameCaseInsensitive = true,
    };

    private readonly HttpClient _httpClient;
    private readonly PgChangeFeedClientOptions _options;
    private readonly HttpClient? _ownedHttpClient;

    /// <summary>
    /// Creates the client over an <see cref="HttpClient"/> it owns, with an
    /// infinite timeout for the long-lived stream, and <see cref="Dispose"/>
    /// disposes it. An <c>https</c> address connects over TLS with the trust
    /// anchors of the operating system, or exactly the certificates of
    /// <see cref="PgChangeFeedClientOptions.TrustAnchorFile"/> when it is set;
    /// chain, validity period and server name are always checked, and a failed
    /// check surfaces as the <see cref="HttpRequestException"/> of the
    /// <see cref="HttpClient"/>. A trust anchor with an <c>http</c> address
    /// throws an <see cref="ArgumentException"/>.
    /// </summary>
    /// <param name="options">The HTTP base URL, the bearer token and the optional trust anchor.</param>
    public PgChangeFeedSseClient(PgChangeFeedClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(options);
        _options = options;
        _ownedHttpClient = TlsTransport.CreateHttpClient(options, Timeout.InfiniteTimeSpan);
        _httpClient = _ownedHttpClient;
    }

    /// <summary>Disposes the <see cref="HttpClient"/> this instance owns, if any (see the options constructor).</summary>
    public void Dispose() => _ownedHttpClient?.Dispose();

    /// <summary>Creates the client over an <see cref="HttpClient"/> the caller owns.</summary>
    /// <param name="httpClient">The client used for the stream request; its timeout must be infinite for a long-lived stream.</param>
    /// <param name="options">The HTTP base URL and the bearer token.</param>
    public PgChangeFeedSseClient(HttpClient httpClient, PgChangeFeedClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(httpClient);
        ArgumentNullException.ThrowIfNull(options);
        _httpClient = httpClient;
        _options = options;
    }

    /// <summary>
    /// Opens <c>GET /changes/stream</c> and yields every <see cref="Change"/>
    /// the server sends from connection time onward. A non-success response
    /// while opening the stream (e.g. a missing/unknown bearer token, or
    /// <c>503</c> when the stream is not available) becomes a typed
    /// <see cref="PgChangeFeedException"/>, the same set
    /// <see cref="PgChangeFeedHttpClient"/> throws. Once the stream is open,
    /// its end (regular server-side close, or an incomplete trailing frame)
    /// ends the enumeration without an exception — the same fire-and-forget
    /// behavior as
    /// <see cref="PgChangeFeed.Client.Grpc.PgChangeFeedGrpcClient.StreamChangesAsync"/>;
    /// a frame whose <c>data:</c> payload cannot be read still throws
    /// <see cref="PgChangeFeedMalformedResponseException"/>, because that is
    /// not a stream end.
    ///
    /// <paramref name="target"/> selects the delivery target of a change: a
    /// set value (sent as the query parameter <c>target</c>) delivers only
    /// changes routed to that target; left <c>null</c> (the default) the
    /// request carries no query and delivers every change. A target no change
    /// carries delivers nothing and raises no error.
    ///
    /// <paramref name="schema"/> and <paramref name="table"/> are each
    /// optional and independent, sent as the query parameters
    /// <c>schema</c> and <c>table</c> only when set: a set
    /// <paramref name="schema"/> without <paramref name="table"/> delivers
    /// every table of that schema, a set <paramref name="table"/> without
    /// <paramref name="schema"/> delivers every table of that name regardless
    /// of schema, both set delivers exactly one table. They combine with
    /// <paramref name="target"/> as a conjunction. Pass the three by name:
    /// they come after <paramref name="cancellationToken"/>.
    /// </summary>
    public async IAsyncEnumerable<Change> StreamChangesAsync(
        [EnumeratorCancellation] CancellationToken cancellationToken = default,
        string? target = null,
        string? schema = null,
        string? table = null)
    {
        var query = new List<string>();
        if (schema is not null)
        {
            query.Add($"schema={Uri.EscapeDataString(schema)}");
        }
        if (table is not null)
        {
            query.Add($"table={Uri.EscapeDataString(table)}");
        }
        if (target is not null)
        {
            query.Add($"target={Uri.EscapeDataString(target)}");
        }
        var path = query.Count == 0
            ? "/changes/stream"
            : "/changes/stream?" + string.Join("&", query);
        var uri = new Uri(_options.Address, path);
        using var request = new HttpRequestMessage(HttpMethod.Get, uri);
        request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", _options.ApiToken);

        using var response = await _httpClient
            .SendAsync(request, HttpCompletionOption.ResponseHeadersRead, cancellationToken)
            .ConfigureAwait(false);

        if (!response.IsSuccessStatusCode)
        {
            var errorBody = await response.Content.ReadAsStringAsync(cancellationToken).ConfigureAwait(false);
            throw BuildException((int)response.StatusCode, errorBody);
        }

        var statusCode = (int)response.StatusCode;
        using var stream = await response.Content.ReadAsStreamAsync(cancellationToken).ConfigureAwait(false);
        using var reader = new StreamReader(stream);

        while (true)
        {
            var frame = await SseFrameParser.ReadFrameAsync(reader, cancellationToken).ConfigureAwait(false);
            if (frame is null)
            {
                yield break;
            }
            yield return ParseChange(frame, statusCode);
        }
    }

    private static Change ParseChange(SseFrame frame, int statusCode)
    {
        Change? change;
        try
        {
            change = JsonSerializer.Deserialize<Change>(frame.Data, JsonOptions);
        }
        catch (JsonException ex)
        {
            throw new PgChangeFeedMalformedResponseException(
                statusCode,
                "PG Change Feed SSE stream delivered a frame whose data payload is not valid " +
                "JSON.",
                ex);
        }

        return change
            ?? throw new PgChangeFeedMalformedResponseException(
                statusCode,
                "PG Change Feed SSE stream delivered a frame whose data payload is empty or " +
                "null.");
    }

    private static PgChangeFeedException BuildException(int statusCode, string body)
    {
        var (message, messageCode) = ErrorBody.Parse(body, JsonOptions);
        return statusCode switch
        {
            400 => new PgChangeFeedBadRequestException(statusCode, message, messageCode),
            401 => new PgChangeFeedUnauthorizedException(statusCode, message, messageCode),
            403 => new PgChangeFeedForbiddenException(statusCode, message, messageCode),
            404 => new PgChangeFeedNotFoundException(statusCode, message, messageCode),
            500 => new PgChangeFeedServerErrorException(statusCode, message, messageCode),
            _ => new PgChangeFeedUnexpectedStatusException(statusCode, message, messageCode),
        };
    }
}
