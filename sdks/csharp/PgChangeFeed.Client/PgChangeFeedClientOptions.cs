namespace PgChangeFeed.Client;

/// <summary>
/// Connection configuration shared by the PG Change Feed clients: the address
/// of the server and the bearer token. The HTTP, gRPC, SSE and NATS clients
/// all authenticate with such a token against a single server address. What a
/// client does with them (which HTTP paths, which gRPC call) is up to the
/// client itself; the options carry no retry or backoff policy.
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

    /// <summary>Creates the options; the token must not be null, empty or whitespace.</summary>
    /// <param name="address">The address of the server.</param>
    /// <param name="apiToken">The bearer token.</param>
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
