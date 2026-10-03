package diagnose_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/diagnose"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeDiagnostics trägt den `DiagnosticsPort` als Fälschung (`ADR-0030`):
// Whitebox-Test des Use Case ohne reale Persistenz. Sie zeichnet die
// aufgerufene Quelle auf.
type fakeDiagnostics struct {
	snapshot outbound.DiagnosticsSnapshot
	err      error
	sources  []model.SourceID
}

func (f *fakeDiagnostics) Read(_ context.Context, source model.SourceID) (outbound.DiagnosticsSnapshot, error) {
	f.sources = append(f.sources, source)
	if f.err != nil {
		return outbound.DiagnosticsSnapshot{}, f.err
	}
	return f.snapshot, nil
}

// TestDiagnoseRuftPortMitDerQuelleAufUndUebersetztUnveraendert trägt den
// Kern der Übersetzung (`ADR-0132` Fitness Function): der Use Case ruft
// `DiagnosticsPort.Read` mit der Quelle aus der Query auf und gibt dessen
// Ergebnis unverändert als `DiagnoseResult` zurück — keine eigene Abfrage,
// keine zweite Sortierung.
//
// Rot färbende Mutation: in `toDiagnoseResult` ein Feld nicht übernehmen
// (z. B. `StorageBytes` weglassen) — dieser Test färbt für das betroffene
// Feld rot.
func TestDiagnoseRuftPortMitDerQuelleAufUndUebersetztUnveraendert(t *testing.T) {
	age := 1.5
	errorClass := "schema"
	errorCode := "PCF-E4003"
	lag := 3.0
	backlog := int64(7)
	estimated := int64(100)
	port := &fakeDiagnostics{snapshot: outbound.DiagnosticsSnapshot{
		HeartbeatAgeSeconds: &age,
		ErrorClass:          &errorClass,
		ErrorCode:           &errorCode,
		CaptureLag:          2.5,
		ConsumerLags:        []outbound.ConsumerLagSnapshot{{ConsumerID: "c-1", Lag: &lag}},
		RetentionBlocker: &outbound.RetentionBlockerSnapshot{
			ConsumerID: "c-1", Name: "Consumer 1", AcknowledgedPosition: 10, Backlog: &backlog,
		},
		StorageBytes: 4096,
		Backfill: []outbound.BackfillTableSnapshot{{
			Schema: "public", Table: "orders", Status: "completed", RowsCopied: 5,
			EstimatedRows: &estimated, WarnEstimatedSize: true, WarnDuration: false, ErrorMessage: "",
		}},
	}}

	result, err := diagnose.NewDiagnoseService(port).Diagnose(context.Background(), inbound.DiagnoseQuery{Source: "src-1"})
	if err != nil {
		t.Fatalf("Diagnose: %v", err)
	}
	if len(port.sources) != 1 || port.sources[0] != "src-1" {
		t.Fatalf("Port-Aufrufe = %+v, wollen genau einen Aufruf mit src-1", port.sources)
	}
	if result.HeartbeatAgeSeconds == nil || *result.HeartbeatAgeSeconds != age {
		t.Fatalf("HeartbeatAgeSeconds = %v, wollen %v", result.HeartbeatAgeSeconds, age)
	}
	if result.ErrorClass == nil || *result.ErrorClass != errorClass {
		t.Fatalf("ErrorClass = %v, wollen %v", result.ErrorClass, errorClass)
	}
	if result.ErrorCode == nil || *result.ErrorCode != errorCode {
		t.Fatalf("ErrorCode = %v, wollen %v", result.ErrorCode, errorCode)
	}
	if result.CaptureLag != 2.5 {
		t.Fatalf("CaptureLag = %v, wollen 2.5", result.CaptureLag)
	}
	if len(result.ConsumerLags) != 1 || result.ConsumerLags[0].ConsumerID != "c-1" || *result.ConsumerLags[0].Lag != lag {
		t.Fatalf("ConsumerLags = %+v", result.ConsumerLags)
	}
	if result.RetentionBlocker == nil || result.RetentionBlocker.ConsumerID != "c-1" ||
		result.RetentionBlocker.Name != "Consumer 1" || result.RetentionBlocker.AcknowledgedPosition != 10 ||
		*result.RetentionBlocker.Backlog != backlog {
		t.Fatalf("RetentionBlocker = %+v", result.RetentionBlocker)
	}
	if result.StorageBytes != 4096 {
		t.Fatalf("StorageBytes = %v, wollen 4096", result.StorageBytes)
	}
	if len(result.Backfill) != 1 || result.Backfill[0].Schema != "public" || result.Backfill[0].Table != "orders" ||
		result.Backfill[0].Status != "completed" || result.Backfill[0].RowsCopied != 5 ||
		*result.Backfill[0].EstimatedRows != estimated || !result.Backfill[0].WarnEstimatedSize || result.Backfill[0].WarnDuration {
		t.Fatalf("Backfill = %+v", result.Backfill)
	}
}

