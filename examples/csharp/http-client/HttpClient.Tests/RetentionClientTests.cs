using System.Net;
using Xunit;

namespace CdcExamples.Http.Tests;

/// <summary>
/// Prüft <c>POST /retention/run</c> — netzlos, über
/// <see cref="FakeHttpMessageHandler"/>. Form-Vorbild:
/// <c>examples/http-client/retention_test.go</c>.
/// </summary>
public class RetentionClientTests
{
    private static Config BaseConfig() => new(
        Addr: "feed:8080", Token: "", AdminToken: "admin-token", Verb: "",
        Source: "", Publication: "", ConsumerId: "", Name: "", Offset: 0,
        Schema: "", Table: "", TableId: "", SchemaVersionId: "", Version: 1,
        From: "", To: "", Limit: "", MinAgeNanos: 0);

    [Fact]
    public async Task RunRetentionAsyncParsesZeroDeletedAsValidOutcome()
    {
        var (httpClient, handler) = TestHttpClientFactory.Create(request =>
        {
            Assert.Equal("/retention/run", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(HttpStatusCode.OK, """{"deleted":0}""");
        });

        var cfg = BaseConfig() with { Source = "quelle-1", MinAgeNanos = 3_600_000_000_000 };
        var resp = await RetentionClient.RunRetentionAsync(httpClient, cfg);

        Assert.Equal(0, resp.Deleted);
        Assert.Contains("\"min_age_nanos\":3600000000000", handler.LastRequestBody);
    }

    [Fact]
    public async Task RunRetentionAsyncFailsOnNon2xx()
    {
        var (httpClient, _) = TestHttpClientFactory.Create(_ => FakeHttpMessageHandler.JsonResponse(
            HttpStatusCode.BadRequest, """{"error":"min_age_nanos darf nicht negativ sein"}"""));

        var cfg = BaseConfig() with { Source = "quelle-1", MinAgeNanos = -1 };
        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => RetentionClient.RunRetentionAsync(httpClient, cfg));

        Assert.Contains("400", ex.Message);
    }
}
