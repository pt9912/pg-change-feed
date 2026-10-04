package bootstrap

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// tlsTestPaar schreibt ein selbstsigniertes Zertifikat samt Schlüssel in ein
// Test-Verzeichnis (keine Datei im Repo, kein Host-Werkzeug) und liefert die
// beiden Pfade und das Zertifikat als Vertrauensanker. `gueltigBis` bestimmt
// das Ablaufdatum.
func tlsTestPaar(t *testing.T, gueltigBis time.Time) (certPath, keyPath string, pool *x509.CertPool) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Schlüssel erzeugen: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "pg-change-feed-test"},
		NotBefore:             gueltigBis.Add(-time.Hour),
		NotAfter:              gueltigBis,
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
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("Schlüssel kodieren: %v", err)
	}
	dir := t.TempDir()
	certPath = filepath.Join(dir, "cert.pem")
	keyPath = filepath.Join(dir, "key.pem")
	if err := os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600); err != nil {
		t.Fatalf("Zertifikat schreiben: %v", err)
	}
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}), 0o600); err != nil {
		t.Fatalf("Schlüssel schreiben: %v", err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("Zertifikat lesen: %v", err)
	}
	pool = x509.NewCertPool()
	pool.AddCert(leaf)
	return certPath, keyPath, pool
}

// tlsZugriffswege trägt die zwei Zugriffswege der Konfiguration: den reinen
// Umgebungsweg und den verdrahteten Weg mit geladener Datei, die `dateiInhalt`
// trägt. Beide rufen `validateTLSPair`.
func tlsZugriffswege(t *testing.T, dateiInhalt string) map[string]func(map[string]string) (Config, error) {
	t.Helper()
	datei := writeConfigFile(t, "source_id: src-datei\npublication: pub-datei\nslot: slot-datei\ntables:\n  public.t1:\n    table_id: tbl-1\n    schema_version: sv-1\n"+dateiInhalt)
	return map[string]func(map[string]string) (Config, error){
		"ConfigFromEnv": func(env map[string]string) (Config, error) {
			return ConfigFromEnv(func(name string) string { return env[name] })
		},
		"ConfigFromEnvAndFile": func(env map[string]string) (Config, error) {
			env["CDC_CONFIG_FILE"] = datei
			return ConfigFromEnvAndFile(func(name string) string { return env[name] })
		},
	}
}

// TestTLSPaarWirdGelesen trägt `LH-FA-SST-011` Happy Path und Boundary der
// Konfiguration: kein Pfad lässt TLS aus, beide Pfade aus der Umgebung werden
// gelesen, und unter geladener Datei schlägt die Umgebungsvariable den
// Datei-Wert Feld für Feld.
// Rot färbende Mutation: in `mergeConfig` die Zuweisung von `cfg.TLSKeyFile`
// streichen — die Fälle des Dateiwegs tragen keinen Schlüssel mehr.
func TestTLSPaarWirdGelesen(t *testing.T) {
	cases := []struct {
		name     string
		datei    string
		env      map[string]string
		wantCert string
		wantKey  string
		nurDatei bool
	}{
		{name: "ungesetzt", wantCert: "", wantKey: ""},
		{name: "beide aus der Umgebung", env: map[string]string{"CDC_TLS_CERT_FILE": "/c.pem", "CDC_TLS_KEY_FILE": "/k.pem"},
			wantCert: "/c.pem", wantKey: "/k.pem"},
		{name: "leere Variablen gelten als ungesetzt", env: map[string]string{"CDC_TLS_CERT_FILE": "", "CDC_TLS_KEY_FILE": ""}},
		{name: "beide aus der Datei", datei: "tls_cert_file: /dc.pem\ntls_key_file: /dk.pem\n", nurDatei: true,
			wantCert: "/dc.pem", wantKey: "/dk.pem"},
		{name: "Zertifikat aus der Datei, Schlüssel aus der Umgebung", datei: "tls_cert_file: /dc.pem\n",
			env: map[string]string{"CDC_TLS_KEY_FILE": "/k.pem"}, nurDatei: true, wantCert: "/dc.pem", wantKey: "/k.pem"},
		{name: "Umgebung schlägt Datei je Feld", datei: "tls_cert_file: /dc.pem\ntls_key_file: /dk.pem\n",
			env: map[string]string{"CDC_TLS_CERT_FILE": "/c.pem"}, nurDatei: true, wantCert: "/c.pem", wantKey: "/dk.pem"},
	}
	for _, tc := range cases {
		for weg, laden := range tlsZugriffswege(t, tc.datei) {
			if tc.nurDatei && weg == "ConfigFromEnv" {
				continue
			}
			t.Run(weg+"/"+tc.name, func(t *testing.T) {
				env := vollständigeEnvOhneDatei()
				for k, v := range tc.env {
					env[k] = v
				}
				cfg, err := laden(env)
				if err != nil {
					t.Fatalf("Konfiguration lädt nicht: %v", err)
				}
				if cfg.TLSCertFile != tc.wantCert || cfg.TLSKeyFile != tc.wantKey {
					t.Fatalf("TLS-Pfade = %q / %q, erwartet %q / %q", cfg.TLSCertFile, cfg.TLSKeyFile, tc.wantCert, tc.wantKey)
				}
			})
		}
	}
}

