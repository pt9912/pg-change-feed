using System.Net;
using Xunit;

namespace CdcExamples.Http.Tests;

/// <summary>
/// Prüft die vier Consumer-Verwaltungs-Fähigkeiten — netzlos, über
/// <see cref="FakeHttpMessageHandler"/>. Form-Vorbild:
/// <c>examples/http-client/consumer_test.go</c>.
/// </summary>
public class ConsumerClientTests
{
    private static Config BaseConfig() => new(
        Addr: "feed:8080", Token: "reader-token", AdminToken: "admin-token", Verb: "",
        Source: "", Publication: "", ConsumerId: "", Name: "", Offset: 0,
        Schema: "", Table: "", TableId: "", SchemaVersionId: "", Version: 1,
        From: "", To: "", Limit: "", MinAgeNanos: 0);

    [Fact]
    public void ConsumerPositionUrlBuilderCarriesConsumerId()
    {
        var got = ConsumerPositionUrlBuilder.Build("feed:8080", "consumer-1");
        Assert.Equal("http://feed:8080/consumers/position?consumer_id=consumer-1", got);
    }

    [Fact]
    public void ConsumerPositionUrlBuilderEscapesReservedCharacters()
    {
        var got = ConsumerPositionUrlBuilder.Build("feed:8080", "consumer & 1");
        Assert.Equal("http://feed:8080/consumers/position?consumer_id=consumer%20%26%201", got);
    }

    [Fact]
    public async Task RegisterConsumerAsyncParsesSuccessResponse()
    {
        var (httpClient, handler) = TestHttpClientFactory.Create(request =>
        {
            Assert.Equal(HttpMethod.Post, request.Method);
            Assert.Equal("/consumers", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.Created,
                """{"consumer_id":"consumer-1","name":"Consumer Eins","already_registered":false}""");
        });

        var cfg = BaseConfig() with { ConsumerId = "consumer-1", Name = "Consumer Eins" };
        var resp = await ConsumerClient.RegisterConsumerAsync(httpClient, cfg);

        Assert.Equal("consumer-1", resp.ConsumerId);
        Assert.False(resp.AlreadyRegistered);
        Assert.Contains("\"consumer_id\":\"consumer-1\"", handler.LastRequestBody);
    }

    [Fact]
    public async Task RegisterConsumerAsyncFailsOnNon2xx()
    {
        var (httpClient, _) = TestHttpClientFactory.Create(_ => FakeHttpMessageHandler.JsonResponse(
            HttpStatusCode.BadRequest, """{"error":"consumer_id und name sind Pflichtfelder"}"""));

        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => ConsumerClient.RegisterConsumerAsync(httpClient, BaseConfig()));

        Assert.Contains("400", ex.Message);
        Assert.Contains("Pflichtfelder", ex.Message);
    }

    [Fact]
    public async Task AcknowledgeConsumerAsyncParsesSuccessResponse()
    {
        var (httpClient, handler) = TestHttpClientFactory.Create(request =>
        {
            Assert.Equal("/consumers/acknowledge", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.OK,
                """{"consumer_id":"consumer-1","source_id":"quelle-1","offset":42}""");
        });

        var cfg = BaseConfig() with { ConsumerId = "consumer-1", Source = "quelle-1", Offset = 42 };
        var resp = await ConsumerClient.AcknowledgeConsumerAsync(httpClient, cfg);

        Assert.Equal(42UL, resp.Offset);
        Assert.Contains("\"offset\":42", handler.LastRequestBody);
    }

    [Fact]
    public async Task AcknowledgeConsumerAsyncFailsOnNon2xx()
    {
        var (httpClient, _) = TestHttpClientFactory.Create(_ => FakeHttpMessageHandler.JsonResponse(
            HttpStatusCode.BadRequest, """{"error":"consumer_id, source und offset sind Pflichtfelder"}"""));

        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => ConsumerClient.AcknowledgeConsumerAsync(httpClient, BaseConfig()));

        Assert.Contains("400", ex.Message);
        Assert.Contains("Pflichtfelder", ex.Message);
    }

    [Fact]
    public async Task ConsumerPositionAsyncParsesUnacknowledgedBoundary()
    {
        var (httpClient, handler) = TestHttpClientFactory.Create(request =>
        {
            Assert.Equal(HttpMethod.Get, request.Method);
            Assert.Equal("consumer_id=consumer-1", request.RequestUri!.Query.TrimStart('?'));
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.OK,
                """{"consumer_id":"consumer-1","source_id":"","offset":0,"acknowledged":false}""");
        });

        var cfg = BaseConfig() with { ConsumerId = "consumer-1" };
        var resp = await ConsumerClient.ConsumerPositionAsync(httpClient, cfg);

        Assert.False(resp.Acknowledged);
        Assert.Null(handler.LastRequestBody);
    }

    [Fact]
    public async Task ConsumerPositionAsyncFailsOnNon2xx()
    {
        var (httpClient, _) = TestHttpClientFactory.Create(_ => FakeHttpMessageHandler.JsonResponse(
            HttpStatusCode.BadRequest, """{"error":"consumer_id ist Pflichtfeld"}"""));

        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => ConsumerClient.ConsumerPositionAsync(httpClient, BaseConfig()));

        Assert.Contains("400", ex.Message);
        Assert.Contains("Pflichtfeld", ex.Message);
    }

    [Fact]
    public async Task RemoveConsumerAsyncParsesIdempotentOutcome()
    {
        var (httpClient, handler) = TestHttpClientFactory.Create(request =>
        {
            Assert.Equal("/consumers/remove", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.OK,
                """{"consumer_id":"nie-registriert","removed":false}""");
        });

        var cfg = BaseConfig() with { ConsumerId = "nie-registriert" };
        var resp = await ConsumerClient.RemoveConsumerAsync(httpClient, cfg);

        Assert.False(resp.Removed);
        Assert.Contains("\"consumer_id\":\"nie-registriert\"", handler.LastRequestBody);
    }

    [Fact]
    public async Task RemoveConsumerAsyncFailsOnNon2xx()
    {
        var (httpClient, _) = TestHttpClientFactory.Create(_ => FakeHttpMessageHandler.JsonResponse(
            HttpStatusCode.BadRequest, """{"error":"consumer_id ist Pflichtfeld"}"""));

        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => ConsumerClient.RemoveConsumerAsync(httpClient, BaseConfig()));

        Assert.Contains("400", ex.Message);
        Assert.Contains("Pflichtfeld", ex.Message);
    }
}
