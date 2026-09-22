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
 * Public entry point for the PG Change Feed NATS full-content stream
 * (`SPEC-024`, `LH-FA-SST-008`): [streamChanges] subscribes to the
 * four-token subject namespace `cdc.stream.<source_id>.<schema>.<table>`
 * (or a wildcard pattern over it, see [buildSourceSubject]/[ALL_SOURCES_SUBJECT])
 * and yields the ten `SPEC-024` message fields as [Change] — the same
 * idiomatic form as the existing gRPC/SSE surfaces
 * ([io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient.streamChanges],
 * [io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient.streamChanges]).
 *
 * Authentication is connection-level, not per-call (`SPEC-024`): the NATS
 * server rejects the connection itself when a server-wide token is
 * configured and the client's token is missing or wrong — there is no
 * per-message header to attach, unlike the HTTP/gRPC/SSE surfaces' bearer
 * token. [PgChangeFeedClientOptions] is still the shared denominator this
 * class reuses: [PgChangeFeedClientOptions.address] becomes the NATS server
 * URL (e.g. `nats://host:4222`) and [PgChangeFeedClientOptions.apiToken]
 * becomes the connection token — the same two-value shape as every other
 * surface, carried over a different transport.
 *
 * **Boundary (`SPEC-024`, `LH-FA-SST-008`):** no delivery guarantee and no
 * in-stream replay (Core NATS, fire-and-forget); a disconnected or
 * slow-reading consumer misses the affected messages permanently. Missed
 * changes remain recoverable through the existing read path
 * ([io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient.readChanges]).
 * A rejected connection (missing/wrong token) or a subscription that fails
 * for any other reason surfaces as whatever exception the underlying
 * `io.nats.client` call throws, propagated unwrapped from [streamChanges]'s
 * [Sequence] once iterated — the same "not a swallowed empty stream"
 * boundary as the gRPC surface's `UNAUTHENTICATED` `StatusException` and the
 * SSE surface's `401`: NATS connection errors are not HTTP status codes or
 * gRPC statuses, and inventing a parallel mapping here would only risk
 * drifting from what `io.nats.client` itself already throws (same reasoning
 * as the C# sibling's `PgChangeFeed.Client.Nats.PgChangeFeedNatsStreamClient`).
 * Only a message whose payload does not parse as the documented `SPEC-024`
 * shape gets its own typed exception, [PgChangeFeedNatsMalformedMessageException] —
 * that is a protocol violation, not a connection failure.
 *
 * Draht-Kenntnis-Vorbild (gelesen, nicht importiert — `ADR-0109`
 * Festlegung 3): `examples/kotlin/nats-stream-client/src/main/kotlin/cdcexamples/natsstream/{Format,Main}.kt`.
 */
class PgChangeFeedNatsStreamClient private constructor(
    private val transport: NatsStreamTransport,
    private val ownedConnection: Connection?,
) : AutoCloseable {

    /**
     * Convenience constructor: connects to [options]'s address with
     * [options]'s token (`SPEC-024`'s connection-level auth) and owns the
     * resulting [Connection]. [close] closes that connection. A rejected
     * connection (missing/wrong server-wide token) throws synchronously from
     * this constructor itself, the same moment `io.nats.client.Nats.connect`
     * throws — this is real `io.nats.client` behavior, not something this
     * class remaps (see the class-level boundary note).
     */
    constructor(options: PgChangeFeedClientOptions) : this(connectOwned(options))

    /**
     * Unpacks the `(transport, connection)` pair [connectOwned] builds into
     * this class's two-argument primary constructor — a single-expression
     * delegate target so [io.nats.client.Nats.connect] runs exactly once.
     * Kept `private`: a bare [Connection]-typed public constructor would
     * collide with this one at the same arity/shape if both existed
     * directly (see [PgChangeFeedNatsStreamClient] taking [Connection]
     * below, which stays distinct only because it carries `ownedConnection = null`
     * rather than delegating through this pair).
     */
    private constructor(owned: Pair<NatsStreamTransport, Connection>) : this(owned.first, owned.second)

    /**
     * Advanced constructor: the [Connection] is injected, not owned — the
     * caller controls connection lifetime and sharing (e.g. one connection
     * behind several subjects/subscriptions, or a connection already
     * authenticated some other way). [close] is then a no-op.
     */
    constructor(connection: Connection) : this(JnatsStreamTransport(connection), ownedConnection = null)

    /**
     * Test-only constructor: injects a [NatsStreamTransport] directly,
     * bypassing connection construction entirely — the same seam
     * [io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient]'s
     * `GrpcStreamTransport` and
     * [io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient]'s
     * `SseTransport` provide for their surfaces, and for the same reason:
     * `io.nats.client.Connection`/`Subscription` need a real server to
     * connect to; faking at this boundary keeps the tests genuinely
     * network-free.
     *
     * `internal` here — as with those two sibling seams — is a
     * **compile-time** Kotlin-compiler visibility boundary against other
     * Kotlin modules' metadata (the Kotlin Gradle plugin's default
     * main/test sourceSet association makes the test sourceSet a friend of
     * `internal` declarations in `main`). It is **not** a JVM bytecode
     * access restriction: in the compiled class file this constructor is an
     * ordinary `public` `<init>` symbol (constructors are always named
     * `<init>` in bytecode and are not covered by Kotlin's `internal`
     * name-mangling) — a Java caller, or reflection from any language, can
     * still invoke it directly given a [NatsStreamTransport] implementation.
     */
    internal constructor(transport: NatsStreamTransport) : this(transport, ownedConnection = null)

    /**
     * Subscribes to [subject] (default: [ALL_SOURCES_SUBJECT], every source
     * and table) and yields every [Change] the NATS server delivers from
     * subscription time onward (`SPEC-024`: fire-and-forget, no replay, one
     * message per row change in commit order). Use [buildSubject]/
     * [buildSourceSubject] to narrow the subject to one table or one
     * source. The returned [Sequence] is cold — nothing subscribes until it
     * is iterated.
     *
     * A subscription rejected or failed by the NATS server (missing/wrong
     * token, `SPEC-024` Negative) ends the sequence with the underlying
     * `io.nats.client` exception itself — see the class-level boundary
     * note. A message whose payload does not parse as the documented
     * `SPEC-024` shape throws [PgChangeFeedNatsMalformedMessageException],
     * because that is a protocol violation, not a stream end.
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
                "PG Change Feed NATS full-content stream delivered a message whose payload is " +
                    "not valid JSON — a protocol violation outside SPEC-024's documented shape.",
                ex,
            )
        }
        return change
            ?: throw PgChangeFeedNatsMalformedMessageException(
                "PG Change Feed NATS full-content stream delivered a message whose payload is " +
                    "empty or null — a protocol violation outside SPEC-024's documented shape.",
            )
    }

    /** Closes the [Connection] this instance owns, if any (see the convenience constructor). */
    override fun close() {
        ownedConnection?.close()
    }

    companion object {
        /**
         * The full-content namespace's root wildcard — every source, every
         * table (`SPEC-024`: `cdc.stream.>`). The default [streamChanges]
         * subscribes to when no narrower subject is supplied.
         */
        const val ALL_SOURCES_SUBJECT: String = "cdc.stream.>"

        private val invalidTokenChars = charArrayOf('.', '*', '>', ' ', '\t', '\n', '\r')
        private val gson = Gson()

        /**
         * Builds the four-token `SPEC-024` subject for one specific table:
         * `cdc.stream.<sourceId>.<schema>.<table>`. Each token is validated
         * to contain none of NATS's own token separator (`.`) or wildcard
         * characters (`*`, `>`) — a token carrying one of these would
         * silently change which subjects the resulting string matches,
         * rather than fail loudly.
         */
        fun buildSubject(sourceId: String, schema: String, table: String): String {
            validateToken(sourceId, "sourceId")
            validateToken(schema, "schema")
            validateToken(table, "table")
            return "cdc.stream.$sourceId.$schema.$table"
        }

        /**
         * Builds the three-token wildcard subject for every table of one
         * source: `cdc.stream.<sourceId>.>` (`SPEC-024`).
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
 * Transport seam between [PgChangeFeedNatsStreamClient] and the actual
 * `io.nats.client` wire — same reasoning as
 * [io.github.pt9912.pgchangefeed.sse.SseTransport]'s KDoc:
 * `io.nats.client.Connection`/`Subscription` carry no pluggable-handler
 * concept, so [PgChangeFeedNatsStreamClient] depends on this seam instead of
 * the real client directly, letting tests inject a network-free fake that
 * hands back canned payloads without any socket at all.
 *
 * [subscribe] returns a `() -> ByteArray?` next-payload supplier — the same
 * shape [io.github.pt9912.pgchangefeed.sse.SseTransport]'s `nextLine`
 * carries, adapted from lines to raw message bytes; a `null` result ends the
 * sequence (subscription closed), any thrown exception propagates from
 * [PgChangeFeedNatsStreamClient.streamChanges]'s [Sequence] unmodified (see
 * that method's boundary note).
 *
 * `internal` here — as with the sibling `SseTransport`/`GrpcStreamTransport`
 * seams — is a **compile-time** Kotlin-compiler visibility boundary against
 * other Kotlin modules' metadata, **not** a JVM bytecode access restriction:
 * see [PgChangeFeedNatsStreamClient]'s test-only constructor KDoc for the
 * full explanation, which applies here identically.
 */
internal fun interface NatsStreamTransport {
    fun subscribe(subject: String): () -> ByteArray?
}

/**
 * The real [NatsStreamTransport]: adapts a caller-supplied
 * `io.nats.client.Connection` (already connected — either owned by
 * [PgChangeFeedNatsStreamClient]'s convenience constructor or injected via
 * its advanced constructor). [io.nats.client.Subscription.nextMessage] is
 * called with [Duration.ZERO], which waits real-unbounded rather than
 * timing out — the same choice
 * `examples/kotlin/nats-stream-client/src/main/kotlin/cdcexamples/natsstream/Main.kt`
 * makes (read as Draht-Kenntnis, not imported — `ADR-0109` Festlegung 3),
 * and matches Core NATS's fire-and-forget semantics: there is no next
 * message to wait *for* once the subscription itself ends, at which point
 * `nextMessage` returns `null`.
 */
internal class JnatsStreamTransport(private val connection: Connection) : NatsStreamTransport {
    override fun subscribe(subject: String): () -> ByteArray? {
        val subscription = connection.subscribe(subject)
        return { subscription.nextMessage(Duration.ZERO)?.data }
    }
}