// TestTLSPaarUnvollstaendigEndetMitConfiguration trägt `LH-FA-SST-011`
// Negative: genau ein Pfad — aus der Umgebung, aus der Datei oder aus je einer
// Quelle ohne Gegenstück — endet mit der Fehlerklasse `configuration` und dem
// Meldungscode `PCF-E2009`; die Zeile nennt den gesetzten und den fehlenden
// Namen. Die Gegenprobe mit beiden Pfaden lädt (`TestTLSPaarWirdGelesen`).
// Rot färbende Mutation: in `validateTLSPair` die Bedingung des ersten Falls
// von `certFile != "" && keyFile == ""` auf `certFile != "" || keyFile == ""`
// ändern — der Fall „nichts gesetzt“ von `TestTLSPaarWirdGelesen` endet mit
// einem Fehler; wird der zweite Fall gestrichen, färbt sich „nur Schlüssel“.
func TestTLSPaarUnvollstaendigEndetMitConfiguration(t *testing.T) {
	cases := []struct {
		name           string
		datei          string
		env            map[string]string
		nurDatei       bool
		gesetzt, fehlt string
	}{
		{name: "nur Zertifikat aus der Umgebung", env: map[string]string{"CDC_TLS_CERT_FILE": "/c.pem"},
			gesetzt: "CDC_TLS_CERT_FILE", fehlt: "CDC_TLS_KEY_FILE"},
		{name: "nur Schlüssel aus der Umgebung", env: map[string]string{"CDC_TLS_KEY_FILE": "/k.pem"},
			gesetzt: "CDC_TLS_KEY_FILE", fehlt: "CDC_TLS_CERT_FILE"},
		{name: "nur Zertifikat aus der Datei", datei: "tls_cert_file: /dc.pem\n", nurDatei: true,
			gesetzt: "CDC_TLS_CERT_FILE", fehlt: "CDC_TLS_KEY_FILE"},
		{name: "nur Schlüssel aus der Datei", datei: "tls_key_file: /dk.pem\n", nurDatei: true,
			gesetzt: "CDC_TLS_KEY_FILE", fehlt: "CDC_TLS_CERT_FILE"},
	}
	for _, tc := range cases {
		for weg, laden := range tlsZugriffswege(t, tc.datei) {
			if tc.nurDatei && weg == "ConfigFromEnv" {
				continue
			}
			t.Run(weg+"/"+tc.name, func(t *testing.T) {
				env := vollständigeEnvOhneDatei()
				for k, v := range tc.env {
					env[k] = v
				}
				_, err := laden(env)
				if err == nil {
					t.Fatal("unvollständiges TLS-Paar lädt ohne Fehler")
				}
				if !errors.Is(err, ErrConfiguration) || !errors.Is(err, ErrTLSPairIncomplete) {
					t.Fatalf("Fehler ist nicht ErrConfiguration und ErrTLSPairIncomplete: %v", err)
				}
				if code, ok := messagecode.From(err); !ok || code != messagecode.TLSPairIncomplete {
					t.Fatalf("Code = %q (%v), erwartet %q", code, ok, messagecode.TLSPairIncomplete)
				}
				for _, want := range []string{"PCF-E2009", "configuration", tc.gesetzt + " gesetzt", tc.fehlt + " fehlt"} {
					if !strings.Contains(err.Error(), want) {
						t.Fatalf("Fehlertext %q trägt %q nicht", err.Error(), want)
					}
				}
			})
		}
	}
}

