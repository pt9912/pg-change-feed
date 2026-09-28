package cdcexamples.http

import java.nio.charset.StandardCharsets

/**
 * UrlEncoding trägt die Prozent-Kodierung nach RFC 3986 (unreserviert:
 * `A-Za-z0-9-_.~`) statt der Formular-Kodierung von `java.net.URLEncoder`
 * (die Leerzeichen als `+` statt `%20` kodiert) — dieselbe Wahl wie das
 * private `encode` in `TablesUrlBuilder.kt`, hier für die Query-Parameter
 * der übrigen neun Fähigkeiten geteilt, damit das Escaping deckungsgleich
 * mit dem C#-Pendant bleibt (`Uri.EscapeDataString`,
 * `examples/csharp/http-client/ConsumerClient.cs` u. a.).
 */
object UrlEncoding {
    fun encode(value: String): String {
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
