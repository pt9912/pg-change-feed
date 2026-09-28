package cdcexamples.http

import com.google.gson.Gson
import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse
import java.time.Duration

/**
 * RequestHelper trägt den gemeinsamen Anfrage/Antwort-Ablauf aller zehn
 * Fähigkeiten dieses Beispiels: JSON-Body kodieren (falls [requestBody]
 * gesetzt ist), das Bearer-Token setzen, senden, und bei [wantStatus] den
 * Antwort-Body typisiert dekodieren. Ein davon abweichender Statuscode ist
 * ein sichtbarer Fehler mit Statuscode und Antworttext, kein stiller
 * Leerwert — dieselbe Form wie [TablesClient.listTables]. Form-Vorbild:
 * `doRequestJSON` in `examples/http-client/request.go`.
 */
object RequestHelper {
    // `sendJson` ist eine öffentliche `inline`-Funktion mit reifiziertem
    // Typparameter (nötig für `T::class.java`) — sie darf deshalb keine
    // `private` Deklaration erreichen; `@PublishedApi internal` hält beide
    // Felder außerhalb der eigentlichen öffentlichen API dieses Beispiels.
    @PublishedApi
    internal val requestTimeout: Duration = Duration.ofSeconds(10)

    @PublishedApi
    internal val gson = Gson()

    inline fun <reified T> sendJson(
        httpClient: HttpClient,
        method: String,
        url: String,
        token: String,
        requestBody: Any?,
        wantStatus: Int,
    ): T {
        val builder = HttpRequest.newBuilder(URI.create(url))
            .header("Authorization", "Bearer $token")
            .timeout(requestTimeout)

        if (requestBody != null) {
            val json = gson.toJson(requestBody)
            builder.header("Content-Type", "application/json")
                .method(method, HttpRequest.BodyPublishers.ofString(json))
        } else {
            builder.method(method, HttpRequest.BodyPublishers.noBody())
        }

        val response = httpClient.send(builder.build(), HttpResponse.BodyHandlers.ofString())
        if (response.statusCode() != wantStatus) {
            throw IllegalStateException("HTTP-Status ${response.statusCode()}: ${response.body()}")
        }
        return gson.fromJson(response.body(), T::class.java)
            ?: throw IllegalStateException("HTTP-Antwort ist leer oder \"null\"")
    }
}
