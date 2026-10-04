// tls.go trägt die TLS-Konfiguration der beiden Driving-Server (`ADR-0150`):
// die Vollständigkeitsprüfung des Paars aus Zertifikat und Schlüssel, die beide
// Konfigurations-Zugriffswege rufen, und den einen Ort, an dem das Paar geladen
// und die `tls.Config` gebaut wird — beide Adapter erhalten sie fertig.
package bootstrap

import (
	"crypto/tls"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
)

// ErrTLSPairIncomplete trägt die Fehlerklasse `configuration` eines
// unvollständigen TLS-Paars: genau einer der beiden Pfade ist gesetzt.
var ErrTLSPairIncomplete = messagecode.New(messagecode.TLSPairIncomplete, "TLS-Zertifikat und -Schlüssel nur teilweise gesetzt")

// ErrTLSPairUnusable trägt die Fehlerklasse `configuration` eines Paars, das
// sich nicht laden lässt: ein Pfad nicht lesbar, keine gültigen PEM-Daten oder
// Zertifikat und Schlüssel passen nicht zusammen.
var ErrTLSPairUnusable = messagecode.New(messagecode.TLSPairUnusable, "TLS-Zertifikat und -Schlüssel nicht ladbar")

// tlsConfigError ist der Fehler einer TLS-Konfiguration: sein Text trägt den
// Code des Sentinels samt Einzelheit, und er ist zugleich ein
// `ErrConfiguration` (Start-Hindernis der Verdrahtung).
type tlsConfigError struct {
	sentinel error
	detail   string
	cause    error
}

func (e *tlsConfigError) Error() string { return e.sentinel.Error() + ": " + e.detail }

func (e *tlsConfigError) Unwrap() []error {
	if e.cause != nil {
		return []error{e.sentinel, ErrConfiguration, e.cause}
	}
	return []error{e.sentinel, ErrConfiguration}
}

// validateTLSPair verlangt, dass Zertifikat und Schlüssel gemeinsam gesetzt
// oder gemeinsam leer sind. Beide Zugriffswege (`ConfigFromEnv`,
// `mergeConfig`) rufen diese eine Funktion; sie liest keine Datei, damit auch
// `--healthcheck` dieselbe Konfiguration lädt, ohne das Paar zu brauchen.
func validateTLSPair(certFile, keyFile string) error {
	switch {
	case certFile != "" && keyFile == "":
		return &tlsConfigError{sentinel: ErrTLSPairIncomplete,
			detail: fmt.Sprintf("%s gesetzt, aber %s fehlt (weder Umgebungsvariable noch Konfigurationsdatei)", envTLSCertFile, envTLSKeyFile)}
	case certFile == "" && keyFile != "":
		return &tlsConfigError{sentinel: ErrTLSPairIncomplete,
			detail: fmt.Sprintf("%s gesetzt, aber %s fehlt (weder Umgebungsvariable noch Konfigurationsdatei)", envTLSKeyFile, envTLSCertFile)}
	}
	return nil
}

// newTLSConfig lädt das Paar und baut die `tls.Config` beider Server: ein
// Zertifikat, Mindestversion TLS 1.2, kein Client-Zertifikat. Ohne Paar ist das
// Ergebnis `nil` (beide Server im Klartext). Ein Ladefehler nennt Pfade und
// Ursache, nie Schlüsselinhalt.
func newTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	if certFile == "" && keyFile == "" {
		return nil, nil
	}
	if err := validateTLSPair(certFile, keyFile); err != nil {
		return nil, err
	}
	pair, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, &tlsConfigError{sentinel: ErrTLSPairUnusable, cause: err,
			detail: fmt.Sprintf("%s=%q, %s=%q: %v", envTLSCertFile, certFile, envTLSKeyFile, keyFile, err)}
	}
	return &tls.Config{
		Certificates: []tls.Certificate{pair},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
