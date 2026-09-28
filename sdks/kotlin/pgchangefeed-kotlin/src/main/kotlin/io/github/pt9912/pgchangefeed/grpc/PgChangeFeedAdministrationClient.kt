package io.github.pt9912.pgchangefeed.grpc

import cdc.administration.v1.AdministrationGrpcKt
import cdc.administration.v1.AdministrationOuterClass.AcknowledgeConsumerRequest
import cdc.administration.v1.AdministrationOuterClass.AcknowledgeConsumerResponse
import cdc.administration.v1.AdministrationOuterClass.DiagnoseRequest
import cdc.administration.v1.AdministrationOuterClass.DiagnoseResponse
import cdc.administration.v1.AdministrationOuterClass.DisableTableRequest
import cdc.administration.v1.AdministrationOuterClass.DisableTableResponse
import cdc.administration.v1.AdministrationOuterClass.EnableTableRequest
import cdc.administration.v1.AdministrationOuterClass.EnableTableResponse
import cdc.administration.v1.AdministrationOuterClass.GetConsumerPositionRequest
import cdc.administration.v1.AdministrationOuterClass.GetConsumerPositionResponse
import cdc.administration.v1.AdministrationOuterClass.GetTableStatusRequest
import cdc.administration.v1.AdministrationOuterClass.GetTableStatusResponse
import cdc.administration.v1.AdministrationOuterClass.ListTablesRequest
import cdc.administration.v1.AdministrationOuterClass.ListTablesResponse
import cdc.administration.v1.AdministrationOuterClass.ReadChangesRequest
import cdc.administration.v1.AdministrationOuterClass.ReadChangesResponse
import cdc.administration.v1.AdministrationOuterClass.RegisterConsumerRequest
import cdc.administration.v1.AdministrationOuterClass.RegisterConsumerResponse
import cdc.administration.v1.AdministrationOuterClass.RemoveConsumerRequest
import cdc.administration.v1.AdministrationOuterClass.RemoveConsumerResponse
import cdc.administration.v1.AdministrationOuterClass.RunRetentionRequest
import cdc.administration.v1.AdministrationOuterClass.RunRetentionResponse
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.grpc.Channel
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import io.grpc.Metadata
import io.grpc.Status
import io.grpc.StatusException

/**
 * Client for the PG Change Feed `Administration` gRPC service: one method
 * per RPC — [registerConsumer], [acknowledgeConsumer], [getConsumerPosition],
 * [removeConsumer], [enableTable], [disableTable], [getTableStatus],
 * [listTables], [runRetention], [readChanges] and [diagnose]. Requests and
 * responses are the generated protobuf messages of `cdc.administration.v1`
 * unchanged — like [PgChangeFeedGrpcClient], there is no separate DTO layer:
 * a protobuf message is already a typed Kotlin class, and a hand-written
 * mirror would only duplicate eleven message shapes without a
 * deserialization need to justify it (unlike
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient], whose DTOs
 * exist because JSON needs a target type to deserialize into). Every
 * non-`OK` gRPC status becomes a typed [PgChangeFeedGrpcException] subtype
 * instead of a result type mixed with the success path, consistent with
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient]'s exception
 * design for the HTTP surface.
 *
 * The bearer token is sent in the `authorization` call metadata entry as
 * `Bearer <token>` on every call, from [PgChangeFeedClientOptions] supplied
 * at construction — there is no global or static state; a process can hold
 * several independently configured instances at once.
 *
 * Because of a name clash between the `.proto` file `administration.proto`
 * and the service `Administration` it defines, `protoc`'s Kotlin/Java codegen
 * places every message type under `AdministrationOuterClass.*`, not
 * `Administration.*`, while the coroutine stub keeps the name
 * `AdministrationGrpcKt.AdministrationCoroutineStub` — the same generated
 * shape `examples/kotlin/grpc-client` already relies on.
 */
