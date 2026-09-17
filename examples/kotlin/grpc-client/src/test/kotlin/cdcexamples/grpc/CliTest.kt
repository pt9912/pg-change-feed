package cdcexamples.grpc

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

/**
 * Prüft den Aufbau der Laufzeit-Konfiguration aus Flags und einem
 * injizierten Umgebungs-Lookup — netzlos, reine Funktion. Form-Vorbild:
 * `examples/csharp/grpc-client/GrpcClient.Tests/CliTests.cs` (`slice-102`),
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
                else -> null
            }
        }
        assertEquals(Config("feed:9090", "reader-token"), got)
    }

    @Test
    fun parseFlagOverridesEnvironment() {
        val got = Cli.parse(
            arrayOf("--addr", "other:9091", "--token=other-token"),
        ) { name -> if (name == "CDC_GRPC_ADDR") "feed:9090" else null }
        assertEquals(Config("other:9091", "other-token"), got)
    }

    @Test
    fun parseNoEnvironmentNoFlagsReturnsEmpty() {
        val got = Cli.parse(emptyArray(), emptyEnv)
        assertEquals(Config("", ""), got)
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
}
