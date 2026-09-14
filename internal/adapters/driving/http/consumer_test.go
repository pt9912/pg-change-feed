package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// newDefaultTestServer verdrahtet einen Server mit den beiden Test-Tokens
// (`testReaderToken`/`testAdminToken`, `server_test.go`) und den in `cfg`
// gesetzten Use Cases — nicht gesetzte Fähigkeiten bleiben `nil` und werden
// von den jeweiligen Tests nicht angesprochen.
func newDefaultTestServer(t *testing.T, cfg Config) *httptest.Server {
	t.Helper()
	cfg.Addr = "unused:0"
	cfg.TokenReader = testReaderToken
	cfg.TokenAdmin = testAdminToken
	srv := New(cfg)
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return ts
}

// doRequest baut und sendet einen Request mit optionalem Bearer-Token und
// optionalem Body gegen den Test-Server — Erweiterung von `postConsumer`
// (`server_test.go`) auf beliebige Methode/Pfad-Kombinationen.
func doRequest(t *testing.T, ts *httptest.Server, method, path, token, body string) *http.Response {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, ts.URL+path, reader)
	if err != nil {
		t.Fatalf("Request bauen: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Request senden: %v", err)
	}
	return resp
}

// fakeAcknowledgeConsumerUseCase trägt eine In-Memory-Fälschung
// (`ADR-0030`): dieselbe Monotonie-Invariante wie der reale
// `ConsumerStatePort` über `model.ConsumerPosition.Advance` — die
// Bestätigungslogik selbst ist bereits am Use Case getestet
// (`internal/application/usecase/acknowledge`), dieser Test belegt nur die
// Übersetzung HTTP↔Use-Case.
type fakeAcknowledgeConsumerUseCase struct {
	positions map[model.ConsumerID]model.ConsumerPosition
}

func newFakeAcknowledgeConsumerUseCase() *fakeAcknowledgeConsumerUseCase {
	return &fakeAcknowledgeConsumerUseCase{positions: map[model.ConsumerID]model.ConsumerPosition{}}
}

func (f *fakeAcknowledgeConsumerUseCase) Acknowledge(_ context.Context, command inbound.AcknowledgeConsumerCommand) (inbound.AcknowledgeConsumerResult, error) {
	if command.Consumer == "" {
		return inbound.AcknowledgeConsumerResult{}, domainerrors.ErrEmptyIdentifier
	}
	position, err := model.NewSourcePosition(command.Position.SourceID, command.Position.Offset)
	if err != nil {
		return inbound.AcknowledgeConsumerResult{}, err
	}
	current, ok := f.positions[command.Consumer]
	if !ok {
		current, _ = model.NewConsumerPosition(command.Consumer)
	}
	carried, err := current.Advance(position)
	if err != nil {
		return inbound.AcknowledgeConsumerResult{}, err
	}
	f.positions[command.Consumer] = carried
	return inbound.AcknowledgeConsumerResult{Position: carried}, nil
}

// fakeAcknowledgeFailingUseCase trägt einen nicht-domänenspezifischen
// Fehler — die Grenze zwischen `400` (Domänen-Invariante) und `500`
// (unerwarteter interner Fehler).
type fakeAcknowledgeFailingUseCase struct{ err error }

func (f fakeAcknowledgeFailingUseCase) Acknowledge(context.Context, inbound.AcknowledgeConsumerCommand) (inbound.AcknowledgeConsumerResult, error) {
	return inbound.AcknowledgeConsumerResult{}, f.err
}

// TestAcknowledgeConsumerAdminTokenBestaetigt trägt den Happy-Path
// (`LH-FA-CON-004`): ein gültiges `admin`-Token bestätigt die Position, die
// Antwort trägt die übersetzten Felder.
func TestAcknowledgeConsumerAdminTokenBestaetigt(t *testing.T) {
	ts := newDefaultTestServer(t, Config{AcknowledgeConsumer: newFakeAcknowledgeConsumerUseCase()})
	resp := doRequest(t, ts, http.MethodPost, "/consumers/acknowledge", testAdminToken,
		`{"consumer_id":"c1","source_id":"src-1","offset":42}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
	var decoded acknowledgeConsumerResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if decoded.ConsumerID != "c1" || decoded.SourceID != "src-1" || decoded.Offset != 42 {
		t.Fatalf("Antwort: %+v", decoded)
	}
}

// TestAcknowledgeConsumerReaderTokenEndetMit403 trägt `ADR-0057` Teilfrage 3:
// `AcknowledgeConsumer` ist ein schreibender Endpunkt, ein `reader`-Token
// erreicht ihn nicht.
func TestAcknowledgeConsumerReaderTokenEndetMit403(t *testing.T) {
	ts := newDefaultTestServer(t, Config{AcknowledgeConsumer: newFakeAcknowledgeConsumerUseCase()})
	resp := doRequest(t, ts, http.MethodPost, "/consumers/acknowledge", testReaderToken,
		`{"consumer_id":"c1","source_id":"src-1","offset":42}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Status: %d (Erwartung: 403)", resp.StatusCode)
	}
}

