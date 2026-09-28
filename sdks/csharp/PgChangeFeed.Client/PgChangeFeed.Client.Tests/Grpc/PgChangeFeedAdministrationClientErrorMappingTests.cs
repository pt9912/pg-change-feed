using Cdc.Administration.V1;
using Grpc.Core;
using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>
/// The gRPC-status-to-exception mapping shared by every RPC of
/// <see cref="PgChangeFeedAdministrationClient"/> — exercised through
/// <c>ListTables</c> as the representative call, since the mapping itself
/// does not depend on which RPC failed.
/// </summary>
public class PgChangeFeedAdministrationClientErrorMappingTests
{
    private static Task<ListTablesResponse> Call(FakeUnaryCallInvoker invoker)
        => AdministrationTestClientFactory.Create(invoker)
            .ListTablesAsync(new ListTablesRequest { Source = "src", Publication = "pub" });

    [Fact]
    public async Task InvalidArgument_ThrowsTypedInvalidArgumentException()
    {
        var invoker = FakeUnaryCallInvoker.WithStatus(new Status(StatusCode.InvalidArgument, "source must not be empty"));

        var ex = await Assert.ThrowsAsync<PgChangeFeedGrpcInvalidArgumentException>(() => Call(invoker));

        Assert.Equal(StatusCode.InvalidArgument, ex.StatusCode);
        Assert.Equal("source must not be empty", ex.Message);
    }

    [Fact]
    public async Task Unauthenticated_ThrowsTypedUnauthenticatedException()
    {
        var invoker = FakeUnaryCallInvoker.WithStatus(new Status(StatusCode.Unauthenticated, "missing or unknown bearer token"));

        var ex = await Assert.ThrowsAsync<PgChangeFeedGrpcUnauthenticatedException>(() => Call(invoker));

        Assert.Equal(StatusCode.Unauthenticated, ex.StatusCode);
    }

    [Fact]
    public async Task PermissionDenied_ThrowsTypedPermissionDeniedException()
    {
        var invoker = FakeUnaryCallInvoker.WithStatus(new Status(StatusCode.PermissionDenied, "Rechtsklasse unzureichend für diese RPC"));

        var ex = await Assert.ThrowsAsync<PgChangeFeedGrpcPermissionDeniedException>(() => Call(invoker));

        Assert.Equal(StatusCode.PermissionDenied, ex.StatusCode);
    }

    [Fact]
    public async Task NotFound_ThrowsTypedNotFoundException()
    {
        var invoker = FakeUnaryCallInvoker.WithStatus(new Status(StatusCode.NotFound, "table does not exist at source"));

        var ex = await Assert.ThrowsAsync<PgChangeFeedGrpcNotFoundException>(() => Call(invoker));

        Assert.Equal(StatusCode.NotFound, ex.StatusCode);
    }

    [Fact]
    public async Task Internal_ThrowsTypedInternalException()
    {
        var invoker = FakeUnaryCallInvoker.WithStatus(new Status(StatusCode.Internal, "interner Fehler"));

        var ex = await Assert.ThrowsAsync<PgChangeFeedGrpcInternalException>(() => Call(invoker));

        Assert.Equal(StatusCode.Internal, ex.StatusCode);
    }

    [Fact]
    public async Task StatusCodeOutsideDocumentedSet_ThrowsUnexpectedStatus()
    {
        var invoker = FakeUnaryCallInvoker.WithStatus(new Status(StatusCode.Unavailable, "server unreachable"));

        var ex = await Assert.ThrowsAsync<PgChangeFeedGrpcUnexpectedStatusException>(() => Call(invoker));

        Assert.Equal(StatusCode.Unavailable, ex.StatusCode);
        Assert.Equal("server unreachable", ex.Message);
    }
}
