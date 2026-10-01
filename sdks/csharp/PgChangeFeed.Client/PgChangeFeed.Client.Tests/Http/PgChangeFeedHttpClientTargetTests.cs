using System.Net;
using PgChangeFeed.Client.Http;
using Xunit;

namespace PgChangeFeed.Client.Tests.Http;

/// <summary>
/// The optional <c>target</c> parameter of <c>ReadChangesAsync</c>: it appears
/// in the query only when set, without it the query is byte-identical to the
/// call before the parameter existed.
/// </summary>
public class PgChangeFeedHttpClientTargetTests
{
    private static string QueryOf(Func<PgChangeFeedHttpClient, Task> call)
    {
        string? seen = null;
        var (client, _) = TestClientFactory.Create(request =>
        {
            seen = request.RequestUri!.Query.TrimStart('?');
            return FakeHttpMessageHandler.JsonResponse(HttpStatusCode.OK, """{"changes":[]}""");
        });
        call(client).GetAwaiter().GetResult();
        return seen!;
    }

    [Fact]
    public void ReadChangesAsync_WithoutTarget_QueryCarriesNoTarget()
    {
        Assert.Equal("source=src", QueryOf(c => c.ReadChangesAsync("src")));
    }

    [Theory]
    [InlineData(null, "source=src")]
    [InlineData("", "source=src&target=")]
    [InlineData("eu", "source=src&target=eu")]
    [InlineData("a&b=c", "source=src&target=a%26b%3Dc")]
    public void ReadChangesAsync_Target_IsSentOnlyWhenSetAndEscaped(string? target, string expectedQuery)
    {
        Assert.Equal(expectedQuery, QueryOf(c => c.ReadChangesAsync("src", target: target)));
    }

    [Fact]
    public void ReadChangesAsync_TargetWithSchemaTableRange_IsAConjunctionOnTheWire()
    {
        var query = QueryOf(c => c.ReadChangesAsync(
            "src", schema: "public", table: "orders", from: 1, to: 10, limit: 5, target: "eu"));

        Assert.Equal("source=src&schema=public&table=orders&target=eu&from=1&to=10&limit=5", query);
    }
}
