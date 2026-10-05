using System.Security.Cryptography.X509Certificates;
using Cdc.Administration.V1;
using Grpc.Core;
using Microsoft.AspNetCore.Server.Kestrel.Core;
using PgChangeFeed.Client.Grpc;
using PgChangeFeed.Client.Http;
using PgChangeFeed.Client.Sse;
using Xunit;

namespace PgChangeFeed.Client.Tests.Tls;

/// <summary>
/// The four surfaces that connect over TLS (HTTP, SSE, gRPC stream, gRPC
/// administration) against a Kestrel server with a certificate created at run
/// time: with the trust anchor the connection works, and a foreign anchor, a
/// server name that is not in the certificate, an expired certificate and no
/// anchor at all (the certificate is unknown to the system) each end with a
/// connection error of the surface, without a plaintext retry.
/// </summary>
public sealed class TlsClientTests : IDisposable
{
    public enum Surface { Http, Sse, GrpcStream, GrpcAdmin }

    private const string ServerName = "localhost";

    private readonly TestCertificates _certificates = new();

    public void Dispose() => _certificates.Dispose();

    public static TheoryData<Surface> Surfaces => new() { Surface.Http, Surface.Sse, Surface.GrpcStream, Surface.GrpcAdmin };

    private static PgChangeFeedClientOptions Options(Uri address, string? anchorFile)
        => new(address, "tls-test-token", anchorFile);

    private static async Task UseAsync(Surface surface, PgChangeFeedClientOptions options)
    {
        switch (surface)
        {
            case Surface.Http:
                using (var http = new PgChangeFeedHttpClient(options))
                {
                    var tables = await http.ListTablesAsync("src", "pub");
                    Assert.Empty(tables.Tables);
                }

                break;
            case Surface.Sse:
                using (var sse = new PgChangeFeedSseClient(options))
                {
                    await foreach (var change in sse.StreamChangesAsync())
                    {
                        Assert.Equal(TlsTestServer.SentinelChangeId, change.ChangeId);
                        return;
                    }

                    Assert.Fail("The SSE stream ended without a change.");
                }

                break;
            case Surface.GrpcStream:
                using (var grpc = new PgChangeFeedGrpcClient(options))
                {
                    await foreach (var change in grpc.StreamChangesAsync())
                    {
                        Assert.Equal(TlsTestServer.SentinelChangeId, change.ChangeId);
                        return;
                    }

                    Assert.Fail("The gRPC stream ended without a change.");
                }

                break;
            default:
                using (var admin = new PgChangeFeedAdministrationClient(options))
                {
                    var tables = await admin.ListTablesAsync(new ListTablesRequest { Source = "src", Publication = "pub" });
                    Assert.Empty(tables.Tables);
                }

                break;
        }
    }

    /// <summary>
    /// The surface fails with its connection error: the <see cref="HttpRequestException"/>
    /// of the HTTP client, the <see cref="RpcException"/> of the stream, the
    /// <see cref="PgChangeFeedGrpcUnexpectedStatusException"/> without a
    /// message code of the administration client; both gRPC errors carry the
    /// <see cref="HttpRequestException"/> of the transport (the server was
    /// never reached).
    /// </summary>
    private static async Task ExpectConnectionFailureAsync(Surface surface, PgChangeFeedClientOptions options)
    {
        switch (surface)
        {
            case Surface.Http:
            case Surface.Sse:
                await Assert.ThrowsAsync<HttpRequestException>(() => UseAsync(surface, options));
                break;
            case Surface.GrpcStream:
                var rpc = await Assert.ThrowsAsync<RpcException>(() => UseAsync(surface, options));
                Assert.IsType<HttpRequestException>(rpc.Status.DebugException);
                break;
            default:
                var admin = await Assert.ThrowsAsync<PgChangeFeedGrpcUnexpectedStatusException>(() => UseAsync(surface, options));
                Assert.IsType<HttpRequestException>(((RpcException)admin.InnerException!).Status.DebugException);
                Assert.Null(admin.MessageCode);
                break;
        }
    }

    private X509Certificate2 ServerCertificate(DateTimeOffset? notBefore = null, DateTimeOffset? notAfter = null)
        => _certificates.Create([ServerName], notBefore: notBefore, notAfter: notAfter);

