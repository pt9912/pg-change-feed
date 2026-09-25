package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.nats.PgChangeFeedNatsStreamClient
import java.net.URI
import kotlin.test.Test
import kotlin.test.assertTrue

/**
 * Real-server phase for the NATS stream client: subscribes the source's
 * namespace against the running feed container and receives a change
 * committed afterwards as a full JSON event; a second connect with a wrong
 * token is rejected by the NATS server — the rejection throws synchronously
 * from the client constructor (real `io.nats.client` behavior), carrying the
 * server's `Authorization` wording.
 */
class NatsRealserverTest {
    @Test
    fun receivesACommittedChangeOverTheStream() {
        val options = PgChangeFeedClientOptions(URI(PhaseEnvironment.natsUrl), PhaseEnvironment.natsStreamToken)
        PgChangeFeedNatsStreamClient(options).use { client ->
            val subject = PgChangeFeedNatsStreamClient.buildSourceSubject(PhaseEnvironment.sourceId)
            println("READY")
            System.out.flush()

            val received = client.streamChanges(subject)
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
    }

    @Test
    fun connectWithWrongTokenIsRejectedByTheServer() {
        val options = PgChangeFeedClientOptions(URI(PhaseEnvironment.natsUrl), "wrong-token")

        val exception = runCatching {
            PgChangeFeedNatsStreamClient(options)
        }.exceptionOrNull() ?: throw IllegalStateException("Der Connect endete ohne Ausnahme, wollen Ablehnung")

        // Die Ablehnungsursache (Server-Rohwortlaut "-ERR Authorization
        // Violation") trägt die Ausnahme-Kette (toString schließt die
        // inneren Meldungen ein); der io.nats-JVM-Wrapper hält sie in eine
        // eigene Connect-Meldung.
        assertTrue(
            exception.toString().contains("Authorization"),
            "wollen die Verbindungsablehnung, erhalten: ${exception::class.simpleName}: ${exception.message}",
        )
        println("REJECTED token-rejected: $exception")
        System.out.flush()
    }
}