package cdcexamples.grpc

import io.grpc.Grpc
import io.grpc.InsecureChannelCredentials
import io.grpc.ManagedChannel
import io.grpc.TlsChannelCredentials
import java.io.File
import java.security.cert.CertificateException
import java.security.cert.CertificateFactory

/**
 * ChannelFactory baut den gRPC-Kanal: nennt `Config.caFile` ein PEM-Zertifikat,
 * spricht der Kanal TLS und vertraut ausschließlich diesem Zertifikat als
 * Anker (Kette und Name werden geprüft, kein Überspringen der Prüfung); ohne
 * Angabe spricht er Klartext. Form-Vorbild: `transportCredentials` in
 * `examples/grpc-client/main.go`.
 */
object ChannelFactory {
    /**
     * create liefert den Kanal zur Konfiguration. Eine nicht lesbare Datei und
     * eine Datei ohne PEM-Zertifikat enden als [IllegalArgumentException] vor
     * dem Verbindungsaufbau; ein DER-Zertifikat gilt als Anker.
     */
    fun create(cfg: Config): ManagedChannel {
        val credentials = if (cfg.caFile.isEmpty()) {
            InsecureChannelCredentials.create()
        } else {
            val anchor = File(cfg.caFile)
            if (!anchor.isFile || !anchor.canRead()) {
                throw IllegalArgumentException("grpc-client: Vertrauensanker ${cfg.caFile} nicht lesbar")
            }
            val certificates = try {
                anchor.inputStream().use { CertificateFactory.getInstance("X.509").generateCertificates(it) }
            } catch (ex: CertificateException) {
                emptyList()
            }
            if (certificates.isEmpty()) {
                throw IllegalArgumentException("grpc-client: Vertrauensanker ${cfg.caFile} enthält kein PEM-Zertifikat")
            }
            TlsChannelCredentials.newBuilder().trustManager(anchor).build()
        }
        return Grpc.newChannelBuilder(cfg.addr, credentials).build()
    }
}
