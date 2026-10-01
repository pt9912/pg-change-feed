package io.github.pt9912.pgchangefeed.integration

import com.google.gson.JsonParser
import java.util.concurrent.CopyOnWriteArrayList
import kotlin.concurrent.thread
import kotlin.test.assertTrue

/** One received change reduced to what the routing phases assert on. */
data class RouteRow(val changeId: String, val table: String, val region: String?, val name: String?) {
    companion object {
        /** Reads `region` and `name` from the JSON text of a row image. */
        fun from(changeId: String, table: String, imageJson: String?): RouteRow {
            if (imageJson.isNullOrEmpty()) {
                return RouteRow(changeId, table, null, null)
            }
            val image = JsonParser.parseString(imageJson).asJsonObject
            fun text(key: String): String? =
                image.get(key)?.takeIf { it.isJsonPrimitive }?.asString
            return RouteRow(changeId, table, text("region"), text("name"))
        }
    }
}

/** Collects the rows of one stream on a daemon thread. */
class RouteCollector {
    private val collected = CopyOnWriteArrayList<RouteRow>()

    @Volatile
    private var stopped = false

    @Volatile
    var failure: Throwable? = null
        private set

    val rows: List<RouteRow> get() = collected.toList()

    /**
     * Runs [consume] until the stream ends or [stop] is called; it hands every
     * received row to the callback. A fault before the stop is kept in
     * [failure]. The read blocks while the stream is idle, so the thread is a
     * daemon: closing the client after the scenario ends the read.
     */
    fun start(consume: (emit: (RouteRow) -> Unit) -> Unit) {
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
 * The shared course of the routing phases. A client with the target
 * `routeTargetA` and a client without a target read the same table; the
 * table carries two routing rules (region A to target A, region B to target
 * B), and the runner commits groups of three changes (no rule, B, A). The
 * client with the target must receive the change for A and no change of any
 * other region; the client without a target must receive all three kinds.
 */
object RouteScenario {
    private const val POSITIVE_DEADLINE_MILLIS = 90_000L

    private fun own(rows: List<RouteRow>): List<RouteRow> =
        rows.filter { it.table == PhaseEnvironment.table && it.name == PhaseEnvironment.sentinel }

    private fun hasAllThree(rows: List<RouteRow>): Boolean {
        val regions = own(rows).map { it.region }.toSet()
        return PhaseEnvironment.routeTargetA in regions &&
            PhaseEnvironment.routeTargetB in regions &&
            PhaseEnvironment.routeRegionNone in regions
    }

    private fun throwOnFailure(vararg collectors: RouteCollector) {
        collectors.forEach { collector ->
            collector.failure?.let { throw IllegalStateException("Ein Stream-Konsument ist ausgefallen.", it) }
        }
    }

    /** Stream surfaces: positive phase, quiet window, then the checks. */
    fun runStreams(targeted: RouteCollector, unfiltered: RouteCollector) {
        println("READY")
        System.out.flush()
        val deadline = System.currentTimeMillis() + POSITIVE_DEADLINE_MILLIS
        while (!(own(targeted.rows).any { it.region == PhaseEnvironment.routeTargetA } && hasAllThree(unfiltered.rows))) {
            throwOnFailure(targeted, unfiltered)
            assertTrue(
                System.currentTimeMillis() < deadline,
                "innerhalb der Frist weder die Change des Ziels ${PhaseEnvironment.routeTargetA} am Client mit Ziel " +
                    "noch alle drei Gruppen am Client ohne Ziel empfangen",
            )
            Thread.sleep(100)
        }
        println("SEEN")
        System.out.flush()

        // The window starts only after the client without a target received the
        // change of the other target: the same delivery path has then
        // demonstrably dispatched it, so its absence at the client with a target
        // is no early cut-off.
        Thread.sleep(PhaseEnvironment.routeQuietSeconds * 1000)
        throwOnFailure(targeted, unfiltered)
        targeted.stop()
        unfiltered.stop()

        evaluate(targeted.rows, unfiltered.rows, PhaseEnvironment.routeQuietSeconds)
    }

    /** Pull surface: the target read equals the region-A part of the unfiltered read. */
    fun runPull(read: (String?) -> List<RouteRow>) {
        println("READY")
        System.out.flush()
        val deadline = System.currentTimeMillis() + POSITIVE_DEADLINE_MILLIS
        while (!hasAllThree(read(null))) {
            assertTrue(
                System.currentTimeMillis() < deadline,
                "innerhalb der Frist nicht alle drei Gruppen über die ungefilterte Lesung gelesen",
            )
            Thread.sleep(300)
        }
        println("SEEN")
        System.out.flush()

        val before = read(null)
        val targeted = read(PhaseEnvironment.routeTargetA)
        val after = read(null)

        fun regionA(rows: List<RouteRow>): Set<String> =
            own(rows).filter { it.region == PhaseEnvironment.routeTargetA }.map { it.changeId }.toSet()
        val targetedOwn = own(targeted).map { it.changeId }.toSet()
        assertTrue(
            targetedOwn.containsAll(regionA(before)),
            "die Lesung mit Ziel enthält nicht jede Change des Ziels, die die Lesung davor ohne Ziel lieferte",
        )
        assertTrue(
            regionA(after).containsAll(targetedOwn),
            "die Lesung mit Ziel enthält eine Change, die die Lesung danach ohne Ziel nicht dem Ziel zuordnet",
        )

        evaluate(targeted, after, 0)
    }

    private fun evaluate(targeted: List<RouteRow>, unfiltered: List<RouteRow>, quietSeconds: Long) {
        val foreign = targeted.filter {
            it.table != PhaseEnvironment.table || it.region != PhaseEnvironment.routeTargetA
        }
        assertTrue(
            foreign.isEmpty(),
            "der Client mit Ziel ${PhaseEnvironment.routeTargetA} empfing fremde Changes: " +
                foreign.joinToString(", ") { "${it.changeId}(region=${it.region})" },
        )
        assertTrue(own(targeted).isNotEmpty(), "der Client mit Ziel empfing keine Change dieser Phase")
        assertTrue(hasAllThree(unfiltered), "der Client ohne Ziel sah nicht alle drei Gruppen")

        val ownUnfiltered = own(unfiltered)
        targeted.forEach { println("RECEIVED_TARGETED change_id=${it.changeId} table=${it.table} region=${it.region}") }
        ownUnfiltered.forEach {
            println("RECEIVED_UNFILTERED change_id=${it.changeId} table=${it.table} region=${it.region}")
        }
        println(
            "ROUTE_RESULT target=${PhaseEnvironment.routeTargetA} targeted=${targeted.size} foreign=0 " +
                "unfiltered=${ownUnfiltered.size} quiet_seconds=$quietSeconds",
        )
        System.out.flush()
    }
}