// TestDiagnoseUebersetztAbwesenheitenAlsNil trägt die drei Abwesenheits-Fälle
// (kein Lebenszeichen, kein Blocker, unbekannte Schätzung/Rückstand) —
// dieselbe Nil-Lesart wie der bestehende CLI-Text.
func TestDiagnoseUebersetztAbwesenheitenAlsNil(t *testing.T) {
	port := &fakeDiagnostics{snapshot: outbound.DiagnosticsSnapshot{
		ConsumerLags: []outbound.ConsumerLagSnapshot{{ConsumerID: "c-1", Lag: nil}},
		Backfill:     []outbound.BackfillTableSnapshot{{Schema: "public", Table: "orders", EstimatedRows: nil}},
	}}

	result, err := diagnose.NewDiagnoseService(port).Diagnose(context.Background(), inbound.DiagnoseQuery{Source: "src-1"})
	if err != nil {
		t.Fatalf("Diagnose: %v", err)
	}
	if result.HeartbeatAgeSeconds != nil || result.ErrorClass != nil {
		t.Fatalf("HeartbeatAgeSeconds/ErrorClass = %v/%v, wollen beide nil", result.HeartbeatAgeSeconds, result.ErrorClass)
	}
	if result.RetentionBlocker != nil {
		t.Fatalf("RetentionBlocker = %+v, wollen nil (kein Blocker)", result.RetentionBlocker)
	}
	if result.ConsumerLags[0].Lag != nil {
		t.Fatalf("ConsumerLags[0].Lag = %v, wollen nil (unbekannt)", result.ConsumerLags[0].Lag)
	}
	if result.Backfill[0].EstimatedRows != nil {
		t.Fatalf("Backfill[0].EstimatedRows = %v, wollen nil (unbekannt)", result.Backfill[0].EstimatedRows)
	}
}

// TestDiagnoseLeereMengenBleibenGesetzt trägt dieselbe Zusage wie
// `TestReadChangesEmptyResultIsSet`: eine leere Menge bleibt eine leere,
// gesetzte Liste, nie `nil`.
func TestDiagnoseLeereMengenBleibenGesetzt(t *testing.T) {
	port := &fakeDiagnostics{}

	result, err := diagnose.NewDiagnoseService(port).Diagnose(context.Background(), inbound.DiagnoseQuery{Source: "src-1"})
	if err != nil {
		t.Fatalf("Diagnose: %v", err)
	}
	if result.ConsumerLags == nil || len(result.ConsumerLags) != 0 {
		t.Fatalf("ConsumerLags = %v, wollen leere, gesetzte Liste", result.ConsumerLags)
	}
	if result.Backfill == nil || len(result.Backfill) != 0 {
		t.Fatalf("Backfill = %v, wollen leere, gesetzte Liste", result.Backfill)
	}
}

// TestDiagnoseRejectsMissingSource trägt dieselbe Grenze wie
// `TestReadChangesRejectsMissingSource`: eine leere Quellen-Kennung endet
// als ungültige Eingabe, der Port wird nicht berührt.
func TestDiagnoseRejectsMissingSource(t *testing.T) {
	port := &fakeDiagnostics{}

	_, err := diagnose.NewDiagnoseService(port).Diagnose(context.Background(), inbound.DiagnoseQuery{})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Fehler = %v, wollen %v", err, domainerrors.ErrEmptyIdentifier)
	}
	if len(port.sources) != 0 {
		t.Fatalf("Port-Aufrufe = %d, wollen 0", len(port.sources))
	}
}

// TestDiagnoseCarriesPortFailure trägt dieselbe Zusage wie
// `TestReadChangesCarriesPortFailure`: ein Port-Fehler kommt unverändert
// zurück.
func TestDiagnoseCarriesPortFailure(t *testing.T) {
	cause := stderrors.New("Verbindung abgelehnt")
	port := &fakeDiagnostics{err: cause}

	_, err := diagnose.NewDiagnoseService(port).Diagnose(context.Background(), inbound.DiagnoseQuery{Source: "src-1"})
	if !stderrors.Is(err, cause) {
		t.Fatalf("Fehler = %v, wollen die Ursache", err)
	}
}

var _ outbound.DiagnosticsPort = (*fakeDiagnostics)(nil)
