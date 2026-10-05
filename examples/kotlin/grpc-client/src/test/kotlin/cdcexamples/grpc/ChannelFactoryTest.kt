package cdcexamples.grpc

import cdc.administration.v1.AdministrationGrpcKt
import cdc.administration.v1.AdministrationOuterClass.ListTablesRequest
import cdc.administration.v1.AdministrationOuterClass.ListTablesResponse
import io.grpc.Grpc
import io.grpc.InsecureServerCredentials
import io.grpc.Server
import io.grpc.ServerCredentials
import io.grpc.StatusException
import io.grpc.TlsServerCredentials
import kotlinx.coroutines.runBlocking
import java.io.File
import java.nio.file.Files
import java.security.KeyStore
import java.util.concurrent.TimeUnit
import javax.net.ssl.KeyManagerFactory
import kotlin.test.AfterTest
import kotlin.test.Test
import kotlin.test.assertContains
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

/**
 * Prüft [ChannelFactory]: die Wahl „Anker gesetzt, dann TLS“ ist an den
 * Optionswert gebunden. Der Kanal läuft gegen einen gRPC-Server auf Loopback
 * (mit und ohne TLS); das Zertifikat erzeugt der Test zur Laufzeit mit dem
 * `keytool` des JDK im Temp-Verzeichnis — kein Zertifikat und kein Schlüssel
 * liegen im Repo.
 */
class ChannelFactoryTest {
    private val dir: File = Files.createTempDirectory("grpc-client-tls").toFile()
    private val servers = mutableListOf<Server>()

    @AfterTest
    fun cleanup() {
        servers.forEach { it.shutdownNow().awaitTermination(5, TimeUnit.SECONDS) }
        dir.deleteRecursively()
    }

    private fun keytool(vararg args: String) {
        val bin = File(System.getProperty("java.home"), "bin/keytool").path
        val process = ProcessBuilder(listOf(bin) + args).redirectErrorStream(true).start()
        val output = process.inputStream.readBytes().decodeToString()
        assertEquals(0, process.waitFor(), "keytool: $output")
    }

    /** Erzeugt ein selbstsigniertes Zertifikat mit dem SAN [san] und gibt Keystore und Anker-PEM zurück. */
    private fun newCertificate(name: String, san: String = "ip:127.0.0.1"): Pair<File, File> {
        val store = File(dir, "$name.p12")
        val pem = File(dir, "$name.pem")
        keytool(
            "-genkeypair", "-alias", name, "-keyalg", "EC", "-groupname", "secp256r1",
            "-dname", "CN=grpc-client-test", "-ext", "san=$san", "-validity", "1",
            "-storetype", "PKCS12", "-keystore", store.path, "-storepass", "changeit",
        )
        keytool("-exportcert", "-rfc", "-alias", name, "-keystore", store.path, "-storepass", "changeit", "-file", pem.path)
        return store to pem
    }

    private fun serverCredentials(store: File?): ServerCredentials {
        if (store == null) {
            return InsecureServerCredentials.create()
        }
        val keyStore = KeyStore.getInstance("PKCS12")
        store.inputStream().use { keyStore.load(it, "changeit".toCharArray()) }
        val factory = KeyManagerFactory.getInstance(KeyManagerFactory.getDefaultAlgorithm())
        factory.init(keyStore, "changeit".toCharArray())
        return TlsServerCredentials.newBuilder().keyManager(*factory.keyManagers).build()
    }

    private fun startServer(store: File?): String {
        val service = object : AdministrationGrpcKt.AdministrationCoroutineImplBase() {
            override suspend fun listTables(request: ListTablesRequest): ListTablesResponse =
                ListTablesResponse.getDefaultInstance()
        }
        val server = Grpc.newServerBuilderForPort(0, serverCredentials(store)).addService(service).build().start()
        servers.add(server)
        return "127.0.0.1:${server.port}"
    }

    private fun config(addr: String, caFile: String = "") = Config(
        addr = addr, token = "", adminToken = "", verb = "list-tables", schema = "", table = "",
        consumerId = "", name = "", offset = 0, tableId = "", schemaVersionId = "", version = 1,
        source = "", publication = "", from = 0, to = 0, limit = 0, minAgeNanos = 0, caFile = caFile,
    )

    /** Ruft ListTables über den von [ChannelFactory] gebauten Kanal auf; null heißt Erfolg. */
    private fun call(cfg: Config): StatusException? {
        val channel = ChannelFactory.create(cfg)
        try {
            val stub = AdministrationGrpcKt.AdministrationCoroutineStub(channel).withDeadlineAfter(5, TimeUnit.SECONDS)
            runBlocking { stub.listTables(ListTablesRequest.getDefaultInstance()) }
            return null
        } catch (ex: StatusException) {
            return ex
        } finally {
            channel.shutdownNow()
        }
    }

    @Test
    fun anchorReachesTlsServer() {
        val (store, pem) = newCertificate("a")
        val addr = startServer(store)
        assertEquals(null, call(config(addr, pem.path)))
    }

    @Test
    fun withoutAnchorPlaintextFailsAgainstTlsServer() {
        val (store, _) = newCertificate("a")
        val addr = startServer(store)
        assertTrue(call(config(addr)) != null, "Klartext-Aufruf gegen den TLS-Server gelang")
    }

    @Test
    fun withoutAnchorReachesPlaintextServer() {
        val addr = startServer(null)
        assertEquals(null, call(config(addr)))
    }

    @Test
    fun anchorFailsAgainstPlaintextServer() {
        val (_, pem) = newCertificate("a")
        val addr = startServer(null)
        assertTrue(call(config(addr, pem.path)) != null, "TLS-Aufruf gegen den Klartext-Server gelang")
    }

    @Test
    fun foreignAnchorFailsChainCheck() {
        val (store, _) = newCertificate("a")
        val (_, foreign) = newCertificate("b")
        val addr = startServer(store)
        assertTrue(call(config(addr, foreign.path)) != null, "Aufruf mit fremdem Anker gelang")
    }

    /** Der Anker ist das Serverzertifikat (Kette stimmt), sein SAN nennt aber nicht die Adresse der Verbindung. */
    @Test
    fun nameOutsideSanFailsAgainstTrustedCertificate() {
        val (store, pem) = newCertificate("n", san = "dns:other.example")
        val addr = startServer(store)
        assertTrue(call(config(addr, pem.path)) != null, "Aufruf an einen Namen außerhalb des SAN gelang")
    }

    @Test
    fun derAnchorIsAcceptedAndReachesTlsServer() {
        val (store, _) = newCertificate("d")
        val der = File(dir, "d.der")
        keytool("-exportcert", "-alias", "d", "-keystore", store.path, "-storepass", "changeit", "-file", der.path)
        val addr = startServer(store)
        assertEquals(null, call(config(addr, der.path)))
    }

    @Test
    fun missingAndPemlessAnchorFailBeforeConnecting() {
        val missing = assertFailsWith<IllegalArgumentException> {
            ChannelFactory.create(config("127.0.0.1:1", File(dir, "fehlt.pem").path))
        }
        assertContains(missing.message ?: "", "nicht lesbar")
        val plain = File(dir, "kein-pem.txt").apply { writeText("kein Zertifikat") }
        val pemless = assertFailsWith<IllegalArgumentException> {
            ChannelFactory.create(config("127.0.0.1:1", plain.path))
        }
        assertContains(pemless.message ?: "", "kein PEM-Zertifikat")
    }
}
