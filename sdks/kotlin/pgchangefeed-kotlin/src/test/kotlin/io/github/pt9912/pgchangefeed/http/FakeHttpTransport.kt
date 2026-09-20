package io.github.pt9912.pgchangefeed.http

/**
 * A network-free [HttpTransport]: every test in this module stays free of
 * real sockets (no real server, no loopback listener) — see
 * [HttpTransport]'s KDoc for why `java.net.http.HttpClient` itself cannot
 * be faked this cheaply. [lastRequest] lets a test assert on the exact
 * method/URL/headers/body the client under test built.
 */
internal class FakeHttpTransport(private val responder: (TransportRequest) -> TransportResponse) : HttpTransport {
    var lastRequest: TransportRequest? = null
        private set

    override fun send(request: TransportRequest): TransportResponse {
        lastRequest = request
        return responder(request)
    }

    companion object {
        fun jsonResponse(statusCode: Int, json: String): (TransportRequest) -> TransportResponse =
            { TransportResponse(statusCode, json) }
    }
}
