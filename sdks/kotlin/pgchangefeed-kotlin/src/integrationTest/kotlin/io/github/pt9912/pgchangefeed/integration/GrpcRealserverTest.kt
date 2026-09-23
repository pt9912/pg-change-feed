package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient
import io.grpc.Status
import io.grpc.StatusException
import java.net.URI
import kotlin.test.Test
import kotlin.test.assertTrue
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.runBlocking

/**
 * Realserver phase for the gRPC stream surface (SPEC-020): opens the server
 * stream against the running feed container and receives a change committed
 * afterwards; a second open with an unknown token is rejected with gRPC
 * status `Unauthenticated` (SPEC-020 Negative), not swallowed as an empty
 * stream. The commit-to-delivery window is fire-and-forget — the runner
 * commits a bounded sequence of unique rows until one arrives; this test
 * keeps receiving until it sees its sentinel.
 */
class GrpcRealserverTest {
    private fun options(token: String) =
        PgChangeFeedClientOptions(URI("http://${PhaseEnvironment.grpcAddr}"), token)

    @Test
    fun receivesACommittedChangeOverTheStream() = runBlocking {
        PgChangeFeedGrpcClient(options(PhaseEnvironment.apiToken)).use { client ->
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

            assertTrue(received.changeId.isNotEmpty(), "change_id leer")
            assertTrue(received.transactionId.isNotEmpty(), "transaction_id leer")
            assertTrue(received.sourceTableId.isNotEmpty(), "source_table_id leer")
            assertTrue(received.schemaVersion.isNotEmpty(), "schema_version leer")
            assertTrue(received.operation == "INSERT", "operation: ${received.operation}")
            assertTrue(received.table == PhaseEnvironment.table, "table: ${received.table}")
            assertTrue(received.schema == "public", "schema: ${received.schema}")
            assertTrue(
                (received.newImage?.toStringUtf8() ?: "").contains(PhaseEnvironment.sentinel),
                "new_image ohne Sentinel",
            )

            println(
                "RECEIVED change_id=${received.changeId} table=${received.table} " +
                    "operation=${received.operation} new_image=${received.newImage?.toStringUtf8()}",
            )
            System.out.flush()
        }
    }

    @Test
    fun openWithUnknownTokenIsRejectedUnauthenticated() {
        PgChangeFeedGrpcClient(options("no-such-token")).use { client ->
            val exception = runBlocking {
                runCatching {
                    client.streamChanges().firstOrNull()
                }.exceptionOrNull()
            } ?: throw IllegalStateException("Der Stream öffnete ohne gültiges Token, wollen Unauthenticated")

            assertTrue(
                exception is StatusException && exception.status.code == Status.Code.UNAUTHENTICATED,
                "wollen Unauthenticated, erhalten: ${exception::class.simpleName}: ${exception.message}",
            )
            println("REJECTED code=Unauthenticated")
            System.out.flush()
        }
    }
}