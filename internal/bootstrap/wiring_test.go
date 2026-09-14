package bootstrap_test

import (
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/bootstrap"
)

// getenv aus einer Tabelle; fehlende Namen lesen die leere Zeichenkette —
// dieselbe Semantik wie os.Getenv am Produktionsaufruf.
func getenv(values map[string]string) func(string) string {
	return func(name string) string { return values[name] }
}

// vollständigeVerdrahtung trägt den Umgebungs-Stand, mit dem die
// Verdrahtung startet; die Tests entfernen oder verletzen je ein Element.
func vollständigeVerdrahtung() map[string]string {
	return map[string]string{
		"CDC_CAPTURE_DSN": "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_ADMIN_DSN":   "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_READER_DSN":  "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_SOURCE_ID":   "src-1",
		"CDC_PUBLICATION": "pub_1",
		"CDC_SLOT":        "slot_1",
		"CDC_TABLES":      "public.t1=tbl-1:sv-1",
	}
}

// TestConfigFromEnvOhneVorbedingung trägt die Verweigerungspfade der
// Umgebungs-Lese: je fehlende Vorbedingung endet die Verdrahtung über die
// Klasse `configuration` (bootstrap.ErrConfiguration), ohne zu starten —
// die Fehlerzeile nennt den ENV-Namen, nie den Wert.
func TestConfigFromEnvOhneVorbedingung(t *testing.T) {
	for _, missing := range []string{
		"CDC_CAPTURE_DSN", "CDC_ADMIN_DSN", "CDC_READER_DSN", "CDC_SOURCE_ID", "CDC_PUBLICATION", "CDC_SLOT",
	} {
		values := vollständigeVerdrahtung()
		delete(values, missing)
		cfg, err := bootstrap.ConfigFromEnv(getenv(values))
		if err == nil {
			t.Fatalf("%s fehlt: Verdrahtung startet ohne Vorbedingung", missing)
		}
		if !errors.Is(err, bootstrap.ErrConfiguration) {
			t.Fatalf("%s fehlt: Fehlerklasse configuration erwartet, erhalten: %v", missing, err)
		}
		if !strings.Contains(err.Error(), missing) {
			t.Fatalf("%s fehlt: Fehlerzeile nennt den ENV-Namen nicht: %v", missing, err)
		}
		if cfg.CaptureDSN != "" || cfg.AdminDSN != "" || cfg.ReaderDSN != "" || len(cfg.Tables) != 0 {
			t.Fatalf("%s fehlt: Verdrahtung liest trotz Verweigerung Teileingabe", missing)
		}
	}
}

// TestConfigFromEnvOhneTabellenAktivierung trägt die Grenze ohne
// Aktivierung: ein Feed ohne aktivierte Tabelle startet nicht.
func TestConfigFromEnvOhneTabellenAktivierung(t *testing.T) {
	for _, raw := range []string{"", "  "} {
		values := vollständigeVerdrahtung()
		values["CDC_TABLES"] = raw
		_, err := bootstrap.ConfigFromEnv(getenv(values))
		if err == nil {
			t.Fatalf("CDC_TABLES %q: Verdrahtung startet ohne Aktivierung", raw)
		}
		if !errors.Is(err, bootstrap.ErrConfiguration) {
			t.Fatalf("CDC_TABLES %q: Fehlerklasse configuration erwartet, erhalten: %v", raw, err)
		}
	}
}

// TestConfigFromEnvTabellenForm trägt die Form-Grenze der
// Aktivierungs-Liste: Einträge ohne Trennung, ohne Bindung, ohne
// Schema-Version-Kennung enden über die Klasse `configuration`.
func TestConfigFromEnvTabellenForm(t *testing.T) {
	for _, raw := range []string{
		"public.t1",                   // ohne Aktivierungs-Trennung
		"public.t1=tbl-1",             // ohne Schema-Version-Trennung
		"public.t1=:sv-1",             // ohne Tabellen-Kennung
		"public.t1=tbl-1:",            // ohne Schema-Version-Kennung
		"=tbl-1:sv-1",                 // ohne qualifizierten Tabellennamen
		"public.t1=tbl-1:sv-1,fehler", // zweiter Eintrag ohne Form
	} {
		values := vollständigeVerdrahtung()
		values["CDC_TABLES"] = raw
		_, err := bootstrap.ConfigFromEnv(getenv(values))
		if err == nil {
			t.Fatalf("CDC_TABLES %q: Verdrahtung startet trotz Formverletzung", raw)
		}
		if !errors.Is(err, bootstrap.ErrConfiguration) {
			t.Fatalf("CDC_TABLES %q: Fehlerklasse configuration erwartet, erhalten: %v", raw, err)
		}
	}
}

