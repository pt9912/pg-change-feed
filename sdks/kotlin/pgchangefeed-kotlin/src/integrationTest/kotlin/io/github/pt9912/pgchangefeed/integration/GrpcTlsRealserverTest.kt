package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.runBlocking
import java.net.URI
import kotlin.test.Test

/**
 * Real-server TLS phase for the gRPC stream client: with the trust anchor the
 * SDK opens the server stream over TLS and receives a change committed
 * afterwards (the runner reads its `change_id` back from `cdc.changes`);
 * without an anchor, with a foreign anchor and with a server name outside the
 * certificate the stream fails the TLS check.
 */
class GrpcTlsRealserverTest {
    private val address get() = URI("https://${PhaseEnvironment.grpcAddr}")

    @Test
    fun receivesACommittedChangeOverTlsWithTheTrustAnchor() = runBlocking {
        val options = PgChangeFeedClientOptions(
            address, PhaseEnvironment.apiToken, TlsScenarios.anchor(PhaseEnvironment.tlsCaFile),
        )
        PgChangeFeedGrpcClient(options).use { client ->
            println("READY")
            System.out.flush()

            val received = client.streamChanges()
                .firstOrNull { change ->
                    change.table == PhaseEnvironment.table &&
                        change.operation == "INSERT" &&
                        (change.newImage?.toStringUtf8() ?: "").contains(PhaseEnvironment.sentinel)
                }
                ?: throw IllegalStateException(
                    "kein Change mit dem Sentinel ${PhaseEnvironment.sentinel} innerhalb der Frist empfangen",
                )

            println(
                "RECEIVED change_id=${received.changeId} table=${received.table} " +
                    "operation=${received.operation} new_image=${received.newImage?.toStringUtf8()}",
            )
            System.out.flush()
        }
    }

    @Test
    fun streamsWithoutAnchorWithForeignAnchorAndWithAWrongServerNameFailTheTlsCheck() {
        TlsScenarios.expectRefusals(address, PhaseEnvironment.apiToken) { options ->
            PgChangeFeedGrpcClient(options).use { client ->
                runBlocking { client.streamChanges().firstOrNull() }
            }
        }
    }
}
