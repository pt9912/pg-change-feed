using System.Net.Http;
using PgChangeFeed.Client.Http;

namespace PgChangeFeed.Client.Tests.Http;

/// <summary>
/// Builds a <see cref="PgChangeFeedHttpClient"/> wired to a
/// <see cref="FakeHttpMessageHandler"/> — every test in this project stays
/// network-free (no real server, no real socket).
/// </summary>
internal static class TestClientFactory
{
    public static readonly Uri ServerAddress = new("http://example.invalid:8080");

    public static (PgChangeFeedHttpClient Client, FakeHttpMessageHandler Handler) Create(
        Func<HttpRequestMessage, HttpResponseMessage> responder, string apiToken = "test-token")
    {
        var handler = new FakeHttpMessageHandler(responder);
        var httpClient = new HttpClient(handler);
        var options = new PgChangeFeedClientOptions(ServerAddress, apiToken);
        return (new PgChangeFeedHttpClient(httpClient, options), handler);
    }
}
