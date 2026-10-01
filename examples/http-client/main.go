// Command http-client ist ein öffentliches Beispiel für den
// Anfrage/Antwort-Zugriff über die HTTP-/JSON-API (`LH-FA-SST-006`): das
// `-verb`-Flag ruft eine von zehn dokumentierten Fähigkeiten real gegen den
// laufenden Feed-Container auf und gibt die Antwort aus. Startform ist
// `go run ./examples/http-client -source <quelle> -publication <publication>`
// (Default-Verb `tables`); der Zugriffs-Abschnitt des Benutzerhandbuchs ist
// `### Zugriff über die HTTP-/JSON-API` (`docs/user/benutzerhandbuch.md`).
//
// Dieses Programm ist das Vorbild (minimal, lesbar); der E2E-Belegträger zu
// LH-FA-SST-006 ist der Wegwerf-Client `tools/harness/httpclient`. Das
// Beispiel trägt keine Zustandsmaschine: es stellt eine Anfrage und endet.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// requestTimeout begrenzt den einzelnen Anfrage/Antwort-Aufruf; die
// Verwaltungs-API antwortet synchron (`ADR-0057`).
const requestTimeout = 10 * time.Second

// config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse der
// HTTP-API, die zwei Token-Klassen (`CDC_API_TOKEN_READER`/
// `CDC_API_TOKEN_ADMIN`) und die Felder aller zehn Fähigkeiten — je Verb
// werden nur die tatsächlich nötigen Felder geprüft (`validate`).
type config struct {
	addr       string
	token      string
	adminToken string
	verb       string

	source      string
	publication string

	consumerID string
	name       string
	offset     uint64

	schema          string
	table           string
	tableID         string
	schemaVersionID string
	version         int64

	target string
	from   string
	to     string
	limit  string

	minAgeNanos int64
}

// knownVerbs trägt die geschlossene Menge der `-verb`-Werte — ein
// unbekannter Wert bricht ab, bevor ein Netzwerkaufruf versucht wird.
var knownVerbs = map[string]bool{
	"tables":            true,
	"changes":           true,
	"register-consumer": true,
	"acknowledge":       true,
	"consumer-position": true,
	"remove-consumer":   true,
	"enable-table":      true,
	"disable-table":     true,
	"table-status":      true,
	"retention-run":     true,
}

func main() {
	cfg := parseFlags()

	if err := validate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "http-client: %v\n", err)
		os.Exit(2)
	}

	client := &http.Client{Timeout: requestTimeout}
	body, err := dispatch(client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "http-client: -verb=%s fehlgeschlagen: %v\n", cfg.verb, err)
		os.Exit(1)
	}
	fmt.Println(body)
}

// parseFlags liest die Flag-Werte und füllt fehlende Felder aus den im
// Handbuch dokumentierten Umgebungsvariablen (`ADR-0076` Festlegung 1).
func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", os.Getenv("CDC_HTTP_ADDR"), "Horch-Adresse der HTTP-API, host:port (Default: CDC_HTTP_ADDR)")
	flag.StringVar(&cfg.token, "token", os.Getenv("CDC_API_TOKEN_READER"), "Bearer-Token der lesenden Rechtsklasse (Default: CDC_API_TOKEN_READER)")
	flag.StringVar(&cfg.adminToken, "admin-token", os.Getenv("CDC_API_TOKEN_ADMIN"), "Bearer-Token der administrativen Rechtsklasse (Default: CDC_API_TOKEN_ADMIN)")
	flag.StringVar(&cfg.verb, "verb", "tables", "Aufgerufene Fähigkeit: tables|changes|register-consumer|acknowledge|consumer-position|remove-consumer|enable-table|disable-table|table-status|retention-run")

	flag.StringVar(&cfg.source, "source", "", "Quelle (source_id)")
	flag.StringVar(&cfg.publication, "publication", "", "Publication")

	flag.StringVar(&cfg.consumerID, "consumer-id", "", "Consumer-Kennung")
	flag.StringVar(&cfg.name, "name", "", "Anzeigename des Consumers (nur register-consumer)")
	flag.Uint64Var(&cfg.offset, "offset", 0, "Bestätigte Position (nur acknowledge)")

	flag.StringVar(&cfg.schema, "schema", "", "Schema der Tabelle")
	flag.StringVar(&cfg.table, "table", "", "Name der Tabelle")
	flag.StringVar(&cfg.tableID, "table-id", "", "Tabellen-Kennung (Default: <schema>.<table>, nur enable-table)")
	flag.StringVar(&cfg.schemaVersionID, "schema-version-id", "", "Schema-Versions-Kennung (Default: <table-id>-v1, nur enable-table)")
	flag.Int64Var(&cfg.version, "version", 1, "Versionsnummer der Tabelle (nur enable-table)")

	flag.StringVar(&cfg.target, "target", "", "Zustellziel einer Change (nur changes, optional; leer = kein Filter)")
	flag.StringVar(&cfg.from, "from", "", "Untere Positions-Grenze, einschließlich (nur changes, optional)")
	flag.StringVar(&cfg.to, "to", "", "Obere Positions-Grenze, ausschließlich (nur changes, optional)")
	flag.StringVar(&cfg.limit, "limit", "", "Maximale Zeilenzahl (nur changes, optional)")

	flag.Int64Var(&cfg.minAgeNanos, "min-age-nanos", 0, "Mindestalter der Retention-Policy in Nanosekunden (nur retention-run)")
	flag.Parse()
	return cfg
}

