// Command httpclient ist ein Wegwerf-Testclient für die HTTP-API-E2E-Belege
// (LH-FA-SST-006, ADR-0057, ADR-0081): er ruft einen administrativen Endpunkt
// (RegisterConsumer) mit einem Admin-Token und zwei lesende Endpunkte
// (ListTables und `GET /changes`) mit einem Reader-Token real per HTTP auf.
// Der Changes-Aufruf wertet die Antwort inhaltlich aus — Form, Filter-Treue,
// Operationswerte und Reihenfolge —, nicht nur ihren Status; jedes Ergebnis
// meldet er über eine benannte Zeile auf stdout. Träger ist
// tools/harness/run-integration-tests.sh — der Aufrufer liest die
// stdout-Zeilen dieses Prozesses.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// readChange trägt die für die inhaltliche Auswertung nötigen Felder eines
// gelesenen Change (`SPEC-022`): Kennung, Klartext-Identität der Tabelle,
// Operation, Commit-Position und das neue Row Image.
type readChange struct {
	CommitPosition int64           `json:"commit_position"`
	ChangeID       string          `json:"change_id"`
	Schema         string          `json:"schema"`
	Table          string          `json:"table"`
	Operation      string          `json:"operation"`
	NewImage       json.RawMessage `json:"new_image"`
}

func main() {
	if len(os.Args) != 13 {
		fmt.Fprintln(os.Stderr, "usage: httpclient <base-url> <admin-token> <reader-token> <consumer-id> <consumer-name> <source> <publication> <schema> <table> <from> <to> <limit>")
		os.Exit(2)
	}
	baseURL := os.Args[1]
	adminToken := os.Args[2]
	readerToken := os.Args[3]
	consumerID := os.Args[4]
	consumerName := os.Args[5]
	source := os.Args[6]
	publication := os.Args[7]
	readSchema := os.Args[8]
	readTable := os.Args[9]
	readFrom := os.Args[10]
	readTo := os.Args[11]
	readLimit := os.Args[12]

	client := &http.Client{Timeout: 10 * time.Second}

	registerBody, err := call(client, http.MethodPost, baseURL+"/consumers", adminToken,
		map[string]string{"consumer_id": consumerID, "name": consumerName}, http.StatusCreated)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: RegisterConsumer (admin) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("REGISTERED body=%s\n", registerBody)

	listURL := fmt.Sprintf("%s/tables?source=%s&publication=%s", baseURL, source, publication)
	listBody, err := call(client, http.MethodGet, listURL, readerToken, nil, http.StatusOK)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: ListTables (reader) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("LISTED body=%s\n", listBody)

	if err := readChanges(client, baseURL, readerToken, source, readSchema, readTable, readFrom, readTo, readLimit); err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: GET /changes (reader) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
}

// readChanges ruft `GET /changes` mit der `reader`-Klasse auf und prüft die
// Antwort inhaltlich (`LH-FA-SST-006`, `LH-FA-REA-001` ff., `ADR-0081`): eine
// gesetzte, nicht leere Changes-Liste; je Eintrag eine nicht leere Kennung,
// ein bekannter Operationswert, eine Position ≥ 1 und die zum Filter
// passende Klartext-Identität; die Einträge in nicht absteigender
// Commit-Position (die deterministische Ordnung des Endpunkts,
// `LH-FA-REA-004`). Ein leerer Parameter lässt den jeweiligen Query-Wert
// weg — `source` bleibt Pflicht.
func readChanges(client *http.Client, baseURL, token, source, schema, table, from, to, limit string) error {
	query := url.Values{}
	query.Set("source", source)
	if schema != "" {
		query.Set("schema", schema)
	}
	if table != "" {
		query.Set("table", table)
	}
	if from != "" {
		query.Set("from", from)
	}
	if to != "" {
		query.Set("to", to)
	}
	if limit != "" {
		query.Set("limit", limit)
	}
	body, err := call(client, http.MethodGet, baseURL+"/changes?"+query.Encode(), token, nil, http.StatusOK)
	if err != nil {
		return err
	}

	var result struct {
		Changes []readChange `json:"changes"`
	}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return fmt.Errorf("Antwort ist kein JSON: %v (%s)", err, body)
	}
	if len(result.Changes) == 0 {
		return fmt.Errorf("leere Changes-Liste für den erwarteten Bestand (%s.%s, Bereich [%s,%s)): %s",
			schema, table, from, to, body)
	}

	var lastPosition int64
	for i, change := range result.Changes {
		if change.ChangeID == "" {
			return fmt.Errorf("Eintrag %d trägt keine change_id: %s", i, body)
		}
		if schema != "" && change.Schema != schema {
			return fmt.Errorf("Eintrag %d trägt schema=%q, der Filter verlangt %q: %s", i, change.Schema, schema, body)
		}
		if table != "" && change.Table != table {
			return fmt.Errorf("Eintrag %d trägt table=%q, der Filter verlangt %q: %s", i, change.Table, table, body)
		}
		switch change.Operation {
		case "INSERT", "UPDATE", "DELETE":
		default:
			return fmt.Errorf("Eintrag %d trägt die unbekannte Operation %q: %s", i, change.Operation, body)
		}
		if change.CommitPosition < 1 {
			return fmt.Errorf("Eintrag %d trägt die Commit-Position %d (Position 0 existiert nicht): %s", i, change.CommitPosition, body)
		}
		if i > 0 && change.CommitPosition < lastPosition {
			return fmt.Errorf("Eintrag %d trägt die Commit-Position %d nach %d — die Ordnung ist absteigend: %s", i, change.CommitPosition, lastPosition, body)
		}
		lastPosition = change.CommitPosition
		fmt.Printf("READ changes=%d table=%s schema=%s change_id=%s operation=%s commit_position=%d new_image=%s\n",
			len(result.Changes), change.Table, change.Schema, change.ChangeID, change.Operation, change.CommitPosition, string(change.NewImage))
	}
	return nil
}

// call führt einen HTTP-Request mit Bearer-Token aus und liefert den
// Response-Body als String, wenn der Statuscode dem erwarteten Wert
// entspricht; sonst einen Fehler mit Statuscode und Body.
func call(client *http.Client, method, url, token string, jsonBody any, wantStatus int) (string, error) {
	var reqBody io.Reader
	if jsonBody != nil {
		encoded, err := json.Marshal(jsonBody)
		if err != nil {
			return "", err
		}
		reqBody = bytes.NewReader(encoded)
	}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if jsonBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != wantStatus {
		return "", fmt.Errorf("status %d (erwartet %d): %s", resp.StatusCode, wantStatus, string(payload))
	}
	return string(payload), nil
}
