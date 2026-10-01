// Command sseclient ist ein Wegwerf-Testclient für den SSE-Stream-E2E-Beleg
// (LH-FA-SST-008): er öffnet den Endpunkt `GET /changes/stream` real
// per HTTP gegen den laufenden Feed-Container, meldet die Bereitschaft über
// die Zeile "READY" auf stdout, empfängt danach einen committeten Change mit
// vollständigem Inhalt ("RECEIVED") und belegt abschließend, dass ein
// Öffnungsversuch ohne gültiges Token mit HTTP-Status 401 abgelehnt wird
// ("REJECTED"). Träger ist tools/harness/run-integration-tests.sh — der
// Aufrufer liest die stdout-Zeilen dieses Prozesses über `docker logs`, nicht
// über einen Exit-Code allein, weil "READY" vor der auslösenden Change
// beobachtbar sein muss. Das Flag `-target` setzt den Query-Parameter `target`
// (Zustellziel), `-count` die Zahl der Changes, die der Client vor dem
// Token-Test empfängt.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// sseChange trägt die für die Auswertung nötigen Felder eines Change-Events
// (SPEC-018): Kennung, Tabelle, Operation und das neue Row Image.
type sseChange struct {
	ChangeID  string          `json:"change_id"`
	Table     string          `json:"table"`
	Operation string          `json:"operation"`
	NewImage  json.RawMessage `json:"new_image"`
}

func main() {
	target := flag.String("target", "", "Zustellziel-Filter der Query (leer = kein Filter)")
	count := flag.Int("count", 1, "Zahl der Changes, die der Client empfängt, bevor er den Token-Test fährt")
	flag.Parse()
	args := flag.Args()
	if (len(args) != 2 && len(args) != 4) || *count < 1 {
		fmt.Fprintln(os.Stderr, "usage: sseclient [-target <ziel>] [-count <n>] <base-url> <token> [<schema> <table>]")
		os.Exit(2)
	}
	baseURL, token := args[0], args[1]
	// Zwei optionale, nachgestellte Argumente tragen das Filterpaar der
	// Query (`ADR-0133`); ohne sie bleibt die Query unverändert leer.
	var schema, table string
	if len(args) == 4 {
		schema, table = args[2], args[3]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	streamURL := baseURL + "/changes/stream"
	if schema != "" || table != "" || *target != "" {
		query := url.Values{}
		if schema != "" {
			query.Set("schema", schema)
		}
		if table != "" {
			query.Set("table", table)
		}
		if *target != "" {
			query.Set("target", *target)
		}
		streamURL += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, streamURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sseclient: Request bauen: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sseclient: Stream öffnen: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "sseclient: Stream-Öffnung endete mit Status %d (Erwartung 200)\n", resp.StatusCode)
		os.Exit(1)
	}

	// Die Bereitschafts-Zeile folgt dem Empfang der Antwort-Header; der
	// Handler abonniert den Broadcaster vor dem Header-Flush, ein danach
	// committeter Change hat damit einen registrierten Empfänger.
	fmt.Println("READY")

	// Ein Scanner über den ganzen Stream: gepufferte Bytes hinter dem ersten
	// Event gehen beim nächsten Lesen nicht verloren.
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for i := 0; i < *count; i++ {
		change, err := naechstesChange(scanner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sseclient: kein Change auf dem Stream innerhalb der Frist: %v\n", err)
			fmt.Println("TIMEOUT")
			os.Exit(1)
		}
		fmt.Printf("RECEIVED change_id=%s table=%s operation=%s new_image=%s\n",
			change.ChangeID, change.Table, change.Operation, string(change.NewImage))
	}

	// Der Stream bleibt bis zum Verbindungsende offen; für den Rest des
	// Prozesses wird die Verbindung geschlossen.
	_ = resp.Body.Close()

	if err := assertUnauthorized(baseURL); err != nil {
		fmt.Fprintf(os.Stderr, "sseclient: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("REJECTED code=401")
}

// naechstesChange liest SSE-Events, bis die data-Zeile eines Events einen
// Change trägt, und liefert ihn zurück.
func naechstesChange(scanner *bufio.Scanner) (sseChange, error) {
	var data string
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "data: "):
			data = strings.TrimPrefix(line, "data: ")
		case line == "":
			if data == "" {
				continue
			}
			var change sseChange
			if err := json.Unmarshal([]byte(data), &change); err != nil {
				return sseChange{}, err
			}
			return change, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return sseChange{}, err
	}
	return sseChange{}, fmt.Errorf("Stream endete ohne Change-Event")
}

// assertUnauthorized öffnet den Stream ohne Authorization-Header und
// erwartet die Ablehnung mit HTTP-Status 401 — sichtbar, nicht still mit
// leeren Daten fortgesetzt (LH-FA-SST-008 Negative).
func assertUnauthorized(baseURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/changes/stream", nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("Stream-Öffnungsversuch ohne Token endete mit Status %d, wollen 401", resp.StatusCode)
	}
	return nil
}
