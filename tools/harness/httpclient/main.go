// Command httpclient ist ein Wegwerf-Testclient für die HTTP-API-E2E-Belege
// (LH-FA-SST-006): er ruft einen administrativen Endpunkt (RegisterConsumer)
// mit einem Admin-Token und drei lesende Endpunkte (ListTables, `GET
// /changes` und `GET /diagnose`) mit einem Reader-Token real per HTTP auf.
// Der Changes-Aufruf wertet die Antwort inhaltlich aus — Form, Filter-Treue,
// Operationswerte und Reihenfolge —, nicht nur ihren Status; jedes Ergebnis
// meldet er über eine benannte Zeile auf stdout. Träger ist
// tools/harness/run-integration-tests.sh — der Aufrufer liest die
// stdout-Zeilen dieses Prozesses.
//
// Zwei weitere Modi (erstes Argument `acknowledge` bzw. `remove`) tragen die
// administrative Consumer-Entfernung als zwei getrennte Aufrufe — der
// Aufrufer prüft den DB-Zustand real dazwischen (Retention-Blocker vor,
// Abwesenheit nach der Entfernung), was innerhalb eines einzigen
// Prozesslaufs nicht beobachtbar wäre.
//
// Zwei lesende Modi tragen den Backfill-Beleg: `changes` liest `GET
// /changes` mit dem reader-Token und gibt je Change eine READ-Zeile mit
// Kennung, Operation, Commit-Position, Row Image und `origin` aus;
// `position` liest `GET /consumers/position` und gibt den Antwort-Body als
// POSITION-Zeile aus.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// readChange trägt die für die inhaltliche Auswertung nötigen Felder eines
// gelesenen Change (`SPEC-022`): Kennung, Klartext-Identität der Tabelle,
// Operation, Commit-Position, das neue Row Image und die Herkunft.
type readChange struct {
	CommitPosition int64           `json:"commit_position"`
	ChangeID       string          `json:"change_id"`
	Schema         string          `json:"schema"`
	Table          string          `json:"table"`
	Operation      string          `json:"operation"`
	NewImage       json.RawMessage `json:"new_image"`
	Origin         string          `json:"origin"`
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "acknowledge" {
		runAcknowledgeFlow(os.Args[2:])
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "remove" {
		runRemoveFlow(os.Args[2:])
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "changes" {
		runChangesFlow(os.Args[2:])
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "position" {
		runPositionFlow(os.Args[2:])
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "fault" {
		runFaultFlow(os.Args[2:])
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "rejected" {
		runRejectedFlow(os.Args[2:])
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "probe" {
		runProbeFlow(os.Args[2:])
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "enable-table" {
		runEnableTableFlow(os.Args[2:])
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "disable-table" {
		runDisableTableFlow(os.Args[2:])
		return
	}
	if len(os.Args) != 13 {
		fmt.Fprintln(os.Stderr, "usage: httpclient <base-url> <admin-token> <reader-token> <consumer-id> <consumer-name> <source> <publication> <schema> <table> <from> <to> <limit>")
		fmt.Fprintln(os.Stderr, "   or: httpclient acknowledge <base-url> <admin-token> <consumer-id> <source-id> <offset>")
		fmt.Fprintln(os.Stderr, "   or: httpclient remove <base-url> <admin-token> <consumer-id>")
		fmt.Fprintln(os.Stderr, "   or: httpclient changes [-target <ziel>] <base-url> <reader-token> <source> <schema> <table> <from> <to>")
		fmt.Fprintln(os.Stderr, "   or: httpclient position <base-url> <reader-token> <consumer-id>")
		fmt.Fprintln(os.Stderr, "   or: httpclient fault <base-url> <reader-token> <source> <error-code|none>")
		fmt.Fprintln(os.Stderr, "   or: httpclient rejected <base-url> <reader-token>")
		fmt.Fprintln(os.Stderr, "   or: httpclient probe <base-url> <source> <publication> <token>...")
		fmt.Fprintln(os.Stderr, "   or: httpclient enable-table <base-url> <admin-token> <source> <schema> <table> <table-id> <schema-version-id> <publication>")
		fmt.Fprintln(os.Stderr, "   or: httpclient disable-table <base-url> <admin-token> <source> <schema> <table> <publication>")
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

	if err := readChanges(client, baseURL, readerToken, source, readSchema, readTable, "", readFrom, readTo, readLimit); err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: GET /changes (reader) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}

	diagnoseBody, err := call(client, http.MethodGet, baseURL+"/diagnose?source="+url.QueryEscape(source), readerToken, nil, http.StatusOK)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: GET /diagnose (reader) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("DIAGNOSED body=%s\n", diagnoseBody)
}

// readChanges ruft `GET /changes` mit der `reader`-Klasse auf und prüft die
// Antwort inhaltlich (`LH-FA-SST-006`): eine
// gesetzte, nicht leere Changes-Liste; je Eintrag eine nicht leere Kennung,
// ein bekannter Operationswert, eine Position ≥ 1 und die zum Filter
// passende Klartext-Identität; die Einträge in nicht absteigender
// Commit-Position (die deterministische Ordnung des Endpunkts)
// und eine Herkunft `wal` oder `backfill`. Ein leerer
// Parameter lässt den jeweiligen Query-Wert weg — `source` bleibt Pflicht.
func readChanges(client *http.Client, baseURL, token, source, schema, table, target, from, to, limit string) error {
	query := url.Values{}
	query.Set("source", source)
	if schema != "" {
		query.Set("schema", schema)
	}
	if table != "" {
		query.Set("table", table)
	}
	if target != "" {
		query.Set("target", target)
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
		switch change.Origin {
		case "wal", "backfill":
		default:
			return fmt.Errorf("Eintrag %d trägt die unbekannte Herkunft %q: %s", i, change.Origin, body)
		}
		lastPosition = change.CommitPosition
		fmt.Printf("READ changes=%d table=%s schema=%s change_id=%s operation=%s commit_position=%d new_image=%s origin=%s\n",
			len(result.Changes), change.Table, change.Schema, change.ChangeID, change.Operation, change.CommitPosition, string(change.NewImage), change.Origin)
	}
	return nil
}

// runChangesFlow trägt `GET /changes` als eigenen Modus ohne Registrierung
// (LH-FA-CAP-009): dieselbe inhaltliche Auswertung wie im
// Hauptmodus, ohne einen Consumer anzulegen.
func runChangesFlow(args []string) {
	flags := flag.NewFlagSet("changes", flag.ContinueOnError)
	target := flags.String("target", "", "Zustellziel-Filter (Query-Parameter target, leer = kein Filter)")
	if err := flags.Parse(args); err != nil {
		os.Exit(2)
	}
	args = flags.Args()
	if len(args) != 7 {
		fmt.Fprintln(os.Stderr, "usage: httpclient changes [-target <ziel>] <base-url> <reader-token> <source> <schema> <table> <from> <to>")
		os.Exit(2)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	if err := readChanges(client, args[0], args[1], args[2], args[3], args[4], *target, args[5], args[6], ""); err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: GET /changes (reader) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
}

// runFaultFlow liest `GET /diagnose` mit dem reader-Token, bis das Feld
// `error_code` der Antwort den erwarteten Wert trägt (`none` heißt: kein
// Fehlerzustand, `null`), und gibt den Körper als FAULT-Zeile aus. Der
// Aufrufer setzt den Fehlerzustand in `cdc.process_heartbeat`; der periodische
// Herzschlag des Feed-Containers löscht ihn wieder, deshalb liest der Client
// wiederholt, bis die Antwort den Zustand trägt. Ohne den erwarteten Wert
// binnen der Frist endet der Prozess mit Ausgang 1.
func runFaultFlow(args []string) {
	if len(args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: httpclient fault <base-url> <reader-token> <source> <error-code|none>")
		os.Exit(2)
	}
	baseURL, token, source, want := args[0], args[1], args[2], args[3]
	client := &http.Client{Timeout: 10 * time.Second}
	var last string
	for attempt := 0; attempt < 60; attempt++ {
		body, err := call(client, http.MethodGet, baseURL+"/diagnose?source="+url.QueryEscape(source), token, nil, http.StatusOK)
		if err != nil {
			fmt.Fprintf(os.Stderr, "httpclient: GET /diagnose (reader) fehlgeschlagen: %v\n", err)
			os.Exit(1)
		}
		last = strings.TrimSpace(body)
		var decoded struct {
			ErrorClass *string `json:"error_class"`
			ErrorCode  *string `json:"error_code"`
		}
		if err := json.Unmarshal([]byte(last), &decoded); err != nil {
			fmt.Fprintf(os.Stderr, "httpclient: GET /diagnose: Antwort nicht lesbar: %v\n", err)
			os.Exit(1)
		}
		if (want == "none" && decoded.ErrorCode == nil && decoded.ErrorClass == nil) ||
			(want != "none" && decoded.ErrorCode != nil && *decoded.ErrorCode == want) {
			fmt.Printf("FAULT want=%s body=%s\n", want, last)
			return
		}
		time.Sleep(250 * time.Millisecond)
	}
	fmt.Fprintf(os.Stderr, "httpclient: GET /diagnose trägt error_code %q nicht binnen der Frist: %s\n", want, last)
	os.Exit(1)
}

// runRejectedFlow ruft zwei Endpunkte so auf, dass der Server ablehnt, und
// gibt Status und Fehlerkörper je als REJECTED-Zeile aus: `GET /changes`
// ohne `source` (`400`, mit Code) und `GET /changes` ohne Token (`401`, ohne
// Code). Der Körper wird unverändert ausgegeben.
func runRejectedFlow(args []string) {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: httpclient rejected <base-url> <reader-token>")
		os.Exit(2)
	}
	baseURL, token := args[0], args[1]
	client := &http.Client{Timeout: 10 * time.Second}
	for _, probe := range []struct {
		name  string
		token string
	}{{"missing-source", token}, {"missing-token", ""}} {
		req, err := http.NewRequest(http.MethodGet, baseURL+"/changes", nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "httpclient: %v\n", err)
			os.Exit(1)
		}
		if probe.token != "" {
			req.Header.Set("Authorization", "Bearer "+probe.token)
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "httpclient: GET /changes (%s) fehlgeschlagen: %v\n", probe.name, err)
			os.Exit(1)
		}
		payload, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "httpclient: GET /changes (%s): Körper nicht lesbar: %v\n", probe.name, err)
			os.Exit(1)
		}
		fmt.Printf("REJECTED probe=%s status=%d body=%s\n", probe.name, resp.StatusCode, strings.TrimSpace(string(payload)))
	}
}

// runPositionFlow liest `GET /consumers/position` mit dem reader-Token und
// gibt den Antwort-Body aus (LH-FA-CON-005): die Anfangsposition eines
// Consumers ohne Bestätigung ist ein Messwert des Aufrufers, keine Annahme
// dieses Clients.
func runPositionFlow(args []string) {
	if len(args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: httpclient position <base-url> <reader-token> <consumer-id>")
		os.Exit(2)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	positionURL := fmt.Sprintf("%s/consumers/position?consumer_id=%s", args[0], url.QueryEscape(args[2]))
	body, err := call(client, http.MethodGet, positionURL, args[1], nil, http.StatusOK)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: GetConsumerPosition (reader) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("POSITION consumer=%s body=%s\n", args[2], strings.TrimSpace(body))
}

// runEnableTableFlow trägt EnableTable real per HTTP mit dem Admin-Token
// (LH-FA-CFG-001): aktiviert die Tabelle über den direkten HTTP-Zugriffsweg,
// ohne die SQL-Antragsqueue (`cdc.enable_table`) zu durchlaufen — der Beleg
// dafür, dass dieser Weg den laufenden Feed-Container ohne Neustart
// erfassend macht.
func runEnableTableFlow(args []string) {
	if len(args) != 8 {
		fmt.Fprintln(os.Stderr, "usage: httpclient enable-table <base-url> <admin-token> <source> <schema> <table> <table-id> <schema-version-id> <publication>")
		os.Exit(2)
	}
	baseURL, adminToken, source, schema, table, tableID, schemaVersionID, publication := args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7]

	client := &http.Client{Timeout: 10 * time.Second}
	body, err := call(client, http.MethodPost, baseURL+"/tables/enable", adminToken, map[string]any{
		"source": source, "schema": schema, "table": table,
		"table_id": tableID, "schema_version_id": schemaVersionID, "version": 1, "publication": publication,
	}, http.StatusCreated)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: EnableTable (admin) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("ENABLED body=%s\n", body)
}

// runProbeFlow ordnet jedes übergebene Token einer Rechtsklasse zu, wie der
// Server sie am Statuscode zeigt (LH-FA-SST-012): `GET /tables` (lesender
// Endpunkt) und `POST /consumers` mit leerem Körper (administrativer
// Endpunkt; die Rechteprüfung steht vor der Eingabeprüfung, es entsteht kein
// Consumer). Je Token eine PROBE-Zeile mit der Position des Tokens und beiden
// Statuscodes — der Wert des Tokens steht nie in der Ausgabe. Eine
// Transportstörung (der Server startet noch) wird bis zu 30 Sekunden
// wiederholt, ein Statuscode nie.
func runProbeFlow(args []string) {
	if len(args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: httpclient probe <base-url> <source> <publication> <token>...")
		os.Exit(2)
	}
	baseURL, source, publication, tokens := args[0], args[1], args[2], args[3:]
	client := &http.Client{Timeout: 10 * time.Second}
	listURL := fmt.Sprintf("%s/tables?source=%s&publication=%s", baseURL, url.QueryEscape(source), url.QueryEscape(publication))
	for i, token := range tokens {
		reader, err := probeStatus(client, http.MethodGet, listURL, token, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "httpclient: probe GET /tables (Token %d): %v\n", i+1, err)
			os.Exit(1)
		}
		admin, err := probeStatus(client, http.MethodPost, baseURL+"/consumers", token, "{}")
		if err != nil {
			fmt.Fprintf(os.Stderr, "httpclient: probe POST /consumers (Token %d): %v\n", i+1, err)
			os.Exit(1)
		}
		fmt.Printf("PROBE index=%d reader_endpoint=%d admin_endpoint=%d\n", i+1, reader, admin)
	}
}

// probeStatus sendet eine Anfrage mit Bearer-Token und liefert den
// Statuscode; ein Transportfehler wird bis zu 30 Sekunden wiederholt.
func probeStatus(client *http.Client, method, target, token, body string) (int, error) {
	deadline := time.Now().Add(30 * time.Second)
	for {
		req, err := http.NewRequest(method, target, strings.NewReader(body))
		if err != nil {
			return 0, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err := client.Do(req)
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			return resp.StatusCode, nil
		}
		if time.Now().After(deadline) {
			return 0, err
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// runDisableTableFlow trägt DisableTable real per HTTP mit dem Admin-Token
// (LH-FA-CFG-002): das Gegenstück zu `runEnableTableFlow`.
func runDisableTableFlow(args []string) {
	if len(args) != 6 {
		fmt.Fprintln(os.Stderr, "usage: httpclient disable-table <base-url> <admin-token> <source> <schema> <table> <publication>")
		os.Exit(2)
	}
	baseURL, adminToken, source, schema, table, publication := args[0], args[1], args[2], args[3], args[4], args[5]

	client := &http.Client{Timeout: 10 * time.Second}
	body, err := call(client, http.MethodPost, baseURL+"/tables/disable", adminToken, map[string]any{
		"source": source, "schema": schema, "table": table, "publication": publication,
	}, http.StatusOK)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: DisableTable (admin) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("DISABLED body=%s\n", body)
}

// runAcknowledgeFlow trägt AcknowledgeConsumer real per HTTP mit dem
// Admin-Token (LH-FA-CON-006-Vorbedingung): der übergebene, zuvor
// registrierte Consumer trägt die Position danach real als
// Retention-Blocker der Quelle fort — der Aufrufer prüft
// das über cdc.retention_blockers, bevor er den Consumer entfernt.
func runAcknowledgeFlow(args []string) {
	if len(args) != 5 {
		fmt.Fprintln(os.Stderr, "usage: httpclient acknowledge <base-url> <admin-token> <consumer-id> <source-id> <offset>")
		os.Exit(2)
	}
	baseURL := args[0]
	adminToken := args[1]
	consumerID := args[2]
	sourceID := args[3]
	offset, err := strconv.ParseUint(args[4], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: offset %q ist keine gültige Zahl: %v\n", args[4], err)
		os.Exit(2)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	if _, err := call(client, http.MethodPost, baseURL+"/consumers/acknowledge", adminToken,
		map[string]any{"consumer_id": consumerID, "source_id": sourceID, "offset": offset}, http.StatusOK); err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: AcknowledgeConsumer (admin) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("ACKNOWLEDGED consumer=%s source=%s offset=%d\n", consumerID, sourceID, offset)
}

// runRemoveFlow trägt RemoveConsumer real per HTTP mit dem Admin-Token
// (LH-FA-CON-006): entfernt den übergebenen, zuvor registrierten Consumer.
func runRemoveFlow(args []string) {
	if len(args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: httpclient remove <base-url> <admin-token> <consumer-id>")
		os.Exit(2)
	}
	baseURL := args[0]
	adminToken := args[1]
	consumerID := args[2]

	client := &http.Client{Timeout: 10 * time.Second}
	removeBody, err := call(client, http.MethodPost, baseURL+"/consumers/remove", adminToken,
		map[string]any{"consumer_id": consumerID}, http.StatusOK)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httpclient: RemoveConsumer (admin) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("REMOVED consumer=%s body=%s\n", consumerID, removeBody)
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
