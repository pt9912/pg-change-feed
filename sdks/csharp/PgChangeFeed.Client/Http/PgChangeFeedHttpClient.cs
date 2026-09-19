using System.Net.Http;
using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;
using PgChangeFeed.Client.Http.Models;

namespace PgChangeFeed.Client.Http;

/// <summary>
/// Public entry point for the PG Change Feed HTTP/JSON API: one method per
/// wire capability — the nine port-covered capabilities of SPEC-018
/// (<c>RegisterConsumer</c>, <c>AcknowledgeConsumer</c>,
/// <c>GetConsumerPosition</c>, <c>RemoveConsumer</c>, <c>EnableTable</c>,
/// <c>DisableTable</c>, <c>GetStatus</c>, <c>ListTables</c>,
/// <c>RunRetention</c>) plus reading persisted changes (<c>ReadChanges</c>,
/// SPEC-022). Requests/responses are typed DTOs that mirror the SPEC-018/
/// SPEC-022 JSON schemas exactly (<c>PgChangeFeed.Client.Http.Models</c>);
/// every non-success response becomes a typed <see cref="PgChangeFeedException"/>
/// subtype instead of a result type mixed with the success path.
///
/// The <see cref="HttpClient"/> is injected, not owned — the caller controls
/// its lifetime, connection pooling and any <c>DelegatingHandler</c>
/// pipeline (proxies, retries, logging); this type never disposes it. The
/// bearer token and server address come from <see cref="PgChangeFeedClientOptions"/>,
/// supplied at construction — no global or static state, a process can hold
/// several independently configured instances at once.
///
/// Draht-Kenntnis-Vorbild (gelesen, nicht importiert — ADR-0106 Festlegung 2):
/// <c>examples/csharp/http-client/TablesClient.cs</c>,
/// <c>TablesUrlBuilder.cs</c>.
/// </summary>
public sealed class PgChangeFeedHttpClient
{
    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNameCaseInsensitive = true,
    };

    private readonly HttpClient _httpClient;
    private readonly PgChangeFeedClientOptions _options;

    public PgChangeFeedHttpClient(HttpClient httpClient, PgChangeFeedClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(httpClient);
        ArgumentNullException.ThrowIfNull(options);
        _httpClient = httpClient;
        _options = options;
    }

    /// <summary><c>RegisterConsumer</c> — <c>POST /consumers</c> (admin, SPEC-018).</summary>
    public Task<RegisterConsumerResponse> RegisterConsumerAsync(
        RegisterConsumerRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return PostAsync<RegisterConsumerRequest, RegisterConsumerResponse>(
            "/consumers", request, cancellationToken);
    }

    /// <summary><c>AcknowledgeConsumer</c> — <c>POST /consumers/acknowledge</c> (admin, SPEC-018).</summary>
    public Task<AcknowledgeConsumerResponse> AcknowledgeConsumerAsync(
        AcknowledgeConsumerRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return PostAsync<AcknowledgeConsumerRequest, AcknowledgeConsumerResponse>(
            "/consumers/acknowledge", request, cancellationToken);
    }

    /// <summary><c>GetConsumerPosition</c> — <c>GET /consumers/position</c> (reader or admin, SPEC-018).</summary>
    public Task<ConsumerPositionResponse> GetConsumerPositionAsync(
        string consumerId, CancellationToken cancellationToken = default)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(consumerId);
        var query = BuildQuery(("consumer_id", consumerId));
        return GetAsync<ConsumerPositionResponse>($"/consumers/position?{query}", cancellationToken);
    }

    /// <summary><c>RemoveConsumer</c> — <c>POST /consumers/remove</c> (admin, SPEC-018).</summary>
    public Task<RemoveConsumerResponse> RemoveConsumerAsync(
        string consumerId, CancellationToken cancellationToken = default)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(consumerId);
        return PostAsync<RemoveConsumerRequest, RemoveConsumerResponse>(
            "/consumers/remove", new RemoveConsumerRequest(consumerId), cancellationToken);
    }

    /// <summary><c>EnableTable</c> — <c>POST /tables/enable</c> (admin, SPEC-018).</summary>
    public Task<EnableTableResponse> EnableTableAsync(
        EnableTableRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return PostAsync<EnableTableRequest, EnableTableResponse>(
            "/tables/enable", request, cancellationToken);
    }

    /// <summary><c>DisableTable</c> — <c>POST /tables/disable</c> (admin, SPEC-018).</summary>
    public Task<DisableTableResponse> DisableTableAsync(
        DisableTableRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return PostAsync<DisableTableRequest, DisableTableResponse>(
            "/tables/disable", request, cancellationToken);
    }

    /// <summary><c>GetStatus</c> — <c>GET /tables/status</c> (reader or admin, SPEC-018).</summary>
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

    /// <summary><c>ListTables</c> — <c>GET /tables</c> (reader or admin, SPEC-018).</summary>
    public Task<ListTablesResponse> ListTablesAsync(
        string source, string publication, CancellationToken cancellationToken = default)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(source);
        ArgumentException.ThrowIfNullOrWhiteSpace(publication);
        var query = BuildQuery(("source", source), ("publication", publication));
        return GetAsync<ListTablesResponse>($"/tables?{query}", cancellationToken);
    }

    /// <summary><c>RunRetention</c> — <c>POST /retention/run</c> (admin, SPEC-018).</summary>
    public Task<RunRetentionResponse> RunRetentionAsync(
        RunRetentionRequest request, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        return PostAsync<RunRetentionRequest, RunRetentionResponse>(
            "/retention/run", request, cancellationToken);
    }

    /// <summary>
    /// <c>ReadChanges</c> — <c>GET /changes</c> (reader or admin, SPEC-022).
    /// <paramref name="source"/> is mandatory; <paramref name="schema"/>/
    /// <paramref name="table"/> are optional and independent;
    /// <paramref name="from"/> is inclusive, <paramref name="to"/> exclusive
    /// (both <c>commit_position</c> values); there is no default
    /// <paramref name="limit"/>.
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
                "valid JSON — a protocol violation outside SPEC-018/SPEC-022's documented shapes.",
                ex);
        }

        return result
            ?? throw new PgChangeFeedMalformedResponseException(
                statusCode,
                $"PG Change Feed HTTP API returned status {statusCode} with an empty or " +
                "null response body — a protocol violation outside SPEC-018/SPEC-022's documented shapes.");
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
