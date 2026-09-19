using System.Net;
using System.Net.Http;
using System.Text;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Http.Models;
using Xunit;

namespace PgChangeFeed.Client.Tests.Http;

/// <summary>
/// The SPEC-018 error-response mapping — the auth boundary (<c>401</c>
/// missing/unknown token, <c>403</c> a <c>reader</c> token against an
/// <c>admin</c> endpoint) plus the remaining documented statuses
/// (<c>400</c>/<c>500</c>) and the defensive fallback for anything outside
/// that closed set.
/// </summary>
public class PgChangeFeedHttpClientAuthBoundaryTests
{
    [Fact]
    public async Task MissingOrUnknownToken_ThrowsUnauthorized()
    {
        var (client, handler) = TestClientFactory.Create(
            _ => FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.Unauthorized, """{"error":"missing or unknown bearer token"}"""),
            apiToken: "unknown-token");

        var ex = await Assert.ThrowsAsync<PgChangeFeedUnauthorizedException>(
            () => client.ListTablesAsync("src", "pub"));

        Assert.Equal(401, ex.StatusCode);
        Assert.Equal("missing or unknown bearer token", ex.Message);
        Assert.Equal("Bearer unknown-token", handler.LastRequest!.Headers.Authorization!.ToString());
    }

    [Fact]
    public async Task ReaderTokenAgainstAdminEndpoint_ThrowsForbidden()
    {
        var (client, handler) = TestClientFactory.Create(
            _ => FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.Forbidden, """{"error":"reader token cannot reach an admin endpoint"}"""),
            apiToken: "reader-token");

        var ex = await Assert.ThrowsAsync<PgChangeFeedForbiddenException>(
            () => client.RegisterConsumerAsync(new RegisterConsumerRequest("c-1", "n-1")));

        Assert.Equal(403, ex.StatusCode);
        Assert.Equal("reader token cannot reach an admin endpoint", ex.Message);
        Assert.Equal("Bearer reader-token", handler.LastRequest!.Headers.Authorization!.ToString());
    }

    [Fact]
    public async Task InvalidRequestBody_ThrowsBadRequest()
    {
        var (client, _) = TestClientFactory.Create(
            _ => FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.BadRequest, """{"error":"consumer_id must not be empty"}"""));

        var ex = await Assert.ThrowsAsync<PgChangeFeedBadRequestException>(
            () => client.RegisterConsumerAsync(new RegisterConsumerRequest("", "n-1")));

        Assert.Equal(400, ex.StatusCode);
        Assert.Equal("consumer_id must not be empty", ex.Message);
    }

    [Fact]
    public async Task UnexpectedInternalError_ThrowsServerError()
    {
        var (client, _) = TestClientFactory.Create(
            _ => FakeHttpMessageHandler.JsonResponse(
                HttpStatusCode.InternalServerError, """{"error":"unexpected internal error"}"""));

        var ex = await Assert.ThrowsAsync<PgChangeFeedServerErrorException>(
            () => client.ListTablesAsync("src", "pub"));

        Assert.Equal(500, ex.StatusCode);
        Assert.Equal("unexpected internal error", ex.Message);
    }

    [Fact]
    public async Task StatusCodeOutsideDocumentedSet_ThrowsUnexpectedStatus()
    {
        var (client, _) = TestClientFactory.Create(
            _ => FakeHttpMessageHandler.JsonResponse(HttpStatusCode.Conflict, """{"error":"unexpected"}"""));

        var ex = await Assert.ThrowsAsync<PgChangeFeedUnexpectedStatusException>(
            () => client.ListTablesAsync("src", "pub"));

        Assert.Equal(409, ex.StatusCode);
    }

    [Fact]
    public async Task NonJsonErrorBody_FallsBackToRawBody()
    {
        var (client, _) = TestClientFactory.Create(_ => new HttpResponseMessage(HttpStatusCode.InternalServerError)
        {
            Content = new StringContent("plain text failure", Encoding.UTF8, "text/plain"),
        });

        var ex = await Assert.ThrowsAsync<PgChangeFeedServerErrorException>(
            () => client.ListTablesAsync("src", "pub"));

        Assert.Equal("plain text failure", ex.Message);
    }

    [Fact]
    public async Task NonJsonSuccessBody_ThrowsMalformedResponse()
    {
        var (client, _) = TestClientFactory.Create(_ => new HttpResponseMessage(HttpStatusCode.OK)
        {
            Content = new StringContent("not json at all", Encoding.UTF8, "application/json"),
        });

        var ex = await Assert.ThrowsAsync<PgChangeFeedMalformedResponseException>(
            () => client.ListTablesAsync("src", "pub"));

        Assert.Equal(200, ex.StatusCode);
        Assert.IsNotType<System.Text.Json.JsonException>(ex);
    }
}
