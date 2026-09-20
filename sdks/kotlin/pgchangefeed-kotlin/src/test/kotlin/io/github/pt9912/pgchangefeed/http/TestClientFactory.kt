package io.github.pt9912.pgchangefeed.http

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import java.net.URI

/**
 * Builds a [PgChangeFeedHttpClient] wired to a [FakeHttpTransport] via the
 * `internal` transport constructor — every test in this module stays
 * network-free (no real server, no real socket, analogous to the C#
 * sibling's `TestClientFactory` + `FakeHttpMessageHandler`).
 */
internal object TestClientFactory {
    val serverAddress: URI = URI.create("http://example.invalid:8080")

    fun create(
        apiToken: String = "test-token",
        responder: (TransportRequest) -> TransportResponse,
    ): Pair<PgChangeFeedHttpClient, FakeHttpTransport> {
        val transport = FakeHttpTransport(responder)
        val options = PgChangeFeedClientOptions(serverAddress, apiToken)
        return PgChangeFeedHttpClient(transport, options) to transport
    }
}
