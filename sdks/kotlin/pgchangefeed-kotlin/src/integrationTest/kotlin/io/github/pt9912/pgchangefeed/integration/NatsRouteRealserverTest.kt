package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.nats.PgChangeFeedNatsStreamClient
import java.net.URI
import kotlin.test.Test

/**
 * Real-server routing phase for the NATS stream client: a subscription to the
 * subject of one delivery target receives the changes routed to that target
 * and no change of another target or without a target (checked over a quiet
 * window), while a subscription to the source's namespace receives all of
 * them.
 */
class NatsRouteRealserverTest {
    @Test
    fun targetSubjectReceivesOnlyItsTargetAndSourceSubjectReceivesAll() {
        val options = PgChangeFeedClientOptions(URI(PhaseEnvironment.natsUrl), PhaseEnvironment.natsStreamToken)
        PgChangeFeedNatsStreamClient(options).use { targetedClient ->
            PgChangeFeedNatsStreamClient(options).use { unfilteredClient ->
                val targetSubject = PgChangeFeedNatsStreamClient.buildTargetSubject(
                    PhaseEnvironment.sourceId,
                    PhaseEnvironment.routeTargetA,
                )
                val sourceSubject = PgChangeFeedNatsStreamClient.buildSourceSubject(PhaseEnvironment.sourceId)
                val targeted = RouteCollector()
                val unfiltered = RouteCollector()
                targeted.start { emit ->
                    targetedClient.streamChanges(targetSubject).forEach { change ->
                        emit(RouteRow.from(change.changeId, change.table, change.newImage?.toString()))
                    }
                }
                unfiltered.start { emit ->
                    unfilteredClient.streamChanges(sourceSubject).forEach { change ->
                        emit(RouteRow.from(change.changeId, change.table, change.newImage?.toString()))
                    }
                }

                RouteScenario.runStreams(targeted, unfiltered)
            }
        }
    }
}
