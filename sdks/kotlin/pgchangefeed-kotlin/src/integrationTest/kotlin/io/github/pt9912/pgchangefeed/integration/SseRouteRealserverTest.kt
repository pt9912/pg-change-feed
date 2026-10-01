package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient
import java.net.URI
import kotlin.test.Test

/**
 * Real-server routing phase for the SSE stream client: a stream opened with
 * `target` receives the changes routed to that target and no change of
 * another target or without a target (checked over a quiet window), while a
 * stream opened without `target` receives all of them.
 */
class SseRouteRealserverTest {
    @Test
    fun streamWithTargetReceivesOnlyItsTargetAndStreamWithoutTargetReceivesAll() {
        val httpClient = java.net.http.HttpClient.newHttpClient()
        val options = PgChangeFeedClientOptions(URI(PhaseEnvironment.httpAddr), PhaseEnvironment.apiToken)
        val client = PgChangeFeedSseClient(httpClient, options)
        val targeted = RouteCollector()
        val unfiltered = RouteCollector()
        targeted.start { emit ->
            client.streamChanges(target = PhaseEnvironment.routeTargetA).forEach { change ->
                emit(RouteRow.from(change.changeId, change.table, change.newImage?.toString()))
            }
        }
        unfiltered.start { emit ->
            client.streamChanges().forEach { change ->
                emit(RouteRow.from(change.changeId, change.table, change.newImage?.toString()))
            }
        }

        RouteScenario.runStreams(targeted, unfiltered)
    }
}
