package io.github.pt9912.pgchangefeed.sse

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.http.TransportRequest
import java.net.URI

/**
 * Builds a [PgChangeFeedSseClient] wired to a [FakeSseTransport] via the
 * `internal` transport constructor — every test in this module stays
 * network-free (no real server, no real socket), analogous to
 * `io.github.pt9912.pgchangefeed.http.TestClientFactory`.
 */
internal object TestClientFactory {
    val serverAddress: URI = URI.create("http://example.invalid:8080")

    fun create(
        apiToken: String = "test-token",
        responder: (TransportRequest) -> SseTransportResponse,
    ): Pair<PgChangeFeedSseClient, FakeSseTransport> {
        val transport = FakeSseTransport(responder)
        val options = PgChangeFeedClientOptions(serverAddress, apiToken)
        return PgChangeFeedSseClient(transport, options) to transport
    }
}
