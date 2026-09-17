package cdcexamples.nats

import java.nio.charset.StandardCharsets

/**
 * ChangesUrl baut die Lese-Adresse der Änderungen einer
 * Quelle/Schema/Tabelle über die HTTP-/JSON-API (`LH-FA-SST-006`): ein
 * `GET` auf den `reader`-Endpunkt `/changes` mit den drei Filtern aus dem
 * Wecksignal-Subjekt (`source`, `schema`, `table`) — der zweiseitige Ablauf
 * aus `ADR-0079`. `addr` ist die Horch-Adresse des Feed-Containers
 * (`CDC_HTTP_ADDR`, Form `host:port`). Reine Funktion, Form-Vorbild:
 * `examples/nats-client/subject.go` (`ChangesURL`).
 *
 * Die Prozent-Kodierung folgt RFC 3986 (unreserviert: `A-Za-z0-9-_.~`) statt
 * der Formular-Kodierung von `java.net.URLEncoder` — dasselbe Muster wie
 * `examples/kotlin/http-client/TablesUrlBuilder.kt` und deckungsgleich mit
 * dem C#-Pendant (`Uri.EscapeDataString`,
 * `examples/csharp/nats-client/ChangesUrlBuilder.cs`).
 */
object ChangesUrl {
    fun build(addr: String, source: String, schema: String, table: String): String {
        val query = listOf("schema" to schema, "source" to source, "table" to table)
            .joinToString("&") { (key, value) -> "$key=${encode(value)}" }
        return "http://$addr/changes?$query"
    }

    private fun encode(value: String): String {
        val sb = StringBuilder()
        for (byte in value.toByteArray(StandardCharsets.UTF_8)) {
            val unsigned = byte.toInt() and 0xFF
            val c = unsigned.toChar()
            val isUnreserved = c in 'A'..'Z' || c in 'a'..'z' || c in '0'..'9' ||
                c == '-' || c == '_' || c == '.' || c == '~'
            if (isUnreserved) {
                sb.append(c)
            } else {
                sb.append('%')
                sb.append(String.format("%02X", unsigned))
            }
        }
        return sb.toString()
    }
}
