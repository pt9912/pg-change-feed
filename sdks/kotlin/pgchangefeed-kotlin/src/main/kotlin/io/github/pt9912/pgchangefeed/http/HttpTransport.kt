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
 * Transport seam between [PgChangeFeedHttpClient] and the actual wire.
 *
 * `java.net.http.HttpClient` (the JDK client `examples/kotlin/http-client`
 * uses, SPEC-023) has **no** pluggable-handler concept — unlike C#'s
 * `HttpMessageHandler` (injected into `System.Net.Http.HttpClient`, faked
 * in the C# SDK's tests via `FakeHttpMessageHandler`) or Python's `httpx`
 * (faked via `httpx.MockTransport`), `HttpClient.send()` always performs
 * real socket I/O against the `URI` on the request. Injecting the real
 * `java.net.http.HttpClient` directly into [PgChangeFeedHttpClient] would
 * therefore force every response-shape/error-mapping test to open a real
 * socket (either a genuine network call or an in-process loopback HTTP
 * server) just to exercise JSON (de)serialization and status-code mapping.
 *
 * This interface is the fix: [PgChangeFeedHttpClient] depends on
 * [HttpTransport], not on `java.net.http.HttpClient` directly. The public
 * constructor (`PgChangeFeedHttpClient(HttpClient, PgChangeFeedClientOptions)`)
 * keeps the same caller-supplies-the-HttpClient parity as the C#/Python
 * SDKs — the caller controls connection pooling, proxies and the
 * `HttpClient`'s lifetime, this type never closes it — by wrapping it in
 * [JdkHttpTransport] internally. The `internal` secondary constructor that
 * takes an [HttpTransport] directly is invisible outside this module, so it
 * adds no public API surface; it exists purely so the test source set
 * (which the Kotlin Gradle plugin's default `main`/`test` association makes
 * a friend of `internal` declarations) can inject a fake [HttpTransport]
 * that returns canned [TransportResponse]s without any socket at all — a
 * genuinely network-free test, stronger than a loopback server.
 */
internal fun interface HttpTransport {
    fun send(request: TransportRequest): TransportResponse
}

/**
 * The real [HttpTransport]: adapts a caller-supplied `java.net.http.HttpClient`
 * (the same JDK client `examples/kotlin/http-client`'s `TablesClient` uses,
 * read as draht-Kenntnis, not imported — `ADR-0109` Festlegung 3).
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
