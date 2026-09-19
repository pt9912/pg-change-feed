using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Http.Models;

/// <summary>
/// The uniform error body of every non-success response
/// (<c>{"error": "&lt;text&gt;"}</c>, SPEC-018) — internal, because it is
/// only ever used to build a typed <see cref="PgChangeFeedException"/>; it is
/// never a return value of a public method.
/// </summary>
internal sealed record ErrorResponse(
    [property: JsonPropertyName("error")] string Error);
