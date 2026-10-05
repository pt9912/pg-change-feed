using System.Security.Cryptography.X509Certificates;

namespace PgChangeFeed.Client;

/// <summary>
/// Connection configuration shared by the PG Change Feed clients: the address
/// of the server, the bearer token and, optionally, a trust anchor for TLS.
/// The HTTP, gRPC, SSE and NATS clients all authenticate with such a token
/// against a single server address. What a client does with them (which HTTP
/// paths, which gRPC call) is up to the client itself; the options carry no
/// retry or backoff policy.
/// </summary>
public sealed class PgChangeFeedClientOptions
{
    /// <summary>
    /// The address of the PG Change Feed server, in the form the client needs:
    /// the HTTP base URL for the HTTP and SSE clients, the gRPC endpoint for
    /// the gRPC client, the NATS URL for the NATS client.
    /// </summary>
    public Uri Address { get; }

    /// <summary>
    /// The bearer token sent as the authorization credential.
    /// </summary>
    public string ApiToken { get; }

    /// <summary>
    /// The path of a PEM file with the certificates a TLS connection trusts
    /// (the certificate of the issuer or of the server), or <c>null</c> when
    /// the trust anchors of the operating system apply. It is used by the
    /// HTTP, SSE and gRPC clients that build their own connection (the
    /// constructors that take only the options) and requires an <c>https</c>
    /// address there. Chain, validity period and server name are always
    /// checked; there is no switch that turns the check off.
    /// </summary>
    public string? TrustAnchorFile { get; }

    internal X509Certificate2Collection? TrustAnchors { get; }

    /// <summary>Creates the options; the token must not be null, empty or whitespace.</summary>
    /// <param name="address">The address of the server.</param>
    /// <param name="apiToken">The bearer token.</param>
    public PgChangeFeedClientOptions(Uri address, string apiToken)
        : this(address, apiToken, null)
    {
    }

    /// <summary>
    /// Creates the options with a trust anchor for TLS. A path that cannot be
    /// read, or a file without a PEM certificate, throws an
    /// <see cref="ArgumentException"/> here, not with the first request.
    /// </summary>
    /// <param name="address">The address of the server.</param>
    /// <param name="apiToken">The bearer token.</param>
    /// <param name="trustAnchorFile">The path of a PEM file with the trusted certificates, or <c>null</c>.</param>
    public PgChangeFeedClientOptions(Uri address, string apiToken, string? trustAnchorFile)
    {
        ArgumentNullException.ThrowIfNull(address);
        if (string.IsNullOrWhiteSpace(apiToken))
        {
            throw new ArgumentException("API token must not be null or empty.", nameof(apiToken));
        }

        Address = address;
        ApiToken = apiToken;
        TrustAnchorFile = trustAnchorFile;
        if (trustAnchorFile is not null)
        {
            TrustAnchors = LoadTrustAnchors(trustAnchorFile);
        }
    }

    private static X509Certificate2Collection LoadTrustAnchors(string path)
    {
        var anchors = new X509Certificate2Collection();
        try
        {
            anchors.ImportFromPemFile(path);
        }
        catch (Exception error) when (error is IOException or UnauthorizedAccessException
            or System.Security.Cryptography.CryptographicException or ArgumentException)
        {
            throw new ArgumentException(
                $"The trust anchor file '{path}' cannot be read as PEM certificates: {error.Message}",
                "trustAnchorFile", error);
        }

        if (anchors.Count == 0)
        {
            throw new ArgumentException(
                $"The trust anchor file '{path}' contains no PEM certificate.", "trustAnchorFile");
        }

        return anchors;
    }
}
