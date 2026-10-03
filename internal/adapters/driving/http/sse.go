package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// sseEventChange ist der `event:`-Typ jedes Change-Events (`SPEC-018`,
// Abschnitt `GET /changes/stream`).
const sseEventChange = "change"

// Die Query-Parameter des SSE-Endpunkts `GET /changes/stream` (`ADR-0133`
// Teilfrage 4): `schema`/`table`/`target`, dieselben Feldnamen wie beim
// gRPC-Stream und bei `GET /changes`, je optional und unabhängig.
const (
	streamChangesParamSchema = "schema"
	streamChangesParamTable  = "table"
	streamChangesParamTarget = "target"
)

// streamChangesParams trägt die geschlossene Parameter-Menge dieses
// Endpunkts; ein Parameter außerhalb der Menge endet mit `400`, dieselbe
// Begründung wie bei `GET /changes` (`ADR-0081` Teilfrage 4) — ein Filter,
// der still ignoriert wird, ändert den zugestellten Ergebnisstand.
var streamChangesParams = map[string]bool{
	streamChangesParamSchema: true,
	streamChangesParamTable:  true,
	streamChangesParamTarget: true,
}

// parseStreamChangesFilter liest den optionalen
// `schema`/`table`/`target`-Filter aus den Query-Parametern; ein Parameter
// außerhalb der geschlossenen Menge endet sichtbar (`ADR-0133` Teilfrage 4).
// Ein `target` wird nicht geprüft: der Vergleich im Speicher trifft einen
// Namen außerhalb des Alphabets nie, der Stream bleibt dann leer.
func parseStreamChangesFilter(values url.Values) (schema, table, target string, err error) {
	for name := range values {
		if !streamChangesParams[name] {
			return "", "", "", rejectParam(messagecode.RejectedParameterUnknown, "unbekannter Query-Parameter %q", name)
		}
	}
	return values.Get(streamChangesParamSchema), values.Get(streamChangesParamTable), values.Get(streamChangesParamTarget), nil
}

// changeSubscriber trägt die eine Fähigkeit, die dieser Adapter vom
// `Broadcaster` braucht (`ADR-0061` Teilfrage 1/2): registrieren und einen
// Kanal samt Abmelde-Funktion erhalten. `*grpcstream.Broadcaster` erfüllt
// dieses Interface strukturell; die Composition Root verdrahtet es. Der
// Adapter kennt den Driven-Broadcaster nicht über einen Import des
// `grpcstream`-Pakets — dieselbe lokal deklarierte Schnittstelle wie im
// gRPC-Adapter (`internal/adapters/driving/grpc/server.go`).
type changeSubscriber interface {
	Subscribe() (<-chan *model.Change, func())
}

// streamChange trägt das SSE-Nachrichtenschema eines Change-Events
// (`SPEC-018`, Abschnitt `GET /changes/stream`): dieselben Felder wie der
// Domain-Typ `model.Change` ohne `Origin`. Die
// Row Images stehen als eingebettete JSON-Werte; ein fehlendes Bild
// (Boundary) wird zu `null`.
type streamChange struct {
	ChangeID      string          `json:"change_id"`
	TransactionID string          `json:"transaction_id"`
	SourceTableID string          `json:"source_table_id"`
	Sequence      int64           `json:"sequence"`
	Operation     string          `json:"operation"`
	OldImage      json.RawMessage `json:"old_image"`
	NewImage      json.RawMessage `json:"new_image"`
	SchemaVersion string          `json:"schema_version"`
	Schema        string          `json:"schema"`
	Table         string          `json:"table"`
}

// rowImage liefert ein Row Image als eingebetteten JSON-Wert; ein leeres
// Bild wird zu `null`, damit das Event auch ohne Bild gültiges JSON bleibt
// (`LH-FA-CAP-008` Boundary).
func rowImage(image []byte) json.RawMessage {
	if len(image) == 0 {
		return json.RawMessage("null")
	}
	return json.RawMessage(image)
}

// toStreamChange übersetzt einen Domain-Change in seine SSE-Nachrichtenform
// (`SPEC-018`): dieselben Felder, Row Images unverändert übernommen.
func toStreamChange(change *model.Change) streamChange {
	return streamChange{
		ChangeID:      string(change.ID),
		TransactionID: string(change.TransactionID),
		SourceTableID: string(change.SourceTableID),
		Sequence:      change.Sequence,
		Operation:     string(change.Operation),
		OldImage:      rowImage(change.OldImage),
		NewImage:      rowImage(change.NewImage),
		SchemaVersion: string(change.SchemaVersion),
		Schema:        change.Schema,
		Table:         change.Table,
	}
}

// streamChangesHandler trägt den SSE-Endpunkt `GET /changes/stream`
// (`ADR-0061` Teilfrage 1/2): er registriert einen Empfänger am
// `Broadcaster` und schreibt jeden eintreffenden Change als
// `text/event-stream`-Event, das er je Event über `http.Flusher` sofort
// ausliefert. Der `Last-Event-ID`-Header wird weder gesendet noch
// ausgewertet — der Stream trägt kein Replay (`ADR-0061` Teilfrage 3);
// verpasste Changes holt ein Consumer über den bestehenden Lesezugriffsweg
// nach (Boundary). Der Handler hält die
// Verbindung bis zum Verbindungsende des Clients offen; sein
// `Subscribe`-Aufruf wird über die Abmelde-Funktion beim Verlassen
// freigegeben. Ohne verdrahteten `Broadcaster` antwortet der Endpunkt mit
// `503`, statt eine leere Verbindung offenzuhalten (`ADR-0061` Teilfrage 5).
// Die Authentifizierung trägt die vorgelagerte `withToken`-Middleware: ein
// Aufruf ohne oder mit unbekanntem Bearer-Token endet mit `401`, bevor
// dieser Handler und damit vor jedem SSE-Event läuft (`ADR-0061` Teilfrage 4).
func streamChangesHandler(subscriber changeSubscriber, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if subscriber == nil {
			writeError(w, http.StatusServiceUnavailable, "ChangeStream ohne Broadcaster verdrahtet", messagecode.WiringPrecondition)
			return
		}
		schema, table, target, err := parseStreamChangesFilter(r.URL.Query())
		if err != nil {
			writeBadRequest(w, err)
			return
		}
		flusher, ok := w.(http.Flusher)
		if !ok {
			writeError(w, http.StatusInternalServerError, "Antwort-Writer trägt kein http.Flusher", messagecode.InternalFallback)
			return
		}
		changes, cancel := subscriber.Subscribe()
		defer cancel()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()

		log.Info(r.Context(), "http: SSE-Stream geöffnet")
		defer log.Info(context.Background(), "http: SSE-Stream beendet")

		for {
			select {
			case <-r.Context().Done():
				return
			case change := <-changes:
				// Derselbe Filter wie beim gRPC-Stream, geprüft nach dem
				// `Subscribe()`-Aufruf (`ADR-0133`): der Broadcaster bleibt
				// unverändert ungefiltert.
				if !change.MatchesFilter(schema, table, target) {
					continue
				}
				payload, err := json.Marshal(toStreamChange(change))
				if err != nil {
					log.Warn(r.Context(), "http: SSE-Change nicht kodierbar", messagecode.LogKey, messagecode.WarnChangeNotEncodable, "error", err)
					return
				}
				if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", sseEventChange, payload); err != nil {
					return
				}
				flusher.Flush()
			}
		}
	})
}