// TestTLSDateiSchluesselGehoerenNichtZurZugangsdatenKlasse trägt die
// Gegenprobe zur Zugangsdaten-Klasse: eine Datei mit `tls_cert_file` und
// `tls_key_file` lädt ohne Fehler, ihre Pfade sind keine Zugangsdaten.
// Rot färbende Mutation: `tls_key_file` in `forbiddenFileCredentialKeys`
// aufnehmen — das Laden endet mit der Zeile der Zugangsdaten-Klasse.
func TestTLSDateiSchluesselGehoerenNichtZurZugangsdatenKlasse(t *testing.T) {
	path := writeConfigFile(t, "tls_cert_file: /c.pem\ntls_key_file: /k.pem\n")
	file, err := ConfigFromFile(path)
	if err != nil {
		t.Fatalf("Datei mit beiden TLS-Schlüsseln lädt nicht: %v", err)
	}
	if file.TLSCertFile != "/c.pem" || file.TLSKeyFile != "/k.pem" {
		t.Fatalf("Datei-Felder = %q / %q", file.TLSCertFile, file.TLSKeyFile)
	}
	for _, key := range []string{"tls_cert_file", "tls_key_file"} {
		for _, forbidden := range forbiddenFileCredentialKeys {
			if forbidden == key {
				t.Fatalf("%q steht in der Zugangsdaten-Klasse", key)
			}
		}
	}
}

// TestNewTLSConfigBautEineKonfiguration trägt den Bau der einen
// `tls.Config`: kein Paar ergibt `nil` (Klartext), ein gültiges Paar ein
// Zertifikat, Mindestversion TLS 1.2 explizit gesetzt und keine
// Client-Zertifikat-Pflicht.
// Rot färbende Mutation: `MinVersion: tls.VersionTLS12` aus `newTLSConfig`
// streichen oder auf `tls.VersionTLS10` setzen — die Prüfung des Feldes wird
// rot; `Certificates` zu `nil` färbt die Zertifikat-Prüfung.
func TestNewTLSConfigBautEineKonfiguration(t *testing.T) {
	if cfg, err := newTLSConfig("", ""); cfg != nil || err != nil {
		t.Fatalf("ohne Paar: %v / %v, erwartet nil / nil", cfg, err)
	}
	certPath, keyPath, _ := tlsTestPaar(t, time.Now().Add(time.Hour))
	cfg, err := newTLSConfig(certPath, keyPath)
	if err != nil {
		t.Fatalf("gültiges Paar: %v", err)
	}
	if cfg.MinVersion != tls.VersionTLS12 {
		t.Fatalf("MinVersion = %#x, erwartet %#x (TLS 1.2)", cfg.MinVersion, tls.VersionTLS12)
	}
	if len(cfg.Certificates) != 1 {
		t.Fatalf("Zertifikate = %d, erwartet 1", len(cfg.Certificates))
	}
	if cfg.ClientAuth != tls.NoClientCert {
		t.Fatalf("ClientAuth = %v, erwartet keine Client-Zertifikat-Pflicht", cfg.ClientAuth)
	}
}

// TestNewTLSConfigMindestversionAmDraht trägt die Mindestversion an der
// Verbindung: ein Client, der höchstens TLS 1.1 spricht, scheitert am
// Server mit der gebauten Konfiguration; die Gegenprobe mit TLS 1.2 gelingt.
// Der Server stammt aus `net/http/httptest` mit derselben `tls.Config` — die
// Aussage gilt für diese Konfiguration, nicht für die zwei Adapter.
// Rot färbende Mutation: `MinVersion` in `newTLSConfig` auf `tls.VersionTLS10`
// setzen — der TLS-1.1-Client verbindet.
func TestNewTLSConfigMindestversionAmDraht(t *testing.T) {
	certPath, keyPath, pool := tlsTestPaar(t, time.Now().Add(time.Hour))
	cfg, err := newTLSConfig(certPath, keyPath)
	if err != nil {
		t.Fatalf("gültiges Paar: %v", err)
	}
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	server.TLS = cfg
	server.StartTLS()
	t.Cleanup(server.Close)

	get := func(min, max uint16) error {
		client := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{RootCAs: pool, MinVersion: min, MaxVersion: max},
			DisableKeepAlives: true,
		}}
		resp, err := client.Get(server.URL)
		if err == nil {
			_ = resp.Body.Close()
		}
		return err
	}
	if err := get(tls.VersionTLS12, tls.VersionTLS12); err != nil {
		t.Fatalf("TLS-1.2-Client: %v", err)
	}
	err = get(tls.VersionTLS10, tls.VersionTLS11)
	if err == nil {
		t.Fatal("TLS-1.1-Client verbindet mit der gebauten Konfiguration")
	}
	t.Logf("TLS-1.1-Client: %v", err)
}

