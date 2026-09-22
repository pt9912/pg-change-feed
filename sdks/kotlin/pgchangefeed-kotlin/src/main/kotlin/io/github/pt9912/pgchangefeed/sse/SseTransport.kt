package io.github.pt9912.pgchangefeed.sse

import io.github.pt9912.pgchangefeed.http.TransportRequest
import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse

/**
 * The raw status code and a lazy line supplier of an opened `GET
 * /changes/stream` response (`SPEC-021`). Unlike
 * [io.github.pt9912.pgchangefeed.http.TransportResponse]'s fully buffered
 * `body: String`, a change stream is unbounded and long-lived — buffering
 * the whole body before returning would never complete on a successful
 * open. [nextLine] carries the same `() -> String?` shape
 * `examples/kotlin/sse-client/Main.kt` already uses over
 * `HttpResponse.BodyHandlers.ofLines()`'s iterator (read as Draht-Kenntnis,
 * not imported — `ADR-0109` Festlegung 3), and is exactly what
 * [SseFrameParser.readFrame] expects.
 */
internal data class SseTransportResponse(
    val statusCode: Int,
    val nextLine: () -> String?,
)

/**
 * Transport seam between [PgChangeFeedSseClient] and the actual wire — same
 * reasoning as [io.github.pt9912.pgchangefeed.http.HttpTransport]'s KDoc:
 * `java.net.http.HttpClient` has no pluggable-handler concept, so
 * [PgChangeFeedSseClient] depends on this seam instead of the JDK client
 * directly, letting tests inject a network-free fake that returns canned
 * [SseTransportResponse]s without any socket at all — the same seam
 * [io.github.pt9912.pgchangefeed.http.HttpTransport] provides for the HTTP
 * client surface, and [io.github.pt9912.pgchangefeed.grpc.GrpcStreamTransport]
 * for the gRPC surface.
 *
 * `internal` here — as with those two sibling seams — is a **compile-time**
 * Kotlin-compiler visibility boundary against other Kotlin modules'
 * metadata (the Kotlin Gradle plugin's default main/test sourceSet
 * association makes the test sourceSet a friend of `internal` declarations
 * in `main`). It is **not** a JVM bytecode access restriction: in the
 * compiled class file this interface, [PgChangeFeedSseClient]'s
 * `internal`-constructor test seam, and [JdkSseTransport] below are
 * ordinary `public` symbols (interfaces are not covered by Kotlin's
 * `internal` name-mangling, and constructors are always named `<init>` in
 * bytecode and are likewise not mangled) — a Java caller, or reflection
 * from any language, can still implement [SseTransport] and invoke that
 * constructor directly.
 */
internal fun interface SseTransport {
    fun open(request: TransportRequest): SseTransportResponse
}

/**
 * The real [SseTransport]: adapts a caller-supplied `java.net.http.HttpClient`
 * (the same JDK client [io.github.pt9912.pgchangefeed.http.JdkHttpTransport]
 * wraps for the request/response HTTP surface) via
 * `HttpResponse.BodyHandlers.ofLines()`. Deliberately carries no
 * `.timeout()` on the built request — unlike
 * [io.github.pt9912.pgchangefeed.http.JdkHttpTransport], which sets a fixed
 * 10-second timeout appropriate for a request/response call, a change
 * stream is long-lived by design and has no natural response deadline
 * (pattern: `examples/kotlin/sse-client/Main.kt`'s own comment on the same
 * choice, read as Draht-Kenntnis, not imported — `ADR-0109` Festlegung 3).
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
