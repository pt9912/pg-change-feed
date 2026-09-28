using System.Net;
using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;

namespace CdcExamples.Http;

/// <summary>
/// RequestHelper trägt den gemeinsamen Anfrage/Antwort-Ablauf aller zehn
/// Fähigkeiten dieses Beispiels: JSON-Body kodieren (falls <paramref
/// name="requestBody"/> gesetzt ist), das Bearer-Token setzen, senden, und
/// bei <paramref name="wantStatus"/> den Antwort-Body typisiert dekodieren.
/// Ein davon abweichender Statuscode ist ein sichtbarer Fehler mit
/// Statuscode und Antworttext, kein stiller Leerwert — dieselbe Form wie
/// <c>TablesClient.ListTablesAsync</c>. Form-Vorbild:
/// <c>doRequestJSON</c> in <c>examples/http-client/request.go</c>.
/// </summary>
public static class RequestHelper
{
    public static async Task<TResponse> SendJsonAsync<TResponse>(
        System.Net.Http.HttpClient httpClient,
        HttpMethod method,
        string url,
        string token,
        object? requestBody,
        HttpStatusCode wantStatus,
        CancellationToken cancellationToken = default)
    {
        using var request = new HttpRequestMessage(method, url);
        if (requestBody is not null)
        {
            var json = JsonSerializer.Serialize(requestBody);
            request.Content = new StringContent(json, Encoding.UTF8, "application/json");
        }
        request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", token);

        using var response = await httpClient.SendAsync(request, cancellationToken).ConfigureAwait(false);
        var body = await response.Content.ReadAsStringAsync(cancellationToken).ConfigureAwait(false);
        if (response.StatusCode != wantStatus)
        {
            throw new InvalidOperationException($"HTTP-Status {(int)response.StatusCode}: {body}");
        }
        return JsonSerializer.Deserialize<TResponse>(body)
            ?? throw new InvalidOperationException("HTTP-Antwort ist leer oder \"null\"");
    }
}
