package cdcexamples.sse

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNull

/**
 * Prüft den Aufbau der Stream-Adresse und das Zerlegen eines SSE-Frames
 * (`LH-FA-SST-008`, `ADR-0061`) — netzlos, reine Funktionen. Form-Vorbild:
 * `examples/sse-client/stream_test.go`,
 * `examples/csharp/sse-client/SseClient.Tests/SseStreamTests.cs`.
 */
class SseStreamTest {
    /**
     * lines liefert eine Zeilenquelle über einen Ausschnitt; jeder Aufruf
     * gibt die nächste Zeile zurück, nach der letzten `null` — dieselbe Form,
     * die `main` aus dem `Stream<String>`-Iterator über den Response-Body
     * bildet.
     */
    private fun lines(vararg all: String): () -> String? {
        var i = 0
        return { if (i < all.size) all[i++] else null }
    }

    @Test
    fun streamUrlBuildsChangesStreamEndpoint() {
        val got = SseStream.streamUrl("feed:8080")
        assertEquals("http://feed:8080/changes/stream", got)
    }

    @Test
    fun readEventReadsNameAndPayload() {
        val ev = SseStream.readEvent(
            lines(
                "event: change",
                """data: {"change_id":"c-1","table":"orders"}""",
                "",
            ),
        )
        assertEquals(Event("change", """{"change_id":"c-1","table":"orders"}"""), ev)
    }

    @Test
    fun readEventStopsAtFrameBoundary() {
        val next = lines(
            "event: change",
            """data: {"change_id":"c-1"}""",
            "",
            "event: change",
            """data: {"change_id":"c-2"}""",
            "",
        )
        val first = SseStream.readEvent(next)
        val second = SseStream.readEvent(next)
        assertEquals(Event("change", """{"change_id":"c-1"}"""), first)
        assertEquals(Event("change", """{"change_id":"c-2"}"""), second)
    }

    @Test
    fun readEventReportsExhaustedSource() {
        val next = lines(
            "event: change",
            """data: {"change_id":"c-1"}""",
        )
        val ev = SseStream.readEvent(next)
        assertNull(ev)
    }
}
