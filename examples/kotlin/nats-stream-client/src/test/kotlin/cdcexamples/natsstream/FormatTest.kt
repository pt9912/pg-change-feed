package cdcexamples.natsstream

import com.google.gson.JsonParser
import kotlin.test.Test
import kotlin.test.assertEquals

/**
 * Prüft [Format.formatChange] — reine Funktion, netzlos testbar, inklusive
 * des `DELETE`-Falls, in dem das Row Image auf dem Draht das JSON-Literal
 * `null` trägt (dieselbe Übersetzung wie `natsstream.rowImage` auf der
 * Erzeugerseite). Form-Vorbild:
 * `examples/kotlin/grpc-client/src/test/kotlin/cdcexamples/grpc/FormatTest.kt`.
 */
class FormatTest {
    private fun json(raw: String) = JsonParser.parseString(raw)

    @Test
    fun formatChangeCarriesIdentityAndPayload() {
        val change = StreamMessage(
            changeId = "c-1",
            transactionId = "t-1",
            sourceTableId = "st-1",
            sequence = 1,
            operation = "INSERT",
            oldImage = json("null"),
            newImage = json("{\"id\":1}"),
            schemaVersion = "v1",
            schema = "public",
            table = "orders",
        )

        val got = Format.formatChange(change)

        assertEquals(
            "nats-stream-client: change_id=c-1 table=public.orders operation=INSERT new_image={\"id\":1}",
            got,
        )
    }

    @Test
    fun formatChangeHandlesMissingNewImage() {
        val change = StreamMessage(
            changeId = "c-2",
            transactionId = "t-1",
            sourceTableId = "st-1",
            sequence = 2,
            operation = "DELETE",
            oldImage = json("{\"id\":1}"),
            newImage = json("null"),
            schemaVersion = "v1",
            schema = "public",
            table = "orders",
        )

        val got = Format.formatChange(change)

        assertEquals(
            "nats-stream-client: change_id=c-2 table=public.orders operation=DELETE new_image=null",
            got,
        )
    }
}
