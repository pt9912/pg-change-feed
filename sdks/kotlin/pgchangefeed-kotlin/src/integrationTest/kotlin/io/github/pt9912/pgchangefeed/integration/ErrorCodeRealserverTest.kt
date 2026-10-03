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
 * Real-server phase for the message code on the typed errors. Enabling a table
 * that does not exist in the source database is answered with HTTP status 404
 * and gRPC status `NOT_FOUND`; both carry the message code the server assigns
 * to that case, and the SDK exposes it as `messageCode` of the typed
 * exception: from the `code` field of the HTTP error body, and from the
 * `google.rpc.ErrorInfo` status detail on the gRPC side. A reader token
 * against the same operation is rejected with HTTP status 403 and gRPC status
 * `PERMISSION_DENIED`, which carry no message code. The unit tests build the
 * wire form themselves; this phase reads it from the real server.
 */
class ErrorCodeRealserverTest {
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

    private fun httpClientFor(token: String) = PgChangeFeedHttpClient(
        java.net.http.HttpClient.newHttpClient(),
        PgChangeFeedClientOptions(URI(PhaseEnvironment.httpAddr), token),
    )

    private fun administrationClientFor(token: String) = PgChangeFeedAdministrationClient(
        PgChangeFeedClientOptions(URI("http://${PhaseEnvironment.grpcAddr}"), token),
    )

    @Test
    fun enablingAMissingTableCarriesTheMessageCodeOverHttpAndGrpc() {
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
        val code = httpException.messageCode
        assertNotNull(code, "HTTP trägt keinen Meldungscode")
        assertTrue(code.startsWith("PCF-E"), "unerwartete Form: $code")
        assertEquals(code, grpcException.messageCode)
        println("RECEIVED code=$code http=${httpException.statusCode} grpc=${grpcException.statusCode}")
        System.out.flush()
    }

    @Test
    fun aRejectedCallCarriesNoMessageCode() {
        val httpException = runCatching { httpClientFor(PhaseEnvironment.readerToken).enableTable(httpEnableRequest()) }
            .exceptionOrNull()
        val grpcException = administrationClientFor(PhaseEnvironment.readerToken).use { client ->
            runCatching { runBlocking { client.enableTable(grpcEnableRequest()) } }.exceptionOrNull()
        }

        assertTrue(httpException is PgChangeFeedForbiddenException, "HTTP: wollen 403, erhalten: $httpException")
        assertTrue(
            grpcException is PgChangeFeedGrpcPermissionDeniedException,
            "gRPC: wollen PERMISSION_DENIED, erhalten: $grpcException",
        )
        assertEquals(403, httpException.statusCode)
        assertEquals(Status.Code.PERMISSION_DENIED, grpcException.statusCode)
        assertNull(httpException.messageCode)
        assertNull(grpcException.messageCode)
        println("REJECTED http=${httpException.statusCode} grpc=${grpcException.statusCode} code=none")
        System.out.flush()
    }
}
