package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient
import java.net.URI
import java.net.http.HttpClient
import kotlin.test.Test
import kotlin.test.assertTrue

/**
 * Real-server phase for the SSE stream client: opens `GET /changes/stream`
 * against the running feed container and receives a change committed
 * afterwards; a second open with an unknown token is rejected with HTTP status
 * 401, as a typed
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedUnauthorizedException] — the
 * same set the HTTP client throws.
 */
class SseRealserverTest {
    @Test
    fun receivesACommittedChangeOverTheStream() {
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

        assertTrue(received.changeId.isNotEmpty(), "change_id leer")
        assertTrue(received.transactionId.isNotEmpty(), "transaction_id leer")
        assertTrue(received.sourceTableId.isNotEmpty(), "source_table_id leer")
        assertTrue(received.schemaVersion.isNotEmpty(), "schema_version leer")
        assertTrue(received.operation == "INSERT", "operation: ${received.operation}")
        assertTrue(received.oldImage == null || received.oldImage!!.isJsonNull, "ein INSERT trägt kein Alt-Bild (gson trägt JSON-null als JsonNull)")
        assertTrue(received.table == PhaseEnvironment.table, "table: ${received.table}")
        assertTrue(received.schema == "public", "schema: ${received.schema}")
        assertTrue(
            (received.newImage.toString()).contains(PhaseEnvironment.sentinel),
            "new_image ohne Sentinel",
        )

        println(
            "RECEIVED change_id=${received.changeId} table=${received.table} " +
                "operation=${received.operation} new_image=${received.newImage}",
        )
        System.out.flush()
    }

    @Test
    fun openWithUnknownTokenIsRejected401() {
        val httpClient = java.net.http.HttpClient.newHttpClient()
        val options = PgChangeFeedClientOptions(URI(PhaseEnvironment.httpAddr), "no-such-token")
        val client = PgChangeFeedSseClient(httpClient, options)

        val exception = runCatching {
            client.streamChanges().firstOrNull()
        }.exceptionOrNull() ?: throw IllegalStateException("Der Aufruf endete ohne Ausnahme, wollen 401")

        assertTrue(
            exception is io.github.pt9912.pgchangefeed.http.PgChangeFeedUnauthorizedException &&
                exception.statusCode == 401,
            "wollen 401, erhalten: ${exception::class.simpleName}: ${exception.message}",
        )
        println("REJECTED status=401")
        System.out.flush()
    }
}