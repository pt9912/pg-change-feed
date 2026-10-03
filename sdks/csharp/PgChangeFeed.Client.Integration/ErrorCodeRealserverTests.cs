using System.Net.Http;
using Grpc.Core;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Grpc;
using PgChangeFeed.Client.Http;
using Xunit;
using GrpcEnableTableRequest = Cdc.Administration.V1.EnableTableRequest;
using HttpEnableTableRequest = PgChangeFeed.Client.Http.Models.EnableTableRequest;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server phase for the message code on the typed errors. Enabling a
/// table that does not exist in the source database is answered with HTTP
/// status 404 and gRPC status <c>NotFound</c>; both carry the message code the
/// server assigns to that case, and the SDK exposes it as
/// <c>MessageCode</c> of the typed exception: from the <c>code</c> field of
/// the HTTP error body, and from the <c>google.rpc.ErrorInfo</c> status detail
/// on the gRPC side. A reader token against the same operation is rejected
/// with HTTP status 403 and gRPC status <c>PermissionDenied</c>, which carry no
/// message code. The unit tests build the wire form themselves; this phase
/// reads it from the real server.
/// </summary>
public sealed class ErrorCodeRealserverTests
{
    private const string MissingTable = "sdk_error_code_missing_table";

    // The server's code for enabling a table that is missing at the source
    // (request: schema "public", table "sdk_error_code_missing_table").
    private const string MissingTableCode = "PCF-E8025";

    private static HttpEnableTableRequest HttpEnableRequest() => new(
        PhaseEnvironment.SourceId, "public", MissingTable, "sdk-error-code", "sdk-error-code-v1", 1,
        PhaseEnvironment.HttpPublication);

    private static GrpcEnableTableRequest GrpcEnableRequest() => new()
    {
        Source = PhaseEnvironment.SourceId,
        Schema = "public",
        Table = MissingTable,
        TableId = "sdk-error-code",
        SchemaVersionId = "sdk-error-code-v1",
        Version = 1,
        Publication = PhaseEnvironment.HttpPublication,
    };

    private static PgChangeFeedHttpClient HttpClientFor(HttpClient httpClient, string token)
        => new(httpClient, new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.HttpAddr), token));

    private static PgChangeFeedAdministrationClient AdministrationClientFor(string token)
        => new(new PgChangeFeedClientOptions(new Uri($"http://{PhaseEnvironment.GrpcAddr}"), token));

    [Fact]
    public async Task EnablingAMissingTableCarriesTheMessageCodeOverHttpAndGrpc()
    {
        PhaseEnvironment.Print("READY");
        using var httpClient = new HttpClient();
        using var administration = AdministrationClientFor(PhaseEnvironment.AdminToken);

        var httpException = await Assert.ThrowsAsync<PgChangeFeedNotFoundException>(() =>
            HttpClientFor(httpClient, PhaseEnvironment.AdminToken)
                .EnableTableAsync(HttpEnableRequest(), PhaseEnvironment.ReceiveCts.Token));
        var grpcException = await Assert.ThrowsAsync<PgChangeFeedGrpcNotFoundException>(() =>
            administration.EnableTableAsync(GrpcEnableRequest(), PhaseEnvironment.ReceiveCts.Token));

        Assert.Equal(404, httpException.StatusCode);
        Assert.Equal(StatusCode.NotFound, grpcException.StatusCode);
        Assert.Equal(MissingTableCode, httpException.MessageCode);
        Assert.Equal(MissingTableCode, grpcException.MessageCode);
        PhaseEnvironment.Print(
            $"RECEIVED code={httpException.MessageCode} http={httpException.StatusCode} grpc={grpcException.StatusCode}");
    }

    [Fact]
    public async Task ARejectedCallCarriesNoMessageCode()
    {
        using var httpClient = new HttpClient();
        using var administration = AdministrationClientFor(PhaseEnvironment.ReaderToken);

        var httpException = await Assert.ThrowsAsync<PgChangeFeedForbiddenException>(() =>
            HttpClientFor(httpClient, PhaseEnvironment.ReaderToken)
                .EnableTableAsync(HttpEnableRequest(), PhaseEnvironment.RejectCts.Token));
        var grpcException = await Assert.ThrowsAsync<PgChangeFeedGrpcPermissionDeniedException>(() =>
            administration.EnableTableAsync(GrpcEnableRequest(), PhaseEnvironment.RejectCts.Token));

        Assert.Equal(403, httpException.StatusCode);
        Assert.Equal(StatusCode.PermissionDenied, grpcException.StatusCode);
        Assert.Null(httpException.MessageCode);
        Assert.Null(grpcException.MessageCode);
        PhaseEnvironment.Print($"REJECTED http={httpException.StatusCode} grpc={grpcException.StatusCode} code=none");
    }
}
