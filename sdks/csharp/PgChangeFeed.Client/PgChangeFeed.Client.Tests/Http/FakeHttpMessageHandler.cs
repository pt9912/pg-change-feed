using System.Net;
using System.Net.Http;
using System.Text;

namespace PgChangeFeed.Client.Tests.Http;

/// <summary>
/// A minimal fake <see cref="HttpMessageHandler"/> for network-free tests of
/// <see cref="PgChangeFeed.Client.Http.PgChangeFeedHttpClient"/> — it never
/// touches a socket; it hands back a caller-supplied response and records
/// the outgoing request (URL, headers, body) for assertions.
/// </summary>
internal sealed class FakeHttpMessageHandler : HttpMessageHandler
{
    private readonly Func<HttpRequestMessage, HttpResponseMessage> _responder;

    public HttpRequestMessage? LastRequest { get; private set; }

    public string? LastRequestBody { get; private set; }

    public FakeHttpMessageHandler(Func<HttpRequestMessage, HttpResponseMessage> responder)
    {
        _responder = responder;
    }

    protected override async Task<HttpResponseMessage> SendAsync(
        HttpRequestMessage request, CancellationToken cancellationToken)
    {
        LastRequest = request;
        LastRequestBody = request.Content is null
            ? null
            : await request.Content.ReadAsStringAsync(cancellationToken).ConfigureAwait(false);
        return _responder(request);
    }

    public static HttpResponseMessage JsonResponse(HttpStatusCode statusCode, string json)
        => new(statusCode)
        {
            Content = new StringContent(json, Encoding.UTF8, "application/json"),
        };
}
