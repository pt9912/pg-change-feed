package io.github.pt9912.pgchangefeed.integration

import com.google.gson.JsonParser
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient
import java.net.URI
import kotlin.test.Test
import kotlin.test.assertFalse
import kotlin.test.assertTrue
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.runBlocking

/**
 * Real-server rule phase for the gRPC stream client: opens the server stream
 * against a table carrying an active `rename_column` rule and receives a
 * change committed afterwards; the row image arrives with the renamed key
 * and without the source key, read through the client's opaque wire form
 * (the raw bytes, parsed as JSON here only to assert on the property under
 * test — the client itself never interprets the row image shape).
 */
class GrpcRuleRealserverTest {
    @Test
    fun receivesARenamedRowImageKeyOverTheStream() = runBlocking {
        val options = PgChangeFeedClientOptions(URI("http://${PhaseEnvironment.grpcAddr}"), PhaseEnvironment.apiToken)
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

            val newImage = JsonParser.parseString(received.newImage?.toStringUtf8() ?: "{}").asJsonObject
            assertTrue(newImage.has(PhaseEnvironment.ruleTargetKey), "Zielschlüssel fehlt: $newImage")
            assertTrue(newImage.get(PhaseEnvironment.ruleTargetKey).asString == PhaseEnvironment.sentinel)
            assertFalse(newImage.has(PhaseEnvironment.ruleSourceKey), "Quellschlüssel noch vorhanden: $newImage")

            println(
                "RECEIVED change_id=${received.changeId} table=${received.table} " +
                    "operation=${received.operation} new_image=${received.newImage?.toStringUtf8()}",
            )
            System.out.flush()
        }
    }
}
