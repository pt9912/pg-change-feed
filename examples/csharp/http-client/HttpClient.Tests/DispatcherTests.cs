using Xunit;

namespace CdcExamples.Http.Tests;

/// <summary>
/// Prüft den defensiven Schutz von <see cref="Dispatcher.DispatchAsync"/>
/// gegen ein unbekanntes Verb — für den Fall eines Aufrufs ohne vorherige
/// <see cref="Validator.Validate"/>. Form-Vorbild:
/// <c>examples/http-client/main_test.go</c>,
/// <c>TestDispatchRejectsUnknownVerb</c>.
/// </summary>
public class DispatcherTests
{
    [Fact]
    public async Task DispatchAsyncRejectsUnknownVerbBeforeAnyNetworkCall()
    {
        using var httpClient = new System.Net.Http.HttpClient();
        var cfg = new Config(
            Addr: "feed:8080", Token: "", AdminToken: "", Verb: "unbekannt",
            Source: "", Publication: "", ConsumerId: "", Name: "", Offset: 0,
            Schema: "", Table: "", TableId: "", SchemaVersionId: "", Version: 1,
            From: "", To: "", Limit: "", MinAgeNanos: 0);

        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => Dispatcher.DispatchAsync(httpClient, cfg));

        Assert.Contains("unbekanntes --verb", ex.Message);
    }
}
