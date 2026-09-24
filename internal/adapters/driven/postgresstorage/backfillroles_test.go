package postgresstorage_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Rollen-Tests dieser Datei führen die drei Backfill-Adapter unter der
// Rolle, die die Verdrahtung ihnen zuweist (`ADR-0113` Festlegung 1,
// `ADR-0047`, `LH-QA-SEC-001`, `LH-QA-SEC-002`): die Annahme als
// `cdc_admin`-Login, der Run-Zustand und der Schreiber als
// `cdc_capture`-Login — beide ohne Superuser-Recht und ohne Eigentum an
// einem `cdc`-Objekt. Die Rollen-Tests in `roles_test.go` führen rohe
// Anweisungen per `SET ROLE` aus; diese Datei belegt, dass die Anweisungen
// der Adapter selbst ihre Rolle tragen.

// login legt eine anmeldefähige Identität `IN ROLE <Gruppenrolle>` an (so
// legt sie `docs/user/benutzerhandbuch.md` §2 an) und liefert ihre DSN; die
// Bereinigung entfernt sie nach den Adaptern des Tests.
func (f *backfillFixture) login(name, group string) string {
	f.t.Helper()
	ctx := context.Background()
	const password = "test-login-password"
	f.exec(fmt.Sprintf(`DROP ROLE IF EXISTS %q`, name))
	f.exec(fmt.Sprintf(`CREATE ROLE %q LOGIN PASSWORD '%s' IN ROLE %q`, name, password, group))
	f.t.Cleanup(func() {
		_, _ = f.pool.Exec(ctx, fmt.Sprintf(`DROP ROLE IF EXISTS %q`, name))
	})
	parsed, err := url.Parse(f.dsn)
	if err != nil {
		f.t.Fatalf("Basis-DSN nicht parsebar: %v", err)
	}
	parsed.User = url.UserPassword(name, password)
	return parsed.String()
}

