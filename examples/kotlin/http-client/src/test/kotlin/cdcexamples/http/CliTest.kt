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

    @Test
    fun parseUsesEnvironmentDefaults() {
        val got = Cli.parse(emptyArray()) { name ->
            when (name) {
                "CDC_HTTP_ADDR" -> "feed:8080"
                "CDC_API_TOKEN_READER" -> "tok-env"
                else -> null
            }
        }
        assertEquals(Config("feed:8080", "tok-env", "", ""), got)
    }

    @Test
    fun parseFlagsOverrideEnvironment() {
        val got = Cli.parse(
            arrayOf("--addr", "override:9090", "--token=tok-flag", "--source", "quelle-1", "--publication", "pub-1"),
        ) { name -> if (name == "CDC_HTTP_ADDR") "feed:8080" else null }
        assertEquals(Config("override:9090", "tok-flag", "quelle-1", "pub-1"), got)
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
