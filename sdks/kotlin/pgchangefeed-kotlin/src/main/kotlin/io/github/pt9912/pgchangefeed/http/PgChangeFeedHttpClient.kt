package io.github.pt9912.pgchangefeed.http

import com.google.gson.Gson
import com.google.gson.JsonSyntaxException
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.http.model.AcknowledgeConsumerRequest
import io.github.pt9912.pgchangefeed.http.model.AcknowledgeConsumerResponse
import io.github.pt9912.pgchangefeed.http.model.ConsumerPositionResponse
import io.github.pt9912.pgchangefeed.http.model.DisableTableRequest
import io.github.pt9912.pgchangefeed.http.model.DisableTableResponse
import io.github.pt9912.pgchangefeed.http.model.EnableTableRequest
import io.github.pt9912.pgchangefeed.http.model.EnableTableResponse
import io.github.pt9912.pgchangefeed.http.model.ErrorResponse
import io.github.pt9912.pgchangefeed.http.model.ListTablesResponse
import io.github.pt9912.pgchangefeed.http.model.ReadChangesResponse
import io.github.pt9912.pgchangefeed.http.model.RegisterConsumerRequest
import io.github.pt9912.pgchangefeed.http.model.RegisterConsumerResponse
import io.github.pt9912.pgchangefeed.http.model.RemoveConsumerResponse
import io.github.pt9912.pgchangefeed.http.model.RunRetentionRequest
import io.github.pt9912.pgchangefeed.http.model.RunRetentionResponse
import io.github.pt9912.pgchangefeed.http.model.TableStatusResponse
import java.net.http.HttpClient
import java.nio.charset.StandardCharsets

/**
 * Public entry point for the PG Change Feed HTTP/JSON API: one method per
 * wire capability — the nine port-covered capabilities of SPEC-018
 * (`RegisterConsumer`, `AcknowledgeConsumer`, `GetConsumerPosition`,
 * `RemoveConsumer`, `EnableTable`, `DisableTable`, `GetStatus`,
 * `ListTables`, `RunRetention`) plus reading persisted changes
 * (`readChanges`, SPEC-022). Requests/responses are typed data classes that
 * mirror the SPEC-018/SPEC-022 JSON schemas exactly
 * (`io.github.pt9912.pgchangefeed.http.model`); every non-success response
 * becomes a typed [PgChangeFeedException] subtype instead of a raw HTTP
 * status code mixed with the success path — see that type's KDoc for the
 * sealed-class design decision.
 *
 * The `java.net.http.HttpClient` passed to the public constructor is
 * injected, not owned — the caller controls its lifetime, connection
 * pooling and any redirect/proxy configuration; this type never closes it.
 * The bearer token and server address come from [PgChangeFeedClientOptions],
 * supplied at construction — no global or static state, a process can hold
 * several independently configured instances at once.
 *
 * Draht-Kenntnis-Vorbild (gelesen, nicht importiert — `ADR-0109`
 * Festlegung 3): `examples/kotlin/http-client/TablesClient.kt`,
 * `TablesUrlBuilder.kt`, and `spec/pflichtenheft.md` SPEC-018/SPEC-022
 * directly.
 */
