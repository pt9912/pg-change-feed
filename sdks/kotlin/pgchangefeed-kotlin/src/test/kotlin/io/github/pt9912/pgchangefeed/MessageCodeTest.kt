package io.github.pt9912.pgchangefeed

import cdc.administration.v1.AdministrationOuterClass.DiagnoseRequest
import cdc.administration.v1.AdministrationOuterClass.DiagnoseResponse
import cdc.administration.v1.AdministrationOuterClass.HeartbeatStatus
import cdc.administration.v1.AdministrationOuterClass.ListTablesRequest
import com.google.protobuf.ByteString
import com.google.protobuf.CodedOutputStream
import io.github.pt9912.pgchangefeed.grpc.AdministrationTestClientFactory
import io.github.pt9912.pgchangefeed.grpc.FakeAdministrationTransport
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcInternalException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcInvalidArgumentException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcNotFoundException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcPermissionDeniedException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcUnauthenticatedException
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcUnexpectedStatusException
import io.github.pt9912.pgchangefeed.http.FakeHttpTransport
import io.github.pt9912.pgchangefeed.http.PgChangeFeedBadRequestException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedForbiddenException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedNotFoundException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedServerErrorException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnauthorizedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnexpectedStatusException
import io.github.pt9912.pgchangefeed.sse.FakeSseTransport
import io.grpc.Metadata
import io.grpc.Status
import io.grpc.StatusException
import kotlinx.coroutines.runBlocking
import java.io.ByteArrayOutputStream
import kotlin.reflect.KClass
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertIs
import kotlin.test.assertNull
import io.github.pt9912.pgchangefeed.http.TestClientFactory as HttpTestClientFactory
import io.github.pt9912.pgchangefeed.sse.TestClientFactory as SseTestClientFactory
import com.google.protobuf.Any as ProtoAny

/**
 * The message code on the typed errors: one input table per wire form (the
 * HTTP error body for the HTTP client and the SSE client, the
 * `grpc-status-details-bin` trailer for the administration client), the same
 * rows in every SDK language. The status detail of the gRPC tests is
 * serialized here with the Protobuf runtime's own writer, not with the reader
 * under test.
 */
class MessageCodeTest {
    private data class HttpRow(
        val status: Int,
        val body: String,
        val expected: KClass<out PgChangeFeedException>,
        val code: String?,
        val text: String,
    )

    private val httpRows = listOf(
        HttpRow(400, """{"error":"x","code":"PCF-E8051"}""", PgChangeFeedBadRequestException::class, "PCF-E8051", "x"),
        HttpRow(400, """{"error":"x"}""", PgChangeFeedBadRequestException::class, null, "x"),
        HttpRow(400, """{"error":"x","code":""}""", PgChangeFeedBadRequestException::class, null, "x"),
        HttpRow(400, """{"error":"x","code":null}""", PgChangeFeedBadRequestException::class, null, "x"),
        HttpRow(400, """{"error":"x","code":5}""", PgChangeFeedBadRequestException::class, null, "x"),
        HttpRow(404, """{"error":"x","code":"PCF-E8025"}""", PgChangeFeedNotFoundException::class, "PCF-E8025", "x"),
        HttpRow(500, """{"error":"x","code":"PCF-E7000"}""", PgChangeFeedServerErrorException::class, "PCF-E7000", "x"),
        HttpRow(503, """{"error":"x","code":"PCF-E2001"}""", PgChangeFeedUnexpectedStatusException::class, "PCF-E2001", "x"),
        HttpRow(401, """{"error":"x"}""", PgChangeFeedUnauthorizedException::class, null, "x"),
        HttpRow(403, """{"error":"x"}""", PgChangeFeedForbiddenException::class, null, "x"),
        HttpRow(400, """{"error":"x","code":"NO-PCF","extra":1}""", PgChangeFeedBadRequestException::class, "NO-PCF", "x"),
        HttpRow(502, "Bad Gateway", PgChangeFeedUnexpectedStatusException::class, null, "Bad Gateway"),
        HttpRow(400, """{"error":5,"code":"PCF-E8051"}""", PgChangeFeedBadRequestException::class, "PCF-E8051", """{"error":5,"code":"PCF-E8051"}"""),
    )

