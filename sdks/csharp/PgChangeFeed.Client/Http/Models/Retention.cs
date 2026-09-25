using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Http.Models;

/// <summary>
/// The retention run to start: deletes the changes of <c>Source</c> that are
/// older than <c>MinAgeNanos</c> and that every consumer with a stored
/// position has passed. <c>Source</c> is mandatory, <c>MinAgeNanos</c> must be
/// &gt;= 0.
/// </summary>
public sealed record RunRetentionRequest(
    [property: JsonPropertyName("source")] string Source,
    [property: JsonPropertyName("min_age_nanos")] long MinAgeNanos);

/// <summary>
/// The number of changes actually deleted.
/// </summary>
public sealed record RunRetentionResponse(
    [property: JsonPropertyName("deleted")] int Deleted);
