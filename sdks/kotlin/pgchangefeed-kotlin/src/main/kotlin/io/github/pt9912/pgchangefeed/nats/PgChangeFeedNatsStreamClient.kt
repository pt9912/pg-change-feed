package io.github.pt9912.pgchangefeed.nats

import com.google.gson.Gson
import com.google.gson.JsonSyntaxException
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.nats.model.Change
import io.nats.client.Connection
import io.nats.client.Nats
import io.nats.client.Options
import java.time.Duration

/**
 * Client for the PG Change Feed live change stream over NATS: [streamChanges]
 * subscribes to the subject namespace `cdc.stream.<source_id>.<schema>.<table>`
 * (or a wildcard pattern over it, see [buildSourceSubject]/[ALL_SOURCES_SUBJECT])
 * and yields the ten message fields as [Change] — the same form as the gRPC
 * and SSE clients
 * ([io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient.streamChanges],
 * [io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient.streamChanges]).
 *
 * Authentication is connection-level, not per call: the NATS server rejects
 * the connection itself when it is configured with a token and the client's
 * token is missing or wrong — there is no per-message header to attach.
 * [PgChangeFeedClientOptions] is reused: [PgChangeFeedClientOptions.address] is
 * the NATS server URL (e.g. `nats://host:4222`) and
 * [PgChangeFeedClientOptions.apiToken] is the connection token.
 *
 * **Limits:** no delivery guarantee and no in-stream replay (fire-and-forget);
 * a disconnected or slow-reading consumer misses the affected messages
 * permanently. Missed changes remain recoverable through the read path
 * ([io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient.readChanges]).
 * A rejected connection (missing/wrong token) or a subscription that fails for
 * any other reason surfaces as whatever exception the underlying
 * `io.nats.client` call throws, propagated unwrapped — the stream is never
 * silently empty, the same as the gRPC client's `UNAUTHENTICATED`
 * `StatusException` and the SSE client's `401`. Only a message whose payload
 * cannot be read gets its own typed exception,
 * [PgChangeFeedNatsMalformedMessageException].
 */
