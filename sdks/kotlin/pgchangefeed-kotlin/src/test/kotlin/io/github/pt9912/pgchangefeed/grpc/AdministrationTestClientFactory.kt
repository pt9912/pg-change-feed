package io.github.pt9912.pgchangefeed.grpc

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import java.net.URI

/**
 * Builds a [PgChangeFeedAdministrationClient] wired to a
 * [FakeAdministrationTransport] — every test in this file group stays
 * network-free (no real server, no real socket).
 */
internal object AdministrationTestClientFactory {
    val serverAddress: URI = URI("http://localhost:50051")

    fun create(transport: FakeAdministrationTransport, apiToken: String = "test-token"): PgChangeFeedAdministrationClient =
        PgChangeFeedAdministrationClient(transport, PgChangeFeedClientOptions(serverAddress, apiToken))
}