// TestNewTLSConfigLadefehlerEndenMitConfiguration trägt `LH-FA-SST-011`
// Negative: ein nicht lesbarer Pfad, Dateien ohne PEM-Daten und ein
// Zertifikat mit fremdem Schlüssel enden mit der Fehlerklasse `configuration`
// und dem Meldungscode `PCF-E2010`; die Zeile nennt die Pfade und kein
// Schlüsselmaterial. Die Gegenprobe ist das gültige Paar
// (`TestNewTLSConfigBautEineKonfiguration`).
// Rot färbende Mutation: in `newTLSConfig` den Fehler von
// `tls.LoadX509KeyPair` verwerfen — jeder Fall liefert eine Konfiguration.
func TestNewTLSConfigLadefehlerEndenMitConfiguration(t *testing.T) {
	certPath, keyPath, _ := tlsTestPaar(t, time.Now().Add(time.Hour))
	_, fremderKey, _ := tlsTestPaar(t, time.Now().Add(time.Hour))
	leer := filepath.Join(t.TempDir(), "leer.pem")
	if err := os.WriteFile(leer, []byte("kein PEM\n"), 0o600); err != nil {
		t.Fatalf("Hilfsdatei schreiben: %v", err)
	}
	fehlt := filepath.Join(t.TempDir(), "fehlt.pem")
	cases := []struct {
		name, cert, key string
	}{
		{"Zertifikat nicht lesbar", fehlt, keyPath},
		{"Schlüssel nicht lesbar", certPath, fehlt},
		{"Zertifikat ohne PEM-Daten", leer, keyPath},
		{"Schlüssel ohne PEM-Daten", certPath, leer},
		{"Schlüssel gehört zu einem anderen Zertifikat", certPath, fremderKey},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := newTLSConfig(tc.cert, tc.key)
			if err == nil || cfg != nil {
				t.Fatalf("nicht ladbares Paar: %v / %v", cfg, err)
			}
			if !errors.Is(err, ErrConfiguration) || !errors.Is(err, ErrTLSPairUnusable) {
				t.Fatalf("Fehler ist nicht ErrConfiguration und ErrTLSPairUnusable: %v", err)
			}
			if code, ok := messagecode.From(err); !ok || code != messagecode.TLSPairUnusable {
				t.Fatalf("Code = %q (%v), erwartet %q", code, ok, messagecode.TLSPairUnusable)
			}
			for _, want := range []string{"PCF-E2010", "configuration", tc.cert, tc.key} {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("Fehlertext %q trägt %q nicht", err.Error(), want)
				}
			}
			if strings.Contains(err.Error(), "PRIVATE KEY") {
				t.Fatalf("Fehlertext trägt Schlüsselmaterial: %q", err.Error())
			}
		})
	}
}

// TestNewTLSConfigLiestDieKetteDerZertifikatsdatei trägt die Form der
// Zertifikatsdatei: eine Datei mit dem Zertifikat und einem weiteren
// Zertifikat dahinter (Kette) lädt, und die Konfiguration trägt beide Glieder.
// Die Bindung liegt an der Eingabe: der Fall mit einem einzelnen Zertifikat
// trägt ein Glied.
func TestNewTLSConfigLiestDieKetteDerZertifikatsdatei(t *testing.T) {
	certPath, keyPath, _ := tlsTestPaar(t, time.Now().Add(time.Hour))
	zweitesPath, _, _ := tlsTestPaar(t, time.Now().Add(time.Hour))
	erstes, err := os.ReadFile(certPath)
	if err != nil {
		t.Fatalf("Zertifikat lesen: %v", err)
	}
	zweites, err := os.ReadFile(zweitesPath)
	if err != nil {
		t.Fatalf("zweites Zertifikat lesen: %v", err)
	}
	kette := filepath.Join(t.TempDir(), "kette.pem")
	if err := os.WriteFile(kette, append(erstes, zweites...), 0o600); err != nil {
		t.Fatalf("Kette schreiben: %v", err)
	}
	einzeln, err := newTLSConfig(certPath, keyPath)
	if err != nil {
		t.Fatalf("einzelnes Zertifikat: %v", err)
	}
	if got := len(einzeln.Certificates[0].Certificate); got != 1 {
		t.Fatalf("einzelnes Zertifikat trägt %d Glieder, erwartet 1", got)
	}
	cfg, err := newTLSConfig(kette, keyPath)
	if err != nil {
		t.Fatalf("Zertifikat mit Kette: %v", err)
	}
	if got := len(cfg.Certificates[0].Certificate); got != 2 {
		t.Fatalf("Zertifikat mit Kette trägt %d Glieder, erwartet 2", got)
	}
}

