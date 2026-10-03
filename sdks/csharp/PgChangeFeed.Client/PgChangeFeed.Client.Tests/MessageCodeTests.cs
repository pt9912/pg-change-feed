using System.Net;
using Cdc.Administration.V1;
using Google.Protobuf;
using Grpc.Core;
using PgChangeFeed.Client.Grpc;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Tests.Grpc;
using PgChangeFeed.Client.Tests.Http;
using Xunit;
using Any = Google.Protobuf.WellKnownTypes.Any;

namespace PgChangeFeed.Client.Tests;

/// <summary>
/// The message code on the typed errors: one input table per wire form (the
/// HTTP error body for the HTTP client and the SSE client, the
/// <c>grpc-status-details-bin</c> trailer for the administration client), the
/// same rows in every SDK language. The status detail of the gRPC tests is
/// serialized here with the Protobuf runtime's own writer, not with the reader
/// under test.
/// </summary>
public class MessageCodeTests
{
    private const string ErrorInfoUrl = "type.googleapis.com/google.rpc.ErrorInfo";
    private const string RetryInfoUrl = "type.googleapis.com/google.rpc.RetryInfo";
    private const string ServerDomain = "pg-change-feed";

    public static TheoryData<int, string, Type, string?, string> HttpRows => new()
    {
        { 400, """{"error":"x","code":"PCF-E8051"}""", typeof(PgChangeFeedBadRequestException), "PCF-E8051", "x" },
        { 400, """{"error":"x"}""", typeof(PgChangeFeedBadRequestException), null, "x" },
        { 400, """{"error":"x","code":""}""", typeof(PgChangeFeedBadRequestException), null, "x" },
        { 400, """{"error":"x","code":null}""", typeof(PgChangeFeedBadRequestException), null, "x" },
        { 400, """{"error":"x","code":5}""", typeof(PgChangeFeedBadRequestException), null, "x" },
        { 404, """{"error":"x","code":"PCF-E8025"}""", typeof(PgChangeFeedNotFoundException), "PCF-E8025", "x" },
        { 500, """{"error":"x","code":"PCF-E7000"}""", typeof(PgChangeFeedServerErrorException), "PCF-E7000", "x" },
        { 503, """{"error":"x","code":"PCF-E2001"}""", typeof(PgChangeFeedUnexpectedStatusException), "PCF-E2001", "x" },
        { 401, """{"error":"x"}""", typeof(PgChangeFeedUnauthorizedException), null, "x" },
        { 403, """{"error":"x"}""", typeof(PgChangeFeedForbiddenException), null, "x" },
        { 400, """{"error":"x","code":"NO-PCF","extra":1}""", typeof(PgChangeFeedBadRequestException), "NO-PCF", "x" },
        { 502, "Bad Gateway", typeof(PgChangeFeedUnexpectedStatusException), null, "Bad Gateway" },
    };

    [Theory]
    [MemberData(nameof(HttpRows))]
    public async Task HttpClient_ErrorCarriesTheMessageCode(int status, string body, Type expected, string? code, string text)
    {
        var (client, _) = Http.TestClientFactory.Create(
            _ => FakeHttpMessageHandler.JsonResponse((HttpStatusCode)status, body));

        var ex = await Assert.ThrowsAnyAsync<PgChangeFeedException>(() => client.ListTablesAsync("src", "pub"));

        Assert.IsType(expected, ex);
        Assert.Equal(code, ex.MessageCode);
        Assert.Equal(text, ex.Message);
        Assert.Equal(status, ex.StatusCode);
    }