    [Theory]
    [MemberData(nameof(Surfaces))]
    public async Task WithTrustAnchor_ConnectsOverTls(Surface surface)
    {
        var certificate = ServerCertificate();
        await using var server = await TlsTestServer.StartAsync(certificate);

        await UseAsync(surface, Options(server.Address(ServerName, tls: true), _certificates.WritePem(certificate)));
    }

    [Theory]
    [MemberData(nameof(Surfaces))]
    public async Task IssuerAsTrustAnchor_ConnectsOverTls(Surface surface)
    {
        var (server, authority) = _certificates.CreateIssued(ServerName);
        await using var testServer = await TlsTestServer.StartAsync(server);

        await UseAsync(surface, Options(testServer.Address(ServerName, tls: true), _certificates.WritePem(authority)));
    }

    [Theory]
    [MemberData(nameof(Surfaces))]
    public async Task IssuerOfAnotherServerAsTrustAnchor_FailsTheConnection(Surface surface)
    {
        var (server, _) = _certificates.CreateIssued(ServerName);
        var (_, otherAuthority) = _certificates.CreateIssued(ServerName);
        await using var testServer = await TlsTestServer.StartAsync(server);

        await ExpectConnectionFailureAsync(surface, Options(
            testServer.Address(ServerName, tls: true), _certificates.WritePem(otherAuthority)));
    }

    [Theory]
    [MemberData(nameof(Surfaces))]
    public async Task AnchorFileWithSeveralCertificates_TrustsTheOneThatMatches(Surface surface)
    {
        var certificate = ServerCertificate();
        var foreign = ServerCertificate();
        await using var server = await TlsTestServer.StartAsync(certificate);

        await UseAsync(surface, Options(
            server.Address(ServerName, tls: true), _certificates.WritePem(foreign, certificate)));
    }

    [Theory]
    [MemberData(nameof(Surfaces))]
    public async Task ForeignAnchor_FailsTheConnection(Surface surface)
    {
        var certificate = ServerCertificate();
        var foreign = ServerCertificate();
        await using var server = await TlsTestServer.StartAsync(certificate);

        await ExpectConnectionFailureAsync(surface, Options(
            server.Address(ServerName, tls: true), _certificates.WritePem(foreign)));
    }

    [Theory]
    [MemberData(nameof(Surfaces))]
    public async Task ServerNameNotInCertificate_FailsTheConnection(Surface surface)
    {
        var certificate = ServerCertificate();
        await using var server = await TlsTestServer.StartAsync(certificate);

        // The anchor is the server certificate itself; the address names the
        // loopback IP, which the certificate does not list.
        await ExpectConnectionFailureAsync(surface, Options(
            server.Address("127.0.0.1", tls: true), _certificates.WritePem(certificate)));
    }

    [Theory]
    [MemberData(nameof(Surfaces))]
    public async Task ExpiredCertificate_FailsTheConnection(Surface surface)
    {
        var expired = ServerCertificate(
            notBefore: DateTimeOffset.UtcNow.AddDays(-3), notAfter: DateTimeOffset.UtcNow.AddDays(-2));
        await using var server = await TlsTestServer.StartAsync(expired);

        await ExpectConnectionFailureAsync(surface, Options(
            server.Address(ServerName, tls: true), _certificates.WritePem(expired)));
    }

    [Theory]
    [MemberData(nameof(Surfaces))]
    public async Task WithoutAnchor_TheSystemTrustDecidesAndRefusesAnUnknownCertificate(Surface surface)
    {
        await using var server = await TlsTestServer.StartAsync(ServerCertificate());

        await ExpectConnectionFailureAsync(surface, Options(server.Address(ServerName, tls: true), null));
    }

    [Theory]
    [MemberData(nameof(Surfaces))]
    public async Task WithoutAnchor_PlaintextAddressWorksUnchanged(Surface surface)
    {
        var protocols = surface is Surface.GrpcStream or Surface.GrpcAdmin ? HttpProtocols.Http2 : HttpProtocols.Http1;
        await using var server = await TlsTestServer.StartAsync(null, protocols);

        await UseAsync(surface, Options(server.Address(ServerName, tls: false), null));
    }

    [Theory]
    [MemberData(nameof(Surfaces))]
    public async Task TrustAnchorWithPlaintextAddress_IsRefusedWhenTheClientIsCreated(Surface surface)
    {
        var certificate = ServerCertificate();
        var options = Options(new Uri("http://localhost:9"), _certificates.WritePem(certificate));

        await Assert.ThrowsAsync<ArgumentException>(() => UseAsync(surface, options));
    }
}
