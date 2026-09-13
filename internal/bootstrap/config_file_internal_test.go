package bootstrap

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// writeConfigFile schreibt den YAML-Inhalt in eine temporäre Datei und
// liefert deren Pfad — Whitebox-Testhelfer, analog zu `vollständigeVerdrahtung`
// in `wiring_test.go`.
func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("Konfigurationsdatei nicht schreibbar: %v", err)
	}
	return path
}

// vollständigeEnvOhneDatei trägt dieselbe vollständige Env-Vorbedingung wie
// `vollständigeVerdrahtung` in `wiring_test.go` (eigene Kopie: das Paket
// dieser Datei ist `bootstrap`, nicht `bootstrap_test`), ohne
// `CDC_CONFIG_FILE`.
func vollständigeEnvOhneDatei() map[string]string {
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

// TestConfigFromFileStriktesDecoding trägt `ADR-0052` Entscheidung 1: ein
// unbekannter Schlüssel bricht das Laden über die Klasse `configuration` ab.
func TestConfigFromFileStriktesDecoding(t *testing.T) {
	path := writeConfigFile(t, "unbekannter_schluessel: wert\n")
	_, err := ConfigFromFile(path)
	if err == nil {
		t.Fatal("unbekannter Schlüssel: Datei lädt trotz striktem Decoding")
	}
	if !errors.Is(err, ErrConfiguration) {
		t.Fatalf("unbekannter Schlüssel: Fehlerklasse configuration erwartet, erhalten: %v", err)
	}
}

// TestConfigFromFileLehntDSNAb trägt `ADR-0052` Entscheidung 6 — die
// wichtigste Einzelentscheidung: jeder der drei DSN-Schlüssel in der Datei
// bricht das Laden ab, unabhängig vom strikten Decoding.
func TestConfigFromFileLehntDSNAb(t *testing.T) {
	for _, key := range []string{"capture_dsn", "admin_dsn", "reader_dsn"} {
		path := writeConfigFile(t, key+": postgres://sollte-nicht-hier-stehen\n")
		_, err := ConfigFromFile(path)
		if err == nil {
			t.Fatalf("%s: Datei mit DSN-Schlüssel lädt", key)
		}
		if !errors.Is(err, ErrConfiguration) {
			t.Fatalf("%s: Fehlerklasse configuration erwartet, erhalten: %v", key, err)
		}
	}
}

// TestConfigFromFileNichtLesbar trägt die Fehlerklasse `configuration` für
// einen nicht existierenden Pfad (`ADR-0052` Entscheidung 5: kein stiller
// Fallback auf Env-only, wenn der Pfad explizit angegeben, aber nicht lesbar
// ist).
func TestConfigFromFileNichtLesbar(t *testing.T) {
	_, err := ConfigFromFile(filepath.Join(t.TempDir(), "existiert-nicht.yaml"))
	if !errors.Is(err, ErrConfiguration) {
		t.Fatalf("Fehlerklasse configuration erwartet, erhalten: %v", err)
	}
}

// TestConfigFromFileLeereDatei trägt die Grenze: eine leere Datei ist kein
// Fehler, sie trägt nur keine Datei-Basis.
func TestConfigFromFileLeereDatei(t *testing.T) {
	path := writeConfigFile(t, "")
	cfg, err := ConfigFromFile(path)
	if err != nil {
		t.Fatalf("leere Datei: unerwarteter Fehler: %v", err)
	}
	if cfg.SourceID != "" || cfg.Publication != "" || cfg.Slot != "" || len(cfg.Tables) != 0 {
		t.Fatalf("leere Datei: unerwartete Felder gesetzt: %+v", cfg)
	}
}

// TestConfigFromFileTabellenMapping trägt das eigene DoD-Item: eine Datei
// mit mehreren Tabellen-Aktivierungen als YAML-Mapping lädt vollständig.
func TestConfigFromFileTabellenMapping(t *testing.T) {
	path := writeConfigFile(t, `
source_id: src-1
publication: pub_1
slot: slot_1
tables:
  public.t1:
    table_id: tbl-1
    schema_version: sv-1
  public.t2:
    table_id: tbl-2
    schema_version: sv-2
`)
	cfg, err := ConfigFromFile(path)
	if err != nil {
		t.Fatalf("vollständige Datei: unerwarteter Fehler: %v", err)
	}
	if len(cfg.Tables) != 2 {
		t.Fatalf("Tabellen-Bindungen: %d (Erwartung: 2)", len(cfg.Tables))
	}
	binding, exists := cfg.Tables["public.t2"]
	if !exists || binding.TableID != "tbl-2" || binding.SchemaVersion != "sv-2" {
		t.Fatalf("Tabellen-Bindung public.t2: %+v", cfg.Tables)
	}
}

// TestConfigFromEnvAndFileLeereEnvVariable trägt die Regressions-Pflicht
// (`ADR-0052` Entscheidung 3): ein bestehender Env-only-Aufruf ohne
// `CDC_CONFIG_FILE` zeigt identisches Verhalten zu `ConfigFromEnv` — bis auf
// die Fehlerklasse selbst wird hier der volle Config-Wert verglichen.
func TestConfigFromEnvAndFileLeereEnvVariable(t *testing.T) {
	values := vollständigeEnvOhneDatei()
	getenv := func(name string) string { return values[name] }

	viaEnvAndFile, err := ConfigFromEnvAndFile(getenv)
	if err != nil {
		t.Fatalf("ConfigFromEnvAndFile ohne CDC_CONFIG_FILE: %v", err)
	}
	viaEnvOnly, err := ConfigFromEnv(getenv)
	if err != nil {
		t.Fatalf("ConfigFromEnv: %v", err)
	}
	if viaEnvAndFile.CaptureDSN != viaEnvOnly.CaptureDSN ||
		viaEnvAndFile.AdminDSN != viaEnvOnly.AdminDSN ||
		viaEnvAndFile.ReaderDSN != viaEnvOnly.ReaderDSN ||
		viaEnvAndFile.Source != viaEnvOnly.Source ||
		viaEnvAndFile.Publication != viaEnvOnly.Publication ||
		viaEnvAndFile.Slot != viaEnvOnly.Slot ||
		viaEnvAndFile.LogLevel != viaEnvOnly.LogLevel ||
		len(viaEnvAndFile.Tables) != len(viaEnvOnly.Tables) {
		t.Fatalf("ConfigFromEnvAndFile weicht ohne CDC_CONFIG_FILE vom Env-only-Pfad ab: %+v vs %+v", viaEnvAndFile, viaEnvOnly)
	}
}

// TestConfigFromEnvAndFileDateiNichtLadbar trägt `ADR-0052` Entscheidung 5:
// eine gesetzte, aber nicht lesbare `CDC_CONFIG_FILE` bricht ab — kein
// stiller Fallback auf Env-only.
func TestConfigFromEnvAndFileDateiNichtLadbar(t *testing.T) {
	values := vollständigeEnvOhneDatei()
	values["CDC_CONFIG_FILE"] = filepath.Join(t.TempDir(), "existiert-nicht.yaml")
	getenv := func(name string) string { return values[name] }

	_, err := ConfigFromEnvAndFile(getenv)
	if !errors.Is(err, ErrConfiguration) {
		t.Fatalf("Fehlerklasse configuration erwartet, erhalten: %v", err)
	}
}

// TestMergeConfigEinzelnesFeldPerEnv trägt das zentrale DoD-Item: nur ein
// einzelnes Feld (`slot`) wird per Env-Var überschrieben, alle übrigen
// Felder stammen unverändert aus der Datei (`ADR-0052` Entscheidung 2).
func TestMergeConfigEinzelnesFeldPerEnv(t *testing.T) {
	path := writeConfigFile(t, `
source_id: src-datei
publication: pub-datei
slot: slot-datei
tables:
  public.t1:
    table_id: tbl-1
    schema_version: sv-1
`)
	values := map[string]string{
		"CDC_CAPTURE_DSN": "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_ADMIN_DSN":   "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_READER_DSN":  "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_CONFIG_FILE": path,
		// Nur CDC_SLOT ist gesetzt — der einzige Wert, der die Datei
		// überschreiben soll.
		"CDC_SLOT": "slot-env",
	}
	getenv := func(name string) string { return values[name] }

	cfg, err := ConfigFromEnvAndFile(getenv)
	if err != nil {
		t.Fatalf("Merge mit einzelnem Env-Override: %v", err)
	}
	if string(cfg.Source) != "src-datei" {
		t.Fatalf("Source: %q, Erwartung aus Datei: src-datei", cfg.Source)
	}
	if cfg.Publication != "pub-datei" {
		t.Fatalf("Publication: %q, Erwartung aus Datei: pub-datei", cfg.Publication)
	}
	if cfg.Slot != "slot-env" {
		t.Fatalf("Slot: %q, Erwartung aus Env-Override: slot-env", cfg.Slot)
	}
	binding, exists := cfg.Tables["public.t1"]
	if !exists || binding.TableID != "tbl-1" || binding.SchemaVersion != "sv-1" {
		t.Fatalf("Tabellen-Bindung aus Datei: %+v", cfg.Tables)
	}
}

// TestMergeConfigTabellenCDCTablesSchlaegtDatei trägt die
// Implementer-Entscheidung zu §6 Risiko 1 des Slice-Plans: eine gesetzte
// `CDC_TABLES` schlägt die gesamte Datei-`tables`-Mapping vollständig, ohne
// Vermischung — die Datei-Tabelle `public.t-datei` erscheint im Ergebnis
// nicht.
func TestMergeConfigTabellenCDCTablesSchlaegtDatei(t *testing.T) {
	path := writeConfigFile(t, `
source_id: src-1
publication: pub-1
slot: slot-1
tables:
  public.t-datei:
    table_id: tbl-datei
    schema_version: sv-datei
`)
	values := map[string]string{
		"CDC_CAPTURE_DSN": "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_ADMIN_DSN":   "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_READER_DSN":  "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_CONFIG_FILE": path,
		"CDC_TABLES":      "public.t-env=tbl-env:sv-env",
	}
	getenv := func(name string) string { return values[name] }

	cfg, err := ConfigFromEnvAndFile(getenv)
	if err != nil {
		t.Fatalf("Merge mit CDC_TABLES gesetzt: %v", err)
	}
	if len(cfg.Tables) != 1 {
		t.Fatalf("Tabellen-Bindungen: %d (Erwartung: 1 — CDC_TABLES schlägt die Datei vollständig)", len(cfg.Tables))
	}
	if _, exists := cfg.Tables["public.t-datei"]; exists {
		t.Fatal("Datei-Tabelle public.t-datei ist trotz gesetzter CDC_TABLES im Ergebnis — Vermischung statt vollständigem Vorrang")
	}
	binding, exists := cfg.Tables["public.t-env"]
	if !exists || binding.TableID != "tbl-env" || binding.SchemaVersion != "sv-env" {
		t.Fatalf("Env-Tabelle public.t-env: %+v", cfg.Tables)
	}
}

// TestMergeConfigFehlendePflichtfelder trägt die Validierung nach dem
// Merge: fehlt ein Pflichtfeld in beiden Quellen (hier: Datei ohne
// `tables`, keine `CDC_TABLES`), bricht der Aufruf über die Klasse
// `configuration` ab.
func TestMergeConfigFehlendePflichtfelder(t *testing.T) {
	path := writeConfigFile(t, `
source_id: src-1
publication: pub-1
slot: slot-1
`)
	values := map[string]string{
		"CDC_CAPTURE_DSN": "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_ADMIN_DSN":   "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_READER_DSN":  "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_CONFIG_FILE": path,
	}
	getenv := func(name string) string { return values[name] }

	_, err := ConfigFromEnvAndFile(getenv)
	if !errors.Is(err, ErrConfiguration) {
		t.Fatalf("Fehlerklasse configuration erwartet, erhalten: %v", err)
	}
}

// TestMergeConfigLogLevelUndWALRetentionAusDatei trägt die beiden
// zusätzlichen Datei-Felder ohne Env-Gegenstück: `log_level` mit
// Env-Override, `wal_retention_warn_bytes`/`wal_retention_error_bytes` ohne
// einen solchen (die Datei ist ihre einzige Herkunft neben dem
// SPEC-013-Startwert selbst).
func TestMergeConfigLogLevelUndWALRetentionAusDatei(t *testing.T) {
	path := writeConfigFile(t, `
source_id: src-1
publication: pub-1
slot: slot-1
tables:
  public.t1:
    table_id: tbl-1
    schema_version: sv-1
log_level: warn
wal_retention_warn_bytes: 111
wal_retention_error_bytes: 222
`)
	values := map[string]string{
		"CDC_CAPTURE_DSN": "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_ADMIN_DSN":   "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_READER_DSN":  "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_CONFIG_FILE": path,
	}
	getenv := func(name string) string { return values[name] }

	cfg, err := ConfigFromEnvAndFile(getenv)
	if err != nil {
		t.Fatalf("Merge mit log_level/wal_retention aus Datei: %v", err)
	}
	if cfg.LogLevel.String() != "WARN" {
		t.Fatalf("LogLevel: %s, Erwartung: WARN", cfg.LogLevel)
	}
	if cfg.WALRetentionWarnBytes != 111 || cfg.WALRetentionErrorBytes != 222 {
		t.Fatalf("WAL-Retention-Overrides: warn=%d error=%d, Erwartung: 111/222", cfg.WALRetentionWarnBytes, cfg.WALRetentionErrorBytes)
	}
}
