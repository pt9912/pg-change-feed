package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// fakeDiagnoseUseCase trägt eine In-Memory-Fälschung des Inbound Ports
// (`ADR-0030`): Whitebox-Test des Adapters ohne reale Persistenz.
type fakeDiagnoseUseCase struct {
	result  inbound.DiagnoseResult
	err     error
	queries []inbound.DiagnoseQuery
}

func (f *fakeDiagnoseUseCase) Diagnose(_ context.Context, query inbound.DiagnoseQuery) (inbound.DiagnoseResult, error) {
	f.queries = append(f.queries, query)
	if f.err != nil {
		return inbound.DiagnoseResult{}, f.err
	}
	return f.result, nil
}

func newDiagnoseServer(t *testing.T, useCase inbound.DiagnoseUseCase) *httptest.Server {
	t.Helper()
	srv := New(Config{
		Addr:        "unused:0",
		TokenReader: testReaderToken,
		TokenAdmin:  testAdminToken,
		Diagnose:    useCase,
	})
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return ts
}

func getDiagnose(t *testing.T, ts *httptest.Server, token, path string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, ts.URL+path, nil)
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

// TestDiagnoseOhneTokenEndetMit401 trägt die Authn-Grenze des Endpunkts
// (`ADR-0132` Teilfrage 4).
func TestDiagnoseOhneTokenEndetMit401(t *testing.T) {
	ts := newDiagnoseServer(t, &fakeDiagnoseUseCase{})
	resp := getDiagnose(t, ts, "", "/diagnose?source=src-1")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Status: %d (Erwartung: 401)", resp.StatusCode)
	}
}

// TestDiagnoseReaderTokenLiestBericht trägt den Happy Path (`ADR-0132`
// Teilfrage 4): das `reader`-Token genügt, die Antwort spiegelt das
// Use-Case-Ergebnis 1:1.
func TestDiagnoseReaderTokenLiestBericht(t *testing.T) {
	age := 1.5
	errorClass := "schema"
	lag := 3.0
	backlog := int64(7)
	estimated := int64(10)
	useCase := &fakeDiagnoseUseCase{result: inbound.DiagnoseResult{
		HeartbeatAgeSeconds: &age,
		ErrorClass:          &errorClass,
		CaptureLag:          2.5,
		ConsumerLags:        []inbound.ConsumerLag{{ConsumerID: "c-1", Lag: &lag}},
		RetentionBlocker: &inbound.RetentionBlocker{
			ConsumerID: "c-1", Name: "Consumer 1", AcknowledgedPosition: 10, Backlog: &backlog,
		},
		StorageBytes: 4096,
		Backfill: []inbound.BackfillTableStatus{{
			Schema: "public", Table: "orders", Status: "completed", RowsCopied: 5,
			EstimatedRows: &estimated, WarnEstimatedSize: true, WarnDuration: false, ErrorMessage: "",
		}},
	}}
	ts := newDiagnoseServer(t, useCase)

	resp := getDiagnose(t, ts, testReaderToken, "/diagnose?source=src-1")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("Status: %d (Erwartung: 200), Body: %s", resp.StatusCode, raw)
	}
	if got := len(useCase.queries); got != 1 || useCase.queries[0].Source != "src-1" {
		t.Fatalf("Use-Case-Aufrufe = %+v, wollen genau einen mit Quelle src-1", useCase.queries)
	}
	var decoded diagnoseResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if decoded.HeartbeatAgeSeconds == nil || *decoded.HeartbeatAgeSeconds != age {
		t.Fatalf("heartbeat_age_seconds = %v, wollen %v", decoded.HeartbeatAgeSeconds, age)
	}
	if decoded.ErrorClass == nil || *decoded.ErrorClass != errorClass {
		t.Fatalf("error_class = %v, wollen %v", decoded.ErrorClass, errorClass)
	}
	if decoded.CaptureLag != 2.5 {
		t.Fatalf("capture_lag = %v, wollen 2.5", decoded.CaptureLag)
	}
	if len(decoded.ConsumerLags) != 1 || decoded.ConsumerLags[0].ConsumerID != "c-1" || *decoded.ConsumerLags[0].Lag != lag {
		t.Fatalf("consumer_lags = %+v", decoded.ConsumerLags)
	}
	if decoded.RetentionBlocker == nil || decoded.RetentionBlocker.ConsumerID != "c-1" || *decoded.RetentionBlocker.Backlog != backlog {
		t.Fatalf("retention_blocker = %+v", decoded.RetentionBlocker)
	}
	if decoded.StorageBytes != 4096 {
		t.Fatalf("storage_bytes = %v, wollen 4096", decoded.StorageBytes)
	}
	if len(decoded.Backfill) != 1 || decoded.Backfill[0].Table != "orders" || *decoded.Backfill[0].EstimatedRows != estimated {
		t.Fatalf("backfill = %+v", decoded.Backfill)
	}
}

