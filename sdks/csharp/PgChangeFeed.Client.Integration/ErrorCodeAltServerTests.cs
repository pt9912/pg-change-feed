using System.Net.Http;
using Cdc.Administration.V1;
using Grpc.Core;
using PgChangeFeed.Client;
using PgChangeFeed.Client.Grpc;
using PgChangeFeed.Client.Http;
using Xunit;
using GrpcEnableTableRequest = Cdc.Administration.V1.EnableTableRequest;
using HttpEnableTableRequest = PgChangeFeed.Client.Http.Models.EnableTableRequest;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server phase against a server older than the message codes. The same
/// calls as <see cref="ErrorCodeRealserverTests"/> end with the same typed
/// errors, but the server sends no <c>code</c> in the HTTP error body, no
/// <c>google.rpc.ErrorInfo</c> in the gRPC status detail and no
/// <c>error_code</c> in the diagnose heartbeat: every message code stays
/// <see langword="null"/>, the error text is unchanged and the SDK does not
/// fail. The error state of the diagnose report is written by the runner
/// without an <c>error_code</c>.
/// </summary>
public sealed class ErrorCodeAltServerTests
{
    private const string MissingTable = "sdk_error_code_missing_table";

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
    public async Task AnOlderServerCarriesNoMessageCodeAndTheSdkStaysIntact()
    {
        PhaseEnvironment.Print("READY");
        using var httpClient = new HttpClient();
        using var admin = AdministrationClientFor(PhaseEnvironment.AdminToken);
        using var reader = AdministrationClientFor(PhaseEnvironment.ReaderToken);

        var httpException = await Assert.ThrowsAsync<PgChangeFeedNotFoundException>(() =>
            HttpClientFor(httpClient, PhaseEnvironment.AdminToken)
                .EnableTableAsync(HttpEnableRequest(), PhaseEnvironment.ReceiveCts.Token));
        var grpcException = await Assert.ThrowsAsync<PgChangeFeedGrpcNotFoundException>(() =>
            admin.EnableTableAsync(GrpcEnableRequest(), PhaseEnvironment.ReceiveCts.Token));
        Assert.Equal(404, httpException.StatusCode);
        Assert.Equal(StatusCode.NotFound, grpcException.StatusCode);
        Assert.Null(httpException.MessageCode);
        Assert.Null(grpcException.MessageCode);
        Assert.False(string.IsNullOrEmpty(httpException.Message));
        Assert.False(string.IsNullOrEmpty(grpcException.Message));
        Assert.Equal(PhaseEnvironment.AltServerHttpText, httpException.Message);

        var httpForbidden = await Assert.ThrowsAsync<PgChangeFeedForbiddenException>(() =>
            HttpClientFor(httpClient, PhaseEnvironment.ReaderToken)
                .EnableTableAsync(HttpEnableRequest(), PhaseEnvironment.RejectCts.Token));
        var grpcForbidden = await Assert.ThrowsAsync<PgChangeFeedGrpcPermissionDeniedException>(() =>
            reader.EnableTableAsync(GrpcEnableRequest(), PhaseEnvironment.RejectCts.Token));
        Assert.Equal(403, httpForbidden.StatusCode);
        Assert.Equal(StatusCode.PermissionDenied, grpcForbidden.StatusCode);
        Assert.Null(httpForbidden.MessageCode);
        Assert.Null(grpcForbidden.MessageCode);

        var normal = await admin.DiagnoseAsync(
            new DiagnoseRequest { Source = PhaseEnvironment.SourceId }, PhaseEnvironment.ReceiveCts.Token);
        Assert.NotNull(normal.Heartbeat);
        Assert.True(normal.Heartbeat.Known);
        Assert.Equal(string.Empty, normal.Heartbeat.ErrorClass);
        Assert.Equal(string.Empty, normal.Heartbeat.ErrorCode);
        PhaseEnvironment.Print("NORMAL_DONE");

        // The runner writes the error state (error_class without error_code)
        // after NORMAL_DONE; the periodic heartbeat resets it, so the call is
        // repeated until the report shows it.
        HeartbeatStatus? failing = null;
        var deadline = DateTime.UtcNow.AddSeconds(60);
        while (DateTime.UtcNow < deadline)
        {
            var report = await admin.DiagnoseAsync(
                new DiagnoseRequest { Source = PhaseEnvironment.SourceId }, PhaseEnvironment.ReceiveCts.Token);
            if (report.Heartbeat?.ErrorClass == "schema")
            {
                failing = report.Heartbeat;
                break;
            }

            await Task.Delay(200);
        }

        Assert.NotNull(failing);
        Assert.Equal("schema", failing.ErrorClass);
        Assert.Equal(string.Empty, failing.ErrorCode);
        PhaseEnvironment.Print(
            $"RECEIVED code=none http=404 grpc=NOT_FOUND text={httpException.Message.Length} diag_error_code=leer");
    }
}
