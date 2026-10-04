package http

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"io"
	"math/big"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

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
// liefert die Adresse; `Close` im Cleanup beendet auch einen offenen
// SSE-Stream.
func startTLSTestServer(t *testing.T, cfg Config, tlsConfig *tls.Config) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listener: %v", err)
	}
	cfg.Addr = listener.Addr().String()
	cfg.TLSConfig = tlsConfig
	cfg.TokenReader = testReaderToken
	cfg.TokenAdmin = testAdminToken
	srv := New(cfg)
	done := make(chan error, 1)
	go func() { done <- srv.serve(listener) }()
	t.Cleanup(func() {
		_ = srv.httpServer.Close()
		if err := <-done; err != nil {
			t.Errorf("serve endete mit Fehler: %v", err)
		}
	})
	return listener.Addr().String()
}

func tlsTestClient(pool *x509.CertPool) *http.Client {
	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}, DisableKeepAlives: true},
	}
}

func tlsTestRequest(t *testing.T, client *http.Client, method, url, token, body string) *http.Response {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatalf("Request bauen: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request senden (%s): %v", url, err)
	}
	return resp
}

// tlsTestConfig trägt die Konfiguration, die `internal/bootstrap` baut: ein
// Zertifikat, Mindestversion TLS 1.2.
func tlsTestConfig(cert tls.Certificate) *tls.Config {
	return &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12}
}

func tlsTestFakes(t *testing.T) Config {
	t.Helper()
	table, err := model.NewSourceTable("tbl-1", "src-1", "public", "orders")
	if err != nil {
		t.Fatalf("Tabelle bauen: %v", err)
	}
	return Config{
		RegisterConsumer: newFakeRegisterConsumerUseCase(),
		ListTables: fakeListTablesUseCase{result: inbound.ListTablesResult{
			Tables: []model.SourceTable{table}, Retained: []model.SourceTable{},
		}},
		Subscriber: newFakeChangeSubscriber(4),
	}
}

// TestServerTLSBedientJedeFaehigkeit trägt den Happy Path von
// `LH-FA-SST-011` am HTTP-Server: ein TLS-Client mit dem Zertifikat als
// Vertrauensanker erreicht einen lesenden Endpunkt, einen administrativen
// Endpunkt und den SSE-Stream samt Event.
// Rot färbende Mutation: in `serve` die Verzweigung auf `TLSConfig != nil`
// durch `false` ersetzen — der Server bedient im Klartext, der Handshake des
// Clients scheitert.
func TestServerTLSBedientJedeFaehigkeit(t *testing.T) {
	cert, pool := neuesTestZertifikat(t)
	cfg := tlsTestFakes(t)
	subscriber := cfg.Subscriber.(*fakeChangeSubscriber)
	addr := startTLSTestServer(t, cfg, tlsTestConfig(cert))
	client := tlsTestClient(pool)

	reader := tlsTestRequest(t, client, http.MethodGet, "https://"+addr+"/tables?source=src-1&publication=cdc_pub", testReaderToken, "")
	defer reader.Body.Close()
	if reader.StatusCode != http.StatusOK {
		t.Fatalf("lesender Endpunkt über TLS: Status %d (Erwartung 200)", reader.StatusCode)
	}
	if reader.TLS == nil || reader.TLS.Version < tls.VersionTLS12 {
		t.Fatalf("Antwort trägt keine TLS-Verbindung ab Version 1.2: %+v", reader.TLS)
	}

	admin := tlsTestRequest(t, client, http.MethodPost, "https://"+addr+"/consumers", testAdminToken, `{"consumer_id":"c1","name":"Consumer 1"}`)
	defer admin.Body.Close()
	if admin.StatusCode != http.StatusCreated {
		t.Fatalf("administrativer Endpunkt über TLS: Status %d (Erwartung 201)", admin.StatusCode)
	}

	stream := tlsTestRequest(t, client, http.MethodGet, "https://"+addr+"/changes/stream", testReaderToken, "")
	defer stream.Body.Close()
	if stream.StatusCode != http.StatusOK {
		t.Fatalf("SSE über TLS: Status %d (Erwartung 200)", stream.StatusCode)
	}
	select {
	case <-subscriber.ready:
	case <-time.After(3 * time.Second):
		t.Fatal("der Stream hat sich nicht am Broadcaster registriert")
	}
	change := neuerStreamTestChange(t)
	subscriber.changes <- &change
	event, data := readSSEEvent(t, bufio.NewReader(stream.Body))
	if event != "change" || !strings.Contains(data, "StreamE2ESentinel") {
		t.Fatalf("SSE-Event über TLS: event=%q data=%q", event, data)
	}
}

