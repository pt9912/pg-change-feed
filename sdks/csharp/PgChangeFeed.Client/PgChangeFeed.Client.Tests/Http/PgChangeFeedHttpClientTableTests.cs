using System.Net;
using System.Net.Http;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Http.Models;
using Xunit;

namespace PgChangeFeed.Client.Tests.Http;

/// <summary>
/// Happy-path (and the table-only <c>404</c>) coverage for the four
/// table-management capabilities (<c>EnableTable</c>,
/// <c>DisableTable</c>, <c>GetStatus</c>, <c>ListTables</c>).
/// </summary>
public class PgChangeFeedHttpClientTableTests
{
    private static EnableTableRequest SampleEnableRequest()
        => new("src", "public", "orders", "t-1", "sv-1", 1, "pub");

    [Fact]
    public async Task EnableTableAsync_HappyPath_ReturnsTypedResponse()
    {
        var (client, handler) = TestClientFactory.Create(request =>
        {
            Assert.Equal("/tables/enable", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.Created,
                """{"table_id":"t-1","source":"src","schema":"public","table":"orders","already_enabled":false}""");
        });

        var response = await client.EnableTableAsync(SampleEnableRequest());

        Assert.Equal("t-1", response.TableId);
        Assert.False(response.AlreadyEnabled);
        Assert.Contains("\"version\":1", handler.LastRequestBody);
    }

    [Fact]
    public async Task EnableTableAsync_TableMissingAtSource_ThrowsNotFound()
    {
        var (client, _) = TestClientFactory.Create(_ => FakeHttpMessageHandler.JsonResponse(
            HttpStatusCode.NotFound, """{"error":"table does not exist at source"}"""));

        var ex = await Assert.ThrowsAsync<PgChangeFeedNotFoundException>(
            () => client.EnableTableAsync(SampleEnableRequest()));

        Assert.Equal(404, ex.StatusCode);
        Assert.Equal("table does not exist at source", ex.Message);
    }

    [Fact]
    public async Task DisableTableAsync_HappyPath_ReturnsTypedResponse()
    {
        var (client, handler) = TestClientFactory.Create(request =>
        {
            Assert.Equal("/tables/disable", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(HttpStatusCode.OK, """{"removed":true,"retained":false}""");
        });

        var response = await client.DisableTableAsync(new DisableTableRequest("src", "public", "orders", "pub"));

        Assert.True(response.Removed);
        Assert.False(response.Retained);
        Assert.Contains("\"publication\":\"pub\"", handler.LastRequestBody);
    }

    [Fact]
    public async Task DisableTableAsync_TableMissingAtSource_ThrowsNotFound()
    {
        var (client, _) = TestClientFactory.Create(_ => FakeHttpMessageHandler.JsonResponse(
            HttpStatusCode.NotFound, """{"error":"table does not exist at source"}"""));

        var ex = await Assert.ThrowsAsync<PgChangeFeedNotFoundException>(
            () => client.DisableTableAsync(new DisableTableRequest("src", "public", "orders", "pub")));

        Assert.Equal(404, ex.StatusCode);
    }

    [Fact]
    public async Task GetStatusAsync_HappyPath_BuildsQueryAndReturnsTypedResponse()
    {
        var (client, handler) = TestClientFactory.Create(request =>
        {
            Assert.Equal("/tables/status", request.RequestUri!.AbsolutePath);
            Assert.Equal(
                "source=src&schema=public&table=orders&publication=pub",
                request.RequestUri!.Query.TrimStart('?'));
            return FakeHttpMessageHandler.JsonResponse(HttpStatusCode.OK, """{"enabled":true,"retained":false}""");
        });

        var response = await client.GetStatusAsync("src", "public", "orders", "pub");

        Assert.True(response.Enabled);
        Assert.False(response.Retained);
        Assert.Null(handler.LastRequestBody);
    }

    [Fact]
    public async Task ListTablesAsync_HappyPath_ReturnsTypedResponse()
    {
        var (client, handler) = TestClientFactory.Create(request =>
        {
            Assert.Equal("/tables", request.RequestUri!.AbsolutePath);
            Assert.Equal("source=src&publication=pub", request.RequestUri!.Query.TrimStart('?'));
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.OK,
                """{"tables":[{"table_id":"t-1","source":"src","schema":"public","table":"orders"}],"retained":[]}""");
        });

        var response = await client.ListTablesAsync("src", "pub");

        Assert.Single(response.Tables);
        Assert.Equal("t-1", response.Tables[0].TableId);
        Assert.Empty(response.Retained);
        Assert.Null(handler.LastRequestBody);
    }
}
