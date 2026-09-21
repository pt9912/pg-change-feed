using System.Text;
using PgChangeFeed.Client.Nats;
using Xunit;

namespace PgChangeFeed.Client.Tests.Nats;

/// <summary>
/// Happy-path behavior of <see cref="PgChangeFeedNatsStreamClient.StreamChangesAsync"/>
/// against a <see cref="FakeNatsClient"/> — no real NATS server involved.
/// Muster: <c>PgChangeFeed.Client.Tests.Grpc.PgChangeFeedGrpcClientTests</c>.
/// </summary>
public class PgChangeFeedNatsStreamClientTests
{
    private static byte[] ChangePayload(string changeId, string operation) => Encoding.UTF8.GetBytes($$"""
        {
          "change_id": "{{changeId}}",
          "transaction_id": "tx-1",
          "source_table_id": "table-1",
          "sequence": 1,
          "operation": "{{operation}}",
          "old_image": null,
          "new_image": {"id": 1},
          "schema_version": "table-1-v1",
          "schema": "public",
          "table": "orders"
        }
        """);

    [Fact]
    public async Task StreamChangesAsync_YieldsEveryMessageInOrder()
    {
        var fake = FakeNatsClient.WithPayloads(
            ChangePayload("change-1", "INSERT"),
            ChangePayload("change-2", "UPDATE"));
        var client = new PgChangeFeedNatsStreamClient(fake);

        var received = new List<string>();
        await foreach (var change in client.StreamChangesAsync())
        {
            received.Add($"{change.ChangeId}:{change.Operation}");
        }

        Assert.Equal(["change-1:INSERT", "change-2:UPDATE"], received);
    }

    [Fact]
    public async Task StreamChangesAsync_DefaultsToAllSourcesSubject()
    {
        var fake = FakeNatsClient.WithPayloads(ChangePayload("change-1", "INSERT"));
        var client = new PgChangeFeedNatsStreamClient(fake);

        await foreach (var _ in client.StreamChangesAsync())
        {
        }

        Assert.Equal(PgChangeFeedNatsStreamClient.AllSourcesSubject, fake.LastSubject);
    }

    [Fact]
    public async Task StreamChangesAsync_SubscribesToTheSuppliedSubject()
    {
        var fake = FakeNatsClient.WithPayloads(ChangePayload("change-1", "INSERT"));
        var client = new PgChangeFeedNatsStreamClient(fake);
        var subject = PgChangeFeedNatsStreamClient.BuildSubject("source-1", "public", "orders");

        await foreach (var _ in client.StreamChangesAsync(subject))
        {
        }

        Assert.Equal("cdc.stream.source-1.public.orders", fake.LastSubject);
    }

    [Fact]
    public async Task StreamChangesAsync_EmptySubscription_YieldsNothing()
    {
        var fake = FakeNatsClient.WithPayloads();
        var client = new PgChangeFeedNatsStreamClient(fake);

        var received = new List<string>();
        await foreach (var change in client.StreamChangesAsync())
        {
            received.Add(change.ChangeId);
        }

        Assert.Empty(received);
    }
}
