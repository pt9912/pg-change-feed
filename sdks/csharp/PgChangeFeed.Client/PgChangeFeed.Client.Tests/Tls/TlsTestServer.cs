using System.Buffers.Binary;
using System.Net;
using System.Security.Cryptography.X509Certificates;
using Cdc.Administration.V1;
using Cdc.Stream.V1;
using Google.Protobuf;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Server.Kestrel.Core;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

namespace PgChangeFeed.Client.Tests.Tls;

/// <summary>
/// A Kestrel server on a loopback port that answers the four surfaces the
/// clients use (HTTP list of tables, SSE stream, gRPC stream, gRPC
/// administration list of tables) with fixed content, over TLS with the given
/// certificate or in plaintext when none is given. It lets the tests read what
/// a client does with a certificate it is shown.
/// </summary>
internal sealed class TlsTestServer : IAsyncDisposable
{
    internal const string SentinelChangeId = "tls-test-change";

    private readonly WebApplication _app;

    private TlsTestServer(WebApplication app, int port)
    {
        _app = app;
        Port = port;
    }

    internal int Port { get; }

    internal Uri Address(string host, bool tls) => new($"{(tls ? "https" : "http")}://{host}:{Port}");

    /// <summary>
    /// Starts the server. With a certificate it serves HTTP/1.1 and HTTP/2
    /// over TLS; without one it serves the given plaintext protocol
    /// (HTTP/2 for gRPC, HTTP/1.1 otherwise).
    /// </summary>
    internal static async Task<TlsTestServer> StartAsync(
        X509Certificate2? certificate, HttpProtocols plaintextProtocols = HttpProtocols.Http1)
    {
        var builder = WebApplication.CreateSlimBuilder();
        builder.Logging.ClearProviders();
        builder.WebHost.ConfigureKestrel(kestrel => kestrel.Listen(IPAddress.Loopback, 0, listen =>
        {
            if (certificate is not null)
            {
                listen.Protocols = HttpProtocols.Http1AndHttp2;
                listen.UseHttps(certificate);
            }
            else
            {
                listen.Protocols = plaintextProtocols;
            }
        }));
        var app = builder.Build();

        app.MapGet("/tables", () => Results.Text("""{"tables":[],"retained":[]}""", "application/json"));
        app.MapGet("/changes/stream", async (HttpContext context) =>
        {
            context.Response.ContentType = "text/event-stream";
            await context.Response.WriteAsync(
                "event: change\n" +
                $$"""data: {"change_id":"{{SentinelChangeId}}","transaction_id":"tx-1","source_table_id":"t-1","sequence":1,"operation":"INSERT","new_image":{"id":1},"schema_version":"v1","schema":"public","table":"orders"}""" +
                "\n\n");
        });
        app.MapPost("/cdc.stream.v1.ChangeStream/StreamChanges", (HttpContext context) =>
            GrpcReplyAsync(context, new Change { ChangeId = SentinelChangeId, Operation = "INSERT", Schema = "public", Table = "orders" }));
        app.MapPost("/cdc.administration.v1.Administration/ListTables", (HttpContext context) =>
            GrpcReplyAsync(context, new ListTablesResponse()));

        await app.StartAsync();
        var address = app.Services.GetRequiredService<Microsoft.AspNetCore.Hosting.Server.IServer>()
            .Features.Get<Microsoft.AspNetCore.Hosting.Server.Features.IServerAddressesFeature>()!
            .Addresses.Single();
        return new TlsTestServer(app, new Uri(address).Port);
    }

    private static async Task GrpcReplyAsync(HttpContext context, IMessage message)
    {
        context.Response.ContentType = "application/grpc";
        var body = message.ToByteArray();
        var frame = new byte[5 + body.Length];
        BinaryPrimitives.WriteUInt32BigEndian(frame.AsSpan(1, 4), (uint)body.Length);
        body.CopyTo(frame, 5);
        await context.Response.Body.WriteAsync(frame);
        context.Response.AppendTrailer("grpc-status", "0");
    }

    public async ValueTask DisposeAsync()
    {
        await _app.StopAsync();
        await _app.DisposeAsync();
    }
}
