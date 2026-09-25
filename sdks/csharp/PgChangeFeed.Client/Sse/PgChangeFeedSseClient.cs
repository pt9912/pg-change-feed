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
/// The <see cref="System.Net.Http.HttpClient"/> is passed in, not owned —
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
public sealed class PgChangeFeedSseClient
{
    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNameCaseInsensitive = true,
    };

    private readonly HttpClient _httpClient;
    private readonly PgChangeFeedClientOptions _options;

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
    /// </summary>
    public async IAsyncEnumerable<Change> StreamChangesAsync(
        [EnumeratorCancellation] CancellationToken cancellationToken = default)
    {
        var uri = new Uri(_options.Address, "/changes/stream");
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
        var message = ExtractErrorMessage(body);
        return statusCode switch
        {
            400 => new PgChangeFeedBadRequestException(statusCode, message),
            401 => new PgChangeFeedUnauthorizedException(statusCode, message),
            403 => new PgChangeFeedForbiddenException(statusCode, message),
            404 => new PgChangeFeedNotFoundException(statusCode, message),
            500 => new PgChangeFeedServerErrorException(statusCode, message),
            _ => new PgChangeFeedUnexpectedStatusException(statusCode, message),
        };
    }

    private static string ExtractErrorMessage(string body)
    {
        try
        {
            var error = JsonSerializer.Deserialize<PgChangeFeed.Client.Http.Models.ErrorResponse>(body, JsonOptions);
            return error?.Error ?? body;
        }
        catch (JsonException)
        {
            return body;
        }
    }
}
