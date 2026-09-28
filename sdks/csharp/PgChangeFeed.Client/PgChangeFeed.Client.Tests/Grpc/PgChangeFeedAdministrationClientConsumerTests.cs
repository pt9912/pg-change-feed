using Cdc.Administration.V1;
using Grpc.Core;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>
/// Happy-path coverage for the four consumer-management RPCs
/// (<c>RegisterConsumer</c>, <c>AcknowledgeConsumer</c>,
/// <c>GetConsumerPosition</c>, <c>RemoveConsumer</c>), plus the bearer-token
/// metadata form and the null-request boundary shared by every method.
/// </summary>
public class PgChangeFeedAdministrationClientConsumerTests
{
    [Fact]
    public async Task RegisterConsumerAsync_HappyPath_ReturnsTypedResponseAndSendsBearerToken()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(
            new RegisterConsumerResponse { ConsumerId = "c-1", Name = "n-1", AlreadyRegistered = false });
        var client = AdministrationTestClientFactory.Create(invoker, apiToken: "admin-token");

        var response = await client.RegisterConsumerAsync(
            new RegisterConsumerRequest { ConsumerId = "c-1", Name = "n-1" });

        Assert.Equal("c-1", response.ConsumerId);
        Assert.False(response.AlreadyRegistered);
        var headers = invoker.LastCallOptions!.Value.Headers;
        var authEntry = Assert.Single(headers!, e => e.Key == "authorization");
        Assert.Equal("Bearer admin-token", authEntry.Value);
    }

    [Fact]
    public async Task RegisterConsumerAsync_NullRequest_ThrowsArgumentNullException()
    {
        var client = AdministrationTestClientFactory.Create(FakeUnaryCallInvoker.WithStatus(new Status(StatusCode.OK, string.Empty)));

        await Assert.ThrowsAsync<ArgumentNullException>(() => client.RegisterConsumerAsync(null!));
    }

    [Fact]
    public async Task AcknowledgeConsumerAsync_HappyPath_ReturnsTypedResponse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(
            new AcknowledgeConsumerResponse { ConsumerId = "c-1", SourceId = "s-1", Offset = 42 });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.AcknowledgeConsumerAsync(
            new AcknowledgeConsumerRequest { ConsumerId = "c-1", SourceId = "s-1", Offset = 42 });

        Assert.Equal("s-1", response.SourceId);
        Assert.Equal(42UL, response.Offset);
        var sent = Assert.IsType<AcknowledgeConsumerRequest>(invoker.LastRequest);
        Assert.Equal(42UL, sent.Offset);
    }

    [Fact]
    public async Task GetConsumerPositionAsync_HappyPath_ReturnsTypedResponse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(
            new GetConsumerPositionResponse { ConsumerId = "c-1", SourceId = "s-1", Offset = 7, Acknowledged = true });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.GetConsumerPositionAsync(new GetConsumerPositionRequest { ConsumerId = "c-1" });

        Assert.Equal(7UL, response.Offset);
        Assert.True(response.Acknowledged);
    }

    [Fact]
    public async Task GetConsumerPositionAsync_NeverAcknowledged_ReportsAcknowledgedFalse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(
            new GetConsumerPositionResponse { ConsumerId = "c-1", SourceId = "s-1", Offset = 0, Acknowledged = false });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.GetConsumerPositionAsync(new GetConsumerPositionRequest { ConsumerId = "c-1" });

        Assert.Equal(0UL, response.Offset);
        Assert.False(response.Acknowledged);
    }

    [Fact]
    public async Task RemoveConsumerAsync_HappyPath_ReturnsTypedResponse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(
            new RemoveConsumerResponse { ConsumerId = "c-1", Removed = true });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.RemoveConsumerAsync(new RemoveConsumerRequest { ConsumerId = "c-1" });

        Assert.True(response.Removed);
    }

    [Fact]
    public async Task RemoveConsumerAsync_NeverRegistered_ReportsRemovedFalse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(
            new RemoveConsumerResponse { ConsumerId = "c-1", Removed = false });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.RemoveConsumerAsync(new RemoveConsumerRequest { ConsumerId = "c-1" });

        Assert.False(response.Removed);
    }
}
