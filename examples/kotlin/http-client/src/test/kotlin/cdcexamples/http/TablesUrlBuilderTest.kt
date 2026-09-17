package cdcexamples.http

import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * Prüft den Aufbau der HTTP-Abfrage (`LH-FA-SST-006`) — netzlos, reine
 * Funktion. Form-Vorbild: `examples/http-client/tables_test.go`,
 * `examples/csharp/http-client/HttpClient.Tests/TablesUrlBuilderTests.cs`.
 */
class TablesUrlBuilderTest {
    @Test
    fun buildCarriesBothRequiredFilters() {
        val got = TablesUrlBuilder.build("feed:8080", "quelle-1", "pub_quelle_1")
        assertEquals("http://feed:8080/tables?publication=pub_quelle_1&source=quelle-1", got)
    }

    @Test
    fun buildKeepsHostBoundary() {
        val got = TablesUrlBuilder.build("localhost:9090", "src-e2e", "cdc_pub")
        assertEquals("http://localhost:9090/tables?publication=cdc_pub&source=src-e2e", got)
    }

    @Test
    fun buildEscapesReservedCharacters() {
        val got = TablesUrlBuilder.build("feed:8080", "quelle & test", "pub/1")
        assertEquals("http://feed:8080/tables?publication=pub%2F1&source=quelle%20%26%20test", got)
    }
}
