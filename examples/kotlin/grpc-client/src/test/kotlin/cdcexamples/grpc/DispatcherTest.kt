package cdcexamples.grpc

import cdc.administration.v1.AdministrationGrpcKt
import io.grpc.ManagedChannelBuilder
import kotlinx.coroutines.runBlocking
import kotlin.test.AfterTest
import kotlin.test.Test
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

/**
 * Prüft den defensiven Schutz von [Dispatcher.dispatchAdmin] gegen ein
 * unbekanntes Verb — für den Fall eines Aufrufs ohne vorherige
 * [Validator.validate]. Der Stub bindet an einen nie kontaktierten Kanal
 * (gRPC-Kanäle verbinden erst beim ersten Aufruf) — kein Netzwerkzugriff.
 * Form-Vorbild:
 * `examples/csharp/grpc-client/GrpcClient.Tests/DispatcherTests.cs`,
 * `examples/grpc-client/main_test.go`, `TestDispatchAdminRejectsUnknownVerb`.
 */
class DispatcherTest {
    private val channel = ManagedChannelBuilder.forTarget("localhost:0").usePlaintext().build()

    @AfterTest
    fun shutdownChannel() {
        channel.shutdownNow()
    }

    @Test
    fun dispatchAdminRejectsUnknownVerbBeforeAnyNetworkCall() {
        val client = AdministrationGrpcKt.AdministrationCoroutineStub(channel)
        val cfg = Config(
            addr = "feed:9090", token = "", adminToken = "", verb = "unbekannt",
            schema = "", table = "", consumerId = "", name = "", offset = 0,
            tableId = "", schemaVersionId = "", version = 1,
            source = "", publication = "", from = 0, to = 0, limit = 0, minAgeNanos = 0,
        )

        val ex = assertFailsWith<IllegalStateException> {
            runBlocking { Dispatcher.dispatchAdmin(client, cfg) }
        }
        assertTrue(ex.message!!.contains("unbekanntes --verb"))
    }
}
