package cdcexamples.sse

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

/**
 * Prüft den Aufbau der Laufzeit-Konfiguration aus Flags und einem
 * injizierten Umgebungs-Lookup — netzlos, reine Funktion. Form-Vorbild:
 * `examples/csharp/sse-client/SseClient.Tests/CliTests.cs`.
 */
class CliTest {
    private val emptyEnv: (String) -> String? = { null }

    @Test
    fun parseUsesEnvironmentDefaults() {
        val got = Cli.parse(emptyArray()) { name ->
            when (name) {
                "CDC_HTTP_ADDR" -> "feed:8080"
                "CDC_API_TOKEN_READER" -> "tok-env"
                else -> null
            }
        }
        assertEquals(Config("feed:8080", "tok-env"), got)
    }

    @Test
    fun parseFlagsOverrideEnvironment() {
        val got = Cli.parse(
            arrayOf("--addr", "override:9090", "--token=tok-flag"),
        ) { name -> if (name == "CDC_HTTP_ADDR") "feed:8080" else null }
        assertEquals(Config("override:9090", "tok-flag"), got)
    }

    @Test
    fun parseTargetFlagFillsTheConfigInBothFlagForms() {
        assertEquals("eu", Cli.parse(arrayOf("--target", "eu"), emptyEnv).target)
        assertEquals("eu", Cli.parse(arrayOf("--target=eu"), emptyEnv).target)
        assertEquals("", Cli.parse(emptyArray(), emptyEnv).target)
    }

    @Test
    fun parseSchemaAndTableFlagsFillTheConfigInBothFlagForms() {
        val got = Cli.parse(arrayOf("--schema", "public", "--table=orders"), emptyEnv)
        assertEquals("public", got.schema)
        assertEquals("orders", got.table)
        val none = Cli.parse(emptyArray(), emptyEnv)
        assertEquals("", none.schema)
        assertEquals("", none.table)
    }

    @Test
    fun flagsWireIntoTheStreamUrl() {
        val base = "http://feed:8080/changes/stream"
        val cases = listOf(
            emptyList<String>() to base,
            listOf("--schema", "public", "--table", "orders", "--target", "eu") to
                "$base?schema=public&table=orders&target=eu",
            listOf("--table=orders") to "$base?table=orders",
            listOf("--target", "eu") to "$base?target=eu",
        )
        for ((flags, want) in cases) {
            val args = (listOf("--addr", "feed:8080") + flags).toTypedArray()
            assertEquals(want, SseStream.streamUrl(Cli.parse(args, emptyEnv)), "flags $flags")
        }
    }

    @Test
    fun parseRejectsUnknownFlag() {
        assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--unknown"), emptyEnv) }
    }

    @Test
    fun parseRejectsMissingFlagValue() {
        assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--addr"), emptyEnv) }
    }
}
