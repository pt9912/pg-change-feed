package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient
import java.net.URI
import kotlin.test.Test

/**
 * Real-server TLS phase for the SSE stream client: with the trust anchor the
 * SDK opens `GET /changes/stream` over TLS and receives a change committed
 * afterwards (the runner reads its `change_id` back from `cdc.changes`);
 * without an anchor, with a foreign anchor and with a server name outside the
 * certificate the stream fails the TLS check.
 */
class SseTlsRealserverTest {
    private val address = URI(PhaseEnvironment.httpAddr)

    @Test
    fun receivesACommittedChangeOverTlsWithTheTrustAnchor() {
        val options = PgChangeFeedClientOptions(
            address, PhaseEnvironment.apiToken, TlsScenarios.anchor(PhaseEnvironment.tlsCaFile),
        )
        val client = PgChangeFeedSseClient(options)
        println("READY")
        System.out.flush()

        val received = client.streamChanges()
            .firstOrNull { change ->
                change.table == PhaseEnvironment.table &&
                    change.operation == "INSERT" &&
                    (change.newImage?.toString() ?: "").contains(PhaseEnvironment.sentinel)
            }
            ?: throw IllegalStateException(
                "kein Event mit dem Sentinel ${PhaseEnvironment.sentinel} innerhalb der Frist empfangen",
            )

        println(
            "RECEIVED change_id=${received.changeId} table=${received.table} " +
                "operation=${received.operation} new_image=${received.newImage}",
        )
        System.out.flush()
    }

    @Test
    fun streamsWithoutAnchorWithForeignAnchorAndWithAWrongServerNameFailTheTlsCheck() {
        TlsScenarios.expectRefusals(address, PhaseEnvironment.apiToken) { options ->
            PgChangeFeedSseClient(options).streamChanges().firstOrNull()
        }
    }
}
