using System.Net;
using System.Net.Http.Headers;

namespace CdcExamples.Nats;

/// <summary>
/// ChangesClient holt die Änderungen eines Quelle/Schema/Tabelle-Filters mit
/// dem <c>reader</c>-Token ab (<c>LH-FA-SST-006</c>, zweiseitiger Ablauf
/// <c>ADR-0079</c>) — Form-Vorbild: <c>examples/nats-client/main.go</c>
/// (<c>fetchChanges</c>). Ein Nicht-200-Status ist ein sichtbarer Fehler,
/// kein stiller Leerwert. Der <see cref="System.Net.Http.HttpClient"/> ist
/// injiziert, damit der Fehlerpfad ohne erreichbaren Host netzlos testbar ist.
/// </summary>
public static class ChangesClient
{
    public static async Task<string> FetchAsync(
        System.Net.Http.HttpClient httpClient,
        Config cfg,
        CancellationToken cancellationToken = default)
    {
        var url = ChangesUrlBuilder.Build(cfg.Addr, cfg.Source, cfg.Schema, cfg.Table);
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
