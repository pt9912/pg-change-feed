package cdcexamples.nats

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotEquals

/**
 * Prüft die Ableitung des Wecksignal-Subjekts aus den drei Bestandteilen
 * (`SPEC-017`) — netzlos, reine Funktion. Form-Vorbild:
 * `examples/nats-client/subject_test.go`.
 */
class SubjectTest {
    @Test
    fun buildDerivesTableGranularSubject() {
        val got = Subject.build("quelle-1", "public", "orders")
        assertEquals("cdc.changes.quelle-1.public.orders", got)
    }

    @Test
    fun buildKeepsDistinctTables() {
        val orders = Subject.build("quelle-1", "public", "orders")
        val invoices = Subject.build("quelle-1", "public", "invoices")
        assertNotEquals(orders, invoices)
    }
}
