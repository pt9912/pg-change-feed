package cdcexamples.http

import java.nio.charset.StandardCharsets

/**
 * TablesUrlBuilder baut die Lese-Adresse der Tabellen-Auflistung über die
 * HTTP-/JSON-API (`LH-FA-SST-006`): ein `GET` auf den `reader`-Endpunkt
 * `/tables` mit seinen zwei Pflichtfeldern `source` und `publication`.
 * `addr` ist die Horch-Adresse des Feed-Containers (`CDC_HTTP_ADDR`, Form
 * `host:port`). Reine Funktion, Form-Vorbild:
 * `examples/http-client/tables.go`.
 *
 * Die Prozent-Kodierung folgt RFC 3986 (unreserviert:
 * `A-Za-z0-9-_.~`) statt der Formular-Kodierung von
 * `java.net.URLEncoder` (die Leerzeichen als `+` statt `%20` kodiert) — das
 * hält das Escaping deckungsgleich mit dem C#-Pendant
 * (`Uri.EscapeDataString`, `examples/csharp/http-client/TablesUrlBuilder.cs`).
 */
object TablesUrlBuilder {
    fun build(addr: String, source: String, publication: String): String {
        val query = listOf("publication" to publication, "source" to source)
            .joinToString("&") { (key, value) -> "$key=${encode(value)}" }
        return "http://$addr/tables?$query"
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
