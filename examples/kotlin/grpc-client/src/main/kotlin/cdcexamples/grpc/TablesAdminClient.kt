package cdcexamples.grpc

import cdc.administration.v1.AdministrationGrpcKt
import cdc.administration.v1.AdministrationOuterClass.DisableTableRequest
import cdc.administration.v1.AdministrationOuterClass.EnableTableRequest
import cdc.administration.v1.AdministrationOuterClass.GetTableStatusRequest
import cdc.administration.v1.AdministrationOuterClass.ListTablesRequest

/**
 * TablesAdminClient ruft die vier Tabellen-Verwaltungs-RPCs des
 * `Administration`-Diensts auf. Form-Vorbild:
 * `examples/csharp/grpc-client/TablesAdminClient.cs`,
 * `examples/grpc-client/tables_admin.go`.
 */
object TablesAdminClient {
    /** enableTable ruft die admin-RPC `EnableTable` auf (`LH-FA-CFG-001`). */
    suspend fun enableTable(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String {
        val resp = withCallTimeout(client).enableTable(
            EnableTableRequest.newBuilder()
                .setSource(cfg.source).setSchema(cfg.schema).setTable(cfg.table)
                .setTableId(cfg.tableId).setSchemaVersionId(cfg.schemaVersionId)
                .setVersion(cfg.version).setPublication(cfg.publication)
                .build(),
            CallMetadata.headers(cfg.adminToken),
        )
        return Format.formatEnableTable(resp)
    }

    /** disableTable ruft die admin-RPC `DisableTable` auf (`LH-FA-CFG-002`). */
    suspend fun disableTable(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String {
        val resp = withCallTimeout(client).disableTable(
            DisableTableRequest.newBuilder()
                .setSource(cfg.source).setSchema(cfg.schema).setTable(cfg.table).setPublication(cfg.publication)
                .build(),
            CallMetadata.headers(cfg.adminToken),
        )
        return Format.formatDisableTable(resp)
    }

    /** getTableStatus ruft die reader-RPC `GetTableStatus` auf (`LH-FA-CFG-003`). */
    suspend fun getTableStatus(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String {
        val resp = withCallTimeout(client).getTableStatus(
            GetTableStatusRequest.newBuilder()
                .setSource(cfg.source).setSchema(cfg.schema).setTable(cfg.table).setPublication(cfg.publication)
                .build(),
            CallMetadata.headers(cfg.token),
        )
        return Format.formatTableStatus(resp)
    }

    /**
     * listTables ruft die reader-RPC `ListTables` auf (`LH-FA-CFG-004`) und
     * listet jede Tabelle der Antwort — `tables` und `retained` je eine
     * `SourceTable`-Zeile.
     */
    suspend fun listTables(client: AdministrationGrpcKt.AdministrationCoroutineStub, cfg: Config): String {
        val resp = withCallTimeout(client).listTables(
            ListTablesRequest.newBuilder().setSource(cfg.source).setPublication(cfg.publication).build(),
            CallMetadata.headers(cfg.token),
        )
        return Format.formatListTables(resp)
    }
}
