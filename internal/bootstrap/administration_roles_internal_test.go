package bootstrap

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgressnapshot"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/systemclock"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/backfill"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/disable"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/enable"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/excludecolumn"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/includecolumn"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// adminPathLoginDSN legt eine anmeldefähige Identität `IN ROLE <Gruppenrolle>`
// an (so legt sie `docs/user/benutzerhandbuch.md` §2 an) und liefert ihre
// DSN; die Bereinigung entfernt sie wieder. Die Identität ist kein
// Superuser und Eigentümer nichts — jedes Recht, das sie hat, kommt aus der
// Gruppenrolle (`tools/schema/nacharbeit-roles.sql`).
func adminPathLoginDSN(t *testing.T, admin *pgxpool.Pool, baseDSN, login, group string) string {
	t.Helper()
	ctx := context.Background()
	const password = "test-login-password"
	if _, err := admin.Exec(ctx, fmt.Sprintf(`DROP ROLE IF EXISTS %q`, login)); err != nil {
		t.Fatalf("Vorab-Aufräumen der Login-Identität %s: %v", login, err)
	}
	if _, err := admin.Exec(ctx, fmt.Sprintf(`CREATE ROLE %q LOGIN PASSWORD '%s' IN ROLE %q`, login, password, group)); err != nil {
		t.Fatalf("Login-Identität %s anlegen: %v", login, err)
	}
	t.Cleanup(func() {
		_, _ = admin.Exec(context.Background(), fmt.Sprintf(`DROP ROLE IF EXISTS %q`, login))
	})
	parsed, err := url.Parse(baseDSN)
	if err != nil {
		t.Fatalf("Basis-DSN nicht parsebar: %v", err)
	}
	parsed.User = url.UserPassword(login, password)
	return parsed.String()
}

