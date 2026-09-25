package io.github.pt9912.pgchangefeed.http

import io.github.pt9912.pgchangefeed.http.model.DisableTableRequest
import io.github.pt9912.pgchangefeed.http.model.EnableTableRequest
import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * Happy-path coverage for the four table-management capabilities
 * (`EnableTable`, `DisableTable`, `GetStatus`, `ListTables`).
 */
class PgChangeFeedHttpClientTableTest {

    @Test
    fun `enableTable posts to tables enable and parses the response`() {
        val (client, transport) = TestClientFactory.create { request ->
            assertEquals("POST", request.method)
            assertEquals("http://example.invalid:8080/tables/enable", request.url)
            FakeHttpTransport.jsonResponse(
                201,
                """{"table_id":"t-1","source":"src","schema":"public","table":"orders","already_enabled":false}""",
            )(request)
        }

        val response = client.enableTable(
            EnableTableRequest(
                source = "src",
                schema = "public",
                table = "orders",
                tableId = "t-1",
                schemaVersionId = "sv-1",
                version = 1L,
                publication = "pub",
            ),
        )

        assertEquals("t-1", response.tableId)
        assertEquals(false, response.alreadyEnabled)
        assertEquals(true, transport.lastRequest?.body?.contains("\"version\":1") == true)
    }

    @Test
    fun `disableTable posts to tables disable and parses the response`() {
        val (client, _) = TestClientFactory.create { request ->
            assertEquals("POST", request.method)
            assertEquals("http://example.invalid:8080/tables/disable", request.url)
            FakeHttpTransport.jsonResponse(200, """{"removed":true,"retained":false}""")(request)
        }

        val response = client.disableTable(DisableTableRequest("src", "public", "orders", "pub"))

        assertEquals(true, response.removed)
        assertEquals(false, response.retained)
    }

    @Test
    fun `getStatus gets tables status with all four query parameters`() {
        val (client, transport) = TestClientFactory.create { request ->
            assertEquals("GET", request.method)
            assertEquals(
                "http://example.invalid:8080/tables/status?source=src&schema=public&table=orders&publication=pub",
                request.url,
            )
            FakeHttpTransport.jsonResponse(200, """{"enabled":true,"retained":false}""")(request)
        }

        val response = client.getStatus("src", "public", "orders", "pub")

        assertEquals(true, response.enabled)
        assertEquals(false, response.retained)
        assertEquals(true, transport.lastRequest != null)
    }

    @Test
    fun `listTables gets tables and parses both lists`() {
        val (client, _) = TestClientFactory.create { request ->
            assertEquals("GET", request.method)
            assertEquals("http://example.invalid:8080/tables?source=src&publication=pub", request.url)
            FakeHttpTransport.jsonResponse(
                200,
                """{"tables":[{"table_id":"t-1","source":"src","schema":"public","table":"orders"}],"retained":[]}""",
            )(request)
        }

        val response = client.listTables("src", "pub")

        assertEquals(1, response.tables.size)
        assertEquals("orders", response.tables[0].table)
        assertEquals(0, response.retained.size)
    }

    @Test
    fun `listTables returns both lists empty without any activation`() {
        val (client, _) = TestClientFactory.create { request ->
            FakeHttpTransport.jsonResponse(200, """{"tables":[],"retained":[]}""")(request)
        }

        val response = client.listTables("src", "pub")

        assertEquals(0, response.tables.size)
        assertEquals(0, response.retained.size)
    }

    @Test
    fun `listTables percent-encodes reserved characters in query parameters`() {
        val (client, _) = TestClientFactory.create { request ->
            assertEquals(
                "http://example.invalid:8080/tables?source=quelle%20%26%20test&publication=pub%2F1",
                request.url,
            )
            FakeHttpTransport.jsonResponse(200, """{"tables":[],"retained":[]}""")(request)
        }

        client.listTables("quelle & test", "pub/1")
    }
}
