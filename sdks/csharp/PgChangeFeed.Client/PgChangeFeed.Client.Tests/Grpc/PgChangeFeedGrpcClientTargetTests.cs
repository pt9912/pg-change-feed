using Cdc.Administration.V1;
using Cdc.Stream.V1;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>
/// The optional <c>target</c> of the gRPC stream request and of the
/// <c>ReadChanges</c> request: set it travels in the request message, left out
/// the message equals the one built before the field existed.
/// </summary>
public class PgChangeFeedGrpcClientTargetTests
{
    private static PgChangeFeedClientOptions Options()
        => new(new Uri("http://localhost:50051"), "reader-token");

    private static async Task<StreamChangesRequest> SentStreamRequest(Func<PgChangeFeedGrpcClient, IAsyncEnumerable<Change>> open)
    {
        var invoker = FakeCallInvoker.WithMessages<Change>();
        using var client = new PgChangeFeedGrpcClient(invoker, Options());
        await foreach (var _ in open(client))
        {
            // draining is enough to trigger the call.
        }
        return Assert.IsType<StreamChangesRequest>(invoker.LastRequest);
    }

    [Fact]
    public async Task StreamChangesAsync_WithoutTarget_RequestEqualsTheEmptyRequest()
    {
        var sent = await SentStreamRequest(c => c.StreamChangesAsync());

        Assert.Equal(new StreamChangesRequest(), sent);
    }

    [Theory]
    [InlineData(null, "")]
    [InlineData("", "")]
    [InlineData("eu", "eu")]
    public async Task StreamChangesAsync_Target_TravelsInTheRequest(string? target, string expected)
    {
        var sent = await SentStreamRequest(c => c.StreamChangesAsync(target: target));

        Assert.Equal(expected, sent.Target);
    }

    [Fact]
    public async Task StreamChangesAsync_TargetWithSchemaAndTable_AllThreeAreSet()
    {
        var sent = await SentStreamRequest(c => c.StreamChangesAsync(schema: "public", table: "orders", target: "eu"));

        Assert.Equal("public", sent.Schema);
        Assert.Equal("orders", sent.Table);
        Assert.Equal("eu", sent.Target);
    }

    [Theory]
    [InlineData("", "")]
    [InlineData("eu", "eu")]
    public async Task AdministrationReadChangesAsync_Target_ReachesTheServerRequestUnchanged(string target, string expected)
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new ReadChangesResponse());
        var client = AdministrationTestClientFactory.Create(invoker);

        await client.ReadChangesAsync(new ReadChangesRequest { Source = "src", Schema = "public", Target = target });

        var sent = Assert.IsType<ReadChangesRequest>(invoker.LastRequest);
        Assert.Equal(expected, sent.Target);
        Assert.Equal("public", sent.Schema);
    }
}
