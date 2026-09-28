package cdcexamples.http

import com.sun.net.httpserver.HttpServer
import java.net.InetSocketAddress
import java.nio.charset.StandardCharsets

/**
 * FakeServer startet einen echten, aber netzlosen In-Prozess-HTTP-Server
 * auf dem Loopback-Interface für die Aufruf-Funktionen dieses Beispiels —
 * dieselbe Grenze wie `httptest.NewServer` in `examples/http-client` (Go).
 * Er zeichnet die letzte Anfrage (Methode, Pfad, Query, Body) auf und
 * antwortet mit dem vom Aufrufer vorgegebenen Statuscode und Body.
 */
class FakeServer private constructor(private val server: HttpServer) : AutoCloseable {
    val addr: String get() = "127.0.0.1:${server.address.port}"

    var lastMethod: String = ""
        private set
    var lastPath: String = ""
        private set
    var lastQuery: String = ""
        private set
    var lastBody: String? = null
        private set

    private var responseStatus = 200
    private var responseBody = ""

    fun respondWith(status: Int, body: String) {
        responseStatus = status
        responseBody = body
    }

    companion object {
        fun start(): FakeServer {
            val server = HttpServer.create(InetSocketAddress("127.0.0.1", 0), 0)
            val fake = FakeServer(server)
            server.createContext("/") { exchange ->
                fake.lastMethod = exchange.requestMethod
                fake.lastPath = exchange.requestURI.path
                fake.lastQuery = exchange.requestURI.query ?: ""
                val rawBody = exchange.requestBody.readBytes()
                fake.lastBody = if (rawBody.isEmpty()) null else rawBody.toString(StandardCharsets.UTF_8)

                val bytes = fake.responseBody.toByteArray(StandardCharsets.UTF_8)
                exchange.responseHeaders.set("Content-Type", "application/json")
                exchange.sendResponseHeaders(fake.responseStatus, bytes.size.toLong())
                exchange.responseBody.use { it.write(bytes) }
            }
            server.start()
            return fake
        }
    }

    override fun close() {
        server.stop(0)
    }
}
