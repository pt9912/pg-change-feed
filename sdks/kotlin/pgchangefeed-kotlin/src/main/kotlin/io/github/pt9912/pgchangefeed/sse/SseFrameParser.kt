package io.github.pt9912.pgchangefeed.sse

/**
 * One complete SSE frame of the change stream — the `event:` name and the
 * `data:` payload accumulated up to the blank line that terminates it.
 */
internal data class SseFrame(val name: String, val data: String)

/**
 * Splits Server-Sent-Events frames from a line source into [SseFrame] values —
 * plain line processing over the JDK standard library, no further dependency.
 *
 * A frame ends with the blank line the server writes after the `event:` and
 * the `data:` line; it separates two consecutive events. When the source is
 * exhausted before a frame is complete ([next] returns `null`), [readFrame]
 * also returns `null`, and a frame already begun is discarded — an incomplete
 * frame is not an event.
 */
internal object SseFrameParser {
    private const val EVENT_PREFIX = "event: "
    private const val DATA_PREFIX = "data: "

    fun readFrame(next: () -> String?): SseFrame? {
        var name: String? = null
        var data: String? = null
        while (true) {
            val line = next() ?: return null
            if (line.isEmpty()) {
                if (name == null && data == null) {
                    continue
                }
                return SseFrame(name ?: "", data ?: "")
            }
            when {
                line.startsWith(EVENT_PREFIX) -> name = line.removePrefix(EVENT_PREFIX)
                line.startsWith(DATA_PREFIX) -> data = line.removePrefix(DATA_PREFIX)
            }
        }
    }
}
