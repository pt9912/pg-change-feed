package io.github.pt9912.pgchangefeed.tls

import cdc.administration.v1.AdministrationGrpcKt
import cdc.administration.v1.AdministrationOuterClass.ListTablesRequest
import cdc.administration.v1.AdministrationOuterClass.ListTablesResponse
import cdc.stream.v1.ChangeStreamGrpcKt
import cdc.stream.v1.Changestream.Change
import cdc.stream.v1.Changestream.StreamChangesRequest
import com.sun.net.httpserver.HttpExchange
import com.sun.net.httpserver.HttpServer
import com.sun.net.httpserver.HttpsConfigurator
import com.sun.net.httpserver.HttpsServer
import io.grpc.Grpc
import io.grpc.InsecureServerCredentials
import io.grpc.Server
import io.grpc.ServerCredentials
import io.grpc.TlsServerCredentials
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flowOf
import java.io.File
import java.net.InetAddress
import java.net.InetSocketAddress
import java.nio.file.Files
import java.security.KeyStore
import java.util.concurrent.TimeUnit
import javax.net.ssl.KeyManagerFactory
import javax.net.ssl.SSLContext
import kotlin.test.assertEquals

/** A certificate created at run time: the key store the server uses and the PEM file that serves as trust anchor. */
internal class TestCertificate(val store: File, val pem: File)

/**
 * Servers on loopback ports for the TLS tests, with and without TLS, and the
 * certificates they present. Certificates are created at run time with the
 * `keytool` of the JDK in a temporary directory that [close] removes; nothing
 * is committed.
 */
internal class TlsFixture : AutoCloseable {
    private val directory: File = Files.createTempDirectory("pgcf-tls").toFile()
    private val httpServers = mutableListOf<HttpServer>()
    private val grpcServers = mutableListOf<Server>()
    private var counter = 0

    override fun close() {
        httpServers.forEach { it.stop(0) }
        grpcServers.forEach { it.shutdownNow().awaitTermination(5, TimeUnit.SECONDS) }
        directory.deleteRecursively()
    }

    /** Creates a self-signed certificate that lists [san] (for example `dns:localhost`), valid from [startDate]. */
    fun newCertificate(san: String = "dns:localhost", startDate: String? = null): TestCertificate {
        val name = "cert${counter++}"
        val store = File(directory, "$name.p12")
        val pem = File(directory, "$name.pem")
        val genArgs = mutableListOf(
            "-genkeypair", "-alias", name, "-keyalg", "EC", "-groupname", "secp256r1",
            "-dname", "CN=pgcf-test-$name", "-ext", "san=$san", "-validity", "1",
            "-storetype", "PKCS12", "-keystore", store.path, "-storepass", PASSWORD,
        )
        if (startDate != null) {
            genArgs += listOf("-startdate", startDate)
        }
        keytool(genArgs)
        keytool(listOf("-exportcert", "-rfc", "-alias", name, "-keystore", store.path, "-storepass", PASSWORD, "-file", pem.path))
        return TestCertificate(store, pem)
    }

    /**
     * Creates a certificate authority and a server certificate for [san] that it issued. The key store is the
     * server's (key, server certificate and authority); the PEM file is the authority certificate, which is the
     * trust anchor — the server certificate is not its own anchor.
     */
    fun newIssuedCertificate(san: String = "dns:localhost"): TestCertificate {
        val name = "issued${counter++}"
        val authorityStore = File(directory, "$name-ca.p12")
        val authorityPem = File(directory, "$name-ca.pem")
        val serverStore = File(directory, "$name.p12")
        val request = File(directory, "$name.csr")
        val signed = File(directory, "$name-signed.pem")
        keytool(
            listOf(
                "-genkeypair", "-alias", "ca", "-keyalg", "EC", "-groupname", "secp256r1", "-dname", "CN=pgcf-test-ca-$name",
                "-ext", "bc:c", "-ext", "ku:c=keyCertSign", "-validity", "1", "-storetype", "PKCS12",
                "-keystore", authorityStore.path, "-storepass", PASSWORD,
            ),
        )
        keytool(listOf("-exportcert", "-rfc", "-alias", "ca", "-keystore", authorityStore.path, "-storepass", PASSWORD, "-file", authorityPem.path))
        keytool(
            listOf(
                "-genkeypair", "-alias", "server", "-keyalg", "EC", "-groupname", "secp256r1", "-dname", "CN=pgcf-test-server-$name",
                "-validity", "1", "-storetype", "PKCS12", "-keystore", serverStore.path, "-storepass", PASSWORD,
            ),
        )
        keytool(listOf("-certreq", "-alias", "server", "-keystore", serverStore.path, "-storepass", PASSWORD, "-file", request.path))
        keytool(
            listOf(
                "-gencert", "-alias", "ca", "-keystore", authorityStore.path, "-storepass", PASSWORD,
                "-infile", request.path, "-outfile", signed.path, "-rfc", "-ext", "san=$san", "-validity", "1",
            ),
        )
        keytool(listOf("-importcert", "-alias", "ca", "-file", authorityPem.path, "-keystore", serverStore.path, "-storepass", PASSWORD, "-noprompt"))
        keytool(listOf("-importcert", "-alias", "server", "-file", signed.path, "-keystore", serverStore.path, "-storepass", PASSWORD))
        return TestCertificate(serverStore, authorityPem)
    }