// TestNewTLSConfigPruftAblaufUndNamenNicht hält die benannte Grenze fest: ein
// abgelaufenes Zertifikat lädt, der Server startet damit; den Ablauf und den
// Namen prüft der Client. Ändert eine spätere Fassung das, färbt dieser Test
// sich und der Handbuchtext zu „ungültig“ gehört mit.
func TestNewTLSConfigPruftAblaufUndNamenNicht(t *testing.T) {
	certPath, keyPath, _ := tlsTestPaar(t, time.Now().Add(-time.Minute))
	if _, err := newTLSConfig(certPath, keyPath); err != nil {
		t.Fatalf("abgelaufenes Zertifikat wird beim Laden abgelehnt: %v", err)
	}
}

// TestRunLaedtTLSPaarVorJederVerbindung trägt `LH-FA-SST-011` Negative am
// Start: ein nicht ladbares Paar beendet `Run` mit `PCF-E2010`, bevor eine
// Datenbankverbindung versucht wird (die DSNs zeigen auf einen geschlossenen
// Port: wäre die Reihenfolge anders, endete `Run` mit der Klasse `storage`)
// und ohne dass auf der HTTP- oder gRPC-Adresse ein Listener steht. Die
// Gegenprobe mit gültigem Paar erreicht die Datenbank und scheitert dort.
// Rot färbende Mutation: in `Run` den Fehler von `newTLSConfig` verwerfen und
// mit `nil` fortfahren — der Lauf endet mit der Klasse `storage` statt mit
// `PCF-E2010`.
func TestRunLaedtTLSPaarVorJederVerbindung(t *testing.T) {
	const nichtErreichbar = "postgres://x:x@127.0.0.1:1/%s?sslmode=disable&connect_timeout=1"
	freiePort := func() string {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("Port reservieren: %v", err)
		}
		addr := l.Addr().String()
		_ = l.Close()
		return addr
	}
	httpAddr, grpcAddr := freiePort(), freiePort()
	certPath, keyPath, _ := tlsTestPaar(t, time.Now().Add(time.Hour))
	fehlt := filepath.Join(t.TempDir(), "fehlt.pem")

	basis := Config{
		CaptureDSN:  strings.Replace(nichtErreichbar, "%s", "capture-db", 1),
		AdminDSN:    strings.Replace(nichtErreichbar, "%s", "admin-db", 1),
		ReaderDSN:   strings.Replace(nichtErreichbar, "%s", "reader-db", 1),
		Source:      model.SourceID("src-1"),
		Publication: "pub_1",
		Slot:        "slot_1",
		HTTPAddr:    httpAddr,
		GRPCAddr:    grpcAddr,
	}

	unbrauchbar := basis
	unbrauchbar.TLSCertFile, unbrauchbar.TLSKeyFile = fehlt, keyPath
	err := Run(context.Background(), unbrauchbar)
	if !errors.Is(err, ErrTLSPairUnusable) || errors.Is(err, outbound.ErrStorage) {
		t.Fatalf("Run mit nicht ladbarem Paar: %v, erwartet ErrTLSPairUnusable ohne Datenbankfehler", err)
	}
	for _, addr := range []string{httpAddr, grpcAddr} {
		if conn, dialErr := net.DialTimeout("tcp", addr, time.Second); dialErr == nil {
			_ = conn.Close()
			t.Fatalf("nach dem Ladefehler hört ein Server auf %s", addr)
		}
	}

	gueltig := basis
	gueltig.TLSCertFile, gueltig.TLSKeyFile = certPath, keyPath
	err = Run(context.Background(), gueltig)
	if !errors.Is(err, outbound.ErrStorage) || errors.Is(err, ErrTLSPairUnusable) {
		t.Fatalf("Run mit gültigem Paar: %v, erwartet den Datenbankfehler der Klasse storage", err)
	}
}