class PgChangeFeedNatsStreamClient private constructor(
    private val transport: NatsStreamTransport,
    private val ownedConnection: Connection?,
) : AutoCloseable {

    /**
     * Convenience constructor: connects to [options]'s address with [options]'s
     * token and owns the resulting [Connection]. [close] closes that
     * connection. A rejected connection (missing/wrong token) throws
     * synchronously from this constructor itself, the moment
     * `io.nats.client.Nats.connect` throws; the class does not remap it (see
     * the class-level note on limits).
     */
    constructor(options: PgChangeFeedClientOptions) : this(connectOwned(options))

    /**
     * Unpacks the `(transport, connection)` pair [connectOwned] builds into the
     * two-argument primary constructor, so that `io.nats.client.Nats.connect`
     * runs exactly once. It is `private` because a public constructor taking a
     * bare pair would be ambiguous next to the [Connection] constructor.
     */
    private constructor(owned: Pair<NatsStreamTransport, Connection>) : this(owned.first, owned.second)

    /**
     * Advanced constructor: the [Connection] is passed in, not owned — the
     * caller controls connection lifetime and sharing (e.g. one connection
     * behind several subjects/subscriptions, or a connection already
     * authenticated some other way). [close] is then a no-op.
     */
    constructor(connection: Connection) : this(JnatsStreamTransport(connection), ownedConnection = null)

    /**
     * Constructor for the test source set: takes a [NatsStreamTransport]
     * directly and builds no connection — the same kind of transport interface
     * the gRPC and SSE clients use, because `io.nats.client.Connection` needs a
     * real server; faking at this boundary keeps the tests network-free.
     *
     * `internal` is a compile-time visibility boundary of the Kotlin compiler
     * (the test source set is a friend of the `internal` declarations of
     * `main`), not a JVM access restriction: in the compiled class file this
     * constructor is an ordinary `public` `<init>`, so a Java caller or
     * reflection could still invoke it with a [NatsStreamTransport].
     */
    internal constructor(transport: NatsStreamTransport) : this(transport, ownedConnection = null)

    /**
     * Subscribes to [subject] (default: [ALL_SOURCES_SUBJECT], every source and
     * table) and yields every [Change] the NATS server delivers from
     * subscription time onward (fire-and-forget, no replay, one message per row
     * change in commit order). Use [buildSubject]/[buildSourceSubject] to
     * narrow the subject to one table or one source. The returned [Sequence] is
     * cold — nothing subscribes until it is iterated.
     *
     * A subscription rejected or failed by the NATS server ends the sequence
     * with the underlying `io.nats.client` exception itself — see the
     * class-level note on limits. A message whose payload cannot be read throws
     * [PgChangeFeedNatsMalformedMessageException], because that is not a stream
     * end.
     */
    fun streamChanges(subject: String = ALL_SOURCES_SUBJECT): Sequence<Change> = sequence {
        val nextPayload = transport.subscribe(subject)
        while (true) {
            val payload = nextPayload() ?: break
            yield(parseChange(payload))
        }
    }

    private fun parseChange(payload: ByteArray): Change {
        val json = String(payload, Charsets.UTF_8)
        val change = try {
            gson.fromJson(json, Change::class.java)
        } catch (ex: JsonSyntaxException) {
            throw PgChangeFeedNatsMalformedMessageException(
                "PG Change Feed NATS stream delivered a message whose payload is " +
                    "not valid JSON.",
                ex,
            )
        }
        return change
            ?: throw PgChangeFeedNatsMalformedMessageException(
                "PG Change Feed NATS stream delivered a message whose payload is " +
                    "empty or null.",
            )
    }

    /** Closes the [Connection] this instance owns, if any (see the convenience constructor). */
    override fun close() {
        ownedConnection?.close()
    }

    companion object {
        /**
         * The root wildcard of the stream namespace — every source, every table
         * (`cdc.stream.>`). The default [streamChanges] subscribes to when no
         * narrower subject is supplied.
         */
        const val ALL_SOURCES_SUBJECT: String = "cdc.stream.>"

        private val invalidTokenChars = charArrayOf('.', '*', '>', ' ', '\t', '\n', '\r')
        private val gson = Gson()

        /**
         * Builds the four-token subject for one specific table:
         * `cdc.stream.<sourceId>.<schema>.<table>`. Each token is validated to
         * contain none of NATS's own token separator (`.`) or wildcard
         * characters (`*`, `>`) — a token carrying one of these would silently
         * change which subjects the resulting string matches, rather than fail
         * loudly.
         */
        fun buildSubject(sourceId: String, schema: String, table: String): String {
            validateToken(sourceId, "sourceId")
            validateToken(schema, "schema")
            validateToken(table, "table")
            return "cdc.stream.$sourceId.$schema.$table"
        }

        /**
         * Builds the three-token wildcard subject for every table of one
         * source: `cdc.stream.<sourceId>.>`.
         */
        fun buildSourceSubject(sourceId: String): String {
            validateToken(sourceId, "sourceId")
            return "cdc.stream.$sourceId.>"
        }

        private fun validateToken(value: String, paramName: String) {
            require(value.isNotBlank()) { "$paramName must not be blank." }
            require(value.none { it in invalidTokenChars }) {
                "$paramName must not contain '.', '*', '>', or whitespace — these are NATS " +
                    "subject separators/wildcards, not part of a token's own value."
            }
        }

        private fun connectOwned(options: PgChangeFeedClientOptions): Pair<NatsStreamTransport, Connection> {
            val jnatsOptions = Options.Builder()
                .server(options.address.toString())
                .token(options.apiToken)
                .build()
            val connection = Nats.connect(jnatsOptions)
            return JnatsStreamTransport(connection) to connection
        }
    }
}

/**
 * Transport interface between [PgChangeFeedNatsStreamClient] and the
 * `io.nats.client` wire — like [io.github.pt9912.pgchangefeed.sse.SseTransport]:
 * `io.nats.client.Connection`/`Subscription` have no pluggable handler, so
 * [PgChangeFeedNatsStreamClient] depends on this interface instead of the real
 * client directly, letting tests inject a network-free fake that hands back
 * canned payloads.
 *
 * [subscribe] returns a `() -> ByteArray?` next-payload supplier; a `null`
 * result ends the sequence (subscription closed), any thrown exception
 * propagates unmodified from [PgChangeFeedNatsStreamClient.streamChanges]'s
 * [Sequence].
 *
 * `internal` is a compile-time visibility boundary of the Kotlin compiler, not
 * a JVM access restriction; see the test-source-set constructor of
 * [PgChangeFeedNatsStreamClient] for the details, which apply here identically.
 */
internal fun interface NatsStreamTransport {
    fun subscribe(subject: String): () -> ByteArray?
}

/**
 * The real [NatsStreamTransport]: adapts an already connected
 * `io.nats.client.Connection` (either owned by [PgChangeFeedNatsStreamClient]'s
 * convenience constructor or passed in through its advanced constructor).
 * [io.nats.client.Subscription.nextMessage] is called with [Duration.ZERO],
 * which waits without a timeout; when the subscription itself ends,
 * `nextMessage` returns `null` and the sequence ends.
 */
internal class JnatsStreamTransport(private val connection: Connection) : NatsStreamTransport {
    override fun subscribe(subject: String): () -> ByteArray? {
        val subscription = connection.subscribe(subject)
        return { subscription.nextMessage(Duration.ZERO)?.data }
    }
}