// Der Weg eines Backfills unter den zwei Rollen: die Annahme legt Run und
// Vermerk als `cdc_admin`-Login an, der Run-Zustand führt ihn als
// `cdc_capture`-Login (`queued` lesen, `running`, Fortschritt, Endzustand
// `failed`/`interrupted`), der Schreiber committet Blöcke und Run-Zeile als
// `cdc_capture`-Login. Rot färbende Mutation: einen Grant der Rollen-Datei
// streichen (`SELECT`/`INSERT` an `cdc_admin` auf `backfill_run` oder
// `administration_request`-`UPDATE`, `SELECT`/`UPDATE` an `cdc_capture` auf
// `backfill_run`, `INSERT` an `cdc_capture` auf `transaction`/`change`) —
// der Test meldet die Anweisung des Adapters, die mit SQLSTATE 42501
// scheitert.
func TestBackfillAdaptersRunUnderTheirRoles(t *testing.T) {
	f := newBackfillFixture(t)
	adminDSN := f.login("pgc_test_bfroles_admin", "cdc_admin")
	captureDSN := f.login("pgc_test_bfroles_capture", "cdc_capture")
	ctx := context.Background()
	base := time.Date(2026, 9, 24, 19, 0, 0, 0, time.UTC)

	// Annahme (cdc_admin), Run-Zustand und Schreiber (cdc_capture), bis zum
	// atomaren Commit.
	w := newWriterRunAs(t, f, "bfr-ok", 3_100_100, adminDSN, captureDSN)
	if status, ok := f.runStatus("bfr-ok"); !ok || status != "running" {
		t.Fatalf("Run nach Annahme und MarkRunning unter den Rollen = %q (%v), erwartet running", status, ok)
	}
	if got := f.requestStatus("bfr-ok"); got != "applied" {
		t.Fatalf("Antragsvermerk der Annahme unter cdc_admin = %q, erwartet applied", got)
	}
	tx, err := w.writer.Begin(ctx, w.run)
	if err != nil {
		t.Fatalf("Begin unter cdc_capture: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })
	if err := tx.AppendBlock(ctx, w.block(t, 1, 2, base)); err != nil {
		t.Fatalf("AppendBlock unter cdc_capture: %v", err)
	}
	if err := tx.Commit(ctx, w.completed(t, 2, base.Add(time.Minute))); err != nil {
		t.Fatalf("Commit unter cdc_capture: %v", err)
	}
	if status, _ := f.runStatus("bfr-ok"); status != "completed" {
		t.Fatalf("Run nach Commit unter cdc_capture = %q, erwartet completed", status)
	}
	if transactions, changes := f.backfillRows("bfr-ok"); transactions != 1 || changes != 2 {
		t.Fatalf("nach dem Commit unter cdc_capture: %d Transaktionen, %d Changes, erwartet 1/2", transactions, changes)
	}

	// Run-Zustand: Queued lesen, Endzustand failed und interrupted.
	admission, err := postgresstorage.NewBackfillAdmission(ctx, adminDSN)
	if err != nil {
		t.Fatalf("NewBackfillAdmission: %v", err)
	}
	t.Cleanup(admission.Close)
	runs, err := postgresstorage.NewBackfillRun(ctx, captureDSN)
	if err != nil {
		t.Fatalf("NewBackfillRun: %v", err)
	}
	t.Cleanup(runs.Close)

	failing := f.admit(admission, "bfr-fail", "bfr_fail", base, model.UnknownRowEstimate())
	interrupted := f.admit(admission, "bfr-int", "bfr_int", base.Add(time.Second), model.UnknownRowEstimate())
	queued, err := runs.Queued(ctx, backfillTestSource)
	if err != nil {
		t.Fatalf("Queued unter cdc_capture: %v", err)
	}
	seen := map[model.BackfillRunID]bool{}
	for _, run := range queued {
		seen[run.ID] = true
	}
	if !seen["bfr-fail"] || !seen["bfr-int"] {
		t.Fatalf("Queued unter cdc_capture lieferte %v, erwartet bfr-fail und bfr-int", seen)
	}
	failingRunning := mustRunning(t, runs, failing, base.Add(2*time.Second))
	failed, err := failingRunning.Fail(model.NewTimePoint(base.Add(3*time.Second).UnixNano()), model.ErrorClassStorage, "Testfehler")
	if err != nil {
		t.Fatalf("Fail (Domäne): %v", err)
	}
	if err := runs.Finish(ctx, failed); err != nil {
		t.Fatalf("Finish unter cdc_capture: %v", err)
	}
	if status, _ := f.runStatus("bfr-fail"); status != "failed" {
		t.Fatalf("Run nach Finish(failed) unter cdc_capture = %q", status)
	}
	mustRunning(t, runs, interrupted, base.Add(4*time.Second))
	if count, err := runs.InterruptRunning(ctx, backfillTestSource, model.NewTimePoint(base.Add(5*time.Second).UnixNano())); err != nil || count < 1 {
		t.Fatalf("InterruptRunning unter cdc_capture = %d, %v, erwartet mindestens einen Run", count, err)
	}
	if status, _ := f.runStatus("bfr-int"); status != "interrupted" {
		t.Fatalf("Run nach InterruptRunning unter cdc_capture = %q", status)
	}
}

// Jeder Adapter scheitert unter der Rolle des anderen mit SQLSTATE 42501:
// die Annahme als `cdc_capture` (kein `INSERT` auf `backfill_run`), der
// Run-Zustand und der Schreiber als `cdc_admin` (kein `UPDATE` auf
// `backfill_run`, kein `INSERT` auf `transaction`/`change`); als `cdc_capture` lässt
// sich eine `backfill_run`-Zeile weder anlegen noch löschen. Rot färbende
// Mutation: `INSERT` oder `DELETE` an `cdc_capture` bzw. `UPDATE` an
// `cdc_admin` auf `cdc.backfill_run` (Rollen-Datei) ergänzen — der
// jeweilige Fall meldet, dass die Anweisung gelang statt zu scheitern.
func TestBackfillAdaptersFailUnderTheOtherRole(t *testing.T) {
	f := newBackfillFixture(t)
	adminDSN := f.login("pgc_test_bfroles_other_admin", "cdc_admin")
	captureDSN := f.login("pgc_test_bfroles_other_capture", "cdc_capture")
	ctx := context.Background()
	base := time.Date(2026, 9, 24, 20, 0, 0, 0, time.UTC)

	// Die Annahme unter cdc_capture: kein INSERT auf backfill_run.
	f.pendingRequest("bfr-wrong-admit", "public", "bfr_wrong_admit")
	wrongAdmission, err := postgresstorage.NewBackfillAdmission(ctx, captureDSN)
	if err != nil {
		t.Fatalf("NewBackfillAdmission mit cdc_capture-Login: %v", err)
	}
	t.Cleanup(wrongAdmission.Close)
	run := queuedRun(t, "bfr-wrong-admit", "bfr_wrong_admit", base, model.UnknownRowEstimate())
	if err := wrongAdmission.Admit(ctx, "bfr-wrong-admit", run); !permissionDenied(err) {
		t.Fatalf("Admit unter cdc_capture: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
	if _, ok := f.runStatus("bfr-wrong-admit"); ok {
		t.Fatal("Admit unter cdc_capture hinterließ eine Run-Zeile")
	}
	if got := f.requestStatus("bfr-wrong-admit"); got != "pending" {
		t.Fatalf("Admit unter cdc_capture verändert den Antrag: %q, erwartet pending", got)
	}

	// Run-Zustand und Schreiber unter cdc_admin: kein UPDATE auf backfill_run
	// und kein INSERT auf transaction/change.
	w := newWriterRunAs(t, f, "bfr-wrong", 3_100_200, adminDSN, captureDSN)
	wrongRuns, err := postgresstorage.NewBackfillRun(ctx, adminDSN)
	if err != nil {
		t.Fatalf("NewBackfillRun mit cdc_admin-Login: %v", err)
	}
	t.Cleanup(wrongRuns.Close)
	interrupted, err := w.run.Interrupt(model.NewTimePoint(base.UnixNano()))
	if err != nil {
		t.Fatalf("Interrupt (Domäne): %v", err)
	}
	if err := wrongRuns.Finish(ctx, interrupted); !permissionDenied(err) {
		t.Fatalf("Finish unter cdc_admin: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
	if status, _ := f.runStatus("bfr-wrong"); status != "running" {
		t.Fatalf("Finish unter cdc_admin verändert den Run: %q, erwartet running", status)
	}
	wrongWriter, err := postgresstorage.NewBackfillWriter(ctx, adminDSN)
	if err != nil {
		t.Fatalf("NewBackfillWriter mit cdc_admin-Login: %v", err)
	}
	t.Cleanup(wrongWriter.Close)
	tx, err := wrongWriter.Begin(ctx, w.run)
	if err == nil {
		t.Cleanup(func() { _ = tx.Rollback(ctx) })
		err = tx.AppendBlock(ctx, w.block(t, 1, 1, base))
	}
	if !permissionDenied(err) {
		t.Fatalf("Schreiber unter cdc_admin (Begin/AppendBlock): erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
	if transactions, changes := f.backfillRows("bfr-wrong"); transactions != 0 || changes != 0 {
		t.Fatalf("Schreiber unter cdc_admin hinterließ %d Transaktionen, %d Changes", transactions, changes)
	}

	// Eine backfill_run-Zeile lässt sich als cdc_capture weder anlegen noch
	// löschen.
	capture, err := pgxpool.New(ctx, captureDSN)
	if err != nil {
		t.Fatalf("Verbindungsaufbau mit cdc_capture-Login: %v", err)
	}
	t.Cleanup(capture.Close)
	if _, err := capture.Exec(ctx,
		`INSERT INTO cdc.backfill_run (run_id, source_id, schema_name, table_name, status)
		 VALUES ('bfr-capture-insert', $1, 'public', 'bfr_capture', 'queued')`, backfillTestSource,
	); !permissionDenied(err) {
		t.Fatalf("INSERT auf backfill_run als cdc_capture-Login: erwartet SQLSTATE 42501, erhalten %v", err)
	}
	if _, err := capture.Exec(ctx, "DELETE FROM cdc.backfill_run WHERE run_id = 'bfr-wrong'"); !permissionDenied(err) {
		t.Fatalf("DELETE auf backfill_run als cdc_capture-Login: erwartet SQLSTATE 42501, erhalten %v", err)
	}
}
