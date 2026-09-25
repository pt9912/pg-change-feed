package io.github.pt9912.pgchangefeed

import java.net.URI

/**
 * Connection configuration shared by the PG Change Feed clients: the address
 * of the server and the bearer token. The HTTP, gRPC, SSE and NATS clients
 * all authenticate with such a token against a single server address. What a
 * client does with them (which HTTP paths, which gRPC call) is up to the
 * client itself; the options carry no retry or backoff policy.
 *
 * @property address the address of the PG Change Feed server, in the form the
 *   client needs: the HTTP base URL for the HTTP and SSE clients, the gRPC
 *   endpoint for the gRPC client, the NATS URL for the NATS client.
 * @property apiToken the bearer token sent as the authorization credential;
 *   must not be blank.
 */
class PgChangeFeedClientOptions(
    val address: URI,
    val apiToken: String,
) {
    init {
        require(apiToken.isNotBlank()) { "API token must not be blank." }
    }
}
