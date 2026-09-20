package io.github.pt9912.pgchangefeed.http

import io.github.pt9912.pgchangefeed.http.model.RunRetentionRequest
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue

/**
 * Happy-path coverage for `RunRetention` (SPEC-018) and `ReadChanges`
 * (SPEC-022, `GET /changes`), including the optional-parameter handling
 * and the embedded-JSON `old_image`/`new_image` shape.
 */
class PgChangeFeedHttpClientRetentionAndChangesTest {

    @Test
    fun `runRetention posts to retention run and parses the response`() {
        val (client, transport) = TestClientFactory.create { request ->
            assertEquals("POST", request.method)
            assertEquals("http://example.invalid:8080/retention/run", request.url)
            FakeHttpTransport.jsonResponse(200, """{"deleted":12}""")(request)
        }

        val response = client.runRetention(RunRetentionRequest("src", 0L))

        assertEquals(12, response.deleted)
        assertEquals(true, transport.lastRequest?.body?.contains("\"min_age_nanos\":0") == true)
    }

    @Test
    fun `readChanges gets changes with only the mandatory source parameter`() {
        val (client, _) = TestClientFactory.create { request ->
            assertEquals("GET", request.method)
            assertEquals("http://example.invalid:8080/changes?source=src", request.url)
            FakeHttpTransport.jsonResponse(200, """{"changes":[]}""")(request)
        }

        val response = client.readChanges(source = "src")

        assertEquals(0, response.changes.size)
    }

    @Test
    fun `readChanges gets changes with every optional filter set`() {
        val (client, _) = TestClientFactory.create { request ->
            assertEquals(
                "http://example.invalid:8080/changes?source=src&schema=public&table=orders&from=1&to=10&limit=5",
                request.url,
            )
            FakeHttpTransport.jsonResponse(200, """{"changes":[]}""")(request)
        }

        client.readChanges(source = "src", schema = "public", table = "orders", from = 1L, to = 10L, limit = 5)
    }

    @Test
    fun `readChanges parses a change with a populated new_image and a JSON-null old_image`() {
        val (client, _) = TestClientFactory.create { request ->
            FakeHttpTransport.jsonResponse(
                200,
                """{"changes":[{"commit_position":1,"change_id":"ch-1","transaction_id":"tx-1",""" +
                    """"source_table_id":"t-1","schema":"public","table":"orders","sequence":1,""" +
                    """"operation":"INSERT","old_image":null,"new_image":{"id":1,"name":"a"},""" +
                    """"schema_version":"sv-1","committed_at":"2026-09-20T00:00:00Z"}]}""",
            )(request)
        }

        val response = client.readChanges(source = "src")

        val change = response.changes.single()
        assertEquals("ch-1", change.changeId)
        assertEquals("INSERT", change.operation)
        assertTrue(change.oldImage?.isJsonNull == true)
        assertEquals("a", change.newImage?.asJsonObject?.get("name")?.asString)
    }

    @Test
    fun `readChanges returns an empty list instead of a 404 on no match`() {
        val (client, _) = TestClientFactory.create { request ->
            FakeHttpTransport.jsonResponse(200, """{"changes":[]}""")(request)
        }

        val response = client.readChanges(source = "src", from = 999L)

        assertEquals(0, response.changes.size)
    }
}
