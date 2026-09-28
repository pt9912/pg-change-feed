package cdcexamples.http

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

/**
 * Prüft den Aufbau der Laufzeit-Konfiguration aus Flags und einem
 * injizierten Umgebungs-Lookup — netzlos, reine Funktion. Form-Vorbild:
 * `examples/csharp/http-client/HttpClient.Tests/CliTests.cs`.
 */
class CliTest {
    private val emptyEnv: (String) -> String? = { null }

    private fun expectedDefaults(
        addr: String = "",
        token: String = "",
        adminToken: String = "",
        verb: String = "tables",
        source: String = "",
        publication: String = "",
    ) = Config(
        addr = addr, token = token, adminToken = adminToken, verb = verb,
        source = source, publication = publication,
        consumerId = "", name = "", offset = 0,
        schema = "", table = "", tableId = "", schemaVersionId = "", version = 1,
        from = "", to = "", limit = "", minAgeNanos = 0,
    )

    @Test
    fun parseUsesEnvironmentDefaults() {
        val got = Cli.parse(emptyArray()) { name ->
            when (name) {
                "CDC_HTTP_ADDR" -> "feed:8080"
                "CDC_API_TOKEN_READER" -> "tok-env"
                "CDC_API_TOKEN_ADMIN" -> "tok-admin-env"
                else -> null
            }
        }
        assertEquals(expectedDefaults(addr = "feed:8080", token = "tok-env", adminToken = "tok-admin-env"), got)
    }

    @Test
    fun parseFlagsOverrideEnvironment() {
        val got = Cli.parse(
            arrayOf("--addr", "override:9090", "--token=tok-flag", "--source", "quelle-1", "--publication", "pub-1"),
        ) { name -> if (name == "CDC_HTTP_ADDR") "feed:8080" else null }
        assertEquals(
            expectedDefaults(addr = "override:9090", token = "tok-flag", source = "quelle-1", publication = "pub-1"),
            got,
        )
    }

    @Test
    fun parseVerbAndCapabilityFieldsFillAllTenFieldsIndependently() {
        val got = Cli.parse(
            arrayOf(
                "--verb", "enable-table", "--admin-token", "admin-tok",
                "--source", "src", "--schema", "public", "--table", "orders",
                "--table-id", "public.orders", "--schema-version-id", "public.orders-v1",
                "--version", "2", "--publication", "pub",
                "--consumer-id", "c-1", "--name", "n-1", "--offset", "42",
                "--from", "1", "--to", "10", "--limit", "5", "--min-age-nanos", "-3",
            ),
            emptyEnv,
        )
        assertEquals(
            Config(
                addr = "", token = "", adminToken = "admin-tok", verb = "enable-table",
                source = "src", publication = "pub",
                consumerId = "c-1", name = "n-1", offset = 42,
                schema = "public", table = "orders", tableId = "public.orders", schemaVersionId = "public.orders-v1", version = 2,
                from = "1", to = "10", limit = "5", minAgeNanos = -3,
            ),
            got,
        )
    }

    @Test
    fun parseInlineFlagValueSplitsOnFirstEquals() {
        val got = Cli.parse(arrayOf("--verb=changes"), emptyEnv)
        assertEquals("changes", got.verb)
    }

    @Test
    fun parseRejectsUnknownFlag() {
        assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--unknown"), emptyEnv) }
    }

    @Test
    fun parseRejectsMissingFlagValue() {
        assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--addr"), emptyEnv) }
    }

    @Test
    fun parseRejectsNonNumericOffset() {
        assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--offset", "nicht-zahl"), emptyEnv) }
    }

    @Test
    fun parseRejectsNegativeOffset() {
        assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--offset", "-1"), emptyEnv) }
    }

    @Test
    fun parseRejectsNonNumericVersion() {
        assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--version", "nicht-zahl"), emptyEnv) }
    }

    @Test
    fun parseRejectsNonNumericMinAgeNanos() {
        assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--min-age-nanos", "nicht-zahl"), emptyEnv) }
    }
}
