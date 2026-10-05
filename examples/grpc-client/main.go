// Command grpc-client ist ein öffentliches Beispiel für die vollständige
// gRPC-Fläche (`LH-FA-SST-006`): das `-verb`-Flag ruft eine von zwölf
// dokumentierten Fähigkeiten real gegen den laufenden Feed-Container auf —
// der Live-Change-Stream (Default-Verb `stream`, unverändert die
// ursprüngliche Aufrufform) und die elf unären RPCs des
// `Administration`-Diensts. Startform ist `go run ./examples/grpc-client`;
// die Zugriffs-Abschnitte des Benutzerhandbuchs sind `### Zugriff über den
// gRPC-Change-Stream` und `### Zugriff über die gRPC-Verwaltungs-API`
// (`docs/user/benutzerhandbuch.md`).
//
// Dieses Programm ist das Vorbild (minimal, lesbar); der E2E-Belegträger ist
// der Wegwerf-Client `tools/harness/grpcadminclient`. Das Beispiel trägt
// keine Zustandsmaschine: die elf RPCs stellen je eine Anfrage und enden,
// der Stream bleibt offen und kennt kein Replay — verpasste Changes holt der
// bestehende Lesezugriffsweg (`ReadChanges`, `GET /changes`) nach.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
)

// authorizationMetadataKey und bearerPrefix tragen dieselbe Wertform wie der
// Auth-Interceptor des Adapters (`SPEC-020`): der `authorization`-Metadata-
// Wert trägt den Token hinter dem `Bearer `-Vorsprung.
const (
	authorizationMetadataKey = "authorization"
	bearerPrefix             = "Bearer "
)

// requestTimeout begrenzt einen einzelnen unären RPC-Aufruf; die
// Administration-API antwortet synchron (dieselbe Frist wie
// `examples/http-client`s `doRequestJSON`). Der Stream-Verb trägt keine
// eigene Frist — er bleibt offen, bis die Verbindung endet.
const requestTimeout = 10 * time.Second

// config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse des
// gRPC-Servers, die zwei Token-Klassen und die Felder aller zwölf
// Fähigkeiten — je Verb werden nur die tatsächlich nötigen Felder geprüft
// (`validate`). `schema`/`table` dienen doppelt: als optionaler Stream-Filter
// (ADR-0133) und als Tabellen-Identität der Verwaltungs-RPCs.
type config struct {
	addr       string
	token      string
	adminToken string
	verb       string
	caFile     string

	schema string
	table  string
	target string

	consumerID string
	name       string
	offset     uint64

	tableID         string
	schemaVersionID string
	version         int64
	source          string
	publication     string

	from  uint64
	to    uint64
	limit int64

	minAgeNanos int64
}

// knownVerbs trägt die geschlossene Menge der `-verb`-Werte — ein
// unbekannter Wert bricht ab, bevor ein Netzwerkaufruf versucht wird.
var knownVerbs = map[string]bool{
	"stream":                true,
	"register-consumer":     true,
	"acknowledge-consumer":  true,
	"get-consumer-position": true,
	"remove-consumer":       true,
	"enable-table":          true,
	"disable-table":         true,
	"get-table-status":      true,
	"list-tables":           true,
	"run-retention":         true,
	"read-changes":          true,
	"diagnose":              true,
}

func main() {
	cfg := parseFlags()
	if err := validate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "grpc-client: %v\n", err)
		os.Exit(2)
	}

	creds, err := transportCredentials(cfg.caFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpc-client: %v\n", err)
		os.Exit(2)
	}
	conn, err := grpc.NewClient(cfg.addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpc-client: Verbindung (%s) fehlgeschlagen: %v\n", cfg.addr, err)
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	if cfg.verb == "stream" {
		runStream(conn, cfg)
		return
	}

	client := administrationv1.NewAdministrationClient(conn)
	out, err := dispatchAdmin(client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpc-client: -verb=%s fehlgeschlagen: %v\n", cfg.verb, err)
		os.Exit(1)
	}
	fmt.Println(out)
}

// transportCredentials wählt den Transport: nennt caFile ein PEM-Zertifikat,
// vertraut der Client ihm als Anker und spricht TLS (Kette und Name werden
// geprüft); ohne Angabe spricht er Klartext. Eine nicht lesbare Datei oder
// eine Datei ohne PEM-Zertifikat ist ein Fehler vor dem Verbindungsaufbau.
func transportCredentials(caFile string) (credentials.TransportCredentials, error) {
	if caFile == "" {
		return insecure.NewCredentials(), nil
	}
	pem, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("Vertrauensanker %s nicht lesbar: %w", caFile, err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("Vertrauensanker %s enthält kein PEM-Zertifikat", caFile)
	}
	return credentials.NewTLS(&tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12}), nil
}

