package io.github.pt9912.pgchangefeed.integration

import com.google.gson.JsonObject
import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient
import java.net.URI
import java.util.concurrent.CopyOnWriteArrayList
import kotlin.concurrent.thread
import kotlin.test.Test
import kotlin.test.assertTrue

/** One received change reduced to what the filter phase asserts on. */
data class FilterRow(val changeId: String, val schema: String, val table: String, val name: String?)

/** Collects the rows of one stream on a daemon thread. */
class FilterCollector {
    private val collected = CopyOnWriteArrayList<FilterRow>()

    @Volatile
    private var stopped = false

    @Volatile
    var failure: Throwable? = null
        private set

    val rows: List<FilterRow> get() = collected.toList()

    /**
     * Runs [consume] until the stream ends or [stop] is called; it hands every
     * received row to the callback. A fault before the stop is kept in
     * [failure]. The read blocks while the stream is idle, so the thread is a
     * daemon.
     */
    fun start(consume: (emit: (FilterRow) -> Unit) -> Unit) {
        thread(isDaemon = true) {
            try {
                consume { row -> collected.add(row) }
            } catch (ex: Throwable) {
                if (!stopped) {
                    failure = ex
                }
            }
        }
    }

    fun stop() {
        stopped = true
    }
}

/**
 * Real-server filter phase for the SSE stream client. Three tables are
 * captured: `A` and `B` in the first schema and a table named like `A` in a
 * second schema. A stream opened with `schema` and `table` of `A` receives
 * exactly the changes of that table, a stream opened with the second `schema`
 * alone receives exactly the changes of that schema, and a stream without a
 * filter receives all three. The first group of three changes proves that
 * every stream is connected (`SEEN`); the runner then commits a second group
 * with its own sentinel, and the quiet window starts only after all three
 * streams hold their rows of that second group (`SEEN_SECOND`).
 */
