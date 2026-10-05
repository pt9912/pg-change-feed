package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient
import io.github.pt9912.pgchangefeed.http.model.RegisterConsumerRequest
import java.net.URI
import java.time.LocalDateTime
import java.time.format.DateTimeFormatter
import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * Real-server TLS phase for the HTTP API client: the feed container serves HTTP
 * over TLS with a certificate the Java runtime does not know. With the trust
 * anchor the SDK registers a disposable consumer with the admin token and lists
 * tables with the reader token (the runner reads the registration back from
 * `cdc.consumer`); without an anchor, with a foreign anchor and with a server
 * name outside the certificate the same call fails the TLS check.
 */
class HttpTlsRealserverTest {
    private val address = URI(PhaseEnvironment.httpAddr)

    @Test
    fun registersAConsumerAndListsTablesOverTlsWithTheTrustAnchor() {
        println("READY")
        System.out.flush()

        val anchor = TlsScenarios.anchor(PhaseEnvironment.tlsCaFile)
        val admin = PgChangeFeedHttpClient(PgChangeFeedClientOptions(address, PhaseEnvironment.adminToken, anchor))
        val consumerId = "kotlin-sdk-tls-" + LocalDateTime.now().format(DateTimeFormatter.ofPattern("yyyyMMddHHmmss"))
        val registered = admin.registerConsumer(RegisterConsumerRequest(consumerId, "Kotlin SDK TLS $consumerId"))
        assertEquals(consumerId, registered.consumerId, "Registrierung traegt eine andere Identitaet zurueck")

        val reader = PgChangeFeedHttpClient(PgChangeFeedClientOptions(address, PhaseEnvironment.readerToken, anchor))
        val tables = reader.listTables(PhaseEnvironment.sourceId, PhaseEnvironment.httpPublication)

        println("RECEIVED consumer_id=${registered.consumerId} tables=${tables.tables.size}")
        System.out.flush()
    }

    @Test
    fun callsWithoutAnchorWithForeignAnchorAndWithAWrongServerNameFailTheTlsCheck() {
        TlsScenarios.expectRefusals(address, PhaseEnvironment.readerToken) { options ->
            PgChangeFeedHttpClient(options).listTables(PhaseEnvironment.sourceId, PhaseEnvironment.httpPublication)
        }
    }
}