    [Theory]
    [MemberData(nameof(HttpRows))]
    public async Task SseClient_ErrorCarriesTheMessageCode(int status, string body, Type expected, string? code, string text)
    {
        var (client, _) = Sse.TestClientFactory.Create(
            _ => FakeHttpMessageHandler.JsonResponse((HttpStatusCode)status, body));

        var ex = await Assert.ThrowsAnyAsync<PgChangeFeedException>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync())
            {
            }
        });

        Assert.IsType(expected, ex);
        Assert.Equal(code, ex.MessageCode);
        Assert.Equal(text, ex.Message);
    }

    [Fact]
    public void HttpException_ConstructedWithoutMessageCode_KeepsWorkingAndHasNone()
    {
        var ex = new PgChangeFeedBadRequestException(400, "x");

        Assert.Equal(400, ex.StatusCode);
        Assert.Equal("x", ex.Message);
        Assert.Null(ex.MessageCode);
    }

    private static byte[] ErrorInfo(string reason, string domain)
    {
        using var stream = new MemoryStream();
        var output = new CodedOutputStream(stream);
        output.WriteTag(1, WireFormat.WireType.LengthDelimited);
        output.WriteString(reason);
        output.WriteTag(2, WireFormat.WireType.LengthDelimited);
        output.WriteString(domain);
        output.Flush();
        return stream.ToArray();
    }

    private static Any Detail(string url, byte[] value) => new() { TypeUrl = url, Value = ByteString.CopyFrom(value) };

    private static byte[] StatusBytes(params Any[] details)
    {
        using var stream = new MemoryStream();
        var output = new CodedOutputStream(stream);
        foreach (var detail in details)
        {
            output.WriteTag(3, WireFormat.WireType.LengthDelimited);
            output.WriteMessage(detail);
        }

        output.Flush();
        return stream.ToArray();
    }

    private static Metadata Trailers(byte[] statusBytes) => new() { { "grpc-status-details-bin", statusBytes } };

    private static Metadata ServerDetail(string reason) =>
        Trailers(StatusBytes(Detail(ErrorInfoUrl, ErrorInfo(reason, ServerDomain))));

    public static TheoryData<StatusCode, Metadata?, Type, string?> GrpcRows => new()
    {
        { StatusCode.InvalidArgument, ServerDetail("PCF-E8051"), typeof(PgChangeFeedGrpcInvalidArgumentException), "PCF-E8051" },
        { StatusCode.NotFound, ServerDetail("PCF-E8025"), typeof(PgChangeFeedGrpcNotFoundException), "PCF-E8025" },
        { StatusCode.Internal, ServerDetail("PCF-E7000"), typeof(PgChangeFeedGrpcInternalException), "PCF-E7000" },
        { StatusCode.InvalidArgument, null, typeof(PgChangeFeedGrpcInvalidArgumentException), null },
        { StatusCode.Unauthenticated, null, typeof(PgChangeFeedGrpcUnauthenticatedException), null },
        { StatusCode.PermissionDenied, null, typeof(PgChangeFeedGrpcPermissionDeniedException), null },
        {
            StatusCode.InvalidArgument,
            Trailers(StatusBytes(Detail(ErrorInfoUrl, ErrorInfo("PCF-E8051", "example.com")))),
            typeof(PgChangeFeedGrpcInvalidArgumentException),
            null
        },
        {
            StatusCode.InvalidArgument,
            Trailers(StatusBytes(Detail(RetryInfoUrl, ErrorInfo("PCF-E8051", ServerDomain)))),
            typeof(PgChangeFeedGrpcInvalidArgumentException),
            null
        },
        {
            StatusCode.InvalidArgument,
            Trailers(StatusBytes(Detail(ErrorInfoUrl, ErrorInfo(string.Empty, ServerDomain)))),
            typeof(PgChangeFeedGrpcInvalidArgumentException),
            null
        },
        { StatusCode.InvalidArgument, Trailers(new byte[] { 0xFF, 0xFF }), typeof(PgChangeFeedGrpcInvalidArgumentException), null },
        {
            StatusCode.InvalidArgument,
            Trailers(StatusBytes(
                Detail(ErrorInfoUrl, ErrorInfo("PCF-E0001", "example.com")),
                Detail(ErrorInfoUrl, ErrorInfo("PCF-E8051", ServerDomain)))),
            typeof(PgChangeFeedGrpcInvalidArgumentException),
            "PCF-E8051"
        },
        { StatusCode.Unavailable, ServerDetail("PCF-E8051"), typeof(PgChangeFeedGrpcUnexpectedStatusException), "PCF-E8051" },
    };

    [Theory]
    [MemberData(nameof(GrpcRows))]
    public async Task AdministrationClient_ErrorCarriesTheMessageCode(StatusCode status, Metadata? trailers, Type expected, string? code)
    {
        var invoker = FakeUnaryCallInvoker.WithStatus(new Status(status, "boom"), trailers);

        var ex = await Assert.ThrowsAnyAsync<PgChangeFeedGrpcException>(
            () => AdministrationTestClientFactory.Create(invoker)
                .ListTablesAsync(new ListTablesRequest { Source = "src", Publication = "pub" }));

        Assert.IsType(expected, ex);
        Assert.Equal(code, ex.MessageCode);
        Assert.Equal(status, ex.StatusCode);
        Assert.Equal("boom", ex.Message);
        Assert.IsType<RpcException>(ex.InnerException);
    }

    [Fact]
    public void GrpcException_ConstructedWithoutMessageCode_KeepsWorkingAndHasNone()
    {
        var cause = new RpcException(new Status(StatusCode.NotFound, "gone"));

        var ex = new PgChangeFeedGrpcNotFoundException("gone", cause);

        Assert.Equal(StatusCode.NotFound, ex.StatusCode);
        Assert.Null(ex.MessageCode);
    }

    [Theory]
    [InlineData("PCF-E4003")]
    [InlineData("")]
    public async Task Diagnose_ErrorCodeArrivesUnchanged(string errorCode)
    {
        var response = new DiagnoseResponse
        {
            Heartbeat = new HeartbeatStatus
            {
                Known = true,
                AgeSeconds = 1.0,
                ErrorClass = errorCode.Length == 0 ? string.Empty : "schema",
                ErrorCode = errorCode,
            },
        };
        var invoker = FakeUnaryCallInvoker.WithResponse(response);

        var result = await AdministrationTestClientFactory.Create(invoker)
            .DiagnoseAsync(new DiagnoseRequest { Source = "src" });

        Assert.Equal(errorCode, result.Heartbeat.ErrorCode);
    }
}
