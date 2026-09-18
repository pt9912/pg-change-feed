package bootstrap

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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

// TestConfigFromFileLehntZugangsdatenAb trägt `ADR-0088` Festlegung 1/4 —
// die zugangsdaten-tragende Klasse: jeder ihrer sieben Schlüssel in der Datei
// bricht das Laden ab, unabhängig vom strikten Decoding. Die Fehlerzeile
// benennt den Schlüssel **und** den Grund; damit ist sie von der
// generischen „unbekannter Schlüssel"-Meldung des strikten Decodings
// unterscheidbar (`ADR-0088` Festlegung 4).
//
// Rot färbende Mutation: einen Schlüssel aus
// `forbiddenFileCredentialKeys` streichen — dann fällt er auf das strikte
// Decoding zurück, die Meldung verliert den Grund und genau dieser Fall
// wird rot.
func TestConfigFromFileLehntZugangsdatenAb(t *testing.T) {
	for _, key := range []string{
		"capture_dsn", "admin_dsn", "reader_dsn",
		"api_token_reader", "api_token_admin", "nats_url", "nats_stream_token",
	} {
		t.Run(key, func(t *testing.T) {
			path := writeConfigFile(t, key+": sollte-nicht-hier-stehen\n")
			_, err := ConfigFromFile(path)
			if err == nil {
				t.Fatalf("%s: Datei mit Zugangsdaten-Schlüssel lädt", key)
			}
			if !errors.Is(err, ErrConfiguration) {
				t.Fatalf("%s: Fehlerklasse configuration erwartet, erhalten: %v", key, err)
			}
			if !strings.Contains(err.Error(), key) {
				t.Fatalf("%s: Fehlerzeile benennt den Schlüssel nicht: %v", key, err)
			}
			if !strings.Contains(err.Error(), "Zugangsdaten bleiben env-var-exklusiv") {
				t.Fatalf("%s: Fehlerzeile benennt den Grund nicht — sie liest sich wie ein unbekannter Schlüssel: %v", key, err)
			}
		})
	}
}

