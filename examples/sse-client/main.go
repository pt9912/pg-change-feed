// Command sse-client ist ein öffentliches Beispiel für den Live-Change-Stream
// über Server-Sent-Events (`ADR-0061`): es öffnet
// den Endpunkt `GET /changes/stream` real gegen den laufenden Feed-Container
// und gibt jedes Event aus. Startform ist `go run ./examples/sse-client`; der
// Zugriffs-Abschnitt des Benutzerhandbuchs ist `### Zugriff über Server-Sent-Events`
// (`docs/user/benutzerhandbuch.md`).
//
// Dieses Programm ist das Vorbild (minimal, lesbar); der E2E-Belegträger ist
// der Wegwerf-Client `tools/harness/sseclient`. Das Beispiel
// trägt keine Zustandsmaschine: der Stream kennt kein Replay, verpasste
// Changes holt der bestehende Lesezugriffsweg nach.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
)

// maxEventBytes trägt die Obergrenze einer gelesenen Zeile. Sie liegt über der
// Default-Grenze von `bufio.Scanner` (`bufio.MaxScanTokenSize`, 64 KiB), weil
// eine `data:`-Zeile den vollständigen Change-Inhalt samt Row Images trägt
// (`ADR-0061`).
const maxEventBytes = 1024 * 1024

// config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse der
// HTTP-API und das Token der lesenden Rechtsklasse. Beide kommen aus denselben
// Umgebungsvariablen, die das Benutzerhandbuch führt (`CDC_HTTP_ADDR`,
// `CDC_API_TOKEN_READER`), und lassen sich per Flag übersteuern (`ADR-0076`
// Festlegung 1).
type config struct {
	addr   string
	token  string
	target string
	schema string
	table  string
}

func main() {
	cfg, err := parseConfig(os.Args[1:], os.Getenv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sse-client: %v\n", err)
		os.Exit(2)
	}
	if cfg.addr == "" {
		fmt.Fprintln(os.Stderr, "sse-client: keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder -addr) ist nötig, um den Stream zu öffnen")
		os.Exit(2)
	}
	if cfg.token == "" {
		fmt.Fprintln(os.Stderr, "sse-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder -token) ist nötig, um über die reader-Rechtsklasse zu lesen")
		os.Exit(2)
	}

	streamURL := cfg.streamURL()
	req, err := http.NewRequest(http.MethodGet, streamURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sse-client: Request bauen: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.token)

	// Der Stream bleibt offen, bis die Verbindung endet: der Aufruf läuft
	// ohne Antwort-Frist (`http.Client.Timeout`).
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sse-client: Stream öffnen (%s): %v\n", streamURL, err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "sse-client: Stream-Öffnung endete mit HTTP-Status %d (Erwartung 200)\n", resp.StatusCode)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), maxEventBytes)
	next := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	for {
		ev, err := readEvent(next)
		if err != nil {
			// Das Ende der Quelle trägt zwei Ausgänge — ein Lesefehler
			// (`scanner.Err`) oder eine geschlossene Verbindung. Beide werden
			// gemeldet; der Prozess endet mit Nicht-Null.
			if serr := scanner.Err(); serr != nil {
				fmt.Fprintf(os.Stderr, "sse-client: Stream endete mit Fehler: %v\n", serr)
			} else {
				fmt.Fprintln(os.Stderr, "sse-client: Stream wurde geschlossen")
			}
			os.Exit(1)
		}
		fmt.Printf("sse-client: %s %s\n", ev.Name, ev.Data)
	}
}

// parseConfig liest die Flag-Werte aus args und füllt fehlende Felder aus den
// im Handbuch dokumentierten Umgebungsvariablen (`ADR-0076` Festlegung 1);
// das Umgebungs-Lookup ist injiziert, damit die Verdrahtung Flag → Konfiguration
// netzlos testbar ist.
func parseConfig(args []string, getenv func(string) string) (config, error) {
	var cfg config
	fs := flag.NewFlagSet("sse-client", flag.ContinueOnError)
	fs.StringVar(&cfg.addr, "addr", getenv("CDC_HTTP_ADDR"), "Horch-Adresse der HTTP-API, host:port (Default: CDC_HTTP_ADDR)")
	fs.StringVar(&cfg.token, "token", getenv("CDC_API_TOKEN_READER"), "Bearer-Token der lesenden Rechtsklasse (Default: CDC_API_TOKEN_READER)")
	fs.StringVar(&cfg.target, "target", "", "Zustellziel einer Change (optional; leer = kein Filter, jede Change)")
	fs.StringVar(&cfg.schema, "schema", "", "Schema einer Change (optional; leer = kein Filter)")
	fs.StringVar(&cfg.table, "table", "", "Tabellenname einer Change (optional; leer = kein Filter)")
	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	return cfg, nil
}

// streamURL baut die Stream-Adresse aus der Konfiguration; die drei Filter
// gehen unverändert an StreamURL.
func (c config) streamURL() string {
	return StreamURL(c.addr, c.target, c.schema, c.table)
}
