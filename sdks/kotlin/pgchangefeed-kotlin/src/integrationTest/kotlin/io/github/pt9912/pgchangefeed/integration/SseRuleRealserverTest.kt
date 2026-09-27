package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient
import java.net.URI
import kotlin.test.Test
import kotlin.test.assertFalse
import kotlin.test.assertTrue

/**
 * Real-server rule phase for the SSE stream client: opens
 * `GET /changes/stream` against a table carrying an active `rename_column`
 * rule and receives a change committed afterwards; the row image arrives
 * with the renamed key and without the source key, read through the
 * client's opaque `JsonElement` model.
 */
class SseRuleRealserverTest {
    @Test
    fun receivesARenamedRowImageKeyOverTheStream() {
        val httpClient = java.net.http.HttpClient.newHttpClient()
        val options = PgChangeFeedClientOptions(URI(PhaseEnvironment.httpAddr), PhaseEnvironment.apiToken)
        val client = PgChangeFeedSseClient(httpClient, options)
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

        val newImage = received.newImage?.asJsonObject
            ?: throw IllegalStateException("new_image fehlt")
        assertTrue(newImage.has(PhaseEnvironment.ruleTargetKey), "Zielschlüssel fehlt: $newImage")
        assertTrue(newImage.get(PhaseEnvironment.ruleTargetKey).asString == PhaseEnvironment.sentinel)
        assertFalse(newImage.has(PhaseEnvironment.ruleSourceKey), "Quellschlüssel noch vorhanden: $newImage")

        println(
            "RECEIVED change_id=${received.changeId} table=${received.table} " +
                "operation=${received.operation} new_image=${received.newImage}",
        )
        System.out.flush()
    }
}