class PgChangeFeedAdministrationClient private constructor(
    private val transport: AdministrationTransport,
    private val options: PgChangeFeedClientOptions,
    private val ownedChannel: ManagedChannel?,
) : AutoCloseable {

    /**
     * Convenience constructor: opens and owns its own [ManagedChannel]
     * against [options]'s address (`usePlaintext()` — the PG Change Feed
     * server speaks plaintext gRPC over HTTP/2). [close] shuts down that
     * channel. A TLS-terminated deployment builds its own [Channel] and uses
     * the advanced constructor below instead.
     *
     * [options].address must carry a resolvable host and port (e.g.
     * `http://pg-change-feed:50051`) — only the authority part is used, the
     * scheme is ignored.
     */
    constructor(options: PgChangeFeedClientOptions) : this(options, buildOwnedChannel(options))

    private constructor(options: PgChangeFeedClientOptions, channel: ManagedChannel) :
        this(GeneratedAdministrationStubTransport(channel), options, channel)

    /**
     * Advanced constructor: the [Channel] is passed in, not owned — the
     * caller controls channel lifetime and sharing (e.g. one channel behind
     * several client surfaces, or a TLS-terminated channel). [close] is then
     * a no-op.
     */
    constructor(channel: Channel, options: PgChangeFeedClientOptions) :
        this(GeneratedAdministrationStubTransport(channel), options, ownedChannel = null)

    /**
     * Constructor for the test source set: takes an [AdministrationTransport]
     * directly and builds no channel — see
     * [PgChangeFeedGrpcClient]'s test-only constructor KDoc for why `internal`
     * here is a compile-time visibility boundary of the Kotlin compiler, not
     * a JVM access restriction.
     */
    internal constructor(transport: AdministrationTransport, options: PgChangeFeedClientOptions) :
        this(transport, options, ownedChannel = null)

    /**
     * Registers a consumer, a named reader whose position the server keeps
     * (admin token). Registering an existing consumer changes nothing; the
     * response reports it with `alreadyRegistered`.
     */
    suspend fun registerConsumer(request: RegisterConsumerRequest): RegisterConsumerResponse =
        callAsync { transport.registerConsumer(request, headers()) }

    /**
     * Stores the position up to which a consumer has processed a source
     * (admin token). Repeating the stored position has no effect; an
     * earlier position, or a position of another source, raises
     * [PgChangeFeedGrpcInvalidArgumentException].
     */
    suspend fun acknowledgeConsumer(request: AcknowledgeConsumerRequest): AcknowledgeConsumerResponse =
        callAsync { transport.acknowledgeConsumer(request, headers()) }

    /**
     * Reads the stored position of a consumer (reader or admin token).
     * `acknowledged` is `false` for a consumer that never acknowledged;
     * `offset` is then the starting position.
     */
    suspend fun getConsumerPosition(request: GetConsumerPositionRequest): GetConsumerPositionResponse =
        callAsync { transport.getConsumerPosition(request, headers()) }

    /**
     * Removes a consumer (admin token); `removed` is `false` for one that
     * was never registered.
     */
    suspend fun removeConsumer(request: RemoveConsumerRequest): RemoveConsumerResponse =
        callAsync { transport.removeConsumer(request, headers()) }

    /**
     * Starts capturing a table of a source (admin token). `alreadyEnabled`
     * is `true` when the table was captured already; a table that does not
     * exist in the source database raises [PgChangeFeedGrpcNotFoundException].
     */
    suspend fun enableTable(request: EnableTableRequest): EnableTableResponse =
        callAsync { transport.enableTable(request, headers()) }

    /**
     * Stops capturing a table (admin token). `retained` is `true` when
     * changes already stored for the table remain readable; a table that
     * does not exist in the source database raises
     * [PgChangeFeedGrpcNotFoundException].
     */
    suspend fun disableTable(request: DisableTableRequest): DisableTableResponse =
        callAsync { transport.disableTable(request, headers()) }

    /**
     * Tells whether a table is captured (`enabled`) or no longer captured
     * with stored changes remaining (`retained`) (reader or admin token). A
     * table that was never enabled reports both as `false`; a table that
     * does not exist in the source database raises
     * [PgChangeFeedGrpcNotFoundException].
     */
    suspend fun getTableStatus(request: GetTableStatusRequest): GetTableStatusResponse =
        callAsync { transport.getTableStatus(request, headers()) }

    /**
     * Lists the captured tables and the tables that are no longer captured
     * but whose stored changes remain (reader or admin token).
     */
    suspend fun listTables(request: ListTablesRequest): ListTablesResponse =
        callAsync { transport.listTables(request, headers()) }

    /**
     * Deletes the stored changes of a source that are older than
     * `minAgeNanos` and that every consumer with a stored position has
     * passed (admin token); `deleted` is the number removed. A
     * `minAgeNanos` of `0` means no minimum age and is valid.
     */
    suspend fun runRetention(request: RunRetentionRequest): RunRetentionResponse =
        callAsync { transport.runRetention(request, headers()) }

    /**
     * Reads a bounded range of stored changes of a source (reader or admin
     * token) — the same filter and range semantics as
     * [io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient.readChanges]:
     * an unset `from`/`to`/`limit` carries `0` (not set).
     */
    suspend fun readChanges(request: ReadChangesRequest): ReadChangesResponse =
        callAsync { transport.readChanges(request, headers()) }

    /**
     * Reads the operational diagnose report of a source (reader or admin
     * token) — the same report as the CLI diagnose mode and `GET /diagnose`.
     * A `known`/`present`/`*Known` field of `false` carries the respective
     * absence case (no heartbeat ever written, no blocking consumer,
     * unknown estimate/backlog).
     */
    suspend fun diagnose(request: DiagnoseRequest): DiagnoseResponse =
        callAsync { transport.diagnose(request, headers()) }

    /** Shuts down the channel this instance owns, if any (see the convenience constructor). */
    override fun close() {
        ownedChannel?.shutdownNow()
    }

    private fun headers(): Metadata = Metadata().apply {
        put(AUTHORIZATION_METADATA_ENTRY, BEARER_PREFIX + options.apiToken)
    }

    private suspend fun <T> callAsync(call: suspend () -> T): T =
        try {
            call()
        } catch (e: StatusException) {
            throw mapException(e)
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

        fun mapException(e: StatusException): PgChangeFeedGrpcException {
            val message = e.status.description ?: ""
            return when (e.status.code) {
                Status.Code.INVALID_ARGUMENT -> PgChangeFeedGrpcInvalidArgumentException(message, e)
                Status.Code.UNAUTHENTICATED -> PgChangeFeedGrpcUnauthenticatedException(message, e)
                Status.Code.PERMISSION_DENIED -> PgChangeFeedGrpcPermissionDeniedException(message, e)
                Status.Code.NOT_FOUND -> PgChangeFeedGrpcNotFoundException(message, e)
                Status.Code.INTERNAL -> PgChangeFeedGrpcInternalException(message, e)
                else -> PgChangeFeedGrpcUnexpectedStatusException(e.status.code, message, e)
            }
        }
    }
}

