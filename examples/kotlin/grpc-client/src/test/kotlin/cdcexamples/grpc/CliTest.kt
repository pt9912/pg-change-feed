package cdcexamples.grpc

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

/**
 * Prüft den Aufbau der Laufzeit-Konfiguration aus Flags und einem
 * injizierten Umgebungs-Lookup — netzlos, reine Funktion. Diese Funktion
 * prüft nur die Flag-Syntax; Rechtsklassen und Pflichtfelder je Verb prüft
 * [ValidatorTest]. Form-Vorbild:
 * `examples/csharp/grpc-client/GrpcClient.Tests/CliTests.cs`,
 * `examples/kotlin/http-client/CliTest.kt`.
 */
class CliTest {
    private val emptyEnv: (String) -> String? = { null }

    @Test
    fun parseUsesEnvironmentDefaults() {
        val got = Cli.parse(emptyArray()) { name ->
            when (name) {
                "CDC_GRPC_ADDR" -> "feed:9090"
                "CDC_API_TOKEN_READER" -> "reader-token"
                "CDC_API_TOKEN_ADMIN" -> "admin-token"
                else -> null
            }
        }
        assertEquals("feed:9090", got.addr)
        assertEquals("reader-token", got.token)
        assertEquals("admin-token", got.adminToken)
    }

    @Test
    fun parseCaFileFromEnvironmentAndFlag() {
        val env: (String) -> String? = { name -> if (name == "CDC_TLS_CA_FILE") "/env.pem" else null }
        assertEquals("", Cli.parse(emptyArray(), emptyEnv).caFile)
        assertEquals("/env.pem", Cli.parse(emptyArray(), env).caFile)
        assertEquals("/flag.pem", Cli.parse(arrayOf("--ca-file=/flag.pem"), env).caFile)
    }

    @Test
    fun parseDefaultVerbIsStream() {
        val got = Cli.parse(emptyArray(), emptyEnv)
        assertEquals("stream", got.verb)
    }

    @Test
    fun parseFlagOverridesEnvironment() {
        val got = Cli.parse(
            arrayOf("--addr", "other:9091", "--token=other-token"),
        ) { name -> if (name == "CDC_GRPC_ADDR") "feed:9090" else null }
        assertEquals("other:9091", got.addr)
        assertEquals("other-token", got.token)
    }

    @Test
    fun parseNoEnvironmentNoFlagsReturnsEmpty() {
        val got = Cli.parse(emptyArray(), emptyEnv)
        assertEquals("", got.addr)
        assertEquals("", got.token)
        assertEquals("", got.adminToken)
    }

    @Test
    fun parseStreamFilterFlags() {
        val got = Cli.parse(arrayOf("--schema=public", "--table=orders"), emptyEnv)
        assertEquals("public", got.schema)
        assertEquals("orders", got.table)
    }

    @Test
    fun parseTargetFlagInBothFlagForms() {
        assertEquals("eu", Cli.parse(arrayOf("--target=eu"), emptyEnv).target)
        assertEquals("eu", Cli.parse(arrayOf("--target", "eu"), emptyEnv).target)
        assertEquals("", Cli.parse(emptyArray(), emptyEnv).target)
    }

    @Test
    fun streamRequestCarriesTargetWithSchemaAndTable() {
        val request = StreamRequest.build(Cli.parse(arrayOf("--schema=public", "--table=orders", "--target=eu"), emptyEnv))
        assertEquals("public", request.schema)
        assertEquals("orders", request.table)
        assertEquals("eu", request.target)
        assertEquals(
            cdc.stream.v1.Changestream.StreamChangesRequest.getDefaultInstance(),
            StreamRequest.build(Cli.parse(emptyArray(), emptyEnv)),
        )
    }

    @Test
    fun readChangesRequestCarriesTarget() {
        val request = ChangesClient.buildRequest(Cli.parse(arrayOf("--source=quelle-1", "--schema=public", "--target=eu"), emptyEnv))
        assertEquals("quelle-1", request.source)
        assertEquals("public", request.schema)
        assertEquals("eu", request.target)
        assertEquals("", ChangesClient.buildRequest(Cli.parse(arrayOf("--source=quelle-1"), emptyEnv)).target)
    }

    @Test
    fun parseEnableTableFields() {
        val got = Cli.parse(
            arrayOf(
                "--verb=enable-table", "--source=quelle-1", "--schema=public", "--table=orders",
                "--table-id=public.orders", "--schema-version-id=public.orders-v2", "--version=2",
                "--publication=pub_quelle_1",
            ),
            emptyEnv,
        )
        assertEquals("enable-table", got.verb)
        assertEquals("quelle-1", got.source)
        assertEquals("public.orders", got.tableId)
        assertEquals("public.orders-v2", got.schemaVersionId)
        assertEquals(2L, got.version)
        assertEquals("pub_quelle_1", got.publication)
    }

    @Test
    fun parseConsumerFields() {
        val got = Cli.parse(arrayOf("--consumer-id=c-1", "--name=Consumer", "--offset=42"), emptyEnv)
        assertEquals("c-1", got.consumerId)
        assertEquals("Consumer", got.name)
        assertEquals(42L, got.offset)
    }

    @Test
    fun parseReadChangesRangeFields() {
        val got = Cli.parse(arrayOf("--from=10", "--to=20", "--limit=5"), emptyEnv)
        assertEquals(10L, got.from)
        assertEquals(20L, got.to)
        assertEquals(5L, got.limit)
    }

    @Test
    fun parseRunRetentionMinAgeNanos() {
        val got = Cli.parse(arrayOf("--min-age-nanos=1000000000"), emptyEnv)
        assertEquals(1_000_000_000L, got.minAgeNanos)
    }

    @Test
    fun parseRejectsUnknownFlag() {
        val ex = assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--bogus"), emptyEnv) }
        assert(ex.message?.contains("unbekanntes Flag") == true)
    }

    @Test
    fun parseRejectsMissingFlagValue() {
        val ex = assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--addr"), emptyEnv) }
        assert(ex.message?.contains("braucht einen Wert") == true)
    }

    @Test
    fun parseRejectsNonIntegerOffset() {
        val ex = assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--offset=abc"), emptyEnv) }
        assert(ex.message?.contains("nicht-negative Ganzzahl") == true)
    }

    @Test
    fun parseRejectsNonIntegerVersion() {
        val ex = assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--version=abc"), emptyEnv) }
        assert(ex.message?.contains("Ganzzahl") == true)
    }
}
