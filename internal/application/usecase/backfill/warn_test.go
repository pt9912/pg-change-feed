package backfill

import (
	"testing"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

func estimateOf(t *testing.T, rows int64) model.RowEstimate {
	t.Helper()
	estimate, err := model.NewRowEstimate(rows)
	if err != nil {
		t.Fatal(err)
	}
	return estimate
}

// TestWarnsEstimatedSizeBoundary trägt Warnung (1) an der Richtgröße: genau
// auf ihr keine Warnung, eine Zeile darüber gesetzt; die bekannte `0` und
// eine unbekannte Schätzung warnen nicht.
func TestWarnsEstimatedSizeBoundary(t *testing.T) {
	cases := []struct {
		name     string
		estimate model.RowEstimate
		want     bool
	}{
		{"unbekannt", model.UnknownRowEstimate(), false},
		{"bekannte Null", estimateOf(t, 0), false},
		{"eine Zeile unter der Richtgröße", estimateOf(t, estimatedRowsGuideline-1), false},
		{"genau auf der Richtgröße", estimateOf(t, estimatedRowsGuideline), false},
		{"eine Zeile über der Richtgröße", estimateOf(t, estimatedRowsGuideline+1), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := warnsEstimatedSize(tc.estimate); got != tc.want {
				t.Fatalf("warnsEstimatedSize = %t, will %t", got, tc.want)
			}
		})
	}
}

// TestWarnsCopyDurationBoundary trägt Warnung (2) an der Toleranz: genau auf
// ihr keine Warnung, eine Nanosekunde darüber gesetzt.
func TestWarnsCopyDurationBoundary(t *testing.T) {
	started := model.NewTimePoint(1_000)
	cases := []struct {
		name    string
		elapsed int64
		want    bool
	}{
		{"null", 0, false},
		{"eine Nanosekunde unter der Toleranz", copyDurationToleranceNanos - 1, false},
		{"genau auf der Toleranz", copyDurationToleranceNanos, false},
		{"eine Nanosekunde über der Toleranz", copyDurationToleranceNanos + 1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := warnsCopyDuration(started, model.NewTimePoint(1_000+tc.elapsed)); got != tc.want {
				t.Fatalf("warnsCopyDuration(%d ns) = %t, will %t", tc.elapsed, got, tc.want)
			}
		})
	}
}

// TestWarnedCopyDuration trägt die Regeln der Anwendung auf einen Run: eine
// gesetzte Warnung bleibt gesetzt, ein Run ohne Beginn der Kopie (`queued`)
// trägt nie eine Warnung, und der Empfänger bleibt unverändert.
func TestWarnedCopyDuration(t *testing.T) {
	late := model.NewTimePoint(1_000 + copyDurationToleranceNanos + 1)
	queued, err := model.NewQueuedBackfillRun("run-1", "src-1", "public", "orders", model.NewTimePoint(500))
	if err != nil {
		t.Fatal(err)
	}
	running, err := queued.Start(model.NewTimePoint(1_000))
	if err != nil {
		t.Fatal(err)
	}

	if got := warnedCopyDuration(running, late); !got.WarnDuration {
		t.Fatalf("laufender Run über der Toleranz trägt keine Warnung: %+v", got)
	}
	if running.WarnDuration {
		t.Fatal("der Empfänger wurde verändert")
	}

	warned := running
	warned.WarnDuration = true
	if got := warnedCopyDuration(warned, model.NewTimePoint(1_001)); !got.WarnDuration {
		t.Fatalf("eine gesetzte Warnung wurde zurückgenommen: %+v", got)
	}

	if got := warnedCopyDuration(queued, late); got.WarnDuration {
		t.Fatalf("ein Run ohne Beginn der Kopie trägt eine Warnung: %+v", got)
	}
}

// TestWarnedEstimatedSize trägt die Anwendung auf den Run: die Warnung
// folgt der Schätzung des Runs und lässt die übrigen Felder unverändert.
func TestWarnedEstimatedSize(t *testing.T) {
	run, err := model.NewQueuedBackfillRun("run-1", "src-1", "public", "orders", model.NewTimePoint(500))
	if err != nil {
		t.Fatal(err)
	}
	over := run.WithEstimatedRows(estimateOf(t, estimatedRowsGuideline+1))
	got := warnedEstimatedSize(over)
	if !got.WarnEstimatedSize {
		t.Fatalf("Schätzung über der Richtgröße ohne Warnung: %+v", got)
	}
	got.WarnEstimatedSize = false
	if got != over {
		t.Fatalf("die Anwendung änderte mehr als die Warnung: %+v ≠ %+v", got, over)
	}
	if got := warnedEstimatedSize(run); got.WarnEstimatedSize {
		t.Fatalf("unbekannte Schätzung trägt eine Warnung: %+v", got)
	}
}
