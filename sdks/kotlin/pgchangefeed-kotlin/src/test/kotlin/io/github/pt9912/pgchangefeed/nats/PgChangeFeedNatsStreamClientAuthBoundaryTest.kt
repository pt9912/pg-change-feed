package io.github.pt9912.pgchangefeed.nats

import java.io.IOException
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertSame

/**
 * The `SPEC-024` connection-level auth boundary and the malformed-payload
 * path when subscribing to the NATS full-content stream — no real NATS
 * server involved (`AGENTS.md` §3.1 in the test run). Muster:
 * `io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClientAuthBoundaryTest`,
 * `io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClientAuthBoundaryTest`,
 * adapted to NATS's connection-level (not per-message) auth boundary — same
 * shape as the C# sibling's `PgChangeFeedNatsStreamClientAuthBoundaryTests`.
 */
class PgChangeFeedNatsStreamClientAuthBoundaryTest {

    // Rot färbende Mutation (real geprüft,
    // slice-sdk-kotlin-nats-stream-client-flaeche): `PgChangeFeedNatsStreamClient.streamChanges`s
    // `val payload = nextPayload() ?: break`-Zeile um ein umgebendes
    // `try { ... } catch (ex: Exception) { break }` ergänzt, das eine
    // Transport-Exception verschluckt statt sie weiterzureichen — dieser
    // Test schlägt dann fehl (kein `IOException` mehr geworfen, stattdessen
    // eine leere Sequence), statt still grün zu bleiben.
    @Test
    fun `a rejected connection propagates the underlying NATS exception unwrapped`() {
        // io.nats.client.Nats.connect throws java.io.IOException (checked)
        // when the server rejects the connection (e.g. a missing/wrong
        // server-wide token) — this test's fake models that failure class
        // without a real server; see FakeNatsStreamTransport.withFailure's
        // KDoc for the acknowledged, still-open gap to a real rejection.
        val authError = IOException("authorization violation")
        val transport = FakeNatsStreamTransport.withFailure(authError)
        val client = PgChangeFeedNatsStreamClient(transport)

        val ex = assertFailsWith<IOException> {
            client.streamChanges().toList()
        }

        assertSame(authError, ex)
    }

    @Test
    fun `a non-JSON payload throws MalformedMessage`() {
        val transport = FakeNatsStreamTransport.withPayloads("not json at all".toByteArray(Charsets.UTF_8))
        val client = PgChangeFeedNatsStreamClient(transport)

        assertFailsWith<PgChangeFeedNatsMalformedMessageException> {
            client.streamChanges().toList()
        }
    }

    @Test
    fun `a JSON null payload throws MalformedMessage`() {
        val transport = FakeNatsStreamTransport.withPayloads("null".toByteArray(Charsets.UTF_8))
        val client = PgChangeFeedNatsStreamClient(transport)

        assertFailsWith<PgChangeFeedNatsMalformedMessageException> {
            client.streamChanges().toList()
        }
    }

    @Test
    fun `an empty subscription yields nothing`() {
        val transport = FakeNatsStreamTransport.withPayloads()
        val client = PgChangeFeedNatsStreamClient(transport)

        val received = client.streamChanges().toList()

        assertEquals(emptyList(), received)
    }
}