// TestConfigFromEnvLiestAktivierung trägt den Lese-Pfad: die
// Aktivierungs-Liste landet vollständig in den Tabellen-Bindungen.
func TestConfigFromEnvLiestAktivierung(t *testing.T) {
	values := vollständigeVerdrahtung()
	values["CDC_TABLES"] = "public.t1=tbl-1:sv-1, public.t2=tbl-2:sv-2"
	cfg, err := bootstrap.ConfigFromEnv(getenv(values))
	if err != nil {
		t.Fatalf("vollständige Vorbedingung: %v", err)
	}
	if cfg.CaptureDSN != values["CDC_CAPTURE_DSN"] ||
		cfg.AdminDSN != values["CDC_ADMIN_DSN"] ||
		cfg.ReaderDSN != values["CDC_READER_DSN"] ||
		string(cfg.Source) != "src-1" ||
		cfg.Publication != "pub_1" ||
		cfg.Slot != "slot_1" {
		t.Fatalf("Verdrahtungs-Eingabe unvollständig gelesen: %+v", cfg)
	}
	binding, exists := cfg.Tables["public.t2"]
	if !exists || binding.TableID != "tbl-2" || binding.SchemaVersion != "sv-2" {
		t.Fatalf("Tabellen-Bindung: %+v", cfg.Tables)
	}
	if len(cfg.Tables) != 2 {
		t.Fatalf("Tabellen-Bindungen: %d (Erwartung: 2)", len(cfg.Tables))
	}
}

// TestConfigFromEnvHTTPAddrBleibtOptional trägt die Additivitäts-Regel
// von `ADR-0057`/`slice-059`: `CDC_HTTP_ADDR` und die beiden Token-
// Umgebungsvariablen sind keine Vorbedingung — eine vollständige
// Verdrahtung ohne die drei neuen Namen startet unverändert
// (Regressionstest), gesetzte Werte landen unverändert in `Config`.
func TestConfigFromEnvHTTPAddrBleibtOptional(t *testing.T) {
	values := vollständigeVerdrahtung()
	cfg, err := bootstrap.ConfigFromEnv(getenv(values))
	if err != nil {
		t.Fatalf("bestehende Verdrahtung ohne CDC_HTTP_ADDR/CDC_API_TOKEN_*: %v", err)
	}
	if cfg.HTTPAddr != "" || cfg.APITokenReader != "" || cfg.APITokenAdmin != "" {
		t.Fatalf("Config trägt trotz ungesetzter Umgebung einen Wert: %+v", cfg)
	}

	values["CDC_HTTP_ADDR"] = ":8080"
	values["CDC_API_TOKEN_READER"] = "reader-token"
	values["CDC_API_TOKEN_ADMIN"] = "admin-token"
	cfg, err = bootstrap.ConfigFromEnv(getenv(values))
	if err != nil {
		t.Fatalf("vollständige Vorbedingung plus HTTP-Verdrahtung: %v", err)
	}
	if cfg.HTTPAddr != ":8080" || cfg.APITokenReader != "reader-token" || cfg.APITokenAdmin != "admin-token" {
		t.Fatalf("Config liest die drei neuen Umgebungsvariablen nicht vollständig: %+v", cfg)
	}
}

// TestConfigFromEnvGRPCAddrBleibtOptional trägt die Additivitäts-Regel von
// `ADR-0060` Teilfrage 6: `CDC_GRPC_ADDR` ist keine Vorbedingung — eine
// vollständige Verdrahtung ohne den Namen startet unverändert
// (Regressionstest), ein gesetzter Wert landet unverändert in `Config`.
func TestConfigFromEnvGRPCAddrBleibtOptional(t *testing.T) {
	values := vollständigeVerdrahtung()
	cfg, err := bootstrap.ConfigFromEnv(getenv(values))
	if err != nil {
		t.Fatalf("bestehende Verdrahtung ohne CDC_GRPC_ADDR: %v", err)
	}
	if cfg.GRPCAddr != "" {
		t.Fatalf("Config trägt trotz ungesetzter Umgebung eine gRPC-Adresse: %+v", cfg)
	}

	values["CDC_GRPC_ADDR"] = ":9090"
	cfg, err = bootstrap.ConfigFromEnv(getenv(values))
	if err != nil {
		t.Fatalf("vollständige Vorbedingung plus gRPC-Verdrahtung: %v", err)
	}
	if cfg.GRPCAddr != ":9090" {
		t.Fatalf("Config liest CDC_GRPC_ADDR nicht: %+v", cfg)
	}
}

// TestConfigFromEnvLogLevel trägt den Default und die erkannten Textformen
// von `CDC_LOG_LEVEL` (`LH-QA-OPS-004`): anders als die fünf
// Vorbedingungen oben bricht ein leerer oder nicht erkannter Wert die
// Verdrahtung nicht ab — er bleibt beim Default `Info`
// (`bootstrap.parseLogLevel`).
func TestConfigFromEnvLogLevel(t *testing.T) {
	for _, testcase := range []struct {
		raw      string
		expected slog.Level
	}{
		{raw: "", expected: slog.LevelInfo},
		{raw: "debug", expected: slog.LevelDebug},
		{raw: "DEBUG", expected: slog.LevelDebug},
		{raw: "warn", expected: slog.LevelWarn},
		{raw: "error", expected: slog.LevelError},
		{raw: "nicht-erkannt", expected: slog.LevelInfo},
	} {
		values := vollständigeVerdrahtung()
		if testcase.raw != "" {
			values["CDC_LOG_LEVEL"] = testcase.raw
		}
		cfg, err := bootstrap.ConfigFromEnv(getenv(values))
		if err != nil {
			t.Fatalf("CDC_LOG_LEVEL %q: vollständige Vorbedingung, aber Fehler: %v", testcase.raw, err)
		}
		if cfg.LogLevel != testcase.expected {
			t.Fatalf("CDC_LOG_LEVEL %q: Level %s, Erwartung %s", testcase.raw, cfg.LogLevel, testcase.expected)
		}
	}
}
