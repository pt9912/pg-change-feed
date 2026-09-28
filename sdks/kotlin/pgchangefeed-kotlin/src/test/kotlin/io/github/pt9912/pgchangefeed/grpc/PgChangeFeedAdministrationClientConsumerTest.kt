package io.github.pt9912.pgchangefeed.grpc

import cdc.administration.v1.AdministrationOuterClass.AcknowledgeConsumerRequest
import cdc.administration.v1.AdministrationOuterClass.AcknowledgeConsumerResponse
import cdc.administration.v1.AdministrationOuterClass.GetConsumerPositionRequest
import cdc.administration.v1.AdministrationOuterClass.GetConsumerPositionResponse
import cdc.administration.v1.AdministrationOuterClass.RegisterConsumerRequest
import cdc.administration.v1.AdministrationOuterClass.RegisterConsumerResponse
import cdc.administration.v1.AdministrationOuterClass.RemoveConsumerRequest
import cdc.administration.v1.AdministrationOuterClass.RemoveConsumerResponse
import io.grpc.Metadata
import kotlinx.coroutines.runBlocking
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertTrue

/**
 * Happy-path coverage for the four consumer-management RPCs
 * (`registerConsumer`, `acknowledgeConsumer`, `getConsumerPosition`,
 * `removeConsumer`), plus the bearer-token metadata form shared by every
 * method and each RPC's boundary report (already-registered, never
 * acknowledged, never registered).
 */
class PgChangeFeedAdministrationClientConsumerTest {
    private val authorizationKey: Metadata.Key<String> =
        Metadata.Key.of("authorization", Metadata.ASCII_STRING_MARSHALLER)

    @Test
    fun `registerConsumer happy path returns typed response and sends bearer token`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            RegisterConsumerResponse.newBuilder().setConsumerId("c-1").setName("n-1").setAlreadyRegistered(false).build(),
        )
        val client = AdministrationTestClientFactory.create(transport, apiToken = "admin-token")

        val response = client.registerConsumer(
            RegisterConsumerRequest.newBuilder().setConsumerId("c-1").setName("n-1").build(),
        )

        assertEquals("c-1", response.consumerId)
        assertFalse(response.alreadyRegistered)
        assertEquals("Bearer admin-token", transport.lastHeaders?.get(authorizationKey))
    }

    @Test
    fun `registerConsumer existing consumer reports alreadyRegistered`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            RegisterConsumerResponse.newBuilder().setConsumerId("c-1").setName("n-1").setAlreadyRegistered(true).build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.registerConsumer(
            RegisterConsumerRequest.newBuilder().setConsumerId("c-1").setName("n-1").build(),
        )

        assertTrue(response.alreadyRegistered)
    }

    @Test
    fun `acknowledgeConsumer happy path returns typed response`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            AcknowledgeConsumerResponse.newBuilder().setConsumerId("c-1").setSourceId("s-1").setOffset(42).build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.acknowledgeConsumer(
            AcknowledgeConsumerRequest.newBuilder().setConsumerId("c-1").setSourceId("s-1").setOffset(42).build(),
        )

        assertEquals("s-1", response.sourceId)
        assertEquals(42L, response.offset)
    }

    @Test
    fun `getConsumerPosition happy path returns typed response`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            GetConsumerPositionResponse.newBuilder()
                .setConsumerId("c-1").setSourceId("s-1").setOffset(7).setAcknowledged(true)
                .build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.getConsumerPosition(GetConsumerPositionRequest.newBuilder().setConsumerId("c-1").build())

        assertEquals(7L, response.offset)
        assertTrue(response.acknowledged)
    }

    @Test
    fun `getConsumerPosition never acknowledged reports acknowledged false`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            GetConsumerPositionResponse.newBuilder()
                .setConsumerId("c-1").setSourceId("s-1").setOffset(0).setAcknowledged(false)
                .build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.getConsumerPosition(GetConsumerPositionRequest.newBuilder().setConsumerId("c-1").build())

        assertEquals(0L, response.offset)
        assertFalse(response.acknowledged)
    }

    @Test
    fun `removeConsumer happy path returns typed response`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            RemoveConsumerResponse.newBuilder().setConsumerId("c-1").setRemoved(true).build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.removeConsumer(RemoveConsumerRequest.newBuilder().setConsumerId("c-1").build())

        assertTrue(response.removed)
    }

    @Test
    fun `removeConsumer never registered reports removed false`() = runBlocking {
        val transport = FakeAdministrationTransport.withResponse(
            RemoveConsumerResponse.newBuilder().setConsumerId("c-1").setRemoved(false).build(),
        )
        val client = AdministrationTestClientFactory.create(transport)

        val response = client.removeConsumer(RemoveConsumerRequest.newBuilder().setConsumerId("c-1").build())

        assertFalse(response.removed)
    }
}
