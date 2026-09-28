package cdcexamples.http

import com.google.gson.annotations.SerializedName
import java.net.http.HttpClient

/**
 * RegisterConsumerRequest/-Response spiegeln den JSON-Vertrag von
 * `POST /consumers` (`LH-FA-CON-001`): beide Request-Felder Pflicht,
 * `alreadyRegistered` trägt die Idempotenz-Antwort ohne eigenen
 * Statuscode. Form-Vorbild: `examples/http-client/consumer.go`.
 */
data class RegisterConsumerRequest(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("name") val name: String,
)

data class RegisterConsumerResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("name") val name: String,
    @SerializedName("already_registered") val alreadyRegistered: Boolean,
)

/**
 * AcknowledgeConsumerRequest/-Response spiegeln `POST /consumers/acknowledge`
 * (`LH-FA-CON-004`).
 */
data class AcknowledgeConsumerRequest(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("source_id") val sourceId: String,
    @SerializedName("offset") val offset: Long,
)

data class AcknowledgeConsumerResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("source_id") val sourceId: String,
    @SerializedName("offset") val offset: Long,
)

/**
 * ConsumerPositionResponse spiegelt `GET /consumers/position`
 * (`LH-FA-CON-005`): `acknowledged` unterscheidet die definierte
 * Anfangsposition (kein Nachweis) von einer echten Bestätigung mit
 * Offset 0.
 */
data class ConsumerPositionResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("source_id") val sourceId: String,
    @SerializedName("offset") val offset: Long,
    @SerializedName("acknowledged") val acknowledged: Boolean,
)

/**
 * RemoveConsumerRequest/-Response spiegeln `POST /consumers/remove`
 * (`LH-FA-CON-006`): `removed` trägt den Idempotenz-Ausgang — ein nie
 * registrierter Consumer meldet `false`, keinen `404`.
 */
data class RemoveConsumerRequest(
    @SerializedName("consumer_id") val consumerId: String,
)

data class RemoveConsumerResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("removed") val removed: Boolean,
)

/**
 * ConsumerPositionUrlBuilder baut die Lese-Adresse von
 * `GET /consumers/position`: ihr einziges Pflichtfeld ist `consumer_id`.
 * Form-Vorbild: `ConsumerPositionURL` in `examples/http-client/consumer.go`.
 */
object ConsumerPositionUrlBuilder {
    fun build(addr: String, consumerId: String): String =
        "http://$addr/consumers/position?consumer_id=${UrlEncoding.encode(consumerId)}"
}

/**
 * ConsumerClient ruft die vier Consumer-Verwaltungs-Fähigkeiten der
 * HTTP-/JSON-API auf. Form-Vorbild: `examples/http-client/consumer.go`.
 */
object ConsumerClient {
    fun registerConsumer(httpClient: HttpClient, cfg: Config): RegisterConsumerResponse =
        RequestHelper.sendJson(
            httpClient, "POST", "http://${cfg.addr}/consumers", cfg.adminToken,
            RegisterConsumerRequest(cfg.consumerId, cfg.name), 201,
        )

    fun acknowledgeConsumer(httpClient: HttpClient, cfg: Config): AcknowledgeConsumerResponse =
        RequestHelper.sendJson(
            httpClient, "POST", "http://${cfg.addr}/consumers/acknowledge", cfg.adminToken,
            AcknowledgeConsumerRequest(cfg.consumerId, cfg.source, cfg.offset), 200,
        )

    fun consumerPosition(httpClient: HttpClient, cfg: Config): ConsumerPositionResponse =
        RequestHelper.sendJson(
            httpClient, "GET", ConsumerPositionUrlBuilder.build(cfg.addr, cfg.consumerId), cfg.token,
            null, 200,
        )

    fun removeConsumer(httpClient: HttpClient, cfg: Config): RemoveConsumerResponse =
        RequestHelper.sendJson(
            httpClient, "POST", "http://${cfg.addr}/consumers/remove", cfg.adminToken,
            RemoveConsumerRequest(cfg.consumerId), 200,
        )
}