    @Test
    fun `http client error carries the message code`() {
        for ((index, row) in httpRows.withIndex()) {
            val (client, _) = HttpTestClientFactory.create(responder = FakeHttpTransport.jsonResponse(row.status, row.body))

            val ex = assertFailsWith<PgChangeFeedException>("row ${index + 1}") { client.listTables("src", "pub") }

            assertEquals(row.expected, ex::class, "row ${index + 1}")
            assertEquals(row.code, ex.messageCode, "row ${index + 1}")
            assertEquals(row.text, ex.message, "row ${index + 1}")
            assertEquals(row.status, ex.statusCode, "row ${index + 1}")
        }
    }

    @Test
    fun `sse client error carries the message code`() {
        for ((index, row) in httpRows.withIndex()) {
            val (client, _) = SseTestClientFactory.create { _ -> FakeSseTransport.lineResponse(row.status, row.body) }

            val ex = assertFailsWith<PgChangeFeedException>("row ${index + 1}") { client.streamChanges().toList() }

            assertEquals(row.expected, ex::class, "row ${index + 1}")
            assertEquals(row.code, ex.messageCode, "row ${index + 1}")
            assertEquals(row.text, ex.message, "row ${index + 1}")
        }
    }

    @Test
    fun `http exception constructed without message code keeps working and has none`() {
        val ex = PgChangeFeedBadRequestException(400, "x")

        assertEquals(400, ex.statusCode)
        assertEquals("x", ex.message)
        assertNull(ex.messageCode)
    }

    private val errorInfoUrl = "type.googleapis.com/google.rpc.ErrorInfo"
    private val retryInfoUrl = "type.googleapis.com/google.rpc.RetryInfo"
    private val serverDomain = "pg-change-feed"

    private fun errorInfo(reason: String, domain: String): ByteString {
        val bytes = ByteArrayOutputStream()
        val output = CodedOutputStream.newInstance(bytes)
        output.writeString(1, reason)
        output.writeString(2, domain)
        output.flush()
        return ByteString.copyFrom(bytes.toByteArray())
    }

    private fun detail(url: String, value: ByteString): ProtoAny =
        ProtoAny.newBuilder().setTypeUrl(url).setValue(value).build()

    private fun statusBytes(vararg details: ProtoAny): ByteArray {
        val bytes = ByteArrayOutputStream()
        val output = CodedOutputStream.newInstance(bytes)
        for (detail in details) {
            output.writeMessage(3, detail)
        }
        output.flush()
        return bytes.toByteArray()
    }

    private fun trailers(statusBytes: ByteArray): Metadata =
        Metadata().apply {
            put(Metadata.Key.of("grpc-status-details-bin", Metadata.BINARY_BYTE_MARSHALLER), statusBytes)
        }

    private fun serverDetail(reason: String): Metadata =
        trailers(statusBytes(detail(errorInfoUrl, errorInfo(reason, serverDomain))))

    private data class GrpcRow(
        val status: Status,
        val trailers: Metadata?,
        val expected: KClass<out PgChangeFeedGrpcException>,
        val code: String?,
    )

