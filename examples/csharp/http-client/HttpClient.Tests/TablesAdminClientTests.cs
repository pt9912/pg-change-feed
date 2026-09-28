using System.Net;
using Xunit;

namespace CdcExamples.Http.Tests;

/// <summary>
/// Prüft die drei Tabellen-Verwaltungs-Fähigkeiten — netzlos, über
/// <see cref="FakeHttpMessageHandler"/>. Form-Vorbild:
/// <c>examples/http-client/tables_admin_test.go</c>.
/// </summary>
public class TablesAdminClientTests
{
    private static Config BaseConfig() => new(
        Addr: "feed:8080", Token: "reader-token", AdminToken: "admin-token", Verb: "",
        Source: "", Publication: "", ConsumerId: "", Name: "", Offset: 0,
        Schema: "", Table: "", TableId: "", SchemaVersionId: "", Version: 1,
        From: "", To: "", Limit: "", MinAgeNanos: 0);

    [Fact]
    public void TableStatusUrlBuilderCarriesAllFourFields()
    {
        var got = TableStatusUrlBuilder.Build("feed:8080", "quelle-1", "public", "orders", "pub_quelle_1");
        Assert.Equal("http://feed:8080/tables/status?publication=pub_quelle_1&schema=public&source=quelle-1&table=orders", got);
    }

    [Fact]
    public async Task EnableTableAsyncParsesSuccessResponse()
    {
        var (httpClient, handler) = TestHttpClientFactory.Create(request =>
        {
            Assert.Equal("/tables/enable", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.Created,
                """{"table_id":"public.orders","source":"quelle-1","schema":"public","table":"orders","already_enabled":false}""");
        });

        var cfg = BaseConfig() with
        {
            Source = "quelle-1", Schema = "public", Table = "orders",
            TableId = "public.orders", SchemaVersionId = "public.orders-v1", Version = 1,
            Publication = "pub_quelle_1",
        };
        var resp = await TablesAdminClient.EnableTableAsync(httpClient, cfg);

        Assert.False(resp.AlreadyEnabled);
        Assert.Contains("\"table_id\":\"public.orders\"", handler.LastRequestBody);
        Assert.Contains("\"schema_version_id\":\"public.orders-v1\"", handler.LastRequestBody);
        Assert.Contains("\"version\":1", handler.LastRequestBody);
    }

    [Fact]
    public async Task EnableTableAsyncFailsOnNon2xx()
    {
        var (httpClient, _) = TestHttpClientFactory.Create(_ => FakeHttpMessageHandler.JsonResponse(
            HttpStatusCode.NotFound, """{"error":"Tabelle existiert an der Quelle nicht"}"""));

        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => TablesAdminClient.EnableTableAsync(httpClient, BaseConfig()));

        Assert.Contains("404", ex.Message);
    }

    [Fact]
    public async Task DisableTableAsyncParsesBothOutcomeFields()
    {
        var (httpClient, handler) = TestHttpClientFactory.Create(request =>
        {
            Assert.Equal("/tables/disable", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(HttpStatusCode.OK, """{"removed":true,"retained":true}""");
        });

        var cfg = BaseConfig() with { Source = "quelle-1", Schema = "public", Table = "orders", Publication = "pub_quelle_1" };
        var resp = await TablesAdminClient.DisableTableAsync(httpClient, cfg);

        Assert.True(resp.Removed);
        Assert.True(resp.Retained);
        Assert.Contains("\"publication\":\"pub_quelle_1\"", handler.LastRequestBody);
    }

    [Fact]
    public async Task DisableTableAsyncFailsOnNon2xx()
    {
        var (httpClient, _) = TestHttpClientFactory.Create(_ => FakeHttpMessageHandler.JsonResponse(
            HttpStatusCode.BadRequest, """{"error":"source, schema, table und publication sind Pflichtfelder"}"""));

        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => TablesAdminClient.DisableTableAsync(httpClient, BaseConfig()));

        Assert.Contains("400", ex.Message);
        Assert.Contains("Pflichtfelder", ex.Message);
    }

    [Fact]
    public async Task TableStatusAsyncParsesNeverEnabledBoundary()
    {
        var (httpClient, handler) = TestHttpClientFactory.Create(request =>
        {
            Assert.Equal(HttpMethod.Get, request.Method);
            return FakeHttpMessageHandler.JsonResponse(HttpStatusCode.OK, """{"enabled":false,"retained":false}""");
        });

        var cfg = BaseConfig() with { Source = "quelle-1", Schema = "public", Table = "nie_aktiviert", Publication = "pub_quelle_1" };
        var resp = await TablesAdminClient.TableStatusAsync(httpClient, cfg);

        Assert.False(resp.Enabled);
        Assert.False(resp.Retained);
        Assert.Null(handler.LastRequestBody);
    }

    [Fact]
    public async Task TableStatusAsyncFailsOnNon2xx()
    {
        var (httpClient, _) = TestHttpClientFactory.Create(_ => FakeHttpMessageHandler.JsonResponse(
            HttpStatusCode.BadRequest, """{"error":"source, schema, table und publication sind Pflichtfelder"}"""));

        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => TablesAdminClient.TableStatusAsync(httpClient, BaseConfig()));

        Assert.Contains("400", ex.Message);
        Assert.Contains("Pflichtfelder", ex.Message);
    }
}
