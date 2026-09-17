using System.Net;
using System.Net.Http.Headers;

namespace CdcExamples.Http;

/// <summary>
/// TablesClient ruft die Tabellen-Auflistung mit dem <c>reader</c>-Token ab
/// (<c>LH-FA-SST-006</c>, <c>ADR-0057</c>) — Form-Vorbild:
/// <c>examples/http-client/main.go</c>, <c>listTables</c>. Ein Nicht-200-
/// Status ist ein sichtbarer Fehler, kein stiller Leerwert. Der
/// <see cref="System.Net.Http.HttpClient"/> ist injiziert, damit der
/// Fehlerpfad ohne erreichbaren Host netzlos testbar ist.
/// </summary>
public static class TablesClient
{
    public static async Task<string> ListTablesAsync(
        System.Net.Http.HttpClient httpClient,
        Config cfg,
        CancellationToken cancellationToken = default)
    {
        var url = TablesUrlBuilder.Build(cfg.Addr, cfg.Source, cfg.Publication);
        using var request = new HttpRequestMessage(HttpMethod.Get, url);
        request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", cfg.Token);

        using var response = await httpClient.SendAsync(request, cancellationToken).ConfigureAwait(false);
        var body = await response.Content.ReadAsStringAsync(cancellationToken).ConfigureAwait(false);
        if (response.StatusCode != HttpStatusCode.OK)
        {
            throw new InvalidOperationException($"HTTP-Status {(int)response.StatusCode}: {body}");
        }
        return body;
    }
}
