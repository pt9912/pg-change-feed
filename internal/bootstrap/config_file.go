// config_file.go trägt den optionalen Datei-Ladepfad der Verdrahtung
// (`ADR-0052`): eine YAML-Konfigurationsdatei ergänzt `ConfigFromEnv`
// additiv, mit Umgebungsvariable-schlägt-Datei-Feld-für-Feld-Precedence
// (`ADR-0052` Entscheidung 2) und env-var-exklusiven DSNs (`ADR-0052`
// Entscheidung 6). Der Zugriffsweg ist ausschließlich `CDC_CONFIG_FILE`
// (`ADR-0052` Entscheidung 5) — ein `--config`-CLI-Flag ist Folgepflicht,
// nicht Teil dieses Standes.
package bootstrap

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// envConfigFile trägt den Pfad der optionalen Konfigurationsdatei
// (`ADR-0052` Entscheidung 5); leer/unbenannt bedeutet: kein Dateizugriff,
// der Env-only-Pfad (`ConfigFromEnv`) bleibt unverändert Default.
const envConfigFile = "CDC_CONFIG_FILE"

// forbiddenFileDSNKeys trägt die drei Schlüssel, die in der
// Konfigurationsdatei nicht vorkommen dürfen (`ADR-0052` Entscheidung 6 —
// die wichtigste Einzelentscheidung dieser ADR: Secrets bleiben
// env-var-exklusiv, `LH-QA-SEC-001`/`002`). Ein Treffer bricht das Laden
// mit einer eigenen, den Secret-Grund benennenden Fehlerzeile ab, statt
// nur als generischer „unbekannter Schlüssel" des strikten Decodings unten
// zu erscheinen.
var forbiddenFileDSNKeys = []string{"capture_dsn", "admin_dsn", "reader_dsn"}

// fileTableBinding trägt eine einzelne Tabellen-Aktivierung der
// Konfigurationsdatei als YAML-Mapping (`ADR-0052` Entscheidung 6) — anders
// als die `CDC_TABLES`-Zeichenkettenform (`parseTables`) erlaubt die
// Datei-Form Kommentare und Gruppierung je Tabelle.
type fileTableBinding struct {
	TableID       string `yaml:"table_id"`
	SchemaVersion string `yaml:"schema_version"`
}

// fileConfig trägt die in der Konfigurationsdatei zulässigen Felder
// (`ADR-0052` Entscheidung 6): nur die nicht credential-tragenden Felder.
// `CaptureDSN`/`AdminDSN`/`ReaderDSN` bleiben env-var-exklusiv und haben
// hier bewusst kein Gegenstück — ein Treffer auf einen der drei
// DSN-Schlüssel wird vor dem Decoding in diesen Typ abgefangen
// (`forbiddenFileDSNKeys`).
type fileConfig struct {
	SourceID    string                      `yaml:"source_id"`
	Publication string                      `yaml:"publication"`
	Slot        string                      `yaml:"slot"`
	Tables      map[string]fileTableBinding `yaml:"tables"`
	LogLevel    string                      `yaml:"log_level"`
	// WALRetentionWarnBytes und WALRetentionErrorBytes tragen dieselben
	// SPEC-013-Overrides wie `Config.WALRetentionWarnBytes`/
	// `WALRetentionErrorBytes` — ohne Umgebungs-Gegenstück (siehe dort):
	// die Datei ist für diese beiden Felder die einzige Herkunft neben dem
	// SPEC-013-Startwert selbst.
	WALRetentionWarnBytes  int64 `yaml:"wal_retention_warn_bytes"`
	WALRetentionErrorBytes int64 `yaml:"wal_retention_error_bytes"`
}

// ConfigFromFile lädt die optionale Konfigurationsdatei. Einer der drei
// DSN-Schlüssel (`forbiddenFileDSNKeys`) wird auf dem roh eingelesenen
// Dokument geprüft und liefert — im aktuellen Kontrollfluss immer zuerst —
// `ErrConfiguration` mit einer eigenen, den Secret-Grund benennenden
// Fehlerzeile (`ADR-0052` Entscheidung 6). Jeder andere unbekannte
// Schlüssel liefert `ErrConfiguration` über striktes YAML-Decoding
// (`yaml.Decoder.KnownFields(true)`, `ADR-0052` Entscheidung 1); da
// `fileConfig` auch keines der drei DSN-Felder deklariert, würde
// `KnownFields` sie ebenfalls ablehnen, sollte der explizite Check je
// entfallen — im jetzigen Kontrollfluss ist dieser Pfad für die drei
// DSN-Schlüssel nicht erreichbar, weil der explizite Check vorher
// zurückkehrt. Eine leere Datei (kein YAML-Dokument, z. B. nur
// Kommentare) liefert die Nullwerte zurück, keinen Fehler — sie trägt
// dann keine Datei-Basis, jedes Feld bleibt der Env-var-Seite von
// `mergeConfig` überlassen.
func ConfigFromFile(path string) (fileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return fileConfig{}, fmt.Errorf("%w: Konfigurationsdatei %q nicht lesbar: %v", ErrConfiguration, path, err)
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return fileConfig{}, fmt.Errorf("%w: Konfigurationsdatei %q trägt kein gültiges YAML: %v", ErrConfiguration, path, err)
	}
	for _, key := range forbiddenFileDSNKeys {
		if _, found := raw[key]; found {
			return fileConfig{}, fmt.Errorf("%w: Konfigurationsdatei %q trägt den Schlüssel %q — DSNs bleiben env-var-exklusiv (ADR-0052 Entscheidung 6)", ErrConfiguration, path, key)
		}
	}

	var parsed fileConfig
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&parsed); err != nil {
		if errors.Is(err, io.EOF) {
			return fileConfig{}, nil
		}
		return fileConfig{}, fmt.Errorf("%w: Konfigurationsdatei %q: %v", ErrConfiguration, path, err)
	}
	return parsed, nil
}

