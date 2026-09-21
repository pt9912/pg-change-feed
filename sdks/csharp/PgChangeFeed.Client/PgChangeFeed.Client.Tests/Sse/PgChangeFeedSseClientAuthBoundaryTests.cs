using System.Net;
using System.Net.Http;
using System.Text;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Tests.Http;
using Xunit;

namespace PgChangeFeed.Client.Tests.Sse;

/// <summary>
/// The SPEC-021 error-response mapping when opening <c>GET /changes/stream</c>
/// fails before any frame is delivered — the same closed set of typed
/// exceptions as <see cref="PgChangeFeed.Client.Http.PgChangeFeedHttpClient"/>
/// (SPEC-018's uniform <c>{"error": "&lt;text&gt;"}</c> body), reused rather
/// than a second hierarchy (<c>PgChangeFeedException</c> and its five
/// concrete subtypes). Muster: <c>PgChangeFeedHttpClientAuthBoundaryTests</c>.
/// </summary>
public class PgChangeFeedSseClientAuthBoundaryTests
{
    [Fact]
    public async Task MissingOrUnknownToken_ThrowsUnauthorized()
    {
        var (client, handler) = TestClientFactory.Create(
            _ => FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.Unauthorized, """{"error":"missing or unknown bearer token"}"""),
            apiToken: "unknown-token");

        var ex = await Assert.ThrowsAsync<PgChangeFeedUnauthorizedException>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync())
            {
            }
        });

        Assert.Equal(401, ex.StatusCode);
        Assert.Equal("missing or unknown bearer token", ex.Message);
        Assert.Equal("Bearer unknown-token", handler.LastRequest!.Headers.Authorization!.ToString());
    }

    [Fact]
    public async Task NoBroadcasterWired_ThrowsUnexpectedStatus()
    {
        // SPEC-021: `503`, wenn CDC_HTTP_ADDR gesetzt, aber kein Broadcaster
        // verdrahtet ist — außerhalb des geschlossenen 400/401/403/404/500-Sets.
        var (client, _) = TestClientFactory.Create(
            _ => FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.ServiceUnavailable, """{"error":"ChangeStream ohne Broadcaster verdrahtet"}"""));

        var ex = await Assert.ThrowsAsync<PgChangeFeedUnexpectedStatusException>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync())
            {
            }
        });

        Assert.Equal(503, ex.StatusCode);
    }

    [Fact]
    public async Task UnexpectedInternalError_ThrowsServerError()
    {
        var (client, _) = TestClientFactory.Create(
            _ => FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.InternalServerError, """{"error":"unexpected internal error"}"""));

        var ex = await Assert.ThrowsAsync<PgChangeFeedServerErrorException>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync())
            {
            }
        });

        Assert.Equal(500, ex.StatusCode);
        Assert.Equal("unexpected internal error", ex.Message);
    }

    [Fact]
    public async Task NonJsonErrorBody_FallsBackToRawBody()
    {
        var (client, _) = TestClientFactory.Create(_ => new HttpResponseMessage(HttpStatusCode.InternalServerError)
        {
            Content = new StringContent("plain text failure", Encoding.UTF8, "text/plain"),
        });

        var ex = await Assert.ThrowsAsync<PgChangeFeedServerErrorException>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync())
            {
            }
        });

        Assert.Equal("plain text failure", ex.Message);
    }

    [Fact]
    public async Task NonJsonFrameData_ThrowsMalformedResponse()
    {
        var (client, _) = TestClientFactory.Create(_ => new HttpResponseMessage(HttpStatusCode.OK)
        {
            Content = new StringContent(
                "event: change\ndata: not json at all\n\n", Encoding.UTF8, "text/event-stream"),
        });

        var ex = await Assert.ThrowsAsync<PgChangeFeedMalformedResponseException>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync())
            {
            }
        });

        Assert.Equal(200, ex.StatusCode);
        Assert.IsNotType<System.Text.Json.JsonException>(ex);
    }
}
