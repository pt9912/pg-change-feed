using System.Net;
using System.Net.Http;
using System.Text;
using Xunit;

namespace PgChangeFeed.Client.Tests.Sse;

/// <summary>
/// The optional <c>target</c> of the SSE stream: it appears as the query
/// parameter <c>target</c> only when set; without it the request URI is the
/// bare stream path.
/// </summary>
public class PgChangeFeedSseClientTargetTests
{
    private static async Task<Uri> RequestUriOf(string? target)
    {
        var (client, handler) = TestClientFactory.Create(_ => new HttpResponseMessage(HttpStatusCode.OK)
        {
            Content = new StringContent(string.Empty, Encoding.UTF8, "text/event-stream"),
        });
        await foreach (var _ in client.StreamChangesAsync(target: target))
        {
            // an empty body ends the enumeration immediately.
        }
        return handler.LastRequest!.RequestUri!;
    }

    [Fact]
    public async Task StreamChangesAsync_WithoutTarget_RequestHasNoQuery()
    {
        var uri = await RequestUriOf(null);

        Assert.Equal("/changes/stream", uri.AbsolutePath);
        Assert.Equal(string.Empty, uri.Query);
    }

    [Theory]
    [InlineData("", "?target=")]
    [InlineData("eu", "?target=eu")]
    [InlineData("a&b=c", "?target=a%26b%3Dc")]
    public async Task StreamChangesAsync_Target_IsSentAsEscapedQueryParameter(string target, string expectedQuery)
    {
        var uri = await RequestUriOf(target);

        Assert.Equal("/changes/stream", uri.AbsolutePath);
        Assert.Equal(expectedQuery, uri.Query);
    }
}
