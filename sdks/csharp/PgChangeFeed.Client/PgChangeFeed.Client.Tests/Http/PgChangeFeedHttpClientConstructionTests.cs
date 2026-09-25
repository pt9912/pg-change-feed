using System.Net.Http;
using PgChangeFeed.Client.Http;
using Xunit;

namespace PgChangeFeed.Client.Tests.Http;

/// <summary>Constructor argument validation — no bearer token, no global state.</summary>
public class PgChangeFeedHttpClientConstructionTests
{
    [Fact]
    public void Constructor_WithNullHttpClient_Throws()
    {
        var options = new PgChangeFeedClientOptions(TestClientFactory.ServerAddress, "token");

        Assert.Throws<ArgumentNullException>(() => new PgChangeFeedHttpClient(null!, options));
    }

    [Fact]
    public void Constructor_WithNullOptions_Throws()
    {
        using var httpClient = new HttpClient();

        Assert.Throws<ArgumentNullException>(() => new PgChangeFeedHttpClient(httpClient, null!));
    }
}
