package grpc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
	"github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// neuesTestZertifikat erzeugt ein selbstsigniertes Zertifikat für `127.0.0.1`
// und `localhost` zur Laufzeit des Tests — keine Datei im Repo, kein
// Host-Werkzeug. Das Zertifikat ist sein eigener Vertrauensanker.
func neuesTestZertifikat(t *testing.T) (tls.Certificate, *x509.CertPool) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Schlüssel erzeugen: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "pg-change-feed-test"},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("Zertifikat erzeugen: %v", err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("Zertifikat lesen: %v", err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(leaf)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}, pool
}

// startTLSTestServer startet den Adapter über `serve` auf einem
// Loopback-Listener (Port 0) mit der übergebenen TLS-Konfiguration und
// liefert Adresse und den Subscriber des Streams.
func startTLSTestServer(t *testing.T, tlsConfig *tls.Config) (string, *fakeSubscriber) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listener: %v", err)
	}
	subscriber := newFakeSubscriber()
	table, err := model.NewSourceTable("tbl-1", "src-1", "public", "orders")
	if err != nil {
		t.Fatalf("Tabelle bauen: %v", err)
	}
	srv := New(Config{
		Addr:        listener.Addr().String(),
		TLSConfig:   tlsConfig,
		TokenReader: testReaderToken,
		TokenAdmin:  testAdminToken,
		Subscriber:  subscriber,
		ListTables: &fakeListTables{result: inbound.ListTablesResult{
			Tables: []model.SourceTable{table}, Retained: []model.SourceTable{},
		}},
	})
	done := make(chan error, 1)
	go func() { done <- srv.serve(listener) }()
	t.Cleanup(func() {
		srv.Shutdown()
		if err := <-done; err != nil {
			t.Errorf("serve endete mit Fehler: %v", err)
		}
	})
	return listener.Addr().String(), subscriber
}

func tlsTestConn(t *testing.T, addr string, creds credentials.TransportCredentials) *grpc.ClientConn {
	t.Helper()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		t.Fatalf("Client-Verbindung bauen: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

func tlsTestListTables(conn *grpc.ClientConn) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, authorizationMetadataKey, bearerPrefix+testReaderToken)
	_, err := administrationv1.NewAdministrationClient(conn).ListTables(ctx, &administrationv1.ListTablesRequest{Source: "src-1", Publication: "cdc_pub"})
	return err
}

// TestServerTLSBedientStreamUndVerwaltung trägt den Happy Path von
// `LH-FA-SST-011` am gRPC-Server: ein TLS-Client mit dem Zertifikat als
// Vertrauensanker bekommt ein Stream-Event und einen Verwaltungs-RPC über
// dieselbe Adresse.
// Rot färbende Mutation: in `New` die Bedingung `cfg.TLSConfig != nil` durch
// `false` ersetzen — der Server bedient im Klartext, der Handshake des
// Clients scheitert.
func TestServerTLSBedientStreamUndVerwaltung(t *testing.T) {
	cert, pool := neuesTestZertifikat(t)
	addr, subscriber := startTLSTestServer(t, &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12})
	conn := tlsTestConn(t, addr, credentials.NewTLS(&tls.Config{RootCAs: pool}))

	if err := tlsTestListTables(conn); err != nil {
		t.Fatalf("Verwaltungs-RPC über TLS: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, authorizationMetadataKey, bearerPrefix+testReaderToken)
	stream, err := streamv1.NewChangeStreamClient(conn).StreamChanges(ctx, &streamv1.StreamChangesRequest{})
	if err != nil {
		t.Fatalf("Stream über TLS öffnen: %v", err)
	}
	select {
	case <-subscriber.ready:
	case <-time.After(3 * time.Second):
		t.Fatal("der Stream hat sich nicht am Broadcaster registriert")
	}
	change, err := model.NewChange("change-1", "tx-1", "table-1", 2, model.OperationInsert,
		[]byte(`{"id":1}`), []byte(`{"id":1,"name":"TLSSentinel"}`), "table-1-v1")
	if err != nil {
		t.Fatalf("Change bauen: %v", err)
	}
	change.Schema = "public"
	change.Table = "orders"
	subscriber.changes <- &change
	msg, err := recvChange(t, stream)
	if err != nil {
		t.Fatalf("Recv über TLS: %v", err)
	}
	if msg.GetChangeId() != "change-1" || string(msg.GetNewImage()) != `{"id":1,"name":"TLSSentinel"}` {
		t.Fatalf("übertragener Change: %+v", msg)
	}
}

// TestServerTLSLehntKlartextClientAb trägt die Boundary von
// `LH-FA-SST-011`: ein Klartext-Client bekommt auf derselben Adresse keine
// Antwort der API — auch nicht mit gültigem Token. Derselbe Server bedient
// im selben Test einen TLS-Client (Gegenprobe: der Server steht). Der
// gemessene Ausgang von `grpc-go` steht im Log des Tests.
// Rot färbende Mutation: in `New` die Bedingung `cfg.TLSConfig != nil` durch
// `false` ersetzen — der Klartext-Client bekommt dann eine Antwort.
func TestServerTLSLehntKlartextClientAb(t *testing.T) {
	cert, pool := neuesTestZertifikat(t)
	addr, _ := startTLSTestServer(t, &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12})

	if err := tlsTestListTables(tlsTestConn(t, addr, credentials.NewTLS(&tls.Config{RootCAs: pool}))); err != nil {
		t.Fatalf("Gegenprobe über TLS: %v", err)
	}

	err := tlsTestListTables(tlsTestConn(t, addr, insecure.NewCredentials()))
	if err == nil {
		t.Fatal("Klartext-Client wurde bedient")
	}
	t.Logf("Klartext-Client: Status %s, Fehler %v", status.Code(err), err)
	if status.Code(err) == codes.OK {
		t.Fatalf("Klartext-Client bekam Status OK: %v", err)
	}
}

// TestServerOhneTLSKonfigurationBedientKlartext trägt die Boundary „ohne die
// Konfiguration unverändert“: ohne `TLSConfig` antwortet derselbe Weg im
// Klartext, und ein TLS-Client scheitert am Handshake.
// Rot färbende Mutation: in `New` die Bedingung `cfg.TLSConfig != nil` durch
// `true` ersetzen — der Server baut Credentials aus `nil`, der Klartext-Aufruf
// scheitert.
func TestServerOhneTLSKonfigurationBedientKlartext(t *testing.T) {
	_, pool := neuesTestZertifikat(t)
	addr, _ := startTLSTestServer(t, nil)

	if err := tlsTestListTables(tlsTestConn(t, addr, insecure.NewCredentials())); err != nil {
		t.Fatalf("Klartext ohne TLS-Konfiguration: %v", err)
	}
	if err := tlsTestListTables(tlsTestConn(t, addr, credentials.NewTLS(&tls.Config{RootCAs: pool}))); err == nil {
		t.Fatal("TLS-Client erreichte einen Klartext-Server")
	}
}
