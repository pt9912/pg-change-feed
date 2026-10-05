package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// writeTestCert erzeugt zur Laufzeit ein selbstsigniertes Zertifikat für
// 127.0.0.1 und schreibt es als PEM in das Temp-Verzeichnis des Tests; der
// Pfad der Zertifikatsdatei und das geladene Paar werden zurückgegeben.
func writeTestCert(t *testing.T) (string, tls.Certificate) {
	t.Helper()
	return writeTestCertSAN(t, nil, []net.IP{net.ParseIP("127.0.0.1")})
}

// writeTestCertSAN erzeugt das Zertifikat mit genau den genannten SAN-Einträgen.
func writeTestCertSAN(t *testing.T, dns []string, ips []net.IP) (string, tls.Certificate) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Schlüssel: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "grpc-client-test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     dns,
		IPAddresses:  ips,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("Zertifikat: %v", err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	path := filepath.Join(t.TempDir(), "anchor.pem")
	if err := os.WriteFile(path, pemBytes, 0o600); err != nil {
		t.Fatalf("Datei: %v", err)
	}
	return path, tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
}

// startHealthServer startet einen gRPC-Server auf Loopback (mit TLS, wenn
// cert nicht nil ist) und gibt seine Adresse zurück.
func startHealthServer(t *testing.T, cert *tls.Certificate) string {
	t.Helper()
	var opts []grpc.ServerOption
	if cert != nil {
		opts = append(opts, grpc.Creds(credentials.NewServerTLSFromCert(cert)))
	}
	srv := grpc.NewServer(opts...)
	healthpb.RegisterHealthServer(srv, health.NewServer())
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)
	return lis.Addr().String()
}

// callHealth baut den Kanal über transportCredentials(caFile) auf und ruft
// einen Health-Check ab.
func callHealth(addr, caFile string) error {
	creds, err := transportCredentials(caFile)
	if err != nil {
		return err
	}
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{})
	return err
}

// TestTransportCredentialsWithAnchorReachesTLSServer bindet die Wahl an den
// Optionswert: mit der Ankerdatei erreicht der Aufruf den TLS-Server.
func TestTransportCredentialsWithAnchorReachesTLSServer(t *testing.T) {
	anchor, cert := writeTestCert(t)
	addr := startHealthServer(t, &cert)
	if err := callHealth(addr, anchor); err != nil {
		t.Fatalf("TLS-Aufruf mit Anker: %v", err)
	}
}

// TestTransportCredentialsWithoutAnchorFailsAgainstTLSServer: ohne Angabe
// bleibt der Transport Klartext und der TLS-Server lehnt ihn ab.
func TestTransportCredentialsWithoutAnchorFailsAgainstTLSServer(t *testing.T) {
	_, cert := writeTestCert(t)
	addr := startHealthServer(t, &cert)
	if err := callHealth(addr, ""); err == nil {
		t.Fatal("Klartext-Aufruf gegen den TLS-Server gelang")
	}
}

// TestTransportCredentialsWithoutAnchorReachesPlaintextServer: ohne Angabe
// erreicht der Aufruf einen Klartext-Server.
func TestTransportCredentialsWithoutAnchorReachesPlaintextServer(t *testing.T) {
	addr := startHealthServer(t, nil)
	if err := callHealth(addr, ""); err != nil {
		t.Fatalf("Klartext-Aufruf: %v", err)
	}
}

// TestTransportCredentialsWithAnchorFailsAgainstPlaintextServer: mit Anker
// spricht der Client TLS und erreicht einen Klartext-Server nicht.
func TestTransportCredentialsWithAnchorFailsAgainstPlaintextServer(t *testing.T) {
	anchor, _ := writeTestCert(t)
	addr := startHealthServer(t, nil)
	if err := callHealth(addr, anchor); err == nil {
		t.Fatal("TLS-Aufruf gegen den Klartext-Server gelang")
	}
}

// TestTransportCredentialsRejectsForeignAnchor: ein Anker, der nicht zum
// Zertifikat des Servers gehört, scheitert an der Prüfung der Kette.
func TestTransportCredentialsRejectsForeignAnchor(t *testing.T) {
	_, cert := writeTestCert(t)
	foreign, _ := writeTestCert(t)
	addr := startHealthServer(t, &cert)
	if err := callHealth(addr, foreign); err == nil {
		t.Fatal("Aufruf mit fremdem Anker gelang")
	}
}

// TestTransportCredentialsRejectsNameOutsideSAN: der Anker ist das
// Serverzertifikat (die Kette stimmt), sein SAN nennt aber nicht die Adresse
// der Verbindung; der Aufruf scheitert an der Namensprüfung.
func TestTransportCredentialsRejectsNameOutsideSAN(t *testing.T) {
	anchor, cert := writeTestCertSAN(t, []string{"other.example"}, nil)
	addr := startHealthServer(t, &cert)
	err := callHealth(addr, anchor)
	if err == nil {
		t.Fatal("Aufruf an einen Namen außerhalb des SAN gelang")
	}
	if !strings.Contains(err.Error(), "x509") {
		t.Fatalf("Fehler nennt keine Zertifikatsprüfung: %v", err)
	}
}

// TestTransportCredentialsRejectsUnreadableAndPEMlessAnchor: eine fehlende
// und eine PEM-lose Datei enden vor dem Verbindungsaufbau mit Fehlertext.
func TestTransportCredentialsRejectsUnreadableAndPEMlessAnchor(t *testing.T) {
	_, err := transportCredentials(filepath.Join(t.TempDir(), "fehlt.pem"))
	if err == nil || !strings.Contains(err.Error(), "nicht lesbar") {
		t.Fatalf("fehlende Datei: %v", err)
	}
	plain := filepath.Join(t.TempDir(), "kein-pem.txt")
	if err := os.WriteFile(plain, []byte("kein Zertifikat"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = transportCredentials(plain)
	if err == nil || !strings.Contains(err.Error(), "kein PEM-Zertifikat") {
		t.Fatalf("PEM-lose Datei: %v", err)
	}
}

// TestTransportCredentialsDefaultIsInsecure: ohne Angabe liefert die Wahl
// die Klartext-Credentials.
func TestTransportCredentialsDefaultIsInsecure(t *testing.T) {
	creds, err := transportCredentials("")
	if err != nil {
		t.Fatal(err)
	}
	if got := creds.Info().SecurityProtocol; got != insecure.NewCredentials().Info().SecurityProtocol {
		t.Fatalf("SecurityProtocol = %q", got)
	}
}
