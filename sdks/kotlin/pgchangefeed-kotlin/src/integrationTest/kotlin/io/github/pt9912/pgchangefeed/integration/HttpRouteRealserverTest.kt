package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient
import java.net.URI
import kotlin.test.Test

/**
 * Real-server routing phase for the HTTP API client: `readChanges` with
 * `target` returns exactly the changes routed to that target (all of region
 * A, and every one the unfiltered read assigns to A), while the call without
 * `target` returns all of them.
 */
class HttpRouteRealserverTest {
    @Test
    fun readWithTargetReturnsOnlyItsTargetAndReadWithoutTargetReturnsAll() {
        val httpClient = java.net.http.HttpClient.newHttpClient()
        val options = PgChangeFeedClientOptions(URI(PhaseEnvironment.httpAddr), PhaseEnvironment.apiToken)
        val client = PgChangeFeedHttpClient(httpClient, options)

        RouteScenario.runPull { target ->
            client.readChanges(PhaseEnvironment.sourceId, "public", PhaseEnvironment.table, target = target)
                .changes
                .map { change -> RouteRow.from(change.changeId, change.table, change.newImage?.toString()) }
        }
    }
}
