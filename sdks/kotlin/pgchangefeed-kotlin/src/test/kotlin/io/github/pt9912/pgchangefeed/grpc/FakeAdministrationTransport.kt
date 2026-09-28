package io.github.pt9912.pgchangefeed.grpc

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
import io.grpc.Metadata
import io.grpc.Status
import io.grpc.StatusException

/**
 * A network-free [AdministrationTransport]: whichever of the eleven RPCs a
 * test calls returns [response] when set, or throws a [StatusException] of
 * [failureStatus] when set instead — the mapping to a typed
 * [PgChangeFeedGrpcException] is [PgChangeFeedAdministrationClient]'s
 * responsibility, not this fake's. [lastRequest] and [lastHeaders] let a
 * test assert on the exact request/metadata a client call built, regardless
 * of which RPC made it — the unary-call counterpart of
 * [FakeGrpcStreamTransport].
 */
internal class FakeAdministrationTransport private constructor(
    private val response: Any?,
    private val failureStatus: Status?,
) : AdministrationTransport {

    var lastRequest: Any? = null
        private set
    var lastHeaders: Metadata? = null
        private set

    private fun record(request: Any, headers: Metadata) {
        lastRequest = request
        lastHeaders = headers
        failureStatus?.let { throw StatusException(it) }
    }

    override suspend fun registerConsumer(request: RegisterConsumerRequest, headers: Metadata): RegisterConsumerResponse {
        record(request, headers)
        return response as RegisterConsumerResponse
    }

    override suspend fun acknowledgeConsumer(
        request: AcknowledgeConsumerRequest,
        headers: Metadata,
    ): AcknowledgeConsumerResponse {
        record(request, headers)
        return response as AcknowledgeConsumerResponse
    }

    override suspend fun getConsumerPosition(
        request: GetConsumerPositionRequest,
        headers: Metadata,
    ): GetConsumerPositionResponse {
        record(request, headers)
        return response as GetConsumerPositionResponse
    }

    override suspend fun removeConsumer(request: RemoveConsumerRequest, headers: Metadata): RemoveConsumerResponse {
        record(request, headers)
        return response as RemoveConsumerResponse
    }

    override suspend fun enableTable(request: EnableTableRequest, headers: Metadata): EnableTableResponse {
        record(request, headers)
        return response as EnableTableResponse
    }

    override suspend fun disableTable(request: DisableTableRequest, headers: Metadata): DisableTableResponse {
        record(request, headers)
        return response as DisableTableResponse
    }

    override suspend fun getTableStatus(request: GetTableStatusRequest, headers: Metadata): GetTableStatusResponse {
        record(request, headers)
        return response as GetTableStatusResponse
    }

    override suspend fun listTables(request: ListTablesRequest, headers: Metadata): ListTablesResponse {
        record(request, headers)
        return response as ListTablesResponse
    }

    override suspend fun runRetention(request: RunRetentionRequest, headers: Metadata): RunRetentionResponse {
        record(request, headers)
        return response as RunRetentionResponse
    }

    override suspend fun readChanges(request: ReadChangesRequest, headers: Metadata): ReadChangesResponse {
        record(request, headers)
        return response as ReadChangesResponse
    }

    override suspend fun diagnose(request: DiagnoseRequest, headers: Metadata): DiagnoseResponse {
        record(request, headers)
        return response as DiagnoseResponse
    }

    companion object {
        /** Builds a fake transport whose call returns [response] — the happy path. */
        fun withResponse(response: Any): FakeAdministrationTransport = FakeAdministrationTransport(response, null)

        /**
         * Builds a fake transport whose call fails immediately with
         * [status], modeling the shape of a server rejection (a
         * [StatusException] from the suspend function itself) — the real
         * rejection is covered by the real-server integration test.
         */
        fun withStatus(status: Status): FakeAdministrationTransport = FakeAdministrationTransport(null, status)
    }
}
