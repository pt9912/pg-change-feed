package outbound

import (
	"context"
	stderrors "errors"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrBackfillStorage trägt die Fehlerklasse `storage` der Backfill-Ports
// (`SPEC-008`, `ADR-0023`): ein Persistenzfehler am Run-Zustand oder am
// Schreiber bleibt über `errors.Is` klassifizierbar, ohne Treibertyp.
var ErrBackfillStorage = stderrors.New("Fehlerklasse storage: Persistenzfehler im Backfill-Speicher")

// BackfillRunPort trägt die Fähigkeit, den Zustand eines angenommenen
// Backfill-Runs fortzuschreiben (`ARC-004`, Fähigkeits-Port je `ADR-0034`;
// `SPEC-029`). Er hat **keine** Operation „anlegen": die Run-Zeile entsteht
// allein über `BackfillAdmissionPort.Admit`. Jede schreibende Operation
// wirkt auf einen Run, der noch nicht endgültig beendet ist, und ändert
// keinen beendeten Run.
//
// Zeitbegrenzung: jede Operation ist adapterseitig zeitbegrenzt. Der Use Case
// hält den Endzustand auf einem vom Abbruch gelösten Kontext fest
// (`context.WithoutCancel`, ohne Frist des Aufrufers); der Adapter beendet
// sich bei einem abgelösten Kontext nicht vorzeitig und begrenzt seine Dauer
// selbst.
type BackfillRunPort interface {
	// Queued liest die `queued`-Runs der Quelle in Antragsreihenfolge
	// (`requested_at`, dann Run-Kennung).
	Queued(ctx context.Context, source model.SourceID) ([]model.BackfillRun, error)

	// MarkRunning setzt den Run auf `running` und hält den Beginn der
	// Kopierdauer fest.
	MarkRunning(ctx context.Context, run model.BackfillRun) error

	// RecordProgress schreibt die Snapshot-Position und den
	// Fortschrittszähler des Runs fort — außerhalb der Daten-Transaktion,
	// damit der Fortschritt für Leser sichtbar ist.
	RecordProgress(ctx context.Context, run model.BackfillRun) error

	// Finish hält den Endzustand eines Runs ohne Daten-Commit fest:
	// `failed` und `interrupted` mit `finished_at` und Fehlertext,
	// `completed` für eine leere Tabelle. Den Endzustand `completed` eines
	// Runs mit Daten schreibt der Schreiber in seinem einen Commit.
	//
	// `Finish` auf einen bereits beendeten Run ist ein wirkungsloser Erfolg:
	// der Endzustand bleibt unverändert, der Aufruf meldet keinen Fehler. Der
	// Use Case setzt das voraus: ist der Ausgang eines Commits unbekannt
	// (Verbindungsabbruch), rollt er zurück und ruft `Finish(failed)`; hat der
	// Commit serverseitig gewirkt, bleibt die Zeile `completed`.
	Finish(ctx context.Context, run model.BackfillRun) error

	// InterruptRunning setzt jeden `running`-Run der Quelle auf
	// `interrupted` und meldet ihre Zahl (Abgleich beim Prozessstart); die
	// Operation ist adapterseitig zeitbegrenzt.
	InterruptRunning(ctx context.Context, source model.SourceID, at model.TimePoint) (int, error)
}
