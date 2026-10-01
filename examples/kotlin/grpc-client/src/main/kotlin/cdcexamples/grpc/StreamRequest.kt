package cdcexamples.grpc

import cdc.stream.v1.Changestream.StreamChangesRequest

/**
 * StreamRequest bildet die Anfrage des Streams aus der Konfiguration:
 * `schema`/`table`/`target` tragen den optionalen, unabhängig setzbaren
 * Filter (`ADR-0133`) — alle leer liefert jeden Change aller aktivierten
 * Tabellen. Form-Vorbild: `streamRequest` in
 * `examples/grpc-client/stream.go`.
 */
object StreamRequest {
    fun build(cfg: Config): StreamChangesRequest =
        StreamChangesRequest.newBuilder()
            .setSchema(cfg.schema).setTable(cfg.table).setTarget(cfg.target)
            .build()
}
