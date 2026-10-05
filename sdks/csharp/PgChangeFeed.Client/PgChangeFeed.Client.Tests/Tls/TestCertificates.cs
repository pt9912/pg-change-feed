using System.Net;
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;

namespace PgChangeFeed.Client.Tests.Tls;

/// <summary>
/// Self-signed certificates created at run time for the TLS tests: nothing
/// is committed, and every PEM file lives in a temporary directory the test
/// removes. Each certificate is its own trust anchor (a CA certificate that
/// is also the server certificate).
/// </summary>
internal sealed class TestCertificates : IDisposable
{
    private readonly string _directory = Directory.CreateTempSubdirectory("pgcf-tls-").FullName;

    /// <summary>
    /// Creates a certificate for the given DNS names and, when
    /// <paramref name="includeLoopbackIp"/> is set, for 127.0.0.1. The
    /// certificate is valid from <paramref name="notBefore"/> to
    /// <paramref name="notAfter"/>.
    /// </summary>
    internal X509Certificate2 Create(
        string[] dnsNames,
        bool includeLoopbackIp = false,
        DateTimeOffset? notBefore = null,
        DateTimeOffset? notAfter = null)
    {
        using var key = ECDsa.Create(ECCurve.NamedCurves.nistP256);
        var request = new CertificateRequest($"CN=pgcf-test-{Guid.NewGuid():N}", key, HashAlgorithmName.SHA256);
        request.CertificateExtensions.Add(new X509BasicConstraintsExtension(true, false, 0, true));
        request.CertificateExtensions.Add(new X509KeyUsageExtension(
            X509KeyUsageFlags.DigitalSignature | X509KeyUsageFlags.KeyCertSign, true));
        request.CertificateExtensions.Add(new X509EnhancedKeyUsageExtension(
            [new Oid("1.3.6.1.5.5.7.3.1")], false));
        var names = new SubjectAlternativeNameBuilder();
        foreach (var name in dnsNames)
        {
            names.AddDnsName(name);
        }

        if (includeLoopbackIp)
        {
            names.AddIpAddress(IPAddress.Loopback);
        }

        request.CertificateExtensions.Add(names.Build());

        using var created = request.CreateSelfSigned(
            notBefore ?? DateTimeOffset.UtcNow.AddMinutes(-5),
            notAfter ?? DateTimeOffset.UtcNow.AddDays(1));
        // Re-import so that the private key is usable by the TLS stack.
        return X509CertificateLoader.LoadPkcs12(created.Export(X509ContentType.Pfx), null);
    }

    /// <summary>
    /// Creates a certificate authority and a server certificate for
    /// <paramref name="dnsName"/> that it issued. The server certificate is not
    /// its own anchor: the authority certificate is.
    /// </summary>
    internal (X509Certificate2 Server, X509Certificate2 Authority) CreateIssued(string dnsName)
    {
        using var authorityKey = ECDsa.Create(ECCurve.NamedCurves.nistP256);
        var authorityRequest = new CertificateRequest(
            $"CN=pgcf-test-ca-{Guid.NewGuid():N}", authorityKey, HashAlgorithmName.SHA256);
        authorityRequest.CertificateExtensions.Add(new X509BasicConstraintsExtension(true, false, 0, true));
        authorityRequest.CertificateExtensions.Add(new X509KeyUsageExtension(X509KeyUsageFlags.KeyCertSign, true));
        using var authority = authorityRequest.CreateSelfSigned(
            DateTimeOffset.UtcNow.AddMinutes(-5), DateTimeOffset.UtcNow.AddDays(1));

        using var serverKey = ECDsa.Create(ECCurve.NamedCurves.nistP256);
        var serverRequest = new CertificateRequest(
            $"CN=pgcf-test-server-{Guid.NewGuid():N}", serverKey, HashAlgorithmName.SHA256);
        serverRequest.CertificateExtensions.Add(new X509BasicConstraintsExtension(false, false, 0, true));
        serverRequest.CertificateExtensions.Add(new X509KeyUsageExtension(X509KeyUsageFlags.DigitalSignature, true));
        serverRequest.CertificateExtensions.Add(new X509EnhancedKeyUsageExtension([new Oid("1.3.6.1.5.5.7.3.1")], false));
        var names = new SubjectAlternativeNameBuilder();
        names.AddDnsName(dnsName);
        serverRequest.CertificateExtensions.Add(names.Build());
        using var issued = serverRequest.Create(
            authority, DateTimeOffset.UtcNow.AddMinutes(-5), DateTimeOffset.UtcNow.AddDays(1), RandomNumberGenerator.GetBytes(8));
        using var server = issued.CopyWithPrivateKey(serverKey);

        return (
            X509CertificateLoader.LoadPkcs12(server.Export(X509ContentType.Pfx), null),
            X509CertificateLoader.LoadCertificate(authority.Export(X509ContentType.Cert)));
    }

    /// <summary>Writes the certificates (public part only) as one PEM file and returns its path.</summary>
    internal string WritePem(params X509Certificate2[] certificates)
    {
        var path = Path.Combine(_directory, $"{Guid.NewGuid():N}.pem");
        File.WriteAllText(path, string.Concat(certificates.Select(c => c.ExportCertificatePem() + "\n")));
        return path;
    }

    /// <summary>Writes a file that holds no PEM certificate and returns its path.</summary>
    internal string WriteText(string text)
    {
        var path = Path.Combine(_directory, $"{Guid.NewGuid():N}.txt");
        File.WriteAllText(path, text);
        return path;
    }

    internal string MissingPath() => Path.Combine(_directory, "missing.pem");

    public void Dispose() => Directory.Delete(_directory, recursive: true);
}
