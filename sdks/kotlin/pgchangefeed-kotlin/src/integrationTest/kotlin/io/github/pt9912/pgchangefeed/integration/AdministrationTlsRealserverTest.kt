package io.github.pt9912.pgchangefeed.integration

import cdc.administration.v1.AdministrationOuterClass
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedAdministrationClient
import kotlinx.coroutines.runBlocking
import java.net.URI
import java.time.LocalDateTime
import java.time.format.DateTimeFormatter
import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * Real-server TLS phase for the gRPC administration client: with the trust
 * anchor the SDK registers a disposable consumer with the admin token and lists
 * tables with the reader token over TLS (the runner reads the registration back
 * from `cdc.consumer`); without an anchor, with a foreign anchor and with a
 * server name outside the certificate the same call fails the TLS check.
 */
class AdministrationTlsRealserverTest {
    private val address get() = URI("https://${PhaseEnvironment.grpcAddr}")

    private fun listTablesRequest() = AdministrationOuterClass.ListTablesRequest.newBuilder()
        .setSource(PhaseEnvironment.sourceId)
        .setPublication(PhaseEnvironment.httpPublication)
        .build()

    @Test
    fun registersAConsumerAndListsTablesOverTlsWithTheTrustAnchor() = runBlocking {
        println("READY")
        System.out.flush()

        val anchor = TlsScenarios.anchor(PhaseEnvironment.tlsCaFile)
        val consumerId = "kotlin-sdk-tls-admin-" + LocalDateTime.now().format(DateTimeFormatter.ofPattern("yyyyMMddHHmmss"))
        val registered = PgChangeFeedAdministrationClient(
            PgChangeFeedClientOptions(address, PhaseEnvironment.adminToken, anchor),
        ).use { admin ->
            admin.registerConsumer(
                AdministrationOuterClass.RegisterConsumerRequest.newBuilder()
                    .setConsumerId(consumerId)
                    .setName("Kotlin SDK TLS $consumerId")
                    .build(),
            )
        }
        assertEquals(consumerId, registered.consumerId, "Registrierung traegt eine andere Identitaet zurueck")

        val tables = PgChangeFeedAdministrationClient(
            PgChangeFeedClientOptions(address, PhaseEnvironment.readerToken, anchor),
        ).use { reader -> reader.listTables(listTablesRequest()) }

        println("RECEIVED consumer_id=${registered.consumerId} tables=${tables.tablesCount}")
        System.out.flush()
    }

    @Test
    fun callsWithoutAnchorWithForeignAnchorAndWithAWrongServerNameFailTheTlsCheck() {
        TlsScenarios.expectRefusals(address, PhaseEnvironment.readerToken) { options ->
            PgChangeFeedAdministrationClient(options).use { client ->
                runBlocking { client.listTables(listTablesRequest()) }
            }
        }
    }
}
