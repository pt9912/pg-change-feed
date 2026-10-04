package bootstrap_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/bootstrap"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
)

// apiTokensZugriffswege trägt die zwei Zugriffswege der Konfiguration: den
// reinen Umgebungsweg und den verdrahteten Weg mit geladener Datei, den der
// Feed-Container fährt. Beide lesen die Token-Variablen über dieselbe Funktion.
func apiTokensZugriffswege(t *testing.T) map[string]func(map[string]string) (bootstrap.Config, error) {
	t.Helper()
	datei := filepath.Join(t.TempDir(), "cdc.yaml")
	inhalt := "source_id: src-datei\npublication: pub-datei\nslot: slot-datei\ntables:\n  public.t1:\n    table_id: tbl-1\n    schema_version: sv-1\n"
	if err := os.WriteFile(datei, []byte(inhalt), 0o600); err != nil {
		t.Fatalf("Konfigurationsdatei schreiben: %v", err)
	}
	return map[string]func(map[string]string) (bootstrap.Config, error){
		"ConfigFromEnv": func(env map[string]string) (bootstrap.Config, error) {
			return bootstrap.ConfigFromEnv(getenv(env))
		},
		"ConfigFromEnvAndFile": func(env map[string]string) (bootstrap.Config, error) {
			env["CDC_CONFIG_FILE"] = datei
			return bootstrap.ConfigFromEnvAndFile(getenv(env))
		},
	}
}

// TestAPITokenListenWerdenGelesen trägt `LH-FA-SST-012` Happy Path und
// Boundary der Konfiguration: die Plural-Variablen liefern je Klasse die
// Liste, der Singular bleibt unverändert daneben, ein Singular mit Komma wird
// nicht zerschnitten, eine leere Zeichenkette gilt wie ungesetzt — auf beiden
// Zugriffswegen, auch unter geladener Datei.
// Rot färbende Mutation: in `mergeConfig` den Aufruf `applyAPITokens` durch
// die zwei Singular-Zuweisungen ersetzen — der Fall von `ConfigFromEnvAndFile`
// trägt keine Liste mehr.
func TestAPITokenListenWerdenGelesen(t *testing.T) {
	cases := []struct {
		name         string
		env          map[string]string
		wantReader   []string
		wantAdmin    []string
		wantSingular [2]string
	}{
		{"ungesetzt", nil, nil, nil, [2]string{"", ""}},
		{"leere Zeichenkette gilt als ungesetzt",
			map[string]string{"CDC_API_TOKENS_READER": "", "CDC_API_TOKENS_ADMIN": ""}, nil, nil, [2]string{"", ""}},
		{"zwei Token je Klasse",
			map[string]string{"CDC_API_TOKENS_READER": "a,b", "CDC_API_TOKENS_ADMIN": "x,y"},
			[]string{"a", "b"}, []string{"x", "y"}, [2]string{"", ""}},
		{"ein Token in der Liste",
			map[string]string{"CDC_API_TOKENS_READER": "a"}, []string{"a"}, nil, [2]string{"", ""}},
		{"Singular neben der Liste",
			map[string]string{"CDC_API_TOKEN_READER": "alt", "CDC_API_TOKENS_READER": "neu", "CDC_API_TOKEN_ADMIN": "admin-alt"},
			[]string{"neu"}, nil, [2]string{"alt", "admin-alt"}},
		{"Singular mit Komma bleibt ein Token",
			map[string]string{"CDC_API_TOKEN_READER": "a,b"}, nil, nil, [2]string{"a,b", ""}},
	}
	for weg, laden := range apiTokensZugriffswege(t) {
		for _, tc := range cases {
			t.Run(weg+"/"+tc.name, func(t *testing.T) {
				env := vollständigeVerdrahtung()
				for k, v := range tc.env {
					env[k] = v
				}
				cfg, err := laden(env)
				if err != nil {
					t.Fatalf("Konfiguration lädt nicht: %v", err)
				}
				if !reflect.DeepEqual(cfg.APITokensReader, tc.wantReader) || !reflect.DeepEqual(cfg.APITokensAdmin, tc.wantAdmin) {
					t.Fatalf("Listen = %q / %q, erwartet %q / %q", cfg.APITokensReader, cfg.APITokensAdmin, tc.wantReader, tc.wantAdmin)
				}
				if cfg.APITokenReader != tc.wantSingular[0] || cfg.APITokenAdmin != tc.wantSingular[1] {
					t.Fatalf("Singular = %q / %q, erwartet %q / %q", cfg.APITokenReader, cfg.APITokenAdmin, tc.wantSingular[0], tc.wantSingular[1])
				}
			})
		}
	}
}

