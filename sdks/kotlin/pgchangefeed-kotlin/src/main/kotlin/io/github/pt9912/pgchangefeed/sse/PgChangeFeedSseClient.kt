package io.github.pt9912.pgchangefeed.sse

import com.google.gson.Gson
import com.google.gson.JsonSyntaxException
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.TlsSupport
import io.github.pt9912.pgchangefeed.http.PgChangeFeedBadRequestException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedForbiddenException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedMalformedResponseException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedNotFoundException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedServerErrorException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnauthorizedException
import io.github.pt9912.pgchangefeed.http.PgChangeFeedUnexpectedStatusException
import io.github.pt9912.pgchangefeed.http.TransportRequest
import io.github.pt9912.pgchangefeed.http.parseErrorBody
import io.github.pt9912.pgchangefeed.http.percentEncode
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
     * Convenience constructor: builds its own `java.net.http.HttpClient` from
     * [options]. An `https` address connects over TLS: the trust anchors of the
     * Java runtime apply, or exactly the certificates of
     * [PgChangeFeedClientOptions.trustAnchorFile] when it is set; chain,
     * validity period and server name are always checked, and a failed check
     * surfaces as the `IOException` of `java.net.http.HttpClient`, never as a
     * plaintext retry. A trust anchor with an `http` address throws an
     * [IllegalArgumentException]. A trust anchor in the options does not apply
     * to an [HttpClient] passed to the constructor above.
     */
    constructor(options: PgChangeFeedClientOptions) :
        this(JdkSseTransport(TlsSupport.httpClient(options)), options)

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
     *
     * [target] selects the delivery target of a change: a set value (sent as
     * the query parameter `target`) delivers only changes routed to that
     * target; left `null` (the default) the request carries no query and
     * delivers every change. A target no change carries delivers nothing and
     * raises no error.
     *
     * [schema] and [table] are each optional and independent, sent as the
     * query parameters `schema` and `table` only when set: a set [schema]
     * without [table] delivers every table of that schema, a set [table]
     * without [schema] delivers every table of that name regardless of
     * schema, both set delivers exactly one table. They combine with [target]
     * as a conjunction. Pass the three by name.
     */
    fun streamChanges(
        target: String? = null,
        schema: String? = null,
        table: String? = null,
    ): Sequence<Change> = sequence {
        val params = listOf("schema" to schema, "table" to table, "target" to target)
            .mapNotNull { (name, value) -> value?.let { "$name=" + percentEncode(it) } }
        val query = if (params.isEmpty()) "" else "?" + params.joinToString("&")
        val request = TransportRequest(
            method = "GET",
            url = options.address.toString().trimEnd('/') + "/changes/stream" + query,
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
        val (message, messageCode) = parseErrorBody(body)
        return when (statusCode) {
            400 -> PgChangeFeedBadRequestException(statusCode, message, messageCode)
            401 -> PgChangeFeedUnauthorizedException(statusCode, message, messageCode)
            403 -> PgChangeFeedForbiddenException(statusCode, message, messageCode)
            404 -> PgChangeFeedNotFoundException(statusCode, message, messageCode)
            500 -> PgChangeFeedServerErrorException(statusCode, message, messageCode)
            else -> PgChangeFeedUnexpectedStatusException(statusCode, message, messageCode)
        }
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
