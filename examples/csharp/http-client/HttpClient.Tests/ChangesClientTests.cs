using System.Net;
using Xunit;

namespace CdcExamples.Http.Tests;

/// <summary>
/// Prüft <c>GET /changes</c> — netzlos, über
/// <see cref="FakeHttpMessageHandler"/>. Form-Vorbild:
/// <c>examples/http-client/changes_test.go</c>.
/// </summary>
public class ChangesClientTests
{
    private static Config BaseConfig() => new(
        Addr: "feed:8080", Token: "reader-token", AdminToken: "", Verb: "",
        Source: "", Publication: "", ConsumerId: "", Name: "", Offset: 0,
        Schema: "", Table: "", TableId: "", SchemaVersionId: "", Version: 1,
        From: "", To: "", Limit: "", MinAgeNanos: 0);

    [Fact]
    public void ChangesUrlBuilderCarriesOnlyRequiredField()
    {
        var got = ChangesUrlBuilder.Build("feed:8080", "quelle-1", "", "", "", "", "");
        Assert.Equal("http://feed:8080/changes?source=quelle-1", got);
    }

    [Fact]
    public void ChangesUrlBuilderCarriesAllOptionalFieldsIndependently()
    {
        var got = ChangesUrlBuilder.Build("feed:8080", "quelle-1", "public", "orders", "10", "20", "5");
        Assert.Equal("http://feed:8080/changes?from=10&limit=5&schema=public&source=quelle-1&table=orders&to=20", got);
    }

    [Fact]
    public void ChangesUrlBuilderOmitsEmptyOptionalValues()
    {
        var got = ChangesUrlBuilder.Build("feed:8080", "quelle-1", "", "orders", "", "20", "");
        Assert.Equal("http://feed:8080/changes?source=quelle-1&table=orders&to=20", got);
    }

    [Fact]
    public async Task ReadChangesAsyncParsesSuccessResponseWithEmbeddedImages()
    {
        var (httpClient, handler) = TestHttpClientFactory.Create(request =>
        {
            Assert.Equal(HttpMethod.Get, request.Method);
            Assert.Equal("source=quelle-1", request.RequestUri!.Query.TrimStart('?'));
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.OK,
                """
                {"changes":[{"commit_position":42,"change_id":"change-1","transaction_id":"tx-1","source_table_id":"t-1","schema":"public","table":"orders","sequence":0,"operation":"INSERT","old_image":null,"new_image":{"name":"erste Zeile"},"schema_version":"sv-1","committed_at":"2026-09-19T00:00:00Z","origin":"wal"}]}
                """);
        });

        var cfg = BaseConfig() with { Source = "quelle-1" };
        var resp = await ChangesClient.ReadChangesAsync(httpClient, cfg);

        Assert.Single(resp.Changes);
        var change = resp.Changes[0];
        Assert.Equal("change-1", change.ChangeId);
        Assert.Equal("wal", change.Origin);
        Assert.Null(change.OldImage);
        Assert.NotNull(change.NewImage);
        Assert.Equal("erste Zeile", change.NewImage!.Value.GetProperty("name").GetString());
        Assert.Null(handler.LastRequestBody);
    }

    [Fact]
    public async Task ReadChangesAsyncFailsOnNon2xx()
    {
        var (httpClient, _) = TestHttpClientFactory.Create(_ => FakeHttpMessageHandler.JsonResponse(
            HttpStatusCode.BadRequest, """{"error":"source ist Pflichtfeld"}"""));

        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => ChangesClient.ReadChangesAsync(httpClient, BaseConfig()));

        Assert.Contains("400", ex.Message);
    }
}
