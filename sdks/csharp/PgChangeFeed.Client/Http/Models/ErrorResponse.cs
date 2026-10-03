using System.Text.Json;
using System.Text.Json.Serialization;

namespace PgChangeFeed.Client.Http.Models;

/// <summary>
/// The error body of every non-success response
/// (<c>{"error": "&lt;text&gt;", "code": "&lt;message code&gt;"}</c>) — internal,
/// because it is only ever used to build a typed
/// <see cref="PgChangeFeedException"/>; it is never a return value of a public
/// method. <c>Code</c> stays a <see cref="JsonElement"/> so that a <c>code</c>
/// of another JSON type does not fail the read of the error text.
/// </summary>
internal sealed record ErrorResponse(
    [property: JsonPropertyName("error")] string Error,
    [property: JsonPropertyName("code")] JsonElement Code = default);
