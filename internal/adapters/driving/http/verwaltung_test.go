package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeEnableTableUseCase trägt eine In-Memory-Fälschung des
// `TableActivationPort`: `tableExists` steuert `inbound.ErrSourceTableMissing`
// (`404`), `enabled` trägt den Idempotenz-Ausgang (`already_enabled`).
type fakeEnableTableUseCase struct {
	tableExists bool
	enabled     map[string]bool
}

func (f fakeEnableTableUseCase) Enable(_ context.Context, command inbound.EnableTableCommand) (inbound.EnableTableResult, error) {
	table, err := model.NewSourceTable(command.TableID, command.Source, command.Schema, command.Table)
	if err != nil {
		return inbound.EnableTableResult{}, err
	}
	if _, err := model.NewSchemaVersion(command.SchemaVersionID, command.TableID, command.Version); err != nil {
		return inbound.EnableTableResult{}, err
	}
	if !f.tableExists {
		return inbound.EnableTableResult{}, fmt.Errorf("%w: %s.%s", inbound.ErrSourceTableMissing, command.Schema, command.Table)
	}
	key := command.Schema + "." + command.Table
	already := f.enabled[key]
	f.enabled[key] = true
	return inbound.EnableTableResult{Table: table, AlreadyEnabled: already}, nil
}

