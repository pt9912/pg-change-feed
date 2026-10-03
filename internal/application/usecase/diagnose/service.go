// Package diagnose trägt den Diagnose-Use-Case (`ADR-0132`): er ist eine
// dünne Fassade über dem `DiagnosticsPort` — keine eigene Abfrage, keine
// zweite Sortierung, kein zweiter Lesepfad.
package diagnose

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// DiagnoseQuery und DiagnoseResult sind die Transport-Typen des
// Diagnose-Use-Cases (`ADR-0042`): ihre Definition liegt am Inbound-Port —
// der Port trägt seinen Vertrag einschließlich der Transport-Typen, Ports
// referenzieren nur die Domain —; die Use-Case-Adressen sind Aliase
// desselben Typs.
type (
	DiagnoseQuery  = inbound.DiagnoseQuery
	DiagnoseResult = inbound.DiagnoseResult
)

// DiagnoseService implementiert `inbound.DiagnoseUseCase` und liest die
// Diagnosesignale über den `DiagnosticsPort`.
type DiagnoseService struct {
	diagnostics outbound.DiagnosticsPort
}

// NewDiagnoseService verdrahtet den Use Case mit dem Diagnostics-Port.
func NewDiagnoseService(diagnostics outbound.DiagnosticsPort) *DiagnoseService {
	return &DiagnoseService{diagnostics: diagnostics}
}

var _ inbound.DiagnoseUseCase = (*DiagnoseService)(nil)

// Diagnose bildet `DiagnoseQuery` auf einen `DiagnosticsPort.Read`-Aufruf ab
// und übersetzt das Ergebnis unverändert in `DiagnoseResult` — keine eigene
// Abfrage, keine zweite Klassifikation. Eine leere Quellen-Kennung ist eine
// ungültige Eingabe und endet über den bestehenden Fehlerpfad, keine stille
// Übernahme über alle Quellen (dieselbe Grenze wie `ReadChangesService`).
func (s *DiagnoseService) Diagnose(ctx context.Context, query DiagnoseQuery) (DiagnoseResult, error) {
	if query.Source == "" {
		return DiagnoseResult{}, domainerrors.ErrEmptyIdentifier
	}
	snapshot, err := s.diagnostics.Read(ctx, query.Source)
	if err != nil {
		return DiagnoseResult{}, err
	}
	return toDiagnoseResult(snapshot), nil
}

// toDiagnoseResult übersetzt den Outbound-Snapshot in die Inbound-Rückgabe —
// zwei Schichten, zwei Typen, 1:1 gemappt, keine inhaltliche Änderung.
func toDiagnoseResult(snapshot outbound.DiagnosticsSnapshot) DiagnoseResult {
	consumerLags := make([]inbound.ConsumerLag, 0, len(snapshot.ConsumerLags))
	for _, lag := range snapshot.ConsumerLags {
		consumerLags = append(consumerLags, inbound.ConsumerLag{ConsumerID: lag.ConsumerID, Lag: lag.Lag})
	}
	var blocker *inbound.RetentionBlocker
	if snapshot.RetentionBlocker != nil {
		blocker = &inbound.RetentionBlocker{
			ConsumerID:           snapshot.RetentionBlocker.ConsumerID,
			Name:                 snapshot.RetentionBlocker.Name,
			AcknowledgedPosition: snapshot.RetentionBlocker.AcknowledgedPosition,
			Backlog:              snapshot.RetentionBlocker.Backlog,
		}
	}
	backfill := make([]inbound.BackfillTableStatus, 0, len(snapshot.Backfill))
	for _, table := range snapshot.Backfill {
		backfill = append(backfill, inbound.BackfillTableStatus{
			Schema:            table.Schema,
			Table:             table.Table,
			Status:            table.Status,
			RowsCopied:        table.RowsCopied,
			EstimatedRows:     table.EstimatedRows,
			WarnEstimatedSize: table.WarnEstimatedSize,
			WarnDuration:      table.WarnDuration,
			ErrorMessage:      table.ErrorMessage,
		})
	}
	return DiagnoseResult{
		HeartbeatAgeSeconds: snapshot.HeartbeatAgeSeconds,
		ErrorClass:          snapshot.ErrorClass,
		ErrorCode:           snapshot.ErrorCode,
		CaptureLag:          snapshot.CaptureLag,
		ConsumerLags:        consumerLags,
		RetentionBlocker:    blocker,
		StorageBytes:        snapshot.StorageBytes,
		Backfill:            backfill,
	}
}
