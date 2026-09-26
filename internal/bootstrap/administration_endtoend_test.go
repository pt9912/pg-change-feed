package bootstrap

import (
	"context"
	stderrors "errors"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/excludecolumn"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/includecolumn"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// TestAdministrationRequestColumnEndToEndAgainstPostgreSQL trägt den realen
// Spaltenausschluss-Weg gegen eine PostgreSQL-Instanz (`LH-FA-CFG-005`,
// `ADR-0059`): der SQL-Funktionsaufruf legt einen Antrag an, der Adapter
// liest ihn, `applyAdministrationRequest` verarbeitet ihn über den Inbound
// Port und der Ergebnis-Vermerk schreibt den Ausgang in den Antrags-Datensatz
// — Happy Path (`applied`, vorhandene Spalte) und Negative-Fall (`failed`
// samt Fehlertext, nicht existierende Spalte, `LH-FA-CFG-005` Negative). Das
// ist der Kompositionsteil, den der Adapter-Test nicht erreichen kann:
// `applyAdministrationRequest` lebt in dieser Composition Root, die
// Spaltenprüfung läuft über den realen `ColumnExclusionPort` des
// Aktivierungs-Adapters gegen den Katalog. Der reale Ende-zu-Ende-Beleg am
// laufenden Feed-Container liegt in `tools/harness/run-integration-tests.sh`
// (`make test-integration`, `LH-FA-CFG-005`).
//
// Kopplung: der Schema-Stand dieses Laufs kommt aus dem d-migrate-Rollout,
// den der Lauf-Aufruf vor dem `internal/bootstrap`-Aufruf anwendet
// (`tools/schema/apply-rollout.sh`) — er trägt `cdc.administration_request`
// und die sieben Antrags-Funktionen
// (`tools/schema/nacharbeit-administration.sql`).
func TestAdministrationRequestColumnEndToEndAgainstPostgreSQL(t *testing.T) {
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	const (
		sourceID  = model.SourceID("src-administration-e2e")
		testTable = "admin_column_e2e"
	)
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Administrations-Quelle (E2E)') ON CONFLICT (source_id) DO NOTHING",
		string(sourceID),
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	// Die Vorbedingungs-Tabelle trägt beide Spalten-Zustände: die eine
	// existiert (Happy Path), die andere nicht (Negative-Fall).
	if _, err := pool.Exec(ctx,
		"CREATE TABLE IF NOT EXISTS public."+testTable+" (id integer PRIMARY KEY, secret text)",
	); err != nil {
		t.Fatalf("Quell-Tabelle anlegen: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DROP TABLE IF EXISTS public."+testTable)
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE source_id = $1", string(sourceID))
	})

	requests, err := postgresstorage.NewAdministrationRequest(ctx, dsn, postgresstorage.WithLog(&recordingLog{}))
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(requests.Close)
	columns, err := postgresstorage.NewTableActivation(ctx, dsn, postgresstorage.WithLog(&recordingLog{}))
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(columns.Close)

	// Die laufende Erfassung desselben Prozesses: die Test-Tabelle ist
	// gebunden, der real verarbeitete Antrag trägt seinen Ausschlussstand
	// nach (`applyAdministrationRequest`, `ADR-0059` Teilfrage 3).
	assembler, err := mapper.NewAssembler(sourceID, map[string]mapper.TableBinding{
		"public." + testTable: {TableID: "tbl-admin-column-e2e", SchemaVersion: "sv-admin-column-e2e"},
	}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}

	deps := administrationDeps{
		requests:       requests,
		excludeColumns: excludecolumn.NewExcludeColumnService(columns),
		includeColumns: includecolumn.NewIncludeColumnService(columns),
		assembler:      assembler,
		log:            &recordingLog{},
	}

	// Happy Path: eine vorhandene Spalte.
	var excludeID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.exclude_column($1, $2, $3, $4)", string(sourceID), "public", testTable, "secret",
	).Scan(&excludeID); err != nil {
		t.Fatalf("cdc.exclude_column: %v", err)
	}
	excludeRequest := readPendingAdministrationRequest(t, ctx, requests, model.AdministrationRequestID(excludeID))
	if excludeRequest.Kind != model.AdministrationRequestExcludeColumn || excludeRequest.Column != "secret" {
		t.Fatalf("exclude_column-Antrag = %+v, wollen Antragsart exclude_column mit Spalte secret", excludeRequest)
	}
	if err := applyAdministrationRequest(ctx, deps, excludeRequest); err != nil {
		t.Fatalf("applyAdministrationRequest(exclude_column, vorhandene Spalte) = %v, wollen nil", err)
	}
	if err := requests.MarkApplied(ctx, excludeRequest.ID); err != nil {
		t.Fatalf("MarkApplied: %v", err)
	}
	if got := readAdministrationRequestStatus(t, ctx, pool, excludeID); got != "applied" {
		t.Fatalf("Status exclude_column-Antrag = %q, wollen applied", got)
	}
	// Die Filterwirkung liegt an der laufenden Bindung: der reale Antrag
	// trägt `secret` in den Ausschlussstand nach, ein folgender Change führt
	// die Spalte nicht mehr.
	if image := assemblerRowImage(t, assembler, 1, "public", testTable); image != `{"id":"1"}` {
		t.Fatalf("Row Image nach dem real verarbeiteten exclude_column-Antrag = %s, wollen ohne den ausgeschlossenen Schlüssel secret", image)
	}

	// Negative-Fall: dieselbe Tabelle, eine nicht existierende Spalte.
	var includeID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.include_column($1, $2, $3, $4)", string(sourceID), "public", testTable, "does_not_exist",
	).Scan(&includeID); err != nil {
		t.Fatalf("cdc.include_column: %v", err)
	}
	includeRequest := readPendingAdministrationRequest(t, ctx, requests, model.AdministrationRequestID(includeID))
	if err := applyAdministrationRequest(ctx, deps, includeRequest); !stderrors.Is(err, inbound.ErrSourceColumnMissing) {
		t.Fatalf("applyAdministrationRequest(include_column, fehlende Spalte) = %v, wollen ErrSourceColumnMissing", err)
	} else if err := requests.MarkFailed(ctx, includeRequest.ID, err.Error()); err != nil {
		t.Fatalf("MarkFailed: %v", err)
	}
	status, message := readAdministrationRequestStatusAndMessage(t, ctx, pool, includeID)
	if status != "failed" {
		t.Fatalf("Status include_column-Antrag = %q, wollen failed", status)
	}
	if !strings.Contains(message, inbound.ErrSourceColumnMissing.Error()) {
		t.Fatalf("Fehlertext = %q, wollen %q enthalten", message, inbound.ErrSourceColumnMissing.Error())
	}
	// Ein gescheiterter Einschluss-Antrag lässt den geführten Ausschluss
	// stehen: der Fehler endet vor dem Assembler-Nachtrag.
	if image := assemblerRowImage(t, assembler, 2, "public", testTable); image != `{"id":"1"}` {
		t.Fatalf("Row Image nach dem gescheiterten include_column-Antrag = %s, wollen ohne den ausgeschlossenen Schlüssel secret", image)
	}
}

