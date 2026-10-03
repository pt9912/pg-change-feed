using System.Text.Json;
using PgChangeFeed.Client.Http.Models;

namespace PgChangeFeed.Client.Http;

/// <summary>
/// Reads the error text and the message code out of the body of a
/// non-success response; shared by the HTTP and the SSE client.
/// </summary>
internal static class ErrorBody
{
    /// <summary>
    /// Returns the error text (the raw body when <c>error</c> is not a JSON
    /// string) and the message code; only a non-empty JSON string in
    /// <c>code</c> counts as a code, whatever the type of <c>error</c>.
    /// Never throws.
    /// </summary>
    internal static (string Message, string? MessageCode) Parse(string body, JsonSerializerOptions options)
    {
        try
        {
            var error = JsonSerializer.Deserialize<ErrorResponse>(body, options);
            if (error is null)
            {
                return (body, null);
            }

            var code = error.Code.ValueKind == JsonValueKind.String ? error.Code.GetString() : null;
            var text = error.Error.ValueKind == JsonValueKind.String ? error.Error.GetString() : null;
            return (text ?? body, string.IsNullOrEmpty(code) ? null : code);
        }
        catch (JsonException)
        {
            return (body, null);
        }
    }
}
