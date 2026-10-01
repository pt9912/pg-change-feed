package cdcexamples.http

import java.net.http.HttpClient
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

/**
 * Prüft `GET /changes` — netzlos, über [FakeServer]. Form-Vorbild:
 * `examples/http-client/changes_test.go`.
 */
class ChangesClientTest {
    private fun baseConfig(addr: String) = Config(
        addr = addr, token = "reader-token", adminToken = "", verb = "",
        source = "", publication = "", consumerId = "", name = "", offset = 0,
        schema = "", table = "", tableId = "", schemaVersionId = "", version = 1,
        from = "", to = "", limit = "", minAgeNanos = 0,
    )

    @Test
    fun changesUrlBuilderCarriesOnlyRequiredField() {
        val got = ChangesUrlBuilder.build("feed:8080", "quelle-1", "", "", "", "", "")
        assertEquals("http://feed:8080/changes?source=quelle-1", got)
    }

    @Test
    fun changesUrlBuilderCarriesAllOptionalFieldsIndependently() {
        val got = ChangesUrlBuilder.build("feed:8080", "quelle-1", "public", "orders", "10", "20", "5")
        assertEquals("http://feed:8080/changes?from=10&limit=5&schema=public&source=quelle-1&table=orders&to=20", got)
    }

    @Test
    fun changesUrlBuilderCarriesTheTargetAsEscapedQueryParameter() {
        assertEquals(
            "http://feed:8080/changes?source=quelle-1&target=eu",
            ChangesUrlBuilder.build("feed:8080", "quelle-1", "", "", "", "", "", "eu"),
        )
        assertEquals(
            "http://feed:8080/changes?source=quelle-1&target=a%26b%3Dc",
            ChangesUrlBuilder.build("feed:8080", "quelle-1", "", "", "", "", "", "a&b=c"),
        )
        assertEquals(
            "http://feed:8080/changes?source=quelle-1",
            ChangesUrlBuilder.build("feed:8080", "quelle-1", "", "", "", "", "", ""),
        )
    }

    @Test
    fun readChangesSendsTheTargetOfTheConfig() {
        FakeServer.start().use { server ->
            server.respondWith(200, """{"changes":[]}""")
            val cfg = baseConfig(server.addr).copy(source = "quelle-1", target = "eu")

            ChangesClient.readChanges(HttpClient.newHttpClient(), cfg)

            assertEquals("source=quelle-1&target=eu", server.lastQuery)
        }
    }

    @Test
    fun changesUrlBuilderOmitsEmptyOptionalValues() {
        val got = ChangesUrlBuilder.build("feed:8080", "quelle-1", "", "orders", "", "20", "")
        assertEquals("http://feed:8080/changes?source=quelle-1&table=orders&to=20", got)
    }

    @Test
    fun readChangesParsesSuccessResponseWithEmbeddedImages() {
        FakeServer.start().use { server ->
            server.respondWith(
                200,
                """{"changes":[{"commit_position":42,"change_id":"change-1","transaction_id":"tx-1","source_table_id":"t-1","schema":"public","table":"orders","sequence":0,"operation":"INSERT","old_image":null,"new_image":{"name":"erste Zeile"},"schema_version":"sv-1","committed_at":"2026-09-19T00:00:00Z","origin":"wal"}]}""",
            )
            val httpClient = HttpClient.newHttpClient()
            val cfg = baseConfig(server.addr).copy(source = "quelle-1")

            val resp = ChangesClient.readChanges(httpClient, cfg)

            assertEquals(1, resp.changes.size)
            val change = resp.changes[0]
            assertEquals("change-1", change.changeId)
            assertEquals("wal", change.origin)
            assertTrue(change.oldImage.isJsonNull)
            assertEquals("erste Zeile", change.newImage.asJsonObject.get("name").asString)
            assertEquals("source=quelle-1", server.lastQuery)
        }
    }

    @Test
    fun readChangesFailsOnNon2xx() {
        FakeServer.start().use { server ->
            server.respondWith(400, """{"error":"source ist Pflichtfeld"}""")
            val httpClient = HttpClient.newHttpClient()

            val ex = assertFailsWith<IllegalStateException> {
                ChangesClient.readChanges(httpClient, baseConfig(server.addr))
            }
            assertTrue(ex.message!!.contains("400"))
        }
    }
}
