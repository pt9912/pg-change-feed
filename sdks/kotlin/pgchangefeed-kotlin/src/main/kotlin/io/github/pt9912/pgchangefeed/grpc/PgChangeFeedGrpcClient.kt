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
 * Public entry point for the PG Change Feed live-change stream (`SPEC-020`,
 * `LH-FA-SST-008`): [streamChanges] opens the `ChangeStream/StreamChanges`
 * server-streaming RPC and yields the generated [Change] message — the ten
 * `SPEC-020` fields (`change_id`, `transaction_id`, `source_table_id`,
 * `sequence`, `operation`, `old_image`, `new_image`, `schema_version`,
 * `schema`, `table`) — unmapped, exactly as the wire defines them. There is
 * no separate DTO layer here (unlike
 * `io.github.pt9912.pgchangefeed.http.model`): the generated stub already
 * is a typed, versioned representation of the wire schema; introducing a
 * second, hand-mapped type would only risk drifting from it.
 *
 * The bearer token is sent in the `authorization` gRPC metadata entry as
 * `Bearer <token>` (`SPEC-020`) on every call — there is no global or
 * static state; a process can hold several independently configured
 * instances at once.
 *
 * **Boundary (`SPEC-020`, `LH-FA-SST-008`):** the stream carries no replay
 * and no table-granular filtering. A consumer that needs either uses the
 * existing read path
 * (`io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient.readChanges`),
 * not this stream.
 *
 * `.proto` liegt NICHT im committeten Baum dieses Pakets — der generierte
 * Coroutine-Stub ([ChangeStreamGrpcKt]) und der Nachrichtentyp ([Change])
 * entstehen im Bau aus dem zusätzlichen, benannten Docker-Bau-Kontext
 * `proto` (`ADR-0109` Festlegung 3, `sdks/kotlin/Dockerfile`).
 *
 * Draht-Kenntnis-Vorbild (gelesen, nicht importiert — `ADR-0109`
 * Festlegung 3): `examples/kotlin/grpc-client/src/main/kotlin/cdcexamples/grpc/Main.kt`.
 */
class PgChangeFeedGrpcClient private constructor(
    private val transport: GrpcStreamTransport,
    private val options: PgChangeFeedClientOptions,
    private val ownedChannel: ManagedChannel?,
) : AutoCloseable {

    /**
     * Convenience constructor: opens and owns its own [ManagedChannel]
     * against [options]'s address (`usePlaintext()` — the PG Change Feed
     * server speaks plaintext gRPC over HTTP/2 by default, the same choice
     * as the reference clients and the C#/Python SDKs). [close] shuts down
     * that channel. A TLS-terminated deployment builds its own [Channel]
     * and uses the advanced constructor below instead.
     *
     * [options].address must carry a resolvable host and port (e.g.
     * `http://pg-change-feed:50051`) — only the authority part is used, the
     * scheme is ignored; this mirrors [options] being the shared connection
     * denominator between the HTTP and gRPC surfaces (`ADR-0109`
     * Festlegung 1).
     */
    constructor(options: PgChangeFeedClientOptions) : this(options, buildOwnedChannel(options))

    private constructor(options: PgChangeFeedClientOptions, channel: ManagedChannel) :
        this(GeneratedStubTransport(channel), options, channel)

    /**
     * Advanced constructor: the [Channel] is injected, not owned — the
     * caller controls channel lifetime and sharing (e.g. one channel behind
     * several client surfaces, or a TLS-terminated channel). [close] is
     * then a no-op.
     */
    constructor(channel: Channel, options: PgChangeFeedClientOptions) :
        this(GeneratedStubTransport(channel), options, ownedChannel = null)

    /**
     * Test-only constructor: injects a [GrpcStreamTransport] directly,
     * bypassing channel construction entirely — the same seam
     * `io.github.pt9912.pgchangefeed.http.HttpTransport` uses for the HTTP
     * client surface, and for the same reason: the generated coroutine stub
     * exposes no pluggable-handler concept the way `Grpc.Core.CallInvoker`
     * does on the C# side; faking at this boundary keeps the tests
     * genuinely network-free without hand-rolling a fake low-level
     * [Channel]/`ClientCall` pair.
     *
     * `internal` here — as with `HttpTransport`'s test-only constructor —
     * is a **compile-time** Kotlin-compiler visibility boundary against
     * other Kotlin modules' metadata (the Kotlin Gradle plugin's default
     * main/test sourceSet association makes the test sourceSet a friend of
     * `internal` declarations in `main`). It is **not** a JVM bytecode
     * access restriction: in the compiled class file this constructor is
     * an ordinary `public` `<init>` symbol (constructors are always named
     * `<init>` in bytecode and are not covered by Kotlin's `internal`
     * name-mangling) — a Java caller, or reflection from any language, can
     * still invoke it directly given a [GrpcStreamTransport] implementation.
     */
    internal constructor(transport: GrpcStreamTransport, options: PgChangeFeedClientOptions) :
        this(transport, options, ownedChannel = null)

    /**
     * `StreamChanges` — opens the server-streaming RPC and yields every
     * [Change] the server sends from connection time onward (`SPEC-020`:
     * fire-and-forget, no replay, one message per row change in commit
     * order). The request carries no filter (`SPEC-020`: table-granular
     * filtering is not part of this version).
     *
     * A missing or invalid bearer token ends the call with gRPC status
     * `UNAUTHENTICATED` (`SPEC-020` Negative) — this surfaces as an
     * `io.grpc.StatusException` from the returned [Flow] itself once
     * collected, not a swallowed empty stream.
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
 * Transport seam between [PgChangeFeedGrpcClient] and the generated
 * coroutine stub (`ChangeStreamGrpcKt.ChangeStreamCoroutineStub`) — see
 * [PgChangeFeedGrpcClient]'s test-only constructor KDoc for why this seam
 * exists and what `internal` does and does not restrict.
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