// TestAcknowledgeConsumerLeereKennungEndetMit400 trägt die
// Domänen-Invariante (`domainerrors.ErrEmptyIdentifier`) als `400`.
func TestAcknowledgeConsumerLeereKennungEndetMit400(t *testing.T) {
	ts := newDefaultTestServer(t, Config{AcknowledgeConsumer: newFakeAcknowledgeConsumerUseCase()})
	resp := doRequest(t, ts, http.MethodPost, "/consumers/acknowledge", testAdminToken,
		`{"consumer_id":"","source_id":"src-1","offset":42}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
}

// TestAcknowledgeConsumerInternerFehlerEndetMit500 trägt die Abgrenzung zur
// Domänen-Invariante: ein nicht-domänenspezifischer Fehler endet als `500`.
func TestAcknowledgeConsumerInternerFehlerEndetMit500(t *testing.T) {
	ts := newDefaultTestServer(t, Config{AcknowledgeConsumer: fakeAcknowledgeFailingUseCase{err: errors.New("speicher nicht erreichbar")}})
	resp := doRequest(t, ts, http.MethodPost, "/consumers/acknowledge", testAdminToken,
		`{"consumer_id":"c1","source_id":"src-1","offset":42}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Status: %d (Erwartung: 500)", resp.StatusCode)
	}
}

// fakeGetConsumerPositionUseCase liest aus derselben Kartenform wie
// `fakeAcknowledgeConsumerUseCase`, unabhängig instanziiert je Test.
type fakeGetConsumerPositionUseCase struct {
	positions map[model.ConsumerID]model.ConsumerPosition
}

func (f fakeGetConsumerPositionUseCase) Position(_ context.Context, query inbound.GetConsumerPositionQuery) (inbound.GetConsumerPositionResult, error) {
	if query.Consumer == "" {
		return inbound.GetConsumerPositionResult{}, domainerrors.ErrEmptyIdentifier
	}
	position, ok := f.positions[query.Consumer]
	if !ok {
		position, _ = model.NewConsumerPosition(query.Consumer)
	}
	return inbound.GetConsumerPositionResult{Position: position}, nil
}

