package cdcexamples.grpc

import cdc.administration.v1.AdministrationGrpcKt
import cdc.administration.v1.AdministrationOuterClass.DiagnoseRequest

/**
 * DiagnoseClient ruft die reader-RPC `Diagnose` des `Administration`-Diensts
 * auf. Form-Vorbild: `examples/csharp/grpc-client/DiagnoseClient.cs`,
 * `examples/grpc-client/diagnose.go`.
 */
object DiagnoseClient {
    /** diagnose ruft `Diagnose` auf (`LH-FA-SST-003`): denselben Bericht wie der CLI-Sondermodus und `GET /diagnose`. */
    suspend fun diagnose(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String {
        val resp = withCallTimeout(client).diagnose(
            DiagnoseRequest.newBuilder().setSource(cfg.source).build(),
            CallMetadata.headers(cfg.token),
        )
        return Format.formatDiagnose(resp)
    }
}
