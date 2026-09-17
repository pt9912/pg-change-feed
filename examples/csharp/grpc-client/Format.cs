using Cdc.Stream.V1;

namespace CdcExamples.Grpc;

/// <summary>
/// Format baut die Ausgabezeile einer empfangenen Stream-Nachricht
/// (<c>LH-FA-SST-008</c>) — reine Funktion, netzlos testbar. Form-Vorbild:
/// <c>examples/grpc-client/format.go</c>.
/// </summary>
public static class Format
{
    /// <summary>
    /// FormatChange baut Tabelle, Operation und das neue Row Image, dazu die
    /// Kennung, an der sich die Nachricht gegen den Lesezugriffsweg
    /// <c>cdc.changes</c> halten lässt (<c>change_id</c>).
    /// </summary>
    public static string FormatChange(Change change)
    {
        return $"grpc-client: change_id={change.ChangeId} table={change.Schema}.{change.Table} " +
               $"operation={change.Operation} new_image={change.NewImage.ToStringUtf8()}";
    }
}
