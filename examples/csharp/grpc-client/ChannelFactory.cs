using System.Net.Security;
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;
using Grpc.Net.Client;

namespace CdcExamples.Grpc;

/// <summary>
/// ChannelFactory baut den gRPC-Kanal: nennt <c>Config.CaFile</c> ein
/// PEM-Zertifikat, spricht der Kanal TLS und vertraut ausschließlich diesem
/// Zertifikat als Anker (Kette und Name werden geprüft, kein Überspringen der
/// Prüfung); ohne Angabe spricht er Klartext (HTTP/2 ohne Transportverschlüsselung).
/// Form-Vorbild: <c>transportCredentials</c> in <c>examples/grpc-client/main.go</c>.
/// </summary>
public static class ChannelFactory
{
    /// <summary>Adresse des Kanals: <c>https://</c> mit Anker, sonst <c>http://</c>.</summary>
    public static string Address(Config cfg) =>
        $"{(cfg.CaFile.Length > 0 ? "https" : "http")}://{cfg.Addr}";

    /// <summary>
    /// Create liefert den Kanal zur Konfiguration. Eine nicht lesbare Datei und
    /// eine Datei ohne PEM-Zertifikat enden als <see cref="ArgumentException"/>
    /// vor dem Verbindungsaufbau.
    /// </summary>
    public static GrpcChannel Create(Config cfg)
    {
        if (cfg.CaFile.Length == 0)
        {
            // HTTP/2 ohne Transportverschlüsselung (h2c) für den `HttpClient`,
            // den `GrpcChannel.ForAddress` intern benutzt.
            AppContext.SetSwitch("System.Net.Http.SocketsHttpHandler.Http2UnencryptedSupport", true);
            return GrpcChannel.ForAddress(Address(cfg));
        }

        return GrpcChannel.ForAddress(
            Address(cfg),
            new GrpcChannelOptions { HttpHandler = TlsHandler(cfg.CaFile) });
    }

    /// <summary>
    /// TlsHandler liest die Anker aus <paramref name="caFile"/> und liefert den
    /// Handler, dessen Zertifikatsprüfung nur diese Anker kennt.
    /// </summary>
    public static SocketsHttpHandler TlsHandler(string caFile)
    {
        var anchors = new X509Certificate2Collection();
        try
        {
            anchors.ImportFromPem(File.ReadAllText(caFile));
        }
        catch (Exception ex) when (ex is IOException or UnauthorizedAccessException)
        {
            throw new ArgumentException($"grpc-client: Vertrauensanker {caFile} nicht lesbar: {ex.Message}");
        }
        catch (CryptographicException)
        {
            anchors.Clear();
        }

        if (anchors.Count == 0)
        {
            throw new ArgumentException($"grpc-client: Vertrauensanker {caFile} enthält kein PEM-Zertifikat");
        }

        var handler = new SocketsHttpHandler();
        handler.SslOptions.RemoteCertificateValidationCallback =
            (_, certificate, _, errors) => Validate(anchors, certificate, errors);
        return handler;
    }

    private static bool Validate(X509Certificate2Collection anchors, X509Certificate? certificate, SslPolicyErrors errors)
    {
        if (certificate is null
            || (errors & (SslPolicyErrors.RemoteCertificateNotAvailable | SslPolicyErrors.RemoteCertificateNameMismatch)) != 0)
        {
            return false;
        }

        using var chain = new X509Chain();
        chain.ChainPolicy.TrustMode = X509ChainTrustMode.CustomRootTrust;
        chain.ChainPolicy.CustomTrustStore.AddRange(anchors);
        chain.ChainPolicy.RevocationMode = X509RevocationMode.NoCheck;
        using var leaf = new X509Certificate2(certificate);
        return chain.Build(leaf);
    }
}
