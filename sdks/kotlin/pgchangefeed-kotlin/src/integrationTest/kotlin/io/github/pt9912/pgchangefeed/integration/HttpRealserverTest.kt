package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnauthorizedException
import io.github.pt9912.pgchangefeed.http.model.RegisterConsumerRequest
import java.net.URI
import java.net.http.HttpClient
import java.time.LocalDateTime
import java.time.format.DateTimeFormatter
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue

/**
 * Real-server phase for the HTTP API client: a roundtrip — the SDK registers a
 * disposable consumer with the admin token and lists tables with the reader
 * token; the runner checks the registration against the SQL read path
 * (`cdc.consumer`); a call with an unknown token is rejected with HTTP status
 * 401.
 */
class HttpRealserverTest {
    @Test
    fun registersAConsumerAndListsTablesOverTheApi() {
        val httpClient = java.net.http.HttpClient.newHttpClient()
        println("READY")
        System.out.flush()

        val adminOptions = PgChangeFeedClientOptions(URI(PhaseEnvironment.httpAddr), PhaseEnvironment.adminToken)
        val adminClient = PgChangeFeedHttpClient(httpClient, adminOptions)
        val consumerId = "kotlin-sdk-e2e-" +
            LocalDateTime.now().format(DateTimeFormatter.ofPattern("yyyyMMddHHmmss"))
        val registered = adminClient.registerConsumer(RegisterConsumerRequest(consumerId, "Kotlin SDK E2E $consumerId"))
        assertEquals(consumerId, registered.consumerId, "Registrierung tragt eine andere Identitaet zurueck")

        val readerOptions = PgChangeFeedClientOptions(URI(PhaseEnvironment.httpAddr), PhaseEnvironment.readerToken)
        val readerClient = PgChangeFeedHttpClient(httpClient, readerOptions)
        val tables = readerClient.listTables(PhaseEnvironment.sourceId, PhaseEnvironment.httpPublication)

        println("RECEIVED consumer_id=${registered.consumerId} tables=${tables.tables.size}")
        System.out.flush()
    }

    @Test
    fun callWithUnknownTokenIsRejected401() {
        val httpClient = java.net.http.HttpClient.newHttpClient()
        val options = PgChangeFeedClientOptions(URI(PhaseEnvironment.httpAddr), "no-such-token")
        val client = PgChangeFeedHttpClient(httpClient, options)

        val exception = runCatching {
            client.listTables(PhaseEnvironment.sourceId, PhaseEnvironment.httpPublication)
        }.exceptionOrNull() ?: throw IllegalStateException("Der Aufruf endete ohne Ausnahme, wollen 401")

        assertTrue(
            exception is PgChangeFeedUnauthorizedException && exception.statusCode == 401,
            "wollen 401, erhalten: ${exception::class.simpleName}: ${exception.message}",
        )
        println("REJECTED status=401")
        System.out.flush()
    }
}