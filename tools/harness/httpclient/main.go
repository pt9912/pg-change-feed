// Command httpclient ist ein Wegwerf-Testclient für den HTTP-API-E2E-Beleg
// (LH-FA-SST-006, ADR-0057): er ruft einen administrativen Endpunkt
// (RegisterConsumer) mit einem Admin-Token und einen lesenden Endpunkt
// (ListTables) mit einem Reader-Token real per HTTP auf und meldet jedes
// Ergebnis über eine benannte Zeile auf stdout. Träger ist
// tools/harness/run-integration-tests.sh — der Aufrufer liest die
// stdout-Zeilen dieses Prozesses.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 8 {
		fmt.Fprintln(os.Stderr, "usage: httpclient <base-url> <admin-token> <reader-token> <consumer-id> <consumer-name> <source> <publication>")
		os.Exit(2)
	}
	baseURL := os.Args[1]
	adminToken := os.Args[2]
	readerToken := os.Args[3]
	consumerID := os.Args[4]
	consumerName := os.Args[5]
	source := os.Args[6]
	publication := os.Args[7]

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
