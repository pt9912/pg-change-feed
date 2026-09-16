package bootstrap_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/bootstrap"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// TestRunScheitertAmErstenKonstruktorOhneErreichbareQuelle trägt den
// netzlos erreichbaren Teil von `Run` (`ADR-0082` §Kontext (3)): die
// Verdrahtung baut den Telemetrie-Adapter, ruft den ersten Konstruktor
// und reicht dessen Fehler durch, statt weiter zu verdrahten.
// `postgresstorage.New` ist `pgxpool.New` plus `pool.Ping` — der
// geschlossene lokale Port scheitert im `Ping`, ohne Dienst und ohne
// Netz (`connect_timeout` begrenzt den Verbindungsversuch am Adapter).
//
// Gebunden an die Eingabeseite: die drei DSNs tragen je einen **eigenen**
// Datenbanknamen, und die Fehlerzeile wird auf den Namen aus
// `cfg.CaptureDSN` geprüft — der Name aus `cfg.AdminDSN` (`admin-db`)
// darf darin nicht vorkommen: der Lauf erreicht den Aktivierungs-Pool
// nicht. Die tragende Bindung an den **ersten** Konstruktor ist der
// Sentinel: ein `Run`, das seinen Fehler verschluckt, endet eine Stufe
// später am Schema-Store — **derselbe** `cfg.CaptureDSN`, eigener
// Sentinel (`outbound.ErrSchemaStoreStorage`) — und färbt diesen Test
// über die Klassen-Prüfung rot. Die Klasse `storage`
// (`outbound.ErrStorage`) kommt vom ersten Konstruktor, nicht von `Run`.
func TestRunScheitertAmErstenKonstruktorOhneErreichbareQuelle(t *testing.T) {
	const (
		nichtErreichbar = "postgres://x:x@127.0.0.1:1/%s?sslmode=disable&connect_timeout=1"
		captureDB       = "capture-db"
		adminDB         = "admin-db"
	)

	cfg := bootstrap.Config{
		CaptureDSN:  fmt.Sprintf(nichtErreichbar, captureDB),
		AdminDSN:    fmt.Sprintf(nichtErreichbar, adminDB),
		ReaderDSN:   fmt.Sprintf(nichtErreichbar, "reader-db"),
		Source:      model.SourceID("src-1"),
		Publication: "pub_1",
		Slot:        "slot_1",
	}

	err := bootstrap.Run(context.Background(), cfg)
	if err == nil {
		t.Fatal("Run ohne erreichbare Quelle: kein Fehler — die Verdrahtung läuft hinter dem Ping-Riegel weiter")
	}
	if !errors.Is(err, outbound.ErrStorage) {
		t.Fatalf("Run-Fehler = %v, wollen die Klasse storage (outbound.ErrStorage) aus dem ersten Konstruktor", err)
	}
	if !strings.Contains(err.Error(), captureDB) {
		t.Fatalf("Run-Fehler = %v, wollen den Datenbanknamen aus cfg.CaptureDSN (%q)", err, captureDB)
	}
	if strings.Contains(err.Error(), adminDB) {
		t.Fatalf("Run-Fehler = %v, wollen keinen Zugriff über den Aktivierungs-Pool (%q)", err, adminDB)
	}
}
