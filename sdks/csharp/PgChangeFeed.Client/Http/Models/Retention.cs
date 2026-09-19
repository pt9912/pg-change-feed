using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Http.Models;

/// <summary>
/// <c>RunRetention</c> request — <c>POST /retention/run</c> (SPEC-018);
/// <c>Source</c> mandatory, <c>MinAgeNanos</c> must be &gt;= 0.
/// </summary>
public sealed record RunRetentionRequest(
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("min_age_nanos")] long MinAgeNanos);

/// <summary>
/// <c>RunRetention</c> response (<c>200</c>) — the number of changes
/// actually deleted (SPEC-018).
/// </summary>
public sealed record RunRetentionResponse(
    [property: JsonPropertyName("deleted")] int Deleted);
