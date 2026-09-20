package io.github.pt9912.pgchangefeed.grpc

import cdc.stream.v1.Changestream.Change
import io.grpc.Metadata
import io.grpc.Status
import io.grpc.StatusException
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow

/**
 * A network-free [GrpcStreamTransport]: it never opens an [io.grpc.Channel]
 * or touches a socket — see [PgChangeFeedGrpcClient]'s test-only
 * constructor KDoc for why `java.net.http`-style socket-free faking is not
 * available at the [io.grpc.Channel] level for the generated coroutine
 * stub. [lastHeaders] lets a test assert on the exact `authorization`
 * metadata entry the client under test built.
 */
internal class FakeGrpcStreamTransport private constructor(
    private val messages: List<Change>,
    private val failureStatus: Status?,
) : GrpcStreamTransport {

    var lastHeaders: Metadata? = null
        private set

    override fun streamChanges(headers: Metadata): Flow<Change> {
        lastHeaders = headers
        return flow {
            messages.forEach { emit(it) }
            if (failureStatus != null) {
                throw StatusException(failureStatus)
            }
        }
    }

    companion object {
        /**
         * Builds a fake transport whose stream yields the given [messages]
         * in order, then completes cleanly — the happy path (`SPEC-020`:
         * one message per row change, in order).
         */
        fun withMessages(vararg messages: Change): FakeGrpcStreamTransport =
            FakeGrpcStreamTransport(messages.toList(), failureStatus = null)

        /**
         * Builds a fake transport whose stream fails immediately with
         * [status], before any message — simulating what a real gRPC
         * channel does when a call is rejected up front (e.g. `SPEC-020`'s
         * `UNAUTHENTICATED` auth boundary). Ausgang bewusst offen
         * (Slice-Plan §6): dieses Fake bildet die Fehler-*Form* nach (ein
         * `StatusException` aus dem Stream selbst), nicht den realen
         * Server-seitigen Ablehnungsmechanismus — ein realer
         * Rundlauf-Beleg bleibt `make test-integration`s
         * `tools/harness/grpcclient` vorbehalten.
         */
        fun withStatus(status: Status): FakeGrpcStreamTransport =
            FakeGrpcStreamTransport(emptyList(), status)
    }
}
