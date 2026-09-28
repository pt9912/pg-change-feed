package cdcexamples.http

import java.net.http.HttpClient
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertTrue

/**
 * Prüft die drei Tabellen-Verwaltungs-Fähigkeiten — netzlos, über
 * [FakeServer]. Form-Vorbild: `examples/http-client/tables_admin_test.go`.
 */
class TablesAdminClientTest {
    private fun baseConfig(addr: String) = Config(
        addr = addr, token = "reader-token", adminToken = "admin-token", verb = "",
        source = "", publication = "", consumerId = "", name = "", offset = 0,
        schema = "", table = "", tableId = "", schemaVersionId = "", version = 1,
        from = "", to = "", limit = "", minAgeNanos = 0,
    )

    @Test
    fun tableStatusUrlBuilderCarriesAllFourFields() {
        val got = TableStatusUrlBuilder.build("feed:8080", "quelle-1", "public", "orders", "pub_quelle_1")
        assertEquals("http://feed:8080/tables/status?publication=pub_quelle_1&schema=public&source=quelle-1&table=orders", got)
    }

    @Test
    fun enableTableParsesSuccessResponse() {
        FakeServer.start().use { server ->
            server.respondWith(
                201,
                """{"table_id":"public.orders","source":"quelle-1","schema":"public","table":"orders","already_enabled":false}""",
            )
            val httpClient = HttpClient.newHttpClient()
            val cfg = baseConfig(server.addr).copy(
                source = "quelle-1", schema = "public", table = "orders",
                tableId = "public.orders", schemaVersionId = "public.orders-v1", version = 1,
                publication = "pub_quelle_1",
            )

            val resp = TablesAdminClient.enableTable(httpClient, cfg)

            assertFalse(resp.alreadyEnabled)
            assertEquals("/tables/enable", server.lastPath)
            assertTrue(server.lastBody!!.contains("\"table_id\":\"public.orders\""))
            assertTrue(server.lastBody!!.contains("\"schema_version_id\":\"public.orders-v1\""))
            assertTrue(server.lastBody!!.contains("\"version\":1"))
        }
    }

    @Test
    fun enableTableFailsOnNon2xx() {
        FakeServer.start().use { server ->
            server.respondWith(404, """{"error":"Tabelle existiert an der Quelle nicht"}""")
            val httpClient = HttpClient.newHttpClient()

            val ex = assertFailsWith<IllegalStateException> {
                TablesAdminClient.enableTable(httpClient, baseConfig(server.addr))
            }
            assertTrue(ex.message!!.contains("404"))
        }
    }

    @Test
    fun disableTableParsesBothOutcomeFields() {
        FakeServer.start().use { server ->
            server.respondWith(200, """{"removed":true,"retained":true}""")
            val httpClient = HttpClient.newHttpClient()
            val cfg = baseConfig(server.addr).copy(source = "quelle-1", schema = "public", table = "orders", publication = "pub_quelle_1")

            val resp = TablesAdminClient.disableTable(httpClient, cfg)

            assertTrue(resp.removed)
            assertTrue(resp.retained)
            assertEquals("/tables/disable", server.lastPath)
            assertTrue(server.lastBody!!.contains("\"publication\":\"pub_quelle_1\""))
        }
    }

    @Test
    fun disableTableFailsOnNon2xx() {
        FakeServer.start().use { server ->
            server.respondWith(400, """{"error":"source, schema, table und publication sind Pflichtfelder"}""")
            val httpClient = HttpClient.newHttpClient()

            val ex = assertFailsWith<IllegalStateException> {
                TablesAdminClient.disableTable(httpClient, baseConfig(server.addr))
            }
            assertTrue(ex.message!!.contains("400"))
            assertTrue(ex.message!!.contains("Pflichtfelder"))
        }
    }

    @Test
    fun tableStatusParsesNeverEnabledBoundary() {
        FakeServer.start().use { server ->
            server.respondWith(200, """{"enabled":false,"retained":false}""")
            val httpClient = HttpClient.newHttpClient()
            val cfg = baseConfig(server.addr).copy(source = "quelle-1", schema = "public", table = "nie_aktiviert", publication = "pub_quelle_1")

            val resp = TablesAdminClient.tableStatus(httpClient, cfg)

            assertFalse(resp.enabled)
            assertFalse(resp.retained)
            assertEquals("GET", server.lastMethod)
        }
    }

    @Test
    fun tableStatusFailsOnNon2xx() {
        FakeServer.start().use { server ->
            server.respondWith(400, """{"error":"source, schema, table und publication sind Pflichtfelder"}""")
            val httpClient = HttpClient.newHttpClient()

            val ex = assertFailsWith<IllegalStateException> {
                TablesAdminClient.tableStatus(httpClient, baseConfig(server.addr))
            }
            assertTrue(ex.message!!.contains("400"))
            assertTrue(ex.message!!.contains("Pflichtfelder"))
        }
    }
}
