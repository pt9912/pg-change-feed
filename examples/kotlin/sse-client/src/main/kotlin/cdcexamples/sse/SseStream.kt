package cdcexamples.sse

/**
 * Event trägt ein Frame des Change-Streams (`LH-FA-SST-008`): den Namen aus
 * der `event:`-Zeile und die Nutzlast aus der `data:`-Zeile. Der Server
 * schreibt je Change ein JSON-Objekt mit den zehn Nachrichtenfeldern in die
 * `data:`-Zeile (`ADR-0061`).
 */
data class Event(val name: String, val data: String)

/**
 * SseStream baut die Stream-Adresse und zerlegt die SSE-Frames — reine
 * Funktionen, netzlos testbar. Form-Vorbild:
 * `examples/sse-client/stream.go`, `examples/csharp/sse-client/SseStream.cs`.
 * **Keine neue Abhängigkeit** (`ADR-0090` Festlegung 1): das Zerlegen eines
 * Frames ist Zeilen-Verarbeitung über die JDK-Standardbibliothek, kein
 * Fremdmodul (etwa `okhttp-sse`) nötig.
 */
object SseStream {
    /**
     * streamUrl baut die Adresse des SSE-Endpunkts. Ein gesetztes [target]
     * erscheint als maskierter Query-Parameter und wählt das Zustellziel einer
     * Change (`ADR-0137`); ein leeres lässt die Adresse ohne Query.
     */
    fun streamUrl(addr: String, target: String = ""): String =
        if (target.isEmpty()) "http://$addr/changes/stream" else "http://$addr/changes/stream?target=${encode(target)}"

    /**
     * Prozent-Kodierung nach RFC 3986 (unreserviert: `A-Za-z0-9-_.~`), wie
     * `Uri.EscapeDataString` im C#-Pendant.
     */
    private fun encode(value: String): String {
        val sb = StringBuilder()
        for (byte in value.toByteArray(Charsets.UTF_8)) {
            val unsigned = byte.toInt() and 0xFF
            val c = unsigned.toChar()
            val isUnreserved = c in 'A'..'Z' || c in 'a'..'z' || c in '0'..'9' ||
                c == '-' || c == '_' || c == '.' || c == '~'
            if (isUnreserved) sb.append(c) else sb.append('%').append(String.format("%02X", unsigned))
        }
        return sb.toString()
    }

    /**
     * readEvent liest ein vollständiges Frame über [next] und liefert es
     * zurück. Ein Frame endet mit der Leerzeile, die der Server nach der
     * `event:`- und der `data:`-Zeile schreibt (`ADR-0061`); sie trennt zwei
     * aufeinanderfolgende Events. Ist die Quelle vor dem Frame-Abschluss
     * erschöpft ([next] liefert `null`), liefert readEvent `null`, und ein
     * begonnenes Frame wird verworfen — ein unvollständiges Frame ist kein
     * Event.
     */
    fun readEvent(next: () -> String?): Event? {
        var name: String? = null
        var data: String? = null
        while (true) {
            val line = next() ?: return null
            if (line.isEmpty()) {
                if (name == null && data == null) {
                    continue
                }
                return Event(name ?: "", data ?: "")
            }
            when {
                line.startsWith("event: ") -> name = line.removePrefix("event: ")
                line.startsWith("data: ") -> data = line.removePrefix("data: ")
            }
        }
    }
}