class PgChangeFeedHttpClient internal constructor(
    private val transport: HttpTransport,
    private val options: PgChangeFeedClientOptions,
) {
    /**
     * Public constructor — wraps the caller-supplied `java.net.http.HttpClient`
     * in the real [JdkHttpTransport]. See [HttpTransport]'s KDoc for why the
     * client depends on a transport seam rather than on `HttpClient` directly.
     */
    constructor(httpClient: HttpClient, options: PgChangeFeedClientOptions) :
        this(JdkHttpTransport(httpClient), options)

    /** `RegisterConsumer` — `POST /consumers` (admin, SPEC-018, `LH-FA-CON-001`). */
    fun registerConsumer(request: RegisterConsumerRequest): RegisterConsumerResponse =
        post("/consumers", gson.toJson(request))

    /** `AcknowledgeConsumer` — `POST /consumers/acknowledge` (admin, SPEC-018, `LH-FA-CON-004`). */
    fun acknowledgeConsumer(request: AcknowledgeConsumerRequest): AcknowledgeConsumerResponse =
        post("/consumers/acknowledge", gson.toJson(request))

    /** `GetConsumerPosition` — `GET /consumers/position` (reader or admin, SPEC-018, `LH-FA-CON-003`/`005`). */
    fun getConsumerPosition(consumerId: String): ConsumerPositionResponse =
        get("/consumers/position?" + buildQuery("consumer_id" to consumerId))

    /** `RemoveConsumer` — `POST /consumers/remove` (admin, SPEC-018, `LH-FA-CON-006`). */
    fun removeConsumer(consumerId: String): RemoveConsumerResponse =
        post("/consumers/remove", gson.toJson(mapOf("consumer_id" to consumerId)))

    /** `EnableTable` — `POST /tables/enable` (admin, SPEC-018, `LH-FA-CFG-001`). */
    fun enableTable(request: EnableTableRequest): EnableTableResponse =
        post("/tables/enable", gson.toJson(request))

    /** `DisableTable` — `POST /tables/disable` (admin, SPEC-018, `LH-FA-CFG-002`). */
    fun disableTable(request: DisableTableRequest): DisableTableResponse =
        post("/tables/disable", gson.toJson(request))

    /** `GetStatus` — `GET /tables/status` (reader or admin, SPEC-018, `LH-FA-CFG-003`). */
    fun getStatus(source: String, schema: String, table: String, publication: String): TableStatusResponse =
        get(
            "/tables/status?" + buildQuery(
                "source" to source,
                "schema" to schema,
                "table" to table,
                "publication" to publication,
            ),
        )

    /** `ListTables` — `GET /tables` (reader or admin, SPEC-018, `LH-FA-CFG-004`). */
    fun listTables(source: String, publication: String): ListTablesResponse =
        get("/tables?" + buildQuery("source" to source, "publication" to publication))

    /** `RunRetention` — `POST /retention/run` (admin, SPEC-018, `LH-FA-RET-002`..`004`). */
    fun runRetention(request: RunRetentionRequest): RunRetentionResponse =
        post("/retention/run", gson.toJson(request))

    /**
     * `ReadChanges` — `GET /changes` (reader or admin, SPEC-022). [source]
     * is mandatory; [schema]/[table] are optional and independent; [from]
     * is inclusive, [to] exclusive (both `commit_position` values); there
     * is no default [limit].
     */
    fun readChanges(
        source: String,
        schema: String? = null,
        table: String? = null,
        from: Long? = null,
        to: Long? = null,
        limit: Int? = null,
    ): ReadChangesResponse {
        val query = buildQuery(
            "source" to source,
            "schema" to schema,
            "table" to table,
            "from" to from?.toString(),
            "to" to to?.toString(),
            "limit" to limit?.toString(),
        )
        return get("/changes?$query")
    }

    private fun buildQuery(vararg params: Pair<String, String?>): String =
        params
            .filter { (_, value) -> value != null }
            .joinToString("&") { (key, value) -> "$key=${encode(value!!)}" }

    private inline fun <reified T> get(pathAndQuery: String): T =
        dispatch("GET", pathAndQuery, body = null)

    private inline fun <reified T> post(path: String, body: String): T =
        dispatch("POST", path, body)

    private inline fun <reified T> dispatch(method: String, pathAndQuery: String, body: String?): T {
        val headers = buildMap {
            put("Authorization", "Bearer ${options.apiToken}")
            if (body != null) {
                put("Content-Type", "application/json; charset=utf-8")
            }
        }
        val request = TransportRequest(
            method = method,
            url = options.address.toString().trimEnd('/') + pathAndQuery,
            headers = headers,
            body = body,
        )
        val response = transport.send(request)
        if (response.statusCode !in 200..299) {
            throw buildException(response.statusCode, response.body)
        }
        return parseSuccessBody(response.statusCode, response.body)
    }

    private inline fun <reified T> parseSuccessBody(statusCode: Int, body: String): T {
        val result = try {
            gson.fromJson(body, T::class.java)
        } catch (ex: JsonSyntaxException) {
            throw PgChangeFeedMalformedResponseException(
                statusCode,
                "PG Change Feed HTTP API returned status $statusCode with a response body that is " +
                    "not valid JSON — a protocol violation outside SPEC-018/SPEC-022's documented shapes.",
                ex,
            )
        }
        return result
            ?: throw PgChangeFeedMalformedResponseException(
                statusCode,
                "PG Change Feed HTTP API returned status $statusCode with an empty or null response " +
                    "body — a protocol violation outside SPEC-018/SPEC-022's documented shapes.",
            )
    }

    private fun buildException(statusCode: Int, body: String): PgChangeFeedException {
        val message = extractErrorMessage(body)
        return when (statusCode) {
            400 -> PgChangeFeedBadRequestException(statusCode, message)
            401 -> PgChangeFeedUnauthorizedException(statusCode, message)
            403 -> PgChangeFeedForbiddenException(statusCode, message)
            404 -> PgChangeFeedNotFoundException(statusCode, message)
            500 -> PgChangeFeedServerErrorException(statusCode, message)
            else -> PgChangeFeedUnexpectedStatusException(statusCode, message)
        }
    }

    private fun extractErrorMessage(body: String): String =
        try {
            gson.fromJson(body, ErrorResponse::class.java)?.error ?: body
        } catch (ex: JsonSyntaxException) {
            body
        }

    /**
     * RFC 3986 percent-encoding (unreserved: `A-Za-z0-9-_.~`) — not
     * `java.net.URLEncoder`'s form-encoding (which encodes a space as `+`
     * instead of `%20`). Same escaping as the C# sibling's
     * `Uri.EscapeDataString` and `examples/kotlin/http-client/TablesUrlBuilder.kt`'s
     * hand-rolled encoder (read as a Vorbild, not imported — `ADR-0109`
     * Festlegung 3).
     */
    private fun encode(value: String): String {
        val builder = StringBuilder()
        for (byte in value.toByteArray(StandardCharsets.UTF_8)) {
            val unsigned = byte.toInt() and 0xFF
            val c = unsigned.toChar()
            val isUnreserved = c in 'A'..'Z' || c in 'a'..'z' || c in '0'..'9' ||
                c == '-' || c == '_' || c == '.' || c == '~'
            if (isUnreserved) {
                builder.append(c)
            } else {
                builder.append('%')
                builder.append(String.format("%02X", unsigned))
            }
        }
        return builder.toString()
    }

    private companion object {
        val gson = Gson()
    }
}
