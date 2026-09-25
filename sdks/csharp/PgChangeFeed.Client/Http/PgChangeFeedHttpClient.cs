using System.Net.Http;
using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;
using PgChangeFeed.Client.Http.Models;

namespace PgChangeFeed.Client.Http;

/// <summary>
/// Client for the PG Change Feed HTTP/JSON API: one method per capability —
/// <c>RegisterConsumerAsync</c>, <c>AcknowledgeConsumerAsync</c>,
/// <c>GetConsumerPositionAsync</c>, <c>RemoveConsumerAsync</c>,
/// <c>EnableTableAsync</c>, <c>DisableTableAsync</c>, <c>GetStatusAsync</c>,
/// <c>ListTablesAsync</c>, <c>RunRetentionAsync</c> and
/// <c>ReadChangesAsync</c>. Requests and responses are typed DTOs that mirror
/// the JSON documents of the API (<c>PgChangeFeed.Client.Http.Models</c>);
/// every non-success response becomes a typed <see cref="PgChangeFeedException"/>
/// subtype instead of a result type mixed with the success path. Connection
/// errors and timeouts of the transport are not converted; they reach the
/// caller as the <see cref="HttpRequestException"/> or
/// <see cref="TaskCanceledException"/> of the <see cref="HttpClient"/>.
///
/// The <see cref="HttpClient"/> is passed in, not owned — the caller controls
/// its lifetime, connection pooling and any <c>DelegatingHandler</c>
/// pipeline (proxies, retries, logging); this type never disposes it. The
/// bearer token and server address come from <see cref="PgChangeFeedClientOptions"/>,
/// supplied at construction — there is no global or static state, so a
/// process can hold several independently configured instances at once.
/// </summary>
public sealed class PgChangeFeedHttpClient
{
    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNameCaseInsensitive = true,
    };

    private readonly HttpClient _httpClient;
    private readonly PgChangeFeedClientOptions _options;

    /// <summary>Creates the client over an <see cref="HttpClient"/> the caller owns.</summary>
    /// <param name="httpClient">The client used for every request.</param>
    /// <param name="options">The server address and the bearer token.</param>
    public PgChangeFeedHttpClient(HttpClient httpClient, PgChangeFeedClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(httpClient);
        ArgumentNullException.ThrowIfNull(options);
        _httpClient = httpClient;
        _options = options;
    }

    /// <summary>
    /// Registers a consumer, a named reader whose position the server keeps
    /// (<c>POST /consumers</c>, admin token). Registering an existing consumer
    /// changes nothing; the response reports it with <c>AlreadyRegistered</c>.
    /// </summary>
    public Task<RegisterConsumerResponse> RegisterConsumerAsync(
        RegisterConsumerRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return PostAsync<RegisterConsumerRequest, RegisterConsumerResponse>(
            "/consumers", request, cancellationToken);
    }

    /// <summary>
    /// Stores the position up to which a consumer has processed a source
    /// (<c>POST /consumers/acknowledge</c>, admin token). Repeating the stored
    /// position has no effect; an earlier position, or a position of another
    /// source, is rejected with <see cref="PgChangeFeedBadRequestException"/>.
    /// </summary>
    public Task<AcknowledgeConsumerResponse> AcknowledgeConsumerAsync(
        AcknowledgeConsumerRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return PostAsync<AcknowledgeConsumerRequest, AcknowledgeConsumerResponse>(
            "/consumers/acknowledge", request, cancellationToken);
    }

    /// <summary>
    /// Reads the stored position of a consumer (<c>GET /consumers/position</c>,
    /// reader or admin token). <c>Acknowledged</c> is <c>false</c> for a
    /// consumer that never acknowledged; <c>Offset</c> is then the starting
    /// position.
    /// </summary>
    public Task<ConsumerPositionResponse> GetConsumerPositionAsync(
        string consumerId, CancellationToken cancellationToken = default)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(consumerId);
        var query = BuildQuery(("consumer_id", consumerId));
        return GetAsync<ConsumerPositionResponse>($"/consumers/position?{query}", cancellationToken);
    }

    /// <summary>
    /// Removes a consumer (<c>POST /consumers/remove</c>, admin token);
    /// <c>Removed</c> is <c>false</c> for one that was never registered.
    /// </summary>
    public Task<RemoveConsumerResponse> RemoveConsumerAsync(
        string consumerId, CancellationToken cancellationToken = default)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(consumerId);
        return PostAsync<RemoveConsumerRequest, RemoveConsumerResponse>(
            "/consumers/remove", new RemoveConsumerRequest(consumerId), cancellationToken);
    }

    /// <summary>
    /// Starts capturing a table of a source (<c>POST /tables/enable</c>, admin
    /// token). <c>AlreadyEnabled</c> is <c>true</c> when the table was captured
    /// already; a table that does not exist in the source database raises
    /// <see cref="PgChangeFeedNotFoundException"/>.
    /// </summary>
    public Task<EnableTableResponse> EnableTableAsync(
        EnableTableRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return PostAsync<EnableTableRequest, EnableTableResponse>(
            "/tables/enable", request, cancellationToken);
    }

    /// <summary>
    /// Stops capturing a table (<c>POST /tables/disable</c>, admin token).
    /// <c>Retained</c> is <c>true</c> when changes already stored for the table
    /// remain readable; a table that does not exist in the source database
    /// raises <see cref="PgChangeFeedNotFoundException"/>.
    /// </summary>
    public Task<DisableTableResponse> DisableTableAsync(
        DisableTableRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return PostAsync<DisableTableRequest, DisableTableResponse>(
            "/tables/disable", request, cancellationToken);
    }

    /// <summary>
    /// Tells whether a table is captured (<c>Enabled</c>) or no longer captured
    /// with stored changes remaining (<c>Retained</c>)
    /// (<c>GET /tables/status</c>, reader or admin token). A table that was never
    /// enabled reports both as <c>false</c>; a table that does not exist in the
    /// source database raises <see cref="PgChangeFeedNotFoundException"/>.
    /// </summary>
    public Task<TableStatusResponse> GetStatusAsync(
        string source, string schema, string table, string publication, CancellationToken cancellationToken = default)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(source);
        ArgumentException.ThrowIfNullOrWhiteSpace(schema);
        ArgumentException.ThrowIfNullOrWhiteSpace(table);
        ArgumentException.ThrowIfNullOrWhiteSpace(publication);
        var query = BuildQuery(
            ("source", source), ("schema", schema), ("table", table), ("publication", publication));
        return GetAsync<TableStatusResponse>($"/tables/status?{query}", cancellationToken);
    }

    /// <summary>
    /// Lists the captured tables and the tables that are no longer captured but
    /// whose stored changes remain (<c>GET /tables</c>, reader or admin token).
    /// </summary>
    public Task<ListTablesResponse> ListTablesAsync(
        string source, string publication, CancellationToken cancellationToken = default)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(source);
        ArgumentException.ThrowIfNullOrWhiteSpace(publication);
        var query = BuildQuery(("source", source), ("publication", publication));
        return GetAsync<ListTablesResponse>($"/tables?{query}", cancellationToken);
    }

    /// <summary>
    /// Deletes the stored changes of a source that are older than
    /// <c>MinAgeNanos</c> and that every consumer with a stored position has
    /// passed (<c>POST /retention/run</c>, admin token); <c>Deleted</c> is the
    /// number removed.
    /// </summary>
    public Task<RunRetentionResponse> RunRetentionAsync(
        RunRetentionRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return PostAsync<RunRetentionRequest, RunRetentionResponse>(
            "/retention/run", request, cancellationToken);
    }

    /// <summary>
    /// Reads stored changes of a source (<c>GET /changes</c>, reader or admin
    /// token). <paramref name="source"/> is mandatory; <paramref name="schema"/>/
    /// <paramref name="table"/> are optional and independent;
    /// <paramref name="from"/> is inclusive, <paramref name="to"/> exclusive
    /// (both <c>commit_position</c> values). <paramref name="limit"/> cuts rows,
    /// not positions, and there is no default limit. A range without changes
    /// returns an empty list.
    /// </summary>
    public Task<ReadChangesResponse> ReadChangesAsync(
        string source,
        string? schema = null,
        string? table = null,
        long? from = null,
        long? to = null,
        int? limit = null,
        CancellationToken cancellationToken = default)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(source);
        var query = BuildQuery(
            ("source", source),
            ("schema", schema),
            ("table", table),
            ("from", from?.ToString()),
            ("to", to?.ToString()),
            ("limit", limit?.ToString()));
        return GetAsync<ReadChangesResponse>($"/changes?{query}", cancellationToken);
    }

    private static string BuildQuery(params (string Key, string? Value)[] parameters)
    {
        var pairs = parameters
            .Where(p => p.Value is not null)
            .Select(p => $"{p.Key}={Uri.EscapeDataString(p.Value!)}");
        return string.Join('&', pairs);
    }

    private Task<TResponse> GetAsync<TResponse>(string pathAndQuery, CancellationToken cancellationToken)
        => SendAsync<TResponse>(HttpMethod.Get, pathAndQuery, content: null, cancellationToken);

    private Task<TResponse> PostAsync<TRequest, TResponse>(
        string path, TRequest request, CancellationToken cancellationToken)
    {
        var json = JsonSerializer.Serialize(request, JsonOptions);
        var content = new StringContent(json, Encoding.UTF8, "application/json");
        return SendAsync<TResponse>(HttpMethod.Post, path, content, cancellationToken);
    }

    private async Task<TResponse> SendAsync<TResponse>(
        HttpMethod method, string pathAndQuery, HttpContent? content, CancellationToken cancellationToken)
    {
        var uri = new Uri(_options.Address, pathAndQuery);
        using var request = new HttpRequestMessage(method, uri) { Content = content };
        request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", _options.ApiToken);

        using var response = await _httpClient.SendAsync(request, cancellationToken).ConfigureAwait(false);
        var body = await response.Content.ReadAsStringAsync(cancellationToken).ConfigureAwait(false);

        if (!response.IsSuccessStatusCode)
        {
            throw BuildException((int)response.StatusCode, body);
        }

        return DeserializeSuccessBody<TResponse>((int)response.StatusCode, body);
    }

    private static TResponse DeserializeSuccessBody<TResponse>(int statusCode, string body)
    {
        TResponse? result;
        try
        {
            result = JsonSerializer.Deserialize<TResponse>(body, JsonOptions);
        }
        catch (JsonException ex)
        {
            throw new PgChangeFeedMalformedResponseException(
                statusCode,
                $"PG Change Feed HTTP API returned status {statusCode} with a response body that is not " +
                "valid JSON.",
                ex);
        }

        return result
            ?? throw new PgChangeFeedMalformedResponseException(
                statusCode,
                $"PG Change Feed HTTP API returned status {statusCode} with an empty or " +
                "null response body.");
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
            var error = JsonSerializer.Deserialize<ErrorResponse>(body, JsonOptions);
            return error?.Error ?? body;
        }
        catch (JsonException)
        {
            return body;
        }
    }
}
