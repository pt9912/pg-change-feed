package io.github.pt9912.pgchangefeed.tls

import cdc.administration.v1.AdministrationOuterClass.ListTablesRequest
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedAdministrationClient
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcUnexpectedStatusException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient
import io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient
import io.grpc.Status
import io.grpc.StatusException
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import java.io.File
import java.io.IOException
import java.net.URI
import kotlin.test.AfterTest
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNull

/**
 * The four surfaces that connect over TLS (HTTP, SSE, gRPC stream, gRPC
 * administration) against servers with a certificate created at run time:
 * with the trust anchor the connection works, and a foreign anchor, a server
 * name that is not in the certificate, an expired certificate and no anchor at
 * all (the certificate is unknown to the Java runtime) each end with the
 * connection error of the surface, without a plaintext retry.
 */
class TlsClientTest {
    private val fixture = TlsFixture()

    @AfterTest
    fun cleanup() = fixture.close()

    private enum class Surface { HTTP, SSE, GRPC_STREAM, GRPC_ADMIN }

    private fun options(scheme: String, host: String, port: Int, anchor: File?) =
        PgChangeFeedClientOptions(URI("$scheme://$host:$port"), "tls-test-token", anchor?.toPath())

    /** Starts the server of the surface with [certificate] (or in plaintext) and returns its port. */
    private fun start(surface: Surface, certificate: TestCertificate?): Int = when (surface) {
        Surface.HTTP, Surface.SSE -> fixture.startHttp(certificate)
        Surface.GRPC_STREAM, Surface.GRPC_ADMIN -> fixture.startGrpc(certificate)
    }

    private fun use(surface: Surface, options: PgChangeFeedClientOptions) {
        when (surface) {
            Surface.HTTP -> {
                val tables = PgChangeFeedHttpClient(options).listTables("src", "pub")
                assertEquals(0, tables.tables.size)
            }
            Surface.SSE -> assertEquals(TlsFixture.SENTINEL, PgChangeFeedSseClient(options).streamChanges().first().changeId)
            Surface.GRPC_STREAM -> PgChangeFeedGrpcClient(options).use { client ->
                assertEquals(TlsFixture.SENTINEL, runBlocking { client.streamChanges().first() }.changeId)
            }
            Surface.GRPC_ADMIN -> PgChangeFeedAdministrationClient(options).use { client ->
                runBlocking { client.listTables(ListTablesRequest.getDefaultInstance()) }
            }
        }
    }

    /**
     * The surface fails with its connection error: the `IOException` of the
     * HTTP client, the `StatusException` with status `UNAVAILABLE` of the
     * stream, the `PgChangeFeedGrpcUnexpectedStatusException` without a
     * message code of the administration client.
     */
    private fun expectConnectionFailure(surface: Surface, options: PgChangeFeedClientOptions) {
        when (surface) {
            Surface.HTTP, Surface.SSE -> assertFailsWith<IOException>("$surface") { use(surface, options) }
            Surface.GRPC_STREAM -> {
                val error = assertFailsWith<StatusException>("$surface") { use(surface, options) }
                assertEquals(Status.Code.UNAVAILABLE, error.status.code, "$surface")
            }
            Surface.GRPC_ADMIN -> {
                val error = assertFailsWith<PgChangeFeedGrpcUnexpectedStatusException>("$surface") { use(surface, options) }
                assertEquals(Status.Code.UNAVAILABLE, error.statusCode, "$surface")
                assertNull(error.messageCode, "$surface")
            }
        }
    }

    @Test
    fun `with a trust anchor every surface connects over TLS`() {
        for (surface in Surface.entries) {
            val certificate = fixture.newCertificate()
            val port = start(surface, certificate)
            use(surface, options("https", "localhost", port, certificate.pem))
        }
    }

    @Test
    fun `the issuer of the server certificate works as trust anchor`() {
        for (surface in Surface.entries) {
            val issued = fixture.newIssuedCertificate()
            val port = start(surface, issued)
            use(surface, options("https", "localhost", port, issued.pem))
        }
    }

    @Test
    fun `the issuer of another server certificate fails the connection`() {
        for (surface in Surface.entries) {
            val port = start(surface, fixture.newIssuedCertificate())
            expectConnectionFailure(surface, options("https", "localhost", port, fixture.newIssuedCertificate().pem))
        }
    }

    @Test
    fun `an anchor file with several certificates trusts the one that matches`() {
        for (surface in Surface.entries) {
            val certificate = fixture.newCertificate()
            val foreign = fixture.newCertificate()
            val port = start(surface, certificate)
            use(surface, options("https", "localhost", port, fixture.concatenate(foreign, certificate)))
        }
    }

    @Test
    fun `a foreign anchor fails the connection`() {
        for (surface in Surface.entries) {
            val port = start(surface, fixture.newCertificate())
            expectConnectionFailure(surface, options("https", "localhost", port, fixture.newCertificate().pem))
        }
    }

    @Test
    fun `a server name that is not in the certificate fails the connection`() {
        for (surface in Surface.entries) {
            val certificate = fixture.newCertificate(san = "dns:other.example")
            val port = start(surface, certificate)
            // The anchor is the server certificate itself, so the chain is right; the address names `localhost`.
            expectConnectionFailure(surface, options("https", "localhost", port, certificate.pem))
        }
    }

    @Test
    fun `an expired certificate fails the connection`() {
        for (surface in Surface.entries) {
            val expired = fixture.newCertificate(startDate = "-3d")
            val port = start(surface, expired)
            expectConnectionFailure(surface, options("https", "localhost", port, expired.pem))
        }
    }

    @Test
    fun `without an anchor the Java runtime decides and refuses an unknown certificate`() {
        for (surface in Surface.entries) {
            val port = start(surface, fixture.newCertificate())
            expectConnectionFailure(surface, options("https", "localhost", port, null))
        }
    }

    @Test
    fun `without an anchor a plaintext address works unchanged`() {
        for (surface in Surface.entries) {
            val port = start(surface, null)
            use(surface, options("http", "localhost", port, null))
        }
    }

    @Test
    fun `a trust anchor with a plaintext address is refused when the client is created`() {
        val certificate = fixture.newCertificate()
        for (surface in Surface.entries) {
            val options = options("http", "localhost", 9, certificate.pem)
            assertFailsWith<IllegalArgumentException>("$surface") { use(surface, options) }
        }
    }

    @Test
    fun `the anchor option is checked when the options are created`() {
        val address = URI("https://feed.example.invalid:8443")
        val withAnchor = PgChangeFeedClientOptions(address, "token", fixture.newCertificate().pem.toPath())
        assertEquals(true, withAnchor.trustAnchorFile != null)
        assertNull(PgChangeFeedClientOptions(address, "token").trustAnchorFile)

        val missing = assertFailsWith<IllegalArgumentException> {
            PgChangeFeedClientOptions(address, "token", fixture.missingFile().toPath())
        }
        assertEquals(true, missing.message!!.contains("missing.pem"))
        assertFailsWith<IllegalArgumentException> {
            PgChangeFeedClientOptions(address, "token", fixture.textFile("this is not a certificate\n").toPath())
        }
        assertFailsWith<IllegalArgumentException> {
            PgChangeFeedClientOptions(address, "token", fixture.textFile("").toPath())
        }
    }
}
