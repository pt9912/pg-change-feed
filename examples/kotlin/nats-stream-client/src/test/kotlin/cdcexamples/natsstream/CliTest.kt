package cdcexamples.natsstream

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

/**
 * Prüft den Aufbau der Laufzeit-Konfiguration aus Flags und einem
 * injizierten Umgebungs-Lookup — netzlos, reine Funktion. Form-Vorbild:
 * `examples/kotlin/grpc-client/src/test/kotlin/cdcexamples/grpc/CliTest.kt`.
 */
class CliTest {
    private val emptyEnv: (String) -> String? = { null }

    @Test
    fun parseUsesEnvironmentDefaults() {
        val got = Cli.parse(emptyArray()) { name ->
            when (name) {
                "CDC_NATS_URL" -> "nats://feed:4222"
                "CDC_NATS_STREAM_TOKEN" -> "stream-token"
                else -> null
            }
        }
        assertEquals(Config("nats://feed:4222", "stream-token"), got)
    }

    @Test
    fun parseFlagOverridesEnvironment() {
        val got = Cli.parse(
            arrayOf("--nats-url", "nats://override:4222", "--token=other-token"),
        ) { name -> if (name == "CDC_NATS_URL") "nats://feed:4222" else null }
        assertEquals(Config("nats://override:4222", "other-token"), got)
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
        val ex = assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--nats-url"), emptyEnv) }
        assert(ex.message?.contains("braucht einen Wert") == true)
    }
}
