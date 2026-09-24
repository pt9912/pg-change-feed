package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Query-Parameter des lesenden Endpunkts `GET /changes` (`SPEC-022`):
// `source` Pflicht, `schema`/`table`/`from`/`to`/`limit` optional und
// unabhängig (`ADR-0081` Teilfrage 2). Die Namen folgen dem API-Vokabular,
// nicht dem Spaltenvokabular der View (`schema`/`table`, nicht
// `schema_name`/`table_name`) — dieselbe Form wie `listTablesResponse`.
const (
	readChangesParamSource = "source"
	readChangesParamSchema = "schema"
	readChangesParamTable  = "table"
	readChangesParamFrom   = "from"
	readChangesParamTo     = "to"
	readChangesParamLimit  = "limit"
)

// readChangesParams trägt die geschlossene Parameter-Menge dieses
// Endpunkts. Anders als die neun übrigen Endpunkte endet ein Parameter
// außerhalb der Menge mit `400`: bei einem Filter-Endpunkt ist das
// Ignorieren keine Auslassung, sondern eine **stille Änderung des
// Ergebnisstands** (`ADR-0081` Teilfrage 4, `LH-FA-SST-006` Negative).
var readChangesParams = map[string]bool{
	readChangesParamSource: true,
	readChangesParamSchema: true,
	readChangesParamTable:  true,
	readChangesParamFrom:   true,
	readChangesParamTo:     true,
	readChangesParamLimit:  true,
}

// readChangeResponse trägt einen Change in der Antwortform von
// `GET /changes` (`SPEC-022`): dieselben Felder wie der Domain-Typ
// `model.Change` (`internal/domain/model/change.go`), ergänzt um die
// Commit-Position seiner Quelltransaktion und deren Commit-Zeitpunkt. Die
// Row Images stehen als eingebettete JSON-Werte; ein fehlendes Bild
// (`LH-FA-CAP-008` Boundary) wird zu `null`. `committed_at` trägt RFC 3339
// mit Nanosekunden in UTC. `origin` steht als letztes Feld und trägt `wal`
// oder `backfill` (`SPEC-002`, `LH-FA-CAP-009`); ein fehlender Wert liest
// als `wal`. Die Live-Wege tragen das Feld nicht.
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

// readChangesResponse trägt den JSON-Response-Body bei Erfolg
// (`SPEC-022`): ohne Treffer eine leere, gesetzte Liste — nie `null`, nie
// `404` (`LH-FA-REA-006` Boundary, `ADR-0081` Teilfrage 4); die
// Initialisierung mit `make(…, 0, …)` trägt das.
type readChangesResponse struct {
	Changes []readChangeResponse `json:"changes"`
}

// toReadChangeResponse übersetzt einen gelesenen Change in seine
// Antwortform (`SPEC-022`); die Commit-Position liegt auf der
// Quelltransaktion (`LH-FA-REA-004.a`).
func toReadChangeResponse(change inbound.ReadChange) readChangeResponse {
	return readChangeResponse{
		CommitPosition: int64(change.Position.Offset),
		ChangeID:       string(change.Change.ID),
		TransactionID:  string(change.Change.TransactionID),
		SourceTableID:  string(change.Change.SourceTableID),
		Schema:         change.Change.Schema,
		Table:          change.Change.Table,
		Sequence:       change.Change.Sequence,
		Operation:      string(change.Change.Operation),
		OldImage:       rowImage(change.Change.OldImage),
		NewImage:       rowImage(change.Change.NewImage),
		SchemaVersion:  string(change.Change.SchemaVersion),
		CommittedAt:    time.Unix(0, change.CommittedAt.UnixNanos).UTC().Format(time.RFC3339Nano),
		Origin:         string(change.Change.Origin.OrDefault()),
	}
}