    /** Writes the PEM files of the certificates, one after the other, into one file. */
    fun concatenate(vararg certificates: TestCertificate): File {
        val file = File(directory, "bundle${counter++}.pem")
        file.writeText(certificates.joinToString("") { it.pem.readText() })
        return file
    }

    fun textFile(text: String): File = File(directory, "text${counter++}.txt").also { it.writeText(text) }

    fun missingFile(): File = File(directory, "missing.pem")

    private fun keytool(args: List<String>) {
        val bin = File(System.getProperty("java.home"), "bin/keytool").path
        val process = ProcessBuilder(listOf(bin) + args).redirectErrorStream(true).start()
        val output = process.inputStream.readBytes().decodeToString()
        assertEquals(0, process.waitFor(), "keytool: $output")
    }

    private fun keyManagers(store: File) = KeyStore.getInstance("PKCS12").let { keyStore ->
        store.inputStream().use { keyStore.load(it, PASSWORD.toCharArray()) }
        KeyManagerFactory.getInstance(KeyManagerFactory.getDefaultAlgorithm()).also { it.init(keyStore, PASSWORD.toCharArray()) }
            .keyManagers
    }

    /** Starts the HTTP and SSE server; with [certificate] over TLS, otherwise in plaintext. Returns the port. */
    fun startHttp(certificate: TestCertificate?): Int {
        val loopback = InetSocketAddress(InetAddress.getLoopbackAddress(), 0)
        val server: HttpServer = if (certificate == null) {
            HttpServer.create(loopback, 0)
        } else {
            HttpsServer.create(loopback, 0).also {
                val context = SSLContext.getInstance("TLS")
                context.init(keyManagers(certificate.store), null, null)
                it.httpsConfigurator = HttpsConfigurator(context)
            }
        }
        server.createContext("/tables") { exchange ->
            reply(exchange, "application/json", """{"tables":[],"retained":[]}""")
        }
        server.createContext("/changes/stream") { exchange ->
            reply(
                exchange,
                "text/event-stream",
                "event: change\ndata: " +
                    """{"change_id":"$SENTINEL","transaction_id":"tx-1","source_table_id":"t-1","sequence":1,""" +
                    """"operation":"INSERT","old_image":null,"new_image":{"id":1},"schema_version":"v1",""" +
                    """"schema":"public","table":"orders"}""" + "\n\n",
            )
        }
        server.start()
        httpServers.add(server)
        return server.address.port
    }

    private fun reply(exchange: HttpExchange, contentType: String, body: String) {
        val bytes = body.toByteArray()
        exchange.responseHeaders.add("Content-Type", contentType)
        exchange.sendResponseHeaders(200, bytes.size.toLong())
        exchange.responseBody.use { it.write(bytes) }
    }

    /** Starts the gRPC server (stream and administration); with [certificate] over TLS, otherwise in plaintext. Returns the port. */
    fun startGrpc(certificate: TestCertificate?): Int {
        val credentials: ServerCredentials = if (certificate == null) {
            InsecureServerCredentials.create()
        } else {
            TlsServerCredentials.newBuilder().keyManager(*keyManagers(certificate.store)).build()
        }
        val stream = object : ChangeStreamGrpcKt.ChangeStreamCoroutineImplBase() {
            override fun streamChanges(request: StreamChangesRequest): Flow<Change> =
                flowOf(Change.newBuilder().setChangeId(SENTINEL).setOperation("INSERT").setTable("orders").build())
        }
        val administration = object : AdministrationGrpcKt.AdministrationCoroutineImplBase() {
            override suspend fun listTables(request: ListTablesRequest): ListTablesResponse =
                ListTablesResponse.getDefaultInstance()
        }
        val server = Grpc.newServerBuilderForPort(0, credentials)
            .addService(stream).addService(administration).build().start()
        grpcServers.add(server)
        return server.port
    }

    companion object {
        const val SENTINEL = "tls-test-change"
        private const val PASSWORD = "changeit"
    }
}
