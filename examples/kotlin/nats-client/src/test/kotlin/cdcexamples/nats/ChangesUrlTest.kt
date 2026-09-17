package cdcexamples.nats

import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * Prüft den Aufbau der HTTP-Abfrage des zweiseitigen Ablaufs
 * (`LH-FA-SST-006`, `ADR-0079`) — netzlos, reine Funktion. Form-Vorbild:
 * `examples/nats-client/subject_test.go` (`TestChangesURL...`),
 * `examples/csharp/nats-client/NatsClient.Tests/ChangesUrlBuilderTests.cs`.
 */
class ChangesUrlTest {
    @Test
    fun buildCarriesAllThreeFiltersOnReadPath() {
        val got = ChangesUrl.build("feed:8080", "quelle-1", "public", "orders")
        assertEquals("http://feed:8080/changes?schema=public&source=quelle-1&table=orders", got)
    }

    @Test
    fun buildKeepsHostBoundary() {
        val got = ChangesUrl.build("localhost:9090", "src-e2e", "public", "feed_e2e")
        assertEquals("http://localhost:9090/changes?schema=public&source=src-e2e&table=feed_e2e", got)
    }

    @Test
    fun buildEscapesReservedCharacters() {
        val got = ChangesUrl.build("feed:8080", "quelle & test", "sch/1", "tab 1")
        assertEquals("http://feed:8080/changes?schema=sch%2F1&source=quelle%20%26%20test&table=tab%201", got)
    }
}
