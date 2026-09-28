using Xunit;

namespace CdcExamples.Grpc.Tests;

/// <summary>
/// Prüft den defensiven Schutz von <see cref="Dispatcher.DispatchAdminAsync"/>
/// gegen ein unbekanntes Verb — für den Fall eines Aufrufs ohne vorherige
/// <see cref="Validator.Validate"/>. Ein <c>null</c>-Client wird dabei nie
/// erreicht. Form-Vorbild: <c>examples/grpc-client/main_test.go</c>,
/// <c>TestDispatchAdminRejectsUnknownVerb</c>.
/// </summary>
public class DispatcherTests
{
    private static Config BaseConfig(string verb) => new(
        Addr: "feed:9090", Token: "", AdminToken: "", Verb: verb,
        Schema: "", Table: "", ConsumerId: "", Name: "", Offset: 0,
        TableId: "", SchemaVersionId: "", Version: 1,
        Source: "", Publication: "", From: 0, To: 0, Limit: 0, MinAgeNanos: 0);

    [Fact]
    public async Task DispatchAdminAsyncRejectsUnknownVerbBeforeAnyNetworkCall()
    {
        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => Dispatcher.DispatchAdminAsync(client: null!, BaseConfig("unbekannt")));

        Assert.Contains("unbekanntes --verb", ex.Message);
    }
}
