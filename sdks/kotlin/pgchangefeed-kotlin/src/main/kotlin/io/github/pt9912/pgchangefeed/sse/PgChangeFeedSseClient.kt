package io.github.pt9912.pgchangefeed.sse

import com.google.gson.Gson
import com.google.gson.JsonSyntaxException
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.http.PgChangeFeedBadRequestException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedForbiddenException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedMalformedResponseException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedNotFoundException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedServerErrorException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnauthorizedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnexpectedStatusException
import io.github.pt9912.pgchangefeed.http.TransportRequest
import io.github.pt9912.pgchangefeed.http.model.ErrorResponse
import io.github.pt9912.pgchangefeed.sse.model.Change
import java.net.http.HttpClient

/**
 * Public entry point for the PG Change Feed live-change stream over
 * Server-Sent-Events (`SPEC-021`, `LH-FA-SST-008`): [streamChanges] opens
 * `GET /changes/stream`, assembles each SSE frame via [SseFrameParser] and
 * yields the ten `SPEC-021` message fields as [Change] — reusing the same
 * typed [PgChangeFeedException] hierarchy
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient] throws (the
 * uniform SPEC-018-style `{"error": "<text>"}` error body), not a second
 * hierarchy.
 *
 * [streamChanges] returns a [Sequence]`<`[Change]`>`, not a
 * `kotlinx.coroutines.flow.Flow` the way
 * [io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient.streamChanges]
 * does: the generated gRPC coroutine stub gives that surface a genuinely
 * *suspending*, non-blocking call to build on. `java.net.http.HttpClient.send()`
 * (the same JDK client
 * [io.github.pt9912.pgchangefeed.http.JdkHttpTransport] wraps) has no such
 * non-blocking primitive — it is an ordinary blocking call, exactly like
 * every method on [io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient].
 * Wrapping that blocking read in a `flow { … }` builder would need either
 * an explicit `flowOn(Dispatchers.IO)` hop or misleadingly imply suspension
 * points this transport does not offer. [Sequence] keeps the same
 * synchronous execution model as the rest of this package's HTTP surface
 * while staying **cold**: nothing runs — no request is sent, no frame is
 * parsed — until the returned [Sequence] is actually iterated (`ADR-0109`
 * Festlegung 1 unifies only the shared connection denominator across
 * surfaces, not each surface's execution model).
 *
 * The bearer token is sent in the `Authorization` header as
 * `Bearer <token>` (`SPEC-021`), same as
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient]. There is no
 * global or static state; a process can hold several independently
 * configured instances at once.
 *
 * **Boundary (`SPEC-021`, `LH-FA-SST-008`):** no delivery guarantee and no
 * in-stream replay — a disconnected or slow-reading consumer misses the
 * affected messages permanently. Missed changes remain recoverable through
 * the existing read path
 * (`io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient.readChanges`).
 *
 * Draht-Kenntnis-Vorbild (gelesen, nicht importiert — `ADR-0109`
 * Festlegung 3): `examples/kotlin/sse-client/src/main/kotlin/cdcexamples/sse/SseStream.kt`,
 * `Main.kt`.
 */
class PgChangeFeedSseClient internal constructor(
    private val transport: SseTransport,
    private val options: PgChangeFeedClientOptions,
) {
    /**
     * Public constructor — wraps the caller-supplied `java.net.http.HttpClient`
     * in the real [JdkSseTransport]. The [HttpClient] is injected, not
     * owned — the caller controls its lifetime, connection pooling and any
     * redirect/proxy configuration; this type never closes it. See
     * [SseTransport]'s KDoc for why the client depends on a transport seam
     * rather than on [HttpClient] directly.
     */
    constructor(httpClient: HttpClient, options: PgChangeFeedClientOptions) :
        this(JdkSseTransport(httpClient), options)

    /**
     * `GET /changes/stream` (`SPEC-021`) — opens the stream and yields
     * every [Change] the server sends from connection time onward
     * (fire-and-forget, no replay, one message per row change in commit
     * order). A non-success response while opening the stream (e.g. a
     * missing/unknown bearer token, or `503` when no `Broadcaster` is
     * wired) becomes a typed [PgChangeFeedException] once the returned
     * [Sequence] is iterated — this call itself never sends a request or
     * throws. Once the stream is open, its end (regular server-side close,
     * or an incomplete trailing frame) ends the sequence without an
     * exception — the same fire-and-forget semantics as
     * [io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient.streamChanges];
     * a frame whose `data:` payload does not parse as the documented
     * `SPEC-021` shape still throws [PgChangeFeedMalformedResponseException],
     * because that is a protocol violation, not a stream end.
     */
    fun streamChanges(): Sequence<Change> = sequence {
        val request = TransportRequest(
            method = "GET",
            url = options.address.toString().trimEnd('/') + "/changes/stream",
            headers = mapOf("Authorization" to "Bearer ${options.apiToken}"),
            body = null,
        )
        val response = transport.open(request)
        if (response.statusCode !in 200..299) {
            throw buildException(response.statusCode, drainToString(response.nextLine))
        }
        while (true) {
            val frame = SseFrameParser.readFrame(response.nextLine) ?: break
            yield(parseChange(frame.data, response.statusCode))
        }
    }

    private fun parseChange(data: String, statusCode: Int): Change {
        val change = try {
            gson.fromJson(data, Change::class.java)
        } catch (ex: JsonSyntaxException) {
            throw PgChangeFeedMalformedResponseException(
                statusCode,
                "PG Change Feed SSE stream delivered a frame whose data payload is not valid " +
                    "JSON — a protocol violation outside SPEC-021's documented shape.",
                ex,
            )
        }
        return change
            ?: throw PgChangeFeedMalformedResponseException(
                statusCode,
                "PG Change Feed SSE stream delivered a frame whose data payload is empty or " +
                    "null — a protocol violation outside SPEC-021's documented shape.",
            )
    }

    private fun buildException(statusCode: Int, body: String): PgChangeFeedException {
        val message = extractErrorMessage(body)
        return when (statusCode) {
            400 -> PgChangeFeedBadRequestException(statusCode, message)
            401 -> PgChangeFeedUnauthorizedException(statusCode, message)
            403 -> PgChangeFeedForbiddenException(statusCode, message)
            404 -> PgChangeFeedNotFoundException(statusCode, message)
            500 -> PgChangeFeedServerErrorException(statusCode, message)
            else -> PgChangeFeedUnexpectedStatusException(statusCode, message)
        }
    }

    private fun extractErrorMessage(body: String): String =
        try {
            gson.fromJson(body, ErrorResponse::class.java)?.error ?: body
        } catch (ex: JsonSyntaxException) {
            body
        }

    /**
     * Reassembles the remaining lines of a non-success response into a
     * single body string — the error-response counterpart of
     * [io.github.pt9912.pgchangefeed.http.TransportResponse]'s already-
     * buffered `body`. A non-success `GET /changes/stream` response
     * carries a single-line `{"error": "<text>"}` body (`SPEC-021`, same
     * uniform shape as SPEC-018), never a stream — joining every remaining
     * line with `\n` round-trips that single line unchanged.
     */
    private fun drainToString(nextLine: () -> String?): String {
        val builder = StringBuilder()
        while (true) {
            val line = nextLine() ?: break
            if (builder.isNotEmpty()) {
                builder.append('\n')
            }
            builder.append(line)
        }
        return builder.toString()
    }

    private companion object {
        val gson = Gson()
    }
}
