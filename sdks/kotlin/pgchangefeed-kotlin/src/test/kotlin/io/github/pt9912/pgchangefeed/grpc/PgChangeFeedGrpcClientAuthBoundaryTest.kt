package io.github.pt9912.pgchangefeed.grpc

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.grpc.Metadata
import io.grpc.Status
import io.grpc.StatusException
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.runBlocking
import java.net.URI
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

/**
 * [PgChangeFeedGrpcClient] against a fake [GrpcStreamTransport] (no real
 * server, no network — `AGENTS.md` §3.1 in the test run): the bearer-token
 * metadata form (`SPEC-020`) and the auth boundary (`UNAUTHENTICATED`
 * surfaces from the stream itself, not a swallowed empty flow) — the
 * consumer-side counterpart of the server-side fitness function `ADR-0060`
 * Teilfrage 4 names (Go-side test vorbild cited in
 * [PgChangeFeedGrpcClientMessageSchemaTest]'s KDoc).
 */
class PgChangeFeedGrpcClientAuthBoundaryTest {
    private val authorizationKey: Metadata.Key<String> =
        Metadata.Key.of("authorization", Metadata.ASCII_STRING_MARSHALLER)

    private fun options(token: String = "reader-token") =
        PgChangeFeedClientOptions(URI("http://localhost:50051"), token)

    @Test
    fun `streamChanges sends bearer token in authorization metadata`() = runBlocking {
        val transport = FakeGrpcStreamTransport.withMessages()
        val client = PgChangeFeedGrpcClient(transport, options("reader-token"))

        client.streamChanges().toList()

        assertEquals("Bearer reader-token", transport.lastHeaders?.get(authorizationKey))
    }

    // Rot färbende Mutation (real geprüft, slice-sdk-kotlin-grpc-client-flaeche):
    // `PgChangeFeedGrpcClient.streamChanges()` um `.catch { }` erweitert, das
    // den `StatusException` still schluckt statt ihn weiterzureichen — dieser
    // Test schlägt dann fehl (kein `StatusException` mehr geworfen), statt
    // still grün zu bleiben.
    @Test
    fun `streamChanges with missing or invalid token throws unauthenticated`() = runBlocking {
        val transport = FakeGrpcStreamTransport.withStatus(
            Status.UNAUTHENTICATED.withDescription("missing or unknown bearer token"),
        )
        val client = PgChangeFeedGrpcClient(transport, options("unknown-token"))

        val ex = assertFailsWith<StatusException> {
            client.streamChanges().toList()
        }

        assertEquals(Status.Code.UNAUTHENTICATED, ex.status.code)
    }
}
