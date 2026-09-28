package cdcexamples.grpc

import cdc.administration.v1.AdministrationGrpcKt
import cdc.administration.v1.AdministrationOuterClass.ReadChangesRequest

/**
 * ChangesClient ruft die reader-RPC `ReadChanges` des
 * `Administration`-Diensts auf. Form-Vorbild:
 * `examples/csharp/grpc-client/ChangesClient.cs`,
 * `examples/grpc-client/changes.go`.
 */
object ChangesClient {
    /**
     * readChanges ruft `ReadChanges` auf (`LH-FA-SST-006`): einen begrenzten
     * Bereich persistierter Änderungen, dieselbe Filter- und
     * Bereichs-Semantik wie `GET /changes` — ein leerer
     * `from`/`to`/`limit`-Wert trägt `0` (nicht gesetzt).
     */
    suspend fun readChanges(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String {
        val resp = withCallTimeout(client).readChanges(
            ReadChangesRequest.newBuilder()
                .setSource(cfg.source).setSchema(cfg.schema).setTable(cfg.table)
                .setFrom(cfg.from).setTo(cfg.to).setLimit(cfg.limit)
                .build(),
            CallMetadata.headers(cfg.token),
        )
        return Format.formatReadChanges(resp)
    }
}
