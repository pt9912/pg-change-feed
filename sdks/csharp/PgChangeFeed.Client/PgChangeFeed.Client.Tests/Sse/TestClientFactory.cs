using System.Net.Http;
using PgChangeFeed.Client.Sse;
using PgChangeFeed.Client.Tests.Http;

namespace PgChangeFeed.Client.Tests.Sse;

/// <summary>
/// Builds a <see cref="PgChangeFeedSseClient"/> wired to a
/// <see cref="FakeHttpMessageHandler"/> (reused from
/// <c>PgChangeFeed.Client.Tests.Http</c> — <c>internal</c> is assembly-scoped,
/// not namespace-scoped, so this test project's existing fake handler
/// applies here unchanged) — every test in this project stays network-free
/// (no real server, no real socket).
/// </summary>
internal static class TestClientFactory
{
    public static readonly Uri ServerAddress = new("http://example.invalid:8080");

    public static (PgChangeFeedSseClient Client, FakeHttpMessageHandler Handler) Create(
        Func<HttpRequestMessage, HttpResponseMessage> responder, string apiToken = "test-token")
    {
        var handler = new FakeHttpMessageHandler(responder);
        var httpClient = new HttpClient(handler);
        var options = new PgChangeFeedClientOptions(ServerAddress, apiToken);
        return (new PgChangeFeedSseClient(httpClient, options), handler);
    }
}
