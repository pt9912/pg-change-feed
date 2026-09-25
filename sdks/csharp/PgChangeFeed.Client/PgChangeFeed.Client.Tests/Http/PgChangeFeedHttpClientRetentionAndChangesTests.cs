using System.Net;
using System.Net.Http;
using PgChangeFeed.Client.Http.Models;
using Xunit;

namespace PgChangeFeed.Client.Tests.Http;

/// <summary>
/// Happy-path coverage for <c>RunRetention</c> and
/// <c>ReadChanges</c> (<c>GET /changes</c>) — including the
/// optional query parameters and the null-vs-embedded-JSON row image form.
/// </summary>
public class PgChangeFeedHttpClientRetentionAndChangesTests
{
    [Fact]
    public async Task RunRetentionAsync_HappyPath_ReturnsTypedResponse()
    {
        var (client, handler) = TestClientFactory.Create(request =>
        {
            Assert.Equal("/retention/run", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(HttpStatusCode.OK, """{"deleted":3}""");
        });

        var response = await client.RunRetentionAsync(new RunRetentionRequest("src", 1_000_000_000));

        Assert.Equal(3, response.Deleted);
        Assert.Contains("\"min_age_nanos\":1000000000", handler.LastRequestBody);
    }

    [Fact]
    public async Task ReadChangesAsync_HappyPath_ReturnsTypedResponseWithImages()
    {
        var (client, handler) = TestClientFactory.Create(request =>
        {
            Assert.Equal(HttpMethod.Get, request.Method);
            Assert.Equal("/changes", request.RequestUri!.AbsolutePath);
            Assert.Equal(
                "source=src&schema=public&table=orders&from=1&to=10&limit=5",
                request.RequestUri!.Query.TrimStart('?'));
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.OK,
                """
                {"changes":[{"commit_position":1,"change_id":"ch-1","transaction_id":"tx-1","source_table_id":"t-1","schema":"public","table":"orders","sequence":0,"operation":"INSERT","old_image":null,"new_image":{"id":1},"schema_version":"sv-1","committed_at":"2026-09-19T00:00:00Z"}]}
                """);
        });

        var response = await client.ReadChangesAsync(
            "src", schema: "public", table: "orders", from: 1, to: 10, limit: 5);

        Assert.Single(response.Changes);
        var change = response.Changes[0];
        Assert.Equal("ch-1", change.ChangeId);
        Assert.Equal("INSERT", change.Operation);
        Assert.Null(change.OldImage);
        Assert.NotNull(change.NewImage);
        Assert.Equal(1, change.NewImage!.Value.GetProperty("id").GetInt32());
        Assert.Null(handler.LastRequestBody);
    }

    // Body shape: the output of the server's GET /changes handler, where
    // `origin` is the last field of each change (Go test
    // TestReadChangesTraegtOriginAlsLetztesFeld,
    // internal/adapters/driving/http/readchanges_test.go). Not captured from
    // a running server.
    private static string ChangeBody(string originField) =>
        """{"changes":[{"commit_position":1,"change_id":"ch-1","transaction_id":"tx-1","source_table_id":"t-1","schema":"public","table":"orders","sequence":0,"operation":"INSERT","old_image":null,"new_image":{"id":1},"schema_version":"sv-1","committed_at":"2026-09-19T00:00:00Z"@ORIGIN@}]}"""
            .Replace("@ORIGIN@", originField);

    [Theory]
    [InlineData(",\"origin\":\"wal\"", "wal")]
    [InlineData(",\"origin\":\"backfill\"", "backfill")]
    [InlineData(",\"origin\":\"future-kind\"", "future-kind")]
    [InlineData("", "wal")]
    [InlineData(",\"origin\":null", "wal")]
    [InlineData(",\"origin\":\"\"", "")]
    public async Task ReadChangesAsync_Origin_ReadsTheServerValueAndDefaultsToWal(
        string originField, string expected)
    {
        var (client, _) = TestClientFactory.Create(_ =>
            FakeHttpMessageHandler.JsonResponse(HttpStatusCode.OK, ChangeBody(originField)));

        var response = await client.ReadChangesAsync("src");

        Assert.Equal(expected, response.Changes[0].Origin);
    }

    [Fact]
    public async Task ReadChangesAsync_NoOptionalParameters_OmitsThemFromQuery()
    {
        var (client, _) = TestClientFactory.Create(request =>
        {
            Assert.Equal("source=src", request.RequestUri!.Query.TrimStart('?'));
            return FakeHttpMessageHandler.JsonResponse(HttpStatusCode.OK, """{"changes":[]}""");
        });

        var response = await client.ReadChangesAsync("src");

        Assert.Empty(response.Changes);
    }
}
