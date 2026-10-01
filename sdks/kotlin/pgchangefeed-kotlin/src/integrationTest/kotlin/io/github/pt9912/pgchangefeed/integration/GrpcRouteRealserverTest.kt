package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient
import java.net.URI
import kotlin.test.Test
import kotlinx.coroutines.runBlocking

/**
 * Real-server routing phase for the gRPC stream client: a stream opened with
 * `target` receives the changes routed to that target and no change of
 * another target or without a target (checked over a quiet window), while a
 * stream opened without `target` receives all of them.
 */
class GrpcRouteRealserverTest {
    @Test
    fun streamWithTargetReceivesOnlyItsTargetAndStreamWithoutTargetReceivesAll() {
        val options = PgChangeFeedClientOptions(URI("http://${PhaseEnvironment.grpcAddr}"), PhaseEnvironment.apiToken)
        PgChangeFeedGrpcClient(options).use { targetedClient ->
            PgChangeFeedGrpcClient(options).use { unfilteredClient ->
                val targeted = RouteCollector()
                val unfiltered = RouteCollector()
                targeted.start { emit ->
                    runBlocking {
                        targetedClient.streamChanges(target = PhaseEnvironment.routeTargetA).collect { change ->
                            emit(RouteRow.from(change.changeId, change.table, change.newImage?.toStringUtf8()))
                        }
                    }
                }
                unfiltered.start { emit ->
                    runBlocking {
                        unfilteredClient.streamChanges().collect { change ->
                            emit(RouteRow.from(change.changeId, change.table, change.newImage?.toStringUtf8()))
                        }
                    }
                }

                RouteScenario.runStreams(targeted, unfiltered)
            }
        }
    }
}
