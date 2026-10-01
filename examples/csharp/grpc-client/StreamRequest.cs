using Cdc.Stream.V1;

namespace CdcExamples.Grpc;

/// <summary>
/// StreamRequest bildet die Anfrage des Streams aus der Konfiguration:
/// <c>Schema</c>/<c>Table</c>/<c>Target</c> tragen den optionalen, unabhängig
/// setzbaren Filter (<c>ADR-0133</c>) — alle leer liefert jeden Change aller
/// aktivierten Tabellen. Form-Vorbild: <c>streamRequest</c> in
/// <c>examples/grpc-client/stream.go</c>.
/// </summary>
public static class StreamRequest
{
    public static StreamChangesRequest Build(Config cfg) =>
        new() { Schema = cfg.Schema, Table = cfg.Table, Target = cfg.Target };
}
