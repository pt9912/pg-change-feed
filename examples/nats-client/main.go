// Command nats-client ist ein öffentliches Beispiel für den zweiseitigen
// Zugriffsweg über das NATS-Wecksignal (LH-FA-SST-007, ADR-0079): es abonniert
// das tabellen-granulare Subjekt `cdc.changes.<source_id>.<schema>.<table>`
// (SPEC-017, ADR-0056) und holt beim Weckruf die Änderung selbst über die
// HTTP-/JSON-API (LH-FA-SST-006) — das Signal trägt per Vertrag einen leeren
// Payload (ADR-0055), es sagt nur „lies erneut über den bestehenden
// Zugriffsweg". Startform ist `go run ./examples/nats-client`; der
// Zugriffs-Abschnitt des Benutzerhandbuchs ist
// `### Zugriff über das NATS-Wecksignal` (`docs/user/benutzerhandbuch.md`).
//
// Dieses Programm ist das Vorbild (minimal, lesbar); der E2E-Belegträger zu
// LH-FA-SST-007 ist der Wegwerf-Client `tools/harness/natssub`. Es ist ein
// Muster, kein Dauerbetriebs-Client: eine Zustandsmaschine (Reconnect,
// Deduplizierung, Rückstand) ist bewusst nicht seine Aufgabe (ADR-0055
// §Konsequenzen).
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// readTimeout begrenzt den einzelnen Lese-Aufruf über die HTTP-API; das
// Warten auf ein Wecksignal ist davon nicht betroffen.
const readTimeout = 10 * time.Second

// config trägt die Laufzeit-Eingabe des Beispiels: die beiden Draht-Adressen,
// das Token der lesenden Rechtsklasse und die drei Bestandteile des
// Subjekts. Adresse und Token kommen aus denselben Umgebungsvariablen, die
// das Benutzerhandbuch führt (`CDC_NATS_URL`, `CDC_HTTP_ADDR`,
// `CDC_API_TOKEN_READER`), und lassen sich per Flag übersteuern (`ADR-0079`
// Festlegung 2). Quelle, Schema und Tabelle sind Flags — das Handbuch führt
// für eine einzelne Tabelle keine `CDC_*`-Variable.
type config struct {
	natsURL string
	addr    string
	token   string
	source  string
	schema  string
	table   string
}

func main() {
	cfg := parseFlags()

	// `CDC_HTTP_ADDR` ungesetzt heißt: die HTTP-API ist deaktiviert
	// (`SPEC-018`). Ohne sie kann der Weckruf keine Änderung holen — das
	// Beispiel scheitert sichtbar, statt still nichts zu tun (`ADR-0079`
	// Festlegung 2).
	if cfg.addr == "" {
		fmt.Fprintln(os.Stderr, "nats-client: keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder -http-addr) ist nötig, um beim Weckruf die Änderung zu holen")
		os.Exit(2)
	}
	if cfg.natsURL == "" {
		fmt.Fprintln(os.Stderr, "nats-client: keine NATS-URL gesetzt — CDC_NATS_URL (oder -nats-url) ist nötig, um das Wecksignal zu abonnieren")
		os.Exit(2)
	}
	if cfg.source == "" || cfg.schema == "" || cfg.table == "" {
		fmt.Fprintln(os.Stderr, "nats-client: -source, -schema und -table sind Pflicht")
		os.Exit(2)
	}
	if cfg.token == "" {
		fmt.Fprintln(os.Stderr, "nats-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder -token) ist nötig, um über die reader-Rechtsklasse zu lesen")
		os.Exit(2)
	}

	subject := Subject(cfg.source, cfg.schema, cfg.table)

	conn, err := nats.Connect(cfg.natsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nats-client: Verbindung zu NATS (%s) fehlgeschlagen: %v\n", cfg.natsURL, err)
		os.Exit(1)
	}
	defer conn.Close()

	sub, err := conn.SubscribeSync(subject)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nats-client: Abonnement auf %s fehlgeschlagen: %v\n", subject, err)
		os.Exit(1)
	}
	defer func() { _ = sub.Unsubscribe() }()
	if err := conn.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "nats-client: Flush nach Abonnement fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("nats-client: lauscht auf %s — beim nächsten Weckruf wird die Änderung über %s geholt\n",
		subject, cfg.addr)

	msg, err := sub.NextMsg(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nats-client: Empfang des Wecksignals fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}

	// Der Payload trägt keine Änderungsdaten (`SPEC-017`); die Änderung kommt
	// ausschließlich über den HTTP-Lesezugriff.
	fmt.Printf("nats-client: Weckruf auf %s (Payload %d Byte) — hole die Änderung über HTTP\n",
		msg.Subject, len(msg.Data))

	body, err := fetchChanges(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nats-client: HTTP-Abfrage fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(body)
}

// parseFlags liest die Flag-Werte und füllt fehlende Felder aus den im
// Handbuch dokumentierten Umgebungsvariablen (`ADR-0079` Festlegung 2).
func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.natsURL, "nats-url", os.Getenv("CDC_NATS_URL"), "NATS-Server-URL des Wecksignals (Default: CDC_NATS_URL)")
	flag.StringVar(&cfg.addr, "http-addr", os.Getenv("CDC_HTTP_ADDR"), "Horch-Adresse der HTTP-API, host:port (Default: CDC_HTTP_ADDR)")
	flag.StringVar(&cfg.token, "token", os.Getenv("CDC_API_TOKEN_READER"), "Bearer-Token der lesenden Rechtsklasse (Default: CDC_API_TOKEN_READER)")
	flag.StringVar(&cfg.source, "source", "", "Quelle (source_id) des Wecksignal-Subjekts")
	flag.StringVar(&cfg.schema, "schema", "", "Schema der Tabelle")
	flag.StringVar(&cfg.table, "table", "", "Tabelle")
	flag.Parse()
	return cfg
}

// fetchChanges holt die Änderungen eines Quelle/Schema/Tabelle-Filters über
// die HTTP-/JSON-API und liefert den Response-Body als Text. Ein Nicht-200-
// Status ist ein sichtbarer Fehler, kein stiller Leerwert.
func fetchChanges(cfg config) (string, error) {
	req, err := http.NewRequest(http.MethodGet, ChangesURL(cfg.addr, cfg.source, cfg.schema, cfg.table), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.token)

	client := &http.Client{Timeout: readTimeout}
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
