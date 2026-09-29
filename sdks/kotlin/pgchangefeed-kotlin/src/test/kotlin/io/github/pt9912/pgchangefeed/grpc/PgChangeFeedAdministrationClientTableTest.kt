package io.github.pt9912.pgchangefeed.grpc

import cdc.administration.v1.AdministrationOuterClass.DisableTableRequest
import cdc.administration.v1.AdministrationOuterClass.DisableTableResponse
import cdc.administration.v1.AdministrationOuterClass.EnableTableRequest
import cdc.administration.v1.AdministrationOuterClass.EnableTableResponse
import cdc.administration.v1.AdministrationOuterClass.GetTableStatusRequest
import cdc.administration.v1.AdministrationOuterClass.GetTableStatusResponse
import cdc.administration.v1.AdministrationOuterClass.ListTablesRequest
import cdc.administration.v1.AdministrationOuterClass.ListTablesResponse
import cdc.administration.v1.AdministrationOuterClass.SourceTable
import io.grpc.Status
import kotlinx.coroutines.runBlocking
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertTrue

/**
 * Happy-path (and the table-only `NotFound`) coverage for the four
 * table-management RPCs (`enableTable`, `disableTable`, `getTableStatus`,
 * `listTables`).
 */
class PgChangeFeedAdministrationClientTableTest {
    private fun sampleEnableRequest(): EnableTableRequest = EnableTableRequest.newBuilder()
        .setSource("src").setSchema("public").setTable("orders")
        .setTableId("t-1").setSchemaVersionId("sv-1").setVersion(1).setPublication("pub")
        .build()

    @Test
    fun `enableTable happy path returns typed response`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            EnableTableResponse.newBuilder()
                .setTableId("t-1").setSource("src").setSchema("public").setTable("orders").setAlreadyEnabled(false)
                .build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.enableTable(sampleEnableRequest())

        assertEquals("t-1", response.tableId)
        assertFalse(response.alreadyEnabled)
    }

    @Test
    fun `enableTable table missing at source throws notFound`() = runBlocking {
        val transport = FakeAdministrationTransport.withStatus(
            Status.NOT_FOUND.withDescription("table does not exist at source"),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val ex = assertFailsWith<PgChangeFeedGrpcNotFoundException> {
            client.enableTable(sampleEnableRequest())
        }

        assertEquals(Status.Code.NOT_FOUND, ex.statusCode)
        assertEquals("table does not exist at source", ex.message)
    }

    @Test
    fun `disableTable happy path returns typed response`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            DisableTableResponse.newBuilder().setRemoved(true).setRetained(false).build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.disableTable(
            DisableTableRequest.newBuilder().setSource("src").setSchema("public").setTable("orders").setPublication("pub").build(),
        )

        assertTrue(response.removed)
        assertFalse(response.retained)
    }

    @Test
    fun `disableTable table missing at source throws notFound`() = runBlocking {
        val transport = FakeAdministrationTransport.withStatus(
            Status.NOT_FOUND.withDescription("table does not exist at source"),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val ex = assertFailsWith<PgChangeFeedGrpcNotFoundException> {
            client.disableTable(
                DisableTableRequest.newBuilder()
                    .setSource("src").setSchema("public").setTable("orders").setPublication("pub").build(),
            )
        }

        assertEquals(Status.Code.NOT_FOUND, ex.statusCode)
        assertEquals("table does not exist at source", ex.message)
    }

    @Test
    fun `getTableStatus happy path returns typed response`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            GetTableStatusResponse.newBuilder().setEnabled(true).setRetained(false).build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.getTableStatus(
            GetTableStatusRequest.newBuilder()
                .setSource("src").setSchema("public").setTable("orders").setPublication("pub").build(),
        )

        assertTrue(response.enabled)
        assertFalse(response.retained)
    }

    @Test
    fun `getTableStatus never enabled reports both false`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            GetTableStatusResponse.newBuilder().setEnabled(false).setRetained(false).build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.getTableStatus(
            GetTableStatusRequest.newBuilder()
                .setSource("src").setSchema("public").setTable("orders").setPublication("pub").build(),
        )

        assertFalse(response.enabled)
        assertFalse(response.retained)
    }

    @Test
    fun `listTables happy path returns typed response`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            ListTablesResponse.newBuilder()
                .addTables(SourceTable.newBuilder().setTableId("t-1").setSource("src").setSchema("public").setTable("orders").build())
                .build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.listTables(ListTablesRequest.newBuilder().setSource("src").setPublication("pub").build())

        assertEquals(1, response.tablesCount)
        assertEquals("t-1", response.tablesList[0].tableId)
        assertEquals(0, response.retainedCount)
    }
}
