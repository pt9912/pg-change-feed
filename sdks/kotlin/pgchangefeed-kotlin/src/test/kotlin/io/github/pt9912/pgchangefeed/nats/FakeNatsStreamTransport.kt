package io.github.pt9912.pgchangefeed.nats

/**
 * A network-free [NatsStreamTransport]: every test in this module stays
 * free of real sockets (no real NATS server) — see [NatsStreamTransport]'s
 * KDoc for why `io.nats.client.Connection`/`Subscription` cannot be faked
 * this cheaply. [lastSubject] lets a test assert on the exact subject the
 * client under test subscribed to. Muster:
 * `io.github.pt9912.pgchangefeed.sse.FakeSseTransport`.
 */
internal class FakeNatsStreamTransport private constructor(
    private val payloads: List<ByteArray>,
    private val failure: Exception?,
) : NatsStreamTransport {

    var lastSubject: String? = null
        private set

    override fun subscribe(subject: String): () -> ByteArray? {
        lastSubject = subject
        val iterator = payloads.iterator()
        return {
            if (iterator.hasNext()) {
                iterator.next()
            } else if (failure != null) {
                throw failure
            } else {
                null
            }
        }
    }

    companion object {
        /**
         * Builds a fake transport whose subscription yields the given
         * [payloads] in order, then ends the sequence cleanly — the happy
         * path (`SPEC-024`: one message per row change, in order).
         */
        fun withPayloads(vararg payloads: ByteArray): FakeNatsStreamTransport =
            FakeNatsStreamTransport(payloads.toList(), failure = null)

        /**
         * Builds a fake transport whose subscription yields no message and
         * then throws [failure] from the next-payload supplier itself —
         * simulating what a real NATS subscription does when the
         * underlying connection was rejected or fails (`SPEC-024`'s
         * connection-level auth boundary: a missing or wrong
         * `CDC_NATS_STREAM_TOKEN`). Ausgang bewusst offen (Slice-Plan §6):
         * dieses Fake bildet die Fehler-*Form* nach (eine Exception aus dem
         * Stream selbst), nicht den realen Server-seitigen
         * Ablehnungsmechanismus — ein realer Rundlauf-Beleg bleibt `make
         * test-integration`s `tools/harness/natsstreamsub` vorbehalten.
         */
        fun withFailure(failure: Exception): FakeNatsStreamTransport =
            FakeNatsStreamTransport(emptyList(), failure)
    }
}