// readPendingAdministrationRequest liest genau den übergebenen Antrag über
// den Port zurück — derselbe Lesezugriffsweg der Administrations-Goroutine.
func readPendingAdministrationRequest(t *testing.T, ctx context.Context, requests *postgresstorage.AdministrationRequestAdapter, id model.AdministrationRequestID) model.AdministrationRequest {
	t.Helper()
	pending, err := requests.ListPending(ctx)
	if err != nil {
		t.Fatalf("ListPending: %v", err)
	}
	for _, request := range pending {
		if request.ID == id {
			return request
		}
	}
	t.Fatalf("ListPending trägt den Antrag %q nicht: %+v", id, pending)
	return model.AdministrationRequest{}
}

func readAdministrationRequestStatus(t *testing.T, ctx context.Context, pool *pgxpool.Pool, id string) string {
	t.Helper()
	var status string
	if err := pool.QueryRow(ctx,
		"SELECT status FROM cdc.administration_request WHERE administration_request_id = $1", id,
	).Scan(&status); err != nil {
		t.Fatalf("Status lesen: %v", err)
	}
	return status
}

func readAdministrationRequestStatusAndMessage(t *testing.T, ctx context.Context, pool *pgxpool.Pool, id string) (string, string) {
	t.Helper()
	var status string
	var message *string
	if err := pool.QueryRow(ctx,
		"SELECT status, error_message FROM cdc.administration_request WHERE administration_request_id = $1", id,
	).Scan(&status, &message); err != nil {
		t.Fatalf("Status/Fehlertext lesen: %v", err)
	}
	if message == nil {
		return status, ""
	}
	return status, *message
}
