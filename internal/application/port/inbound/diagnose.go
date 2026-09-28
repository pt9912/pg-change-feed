// Der Diagnose-Use-Case trägt diese Datei an einer Stelle (`ADR-0132`): die
// Driving-Adapter (CLI, HTTP, gRPC) lesen die Diagnosesignale einer Quelle
// über ihn; die Orchestrierung liegt im Application Service.

package inbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// DiagnoseQuery trägt die Eingabe des Diagnose-Lesens: die Quelle
// (**Pflicht**, `ADR-0132`).
type DiagnoseQuery struct {
	Source model.SourceID
}

// ConsumerLag trägt den Verarbeitungsrückstand eines Consumers mit
// mindestens einer bestätigten Position; `Lag == nil` bedeutet, dass die
// gebundene Quelle noch nie eine Transaktion trug.
type ConsumerLag struct {
	ConsumerID string
	Lag        *float64
}

// RetentionBlocker trägt den aktuell die Löschung blockierenden Consumer
// einer Quelle; `Backlog == nil` bedeutet, dass die Quelle noch nie eine
// Transaktion trug.
type RetentionBlocker struct {
	ConsumerID           string
	Name                 string
	AcknowledgedPosition int64
	Backlog              *int64
}

// BackfillTableStatus trägt den zuletzt beantragten Backfill-Run einer
// Tabelle; `EstimatedRows == nil` bedeutet „unbekannt", nie `0`.
type BackfillTableStatus struct {
	Schema            string
	Table             string
	Status            string
	RowsCopied        int64
	EstimatedRows     *int64
	WarnEstimatedSize bool
	WarnDuration      bool
	ErrorMessage      string
}

// DiagnoseResult trägt die sechs Diagnose-Signalgruppen einer Quelle
// (`ADR-0132`), 1:1 aus dem bestehenden CLI-Text-Bericht übernommen.
// `HeartbeatAgeSeconds == nil` bedeutet „kein Lebenszeichen — Instanz hat
// noch nie geschlagen"; in diesem Fall trägt `ErrorClass` ebenfalls `nil`,
// gelesen als „unbekannt", nicht als „Normalbetrieb" — die Unterscheidung
// liegt an `HeartbeatAgeSeconds`, nicht an einem dritten Feld.
// `RetentionBlocker == nil` bedeutet „kein Blocker".
type DiagnoseResult struct {
	HeartbeatAgeSeconds *float64
	ErrorClass          *string
	CaptureLag          float64
	ConsumerLags        []ConsumerLag
	RetentionBlocker    *RetentionBlocker
	StorageBytes        float64
	Backfill            []BackfillTableStatus
}

// DiagnoseUseCase liest die Diagnosesignale einer Quelle über den
// `DiagnosticsPort` (`ADR-0132`): er ist eine Übersetzung, keine eigene
// Abfrage — alle drei Zugriffswege (CLI, HTTP, gRPC) laufen über ihn, kein
// zweiter Domänenpfad.
type DiagnoseUseCase interface {
	Diagnose(ctx context.Context, query DiagnoseQuery) (DiagnoseResult, error)
}
