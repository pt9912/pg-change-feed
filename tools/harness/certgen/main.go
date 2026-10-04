// Command certgen erzeugt ein selbstsigniertes TLS-Zertifikat samt privatem
// Schlüssel für die TLS-Phase von tools/harness/run-integration-tests.sh
// (LH-FA-SST-011). Aufruf: `certgen <verzeichnis> <name> <alternativname>...`
// schreibt `<verzeichnis>/<name>.pem` (Zertifikat) und
// `<verzeichnis>/<name>-key.pem` (Schlüssel, PKCS#8) — jeder Alternativname
// ist ein DNS-Name oder eine IP-Adresse. Das Zertifikat ist sein eigener
// Vertrauensanker; es gilt 24 Stunden und wird nur in das übergebene
// Verzeichnis geschrieben, nie in den Arbeitsbaum. Nur Standardbibliothek.
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: certgen <verzeichnis> <name> <alternativname>...")
		os.Exit(2)
	}
	dir, name, hosts := os.Args[1], os.Args[2], os.Args[3:]
	if err := generate(dir, name, hosts); err != nil {
		fmt.Fprintf(os.Stderr, "certgen: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("CERTGEN cert=%s key=%s hosts=%d\n",
		filepath.Join(dir, name+".pem"), filepath.Join(dir, name+"-key.pem"), len(hosts))
}

func generate(dir, name string, hosts []string) error {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("Schlüssel erzeugen: %w", err)
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 62))
	if err != nil {
		return fmt.Errorf("Seriennummer erzeugen: %w", err)
	}
	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: name},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	for _, host := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, host)
		}
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("Zertifikat erzeugen: %w", err)
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return fmt.Errorf("Schlüssel kodieren: %w", err)
	}
	if err := writePEM(filepath.Join(dir, name+".pem"), "CERTIFICATE", der, 0o644); err != nil {
		return err
	}
	return writePEM(filepath.Join(dir, name+"-key.pem"), "PRIVATE KEY", keyDER, 0o600)
}

func writePEM(path, blockType string, der []byte, mode os.FileMode) error {
	out := pem.EncodeToMemory(&pem.Block{Type: blockType, Bytes: der})
	if err := os.WriteFile(path, out, mode); err != nil {
		return fmt.Errorf("%s schreiben: %w", path, err)
	}
	return nil
}
