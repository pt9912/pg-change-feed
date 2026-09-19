namespace PgChangeFeed.Client;

/// <summary>
/// Shared connection configuration for PG Change Feed client surfaces. The
/// HTTP API (SPEC-018) and the gRPC change stream (SPEC-020) both authenticate
/// with a bearer token against a single server address — this is the one
/// unstrittige, shared configuration denominator identified while writing
/// this project skeleton (ADR-0106 Festlegung 1: "HTTP und gRPC teilen
/// Auth-Header-Form (Bearer Token) und Grundkonfiguration (Adresse,
/// Token)"). Surface-specific behavior (which HTTP paths, which gRPC stub,
/// retry/backoff policy) is deliberately NOT part of this class — it is
/// added by the follow-up slices that build the actual client surfaces.
/// </summary>
public sealed class PgChangeFeedClientOptions
{
    /// <summary>
    /// The base address of the PG Change Feed server (HTTP or gRPC endpoint,
    /// depending on the surface that consumes these options).
    /// </summary>
    public Uri Address { get; }

    /// <summary>
    /// The bearer token sent as an authorization credential (SPEC-018).
    /// </summary>
    public string ApiToken { get; }

    public PgChangeFeedClientOptions(Uri address, string apiToken)
    {
        ArgumentNullException.ThrowIfNull(address);
        if (string.IsNullOrWhiteSpace(apiToken))
        {
            throw new ArgumentException("API token must not be null or empty.", nameof(apiToken));
        }

        Address = address;
        ApiToken = apiToken;
    }
}
