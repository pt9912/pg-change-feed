using System.Text.Json;
using System.Text.Json.Serialization;

namespace CdcExamples.NatsStream;

/// <summary>
/// StreamMessage trägt dasselbe Nachrichtenschema wie der Publisher
/// (<c>internal/adapters/driven/natsstream</c>, <c>SPEC-024</c>) — hier
/// eigenständig geführt: dieses Beispiel importiert keinen privaten
/// Paketbaum dieses Repositories (<c>SPEC-023</c>).
/// </summary>
public sealed record StreamMessage(
    [property: JsonPropertyName("change_id")] string ChangeId,
    [property: JsonPropertyName("transaction_id")] string TransactionId,
    [property: JsonPropertyName("source_table_id")] string SourceTableId,
    [property: JsonPropertyName("sequence")] long Sequence,
    [property: JsonPropertyName("operation")] string Operation,
    [property: JsonPropertyName("old_image")] JsonElement OldImage,
    [property: JsonPropertyName("new_image")] JsonElement NewImage,
    [property: JsonPropertyName("schema_version")] string SchemaVersion,
    [property: JsonPropertyName("schema")] string Schema,
    [property: JsonPropertyName("table")] string Table);

/// <summary>
/// Format baut die Ausgabezeile einer empfangenen Stream-Nachricht
/// (<c>LH-FA-SST-008</c>): Tabelle, Operation und das neue Row Image, dazu
/// die Kennung, an der sich die Nachricht gegen den Lesezugriffsweg
/// <c>cdc.changes</c> halten lässt (<c>change_id</c>) — dieselbe Form wie
/// beim gRPC-Beispiel (<c>examples/csharp/grpc-client/Format.cs</c>).
/// </summary>
public static class Format
{
    public static string FormatChange(StreamMessage change) =>
        $"nats-stream-client: change_id={change.ChangeId} table={change.Schema}.{change.Table} " +
        $"operation={change.Operation} new_image={change.NewImage.GetRawText()}";
}
