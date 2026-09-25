package io.github.pt9912.pgchangefeed.grpc

import cdc.stream.v1.ChangeStreamGrpcKt
import cdc.stream.v1.Changestream.Change
import cdc.stream.v1.Changestream.StreamChangesRequest
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.grpc.Channel
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import io.grpc.Metadata
import kotlinx.coroutines.flow.Flow

/**
 * Client for the PG Change Feed live change stream over gRPC: [streamChanges]
 * opens the `ChangeStream/StreamChanges` server-streaming call and yields the
 * generated [Change] message with the ten fields `change_id`,
 * `transaction_id`, `source_table_id`, `sequence`, `operation`, `old_image`,
 * `new_image`, `schema_version`, `schema` and `table`, exactly as the server
 * sends them. There is no separate DTO layer (unlike
 * `io.github.pt9912.pgchangefeed.http.model`): the generated message already
 * is the typed form of the stream schema.
 *
 * The bearer token is sent in the `authorization` call metadata entry as
 * `Bearer <token>` on every call — there is no global or static state; a
 * process can hold several independently configured instances at once.
 *
 * **Limits:** the stream carries no replay and cannot be filtered by table. A
 * consumer that needs either uses the read path
 * (`io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient.readChanges`),
 * not this stream.
 *
 * The `.proto` file is not part of the committed tree of this package: the
 * generated coroutine stub ([ChangeStreamGrpcKt]) and the message type
 * ([Change]) are generated during the build from the additional named Docker
 * build context `proto` (`sdks/kotlin/Dockerfile`).
 */
class PgChangeFeedGrpcClient private constructor(
    private val transport: GrpcStreamTransport,
    private val options: PgChangeFeedClientOptions,
    private val ownedChannel: ManagedChannel?,
) : AutoCloseable {

    /**
     * Convenience constructor: opens and owns its own [ManagedChannel] against
     * [options]'s address (`usePlaintext()` — the PG Change Feed server speaks
     * plaintext gRPC over HTTP/2). [close] shuts down that channel. A
     * TLS-terminated deployment builds its own [Channel] and uses the advanced
     * constructor below instead.
     *
     * [options].address must carry a resolvable host and port (e.g.
     * `http://pg-change-feed:50051`) — only the authority part is used, the
     * scheme is ignored.
     */
    constructor(options: PgChangeFeedClientOptions) : this(options, buildOwnedChannel(options))

    private constructor(options: PgChangeFeedClientOptions, channel: ManagedChannel) :
        this(GeneratedStubTransport(channel), options, channel)

    /**
     * Advanced constructor: the [Channel] is passed in, not owned — the caller
     * controls channel lifetime and sharing (e.g. one channel behind several
     * clients, or a TLS-terminated channel). [close] is then a no-op.
     */
    constructor(channel: Channel, options: PgChangeFeedClientOptions) :
        this(GeneratedStubTransport(channel), options, ownedChannel = null)

    /**
     * Constructor for the test source set: takes a [GrpcStreamTransport]
     * directly and builds no channel — the same transport interface the HTTP
     * client uses (`io.github.pt9912.pgchangefeed.http.HttpTransport`), so
     * tests need no network and no hand-rolled fake [Channel].
     *
     * `internal` is a compile-time visibility boundary of the Kotlin compiler
     * (the test source set is a friend of the `internal` declarations of
     * `main`), not a JVM access restriction: in the compiled class file this
     * constructor is an ordinary `public` `<init>`, so a Java caller or
     * reflection could still invoke it with a [GrpcStreamTransport].
     */
    internal constructor(transport: GrpcStreamTransport, options: PgChangeFeedClientOptions) :
        this(transport, options, ownedChannel = null)

    /**
     * Opens the server-streaming call and yields every [Change] the server
     * sends from connection time onward: fire-and-forget, no replay, one
     * message per row change in commit order. The request carries no filter.
     *
     * A missing or invalid bearer token ends the call with gRPC status
     * `UNAUTHENTICATED` — this surfaces as an `io.grpc.StatusException` from
     * the returned [Flow] once collected, not as a silently empty stream.
     */
    fun streamChanges(): Flow<Change> {
        val headers = Metadata().apply {
            put(AUTHORIZATION_METADATA_ENTRY, BEARER_PREFIX + options.apiToken)
        }
        return transport.streamChanges(headers)
    }

    /** Shuts down the channel this instance owns, if any (see the convenience constructor). */
    override fun close() {
        ownedChannel?.shutdownNow()
    }

    private companion object {
        const val AUTHORIZATION_METADATA_KEY = "authorization"
        const val BEARER_PREFIX = "Bearer "
        val AUTHORIZATION_METADATA_ENTRY: Metadata.Key<String> =
            Metadata.Key.of(AUTHORIZATION_METADATA_KEY, Metadata.ASCII_STRING_MARSHALLER)

        fun buildOwnedChannel(options: PgChangeFeedClientOptions): ManagedChannel {
            val host = options.address.host
            val port = options.address.port
            require(host != null && port != -1) {
                "options.address must carry a resolvable host and port (got '${options.address}')"
            }
            return ManagedChannelBuilder.forAddress(host, port).usePlaintext().build()
        }
    }
}

/**
 * Transport interface between [PgChangeFeedGrpcClient] and the generated
 * coroutine stub (`ChangeStreamGrpcKt.ChangeStreamCoroutineStub`); see the
 * test-source-set constructor of [PgChangeFeedGrpcClient] for why it exists
 * and what `internal` does and does not restrict.
 */
internal fun interface GrpcStreamTransport {
    fun streamChanges(headers: Metadata): Flow<Change>
}

/** The real [GrpcStreamTransport]: wraps the generated coroutine stub. */
internal class GeneratedStubTransport(channel: Channel) : GrpcStreamTransport {
    private val stub = ChangeStreamGrpcKt.ChangeStreamCoroutineStub(channel)

    override fun streamChanges(headers: Metadata): Flow<Change> =
        stub.streamChanges(StreamChangesRequest.getDefaultInstance(), headers)
}
