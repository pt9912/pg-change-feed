using System.Net;
using System.Text.Json.Serialization;
using Xunit;

namespace CdcExamples.Http.Tests;

/// <summary>
/// Prüft den gemeinsamen Anfrage/Antwort-Ablauf aller zehn Fähigkeiten —
/// netzlos, über <see cref="FakeHttpMessageHandler"/>. Form-Vorbild:
/// <c>examples/http-client/request_test.go</c>.
/// </summary>
public class RequestHelperTests
{
    private sealed record Echo([property: JsonPropertyName("field")] string Field);

    [Fact]
    public async Task SendJsonAsyncDecodesSuccessBodyAndSetsBearerToken()
    {
        var (httpClient, handler) = TestHttpClientFactory.Create(request =>
        {
            Assert.Equal("Bearer test-token", request.Headers.Authorization?.ToString());
            return FakeHttpMessageHandler.JsonResponse(HttpStatusCode.OK, """{"field":"echo-wert"}""");
        });

        var resp = await RequestHelper.SendJsonAsync<Echo>(
            httpClient, HttpMethod.Post, "http://feed:8080/echo", "test-token",
            new Echo("wert"), HttpStatusCode.OK);

        Assert.Equal("echo-wert", resp.Field);
        Assert.Contains("\"field\":\"wert\"", handler.LastRequestBody);
    }

    [Fact]
    public async Task SendJsonAsyncSendsNoBodyWithoutRequestBody()
    {
        var (httpClient, handler) = TestHttpClientFactory.Create(_ =>
            FakeHttpMessageHandler.JsonResponse(HttpStatusCode.OK, """{"field":"wert"}"""));

        await RequestHelper.SendJsonAsync<Echo>(
            httpClient, HttpMethod.Get, "http://feed:8080/echo", "test-token", null, HttpStatusCode.OK);

        Assert.Null(handler.LastRequestBody);
    }

    [Fact]
    public async Task SendJsonAsyncFailsOnUnexpectedStatus()
    {
        var (httpClient, _) = TestHttpClientFactory.Create(_ =>
            FakeHttpMessageHandler.JsonResponse(HttpStatusCode.Forbidden, """{"error":"kein Zugriff"}"""));

        var ex = await Assert.ThrowsAsync<InvalidOperationException>(
            () => RequestHelper.SendJsonAsync<Echo>(
                httpClient, HttpMethod.Get, "http://feed:8080/echo", "test-token", null, HttpStatusCode.OK));

        Assert.Contains("403", ex.Message);
        Assert.Contains("kein Zugriff", ex.Message);
    }

    [Fact]
    public async Task SendJsonAsyncFailsOnUnreachableHost()
    {
        using var httpClient = new System.Net.Http.HttpClient { Timeout = TimeSpan.FromSeconds(2) };

        await Assert.ThrowsAnyAsync<Exception>(() => RequestHelper.SendJsonAsync<Echo>(
            httpClient, HttpMethod.Get, "http://127.0.0.1:1/nichts", "test-token", null, HttpStatusCode.OK));
    }
}
