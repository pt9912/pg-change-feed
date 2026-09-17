package cdcexamples.grpc

import cdc.stream.v1.Changestream.Change
import com.google.protobuf.ByteString
import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * Prüft [Format.formatChange] — reine Funktion, netzlos testbar. Form-Vorbild:
 * `examples/csharp/grpc-client/GrpcClient.Tests/FormatTests.cs` (`slice-102`).
 */
class FormatTest {
    @Test
    fun formatChangeBuildsExpectedLine() {
        val change = Change.newBuilder()
            .setChangeId("chg-1")
            .setSchema("public")
            .setTable("widgets")
            .setOperation("INSERT")
            .setNewImage(ByteString.copyFromUtf8("{\"id\":1}"))
            .build()

        val got = Format.formatChange(change)

        assertEquals(
            "grpc-client: change_id=chg-1 table=public.widgets operation=INSERT new_image={\"id\":1}",
            got,
        )
    }
}