// TestAdministrationPathRunsUnderLeastPrivilegeLogins belegt den ganzen
// Antragsverarbeitungs-Pfad der Administrations-Goroutine (`ADR-0050`,
// `processAdministrationRequests`) unter den Rollen, die die Verdrahtung ihm
// zuweist (`ADR-0047`, `LH-QA-SEC-001`, `LH-QA-SEC-002`): die Antrags-Queue
// und die Aktivierung über eine `cdc_admin`-Login-Identität, der
// Schema-Speicher über eine `cdc_capture`-Login-Identität — beide ohne
// Superuser-Recht und ohne Eigentum an einem `cdc`-Objekt. Der Lauf zieht je
// einen Antrag jeder der fünf Antragsarten durch `ListPending`, den Use Case,
// die laufende `Assembler`-Bindung und den Vermerk `applied`; ein
// fünfter Antrag (`include_column` auf eine fehlende Spalte) endet im Vermerk
// `failed` samt Fehlertext.
//
// Die Superuser-Schwester `TestAdministrationRequestColumnEndToEndAgainstPostgreSQL`
// belegt die Verarbeitungslogik; dieser Test belegt, dass jede Anweisung des
// Pfads ihre Rolle trägt — ein Recht, das die Rolle nicht hat (Grants in
// `tools/schema/nacharbeit-roles.sql`), scheitert dort mit SQLSTATE 42501.
//
// Rot färbende Mutation: den Grant `SELECT, UPDATE ON cdc.administration_request
// TO cdc_admin` aus der Rollout-Datei streichen (oder `UPDATE` allein), oder
// `SELECT, INSERT ON cdc.backfill_run TO cdc_admin` — der
// Test meldet die Anweisung, die an dem fehlenden Recht scheitert.
func TestAdministrationPathRunsUnderLeastPrivilegeLogins(t *testing.T) {
	baseDSN := os.Getenv("CDC_STORE_TEST_DSN")
	if baseDSN == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	ctx := context.Background()
	admin, err := pgxpool.New(ctx, baseDSN)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(admin.Close)

	const (
		sourceID    = model.SourceID("src-admin-roles")
		testTable   = "admin_roles_path"
		publication = "pgc_admin_roles_pub"
	)
	if _, err := admin.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Administrationspfad unter Rollen') ON CONFLICT (source_id) DO NOTHING",
		string(sourceID),
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	// Betriebs-Vorbedingung je aktivierter Tabelle (`nacharbeit-roles.sql`,
	// Grenze): die Quelltabelle gehört `cdc_admin`.
	if _, err := admin.Exec(ctx, "CREATE TABLE IF NOT EXISTS public."+testTable+" (id integer PRIMARY KEY, secret text)"); err != nil {
		t.Fatalf("Quell-Tabelle anlegen: %v", err)
	}
	if _, err := admin.Exec(ctx, "ALTER TABLE public."+testTable+" OWNER TO cdc_admin"); err != nil {
		t.Fatalf("Eigentümer-Übertragung: %v", err)
	}
	t.Cleanup(func() {
		cleanup := context.Background()
		_, _ = admin.Exec(cleanup, "DROP PUBLICATION IF EXISTS "+publication)
		_, _ = admin.Exec(cleanup, "DELETE FROM cdc.administration_request WHERE source_id = $1", string(sourceID))
		_, _ = admin.Exec(cleanup, "DELETE FROM cdc.schema_version WHERE source_table_id IN (SELECT source_table_id FROM cdc.source_table WHERE source_id = $1)", string(sourceID))
		_, _ = admin.Exec(cleanup, "DELETE FROM cdc.source_table WHERE source_id = $1", string(sourceID))
		_, _ = admin.Exec(cleanup, "DROP TABLE IF EXISTS public."+testTable)
	})

	adminDSN := adminPathLoginDSN(t, admin, baseDSN, "pgc_test_adminpath_admin", "cdc_admin")
	captureDSN := adminPathLoginDSN(t, admin, baseDSN, "pgc_test_adminpath_capture", "cdc_capture")

	requests, err := postgresstorage.NewAdministrationRequest(ctx, adminDSN, postgresstorage.WithLog(&recordingLog{}))
	if err != nil {
		t.Fatalf("NewAdministrationRequest mit cdc_admin-Login: %v", err)
	}
	t.Cleanup(requests.Close)
	activation, err := postgresstorage.NewTableActivation(ctx, adminDSN, postgresstorage.WithLog(&recordingLog{}))
	if err != nil {
		t.Fatalf("NewTableActivation mit cdc_admin-Login: %v", err)
	}
	t.Cleanup(activation.Close)
	schemaStore, err := postgresstorage.NewSchemaStore(ctx, captureDSN, postgresstorage.WithLog(&recordingLog{}))
	if err != nil {
		t.Fatalf("NewSchemaStore mit cdc_capture-Login: %v", err)
	}
	t.Cleanup(schemaStore.Close)
	// Der Backfill-Antragszweig läuft mit denselben Rollen, die die
	// Verdrahtung zuweist (`ADR-0113` Festlegung 1): die Annahme unter dem
	// `cdc_admin`-Login, Run-Zustand, Schreiber und Snapshot-Adapter unter
	// dem `cdc_capture`-Login.
	backfillAdmission, err := postgresstorage.NewBackfillAdmission(ctx, adminDSN, postgresstorage.WithLog(&recordingLog{}))
	if err != nil {
		t.Fatalf("NewBackfillAdmission mit cdc_admin-Login: %v", err)
	}
	t.Cleanup(backfillAdmission.Close)
	backfillRuns, err := postgresstorage.NewBackfillRun(ctx, captureDSN, postgresstorage.WithLog(&recordingLog{}))
	if err != nil {
		t.Fatalf("NewBackfillRun mit cdc_capture-Login: %v", err)
	}
	t.Cleanup(backfillRuns.Close)
	backfillWriter, err := postgresstorage.NewBackfillWriter(ctx, captureDSN, postgresstorage.WithLog(&recordingLog{}))
	if err != nil {
		t.Fatalf("NewBackfillWriter mit cdc_capture-Login: %v", err)
	}
	t.Cleanup(backfillWriter.Close)
	backfillSnapshot, err := postgressnapshot.New(captureDSN)
	if err != nil {
		t.Fatalf("postgressnapshot.New mit cdc_capture-Login: %v", err)
	}
	t.Cleanup(func() { _, _ = admin.Exec(context.Background(), "DELETE FROM cdc.backfill_run WHERE source_id = $1", string(sourceID)) })

	assembler, err := mapper.NewAssembler(sourceID, map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	log := &recordingLog{}
	backfillWake := newBackfillWake()
	deps := administrationDeps{
		requests:       requests,
		activation:     activation,
		enableTables:   enable.NewEnableTableService(activation),
		disableTables:  disable.NewDisableTableService(activation),
		excludeColumns: excludecolumn.NewExcludeColumnService(activation),
		includeColumns: includecolumn.NewIncludeColumnService(activation),
		backfill: backfill.NewBackfillTableService(backfill.Ports{
			Activation: activation,
			Exclusion:  activation,
			Schemas:    schemaStore,
			Snapshot:   backfillSnapshot,
			Admission:  backfillAdmission,
			Runs:       backfillRuns,
			Writer:     backfillWriter,
			Clock:      systemclock.New(),
		}),
		backfillWake:    backfillWake,
		schemaStore:     schemaStore,
		columnExclusion: activation,
		assembler:       assembler,
		publication:     publication,
		log:             log,
	}

	// request legt den Antrag über die SQL-Administration an, lässt die
	// Goroutine-Schleife einmal über ihn laufen und liest Status und
	// Fehlertext als Superuser zurück.
	requestOn := func(table, function, column string) (id, status, message string) {
		t.Helper()
		args := []any{string(sourceID), "public", table}
		call := "SELECT cdc." + function + "($1, $2, $3)"
		if column != "" {
			call = "SELECT cdc." + function + "($1, $2, $3, $4)"
			args = append(args, column)
		}
		if err := admin.QueryRow(ctx, call, args...).Scan(&id); err != nil {
			t.Fatalf("cdc.%s: %v", function, err)
		}
		processAdministrationRequests(ctx, deps)
		if err := admin.QueryRow(ctx,
			"SELECT status, COALESCE(error_message, '') FROM cdc.administration_request WHERE administration_request_id = $1", id,
		).Scan(&status, &message); err != nil {
			t.Fatalf("Status lesen: %v", err)
		}
		return id, status, message
	}
	expectApplied := func(function, column string) {
		t.Helper()
		if _, status, message := requestOn(testTable, function, column); status != "applied" {
			t.Fatalf("%s unter cdc_admin-/cdc_capture-Login: Status %q (%s), erwartet applied — Log: %v", function, status, message, log.messages)
		}
	}

	expectApplied("enable_table", "")
	tableID := ""
	if err := admin.QueryRow(ctx,
		"SELECT source_table_id FROM cdc.source_table WHERE source_id = $1 AND schema_name = 'public' AND table_name = $2",
		string(sourceID), testTable,
	).Scan(&tableID); err != nil {
		t.Fatalf("Bindungs-Zeile nach enable_table: %v", err)
	}
	var member int
	if err := admin.QueryRow(ctx,
		"SELECT count(*) FROM pg_publication_tables WHERE pubname = $1 AND schemaname = 'public' AND tablename = $2",
		publication, testTable,
	).Scan(&member); err != nil || member != 1 {
		t.Fatalf("Publication-Mitgliedschaft nach enable_table: %d (%v), erwartet 1", member, err)
	}

	expectApplied("exclude_column", "secret")
	expectApplied("include_column", "secret")

	// Ein fehlgeschlagener Antrag endet im Vermerk failed samt Fehlertext.
	if _, status, message := requestOn(testTable, "include_column", "does_not_exist"); status != "failed" || message == "" {
		t.Fatalf("include_column auf fehlende Spalte: Status %q, Fehlertext %q, erwartet failed samt Text", status, message)
	}

	// Der Backfill-Antrag (`LH-FA-CAP-009`): die Annahme legt unter dem
	// `cdc_admin`-Login die Run-Zeile `queued` an und vermerkt den Antrag
	// `applied`, die Zeilenzahl wird unter dem `cdc_capture`-Login geschätzt
	// (unbekannt bleibt NULL, nie 0: die Tabelle ist nie analysiert), der
	// Worker wird geweckt. Ein zweiter Antrag derselben Tabelle bei aktivem
	// Run endet `failed` samt Text, ohne zweite Run-Zeile und ohne Signal; ein
	// Antrag für eine nicht aktivierte Tabelle endet `failed` ohne Run-Zeile.
	backfillID, status, message := requestOn(testTable, "backfill_table", "")
	if status != "applied" {
		t.Fatalf("backfill_table unter cdc_admin-/cdc_capture-Login: Status %q (%s), erwartet applied — Log: %v", status, message, log.messages)
	}
	var runStatus string
	var estimated *int64
	if err := admin.QueryRow(ctx, "SELECT status, estimated_rows FROM cdc.backfill_run WHERE run_id = $1", backfillID).Scan(&runStatus, &estimated); err != nil {
		t.Fatalf("Run-Zeile nach backfill_table: %v", err)
	}
	if runStatus != "queued" || estimated != nil {
		t.Fatalf("Run-Zeile nach der Annahme: Status %q, geschätzte Zeilenzahl %v — erwartet queued und unbekannt (NULL)", runStatus, estimated)
	}
	if len(backfillWake) != 1 {
		t.Fatalf("Wecksignale nach der Annahme = %d, erwartet 1", len(backfillWake))
	}
	<-backfillWake

	secondID, status, message := requestOn(testTable, "backfill_table", "")
	if status != "failed" || !strings.Contains(message, "aktiver Backfill-Run") {
		t.Fatalf("zweiter backfill_table bei aktivem Run: Status %q, Fehlertext %q, erwartet failed mit dem Grund „aktiver Backfill-Run“", status, message)
	}
	var runRows int
	if err := admin.QueryRow(ctx, "SELECT count(*) FROM cdc.backfill_run WHERE source_id = $1 AND table_name = $2", string(sourceID), testTable).Scan(&runRows); err != nil || runRows != 1 {
		t.Fatalf("Run-Zeilen der Tabelle nach dem abgelehnten Antrag: %d (%v), erwartet 1", runRows, err)
	}
	if err := admin.QueryRow(ctx, "SELECT count(*) FROM cdc.backfill_run WHERE run_id = $1", secondID).Scan(&runRows); err != nil || runRows != 0 {
		t.Fatalf("Run-Zeile des abgelehnten Antrags: %d (%v), erwartet 0", runRows, err)
	}
	if len(backfillWake) != 0 {
		t.Fatal("ein abgelehnter Backfill-Antrag weckte den Worker")
	}

	notActivatedID, status, message := requestOn("admin_roles_not_activated", "backfill_table", "")
	if status != "failed" || !strings.Contains(message, "nicht aktiviert") {
		t.Fatalf("backfill_table auf nicht aktivierte Tabelle: Status %q, Fehlertext %q, erwartet failed mit dem Grund „nicht aktiviert“", status, message)
	}
	if err := admin.QueryRow(ctx, "SELECT count(*) FROM cdc.backfill_run WHERE run_id = $1", notActivatedID).Scan(&runRows); err != nil || runRows != 0 {
		t.Fatalf("Run-Zeile für die nicht aktivierte Tabelle: %d (%v), erwartet 0", runRows, err)
	}

	expectApplied("disable_table", "")
	if err := admin.QueryRow(ctx,
		"SELECT count(*) FROM pg_publication_tables WHERE pubname = $1 AND tablename = $2", publication, testTable,
	).Scan(&member); err != nil || member != 0 {
		t.Fatalf("Publication-Mitgliedschaft nach disable_table: %d (%v), erwartet 0", member, err)
	}
	if strings.Contains(strings.Join(log.messages, "\n"), "Anträge lesen fehlgeschlagen") {
		t.Fatalf("die Schleife meldet einen Lesefehler der Antrags-Queue: %v", log.messages)
	}
}
