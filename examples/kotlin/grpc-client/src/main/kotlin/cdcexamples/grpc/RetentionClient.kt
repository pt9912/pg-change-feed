package cdcexamples.grpc

import cdc.administration.v1.AdministrationGrpcKt
import cdc.administration.v1.AdministrationOuterClass.RunRetentionRequest

/**
 * RetentionClient ruft die admin-RPC `RunRetention` des
 * `Administration`-Diensts auf. Form-Vorbild:
 * `examples/csharp/grpc-client/RetentionClient.cs`,
 * `examples/grpc-client/retention.go`.
 */
object RetentionClient {
    /** runRetention ruft `RunRetention` auf (`LH-FA-RET-002`): `minAgeNanos` 0 heißt „kein zeitliches Mindestalter" und ist gültig. */
    suspend fun runRetention(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String {
        val resp = withCallTimeout(client).runRetention(
            RunRetentionRequest.newBuilder().setSource(cfg.source).setMinAgeNanos(cfg.minAgeNanos).build(),
            CallMetadata.headers(cfg.adminToken),
        )
        return Format.formatRunRetention(resp)
    }
}