// TestGetConsumerPositionReaderTokenLiefertPosition trägt `ADR-0057`
// Teilfrage 3: `GetConsumerPosition` ist ein lesender Endpunkt, ein
// `reader`-Token erreicht ihn.
func TestGetConsumerPositionReaderTokenLiefertPosition(t *testing.T) {
	consumer := model.ConsumerID("c1")
	position, err := model.NewSourcePosition("src-1", 7)
	if err != nil {
		t.Fatalf("Position bauen: %v", err)
	}
	useCase := fakeGetConsumerPositionUseCase{positions: map[model.ConsumerID]model.ConsumerPosition{
		consumer: {ConsumerID: consumer, Position: position},
	}}
	ts := newDefaultTestServer(t, Config{GetConsumerPosition: useCase})
	resp := doRequest(t, ts, http.MethodGet, "/consumers/position?consumer_id=c1", testReaderToken, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
	var decoded getConsumerPositionResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if decoded.ConsumerID != "c1" || decoded.SourceID != "src-1" || decoded.Offset != 7 || !decoded.Acknowledged {
		t.Fatalf("Antwort: %+v", decoded)
	}
}

// TestGetConsumerPositionAdminTokenLiefertPosition trägt die Hierarchie aus
// `ADR-0057` Teilfrage 3: ein `admin`-Token erreicht auch lesende Endpunkte.
func TestGetConsumerPositionAdminTokenLiefertPosition(t *testing.T) {
	useCase := fakeGetConsumerPositionUseCase{positions: map[model.ConsumerID]model.ConsumerPosition{}}
	ts := newDefaultTestServer(t, Config{GetConsumerPosition: useCase})
	resp := doRequest(t, ts, http.MethodGet, "/consumers/position?consumer_id=nie-bestaetigt", testAdminToken, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
	var decoded getConsumerPositionResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if decoded.Acknowledged {
		t.Fatalf("Antwort: %+v (Erwartung: acknowledged=false, LH-FA-CON-005 Boundary)", decoded)
	}
}

// TestGetConsumerPositionOhneKennungEndetMit400 trägt die Formgrenze des
// Query-Parameters.
func TestGetConsumerPositionOhneKennungEndetMit400(t *testing.T) {
	ts := newDefaultTestServer(t, Config{GetConsumerPosition: fakeGetConsumerPositionUseCase{positions: map[model.ConsumerID]model.ConsumerPosition{}}})
	resp := doRequest(t, ts, http.MethodGet, "/consumers/position", testAdminToken, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
}

// fakeRemoveConsumerUseCase trägt eine In-Memory-Fälschung der
// Registrierung: die Rückkehr meldet, ob der Consumer entfernt war
// (`LH-FA-CON-006` Idempotenz).
type fakeRemoveConsumerUseCase struct {
	registered map[model.ConsumerID]bool
}

func (f fakeRemoveConsumerUseCase) Remove(_ context.Context, command inbound.RemoveConsumerCommand) (inbound.RemoveConsumerResult, error) {
	if command.Consumer == "" {
		return inbound.RemoveConsumerResult{}, domainerrors.ErrEmptyIdentifier
	}
	removed := f.registered[command.Consumer]
	delete(f.registered, command.Consumer)
	return inbound.RemoveConsumerResult{Removed: removed}, nil
}

// TestRemoveConsumerAdminTokenEntfernt trägt den Happy-Path
// (`LH-FA-CON-006`): ein registrierter Consumer meldet `removed=true`.
func TestRemoveConsumerAdminTokenEntfernt(t *testing.T) {
	useCase := fakeRemoveConsumerUseCase{registered: map[model.ConsumerID]bool{"c1": true}}
	ts := newDefaultTestServer(t, Config{RemoveConsumer: useCase})
	resp := doRequest(t, ts, http.MethodPost, "/consumers/remove", testAdminToken, `{"consumer_id":"c1"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
	var decoded removeConsumerResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if decoded.ConsumerID != "c1" || !decoded.Removed {
		t.Fatalf("Antwort: %+v", decoded)
	}
}

// TestRemoveConsumerReaderTokenEndetMit403 trägt `ADR-0057` Teilfrage 3:
// `RemoveConsumer` ist ein administrativer Endpunkt.
func TestRemoveConsumerReaderTokenEndetMit403(t *testing.T) {
	useCase := fakeRemoveConsumerUseCase{registered: map[model.ConsumerID]bool{}}
	ts := newDefaultTestServer(t, Config{RemoveConsumer: useCase})
	resp := doRequest(t, ts, http.MethodPost, "/consumers/remove", testReaderToken, `{"consumer_id":"c1"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Status: %d (Erwartung: 403)", resp.StatusCode)
	}
}

// TestRemoveConsumerUnbekannteKennungBleibtIdempotent trägt `slice-060` §2
// DoD explizit: ein nie registrierter Consumer meldet `removed=false` mit
// `200`, nicht `404` — derselbe idempotente Ausgang wie CLI/SQL
// (`LH-FA-SST-006` Boundary: fachlich gleichwertiges Ergebnis über alle
// Zugriffswege).
func TestRemoveConsumerUnbekannteKennungBleibtIdempotent(t *testing.T) {
	useCase := fakeRemoveConsumerUseCase{registered: map[model.ConsumerID]bool{}}
	ts := newDefaultTestServer(t, Config{RemoveConsumer: useCase})
	resp := doRequest(t, ts, http.MethodPost, "/consumers/remove", testAdminToken, `{"consumer_id":"nie-registriert"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200, kein 404)", resp.StatusCode)
	}
	var decoded removeConsumerResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if decoded.Removed {
		t.Fatalf("Antwort: %+v (Erwartung: removed=false)", decoded)
	}
}
