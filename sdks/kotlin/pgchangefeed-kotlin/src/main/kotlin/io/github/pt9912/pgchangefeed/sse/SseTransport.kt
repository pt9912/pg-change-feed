package io.github.pt9912.pgchangefeed.sse

import io.github.pt9912.pgchangefeed.http.TransportRequest
import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse

/**
 * The raw status code and a lazy line supplier of an opened
 * `GET /changes/stream` response. Unlike
 * [io.github.pt9912.pgchangefeed.http.TransportResponse]'s fully buffered
 * `body: String`, a change stream is unbounded and long-lived — buffering the
 * whole body before returning would never complete on a successful open.
 * [nextLine] is the `() -> String?` shape [SseFrameParser.readFrame] expects,
 * backed by the iterator of `HttpResponse.BodyHandlers.ofLines()`.
 */
internal data class SseTransportResponse(
    val statusCode: Int,
    val nextLine: () -> String?,
)

/**
 * Transport interface between [PgChangeFeedSseClient] and the wire — like
 * [io.github.pt9912.pgchangefeed.http.HttpTransport]: `java.net.http.HttpClient`
 * has no pluggable handler, so [PgChangeFeedSseClient] depends on this
 * interface instead of the JDK client directly, letting tests inject a
 * network-free fake that returns canned [SseTransportResponse]s.
 *
 * `internal` is a compile-time visibility boundary of the Kotlin compiler (the
 * test source set is a friend of the `internal` declarations of `main`), not a
 * JVM access restriction: in the compiled class file this interface, the
 * `internal` constructor of [PgChangeFeedSseClient] and [JdkSseTransport] are
 * ordinary `public` symbols, so a Java caller or reflection could still
 * implement [SseTransport] and use that constructor.
 */
internal fun interface SseTransport {
    fun open(request: TransportRequest): SseTransportResponse
}

/**
 * The real [SseTransport]: adapts a caller-supplied `java.net.http.HttpClient`
 * via `HttpResponse.BodyHandlers.ofLines()`. The built request carries no
 * `.timeout()` — unlike [io.github.pt9912.pgchangefeed.http.JdkHttpTransport],
 * whose fixed 10-second timeout suits a request/response call, a change stream
 * is long-lived by design and has no natural response deadline.
 */
internal class JdkSseTransport(private val httpClient: HttpClient) : SseTransport {
    override fun open(request: TransportRequest): SseTransportResponse {
        val builder = HttpRequest.newBuilder(URI.create(request.url)).GET()
        request.headers.forEach { (name, value) -> builder.header(name, value) }

        val response = httpClient.send(builder.build(), HttpResponse.BodyHandlers.ofLines())
        val lines = response.body().iterator()
        val nextLine: () -> String? = { if (lines.hasNext()) lines.next() else null }
        return SseTransportResponse(response.statusCode(), nextLine)
    }
}