// validate prüft cfg gegen die Pflichtfelder des gewählten Verbs — vor jedem
// Netzwerkaufruf. Jedes Verb braucht die Horch-Adresse und die zu seiner
// Rechtsklasse passende Token-Variable; die übrigen Felder folgen dem
// jeweiligen Endpunkt-Vertrag.
func validate(cfg config) error {
	if !knownVerbs[cfg.verb] {
		return fmt.Errorf("unbekanntes -verb %q", cfg.verb)
	}
	if cfg.addr == "" {
		return fmt.Errorf("keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder -addr) ist nötig, um die Verwaltungs-API zu erreichen")
	}

	switch cfg.verb {
	case "tables":
		if err := requireReaderToken(cfg); err != nil {
			return err
		}
		if cfg.source == "" || cfg.publication == "" {
			return fmt.Errorf("-source und -publication sind Pflicht — sie sind die zwei Pflichtfelder von GET /tables")
		}
	case "changes":
		if err := requireReaderToken(cfg); err != nil {
			return err
		}
		if cfg.source == "" {
			return fmt.Errorf("-source ist Pflicht — es ist das einzige Pflichtfeld von GET /changes")
		}
	case "register-consumer":
		if err := requireAdminToken(cfg); err != nil {
			return err
		}
		if cfg.consumerID == "" || cfg.name == "" {
			return fmt.Errorf("-consumer-id und -name sind Pflicht — sie sind die zwei Pflichtfelder von POST /consumers")
		}
	case "acknowledge":
		if err := requireAdminToken(cfg); err != nil {
			return err
		}
		if cfg.consumerID == "" || cfg.source == "" {
			return fmt.Errorf("-consumer-id und -source sind Pflicht — sie sind zwei der drei Pflichtfelder von POST /consumers/acknowledge")
		}
	case "consumer-position":
		if err := requireReaderToken(cfg); err != nil {
			return err
		}
		if cfg.consumerID == "" {
			return fmt.Errorf("-consumer-id ist Pflicht — es ist das einzige Pflichtfeld von GET /consumers/position")
		}
	case "remove-consumer":
		if err := requireAdminToken(cfg); err != nil {
			return err
		}
		if cfg.consumerID == "" {
			return fmt.Errorf("-consumer-id ist Pflicht — es ist das einzige Pflichtfeld von POST /consumers/remove")
		}
	case "enable-table":
		if err := requireAdminToken(cfg); err != nil {
			return err
		}
		if cfg.source == "" || cfg.schema == "" || cfg.table == "" || cfg.publication == "" {
			return fmt.Errorf("-source, -schema, -table und -publication sind Pflicht — table-id und schema-version-id haben einen Default")
		}
	case "disable-table":
		if err := requireAdminToken(cfg); err != nil {
			return err
		}
		if cfg.source == "" || cfg.schema == "" || cfg.table == "" || cfg.publication == "" {
			return fmt.Errorf("-source, -schema, -table und -publication sind Pflicht — sie sind die vier Pflichtfelder von POST /tables/disable")
		}
	case "table-status":
		if err := requireReaderToken(cfg); err != nil {
			return err
		}
		if cfg.source == "" || cfg.schema == "" || cfg.table == "" || cfg.publication == "" {
			return fmt.Errorf("-source, -schema, -table und -publication sind Pflicht — sie sind die vier Pflichtfelder von GET /tables/status")
		}
	case "retention-run":
		if err := requireAdminToken(cfg); err != nil {
			return err
		}
		if cfg.source == "" {
			return fmt.Errorf("-source ist Pflicht — es ist eines der zwei Pflichtfelder von POST /retention/run")
		}
	}
	return nil
}

