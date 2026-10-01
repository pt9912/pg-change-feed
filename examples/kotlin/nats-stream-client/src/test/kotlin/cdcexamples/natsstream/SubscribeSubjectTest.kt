package cdcexamples.natsstream

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

/**
 * Prüft die Ableitung des Abonnement-Subjekts aus `--source` und
 * `--target` — netzlos, reine Funktion. Form-Vorbild:
 * `examples/nats-stream-client/subject_test.go`.
 */
class SubscribeSubjectTest {
    @Test
    fun withoutSourceAndTargetIsTheRootWildcard() {
        assertEquals("cdc.stream.>", SubscribeSubject.resolve("", ""))
    }

    @Test
    fun sourceAndTargetIsTheRouteSubject() {
        assertEquals("cdc.route.quelle-1.eu", SubscribeSubject.resolve("quelle-1", "eu"))
    }

    @Test
    fun halfOrInvalidInputIsRejected() {
        val inputs = listOf(
            "quelle-1" to "", "" to "eu",
            "quelle-1" to "   ", "quelle-1" to "a.b", "quelle-1" to "*", "quelle-1" to ">", "quelle-1" to "a b",
            "a.b" to "eu", "*" to "eu", ">" to "eu",
        )
        for ((source, target) in inputs) {
            assertFailsWith<IllegalArgumentException>("$source/$target") { SubscribeSubject.resolve(source, target) }
        }
    }

    @Test
    fun parseSourceAndTargetFlagsFillTheConfig() {
        val cfg = Cli.parse(arrayOf("--source=quelle-1", "--target", "eu")) { null }
        assertEquals("quelle-1", cfg.source)
        assertEquals("eu", cfg.target)
        assertEquals("", Cli.parse(emptyArray()) { null }.target)
    }
}
