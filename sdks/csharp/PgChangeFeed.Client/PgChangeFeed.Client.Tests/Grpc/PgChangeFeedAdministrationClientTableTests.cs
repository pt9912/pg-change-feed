using Cdc.Administration.V1;
using Grpc.Core;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>
/// Happy-path (and the table-only <c>NotFound</c>) coverage for the four
/// table-management RPCs (<c>EnableTable</c>, <c>DisableTable</c>,
/// <c>GetTableStatus</c>, <c>ListTables</c>).
/// </summary>
public class PgChangeFeedAdministrationClientTableTests
{
    private static EnableTableRequest SampleEnableRequest() => new()
    {
        Source = "src", Schema = "public", Table = "orders",
        TableId = "t-1", SchemaVersionId = "sv-1", Version = 1, Publication = "pub",
    };

    [Fact]
    public async Task EnableTableAsync_HappyPath_ReturnsTypedResponse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new EnableTableResponse
        {
            TableId = "t-1", Source = "src", Schema = "public", Table = "orders", AlreadyEnabled = false,
        });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.EnableTableAsync(SampleEnableRequest());

        Assert.Equal("t-1", response.TableId);
        Assert.False(response.AlreadyEnabled);
    }

    [Fact]
    public async Task EnableTableAsync_TableMissingAtSource_ThrowsNotFound()
    {
        var invoker = FakeUnaryCallInvoker.WithStatus(
            new Status(StatusCode.NotFound, "table does not exist at source"));
        var client = AdministrationTestClientFactory.Create(invoker);

        var ex = await Assert.ThrowsAsync<PgChangeFeedGrpcNotFoundException>(
            () => client.EnableTableAsync(SampleEnableRequest()));

        Assert.Equal(StatusCode.NotFound, ex.StatusCode);
        Assert.Equal("table does not exist at source", ex.Message);
        Assert.IsType<RpcException>(ex.InnerException);
    }

    [Fact]
    public async Task DisableTableAsync_HappyPath_ReturnsTypedResponse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new DisableTableResponse { Removed = true, Retained = false });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.DisableTableAsync(
            new DisableTableRequest { Source = "src", Schema = "public", Table = "orders", Publication = "pub" });

        Assert.True(response.Removed);
        Assert.False(response.Retained);
    }

    [Fact]
    public async Task DisableTableAsync_TableMissingAtSource_ThrowsNotFound()
    {
        var invoker = FakeUnaryCallInvoker.WithStatus(
            new Status(StatusCode.NotFound, "table does not exist at source"));
        var client = AdministrationTestClientFactory.Create(invoker);

        await Assert.ThrowsAsync<PgChangeFeedGrpcNotFoundException>(() => client.DisableTableAsync(
            new DisableTableRequest { Source = "src", Schema = "public", Table = "orders", Publication = "pub" }));
    }

    [Fact]
    public async Task GetTableStatusAsync_HappyPath_ReturnsTypedResponse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new GetTableStatusResponse { Enabled = true, Retained = false });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.GetTableStatusAsync(
            new GetTableStatusRequest { Source = "src", Schema = "public", Table = "orders", Publication = "pub" });

        Assert.True(response.Enabled);
        Assert.False(response.Retained);
    }

    [Fact]
    public async Task GetTableStatusAsync_NeverEnabled_ReportsBothFalse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new GetTableStatusResponse { Enabled = false, Retained = false });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.GetTableStatusAsync(
            new GetTableStatusRequest { Source = "src", Schema = "public", Table = "orders", Publication = "pub" });

        Assert.False(response.Enabled);
        Assert.False(response.Retained);
    }

    [Fact]
    public async Task ListTablesAsync_HappyPath_ReturnsTypedResponse()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse(new ListTablesResponse
        {
            Tables = { new SourceTable { TableId = "t-1", Source = "src", Schema = "public", Table = "orders" } },
        });
        var client = AdministrationTestClientFactory.Create(invoker);

        var response = await client.ListTablesAsync(new ListTablesRequest { Source = "src", Publication = "pub" });

        Assert.Single(response.Tables);
        Assert.Equal("t-1", response.Tables[0].TableId);
        Assert.Empty(response.Retained);
    }
}
