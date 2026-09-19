using System.Net;
using System.Net.Http;
using PgChangeFeed.Client.Http.Models;
using Xunit;

namespace PgChangeFeed.Client.Tests.Http;

/// <summary>
/// Happy-path coverage for the four consumer-management capabilities of
/// SPEC-018 (<c>RegisterConsumer</c>, <c>AcknowledgeConsumer</c>,
/// <c>GetConsumerPosition</c>, <c>RemoveConsumer</c>).
/// </summary>
public class PgChangeFeedHttpClientConsumerTests
{
    [Fact]
    public async Task RegisterConsumerAsync_HappyPath_ReturnsTypedResponse()
    {
        var (client, handler) = TestClientFactory.Create(request =>
        {
            Assert.Equal(HttpMethod.Post, request.Method);
            Assert.Equal("/consumers", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.Created,
                """{"consumer_id":"c-1","name":"n-1","already_registered":false}""");
        });

        var response = await client.RegisterConsumerAsync(new RegisterConsumerRequest("c-1", "n-1"));

        Assert.Equal("c-1", response.ConsumerId);
        Assert.Equal("n-1", response.Name);
        Assert.False(response.AlreadyRegistered);
        Assert.Contains("\"consumer_id\":\"c-1\"", handler.LastRequestBody);
        Assert.Contains("\"name\":\"n-1\"", handler.LastRequestBody);
    }

    [Fact]
    public async Task AcknowledgeConsumerAsync_HappyPath_ReturnsTypedResponse()
    {
        var (client, handler) = TestClientFactory.Create(request =>
        {
            Assert.Equal("/consumers/acknowledge", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.OK,
                """{"consumer_id":"c-1","source_id":"s-1","offset":42}""");
        });

        var response = await client.AcknowledgeConsumerAsync(new AcknowledgeConsumerRequest("c-1", "s-1", 42));

        Assert.Equal("c-1", response.ConsumerId);
        Assert.Equal("s-1", response.SourceId);
        Assert.Equal(42UL, response.Offset);
        Assert.Contains("\"offset\":42", handler.LastRequestBody);
    }

    [Fact]
    public async Task GetConsumerPositionAsync_HappyPath_BuildsQueryAndReturnsTypedResponse()
    {
        var (client, handler) = TestClientFactory.Create(request =>
        {
            Assert.Equal(HttpMethod.Get, request.Method);
            Assert.Equal("/consumers/position", request.RequestUri!.AbsolutePath);
            Assert.Equal("consumer_id=c-1", request.RequestUri!.Query.TrimStart('?'));
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.OK,
                """{"consumer_id":"c-1","source_id":"s-1","offset":7,"acknowledged":true}""");
        });

        var response = await client.GetConsumerPositionAsync("c-1");

        Assert.Equal("c-1", response.ConsumerId);
        Assert.Equal(7UL, response.Offset);
        Assert.True(response.Acknowledged);
        Assert.Null(handler.LastRequestBody);
    }

    [Fact]
    public async Task RemoveConsumerAsync_HappyPath_ReturnsTypedResponse()
    {
        var (client, handler) = TestClientFactory.Create(request =>
        {
            Assert.Equal("/consumers/remove", request.RequestUri!.AbsolutePath);
            return FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.OK,
                """{"consumer_id":"c-1","removed":true}""");
        });

        var response = await client.RemoveConsumerAsync("c-1");

        Assert.True(response.Removed);
        Assert.Contains("\"consumer_id\":\"c-1\"", handler.LastRequestBody);
    }
}