func requireReaderToken(cfg config) error {
	if cfg.token == "" {
		return fmt.Errorf("kein Token gesetzt — CDC_API_TOKEN_READER (oder -token) ist nötig, um über die reader-Rechtsklasse zu lesen")
	}
	return nil
}

func requireAdminToken(cfg config) error {
	if cfg.adminToken == "" {
		return fmt.Errorf("kein Admin-Token gesetzt — CDC_API_TOKEN_ADMIN (oder -admin-token) ist nötig, um über die admin-Rechtsklasse zu schreiben")
	}
	return nil
}

// dispatch ruft die zum Verb gehörende Fähigkeit auf und liefert den
// Antwort-Body als druckbaren Text. `tables` liest die rohe Server-Antwort
// (`listTables`); die übrigen neun Fähigkeiten drucken ihre typisierte
// Antwort eingerückt.
func dispatch(client *http.Client, cfg config) (string, error) {
	switch cfg.verb {
	case "tables":
		return listTables(client, cfg)
	case "changes":
		resp, err := readChanges(client, cfg)
		return marshalIndent(resp, err)
	case "register-consumer":
		resp, err := registerConsumer(client, cfg)
		return marshalIndent(resp, err)
	case "acknowledge":
		resp, err := acknowledgeConsumer(client, cfg)
		return marshalIndent(resp, err)
	case "consumer-position":
		resp, err := consumerPosition(client, cfg)
		return marshalIndent(resp, err)
	case "remove-consumer":
		resp, err := removeConsumer(client, cfg)
		return marshalIndent(resp, err)
	case "enable-table":
		if cfg.tableID == "" {
			cfg.tableID = cfg.schema + "." + cfg.table
		}
		if cfg.schemaVersionID == "" {
			cfg.schemaVersionID = cfg.tableID + "-v1"
		}
		resp, err := enableTable(client, cfg)
		return marshalIndent(resp, err)
	case "disable-table":
		resp, err := disableTable(client, cfg)
		return marshalIndent(resp, err)
	case "table-status":
		resp, err := tableStatus(client, cfg)
		return marshalIndent(resp, err)
	case "retention-run":
		resp, err := runRetention(client, cfg)
		return marshalIndent(resp, err)
	default:
		return "", fmt.Errorf("unbekanntes -verb %q", cfg.verb)
	}
}

// marshalIndent gibt einen bereits gelaufenen Aufruf-Fehler unverändert
// weiter und kodiert sonst die typisierte Antwort eingerückt.
func marshalIndent(v any, err error) (string, error) {
	if err != nil {
		return "", err
	}
	encoded, encErr := json.MarshalIndent(v, "", "  ")
	if encErr != nil {
		return "", encErr
	}
	return string(encoded), nil
}

// listTables ruft die Tabellen-Auflistung mit dem `reader`-Token ab und
// liefert den Response-Body als Text — die einzige Fähigkeit ohne
// typisierte Antwort (`tables.go` trägt nur die URL-Bau-Funktion
// `TablesURL`). Der `*http.Client` kommt von `dispatch`, wie bei den neun
// übrigen Verben — kein eigener zweiter Client.
func listTables(client *http.Client, cfg config) (string, error) {
	req, err := http.NewRequest(http.MethodGet, TablesURL(cfg.addr, cfg.source, cfg.publication), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP-Status %d: %s", resp.StatusCode, body)
	}
	return string(body), nil
}
