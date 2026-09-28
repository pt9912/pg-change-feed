package cdcexamples.http

import com.google.gson.annotations.SerializedName
import java.net.http.HttpClient

/**
 * RunRetentionRequest/-Response spiegeln `POST /retention/run`
 * (`LH-FA-RET-002`): `minAgeNanos` 0 heißt „kein zeitliches Mindestalter"
 * und ist gültig. Form-Vorbild: `examples/http-client/retention.go`.
 */
data class RunRetentionRequest(
    @SerializedName("source") val source: String,
    @SerializedName("min_age_nanos") val minAgeNanos: Long,
)

data class RunRetentionResponse(
    @SerializedName("deleted") val deleted: Int,
)

/**
 * RetentionClient ruft den `admin`-Endpunkt `POST /retention/run` auf.
 * Form-Vorbild: `examples/http-client/retention.go`.
 */
object RetentionClient {
    fun runRetention(httpClient: HttpClient, cfg: Config): RunRetentionResponse =
        RequestHelper.sendJson(
            httpClient, "POST", "http://${cfg.addr}/retention/run", cfg.adminToken,
            RunRetentionRequest(cfg.source, cfg.minAgeNanos), 200,
        )
}
