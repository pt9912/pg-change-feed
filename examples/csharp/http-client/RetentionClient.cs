using System.Net;
using System.Text.Json.Serialization;

namespace CdcExamples.Http;

/// <summary>
/// RunRetentionRequest/-Response spiegeln <c>POST /retention/run</c>
/// (<c>LH-FA-RET-002</c>): <c>MinAgeNanos</c> 0 heißt „kein zeitliches
/// Mindestalter" und ist gültig. Form-Vorbild:
/// <c>examples/http-client/retention.go</c>.
/// </summary>
public sealed record RunRetentionRequest(
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("min_age_nanos")] long MinAgeNanos);

public sealed record RunRetentionResponse(
    [property: JsonPropertyName("deleted")] int Deleted);

/// <summary>
/// RetentionClient ruft den <c>admin</c>-Endpunkt <c>POST /retention/run</c>
/// auf. Form-Vorbild: <c>examples/http-client/retention.go</c>.
/// </summary>
public static class RetentionClient
{
    public static Task<RunRetentionResponse> RunRetentionAsync(
        System.Net.Http.HttpClient httpClient, Config cfg, CancellationToken cancellationToken = default) =>
        RequestHelper.SendJsonAsync<RunRetentionResponse>(
            httpClient, HttpMethod.Post, $"http://{cfg.Addr}/retention/run", cfg.AdminToken,
            new RunRetentionRequest(cfg.Source, cfg.MinAgeNanos), HttpStatusCode.OK, cancellationToken);
}
