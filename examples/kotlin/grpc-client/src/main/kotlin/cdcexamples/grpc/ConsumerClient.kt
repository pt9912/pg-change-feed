package cdcexamples.grpc

import cdc.administration.v1.AdministrationGrpcKt
import cdc.administration.v1.AdministrationOuterClass.AcknowledgeConsumerRequest
import cdc.administration.v1.AdministrationOuterClass.GetConsumerPositionRequest
import cdc.administration.v1.AdministrationOuterClass.RegisterConsumerRequest
import cdc.administration.v1.AdministrationOuterClass.RemoveConsumerRequest

/**
 * ConsumerClient ruft die vier Consumer-Verwaltungs-RPCs des
 * `Administration`-Diensts auf. Form-Vorbild:
 * `examples/csharp/grpc-client/ConsumerClient.cs`,
 * `examples/grpc-client/consumer.go`.
 */
object ConsumerClient {
    /** registerConsumer ruft die admin-RPC `RegisterConsumer` auf (`LH-FA-CON-001`). */
    suspend fun registerConsumer(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String {
        val resp = withCallTimeout(client).registerConsumer(
            RegisterConsumerRequest.newBuilder().setConsumerId(cfg.consumerId).setName(cfg.name).build(),
            CallMetadata.headers(cfg.adminToken),
        )
        return Format.formatRegisterConsumer(resp)
    }

    /** acknowledgeConsumer ruft die admin-RPC `AcknowledgeConsumer` auf (`LH-FA-CON-004`). */
    suspend fun acknowledgeConsumer(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String {
        val resp = withCallTimeout(client).acknowledgeConsumer(
            AcknowledgeConsumerRequest.newBuilder()
                .setConsumerId(cfg.consumerId).setSourceId(cfg.source).setOffset(cfg.offset)
                .build(),
            CallMetadata.headers(cfg.adminToken),
        )
        return Format.formatAcknowledgeConsumer(resp)
    }

    /** getConsumerPosition ruft die reader-RPC `GetConsumerPosition` auf (`LH-FA-CON-005`). */
    suspend fun getConsumerPosition(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String {
        val resp = withCallTimeout(client).getConsumerPosition(
            GetConsumerPositionRequest.newBuilder().setConsumerId(cfg.consumerId).build(),
            CallMetadata.headers(cfg.token),
        )
        return Format.formatConsumerPosition(resp)
    }

    /** removeConsumer ruft die admin-RPC `RemoveConsumer` auf (`LH-FA-CON-006`). */
    suspend fun removeConsumer(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String {
        val resp = withCallTimeout(client).removeConsumer(
            RemoveConsumerRequest.newBuilder().setConsumerId(cfg.consumerId).build(),
            CallMetadata.headers(cfg.adminToken),
        )
        return Format.formatRemoveConsumer(resp)
    }
}
