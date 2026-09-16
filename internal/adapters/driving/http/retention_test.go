package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// fakeRunRetentionUseCase trägt eine In-Memory-Fälschung des
// `RunRetentionUseCase` — die Bereinigungslogik selbst ist bereits am Use
// Case getestet (`internal/application/usecase/retention`), dieser Test
// belegt nur die Übersetzung HTTP↔Use-Case.
type fakeRunRetentionUseCase struct {
	deleted int
}

func (f fakeRunRetentionUseCase) Run(_ context.Context, command inbound.RunRetentionCommand) (inbound.RunRetentionResult, error) {
	if command.Source == "" {
		return inbound.RunRetentionResult{}, domainerrors.ErrEmptyIdentifier
	}
	return inbound.RunRetentionResult{Deleted: f.deleted}, nil
}

type fakeRunRetentionFailingUseCase struct{ err error }

func (f fakeRunRetentionFailingUseCase) Run(context.Context, inbound.RunRetentionCommand) (inbound.RunRetentionResult, error) {
	return inbound.RunRetentionResult{}, f.err
}

// TestRunRetentionAdminTokenLoescht trägt den Happy-Path
// (`LH-FA-RET-002`…`004`).
func TestRunRetentionAdminTokenLoescht(t *testing.T) {
	ts := newDefaultTestServer(t, Config{RunRetention: fakeRunRetentionUseCase{deleted: 3}})
	resp := doRequest(t, ts, http.MethodPost, "/retention/run", testAdminToken,
		`{"source":"src-1","min_age_nanos":3600000000000}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
	var decoded runRetentionResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if decoded.Deleted != 3 {
		t.Fatalf("Antwort: %+v", decoded)
	}
}

// TestRunRetentionReaderTokenEndetMit403 trägt `ADR-0057` Teilfrage 3:
// `RunRetention` ist ein administrativer Endpunkt.
func TestRunRetentionReaderTokenEndetMit403(t *testing.T) {
	ts := newDefaultTestServer(t, Config{RunRetention: fakeRunRetentionUseCase{}})
	resp := doRequest(t, ts, http.MethodPost, "/retention/run", testReaderToken,
		`{"source":"src-1","min_age_nanos":0}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Status: %d (Erwartung: 403)", resp.StatusCode)
	}
}

// TestRunRetentionNegativeDauerEndetMit400 trägt die Domänen-Invariante
// (`domainerrors.ErrNegativeDuration`, über `model.NewDuration`) — die
// Policy entsteht vor dem Use-Case-Aufruf über den Domänen-Konstruktor.
func TestRunRetentionNegativeDauerEndetMit400(t *testing.T) {
	ts := newDefaultTestServer(t, Config{RunRetention: fakeRunRetentionUseCase{}})
	resp := doRequest(t, ts, http.MethodPost, "/retention/run", testAdminToken,
		`{"source":"src-1","min_age_nanos":-1}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
}

// TestRunRetentionLeereQuelleEndetMit400 trägt die Domänen-Invariante
// (`domainerrors.ErrEmptyIdentifier`) aus dem Use Case selbst.
func TestRunRetentionLeereQuelleEndetMit400(t *testing.T) {
	ts := newDefaultTestServer(t, Config{RunRetention: fakeRunRetentionUseCase{}})
	resp := doRequest(t, ts, http.MethodPost, "/retention/run", testAdminToken,
		`{"source":"","min_age_nanos":0}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
}

// TestRunRetentionInternerFehlerEndetMit500 trägt die Abgrenzung zur
// Domänen-Invariante: ein nicht-domänenspezifischer Fehler (z. B. eine
// Speicherstörung) endet als `500`, nicht als `400`.
func TestRunRetentionInternerFehlerEndetMit500(t *testing.T) {
	ts := newDefaultTestServer(t, Config{RunRetention: fakeRunRetentionFailingUseCase{err: errors.New("speicher nicht erreichbar")}})
	resp := doRequest(t, ts, http.MethodPost, "/retention/run", testAdminToken,
		`{"source":"src-1","min_age_nanos":0}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Status: %d (Erwartung: 500)", resp.StatusCode)
	}
}

// TestRunRetentionUngueltigesJSONEndetMit400 trägt die Formgrenze des
// Request-Bodys an der Eingabeseite: derselbe Use Case liefert für einen
// **erreichbaren** Aufruf `404` (`inbound.ErrSourceTableMissing`), der nicht
// dekodierbare Body endet dagegen mit `400` — der Status folgt dem Body,
// nicht dem Fake (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`).
// Rot färbende Mutation: den `Decode`-Fehlerzweig fallenlassen und mit dem
// Nullwert weiterlaufen — dann liefert der Gegenproben-Fake `404` statt `400`.
func TestRunRetentionUngueltigesJSONEndetMit400(t *testing.T) {
	useCase := fakeRunRetentionFailingUseCase{err: inbound.ErrSourceTableMissing}
	ts := newDefaultTestServer(t, Config{RunRetention: useCase})

	resp := doRequest(t, ts, http.MethodPost, "/retention/run", testAdminToken, `{nicht-json`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400 für einen nicht dekodierbaren Body)", resp.StatusCode)
	}

	gegenprobe := doRequest(t, ts, http.MethodPost, "/retention/run", testAdminToken,
		`{"source":"src-1","min_age_nanos":0}`)
	defer gegenprobe.Body.Close()
	if gegenprobe.StatusCode != http.StatusNotFound {
		t.Fatalf("Gegenprobe-Status: %d (Erwartung: 404 — dieser Fake wird für einen gültigen Body erreicht)", gegenprobe.StatusCode)
	}
}
