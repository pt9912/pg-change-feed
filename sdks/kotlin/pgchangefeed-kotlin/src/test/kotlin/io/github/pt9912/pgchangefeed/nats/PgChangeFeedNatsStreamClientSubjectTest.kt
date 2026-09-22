package io.github.pt9912.pgchangefeed.nats

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

/**
 * Subject-Formatierung netzlos geprüft (`SPEC-024`): [PgChangeFeedNatsStreamClient.buildSubject]/
 * [PgChangeFeedNatsStreamClient.buildSourceSubject] bauen den
 * dokumentierten vier- bzw. dreistufigen Subjekt-Aufbau, [PgChangeFeedNatsStreamClient.streamChanges]
 * lässt seinen Default ([PgChangeFeedNatsStreamClient.ALL_SOURCES_SUBJECT])
 * unverändert am Transport ankommen. Muster: der C# Sibling
 * `PgChangeFeed.Client.Tests.Nats.PgChangeFeedNatsStreamClientTests`'
 * Subjekt-Assertions.
 */
class PgChangeFeedNatsStreamClientSubjectTest {

    @Test
    fun `buildSubject formats the four-token subject`() {
        assertEquals(
            "cdc.stream.source-1.public.orders",
            PgChangeFeedNatsStreamClient.buildSubject("source-1", "public", "orders"),
        )
    }

    @Test
    fun `buildSourceSubject formats the three-token wildcard subject`() {
        assertEquals(
            "cdc.stream.source-1.>",
            PgChangeFeedNatsStreamClient.buildSourceSubject("source-1"),
        )
    }

    @Test
    fun `buildSubject rejects a token carrying a subject separator`() {
        assertFailsWith<IllegalArgumentException> {
            PgChangeFeedNatsStreamClient.buildSubject("source.1", "public", "orders")
        }
    }

    @Test
    fun `buildSubject rejects a token carrying a wildcard character`() {
        assertFailsWith<IllegalArgumentException> {
            PgChangeFeedNatsStreamClient.buildSubject("source-1", "public", "orders*")
        }
    }

    @Test
    fun `buildSubject rejects a blank token`() {
        assertFailsWith<IllegalArgumentException> {
            PgChangeFeedNatsStreamClient.buildSubject("source-1", "  ", "orders")
        }
    }

    @Test
    fun `streamChanges defaults to the all-sources wildcard subject`() {
        val transport = FakeNatsStreamTransport.withPayloads()
        val client = PgChangeFeedNatsStreamClient(transport)

        client.streamChanges().toList()

        assertEquals(PgChangeFeedNatsStreamClient.ALL_SOURCES_SUBJECT, transport.lastSubject)
    }

    @Test
    fun `streamChanges subscribes to the supplied subject`() {
        val transport = FakeNatsStreamTransport.withPayloads()
        val client = PgChangeFeedNatsStreamClient(transport)
        val subject = PgChangeFeedNatsStreamClient.buildSubject("source-1", "public", "orders")

        client.streamChanges(subject).toList()

        assertEquals("cdc.stream.source-1.public.orders", transport.lastSubject)
    }
}
