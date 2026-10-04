package io.github.pt9912.pgchangefeed.integration

import cdc.administration.v1.AdministrationOuterClass
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedAdministrationClient
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcNotFoundException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcPermissionDeniedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedForbiddenException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient
import io.github.pt9912.pgchangefeed.http.PgChangeFeedNotFoundException
import io.github.pt9912.pgchangefeed.http.model.EnableTableRequest
import io.grpc.Status
import kotlinx.coroutines.runBlocking
import java.net.URI
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

/**
 * Real-server phase against a server older than the message codes. The same
 * calls as [ErrorCodeRealserverTest] end with the same typed errors, but the
 * server sends no `code` in the HTTP error body, no `google.rpc.ErrorInfo` in
 * the gRPC status detail and no `error_code` in the diagnose heartbeat: every
 * `messageCode` stays `null`, the error text is unchanged and the SDK does not
 * fail. The error state of the diagnose report is written by the runner
 * without an `error_code`.
 */
class ErrorCodeAltServerTest {
    private val missingTable = "sdk_error_code_missing_table"

    private fun httpEnableRequest() = EnableTableRequest(
        PhaseEnvironment.sourceId, "public", missingTable, "sdk-error-code", "sdk-error-code-v1", 1,
        PhaseEnvironment.httpPublication,
    )

    private fun grpcEnableRequest(): AdministrationOuterClass.EnableTableRequest =
        AdministrationOuterClass.EnableTableRequest.newBuilder()
            .setSource(PhaseEnvironment.sourceId)
            .setSchema("public")
            .setTable(missingTable)
            .setTableId("sdk-error-code")
            .setSchemaVersionId("sdk-error-code-v1")
            .setVersion(1)
            .setPublication(PhaseEnvironment.httpPublication)
            .build()

    private fun diagnoseRequest(): AdministrationOuterClass.DiagnoseRequest =
        AdministrationOuterClass.DiagnoseRequest.newBuilder().setSource(PhaseEnvironment.sourceId).build()

    private fun httpClientFor(token: String) = PgChangeFeedHttpClient(
        java.net.http.HttpClient.newHttpClient(),
        PgChangeFeedClientOptions(URI(PhaseEnvironment.httpAddr), token),
    )

    private fun administrationClientFor(token: String) = PgChangeFeedAdministrationClient(
        PgChangeFeedClientOptions(URI("http://${PhaseEnvironment.grpcAddr}"), token),
    )

    @Test
    fun anOlderServerCarriesNoMessageCodeAndTheSdkStaysIntact() {
        println("READY")
        System.out.flush()

        val httpException = runCatching { httpClientFor(PhaseEnvironment.adminToken).enableTable(httpEnableRequest()) }
            .exceptionOrNull()
        val grpcException = administrationClientFor(PhaseEnvironment.adminToken).use { client ->
            runCatching { runBlocking { client.enableTable(grpcEnableRequest()) } }.exceptionOrNull()
        }
        assertTrue(httpException is PgChangeFeedNotFoundException, "HTTP: wollen 404, erhalten: $httpException")
        assertTrue(grpcException is PgChangeFeedGrpcNotFoundException, "gRPC: wollen NOT_FOUND, erhalten: $grpcException")
        assertEquals(404, httpException.statusCode)
        assertEquals(Status.Code.NOT_FOUND, grpcException.statusCode)
        assertNull(httpException.messageCode)
        assertNull(grpcException.messageCode)
        assertTrue(!httpException.message.isNullOrEmpty())
        assertTrue(!grpcException.message.isNullOrEmpty())
        assertEquals(PhaseEnvironment.altServerHttpText, httpException.message)

        val httpForbidden = runCatching { httpClientFor(PhaseEnvironment.readerToken).enableTable(httpEnableRequest()) }
            .exceptionOrNull()
        val grpcForbidden = administrationClientFor(PhaseEnvironment.readerToken).use { client ->
            runCatching { runBlocking { client.enableTable(grpcEnableRequest()) } }.exceptionOrNull()
        }
        assertTrue(httpForbidden is PgChangeFeedForbiddenException, "HTTP: wollen 403, erhalten: $httpForbidden")
        assertTrue(
            grpcForbidden is PgChangeFeedGrpcPermissionDeniedException,
            "gRPC: wollen PERMISSION_DENIED, erhalten: $grpcForbidden",
        )
        assertEquals(403, httpForbidden.statusCode)
        assertEquals(Status.Code.PERMISSION_DENIED, grpcForbidden.statusCode)
        assertNull(httpForbidden.messageCode)
        assertNull(grpcForbidden.messageCode)

        administrationClientFor(PhaseEnvironment.adminToken).use { client ->
            val normal = runBlocking { client.diagnose(diagnoseRequest()) }
            assertTrue(normal.heartbeat.known)
            assertEquals("", normal.heartbeat.errorClass)
            assertEquals("", normal.heartbeat.errorCode)
            println("NORMAL_DONE")
            System.out.flush()

            // The runner writes the error state (error_class without
            // error_code) after NORMAL_DONE; the periodic heartbeat resets
            // it, so the call is repeated until the report shows it.
            var failing: AdministrationOuterClass.HeartbeatStatus? = null
            val deadline = System.nanoTime() + 60_000_000_000L
            while (System.nanoTime() < deadline) {
                val report = runBlocking { client.diagnose(diagnoseRequest()) }
                if (report.heartbeat.errorClass == "schema") {
                    failing = report.heartbeat
                    break
                }
                Thread.sleep(200)
            }
            assertNotNull(failing, "kein Fehlerzustand im Diagnose-Bericht innerhalb von 60 s")
            assertEquals("schema", failing.errorClass)
            assertEquals("", failing.errorCode)
        }

        println("RECEIVED code=none http=404 grpc=NOT_FOUND text=${httpException.message?.length} diag_error_code=leer")
        System.out.flush()
    }
}