// TestAPITokenListeUngueltigEndetMitConfiguration trägt `LH-FA-SST-012`
// Negative: ein leeres Element, ein Element mit Leerraum und eine Liste aus nur
// Trennzeichen enden mit der Fehlerklasse `configuration` und dem Meldungscode
// `PCF-E2008`, die Zeile nennt Variable und Position des Elements und nie den
// Wert. Die Gegenprobe `a,b` ist gültig (siehe `TestAPITokenListenWerdenGelesen`),
// jeder Fall hängt damit an seiner Eingabe.
// Rot färbende Mutation: in `parseAPITokenList` die Prüfung mit
// `unicode.IsSpace` entfernen — die Fälle „Leerzeichen“ und „Tabulator“ laden
// ohne Fehler.
func TestAPITokenListeUngueltigEndetMitConfiguration(t *testing.T) {
	cases := []struct {
		name   string
		raw    string
		detail string
	}{
		{"leeres Element in der Mitte", "geheim1,,geheim2", "Element 2 ist leer"},
		{"leeres letztes Element", "geheim1,", "Element 2 ist leer"},
		{"leeres erstes Element", ",geheim1", "Element 1 ist leer"},
		{"nur Trennzeichen", ",", "Element 1 ist leer"},
		{"Leerzeichen nach dem Komma", "geheim1, geheim2", "Element 2 enthält Leerraum"},
		{"Leerzeichen am Anfang", " geheim1", "Element 1 enthält Leerraum"},
		{"Tabulator im Element", "geheim1\tgeheim2", "Element 1 enthält Leerraum"},
		{"Zeilenumbruch am Ende", "geheim1,geheim2\n", "Element 2 enthält Leerraum"},
	}
	for weg, laden := range apiTokensZugriffswege(t) {
		for _, variable := range []string{"CDC_API_TOKENS_READER", "CDC_API_TOKENS_ADMIN"} {
			for _, tc := range cases {
				t.Run(weg+"/"+variable+"/"+tc.name, func(t *testing.T) {
					env := vollständigeVerdrahtung()
					env[variable] = tc.raw
					_, err := laden(env)
					if err == nil {
						t.Fatalf("Liste %q lädt ohne Fehler", tc.raw)
					}
					if !errors.Is(err, bootstrap.ErrConfiguration) || !errors.Is(err, bootstrap.ErrAPITokenList) {
						t.Fatalf("Fehler ist nicht ErrConfiguration und ErrAPITokenList: %v", err)
					}
					if code, ok := messagecode.From(err); !ok || code != messagecode.APITokenListInvalid {
						t.Fatalf("Code = %q (%v), erwartet %q", code, ok, messagecode.APITokenListInvalid)
					}
					text := err.Error()
					for _, want := range []string{"PCF-E2008", "configuration", variable + ": " + tc.detail} {
						if !strings.Contains(text, want) {
							t.Fatalf("Fehlertext %q trägt %q nicht", text, want)
						}
					}
					if strings.Contains(text, "geheim") {
						t.Fatalf("Fehlertext trägt einen Token-Wert: %q", text)
					}
				})
			}
		}
	}
}