/**
 * Transport interface between [PgChangeFeedAdministrationClient] and the
 * generated coroutine stub (`AdministrationGrpcKt.AdministrationCoroutineStub`)
 * — one suspend function per RPC, the unary-call counterpart of
 * [GrpcStreamTransport]; see [PgChangeFeedGrpcClient]'s test-only constructor
 * KDoc for why the generated stub itself is not fakeable directly.
 */
internal interface AdministrationTransport {
    suspend fun registerConsumer(request: RegisterConsumerRequest, headers: Metadata): RegisterConsumerResponse
    suspend fun acknowledgeConsumer(request: AcknowledgeConsumerRequest, headers: Metadata): AcknowledgeConsumerResponse
    suspend fun getConsumerPosition(request: GetConsumerPositionRequest, headers: Metadata): GetConsumerPositionResponse
    suspend fun removeConsumer(request: RemoveConsumerRequest, headers: Metadata): RemoveConsumerResponse
    suspend fun enableTable(request: EnableTableRequest, headers: Metadata): EnableTableResponse
    suspend fun disableTable(request: DisableTableRequest, headers: Metadata): DisableTableResponse
    suspend fun getTableStatus(request: GetTableStatusRequest, headers: Metadata): GetTableStatusResponse
    suspend fun listTables(request: ListTablesRequest, headers: Metadata): ListTablesResponse
    suspend fun runRetention(request: RunRetentionRequest, headers: Metadata): RunRetentionResponse
    suspend fun readChanges(request: ReadChangesRequest, headers: Metadata): ReadChangesResponse
    suspend fun diagnose(request: DiagnoseRequest, headers: Metadata): DiagnoseResponse
}

/** The real [AdministrationTransport]: wraps the generated coroutine stub. */
internal class GeneratedAdministrationStubTransport(channel: Channel) : AdministrationTransport {
    private val stub = AdministrationGrpcKt.AdministrationCoroutineStub(channel)

    override suspend fun registerConsumer(request: RegisterConsumerRequest, headers: Metadata) =
        stub.registerConsumer(request, headers)

    override suspend fun acknowledgeConsumer(request: AcknowledgeConsumerRequest, headers: Metadata) =
        stub.acknowledgeConsumer(request, headers)

    override suspend fun getConsumerPosition(request: GetConsumerPositionRequest, headers: Metadata) =
        stub.getConsumerPosition(request, headers)

    override suspend fun removeConsumer(request: RemoveConsumerRequest, headers: Metadata) =
        stub.removeConsumer(request, headers)

    override suspend fun enableTable(request: EnableTableRequest, headers: Metadata) =
        stub.enableTable(request, headers)

    override suspend fun disableTable(request: DisableTableRequest, headers: Metadata) =
        stub.disableTable(request, headers)

    override suspend fun getTableStatus(request: GetTableStatusRequest, headers: Metadata) =
        stub.getTableStatus(request, headers)

    override suspend fun listTables(request: ListTablesRequest, headers: Metadata) =
        stub.listTables(request, headers)

    override suspend fun runRetention(request: RunRetentionRequest, headers: Metadata) =
        stub.runRetention(request, headers)

    override suspend fun readChanges(request: ReadChangesRequest, headers: Metadata) =
        stub.readChanges(request, headers)

    override suspend fun diagnose(request: DiagnoseRequest, headers: Metadata) =
        stub.diagnose(request, headers)
}
