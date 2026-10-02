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
    private static async Task<Uri> RequestUriOf(
        string? target, string? schema = null, string? table = null)
    {
        var (client, handler) = TestClientFactory.Create(_ => new HttpResponseMessage(HttpStatusCode.OK)
        {
            Content = new StringContent(string.Empty, Encoding.UTF8, "text/event-stream"),
        });
        await foreach (var _ in client.StreamChangesAsync(target: target, schema: schema, table: table))
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
    [InlineData("a+b", "?target=a%2Bb")]
    [InlineData("100%", "?target=100%25")]
    [InlineData("ü", "?target=%C3%BC")]
    [InlineData("a b", "?target=a%20b")] // the Python client sends a space as "+"; the server decodes both alike
    public async Task StreamChangesAsync_Target_IsSentAsEscapedQueryParameter(string target, string expectedQuery)
    {
        var uri = await RequestUriOf(target);

        Assert.Equal("/changes/stream", uri.AbsolutePath);
        Assert.Equal(expectedQuery, uri.Query);
    }

    [Fact]
    public async Task StreamChangesAsync_WithoutFilters_RequestHasNoQuery()
    {
        var uri = await RequestUriOf(null, null, null);

        Assert.Equal("/changes/stream", uri.AbsolutePath);
        Assert.Equal(string.Empty, uri.Query);
    }

    [Theory]
    [InlineData("", "?schema=")]
    [InlineData("eu", "?schema=eu")]
    [InlineData("a&b=c", "?schema=a%26b%3Dc")]
    [InlineData("a+b", "?schema=a%2Bb")]
    [InlineData("100%", "?schema=100%25")]
    [InlineData("ü", "?schema=%C3%BC")]
    [InlineData("a b", "?schema=a%20b")] // the Python client sends a space as "+"; the server decodes both alike
    public async Task StreamChangesAsync_Schema_IsSentAsEscapedQueryParameter(string schema, string expectedQuery)
    {
        var uri = await RequestUriOf(null, schema: schema);

        Assert.Equal("/changes/stream", uri.AbsolutePath);
        Assert.Equal(expectedQuery, uri.Query);
    }

    [Theory]
    [InlineData("", "?table=")]
    [InlineData("eu", "?table=eu")]
    [InlineData("a&b=c", "?table=a%26b%3Dc")]
    [InlineData("a+b", "?table=a%2Bb")]
    [InlineData("100%", "?table=100%25")]
    [InlineData("ü", "?table=%C3%BC")]
    [InlineData("a b", "?table=a%20b")] // the Python client sends a space as "+"; the server decodes both alike
    public async Task StreamChangesAsync_Table_IsSentAsEscapedQueryParameter(string table, string expectedQuery)
    {
        var uri = await RequestUriOf(null, table: table);

        Assert.Equal("/changes/stream", uri.AbsolutePath);
        Assert.Equal(expectedQuery, uri.Query);
    }

    [Fact]
    public async Task StreamChangesAsync_SchemaAndTable_AreSentTogether()
    {
        var uri = await RequestUriOf(null, schema: "public", table: "orders");

        Assert.Equal("?schema=public&table=orders", uri.Query);
    }

    [Fact]
    public async Task StreamChangesAsync_TableAndTarget_AreSentAsConjunction()
    {
        var uri = await RequestUriOf("eu", table: "orders");

        Assert.Equal("?table=orders&target=eu", uri.Query);
    }

    [Fact]
    public async Task StreamChangesAsync_SchemaTableAndTarget_AreSentAsConjunction()
    {
        var uri = await RequestUriOf("eu", schema: "public", table: "orders");

        Assert.Equal("?schema=public&table=orders&target=eu", uri.Query);
    }
}
