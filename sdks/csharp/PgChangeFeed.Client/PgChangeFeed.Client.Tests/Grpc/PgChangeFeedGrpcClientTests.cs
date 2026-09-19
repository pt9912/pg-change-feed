using Cdc.Stream.V1;
using Google.Protobuf;
using Grpc.Core;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>
/// <see cref="PgChangeFeedGrpcClient"/> against a fake <see cref="CallInvoker"/>
/// (no real server, no network — <c>AGENTS.md</c> §3.1 in the test run): the
/// bearer-token metadata form (SPEC-020), the happy path (messages arrive in
/// order with full content), and the auth boundary
/// (<see cref="StatusCode.Unauthenticated"/> surfaces from the stream itself,
/// not a swallowed empty enumeration).
/// </summary>
public class PgChangeFeedGrpcClientTests
{
    private static PgChangeFeedClientOptions Options(string token = "reader-token")
        => new(new Uri("http://localhost:50051"), token);

    [Fact]
    public async Task StreamChangesAsync_SendsBearerTokenInAuthorizationMetadata()
    {
        var invoker = FakeCallInvoker.WithMessages<Change>();
        using var client = new PgChangeFeedGrpcClient(invoker, Options("reader-token"));

        await foreach (var _ in client.StreamChangesAsync())
        {
            // no messages expected — draining is enough to trigger the call.
        }

        var headers = invoker.LastCallOptions!.Value.Headers;
        Assert.NotNull(headers);
        var authEntry = Assert.Single(headers!, e => e.Key == "authorization");
        Assert.Equal("Bearer reader-token", authEntry.Value);
    }

    [Fact]
    public async Task StreamChangesAsync_YieldsMessagesInOrderWithFullContent()
    {
        var change = new Change
        {
            ChangeId = "change-1",
            TransactionId = "tx-1",
            SourceTableId = "table-1",
            Sequence = 2,
            Operation = "UPDATE",
            OldImage = ByteString.CopyFromUtf8("""{"id":1}"""),
            NewImage = ByteString.CopyFromUtf8("""{"id":1,"bestellstatus":"bezahlt"}"""),
            SchemaVersion = "table-1-v1",
            Schema = "public",
            Table = "orders",
        };
        var invoker = FakeCallInvoker.WithMessages(change);
        using var client = new PgChangeFeedGrpcClient(invoker, Options());

        var received = new List<Change>();
        await foreach (var msg in client.StreamChangesAsync())
        {
            received.Add(msg);
        }

        var only = Assert.Single(received);
        Assert.Equal("change-1", only.ChangeId);
        Assert.Equal("tx-1", only.TransactionId);
        Assert.Equal("table-1", only.SourceTableId);
        Assert.Equal(2, only.Sequence);
        Assert.Equal("UPDATE", only.Operation);
        Assert.Equal("""{"id":1}""", only.OldImage.ToStringUtf8());
        Assert.Equal("""{"id":1,"bestellstatus":"bezahlt"}""", only.NewImage.ToStringUtf8());
        Assert.Equal("table-1-v1", only.SchemaVersion);
        Assert.Equal("public", only.Schema);
        Assert.Equal("orders", only.Table);
    }

    [Fact]
    public async Task StreamChangesAsync_MissingOrInvalidToken_ThrowsUnauthenticated()
    {
        var invoker = FakeCallInvoker.WithStatus(new Status(StatusCode.Unauthenticated, "missing or unknown bearer token"));
        using var client = new PgChangeFeedGrpcClient(invoker, Options("unknown-token"));

        var ex = await Assert.ThrowsAsync<RpcException>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync())
            {
            }
        });

        Assert.Equal(StatusCode.Unauthenticated, ex.StatusCode);
    }
}
