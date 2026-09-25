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
 * Client for the PG Change Feed HTTP/JSON API: one method per capability —
 * `registerConsumer`, `acknowledgeConsumer`, `getConsumerPosition`,
 * `removeConsumer`, `enableTable`, `disableTable`, `getStatus`, `listTables`,
 * `runRetention` and `readChanges`. Requests and responses are typed data
 * classes that mirror the JSON documents of the API
 * (`io.github.pt9912.pgchangefeed.http.model`); every non-success response
 * becomes a typed [PgChangeFeedException] subtype instead of a raw HTTP status
 * code mixed with the success path. Connection errors and timeouts of the
 * transport are not converted; they reach the caller as the `IOException` of
 * `java.net.http.HttpClient`.
 *
 * The `java.net.http.HttpClient` passed to the public constructor is passed
 * in, not owned — the caller controls its lifetime, connection pooling and any
 * redirect/proxy configuration; this class never closes it. The bearer token
 * and server address come from [PgChangeFeedClientOptions], supplied at
 * construction — there is no global or static state, so a process can hold
 * several independently configured instances at once.
 */
class PgChangeFeedHttpClient internal constructor(
    private val transport: HttpTransport,
    private val options: PgChangeFeedClientOptions,
) {
    /**
     * Public constructor — wraps the caller's `java.net.http.HttpClient` in
     * the real [JdkHttpTransport]; see [HttpTransport] for why the client
     * depends on a transport interface rather than on `HttpClient` directly.
     */
    constructor(httpClient: HttpClient, options: PgChangeFeedClientOptions) :
        this(JdkHttpTransport(httpClient), options)

    /**
     * Registers a consumer, a named reader whose position the server keeps
     * (`POST /consumers`, admin token). Registering an existing consumer
     * changes nothing; the response reports it with `alreadyRegistered`.
     */
    fun registerConsumer(request: RegisterConsumerRequest): RegisterConsumerResponse =
        post("/consumers", gson.toJson(request))

    /**
     * Stores the position up to which a consumer has processed a source
     * (`POST /consumers/acknowledge`, admin token). Repeating the stored
     * position has no effect; an earlier position, or a position of another
     * source, is rejected with [PgChangeFeedBadRequestException].
     */
    fun acknowledgeConsumer(request: AcknowledgeConsumerRequest): AcknowledgeConsumerResponse =
        post("/consumers/acknowledge", gson.toJson(request))

    /**
     * Reads the stored position of a consumer (`GET /consumers/position`,
     * reader or admin token). `acknowledged` is `false` for a consumer that
     * never acknowledged; `offset` is then the starting position.
     */
    fun getConsumerPosition(consumerId: String): ConsumerPositionResponse =
        get("/consumers/position?" + buildQuery("consumer_id" to consumerId))

    /**
     * Removes a consumer (`POST /consumers/remove`, admin token); `removed` is
     * `false` for one that was never registered.
     */
    fun removeConsumer(consumerId: String): RemoveConsumerResponse =
        post("/consumers/remove", gson.toJson(mapOf("consumer_id" to consumerId)))

    /**
     * Starts capturing a table of a source (`POST /tables/enable`, admin
     * token). `alreadyEnabled` is `true` when the table was captured already; a
     * table that does not exist in the source database raises
     * [PgChangeFeedNotFoundException].
     */
    fun enableTable(request: EnableTableRequest): EnableTableResponse =
        post("/tables/enable", gson.toJson(request))

    /**
     * Stops capturing a table (`POST /tables/disable`, admin token). `retained`
     * is `true` when changes already stored for the table remain readable; a
     * table that does not exist in the source database raises
     * [PgChangeFeedNotFoundException].
     */
    fun disableTable(request: DisableTableRequest): DisableTableResponse =
        post("/tables/disable", gson.toJson(request))

    /**
     * Tells whether a table is captured (`enabled`) or no longer captured with
     * stored changes remaining (`retained`) (`GET /tables/status`, reader or
     * admin token). A table that was never enabled reports both as `false`; a
     * table that does not exist in the source database raises
     * [PgChangeFeedNotFoundException].
     */
    fun getStatus(source: String, schema: String, table: String, publication: String): TableStatusResponse =
        get(
            "/tables/status?" + buildQuery(
                "source" to source,
                "schema" to schema,
                "table" to table,
                "publication" to publication,
            ),
        )

    /**
     * Lists the captured tables and the tables that are no longer captured but
     * whose stored changes remain (`GET /tables`, reader or admin token).
     */
    fun listTables(source: String, publication: String): ListTablesResponse =
        get("/tables?" + buildQuery("source" to source, "publication" to publication))

    /**
     * Deletes the stored changes of a source that are older than `minAgeNanos`
     * and that every consumer with a stored position has passed
     * (`POST /retention/run`, admin token); `deleted` is the number removed.
     */
    fun runRetention(request: RunRetentionRequest): RunRetentionResponse =
        post("/retention/run", gson.toJson(request))

    /**
     * Reads stored changes of a source (`GET /changes`, reader or admin
     * token). [source] is mandatory; [schema]/[table] are optional and
     * independent; [from] is inclusive, [to] exclusive (both `commit_position`
     * values). [limit] cuts rows, not positions, and there is no default
     * limit. A range without changes returns an empty list.
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
                    "not valid JSON.",
                ex,
            )
        }
        return result
            ?: throw PgChangeFeedMalformedResponseException(
                statusCode,
                "PG Change Feed HTTP API returned status $statusCode with an empty or null response " +
                    "body.",
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
     * instead of `%20`).
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
