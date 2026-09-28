using System.Net;
using System.Text;

namespace CdcExamples.Http.Tests;

/// <summary>
/// Ein minimaler Fake-<see cref="HttpMessageHandler"/> für netzlose Tests der
/// Aufruf-Funktionen dieses Beispiels — er berührt nie einen echten Socket;
/// er liefert eine vom Aufrufer vorgegebene Antwort und zeichnet die
/// gesendete Anfrage (URL, Body) für Assertions auf. Form-Vorbild:
/// <c>sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Http/FakeHttpMessageHandler.cs</c>.
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

    public static HttpResponseMessage JsonResponse(HttpStatusCode statusCode, string json) => new(statusCode)
    {
        Content = new StringContent(json, Encoding.UTF8, "application/json"),
    };
}
