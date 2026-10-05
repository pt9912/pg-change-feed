package io.github.pt9912.pgchangefeed

import io.grpc.ChannelCredentials
import io.grpc.Grpc
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import io.grpc.TlsChannelCredentials
import java.net.Socket
import java.net.http.HttpClient
import java.security.KeyStore
import java.security.cert.CertificateException
import java.security.cert.CertificateExpiredException
import java.security.cert.CertificateNotYetValidException
import java.security.cert.X509Certificate
import javax.net.ssl.SSLContext
import javax.net.ssl.SSLEngine
import javax.net.ssl.TrustManager
import javax.net.ssl.TrustManagerFactory
import javax.net.ssl.X509ExtendedTrustManager

/**
 * Builds the TLS pieces of the connections the clients open themselves.
 * Without a trust anchor the trust anchors of the Java runtime apply; with one,
 * the server certificate must chain to exactly the anchor certificates. The
 * server name (the host of the address) and the validity period are checked in
 * both cases, and nothing here can turn the check off.
 */
internal object TlsSupport {
    /** Fails when the options name a trust anchor but the address does not ask for TLS. */
    fun requireTlsAddressForAnchor(options: PgChangeFeedClientOptions) {
        require(options.trustAnchors == null || isHttps(options)) {
            "A trust anchor requires an https address, got '${options.address}'."
        }
    }

    fun isHttps(options: PgChangeFeedClientOptions): Boolean =
        options.address.scheme.equals("https", ignoreCase = true)

    /** The `java.net.http.HttpClient` for the HTTP and SSE clients that build their own connection. */
    fun httpClient(options: PgChangeFeedClientOptions): HttpClient {
        requireTlsAddressForAnchor(options)
        val anchors = options.trustAnchors ?: return HttpClient.newHttpClient()
        val context = SSLContext.getInstance("TLS")
        context.init(null, trustManagers(anchors), null)
        return HttpClient.newBuilder().sslContext(context).build()
    }

    /**
     * The channel the gRPC clients open and own against the options' address:
     * over TLS for an `https` address (the anchors when set, otherwise the
     * runtime default), in plaintext otherwise. The address must carry a host
     * and a port.
     */
    fun ownedChannel(options: PgChangeFeedClientOptions): ManagedChannel {
        val host = options.address.host
        val port = options.address.port
        require(host != null && port != -1) {
            "options.address must carry a resolvable host and port (got '${options.address}')"
        }
        requireTlsAddressForAnchor(options)
        if (isHttps(options)) {
            return Grpc.newChannelBuilderForAddress(host, port, channelCredentials(options)).build()
        }
        return ManagedChannelBuilder.forAddress(host, port).usePlaintext().build()
    }

    private fun channelCredentials(options: PgChangeFeedClientOptions): ChannelCredentials {
        val anchors = options.trustAnchors ?: return TlsChannelCredentials.create()
        return TlsChannelCredentials.newBuilder().trustManager(*trustManagers(anchors)).build()
    }

    private fun trustManagers(anchors: List<X509Certificate>): Array<TrustManager> {
        val store = KeyStore.getInstance(KeyStore.getDefaultType())
        store.load(null, null)
        anchors.forEachIndexed { index, certificate -> store.setCertificateEntry("anchor-$index", certificate) }
        val factory = TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm())
        factory.init(store)
        return factory.trustManagers.map { manager ->
            if (manager is X509ExtendedTrustManager) ValidityCheckingTrustManager(manager) else manager
        }.toTypedArray()
    }
}

/**
 * The trust manager of the Java runtime does not check the validity period of
 * a server certificate that is itself one of the anchors (a pinned server
 * certificate). This wrapper delegates the whole check and then requires the
 * presented server certificate to be inside its validity period, so an expired
 * certificate fails whether the anchor is the issuer or the server itself.
 */
private class ValidityCheckingTrustManager(private val delegate: X509ExtendedTrustManager) : X509ExtendedTrustManager() {
    override fun checkServerTrusted(chain: Array<out X509Certificate>, authType: String) {
        delegate.checkServerTrusted(chain, authType)
        checkValidity(chain)
    }

    override fun checkServerTrusted(chain: Array<out X509Certificate>, authType: String, socket: Socket) {
        delegate.checkServerTrusted(chain, authType, socket)
        checkValidity(chain)
    }

    override fun checkServerTrusted(chain: Array<out X509Certificate>, authType: String, engine: SSLEngine) {
        delegate.checkServerTrusted(chain, authType, engine)
        checkValidity(chain)
    }

    override fun checkClientTrusted(chain: Array<out X509Certificate>, authType: String) =
        delegate.checkClientTrusted(chain, authType)

    override fun checkClientTrusted(chain: Array<out X509Certificate>, authType: String, socket: Socket) =
        delegate.checkClientTrusted(chain, authType, socket)

    override fun checkClientTrusted(chain: Array<out X509Certificate>, authType: String, engine: SSLEngine) =
        delegate.checkClientTrusted(chain, authType, engine)

    override fun getAcceptedIssuers(): Array<X509Certificate> = delegate.acceptedIssuers

    private fun checkValidity(chain: Array<out X509Certificate>) {
        try {
            chain.firstOrNull()?.checkValidity()
        } catch (e: CertificateExpiredException) {
            throw CertificateException("The server certificate is outside its validity period: ${e.message}", e)
        } catch (e: CertificateNotYetValidException) {
            throw CertificateException("The server certificate is outside its validity period: ${e.message}", e)
        }
    }
}
