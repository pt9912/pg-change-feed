package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrDiagnosticsStorage trägt einen Code der Klasse `internal` dieses Ports
// (`ADR-0132`): ein Lesefehler an einer der Diagnose-Views bleibt über
// `errors.Is` klassifizierbar, ohne Treibertyp — dieselbe Übersetzungsform
// wie `ErrHeartbeatStorage`/`ErrStorage`.
var ErrDiagnosticsStorage = messagecode.New(messagecode.DiagnosticsReadFailed, "Persistenzfehler beim Lesen der Diagnose-Views")

// ConsumerLagSnapshot trägt den Verarbeitungsrückstand eines Consumers mit
// mindestens einer bestätigten Position; `Lag == nil` bedeutet, dass die
// gebundene Quelle noch nie eine Transaktion trug (`ADR-0132`).
type ConsumerLagSnapshot struct {
	ConsumerID string
	Lag        *float64
}

// RetentionBlockerSnapshot trägt den aktuell die Löschung blockierenden
// Consumer einer Quelle; `Backlog == nil` bedeutet, dass die Quelle noch nie
// eine Transaktion trug (`ADR-0132`).
type RetentionBlockerSnapshot struct {
	ConsumerID           string
	Name                 string
	AcknowledgedPosition int64
	Backlog              *int64
}

// BackfillTableSnapshot trägt den zuletzt beantragten Backfill-Run einer
// Tabelle; `EstimatedRows == nil` bedeutet „unbekannt", nie `0` (`ADR-0132`).
type BackfillTableSnapshot struct {
	Schema            string
	Table             string
	Status            string
	RowsCopied        int64
	EstimatedRows     *int64
	WarnEstimatedSize bool
	WarnDuration      bool
	ErrorMessage      string
}

// DiagnosticsSnapshot bündelt die sechs Diagnose-Signalgruppen einer Quelle
// (`ADR-0132`): dieselben Felder, die der bestehende CLI-Sondermodus bereits
// als Text ausgibt. `HeartbeatAgeSeconds == nil` bedeutet „kein
// Lebenszeichen"; in diesem Fall trägt `ErrorClass` ebenfalls `nil`, gelesen
// als „unbekannt", nicht als „Normalbetrieb" — die Unterscheidung liegt an
// `HeartbeatAgeSeconds`, nicht an einem dritten Feld. `ErrorCode` ist der
// Meldungscode des Fehlerzustands und `nil` im Normalbetrieb sowie bei einem
// Fehlerzustand ohne Code. `RetentionBlocker == nil` bedeutet „kein Blocker".
type DiagnosticsSnapshot struct {
	HeartbeatAgeSeconds *float64
	ErrorClass          *string
	ErrorCode           *string
	CaptureLag          float64
	ConsumerLags        []ConsumerLagSnapshot
	RetentionBlocker    *RetentionBlockerSnapshot
	StorageBytes        float64
	Backfill            []BackfillTableSnapshot
}

// DiagnosticsPort trägt den Lesezugriff auf die Diagnose-Views einer Quelle
// (`ADR-0132`): eine Fähigkeit, ein Aufruf, ein Bericht — dieselbe Bündelung,
// die der bestehende CLI-Sondermodus bereits als Text ausgibt, hier
// strukturiert statt direkt gedruckt. Die erste Produktionsimplementierung
// ist der `postgresstorage`-Adapter (`cdc_reader`-Rolle); Tests setzen
// Fake-Ports ein.
type DiagnosticsPort interface {
	// Read liest die sechs Signalgruppen der Quelle in einem Aufruf. Ein
	// Lesefehler an einer der Views geht als `ErrDiagnosticsStorage`
	// zurück.
	Read(ctx context.Context, source model.SourceID) (DiagnosticsSnapshot, error)
}
