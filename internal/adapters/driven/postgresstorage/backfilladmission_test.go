package postgresstorage_test

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

func newBackfillAdmission(t *testing.T, f *backfillFixture) *postgresstorage.BackfillAdmissionAdapter {
	t.Helper()
	admission, err := postgresstorage.NewBackfillAdmission(context.Background(), f.dsn)
	if err != nil {
		t.Fatalf("NewBackfillAdmission: %v", err)
	}
	t.Cleanup(admission.Close)
	return admission
}

// Die Annahme hinterlässt Run-Zeile `queued` und Antragsvermerk `applied`
// **zugleich** (`ADR-0113` Festlegung 1); die Run-Zeile trägt die Werte des
// Runs und die Defaults der übrigen Spalten (`SPEC-029`).
func TestBackfillAdmitLeavesQueuedRunAndAppliedRequestTogether(t *testing.T) {
	f := newBackfillFixture(t)
	admission := newBackfillAdmission(t, f)
	requestedAt := time.Date(2026, 9, 24, 10, 0, 0, 123_000_000, time.UTC)

	f.pendingRequest("adm-happy", "public", "adm_happy")
	run := queuedRun(t, "adm-happy", "adm_happy", requestedAt, knownEstimate(t, 500))
	run.WarnEstimatedSize = true
	if err := admission.Admit(context.Background(), "adm-happy", run); err != nil {
		t.Fatalf("Admit: %v", err)
	}

	if got := f.requestStatus("adm-happy"); got != "applied" {
		t.Fatalf("Antragsstatus = %q, erwartet applied", got)
	}
	var (
		status, source, schema, table string
		requested                     time.Time
		started, finished             *time.Time
		position, estimate            *int64
		rowsCopied                    int64
		warnSize, warnDuration        bool
		errorMessage                  *string
	)
	f.scan(`SELECT status, source_id, schema_name, table_name, requested_at, started_at, finished_at,
	               snapshot_position, rows_copied, estimated_rows, warn_estimated_size, warn_duration, error_message
	        FROM cdc.backfill_run WHERE run_id = 'adm-happy'`, nil,
		&status, &source, &schema, &table, &requested, &started, &finished,
		&position, &rowsCopied, &estimate, &warnSize, &warnDuration, &errorMessage)
	if status != "queued" || source != backfillTestSource || schema != "public" || table != "adm_happy" {
		t.Fatalf("Run-Zeile = %s %s %s.%s", status, source, schema, table)
	}
	if !requested.Equal(requestedAt) {
		t.Fatalf("requested_at = %v, erwartet %v", requested, requestedAt)
	}
	if started != nil || finished != nil || position != nil || errorMessage != nil || rowsCopied != 0 || warnDuration {
		t.Fatalf("Defaults verletzt: started %v finished %v position %v error %v rows %d warnDuration %v",
			started, finished, position, errorMessage, rowsCopied, warnDuration)
	}
	if estimate == nil || *estimate != 500 || !warnSize {
		t.Fatalf("estimated_rows %v, warn_estimated_size %v — erwartet 500/true", estimate, warnSize)
	}
}

// Ein Antragsvermerk ohne `pending`-Zeile (Rollback-Fall) hinterlässt weder
// Run-Zeile noch Vermerk: kein Antrag, ein bereits `applied` und ein bereits
// `failed` vermerkter Antrag. Rot färbende Mutation: die Prüfung
// `RowsAffected() != 1` in `Admit` entfernen — die Run-Zeile bliebe stehen.
func TestBackfillAdmitRollsBackWhenTheRequestIsNotPending(t *testing.T) {
	f := newBackfillFixture(t)
	admission := newBackfillAdmission(t, f)
	now := time.Now().UTC()

	cases := []struct {
		name, table, wantStatus string
		prepare                 func(id string)
	}{
		{"Antrag besteht nicht", "adm_missing", "", func(string) {}},
		{"Antrag bereits applied", "adm_applied", "applied", func(id string) {
			f.pendingRequest(id, "public", "adm_applied")
			f.exec("UPDATE cdc.administration_request SET status = 'applied' WHERE administration_request_id = $1", id)
		}},
		{"Antrag bereits failed", "adm_failed", "failed", func(id string) {
			f.pendingRequest(id, "public", "adm_failed")
			f.exec("UPDATE cdc.administration_request SET status = 'failed', error_message = 'zuvor' WHERE administration_request_id = $1", id)
		}},
	}
	for _, c := range cases {
		id := "adm-notpending-" + c.table
		c.prepare(id)
		run := queuedRun(t, id, c.table, now, model.UnknownRowEstimate())
		err := admission.Admit(context.Background(), model.AdministrationRequestID(id), run)
		if !stderrors.Is(err, outbound.ErrBackfillRequestNotPending) {
			t.Fatalf("%s: Fehler = %v, erwartet ErrBackfillRequestNotPending", c.name, err)
		}
		if _, exists := f.runStatus(id); exists {
			t.Fatalf("%s: die Run-Zeile blieb stehen", c.name)
		}
		if c.wantStatus != "" && f.requestStatus(id) != c.wantStatus {
			t.Fatalf("%s: Antragsstatus = %q, erwartet unverändert %q", c.name, f.requestStatus(id), c.wantStatus)
		}
	}
}

