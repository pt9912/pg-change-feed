package io.github.pt9912.pgchangefeed.integration

import io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions
import io.grpc.StatusException
import java.net.URI
import java.nio.file.Path
import javax.net.ssl.SSLException
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

/**
 * The three TLS refusals every TLS phase proves against the running feed
 * container: no trust anchor (the certificate of the server is unknown to the
 * Java runtime), a foreign trust anchor, and a server name that is not in the
 * certificate (the container name, which resolves in the Docker network but is
 * not listed in the certificate; the anchor is the right one). Each must end
 * with a failed TLS check — a refusal for any other reason (an unreachable
 * name, a refused port) does not count.
 */
object TlsScenarios {
    fun withMismatchHost(address: URI): URI =
        URI(address.scheme, address.userInfo, PhaseEnvironment.tlsMismatchHost, address.port, address.path, address.query, address.fragment)

    fun anchor(file: String): Path = Path.of(file)

    /**
     * Runs [call] with the three option sets and asserts a TLS check failure for each; prints the runner marker afterwards.
     */
    fun expectRefusals(address: URI, token: String, call: (PgChangeFeedClientOptions) -> Unit) {
        val cases = listOf(
            "no_anchor" to PgChangeFeedClientOptions(address, token),
            "foreign_anchor" to PgChangeFeedClientOptions(address, token, anchor(PhaseEnvironment.tlsForeignCaFile)),
            "name_mismatch" to PgChangeFeedClientOptions(withMismatchHost(address), token, anchor(PhaseEnvironment.tlsCaFile)),
        )
        for ((name, options) in cases) {
            val error = runCatching { call(options) }.exceptionOrNull()
            assertNotNull(error, "$name: the call succeeded")
            assertTrue(isVerificationFailure(error), "$name: not a TLS check failure: $error")
        }
        println("REJECTED tls no_anchor=ok foreign_anchor=ok name_mismatch=ok")
        System.out.flush()
    }

    /** True when the cause chain holds the `SSLException` of a failed TLS check (the gRPC transport reports it as an `UNAVAILABLE` status with the description `ssl exception`). */
    fun isVerificationFailure(error: Throwable?): Boolean {
        var current = error
        while (current != null) {
            if (current is SSLException) {
                return true
            }
            if (current is StatusException && current.status.description?.contains("ssl", ignoreCase = true) == true) {
                return true
            }
            current = current.cause
        }
        return false
    }
}
