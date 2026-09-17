// Command http-client ist ein öffentliches Beispiel für den
// Anfrage/Antwort-Zugriff über die HTTP-/JSON-API (`LH-FA-SST-006`,
// `ADR-0057`, `ADR-0076`): es ruft den `reader`-Endpunkt `GET /tables` real
// gegen den laufenden Feed-Container auf und gibt die Antwort aus. Startform
// ist `go run ./examples/http-client -source <quelle> -publication
// <publication>`; der Zugriffs-Abschnitt des Benutzerhandbuchs ist
// `### Zugriff über die HTTP-/JSON-API` (`docs/user/benutzerhandbuch.md`).
//
// Dieses Programm ist das Vorbild (minimal, lesbar); der E2E-Belegträger zu
// LH-FA-SST-006 ist der Wegwerf-Client `tools/harness/httpclient`. Das Beispiel
// trägt keine Zustandsmaschine: es stellt eine Anfrage und endet.
package main

import (
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
// HTTP-API, das Token der lesenden Rechtsklasse und die zwei Pflichtfelder
// des aufgerufenen Endpunkts. Adresse und Token kommen aus denselben
// Umgebungsvariablen, die das Benutzerhandbuch führt (`CDC_HTTP_ADDR`,
// `CDC_API_TOKEN_READER`), und lassen sich per Flag übersteuern (`ADR-0076`
// Festlegung 1). Für Quelle und Publication führt das Handbuch keine
// `CDC_*`-Variable — sie sind Flags.
type config struct {
	addr        string
	token       string
	source      string
	publication string
}

func main() {
	cfg := parseFlags()

	// `CDC_HTTP_ADDR` ungesetzt heißt: die HTTP-API ist deaktiviert — die
	// Zeile dazu steht in §5 *Konfiguration* des Handbuchs. Es gibt dann
	// keinen Endpunkt, und das Beispiel endet mit dieser Meldung.
	if cfg.addr == "" {
		fmt.Fprintln(os.Stderr, "http-client: keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder -addr) ist nötig, um die Verwaltungs-API zu erreichen")
		os.Exit(2)
	}
	if cfg.token == "" {
		fmt.Fprintln(os.Stderr, "http-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder -token) ist nötig, um über die reader-Rechtsklasse zu lesen")
		os.Exit(2)
	}
	if cfg.source == "" || cfg.publication == "" {
		fmt.Fprintln(os.Stderr, "http-client: -source und -publication sind Pflicht — sie sind die zwei Pflichtfelder von GET /tables")
		os.Exit(2)
	}

	body, err := listTables(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "http-client: GET /tables fehlgeschlagen: %v\n", err)
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
	flag.StringVar(&cfg.source, "source", "", "Quelle (source_id), deren Tabellen aufgelistet werden")
	flag.StringVar(&cfg.publication, "publication", "", "Publication, gegen die der Tabellen-Status gelesen wird")
	flag.Parse()
	return cfg
}

// listTables ruft die Tabellen-Auflistung mit dem `reader`-Token ab und
// liefert den Response-Body als Text. Ein Nicht-200-Status ist ein sichtbarer
// Fehler, kein stiller Leerwert.
func listTables(cfg config) (string, error) {
	req, err := http.NewRequest(http.MethodGet, TablesURL(cfg.addr, cfg.source, cfg.publication), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.token)

	client := &http.Client{Timeout: requestTimeout}
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
