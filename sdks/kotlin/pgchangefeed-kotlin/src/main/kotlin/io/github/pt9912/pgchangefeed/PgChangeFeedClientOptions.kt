package io.github.pt9912.pgchangefeed

import java.io.IOException
import java.net.URI
import java.nio.file.Files
import java.nio.file.Path
import java.security.GeneralSecurityException
import java.security.cert.CertificateFactory
import java.security.cert.X509Certificate

/**
 * Connection configuration shared by the PG Change Feed clients: the address
 * of the server, the bearer token and, optionally, a trust anchor for TLS.
 * The HTTP, gRPC, SSE and NATS clients all authenticate with such a token
 * against a single server address. What a client does with them (which HTTP
 * paths, which gRPC call) is up to the client itself; the options carry no
 * retry or backoff policy.
 *
 * @property address the address of the PG Change Feed server, in the form the
 *   client needs: the HTTP base URL for the HTTP and SSE clients, the gRPC
 *   endpoint for the gRPC client, the NATS URL for the NATS client.
 * @property apiToken the bearer token sent as the authorization credential;
 *   must not be blank.
 * @property trustAnchorFile the path of a PEM file with the certificates a TLS
 *   connection trusts (the certificate of the issuer or of the server), or
 *   `null` when the trust anchors of the Java runtime apply. The HTTP, SSE and
 *   gRPC clients that build their own connection (the constructors that take
 *   only the options) use it and require an `https` address then. Chain,
 *   validity period and server name are always checked; there is no switch
 *   that turns the check off. A path that cannot be read, or a file without a
 *   PEM certificate, throws an [IllegalArgumentException] when the options are
 *   created.
 */
class PgChangeFeedClientOptions(
    val address: URI,
    val apiToken: String,
    val trustAnchorFile: Path?,
) {
    /** Options without a trust anchor: the trust anchors of the Java runtime apply to TLS connections. */
    constructor(address: URI, apiToken: String) : this(address, apiToken, null)

    internal val trustAnchors: List<X509Certificate>? = trustAnchorFile?.let { loadTrustAnchors(it) }

    init {
        require(apiToken.isNotBlank()) { "API token must not be blank." }
    }

    private companion object {
        fun loadTrustAnchors(path: Path): List<X509Certificate> {
            val certificates = try {
                Files.newInputStream(path).use { CertificateFactory.getInstance("X.509").generateCertificates(it) }
            } catch (e: IOException) {
                throw IllegalArgumentException("The trust anchor file '$path' cannot be read: ${e.message}", e)
            } catch (e: GeneralSecurityException) {
                throw IllegalArgumentException("The trust anchor file '$path' holds no PEM certificate: ${e.message}", e)
            }
            require(certificates.isNotEmpty()) { "The trust anchor file '$path' contains no PEM certificate." }
            return certificates.map { it as X509Certificate }
        }
    }
}