// TestEnableTableAdminTokenAktiviert trägt den Happy-Path
// (`LH-FA-CFG-001`).
func TestEnableTableAdminTokenAktiviert(t *testing.T) {
	useCase := fakeEnableTableUseCase{tableExists: true, enabled: map[string]bool{}}
	ts := newDefaultTestServer(t, Config{EnableTable: useCase})
	resp := doRequest(t, ts, http.MethodPost, "/tables/enable", testAdminToken,
		`{"source":"src-1","schema":"public","table":"orders","table_id":"tbl-1","schema_version_id":"sv-1","version":1,"publication":"cdc_pub"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Status: %d (Erwartung: 201)", resp.StatusCode)
	}
	var decoded enableTableResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if decoded.TableID != "tbl-1" || decoded.Schema != "public" || decoded.Table != "orders" || decoded.AlreadyEnabled {
		t.Fatalf("Antwort: %+v", decoded)
	}
}

// TestEnableTableReaderTokenEndetMit403 trägt `ADR-0057` Teilfrage 3:
// `EnableTable` ist ein administrativer Endpunkt.
func TestEnableTableReaderTokenEndetMit403(t *testing.T) {
	useCase := fakeEnableTableUseCase{tableExists: true, enabled: map[string]bool{}}
	ts := newDefaultTestServer(t, Config{EnableTable: useCase})
	resp := doRequest(t, ts, http.MethodPost, "/tables/enable", testReaderToken,
		`{"source":"src-1","schema":"public","table":"orders","table_id":"tbl-1","schema_version_id":"sv-1","version":1,"publication":"cdc_pub"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Status: %d (Erwartung: 403)", resp.StatusCode)
	}
}

// TestEnableTableFehlendeTabelleEndetMit404 trägt die einheitliche
// Fehler-Mapping-Regel (`slice-060` §2 DoD): eine an der Quelle fehlende
// physische Tabelle (`inbound.ErrSourceTableMissing`) endet mit `404`, nicht
// `400`/`500`.
func TestEnableTableFehlendeTabelleEndetMit404(t *testing.T) {
	useCase := fakeEnableTableUseCase{tableExists: false, enabled: map[string]bool{}}
	ts := newDefaultTestServer(t, Config{EnableTable: useCase})
	resp := doRequest(t, ts, http.MethodPost, "/tables/enable", testAdminToken,
		`{"source":"src-1","schema":"public","table":"nicht_vorhanden","table_id":"tbl-1","schema_version_id":"sv-1","version":1,"publication":"cdc_pub"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Status: %d (Erwartung: 404)", resp.StatusCode)
	}
}

// TestEnableTableUngueltigeVersionEndetMit400 trägt die Domänen-Invariante
// (`domainerrors.ErrNonPositiveVersion`, über `model.NewSchemaVersion`).
func TestEnableTableUngueltigeVersionEndetMit400(t *testing.T) {
	useCase := fakeEnableTableUseCase{tableExists: true, enabled: map[string]bool{}}
	ts := newDefaultTestServer(t, Config{EnableTable: useCase})
	resp := doRequest(t, ts, http.MethodPost, "/tables/enable", testAdminToken,
		`{"source":"src-1","schema":"public","table":"orders","table_id":"tbl-1","schema_version_id":"sv-1","version":0,"publication":"cdc_pub"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
}

// fakeDisableTableUseCase trägt eine In-Memory-Fälschung analog zur
// realen `DisableTableService`-Zeilen-Ausgangs-Unterscheidung
// (`removed`/`retained`, `LH-FA-CFG-002` Out-of-Scope).
type fakeDisableTableUseCase struct {
	tableExists bool
	result      inbound.DisableTableResult
}

func (f fakeDisableTableUseCase) Disable(_ context.Context, command inbound.DisableTableCommand) (inbound.DisableTableResult, error) {
	if command.Source == "" || command.Schema == "" || command.Table == "" || command.Publication == "" {
		return inbound.DisableTableResult{}, domainerrors.ErrEmptyIdentifier
	}
	if !f.tableExists {
		return inbound.DisableTableResult{}, fmt.Errorf("%w: %s.%s", inbound.ErrSourceTableMissing, command.Schema, command.Table)
	}
	return f.result, nil
}

// TestDisableTableAdminTokenDeaktiviert trägt den Happy-Path
// (`LH-FA-CFG-002`).
func TestDisableTableAdminTokenDeaktiviert(t *testing.T) {
	useCase := fakeDisableTableUseCase{tableExists: true, result: inbound.DisableTableResult{Removed: true}}
	ts := newDefaultTestServer(t, Config{DisableTable: useCase})
	resp := doRequest(t, ts, http.MethodPost, "/tables/disable", testAdminToken,
		`{"source":"src-1","schema":"public","table":"orders","publication":"cdc_pub"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
	var decoded disableTableResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if !decoded.Removed || decoded.Retained {
		t.Fatalf("Antwort: %+v", decoded)
	}
}

// TestDisableTableReaderTokenEndetMit403 trägt `ADR-0057` Teilfrage 3:
// `DisableTable` ist ein administrativer Endpunkt.
func TestDisableTableReaderTokenEndetMit403(t *testing.T) {
	useCase := fakeDisableTableUseCase{tableExists: true}
	ts := newDefaultTestServer(t, Config{DisableTable: useCase})
	resp := doRequest(t, ts, http.MethodPost, "/tables/disable", testReaderToken,
		`{"source":"src-1","schema":"public","table":"orders","publication":"cdc_pub"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Status: %d (Erwartung: 403)", resp.StatusCode)
	}
}

// TestDisableTableFehlendeTabelleEndetMit404 trägt dieselbe Fehler-Mapping-
// Regel wie `TestEnableTableFehlendeTabelleEndetMit404`.
func TestDisableTableFehlendeTabelleEndetMit404(t *testing.T) {
	useCase := fakeDisableTableUseCase{tableExists: false}
	ts := newDefaultTestServer(t, Config{DisableTable: useCase})
	resp := doRequest(t, ts, http.MethodPost, "/tables/disable", testAdminToken,
		`{"source":"src-1","schema":"public","table":"nicht_vorhanden","publication":"cdc_pub"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Status: %d (Erwartung: 404)", resp.StatusCode)
	}
}

// fakeGetStatusUseCase trägt eine In-Memory-Fälschung analog zur realen
// `GetStatusService`-Fallunterscheidung.
type fakeGetStatusUseCase struct {
	tableExists bool
	result      inbound.GetStatusResult
}

func (f fakeGetStatusUseCase) Status(_ context.Context, query inbound.GetStatusQuery) (inbound.GetStatusResult, error) {
	if query.Source == "" || query.Schema == "" || query.Table == "" || query.Publication == "" {
		return inbound.GetStatusResult{}, domainerrors.ErrEmptyIdentifier
	}
	if !f.tableExists {
		return inbound.GetStatusResult{}, fmt.Errorf("%w: %s.%s", inbound.ErrSourceTableMissing, query.Schema, query.Table)
	}
	return f.result, nil
}

// TestGetStatusReaderTokenLiefertStatus trägt `ADR-0057` Teilfrage 3:
// `GetStatus` ist ein lesender Endpunkt.
func TestGetStatusReaderTokenLiefertStatus(t *testing.T) {
	useCase := fakeGetStatusUseCase{tableExists: true, result: inbound.GetStatusResult{Enabled: true}}
	ts := newDefaultTestServer(t, Config{GetStatus: useCase})
	resp := doRequest(t, ts, http.MethodGet,
		"/tables/status?source=src-1&schema=public&table=orders&publication=cdc_pub", testReaderToken, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
	var decoded getStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if !decoded.Enabled || decoded.Retained {
		t.Fatalf("Antwort: %+v", decoded)
	}
}

// TestGetStatusFehlendeParameterEndetMit400 trägt die Formgrenze der
// Query-Parameter.
func TestGetStatusFehlendeParameterEndetMit400(t *testing.T) {
	useCase := fakeGetStatusUseCase{tableExists: true}
	ts := newDefaultTestServer(t, Config{GetStatus: useCase})
	resp := doRequest(t, ts, http.MethodGet, "/tables/status?source=src-1", testAdminToken, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
}

// fakeListTablesUseCase trägt eine In-Memory-Fälschung analog zur realen
// `ListTablesService`-Trennung nach Erfassungs-Zustand.
type fakeListTablesUseCase struct {
	result inbound.ListTablesResult
}

func (f fakeListTablesUseCase) ListTables(_ context.Context, query inbound.ListTablesQuery) (inbound.ListTablesResult, error) {
	if query.Source == "" || query.Publication == "" {
		return inbound.ListTablesResult{}, domainerrors.ErrEmptyIdentifier
	}
	return f.result, nil
}

// TestListTablesReaderTokenLiefertListe trägt `ADR-0057` Teilfrage 3:
// `ListTables` ist ein lesender Endpunkt.
func TestListTablesReaderTokenLiefertListe(t *testing.T) {
	table, err := model.NewSourceTable("tbl-1", "src-1", "public", "orders")
	if err != nil {
		t.Fatalf("Tabelle bauen: %v", err)
	}
	useCase := fakeListTablesUseCase{result: inbound.ListTablesResult{Tables: []model.SourceTable{table}, Retained: []model.SourceTable{}}}
	ts := newDefaultTestServer(t, Config{ListTables: useCase})
	resp := doRequest(t, ts, http.MethodGet, "/tables?source=src-1&publication=cdc_pub", testReaderToken, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
	var decoded listTablesResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("Antwort dekodieren: %v", err)
	}
	if len(decoded.Tables) != 1 || decoded.Tables[0].TableID != "tbl-1" || len(decoded.Retained) != 0 {
		t.Fatalf("Antwort: %+v", decoded)
	}
}

// TestListTablesFehlendeParameterEndetMit400 trägt die Formgrenze der
// Query-Parameter.
func TestListTablesFehlendeParameterEndetMit400(t *testing.T) {
	useCase := fakeListTablesUseCase{}
	ts := newDefaultTestServer(t, Config{ListTables: useCase})
	resp := doRequest(t, ts, http.MethodGet, "/tables", testAdminToken, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Status: %d (Erwartung: 400)", resp.StatusCode)
	}
}