// TestDiagnoseAdminTokenLiest trägt die Abdeckung der niedrigeren Klasse
// durch das `admin`-Token.
func TestDiagnoseAdminTokenLiest(t *testing.T) {
	ts := newDiagnoseServer(t, &fakeDiagnoseUseCase{})
	resp := getDiagnose(t, ts, testAdminToken, "/diagnose?source=src-1")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
}

// TestDiagnoseFehlendeQuelleEndetMit400 trägt das Pflichtfeld.
func TestDiagnoseFehlendeQuelleEndetMit400(t *testing.T) {
	useCase := &fakeDiagnoseUseCase{}
	ts := newDiagnoseServer(t, useCase)
	resp := getDiagnose(t, ts, testReaderToken, "/diagnose")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
	if len(useCase.queries) != 0 {
		t.Fatalf("Use-Case-Aufrufe = %d, wollen 0", len(useCase.queries))
	}
}

// TestDiagnoseUnbekannterParameterEndetMit400 trägt die geschlossene
// Parameter-Menge (`ADR-0132` Teilfrage 4).
func TestDiagnoseUnbekannterParameterEndetMit400(t *testing.T) {
	useCase := &fakeDiagnoseUseCase{}
	ts := newDiagnoseServer(t, useCase)
	resp := getDiagnose(t, ts, testReaderToken, "/diagnose?source=src-1&quelle=src-1")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
	raw, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(raw), "quelle") {
		t.Fatalf("Fehlertext nennt den Parameter nicht: %s", raw)
	}
	if len(useCase.queries) != 0 {
		t.Fatalf("Use-Case-Aufrufe = %d, wollen 0", len(useCase.queries))
	}
}

// TestDiagnoseSpeicherfehlerEndetMit500 trägt die Abgrenzung zum
// Port-Kontrakt: ein Fehler der Klasse `storage` endet als `500`.
func TestDiagnoseSpeicherfehlerEndetMit500(t *testing.T) {
	ts := newDiagnoseServer(t, &fakeDiagnoseUseCase{err: outbound.ErrDiagnosticsStorage})
	resp := getDiagnose(t, ts, testReaderToken, "/diagnose?source=src-1")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Status: %d (Erwartung: 500)", resp.StatusCode)
	}
}

// TestDiagnoseLeereMengenSindNieNull trägt dieselbe Zusage wie
// `readChangesResponse.Changes`: eine leere Menge trägt `[]`, nie `null`.
func TestDiagnoseLeereMengenSindNieNull(t *testing.T) {
	ts := newDiagnoseServer(t, &fakeDiagnoseUseCase{})
	resp := getDiagnose(t, ts, testReaderToken, "/diagnose?source=src-1")
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(raw), `"consumer_lags":null`) || strings.Contains(string(raw), `"backfill":null`) {
		t.Fatalf("Body = %s, wollen leere, gesetzte Listen statt null", raw)
	}
}

// Der Fake erfüllt den Inbound Port.
var _ inbound.DiagnoseUseCase = (*fakeDiagnoseUseCase)(nil)