// Ein zweiter Antrag bei aktivem Run (`queued` und `running`) endet als
// `ErrBackfillRunActive` ohne zweite Zeile und ohne Vermerk: der Antrag
// bleibt `pending`. Ein beendeter Run ist nicht aktiv, eine andere Tabelle
// ist unabhängig. Rot färbende Mutation: in `SelectActiveBackfillRun` den
// Status-Filter auf `status = 'queued'` verengen — der `running`-Fall gälte
// als frei.
func TestBackfillAdmitRefusesASecondActiveRunOfTheSameTable(t *testing.T) {
	f := newBackfillFixture(t)
	admission := newBackfillAdmission(t, f)
	runs, err := postgresstorage.NewBackfillRun(context.Background(), f.dsn)
	if err != nil {
		t.Fatalf("NewBackfillRun: %v", err)
	}
	t.Cleanup(runs.Close)
	now := time.Now().UTC()

	first := f.admit(admission, "adm-active-1", "adm_active", now, model.UnknownRowEstimate())

	assertRefused := func(id string) {
		t.Helper()
		f.pendingRequest(id, "public", "adm_active")
		err := admission.Admit(context.Background(), model.AdministrationRequestID(id), queuedRun(t, id, "adm_active", now, model.UnknownRowEstimate()))
		if !stderrors.Is(err, domainerrors.ErrBackfillRunActive) {
			t.Fatalf("Admit %s: Fehler = %v, erwartet ErrBackfillRunActive", id, err)
		}
		if _, exists := f.runStatus(id); exists {
			t.Fatalf("Admit %s: zweite Run-Zeile entstanden", id)
		}
		if got := f.requestStatus(id); got != "pending" {
			t.Fatalf("Admit %s: Antragsstatus = %q, erwartet pending", id, got)
		}
	}
	assertRefused("adm-active-2")

	started, err := first.Start(model.NewTimePoint(now.UnixNano()))
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := runs.MarkRunning(context.Background(), started); err != nil {
		t.Fatalf("MarkRunning: %v", err)
	}
	assertRefused("adm-active-3")

	f.admit(admission, "adm-active-other", "adm_active_other", now, model.UnknownRowEstimate())

	failed, err := started.Fail(model.NewTimePoint(now.UnixNano()), model.ErrorClassTransient, "beendet")
	if err != nil {
		t.Fatalf("Fail: %v", err)
	}
	if err := runs.Finish(context.Background(), failed); err != nil {
		t.Fatalf("Finish: %v", err)
	}
	f.admit(admission, "adm-active-4", "adm_active", now, model.UnknownRowEstimate())
}

// `estimated_rows` ist NULL, wenn die Schätzung unbekannt ist — nie 0; die
// bekannte Schätzung 0 (analysierte, leere Tabelle) bleibt 0 (`SPEC-029`).
// Beide Werte lesen über den Run-Zustands-Adapter wieder als das, was sie
// sind. Rot färbende Mutation: `mapper.EstimateArgument` liefert für
// „unbekannt“ `int64(0)` statt `nil` — die NULL-Prüfung meldet 0.
func TestBackfillAdmitStoresUnknownEstimateAsNullAndNotAsZero(t *testing.T) {
	f := newBackfillFixture(t)
	admission := newBackfillAdmission(t, f)
	runs, err := postgresstorage.NewBackfillRun(context.Background(), f.dsn)
	if err != nil {
		t.Fatalf("NewBackfillRun: %v", err)
	}
	t.Cleanup(runs.Close)
	now := time.Now().UTC()

	f.admit(admission, "adm-est-unknown", "adm_est_unknown", now, model.UnknownRowEstimate())
	f.admit(admission, "adm-est-zero", "adm_est_zero", now, knownEstimate(t, 0))

	if n := f.count("SELECT count(*) FROM cdc.backfill_run WHERE run_id = 'adm-est-unknown' AND estimated_rows IS NULL"); n != 1 {
		t.Fatal("unbekannte Schätzung ist nicht als NULL gespeichert")
	}
	if n := f.count("SELECT count(*) FROM cdc.backfill_run WHERE run_id = 'adm-est-zero' AND estimated_rows = 0"); n != 1 {
		t.Fatal("bekannte Schätzung 0 ist nicht als 0 gespeichert")
	}

	queued, err := runs.Queued(context.Background(), backfillTestSource)
	if err != nil {
		t.Fatalf("Queued: %v", err)
	}
	found := map[model.BackfillRunID]model.RowEstimate{}
	for _, run := range queued {
		found[run.ID] = run.EstimatedRows
	}
	if _, known := found["adm-est-unknown"].Rows(); known {
		t.Fatal("NULL liest als bekannte Schätzung")
	}
	if rows, known := found["adm-est-zero"].Rows(); !known || rows != 0 {
		t.Fatalf("Schätzung 0 liest als %d, bekannt %v", rows, known)
	}
}
