package io.github.pt9912.pgchangefeed.http

import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse
import java.time.Duration

/**
 * A single, already-addressed HTTP request — method, absolute URL, request
 * headers (including `Authorization`) and an optional JSON body. Built by
 * [PgChangeFeedHttpClient], sent by an [HttpTransport] implementation.
 */
internal data class TransportRequest(
    val method: String,
    val url: String,
    val headers: Map<String, String>,
    val body: String?,
)

/** The raw status code and body text of an HTTP response. */
internal data class TransportResponse(
    val statusCode: Int,
    val body: String,
)

/**
 * Transport seam between [PgChangeFeedHttpClient] and the wire.
 *
 * `java.net.http.HttpClient` has no pluggable handler: `send()` always
 * performs real socket I/O. [PgChangeFeedHttpClient] therefore depends on this
 * interface instead of on `HttpClient` directly, so that tests can exercise
 * JSON handling and status-code mapping against canned [TransportResponse]s
 * without any socket.
 *
 * The public constructor `PgChangeFeedHttpClient(HttpClient, PgChangeFeedClientOptions)`
 * takes the caller's `HttpClient` — the caller controls connection pooling,
 * proxies and the client's lifetime, this class never closes it — and wraps it
 * in [JdkHttpTransport]. The `internal` constructor that takes an
 * [HttpTransport] directly is meant for the test source set. `internal` is a
 * compile-time visibility boundary of the Kotlin compiler, not a JVM access
 * restriction: in the compiled class file the constructor and [HttpTransport]
 * are ordinary `public` symbols, so a Java caller or reflection could still use
 * them.
 */
internal fun interface HttpTransport {
    fun send(request: TransportRequest): TransportResponse
}

/**
 * The real [HttpTransport]: adapts a caller-supplied `java.net.http.HttpClient`.
 * Every request carries a timeout of ten seconds.
 */
internal class JdkHttpTransport(private val httpClient: HttpClient) : HttpTransport {
    override fun send(request: TransportRequest): TransportResponse {
        val bodyPublisher = if (request.body != null) {
            HttpRequest.BodyPublishers.ofString(request.body)
        } else {
            HttpRequest.BodyPublishers.noBody()
        }
        val builder = HttpRequest.newBuilder(URI.create(request.url))
            .timeout(requestTimeout)
            .method(request.method, bodyPublisher)
        request.headers.forEach { (name, value) -> builder.header(name, value) }

        val response = httpClient.send(builder.build(), HttpResponse.BodyHandlers.ofString())
        return TransportResponse(response.statusCode(), response.body())
    }

    private companion object {
        val requestTimeout: Duration = Duration.ofSeconds(10)
    }
}
