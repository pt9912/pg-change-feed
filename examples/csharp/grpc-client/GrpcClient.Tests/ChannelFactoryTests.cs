using System.Net;
using System.Net.Security;
using System.Net.Sockets;
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;
using Xunit;

namespace CdcExamples.Grpc.Tests;

/// <summary>
/// Prüft <see cref="ChannelFactory"/>: die Wahl „Anker gesetzt, dann TLS“ ist
/// an den Optionswert gebunden. Der Handler läuft gegen einen TLS-Server auf
/// Loopback (<c>TcpListener</c> und <c>SslStream</c>) mit einem im Test
/// erzeugten Zertifikat; kein Zertifikat und kein Schlüssel liegen im Repo.
/// </summary>
public sealed class ChannelFactoryTests : IDisposable
{
    private readonly string _dir = Directory.CreateTempSubdirectory("grpc-client-tls").FullName;

    public void Dispose() => Directory.Delete(_dir, recursive: true);

    private static X509Certificate2 NewCertificate()
    {
        using var key = ECDsa.Create(ECCurve.NamedCurves.nistP256);
        var request = new CertificateRequest("CN=grpc-client-test", key, HashAlgorithmName.SHA256);
        var san = new SubjectAlternativeNameBuilder();
        san.AddIpAddress(IPAddress.Loopback);
        request.CertificateExtensions.Add(san.Build());
        using var created = request.CreateSelfSigned(DateTimeOffset.UtcNow.AddHours(-1), DateTimeOffset.UtcNow.AddHours(1));
        return X509CertificateLoader.LoadPkcs12(created.Export(X509ContentType.Pfx), null);
    }

    private string WritePem(X509Certificate2 cert)
    {
        var path = Path.Combine(_dir, $"{Guid.NewGuid():N}.pem");
        File.WriteAllText(path, cert.ExportCertificatePem());
        return path;
    }

    /// <summary>Startet einen TLS-Server, der jede Anfrage mit HTTP 200 beantwortet.</summary>
    private static (int Port, CancellationTokenSource Stop) StartServer(X509Certificate2 cert)
    {
        var listener = new TcpListener(IPAddress.Loopback, 0);
        listener.Start();
        var stop = new CancellationTokenSource();
        stop.Token.Register(listener.Stop);
        _ = Task.Run(async () =>
        {
            while (!stop.IsCancellationRequested)
            {
                try
                {
                    using var tcp = await listener.AcceptTcpClientAsync().ConfigureAwait(false);
                    using var ssl = new SslStream(tcp.GetStream());
                    await ssl.AuthenticateAsServerAsync(cert).ConfigureAwait(false);
                    var buffer = new byte[4096];
                    _ = await ssl.ReadAsync(buffer).ConfigureAwait(false);
                    await ssl.WriteAsync("HTTP/1.1 200 OK\r\nContent-Length: 0\r\nConnection: close\r\n\r\n"u8.ToArray()).ConfigureAwait(false);
                    await ssl.FlushAsync().ConfigureAwait(false);
                }
                catch (Exception) when (!stop.IsCancellationRequested)
                {
                    // ein abgelehnter Handshake endet nur diese Verbindung
                }
                catch (Exception)
                {
                    return;
                }
            }
        });
        return (((IPEndPoint)listener.LocalEndpoint).Port, stop);
    }

    private static async Task<HttpStatusCode> GetAsync(SocketsHttpHandler handler, int port)
    {
        using var client = new HttpClient(handler);
        using var response = await client.GetAsync($"https://127.0.0.1:{port}/").ConfigureAwait(false);
        return response.StatusCode;
    }

    [Fact]
    public async Task TlsHandler_WithAnchor_ReachesTlsServer()
    {
        using var cert = NewCertificate();
        var (port, stop) = StartServer(cert);
        try
        {
            var handler = ChannelFactory.TlsHandler(WritePem(cert));
            Assert.Equal(HttpStatusCode.OK, await GetAsync(handler, port));
        }
        finally
        {
            stop.Cancel();
        }
    }

    [Fact]
    public async Task TlsHandler_WithForeignAnchor_FailsHandshake()
    {
        using var serverCert = NewCertificate();
        using var foreign = NewCertificate();
        var (port, stop) = StartServer(serverCert);
        try
        {
            var handler = ChannelFactory.TlsHandler(WritePem(foreign));
            await Assert.ThrowsAsync<HttpRequestException>(() => GetAsync(handler, port));
        }
        finally
        {
            stop.Cancel();
        }
    }

    [Fact]
    public void Address_UsesHttpsOnlyWithAnchor()
    {
        var plain = new Config("feed:9090", "", "", "stream", "", "", "", "", 0, "", "", 1, "", "", 0, 0, 0, 0);
        Assert.Equal("http://feed:9090", ChannelFactory.Address(plain));
        Assert.Equal("https://feed:9090", ChannelFactory.Address(plain with { CaFile = "/anker.pem" }));
    }

    [Fact]
    public void TlsHandler_RejectsMissingFile()
    {
        var ex = Assert.Throws<ArgumentException>(() => ChannelFactory.TlsHandler(Path.Combine(_dir, "fehlt.pem")));
        Assert.Contains("nicht lesbar", ex.Message);
    }

    [Fact]
    public void TlsHandler_RejectsFileWithoutPem()
    {
        var path = Path.Combine(_dir, "kein-pem.txt");
        File.WriteAllText(path, "kein Zertifikat");
        var ex = Assert.Throws<ArgumentException>(() => ChannelFactory.TlsHandler(path));
        Assert.Contains("kein PEM-Zertifikat", ex.Message);
    }

    [Fact]
    public void Create_WithoutAnchor_BuildsPlaintextChannel()
    {
        var plain = new Config("feed:9090", "", "", "stream", "", "", "", "", 0, "", "", 1, "", "", 0, 0, 0, 0);
        using var channel = ChannelFactory.Create(plain);
        Assert.Equal("feed:9090", channel.Target);
    }

    [Fact]
    public void Create_WithAnchor_ReadsTheFile()
    {
        using var cert = NewCertificate();
        var cfg = new Config("feed:9090", "", "", "stream", "", "", "", "", 0, "", "", 1, "", "", 0, 0, 0, 0, CaFile: WritePem(cert));
        using var channel = ChannelFactory.Create(cfg);
        Assert.Equal("feed:9090", channel.Target);
        Assert.Throws<ArgumentException>(() => ChannelFactory.Create(cfg with { CaFile = Path.Combine(_dir, "fehlt.pem") }));
    }
}
