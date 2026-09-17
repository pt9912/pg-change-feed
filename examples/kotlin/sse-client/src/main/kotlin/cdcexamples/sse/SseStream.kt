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
    fun streamUrl(addr: String): String = "http://$addr/changes/stream"

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
