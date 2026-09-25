package io.github.pt9912.pgchangefeed.sse

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNull

/**
 * [SseFrameParser.readFrame] against a fake line source — network-free:
 * every frame boundary/truncation edge case a real `text/event-stream` body
 * can present.
 */
class SseFrameParserTest {
    private fun linesOf(vararg lines: String): () -> String? {
        val iterator = lines.iterator()
        return { if (iterator.hasNext()) iterator.next() else null }
    }

    @Test
    fun `a complete frame yields its event name and data`() {
        val next = linesOf("event: change", """data: {"a":1}""", "")

        val frame = SseFrameParser.readFrame(next)

        assertEquals(SseFrame("change", """{"a":1}"""), frame)
    }

    @Test
    fun `two consecutive frames are read independently`() {
        val next = linesOf(
            "event: change",
            """data: {"a":1}""",
            "",
            "event: change",
            """data: {"a":2}""",
            "",
        )

        val first = SseFrameParser.readFrame(next)
        val second = SseFrameParser.readFrame(next)

        assertEquals(SseFrame("change", """{"a":1}"""), first)
        assertEquals(SseFrame("change", """{"a":2}"""), second)
    }

    @Test
    fun `leading blank lines before a frame are skipped`() {
        val next = linesOf("", "", "event: change", """data: {"a":1}""", "")

        val frame = SseFrameParser.readFrame(next)

        assertEquals(SseFrame("change", """{"a":1}"""), frame)
    }

    @Test
    fun `an exhausted source before any frame returns null`() {
        val next = linesOf()

        val frame = SseFrameParser.readFrame(next)

        assertNull(frame)
    }

    // Mutation that turns this test red (checked for real):
    // `SseFrameParser.readFrame` um
    // ein zusätzliches `return null` direkt vor dem regulären
    // `return SseFrame(...)` ergänzt, sodass ein vollständiges,
    // abgeschlossenes Frame trotz Leerzeile als `null` statt als Event
    // zurückkommt — dieser Test schlägt dann fehl (erhält `null` statt des
    // erwarteten Frames), statt still grün zu bleiben.
    @Test
    fun `an incomplete trailing frame without its terminating blank line is discarded`() {
        val next = linesOf("event: change", """data: {"a":1}""")

        val frame = SseFrameParser.readFrame(next)

        assertNull(frame)
    }

    @Test
    fun `a line outside the event and data prefixes is ignored`() {
        val next = linesOf(": keep-alive comment", "event: change", """data: {"a":1}""", "")

        val frame = SseFrameParser.readFrame(next)

        assertEquals(SseFrame("change", """{"a":1}"""), frame)
    }

    @Test
    fun `a frame with only a data line and no event name defaults the name to empty`() {
        val next = linesOf("""data: {"a":1}""", "")

        val frame = SseFrameParser.readFrame(next)

        assertEquals(SseFrame("", """{"a":1}"""), frame)
    }
}
