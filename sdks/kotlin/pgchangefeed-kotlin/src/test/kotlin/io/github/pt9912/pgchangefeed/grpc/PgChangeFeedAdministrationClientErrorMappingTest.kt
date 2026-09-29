package io.github.pt9912.pgchangefeed.grpc

import cdc.administration.v1.AdministrationOuterClass.ListTablesRequest
import io.grpc.Status
import kotlinx.coroutines.runBlocking
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

/**
 * The gRPC-status-to-exception mapping shared by every RPC of
 * [PgChangeFeedAdministrationClient] — exercised through `listTables` as the
 * representative call, since the mapping itself does not depend on which RPC
 * failed.
 */
class PgChangeFeedAdministrationClientErrorMappingTest {
    // Mutation that turns every test of this file red (checked for real):
    // `mapException` changed to always return `PgChangeFeedGrpcInternalException`
    // regardless of the gRPC status — each assertion on the specific typed
    // exception below fails, instead of staying silently green.
    private suspend fun call(transport: FakeAdministrationTransport) =
        AdministrationTestClientFactory.create(transport)
            .listTables(ListTablesRequest.newBuilder().setSource("src").setPublication("pub").build())

    @Test
    fun `invalidArgument throws typed invalidArgument exception`() = runBlocking {
        val transport = FakeAdministrationTransport.withStatus(
            Status.INVALID_ARGUMENT.withDescription("source must not be empty"),
        )

        val ex = assertFailsWith<PgChangeFeedGrpcInvalidArgumentException> { call(transport) }

        assertEquals(Status.Code.INVALID_ARGUMENT, ex.statusCode)
        assertEquals("source must not be empty", ex.message)
    }

    @Test
    fun `unauthenticated throws typed unauthenticated exception`() = runBlocking {
        val transport = FakeAdministrationTransport.withStatus(
            Status.UNAUTHENTICATED.withDescription("missing or unknown bearer token"),
        )

        val ex = assertFailsWith<PgChangeFeedGrpcUnauthenticatedException> { call(transport) }

        assertEquals(Status.Code.UNAUTHENTICATED, ex.statusCode)
    }

    @Test
    fun `permissionDenied throws typed permissionDenied exception`() = runBlocking {
        val transport = FakeAdministrationTransport.withStatus(
            Status.PERMISSION_DENIED.withDescription("insufficient role for this rpc"),
        )

        val ex = assertFailsWith<PgChangeFeedGrpcPermissionDeniedException> { call(transport) }

        assertEquals(Status.Code.PERMISSION_DENIED, ex.statusCode)
    }

    @Test
    fun `notFound throws typed notFound exception`() = runBlocking {
        val transport = FakeAdministrationTransport.withStatus(
            Status.NOT_FOUND.withDescription("table does not exist at source"),
        )

        val ex = assertFailsWith<PgChangeFeedGrpcNotFoundException> { call(transport) }

        assertEquals(Status.Code.NOT_FOUND, ex.statusCode)
    }

    @Test
    fun `internal throws typed internal exception`() = runBlocking {
        val transport = FakeAdministrationTransport.withStatus(
            Status.INTERNAL.withDescription("internal error"),
        )

        val ex = assertFailsWith<PgChangeFeedGrpcInternalException> { call(transport) }

        assertEquals(Status.Code.INTERNAL, ex.statusCode)
    }

    @Test
    fun `status code outside documented set throws unexpectedStatus`() = runBlocking {
        val transport = FakeAdministrationTransport.withStatus(
            Status.UNAVAILABLE.withDescription("server unreachable"),
        )

        val ex = assertFailsWith<PgChangeFeedGrpcUnexpectedStatusException> { call(transport) }

        assertEquals(Status.Code.UNAVAILABLE, ex.statusCode)
        assertEquals("server unreachable", ex.message)
    }
}
