using Xunit;

namespace CdcExamples.Http.Tests;

/// <summary>
/// Prüft den Fehlerpfad ohne erreichbaren Host — netzlos: die Verbindung zu
/// einem Loopback-Port ohne Listener scheitert sofort (Connection refused),
/// ohne echtes Netz zu brauchen. Dieselbe Grenze wie beim
/// <c>natsnotify</c>-Adapter (Verbindungsfehler bleiben auf Loopback, auch
/// unter <c>--network none</c>, siehe <c>harness/README.md</c> §Sensors,
/// <c>make test-notify</c>).
/// </summary>
public class TablesClientTests
{
    [Fact]
    public async Task ListTablesAsyncFailsOnUnreachableHost()
    {
        using var httpClient = new System.Net.Http.HttpClient { Timeout = TimeSpan.FromSeconds(2) };
        var cfg = new Config("127.0.0.1:1", "token", "quelle", "pub");

        await Assert.ThrowsAnyAsync<Exception>(() => TablesClient.ListTablesAsync(httpClient, cfg));
    }
}
