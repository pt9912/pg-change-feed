using System.Net.Http.Headers;
using System.Runtime.CompilerServices;
using System.Text.Json;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Sse.Models;

namespace PgChangeFeed.Client.Sse;

/// <summary>
/// Public entry point for the PG Change Feed live-change stream over
/// Server-Sent-Events (<c>SPEC-021</c>, <c>LH-FA-SST-008</c>): one method
/// opens <c>GET /changes/stream</c>, assembles each SSE frame via
/// <see cref="SseFrameParser"/> and yields the ten SPEC-021 message fields
/// as <see cref="Change"/> — the same idiomatic form as the existing gRPC
/// surface (<see cref="PgChangeFeed.Client.Grpc.PgChangeFeedGrpcClient.StreamChangesAsync"/>).
///
/// The <see cref="System.Net.Http.HttpClient"/> is injected, not owned —
/// same as <see cref="PgChangeFeedHttpClient"/>, with one extra obligation
/// for this surface: <see cref="System.Net.Http.HttpClient.Timeout"/> covers
/// the *entire* request/response lifetime for a streamed response, not just
/// the time to the response headers — a long-lived stream therefore needs
/// <c>Timeout = <see cref="Timeout.InfiniteTimeSpan"/></c> on the injected
/// client (pattern: <c>examples/csharp/sse-client/Program.cs</c>), otherwise
/// the connection ends after the default timeout even without a connection
/// problem. The bearer token is sent in the <c>Authorization</c> header as
/// <c>Bearer &lt;token&gt;</c> (SPEC-021), same as <see cref="PgChangeFeedHttpClient"/>.
/// No global or static state — a process can hold several independently
/// configured instances at once.
///
/// <b>Boundary (SPEC-021, LH-FA-SST-008):</b> no delivery guarantee and no
/// in-stream replay; a disconnected or slow-reading consumer misses the
/// affected messages permanently. Missed changes remain recoverable through
/// the existing read path (<see cref="PgChangeFeedHttpClient.ReadChangesAsync"/>).
///
/// Draht-Kenntnis-Vorbild (gelesen, nicht importiert — ADR-0106 Festlegung 2):
/// <c>examples/csharp/sse-client/SseStream.cs</c>, <c>Program.cs</c>.
/// </summary>
public sealed class PgChangeFeedSseClient
{
    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNameCaseInsensitive = true,
    };

    private readonly HttpClient _httpClient;
    private readonly PgChangeFeedClientOptions _options;

    public PgChangeFeedSseClient(HttpClient httpClient, PgChangeFeedClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(httpClient);
        ArgumentNullException.ThrowIfNull(options);
        _httpClient = httpClient;
        _options = options;
    }

    /// <summary>
    /// <c>GET /changes/stream</c> (SPEC-021) — opens the stream and yields
    /// every <see cref="Change"/> the server sends from connection time
    /// onward. A non-success response while opening the stream (e.g. a
    /// missing/unknown bearer token, or <c>503</c> when no
    /// <c>Broadcaster</c> is wired) becomes a typed
    /// <see cref="PgChangeFeedException"/>, the same closed set
    /// <see cref="PgChangeFeedHttpClient"/> throws. Once the stream is open,
    /// its end (regular server-side close, or an incomplete trailing frame)
    /// ends the enumeration without an exception — the same fire-and-forget
    /// semantics as
    /// <see cref="PgChangeFeed.Client.Grpc.PgChangeFeedGrpcClient.StreamChangesAsync"/>;
    /// a frame whose <c>data:</c> payload does not parse as the documented
    /// SPEC-021 shape still throws <see cref="PgChangeFeedMalformedResponseException"/>,
    /// because that is a protocol violation, not a stream end.
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
                "JSON — a protocol violation outside SPEC-021's documented shape.",
                ex);
        }

        return change
            ?? throw new PgChangeFeedMalformedResponseException(
                statusCode,
                "PG Change Feed SSE stream delivered a frame whose data payload is empty or " +
                "null — a protocol violation outside SPEC-021's documented shape.");
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
