package cdcexamples.grpc

import cdc.stream.v1.Changestream.Change

/**
 * Format baut die Ausgabezeile einer empfangenen Stream-Nachricht
 * (`LH-FA-SST-008`) — reine Funktion, netzlos testbar. Form-Vorbild:
 * `examples/csharp/grpc-client/Format.cs` (`slice-102`),
 * `examples/grpc-client/format.go`.
 *
 * `Change` liegt verschachtelt unter dem generierten Java-Deskriptor
 * `cdc.stream.v1.Changestream` — die `.proto` setzt kein
 * `option java_multiple_files`; anders als beim C#-Codegen (flache Klassen
 * im Namensraum `Cdc.Stream.V1`) verschachtelt der Java-/Kotlin-Codegen
 * Nachrichtentypen standardmäßig unter einer Datei-Außenklasse, deren Name
 * sich aus dem `.proto`-Dateinamen ableitet (`changestream.proto` ->
 * `Changestream`). Der Import macht das für den Aufrufer unsichtbar.
 */
object Format {
    fun formatChange(change: Change): String {
        return "grpc-client: change_id=${change.changeId} table=${change.schema}.${change.table} " +
            "operation=${change.operation} new_image=${change.newImage.toStringUtf8()}"
    }
}
