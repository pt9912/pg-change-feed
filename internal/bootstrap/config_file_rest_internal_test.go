package bootstrap

import (
	"errors"
	"strings"
	"testing"
)

// Diese Datei trägt den netzlos prüfbaren Rest der Datei-Konfiguration
// (`ADR-0052`), der neben den Happy Paths in
// `config_file_internal_test.go` offenblieb: die Fehlerzweige von
// `ConfigFromFile` (ungültiges YAML), von `mergeConfig` (die sechs
// Pflichtfelder über beide Quellen hinweg) und von `mergeTables` (eine
// Datei-Aktivierung ohne beide Felder). Kein Fall braucht einen Dienst —
// dieselbe Fläche, die `ADR-0082` als Cluster D2 führt.

// TestConfigFromFileLehntUngueltigesYAMLAb trägt die Fehlerklasse
// `configuration` für ein Dokument, das die YAML-Grammatik bricht — sie
// greift vor dem strikten Decoding und vor dem DSN-Schlüssel-Check
// (`ADR-0052` Entscheidung 1). Ein unvollständiger Flow-Sequenz-Start ist
// ein solcher Grammatik-Bruch: er scheitert am `yaml.Unmarshal` der roh
// eingelesenen Bytes, nicht erst am `KnownFields`-Decoder weiter unten.
//
// Rot färbende Mutation: den `yaml.Unmarshal`-Fehlerzweig auf `return
// fileConfig{}, nil` umstellen — dann lädt diese Datei still als
// Nullwert und der Test bricht.
func TestConfigFromFileLehntUngueltigesYAMLAb(t *testing.T) {
	path := writeConfigFile(t, "source_id: [1, 2\n")
	_, err := ConfigFromFile(path)
	if err == nil {
		t.Fatal("ungültiges YAML: Datei lädt ohne Fehler")
	}
	if !errors.Is(err, ErrConfiguration) {
		t.Fatalf("ungültiges YAML: Fehlerklasse configuration erwartet, erhalten: %v", err)
	}
	if !strings.Contains(err.Error(), "kein gültiges YAML") {
		t.Fatalf("ungültiges YAML: Fehlerzeile benennt den YAML-Grund nicht: %v", err)
	}
}

// TestMergeConfigDSNsBleibenEnvVarExklusiv trägt `ADR-0052` Entscheidung 6
// — die wichtigste Einzelentscheidung dieser ADR (`LH-QA-SEC-001`/`002`):
// die drei DSNs haben **kein** Datei-Gegenstück. Eine vollständige Datei
// füllt sie deshalb nicht auf; fehlt die jeweilige Umgebungsvariable,
// bricht der Merge mit einer Zeile ab, die genau diese Variable benennt —
// statt still auf einen Datei-Wert auszuweichen, den es nicht geben darf.
//
// Rot färbende Mutation: in `mergeConfig` die drei
// DSN-Zuweisungen aus `file`-Feldern speisen (die es nicht gibt) oder ihre
// Leer-Prüfung entfernen — dann läuft der Merge mit leerer DSN durch und
// der Test bricht.
func TestMergeConfigDSNsBleibenEnvVarExklusiv(t *testing.T) {
	path := writeConfigFile(t, `
source_id: src-datei
publication: pub-datei
slot: slot-datei
tables:
  public.t1:
    table_id: tbl-1
    schema_version: sv-1
`)

	for _, missing := range []string{"CDC_CAPTURE_DSN", "CDC_ADMIN_DSN", "CDC_READER_DSN"} {
		t.Run(missing+" fehlt", func(t *testing.T) {
			values := map[string]string{
				"CDC_CAPTURE_DSN": "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
				"CDC_ADMIN_DSN":   "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
				"CDC_READER_DSN":  "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
				"CDC_CONFIG_FILE": path,
			}
			delete(values, missing)
			getenv := func(name string) string { return values[name] }

			_, err := ConfigFromEnvAndFile(getenv)
			if !errors.Is(err, ErrConfiguration) {
				t.Fatalf("%s aus Env entfernt: Fehlerklasse configuration erwartet, erhalten: %v", missing, err)
			}
			if !strings.Contains(err.Error(), missing) {
				t.Fatalf("%s aus Env entfernt: Fehlerzeile benennt %s nicht: %v", missing, missing, err)
			}
		})
	}
}

