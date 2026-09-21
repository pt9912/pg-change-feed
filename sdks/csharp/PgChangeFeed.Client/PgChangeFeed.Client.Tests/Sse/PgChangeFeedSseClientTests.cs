using System.Net;
using System.Net.Http;
using System.Text;
using PgChangeFeed.Client.Sse.Models;
using Xunit;

namespace PgChangeFeed.Client.Tests.Sse;

/// <summary>
/// <see cref="PgChangeFeed.Client.Sse.PgChangeFeedSseClient"/> against a
/// fake <see cref="HttpMessageHandler"/> (no real server, no network —
/// <c>AGENTS.md</c> §3.1 in the test run): the bearer-token header form
/// (SPEC-021), the happy path (several frames arrive in order with full
/// content), and the boundary around an incomplete final frame.
/// </summary>
public class PgChangeFeedSseClientTests
{
    private static string SseFrame(string changeId, string table) =>
        "event: change\n" +
        $$"""data: {"change_id":"{{changeId}}","transaction_id":"tx-1","source_table_id":"table-1","sequence":2,"operation":"UPDATE","old_image":{"id":1},"new_image":{"id":1,"bestellstatus":"bezahlt"},"schema_version":"table-1-v1","schema":"public","table":"{{table}}"}""" +
        "\n\n";

    private static HttpResponseMessage EventStreamResponse(string body) => new(HttpStatusCode.OK)
    {
        Content = new StringContent(body, Encoding.UTF8, "text/event-stream"),
    };

    [Fact]
    public async Task StreamChangesAsync_SendsBearerTokenInAuthorizationHeader()
    {
        var (client, handler) = TestClientFactory.Create(
            _ => EventStreamResponse(SseFrame("c-1", "orders")),
            apiToken: "reader-token");

        var received = new List<Change>();
        await foreach (var change in client.StreamChangesAsync())
        {
            received.Add(change);
        }

        Assert.Equal("Bearer reader-token", handler.LastRequest!.Headers.Authorization!.ToString());
        Assert.Single(received);
    }

    [Fact]
    public async Task StreamChangesAsync_YieldsFramesInOrderWithFullContent()
    {
        var (client, _) = TestClientFactory.Create(
            _ => EventStreamResponse(SseFrame("c-1", "orders") + SseFrame("c-2", "orders")));

        var received = new List<Change>();
        await foreach (var change in client.StreamChangesAsync())
        {
            received.Add(change);
        }

        Assert.Equal(2, received.Count);
        var first = received[0];
        Assert.Equal("c-1", first.ChangeId);
        Assert.Equal("c-2", received[1].ChangeId);
        Assert.Equal("tx-1", first.TransactionId);
        Assert.Equal("table-1", first.SourceTableId);
        Assert.Equal(2, first.Sequence);
        Assert.Equal("UPDATE", first.Operation);
        Assert.Equal("""{"id":1}""", first.OldImage!.Value.GetRawText());
        Assert.Equal("""{"id":1,"bestellstatus":"bezahlt"}""", first.NewImage!.Value.GetRawText());
        Assert.Equal("table-1-v1", first.SchemaVersion);
        Assert.Equal("public", first.Schema);
        Assert.Equal("orders", first.Table);
    }

    [Fact]
    public async Task StreamChangesAsync_IncompleteFinalFrame_EndsEnumerationWithoutException()
    {
        var (client, _) = TestClientFactory.Create(
            _ => EventStreamResponse("event: change\ndata: {\"change_id\":\"c-1\"}"));

        var received = new List<Change>();
        await foreach (var change in client.StreamChangesAsync())
        {
            received.Add(change);
        }

        Assert.Empty(received);
    }

    [Fact]
    public async Task StreamChangesAsync_EmptyStream_EndsEnumerationImmediately()
    {
        var (client, _) = TestClientFactory.Create(_ => EventStreamResponse(""));

        var received = new List<Change>();
        await foreach (var change in client.StreamChangesAsync())
        {
            received.Add(change);
        }

        Assert.Empty(received);
    }
}
