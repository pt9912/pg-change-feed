using System.Net.Security;
using System.Security.Cryptography.X509Certificates;

namespace PgChangeFeed.Client;

/// <summary>
/// Builds the connections the clients open themselves. Without a trust anchor
/// the handler keeps the default trust of the runtime; with one, the server
/// certificate must chain to exactly the anchor certificates, and the server
/// name and the validity period are checked as always.
/// </summary>
internal static class TlsTransport
{
    /// <summary>
    /// Fails when the options name a trust anchor but the address does not ask
    /// for TLS: an anchor never goes with a plaintext connection.
    /// </summary>
    internal static void RequireTlsAddressForAnchor(PgChangeFeedClientOptions options)
    {
        if (options.TrustAnchors is not null
            && !string.Equals(options.Address.Scheme, Uri.UriSchemeHttps, StringComparison.OrdinalIgnoreCase))
        {
            throw new ArgumentException(
                $"A trust anchor requires an https address, got '{options.Address}'.", nameof(options));
        }
    }

    internal static SocketsHttpHandler CreateHandler(PgChangeFeedClientOptions options, bool forGrpc)
    {
        RequireTlsAddressForAnchor(options);
        var handler = new SocketsHttpHandler();
        if (forGrpc)
        {
            // The connection settings gRPC channels use when they build their own handler.
            handler.PooledConnectionIdleTimeout = Timeout.InfiniteTimeSpan;
            handler.KeepAlivePingDelay = TimeSpan.FromSeconds(60);
            handler.KeepAlivePingTimeout = TimeSpan.FromSeconds(30);
            handler.EnableMultipleHttp2Connections = true;
        }

        var anchors = options.TrustAnchors;
        if (anchors is not null)
        {
            handler.SslOptions.RemoteCertificateValidationCallback =
                (_, certificate, chain, errors) => ChainsToAnchors(anchors, certificate, chain, errors);
        }

        return handler;
    }

    internal static HttpClient CreateHttpClient(PgChangeFeedClientOptions options, TimeSpan timeout, bool forGrpc = false)
        => new(CreateHandler(options, forGrpc), disposeHandler: true) { Timeout = timeout };

    /// <summary>
    /// The server name and a presented certificate are required as the
    /// platform reports them; the chain is built again against the anchors
    /// only, so the operating system's roots do not count.
    /// </summary>
    private static bool ChainsToAnchors(
        X509Certificate2Collection anchors,
        X509Certificate? certificate,
        X509Chain? presented,
        SslPolicyErrors errors)
    {
        if (certificate is null)
        {
            return false;
        }

        if ((errors & (SslPolicyErrors.RemoteCertificateNameMismatch
            | SslPolicyErrors.RemoteCertificateNotAvailable)) != 0)
        {
            return false;
        }

        using var leaf = new X509Certificate2(certificate);
        using var chain = new X509Chain();
        chain.ChainPolicy.TrustMode = X509ChainTrustMode.CustomRootTrust;
        chain.ChainPolicy.CustomTrustStore.AddRange(anchors);
        chain.ChainPolicy.RevocationMode = X509RevocationMode.NoCheck;
        if (presented is not null)
        {
            foreach (var element in presented.ChainElements)
            {
                chain.ChainPolicy.ExtraStore.Add(element.Certificate);
            }
        }

        return chain.Build(leaf);
    }
}