// TestMergeConfigPflichtfelderAusBeidenQuellenBelegt trägt die
// Merge-Validierung für die drei nicht-credential-tragenden Pflichtfelder
// (`ADR-0052` Entscheidung 2): `source_id`/`publication`/`slot` dürfen aus
// der Datei kommen — fehlen sie in **beiden** Quellen, benennt die
// Fehlerzeile die Umgebungsvariable, die die Datei hätte ergänzen können.
// Der Fall ist von `TestMergeConfigFehlendePflichtfelder` (dort fehlt die
// gesamte Tabellen-Aktivierung) verschieden: hier ist die Aktivierung
// vorhanden und genau ein Feld gerissen.
func TestMergeConfigPflichtfelderAusBeidenQuellenBelegt(t *testing.T) {
	cases := []struct {
		env    string
		yaml   string
		feldNA string
	}{
		{"CDC_SOURCE_ID", "publication: pub-1\nslot: slot-1\n", "source_id"},
		{"CDC_PUBLICATION", "source_id: src-1\nslot: slot-1\n", "publication"},
		{"CDC_SLOT", "source_id: src-1\npublication: pub-1\n", "slot"},
	}

	for _, c := range cases {
		t.Run(c.env+" fehlt in beiden Quellen", func(t *testing.T) {
			path := writeConfigFile(t, c.yaml+`tables:
  public.t1:
    table_id: tbl-1
    schema_version: sv-1
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
				t.Fatalf("%s in beiden Quellen leer: Fehlerklasse configuration erwartet, erhalten: %v", c.env, err)
			}
			if !strings.Contains(err.Error(), c.env) {
				t.Fatalf("%s in beiden Quellen leer: Fehlerzeile benennt %s nicht: %v", c.env, c.env, err)
			}
			if !strings.Contains(err.Error(), c.feldNA) {
				t.Fatalf("%s in beiden Quellen leer: Fehlerzeile benennt das Datei-Feld %q nicht: %v", c.env, c.feldNA, err)
			}
		})
	}
}

// TestMergeConfigDateiAktivierungBrauchtBeideFelder trägt die
// Form-Vorbedingung der Datei-Aktivierung (`SPEC-016`, `ADR-0052`
// Entscheidung 6): die YAML-Mapping-Form hat gegenüber `CDC_TABLES` zwei
// Felder, und eine Aktivierung ohne beide ist keine — sie wird nicht als
// halbe Bindung übernommen, sondern benennt den qualifizierten Namen und
// bricht ab. Ohne diese Prüfung entstünde eine Bindung mit leerer
// `SchemaVersion`, die der `Assembler` nicht interpretieren kann.
//
// Rot färbende Mutation: die Prüfung `binding.TableID == "" ||
// binding.SchemaVersion == ""` entfernen — dann lädt die halbe Aktivierung
// und der Test bricht.
func TestMergeConfigDateiAktivierungBrauchtBeideFelder(t *testing.T) {
	path := writeConfigFile(t, `
source_id: src-1
publication: pub-1
slot: slot-1
tables:
  public.halbe_aktivierung:
    table_id: tbl-1
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
		t.Fatalf("Aktivierung ohne schema_version: Fehlerklasse configuration erwartet, erhalten: %v", err)
	}
	if !strings.Contains(err.Error(), "public.halbe_aktivierung") {
		t.Fatalf("Aktivierung ohne schema_version: Fehlerzeile benennt den qualifizierten Namen nicht: %v", err)
	}
}
