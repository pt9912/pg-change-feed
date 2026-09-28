namespace CdcExamples.Http.Tests;

/// <summary>
/// Baut einen <see cref="System.Net.Http.HttpClient"/> über einen
/// <see cref="FakeHttpMessageHandler"/> — jeder Test bleibt netzlos (kein
/// echter Server, kein echter Socket).
/// </summary>
internal static class TestHttpClientFactory
{
    public static (System.Net.Http.HttpClient HttpClient, FakeHttpMessageHandler Handler) Create(
        Func<HttpRequestMessage, HttpResponseMessage> responder)
    {
        var handler = new FakeHttpMessageHandler(responder);
        return (new System.Net.Http.HttpClient(handler), handler);
    }
}
