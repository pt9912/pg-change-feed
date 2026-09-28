using Cdc.Administration.V1;
using Google.Protobuf;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>
/// Happy-path coverage for <c>RunRetention</c> and <c>ReadChanges</c>,
/// including the zero-value boundary of each request's optional fields.
/// </summary>
public class PgChangeFeedAdministrationClientRetentionAndChangesTests
{
    [Fact]
    public async Task RunRetentionAsync_HappyPath_ReturnsTypedResponse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new RunRetentionResponse { Deleted = 3 });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.RunRetentionAsync(new RunRetentionRequest { Source = "src", MinAgeNanos = 1_000_000_000 });

        Assert.Equal(3, response.Deleted);
    }

    [Fact]
    public async Task RunRetentionAsync_ZeroMinAge_IsValid()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new RunRetentionResponse { Deleted = 0 });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.RunRetentionAsync(new RunRetentionRequest { Source = "src", MinAgeNanos = 0 });

        Assert.Equal(0, response.Deleted);
        var sent = Assert.IsType<RunRetentionRequest>(invoker.LastRequest);
        Assert.Equal(0, sent.MinAgeNanos);
    }

    [Fact]
    public async Task ReadChangesAsync_HappyPath_ReturnsTypedResponseWithImages()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new ReadChangesResponse
        {
            Changes =
            {
                new ChangeRecord
                {
                    CommitPosition = 1, ChangeId = "ch-1", TransactionId = "tx-1", SourceTableId = "t-1",
                    Schema = "public", Table = "orders", Sequence = 0, Operation = "INSERT",
                    NewImage = ByteString.CopyFromUtf8("""{"id":1}"""),
                    SchemaVersion = "sv-1", CommittedAt = "2026-09-19T00:00:00Z", Origin = "wal",
                },
            },
        });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.ReadChangesAsync(
            new ReadChangesRequest { Source = "src", Schema = "public", Table = "orders", From = 1, To = 10, Limit = 5 });

        var change = Assert.Single(response.Changes);
        Assert.Equal("ch-1", change.ChangeId);
        Assert.Equal("INSERT", change.Operation);
        Assert.Equal("wal", change.Origin);
        Assert.Equal(0, change.OldImage.Length);
    }

    [Fact]
    public async Task ReadChangesAsync_NoOptionalFields_LeavesThemAtZero()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new ReadChangesResponse());
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.ReadChangesAsync(new ReadChangesRequest { Source = "src" });

        Assert.Empty(response.Changes);
        var sent = Assert.IsType<ReadChangesRequest>(invoker.LastRequest);
        Assert.Equal(string.Empty, sent.Schema);
        Assert.Equal(string.Empty, sent.Table);
        Assert.Equal(0UL, sent.From);
        Assert.Equal(0UL, sent.To);
        Assert.Equal(0, sent.Limit);
    }
}