// TestServerTLSLehntKlartextClientAb trägt die Boundary von
// `LH-FA-SST-011`: ein Klartext-Client bekommt auf derselben Adresse keine
// Antwort der API — auch nicht mit gültigem Token. Dieselbe Adresse bedient
// einen TLS-Client im selben Test (Gegenprobe: der Server steht). Das
// gemessene Verhalten von `net/http` steht im Log des Tests.
// Rot färbende Mutation: in `serve` die Verzweigung auf `TLSConfig != nil`
// durch `false` ersetzen — der Klartext-Client bekommt dann `200`.
func TestServerTLSLehntKlartextClientAb(t *testing.T) {
	cert, pool := neuesTestZertifikat(t)
	addr := startTLSTestServer(t, tlsTestFakes(t), tlsTestConfig(cert))

	secure := tlsTestRequest(t, tlsTestClient(pool), http.MethodGet, "https://"+addr+"/tables?source=src-1&publication=cdc_pub", testReaderToken, "")
	_ = secure.Body.Close()
	if secure.StatusCode != http.StatusOK {
		t.Fatalf("Gegenprobe über TLS: Status %d (Erwartung 200)", secure.StatusCode)
	}

	plain := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{DisableKeepAlives: true}}
	req, err := http.NewRequest(http.MethodGet, "http://"+addr+"/tables?source=src-1&publication=cdc_pub", nil)
	if err != nil {
		t.Fatalf("Request bauen: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+testReaderToken)
	resp, err := plain.Do(req)
	if err != nil {
		t.Logf("Klartext-Client: Transportfehler %v", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Klartext-Client: Status %d, Körper %q", resp.StatusCode, strings.TrimSpace(string(body)))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		t.Fatalf("Klartext-Client wurde bedient: Status %d", resp.StatusCode)
	}
	if strings.Contains(string(body), `"tables"`) {
		t.Fatalf("Klartext-Client bekam Inhalt der API: %q", body)
	}
}

// TestServerOhneTLSKonfigurationBedientKlartext trägt die Boundary „ohne die
// Konfiguration unverändert“: ohne `TLSConfig` antwortet derselbe Weg im
// Klartext, und ein TLS-Client scheitert am Handshake.
// Rot färbende Mutation: in `serve` die Verzweigung auf `TLSConfig != nil`
// durch `true` ersetzen — der Server verlangt TLS ohne Zertifikat.
func TestServerOhneTLSKonfigurationBedientKlartext(t *testing.T) {
	_, pool := neuesTestZertifikat(t)
	addr := startTLSTestServer(t, tlsTestFakes(t), nil)

	plain := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{DisableKeepAlives: true}}
	resp := tlsTestRequest(t, plain, http.MethodGet, "http://"+addr+"/tables?source=src-1&publication=cdc_pub", testReaderToken, "")
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Klartext ohne TLS-Konfiguration: Status %d (Erwartung 200)", resp.StatusCode)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://"+addr+"/tables", nil)
	if err != nil {
		t.Fatalf("Request bauen: %v", err)
	}
	if _, err := tlsTestClient(pool).Do(req); err == nil {
		t.Fatal("TLS-Client erreichte einen Klartext-Server")
	}
}
