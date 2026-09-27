package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient
import io.github.pt9912.pgchangefeed.http.model.Change
import java.net.URI
import kotlin.test.Test
import kotlin.test.assertFalse
import kotlin.test.assertTrue

/**
 * Real-server rule phase for the HTTP API client: reads `GET /changes`
 * (source and table as filter) against a table carrying an active
 * `rename_column` rule; the row image of the change committed while the
 * test polls arrives with the renamed key and without the source key, read
 * through the client's opaque `JsonElement` model.
 */
class HttpRuleRealserverTest {
    @Test
    fun readsARenamedRowImageKeyOverTheApi() {
        val httpClient = java.net.http.HttpClient.newHttpClient()
        val options = PgChangeFeedClientOptions(URI(PhaseEnvironment.httpAddr), PhaseEnvironment.apiToken)
        val client = PgChangeFeedHttpClient(httpClient, options)
        println("READY")
        System.out.flush()

        val received = receiveSentinel(client)

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

    private fun receiveSentinel(client: PgChangeFeedHttpClient): Change {
        val deadline = System.currentTimeMillis() + 90_000
        while (System.currentTimeMillis() < deadline) {
            val response = client.readChanges(PhaseEnvironment.sourceId, "public", PhaseEnvironment.table)
            val match = response.changes.firstOrNull { change ->
                change.table == PhaseEnvironment.table &&
                    change.operation == "INSERT" &&
                    (change.newImage?.toString() ?: "").contains(PhaseEnvironment.sentinel)
            }
            if (match != null) {
                return match
            }
            Thread.sleep(500)
        }
        throw IllegalStateException(
            "kein Change mit dem Sentinel ${PhaseEnvironment.sentinel} innerhalb der Frist über GET /changes gelesen",
        )
    }
}