    private val grpcRows
        get() = listOf(
            GrpcRow(Status.INVALID_ARGUMENT, serverDetail("PCF-E8051"), PgChangeFeedGrpcInvalidArgumentException::class, "PCF-E8051"),
            GrpcRow(Status.NOT_FOUND, serverDetail("PCF-E8025"), PgChangeFeedGrpcNotFoundException::class, "PCF-E8025"),
            GrpcRow(Status.INTERNAL, serverDetail("PCF-E7000"), PgChangeFeedGrpcInternalException::class, "PCF-E7000"),
            GrpcRow(Status.INVALID_ARGUMENT, null, PgChangeFeedGrpcInvalidArgumentException::class, null),
            GrpcRow(Status.UNAUTHENTICATED, null, PgChangeFeedGrpcUnauthenticatedException::class, null),
            GrpcRow(Status.PERMISSION_DENIED, null, PgChangeFeedGrpcPermissionDeniedException::class, null),
            GrpcRow(
                Status.INVALID_ARGUMENT,
                trailers(statusBytes(detail(errorInfoUrl, errorInfo("PCF-E8051", "example.com")))),
                PgChangeFeedGrpcInvalidArgumentException::class,
                null,
            ),
            GrpcRow(
                Status.INVALID_ARGUMENT,
                trailers(statusBytes(detail(retryInfoUrl, errorInfo("PCF-E8051", serverDomain)))),
                PgChangeFeedGrpcInvalidArgumentException::class,
                null,
            ),
            GrpcRow(
                Status.INVALID_ARGUMENT,
                trailers(statusBytes(detail(errorInfoUrl, errorInfo("", serverDomain)))),
                PgChangeFeedGrpcInvalidArgumentException::class,
                null,
            ),
            GrpcRow(
                Status.INVALID_ARGUMENT,
                trailers(byteArrayOf(0xFF.toByte(), 0xFF.toByte())),
                PgChangeFeedGrpcInvalidArgumentException::class,
                null,
            ),
            GrpcRow(
                Status.INVALID_ARGUMENT,
                trailers(
                    statusBytes(
                        detail(errorInfoUrl, errorInfo("PCF-E0001", "example.com")),
                        detail(errorInfoUrl, errorInfo("PCF-E8051", serverDomain)),
                    ),
                ),
                PgChangeFeedGrpcInvalidArgumentException::class,
                "PCF-E8051",
            ),
            GrpcRow(Status.UNAVAILABLE, serverDetail("PCF-E8051"), PgChangeFeedGrpcUnexpectedStatusException::class, "PCF-E8051"),
        )

    @Test
    fun `administration client error carries the message code`() {
        for ((index, row) in grpcRows.withIndex()) {
            val transport = FakeAdministrationTransport.withStatus(row.status.withDescription("boom"), row.trailers)

            val ex = assertFailsWith<PgChangeFeedGrpcException>("row ${index + 1}") {
                runBlocking {
                    AdministrationTestClientFactory.create(transport)
                        .listTables(ListTablesRequest.newBuilder().setSource("src").setPublication("pub").build())
                }
            }

            assertEquals(row.expected, ex::class, "row ${index + 1}")
            assertEquals(row.code, ex.messageCode, "row ${index + 1}")
            assertEquals(row.status.code, ex.statusCode, "row ${index + 1}")
            assertEquals("boom", ex.message, "row ${index + 1}")
            assertIs<StatusException>(ex.cause, "row ${index + 1}")
        }
    }

    @Test
    fun `grpc exception constructed without message code keeps working and has none`() {
        val cause = StatusException(Status.NOT_FOUND)

        val ex = PgChangeFeedGrpcNotFoundException("gone", cause)

        assertEquals(Status.Code.NOT_FOUND, ex.statusCode)
        assertNull(ex.messageCode)
    }

    @Test
    fun `diagnose error code arrives unchanged`() = runBlocking {
        for (errorCode in listOf("PCF-E4003", "")) {
            val transport = FakeAdministrationTransport.withResponse(
                DiagnoseResponse.newBuilder()
                    .setHeartbeat(
                        HeartbeatStatus.newBuilder()
                            .setKnown(true).setAgeSeconds(1.0)
                            .setErrorClass(if (errorCode.isEmpty()) "" else "schema")
                            .setErrorCode(errorCode)
                            .build(),
                    )
                    .build(),
            )

            val response = AdministrationTestClientFactory.create(transport)
                .diagnose(DiagnoseRequest.newBuilder().setSource("src").build())

            assertEquals(errorCode, response.heartbeat.errorCode)
        }
    }
}
