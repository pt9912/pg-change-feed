package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// doRequestJSON trägt den gemeinsamen Anfrage/Antwort-Ablauf aller zehn
// Fähigkeiten dieses Beispiels: JSON-Body kodieren (falls `reqBody` gesetzt
// ist), das Bearer-Token setzen, senden, und bei `wantStatus` den Antwort-Body
// nach `out` dekodieren. Ein davon abweichender Statuscode ist ein sichtbarer
// Fehler mit Statuscode und Antworttext, kein stiller Leerwert — dieselbe
// Form wie `listTables` in `tables.go`.
func doRequestJSON(client *http.Client, method, url, token string, reqBody, out any, wantStatus int) error {
	var bodyReader io.Reader
	if reqBody != nil {
		encoded, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(encoded)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return err
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != wantStatus {
		return fmt.Errorf("HTTP-Status %d: %s", resp.StatusCode, respBody)
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(respBody, out)
}
