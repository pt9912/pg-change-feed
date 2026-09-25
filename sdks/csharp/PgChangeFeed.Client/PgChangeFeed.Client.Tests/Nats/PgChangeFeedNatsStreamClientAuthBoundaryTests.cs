using System.Text;
using NATS.Client.Core;
using PgChangeFeed.Client.Nats;
using Xunit;

namespace PgChangeFeed.Client.Tests.Nats;

/// <summary>
/// The connection-level auth boundary and the malformed-payload path when
/// subscribing to the NATS stream — no real NATS server involved. Pattern:
/// <c>PgChangeFeed.Client.Tests.Sse.PgChangeFeedSseClientAuthBoundaryTests</c>,
/// adapted to NATS's connection-level (not per-message) auth boundary.
/// </summary>
public class PgChangeFeedNatsStreamClientAuthBoundaryTests
{
    [Fact]
    public async Task RejectedConnection_PropagatesNatsServerExceptionUnwrapped()
    {
        // A missing or wrong stream token is rejected by the NATS server
        // itself, at the connection level — the real
        // NATS.Net client surfaces this as a NatsServerException with
        // IsAuthError=true, thrown from the subscription enumeration itself,
        // not as a swallowed empty stream.
        var authError = new NatsServerException("authorization violation");
        var fake = FakeNatsClient.WithFailure(authError);
        var client = new PgChangeFeedNatsStreamClient(fake);

        var ex = await Assert.ThrowsAsync<NatsServerException>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync())
            {
            }
        });

        Assert.Same(authError, ex);
        Assert.True(ex.IsAuthError);
    }

    [Fact]
    public async Task NonJsonPayload_ThrowsMalformedMessage()
    {
        var fake = FakeNatsClient.WithPayloads(Encoding.UTF8.GetBytes("not json at all"));
        var client = new PgChangeFeedNatsStreamClient(fake);

        var ex = await Assert.ThrowsAsync<PgChangeFeedNatsMalformedMessageException>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync())
            {
            }
        });

        Assert.IsNotType<System.Text.Json.JsonException>(ex);
    }

    [Fact]
    public async Task JsonNullPayload_ThrowsMalformedMessage()
    {
        var fake = FakeNatsClient.WithPayloads(Encoding.UTF8.GetBytes("null"));
        var client = new PgChangeFeedNatsStreamClient(fake);

        await Assert.ThrowsAsync<PgChangeFeedNatsMalformedMessageException>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync())
            {
            }
        });
    }
}
