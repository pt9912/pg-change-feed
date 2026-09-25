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
 * Client for the PG Change Feed live change stream over Server-Sent Events:
 * [streamChanges] opens `GET /changes/stream`, assembles each SSE frame via
 * [SseFrameParser] and yields the ten message fields as [Change]. Errors are
 * the typed [PgChangeFeedException] set that
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient] throws (the
 * `{"error": "<text>"}` error body), not a second hierarchy.
 *
 * [streamChanges] returns a [Sequence]`<`[Change]`>`, not a
 * `kotlinx.coroutines.flow.Flow` like
 * [io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient.streamChanges]:
 * the generated gRPC coroutine stub offers a suspending, non-blocking call,
 * while `java.net.http.HttpClient.send()` (the JDK client
 * [io.github.pt9912.pgchangefeed.http.JdkHttpTransport] wraps) is an ordinary
 * blocking call, like every method of
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient]. A [Sequence]
 * keeps the synchronous execution model of the HTTP client while staying
 * **cold**: nothing runs — no request is sent, no frame is parsed — until the
 * returned [Sequence] is iterated.
 *
 * The bearer token is sent in the `Authorization` header as
 * `Bearer <token>`, same as
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient]. There is no
 * global or static state; a process can hold several independently configured
 * instances at once.
 *
 * **Limits:** no delivery guarantee and no in-stream replay — a disconnected
 * or slow-reading consumer misses the affected messages permanently. Missed
 * changes remain recoverable through the read path
 * (`io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient.readChanges`).
 */
class PgChangeFeedSseClient internal constructor(
    private val transport: SseTransport,
    private val options: PgChangeFeedClientOptions,
) {
    /**
     * Public constructor — wraps the caller's `java.net.http.HttpClient` in the
     * real [JdkSseTransport]. The [HttpClient] is passed in, not owned — the
     * caller controls its lifetime, connection pooling and any redirect/proxy
     * configuration; this class never closes it. See [SseTransport] for why the
     * client depends on a transport interface rather than on [HttpClient]
     * directly.
     */
    constructor(httpClient: HttpClient, options: PgChangeFeedClientOptions) :
        this(JdkSseTransport(httpClient), options)

    /**
     * Opens `GET /changes/stream` and yields every [Change] the server sends
     * from connection time onward (fire-and-forget, no replay, one message per
     * row change in commit order). A non-success response while opening the
     * stream (e.g. a missing/unknown bearer token, or `503` when the stream is
     * not available) becomes a typed [PgChangeFeedException] once the returned
     * [Sequence] is iterated — this call itself never sends a request or
     * throws. Once the stream is open, its end (regular server-side close, or
     * an incomplete trailing frame) ends the sequence without an exception —
     * the same fire-and-forget behavior as
     * [io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient.streamChanges];
     * a frame whose `data:` payload cannot be read still throws
     * [PgChangeFeedMalformedResponseException], because that is not a stream
     * end.
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
                    "JSON.",
                ex,
            )
        }
        return change
            ?: throw PgChangeFeedMalformedResponseException(
                statusCode,
                "PG Change Feed SSE stream delivered a frame whose data payload is empty or " +
                    "null.",
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
     * Reassembles the remaining lines of a non-success response into a single
     * body string — the counterpart of
     * [io.github.pt9912.pgchangefeed.http.TransportResponse]'s already
     * buffered `body`. A non-success `GET /changes/stream` response carries a
     * single-line `{"error": "<text>"}` body, never a stream — joining every
     * remaining line with `\n` round-trips that single line unchanged.
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