class SseFilterRealserverTest {
    @Test
    fun streamsWithSchemaAndTableFilterReceiveOnlyTheirSelectionAndStreamWithoutFilterReceivesAll() {
        val httpClient = java.net.http.HttpClient.newHttpClient()
        val options = PgChangeFeedClientOptions(URI(PhaseEnvironment.httpAddr), PhaseEnvironment.apiToken)
        val client = PgChangeFeedSseClient(httpClient, options)
        val f1 = FilterCollector()
        val f2 = FilterCollector()
        val unfiltered = FilterCollector()

        fun row(change: io.github.pt9912.pgchangefeed.sse.model.Change): FilterRow {
            val image = change.newImage
            val name = (image as? JsonObject)?.get("name")?.takeIf { it.isJsonPrimitive }?.asString
            return FilterRow(change.changeId, change.schema, change.table, name)
        }
        f1.start { emit ->
            client.streamChanges(schema = PhaseEnvironment.filterSchemaA, table = PhaseEnvironment.filterTableA)
                .forEach { emit(row(it)) }
        }
        f2.start { emit ->
            client.streamChanges(schema = PhaseEnvironment.filterSchemaOther).forEach { emit(row(it)) }
        }
        unfiltered.start { emit -> client.streamChanges().forEach { emit(row(it)) } }

        println("READY")
        System.out.flush()
        val deadline = System.currentTimeMillis() + POSITIVE_DEADLINE_MILLIS
        val first = PhaseEnvironment.sentinel
        while (!(own(f1.rows, first).any { isTableA(it) } && own(f2.rows, first).any { isOtherSchema(it) } && hasAllThree(unfiltered.rows, PhaseEnvironment.sentinel))) {
            throwOnFailure(f1, f2, unfiltered)
            assertTrue(
                System.currentTimeMillis() < deadline,
                "innerhalb der Frist weder die Change der Tabelle am Client mit Schema und Tabelle, " +
                    "noch die des zweiten Schemas am Client mit Schema, noch alle drei Tabellen am Client ohne Filter empfangen",
            )
            Thread.sleep(100)
        }
        println("SEEN")
        System.out.flush()

        // Every stream received a change, so every connection stands. The
        // second group is committed after that point and reaches all three
        // streams; the window starts once each holds its rows of that group,
        // so the absence of foreign changes at the filtered clients is no
        // early cut-off.
        val second = PhaseEnvironment.filterSentinelSecond
        val secondDeadline = System.currentTimeMillis() + POSITIVE_DEADLINE_MILLIS
        while (!(own(f1.rows, second).any { isTableA(it) } && own(f2.rows, second).any { isOtherSchema(it) } && hasAllThree(unfiltered.rows, second))) {
            throwOnFailure(f1, f2, unfiltered)
            assertTrue(
                System.currentTimeMillis() < secondDeadline,
                "innerhalb der Frist weder die Change der Tabelle der zweiten Gruppe am Client mit Schema und Tabelle, " +
                    "noch die des zweiten Schemas am Client mit Schema, noch alle drei der zweiten Gruppe am Client ohne Filter empfangen",
            )
            Thread.sleep(100)
        }
        println("SEEN_SECOND")
        System.out.flush()

        Thread.sleep(PhaseEnvironment.filterQuietSeconds * 1000)
        throwOnFailure(f1, f2, unfiltered)
        f1.stop()
        f2.stop()
        unfiltered.stop()

        val f1Rows = f1.rows
        val f2Rows = f2.rows
        val ownUnfiltered = ownAny(unfiltered.rows)
        val f1Foreign = f1Rows.filterNot { isTableA(it) }
        val f2Foreign = f2Rows.filterNot { isOtherSchema(it) }
        f1Rows.forEach { println("RECEIVED_F1 change_id=${it.changeId} schema=${it.schema} table=${it.table}") }
        f2Rows.forEach { println("RECEIVED_F2 change_id=${it.changeId} schema=${it.schema} table=${it.table}") }
        ownUnfiltered.forEach { println("RECEIVED_U change_id=${it.changeId} schema=${it.schema} table=${it.table}") }
        // The line carries the measured counts of foreign changes before the
        // checks below, so a foreign receipt is visible in the line itself.
        println(
            "FILTER_RESULT f1=${f1Rows.size} f1_foreign=${f1Foreign.size} f2=${f2Rows.size} " +
                "f2_foreign=${f2Foreign.size} unfiltered=${ownUnfiltered.size} " +
                "quiet_seconds=${PhaseEnvironment.filterQuietSeconds}",
        )
        System.out.flush()

        assertTrue(
            f1Foreign.isEmpty(),
            "der Client mit Schema und Tabelle empfing fremde Changes: " +
                f1Foreign.joinToString(", ") { "${it.changeId}(${it.schema}.${it.table})" },
        )
        assertTrue(
            f2Foreign.isEmpty(),
            "der Client mit Schema allein empfing fremde Changes: " +
                f2Foreign.joinToString(", ") { "${it.changeId}(${it.schema}.${it.table})" },
        )
        assertTrue(ownAny(f1Rows).isNotEmpty(), "der Client mit Schema und Tabelle empfing keine Change dieser Phase")
        assertTrue(ownAny(f2Rows).isNotEmpty(), "der Client mit Schema allein empfing keine Change dieser Phase")
        assertTrue(
            hasAllThree(unfiltered.rows, PhaseEnvironment.sentinel),
            "der Client ohne Filter sah nicht alle drei Tabellen der ersten Gruppe",
        )
        assertTrue(hasAllThree(unfiltered.rows, second), "der Client ohne Filter sah nicht alle drei Tabellen der zweiten Gruppe")
    }

    private companion object {
        const val POSITIVE_DEADLINE_MILLIS = 90_000L

        fun own(rows: List<FilterRow>, sentinel: String): List<FilterRow> = rows.filter { it.name == sentinel }

        fun ownAny(rows: List<FilterRow>): List<FilterRow> =
            rows.filter { it.name == PhaseEnvironment.sentinel || it.name == PhaseEnvironment.filterSentinelSecond }

        fun isTableA(r: FilterRow): Boolean =
            r.schema == PhaseEnvironment.filterSchemaA && r.table == PhaseEnvironment.filterTableA

        fun isTableB(r: FilterRow): Boolean =
            r.schema == PhaseEnvironment.filterSchemaA && r.table == PhaseEnvironment.filterTableB

        fun isOtherSchema(r: FilterRow): Boolean = r.schema == PhaseEnvironment.filterSchemaOther

        fun hasAllThree(rows: List<FilterRow>, sentinel: String): Boolean {
            val mine = own(rows, sentinel)
            return mine.any { isTableA(it) } && mine.any { isTableB(it) } && mine.any { isOtherSchema(it) }
        }

        fun throwOnFailure(vararg collectors: FilterCollector) {
            collectors.forEach { collector ->
                collector.failure?.let { throw IllegalStateException("Ein Stream-Konsument ist ausgefallen.", it) }
            }
        }
    }
}