// parseFlags liest die Flag-Werte und füllt fehlende Felder aus den im
// Handbuch dokumentierten Umgebungsvariablen (`ADR-0076` Festlegung 1).
func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", os.Getenv("CDC_GRPC_ADDR"), "Horch-Adresse des gRPC-Servers, host:port (Default: CDC_GRPC_ADDR)")
	flag.StringVar(&cfg.token, "token", os.Getenv("CDC_API_TOKEN_READER"), "Bearer-Token der lesenden Rechtsklasse (Default: CDC_API_TOKEN_READER)")
	flag.StringVar(&cfg.adminToken, "admin-token", os.Getenv("CDC_API_TOKEN_ADMIN"), "Bearer-Token der administrativen Rechtsklasse (Default: CDC_API_TOKEN_ADMIN)")
	flag.StringVar(&cfg.caFile, "ca-file", os.Getenv("CDC_TLS_CA_FILE"), "Zertifikat (PEM) als Vertrauensanker; gesetzt verbindet das Programm über TLS, ungesetzt im Klartext (Default: CDC_TLS_CA_FILE)")
	flag.StringVar(&cfg.verb, "verb", "stream", "Aufgerufene Fähigkeit: stream|register-consumer|acknowledge-consumer|get-consumer-position|remove-consumer|enable-table|disable-table|get-table-status|list-tables|run-retention|read-changes|diagnose")

	flag.StringVar(&cfg.schema, "schema", "", "Schema (Stream-Filter, optional; sonst Pflichtfeld der Tabellen-RPCs)")
	flag.StringVar(&cfg.table, "table", "", "Tabellenname (Stream-Filter, optional; sonst Pflichtfeld der Tabellen-RPCs)")

	flag.StringVar(&cfg.target, "target", "", "Zustellziel einer Change (Stream-Filter und read-changes, optional; leer = kein Filter)")

	flag.StringVar(&cfg.consumerID, "consumer-id", "", "Consumer-Kennung")
	flag.StringVar(&cfg.name, "name", "", "Anzeigename des Consumers (nur register-consumer)")
	flag.Uint64Var(&cfg.offset, "offset", 0, "Bestätigte Position (nur acknowledge-consumer)")

	flag.StringVar(&cfg.tableID, "table-id", "", "Tabellen-Kennung (Default: <schema>.<table>, nur enable-table)")
	flag.StringVar(&cfg.schemaVersionID, "schema-version-id", "", "Schema-Versions-Kennung (Default: <table-id>-v1, nur enable-table)")
	flag.Int64Var(&cfg.version, "version", 1, "Versionsnummer der Tabelle (nur enable-table)")
	flag.StringVar(&cfg.source, "source", "", "Quelle (source_id)")
	flag.StringVar(&cfg.publication, "publication", "", "Publication")

	flag.Uint64Var(&cfg.from, "from", 0, "Untere Positions-Grenze, einschließlich (nur read-changes, optional)")
	flag.Uint64Var(&cfg.to, "to", 0, "Obere Positions-Grenze, ausschließlich (nur read-changes, optional)")
	flag.Int64Var(&cfg.limit, "limit", 0, "Maximale Zeilenzahl (nur read-changes, optional)")

	flag.Int64Var(&cfg.minAgeNanos, "min-age-nanos", 0, "Mindestalter der Retention-Policy in Nanosekunden (nur run-retention)")
	flag.Parse()
	return cfg
}

