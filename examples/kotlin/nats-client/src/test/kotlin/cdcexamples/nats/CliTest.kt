package cdcexamples.nats

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

/**
 * Prüft den Aufbau der Laufzeit-Konfiguration aus Flags und einem
 * injizierten Umgebungs-Lookup — netzlos, reine Funktion. Form-Vorbild:
 * `examples/csharp/nats-client/NatsClient.Tests/CliTests.cs`.
 */
class CliTest {
    private val emptyEnv: (String) -> String? = { null }

    @Test
    fun parseUsesEnvironmentDefaults() {
        val got = Cli.parse(emptyArray()) { name ->
            when (name) {
                "CDC_NATS_URL" -> "nats://feed:4222"
                "CDC_HTTP_ADDR" -> "feed:8080"
                "CDC_API_TOKEN_READER" -> "tok-env"
                else -> null
            }
        }
        assertEquals(Config("nats://feed:4222", "feed:8080", "tok-env", "", "", ""), got)
    }

    @Test
    fun parseFlagsOverrideEnvironment() {
        val got = Cli.parse(
            arrayOf(
                "--nats-url", "nats://override:4222",
                "--addr", "override:9090",
                "--token=tok-flag",
                "--source", "quelle-1",
                "--schema", "public",
                "--table", "orders",
            ),
        ) { name -> if (name == "CDC_NATS_URL") "nats://feed:4222" else null }
        assertEquals(Config("nats://override:4222", "override:9090", "tok-flag", "quelle-1", "public", "orders"), got)
    }

    @Test
    fun parseRejectsUnknownFlag() {
        assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--unknown"), emptyEnv) }
    }

    @Test
    fun parseRejectsMissingFlagValue() {
        assertFailsWith<IllegalArgumentException> { Cli.parse(arrayOf("--nats-url"), emptyEnv) }
    }
}
