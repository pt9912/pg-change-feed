package io.github.pt9912.pgchangefeed.sse

/**
 * One complete SSE frame of the `SPEC-021` change stream — the `event:`
 * name and the `data:` payload accumulated up to the blank line that
 * terminates it.
 */
internal data class SseFrame(val name: String, val data: String)

/**
 * Zerlegt Server-Sent-Events-Frames aus einer Zeilenquelle in [SseFrame]-
 * Werte — reine Zeilen-Verarbeitung über die JDK-Standardbibliothek, kein
 * Fremdmodul nötig. Form-Vorbild (gelesen, nicht importiert — `ADR-0109`
 * Festlegung 3): `examples/kotlin/sse-client/src/main/kotlin/cdcexamples/sse/SseStream.kt`'s
 * `readEvent`, bereits real erprobt.
 *
 * Ein Frame endet mit der Leerzeile, die der Server nach der `event:`- und
 * der `data:`-Zeile schreibt (`SPEC-021`); sie trennt zwei aufeinander-
 * folgende Events. Ist die Quelle vor dem Frame-Abschluss erschöpft
 * ([next] liefert `null`), liefert [readFrame] ebenfalls `null`, und ein
 * begonnenes Frame wird verworfen — ein unvollständiges Frame ist kein
 * Event.
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