// validate prüft cfg gegen die Pflichtfelder des gewählten Verbs — vor jedem
// Netzwerkaufruf. Jedes Verb braucht die Horch-Adresse und die zu seiner
// Rechtsklasse passende Token-Variable (Rechtsklassen-Tabelle: `ADR-0130`).
func validate(cfg config) error {
	if !knownVerbs[cfg.verb] {
		return fmt.Errorf("unbekanntes -verb %q", cfg.verb)
	}
	if cfg.addr == "" {
		return fmt.Errorf("keine gRPC-Adresse gesetzt — CDC_GRPC_ADDR (oder -addr) ist nötig, um den Server zu erreichen")
	}

	switch cfg.verb {
	case "stream":
		return requireReaderToken(cfg)
	case "register-consumer":
		if err := requireAdminToken(cfg); err != nil {
			return err
		}
		if cfg.consumerID == "" || cfg.name == "" {
			return fmt.Errorf("-consumer-id und -name sind Pflicht — sie sind die zwei Pflichtfelder von RegisterConsumer")
		}
	case "acknowledge-consumer":
		if err := requireAdminToken(cfg); err != nil {
			return err
		}
		if cfg.consumerID == "" || cfg.source == "" {
			return fmt.Errorf("-consumer-id und -source sind Pflicht — sie sind zwei der drei Pflichtfelder von AcknowledgeConsumer")
		}
	case "get-consumer-position":
		if err := requireReaderToken(cfg); err != nil {
			return err
		}
		if cfg.consumerID == "" {
			return fmt.Errorf("-consumer-id ist Pflicht — es ist das einzige Pflichtfeld von GetConsumerPosition")
		}
	case "remove-consumer":
		if err := requireAdminToken(cfg); err != nil {
			return err
		}
		if cfg.consumerID == "" {
			return fmt.Errorf("-consumer-id ist Pflicht — es ist das einzige Pflichtfeld von RemoveConsumer")
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
			return fmt.Errorf("-source, -schema, -table und -publication sind Pflicht — sie sind die vier Pflichtfelder von DisableTable")
		}
	case "get-table-status":
		if err := requireReaderToken(cfg); err != nil {
			return err
		}
		if cfg.source == "" || cfg.schema == "" || cfg.table == "" || cfg.publication == "" {
			return fmt.Errorf("-source, -schema, -table und -publication sind Pflicht — sie sind die vier Pflichtfelder von GetTableStatus")
		}
	case "list-tables":
		if err := requireReaderToken(cfg); err != nil {
			return err
		}
		if cfg.source == "" || cfg.publication == "" {
			return fmt.Errorf("-source und -publication sind Pflicht — sie sind die zwei Pflichtfelder von ListTables")
		}
	case "run-retention":
		if err := requireAdminToken(cfg); err != nil {
			return err
		}
		if cfg.source == "" {
			return fmt.Errorf("-source ist Pflicht — es ist eines der zwei Pflichtfelder von RunRetention")
		}
	case "read-changes":
		if err := requireReaderToken(cfg); err != nil {
			return err
		}
		if cfg.source == "" {
			return fmt.Errorf("-source ist Pflicht — es ist das einzige Pflichtfeld von ReadChanges")
		}
	case "diagnose":
		if err := requireReaderToken(cfg); err != nil {
			return err
		}
		if cfg.source == "" {
			return fmt.Errorf("-source ist Pflicht — es ist das einzige Pflichtfeld von Diagnose")
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

// dispatchAdmin ruft die zum Verb gehörende Administration-RPC auf und
// liefert die formatierte Ausgabe. `stream` läuft nicht hier durch — er hat
// eine eigene, dauerhafte Schleife (`runStream`, `stream.go`).
func dispatchAdmin(client administrationv1.AdministrationClient, cfg config) (string, error) {
	switch cfg.verb {
	case "register-consumer":
		return registerConsumer(client, cfg)
	case "acknowledge-consumer":
		return acknowledgeConsumer(client, cfg)
	case "get-consumer-position":
		return getConsumerPosition(client, cfg)
	case "remove-consumer":
		return removeConsumer(client, cfg)
	case "enable-table":
		if cfg.tableID == "" {
			cfg.tableID = cfg.schema + "." + cfg.table
		}
		if cfg.schemaVersionID == "" {
			cfg.schemaVersionID = cfg.tableID + "-v1"
		}
		return enableTable(client, cfg)
	case "disable-table":
		return disableTable(client, cfg)
	case "get-table-status":
		return getTableStatus(client, cfg)
	case "list-tables":
		return listTables(client, cfg)
	case "run-retention":
		return runRetention(client, cfg)
	case "read-changes":
		return readChanges(client, cfg)
	case "diagnose":
		return diagnose(client, cfg)
	default:
		return "", fmt.Errorf("unbekanntes -verb %q", cfg.verb)
	}
}

// callCtx trägt die Aufruf-Frist und den `authorization`-Metadata-Eintrag in
// der `Bearer`-Wertform — der gemeinsame Ablauf aller elf Administration-RPCs.
func callCtx(token string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	return metadata.AppendToOutgoingContext(ctx, authorizationMetadataKey, bearerPrefix+token), cancel
}
