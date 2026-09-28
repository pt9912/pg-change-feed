package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// readChangeResponse spiegelt einen Change in der Antwortform von
// `GET /changes` (`LH-FA-SST-006`, `internal/adapters/driving/http/readchanges.go`).
// Die Row Images stehen als eingebettete JSON-Werte; ein fehlendes Bild wird
// zu `null`.
type readChangeResponse struct {
	CommitPosition int64           `json:"commit_position"`
	ChangeID       string          `json:"change_id"`
	TransactionID  string          `json:"transaction_id"`
	SourceTableID  string          `json:"source_table_id"`
	Schema         string          `json:"schema"`
	Table          string          `json:"table"`
	Sequence       int64           `json:"sequence"`
	Operation      string          `json:"operation"`
	OldImage       json.RawMessage `json:"old_image"`
	NewImage       json.RawMessage `json:"new_image"`
	SchemaVersion  string          `json:"schema_version"`
	CommittedAt    string          `json:"committed_at"`
	Origin         string          `json:"origin"`
}

// readChangesResponse trägt den JSON-Response-Body bei Erfolg: ohne Treffer
// eine leere, gesetzte Liste.
type readChangesResponse struct {
	Changes []readChangeResponse `json:"changes"`
}

// ChangesURL baut die Lese-Adresse von `GET /changes`: `source` ist Pflicht,
// die übrigen fünf Parameter sind unabhängig optional — ein leerer Wert
// bleibt weg statt als leerer Query-Parameter zu erscheinen.
func ChangesURL(addr, source, schema, table, from, to, limit string) string {
	u := url.URL{Scheme: "http", Host: addr, Path: "/changes"}
	q := u.Query()
	q.Set("source", source)
	setIfNotEmpty(q, "schema", schema)
	setIfNotEmpty(q, "table", table)
	setIfNotEmpty(q, "from", from)
	setIfNotEmpty(q, "to", to)
	setIfNotEmpty(q, "limit", limit)
	u.RawQuery = q.Encode()
	return u.String()
}

// setIfNotEmpty setzt einen Query-Parameter nur, wenn ein Wert vorliegt.
// `GET /changes` lässt nur einen unbekannten Parameter**namen** mit `400`
// enden; ein leerer Wert eines bekannten optionalen Parameters (`schema`,
// `table`, `from`, `to`, `limit`) bleibt für den Server unberücksichtigt und
// endet mit `200` — dieser Client sendet ihn trotzdem gar nicht erst.
func setIfNotEmpty(q url.Values, key, value string) {
	if value != "" {
		q.Set(key, value)
	}
}

// readChanges ruft den `reader`-Endpunkt `GET /changes` auf.
func readChanges(client *http.Client, cfg config) (readChangesResponse, error) {
	var resp readChangesResponse
	err := doRequestJSON(client, http.MethodGet, ChangesURL(cfg.addr, cfg.source, cfg.schema, cfg.table, cfg.from, cfg.to, cfg.limit), cfg.token, nil, &resp, http.StatusOK)
	return resp, err
}
