package io.github.pt9912.pgchangefeed.sse

import io.github.pt9912.pgchangefeed.http.TransportRequest

/**
 * A network-free [SseTransport]: every test in this module stays free of
 * real sockets (no real server, no loopback listener) — see
 * [SseTransport]'s KDoc for why `java.net.http.HttpClient` itself cannot be
 * faked this cheaply. [lastRequest] lets a test assert on the exact
 * method/URL/headers the client under test built.
 */
internal class FakeSseTransport(private val responder: (TransportRequest) -> SseTransportResponse) : SseTransport {
    var lastRequest: TransportRequest? = null
        private set

    override fun open(request: TransportRequest): SseTransportResponse {
        lastRequest = request
        return responder(request)
    }

    companion object {
        /**
         * Builds an [SseTransportResponse] whose line supplier replays
         * [body] split on `\n` — the same shape
         * `HttpResponse.BodyHandlers.ofLines()` would hand [JdkSseTransport],
         * without any socket.
         */
        fun lineResponse(statusCode: Int, body: String): SseTransportResponse {
            val lines = body.split("\n").iterator()
            return SseTransportResponse(statusCode) { if (lines.hasNext()) lines.next() else null }
        }
    }
}