// TestConfigFromFileUnbekannterSchluesselOhneZugangsdatenGrund trägt die
// Abgrenzung der zwei Zustände aus `ADR-0088` Festlegung 4: ein
// Tippfehler-Schlüssel bleibt „unbekannt" und trägt die Begründung der
// Zugangsdaten-Klasse **nicht** — die zwei Aussagen *unbekannt* und
// *unzulässig* sind an der Meldung unterscheidbar.
func TestConfigFromFileUnbekannterSchluesselOhneZugangsdatenGrund(t *testing.T) {
	path := writeConfigFile(t, "http_adr: ':8090'\n")
	_, err := ConfigFromFile(path)
	if !errors.Is(err, ErrConfiguration) {
		t.Fatalf("unbekannter Schlüssel: Fehlerklasse configuration erwartet, erhalten: %v", err)
	}
	if strings.Contains(err.Error(), "Zugangsdaten bleiben env-var-exklusiv") {
		t.Fatalf("unbekannter Schlüssel: Fehlerzeile trägt die Begründung der Zugangsdaten-Klasse: %v", err)
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
// `tables`-Merge-Precedence (`SPEC-016`): eine gesetzte `CDC_TABLES`
// schlägt die gesamte Datei-`tables`-Mapping vollständig, ohne
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

// dateiOhneAdressen trägt eine gültige Konfigurationsdatei ohne
// Oberflächen-Adressen — die Grundlage der Durchleitungs-Tests: jede
// gesetzte Oberflächen-Variable lässt sich damit eindeutig auf ihre
// Env-Herkunft zurückführen.
const dateiOhneAdressen = `
source_id: src-1
publication: pub_1
slot: slot_1
tables:
  public.t1:
    table_id: tbl-1
    schema_version: sv-1
`

// envMitDatei trägt die vollständigen Vorbedingungen inklusive
// `CDC_CONFIG_FILE`; `zusatz` legt die jeweils geprüfte Variable darüber.
func envMitDatei(path string, zusatz map[string]string) map[string]string {
	values := map[string]string{
		"CDC_CAPTURE_DSN": "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_ADMIN_DSN":   "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_READER_DSN":  "postgres://postgres:postgres@quelle:5432/cdc?sslmode=disable",
		"CDC_CONFIG_FILE": path,
	}
	for name, wert := range zusatz {
		values[name] = wert
	}
	return values
}

// TestMergeConfigOberflaechenVariablenAusEnvUnterDatei trägt `ADR-0088`
// Festlegung 3 — die Durchleitung: unter gesetzter `CDC_CONFIG_FILE` wirkt
// jede der fünf Oberflächen-Variablen aus ihrer Env-Herkunft. Je Variable
// ein eigener Fall mit genau einer gesetzten Variablen, damit der Wert
// eindeutig aus ihr stammt und nicht aus einer Nachbar-Variablen.
//
// Rot färbende Mutation: die Zuweisung der jeweiligen Variablen in
// `mergeConfig` streichen — dann bleibt ihr `Config`-Feld leer und genau
// dieser Fall wird rot.
func TestMergeConfigOberflaechenVariablenAusEnvUnterDatei(t *testing.T) {
	path := writeConfigFile(t, dateiOhneAdressen)
	cases := []struct {
		env  string
		wert string
		lies func(Config) string
	}{
		{"CDC_NATS_URL", "nats://nats:4222", func(c Config) string { return c.NatsURL }},
		{"CDC_HTTP_ADDR", ":8090", func(c Config) string { return c.HTTPAddr }},
		{"CDC_GRPC_ADDR", ":9090", func(c Config) string { return c.GRPCAddr }},
		{"CDC_API_TOKEN_READER", "token-reader-env", func(c Config) string { return c.APITokenReader }},
		{"CDC_API_TOKEN_ADMIN", "token-admin-env", func(c Config) string { return c.APITokenAdmin }},
	}
	// CDC_NATS_STREAM_TOKEN braucht zusätzlich eine gesetzte CDC_NATS_URL
	// (`validateNatsStreamTokenRequiresURL`) — ein eigener Fall statt eines
	// Eintrags in der obigen Liste, die je Fall nur eine Variable setzt.
	t.Run("CDC_NATS_STREAM_TOKEN", func(t *testing.T) {
		values := envMitDatei(path, map[string]string{
			"CDC_NATS_URL":          "nats://nats:4222",
			"CDC_NATS_STREAM_TOKEN": "s3cr3t",
		})
		cfg, err := ConfigFromEnvAndFile(func(name string) string { return values[name] })
		if err != nil {
			t.Fatalf("CDC_NATS_STREAM_TOKEN gesetzt unter Datei: %v", err)
		}
		if cfg.NatsStreamToken != "s3cr3t" {
			t.Fatalf("NatsStreamToken: Feld trägt %q, Erwartung aus der Env-Herkunft: %q", cfg.NatsStreamToken, "s3cr3t")
		}
	})
	for _, c := range cases {
		t.Run(c.env, func(t *testing.T) {
			values := envMitDatei(path, map[string]string{c.env: c.wert})
			cfg, err := ConfigFromEnvAndFile(func(name string) string { return values[name] })
			if err != nil {
				t.Fatalf("%s gesetzt unter Datei: %v", c.env, err)
			}
			if got := c.lies(cfg); got != c.wert {
				t.Fatalf("%s: Feld trägt %q, Erwartung aus der Env-Herkunft: %q", c.env, got, c.wert)
			}
		})
	}
}

// TestMergeConfigAdressFelderPrecedence trägt `ADR-0088` Festlegung 2/3 für
// die zwei neuen Datei-Felder in beiden Richtungen: eine gesetzte
// Umgebungsvariable schlägt den Datei-Wert, eine leere lässt ihn stehen.
// Der Datei-Wert ist in jedem Fall die Basis, auf der die Vorrang-Regel
// greift — deshalb steht er in jedem Fall in der Datei.
//
// Rot färbende Mutation: `overrideString` bei den zwei Feldern durch ein
// direktes `getenv(...)` ersetzen — dann fällt der Datei-Wert bei leerer
// Umgebungsvariable weg und der jeweilige Fall wird rot.
func TestMergeConfigAdressFelderPrecedence(t *testing.T) {
	path := writeConfigFile(t, `
source_id: src-1
publication: pub_1
slot: slot_1
tables:
  public.t1:
    table_id: tbl-1
    schema_version: sv-1
http_addr: ":8090"
grpc_addr: ":9090"
`)
	cases := []struct {
		name string
		env  string
		wert string
		lies func(Config) string
		want string
	}{
		{"http_addr ohne Env-Variable", "", "", func(c Config) string { return c.HTTPAddr }, ":8090"},
		{"http_addr mit gesetzter Env-Variable", "CDC_HTTP_ADDR", ":18090", func(c Config) string { return c.HTTPAddr }, ":18090"},
		{"grpc_addr ohne Env-Variable", "", "", func(c Config) string { return c.GRPCAddr }, ":9090"},
		{"grpc_addr mit gesetzter Env-Variable", "CDC_GRPC_ADDR", ":19090", func(c Config) string { return c.GRPCAddr }, ":19090"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			zusatz := map[string]string{}
			if c.env != "" {
				zusatz[c.env] = c.wert
			}
			values := envMitDatei(path, zusatz)
			cfg, err := ConfigFromEnvAndFile(func(name string) string { return values[name] })
			if err != nil {
				t.Fatalf("%s: %v", c.name, err)
			}
			if got := c.lies(cfg); got != c.want {
				t.Fatalf("%s: Feld trägt %q, Erwartung: %q", c.name, got, c.want)
			}
		})
	}
}

// TestMergeConfigOberflaechenUnterDateiAktiv trägt den zweiten Teil von
// `ADR-0088` Festlegung 3 — nicht „im Struct gesetzt", sondern an den
// Prädikaten, die über den Start der Oberflächen entscheiden. Der Zug fährt
// beide Richtungen: mit einer Adresse — aus der Datei **oder** aus der
// Umgebung — stehen die Prädikate auf „an"; trägt keine der beiden Quellen
// eine Adresse, bleiben alle Oberflächen aus.
//
// Zwei benannte Grenzen. (1) Gebunden ist `changeStreamEnabled`s Körper —
// `Run` entscheidet über die Streaming-Fähigkeit mit demselben Aufruf; die
// Leer-Prüfungen der fünf Felder wertet dieser Test dagegen als **eigene**
// Ausdrücke über `cfg` aus, die `!= ""`-Grenzen in `Run`s drei
// Start-Zweigen (`NatsURL`, `HTTPAddr`, `GRPCAddr`) selbst sind damit
// nicht gebunden. (2) Der Nachweis endet vor einem laufenden Server: `Run`
// konstruiert den Store vor dem HTTP-Server und braucht dafür eine
// erreichbare PostgreSQL-Instanz. Für die Datei-Herkunft trägt ihn kein
// Lauf mit realem Server: weder `compose.yaml` noch `tools/` noch `test/`
// setzen `CDC_CONFIG_FILE`.
func TestMergeConfigOberflaechenUnterDateiAktiv(t *testing.T) {
	t.Run("Adressen aus der Datei, Tokens und NATS aus der Umgebung: alle Oberflächen an", func(t *testing.T) {
		path := writeConfigFile(t, `
source_id: src-1
publication: pub_1
slot: slot_1
tables:
  public.t1:
    table_id: tbl-1
    schema_version: sv-1
http_addr: ":8090"
grpc_addr: ":9090"
`)
		values := envMitDatei(path, map[string]string{
			"CDC_NATS_URL":         "nats://nats:4222",
			"CDC_API_TOKEN_READER": "token-reader-env",
			"CDC_API_TOKEN_ADMIN":  "token-admin-env",
		})
		cfg, err := ConfigFromEnvAndFile(func(name string) string { return values[name] })
		if err != nil {
			t.Fatalf("Datei mit Adressen: %v", err)
		}
		if !changeStreamEnabled(cfg.GRPCAddr, cfg.HTTPAddr, natsStreamEnabled(cfg.NatsURL, cfg.NatsStreamToken)) {
			t.Fatalf("changeStreamEnabled(%q, %q) = false — der Broadcaster entstünde nicht", cfg.GRPCAddr, cfg.HTTPAddr)
		}
		if cfg.HTTPAddr == "" || cfg.GRPCAddr == "" {
			t.Fatalf("Adressen unter geladener Datei leer: http=%q grpc=%q — kein Server würde starten", cfg.HTTPAddr, cfg.GRPCAddr)
		}
		if cfg.NatsURL == "" {
			t.Fatal("NatsURL unter geladener Datei leer — die NATS-Verbindung entstünde nicht")
		}
		if cfg.APITokenReader == "" || cfg.APITokenAdmin == "" {
			t.Fatalf("Token-Klassen unter geladener Datei leer: reader=%q admin=%q", cfg.APITokenReader, cfg.APITokenAdmin)
		}
	})

	t.Run("Adressen nur aus der Umgebung: dieselben Prädikate an", func(t *testing.T) {
		path := writeConfigFile(t, dateiOhneAdressen)
		values := envMitDatei(path, map[string]string{
			"CDC_HTTP_ADDR": ":8090",
			"CDC_GRPC_ADDR": ":9090",
		})
		cfg, err := ConfigFromEnvAndFile(func(name string) string { return values[name] })
		if err != nil {
			t.Fatalf("Adressen aus der Umgebung: %v", err)
		}
		if !changeStreamEnabled(cfg.GRPCAddr, cfg.HTTPAddr, natsStreamEnabled(cfg.NatsURL, cfg.NatsStreamToken)) {
			t.Fatalf("changeStreamEnabled(%q, %q) = false — der Broadcaster entstünde nicht", cfg.GRPCAddr, cfg.HTTPAddr)
		}
		if cfg.HTTPAddr != ":8090" || cfg.GRPCAddr != ":9090" {
			t.Fatalf("Adressen aus der Env-Herkunft: http=%q grpc=%q", cfg.HTTPAddr, cfg.GRPCAddr)
		}
	})

	t.Run("keine Adresse in beiden Quellen: alle Oberflächen aus", func(t *testing.T) {
		path := writeConfigFile(t, dateiOhneAdressen)
		values := envMitDatei(path, nil)
		cfg, err := ConfigFromEnvAndFile(func(name string) string { return values[name] })
		if err != nil {
			t.Fatalf("Datei ohne Adressen: %v", err)
		}
		if changeStreamEnabled(cfg.GRPCAddr, cfg.HTTPAddr, natsStreamEnabled(cfg.NatsURL, cfg.NatsStreamToken)) {
			t.Fatalf("changeStreamEnabled(%q, %q) = true ohne Adresse in beiden Quellen", cfg.GRPCAddr, cfg.HTTPAddr)
		}
		if cfg.HTTPAddr != "" || cfg.GRPCAddr != "" || cfg.NatsURL != "" || cfg.NatsStreamToken != "" ||
			cfg.APITokenReader != "" || cfg.APITokenAdmin != "" {
			t.Fatalf("Oberflächen-Felder ohne Herkunft gesetzt: %+v", cfg)
		}
	})
}