// ConfigFromEnvAndFile ist der verdrahtete Zugriffsweg (`ADR-0052`
// Entscheidung 5, verdrahtet in `cmd/pg-change-feed/main.go`):
// `CDC_CONFIG_FILE` leer/unbenannt trägt exakt den heutigen
// `ConfigFromEnv`-Pfad, unverändert — jeder bestehende Env-only-Aufruf
// zeigt identisches Verhalten wie vor diesem Slice. Ist die Variable
// gesetzt, aber die Datei unter diesem Pfad nicht ladbar, ist das ein
// `ErrConfiguration`-Fehler — kein stiller Fallback auf Env-only
// (`ADR-0052` Entscheidung 5).
func ConfigFromEnvAndFile(getenv func(string) string) (Config, error) {
	path := getenv(envConfigFile)
	if path == "" {
		return ConfigFromEnv(getenv)
	}
	file, err := ConfigFromFile(path)
	if err != nil {
		return Config{}, err
	}
	return mergeConfig(file, getenv)
}

// overrideString trägt die Feld-für-Feld-Precedence (`ADR-0052`
// Entscheidung 2): eine gesetzte Umgebungsvariable überschreibt den
// Datei-Wert; eine leere Umgebungsvariable lässt den Datei-Wert
// unverändert.
func overrideString(fileValue, envValue string) string {
	if envValue != "" {
		return envValue
	}
	return fileValue
}

// mergeConfig überschreibt die Datei-Basis Feld für Feld mit jeder
// gesetzten Umgebungsvariable (`ADR-0052` Entscheidung 2) und validiert die
// Vorbedingung danach — dieselben sechs Pflichtfelder wie `ConfigFromEnv`,
// hier über beide Quellen hinweg geprüft. Die drei DSNs bleiben
// env-var-exklusiv (`ADR-0052` Entscheidung 6): sie haben kein
// Datei-Gegenstück und werden unverändert direkt gelesen, wie in
// `ConfigFromEnv`.
func mergeConfig(file fileConfig, getenv func(string) string) (Config, error) {
	cfg := Config{
		CaptureDSN: getenv(envCaptureDSN),
		AdminDSN:   getenv(envAdminDSN),
		ReaderDSN:  getenv(envReaderDSN),
	}
	if cfg.CaptureDSN == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt", ErrConfiguration, envCaptureDSN)
	}
	if cfg.AdminDSN == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt", ErrConfiguration, envAdminDSN)
	}
	if cfg.ReaderDSN == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt", ErrConfiguration, envReaderDSN)
	}

	cfg.Source = model.SourceID(overrideString(file.SourceID, getenv(envSource)))
	if cfg.Source == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt (weder Umgebungsvariable noch source_id in der Konfigurationsdatei)", ErrConfiguration, envSource)
	}
	cfg.Publication = overrideString(file.Publication, getenv(envPublication))
	if cfg.Publication == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt (weder Umgebungsvariable noch publication in der Konfigurationsdatei)", ErrConfiguration, envPublication)
	}
	cfg.Slot = overrideString(file.Slot, getenv(envSlot))
	if cfg.Slot == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt (weder Umgebungsvariable noch slot in der Konfigurationsdatei)", ErrConfiguration, envSlot)
	}

	tables, err := mergeTables(file.Tables, getenv(envTables))
	if err != nil {
		return Config{}, err
	}
	cfg.Tables = tables

	logLevelRaw := getenv(envLogLevel)
	if logLevelRaw == "" {
		logLevelRaw = file.LogLevel
	}
	cfg.LogLevel = parseLogLevel(logLevelRaw)

	cfg.WALRetentionWarnBytes = file.WALRetentionWarnBytes
	cfg.WALRetentionErrorBytes = file.WALRetentionErrorBytes

	return cfg, nil
}

// mergeTables trägt die `tables`-Merge-Precedence (`SPEC-016`): eine
// gesetzte `CDC_TABLES` schlägt die gesamte Datei-`tables`-Mapping
// vollständig — keine Vermischung einzelner Tabellen aus beiden Quellen
// innerhalb derselben Liste. Die Feld-für-Feld-Precedence aus `ADR-0052`
// Entscheidung 2 behandelt `tables` damit als ein Feld (die ganze
// Aktivierungsliste), nicht als Menge einzeln überschreibbarer Einträge.
func mergeTables(fileTables map[string]fileTableBinding, envRaw string) (map[string]mapper.TableBinding, error) {
	if strings.TrimSpace(envRaw) != "" {
		return parseTables(envRaw)
	}
	if len(fileTables) == 0 {
		return nil, fmt.Errorf("%w: %s trägt keine Tabellen-Aktivierung (weder Umgebungsvariable noch tables in der Konfigurationsdatei)", ErrConfiguration, envTables)
	}
	tables := make(map[string]mapper.TableBinding, len(fileTables))
	for qualified, binding := range fileTables {
		if binding.TableID == "" || binding.SchemaVersion == "" {
			return nil, fmt.Errorf("%w: Tabellen-Aktivierung %q in der Konfigurationsdatei trägt nicht beide Felder table_id und schema_version", ErrConfiguration, qualified)
		}
		tables[qualified] = mapper.TableBinding{
			TableID:       model.SourceTableID(binding.TableID),
			SchemaVersion: model.SchemaVersionID(binding.SchemaVersion),
		}
	}
	return tables, nil
}
