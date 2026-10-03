using System.Text.Json;
using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Http.Models;

/// <summary>
/// The error body of every non-success response
/// (<c>{"error": "&lt;text&gt;", "code": "&lt;message code&gt;"}</c>) — internal,
/// because it is only ever used to build a typed
/// <see cref="PgChangeFeedException"/>; it is never a return value of a public
/// method. <c>Error</c> and <c>Code</c> stay <see cref="JsonElement"/>s so that a
/// field of another JSON type fails neither the read of the other field nor of the body.
/// </summary>
internal sealed record ErrorResponse(
    [property: JsonPropertyName("error")] JsonElement Error = default,
    [property: JsonPropertyName("code")] JsonElement Code = default);