// readChangesHandler trägt den lesenden Endpunkt `GET /changes`
// (`LH-FA-SST-006`, `LH-FA-REA-001` ff., `ADR-0081`): er übersetzt die
// Query-Parameter in `inbound.ReadChangesQuery` und ruft den Inbound Port —
// dieselbe Delegation an denselben `ChangeStorePort` wie der
// View-Direktzugriff, kein zweiter Lesepfad. Die Authn trägt die
// vorgelagerte `withToken`-Middleware (`reader`-Rechtsklasse).
//
// Die Adapter-Ebene prüft die Parameter-Form (unbekannter Parameter,
// nicht lesbare Zahl, fehlende Quelle, Position `< 1`); die Grenzen des
// Lese-Vertrags (`limit < 1`, `from > to`) trägt der Port-Kontrakt, ihre
// Fehler kommen über den Use Case zurück und laufen durch
// `writeDomainError` (`400`).
func readChangesHandler(useCase inbound.ReadChangesUseCase, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query, err := parseReadChangesQuery(r.URL.Query())
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		result, err := useCase.ReadChanges(r.Context(), query)
		if err != nil {
			writeDomainError(r.Context(), w, log, "ReadChanges", err)
			return
		}
		changes := make([]readChangeResponse, 0, len(result.Changes))
		for _, change := range result.Changes {
			changes = append(changes, toReadChangeResponse(change))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(readChangesResponse{Changes: changes})
	})
}

// parseReadChangesQuery übersetzt die Query-Parameter in
// `inbound.ReadChangesQuery`. Ein Parameter außerhalb der geschlossenen
// Menge endet sichtbar (`ADR-0081` Teilfrage 4).
func parseReadChangesQuery(values url.Values) (inbound.ReadChangesQuery, error) {
	for name := range values {
		if !readChangesParams[name] {
			return inbound.ReadChangesQuery{}, fmt.Errorf("unbekannter Query-Parameter %q", name)
		}
	}
	source := values.Get(readChangesParamSource)
	if source == "" {
		return inbound.ReadChangesQuery{}, fmt.Errorf("%s ist Pflichtfeld", readChangesParamSource)
	}
	start, err := readChangesPosition(source, readChangesParamFrom, values.Get(readChangesParamFrom))
	if err != nil {
		return inbound.ReadChangesQuery{}, err
	}
	end, err := readChangesPosition(source, readChangesParamTo, values.Get(readChangesParamTo))
	if err != nil {
		return inbound.ReadChangesQuery{}, err
	}
	limit, err := readChangesLimit(values.Get(readChangesParamLimit))
	if err != nil {
		return inbound.ReadChangesQuery{}, err
	}
	return inbound.ReadChangesQuery{
		Source: model.SourceID(source),
		Schema: values.Get(readChangesParamSchema),
		Table:  values.Get(readChangesParamTable),
		Start:  start,
		End:    end,
		Limit:  limit,
	}, nil
}

// readChangesPosition liest eine Positions-Grenze aus einem Query-Parameter
// (`SPEC-003`). Ein leerer Wert trägt keine Grenze; eine nicht lesbare Zahl
// und ein Wert unter 1 — eine Position 0 existiert nicht — enden sichtbar,
// nicht als stillschweigend nicht eingegrenzter Aufruf.
func readChangesPosition(source, name, raw string) (*model.SourcePosition, error) {
	if raw == "" {
		return nil, nil
	}
	offset, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s ist keine Ganzzahl: %q", name, raw)
	}
	if offset < 1 {
		return nil, domainerrors.ErrInvalidPosition
	}
	position, err := model.NewSourcePosition(model.SourceID(source), uint64(offset))
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// readChangesLimit liest das optionale Limit (`LH-FA-REA-003`). Ein leerer
// Wert trägt kein Limit — der Aufruf liest unbegrenzt, es gibt **kein**
// Default-Limit (`ADR-0081` Teilfrage 2). Eine nicht lesbare Zahl endet
// sichtbar; die Grenze `≥ 1` trägt der Port-Kontrakt.
func readChangesLimit(raw string) (*int, error) {
	if raw == "" {
		return nil, nil
	}
	limit, err := strconv.Atoi(raw)
	if err != nil {
		return nil, fmt.Errorf("%s ist keine Ganzzahl: %q", readChangesParamLimit, raw)
	}
	return &limit, nil
}
